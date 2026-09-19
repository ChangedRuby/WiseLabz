package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/crypto"
)

// ConnectorRecord represents a row in the connectors table.
type ConnectorRecord struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	Type          string `json:"type"`
	URL           string `json:"url"`
	Owner         string `json:"owner"`
	VerifyTLS     bool   `json:"verifyTls"`
	ConfigData    string `json:"configData"`
	Enabled       bool   `json:"enabled"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
	LastSyncAt    string `json:"lastSyncAt"`
	// ScheduleSeconds is the auto-sync cadence; nil means manual-only (never
	// scheduled). LastSyncDurationMs is nil until the first sync completes.
	ScheduleSeconds    *int `json:"scheduleSeconds"`
	LastSyncDurationMs *int `json:"lastSyncDurationMs"`
	// NextRunAt is "" when there is no scheduled next run, mirroring LastSyncAt's
	// empty-string-means-null convention rather than a pointer.
	NextRunAt     string `json:"nextRunAt"`
	LastSyncError string `json:"lastSyncError"`
	RetryCount    int    `json:"retryCount"`
	// CredentialExpiresAt is "" when the connector's stored credentials have
	// no known expiry. When set (RFC3339) and in the past, the sync engine
	// refuses to Fetch until refreshed (via CredentialRefresher) or updated.
	CredentialExpiresAt string `json:"credentialExpiresAt"`
	// SecretRotatedAt is when a secret-typed config field's value was last
	// actually changed (not just resubmitted or renamed). Set on create and
	// bumped by UpdateConnector's caller when SecretFieldsChanged reports a
	// real change. Read-only over the API.
	SecretRotatedAt string `json:"secretRotatedAt"`
	// UserExpiresAt is an optional operator-set expiry (RFC3339) the
	// credential_rotation quality check treats as an upper bound on the
	// rotation due date, independent of RotationMaxAgeDays.
	UserExpiresAt string `json:"userExpiresAt"`
	// RotationMaxAgeDays overrides the global rotation.max_age_days config
	// for this connector. nil means "use the global default".
	RotationMaxAgeDays *int   `json:"rotationMaxAgeDays"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

// IsCredentialExpired reports whether the connector's credentials have a
// known expiry that has passed as of now.
func (c *ConnectorRecord) IsCredentialExpired(now time.Time) bool {
	if c.CredentialExpiresAt == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, c.CredentialExpiresAt)
	if err != nil {
		return false
	}
	return now.After(t)
}

// connectorColumns is the shared column list for every connector SELECT.
const connectorColumns = `id, name, category, type, url, owner, verify_tls, config_data, enabled, status, status_message,
	last_sync_at, schedule_seconds, next_run_at, last_sync_duration_ms, last_sync_error, retry_count, credential_expires_at,
	secret_rotated_at, user_expires_at, rotation_max_age_days, created_at, updated_at`

// SnapshotRecord represents a row in the service_snapshots table.
type SnapshotRecord struct {
	ID          string `json:"id"`
	ConnectorID string `json:"connectorId"`
	Data        string `json:"data"`
	FetchedAt   string `json:"fetchedAt"`
}

// --- Connector CRUD ---

// CreateConnector inserts a new connector.
func (s *Store) CreateConnector(ctx context.Context, c *ConnectorRecord) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if c.CreatedAt == "" {
		c.CreatedAt = now
	}
	if c.UpdatedAt == "" {
		c.UpdatedAt = now
	}
	if c.Status == "" {
		c.Status = "unknown"
	}
	if c.ConfigData == "" {
		c.ConfigData = "{}"
	}
	if c.SecretRotatedAt == "" {
		c.SecretRotatedAt = c.CreatedAt
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO connectors (id, name, category, type, url, owner, verify_tls, config_data, enabled, status, status_message, last_sync_at,
			schedule_seconds, next_run_at, last_sync_duration_ms, last_sync_error, retry_count, credential_expires_at,
			secret_rotated_at, user_expires_at, rotation_max_age_days, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, c.ID, c.Name, c.Category, c.Type, c.URL, nilToStr(c.Owner), boolToInt(c.VerifyTLS), c.ConfigData,
		boolToInt(c.Enabled), c.Status, c.StatusMessage, nilToStr(c.LastSyncAt),
		// database/sql converts a nil *int argument to SQL NULL automatically.
		c.ScheduleSeconds, nilToStr(c.NextRunAt), c.LastSyncDurationMs, c.LastSyncError, c.RetryCount,
		nilToStr(c.CredentialExpiresAt), c.SecretRotatedAt, nilToStr(c.UserExpiresAt), c.RotationMaxAgeDays,
		c.CreatedAt, c.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrConflict
		}
		return fmt.Errorf("create connector: %w", err)
	}
	return nil
}

// ExistingConnectorIDs returns the subset of ids that already exist as
// connectors, in one query — used by backup import to check N records
// without N round-trips.
func (s *Store) ExistingConnectorIDs(ctx context.Context, ids []string) (map[string]bool, error) {
	return existingIDs(ctx, s.db, "connectors", ids)
}

// GetConnector retrieves a connector by ID.
func (s *Store) GetConnector(ctx context.Context, id string) (*ConnectorRecord, error) {
	c := &ConnectorRecord{}
	var verifyTLS, enabled int
	var owner, lastSyncAt, nextRunAt, lastSyncError, credentialExpiresAt sql.NullString
	var scheduleSeconds, lastSyncDurationMs sql.NullInt64
	var secretRotatedAt, userExpiresAt sql.NullString
	var rotationMaxAgeDays sql.NullInt64
	err := s.db.QueryRowContext(ctx, `SELECT `+connectorColumns+` FROM connectors WHERE id = ?
	`, id).Scan(&c.ID, &c.Name, &c.Category, &c.Type, &c.URL, &owner, &verifyTLS, &c.ConfigData,
		&enabled, &c.Status, &c.StatusMessage, &lastSyncAt,
		&scheduleSeconds, &nextRunAt, &lastSyncDurationMs, &lastSyncError, &c.RetryCount, &credentialExpiresAt,
		&secretRotatedAt, &userExpiresAt, &rotationMaxAgeDays,
		&c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get connector: %w", err)
	}
	c.VerifyTLS = verifyTLS != 0
	c.Owner = nullStrToStr(owner)
	c.Enabled = enabled != 0
	c.LastSyncAt = nullStrToStr(lastSyncAt)
	c.NextRunAt = nullStrToStr(nextRunAt)
	c.LastSyncError = nullStrToStr(lastSyncError)
	c.ScheduleSeconds = nullInt64ToIntPtr(scheduleSeconds)
	c.LastSyncDurationMs = nullInt64ToIntPtr(lastSyncDurationMs)
	c.CredentialExpiresAt = nullStrToStr(credentialExpiresAt)
	c.SecretRotatedAt = nullStrToStr(secretRotatedAt)
	c.UserExpiresAt = nullStrToStr(userExpiresAt)
	c.RotationMaxAgeDays = nullInt64ToIntPtr(rotationMaxAgeDays)
	return c, nil
}

// ListConnectorsByID returns the requested connectors keyed by ID, for bulk
// actions (bulk-sync/bulk-reauth/bulk-restart) that need every target's full
// record in one round trip.
func (s *Store) ListConnectorsByID(ctx context.Context, ids []string) (map[string]ConnectorRecord, error) {
	connectors := make(map[string]ConnectorRecord, len(ids))
	if len(ids) == 0 {
		return connectors, nil
	}

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.QueryContext(ctx, `SELECT `+connectorColumns+` FROM connectors WHERE id IN (`+placeholders(len(ids))+`)`, args...)
	if err != nil {
		return nil, fmt.Errorf("list connectors by id: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var c ConnectorRecord
		var verifyTLS, enabled int
		var owner, lastSyncAt, nextRunAt, lastSyncError, credentialExpiresAt sql.NullString
		var scheduleSeconds, lastSyncDurationMs sql.NullInt64
		var secretRotatedAt, userExpiresAt sql.NullString
		var rotationMaxAgeDays sql.NullInt64
		if err := rows.Scan(&c.ID, &c.Name, &c.Category, &c.Type, &c.URL, &owner, &verifyTLS, &c.ConfigData,
			&enabled, &c.Status, &c.StatusMessage, &lastSyncAt,
			&scheduleSeconds, &nextRunAt, &lastSyncDurationMs, &lastSyncError, &c.RetryCount, &credentialExpiresAt,
			&secretRotatedAt, &userExpiresAt, &rotationMaxAgeDays,
			&c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		c.VerifyTLS = verifyTLS != 0
		c.Owner = nullStrToStr(owner)
		c.Enabled = enabled != 0
		c.LastSyncAt = nullStrToStr(lastSyncAt)
		c.NextRunAt = nullStrToStr(nextRunAt)
		c.LastSyncError = nullStrToStr(lastSyncError)
		c.ScheduleSeconds = nullInt64ToIntPtr(scheduleSeconds)
		c.LastSyncDurationMs = nullInt64ToIntPtr(lastSyncDurationMs)
		c.CredentialExpiresAt = nullStrToStr(credentialExpiresAt)
		c.SecretRotatedAt = nullStrToStr(secretRotatedAt)
		c.UserExpiresAt = nullStrToStr(userExpiresAt)
		c.RotationMaxAgeDays = nullInt64ToIntPtr(rotationMaxAgeDays)
		connectors[c.ID] = c
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connectors: %w", err)
	}
	return connectors, nil
}

// UpdateConnector updates fields on an existing connector.
func (s *Store) UpdateConnector(ctx context.Context, id string, updates map[string]any) error {
	now := time.Now().UTC().Format(time.RFC3339)
	args := []any{now}
	var parts []string

	for k, v := range updates {
		switch k {
		case "name":
			parts = append(parts, "name = ?")
			args = append(args, v)
		case "url":
			parts = append(parts, "url = ?")
			args = append(args, v)
		case "owner":
			parts = append(parts, "owner = ?")
			args = append(args, v)
		case "verify_tls":
			parts = append(parts, "verify_tls = ?")
			args = append(args, boolToInt(v.(bool)))
		case "config_data":
			parts = append(parts, "config_data = ?")
			args = append(args, v)
		case "enabled":
			parts = append(parts, "enabled = ?")
			args = append(args, boolToInt(v.(bool)))
		case "status":
			parts = append(parts, "status = ?")
			args = append(args, v)
		case "status_message":
			parts = append(parts, "status_message = ?")
			args = append(args, v)
		case "last_sync_at":
			parts = append(parts, "last_sync_at = ?")
			args = append(args, v)
		case "category":
			parts = append(parts, "category = ?")
			args = append(args, v)
		case "type":
			parts = append(parts, "type = ?")
			args = append(args, v)
		case "schedule_seconds":
			parts = append(parts, "schedule_seconds = ?")
			args = append(args, v)
		case "next_run_at":
			parts = append(parts, "next_run_at = ?")
			args = append(args, v)
		case "last_sync_duration_ms":
			parts = append(parts, "last_sync_duration_ms = ?")
			args = append(args, v)
		case "last_sync_error":
			parts = append(parts, "last_sync_error = ?")
			args = append(args, v)
		case "retry_count":
			parts = append(parts, "retry_count = ?")
			args = append(args, v)
		case "credential_expires_at":
			parts = append(parts, "credential_expires_at = ?")
			args = append(args, v)
		case "secret_rotated_at":
			parts = append(parts, "secret_rotated_at = ?")
			args = append(args, v)
		case "user_expires_at":
			parts = append(parts, "user_expires_at = ?")
			args = append(args, v)
		case "rotation_max_age_days":
			parts = append(parts, "rotation_max_age_days = ?")
			args = append(args, v)
		}
	}

	query := "UPDATE connectors SET updated_at = ?"
	if len(parts) > 0 {
		query += ", " + strings.Join(parts, ", ")
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update connector: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteConnector deletes a connector by ID.
func (s *Store) DeleteConnector(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM connectors WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete connector: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ListConnectors returns a paginated list of connectors.
func (s *Store) ListConnectors(ctx context.Context, offset, limit int) ([]ConnectorRecord, int, error) {
	return paginatedQuery(ctx, s.db, "connectors", connectorColumns, "", nil, "created_at DESC", limit, offset, scanConnector)
}

// ListConnectorsByCategory returns connectors filtered by category.
func (s *Store) ListConnectorsByCategory(ctx context.Context, category string, offset, limit int) ([]ConnectorRecord, int, error) {
	return paginatedQuery(ctx, s.db, "connectors", connectorColumns, "WHERE category = ?", []any{category}, "created_at DESC", limit, offset, scanConnector)
}

// ConnectorWithRole is a connector annotated with the requesting user's role.
type ConnectorWithRole struct {
	ConnectorRecord
	Role string
}

// connectorListColumns is connectorColumns qualified for the grant join, with
// config_data blanked: list views never need the (encrypted) credentials.
var connectorListColumns = func() string {
	cols := strings.Split(connectorColumns, ",")
	for i, c := range cols {
		c = strings.TrimSpace(c)
		if c == "config_data" {
			cols[i] = "'' AS config_data"
		} else {
			cols[i] = "c." + c
		}
	}
	return strings.Join(cols, ", ") + ", r.role"
}()

// ListConnectorsForUser returns one page of the connectors the user holds a
// grant on (optionally within a category), newest first, with the user's role
// and the total number of visible connectors. ConfigData is left empty.
func (s *Store) ListConnectorsForUser(ctx context.Context, userID, category string, offset, limit int) ([]ConnectorWithRole, int, error) {
	where, args := "WHERE r.user_id = ?", []any{userID}
	if category != "" {
		where += " AND c.category = ?"
		args = append(args, category)
	}
	return paginatedQuery(ctx, s.db, "connectors c JOIN user_connector_roles r ON r.connector_id = c.id",
		connectorListColumns, where, args, "c.created_at DESC", limit, offset,
		func(row rowScanner) (ConnectorWithRole, error) {
			var out ConnectorWithRole
			c, err := scanConnector(&roleScanner{row: row, role: &out.Role})
			out.ConnectorRecord = c
			return out, err
		})
}

// roleScanner appends a trailing role column to a scanConnector scan.
type roleScanner struct {
	row  rowScanner
	role *string
}

func (r *roleScanner) Scan(dest ...any) error {
	return r.row.Scan(append(dest, r.role)...)
}

// ListAllConnectors returns all connectors (no pagination).
func (s *Store) ListAllConnectors(ctx context.Context) ([]ConnectorRecord, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+connectorColumns+`
		FROM connectors ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list all connectors: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	return scanConnectorRows(rows)
}

// ConnectorName is the id/name pair of a connector, for views that only render
// labels and must not load config_data.
type ConnectorName struct {
	ID   string
	Name string
}

// ListConnectorNames returns every connector's id and name (newest first, same
// order as ListAllConnectors) without loading its config.
func (s *Store) ListConnectorNames(ctx context.Context) ([]ConnectorName, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name FROM connectors ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list connector names: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	out := []ConnectorName{}
	for rows.Next() {
		var c ConnectorName
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector names: %w", err)
	}
	return out, nil
}

// ListDueConnectors returns enabled connectors with a schedule whose next run
// is due (next_run_at <= now, or unset — e.g. right after schedule was first
// configured). Ordered soonest-first. Never returns a nil slice.
func (s *Store) ListDueConnectors(ctx context.Context, now string, limit int) ([]ConnectorRecord, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+connectorColumns+`
		FROM connectors
		WHERE enabled = 1 AND schedule_seconds IS NOT NULL AND (next_run_at IS NULL OR next_run_at <= ?)
		AND NOT EXISTS (SELECT 1 FROM maintenance_windows mw WHERE mw.connector_id = connectors.id AND mw.ends_at > ?)
		ORDER BY next_run_at ASC LIMIT ?
	`, now, now, limit)
	if err != nil {
		return nil, fmt.Errorf("list due connectors: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	return scanConnectorRows(rows)
}

// ClaimDueConnector atomically defers a due connector before fetching it.
// The lease expires after the sync deadline so a process crash cannot leave
// the connector permanently unscheduled. Completion sets the normal cadence.
func (s *Store) ClaimDueConnector(ctx context.Context, id, now, leaseUntil string) (bool, error) {
	result, err := s.db.ExecContext(ctx, `UPDATE connectors SET next_run_at = ?
 WHERE id = ? AND enabled = 1 AND schedule_seconds IS NOT NULL
 AND (next_run_at IS NULL OR next_run_at <= ?)
 AND NOT EXISTS (SELECT 1 FROM maintenance_windows mw WHERE mw.connector_id = connectors.id AND mw.ends_at > ?)
 `, leaseUntil, id, now, now)
	if err != nil {
		return false, fmt.Errorf("claim due connector: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("claim due connector rows: %w", err)
	}
	return n == 1, nil
}

// --- Snapshot operations ---

// CreateSnapshot inserts a new service snapshot.
func (s *Store) CreateSnapshot(ctx context.Context, sn *SnapshotRecord) error {
	if sn.ID == "" {
		sn.ID = uuid.New().String()
	}
	if sn.FetchedAt == "" {
		sn.FetchedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO service_snapshots (id, connector_id, data, fetched_at)
		VALUES (?, ?, ?, ?)
	`, sn.ID, sn.ConnectorID, sn.Data, sn.FetchedAt)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}
	return nil
}

// GetLatestSnapshot returns the most recent snapshot for a connector.
func (s *Store) GetLatestSnapshot(ctx context.Context, connectorID string) (*SnapshotRecord, error) {
	sn := &SnapshotRecord{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, connector_id, data, fetched_at
		FROM service_snapshots WHERE connector_id = ?
		ORDER BY fetched_at DESC LIMIT 1
	`, connectorID).Scan(&sn.ID, &sn.ConnectorID, &sn.Data, &sn.FetchedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get latest snapshot: %w", err)
	}
	return sn, nil
}

// GetSnapshotsByConnector returns snapshots for a connector (newest first).
func (s *Store) GetSnapshotsByConnector(ctx context.Context, connectorID string, limit int) ([]SnapshotRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, connector_id, data, fetched_at
		FROM service_snapshots WHERE connector_id = ?
		ORDER BY fetched_at DESC LIMIT ?
	`, connectorID, limit)
	if err != nil {
		return nil, fmt.Errorf("get snapshots: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var snapshots []SnapshotRecord
	for rows.Next() {
		var sn SnapshotRecord
		if err := rows.Scan(&sn.ID, &sn.ConnectorID, &sn.Data, &sn.FetchedAt); err != nil {
			return nil, fmt.Errorf("scan snapshot: %w", err)
		}
		snapshots = append(snapshots, sn)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate snapshots: %w", err)
	}
	if snapshots == nil {
		snapshots = []SnapshotRecord{}
	}
	return snapshots, nil
}

// CountSnapshotsByConnector returns the number of snapshots for a connector.
func (s *Store) CountSnapshotsByConnector(ctx context.Context, connectorID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM service_snapshots WHERE connector_id = ?`, connectorID).Scan(&count)
	return count, err
}

// CountDocsByConnector returns the number of docs linked to a connector.
func (s *Store) CountDocsByConnector(ctx context.Context, connectorID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM docs WHERE service_id = ?`, connectorID).Scan(&count)
	return count, err
}

// CountChangesByConnector returns the number of changes for a connector.
func (s *Store) CountChangesByConnector(ctx context.Context, connectorID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM changes WHERE service_id = ?`, connectorID).Scan(&count)
	return count, err
}

// CountAlertsByConnector returns the number of alerts for a connector.
func (s *Store) CountAlertsByConnector(ctx context.Context, connectorID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM alerts WHERE service_id = ?`, connectorID).Scan(&count)
	return count, err
}

// --- helpers ---

func scanConnector(row rowScanner) (ConnectorRecord, error) {
	var c ConnectorRecord
	var verifyTLS, enabled int
	var owner, lastSyncAt, nextRunAt, lastSyncError, credentialExpiresAt sql.NullString
	var scheduleSeconds, lastSyncDurationMs sql.NullInt64
	var secretRotatedAt, userExpiresAt sql.NullString
	var rotationMaxAgeDays sql.NullInt64
	if err := row.Scan(&c.ID, &c.Name, &c.Category, &c.Type, &c.URL, &owner, &verifyTLS, &c.ConfigData,
		&enabled, &c.Status, &c.StatusMessage, &lastSyncAt,
		&scheduleSeconds, &nextRunAt, &lastSyncDurationMs, &lastSyncError, &c.RetryCount, &credentialExpiresAt,
		&secretRotatedAt, &userExpiresAt, &rotationMaxAgeDays,
		&c.CreatedAt, &c.UpdatedAt); err != nil {
		return ConnectorRecord{}, err
	}
	c.VerifyTLS = verifyTLS != 0
	c.Owner = nullStrToStr(owner)
	c.Enabled = enabled != 0
	c.LastSyncAt = nullStrToStr(lastSyncAt)
	c.NextRunAt = nullStrToStr(nextRunAt)
	c.LastSyncError = nullStrToStr(lastSyncError)
	c.ScheduleSeconds = nullInt64ToIntPtr(scheduleSeconds)
	c.LastSyncDurationMs = nullInt64ToIntPtr(lastSyncDurationMs)
	c.CredentialExpiresAt = nullStrToStr(credentialExpiresAt)
	c.SecretRotatedAt = nullStrToStr(secretRotatedAt)
	c.UserExpiresAt = nullStrToStr(userExpiresAt)
	c.RotationMaxAgeDays = nullInt64ToIntPtr(rotationMaxAgeDays)
	return c, nil
}

// scanConnectorRows scans every row of a connectors query (non-paginated
// callers that don't go through paginatedQuery). Never returns a nil slice.
func scanConnectorRows(rows *sql.Rows) ([]ConnectorRecord, error) {
	connectors := make([]ConnectorRecord, 0)
	for rows.Next() {
		c, err := scanConnector(rows)
		if err != nil {
			return nil, fmt.Errorf("scan connector: %w", err)
		}
		connectors = append(connectors, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connectors: %w", err)
	}
	return connectors, nil
}

// nullInt64ToIntPtr converts a nullable DB integer to *int (nil when NULL).
func nullInt64ToIntPtr(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return &v
}

func nullStrToStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nilToStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// IsSecretFieldType reports whether a SchemaField.Type holds a value that
// must be encrypted at rest in config_data. "password" is the only kind
// today; "secret" (a future multi-line paste field for PEM certs/SSH keys)
// must be treated the same way, so check this helper rather than comparing
// against "password" directly.
func IsSecretFieldType(t string) bool {
	return t == "password" || t == "secret"
}

// ParseConnectorConfig parses the config_data JSON string into a map,
// decrypting the value of any field the connector type's schema marks as
// secret-bearing (see IsSecretFieldType). A stored value that fails to
// decrypt is assumed to be legacy plaintext (saved before encryption was
// added) and is returned as-is; MarshalConnectorConfig will encrypt it on
// the connector's next save.
func ParseConnectorConfig(connType, data, encKeyB64 string) (map[string]any, error) {
	var cfg map[string]any
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, fmt.Errorf("parse connector config: %w", err)
	}
	if cfg == nil {
		cfg = map[string]any{}
	}

	schema, err := connector.GetTypeSchema(connType)
	if err != nil {
		// Unknown/unregistered connector type: nothing known to decrypt.
		return cfg, nil //nolint:nilerr
	}

	var key []byte
	for _, f := range schema.Fields {
		if !IsSecretFieldType(f.Type) {
			continue
		}
		raw, ok := cfg[f.Key].(string)
		if !ok || raw == "" {
			continue
		}
		if key == nil {
			key, err = crypto.DecodeKey(encKeyB64)
			if err != nil {
				return nil, fmt.Errorf("decode encryption key: %w", err)
			}
		}
		if plaintext, err := crypto.Decrypt(raw, key); err == nil {
			cfg[f.Key] = plaintext
		}
		// else: not valid ciphertext, treat as legacy plaintext and leave as-is.
	}
	return cfg, nil
}

// SecretFieldsChanged reports whether any secret-typed config field's
// plaintext value differs between the connector's currently stored config
// (oldConfigData, as read from connectors.config_data) and newConfig (the
// plaintext config a request is about to save). A rename-only edit or a
// resubmitted-but-unchanged secret must report false; a newly set, changed,
// or cleared secret field reports true.
func SecretFieldsChanged(connType, oldConfigData string, newConfig map[string]any, encKeyB64 string) (bool, error) {
	oldConfig, err := ParseConnectorConfig(connType, oldConfigData, encKeyB64)
	if err != nil {
		return false, fmt.Errorf("parse existing connector config: %w", err)
	}
	schema, err := connector.GetTypeSchema(connType)
	if err != nil {
		// Unknown/unregistered connector type: no secret fields are known,
		// so nothing can have changed.
		return false, nil //nolint:nilerr
	}
	for _, f := range schema.Fields {
		if !IsSecretFieldType(f.Type) {
			continue
		}
		oldVal, _ := oldConfig[f.Key].(string)
		newVal, _ := newConfig[f.Key].(string)
		if oldVal != newVal {
			return true, nil
		}
	}
	return false, nil
}

// MarshalConnectorConfig marshals a config map to a JSON string, encrypting
// the value of any field the connector type's schema marks as secret-bearing
// (see IsSecretFieldType) before marshaling.
func MarshalConnectorConfig(connType string, cfg map[string]any, encKeyB64 string) (string, error) {
	schema, err := connector.GetTypeSchema(connType)
	if err == nil {
		var key []byte
		toEncrypt := make(map[string]any, len(cfg))
		for k, v := range cfg {
			toEncrypt[k] = v
		}
		for _, f := range schema.Fields {
			if !IsSecretFieldType(f.Type) {
				continue
			}
			raw, ok := toEncrypt[f.Key].(string)
			if !ok || raw == "" {
				continue
			}
			if key == nil {
				key, err = crypto.DecodeKey(encKeyB64)
				if err != nil {
					return "", fmt.Errorf("decode encryption key: %w", err)
				}
			}
			encrypted, err := crypto.Encrypt(raw, key)
			if err != nil {
				return "", fmt.Errorf("encrypt connector config field %q: %w", f.Key, err)
			}
			toEncrypt[f.Key] = encrypted
		}
		cfg = toEncrypt
	}
	// Unknown/unregistered connector type: nothing known to encrypt, marshal as-is.

	b, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal connector config: %w", err)
	}
	return string(b), nil
}
