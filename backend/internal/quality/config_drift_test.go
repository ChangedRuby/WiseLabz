package quality

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// createDriftSnapshot creates a snapshot with an explicit fetchedAt so
// GetLatestSnapshot's ORDER BY fetched_at DESC has an unambiguous winner
// even when snapshots are created within the same wall-clock second.
func createDriftSnapshot(t *testing.T, s *store.Store, connectorID, content, fetchedAt string) *store.SnapshotRecord {
	t.Helper()
	data, err := json.Marshal(connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "Firewall Rules", Content: content}},
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	snap := &store.SnapshotRecord{ConnectorID: connectorID, Data: string(data), FetchedAt: fetchedAt}
	if err := s.CreateSnapshot(context.Background(), snap); err != nil {
		t.Fatalf("CreateSnapshot() error: %v", err)
	}
	return snap
}

func TestCheckConfigDriftNoFalsePositiveWithoutPin(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	conn := createConnector(t, s, "platform-team")
	createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | false |", "2020-01-01T00:00:00Z")

	checker := NewChecker(s, nil, nil, RotationConfig{})
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatalf("RunForConnector() error: %v", err)
	}
	if got := findings(t, s, conn.ID, "config_drift", "open"); len(got) != 0 {
		t.Fatalf("open config_drift findings without a pinned golden snapshot = %d, want 0", len(got))
	}
}

func TestCheckConfigDriftNoFalsePositiveWhenUnchanged(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	conn := createConnector(t, s, "platform-team")
	golden := createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | false |", "2020-01-01T00:00:00Z")
	if err := s.PinGoldenSnapshot(ctx, &store.GoldenSnapshotRecord{ConnectorID: conn.ID, SnapshotID: golden.ID, PinnedBy: "user-1"}); err != nil {
		t.Fatalf("PinGoldenSnapshot() error: %v", err)
	}
	// A later snapshot with identical section content: no drift.
	createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | false |", "2021-01-01T00:00:00Z")

	checker := NewChecker(s, nil, nil, RotationConfig{})
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatalf("RunForConnector() error: %v", err)
	}
	if got := findings(t, s, conn.ID, "config_drift", "open"); len(got) != 0 {
		t.Fatalf("open config_drift findings for unchanged snapshot = %d, want 0", len(got))
	}
}

func TestCheckConfigDriftDetectsAndAutoResolves(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	conn := createConnector(t, s, "platform-team")
	golden := createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | false |", "2020-01-01T00:00:00Z")
	if err := s.PinGoldenSnapshot(ctx, &store.GoldenSnapshotRecord{ConnectorID: conn.ID, SnapshotID: golden.ID, PinnedBy: "user-1"}); err != nil {
		t.Fatalf("PinGoldenSnapshot() error: %v", err)
	}
	// A later snapshot deviates from golden: the rule became enabled.
	createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | true |", "2021-01-01T00:00:00Z")

	checker := NewChecker(s, nil, nil, RotationConfig{})
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	open := findings(t, s, conn.ID, "config_drift", "open")
	if len(open) != 1 {
		t.Fatalf("open config_drift findings = %d, want 1: %#v", len(open), open)
	}
	if open[0].Severity == "" || open[0].Title == "" || open[0].Description == "" {
		t.Fatalf("config_drift finding = %#v, want populated severity/title/description", open[0])
	}

	// Re-pinning the drifted snapshot as the new golden baseline resolves
	// the finding: it's no longer drift once it's the baseline itself.
	latest, err := s.GetLatestSnapshot(ctx, conn.ID)
	if err != nil {
		t.Fatalf("GetLatestSnapshot() error: %v", err)
	}
	if err := s.PinGoldenSnapshot(ctx, &store.GoldenSnapshotRecord{ConnectorID: conn.ID, SnapshotID: latest.ID, PinnedBy: "user-2"}); err != nil {
		t.Fatalf("PinGoldenSnapshot() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, conn.ID, "config_drift", "open"); len(got) != 0 {
		t.Fatalf("open config_drift findings after re-pin = %d, want 0", len(got))
	}
	if got := findings(t, s, conn.ID, "config_drift", "resolved"); len(got) != 1 {
		t.Fatalf("resolved config_drift findings = %d, want 1", len(got))
	}
}

func TestCheckConfigDriftResolvesOnUnpin(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	conn := createConnector(t, s, "platform-team")
	golden := createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | false |", "2020-01-01T00:00:00Z")
	if err := s.PinGoldenSnapshot(ctx, &store.GoldenSnapshotRecord{ConnectorID: conn.ID, SnapshotID: golden.ID, PinnedBy: "user-1"}); err != nil {
		t.Fatalf("PinGoldenSnapshot() error: %v", err)
	}
	createDriftSnapshot(t, s, conn.ID, "| Allow SSH | ... | true |", "2021-01-01T00:00:00Z")

	checker := NewChecker(s, nil, nil, RotationConfig{})
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	if got := findings(t, s, conn.ID, "config_drift", "open"); len(got) != 1 {
		t.Fatalf("open config_drift findings before unpin = %d, want 1", len(got))
	}

	if err := s.UnpinGoldenSnapshot(ctx, conn.ID); err != nil {
		t.Fatalf("UnpinGoldenSnapshot() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatalf("RunForConnector() after unpin error: %v", err)
	}
	if got := findings(t, s, conn.ID, "config_drift", "open"); len(got) != 0 {
		t.Fatalf("open config_drift findings after unpin = %d, want 0", len(got))
	}
}
