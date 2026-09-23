package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestJobHealthGetMissingReturnsNoRows(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	if _, err := s.GetJobHealth(ctx, "nope"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetJobHealth(missing) error = %v, want sql.ErrNoRows", err)
	}
}

func TestJobHealthUpsertInsertsThenUpdates(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	rec := JobHealthRecord{
		Name:       "backup",
		CronExpr:   "0 3 * * *",
		LastStatus: "ok",
		LastRunAt:  "2026-01-01T00:00:00Z",
	}
	if err := s.UpsertJobHealth(ctx, rec); err != nil {
		t.Fatalf("UpsertJobHealth(insert) error: %v", err)
	}

	got, err := s.GetJobHealth(ctx, "backup")
	if err != nil {
		t.Fatalf("GetJobHealth() error: %v", err)
	}
	if got.LastStatus != "ok" || got.LastRunAt != "2026-01-01T00:00:00Z" || got.LastError != "" {
		t.Fatalf("GetJobHealth() = %+v, want status ok with no error", got)
	}
	if got.LastSuccessAt != "" || got.LastFailureAt != "" {
		t.Fatalf("GetJobHealth() = %+v, want empty success/failure timestamps (never set)", got)
	}
	if got.UpdatedAt == "" {
		t.Fatal("GetJobHealth() UpdatedAt defaulted to empty, want a timestamp")
	}

	// Update: same name, new status/error — must replace, not duplicate.
	rec2 := JobHealthRecord{
		Name:          "backup",
		CronExpr:      "0 4 * * *",
		LastStatus:    "failing",
		LastError:     "disk full",
		LastRunAt:     "2026-01-02T00:00:00Z",
		LastFailureAt: "2026-01-02T00:00:00Z",
	}
	if err := s.UpsertJobHealth(ctx, rec2); err != nil {
		t.Fatalf("UpsertJobHealth(update) error: %v", err)
	}
	got, err = s.GetJobHealth(ctx, "backup")
	if err != nil {
		t.Fatalf("GetJobHealth() after update error: %v", err)
	}
	if got.LastStatus != "failing" || got.LastError != "disk full" || got.CronExpr != "0 4 * * *" {
		t.Fatalf("GetJobHealth() after update = %+v, want the new failing row", got)
	}

	recs, err := s.ListJobHealth(ctx)
	if err != nil {
		t.Fatalf("ListJobHealth() error: %v", err)
	}
	if len(recs) != 1 {
		t.Fatalf("ListJobHealth() = %d rows, want exactly 1 (update must not duplicate)", len(recs))
	}
}

func TestJobHealthListOrderedByName(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	for _, name := range []string{"sync", "backup", "digest"} {
		if err := s.UpsertJobHealth(ctx, JobHealthRecord{Name: name, CronExpr: "* * * * *", LastStatus: "ok"}); err != nil {
			t.Fatalf("UpsertJobHealth(%s) error: %v", name, err)
		}
	}

	recs, err := s.ListJobHealth(ctx)
	if err != nil {
		t.Fatalf("ListJobHealth() error: %v", err)
	}
	if len(recs) != 3 {
		t.Fatalf("ListJobHealth() = %d rows, want 3", len(recs))
	}
	want := []string{"backup", "digest", "sync"}
	for i, rec := range recs {
		if rec.Name != want[i] {
			t.Fatalf("ListJobHealth()[%d].Name = %q, want %q (alphabetical)", i, rec.Name, want[i])
		}
	}
}

func TestJobHealthDelete(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	if err := s.UpsertJobHealth(ctx, JobHealthRecord{Name: "report:weekly", CronExpr: "0 0 * * 1", LastStatus: "ok"}); err != nil {
		t.Fatalf("UpsertJobHealth() error: %v", err)
	}
	if err := s.DeleteJobHealth(ctx, "report:weekly"); err != nil {
		t.Fatalf("DeleteJobHealth() error: %v", err)
	}
	if _, err := s.GetJobHealth(ctx, "report:weekly"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("GetJobHealth() after delete error = %v, want sql.ErrNoRows", err)
	}

	// Deleting a name that never existed is a no-op, not an error.
	if err := s.DeleteJobHealth(ctx, "never-existed"); err != nil {
		t.Fatalf("DeleteJobHealth(never-existed) error: %v", err)
	}
}
