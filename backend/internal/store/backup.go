package store

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// BackupSchedule represents the scheduled backup configuration (single-row
// table with id='default').
type BackupSchedule struct {
	CronExpr    string
	MaxBackups  int
	MaxAgeHours int
	Enabled     bool
	UpdatedAt   string
}

// GetBackupSchedule retrieves the backup schedule from the database.
// Returns ErrNotFound-style behavior (wrapped error) if no row exists.
func (s *Store) GetBackupSchedule(ctx context.Context) (BackupSchedule, error) {
	var sched BackupSchedule
	err := s.db.QueryRowContext(ctx, `
		SELECT cron_expr, max_backups, max_age_hours, enabled, updated_at
		FROM backup_schedule WHERE id = 'default'
	`).Scan(&sched.CronExpr, &sched.MaxBackups, &sched.MaxAgeHours, &sched.Enabled, &sched.UpdatedAt)
	if err != nil {
		return sched, fmt.Errorf("get backup schedule: %w", err)
	}
	return sched, nil
}

// UpsertBackupSchedule inserts or updates the backup schedule using INSERT OR
// REPLACE / ON CONFLICT semantics (dialect-neutral). The id is always 'default'.
func (s *Store) UpsertBackupSchedule(ctx context.Context, sched BackupSchedule) error {
	if sched.UpdatedAt == "" {
		sched.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO backup_schedule (id, cron_expr, max_backups, max_age_hours, enabled, updated_at)
		VALUES ('default', ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			cron_expr = excluded.cron_expr,
			max_backups = excluded.max_backups,
			max_age_hours = excluded.max_age_hours,
			enabled = excluded.enabled,
			updated_at = excluded.updated_at
	`, sched.CronExpr, sched.MaxBackups, sched.MaxAgeHours, boolToInt(sched.Enabled), sched.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert backup schedule: %w", err)
	}
	return nil
}

// BackupRun represents one backup file that was created.
type BackupRun struct {
	ID          string
	TriggeredBy string // "schedule" or "manual"
	FilePath    string
	SizeBytes   int64
	CreatedAt   string
}

// CreateBackupRun inserts a backup run into the database.
func (s *Store) CreateBackupRun(ctx context.Context, run BackupRun) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO backup_runs (id, triggered_by, file_path, size_bytes, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, run.ID, run.TriggeredBy, run.FilePath, run.SizeBytes, run.CreatedAt)
	if err != nil {
		return fmt.Errorf("create backup run: %w", err)
	}
	return nil
}

func scanBackupRun(row rowScanner) (BackupRun, error) {
	var run BackupRun
	err := row.Scan(&run.ID, &run.TriggeredBy, &run.FilePath, &run.SizeBytes, &run.CreatedAt)
	return run, err
}

// ListBackupRuns returns a paginated list of backup runs, ordered by created_at DESC.
// Returns (runs, total_count, error).
func (s *Store) ListBackupRuns(ctx context.Context, limit, offset int) ([]BackupRun, int, error) {
	return paginatedQuery(ctx, s.db, "backup_runs",
		"id, triggered_by, file_path, size_bytes, created_at",
		"", nil, "created_at DESC", limit, offset, scanBackupRun)
}

// PruneBackupRuns deletes backup runs that are either:
// 1. Beyond maxBackups (kept by recency, i.e., NOT in the top-N most recent), OR
// 2. Older than maxAgeCutoff (as an RFC3339 string, e.g. from time.Now().Add(-duration).Format(time.RFC3339))
//
// Returns the deleted rows so the caller can also unlink their files.
// A maxBackups <= 0 is treated as "no limit by count"; a maxAgeCutoff of empty
// string is treated as "no limit by age".
func (s *Store) PruneBackupRuns(ctx context.Context, maxBackups int, maxAgeCutoff string) ([]BackupRun, error) {
	// Build a query to find victim IDs:
	// - If maxBackups > 0: rows NOT in the top-N most recent (by created_at DESC)
	// - If maxAgeCutoff != "": rows older than that timestamp
	// Both conditions are OR'd together, and duplicates are removed.

	var victims []BackupRun

	// If both limits are disabled, nothing to prune
	if maxBackups <= 0 && maxAgeCutoff == "" {
		return victims, nil
	}

	// Build the WHERE clause to find victim IDs
	conditions := []string{}
	var args []any

	if maxBackups > 0 {
		// rows NOT in the top-N most recent (i.e., beyond the keep limit)
		conditions = append(conditions, `id NOT IN (
			SELECT id FROM backup_runs ORDER BY created_at DESC LIMIT ?
		)`)
		args = append(args, maxBackups)
	}

	if maxAgeCutoff != "" {
		// rows older than the cutoff
		conditions = append(conditions, "created_at < ?")
		args = append(args, maxAgeCutoff)
	}

	// Join conditions with OR
	where := ""
	if len(conditions) > 0 {
		where = "WHERE (" + strings.Join(conditions, " OR ") + ")"
	} else {
		return victims, nil // No conditions, nothing to prune
	}

	// Select victim rows first (so we can return them)
	selectQuery := `
		SELECT id, triggered_by, file_path, size_bytes, created_at
		FROM backup_runs
		` + where

	rows, err := s.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return victims, fmt.Errorf("select backup run victims: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var victimIDs []any
	for rows.Next() {
		var run BackupRun
		if err := rows.Scan(&run.ID, &run.TriggeredBy, &run.FilePath, &run.SizeBytes, &run.CreatedAt); err != nil {
			return victims, fmt.Errorf("scan victim: %w", err)
		}
		victims = append(victims, run)
		victimIDs = append(victimIDs, run.ID)
	}
	if err := rows.Err(); err != nil {
		return victims, fmt.Errorf("victim iteration: %w", err)
	}

	// If no victims, return early
	if len(victimIDs) == 0 {
		return victims, nil
	}

	// Delete the victim rows
	deleteQuery := `DELETE FROM backup_runs WHERE id IN (` +
		placeholders(len(victimIDs)) + `)`
	if _, err := s.db.ExecContext(ctx, deleteQuery, victimIDs...); err != nil {
		return victims, fmt.Errorf("delete backup run victims: %w", err)
	}

	return victims, nil
}

// placeholders returns a comma-separated string of ? placeholders for SQL IN clauses.
func placeholders(n int) string {
	if n == 0 {
		return ""
	}
	result := "?"
	for i := 1; i < n; i++ {
		result += ", ?"
	}
	return result
}

// existingIDs returns the subset of ids already present in table's id
// column, in one round-trip. table is always an internal literal (never
// user input), so it's safe to concatenate directly into the query.
func existingIDs(ctx context.Context, db DBTX, table string, ids []string) (map[string]bool, error) {
	existing := make(map[string]bool, len(ids))
	if len(ids) == 0 {
		return existing, nil
	}
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := db.QueryContext(ctx, `SELECT id FROM `+table+` WHERE id IN (`+placeholders(len(ids))+`)`, args...)
	if err != nil {
		return nil, fmt.Errorf("check existing %s: %w", table, err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan existing %s id: %w", table, err)
		}
		existing[id] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate existing %s: %w", table, err)
	}
	return existing, nil
}
