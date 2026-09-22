package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// GoldenSnapshotRecord represents a row in the golden_snapshots table: the
// snapshot pinned as a connector's known-good configuration baseline, so
// later syncs can be checked for drift against it (#275).
type GoldenSnapshotRecord struct {
	ConnectorID string `json:"connectorId"`
	SnapshotID  string `json:"snapshotId"`
	PinnedBy    string `json:"pinnedBy"`
	PinnedAt    string `json:"pinnedAt"`
}

// PinGoldenSnapshot sets (or replaces) the connector's golden snapshot.
func (s *Store) PinGoldenSnapshot(ctx context.Context, g *GoldenSnapshotRecord) error {
	if g.PinnedAt == "" {
		g.PinnedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO golden_snapshots (connector_id, snapshot_id, pinned_by, pinned_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(connector_id) DO UPDATE SET
			snapshot_id = excluded.snapshot_id,
			pinned_by = excluded.pinned_by,
			pinned_at = excluded.pinned_at
	`, g.ConnectorID, g.SnapshotID, g.PinnedBy, g.PinnedAt)
	if err != nil {
		return fmt.Errorf("pin golden snapshot: %w", err)
	}
	return nil
}

// GetGoldenSnapshot returns the connector's pinned golden snapshot, or
// ErrNotFound if none is pinned.
func (s *Store) GetGoldenSnapshot(ctx context.Context, connectorID string) (*GoldenSnapshotRecord, error) {
	g := &GoldenSnapshotRecord{}
	err := s.db.QueryRowContext(ctx, `
		SELECT connector_id, snapshot_id, pinned_by, pinned_at
		FROM golden_snapshots WHERE connector_id = ?
	`, connectorID).Scan(&g.ConnectorID, &g.SnapshotID, &g.PinnedBy, &g.PinnedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get golden snapshot: %w", err)
	}
	return g, nil
}

// UnpinGoldenSnapshot removes the connector's golden snapshot, if any. A
// no-op (nil error) when none is pinned.
func (s *Store) UnpinGoldenSnapshot(ctx context.Context, connectorID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM golden_snapshots WHERE connector_id = ?`, connectorID)
	if err != nil {
		return fmt.Errorf("unpin golden snapshot: %w", err)
	}
	return nil
}

// GetSnapshotByID returns a single snapshot by ID, or ErrNotFound.
func (s *Store) GetSnapshotByID(ctx context.Context, id string) (*SnapshotRecord, error) {
	sn := &SnapshotRecord{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, connector_id, data, fetched_at
		FROM service_snapshots WHERE id = ?
	`, id).Scan(&sn.ID, &sn.ConnectorID, &sn.Data, &sn.FetchedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get snapshot by id: %w", err)
	}
	return sn, nil
}
