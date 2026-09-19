package store

import (
	"context"
	"testing"
	"time"
)

func seedConnectorForChanges(t *testing.T, s *Store) string {
	t.Helper()
	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	return c.ID
}

func TestChangeRelatedServiceIDsAndPatternIDRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	serviceID := seedConnectorForChanges(t, s)

	c := &ChangeRecord{
		ServiceID:         serviceID,
		ChangeType:        "modified",
		Severity:          "warning",
		Summary:           "Section modified: VMs",
		RelatedServiceIDs: `["dep-1","dep-2"]`,
		PatternID:         "abc123",
	}
	if err := s.CreateChange(ctx, c); err != nil {
		t.Fatalf("CreateChange() error: %v", err)
	}

	got, err := s.GetChange(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetChange() error: %v", err)
	}
	if got.RelatedServiceIDs != `["dep-1","dep-2"]` {
		t.Errorf("RelatedServiceIDs = %q, want the seeded JSON", got.RelatedServiceIDs)
	}
	if got.PatternID != "abc123" {
		t.Errorf("PatternID = %q, want abc123", got.PatternID)
	}
}

func TestChangeRelatedServiceIDsDefaultsToEmptyArray(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	serviceID := seedConnectorForChanges(t, s)

	c := &ChangeRecord{ServiceID: serviceID, ChangeType: "added", Severity: "info", Summary: "x"}
	if err := s.CreateChange(ctx, c); err != nil {
		t.Fatalf("CreateChange() error: %v", err)
	}
	got, err := s.GetChange(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetChange() error: %v", err)
	}
	if got.RelatedServiceIDs != "[]" {
		t.Errorf("RelatedServiceIDs = %q, want [] default", got.RelatedServiceIDs)
	}
	if got.PatternID != "" {
		t.Errorf("PatternID = %q, want empty when not set", got.PatternID)
	}
}

func TestCountRecentChangesByPattern(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	serviceID := seedConnectorForChanges(t, s)

	old := &ChangeRecord{
		ServiceID: serviceID, ChangeType: "modified", Severity: "info", Summary: "x",
		PatternID: "pattern-a", DetectedAt: time.Now().Add(-2 * time.Hour).UTC().Format(time.RFC3339),
	}
	if err := s.CreateChange(ctx, old); err != nil {
		t.Fatalf("CreateChange(old): %v", err)
	}
	recent := &ChangeRecord{
		ServiceID: serviceID, ChangeType: "modified", Severity: "info", Summary: "x",
		PatternID: "pattern-a", DetectedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := s.CreateChange(ctx, recent); err != nil {
		t.Fatalf("CreateChange(recent): %v", err)
	}

	since := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)

	// Excludes the old (outside window) and the recent change itself.
	count, err := s.CountRecentChangesByPattern(ctx, serviceID, "pattern-a", since, recent.ID)
	if err != nil {
		t.Fatalf("CountRecentChangesByPattern: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0 (only match is excluded and the old one is out of window)", count)
	}

	// Without excluding the recent change, it counts.
	count, err = s.CountRecentChangesByPattern(ctx, serviceID, "pattern-a", since, "")
	if err != nil {
		t.Fatalf("CountRecentChangesByPattern: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}

	// Different pattern, or empty pattern, never matches.
	if count, err := s.CountRecentChangesByPattern(ctx, serviceID, "pattern-b", since, ""); err != nil || count != 0 {
		t.Errorf("CountRecentChangesByPattern(other pattern) = (%d, %v), want (0, nil)", count, err)
	}
	if count, err := s.CountRecentChangesByPattern(ctx, serviceID, "", since, ""); err != nil || count != 0 {
		t.Errorf("CountRecentChangesByPattern(empty pattern) = (%d, %v), want (0, nil)", count, err)
	}
}

func TestCountRecentChangePatterns(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	serviceID := seedConnectorForChanges(t, s)
	otherID := seedConnectorForChanges(t, s)
	now := time.Now().UTC()
	records := []*ChangeRecord{
		{ServiceID: serviceID, PatternID: "a", DetectedAt: now.Format(time.RFC3339)},
		{ServiceID: serviceID, PatternID: "a", DetectedAt: now.Format(time.RFC3339)},
		{ServiceID: serviceID, PatternID: "b", DetectedAt: now.Format(time.RFC3339)},
		{ServiceID: serviceID, PatternID: "a", DetectedAt: now.Add(-2 * time.Hour).Format(time.RFC3339)},
		{ServiceID: otherID, PatternID: "a", DetectedAt: now.Format(time.RFC3339)},
		{ServiceID: serviceID, PatternID: "", DetectedAt: now.Format(time.RFC3339)},
	}
	for _, record := range records {
		record.ChangeType = "modified"
		record.Severity = "info"
	}
	if err := s.CreateChanges(ctx, records); err != nil {
		t.Fatal(err)
	}
	counts, err := s.CountRecentChangePatterns(ctx, serviceID, now.Add(-time.Hour).Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
	if len(counts) != 2 || counts["a"] != 2 || counts["b"] != 1 {
		t.Fatalf("unexpected grouped counts: %v", counts)
	}
}
