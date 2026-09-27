package store

import (
	"context"
	"testing"
	"time"
)

func createTestConnector(ctx context.Context, t *testing.T, s *Store) string {
	t.Helper()
	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c.ID
}

func TestSyncRunCreateAndList(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)
	if s.driver == "sqlite" {
		if _, err := s.db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
			t.Fatal(err)
		}
	}
	snapshot := &SnapshotRecord{ConnectorID: connectorID, Data: "{}"}
	if err := s.CreateSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}

	durationMs := 1234
	if err := s.CreateSyncRun(ctx, &SyncRunRecord{
		ConnectorID:  connectorID,
		FinishedAt:   time.Now().UTC().Format(time.RFC3339),
		DurationMs:   &durationMs,
		Status:       SyncRunStatusSuccess,
		Attempt:      1,
		ChangesCount: 3,
		AlertsCount:  1,
		SnapshotID:   &snapshot.ID,
	}); err != nil {
		t.Fatalf("CreateSyncRun() error: %v", err)
	}

	runs, err := s.ListSyncRunsByConnector(ctx, connectorID, 20)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if len(runs) != 1 {
		t.Fatalf("len(runs) = %d, want 1", len(runs))
	}
	r := runs[0]
	if r.ConnectorID != connectorID || r.Status != SyncRunStatusSuccess || r.ChangesCount != 3 || r.AlertsCount != 1 {
		t.Fatalf("runs[0] = %+v, want connectorID=%s status=success changes=3 alerts=1", r, connectorID)
	}
	if r.DurationMs == nil || *r.DurationMs != durationMs {
		t.Fatalf("runs[0].DurationMs = %v, want %d", r.DurationMs, durationMs)
	}
	if r.Attempt != 1 {
		t.Fatalf("runs[0].Attempt = %d, want 1", r.Attempt)
	}
	if r.SnapshotID == nil || *r.SnapshotID != snapshot.ID {
		t.Fatalf("runs[0].SnapshotID = %v, want %s", r.SnapshotID, snapshot.ID)
	}
	// Deleting a snapshot clears the optional FK without removing the run.
	if _, err := s.db.ExecContext(ctx, `DELETE FROM service_snapshots WHERE id = ?`, snapshot.ID); err != nil {
		t.Fatal(err)
	}
	runs, err = s.ListSyncRunsByConnector(ctx, connectorID, 20)
	if err != nil || len(runs) != 1 || runs[0].SnapshotID != nil {
		t.Fatalf("runs after snapshot deletion = %+v, err %v", runs, err)
	}
}

func TestSyncRunDefaultsIDAndAttempt(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	r := &SyncRunRecord{ConnectorID: connectorID, Status: SyncRunStatusError, Error: "boom"}
	if err := s.CreateSyncRun(ctx, r); err != nil {
		t.Fatalf("CreateSyncRun() error: %v", err)
	}
	if r.ID == "" {
		t.Fatal("expected CreateSyncRun to assign an ID")
	}
	if r.Attempt != 1 {
		t.Fatalf("Attempt = %d, want default 1", r.Attempt)
	}

	runs, err := s.ListSyncRunsByConnector(ctx, connectorID, 20)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if len(runs) != 1 || runs[0].Error != "boom" || runs[0].Status != SyncRunStatusError {
		t.Fatalf("runs = %+v, want 1 error run with error=boom", runs)
	}
	if runs[0].DurationMs != nil {
		t.Fatalf("DurationMs = %v, want nil (unset)", runs[0].DurationMs)
	}
}

func TestSyncRunListOrderedNewestFirstAndLimited(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	base := time.Now().UTC()
	for i, offset := range []time.Duration{-2 * time.Hour, -time.Hour, 0} {
		if err := s.CreateSyncRun(ctx, &SyncRunRecord{
			ConnectorID: connectorID,
			StartedAt:   base.Add(offset).Format(time.RFC3339),
			Status:      SyncRunStatusSuccess,
			Attempt:     i + 1,
		}); err != nil {
			t.Fatalf("CreateSyncRun(%d) error: %v", i, err)
		}
	}

	runs, err := s.ListSyncRunsByConnector(ctx, connectorID, 2)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if len(runs) != 2 {
		t.Fatalf("len(runs) = %d, want 2 (limit)", len(runs))
	}
	if runs[0].Attempt != 3 || runs[1].Attempt != 2 {
		t.Fatalf("runs attempts = [%d,%d], want [3,2] (newest first)", runs[0].Attempt, runs[1].Attempt)
	}
}

func TestSyncRunListByConnectorNeverReturnsNil(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	runs, err := s.ListSyncRunsByConnector(ctx, connectorID, 20)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if runs == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(runs) != 0 {
		t.Fatalf("len(runs) = %d, want 0", len(runs))
	}
}
