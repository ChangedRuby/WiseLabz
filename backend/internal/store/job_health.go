package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// JobHealthRecord represents a row in the job_health table: the persisted
// ok/failing status of one scheduled job (see internal/scheduler), keyed by
// job name. Unlike the in-memory flag docexport used to keep, this survives
// a restart, so a job that was already failing doesn't send a duplicate
// "failing" notification just because the process restarted.
type JobHealthRecord struct {
	Name          string `json:"name"`
	CronExpr      string `json:"cronExpr"`
	LastStatus    string `json:"lastStatus"` // "ok" or "failing"
	LastError     string `json:"lastError,omitempty"`
	LastRunAt     string `json:"lastRunAt,omitempty"`
	LastSuccessAt string `json:"lastSuccessAt,omitempty"`
	LastFailureAt string `json:"lastFailureAt,omitempty"`
	UpdatedAt     string `json:"updatedAt"`
}

// scanJobHealth scans one job_health row, translating NULL timestamp
// columns (a job that has never succeeded/failed yet) to "".
func scanJobHealth(scan func(dest ...any) error) (JobHealthRecord, error) {
	var rec JobHealthRecord
	var lastRunAt, lastSuccessAt, lastFailureAt sql.NullString
	if err := scan(&rec.Name, &rec.CronExpr, &rec.LastStatus, &rec.LastError,
		&lastRunAt, &lastSuccessAt, &lastFailureAt, &rec.UpdatedAt); err != nil {
		return rec, err
	}
	rec.LastRunAt = lastRunAt.String
	rec.LastSuccessAt = lastSuccessAt.String
	rec.LastFailureAt = lastFailureAt.String
	return rec, nil
}

// GetJobHealth retrieves the persisted health record for the job named name.
// Returns sql.ErrNoRows if no run has ever been recorded for it.
func (s *Store) GetJobHealth(ctx context.Context, name string) (JobHealthRecord, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT name, cron_expr, last_status, last_error, last_run_at, last_success_at, last_failure_at, updated_at
		FROM job_health WHERE name = ?
	`, name)
	rec, err := scanJobHealth(row.Scan)
	if err != nil {
		return rec, fmt.Errorf("get job health %q: %w", name, err)
	}
	return rec, nil
}

// UpsertJobHealth inserts or updates the health row for rec.Name using
// INSERT ... ON CONFLICT semantics (dialect-neutral, matching
// UpsertRetentionSettings/UpsertBackupSchedule). UpdatedAt defaults to now
// (UTC) when unset. Empty timestamp fields are stored as NULL.
func (s *Store) UpsertJobHealth(ctx context.Context, rec JobHealthRecord) error {
	if rec.UpdatedAt == "" {
		rec.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	strOrNull := func(v string) any {
		if v == "" {
			return nil
		}
		return v
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO job_health (name, cron_expr, last_status, last_error, last_run_at, last_success_at, last_failure_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			cron_expr = excluded.cron_expr,
			last_status = excluded.last_status,
			last_error = excluded.last_error,
			last_run_at = excluded.last_run_at,
			last_success_at = excluded.last_success_at,
			last_failure_at = excluded.last_failure_at,
			updated_at = excluded.updated_at
	`, rec.Name, rec.CronExpr, rec.LastStatus, rec.LastError,
		strOrNull(rec.LastRunAt), strOrNull(rec.LastSuccessAt), strOrNull(rec.LastFailureAt), rec.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert job health %q: %w", rec.Name, err)
	}
	return nil
}

// ListJobHealth returns every persisted job health row, ordered by name for
// stable output.
func (s *Store) ListJobHealth(ctx context.Context) ([]JobHealthRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT name, cron_expr, last_status, last_error, last_run_at, last_success_at, last_failure_at, updated_at
		FROM job_health ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list job health: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	var recs []JobHealthRecord
	for rows.Next() {
		rec, err := scanJobHealth(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("scan job health: %w", err)
		}
		recs = append(recs, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate job health: %w", err)
	}
	return recs, nil
}

// DeleteJobHealth removes the health row for name, if any. Used when a job
// is permanently unregistered (e.g. a deleted scheduled report definition
// in a later PR) so stale rows don't linger in ListJobHealth/the jobs API.
func (s *Store) DeleteJobHealth(ctx context.Context, name string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM job_health WHERE name = ?`, name); err != nil {
		return fmt.Errorf("delete job health %q: %w", name, err)
	}
	return nil
}
