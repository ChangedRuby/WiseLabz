package store

import (
	"database/sql"
	"io/fs"
	"log/slog"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func migrationFiles(t *testing.T, fsys fs.FS, dir string) []string {
	t.Helper()
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names
}

// TestMigrationVersionParity fails when the sqlite and postgres migration
// trees drift: every version needs the same name plus up and down files in both.
func TestMigrationVersionParity(t *testing.T) {
	sqliteFiles := migrationFiles(t, sqliteMigrations, "migrations/sqlite")
	postgresFiles := migrationFiles(t, postgresMigrations, "migrations/postgres")
	if len(sqliteFiles) == 0 {
		t.Fatal("no sqlite migrations found")
	}
	if !reflect.DeepEqual(sqliteFiles, postgresFiles) {
		t.Errorf("migration trees differ\nsqlite:   %v\npostgres: %v", sqliteFiles, postgresFiles)
	}
	ups, downs := 0, 0
	for _, n := range sqliteFiles {
		switch {
		case strings.HasSuffix(n, ".up.sql"):
			ups++
		case strings.HasSuffix(n, ".down.sql"):
			downs++
		default:
			t.Errorf("unexpected migration file %q", n)
		}
	}
	if ups != downs {
		t.Errorf("%d up migrations but %d down migrations", ups, downs)
	}
}

// schemaColumns returns "table.column" for every user table column.
func sqliteSchemaColumns(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT m.name, p.name FROM sqlite_master m, pragma_table_info(m.name) p
		WHERE m.type='table' AND m.name NOT LIKE 'sqlite_%' AND m.name <> 'schema_migrations'`)
	return collectColumns(t, rows, err)
}

func postgresSchemaColumns(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT table_name, column_name FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name <> 'schema_migrations'`)
	return collectColumns(t, rows, err)
}

func collectColumns(t *testing.T, rows *sql.Rows, err error) []string {
	t.Helper()
	if err != nil {
		t.Fatalf("query schema: %v", err)
	}
	defer rows.Close() //nolint:errcheck
	var out []string
	for rows.Next() {
		var tbl, col string
		if err := rows.Scan(&tbl, &col); err != nil {
			t.Fatalf("scan: %v", err)
		}
		out = append(out, tbl+"."+col)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	sort.Strings(out)
	return out
}

// TestMigrationSchemaParity migrates both dialects and compares the resulting
// table/column sets. Opt-in via WISELABZ_TEST_POSTGRES_DSN like the other
// postgres tests.
func TestMigrationSchemaParity(t *testing.T) {
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set; skipping schema parity test")
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	pg := newPostgresTestStore(t, dsn, logger).db.(pgPlaceholderDB).DB

	sqliteDB, err := sql.Open("sqlite", "file:"+t.TempDir()+"/parity.db?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = sqliteDB.Close() })
	if err := RunMigrations(sqliteDB, "sqlite", logger); err != nil {
		t.Fatalf("sqlite migrations: %v", err)
	}

	assertShareLinkRetentionIndexes(t, pg, "postgres", true)
	assertShareLinkRetentionIndexes(t, sqliteDB, "sqlite", true)

	got, want := postgresSchemaColumns(t, pg), sqliteSchemaColumns(t, sqliteDB)
	if !reflect.DeepEqual(got, want) {
		gs, ws := map[string]bool{}, map[string]bool{}
		for _, c := range got {
			gs[c] = true
		}
		for _, c := range want {
			ws[c] = true
		}
		for _, c := range want {
			if !gs[c] {
				t.Errorf("missing in postgres: %s", c)
			}
		}
		for _, c := range got {
			if !ws[c] {
				t.Errorf("missing in sqlite: %s", c)
			}
		}
	}
}
