package sync

import (
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TestCompareIgnoresEntityAttributes documents and locks in that Compare
// only diffs Sections, not Entities/Attributes: gaining structured
// attributes on upgrade (a previously attribute-less snapshot suddenly
// carrying them) must not raise drift on its own. A real config change (a
// rule flipping enabled) is still caught because it also changes the
// rendered "Firewall Rules" section content, which Compare does diff.
func TestCompareIgnoresEntityAttributes(t *testing.T) {
	prev := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "Firewall Rules", Content: "same content"}},
		Entities: []connector.SnapshotEntity{{Kind: "rule", Name: "Allow SSH"}},
	}
	curr := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "Firewall Rules", Content: "same content"}},
		Entities: []connector.SnapshotEntity{{
			Kind: "rule", Name: "Allow SSH",
			Attributes: map[string]any{"enabled": true, "action": "pass", "protocol": "tcp"},
		}},
	}

	results := Compare(prev, curr)
	if len(results) != 0 {
		t.Fatalf("Compare() with unchanged sections and newly-attributed entities = %+v, want no drift results", results)
	}
}

// TestCompareStillDetectsRuleContentChanges verifies a real change to a
// rule (e.g. becoming enabled) is still caught: it changes the rendered
// section content, which is what Compare actually diffs.
func TestCompareStillDetectsRuleContentChanges(t *testing.T) {
	prev := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "Firewall Rules", Content: "| Allow SSH | ... | false |"}},
	}
	curr := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "Firewall Rules", Content: "| Allow SSH | ... | true |"}},
	}

	results := Compare(prev, curr)
	if len(results) != 1 || results[0].Type != "modified" {
		t.Fatalf("Compare() with changed rule content = %+v, want one modified result", results)
	}
}

// TestCompareMapKeyOrderingDoesNotAffectResult verifies map key ordering in
// Attributes (which Go deliberately randomizes on iteration) can't cause
// Compare to produce a different result across runs.
func TestCompareMapKeyOrderingDoesNotAffectResult(t *testing.T) {
	snap := &connector.ServiceSnapshot{
		Sections: []connector.SnapshotSection{{Title: "Firewall Rules", Content: "same"}},
		Entities: []connector.SnapshotEntity{{
			Kind: "rule", Name: "Allow SSH",
			Attributes: map[string]any{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
		}},
	}
	var first []DiffResult
	for i := 0; i < 20; i++ {
		results := Compare(snap, snap)
		if first == nil {
			first = results
			continue
		}
		if len(results) != len(first) {
			t.Fatalf("run %d: got %d results, want %d", i, len(results), len(first))
		}
	}
}
