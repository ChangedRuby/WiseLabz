package store

import (
	"context"
	"testing"
	"time"
)

// TestRecordHealthCheckDefaults verifies RecordHealthCheck fills in a UUID
// and CheckedAt when left unset, and persists the row.
func TestRecordHealthCheckDefaults(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	latency := int64(42)
	hc := &HealthCheckRecord{ConnectorID: connectorID, Status: "online", Message: "Healthy", LatencyMs: &latency}
	if err := s.RecordHealthCheck(ctx, hc); err != nil {
		t.Fatalf("RecordHealthCheck() error: %v", err)
	}
	if hc.ID == "" {
		t.Error("RecordHealthCheck() left ID empty")
	}
	if hc.CheckedAt == "" {
		t.Error("RecordHealthCheck() left CheckedAt empty")
	}

	stats, err := s.GetConnectorUptime(ctx, connectorID, time.Now().UTC().Add(-time.Hour), time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}
	if stats.CheckCount != 1 {
		t.Fatalf("CheckCount = %d, want 1", stats.CheckCount)
	}
}

// TestGetConnectorUptimeDeterministicOutage builds a fixed, known timeline of
// checks — one hour apart over a 6h window, with a single 2h outage in the
// middle — and asserts the exact availability % and MTTR that timeline
// implies, rather than just checking they're "reasonable".
//
// Timeline (checked_at, status), window = [t0, t0+6h):
//
//	t0+0h online   -> up   for 1h (until t0+1h)
//	t0+1h online   -> up   for 1h (until t0+2h)
//	t0+2h offline  -> down for 1h (until t0+3h)   <- outage starts
//	t0+3h offline  -> down for 1h (until t0+4h)
//	t0+4h online   -> up   for 2h (until t0+6h, window end)  <- recovers here
//
// Up total = 1+1+2 = 4h, down total = 1+1 = 2h, over a 6h window:
// availability = 4/6 * 100 = 66.666...%. The one outage ran from t0+2h to
// t0+4h, a 2h recovery time, so MTTR = 2h = 7200s.
func TestGetConnectorUptimeDeterministicOutage(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	checks := []struct {
		offset time.Duration
		status string
	}{
		{0, "online"},
		{1 * time.Hour, "online"},
		{2 * time.Hour, "offline"},
		{3 * time.Hour, "offline"},
		{4 * time.Hour, "online"},
	}
	for _, c := range checks {
		hc := &HealthCheckRecord{
			ConnectorID: connectorID,
			Status:      c.status,
			CheckedAt:   t0.Add(c.offset).Format(time.RFC3339),
		}
		if err := s.RecordHealthCheck(ctx, hc); err != nil {
			t.Fatalf("RecordHealthCheck(%v) error: %v", c.offset, err)
		}
	}

	stats, err := s.GetConnectorUptime(ctx, connectorID, t0, t0.Add(6*time.Hour))
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}

	if stats.CheckCount != 5 {
		t.Errorf("CheckCount = %d, want 5", stats.CheckCount)
	}
	wantAvailability := float64(4) / float64(6) * 100
	if diff := stats.AvailabilityPct - wantAvailability; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("AvailabilityPct = %v, want %v", stats.AvailabilityPct, wantAvailability)
	}
	if stats.OutageCount != 1 {
		t.Fatalf("OutageCount = %d, want 1", stats.OutageCount)
	}
	if stats.MTTRSeconds != (2 * time.Hour).Seconds() {
		t.Errorf("MTTRSeconds = %v, want %v", stats.MTTRSeconds, (2 * time.Hour).Seconds())
	}
}

// TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR verifies an outage
// still ongoing at the window end lowers availability but is excluded from
// the MTTR mean (it has no recovery time yet).
func TestGetConnectorUptimeUnresolvedOutageExcludedFromMTTR(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, c := range []struct {
		offset time.Duration
		status string
	}{
		{0, "online"},
		{1 * time.Hour, "offline"}, // never recovers within the window
	} {
		hc := &HealthCheckRecord{ConnectorID: connectorID, Status: c.status, CheckedAt: t0.Add(c.offset).Format(time.RFC3339)}
		if err := s.RecordHealthCheck(ctx, hc); err != nil {
			t.Fatalf("RecordHealthCheck(%v) error: %v", c.offset, err)
		}
	}

	stats, err := s.GetConnectorUptime(ctx, connectorID, t0, t0.Add(4*time.Hour))
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}
	// up for 1h (0-1h), down for 3h (1h-4h) => 25% available.
	if diff := stats.AvailabilityPct - 25.0; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("AvailabilityPct = %v, want 25", stats.AvailabilityPct)
	}
	if stats.OutageCount != 0 {
		t.Errorf("OutageCount = %d, want 0 (unresolved outage excluded)", stats.OutageCount)
	}
	if stats.MTTRSeconds != 0 {
		t.Errorf("MTTRSeconds = %v, want 0 (no resolved outage)", stats.MTTRSeconds)
	}
}

// TestGetConnectorUptimeNoData verifies an empty window returns zero values
// without error, rather than dividing by zero.
func TestGetConnectorUptimeNoData(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	stats, err := s.GetConnectorUptime(ctx, connectorID, time.Now().UTC().Add(-time.Hour), time.Now().UTC())
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}
	if stats.CheckCount != 0 || stats.AvailabilityPct != 0 || stats.MTTRSeconds != 0 || stats.OutageCount != 0 {
		t.Errorf("GetConnectorUptime() with no data = %+v, want all zero", stats)
	}
}

// TestDeleteOldHealthChecks verifies cutoff-based pruning with no
// "keep latest" guard (unlike snapshots): every row older than cutoff goes,
// including the connector's most recent check if it's old enough.
func TestDeleteOldHealthChecks(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	old := time.Now().UTC().AddDate(0, 0, -100).Format(time.RFC3339)
	recent := time.Now().UTC().Format(time.RFC3339)
	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	if err := s.RecordHealthCheck(ctx, &HealthCheckRecord{ConnectorID: connectorID, Status: "offline", CheckedAt: old}); err != nil {
		t.Fatalf("RecordHealthCheck(old) error: %v", err)
	}
	if err := s.RecordHealthCheck(ctx, &HealthCheckRecord{ConnectorID: connectorID, Status: "online", CheckedAt: recent}); err != nil {
		t.Fatalf("RecordHealthCheck(recent) error: %v", err)
	}

	n, err := s.DeleteOldHealthChecks(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldHealthChecks() error: %v", err)
	}
	if n != 1 {
		t.Fatalf("DeleteOldHealthChecks() deleted %d rows, want 1", n)
	}

	stats, err := s.GetConnectorUptime(ctx, connectorID, time.Now().UTC().AddDate(0, 0, -200), time.Now().UTC().AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("GetConnectorUptime() error: %v", err)
	}
	if stats.CheckCount != 1 {
		t.Fatalf("CheckCount after cleanup = %d, want 1", stats.CheckCount)
	}

	// Idempotent: running again with the same cutoff deletes nothing further.
	n, err = s.DeleteOldHealthChecks(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldHealthChecks() second call error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldHealthChecks() second call deleted %d rows, want 0", n)
	}
}
