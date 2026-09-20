package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SnapshotRecord represents a row in the service_snapshots table.
type SnapshotRecord struct {
	ID          string `json:"id"`
	ConnectorID string `json:"connectorId"`
	Data        string `json:"data"`
	FetchedAt   string `json:"fetchedAt"`
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
