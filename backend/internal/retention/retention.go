// Package retention runs the background job that bounds the growth of
// historical WiseLabz data (service snapshots, doc revisions, alerts, health
// checks, and sync run history) by deleting rows past a configurable
// per-category age.
package retention

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// Fixed retention windows (days) for tables that have no configurable
// setting.
const (
	staleSessionDays = 30
	notificationDays = 90
	changeDays       = 180
	shareLinkDays    = 30
	chatDays         = 90
)

// RunCleanupOnce performs one cleanup pass: for every category whose *Days
// config value is > 0, deletes rows older than the cutoff. A category with
// Days <= 0 is skipped (retention disabled). Errors in one category are
// logged and do not stop the others from running, but are joined into the
// returned error so the scheduler's job health (#384) reflects a partial
// cleanup pass.
func RunCleanupOnce(ctx context.Context, s *store.Store, cfg store.RetentionSettings, logger *slog.Logger) error {
	cutoff := func(days int) string {
		return time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339)
	}
	var errs []error

	if cfg.SnapshotDays > 0 {
		n, err := s.DeleteOldSnapshots(ctx, cutoff(cfg.SnapshotDays))
		if err != nil {
			logger.Error("delete old snapshots", "error", err)
			errs = append(errs, fmt.Errorf("delete old snapshots: %w", err))
		} else if n > 0 {
			logger.Info("Purged old snapshots", "count", n)
		}
	}

	if cfg.DocVersionDays > 0 {
		n, err := s.DeleteOldDocVersions(ctx, cutoff(cfg.DocVersionDays))
		if err != nil {
			logger.Error("delete old doc versions", "error", err)
			errs = append(errs, fmt.Errorf("delete old doc versions: %w", err))
		} else if n > 0 {
			logger.Info("Purged old doc versions", "count", n)
		}
	}

	if cfg.AlertDays > 0 {
		n, err := s.DeleteOldAlerts(ctx, cutoff(cfg.AlertDays))
		if err != nil {
			logger.Error("delete old alerts", "error", err)
			errs = append(errs, fmt.Errorf("delete old alerts: %w", err))
		} else if n > 0 {
			logger.Info("Purged old alerts", "count", n)
		}
	}

	if cfg.SyncRunDays > 0 {
		n, err := s.DeleteOldSyncRuns(ctx, cutoff(cfg.SyncRunDays))
		if err != nil {
			logger.Error("delete old sync runs", "error", err)
			errs = append(errs, fmt.Errorf("delete old sync runs: %w", err))
		} else if n > 0 {
			logger.Info("Purged old sync runs", "count", n)
		}
	}

	if cfg.AuditDays > 0 {
		n, err := s.DeleteOldAuditRecords(ctx, cutoff(cfg.AuditDays))
		if err != nil {
			logger.Error("delete old audit records", "error", err)
			errs = append(errs, fmt.Errorf("delete old audit records: %w", err))
		} else if n > 0 {
			logger.Info("Purged old audit records", "count", n)
		}
	}

	if cfg.HealthCheckDays > 0 {
		n, err := s.DeleteOldHealthChecks(ctx, cutoff(cfg.HealthCheckDays))
		if err != nil {
			logger.Error("delete old health checks", "error", err)
			errs = append(errs, fmt.Errorf("delete old health checks: %w", err))
		} else if n > 0 {
			logger.Info("Purged old health checks", "count", n)
		}
	}

	if cfg.ReportDays > 0 {
		n, err := s.DeleteOldReports(ctx, cutoff(cfg.ReportDays))
		if err != nil {
			logger.Error("delete old reports", "error", err)
			errs = append(errs, fmt.Errorf("delete old reports: %w", err))
		} else if n > 0 {
			logger.Info("Purged old reports", "count", n)
		}
	}

	// Tables without a configurable *Days setting use fixed windows.
	fixed := []struct {
		name   string
		days   int
		delete func(context.Context, string) (int64, error)
	}{
		{"sessions", staleSessionDays, s.DeleteStaleSessions},
		{"notification deliveries", notificationDays, s.DeleteOldDeliveries},
		{"notifications", notificationDays, s.DeleteOldNotifications},
		{"changes", changeDays, s.DeleteOldChanges},
		{"share links", shareLinkDays, s.DeleteExpiredShareLinks},
		{"chat conversations", chatDays, s.DeleteOldChatConversations},
	}
	for _, f := range fixed {
		n, err := f.delete(ctx, cutoff(f.days))
		if err != nil {
			logger.Error("delete old "+f.name, "error", err)
			errs = append(errs, fmt.Errorf("delete old %s: %w", f.name, err))
		} else if n > 0 {
			logger.Info("Purged old "+f.name, "count", n)
		}
	}

	return errors.Join(errs...)
}
