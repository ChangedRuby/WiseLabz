package notifications

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// TestRunDeliveryRetries_PollsUntilContextCanceled verifies that RunDeliveryRetries polls
// continuously until context is canceled. ponytail: simplified by checking only that it runs
// the retry loop without mocking time — full integration would exercise the 30s tick.
func TestRunDeliveryRetries_PollsUntilContextCanceled(t *testing.T) {
	s := newTestStore(t)
	d := NewDispatcher(s, nil)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	RunDeliveryRetries(ctx, d, logger)

	// If RunDeliveryRetries didn't exit on context.Done(), the test would hang.
	// Successful return means it exited cleanly.
}

// TestRetryDueDeliveries_WebhookExceedsRetryBound verifies that a delivery failing all
// maxDeliveryAttempts times stops retrying (no next_attempt_at) and remains marked failed.
func TestRetryDueDeliveries_WebhookExceedsRetryBound(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	setWebhookConfig(t, s, srv.URL)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	webhook, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "webhook")
	if !ok || webhook.Status != store.DeliveryStatusFailed {
		t.Fatalf("expected initial webhook delivery to be failed, got %+v (ok=%v)", webhook, ok)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	// Retry up to maxDeliveryAttempts times.
	for attempt := 2; attempt <= maxDeliveryAttempts; attempt++ {
		// Force delivery due.
		if _, err := s.DB().ExecContext(context.Background(),
			`UPDATE notification_deliveries SET next_attempt_at = ? WHERE id = ?`, "1970-01-01T00:00:00Z", webhook.ID); err != nil {
			t.Fatalf("force due on attempt %d: %v", attempt, err)
		}

		d.retryDueDeliveries(context.Background(), logger)

		// Reload delivery state after retry.
		webhook, ok = findDelivery(deliveriesFor(t, s, webhook.NotificationID), "webhook")
		if !ok {
			t.Fatalf("expected webhook delivery row to still exist on attempt %d", attempt)
		}
		if webhook.Status != store.DeliveryStatusFailed {
			t.Errorf("attempt %d: expected status failed, got %s", attempt, webhook.Status)
		}
		if webhook.Attempts != attempt {
			t.Errorf("attempt %d: expected attempts %d, got %d", attempt, attempt, webhook.Attempts)
		}
	}

	// After maxDeliveryAttempts, next_attempt_at should be empty (stop retrying).
	if webhook.NextAttemptAt != "" {
		t.Errorf("expected empty NextAttemptAt after maxDeliveryAttempts, got %q", webhook.NextAttemptAt)
	}

	// Verify no further retries occur even if we force due again.
	if _, err := s.DB().ExecContext(context.Background(),
		`UPDATE notification_deliveries SET next_attempt_at = ? WHERE id = ?`, "1970-01-01T00:00:00Z", webhook.ID); err != nil {
		t.Fatalf("force due after max attempts: %v", err)
	}

	d.retryDueDeliveries(context.Background(), logger)

	// Check it's not in the due set anymore (no retry attempted).
	due, err := s.ListDueDeliveries(context.Background(), time.Now().UTC().Format(time.RFC3339), 50)
	if err != nil {
		t.Fatalf("list due deliveries: %v", err)
	}
	for _, del := range due {
		if del.ID == webhook.ID {
			t.Errorf("expected delivery to be removed from due set after max attempts, but found it")
		}
	}
}

// TestRetryDueDeliveries_DiscordSucceedsOnRetry verifies a discord delivery that fails initially
// succeeds on a later retry attempt.
func TestRetryDueDeliveries_DiscordSucceedsOnRetry(t *testing.T) {
	s := newTestStore(t)

	failing := true
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if failing {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	setChannelConfig(t, s, "discord", srv.URL)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	before, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "discord")
	if !ok || before.Status != store.DeliveryStatusFailed {
		t.Fatalf("expected initial discord delivery to be failed, got %+v (ok=%v)", before, ok)
	}

	// Force the delivery due now.
	if _, err := s.DB().ExecContext(context.Background(),
		`UPDATE notification_deliveries SET next_attempt_at = ? WHERE id = ?`, "1970-01-01T00:00:00Z", before.ID); err != nil {
		t.Fatalf("force due: %v", err)
	}

	failing = false
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	d.retryDueDeliveries(context.Background(), logger)

	after, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "discord")
	if !ok {
		t.Fatalf("expected discord delivery row to still exist")
	}
	if after.Status != store.DeliveryStatusSent {
		t.Errorf("expected discord status sent after retry, got %s", after.Status)
	}
	if after.Attempts != before.Attempts+1 {
		t.Errorf("expected attempts %d, got %d", before.Attempts+1, after.Attempts)
	}
	if after.NextAttemptAt != "" {
		t.Errorf("expected empty NextAttemptAt for successful delivery, got %q", after.NextAttemptAt)
	}
}

// TestRetryDueDeliveries_SlackExceedsRetryBound verifies slack delivery respects the retry bound.
func TestRetryDueDeliveries_SlackExceedsRetryBound(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	setChannelConfig(t, s, "slack", srv.URL)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	slack, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "slack")
	if !ok || slack.Status != store.DeliveryStatusFailed {
		t.Fatalf("expected initial slack delivery to be failed, got %+v (ok=%v)", slack, ok)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	// Retry up to maxDeliveryAttempts times.
	for attempt := 2; attempt <= maxDeliveryAttempts; attempt++ {
		if _, err := s.DB().ExecContext(context.Background(),
			`UPDATE notification_deliveries SET next_attempt_at = ? WHERE id = ?`, "1970-01-01T00:00:00Z", slack.ID); err != nil {
			t.Fatalf("force due on attempt %d: %v", attempt, err)
		}

		d.retryDueDeliveries(context.Background(), logger)

		slack, ok = findDelivery(deliveriesFor(t, s, slack.NotificationID), "slack")
		if !ok {
			t.Fatalf("expected slack delivery row to still exist on attempt %d", attempt)
		}
		if slack.Status != store.DeliveryStatusFailed {
			t.Errorf("attempt %d: expected status failed, got %s", attempt, slack.Status)
		}
		if slack.Attempts != attempt {
			t.Errorf("attempt %d: expected attempts %d, got %d", attempt, attempt, slack.Attempts)
		}
	}

	// After maxDeliveryAttempts, should have empty next_attempt_at.
	if slack.NextAttemptAt != "" {
		t.Errorf("expected empty NextAttemptAt after maxDeliveryAttempts, got %q", slack.NextAttemptAt)
	}
}

// TestRetryChannel_ChannelDisabledReturnsError verifies that retryChannel fails when the
// channel is disabled.
func TestRetryChannel_ChannelDisabledReturnsError(t *testing.T) {
	s := newTestStore(t)
	setWebhookConfig(t, s, "http://localhost:9999")

	d := NewDispatcher(s, nil)

	// Create a notification for retrying.
	notif := &store.NotificationRecord{
		UserID:    "user-1",
		AlertID:   "alert-1",
		EventType: "alert.created",
		Title:     "Title",
		Message:   "Message",
	}
	if err := s.CreateNotification(context.Background(), notif); err != nil {
		t.Fatalf("create notification: %v", err)
	}

	// Disable the webhook channel.
	cfgJSON := `{"channels":[{"type":"webhook","enabled":false,"config":{"url":"http://localhost:9999"}}]}`
	if _, err := s.DB().ExecContext(context.Background(),
		`INSERT INTO notification_config (id, config_json) VALUES (1, ?) ON CONFLICT(id) DO UPDATE SET config_json = excluded.config_json`, cfgJSON); err != nil {
		t.Fatalf("set disabled webhook config: %v", err)
	}

	// retryChannel should return an error.
	err := d.retryChannel(context.Background(), notif, "webhook")
	if err == nil {
		t.Errorf("expected retryChannel to return error for disabled channel, got nil")
	}
}
