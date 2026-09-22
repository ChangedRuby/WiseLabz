package store

import (
	"context"
	"errors"
	"testing"
)

func mustCreateGoldenSnapshotConnector(t *testing.T, s *Store) string {
	t.Helper()
	ctx := context.Background()
	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com", Enabled: true}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c.ID
}

func TestPinGoldenSnapshotRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateGoldenSnapshotConnector(t, s)

	if _, err := s.GetGoldenSnapshot(ctx, connID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetGoldenSnapshot() before pin error = %v, want ErrNotFound", err)
	}

	snap := &SnapshotRecord{ConnectorID: connID, Data: `{"serviceName":"svc"}`}
	if err := s.CreateSnapshot(ctx, snap); err != nil {
		t.Fatalf("CreateSnapshot() error: %v", err)
	}

	if err := s.PinGoldenSnapshot(ctx, &GoldenSnapshotRecord{ConnectorID: connID, SnapshotID: snap.ID, PinnedBy: "user-1"}); err != nil {
		t.Fatalf("PinGoldenSnapshot() error: %v", err)
	}

	got, err := s.GetGoldenSnapshot(ctx, connID)
	if err != nil {
		t.Fatalf("GetGoldenSnapshot() error: %v", err)
	}
	if got.SnapshotID != snap.ID || got.PinnedBy != "user-1" || got.PinnedAt == "" {
		t.Fatalf("GetGoldenSnapshot() = %+v, want snapshot %s pinned by user-1", got, snap.ID)
	}

	// Re-pinning a different snapshot replaces the previous pin rather than
	// erroring, so a connector always has at most one golden snapshot.
	snap2 := &SnapshotRecord{ConnectorID: connID, Data: `{"serviceName":"svc2"}`}
	if err := s.CreateSnapshot(ctx, snap2); err != nil {
		t.Fatalf("CreateSnapshot() error: %v", err)
	}
	if err := s.PinGoldenSnapshot(ctx, &GoldenSnapshotRecord{ConnectorID: connID, SnapshotID: snap2.ID, PinnedBy: "user-2"}); err != nil {
		t.Fatalf("PinGoldenSnapshot() replace error: %v", err)
	}
	got, err = s.GetGoldenSnapshot(ctx, connID)
	if err != nil {
		t.Fatalf("GetGoldenSnapshot() after replace error: %v", err)
	}
	if got.SnapshotID != snap2.ID || got.PinnedBy != "user-2" {
		t.Fatalf("GetGoldenSnapshot() after replace = %+v, want snapshot %s pinned by user-2", got, snap2.ID)
	}

	if err := s.UnpinGoldenSnapshot(ctx, connID); err != nil {
		t.Fatalf("UnpinGoldenSnapshot() error: %v", err)
	}
	if _, err := s.GetGoldenSnapshot(ctx, connID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetGoldenSnapshot() after unpin error = %v, want ErrNotFound", err)
	}

	// Unpinning again is a no-op, not an error.
	if err := s.UnpinGoldenSnapshot(ctx, connID); err != nil {
		t.Fatalf("UnpinGoldenSnapshot() second call error: %v", err)
	}
}

func TestGetSnapshotByID(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateGoldenSnapshotConnector(t, s)

	snap := &SnapshotRecord{ConnectorID: connID, Data: `{"serviceName":"svc"}`}
	if err := s.CreateSnapshot(ctx, snap); err != nil {
		t.Fatalf("CreateSnapshot() error: %v", err)
	}

	got, err := s.GetSnapshotByID(ctx, snap.ID)
	if err != nil {
		t.Fatalf("GetSnapshotByID() error: %v", err)
	}
	if got.ID != snap.ID || got.ConnectorID != connID {
		t.Fatalf("GetSnapshotByID() = %+v, want %+v", got, snap)
	}

	if _, err := s.GetSnapshotByID(ctx, "does-not-exist"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetSnapshotByID() unknown id error = %v, want ErrNotFound", err)
	}
}
