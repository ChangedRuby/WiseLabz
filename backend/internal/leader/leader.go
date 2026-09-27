// Package leader elects one PostgreSQL replica to run background jobs.
package leader

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

// lockKey is a stable application-wide advisory-lock namespace.
const lockKey int64 = 0x776973656c61627a // "wiselabz"

// Election is the lifecycle's leader-election contract.
type Election interface {
	Campaign(context.Context) error
	Watch(context.Context) <-chan error
	Close() error
}

// Elector holds the advisory lock on one pinned database session.
type Elector struct {
	db       *sql.DB
	interval time.Duration
	mu       sync.Mutex
	conn     *sql.Conn
}

func New(db *sql.DB, interval time.Duration) *Elector {
	return &Elector{db: db, interval: interval}
}

// Campaign polls until this session acquires the lock or ctx ends.
func (e *Elector) Campaign(ctx context.Context) error {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()
	for {
		conn, err := e.db.Conn(ctx)
		if err == nil {
			var held bool
			err = conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockKey).Scan(&held)
			if err == nil && held {
				e.mu.Lock()
				e.conn = conn
				e.mu.Unlock()
				return nil
			}
			_ = conn.Close()
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// Watch checks the pinned session. A failed ping means its lock is no longer safe.
func (e *Elector) Watch(ctx context.Context) <-chan error {
	lost := make(chan error, 1)
	go func() {
		defer close(lost)
		ticker := time.NewTicker(e.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
			e.mu.Lock()
			conn := e.conn
			e.mu.Unlock()
			if conn == nil {
				lost <- errors.New("leader session closed")
				return
			}
			pingCtx, cancel := context.WithTimeout(ctx, e.interval)
			err := conn.PingContext(pingCtx)
			cancel()
			if err != nil {
				if ctx.Err() == nil {
					lost <- fmt.Errorf("leader session lost: %w", err)
				}
				return
			}
		}
	}()
	return lost
}

// Close releases the lock before returning its session to the sql.DB pool.
func (e *Elector) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.conn == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var released bool
	err := e.conn.QueryRowContext(ctx, "SELECT pg_advisory_unlock($1)", lockKey).Scan(&released)
	if closeErr := e.conn.Close(); err == nil {
		err = closeErr
	}
	e.conn = nil
	return err
}

type Noop struct{}

func (Noop) Campaign(context.Context) error     { return nil }
func (Noop) Watch(context.Context) <-chan error { return nil }
func (Noop) Close() error                       { return nil }
