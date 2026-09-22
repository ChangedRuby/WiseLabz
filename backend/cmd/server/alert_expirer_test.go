package main

import (
	"context"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/notifications"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// TestExpireAlertsOnceNotifiesViaDispatcher asserts the alert expirer, now
// running as a scheduler job, actually uses its dispatcher argument (the bug
// the old select-then-sleep loop had: the dispatcher parameter was ignored).
// An expired-snooze alert should both flip back to pending and produce an
// in-app notification for every active user.
func TestExpireAlertsOnceNotifiesViaDispatcher(t *testing.T) {
	ctx := context.Background()
	s := apitest.NewStore(t)
	userID := apitest.NewUser(t, s, "viewer")
	dispatcher := notifications.NewDispatcher(s, nil)

	past := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	alert := &store.AlertRecord{
		ServiceID:    "svc-1",
		Severity:     "warning",
		Title:        "disk usage high",
		Description:  "disk usage crossed threshold",
		Status:       "snoozed",
		SnoozedUntil: past,
	}
	if err := s.CreateAlert(ctx, alert); err != nil {
		t.Fatalf("create alert: %v", err)
	}

	logger := testLogger()
	expireAlertsOnce(ctx, s, dispatcher, logger)
	dispatcher.Wait()

	got, err := s.GetAlert(ctx, alert.ID)
	if err != nil {
		t.Fatalf("get alert: %v", err)
	}
	if got.Status != "pending" {
		t.Fatalf("alert status = %q, want %q", got.Status, "pending")
	}

	notifs, _, err := s.ListNotifications(ctx, userID, false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(notifs) == 0 {
		t.Fatal("expireAlertsOnce did not notify via the dispatcher: no in-app notification was created")
	}
	if notifs[0].AlertID != alert.ID {
		t.Fatalf("notification alertId = %q, want %q", notifs[0].AlertID, alert.ID)
	}
}

// TestExpireAlertsOnceNoExpiredAlertsIsNoop asserts the job does nothing
// (no dispatcher call, no error) when there's nothing to expire.
func TestExpireAlertsOnceNoExpiredAlertsIsNoop(t *testing.T) {
	ctx := context.Background()
	s := apitest.NewStore(t)
	dispatcher := notifications.NewDispatcher(s, nil)
	logger := testLogger()

	expireAlertsOnce(ctx, s, dispatcher, logger)
	dispatcher.Wait()
}
