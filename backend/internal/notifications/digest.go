package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// digestSendHour is the local hour (in each user's digest_timezone) at which
// digest sweeps become eligible to send. The sweep itself runs hourly; this
// narrows delivery to once per eligible cadence window per user.
const digestSendHour = 8

// digestDue is a pure function (no DB, no real clock) so digest eligibility
// is unit-testable in isolation. now and the parsed lastSentAt are compared
// as elapsed time, not calendar day-of-week, so a late or restarted sweep
// still catches up correctly instead of silently missing a week.
func digestDue(cadence, lastSentAt string, localHour int, now time.Time) bool {
	if cadence == "off" {
		return false
	}
	if localHour != digestSendHour {
		return false
	}
	if cadence == "daily" {
		return true
	}
	if cadence == "weekly" {
		if lastSentAt == "" {
			return true
		}
		last, err := time.Parse(time.RFC3339, lastSentAt)
		if err != nil {
			return true
		}
		return now.Sub(last) >= 7*24*time.Hour
	}
	return false
}

// RunDigestSweep runs one pass of the hourly digest job: for every user with
// digest_cadence != "off" whose local time is in the send window, it collects
// notifications accumulated since their last digest, sends a single summary
// through their configured channels, and advances digest_last_sent_at.
func (d *Dispatcher) RunDigestSweep(ctx context.Context, now time.Time, logger *slog.Logger) {
	users, _, err := d.store.ListUsers(ctx, 0, 10000)
	if err != nil {
		logger.Error("digest sweep: failed to list users", "error", err)
		return
	}

	channels := d.loadChannels(ctx)
	routes := d.loadRouting(ctx)

	for _, u := range users {
		if u.Disabled || u.DigestCadence == "off" {
			continue
		}
		loc, err := time.LoadLocation(u.DigestTimezone)
		if err != nil {
			logger.Error("digest sweep: invalid timezone", "userID", u.ID, "timezone", u.DigestTimezone, "error", err)
			continue
		}
		localHour := now.In(loc).Hour()
		if !digestDue(u.DigestCadence, u.DigestLastSentAt, localHour, now) {
			continue
		}

		notifications, err := d.store.ListNotificationsSince(ctx, u.ID, u.DigestLastSentAt, []string{"alert.created", "finding.created", EventSystemJobFailed})
		if err != nil {
			logger.Error("digest sweep: failed to list notifications", "userID", u.ID, "error", err)
			continue
		}

		nowStr := now.Format(time.RFC3339)
		if len(notifications) == 0 {
			if err := d.store.UpdateUser(ctx, u.ID, map[string]any{"digest_last_sent_at": nowStr}); err != nil {
				logger.Error("digest sweep: failed to advance watermark", "userID", u.ID, "error", err)
			}
			continue
		}

		title, message := formatDigest(notifications)
		d.notifyAlert(ctx, channels, routes, "", u.ID, "digest.summary", "", "", title, message, false)

		if err := d.store.UpdateUser(ctx, u.ID, map[string]any{"digest_last_sent_at": nowStr}); err != nil {
			logger.Error("digest sweep: failed to advance watermark", "userID", u.ID, "error", err)
		}
	}
}

// formatDigest builds a plain-text summary: count by severity/event type,
// then up to N individual titles. No template engine — the repo has none,
// and this is short enough that string building is simpler.
func formatDigest(notifications []store.NotificationRecord) (title, message string) {
	const maxTitles = 10
	title = fmt.Sprintf("Digest: %d new notification(s)", len(notifications))

	var b strings.Builder
	fmt.Fprintf(&b, "%d new notification(s) since your last digest:\n\n", len(notifications))
	for i, n := range notifications {
		if i >= maxTitles {
			fmt.Fprintf(&b, "…and %d more\n", len(notifications)-maxTitles)
			break
		}
		fmt.Fprintf(&b, "- %s\n", n.Title)
	}
	return title, b.String()
}
