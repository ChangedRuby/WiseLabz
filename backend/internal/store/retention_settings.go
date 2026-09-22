package store

import (
	"context"
	"fmt"
	"time"
)

// RetentionSettings represents the data-retention cleanup configuration
// (single-row table with id='default'). Mirrors BackupSchedule's pattern.
type RetentionSettings struct {
	SnapshotDays    int
	DocVersionDays  int
	AlertDays       int
	SyncRunDays     int
	AuditDays       int
	HealthCheckDays int
	CronExpr        string
	UpdatedAt       string
}

// GetRetentionSettings retrieves the retention settings from the database.
func (s *Store) GetRetentionSettings(ctx context.Context) (RetentionSettings, error) {
	var rs RetentionSettings
	err := s.db.QueryRowContext(ctx, `
		SELECT snapshot_days, doc_version_days, alert_days, sync_run_days, audit_days, health_check_days, cron_expr, updated_at
		FROM retention_settings WHERE id = 'default'
	`).Scan(&rs.SnapshotDays, &rs.DocVersionDays, &rs.AlertDays, &rs.SyncRunDays, &rs.AuditDays, &rs.HealthCheckDays, &rs.CronExpr, &rs.UpdatedAt)
	if err != nil {
		return rs, fmt.Errorf("get retention settings: %w", err)
	}
	return rs, nil
}

// UpsertRetentionSettings inserts or updates the retention settings using
// INSERT OR REPLACE / ON CONFLICT semantics (dialect-neutral). The id is
// always 'default'.
func (s *Store) UpsertRetentionSettings(ctx context.Context, rs RetentionSettings) error {
	if rs.UpdatedAt == "" {
		rs.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO retention_settings (id, snapshot_days, doc_version_days, alert_days, sync_run_days, audit_days, health_check_days, cron_expr, updated_at)
		VALUES ('default', ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			snapshot_days = excluded.snapshot_days,
			doc_version_days = excluded.doc_version_days,
			alert_days = excluded.alert_days,
			sync_run_days = excluded.sync_run_days,
			audit_days = excluded.audit_days,
			health_check_days = excluded.health_check_days,
			cron_expr = excluded.cron_expr,
			updated_at = excluded.updated_at
	`, rs.SnapshotDays, rs.DocVersionDays, rs.AlertDays, rs.SyncRunDays, rs.AuditDays, rs.HealthCheckDays, rs.CronExpr, rs.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert retention settings: %w", err)
	}
	return nil
}
