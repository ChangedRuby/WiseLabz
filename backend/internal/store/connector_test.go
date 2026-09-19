package store

import (
	"context"
	"testing"
	"time"
)

// TestParseConnectorConfigNullDataDoesNotPanic is a regression test for
// #196: json.Unmarshal("null", &cfg) succeeds but leaves cfg == nil, and
// callers (sync engine, connector health-check handler) write into the
// returned map directly, which panics with "assignment to entry in nil
// map". ParseConnectorConfig must hand back a non-nil, writable map even
// when config_data is the literal string "null".
func TestParseConnectorConfigNullDataDoesNotPanic(t *testing.T) {
	cfg, err := ParseConnectorConfig("unknown-type", "null", "")
	if err != nil {
		t.Fatalf("ParseConnectorConfig() error: %v", err)
	}
	if cfg == nil {
		t.Fatal("ParseConnectorConfig() returned a nil map for null config_data")
	}
	cfg["url"] = "https://example.com" // would panic on a nil map
	if cfg["url"] != "https://example.com" {
		t.Fatalf("cfg[\"url\"] = %v, want https://example.com", cfg["url"])
	}
}

func TestConnectorOwnerRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com", Owner: "Platform"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	got, err := s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() error: %v", err)
	}
	if got.Owner != "Platform" {
		t.Fatalf("Owner = %q, want Platform", got.Owner)
	}
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{"owner": "Operations"}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}
	got, err = s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() after update error: %v", err)
	}
	if got.Owner != "Operations" {
		t.Fatalf("Owner = %q, want Operations", got.Owner)
	}
}

func TestConnectorScheduleFieldsDefaultNull(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	got, err := s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() error: %v", err)
	}
	if got.ScheduleSeconds != nil {
		t.Fatalf("ScheduleSeconds = %v, want nil (manual only)", got.ScheduleSeconds)
	}
	if got.NextRunAt != "" {
		t.Fatalf("NextRunAt = %q, want empty", got.NextRunAt)
	}
	if got.LastSyncDurationMs != nil {
		t.Fatalf("LastSyncDurationMs = %v, want nil", got.LastSyncDurationMs)
	}
	if got.LastSyncError != "" {
		t.Fatalf("LastSyncError = %q, want empty", got.LastSyncError)
	}
	if got.RetryCount != 0 {
		t.Fatalf("RetryCount = %d, want 0", got.RetryCount)
	}
}

func TestConnectorScheduleFieldsRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	schedule := 1800
	durationMs := 42
	next := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{
		"schedule_seconds":      schedule,
		"next_run_at":           next,
		"last_sync_duration_ms": durationMs,
		"last_sync_error":       "boom",
		"retry_count":           2,
	}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}

	got, err := s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() error: %v", err)
	}
	if got.ScheduleSeconds == nil || *got.ScheduleSeconds != schedule {
		t.Fatalf("ScheduleSeconds = %v, want %d", got.ScheduleSeconds, schedule)
	}
	if got.NextRunAt != next {
		t.Fatalf("NextRunAt = %q, want %q", got.NextRunAt, next)
	}
	if got.LastSyncDurationMs == nil || *got.LastSyncDurationMs != durationMs {
		t.Fatalf("LastSyncDurationMs = %v, want %d", got.LastSyncDurationMs, durationMs)
	}
	if got.LastSyncError != "boom" {
		t.Fatalf("LastSyncError = %q, want boom", got.LastSyncError)
	}
	if got.RetryCount != 2 {
		t.Fatalf("RetryCount = %d, want 2", got.RetryCount)
	}

	// Clearing scheduleSeconds/next_run_at back to null (e.g. disabling the
	// schedule) must round-trip to nil, not a zero value.
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{
		"schedule_seconds": nil,
		"next_run_at":      nil,
	}); err != nil {
		t.Fatalf("UpdateConnector(clear) error: %v", err)
	}
	got, err = s.GetConnector(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetConnector() error: %v", err)
	}
	if got.ScheduleSeconds != nil {
		t.Fatalf("ScheduleSeconds after clear = %v, want nil", got.ScheduleSeconds)
	}
	if got.NextRunAt != "" {
		t.Fatalf("NextRunAt after clear = %q, want empty", got.NextRunAt)
	}
}

func TestListDueConnectors(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	now := time.Now().UTC()
	past := now.Add(-time.Hour).Format(time.RFC3339)
	future := now.Add(time.Hour).Format(time.RFC3339)
	nowStr := now.Format(time.RFC3339)

	mustCreate := func(name string, enabled bool) string {
		c := &ConnectorRecord{Name: name, Category: "virtualization", Type: "proxmox", URL: "https://example.com", Enabled: enabled}
		if err := s.CreateConnector(ctx, c); err != nil {
			t.Fatalf("CreateConnector(%s) error: %v", name, err)
		}
		return c.ID
	}
	schedule := 60

	dueID := mustCreate("due-past", true)
	if err := s.UpdateConnector(ctx, dueID, map[string]any{"schedule_seconds": schedule, "next_run_at": past}); err != nil {
		t.Fatalf("UpdateConnector(due) error: %v", err)
	}

	dueNullID := mustCreate("due-null-next-run", true)
	if err := s.UpdateConnector(ctx, dueNullID, map[string]any{"schedule_seconds": schedule}); err != nil {
		t.Fatalf("UpdateConnector(due-null) error: %v", err)
	}

	notDueFutureID := mustCreate("not-due-future", true)
	if err := s.UpdateConnector(ctx, notDueFutureID, map[string]any{"schedule_seconds": schedule, "next_run_at": future}); err != nil {
		t.Fatalf("UpdateConnector(not-due-future) error: %v", err)
	}

	_ = mustCreate("manual-only", true) // schedule_seconds left NULL

	disabledID := mustCreate("disabled", false)
	if err := s.UpdateConnector(ctx, disabledID, map[string]any{"schedule_seconds": schedule, "next_run_at": past}); err != nil {
		t.Fatalf("UpdateConnector(disabled) error: %v", err)
	}

	due, err := s.ListDueConnectors(ctx, nowStr, 20)
	if err != nil {
		t.Fatalf("ListDueConnectors() error: %v", err)
	}
	gotIDs := map[string]bool{}
	for _, c := range due {
		gotIDs[c.ID] = true
	}
	if len(due) != 2 || !gotIDs[dueID] || !gotIDs[dueNullID] {
		t.Fatalf("ListDueConnectors() = %+v, want exactly [dueID, dueNullID]", due)
	}
}

// TestListDueConnectorsExcludesMaintenanceWindow verifies a connector with an
// active maintenance window is skipped by the scheduler query even though it
// is otherwise due, per #236.
func TestListDueConnectorsExcludesMaintenanceWindow(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	now := time.Now().UTC()
	past := now.Add(-time.Hour).Format(time.RFC3339)
	nowStr := now.Format(time.RFC3339)
	schedule := 60

	mustCreate := func(name string) string {
		c := &ConnectorRecord{Name: name, Category: "virtualization", Type: "proxmox", URL: "https://example.com", Enabled: true}
		if err := s.CreateConnector(ctx, c); err != nil {
			t.Fatalf("CreateConnector(%s) error: %v", name, err)
		}
		if err := s.UpdateConnector(ctx, c.ID, map[string]any{"schedule_seconds": schedule, "next_run_at": past}); err != nil {
			t.Fatalf("UpdateConnector(%s) error: %v", name, err)
		}
		return c.ID
	}

	windowedID := mustCreate("windowed")
	dueID := mustCreate("plain-due")

	if err := s.CreateMaintenanceWindow(ctx, &MaintenanceWindowRecord{
		ConnectorID: windowedID, StartsAt: nowStr, EndsAt: now.Add(time.Hour).Format(time.RFC3339), CreatedBy: "user-1",
	}); err != nil {
		t.Fatalf("CreateMaintenanceWindow() error: %v", err)
	}

	due, err := s.ListDueConnectors(ctx, nowStr, 20)
	if err != nil {
		t.Fatalf("ListDueConnectors() error: %v", err)
	}
	gotIDs := map[string]bool{}
	for _, c := range due {
		gotIDs[c.ID] = true
	}
	if gotIDs[windowedID] {
		t.Fatalf("ListDueConnectors() included windowed connector %s, want excluded", windowedID)
	}
	if !gotIDs[dueID] {
		t.Fatalf("ListDueConnectors() = %+v, want plain due connector %s included", due, dueID)
	}
}

// TestListDueConnectorsIncludesExpiredMaintenanceWindow verifies a connector
// whose maintenance window has already ended is due again, matching
// GetActiveMaintenanceWindow's own ends_at > now cutoff.
func TestListDueConnectorsIncludesExpiredMaintenanceWindow(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	now := time.Now().UTC()
	past := now.Add(-time.Hour).Format(time.RFC3339)
	nowStr := now.Format(time.RFC3339)
	schedule := 60

	c := &ConnectorRecord{Name: "was-windowed", Category: "virtualization", Type: "proxmox", URL: "https://example.com", Enabled: true}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{"schedule_seconds": schedule, "next_run_at": past}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}
	if err := s.CreateMaintenanceWindow(ctx, &MaintenanceWindowRecord{
		ConnectorID: c.ID, StartsAt: past, EndsAt: past, CreatedBy: "user-1",
	}); err != nil {
		t.Fatalf("CreateMaintenanceWindow() error: %v", err)
	}

	due, err := s.ListDueConnectors(ctx, nowStr, 20)
	if err != nil {
		t.Fatalf("ListDueConnectors() error: %v", err)
	}
	found := false
	for _, d := range due {
		if d.ID == c.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("ListDueConnectors() = %+v, want connector %s with an expired window included", due, c.ID)
	}
}

func TestClaimDueConnector(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	now := time.Now().UTC()
	schedule := 30
	for _, tc := range []struct {
		name                            string
		enabled, scheduled, maintenance bool
		next                            string
		want                            bool
	}{
		{name: "unset", enabled: true, scheduled: true, want: true},
		{name: "due", enabled: true, scheduled: true, next: now.Add(-time.Minute).Format(time.RFC3339), want: true},
		{name: "future", enabled: true, scheduled: true, next: now.Add(time.Hour).Format(time.RFC3339)},
		{name: "disabled", scheduled: true},
		{name: "manual", enabled: true},
		{name: "maintenance", enabled: true, scheduled: true, maintenance: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := &ConnectorRecord{Name: tc.name, Category: "networking", Type: "test", Enabled: tc.enabled, NextRunAt: tc.next}
			if tc.scheduled {
				c.ScheduleSeconds = &schedule
			}
			if err := s.CreateConnector(ctx, c); err != nil {
				t.Fatal(err)
			}
			if tc.maintenance {
				if err := s.CreateMaintenanceWindow(ctx, &MaintenanceWindowRecord{ConnectorID: c.ID, StartsAt: now.Format(time.RFC3339), EndsAt: now.Add(time.Hour).Format(time.RFC3339), CreatedBy: "user-1"}); err != nil {
					t.Fatal(err)
				}
			}
			lease := now.Add(6 * time.Minute).Format(time.RFC3339)
			claimed, err := s.ClaimDueConnector(ctx, c.ID, now.Format(time.RFC3339), lease)
			if err != nil || claimed != tc.want {
				t.Fatalf("claimed = %v, error = %v; want %v", claimed, err, tc.want)
			}
			claimed, err = s.ClaimDueConnector(ctx, c.ID, now.Format(time.RFC3339), lease)
			if err != nil || claimed {
				t.Fatalf("second claim = %v, error = %v", claimed, err)
			}
			if tc.want {
				claimed, err = s.ClaimDueConnector(ctx, c.ID, now.Add(7*time.Minute).Format(time.RFC3339), now.Add(13*time.Minute).Format(time.RFC3339))
				if err != nil || !claimed {
					t.Fatalf("expired lease claim = %v, error = %v", claimed, err)
				}
			}
		})
	}
}
