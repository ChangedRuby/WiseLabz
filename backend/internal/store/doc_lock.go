package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/WiseLabz/wiselabz/internal/ws"
)

// Doc lock tuning: TTL is how long an acquired lock stays valid without
// renewal; heartbeat is how often the editor UI is expected to re-acquire
// (renew) it. Locks are advisory-only — see AcquireDocLock.
const (
	DocLockTTL       = 5 * time.Minute
	DocLockHeartbeat = 30 * time.Second
)

// ErrLockHeldByOther is returned by AcquireDocLock when a live lock is held
// by a different user.
var ErrLockHeldByOther = errors.New("doc lock held by another user")

// DocLockRecord represents a row in the doc_locks table: an advisory,
// single-writer editing lock on a doc.
type DocLockRecord struct {
	DocID      string `json:"docId"`
	UserID     string `json:"userId"`
	AcquiredAt string `json:"acquiredAt"`
	ExpiresAt  string `json:"expiresAt"`
}

// --- Doc locks (advisory presence/soft-lock, issue #92) ---

// AcquireDocLock takes (or renews) the advisory lock on a doc for userID.
// Succeeds when the doc is unlocked, already expired, or already held by
// userID (renewal); otherwise returns the live lock plus ErrLockHeldByOther.
// This never blocks Save — it's a UI nudge only.
func (s *Store) AcquireDocLock(ctx context.Context, docID, userID string) (*DocLockRecord, error) {
	now := time.Now().UTC()
	acquiredAt := now.Format(time.RFC3339)
	expiresAt := now.Add(DocLockTTL).Format(time.RFC3339)

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO doc_locks (doc_id, user_id, acquired_at, expires_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(doc_id) DO UPDATE SET
			user_id = excluded.user_id,
			acquired_at = excluded.acquired_at,
			expires_at = excluded.expires_at
		WHERE doc_locks.expires_at < ? OR doc_locks.user_id = ?
	`, docID, userID, acquiredAt, expiresAt, acquiredAt, userID)
	if err != nil {
		return nil, fmt.Errorf("acquire doc lock: %w", err)
	}

	lock, err := s.GetDocLock(ctx, docID)
	if err != nil {
		return nil, fmt.Errorf("acquire doc lock: %w", err)
	}
	if lock.UserID != userID {
		return lock, ErrLockHeldByOther
	}
	return lock, nil
}

// ReleaseDocLock releases the lock on docID if held by userID; a no-op
// otherwise (not held, held by someone else, or already expired).
func (s *Store) ReleaseDocLock(ctx context.Context, docID, userID string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM doc_locks WHERE doc_id = ? AND user_id = ?`, docID, userID); err != nil {
		return fmt.Errorf("release doc lock: %w", err)
	}
	return nil
}

// GetDocLock returns the live (non-expired) lock on docID, or ErrNotFound if
// unlocked/expired. Expiry is checked here (lazy), not dependent on a sweep.
func (s *Store) GetDocLock(ctx context.Context, docID string) (*DocLockRecord, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	l := &DocLockRecord{}
	err := s.db.QueryRowContext(ctx, `
		SELECT doc_id, user_id, acquired_at, expires_at
		FROM doc_locks WHERE doc_id = ? AND expires_at > ?
	`, docID, now).Scan(&l.DocID, &l.UserID, &l.AcquiredAt, &l.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get doc lock: %w", err)
	}
	return l, nil
}

// RunDocLockSweep periodically clears expired doc locks and broadcasts their
// expiry so idle viewers' presence banners clear without polling. Expiry
// itself is already enforced lazily by GetDocLock; this sweep is purely row
// hygiene plus the WS notification. Mirrors quality.RunStaleSweep.
func RunDocLockSweep(ctx context.Context, s *Store, hub *ws.Hub, interval time.Duration, logger *slog.Logger) {
	logger.Info("Doc lock sweep started")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		runDocLockSweep(ctx, s, hub, logger)
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
		}
	}
}

func runDocLockSweep(ctx context.Context, s *Store, hub *ws.Hub, logger *slog.Logger) {
	now := time.Now().UTC().Format(time.RFC3339)
	rows, err := s.db.QueryContext(ctx, `SELECT doc_id, user_id FROM doc_locks WHERE expires_at <= ?`, now)
	if err != nil {
		logger.Error("doc lock sweep: list expired", "error", err)
		return
	}
	var expired []DocLockRecord
	for rows.Next() {
		var l DocLockRecord
		if err := rows.Scan(&l.DocID, &l.UserID); err != nil {
			rows.Close() //nolint:errcheck
			logger.Error("doc lock sweep: scan", "error", err)
			return
		}
		expired = append(expired, l)
	}
	rows.Close() //nolint:errcheck
	if err := rows.Err(); err != nil {
		logger.Error("doc lock sweep: iterate expired", "error", err)
		return
	}

	if len(expired) == 0 {
		return
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM doc_locks WHERE expires_at <= ?`, now); err != nil {
		logger.Error("doc lock sweep: delete", "error", err)
		return
	}
	for _, l := range expired {
		if hub != nil {
			hub.Broadcast(ws.EventDocLockExpired, map[string]any{"docId": l.DocID, "userId": l.UserID})
		}
	}
}
