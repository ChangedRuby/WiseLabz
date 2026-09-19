package store

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/source/file" // required for file:// migration sources
)

//go:embed migrations/sqlite/*.sql
var sqliteMigrations embed.FS

//go:embed migrations/postgres/*.sql
var postgresMigrations embed.FS

// RunMigrations runs all pending database migrations for the given driver.
// Supports "sqlite" and "postgres".
func RunMigrations(db *sql.DB, driver string, logger *slog.Logger) error {
	logger.Info("running database migrations")

	m, err := newMigrator(db, driver)
	if err != nil {
		return err
	}

	if driver == "sqlite" {
		return runSQLiteWithForeignKeysOff(db, func() error {
			if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("run migrations: %w", err)
			}
			return nil
		})
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

// runSQLiteWithForeignKeysOff runs fn with foreign key enforcement disabled,
// then verifies no dangling references were introduced and re-enables it.
//
// Table-rebuild migrations DROP a parent table, which SQLite treats as an
// implicit DELETE that fires ON DELETE CASCADE / SET NULL on child rows.
// PRAGMA foreign_keys is a no-op inside a transaction and golang-migrate wraps
// each migration in one, so it must be toggled here, before the migration
// transaction begins. OpenDB pins the SQLite pool to a single
// connection, so the pragma applies to the connection the migrations run on.
func runSQLiteWithForeignKeysOff(db *sql.DB, fn func() error) (err error) {
	var enabled int
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&enabled); err != nil {
		return fmt.Errorf("read foreign_keys pragma: %w", err)
	}
	if enabled == 0 {
		return fn()
	}
	if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("disable foreign keys: %w", err)
	}
	defer func() {
		if _, rerr := db.Exec("PRAGMA foreign_keys = ON"); rerr != nil && err == nil {
			err = fmt.Errorf("re-enable foreign keys: %w", rerr)
		}
	}()

	if err := fn(); err != nil {
		return err
	}

	rows, err := db.Query("PRAGMA foreign_key_check")
	if err != nil {
		return fmt.Errorf("foreign key check: %w", err)
	}
	defer rows.Close() //nolint:errcheck
	if rows.Next() {
		return errors.New("foreign key check failed after migrations: dangling references found")
	}
	return rows.Err()
}

// RunMigrationsDown rolls back the most recently applied migration for the
// given driver. Supports "sqlite" and "postgres".
func RunMigrationsDown(db *sql.DB, driver string, logger *slog.Logger) error {
	logger.Info("rolling back database migration")

	m, err := newMigrator(db, driver)
	if err != nil {
		return err
	}

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run down migration: %w", err)
	}

	return nil
}

// newMigrator builds a golang-migrate instance for the given driver, reusing
// the caller-owned *sql.DB connection.
//
// golang-migrate's Close() calls database.Close() on the underlying driver
// instance, which is owned by the caller. Never call Close() on the returned
// migrator — doing so breaks idempotent migration re-runs.
func newMigrator(db *sql.DB, driver string) (*migrate.Migrate, error) {
	switch driver {
	case "sqlite":
		sub, err := fs.Sub(sqliteMigrations, "migrations/sqlite")
		if err != nil {
			return nil, fmt.Errorf("read sqlite migrations: %w", err)
		}

		src, err := iofs.New(sub, ".")
		if err != nil {
			return nil, fmt.Errorf("create migration source: %w", err)
		}

		dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return nil, fmt.Errorf("create migration driver: %w", err)
		}

		m, err := migrate.NewWithInstance("iofs", src, "sqlite3", dbDriver)
		if err != nil {
			return nil, fmt.Errorf("create migrator: %w", err)
		}
		return m, nil
	case "postgres":
		sub, err := fs.Sub(postgresMigrations, "migrations/postgres")
		if err != nil {
			return nil, fmt.Errorf("read postgres migrations: %w", err)
		}

		src, err := iofs.New(sub, ".")
		if err != nil {
			return nil, fmt.Errorf("create migration source: %w", err)
		}

		dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("create migration driver: %w", err)
		}

		m, err := migrate.NewWithInstance("iofs", src, "postgres", dbDriver)
		if err != nil {
			return nil, fmt.Errorf("create migrator: %w", err)
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}
