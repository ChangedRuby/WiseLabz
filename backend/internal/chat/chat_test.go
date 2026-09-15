package chat

import "testing"

func TestSplitSections(t *testing.T) {
	content := "# Title\n\nintro\n\n## Overview\n\nThis is the overview.\n\n## Dependencies\n\n- foo\n- bar\n"
	sections := SplitSections(content)
	if len(sections) != 3 {
		t.Fatalf("expected 3 sections (preamble + 2 headings), got %d: %+v", len(sections), sections)
	}
	if sections[1].Key != "Overview" || sections[1].Content != "This is the overview." {
		t.Fatalf("unexpected section 1: %+v", sections[1])
	}
	if sections[2].Key != "Dependencies" {
		t.Fatalf("unexpected section 2: %+v", sections[2])
	}
}

func TestCosineSimilarityRanksClosestVectorHighest(t *testing.T) {
	query := []float32{1, 0}
	near := []float32{0.9, 0.1}
	far := []float32{0, 1}

	scoreClose := cosineSimilarity(query, near)
	scoreFar := cosineSimilarity(query, far)
	if scoreClose <= scoreFar {
		t.Fatalf("expected close vector to score higher: close=%v far=%v", scoreClose, scoreFar)
	}
}

func TestPackUnpackVectorRoundTrips(t *testing.T) {
	v := []float32{1.5, -2.25, 0, 3.125}
	got := unpackVector(packVector(v))
	if len(got) != len(v) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(v))
	}
	for i := range v {
		if got[i] != v[i] {
			t.Fatalf("index %d: got %v want %v", i, got[i], v[i])
		}
	}
}
