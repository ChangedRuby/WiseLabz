package store

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDeleteOldSnapshotsProtectsLatest(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	// Only snapshot for this connector, far in the past: must never be deleted,
	// it's what GetLatestSnapshot serves as "current data" for the service.
	if err := s.CreateSnapshot(ctx, &SnapshotRecord{ConnectorID: connectorID, Data: "{}", FetchedAt: old}); err != nil {
		t.Fatalf("CreateSnapshot(only) error: %v", err)
	}

	n, err := s.DeleteOldSnapshots(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldSnapshots() error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldSnapshots() deleted %d rows, want 0 (only/latest snapshot protected)", n)
	}
	if _, err := s.GetLatestSnapshot(ctx, connectorID); err != nil {
		t.Fatalf("GetLatestSnapshot() after cleanup error: %v", err)
	}

	// Add a newer snapshot; now the old one is no longer latest and should be purged.
	recent := time.Now().UTC().Format(time.RFC3339)
	if err := s.CreateSnapshot(ctx, &SnapshotRecord{ConnectorID: connectorID, Data: "{}", FetchedAt: recent}); err != nil {
		t.Fatalf("CreateSnapshot(recent) error: %v", err)
	}

	n, err = s.DeleteOldSnapshots(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldSnapshots() error: %v", err)
	}
	if n != 1 {
		t.Fatalf("DeleteOldSnapshots() deleted %d rows, want 1 (superseded old snapshot)", n)
	}

	snaps, err := s.GetSnapshotsByConnector(ctx, connectorID, 20)
	if err != nil {
		t.Fatalf("GetSnapshotsByConnector() error: %v", err)
	}
	if len(snaps) != 1 || snaps[0].FetchedAt != recent {
		t.Fatalf("snapshots after cleanup = %+v, want only the recent one", snaps)
	}

	// Idempotency: running again with the same cutoff deletes nothing.
	n, err = s.DeleteOldSnapshots(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldSnapshots() second call error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldSnapshots() second call deleted %d rows, want 0", n)
	}
}

func TestDeleteOldDocVersionsProtectsCurrent(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	d := &DocRecord{Title: "Test Doc", Content: "v2", CurrentVersion: 2}
	if err := s.CreateDoc(ctx, d); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	// Superseded revision, old: should be deleted.
	if err := s.CreateDocVersion(ctx, &DocVersionRecord{DocID: d.ID, Rev: 1, Content: "v1", Trigger: "manual", CreatedAt: old}); err != nil {
		t.Fatalf("CreateDocVersion(rev1) error: %v", err)
	}
	// Current revision, also old: must never be deleted.
	if err := s.CreateDocVersion(ctx, &DocVersionRecord{DocID: d.ID, Rev: 2, Content: "v2", Trigger: "manual", CreatedAt: old}); err != nil {
		t.Fatalf("CreateDocVersion(rev2) error: %v", err)
	}

	n, err := s.DeleteOldDocVersions(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldDocVersions() error: %v", err)
	}
	if n != 1 {
		t.Fatalf("DeleteOldDocVersions() deleted %d rows, want 1 (superseded rev only)", n)
	}

	versions, err := s.GetDocVersions(ctx, d.ID)
	if err != nil {
		t.Fatalf("GetDocVersions() error: %v", err)
	}
	if len(versions) != 1 || versions[0].Rev != 2 {
		t.Fatalf("versions after cleanup = %+v, want only rev 2 (current)", versions)
	}

	// Idempotency.
	n, err = s.DeleteOldDocVersions(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldDocVersions() second call error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldDocVersions() second call deleted %d rows, want 0", n)
	}
}

func TestDeleteOldAlertsSkipsActiveStatuses(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	mustCreate := func(status string) string {
		a := &AlertRecord{ServiceID: connectorID, Severity: "info", Title: "t", Status: status, CreatedAt: old}
		if err := s.CreateAlert(ctx, a); err != nil {
			t.Fatalf("CreateAlert(%s) error: %v", status, err)
		}
		return a.ID
	}

	resolvedID := mustCreate("resolved")
	pendingID := mustCreate("pending")
	snoozedID := mustCreate("snoozed")
	dismissedID := mustCreate("dismissed")

	n, err := s.DeleteOldAlerts(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldAlerts() error: %v", err)
	}
	if n != 2 {
		t.Fatalf("DeleteOldAlerts() deleted %d rows, want 2 (resolved + dismissed)", n)
	}

	if _, err := s.GetAlert(ctx, resolvedID); err == nil {
		t.Fatal("resolved alert should have been deleted")
	}
	if _, err := s.GetAlert(ctx, dismissedID); err == nil {
		t.Fatal("dismissed alert should have been deleted")
	}
	if _, err := s.GetAlert(ctx, pendingID); err != nil {
		t.Fatalf("pending alert should survive, GetAlert() error: %v", err)
	}
	if _, err := s.GetAlert(ctx, snoozedID); err != nil {
		t.Fatalf("snoozed alert should survive, GetAlert() error: %v", err)
	}

	// Idempotency.
	n, err = s.DeleteOldAlerts(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldAlerts() second call error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldAlerts() second call deleted %d rows, want 0", n)
	}
}

func TestDeleteOldSyncRuns(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	recent := time.Now().UTC().Format(time.RFC3339)
	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	if err := s.CreateSyncRun(ctx, &SyncRunRecord{ConnectorID: connectorID, StartedAt: old, Status: SyncRunStatusSuccess}); err != nil {
		t.Fatalf("CreateSyncRun(old) error: %v", err)
	}
	if err := s.CreateSyncRun(ctx, &SyncRunRecord{ConnectorID: connectorID, StartedAt: recent, Status: SyncRunStatusSuccess}); err != nil {
		t.Fatalf("CreateSyncRun(recent) error: %v", err)
	}

	n, err := s.DeleteOldSyncRuns(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldSyncRuns() error: %v", err)
	}
	if n != 1 {
		t.Fatalf("DeleteOldSyncRuns() deleted %d rows, want 1", n)
	}

	runs, err := s.ListSyncRunsByConnector(ctx, connectorID, 20)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnector() error: %v", err)
	}
	if len(runs) != 1 || runs[0].StartedAt != recent {
		t.Fatalf("runs after cleanup = %+v, want only the recent run", runs)
	}

	// Idempotency.
	n, err = s.DeleteOldSyncRuns(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldSyncRuns() second call error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldSyncRuns() second call deleted %d rows, want 0", n)
	}
}

func TestDeleteOldAuditRecords(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	old := time.Now().UTC().AddDate(0, 0, -365).Format(time.RFC3339)
	recent := time.Now().UTC().Format(time.RFC3339)
	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	if err := s.CreateAuditRecord(ctx, &AuditRecord{ActorUserID: "u1", ActorRole: "operator", Action: "test.old", CreatedAt: old}); err != nil {
		t.Fatalf("CreateAuditRecord(old) error: %v", err)
	}
	if err := s.CreateAuditRecord(ctx, &AuditRecord{ActorUserID: "u1", ActorRole: "operator", Action: "test.recent", CreatedAt: recent}); err != nil {
		t.Fatalf("CreateAuditRecord(recent) error: %v", err)
	}

	n, err := s.DeleteOldAuditRecords(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldAuditRecords() error: %v", err)
	}
	if n != 1 {
		t.Fatalf("DeleteOldAuditRecords() deleted %d rows, want 1", n)
	}

	records, total, err := s.ListAuditRecords(ctx, "", "", "", "", 0, 20)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if total != 1 || len(records) != 1 || records[0].Action != "test.recent" {
		t.Fatalf("records after cleanup = %+v, want only the recent record", records)
	}

	// Idempotency.
	n, err = s.DeleteOldAuditRecords(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldAuditRecords() second call error: %v", err)
	}
	if n != 0 {
		t.Fatalf("DeleteOldAuditRecords() second call deleted %d rows, want 0", n)
	}
}

func TestDeleteOldSnapshotsBatches(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)

	base := time.Now().UTC().AddDate(0, 0, -365)
	total := retentionBatchSize*2 + 50
	for i := 0; i < total; i++ {
		ts := base.Add(time.Duration(i) * time.Second).Format(time.RFC3339)
		if err := s.CreateSnapshot(ctx, &SnapshotRecord{ConnectorID: connectorID, Data: "{}", FetchedAt: ts}); err != nil {
			t.Fatalf("CreateSnapshot(%d) error: %v", i, err)
		}
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)
	n, err := s.DeleteOldSnapshots(ctx, cutoff)
	if err != nil {
		t.Fatalf("DeleteOldSnapshots() error: %v", err)
	}
	if want := int64(total - 1); n != want {
		t.Fatalf("DeleteOldSnapshots() deleted %d rows, want %d (all but latest)", n, want)
	}
}

func TestDeleteOldSnapshotsStopsOnCancelledContext(t *testing.T) {
	s := newDocTestStore(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := s.DeleteOldSnapshots(ctx, time.Now().UTC().Format(time.RFC3339)); err == nil {
		t.Fatal("DeleteOldSnapshots() with cancelled ctx returned nil error")
	}
}

func TestUnsnoozeExpiredAlerts(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	if n, err := s.UnsnoozeExpiredAlerts(ctx); err != nil || n != 0 {
		t.Fatalf("UnsnoozeExpiredAlerts() on empty table = %d, %v; want 0, nil", n, err)
	}
}

// assertQueryPlanUsesIndex fails unless SQLite's planner reports a SEARCH
// using the named index for the given query. An index nothing uses is pure
// write overhead, so every index migration carries one of these.
func assertQueryPlanUsesIndex(t *testing.T, s *Store, index, query string, args ...any) {
	t.Helper()
	rows, err := s.DB().QueryContext(context.Background(), "EXPLAIN QUERY PLAN "+query, args...)
	if err != nil {
		t.Fatalf("EXPLAIN QUERY PLAN error: %v", err)
	}
	defer rows.Close() //nolint:errcheck

	var plan []string
	for rows.Next() {
		var id, parent, unused int
		var detail string
		if err := rows.Scan(&id, &parent, &unused, &detail); err != nil {
			t.Fatalf("scan query plan: %v", err)
		}
		plan = append(plan, detail)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate query plan: %v", err)
	}
	for _, detail := range plan {
		if strings.Contains(detail, "SEARCH") && strings.Contains(detail, index) {
			return
		}
	}
	t.Errorf("query plan did not SEARCH using %s: %v", index, plan)
}

// TestDeleteStaleSessionsUsesIndex proves 000032's idx_sessions_last_seen
// covers the retention sweep instead of leaving it to scan sessions (#299).
func TestDeleteStaleSessionsUsesIndex(t *testing.T) {
	skipOnPostgres(t, "EXPLAIN QUERY PLAN is SQLite-only")
	s := newDocTestStore(t)
	assertQueryPlanUsesIndex(t, s, "idx_sessions_last_seen", fmt.Sprintf(
		`DELETE FROM sessions WHERE id IN (SELECT t.id FROM sessions t WHERE t.last_seen_at < ? LIMIT %d)`,
		retentionBatchSize), time.Now().UTC().Format(time.RFC3339))
}
