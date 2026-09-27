package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SyncRunStatus is the outcome of one sync attempt.
type SyncRunStatus string

// Sync run lifecycle states.
const (
	SyncRunStatusSuccess SyncRunStatus = "success"
	SyncRunStatusError   SyncRunStatus = "error"
	SyncRunStatusSkipped SyncRunStatus = "skipped"
)

// SyncRunRecord represents a row in the sync_runs table: one historical sync
// attempt (manual or scheduled) for a connector.
type SyncRunRecord struct {
	ID          string        `json:"id"`
	ConnectorID string        `json:"connectorId"`
	StartedAt   string        `json:"startedAt"`
	FinishedAt  string        `json:"finishedAt,omitempty"`
	DurationMs  *int          `json:"durationMs,omitempty"`
	Status      SyncRunStatus `json:"status"`
	Error       string        `json:"error,omitempty"`
	// Attempt is the retry attempt number for this run; 1 for the first try.
	Attempt      int     `json:"attempt"`
	ChangesCount int     `json:"changesCount"`
	AlertsCount  int     `json:"alertsCount"`
	SnapshotID   *string `json:"snapshotId"`
}

// CreateSyncRun inserts a new sync run row.
func (s *Store) CreateSyncRun(ctx context.Context, r *SyncRunRecord) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.StartedAt == "" {
		r.StartedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if r.Attempt == 0 {
		r.Attempt = 1
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sync_runs (id, connector_id, started_at, finished_at, duration_ms, status, error, attempt, changes_count, alerts_count, snapshot_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.ID, r.ConnectorID, r.StartedAt, nilToStr(r.FinishedAt), r.DurationMs, r.Status, r.Error, r.Attempt, r.ChangesCount, r.AlertsCount, r.SnapshotID)
	if err != nil {
		return fmt.Errorf("create sync run: %w", err)
	}
	return nil
}

// syncRunColumns is the shared column list for every sync_runs SELECT.
const syncRunColumns = `id, connector_id, started_at, finished_at, duration_ms, status, error, attempt, changes_count, alerts_count, snapshot_id`

// ListSyncRunsByConnector returns sync runs for a connector, newest first.
// Never returns a nil slice.
func (s *Store) ListSyncRunsByConnector(ctx context.Context, connectorID string, limit int) ([]SyncRunRecord, error) {
	query := `SELECT ` + syncRunColumns + ` FROM sync_runs WHERE connector_id = ?
		ORDER BY started_at DESC, finished_at DESC, id DESC LIMIT ?`

	return scanAll(ctx, s.db, "sync runs", query, []any{connectorID, limit}, scanSyncRun)
}

// ListSyncRunsByConnectorKeyset returns one keyset (cursor) page of a
// connector's sync runs, newest first. Rows strictly before cur in
// (started_at, id) order are returned; a zero cur starts at the newest run.
// total counts every run for the connector, ignoring cur.
//
// Keyset mode tie-breaks on id alone, where the offset listing also considers
// finished_at; that extra key cannot participate in a cursor (it is nullable
// and not unique), and ordering by (started_at, id) is what makes the page
// boundary exact.
func (s *Store) ListSyncRunsByConnectorKeyset(ctx context.Context, connectorID string, cur Keyset, limit int) ([]SyncRunRecord, int, error) {
	return keysetQuery(ctx, s.db, "sync_runs", syncRunColumns, "WHERE connector_id = ?", []any{connectorID},
		"started_at", cur, limit, scanSyncRun)
}

// FailedSyncRun is one failed sync run enriched with its connector's name
// (sync_runs only stores connector_id), for cross-connector diagnostics
// views.
type FailedSyncRun struct {
	ConnectorID   string `json:"connectorId"`
	ConnectorName string `json:"connectorName"`
	StartedAt     string `json:"startedAt"`
	Error         string `json:"error"`
}

// ListRecentFailedSyncRuns returns the most recent failed sync runs across
// all connectors, newest first. Never returns a nil slice.
func (s *Store) ListRecentFailedSyncRuns(ctx context.Context, limit int) ([]FailedSyncRun, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT sr.connector_id, c.name, sr.started_at, sr.error
		FROM sync_runs sr
		JOIN connectors c ON c.id = sr.connector_id
		WHERE sr.status = 'error'
		ORDER BY sr.started_at DESC LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list recent failed sync runs: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	runs := []FailedSyncRun{}
	for rows.Next() {
		var r FailedSyncRun
		if err := rows.Scan(&r.ConnectorID, &r.ConnectorName, &r.StartedAt, &r.Error); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		runs = append(runs, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate failed sync runs: %w", err)
	}
	return runs, nil
}

func scanSyncRun(row rowScanner) (SyncRunRecord, error) {
	var r SyncRunRecord
	var finishedAt sql.NullString
	var durationMs sql.NullInt64
	var snapshotID sql.NullString
	err := row.Scan(&r.ID, &r.ConnectorID, &r.StartedAt, &finishedAt, &durationMs, &r.Status, &r.Error, &r.Attempt, &r.ChangesCount, &r.AlertsCount, &snapshotID)
	if errors.Is(err, sql.ErrNoRows) {
		return SyncRunRecord{}, ErrNotFound
	}
	if err != nil {
		return SyncRunRecord{}, err
	}
	r.FinishedAt = nullStrToStr(finishedAt)
	r.DurationMs = nullInt64ToIntPtr(durationMs)
	if snapshotID.Valid {
		r.SnapshotID = &snapshotID.String
	}
	return r, nil
}
