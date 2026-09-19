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
	"api_keys", "compliance_rules",
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
	var disabledRules int
	if err := db.QueryRow("SELECT COUNT(*) FROM compliance_rules WHERE enabled = 0").Scan(&disabledRules); err != nil || disabledRules != 5 {
		t.Errorf("disabled seeded compliance rules = %d, %v; want 5, nil", disabledRules, err)
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

	// Roll back the newest migration (hot_query_indexes) first; it must drop
	// only its indexes.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() hot_query_indexes error: %v", err)
	}
	var hotIdx string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_alerts_status_created'").Scan(&hotIdx); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("idx_alerts_status_created should not exist after rolling back its migration (err=%v)", err)
	}

	// Roll back the next migration (snapshot_fetched_at_index) first; it
	// must drop only its index.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() snapshot_fetched_at_index error: %v", err)
	}
	var idxName string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_snapshots_fetched_at'").Scan(&idxName); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("idx_snapshots_fetched_at should not exist after rolling back its migration (err=%v)", err)
	}

	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() error: %v", err)
	}

	// Rolling back only the latest migration (user_digest_prefs) should drop
	// the digest columns without touching compliance_rules.
	if hasColumn(t, db, "sqlite", "users", "digest_cadence") {
		t.Error("users.digest_cadence should not exist after rolling back user_digest_prefs")
	}
	var complianceCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM compliance_rules").Scan(&complianceCount); err != nil {
		t.Errorf("compliance_rules should still exist after rolling back only user_digest_prefs: %v", err)
	}
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() second call (compliance_rules) error: %v", err)
	}

	var count int
	for _, table := range tablesCreatedByMigrations {
		if table == "compliance_rules" {
			continue
		}
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Errorf("table %s should still exist after rolling back only the last migration: %v", table, err)
		}
	}
	// Rolling back only the latest migration (compliance_rules) should leave
	// the earlier rotation migration intact.
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='compliance_rules'").Scan(&name)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("compliance_rules should not exist after rolling back its migration (err=%v)", err)
	}
	// Verify the latest migrations can be reapplied cleanly after their down
	// path. RunMigrations reapplies every pending migration, so this brings
	// back both compliance_rules and user_digest_prefs.
	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() after compliance down error: %v", err)
	}
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='compliance_rules'").Scan(&name); err != nil {
		t.Fatalf("compliance_rules missing after reapply: %v", err)
	}
	if !hasColumn(t, db, "sqlite", "users", "digest_cadence") {
		t.Fatal("users.digest_cadence missing after reapply")
	}
	// Four down calls strip the reapplied migrations (hot indexes, snapshot
	// index, digest prefs, compliance_rules) back off, returning to the same "compliance_rules and
	// user_digest_prefs absent" state as before the reapply.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() after reapply, hot indexes call error: %v", err)
	}
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() after reapply, index call error: %v", err)
	}
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() after reapply, first call error: %v", err)
	}
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() after reapply, second call error: %v", err)
	}
	var connectorsSchemaAfterFirst string
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='connectors'").Scan(&connectorsSchemaAfterFirst); err != nil {
		t.Fatalf("query sqlite_master for connectors: %v", err)
	}
	if !strings.Contains(connectorsSchemaAfterFirst, "secret_rotated_at") {
		t.Error("connectors.secret_rotated_at should remain after rolling back compliance_rules")
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='share_links'").Scan(&name)
	if err != nil {
		t.Errorf("share_links should still exist after rolling back only the latest migration (err=%v)", err)
	}

	// Rolling back a second time (credential_rotation) should drop rotation
	// columns without touching share_links.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() second call error: %v", err)
	}
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='connectors'").Scan(&connectorsSchemaAfterFirst); err != nil {
		t.Fatalf("query sqlite_master for connectors: %v", err)
	}
	if strings.Contains(connectorsSchemaAfterFirst, "secret_rotated_at") {
		t.Error("connectors.secret_rotated_at should not exist after rolling back credential_rotation")
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user_connector_roles'").Scan(&name)
	if err != nil {
		t.Errorf("user_connector_roles should still exist after rolling back share_links (err=%v)", err)
	}

	// Rolling back a third time (share_links) should drop share_links without
	// touching connector permissions.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() third call error: %v", err)
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='share_links'").Scan(&name)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("share_links should not exist after rolling back its migration (err=%v)", err)
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user_connector_roles'").Scan(&name)
	if err != nil {
		t.Errorf("user_connector_roles should still exist after rolling back share_links (err=%v)", err)
	}

	// Rolling back a fourth time (connector_permissions) should drop roles and
	// restore users.role, without touching earlier migrations.
	if err := RunMigrationsDown(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrationsDown() fourth call error: %v", err)
	}
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='user_connector_roles'").Scan(&name)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("user_connector_roles should not exist after rolling back connector_permissions (err=%v)", err)
	}
	var usersSchema string
	if err := db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='users'").Scan(&usersSchema); err != nil {
		t.Fatalf("query sqlite_master for users: %v", err)
	}
	if !strings.Contains(usersSchema, "role") || strings.Contains(usersSchema, "instance_admin_role") {
		t.Errorf("users.role should be restored after rolling back connector_permissions, got schema: %s", usersSchema)
	}

	// maintenance_windows is from the migration before connector_permissions,
	// so rolling back the last three migrations must leave it in place.
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='maintenance_windows'").Scan(&name)
	if err != nil {
		t.Errorf("maintenance_windows should still exist after rolling back the last three migrations (err=%v)", err)
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
	// Only the newest migration (hot_query_indexes) is rolled back: its indexes
	// must be gone while snapshot_fetched_at_index and ai_config_providers from
	// earlier migrations stay.
	var name string
	err = db.QueryRow(`SELECT indexname FROM pg_indexes WHERE schemaname = current_schema() AND indexname = 'idx_alerts_status_created'`).Scan(&name)
	if err == nil {
		t.Error("idx_alerts_status_created should not exist after rolling back its migration")
	} else if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("query pg_indexes for idx_alerts_status_created: %v", err)
	}
	if err := db.QueryRow(`SELECT indexname FROM pg_indexes WHERE schemaname = current_schema() AND indexname = 'idx_snapshots_fetched_at'`).Scan(&name); err != nil {
		t.Errorf("idx_snapshots_fetched_at should still exist after rolling back only the last migration: %v", err)
	}
	if err := db.QueryRow(`SELECT table_name FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'ai_config_providers'`).Scan(&name); err != nil {
		t.Errorf("ai_config_providers should still exist after rolling back only the last migration: %v", err)
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

// TestRunMigrationsPreservesRowsWithForeignKeys upgrades a populated database
// (opened the way production does, with foreign_keys=ON) across the
// DROP TABLE rebuild migrations and checks that no ON DELETE CASCADE / SET NULL
// action wiped or nulled child rows (#302).
func TestRunMigrationsPreservesRowsWithForeignKeys(t *testing.T) {
	db, err := OpenDB("sqlite", "file:"+t.TempDir()+"/upgrade.db")
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer db.Close() //nolint:errcheck

	m, err := newMigrator(db, "sqlite")
	if err != nil {
		t.Fatalf("newMigrator: %v", err)
	}
	// Last version before the first rebuild migration (000022).
	if err := m.Migrate(21); err != nil {
		t.Fatalf("migrate to 21: %v", err)
	}

	seed := []string{
		`INSERT INTO users (id, username, role, created_at) VALUES ('u1', 'alice', 'operator', 'now')`,
		`INSERT INTO sessions (id, user_id, token_hash, created_at, last_seen_at) VALUES ('s1', 'u1', 'h', 'now', 'now')`,
		`INSERT INTO connectors (id, name, category, type, url, created_at, updated_at) VALUES ('c1', 'pve', 'virtualization', 'proxmox', 'http://x', 'now', 'now')`,
		`INSERT INTO service_snapshots (id, connector_id, data, fetched_at) VALUES ('n1', 'c1', '{}', 'now')`,
		`INSERT INTO sync_runs (id, connector_id, started_at, status) VALUES ('r1', 'c1', 'now', 'success')`,
		`INSERT INTO docs (id, title, kind, service_id, created_at, updated_at) VALUES ('d1', 'doc', 'service', 'c1', 'now', 'now')`,
		`INSERT INTO quality_findings (id, connector_id, doc_id, check_type, severity, title, first_detected_at, last_seen_at) VALUES ('q1', 'c1', 'd1', 'stale', 'info', 't', 'now', 'now')`,
	}
	for _, q := range seed {
		if _, err := db.Exec(q); err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations: %v", err)
	}

	for table, want := range map[string]int{
		"users": 1, "sessions": 1, "connectors": 1, "service_snapshots": 1,
		"sync_runs": 1, "docs": 1, "quality_findings": 1, "user_connector_roles": 1,
	} {
		var got int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&got); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if got != want {
			t.Errorf("%s rows = %d, want %d", table, got, want)
		}
	}
	var serviceID sql.NullString
	if err := db.QueryRow("SELECT service_id FROM docs WHERE id = 'd1'").Scan(&serviceID); err != nil || !serviceID.Valid {
		t.Errorf("docs.service_id = %v, %v; want it preserved", serviceID, err)
	}
	var docID sql.NullString
	if err := db.QueryRow("SELECT doc_id FROM quality_findings WHERE id = 'q1'").Scan(&docID); err != nil || !docID.Valid {
		t.Errorf("quality_findings.doc_id = %v, %v; want it preserved", docID, err)
	}

	var fk int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&fk); err != nil || fk != 1 {
		t.Errorf("foreign_keys after migrate = %d, %v; want 1", fk, err)
	}
}

func TestGetMigrationStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/status.db?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close() //nolint:errcheck
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	st, err := GetMigrationStatus(db, "sqlite")
	if err != nil {
		t.Fatalf("GetMigrationStatus() before migrate: %v", err)
	}
	if st.Current != 0 || st.Latest == 0 || !st.Pending() {
		t.Errorf("before migrate = %+v, want current 0 and pending", st)
	}

	if err := RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("RunMigrations() error: %v", err)
	}
	st, err = GetMigrationStatus(db, "sqlite")
	if err != nil {
		t.Fatalf("GetMigrationStatus() after migrate: %v", err)
	}
	if st.Current != st.Latest || st.Dirty || st.Pending() {
		t.Errorf("after migrate = %+v, want current == latest, clean", st)
	}
}
