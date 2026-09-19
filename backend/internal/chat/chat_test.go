package chat

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

type stubEmbedder struct {
	err error
}

func (e stubEmbedder) Name() string { return "stub" }
func (e stubEmbedder) Embed(_ context.Context, texts []string) ([][]float32, error) {
	if e.err != nil {
		return nil, e.err
	}
	out := make([][]float32, len(texts))
	for i := range texts {
		out[i] = []float32{1, float32(i)}
	}
	return out, nil
}

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

func TestSyncDocEmbeddingsKeepsOldRowsWhenEmbedFails(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/test.db?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	if err := store.RunMigrations(db, "sqlite", slog.New(slog.NewTextHandler(io.Discard, nil))); err != nil {
		t.Fatalf("migrations: %v", err)
	}
	s := store.New(db, "sqlite")
	if err := s.CreateDoc(ctx, &store.DocRecord{ID: "d1", Title: "d1", Kind: "lab", Content: "x"}); err != nil {
		t.Fatalf("CreateDoc: %v", err)
	}

	if err := SyncDocEmbeddings(ctx, s, stubEmbedder{}, "m", "d1", "## A\n\none\n\n## B\n\ntwo\n"); err != nil {
		t.Fatalf("initial sync: %v", err)
	}
	if err := SyncDocEmbeddings(ctx, s, stubEmbedder{err: errors.New("boom")}, "m", "d1", "## C\n\nthree\n"); err == nil {
		t.Fatal("expected embed error")
	}
	rows, err := s.ListDocSectionEmbeddings(ctx, "", "d1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected old 2 rows preserved after failed embed, got %d", len(rows))
	}

	if err := SyncDocEmbeddings(ctx, s, stubEmbedder{}, "m", "d1", "## C\n\nthree\n"); err != nil {
		t.Fatalf("resync: %v", err)
	}
	rows, _ = s.ListDocSectionEmbeddings(ctx, "", "d1")
	if len(rows) != 1 || rows[0].SectionKey != "C" {
		t.Fatalf("expected replaced single row C, got %+v", rows)
	}
}
