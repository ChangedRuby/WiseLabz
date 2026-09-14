package sync

import (
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestCompareOrderIsDeterministic(t *testing.T) {
	prev := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{
			{Title: "Z", Content: "old-z"},
			{Title: "A", Content: "old-a"},
			{Title: "M", Content: "old-m"},
		},
	}
	curr := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{
			{Title: "Z", Content: "new-z"},
			{Title: "A", Content: "new-a"},
			{Title: "M", Content: "new-m"},
		},
	}

	var firstOrder []string
	for i := 0; i < 20; i++ {
		results := Compare(prev, curr)
		order := make([]string, len(results))
		for j, r := range results {
			order[j] = r.Summary
		}
		if firstOrder == nil {
			firstOrder = order
			continue
		}
		if len(order) != len(firstOrder) {
			t.Fatalf("run %d: got %d results, want %d", i, len(order), len(firstOrder))
		}
		for j := range order {
			if order[j] != firstOrder[j] {
				t.Fatalf("run %d: order = %v, want %v (non-deterministic)", i, order, firstOrder)
			}
		}
	}
}

func TestCompareTagsRelatedServiceIDs(t *testing.T) {
	now := time.Now()
	prev := &connector.ServiceSnapshot{
		Sections:  []connector.SnapshotSection{{Title: "VMs", Content: "vm1"}},
		FetchedAt: now,
	}
	curr := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "VMs", Content: "vm1\nvm2"}},
		Dependencies: []connector.ServiceDependency{
			{Kind: "host", Name: "pve1", Ref: "connector-a"},
			{Kind: "network", Name: "lan", Ref: "connector-b"},
			{Kind: "storage", Name: "unresolved"}, // no Ref: excluded
		},
		FetchedAt: now,
	}

	results := Compare(prev, curr)
	if len(results) != 1 {
		t.Fatalf("Compare() = %d results, want 1", len(results))
	}
	got := results[0].RelatedServiceIDs
	if len(got) != 2 || got[0] != "connector-a" || got[1] != "connector-b" {
		t.Errorf("RelatedServiceIDs = %v, want [connector-a connector-b]", got)
	}
}

func TestCompareDeduplicatesRelatedServiceIDs(t *testing.T) {
	curr := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "New", Content: "x"}},
		Dependencies: []connector.ServiceDependency{
			{Kind: "host", Name: "a", Ref: "svc-1"},
			{Kind: "network", Name: "b", Ref: "svc-1"},
		},
	}
	results := Compare(&connector.ServiceSnapshot{}, curr)
	if len(results) != 1 {
		t.Fatalf("Compare() = %d results, want 1", len(results))
	}
	if got := results[0].RelatedServiceIDs; len(got) != 1 || got[0] != "svc-1" {
		t.Errorf("RelatedServiceIDs = %v, want [svc-1] (deduplicated)", got)
	}
}

func TestCompareRemovedSectionUsesPrevDependencies(t *testing.T) {
	prev := &connector.ServiceSnapshot{
		Sections:     []connector.SnapshotSection{{Title: "Gone", Content: "x"}},
		Dependencies: []connector.ServiceDependency{{Kind: "host", Name: "a", Ref: "svc-9"}},
	}
	results := Compare(prev, &connector.ServiceSnapshot{})
	if len(results) != 1 || results[0].Type != "removed" {
		t.Fatalf("Compare() = %+v, want one removed result", results)
	}
	if got := results[0].RelatedServiceIDs; len(got) != 1 || got[0] != "svc-9" {
		t.Errorf("RelatedServiceIDs = %v, want [svc-9]", got)
	}
}
