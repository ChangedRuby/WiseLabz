package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

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
