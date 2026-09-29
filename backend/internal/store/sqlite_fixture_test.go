package store

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
)

var (
	sqliteTemplateOnce sync.Once
	sqliteTemplateDB   []byte
	sqliteTemplateErr  error
)

// migratedSQLite writes a fresh copy of a fully migrated SQLite database into
// t.TempDir() and returns its path. Callers open it with their own DSN
// options and pool settings, exactly as they previously opened an empty file.
func migratedSQLite(t testing.TB) string {
	t.Helper()
	sqliteTemplateOnce.Do(func() { sqliteTemplateDB, sqliteTemplateErr = buildSQLiteTemplate() })
	if sqliteTemplateErr != nil {
		t.Fatalf("build migrated sqlite template: %v", sqliteTemplateErr)
	}
	path := filepath.Join(t.TempDir(), "test.db")
	if err := os.WriteFile(path, sqliteTemplateDB, 0o600); err != nil {
		t.Fatalf("write migrated sqlite copy: %v", err)
	}
	return path
}

func buildSQLiteTemplate() ([]byte, error) {
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
	if err := RunMigrations(db, "sqlite", logger); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	if err := db.Close(); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// MigratedSQLiteForTest lets external store_test tests share the package template.
// This bridge exists only in test binaries.
func MigratedSQLiteForTest(t testing.TB) string { return migratedSQLite(t) }
