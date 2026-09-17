package store

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

// tablesCreatedByMigrations lists every table the initial migration creates,
// shared by the SQLite and PostgreSQL migration tests.
var tablesCreatedByMigrations = []string{
	"users", "sessions", "oidc_provider_flags", "connectors",
	"service_snapshots", "docs", "doc_versions", "templates",
	"template_sections", "changes", "alerts", "dashboard_layouts",
	"auth_config", "ai_config", "notification_config", "in_app_notifications",
	"quality_findings", "runbooks", "oidc_identities", "doc_locks",
	"api_keys",
}

func TestRunMigrations(t *testing.T) {
	// Use file-based SQLite so golang-migrate can track schema version.
	// :memory: won't work because golang-migrate uses a separate connection
	// for the schema_migrations table.
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}

	// Verify tables exist by querying each one
	for _, table := range tablesCreatedByMigrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Errorf("table %s does not exist or is not queryable: %v", table, err)
		}
	}

	// Verify idempotent — running again should be no-op
	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() second run error: %v", err)
	}

	var name string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='maintenance_windows'").Scan(&name); err != nil {
		t.Errorf("maintenance_windows table missing after migrations: %v", err)
	}
}

// TestRunMigrationsPostgres runs the postgres migration path against a real
// PostgreSQL instance. Opt-in: set WISELABZ_TEST_POSTGRES_DSN (e.g.
// "postgres://wiselabz:wiselabz@localhost:5432/wiselabz?sslmode=disable")
// to a database that RunMigrations is allowed to create tables in. Skipped
// otherwise, so `go test ./...` needs no Postgres instance by default.
func TestRunMigrationsPostgres(t *testing.T) {
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set; skipping postgres migration test")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	if err := RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}

	for _, table := range tablesCreatedByMigrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Errorf("table %s does not exist or is not queryable: %v", table, err)
		}
	}

	// Verify idempotent — running again should be no-op
	if err := RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations() second run error: %v", err)
	}
}

func TestRunMigrationsDown(t *testing.T) {
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}

	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() error: %v", err)
	}

	var count int
	for _, table := range tablesCreatedByMigrations {
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Errorf("table %s should still exist after rolling back only the last migration: %v", table, err)
		}
	}
	// Rolling back only the latest migration (share_links) should drop
	// share_links, without touching connector_permissions or anything earlier.
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='share_links'").Scan(&name)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("share_links should not exist after rolling back its migration (err=%v)", err)
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user_connector_roles'").Scan(&name)
	if err != nil {
		t.Errorf("user_connector_roles should still exist after rolling back only the latest migration (err=%v)", err)
	}

	// Rolling back a second time (connector_permissions) should drop
	// user_connector_roles and restore users.role, without touching earlier
	// migrations' tables/columns.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() second call error: %v", err)
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user_connector_roles'").Scan(&name)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("user_connector_roles should not exist after rolling back its migration (err=%v)", err)
	}
	var usersSchema string
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='users'").Scan(&usersSchema); err != nil {
		t.Fatalf("query sqlite_master for users: %v", err)
	}
	if !strings.Contains(usersSchema, "role") || strings.Contains(usersSchema, "instance_admin_role") {
		t.Errorf("users.role should be restored after rolling back connector_permissions, got schema: %s", usersSchema)
	}

	// maintenance_windows is from the migration before connector_permissions,
	// so rolling back the last two migrations must leave it in place.
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='maintenance_windows'").Scan(&name)
	if err != nil {
		t.Errorf("maintenance_windows should still exist after rolling back the last two migrations (err=%v)", err)
	}

	// dns connector category is from an earlier migration, so rolling back
	// only the latest migration must leave it in place.
	var connectorsSchema string
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='connectors'").Scan(&connectorsSchema); err != nil {
		t.Fatalf("query sqlite_master for connectors: %v", err)
	}
	if !strings.Contains(connectorsSchema, "'dns'") {
		t.Error("connectors.category should still allow 'dns' after rolling back only the latest migration")
	}

	// ai_config_providers is from an earlier migration, so rolling back only
	// the latest migration must leave it in place.
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='ai_config_providers'").Scan(&name)
	if err != nil {
		t.Errorf("ai_config_providers should still exist after rolling back only the latest migration (err=%v)", err)
	}

	var changesSchema string
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='changes'").Scan(&changesSchema); err != nil {
		t.Fatalf("query sqlite_master for changes: %v", err)
	}
	if !strings.Contains(changesSchema, "narration") {
		t.Error("changes should still have narration from an earlier migration")
	}
	for _, table := range []string{"doc_section_embeddings", "chat_conversations", "chat_messages"} {
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %s from an earlier migration should still exist (err=%v)", table, err)
		}
	}
	var aiConfigSchema string
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='ai_config'").Scan(&aiConfigSchema); err != nil {
		t.Fatalf("query sqlite_master for ai_config: %v", err)
	}
	if !strings.Contains(aiConfigSchema, "embed_provider") {
		t.Error("ai_config should still have embed_provider from an earlier migration")
	}
}

// TestRunMigrationsDownPostgres mirrors TestRunMigrationsPostgres but for the
// down path. Opt-in via WISELABZ_TEST_POSTGRES_DSN; skipped otherwise.
func TestRunMigrationsDownPostgres(t *testing.T) {
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set; skipping postgres migration test")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	if err := RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}

	if err := RunMigrationsDown(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrationsDown() error: %v", err)
	}

	var count int
	for _, table := range tablesCreatedByMigrations {
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Errorf("table %s should still exist after rolling back only the last migration: %v", table, err)
		}
	}
	if !hasColumn(t, db, "postgres", "changes", "narration") {
		t.Error("changes.narration should still exist from an earlier migration")
	}
	var name string
	err = db.QueryRow(`SELECT table_name FROM information_schema.tables WHERE table_name = 'ai_config_providers'`).Scan(&name)
	if err == nil {
		t.Error("ai_config_providers should not exist after rolling back its migration")
	} else if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("query information_schema.tables for ai_config_providers: %v", err)
	}
}

func TestSessionAuthProviderMigration(t *testing.T) {
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}
	if !hasColumn(t, db, "sqlite", "sessions", "auth_provider_id") {
		t.Fatal("sessions.auth_provider_id column missing after migrations")
	}
}

func hasColumn(t *testing.T, db *sql.DB, driver, table, column string) bool {
	t.Helper()
	if driver == "postgres" {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*)
			FROM information_schema.columns
			WHERE table_name = $1 AND column_name = $2
		`, table, column).Scan(&count)
		if err != nil {
			t.Fatalf("information_schema.columns(%s.%s): %v", table, column, err)
		}
		return count > 0
	}

	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		t.Fatalf("PRAGMA table_info(%s): %v", table, err)
	}
	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			t.Fatalf("scan table info: %v", err)
		}
		if name == column {
			return true
		}
	}
	return false
}

func TestRunMigrationsDownUnsupportedDriver(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	err = RunMigrationsDown(db, "mysql", logger)
	if err == nil {
		t.Error("expected error for unsupported driver, got nil")
	}
}

func TestRunMigrationsUnsupportedDriver(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	err = RunMigrations(db, "mysql", logger)
	if err == nil {
		t.Error("expected error for unsupported driver, got nil")
	}
}
