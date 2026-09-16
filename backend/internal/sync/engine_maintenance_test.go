package sync

import (
	"context"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// driftingSnapshot removes a section vs the seeded previous snapshot, which
// Compare treats as a "warning" (non-info) change — the kind maintenance
// windows are meant to suppress.
func driftingSnapshot() *connector.ServiceSnapshot {
	return &connector.ServiceSnapshot{
		ServiceName: "svc",
		Sections:    []connector.SnapshotSection{},
		FetchedAt:   time.Now(),
	}
}

func setupMaintenanceTestConnector(t *testing.T, s *store.Store, connType string) *store.ConnectorRecord {
	t.Helper()
	ctx := context.Background()
	connector.Register(
		connector.TypeSchema{Type: connType, Category: "test", Name: "Fake"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: driftingSnapshot()}, nil
		},
	)

	conn := &store.ConnectorRecord{Name: "svc", Category: "networking", Type: connType, Enabled: true}
	if err := s.CreateConnector(ctx, conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	prevSnapshot := `{"serviceName":"svc","sections":[{"title":"Ports","content":"22,80"}]}`
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{
		ConnectorID: conn.ID, Data: prevSnapshot, FetchedAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}
	return conn
}

// TestRunSyncSuppressesChangesDuringMaintenanceWindow verifies an active
// maintenance window suppresses change/alert creation on both scheduled and
// manual sync (RunSync is the shared path both use), while the new snapshot
// is still saved so sync history isn't lost.
func TestRunSyncSuppressesChangesDuringMaintenanceWindow(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	conn := setupMaintenanceTestConnector(t, s, "sync_test_maint_active")

	now := time.Now().UTC()
	if err := s.CreateMaintenanceWindow(ctx, &store.MaintenanceWindowRecord{
		ConnectorID: conn.ID, StartsAt: now.Format(time.RFC3339),
		EndsAt: now.Add(time.Hour).Format(time.RFC3339), CreatedBy: "user-1",
	}); err != nil {
		t.Fatalf("CreateMaintenanceWindow: %v", err)
	}

	countBefore, err := s.CountSnapshotsByConnector(ctx, conn.ID)
	if err != nil {
		t.Fatalf("CountSnapshotsByConnector: %v", err)
	}

	notifier := &fakeNotifier{}
	engine := NewEngine(s, nil, notifier, nil, "")

	result, err := engine.RunSync(ctx, conn.ID, "job1")
	if err != nil {
		t.Fatalf("RunSync: %v", err)
	}
	if result.ChangesCount != 0 {
		t.Fatalf("ChangesCount = %d, want 0 during maintenance window", result.ChangesCount)
	}
	if result.AlertsCount != 0 {
		t.Fatalf("AlertsCount = %d, want 0 during maintenance window", result.AlertsCount)
	}
	if got := notifier.count(); got != 0 {
		t.Fatalf("notifier called %d times, want 0 during maintenance window", got)
	}

	countAfter, err := s.CountSnapshotsByConnector(ctx, conn.ID)
	if err != nil {
		t.Fatalf("CountSnapshotsByConnector: %v", err)
	}
	if countAfter != countBefore+1 {
		t.Fatalf("snapshot count = %d, want %d (snapshot still saved during maintenance window)", countAfter, countBefore+1)
	}
}

// TestRunSyncNoMaintenanceWindowBehavesNormally is the control: with no
// window at all, drift still produces a change and alert.
func TestRunSyncNoMaintenanceWindowBehavesNormally(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	conn := setupMaintenanceTestConnector(t, s, "sync_test_maint_none")

	engine := NewEngine(s, nil, &fakeNotifier{}, nil, "")
	result, err := engine.RunSync(ctx, conn.ID, "job1")
	if err != nil {
		t.Fatalf("RunSync: %v", err)
	}
	if result.ChangesCount != 1 {
		t.Fatalf("ChangesCount = %d, want 1 with no maintenance window", result.ChangesCount)
	}
	if result.AlertsCount != 1 {
		t.Fatalf("AlertsCount = %d, want 1 with no maintenance window", result.AlertsCount)
	}
}

// TestRunSyncExpiredMaintenanceWindowBehavesNormally verifies a window that
// has already ended no longer suppresses drift output.
func TestRunSyncExpiredMaintenanceWindowBehavesNormally(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	conn := setupMaintenanceTestConnector(t, s, "sync_test_maint_expired")

	past := time.Now().UTC().Add(-2 * time.Hour)
	if err := s.CreateMaintenanceWindow(ctx, &store.MaintenanceWindowRecord{
		ConnectorID: conn.ID, StartsAt: past.Format(time.RFC3339),
		EndsAt: past.Add(time.Hour).Format(time.RFC3339), CreatedBy: "user-1", // ended an hour ago
	}); err != nil {
		t.Fatalf("CreateMaintenanceWindow: %v", err)
	}

	engine := NewEngine(s, nil, &fakeNotifier{}, nil, "")
	result, err := engine.RunSync(ctx, conn.ID, "job1")
	if err != nil {
		t.Fatalf("RunSync: %v", err)
	}
	if result.ChangesCount != 1 {
		t.Fatalf("ChangesCount = %d, want 1 with an expired maintenance window", result.ChangesCount)
	}
	if result.AlertsCount != 1 {
		t.Fatalf("AlertsCount = %d, want 1 with an expired maintenance window", result.AlertsCount)
	}
}
