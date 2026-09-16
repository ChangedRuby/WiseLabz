package store

import (
	"context"
	"testing"
	"time"
)

func mustCreateMaintenanceConnector(t *testing.T, s *Store) string {
	t.Helper()
	ctx := context.Background()
	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com", Enabled: true}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c.ID
}

func TestCreateAndGetActiveMaintenanceWindow(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateMaintenanceConnector(t, s)

	got, err := s.GetActiveMaintenanceWindow(ctx, connID)
	if err != nil {
		t.Fatalf("GetActiveMaintenanceWindow() error: %v", err)
	}
	if got != nil {
		t.Fatalf("GetActiveMaintenanceWindow() = %+v, want nil (none created yet)", got)
	}

	now := time.Now().UTC()
	m := &MaintenanceWindowRecord{
		ConnectorID: connID,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(time.Hour).Format(time.RFC3339),
		CreatedBy:   "user-1",
	}
	if err := s.CreateMaintenanceWindow(ctx, m); err != nil {
		t.Fatalf("CreateMaintenanceWindow() error: %v", err)
	}
	if m.ID == "" {
		t.Fatal("CreateMaintenanceWindow() left ID empty")
	}

	got, err = s.GetActiveMaintenanceWindow(ctx, connID)
	if err != nil {
		t.Fatalf("GetActiveMaintenanceWindow() error: %v", err)
	}
	if got == nil || got.ID != m.ID {
		t.Fatalf("GetActiveMaintenanceWindow() = %+v, want %+v", got, m)
	}
	if got.CreatedBy != "user-1" {
		t.Fatalf("CreatedBy = %q, want user-1", got.CreatedBy)
	}
}

func TestGetActiveMaintenanceWindowExpired(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateMaintenanceConnector(t, s)

	now := time.Now().UTC()
	expired := &MaintenanceWindowRecord{
		ConnectorID: connID,
		StartsAt:    now.Add(-2 * time.Hour).Format(time.RFC3339),
		EndsAt:      now.Add(-time.Minute).Format(time.RFC3339), // just expired
		CreatedBy:   "user-1",
	}
	if err := s.CreateMaintenanceWindow(ctx, expired); err != nil {
		t.Fatalf("CreateMaintenanceWindow() error: %v", err)
	}

	got, err := s.GetActiveMaintenanceWindow(ctx, connID)
	if err != nil {
		t.Fatalf("GetActiveMaintenanceWindow() error: %v", err)
	}
	if got != nil {
		t.Fatalf("GetActiveMaintenanceWindow() = %+v, want nil (window expired)", got)
	}
}

func TestGetActiveMaintenanceWindowOverlapping(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateMaintenanceConnector(t, s)

	now := time.Now().UTC()
	sooner := &MaintenanceWindowRecord{
		ConnectorID: connID,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(30 * time.Minute).Format(time.RFC3339),
		CreatedBy:   "user-1",
	}
	later := &MaintenanceWindowRecord{
		ConnectorID: connID,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(2 * time.Hour).Format(time.RFC3339),
		CreatedBy:   "user-2",
	}
	if err := s.CreateMaintenanceWindow(ctx, sooner); err != nil {
		t.Fatalf("CreateMaintenanceWindow(sooner) error: %v", err)
	}
	if err := s.CreateMaintenanceWindow(ctx, later); err != nil {
		t.Fatalf("CreateMaintenanceWindow(later) error: %v", err)
	}

	got, err := s.GetActiveMaintenanceWindow(ctx, connID)
	if err != nil {
		t.Fatalf("GetActiveMaintenanceWindow() error: %v", err)
	}
	if got == nil || got.ID != later.ID {
		t.Fatalf("GetActiveMaintenanceWindow() = %+v, want the later-expiring window %+v", got, later)
	}
}

func TestCloseMaintenanceWindow(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateMaintenanceConnector(t, s)

	now := time.Now().UTC()
	m := &MaintenanceWindowRecord{
		ConnectorID: connID,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(time.Hour).Format(time.RFC3339),
		CreatedBy:   "user-1",
	}
	if err := s.CreateMaintenanceWindow(ctx, m); err != nil {
		t.Fatalf("CreateMaintenanceWindow() error: %v", err)
	}

	if err := s.CloseMaintenanceWindow(ctx, m.ID); err != nil {
		t.Fatalf("CloseMaintenanceWindow() error: %v", err)
	}

	got, err := s.GetActiveMaintenanceWindow(ctx, connID)
	if err != nil {
		t.Fatalf("GetActiveMaintenanceWindow() error: %v", err)
	}
	if got != nil {
		t.Fatalf("GetActiveMaintenanceWindow() = %+v, want nil after close", got)
	}
}

func TestCloseMaintenanceWindowAlreadyClosed(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connID := mustCreateMaintenanceConnector(t, s)

	now := time.Now().UTC()
	m := &MaintenanceWindowRecord{
		ConnectorID: connID,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(time.Hour).Format(time.RFC3339),
		CreatedBy:   "user-1",
	}
	if err := s.CreateMaintenanceWindow(ctx, m); err != nil {
		t.Fatalf("CreateMaintenanceWindow() error: %v", err)
	}
	if err := s.CloseMaintenanceWindow(ctx, m.ID); err != nil {
		t.Fatalf("CloseMaintenanceWindow() (first) error: %v", err)
	}

	// Closing an already-closed window is a no-op success, not an error —
	// the row still exists, so RowsAffected is 1 even though ends_at was
	// already in the past.
	if err := s.CloseMaintenanceWindow(ctx, m.ID); err != nil {
		t.Fatalf("CloseMaintenanceWindow() (second) error: %v", err)
	}
}

func TestCloseMaintenanceWindowNotFound(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	if err := s.CloseMaintenanceWindow(ctx, "does-not-exist"); err != ErrNotFound {
		t.Fatalf("CloseMaintenanceWindow() error = %v, want ErrNotFound", err)
	}
}

func TestListActiveMaintenanceWindows(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connA := mustCreateMaintenanceConnector(t, s)
	connB := mustCreateMaintenanceConnector(t, s)

	now := time.Now().UTC()
	active := &MaintenanceWindowRecord{
		ConnectorID: connA, StartsAt: now.Format(time.RFC3339),
		EndsAt: now.Add(time.Hour).Format(time.RFC3339), CreatedBy: "user-1",
	}
	expired := &MaintenanceWindowRecord{
		ConnectorID: connB, StartsAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
		EndsAt: now.Add(-time.Hour).Format(time.RFC3339), CreatedBy: "user-1",
	}
	if err := s.CreateMaintenanceWindow(ctx, active); err != nil {
		t.Fatalf("CreateMaintenanceWindow(active) error: %v", err)
	}
	if err := s.CreateMaintenanceWindow(ctx, expired); err != nil {
		t.Fatalf("CreateMaintenanceWindow(expired) error: %v", err)
	}

	windows, err := s.ListActiveMaintenanceWindows(ctx)
	if err != nil {
		t.Fatalf("ListActiveMaintenanceWindows() error: %v", err)
	}
	if len(windows) != 1 || windows[0].ID != active.ID {
		t.Fatalf("ListActiveMaintenanceWindows() = %+v, want exactly [active]", windows)
	}
}

func TestListActiveMaintenanceWindowsEmpty(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	windows, err := s.ListActiveMaintenanceWindows(ctx)
	if err != nil {
		t.Fatalf("ListActiveMaintenanceWindows() error: %v", err)
	}
	if windows == nil {
		t.Fatal("ListActiveMaintenanceWindows() returned nil slice, want empty non-nil")
	}
	if len(windows) != 0 {
		t.Fatalf("ListActiveMaintenanceWindows() = %+v, want empty", windows)
	}
}
