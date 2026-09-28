// Package storetest hands tests a fresh, already-migrated SQLite database
// file. Running every migration per test dominated backend test time under
// -race (modernc.org/sqlite is pure Go, so the race detector instruments the
// whole engine: ~1s per test vs ~0.04s without -race). Instead, migrations
// run once per test binary and each test gets its own byte-for-byte copy of
// the result, so tests stay fully isolated. This package is test-only
// despite the non-_test name.
package storetest

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"

	_ "modernc.org/sqlite" // sqlite driver registration

	"github.com/WiseLabz/wiselabz/internal/store"
)

var (
	templateOnce sync.Once
	templateDB   []byte
	templateErr  error
)

// MigratedSQLite writes a fresh copy of a fully migrated SQLite database into
// t.TempDir() and returns its path. Callers open it with their own DSN
// options and pool settings, exactly as they previously opened an empty file.
func MigratedSQLite(t testing.TB) string {
	t.Helper()
	templateOnce.Do(func() { templateDB, templateErr = buildTemplate() })
	if templateErr != nil {
		t.Fatalf("build migrated sqlite template: %v", templateErr)
	}
	path := filepath.Join(t.TempDir(), "test.db")
	if err := os.WriteFile(path, templateDB, 0o600); err != nil {
		t.Fatalf("write migrated sqlite copy: %v", err)
	}
	return path
}

func buildTemplate() ([]byte, error) {
	dir, err := os.MkdirTemp("", "storetest-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) //nolint:errcheck
	path := filepath.Join(dir, "template.db")

	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := store.RunMigrations(db, "sqlite", logger); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	if err := db.Close(); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}
