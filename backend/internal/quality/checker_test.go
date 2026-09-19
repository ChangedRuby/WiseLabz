package quality

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"

	_ "modernc.org/sqlite"
)

func createComplianceRule(t *testing.T, s *store.Store, connectorType, title string) *store.ComplianceRuleRecord {
	t.Helper()
	rule := &store.ComplianceRuleRecord{
		Name: title, ConnectorType: connectorType, EntityKind: "vm",
		Conditions: `[{"attribute":"firewall_enabled","op":"eq","value":false}]`,
		Severity:   "critical", Title: title, RemediationLink: "https://example.test/remediate", Enabled: true,
	}
	if err := s.CreateComplianceRule(context.Background(), rule); err != nil {
		t.Fatalf("CreateComplianceRule() error: %v", err)
	}
	return rule
}

func createComplianceSnapshot(t *testing.T, s *store.Store, connectorID string, firewallEnabled bool) {
	t.Helper()
	data, err := json.Marshal(connector.ServiceSnapshot{Entities: []connector.SnapshotEntity{{
		Kind: "vm", Name: "vm-01", Attributes: map[string]any{"firewall_enabled": firewallEnabled},
	}}})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	fetchedAt := "2020-01-01T00:00:00Z"
	if firewallEnabled {
		fetchedAt = "2030-01-01T00:00:00Z"
	}
	if err := s.CreateSnapshot(context.Background(), &store.SnapshotRecord{ConnectorID: connectorID, Data: string(data), FetchedAt: fetchedAt}); err != nil {
		t.Fatalf("CreateSnapshot() error: %v", err)
	}
}

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/quality.db?cache=shared")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := store.RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
	return store.New(db, "sqlite")
}

func createConnector(t *testing.T, s *store.Store, owner string) *store.ConnectorRecord {
	t.Helper()
	connector := &store.ConnectorRecord{
		Name:     "Test connector",
		Category: "virtualization",
		Type:     "proxmox",
		URL:      "https://example.test",
		Owner:    owner,
	}
	if err := s.CreateConnector(context.Background(), connector); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return connector
}

func findings(t *testing.T, s *store.Store, connectorID, checkType, status string) []store.QualityFindingRecord {
	t.Helper()
	items, _, err := s.ListQualityFindings(context.Background(), connectorID, checkType, status, "", 0, 20)
	if err != nil {
		t.Fatalf("ListQualityFindings() error: %v", err)
	}
	return items
}

func TestCheckStaleDetectsAndAutoResolves(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	stale := &store.DocRecord{
		ID:        "stale-doc",
		Title:     "Old runbook",
		ServiceID: connector.ID,
		Content:   strings.Repeat("useful content ", 4),
		UpdatedAt: time.Now().UTC().Add(-StaleThreshold - time.Hour).Format(time.RFC3339),
	}
	fresh := &store.DocRecord{
		ID:        "fresh-doc",
		Title:     "Fresh runbook",
		ServiceID: connector.ID,
		Content:   strings.Repeat("useful content ", 4),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	for _, doc := range []*store.DocRecord{stale, fresh} {
		if err := s.CreateDoc(ctx, doc); err != nil {
			t.Fatalf("CreateDoc(%s) error: %v", doc.ID, err)
		}
	}

	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	open := findings(t, s, connector.ID, "stale", "open")
	if len(open) != 1 || open[0].DocID != stale.ID || open[0].RemediationLink != "/docs/"+stale.ID {
		t.Fatalf("stale finding = %#v, want stale doc target", open)
	}

	if err := s.UpdateDoc(ctx, stale.ID, stale.Content, nil); err != nil {
		t.Fatalf("UpdateDoc() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, connector.ID, "stale", "open"); len(got) != 0 {
		t.Fatalf("open stale findings = %d, want 0", len(got))
	}
	if got := findings(t, s, connector.ID, "stale", "resolved"); len(got) != 1 {
		t.Fatalf("resolved stale findings = %d, want 1", len(got))
	}
}

func TestCheckEmptyDetectsAndAutoResolves(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	empty := &store.DocRecord{ID: "empty-doc", Title: "Empty runbook", ServiceID: connector.ID, Content: " \n\t "}
	healthy := &store.DocRecord{ID: "healthy-doc", Title: "Healthy runbook", ServiceID: connector.ID, Content: strings.Repeat("documented ", 5)}
	for _, doc := range []*store.DocRecord{empty, healthy} {
		if err := s.CreateDoc(ctx, doc); err != nil {
			t.Fatalf("CreateDoc(%s) error: %v", doc.ID, err)
		}
	}

	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	open := findings(t, s, connector.ID, "empty", "open")
	if len(open) != 1 || open[0].DocID != empty.ID || open[0].RemediationLink != "/docs/"+empty.ID+"/edit" {
		t.Fatalf("empty finding = %#v, want empty doc target", open)
	}

	if err := s.UpdateDoc(ctx, empty.ID, strings.Repeat("documented ", 5), nil); err != nil {
		t.Fatalf("UpdateDoc() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, connector.ID, "empty", "open"); len(got) != 0 {
		t.Fatalf("open empty findings = %d, want 0", len(got))
	}
	if got := findings(t, s, connector.ID, "empty", "resolved"); len(got) != 1 {
		t.Fatalf("resolved empty findings = %d, want 1", len(got))
	}
}

func TestCheckFailingDetectsAndAutoResolves(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	base := time.Now().UTC().Add(-time.Hour)
	for i := 0; i < ConsecutiveFailuresThreshold; i++ {
		run := &store.SyncRunRecord{
			ConnectorID: connector.ID,
			StartedAt:   base.Add(time.Duration(i) * time.Minute).Format(time.RFC3339),
			Status:      store.SyncRunStatusError,
			Error:       "connection failed",
		}
		if err := s.CreateSyncRun(ctx, run); err != nil {
			t.Fatalf("CreateSyncRun(%d) error: %v", i, err)
		}
	}

	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	open := findings(t, s, connector.ID, "failing", "open")
	if len(open) != 1 || open[0].Severity != "critical" || open[0].RemediationLink != "/services/"+connector.ID {
		t.Fatalf("failing finding = %#v, want critical service target", open)
	}

	if err := s.CreateSyncRun(ctx, &store.SyncRunRecord{
		ConnectorID: connector.ID,
		StartedAt:   time.Now().UTC().Format(time.RFC3339),
		Status:      store.SyncRunStatusSuccess,
	}); err != nil {
		t.Fatalf("CreateSyncRun(success) error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, connector.ID, "failing", "open"); len(got) != 0 {
		t.Fatalf("open failing findings = %d, want 0", len(got))
	}
	if got := findings(t, s, connector.ID, "failing", "resolved"); len(got) != 1 {
		t.Fatalf("resolved failing findings = %d, want 1", len(got))
	}
}

func TestCheckOwnershipDetectsAndAutoResolves(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, " \t")
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})

	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() detect error: %v", err)
	}
	open := findings(t, s, connector.ID, "ownership_incomplete", "open")
	if len(open) != 1 || open[0].Severity != "info" || open[0].RemediationLink != "/connectors/"+connector.ID+"/edit" {
		t.Fatalf("ownership finding = %#v, want info connector target", open)
	}

	if err := s.UpdateConnector(ctx, connector.ID, map[string]any{"owner": "platform-team"}); err != nil {
		t.Fatalf("UpdateConnector(owner) error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, connector.ID, "ownership_incomplete", "open"); len(got) != 0 {
		t.Fatalf("open ownership findings = %d, want 0", len(got))
	}
	if got := findings(t, s, connector.ID, "ownership_incomplete", "resolved"); len(got) != 1 {
		t.Fatalf("resolved ownership findings = %d, want 1", len(got))
	}
}

func TestRunForConnectorSkipsDocWithMalformedTimestamp(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "")
	if err := s.CreateDoc(ctx, &store.DocRecord{
		Title:     "Broken timestamp",
		ServiceID: connector.ID,
		Content:   strings.Repeat("documented ", 5),
		UpdatedAt: "not-a-timestamp",
	}); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}
	if err := s.CreateDoc(ctx, &store.DocRecord{
		Title:     "Forgotten runbook",
		ServiceID: connector.ID,
		Content:   strings.Repeat("documented ", 5),
		UpdatedAt: time.Now().UTC().Add(-StaleThreshold - time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	if err := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14}).RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("RunForConnector() error = %v, want nil", err)
	}
	if got := findings(t, s, connector.ID, "stale", "open"); len(got) != 1 {
		t.Fatalf("stale findings = %d, want 1 (malformed doc skipped, other doc still checked)", len(got))
	}
	if got := findings(t, s, connector.ID, "ownership_incomplete", "open"); len(got) != 1 {
		t.Fatalf("ownership findings = %d, want 1", len(got))
	}
}

func TestRunStaleSweepOnceCoversConnectorsWithNoRecentSync(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	if err := s.CreateDoc(ctx, &store.DocRecord{
		Title:     "Forgotten runbook",
		ServiceID: connector.ID,
		Content:   strings.Repeat("documented ", 5),
		UpdatedAt: time.Now().UTC().Add(-StaleThreshold - time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("CreateDoc() error: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	RunStaleSweepOnce(ctx, s, nil, nil, logger)

	if got := findings(t, s, connector.ID, "stale", "open"); len(got) != 1 {
		t.Fatalf("open stale findings = %d, want 1", len(got))
	}
}

func TestComplianceRulesDetectResolveAndEvaluateOnSave(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	matched := createConnector(t, s, "owner")
	other := createConnector(t, s, "owner")
	other.Type = "docker"
	if err := s.UpdateConnector(ctx, other.ID, map[string]any{"type": other.Type}); err != nil {
		t.Fatalf("UpdateConnector(type) error: %v", err)
	}
	createComplianceSnapshot(t, s, matched.ID, false)
	createComplianceSnapshot(t, s, other.ID, false)
	rule := createComplianceRule(t, s, "proxmox", "Firewall disabled")
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})

	if err := checker.EvaluateRule(ctx, rule.ID); err != nil {
		t.Fatalf("EvaluateRule() error: %v", err)
	}
	open := findings(t, s, matched.ID, "compliance", "open")
	if len(open) != 1 || open[0].RuleID != rule.ID || open[0].DetectedCount != 1 || open[0].Description != "Violating entities: vm-01." {
		t.Fatalf("open compliance finding = %#v", open)
	}
	if got := findings(t, s, other.ID, "compliance", "open"); len(got) != 0 {
		t.Fatalf("wrong connector type findings = %#v, want none", got)
	}

	if err := checker.RunForConnector(ctx, matched.ID); err != nil {
		t.Fatalf("RunForConnector() repeat error: %v", err)
	}
	open = findings(t, s, matched.ID, "compliance", "open")
	if len(open) != 1 || open[0].DetectedCount != 2 {
		t.Fatalf("re-detected finding = %#v, want one with count 2", open)
	}

	second := createComplianceRule(t, s, "proxmox", "Second firewall rule")
	if err := checker.RunForConnector(ctx, matched.ID); err != nil {
		t.Fatalf("RunForConnector() second rule error: %v", err)
	}
	open = findings(t, s, matched.ID, "compliance", "open")
	if len(open) != 2 {
		t.Fatalf("open compliance findings = %#v, want two rules", open)
	}
	if err := checker.ResolveRule(ctx, second.ID); err != nil {
		t.Fatalf("ResolveRule() error: %v", err)
	}
	if got := findings(t, s, matched.ID, "compliance", "open"); len(got) != 1 {
		t.Fatalf("open findings after ResolveRule = %#v, want one", got)
	}

	createComplianceSnapshot(t, s, matched.ID, true)
	if err := checker.RunForConnector(ctx, matched.ID); err != nil {
		t.Fatalf("RunForConnector() resolve error: %v", err)
	}
	if got := findings(t, s, matched.ID, "compliance", "open"); len(got) != 0 {
		t.Fatalf("open compliance findings after fix = %#v, want none", got)
	}
}

func TestQualityThresholdBoundaries(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	fixedNow := time.Date(2026, time.September, 5, 12, 0, 0, 0, time.UTC)
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	checker.now = func() time.Time { return fixedNow }

	tests := []struct {
		name      string
		content   string
		updatedAt time.Time
		checkType string
		wantOpen  bool
	}{
		{name: "exact stale threshold is fresh", content: strings.Repeat("x", EmptyContentMinChars), updatedAt: fixedNow.Add(-StaleThreshold), checkType: "stale"},
		{name: "past stale threshold is stale", content: strings.Repeat("x", EmptyContentMinChars), updatedAt: fixedNow.Add(-StaleThreshold - time.Nanosecond), checkType: "stale", wantOpen: true},
		{name: "exact empty threshold is populated", content: strings.Repeat("x", EmptyContentMinChars), updatedAt: fixedNow, checkType: "empty"},
		{name: "below empty threshold is empty", content: strings.Repeat("x", EmptyContentMinChars-1), updatedAt: fixedNow, checkType: "empty", wantOpen: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connector := createConnector(t, s, "platform-team")
			doc := &store.DocRecord{Title: tt.name, ServiceID: connector.ID, Content: tt.content, UpdatedAt: tt.updatedAt.Format(time.RFC3339Nano)}
			if err := s.CreateDoc(ctx, doc); err != nil {
				t.Fatalf("CreateDoc() error: %v", err)
			}
			if err := checker.RunForConnector(ctx, connector.ID); err != nil {
				t.Fatalf("RunForConnector() error: %v", err)
			}
			got := findings(t, s, connector.ID, tt.checkType, "open")
			if (len(got) == 1) != tt.wantOpen {
				t.Fatalf("open %s findings = %+v, want open %v", tt.checkType, got, tt.wantOpen)
			}
		})
	}
}

func TestStaleFindingMovesToRemainingStaleDoc(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	for _, doc := range []*store.DocRecord{
		{ID: "oldest", Title: "Oldest", ServiceID: connector.ID, Content: strings.Repeat("x", EmptyContentMinChars), UpdatedAt: time.Now().Add(-40 * 24 * time.Hour).Format(time.RFC3339)},
		{ID: "still-stale", Title: "Still stale", ServiceID: connector.ID, Content: strings.Repeat("x", EmptyContentMinChars), UpdatedAt: time.Now().Add(-35 * 24 * time.Hour).Format(time.RFC3339)},
	} {
		if err := s.CreateDoc(ctx, doc); err != nil {
			t.Fatalf("CreateDoc(%s) error: %v", doc.ID, err)
		}
	}
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("first RunForConnector() error: %v", err)
	}
	if err := s.UpdateDoc(ctx, "oldest", strings.Repeat("x", EmptyContentMinChars), nil); err != nil {
		t.Fatalf("UpdateDoc() error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("second RunForConnector() error: %v", err)
	}
	open := findings(t, s, connector.ID, "stale", "open")
	if len(open) != 1 || open[0].DocID != "still-stale" || open[0].DetectedCount != 2 {
		t.Fatalf("remaining stale finding = %+v", open)
	}
}

func TestFailingThresholdAndInterruptedStreak(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "platform-team")
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	base := time.Now().UTC().Truncate(time.Second)
	addRun := func(index int, status store.SyncRunStatus) {
		t.Helper()
		stamp := base.Add(time.Duration(index) * time.Millisecond).Format(time.RFC3339Nano)
		if err := s.CreateSyncRun(ctx, &store.SyncRunRecord{ConnectorID: connector.ID, StartedAt: stamp, FinishedAt: stamp, Status: status}); err != nil {
			t.Fatalf("CreateSyncRun(%d) error: %v", index, err)
		}
	}

	addRun(1, store.SyncRunStatusError)
	addRun(2, store.SyncRunStatusError)
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("two failures check error: %v", err)
	}
	if got := findings(t, s, connector.ID, "failing", "open"); len(got) != 0 {
		t.Fatalf("finding after two failures = %+v", got)
	}
	addRun(3, store.SyncRunStatusError)
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("three failures check error: %v", err)
	}
	addRun(4, store.SyncRunStatusError)
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("four failures check error: %v", err)
	}
	open := findings(t, s, connector.ID, "failing", "open")
	if len(open) != 1 || open[0].DetectedCount != 2 {
		t.Fatalf("finding after four failures = %+v", open)
	}
	addRun(5, store.SyncRunStatusSuccess)
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("interrupted streak check error: %v", err)
	}
	if got := findings(t, s, connector.ID, "failing", "open"); len(got) != 0 {
		t.Fatalf("finding after same-second success = %+v", got)
	}
}

func TestManualResolveReopensWhileConditionPersists(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	connector := createConnector(t, s, "")
	checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("first RunForConnector() error: %v", err)
	}
	first := findings(t, s, connector.ID, "ownership_incomplete", "open")[0]
	if err := s.UpdateQualityFindingStatus(ctx, first.ID, "resolved"); err != nil {
		t.Fatalf("manual resolve error: %v", err)
	}
	if err := checker.RunForConnector(ctx, connector.ID); err != nil {
		t.Fatalf("second RunForConnector() error: %v", err)
	}
	open := findings(t, s, connector.ID, "ownership_incomplete", "open")
	if len(open) != 1 || open[0].ID == first.ID || open[0].DetectedCount != 1 {
		t.Fatalf("reopened finding = %+v, previous ID %s", open, first.ID)
	}
}

// snapshotChangingNotifier simulates a sync completing between rule evaluations.
type snapshotChangingNotifier struct {
	change func()
}

func (n snapshotChangingNotifier) NotifyFindingCreated(context.Context, string, string, string) {
	n.change()
}

func TestComplianceRulesShareSnapshotWithinRun(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	conn := createConnector(t, s, "owner")
	createComplianceSnapshot(t, s, conn.ID, false)
	createComplianceRule(t, s, "proxmox", "First rule")
	createComplianceRule(t, s, "proxmox", "Second rule")
	changed := false
	notifier := snapshotChangingNotifier{change: func() {
		if !changed {
			createComplianceSnapshot(t, s, conn.ID, true)
			changed = true
		}
	}}
	checker := NewChecker(s, nil, notifier, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatal(err)
	}
	if got := findings(t, s, conn.ID, "compliance", "open"); len(got) != 2 {
		t.Fatalf("open findings = %d, want both rules evaluated against the same snapshot", len(got))
	}
	if err := checker.RunForConnector(ctx, conn.ID); err != nil {
		t.Fatal(err)
	}
	if got := findings(t, s, conn.ID, "compliance", "open"); len(got) != 0 {
		t.Fatalf("open findings = %d, want next run to use the newer snapshot", len(got))
	}
}

func TestComplianceSnapshotUnavailablePreservesFindings(t *testing.T) {
	for _, tc := range []struct {
		name string
		data string
		want int
	}{
		{name: "missing", want: 2},
		{name: "malformed", data: "{", want: 2},
		{name: "empty", data: `{"entities":[]}`, want: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := newTestStore(t)
			conn := createConnector(t, s, "owner")
			rules := []*store.ComplianceRuleRecord{
				createComplianceRule(t, s, "proxmox", "First rule"),
				createComplianceRule(t, s, "proxmox", "Second rule"),
			}
			for _, rule := range rules {
				if err := s.UpsertQualityFinding(ctx, &store.QualityFindingRecord{
					ID: rule.ID, ConnectorID: conn.ID, RuleID: rule.ID,
					CheckType: "compliance", Severity: "critical", Title: rule.Title,
				}); err != nil {
					t.Fatal(err)
				}
			}
			if tc.data != "" {
				if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{ConnectorID: conn.ID, Data: tc.data}); err != nil {
					t.Fatal(err)
				}
			}
			checker := NewChecker(s, nil, nil, RotationConfig{MaxAgeDays: 90, WarnDays: 14})
			if err := checker.RunForConnector(ctx, conn.ID); err != nil {
				t.Fatal(err)
			}
			if got := findings(t, s, conn.ID, "compliance", "open"); len(got) != tc.want {
				t.Fatalf("open findings after connector check = %d, want %d", len(got), tc.want)
			}
			for _, rule := range rules {
				if err := checker.EvaluateRule(ctx, rule.ID); err != nil {
					t.Fatal(err)
				}
			}
			if got := findings(t, s, conn.ID, "compliance", "open"); len(got) != tc.want {
				t.Fatalf("open findings after rule check = %d, want %d", len(got), tc.want)
			}
		})
	}
}
