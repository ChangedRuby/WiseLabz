// Package store provides the central data store for WiseLabz.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
)

// DBTX is the subset of *sql.DB used by the store package and by API
// handlers that issue raw SQL via Store.DB(). It exists so that a
// PostgreSQL connection can be wrapped with placeholder rewriting (see
// pgdb.go) without changing any call site.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
	PingContext(ctx context.Context) error
	Close() error
}

// Store is the central data access layer, holding the database connection
// and all repository implementations.
type Store struct {
	db                   DBTX
	userStatusMu         sync.Mutex
	userStatusGeneration uint64
	userStatuses         map[string]userStatus
}

type userStatus struct {
	role     string
	disabled bool
	expires  time.Time
}

// invalidateUserStatuses also prevents an in-flight read from repopulating
// the cache with a value read before a mutation committed.
func (s *Store) invalidateUserStatuses() {
	s.userStatusMu.Lock()
	defer s.userStatusMu.Unlock()
	s.userStatusGeneration++
	s.userStatuses = nil
}

type transactionDB struct {
	*sql.Tx
}

func (t transactionDB) PingContext(ctx context.Context) error {
	var one int
	return t.QueryRowContext(ctx, "SELECT 1").Scan(&one)
}

func (t transactionDB) Close() error {
	return t.Rollback()
}

type pgTransactionDB struct {
	transactionDB
}

func (p pgTransactionDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.Tx.ExecContext(ctx, rewritePlaceholders(query), args...)
}

func (p pgTransactionDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return p.Tx.QueryContext(ctx, rewritePlaceholders(query), args...)
}

func (p pgTransactionDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return p.Tx.QueryRowContext(ctx, rewritePlaceholders(query), args...)
}

func (p pgTransactionDB) QueryRow(query string, args ...any) *sql.Row {
	return p.Tx.QueryRow(rewritePlaceholders(query), args...)
}

// New creates a new Store with the given database connection. driver
// selects the placeholder dialect: "postgres" wraps db so that `?`
// placeholders are rewritten to `$1, $2, ...`; any other value (including
// "sqlite") uses db directly.
func New(db *sql.DB, driver string) *Store {
	if driver == "postgres" {
		return &Store{db: pgPlaceholderDB{db}}
	}
	return &Store{db: db}
}

// WithinTransaction runs fn against a transaction-bound Store and commits only on success.
func (s *Store) WithinTransaction(ctx context.Context, fn func(*Store) error) error {
	var (
		tx       *sql.Tx
		storeDB  DBTX
		err      error
		postgres bool
	)
	switch db := s.db.(type) {
	case *sql.DB:
		tx, err = db.BeginTx(ctx, nil)
	case pgPlaceholderDB:
		postgres = true
		tx, err = db.BeginTx(ctx, nil)
	default:
		return fmt.Errorf("begin transaction: unsupported database wrapper %T", s.db)
	}
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	storeDB = transactionDB{Tx: tx}
	if postgres {
		storeDB = pgTransactionDB{transactionDB{Tx: tx}}
	}
	if err := fn(&Store{db: storeDB}); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	s.invalidateUserStatuses()
	return nil
}

// DB returns the underlying database connection for direct queries.
func (s *Store) DB() DBTX {
	return s.db
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// Ping checks the database connection health.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// initSingletons ensures singleton config rows exist (auth_config, ai_config, notification_config).
// Called after migrations run.
func (s *Store) initSingletons(ctx context.Context) error {
	rows := []string{
		`INSERT INTO auth_config (id, local_enabled, access_token_ttl, refresh_token_ttl, step_up_for_destructive)
		 VALUES (1, 1, 900, 604800, 1) ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO ai_config (id, enabled, provider, model, api_key_encrypted, base_url, mode, embed_provider, embed_model, embed_api_key_encrypted, embed_base_url)
		 VALUES (1, 0, NULL, NULL, '', NULL, 'suggest_only', NULL, NULL, '', NULL) ON CONFLICT (id) DO NOTHING`,
		`INSERT INTO notification_config (id, config_json)
		 VALUES (1, '{}') ON CONFLICT (id) DO NOTHING`,
	}

	for _, q := range rows {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("init singleton: %w", err)
		}
	}
	return nil
}

// Init bootstraps the store after migrations: seeds singletons and creates the admin user if needed.
func (s *Store) Init(ctx context.Context, adminPassword string) error {
	if err := s.initSingletons(ctx); err != nil {
		return err
	}

	// Check if any users exist
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return fmt.Errorf("count users: %w", err)
	}

	if count == 0 {
		if adminPassword == "" {
			return fmt.Errorf("no users exist and WISELABZ_ADMIN_PASSWORD is not set: refusing to start without an admin bootstrap password")
		}

		hash, err := auth.HashPassword(adminPassword)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}

		admin := &User{
			Username:                   "admin",
			DisplayName:                "Administrator",
			InstanceAdminRole:          "admin",
			AuthSource:                 "local",
			PasswordHash:               hash,
			CanManageDashboardDefaults: true,
		}
		if err := s.CreateUser(ctx, admin); err != nil {
			return fmt.Errorf("seed admin user: %w", err)
		}
	}

	return nil
}
