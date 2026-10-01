package store

import (
	"context"
	"testing"
)

func TestDashboardAggregatesScopedToConnectors(t *testing.T) {
	s := newDocTestStore(t)
	ctx := context.Background()
	a := newTestConnector(t, s, "conn-a")
	b := newTestConnector(t, s, "conn-b")
	if _, err := s.db.ExecContext(ctx, `UPDATE connectors SET status = 'online', last_sync_at = ? WHERE id = ?`, "2026-01-01T00:00:00Z", a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE connectors SET status = 'offline', last_sync_at = ? WHERE id = ?`, "2026-02-01T00:00:00Z", b.ID); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{a.ID, b.ID} {
		if err := s.CreateChange(ctx, &ChangeRecord{ServiceID: id, ChangeType: "config", Severity: "info", Summary: "chg-" + id, Diff: "[]", AffectedDocIDs: "[]", DetectedAt: "2026-03-01T00:00:00Z"}); err != nil {
			t.Fatal(err)
		}
		if err := s.CreateAlert(ctx, &AlertRecord{ServiceID: id, Severity: "info", Title: "t", Status: "pending"}); err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range []struct {
		name            string
		ids             []string
		online, offline int
		alerts, changes int
		lastSync        string
	}{
		{"only A", []string{a.ID}, 1, 0, 1, 1, "2026-01-01T00:00:00Z"},
		{"A and B", []string{a.ID, b.ID}, 1, 1, 2, 2, "2026-02-01T00:00:00Z"},
		{"empty set", nil, 0, 0, 0, 0, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			counts, err := s.CountConnectorsByStatus(ctx, tc.ids)
			if err != nil {
				t.Fatal(err)
			}
			if counts["online"] != tc.online || counts["offline"] != tc.offline {
				t.Errorf("counts = %v", counts)
			}
			if n, err := s.CountAlertsPending(ctx, tc.ids); err != nil || n != tc.alerts {
				t.Errorf("alerts = %d, %v; want %d", n, err, tc.alerts)
			}
			if got, err := s.GetLatestChangeSummaries(ctx, tc.ids, 5, ""); err != nil || len(got) != tc.changes {
				t.Errorf("changes = %d, %v; want %d", len(got), err, tc.changes)
			}
			if got, err := s.GetLastSyncTimestamp(ctx, tc.ids); err != nil || got != tc.lastSync {
				t.Errorf("lastSync = %q, %v; want %q", got, err, tc.lastSync)
			}
		})
	}
}
