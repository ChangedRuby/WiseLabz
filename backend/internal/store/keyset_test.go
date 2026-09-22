package store

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// traverse walks every keyset page of list and returns the ids in order. It
// fails if the traversal does not terminate within a generous bound, which is
// what an off-by-one in the cursor comparison looks like.
func traverse[T any](t *testing.T, pageSize int, key func(T) (string, string), list func(cur Keyset) ([]T, int, error)) []string {
	t.Helper()

	var ids []string
	cur := Keyset{}
	for page := 0; ; page++ {
		if page > 100 {
			t.Fatalf("keyset traversal did not terminate after %d pages (ids so far: %d)", page, len(ids))
		}
		items, _, err := list(cur)
		if err != nil {
			t.Fatalf("keyset page %d: %v", page, err)
		}
		for _, item := range items {
			sort, id := key(item)
			ids = append(ids, id)
			cur = Keyset{Sort: sort, ID: id}
		}
		if len(items) < pageSize {
			return ids
		}
	}
}

// assertSameSet fails when the keyset traversal and the offset listing did not
// visit exactly the same rows in exactly the same order: a duplicate or a gap
// across page boundaries shows up here.
func assertSameSet(t *testing.T, what string, keyset, offset []string) {
	t.Helper()

	if len(keyset) != len(offset) {
		t.Fatalf("%s: keyset visited %d rows, offset listing has %d", what, len(keyset), len(offset))
	}
	seen := make(map[string]bool, len(keyset))
	for i, id := range keyset {
		if seen[id] {
			t.Fatalf("%s: id %s visited twice by the keyset traversal", what, id)
		}
		seen[id] = true
		if id != offset[i] {
			t.Errorf("%s: row %d = %s, offset listing has %s", what, i, id, offset[i])
		}
	}
}

func TestListAuditRecordsKeysetTraversal(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	// 17 records over 4 distinct timestamps, so most pages of 5 land inside a
	// run of rows sharing created_at — the tie case the (created_at, id) row
	// comparison exists for.
	base := time.Now().UTC().Truncate(time.Second)
	for i := range 17 {
		if err := s.CreateAuditRecord(ctx, &AuditRecord{
			ActorUserID: "user-1", Action: "connector.create", TargetType: "connector",
			TargetID:  fmt.Sprintf("conn-%02d", i),
			CreatedAt: base.Add(-time.Duration(i/5) * time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("CreateAuditRecord(%d) error: %v", i, err)
		}
	}

	const pageSize = 5
	got := traverse(t, pageSize, func(a AuditRecord) (string, string) { return a.CreatedAt, a.ID },
		func(cur Keyset) ([]AuditRecord, int, error) {
			return s.ListAuditRecordsKeyset(ctx, "", "", "", "", cur, pageSize)
		})

	// The offset listing orders by created_at only, so compare against the
	// keyset listing's own first-page-at-a-time ordering instead: read every
	// row in one keyset page and use that as the reference order.
	all, total, err := s.ListAuditRecordsKeyset(ctx, "", "", "", "", Keyset{}, 100)
	if err != nil {
		t.Fatalf("ListAuditRecordsKeyset(all) error: %v", err)
	}
	if total != 17 {
		t.Errorf("total = %d, want 17 (the full filtered count, not the page size)", total)
	}
	want := make([]string, len(all))
	for i, a := range all {
		want[i] = a.ID
	}
	assertSameSet(t, "audit records", got, want)
}

func TestListAuditRecordsKeysetHonoursFilters(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)

	base := time.Now().UTC().Truncate(time.Second)
	for i := range 6 {
		action := "connector.create"
		if i%2 == 0 {
			action = "doc.restore"
		}
		if err := s.CreateAuditRecord(ctx, &AuditRecord{
			ActorUserID: "user-1", Action: action, TargetType: "connector",
			TargetID: fmt.Sprintf("t-%d", i), CreatedAt: base.Add(-time.Duration(i) * time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("CreateAuditRecord(%d) error: %v", i, err)
		}
	}

	got := traverse(t, 2, func(a AuditRecord) (string, string) { return a.CreatedAt, a.ID },
		func(cur Keyset) ([]AuditRecord, int, error) {
			return s.ListAuditRecordsKeyset(ctx, "doc.restore", "", "", "", cur, 2)
		})
	if len(got) != 3 {
		t.Fatalf("filtered keyset traversal visited %d records, want 3", len(got))
	}
}

func TestListChangesKeysetTraversal(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	serviceID := seedConnectorForChanges(t, s)

	base := time.Now().UTC().Truncate(time.Second)
	for i := range 14 {
		c := &ChangeRecord{
			ServiceID: serviceID, ChangeType: "modified", Severity: "warning",
			Summary:    fmt.Sprintf("change %02d", i),
			DetectedAt: base.Add(-time.Duration(i/4) * time.Minute).Format(time.RFC3339),
		}
		if err := s.CreateChange(ctx, c); err != nil {
			t.Fatalf("CreateChange(%d) error: %v", i, err)
		}
	}

	const pageSize = 4
	got := traverse(t, pageSize, func(c ChangeRecord) (string, string) { return c.DetectedAt, c.ID },
		func(cur Keyset) ([]ChangeRecord, int, error) {
			return s.ListChangesKeyset(ctx, serviceID, "", cur, pageSize)
		})

	all, total, err := s.ListChangesKeyset(ctx, serviceID, "", Keyset{}, 100)
	if err != nil {
		t.Fatalf("ListChangesKeyset(all) error: %v", err)
	}
	if total != 14 {
		t.Errorf("total = %d, want 14", total)
	}
	want := make([]string, len(all))
	for i, c := range all {
		want[i] = c.ID
	}
	assertSameSet(t, "changes", got, want)
}

func TestListSyncRunsByConnectorKeysetTraversal(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)
	otherID := createTestConnector(ctx, t, s)

	base := time.Now().UTC().Truncate(time.Second)
	for i := range 11 {
		if err := s.CreateSyncRun(ctx, &SyncRunRecord{
			ConnectorID: connectorID, Status: SyncRunStatusSuccess,
			StartedAt: base.Add(-time.Duration(i/3) * time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("CreateSyncRun(%d) error: %v", i, err)
		}
	}
	if err := s.CreateSyncRun(ctx, &SyncRunRecord{ConnectorID: otherID, Status: SyncRunStatusSuccess}); err != nil {
		t.Fatalf("CreateSyncRun(other) error: %v", err)
	}

	const pageSize = 3
	got := traverse(t, pageSize, func(r SyncRunRecord) (string, string) { return r.StartedAt, r.ID },
		func(cur Keyset) ([]SyncRunRecord, int, error) {
			return s.ListSyncRunsByConnectorKeyset(ctx, connectorID, cur, pageSize)
		})

	all, total, err := s.ListSyncRunsByConnectorKeyset(ctx, connectorID, Keyset{}, 100)
	if err != nil {
		t.Fatalf("ListSyncRunsByConnectorKeyset(all) error: %v", err)
	}
	if total != 11 {
		t.Errorf("total = %d, want 11 (other connectors excluded)", total)
	}
	want := make([]string, len(all))
	for i, r := range all {
		want[i] = r.ID
	}
	assertSameSet(t, "sync runs", got, want)
}

// TestKeysetQueriesUseCoveringIndexes is the proof for migration 000033: every
// keyset query must resolve entirely from an index, with no per-page sort. A
// "USE TEMP B-TREE FOR ... ORDER BY" line in the plan means the planner is
// sorting each page, which is the cost keyset pagination is meant to remove.
//
// SQLite-only: EXPLAIN QUERY PLAN has no PostgreSQL equivalent, and the
// indexes are identical in both dialects (see the 000033 migration pair).
func TestKeysetQueriesUseCoveringIndexes(t *testing.T) {
	if os.Getenv("WISELABZ_TEST_POSTGRES_DSN") != "" {
		t.Skip("EXPLAIN QUERY PLAN is SQLite-only")
	}
	ctx := context.Background()
	s := newDocTestStore(t)

	tests := []struct {
		name  string
		query string
		args  []any
		index string
	}{
		{
			name:  "audit unfiltered",
			query: `SELECT ` + auditColumns + ` FROM audit_log WHERE 1=1 AND (created_at, id) < (?, ?) ORDER BY created_at DESC, id DESC LIMIT ?`,
			args:  []any{"t", "i", 10},
			index: "idx_audit_log_created_id",
		},
		{
			name:  "audit by action",
			query: `SELECT ` + auditColumns + ` FROM audit_log WHERE 1=1 AND action = ? AND (created_at, id) < (?, ?) ORDER BY created_at DESC, id DESC LIMIT ?`,
			args:  []any{"a", "t", "i", 10},
			index: "idx_audit_log_action_created_id",
		},
		{
			name:  "audit by target type",
			query: `SELECT ` + auditColumns + ` FROM audit_log WHERE 1=1 AND target_type = ? AND (created_at, id) < (?, ?) ORDER BY created_at DESC, id DESC LIMIT ?`,
			args:  []any{"connector", "t", "i", 10},
			index: "idx_audit_log_target_type_created_id",
		},
		{
			name:  "changes unfiltered",
			query: `SELECT ` + changeColumns + ` FROM changes WHERE 1=1 AND (detected_at, id) < (?, ?) ORDER BY detected_at DESC, id DESC LIMIT ?`,
			args:  []any{"t", "i", 10},
			index: "idx_changes_detected_id",
		},
		{
			name:  "changes by service",
			query: `SELECT ` + changeColumns + ` FROM changes WHERE 1=1 AND service_id = ? AND (detected_at, id) < (?, ?) ORDER BY detected_at DESC, id DESC LIMIT ?`,
			args:  []any{"svc", "t", "i", 10},
			index: "idx_changes_service_detected_id",
		},
		{
			name:  "changes by severity",
			query: `SELECT ` + changeColumns + ` FROM changes WHERE 1=1 AND severity = ? AND (detected_at, id) < (?, ?) ORDER BY detected_at DESC, id DESC LIMIT ?`,
			args:  []any{"warning", "t", "i", 10},
			index: "idx_changes_severity_detected_id",
		},
		{
			name:  "sync runs by connector",
			query: `SELECT ` + syncRunColumns + ` FROM sync_runs WHERE connector_id = ? AND (started_at, id) < (?, ?) ORDER BY started_at DESC, id DESC LIMIT ?`,
			args:  []any{"conn", "t", "i", 10},
			index: "idx_sync_runs_connector_started_id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plan := queryPlan(ctx, t, s, tc.query, tc.args)
			if !strings.Contains(plan, tc.index) {
				t.Errorf("plan does not use %s:\n%s", tc.index, plan)
			}
			if strings.Contains(plan, "TEMP B-TREE") {
				t.Errorf("plan sorts each page instead of reading it in index order:\n%s", plan)
			}
		})
	}
}

func queryPlan(ctx context.Context, t *testing.T, s *Store, query string, args []any) string {
	t.Helper()

	rows, err := s.db.QueryContext(ctx, "EXPLAIN QUERY PLAN "+query, args...)
	if err != nil {
		t.Fatalf("EXPLAIN QUERY PLAN: %v", err)
	}
	defer rows.Close() //nolint:errcheck

	var plan strings.Builder
	for rows.Next() {
		var id, parent, notUsed int
		var detail string
		if err := rows.Scan(&id, &parent, &notUsed, &detail); err != nil {
			t.Fatalf("scan plan row: %v", err)
		}
		plan.WriteString(detail)
		plan.WriteByte('\n')
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate plan: %v", err)
	}
	return plan.String()
}
