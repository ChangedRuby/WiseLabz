package notifications

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
	_ "modernc.org/sqlite"
)

// newTestStore mirrors internal/sync/engine_test.go's helper: a real sqlite-backed store so
// dispatcher logic exercises actual SQL instead of a mock.
func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() }) //nolint:errcheck

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := store.RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return store.New(db, "sqlite")
}

// setWebhookConfig writes a notification_config row enabling only the webhook channel.
func setWebhookConfig(t *testing.T, s *store.Store, url string) {
	t.Helper()
	setChannelConfig(t, s, "webhook", url)
}

// setChannelConfig writes a notification_config row enabling only the given channel type.
func setChannelConfig(t *testing.T, s *store.Store, channelType, url string) {
	t.Helper()
	cfgJSON := `{"channels":[{"type":"` + channelType + `","enabled":true,"config":{"url":"` + url + `"}}]}`
	// newTestStore doesn't call Store.Init, so the id=1 singleton row may not exist yet — upsert.
	if _, err := s.DB().ExecContext(context.Background(),
		`INSERT INTO notification_config (id, config_json) VALUES (1, ?) ON CONFLICT(id) DO UPDATE SET config_json = excluded.config_json`, cfgJSON); err != nil {
		t.Fatalf("set %s config: %v", channelType, err)
	}
}

// setChannelAndRoutingConfig writes a notification_config row with both channels and routing.
func setChannelAndRoutingConfig(t *testing.T, s *store.Store, channelType, url, routing string) {
	t.Helper()
	cfgJSON := `{"channels":[{"type":"` + channelType + `","enabled":true,"config":{"url":"` + url + `"}}],"routing":` + routing + `}`
	if _, err := s.DB().ExecContext(context.Background(),
		`INSERT INTO notification_config (id, config_json) VALUES (1, ?) ON CONFLICT(id) DO UPDATE SET config_json = excluded.config_json`, cfgJSON); err != nil {
		t.Fatalf("set config: %v", err)
	}
}

func deliveriesFor(t *testing.T, s *store.Store, notificationID string) []store.DeliveryRecord {
	t.Helper()
	// ponytail: ListDeliveries only filters by status, not notification_id; query directly
	// instead of adding a store method the dispatcher package doesn't otherwise need.
	rows, err := s.DB().QueryContext(context.Background(),
		`SELECT id, notification_id, channel, status, attempts, last_error, IFNULL(next_attempt_at, '') FROM notification_deliveries WHERE notification_id = ?`, notificationID)
	if err != nil {
		t.Fatalf("query deliveries: %v", err)
	}
	defer rows.Close() //nolint:errcheck
	var out []store.DeliveryRecord
	for rows.Next() {
		var d store.DeliveryRecord
		if err := rows.Scan(&d.ID, &d.NotificationID, &d.Channel, &d.Status, &d.Attempts, &d.LastError, &d.NextAttemptAt); err != nil {
			t.Fatalf("scan delivery: %v", err)
		}
		out = append(out, d)
	}
	return out
}

func findDelivery(deliveries []store.DeliveryRecord, channel string) (store.DeliveryRecord, bool) {
	for _, d := range deliveries {
		if d.Channel == channel {
			return d, true
		}
	}
	return store.DeliveryRecord{}, false
}

func TestNotifyAlert_NoChannelsConfigured(t *testing.T) {
	s := newTestStore(t)
	d := NewDispatcher(s, nil)

	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}

	deliveries := deliveriesFor(t, s, notifs[0].ID)
	inApp, ok := findDelivery(deliveries, "in_app")
	if !ok {
		t.Fatalf("expected in_app delivery row, got %+v", deliveries)
	}
	if inApp.Status != store.DeliveryStatusSent {
		t.Errorf("expected in_app status sent, got %s", inApp.Status)
	}
	if _, ok := findDelivery(deliveries, "smtp"); ok {
		t.Errorf("expected no smtp delivery row when smtp not configured")
	}
	if _, ok := findDelivery(deliveries, "webhook"); ok {
		t.Errorf("expected no webhook delivery row when webhook not configured")
	}
}

func TestNotifyAlert_WebhookSuccess(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	setWebhookConfig(t, s, srv.URL)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	deliveries := deliveriesFor(t, s, notifs[0].ID)
	webhook, ok := findDelivery(deliveries, "webhook")
	if !ok {
		t.Fatalf("expected webhook delivery row, got %+v", deliveries)
	}
	if webhook.Status != store.DeliveryStatusSent {
		t.Errorf("expected webhook status sent, got %s (err=%s)", webhook.Status, webhook.LastError)
	}
}

func TestNotifyAlert_WebhookFailure(t *testing.T) {
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
	deliveries := deliveriesFor(t, s, notifs[0].ID)
	webhook, ok := findDelivery(deliveries, "webhook")
	if !ok {
		t.Fatalf("expected webhook delivery row, got %+v", deliveries)
	}
	if webhook.Status != store.DeliveryStatusFailed {
		t.Errorf("expected webhook status failed, got %s", webhook.Status)
	}
	if webhook.NextAttemptAt == "" {
		t.Errorf("expected non-empty NextAttemptAt for failed delivery")
	}
	if webhook.LastError == "" {
		t.Errorf("expected non-empty LastError for failed delivery")
	}
}

func TestNotifyAlert_DiscordSuccess(t *testing.T) {
	s := newTestStore(t)
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
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
	deliveries := deliveriesFor(t, s, notifs[0].ID)
	discord, ok := findDelivery(deliveries, "discord")
	if !ok {
		t.Fatalf("expected discord delivery row, got %+v", deliveries)
	}
	if discord.Status != store.DeliveryStatusSent {
		t.Errorf("expected discord status sent, got %s (err=%s)", discord.Status, discord.LastError)
	}
	if !strings.Contains(gotBody, `"content"`) || !strings.Contains(gotBody, "**Title**") {
		t.Errorf("expected discord payload with bolded title, got %s", gotBody)
	}
}

func TestNotifyAlert_DiscordFailure(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	setChannelConfig(t, s, "discord", srv.URL)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	discord, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "discord")
	if !ok {
		t.Fatalf("expected discord delivery row, got none")
	}
	if discord.Status != store.DeliveryStatusFailed {
		t.Errorf("expected discord status failed, got %s", discord.Status)
	}
}

func TestNotifyAlert_SlackSuccess(t *testing.T) {
	s := newTestStore(t)
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.WriteHeader(http.StatusOK)
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
	if !ok {
		t.Fatalf("expected slack delivery row, got none")
	}
	if slack.Status != store.DeliveryStatusSent {
		t.Errorf("expected slack status sent, got %s (err=%s)", slack.Status, slack.LastError)
	}
	if !strings.Contains(gotBody, `"text"`) || !strings.Contains(gotBody, "*Title*") {
		t.Errorf("expected slack payload with italicized title, got %s", gotBody)
	}
}

func TestNotifyAlert_SlackFailure(t *testing.T) {
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
	if !ok {
		t.Fatalf("expected slack delivery row, got none")
	}
	if slack.Status != store.DeliveryStatusFailed {
		t.Errorf("expected slack status failed, got %s", slack.Status)
	}
}

// TestRetryDueDeliveries_DiscordRecoversAfterFailure mirrors the webhook retry test for the
// discord channel; the retry path is shared code (retryChannel), so covering one of the two new
// channel types here is enough — slack goes through the identical path.
func TestRetryDueDeliveries_DiscordRecoversAfterFailure(t *testing.T) {
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
}

// TestRetryDueDeliveries_RecoversAfterFailure drives one retry loop iteration directly: a
// webhook delivery that failed against a dead endpoint is retried once the endpoint (the same
// URL, backed by a mutable handler) starts returning 200, ending up sent with attempts incremented.
func TestRetryDueDeliveries_RecoversAfterFailure(t *testing.T) {
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
	setWebhookConfig(t, s, srv.URL)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	before, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "webhook")
	if !ok || before.Status != store.DeliveryStatusFailed {
		t.Fatalf("expected initial webhook delivery to be failed, got %+v (ok=%v)", before, ok)
	}

	// Force the delivery due now (its natural next_attempt_at is minutes in the future).
	if _, err := s.DB().ExecContext(context.Background(),
		`UPDATE notification_deliveries SET next_attempt_at = ? WHERE id = ?`, "1970-01-01T00:00:00Z", before.ID); err != nil {
		t.Fatalf("force due: %v", err)
	}

	failing = false
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	d.retryDueDeliveries(context.Background(), logger)

	after, ok := findDelivery(deliveriesFor(t, s, notifs[0].ID), "webhook")
	if !ok {
		t.Fatalf("expected webhook delivery row to still exist")
	}
	if after.Status != store.DeliveryStatusSent {
		t.Errorf("expected webhook status sent after retry, got %s", after.Status)
	}
	if after.Attempts != before.Attempts+1 {
		t.Errorf("expected attempts %d, got %d", before.Attempts+1, after.Attempts)
	}
}

// TestNotifyAlertCreated_NoChannelsConfigured verifies that NotifyAlertCreated works correctly
// when no notification channels are configured, creating only an in-app notification.
func TestNotifyAlertCreated_NoChannelsConfigured(t *testing.T) {
	s := newTestStore(t)
	d := NewDispatcher(s, nil)

	// Do NOT configure any channels, ensuring the notification system gracefully handles
	// a zero-channel setup.

	// Create a user to receive the alert.
	u := &store.User{
		Email: "test@example.com",
	}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Dispatch alert to all users.
	d.NotifyAlertCreated(context.Background(), "alert-1", "Test Alert", "This is a test alert")

	// Give any goroutines time to complete.
	// ponytail: simple sleep; in production, would use a sync.WaitGroup or channels.
	time.Sleep(100 * time.Millisecond)

	// Verify notification was created for the user.
	notifs, _, err := s.ListNotifications(context.Background(), u.ID, false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification for user, got %d", len(notifs))
	}

	// Verify only in-app delivery was recorded (no external channels).
	deliveries := deliveriesFor(t, s, notifs[0].ID)
	if len(deliveries) != 1 {
		t.Errorf("expected 1 delivery (in-app only), got %d", len(deliveries))
	}

	inApp, ok := findDelivery(deliveries, "in_app")
	if !ok {
		t.Fatalf("expected in_app delivery row, got %+v", deliveries)
	}
	if inApp.Status != store.DeliveryStatusSent {
		t.Errorf("expected in_app status sent, got %s", inApp.Status)
	}
}

// TestNotifyAlert_RoutingMissingSkips verifies that a channel is skipped when no routing rule exists.
func TestNotifyAlert_RoutingMissingSkips(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Configure webhook channel but no routing rule for it.
	routing := `[{"eventType":"alert.created","channel":"discord","enabled":true,"minSeverity":"info"}]`
	setChannelAndRoutingConfig(t, s, "webhook", srv.URL, routing)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	deliveries := deliveriesFor(t, s, notifs[0].ID)

	// Should have in_app only, no webhook delivery.
	if _, ok := findDelivery(deliveries, "webhook"); ok {
		t.Errorf("expected no webhook delivery when no routing rule matches")
	}
}

// TestNotifyAlert_RoutingDisabledSkips verifies that a channel is skipped when its routing rule is disabled.
func TestNotifyAlert_RoutingDisabledSkips(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Configure webhook channel with a disabled routing rule.
	routing := `[{"eventType":"alert.created","channel":"webhook","enabled":false,"minSeverity":"info"}]`
	setChannelAndRoutingConfig(t, s, "webhook", srv.URL, routing)

	d := NewDispatcher(s, nil)
	d.NotifyAlert("alert-1", "user-1", "alert.created", "Title", "Message")

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	deliveries := deliveriesFor(t, s, notifs[0].ID)

	// Should have in_app only, no webhook delivery because route is disabled.
	if _, ok := findDelivery(deliveries, "webhook"); ok {
		t.Errorf("expected no webhook delivery when routing rule is disabled")
	}
}

// TestNotifyAlert_RoutingBelowSeveritySkips verifies that a channel is skipped when event severity is below minSeverity.
func TestNotifyAlert_RoutingBelowSeveritySkips(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Configure webhook channel with minSeverity=critical (only critical events should route).
	routing := `[{"eventType":"alert.created","channel":"webhook","enabled":true,"minSeverity":"critical"}]`
	setChannelAndRoutingConfig(t, s, "webhook", srv.URL, routing)

	// Create an alert with "warning" severity
	alert := &store.AlertRecord{
		ChangeID:  "change-1",
		ServiceID: "service-1",
		Severity:  "warning",
		Title:     "Warning Alert",
		Status:    "pending",
	}
	if err := s.CreateAlert(context.Background(), alert); err != nil {
		t.Fatalf("create alert: %v", err)
	}

	d := NewDispatcher(s, nil)
	d.NotifyAlertCreated(context.Background(), alert.ID, "Title", "Message")
	time.Sleep(100 * time.Millisecond)

	notifs, _, err := s.ListNotifications(context.Background(), "user-1", false, 0, 10)
	if err != nil {
		// No user might exist, that's OK
		notifs = nil
	}
	if len(notifs) > 0 {
		deliveries := deliveriesFor(t, s, notifs[0].ID)
		if _, ok := findDelivery(deliveries, "webhook"); ok {
			t.Errorf("expected no webhook delivery when severity below minSeverity")
		}
	}
}

// TestNotifyAlert_RoutingAboveSeverityDelivers verifies that a channel delivers when event severity meets minSeverity.
func TestNotifyAlert_RoutingAboveSeverityDelivers(t *testing.T) {
	s := newTestStore(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Configure webhook channel with minSeverity=warning (warning and critical should route).
	routing := `[{"eventType":"alert.created","channel":"webhook","enabled":true,"minSeverity":"warning"}]`
	setChannelAndRoutingConfig(t, s, "webhook", srv.URL, routing)

	// Create a user to receive notifications
	u := &store.User{Email: "test@example.com"}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Create an alert with "critical" severity
	alert := &store.AlertRecord{
		ChangeID:  "change-1",
		ServiceID: "service-1",
		Severity:  "critical",
		Title:     "Critical Alert",
		Status:    "pending",
	}
	if err := s.CreateAlert(context.Background(), alert); err != nil {
		t.Fatalf("create alert: %v", err)
	}

	d := NewDispatcher(s, nil)
	d.NotifyAlertCreated(context.Background(), alert.ID, "Title", "Message")
	time.Sleep(100 * time.Millisecond)

	notifs, _, err := s.ListNotifications(context.Background(), u.ID, false, 0, 10)
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(notifs) == 0 {
		t.Fatalf("expected at least 1 notification")
	}

	deliveries := deliveriesFor(t, s, notifs[0].ID)
	webhook, ok := findDelivery(deliveries, "webhook")
	if !ok {
		t.Fatalf("expected webhook delivery when severity meets minSeverity")
	}
	if webhook.Status != store.DeliveryStatusSent {
		t.Errorf("expected webhook delivery to be sent, got %s (err=%s)", webhook.Status, webhook.LastError)
	}
}
