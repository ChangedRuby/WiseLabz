package store

import (
	"context"
	"errors"
	"testing"
)

// newCascadeTestStore opens SQLite via OpenDB (production's own path, with
// the foreign_keys pragma on) instead of newDocTestStore, so ON DELETE
// CASCADE actually fires — same reasoning as
// TestDeleteUserCascadesMFAFactorsAndRecoveryCodes in mfa_test.go.
func newCascadeTestStore(t *testing.T) *Store {
	t.Helper()
	db, err := OpenDB("sqlite", "file:"+migratedSQLite(t))
	if err != nil {
		t.Fatalf("OpenDB() error: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return New(db, "sqlite")
}

func TestRunbookRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	r := &RunbookRecord{
		Title:       "Restart the service",
		Body:        "Step 1: check logs.",
		TargetType:  "change_type",
		TargetValue: "config_update",
	}
	created, err := s.CreateRunbook(ctx, r)
	if err != nil {
		t.Fatalf("CreateRunbook() error: %v", err)
	}

	got, err := s.GetRunbook(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRunbook() error: %v", err)
	}
	if got.Body != "Step 1: check logs." {
		t.Fatalf("Body = %q, want %q", got.Body, "Step 1: check logs.")
	}

	updated, err := s.UpdateRunbook(ctx, created.ID, map[string]any{"body": "Step 1: check logs. Step 2: restart."})
	if err != nil {
		t.Fatalf("UpdateRunbook() error: %v", err)
	}
	if updated.Body != "Step 1: check logs. Step 2: restart." {
		t.Fatalf("updated.Body = %q, want %q", updated.Body, "Step 1: check logs. Step 2: restart.")
	}

	got, err = s.GetRunbook(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRunbook() after update error: %v", err)
	}
	if got.Body != "Step 1: check logs. Step 2: restart." {
		t.Fatalf("Body after update = %q, want %q", got.Body, "Step 1: check logs. Step 2: restart.")
	}

	if err := s.DeleteRunbook(ctx, created.ID); err != nil {
		t.Fatalf("DeleteRunbook() error: %v", err)
	}
	if _, err := s.GetRunbook(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetRunbook() after delete error = %v, want ErrNotFound", err)
	}
}

func TestRunbookStepsRoundTrip(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	created, saved, err := s.CreateRunbookWithSteps(ctx, &RunbookRecord{
		Title:       "Restart the service",
		TargetType:  "change_type",
		TargetValue: "config_update",
	}, []*RunbookStepRecord{
		{Title: "Restart VM", ConnectorID: c.ID, Verb: "restart", EntityRef: "100"},
		{Title: "Start VM", ConnectorID: c.ID, Verb: "start"},
	})
	if err != nil {
		t.Fatalf("CreateRunbookWithSteps() error: %v", err)
	}
	if len(saved) != 2 {
		t.Fatalf("len(saved) = %d, want 2", len(saved))
	}
	if saved[0].Position != 0 || saved[1].Position != 1 {
		t.Fatalf("positions = %d, %d, want 0, 1", saved[0].Position, saved[1].Position)
	}
	firstStepID := saved[0].ID

	got, err := s.ListRunbookStepsFor(ctx, created.ID)
	if err != nil {
		t.Fatalf("ListRunbookStepsFor() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}

	step, err := s.GetRunbookStep(ctx, created.ID, firstStepID)
	if err != nil {
		t.Fatalf("GetRunbookStep() error: %v", err)
	}
	if step.Title != "Restart VM" || step.EntityRef != "100" {
		t.Fatalf("step = %+v, want Restart VM / entityRef 100", step)
	}

	// UpdateRunbookWithSteps with replaceSteps=false must leave steps
	// untouched (absent "steps" key semantics at the handler layer).
	_, unchanged, err := s.UpdateRunbookWithSteps(ctx, created.ID, map[string]any{"title": "renamed"}, nil, false)
	if err != nil {
		t.Fatalf("UpdateRunbookWithSteps(replaceSteps=false) error: %v", err)
	}
	if len(unchanged) != 2 || unchanged[0].ID != firstStepID {
		t.Fatalf("unchanged steps = %+v, want the original 2 steps preserved", unchanged)
	}

	// Replacing with one kept step (by ID) and one new step should keep the
	// kept step's ID and mint a fresh one for the new step, reordered.
	_, replaced, err := s.UpdateRunbookWithSteps(ctx, created.ID, nil, []*RunbookStepRecord{
		{Title: "New first step", ConnectorID: c.ID, Verb: "stop"},
		{ID: firstStepID, Title: "Restart VM renamed", ConnectorID: c.ID, Verb: "restart", EntityRef: "100"},
	}, true)
	if err != nil {
		t.Fatalf("UpdateRunbookWithSteps(replaceSteps=true) error: %v", err)
	}
	if len(replaced) != 2 {
		t.Fatalf("len(replaced) = %d, want 2", len(replaced))
	}
	if replaced[0].ID == firstStepID {
		t.Fatalf("replaced[0].ID = %q should be a fresh ID, not the kept step's ID", replaced[0].ID)
	}
	if replaced[1].ID != firstStepID || replaced[1].Title != "Restart VM renamed" {
		t.Fatalf("replaced[1] = %+v, want kept step ID %q with updated title", replaced[1], firstStepID)
	}

	// A step ID that never belonged to this runbook must not be "kept" —
	// it gets replaced with a fresh ID instead.
	_, foreign, err := s.UpdateRunbookWithSteps(ctx, created.ID, nil, []*RunbookStepRecord{
		{ID: "not-a-real-step-id", Title: "Sneaky", ConnectorID: c.ID, Verb: "start"},
	}, true)
	if err != nil {
		t.Fatalf("UpdateRunbookWithSteps() error: %v", err)
	}
	if len(foreign) != 1 || foreign[0].ID == "not-a-real-step-id" {
		t.Fatalf("foreign = %+v, want a freshly generated ID", foreign)
	}
}

func TestRunbookStepsCascadeOnRunbookDelete(t *testing.T) {
	ctx := context.Background()
	s := newCascadeTestStore(t)

	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	created, _, err := s.CreateRunbookWithSteps(ctx, &RunbookRecord{
		Title: "R", TargetType: "change_type", TargetValue: "x",
	}, []*RunbookStepRecord{{Title: "Restart", ConnectorID: c.ID, Verb: "restart"}})
	if err != nil {
		t.Fatalf("CreateRunbookWithSteps() error: %v", err)
	}

	if err := s.DeleteRunbook(ctx, created.ID); err != nil {
		t.Fatalf("DeleteRunbook() error: %v", err)
	}

	steps, err := s.ListRunbookStepsFor(ctx, created.ID)
	if err != nil {
		t.Fatalf("ListRunbookStepsFor() after delete error: %v", err)
	}
	if len(steps) != 0 {
		t.Fatalf("len(steps) after runbook delete = %d, want 0 (cascade)", len(steps))
	}
}

func TestRunbookStepsCascadeOnConnectorDelete(t *testing.T) {
	ctx := context.Background()
	s := newCascadeTestStore(t)

	c := &ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	created, _, err := s.CreateRunbookWithSteps(ctx, &RunbookRecord{
		Title: "R", TargetType: "change_type", TargetValue: "x",
	}, []*RunbookStepRecord{{Title: "Restart", ConnectorID: c.ID, Verb: "restart"}})
	if err != nil {
		t.Fatalf("CreateRunbookWithSteps() error: %v", err)
	}

	if err := s.DeleteConnector(ctx, c.ID); err != nil {
		t.Fatalf("DeleteConnector() error: %v", err)
	}

	steps, err := s.ListRunbookStepsFor(ctx, created.ID)
	if err != nil {
		t.Fatalf("ListRunbookStepsFor() after connector delete error: %v", err)
	}
	if len(steps) != 0 {
		t.Fatalf("len(steps) after connector delete = %d, want 0 (cascade)", len(steps))
	}

	// The runbook itself must survive — only its step referencing the
	// deleted connector is gone.
	if _, err := s.GetRunbook(ctx, created.ID); err != nil {
		t.Fatalf("GetRunbook() after connector delete error: %v", err)
	}
}

func TestGetRunbookByTarget(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	r := &RunbookRecord{
		Title:       "High severity triage",
		Body:        "Page the on-call engineer.",
		TargetType:  "alert_severity",
		TargetValue: "critical",
	}
	if _, err := s.CreateRunbook(ctx, r); err != nil {
		t.Fatalf("CreateRunbook() error: %v", err)
	}

	got, err := s.GetRunbookByTarget(ctx, "alert_severity", "critical")
	if err != nil {
		t.Fatalf("GetRunbookByTarget() error: %v", err)
	}
	if got.Title != "High severity triage" {
		t.Fatalf("Title = %q, want %q", got.Title, "High severity triage")
	}

	if _, err := s.GetRunbookByTarget(ctx, "alert_severity", "low"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetRunbookByTarget(not found) error = %v, want ErrNotFound", err)
	}
}
