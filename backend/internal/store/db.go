package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	_ "modernc.org/sqlite"             // SQLite driver (pure Go, no CGO)
)

// PoolConfig holds Postgres connection pool settings. Zero or negative
// values fall back to the defaults. SQLite ignores it (single connection).
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (p PoolConfig) withDefaults() PoolConfig {
	if p.MaxOpenConns <= 0 {
		p.MaxOpenConns = 20
	}
	if p.MaxIdleConns <= 0 {
		p.MaxIdleConns = 5
	}
	if p.MaxIdleConns > p.MaxOpenConns {
		p.MaxIdleConns = p.MaxOpenConns
	}
	if p.ConnMaxLifetime <= 0 {
		p.ConnMaxLifetime = 30 * time.Minute
	}
	if p.ConnMaxIdleTime <= 0 {
		p.ConnMaxIdleTime = 5 * time.Minute
	}
	return p
}

// OpenDB opens a database connection based on the driver and DSN.
// Supports "sqlite" and "postgres". An optional PoolConfig tunes the
// Postgres pool; omitted, defaults apply.
func OpenDB(driver, dsn string, pool ...PoolConfig) (*sql.DB, error) {
	var pc PoolConfig
	if len(pool) > 0 {
		pc = pool[0]
	}
	pc = pc.withDefaults()
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

	// Configure the pool before Ping so the health check uses the final limits.
	if driver == "postgres" {
		// PostgreSQL supports concurrent writers; size the pool accordingly.
		db.SetMaxOpenConns(pc.MaxOpenConns)
		db.SetMaxIdleConns(pc.MaxIdleConns)
		db.SetConnMaxLifetime(pc.ConnMaxLifetime)
		db.SetConnMaxIdleTime(pc.ConnMaxIdleTime)
	} else {
		// SQLite is single-writer — keep the pool at 1 to avoid SQLITE_BUSY.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	}

	if err := db.Ping(); err != nil {
		// Best-effort close on ping failure; the connection pool is not yet healthy.
		if closeErr := db.Close(); closeErr != nil {
			err = fmt.Errorf("ping database: %w (close: %w)", err, closeErr)
		}
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
