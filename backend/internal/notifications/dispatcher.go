// Package notifications provides notification dispatching (in-app, SMTP stub, webhook, Discord,
// Slack) with per-channel delivery tracking and bounded retry for failed deliveries.
package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/WiseLabz/wiselabz/internal/crypto"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// retrySchedule is the backoff delay before each retry attempt (index 0 = delay before the
// 2nd attempt, etc).
// ponytail: fixed schedule, not exponential-from-config; add jitter/config if a real deployment needs it.
var retrySchedule = []time.Duration{time.Minute, 5 * time.Minute, 15 * time.Minute, 30 * time.Minute, time.Hour}

// maxDeliveryAttempts caps retries; once exhausted a failed delivery stays failed with no next_attempt_at.
const maxDeliveryAttempts = 5 // len(retrySchedule)

const maxConcurrentNotifications = 8

// Dispatcher routes alert events to notification channels based on config.
type Dispatcher struct {
	store     *store.Store
	hub       *ws.Hub
	fanoutSem chan struct{}
	encKey    []byte // decodes per-channel signing secrets; nil disables signing
}

// NewDispatcher creates a new notification dispatcher.
func NewDispatcher(s *store.Store, hub *ws.Hub) *Dispatcher {
	return &Dispatcher{store: s, hub: hub, fanoutSem: make(chan struct{}, maxConcurrentNotifications)}
}

// SetEncryptionKey installs the base64 key used to decrypt per-channel webhook signing
// secrets. An invalid key leaves signing disabled and is logged.
func (d *Dispatcher) SetEncryptionKey(b64 string) {
	key, err := crypto.DecodeKey(b64)
	if err != nil {
		slog.Error("notification signing disabled: invalid encryption key", "error", err)
		return
	}
	d.encKey = key
}

// maxNotifyUsers bounds the single-page user fetch in NotifyAlertCreated.
// ponytail: fine for a self-hosted ops tool's user count; paginate if that changes.
const maxNotifyUsers = 1000

// NotifyAlertCreated dispatches a newly created alert to every active user.
func (d *Dispatcher) NotifyAlertCreated(ctx context.Context, alertID, title, message string) {
	users, _, err := d.store.ListUsers(ctx, 0, maxNotifyUsers)
	if err != nil {
		slog.Error("failed to list users for alert notification", "error", err)
		return
	}
	channels := d.loadChannels(ctx)
	routes := d.loadRouting(ctx)
	severity := ""
	connectorID := ""
	if alertID != "" {
		if alert, err := d.store.GetAlert(ctx, alertID); err != nil {
			slog.Error("failed to fetch alert for notification", "error", err, "alertID", alertID)
		} else {
			severity = alert.Severity
			connectorID = alert.ServiceID
		}
	}
	go d.notifyAlertCreated(users, channels, routes, alertID, severity, connectorID, title, message)
}

// NotifyAlertsCreated reuses users, channels, and routing for a committed sync batch.
// Alert metadata is supplied by the caller, avoiding a lookup for every alert.
func (d *Dispatcher) NotifyAlertsCreated(ctx context.Context, alerts []store.AlertRecord) {
	if len(alerts) == 0 {
		return
	}
	users, _, err := d.store.ListUsers(ctx, 0, maxNotifyUsers)
	if err != nil {
		slog.Error("failed to list users for alert notifications", "error", err)
		return
	}
	channels := d.loadChannels(ctx)
	routes := d.loadRouting(ctx)
	batch := append([]store.AlertRecord(nil), alerts...)
	go func() {
		for _, alert := range batch {
			d.notifyAlertCreated(users, channels, routes, alert.ID, alert.Severity, alert.ServiceID, alert.Title, alert.Description)
		}
	}()
}

func (d *Dispatcher) notifyAlertCreated(users []store.User, channels []channelCfg, routes []routeCfg, alertID, severity, connectorID, title, message string) {
	for _, u := range users {
		if u.Disabled {
			continue
		}
		d.fanoutSem <- struct{}{}
		go func(userID string, skipExternal bool) {
			defer func() { <-d.fanoutSem }()
			d.notifyAlert(context.Background(), channels, routes, alertID, userID, "alert.created", severity, connectorID, title, message, skipExternal)
		}(u.ID, u.DigestCadence != "off")
	}
}

// NotifyFindingCreated dispatches a quality finding notification (a newly
// opened finding, or one whose severity just escalated — the caller, quality
// Checker, decides when that's the case) to every active user, generic
// across every finding check_type.
func (d *Dispatcher) NotifyFindingCreated(ctx context.Context, findingID, title, message string) {
	users, _, err := d.store.ListUsers(ctx, 0, maxNotifyUsers)
	if err != nil {
		slog.Error("failed to list users for finding notification", "error", err, "findingID", findingID)
		return
	}
	channels := d.loadChannels(ctx)
	routes := d.loadRouting(ctx)
	severity := ""
	connectorID := ""
	if findingID != "" {
		if finding, err := d.store.GetQualityFinding(ctx, findingID); err != nil {
			slog.Error("failed to fetch finding for notification", "error", err, "findingID", findingID)
		} else {
			severity = finding.Severity
			connectorID = finding.ConnectorID
		}
	}
	go d.notifyFindingCreated(users, channels, routes, severity, connectorID, title, message)
}

// notifyFindingCreated fans out like notifyAlertCreated. It reuses notifyAlert
// with an empty alertID: a finding isn't an alert, and ponytail: the in-app
// notification/WS payload has no finding deep-link yet — add one (a real
// findingId column) if the UI needs to navigate straight to it.
func (d *Dispatcher) notifyFindingCreated(users []store.User, channels []channelCfg, routes []routeCfg, severity, connectorID, title, message string) {
	for _, u := range users {
		if u.Disabled {
			continue
		}
		d.fanoutSem <- struct{}{}
		go func(userID string, skipExternal bool) {
			defer func() { <-d.fanoutSem }()
			d.notifyAlert(context.Background(), channels, routes, "", userID, "finding.created", severity, connectorID, title, message, skipExternal)
		}(u.ID, u.DigestCadence != "off")
	}
}

// channelCfg is one entry of the notification_config "channels" array.
type channelCfg struct {
	Type    string         `json:"type"`
	Enabled bool           `json:"enabled"`
	Config  map[string]any `json:"config"`
}

// routeCfg is one entry of the notification_config "routing" array.
type routeCfg struct {
	EventType         string `json:"eventType"`
	Channel           string `json:"channel"`
	Enabled           bool   `json:"enabled"`
	MinSeverity       string `json:"minSeverity"`
	ConnectorCategory string `json:"connectorCategory"`
	ConnectorID       string `json:"connectorId"`
}

// loadChannels reads the configured notification channels. Returns nil (no channels) on any
// read/parse error so callers just skip optional channels — in-app delivery never depends on this.
func (d *Dispatcher) loadChannels(ctx context.Context) []channelCfg {
	var raw string
	if err := d.store.DB().QueryRowContext(ctx, `SELECT config_json FROM notification_config WHERE id = 1`).Scan(&raw); err != nil || raw == "" {
		return nil
	}
	var doc struct {
		Channels []channelCfg `json:"channels"`
	}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil
	}
	return doc.Channels
}

// loadRouting reads the configured notification routing rules. Returns nil on any
// read/parse error so callers just skip optional routing gates — in-app delivery never depends on this.
func (d *Dispatcher) loadRouting(ctx context.Context) []routeCfg {
	var raw string
	if err := d.store.DB().QueryRowContext(ctx, `SELECT config_json FROM notification_config WHERE id = 1`).Scan(&raw); err != nil || raw == "" {
		return nil
	}
	var doc struct {
		Routing []routeCfg `json:"routing"`
	}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil
	}
	return doc.Routing
}

// channel returns the config for the given channel type and whether it's enabled.
func (d *Dispatcher) channel(ctx context.Context, typ string) (channelCfg, bool) {
	for _, c := range d.loadChannels(ctx) {
		if c.Type == typ {
			return c, c.Enabled
		}
	}
	return channelCfg{}, false
}

// NotifyAlert dispatches an alert to all configured channels, recording a delivery row per
// channel attempted (see store.DeliveryRecord).
func (d *Dispatcher) NotifyAlert(alertID, userID, eventType, title, message string) {
	ctx := context.Background()
	d.notifyAlert(ctx, d.loadChannels(ctx), d.loadRouting(ctx), alertID, userID, eventType, "", "", title, message, false)
}

func (d *Dispatcher) notifyAlert(ctx context.Context, channels []channelCfg, routes []routeCfg, alertID, userID, eventType, severity, connectorID, title, message string, skipExternal bool) {
	notifID, ok := d.sendInApp(ctx, userID, alertID, eventType, title, message)
	if !ok {
		return
	}
	if skipExternal {
		return
	}
	d.notifyExternalChannels(ctx, notifID, channels, routes, eventType, severity, connectorID, userID, title, message)
}

func (d *Dispatcher) sendInApp(ctx context.Context, userID, alertID, eventType, title, message string) (string, bool) {
	notifID, err := d.createInApp(ctx, userID, alertID, eventType, title, message)
	if err != nil {
		slog.Error("failed to create in-app notification", "error", err)
		return "", false
	}
	d.recordDelivery(ctx, notifID, "in_app", store.DeliveryStatusSent, "")
	if d.hub != nil {
		d.hub.BroadcastToUser(userID, eventType, map[string]any{
			"alertId": alertID,
			"title":   title,
			"message": message,
		})
	}
	return notifID, true
}

func (d *Dispatcher) notifyExternalChannels(ctx context.Context, notifID string, channels []channelCfg, routes []routeCfg, eventType, severity, connectorID, userID, title, message string) {
	connectorCategory := d.notificationConnectorCategory(ctx, connectorID)
	shouldSkip := func(channel string) bool {
		return shouldSkipRoute(routes, eventType, channel, severity, connectorCategory, connectorID)
	}
	if !shouldSkip("smtp") {
		if _, enabled := findChannel(channels, "smtp"); enabled {
			// ponytail: SMTP delivery is a stub — always "succeeds" and only logs. Real sending
			// (user email lookup, SMTP auth/TLS, credential decryption) is out of scope for issue #18;
			// still recorded as a real delivery row so retry/observability plumbing already covers it
			// once real sending lands.
			slog.Info("SMTP notification (stub)", "userID", userID, "title", title)
			d.recordDelivery(ctx, notifID, "smtp", store.DeliveryStatusSent, "")
		}
	}
	for _, delivery := range []struct {
		channel string
		payload any
	}{{"webhook", webhookPayload(title, message)}, {"discord", discordPayload(title, message)}, {"slack", slackPayload(title, message)}} {
		if !shouldSkip(delivery.channel) {
			if cfg, enabled := findChannel(channels, delivery.channel); enabled {
				d.attemptChannel(ctx, notifID, delivery.channel, cfg, delivery.payload)
			}
		}
	}
}

func (d *Dispatcher) notificationConnectorCategory(ctx context.Context, connectorID string) string {
	if connectorID != "" {
		if connector, err := d.store.GetConnector(ctx, connectorID); err != nil {
			slog.Error("failed to fetch connector for notification", "error", err, "connectorID", logsafe.Sanitize(connectorID))
			// Treat as unknown category; routes with a category filter just won't match.
		} else {
			return connector.Category
		}
	}
	return ""
}

func shouldSkipRoute(routes []routeCfg, eventType, channel, severity, connectorCategory, connectorID string) bool {
	if len(routes) == 0 {
		// No routing configured; proceed with normal channel-based delivery.
		return false
	}
	route, foundRoute := findRoute(routes, eventType, channel)
	if !foundRoute || !route.Enabled || severityRank(severity) > severityRank(route.MinSeverity) {
		return true
	}
	return (route.ConnectorCategory != "" && route.ConnectorCategory != connectorCategory) ||
		(route.ConnectorID != "" && route.ConnectorID != connectorID)
}

func findChannel(channels []channelCfg, typ string) (channelCfg, bool) {
	for _, c := range channels {
		if c.Type == typ {
			return c, c.Enabled
		}
	}
	return channelCfg{}, false
}

// findRoute returns the first routing rule matching both eventType and channel, and whether one was found.
func findRoute(routes []routeCfg, eventType, channel string) (routeCfg, bool) {
	for _, r := range routes {
		if r.EventType == eventType && r.Channel == channel {
			return r, true
		}
	}
	return routeCfg{}, false
}

// severityRank returns a numeric rank for routing: critical=0, warning=1, info=2, unknown=3.
// Lower rank = more severe. This must match the severity hierarchy elsewhere in the codebase.
func severityRank(severity string) int {
	switch severity {
	case "critical":
		return 0
	case "warning":
		return 1
	case "info":
		return 2
	default:
		return 3
	}
}

// createInApp creates the in-app notification row and returns its ID.
func (d *Dispatcher) createInApp(ctx context.Context, userID, alertID, eventType, title, message string) (string, error) {
	n := &store.NotificationRecord{
		UserID:    userID,
		AlertID:   alertID,
		EventType: eventType,
		Title:     title,
		Message:   message,
	}
	if err := d.store.CreateNotification(ctx, n); err != nil {
		return "", err
	}
	return n.ID, nil
}

// recordDelivery persists the outcome of a first delivery attempt for one channel.
func (d *Dispatcher) recordDelivery(ctx context.Context, notificationID, channel string, status store.DeliveryStatus, lastError string) {
	rec := &store.DeliveryRecord{
		NotificationID: notificationID,
		Channel:        channel,
		Status:         status,
		Attempts:       1,
		LastError:      lastError,
	}
	if status == store.DeliveryStatusFailed {
		rec.NextAttemptAt = time.Now().UTC().Add(retrySchedule[0]).Format(time.RFC3339)
	}
	if err := d.store.CreateDelivery(ctx, rec); err != nil {
		slog.Error("record delivery", "error", err, "channel", channel)
	}
}

// attemptChannel sends payload to the channel's configured URL over the shared webhook transport
// and records the resulting delivery status under channelType.
func (d *Dispatcher) attemptChannel(ctx context.Context, notificationID, channelType string, cfg channelCfg, payload any) {
	url, _ := cfg.Config["url"].(string)
	if url == "" {
		d.recordDelivery(ctx, notificationID, channelType, store.DeliveryStatusFailed, channelType+" url not configured")
		return
	}
	if err := sendWebhook(ctx, url, d.signingSecret(cfg), payload); err != nil {
		d.recordDelivery(ctx, notificationID, channelType, store.DeliveryStatusFailed, err.Error())
		return
	}
	d.recordDelivery(ctx, notificationID, channelType, store.DeliveryStatusSent, "")
}

// webhookPayload shapes a title/message pair into the generic webhook body.
func webhookPayload(title, message string) any {
	return map[string]string{"title": title, "message": message}
}

// discordPayload shapes a title/message pair into a Discord incoming-webhook body.
func discordPayload(title, message string) any {
	return map[string]string{"content": fmt.Sprintf("**%s**\n%s", title, message)}
}

// slackPayload shapes a title/message pair into a Slack incoming-webhook body.
func slackPayload(title, message string) any {
	return map[string]string{"text": fmt.Sprintf("*%s*\n%s", title, message)}
}
