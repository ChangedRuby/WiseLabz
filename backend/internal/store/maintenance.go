package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MaintenanceWindowRecord represents a row in the maintenance_windows table:
// a time-boxed, per-connector suppression of scheduled sync and drift
// alerting/change-record creation (#236).
type MaintenanceWindowRecord struct {
	ID          string `json:"id"`
	ConnectorID string `json:"connectorId"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
}

const maintenanceWindowColumns = `id, connector_id, starts_at, ends_at, created_by, created_at`

// CreateMaintenanceWindow inserts a new maintenance window.
func (s *Store) CreateMaintenanceWindow(ctx context.Context, m *MaintenanceWindowRecord) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.CreatedAt == "" {
		m.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO maintenance_windows (id, connector_id, starts_at, ends_at, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, m.ID, m.ConnectorID, m.StartsAt, m.EndsAt, m.CreatedBy, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("create maintenance window: %w", err)
	}
	return nil
}

// GetActiveMaintenanceWindow returns the connector's active maintenance
// window (ends_at > now), or nil if none. When multiple windows overlap
// (e.g. a new one opened before an earlier one expired), the
// latest-expiring one is returned.
func (s *Store) GetActiveMaintenanceWindow(ctx context.Context, connectorID string) (*MaintenanceWindowRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	m, err := scanMaintenanceWindow(s.db.QueryRowContext(ctx, `
		SELECT `+maintenanceWindowColumns+` FROM maintenance_windows
		WHERE connector_id = ? AND ends_at > ?
		ORDER BY ends_at DESC LIMIT 1
	`, connectorID, now))
	if errors.Is(err, ErrNotFound) {
		return nil, nil //nolint:nilnil
	}
	if err != nil {
		return nil, fmt.Errorf("get active maintenance window: %w", err)
	}
	return m, nil
}

// CloseMaintenanceWindow ends a maintenance window early by setting
// ends_at to now. A no-op on an already-closed (or nonexistent) window
// returns ErrNotFound only when the id doesn't exist at all.
func (s *Store) CloseMaintenanceWindow(ctx context.Context, id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := s.db.ExecContext(ctx, `UPDATE maintenance_windows SET ends_at = ? WHERE id = ?`, now, id)
	if err != nil {
		return fmt.Errorf("close maintenance window: %w", err)
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

// ListActiveMaintenanceWindows returns every currently active maintenance
// window (ends_at > now), across all connectors — used by the ServicesPage
// badge to avoid an N+1 GetActiveMaintenanceWindow call per row. Never
// returns a nil slice.
func (s *Store) ListActiveMaintenanceWindows(ctx context.Context) ([]MaintenanceWindowRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+maintenanceWindowColumns+` FROM maintenance_windows
		WHERE ends_at > ?
		ORDER BY ends_at ASC
	`, now)
	if err != nil {
		return nil, fmt.Errorf("list active maintenance windows: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	windows := make([]MaintenanceWindowRecord, 0)
	for rows.Next() {
		m, err := scanMaintenanceWindow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan maintenance window: %w", err)
		}
		windows = append(windows, *m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate maintenance windows: %w", err)
	}
	return windows, nil
}

func scanMaintenanceWindow(row rowScanner) (*MaintenanceWindowRecord, error) {
	var m MaintenanceWindowRecord
	if err := row.Scan(&m.ID, &m.ConnectorID, &m.StartsAt, &m.EndsAt, &m.CreatedBy, &m.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}
