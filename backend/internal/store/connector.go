package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
