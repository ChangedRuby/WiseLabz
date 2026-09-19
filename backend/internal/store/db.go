package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	_ "modernc.org/sqlite"             // SQLite driver (pure Go, no CGO)
)

// OpenDB opens a database connection based on the driver and DSN.
// Supports "sqlite" and "postgres".
func OpenDB(driver, dsn string) (*sql.DB, error) {
	if driver == "" {
		driver = "sqlite"
	}

	// "postgres" (our config-facing driver name) maps to the "pgx" sql.DB
	// driver registered above. golang-migrate's postgres support package
	// also transitively blank-imports lib/pq, which separately registers
	// itself under the name "postgres" — sql.Open("postgres", ...) would
	// silently use lib/pq instead of the pgx driver we depend on if we
	// didn't translate the name here.
	sqlDriver := driver
	if driver == "postgres" {
		sqlDriver = "pgx"
	} else {
		// WAL lets external readers (backup copy, sqlite3 CLI) coexist with
		// the writer; busy_timeout makes contention wait instead of failing
		// with SQLITE_BUSY.
		for _, p := range []struct{ name, value string }{
			{"foreign_keys", "1"},
			{"journal_mode", "WAL"},
			{"busy_timeout", "5000"},
			{"synchronous", "NORMAL"},
		} {
			if strings.Contains(dsn, "_pragma="+p.name) {
				continue
			}
			separator := "?"
			if strings.Contains(dsn, "?") {
				separator = "&"
			}
			dsn += separator + "_pragma=" + p.name + "(" + p.value + ")"
		}
	}

	db, err := sql.Open(sqlDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		// Best-effort close on ping failure; the connection pool is not yet healthy.
		if closeErr := db.Close(); closeErr != nil {
			err = fmt.Errorf("ping database: %w (close: %w)", err, closeErr)
		}
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if driver == "postgres" {
		// PostgreSQL supports concurrent writers; size the pool accordingly.
		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(5)
	} else {
		// SQLite is single-writer — keep the pool at 1 to avoid SQLITE_BUSY.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	}

	return db, nil
}
