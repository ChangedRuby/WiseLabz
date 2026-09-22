package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/store"

	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

func newSeededStore(t *testing.T, dsn string) *store.Store {
	t.Helper()
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
	if err := s.CreateDoc(context.Background(), &store.DocRecord{ID: "doc-1", Title: "Runbook", Content: "steps"}); err != nil {
		t.Fatalf("create doc: %v", err)
	}
	return s
}

func TestRunVerifyPassAndFail(t *testing.T) {
	ctx := context.Background()
	srcDir := t.TempDir()
	s := newSeededStore(t, "file:"+srcDir+"/src.db?cache=shared")

	run, err := backup.ExportToFile(ctx, s, srcDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	if code := runVerify([]string{"-file", run.FilePath}); code != 0 {
		t.Errorf("runVerify(intact bundle) = %d, want 0", code)
	}

	// Corrupt the bundle in place and confirm verify now fails.
	data, err := os.ReadFile(run.FilePath)
	if err != nil {
		t.Fatalf("read bundle: %v", err)
	}
	data[len(data)/2] ^= 0xFF
	if err := os.WriteFile(run.FilePath, data, 0o600); err != nil {
		t.Fatalf("corrupt bundle: %v", err)
	}
	if code := runVerify([]string{"-file", run.FilePath}); code != 1 {
		t.Errorf("runVerify(corrupted bundle) = %d, want 1", code)
	}
}

func TestRunVerifyDefaultsToLatestInDir(t *testing.T) {
	ctx := context.Background()
	srcDir := t.TempDir()
	s := newSeededStore(t, "file:"+srcDir+"/src.db?cache=shared")

	if _, err := backup.ExportToFile(ctx, s, srcDir); err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	if code := runVerify([]string{"-dir", srcDir}); code != 0 {
		t.Errorf("runVerify(-dir) = %d, want 0", code)
	}
}

func TestRunVerifyRequiresABundle(t *testing.T) {
	if code := runVerify([]string{"-dir", t.TempDir()}); code != 1 {
		t.Errorf("runVerify(empty dir) = %d, want 1", code)
	}
}

func TestRunRestoreImportsIntoConfiguredDatabase(t *testing.T) {
	ctx := context.Background()
	srcDir := t.TempDir()
	s := newSeededStore(t, "file:"+srcDir+"/src.db?cache=shared")

	run, err := backup.ExportToFile(ctx, s, srcDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	targetDir := t.TempDir()
	dsn := fmt.Sprintf("file:%s/target.db?cache=shared", targetDir)
	t.Setenv("WISELABZ_DB_DRIVER", "sqlite")
	t.Setenv("WISELABZ_DB_DSN", dsn)

	if code := runRestore([]string{"-file", run.FilePath, "-yes"}); code != 0 {
		t.Fatalf("runRestore = %d, want 0", code)
	}

	// Confirm the doc actually landed in the target database.
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open target db: %v", err)
	}
	defer func() { _ = db.Close() }()

	target := store.New(db, "sqlite")
	doc, err := target.GetDoc(ctx, "doc-1")
	if err != nil {
		t.Fatalf("GetDoc after restore: %v", err)
	}
	if doc.Title != "Runbook" {
		t.Errorf("restored doc title = %q, want Runbook", doc.Title)
	}
}

func TestRunRestoreRejectsCorruptedBundle(t *testing.T) {
	ctx := context.Background()
	srcDir := t.TempDir()
	s := newSeededStore(t, "file:"+srcDir+"/src.db?cache=shared")

	run, err := backup.ExportToFile(ctx, s, srcDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}
	data, err := os.ReadFile(run.FilePath)
	if err != nil {
		t.Fatalf("read bundle: %v", err)
	}
	data[len(data)/2] ^= 0xFF
	if err := os.WriteFile(run.FilePath, data, 0o600); err != nil {
		t.Fatalf("corrupt bundle: %v", err)
	}

	targetDir := t.TempDir()
	t.Setenv("WISELABZ_DB_DRIVER", "sqlite")
	t.Setenv("WISELABZ_DB_DSN", fmt.Sprintf("file:%s/target.db?cache=shared", targetDir))

	if code := runRestore([]string{"-file", run.FilePath, "-yes"}); code != 1 {
		t.Errorf("runRestore(corrupted bundle) = %d, want 1", code)
	}
}

func TestRunRestoreRequiresFileFlag(t *testing.T) {
	if code := runRestore(nil); code != 1 {
		t.Errorf("runRestore(no -file) = %d, want 1", code)
	}
}
