package docexport_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/WiseLabz/wiselabz/internal/docexport"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/store"

	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
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
		t.Fatalf("run migrations: %v", err)
	}

	s := store.New(db, "sqlite")
	if err := s.Init(context.Background(), "admin-seed-pw-1234"); err != nil {
		t.Fatalf("store init: %v", err)
	}
	return s
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %q: %v", path, err)
	}
	return string(data)
}

func TestExportAllWritesExpectedFilesAndContent(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Kind: "service", Content: "# Runbook\n\nsteps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-2", Title: "Lab Topology", Kind: "lab", Content: "# Lab Topology\n\ngraph"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	dir := t.TempDir()
	e := docexport.NewExporter(s)
	res, err := e.ExportAll(ctx, dir)
	if err != nil {
		t.Fatalf("ExportAll: %v", err)
	}

	if res.Count != 2 {
		t.Fatalf("Count = %d, want 2", res.Count)
	}
	if res.Dir != dir {
		t.Errorf("Dir = %q, want %q", res.Dir, dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("dir has %d entries, want 2", len(entries))
	}

	wantContent := map[string]string{
		"runbook-doc-1.md":      "# Runbook\n\nsteps",
		"lab-topology-doc-2.md": "# Lab Topology\n\ngraph",
	}
	for name, want := range wantContent {
		got := readFile(t, filepath.Join(dir, name))
		if got != want {
			t.Errorf("content of %q = %q, want %q", name, got, want)
		}
	}
}

func TestExportAllPrunesStaleFilesOnRerun(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "v1"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	dir := t.TempDir()
	e := docexport.NewExporter(s)
	if _, err := e.ExportAll(ctx, dir); err != nil {
		t.Fatalf("first ExportAll: %v", err)
	}

	// Delete the doc and add a different one; the old doc's file should be
	// removed on the next run instead of lingering as a stale copy.
	if err := s.DeleteDoc(ctx, "doc-1"); err != nil {
		t.Fatalf("delete doc: %v", err)
	}
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-2", Title: "New Doc", Content: "v2"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	res, err := e.ExportAll(ctx, dir)
	if err != nil {
		t.Fatalf("second ExportAll: %v", err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "runbook-doc-1.md" {
		t.Fatalf("Removed = %v, want [runbook-doc-1.md]", res.Removed)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "new-doc-doc-2.md" {
		t.Fatalf("dir entries = %v, want [new-doc-doc-2.md]", entries)
	}
}

func TestExportAllEmptyDirName(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	e := docexport.NewExporter(s)
	if _, err := e.ExportAll(ctx, ""); err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestRunExportOnceCreatesDirectory(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "steps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	dir := filepath.Join(t.TempDir(), "nested", "export")
	e := docexport.NewExporter(s)
	docexport.RunExportOnce(ctx, e, dir, slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("dir has %d entries, want 1", len(entries))
	}
}

func TestDocExportDefaultCronExprIsValid(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	r := scheduler.New(logger)

	if _, err := r.AddJob("docexport", docexport.DefaultCronExpr, func(_ context.Context) {}); err != nil {
		t.Fatalf("AddJob(docexport.DefaultCronExpr) error: %v", err)
	}
	if r.EntryCount() != 1 {
		t.Fatalf("EntryCount() = %d, want 1", r.EntryCount())
	}
}

// TestScheduledExportRunsAndFires registers the doc export job on a fast
// cron cadence (mirroring how cmd/server/main.go wires it against
// cfg.DocExport.CronExpr) and confirms it actually fires and exports.
func TestScheduledExportRunsAndFires(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "steps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	dir := t.TempDir()
	e := docexport.NewExporter(s)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	r := scheduler.New(logger)

	if _, err := r.AddJob("docexport", "*/1 * * * * *", func(jobCtx context.Context) {
		docexport.RunExportOnce(jobCtx, e, dir, logger)
	}); err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	runCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r.Start(runCtx)

	deadline := time.Now().Add(2 * time.Second)
	var entries []os.DirEntry
	for time.Now().Before(deadline) {
		var err error
		entries, err = os.ReadDir(dir)
		if err == nil && len(entries) == 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	r.Stop()

	if len(entries) != 1 {
		t.Fatalf("dir has %d entries after scheduled run, want 1", len(entries))
	}
}
