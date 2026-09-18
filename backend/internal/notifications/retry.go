package notifications

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

const retryTick = 30 * time.Second

// dueDeliveriesLimit bounds how many due deliveries are retried per tick.
// ponytail: fine for a self-hosted ops tool's notification volume; paginate if that changes.
const dueDeliveriesLimit = 50

// RunDeliveryRetries polls for failed deliveries whose next_attempt_at is due and retries them,
// modeled on runAlertExpirer in cmd/server/main.go.
func RunDeliveryRetries(ctx context.Context, d *Dispatcher, logger *slog.Logger) {
	logger.Info("Notification delivery retrier started")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d.retryDueDeliveries(ctx, logger)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(retryTick):
		}
	}
}

// retryDueDeliveries retries every delivery whose next_attempt_at has passed.
func (d *Dispatcher) retryDueDeliveries(ctx context.Context, logger *slog.Logger) {
	now := time.Now().UTC().Format(time.RFC3339)
	due, err := d.store.ListDueDeliveries(ctx, now, dueDeliveriesLimit)
	if err != nil {
		logger.Error("list due deliveries", "error", err)
		return
	}
	for _, del := range due {
		notif, err := d.store.GetNotificationByID(ctx, del.NotificationID)
		if err != nil {
			logger.Error("get notification for retry", "error", err, "deliveryId", del.ID)
			continue
		}

		var sendErr error
		switch del.Channel {
		case "webhook":
			sendErr = d.retryChannel(ctx, notif, "webhook", webhookPayload)
		case "discord":
			sendErr = d.retryChannel(ctx, notif, "discord", discordPayload)
		case "slack":
			sendErr = d.retryChannel(ctx, notif, "slack", slackPayload)
		default:
			// ponytail: only webhook/discord/slack are retried today; in_app never fails and
			// smtp is a stub that always succeeds, so neither lands in the failed+due set.
			continue
		}

		attempts := del.Attempts + 1
		if sendErr == nil {
			if err := d.store.UpdateDeliveryResult(ctx, del.ID, store.DeliveryStatusSent, attempts, "", ""); err != nil {
				logger.Error("update delivery result", "error", err, "deliveryId", del.ID)
			}
			continue
		}

		var next string
		if attempts < maxDeliveryAttempts {
			next = time.Now().UTC().Add(retrySchedule[attempts-1]).Format(time.RFC3339)
		}
		if err := d.store.UpdateDeliveryResult(ctx, del.ID, store.DeliveryStatusFailed, attempts, sendErr.Error(), next); err != nil {
			logger.Error("update delivery result", "error", err, "deliveryId", del.ID)
		}
	}
}

// retryChannel re-sends a delivery for the given channel type using the current channel config.
// If the channel was disabled or removed since the original attempt, retrying fails without a
// network call.
func (d *Dispatcher) retryChannel(ctx context.Context, notif *store.NotificationRecord, channelType string, payloadFn func(title, message string) any) error {
	cfg, enabled := d.channel(ctx, channelType)
	if !enabled {
		return errors.New(channelType + " channel disabled or removed")
	}
	url, _ := cfg.Config["url"].(string)
	if url == "" {
		return errors.New(channelType + " url not configured")
	}
	return sendWebhook(ctx, url, d.signingSecret(cfg), payloadFn(notif.Title, notif.Message))
}
