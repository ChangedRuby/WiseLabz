package quality

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// fakeNotifier records every NotifyFindingCreated call, for dedup assertions.
type fakeNotifier struct {
	calls  int
	titles []string
}

func (f *fakeNotifier) NotifyFindingCreated(_ context.Context, _, title, _ string) {
	f.calls++
	f.titles = append(f.titles, title)
}

func TestFindingNotificationOpenDispatchesOnce(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "") // triggers ownership_incomplete (info)
	notifier := &fakeNotifier{}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})

	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() error: %v", err)
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls = %d, want 1 for a newly opened finding", notifier.calls)
	}
}

func TestFindingNotificationRedetectAtSameSeverityDoesNotNotify(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "")
	notifier := &fakeNotifier{}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})

	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() first error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() second error: %v", err)
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls = %d, want 1 (re-detection at the same severity must not notify again)", notifier.calls)
	}
}

func TestFindingNotificationEscalationNotifiesAgain(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()
	s := newTestStore(t)
	// 80 days old -> warning under a 90-day/14-day-warn policy.
	c := createConnectorWithRotation(t, s, "proxmox", now.AddDate(0, 0, -80), "", nil)
	notifier := &fakeNotifier{}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	checker.now = func() time.Time { return now }

	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() warning error: %v", err)
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls after warning = %d, want 1", notifier.calls)
	}

	// Advance past the due date: warning -> critical must notify again.
	checker.now = func() time.Time { return now.AddDate(0, 0, 11) }
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() critical error: %v", err)
	}
	if notifier.calls != 2 {
		t.Fatalf("notifier.calls after escalation to critical = %d, want 2", notifier.calls)
	}
}

func TestFindingNotificationDeescalationDoesNotNotify(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()
	s := newTestStore(t)
	// 91 days old -> critical.
	c := createConnectorWithRotation(t, s, "proxmox", now.AddDate(0, 0, -91), "", nil)
	notifier := &fakeNotifier{}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	checker.now = func() time.Time { return now }

	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() critical error: %v", err)
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls after critical = %d, want 1", notifier.calls)
	}

	// A per-connector override widens the window: critical -> warning must
	// not notify again (de-escalation).
	wideMaxAge := 365
	if err := s.UpdateConnector(ctx, c.ID, map[string]any{"rotation_max_age_days": &wideMaxAge}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() after widen error: %v", err)
	}
	open := findings(t, s, c.ID, "credential_rotation", "open")
	if len(open) != 0 {
		t.Fatalf("open findings after widening the window = %d, want 0 (past the 91-day mark but well under 365)", len(open))
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls after de-escalation/resolve = %d, want still 1", notifier.calls)
	}
}

func TestFindingNotificationResolveThenReopenNotifiesAgain(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()
	s := newTestStore(t)
	c := createConnectorWithRotation(t, s, "proxmox", now.AddDate(0, 0, -91), "", nil)
	notifier := &fakeNotifier{}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	checker.now = func() time.Time { return now }

	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls after detect = %d, want 1", notifier.calls)
	}

	if err := s.UpdateConnector(ctx, c.ID, map[string]any{"secret_rotated_at": now.Format(time.RFC3339)}); err != nil {
		t.Fatalf("UpdateConnector() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if notifier.calls != 1 {
		t.Fatalf("notifier.calls after resolve = %d, want still 1", notifier.calls)
	}

	checker.now = func() time.Time { return now.AddDate(0, 0, 91) }
	if err := checker.RunForConnector(ctx, c.ID); err != nil {
		t.Fatalf("RunForConnector() reopen error: %v", err)
	}
	if notifier.calls != 2 {
		t.Fatalf("notifier.calls after resolve-then-reopen = %d, want 2", notifier.calls)
	}
}

func TestFindingNotificationExistingCheckTypeAlsoNotifies(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	if err := s.CreateDoc(ctx, &store.DocRecord{
		Title: "Empty", ServiceID: connector.ID, Content: "x",
	}); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}
	notifier := &fakeNotifier{}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})

	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() error: %v", err)
	}
	found := false
	for _, title := range notifier.titles {
		if strings.Contains(title, "empty") || strings.Contains(title, "Empty") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a notification for the pre-existing 'empty' check type, got titles %v", notifier.titles)
	}
}
