package retention

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"

	_ "modernc.org/sqlite"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := store.RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}

	return store.New(db, "sqlite")
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// TestRunCleanupSkipsDisabledCategories verifies a *Days value of 0 leaves
// that category untouched, while an enabled category is still cleaned up.
func TestRunCleanupSkipsDisabledCategories(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	c := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)

	// sync_runs: retention enabled, should be purged.
	if err := s.CreateSyncRun(ctx, &store.SyncRunRecord{ConnectorID: c.ID, StartedAt: old, Status: store.SyncRunStatusSuccess}); err != nil {
		t.Fatalf("CreateSyncRun() error: %v", err)
	}
	// alerts: retention disabled (AlertDays: 0), should survive even though old + resolved.
	alert := &store.AlertRecord{ServiceID: c.ID, Severity: "info", Title: "t", Status: "resolved", CreatedAt: old}
	if err := s.CreateAlert(ctx, alert); err != nil {
		t.Fatalf("CreateAlert() error: %v", err)
	}
	// audit_log: retention enabled, should be purged.
	if err := s.CreateAuditRecord(ctx, &store.AuditRecord{ActorUserID: "u1", ActorRole: "operator", Action: "test.action", CreatedAt: old}); err != nil {
		t.Fatalf("CreateAuditRecord() error: %v", err)
	}

	cfg := store.RetentionSettings{
		SnapshotDays:   0,
		DocVersionDays: 0,
		AlertDays:      0, // disabled
		SyncRunDays:    30,
		AuditDays:      30,
		CronExpr:       "0 0 * * *",
	}

	RunCleanupOnce(ctx, s, cfg, testLogger())

	runs, err := s.ListSyncRunsByConnector(ctx, c.ID, 20)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if len(runs) != 0 {
		t.Fatalf("sync runs after cleanup = %d, want 0 (SyncRunDays enabled)", len(runs))
	}

	if _, err := s.GetAlert(ctx, alert.ID); err != nil {
		t.Fatalf("alert should survive when AlertDays=0 (disabled), GetAlert() error: %v", err)
	}

	_, auditTotal, err := s.ListAuditRecords(ctx, "", "", "", "", 0, 20)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if auditTotal != 0 {
		t.Fatalf("audit records after cleanup = %d, want 0 (AuditDays enabled)", auditTotal)
	}
}

// TestRunCleanupPrunesOldHealthChecks verifies health_checks rows older than
// the configured window are purged while recent ones survive, and that
// HealthCheckDays=0 disables cleanup for the category like the other
// *Days settings.
func TestRunCleanupPrunesOldHealthChecks(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	c := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	recent := time.Now().UTC().Format(time.RFC3339)

	oldCheck := &store.HealthCheckRecord{ConnectorID: c.ID, Status: "offline", CheckedAt: old}
	if err := s.RecordHealthCheck(ctx, oldCheck); err != nil {
		t.Fatalf("RecordHealthCheck() old error: %v", err)
	}
	recentCheck := &store.HealthCheckRecord{ConnectorID: c.ID, Status: "online", CheckedAt: recent}
	if err := s.RecordHealthCheck(ctx, recentCheck); err != nil {
		t.Fatalf("RecordHealthCheck() recent error: %v", err)
	}

	// HealthCheckDays=0 disables cleanup: nothing purged.
	RunCleanupOnce(ctx, s, store.RetentionSettings{HealthCheckDays: 0}, testLogger())
	stats, err := s.GetConnectorUptime(ctx, c.ID, time.Now().UTC().AddDate(-1, 0, -1), time.Now().UTC().AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}
	if stats.CheckCount != 2 {
		t.Fatalf("check count after disabled cleanup = %d, want 2", stats.CheckCount)
	}

	// HealthCheckDays=30 purges the year-old row, keeps the recent one.
	RunCleanupOnce(ctx, s, store.RetentionSettings{HealthCheckDays: 30}, testLogger())
	stats, err = s.GetConnectorUptime(ctx, c.ID, time.Now().UTC().AddDate(-1, 0, -1), time.Now().UTC().AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}
	if stats.CheckCount != 1 {
		t.Fatalf("check count after cleanup = %d, want 1", stats.CheckCount)
	}
}

// TestRunCleanupAllDBErrors verifies that when every delete call errors
// (closed DB), RunCleanupOnce logs and returns without panicking.
func TestRunCleanupAllDBErrors(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	if err := s.DB().Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	cfg := store.RetentionSettings{
		SnapshotDays: 30, DocVersionDays: 30, AlertDays: 30, SyncRunDays: 30, AuditDays: 30,
	}

	RunCleanupOnce(ctx, s, cfg, testLogger()) // must not panic
}

// TestRunCleanupPartialFailure verifies that an error in one category
// (sync_runs table dropped, simulating a DB error) does not stop the other
// enabled categories from running.
func TestRunCleanupPartialFailure(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	c := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	if err := s.CreateAuditRecord(ctx, &store.AuditRecord{ActorUserID: "u1", ActorRole: "operator", Action: "test.action", CreatedAt: old}); err != nil {
		t.Fatalf("CreateAuditRecord() error: %v", err)
	}

	if _, err := s.DB().ExecContext(ctx, `DROP TABLE sync_runs`); err != nil {
		t.Fatalf("drop sync_runs table: %v", err)
	}

	cfg := store.RetentionSettings{SyncRunDays: 30, AuditDays: 30}

	RunCleanupOnce(ctx, s, cfg, testLogger()) // sync_runs errors, audit_log must still run

	_, auditTotal, err := s.ListAuditRecords(ctx, "", "", "", "", 0, 20)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if auditTotal != 0 {
		t.Fatalf("audit records after cleanup = %d, want 0 (audit cleanup must still run after sync_runs error)", auditTotal)
	}
}

// TestRunCleanupIdempotent verifies a second pass with the same config
// deletes nothing further.
func TestRunCleanupIdempotent(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	c := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	if err := s.CreateSyncRun(ctx, &store.SyncRunRecord{ConnectorID: c.ID, StartedAt: old, Status: store.SyncRunStatusSuccess}); err != nil {
		t.Fatalf("CreateSyncRun() error: %v", err)
	}

	cfg := store.RetentionSettings{SyncRunDays: 30, CronExpr: "@daily"}

	RunCleanupOnce(ctx, s, cfg, testLogger())
	RunCleanupOnce(ctx, s, cfg, testLogger()) // must not error or panic on an already-clean table

	runs, err := s.ListSyncRunsByConnector(ctx, c.ID, 20)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if len(runs) != 0 {
		t.Fatalf("sync runs after second cleanup = %d, want 0", len(runs))
	}
}
