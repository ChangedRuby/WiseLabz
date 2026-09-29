package store_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func newBackupTestStore(t *testing.T) *store.Store {
	t.Helper()
	dsn := "file:" + store.MigratedSQLiteForTest(t) + "?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	s := store.New(db, "sqlite")
	if err := s.Init(context.Background(), "admin-seed-pw-1234"); err != nil {
		t.Fatalf("store init: %v", err)
	}
	return s
}

// TestGetBackupScheduleWhenNotExists returns an error when the schedule row doesn't exist yet.
func TestGetBackupScheduleWhenNotExists(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	_, err := s.GetBackupSchedule(ctx)
	if err == nil {
		t.Fatal("expected error when schedule doesn't exist, got nil")
	}
}

// TestUpsertBackupScheduleCreatesRow inserts a new schedule.
func TestUpsertBackupScheduleCreatesRow(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	sched := store.BackupSchedule{
		CronExpr:    "0 3 * * *",
		MaxBackups:  14,
		MaxAgeHours: 720,
		Enabled:     true,
	}

	if err := s.UpsertBackupSchedule(ctx, sched); err != nil {
		t.Fatalf("UpsertBackupSchedule: %v", err)
	}

	// Verify we can read it back
	retrieved, err := s.GetBackupSchedule(ctx)
	if err != nil {
		t.Fatalf("GetBackupSchedule after insert: %v", err)
	}

	if retrieved.CronExpr != sched.CronExpr {
		t.Errorf("CronExpr mismatch: got %s, want %s", retrieved.CronExpr, sched.CronExpr)
	}
	if retrieved.MaxBackups != sched.MaxBackups {
		t.Errorf("MaxBackups mismatch: got %d, want %d", retrieved.MaxBackups, sched.MaxBackups)
	}
	if retrieved.MaxAgeHours != sched.MaxAgeHours {
		t.Errorf("MaxAgeHours mismatch: got %d, want %d", retrieved.MaxAgeHours, sched.MaxAgeHours)
	}
	if retrieved.Enabled != sched.Enabled {
		t.Errorf("Enabled mismatch: got %v, want %v", retrieved.Enabled, sched.Enabled)
	}
}

// TestUpsertBackupScheduleUpdatesRow calls upsert twice and verifies the second updates.
func TestUpsertBackupScheduleUpdatesRow(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	first := store.BackupSchedule{
		CronExpr:    "0 3 * * *",
		MaxBackups:  14,
		MaxAgeHours: 720,
		Enabled:     true,
	}
	if err := s.UpsertBackupSchedule(ctx, first); err != nil {
		t.Fatalf("first UpsertBackupSchedule: %v", err)
	}

	second := store.BackupSchedule{
		CronExpr:    "0 4 * * *",
		MaxBackups:  7,
		MaxAgeHours: 360,
		Enabled:     false,
	}
	if err := s.UpsertBackupSchedule(ctx, second); err != nil {
		t.Fatalf("second UpsertBackupSchedule: %v", err)
	}

	retrieved, err := s.GetBackupSchedule(ctx)
	if err != nil {
		t.Fatalf("GetBackupSchedule after update: %v", err)
	}

	if retrieved.CronExpr != second.CronExpr {
		t.Errorf("CronExpr not updated: got %s, want %s", retrieved.CronExpr, second.CronExpr)
	}
	if retrieved.MaxBackups != second.MaxBackups {
		t.Errorf("MaxBackups not updated: got %d, want %d", retrieved.MaxBackups, second.MaxBackups)
	}
	if retrieved.Enabled != second.Enabled {
		t.Errorf("Enabled not updated: got %v, want %v", retrieved.Enabled, second.Enabled)
	}
}

// TestCreateBackupRun inserts a backup run.
func TestCreateBackupRun(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	run := store.BackupRun{
		ID:          uuid.New().String(),
		TriggeredBy: "schedule",
		FilePath:    "/path/to/backup.json",
		SizeBytes:   12345,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.CreateBackupRun(ctx, run); err != nil {
		t.Fatalf("CreateBackupRun: %v", err)
	}

	// Verify we can list it
	runs, total, err := s.ListBackupRuns(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns: %v", err)
	}

	if total != 1 {
		t.Errorf("total count: got %d, want 1", total)
	}
	if len(runs) != 1 {
		t.Errorf("runs length: got %d, want 1", len(runs))
	}
	if runs[0].ID != run.ID {
		t.Errorf("ID mismatch: got %s, want %s", runs[0].ID, run.ID)
	}
}

// TestListBackupRunsPaginated tests pagination.
func TestListBackupRunsPaginated(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	// Create 5 runs
	for i := 0; i < 5; i++ {
		run := store.BackupRun{
			ID:          uuid.New().String(),
			TriggeredBy: "schedule",
			FilePath:    "/path/to/backup.json",
			SizeBytes:   int64(i) * 1000,
			CreatedAt:   time.Now().UTC().Add(time.Duration(-i) * time.Minute).Format(time.RFC3339),
		}
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun %d: %v", i, err)
		}
	}

	// List first 2
	runs, total, err := s.ListBackupRuns(ctx, 2, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns page 1: %v", err)
	}

	if total != 5 {
		t.Errorf("total: got %d, want 5", total)
	}
	if len(runs) != 2 {
		t.Errorf("page 1 length: got %d, want 2", len(runs))
	}

	// List next 2
	runs2, _, err := s.ListBackupRuns(ctx, 2, 2)
	if err != nil {
		t.Fatalf("ListBackupRuns page 2: %v", err)
	}

	if len(runs2) != 2 {
		t.Errorf("page 2 length: got %d, want 2", len(runs2))
	}

	// Verify they're ordered by created_at DESC (most recent first)
	if runs[0].CreatedAt < runs[1].CreatedAt {
		t.Error("runs not sorted DESC by created_at")
	}
}

// TestPruneBackupRunsByCount tests pruning by max backup count.
func TestPruneBackupRunsByCount(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	// Create 5 runs with staggered timestamps
	for i := 0; i < 5; i++ {
		runID := uuid.New().String()
		run := store.BackupRun{
			ID:          runID,
			TriggeredBy: "schedule",
			FilePath:    "/path/to/backup-" + runID + ".json",
			SizeBytes:   int64(i) * 1000,
			CreatedAt:   time.Now().UTC().Add(time.Duration(-i*10) * time.Minute).Format(time.RFC3339),
		}
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun %d: %v", i, err)
		}
	}

	// Prune to keep only 2, no age limit
	pruned, err := s.PruneBackupRuns(ctx, 2, "")
	if err != nil {
		t.Fatalf("PruneBackupRuns: %v", err)
	}

	if len(pruned) != 3 {
		t.Errorf("pruned count: got %d, want 3", len(pruned))
	}

	// Verify remaining runs are the 2 most recent
	remaining, total, err := s.ListBackupRuns(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns after prune: %v", err)
	}

	if total != 2 {
		t.Errorf("total remaining: got %d, want 2", total)
	}
	if len(remaining) != 2 {
		t.Errorf("remaining runs: got %d, want 2", len(remaining))
	}
}

// TestPruneBackupRunsByAge tests pruning by age cutoff.
func TestPruneBackupRunsByAge(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	// Create 3 runs: one old, one recent, one in between
	oldTime := time.Now().UTC().Add(-24 * time.Hour)
	midTime := time.Now().UTC().Add(-12 * time.Hour)
	newTime := time.Now().UTC()

	oldRun := store.BackupRun{
		ID:          uuid.New().String(),
		TriggeredBy: "schedule",
		FilePath:    "/path/to/old.json",
		SizeBytes:   1000,
		CreatedAt:   oldTime.Format(time.RFC3339),
	}
	midRun := store.BackupRun{
		ID:          uuid.New().String(),
		TriggeredBy: "schedule",
		FilePath:    "/path/to/mid.json",
		SizeBytes:   2000,
		CreatedAt:   midTime.Format(time.RFC3339),
	}
	newRun := store.BackupRun{
		ID:          uuid.New().String(),
		TriggeredBy: "schedule",
		FilePath:    "/path/to/new.json",
		SizeBytes:   3000,
		CreatedAt:   newTime.Format(time.RFC3339),
	}

	for _, run := range []store.BackupRun{oldRun, midRun, newRun} {
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun: %v", err)
		}
	}

	// Prune runs older than 18 hours ago
	cutoff := time.Now().UTC().Add(-18 * time.Hour).Format(time.RFC3339)
	pruned, err := s.PruneBackupRuns(ctx, 0, cutoff) // maxBackups=0 means no limit by count
	if err != nil {
		t.Fatalf("PruneBackupRuns: %v", err)
	}

	if len(pruned) != 1 {
		t.Errorf("pruned count: got %d, want 1", len(pruned))
	}
	if pruned[0].ID != oldRun.ID {
		t.Error("wrong run was pruned (should be old)")
	}

	// Verify mid and new runs remain
	remaining, _, err := s.ListBackupRuns(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns after prune: %v", err)
	}

	if len(remaining) != 2 {
		t.Errorf("remaining runs: got %d, want 2", len(remaining))
	}
}

// TestPruneBackupRunsNoLimits returns empty list when no limits are set.
func TestPruneBackupRunsNoLimits(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	// Create 3 runs
	for i := 0; i < 3; i++ {
		run := store.BackupRun{
			ID:          uuid.New().String(),
			TriggeredBy: "schedule",
			FilePath:    "/path/to/backup.json",
			SizeBytes:   1000,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		}
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun: %v", err)
		}
	}

	// Prune with no limits
	pruned, err := s.PruneBackupRuns(ctx, 0, "")
	if err != nil {
		t.Fatalf("PruneBackupRuns: %v", err)
	}

	if len(pruned) != 0 {
		t.Errorf("pruned count: got %d, want 0", len(pruned))
	}

	// Verify all runs still exist
	_, total, err := s.ListBackupRuns(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns: %v", err)
	}

	if total != 3 {
		t.Errorf("total: got %d, want 3", total)
	}
}

// TestPruneBackupRunsNegativeMaxBackups verifies a negative maxBackups is
// treated the same as 0 — "no limit by count" — rather than deleting
// everything or panicking on a negative SQL LIMIT.
func TestPruneBackupRunsNegativeMaxBackups(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	for i := 0; i < 3; i++ {
		run := store.BackupRun{
			ID:          uuid.New().String(),
			TriggeredBy: "schedule",
			FilePath:    "/path/to/backup.json",
			SizeBytes:   1000,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		}
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun: %v", err)
		}
	}

	pruned, err := s.PruneBackupRuns(ctx, -5, "")
	if err != nil {
		t.Fatalf("PruneBackupRuns with negative maxBackups: %v", err)
	}
	if len(pruned) != 0 {
		t.Errorf("pruned count: got %d, want 0 (negative maxBackups must disable count-based pruning)", len(pruned))
	}

	_, total, err := s.ListBackupRuns(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns: %v", err)
	}
	if total != 3 {
		t.Errorf("total: got %d, want 3 (nothing should have been pruned)", total)
	}
}

// TestUpsertBackupSchedulePostgresParity runs the same insert-then-update
// round trip as TestUpsertBackupScheduleUpdatesRow against a real Postgres
// database, so the ON CONFLICT upsert (rewritten from `?` to `$1..$n` for
// postgres, see pgdb.go) is verified on both drivers rather than sqlite only.
func TestUpsertBackupSchedulePostgresParity(t *testing.T) {
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set; skipping postgres backup schedule test")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer db.Close() //nolint:errcheck
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := store.RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations(postgres) error: %v", err)
	}
	s := store.New(db, "postgres")

	first := store.BackupSchedule{CronExpr: "0 3 * * *", MaxBackups: 14, MaxAgeHours: 720, Enabled: true}
	if err := s.UpsertBackupSchedule(context.Background(), first); err != nil {
		t.Fatalf("first postgres upsert error: %v", err)
	}
	second := store.BackupSchedule{CronExpr: "0 4 * * *", MaxBackups: 7, MaxAgeHours: 360, Enabled: false}
	if err := s.UpsertBackupSchedule(context.Background(), second); err != nil {
		t.Fatalf("second postgres upsert error: %v", err)
	}

	got, err := s.GetBackupSchedule(context.Background())
	if err != nil {
		t.Fatalf("GetBackupSchedule(postgres) error: %v", err)
	}
	if got.CronExpr != second.CronExpr || got.MaxBackups != second.MaxBackups || got.Enabled != second.Enabled {
		t.Fatalf("postgres upserted schedule = %+v, want %+v", got, second)
	}
}

// TestPruneBackupRunsCombinedLimits prunes using both count and age limits.
func TestPruneBackupRunsCombinedLimits(t *testing.T) {
	ctx := context.Background()
	s := newBackupTestStore(t)

	// Create 5 runs: 2 old, 3 new
	oldTime := time.Now().UTC().Add(-24 * time.Hour)
	newTime := time.Now().UTC()

	for i := 0; i < 2; i++ {
		run := store.BackupRun{
			ID:          uuid.New().String(),
			TriggeredBy: "schedule",
			FilePath:    "/path/to/old-" + string(rune(i)) + ".json",
			SizeBytes:   1000,
			CreatedAt:   oldTime.Format(time.RFC3339),
		}
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun old: %v", err)
		}
	}

	for i := 0; i < 3; i++ {
		run := store.BackupRun{
			ID:          uuid.New().String(),
			TriggeredBy: "schedule",
			FilePath:    "/path/to/new-" + string(rune(i)) + ".json",
			SizeBytes:   2000,
			CreatedAt:   newTime.Format(time.RFC3339),
		}
		if err := s.CreateBackupRun(ctx, run); err != nil {
			t.Fatalf("CreateBackupRun new: %v", err)
		}
	}

	// Prune: keep only 2 most recent, or anything newer than 12 hours ago
	cutoff := time.Now().UTC().Add(-12 * time.Hour).Format(time.RFC3339)
	pruned, err := s.PruneBackupRuns(ctx, 2, cutoff)
	if err != nil {
		t.Fatalf("PruneBackupRuns: %v", err)
	}

	// Should prune: 2 old runs (beyond age limit) + 1 new run (beyond count limit)
	if len(pruned) != 3 {
		t.Errorf("pruned count: got %d, want 3", len(pruned))
	}
}
