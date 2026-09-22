package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// HealthCheckRecord represents a row in the health_checks table: one
// connector/health.go check result, persisted as a time-series point rather
// than overwriting the connector's latest status.
type HealthCheckRecord struct {
	ID          string `json:"id"`
	ConnectorID string `json:"connectorId"`
	Status      string `json:"status"` // "online", "degraded", or "offline"
	Message     string `json:"message"`
	LatencyMs   *int64 `json:"latencyMs,omitempty"`
	CheckedAt   string `json:"checkedAt"`
}

// RecordHealthCheck inserts a health_checks row. ID and CheckedAt default to
// a fresh UUID and now (UTC) respectively when unset.
func (s *Store) RecordHealthCheck(ctx context.Context, hc *HealthCheckRecord) error {
	if hc.ID == "" {
		hc.ID = uuid.New().String()
	}
	if hc.CheckedAt == "" {
		hc.CheckedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO health_checks (id, connector_id, status, message, latency_ms, checked_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, hc.ID, hc.ConnectorID, hc.Status, hc.Message, hc.LatencyMs, hc.CheckedAt)
	if err != nil {
		return fmt.Errorf("record health check: %w", err)
	}
	return nil
}

// UptimeStats summarizes health_checks for one connector over a window,
// time-weighted between consecutive checks (see GetConnectorUptime).
type UptimeStats struct {
	ConnectorID     string  `json:"connectorId"`
	WindowStart     string  `json:"windowStart"`
	WindowEnd       string  `json:"windowEnd"`
	CheckCount      int     `json:"checkCount"`
	AvailabilityPct float64 `json:"availabilityPct"` // 0 when CheckCount == 0
	MTTRSeconds     float64 `json:"mttrSeconds"`     // mean time to recovery; 0 when no resolved outage
	OutageCount     int     `json:"outageCount"`     // resolved outages the MTTR mean is over
}

// GetConnectorUptime computes availability % and MTTR for connectorID over
// [since, until) from the health_checks time series.
//
// Availability is time-weighted, not a simple pass/fail ratio over samples:
// each check's status is assumed to hold from its checked_at until the next
// check (or until `until`, for the last one), and the fraction of that
// window spent in a non-"offline" status is the availability. This matches
// how checks are actually taken — on an interval, not continuously — without
// letting an uneven sampling rate skew the result.
//
// MTTR (mean time to recovery) is the mean, over every outage that both
// started and recovered inside the window, of the time from when a check
// first reported "offline" to the next check that reported something else.
// An outage still ongoing at `until` has no recovery time yet and is
// excluded from the mean, consistent with the usual MTTR definition (mean
// time to recover from *resolved* incidents).
func (s *Store) GetConnectorUptime(ctx context.Context, connectorID string, since, until time.Time) (UptimeStats, error) {
	sinceStr := since.UTC().Format(time.RFC3339)
	untilStr := until.UTC().Format(time.RFC3339)
	stats := UptimeStats{
		ConnectorID: connectorID,
		WindowStart: sinceStr,
		WindowEnd:   untilStr,
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT status, checked_at FROM health_checks
		WHERE connector_id = ? AND checked_at >= ? AND checked_at < ?
		ORDER BY checked_at ASC
	`, connectorID, sinceStr, untilStr)
	if err != nil {
		return stats, fmt.Errorf("get connector uptime: %w", err)
	}
	defer rows.Close() //nolint:errcheck

	type point struct {
		status    string
		checkedAt time.Time
	}
	var points []point
	for rows.Next() {
		var status, checkedAt string
		if err := rows.Scan(&status, &checkedAt); err != nil {
			return stats, fmt.Errorf("scan health check: %w", err)
		}
		ts, err := time.Parse(time.RFC3339, checkedAt)
		if err != nil {
			return stats, fmt.Errorf("parse checked_at %q: %w", checkedAt, err)
		}
		points = append(points, point{status: status, checkedAt: ts})
	}
	if err := rows.Err(); err != nil {
		return stats, fmt.Errorf("iterate health checks: %w", err)
	}

	stats.CheckCount = len(points)
	if len(points) == 0 {
		return stats, nil
	}

	var upDuration, totalDuration time.Duration
	var outageStart time.Time
	inOutage := false
	var recoveries []time.Duration

	for i, p := range points {
		windowEnd := until
		if i+1 < len(points) {
			windowEnd = points[i+1].checkedAt
		}
		segment := windowEnd.Sub(p.checkedAt)
		if segment < 0 {
			segment = 0
		}
		totalDuration += segment
		if p.status != "offline" {
			upDuration += segment
		}

		if p.status == "offline" && !inOutage {
			inOutage = true
			outageStart = p.checkedAt
		} else if p.status != "offline" && inOutage {
			recoveries = append(recoveries, p.checkedAt.Sub(outageStart))
			inOutage = false
		}
	}

	if totalDuration > 0 {
		stats.AvailabilityPct = float64(upDuration) / float64(totalDuration) * 100
	}
	stats.OutageCount = len(recoveries)
	if len(recoveries) > 0 {
		var sum time.Duration
		for _, d := range recoveries {
			sum += d
		}
		stats.MTTRSeconds = (sum / time.Duration(len(recoveries))).Seconds()
	}

	return stats, nil
}

// DeleteOldHealthChecks removes health_checks rows checked before cutoff.
// There is no "keep latest" guard, unlike snapshots: a connector's current
// status lives on the connectors row itself, not derived from this table.
func (s *Store) DeleteOldHealthChecks(ctx context.Context, cutoff string) (int64, error) {
	n, err := s.batchDelete(ctx, "health_checks", `t.checked_at < ?`, cutoff)
	if err != nil {
		return n, fmt.Errorf("delete old health checks: %w", err)
	}
	return n, nil
}
