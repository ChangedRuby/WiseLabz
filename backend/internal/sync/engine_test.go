package sync

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
	_ "modernc.org/sqlite"
)

// fakeNotifier records NotifyAlertCreated calls.
type fakeNotifier struct {
	mu    sync.Mutex
	calls []string // alert IDs
}

type fakeQualityChecker struct {
	mu    sync.Mutex
	calls []string
	err   error
	order *[]string
}

func (f *fakeQualityChecker) RunForConnector(_ context.Context, connectorID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, connectorID)
	if f.order != nil {
		*f.order = append(*f.order, "qualityChecker")
	}
	return f.err
}

func (f *fakeQualityChecker) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

// fakeDocRegenerator records RegenerateForConnector calls, and optionally
// appends to a shared order log (used to assert call ordering relative to
// the quality checker).
type fakeDocRegenerator struct {
	mu    sync.Mutex
	calls []string
	err   error
	order *[]string
}

func (f *fakeDocRegenerator) RegenerateForConnector(_ context.Context, connectorID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, connectorID)
	if f.order != nil {
		*f.order = append(*f.order, "docRegenerator")
	}
	return f.err
}

func (f *fakeDocRegenerator) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func (f *fakeNotifier) NotifyAlertsCreated(_ context.Context, alerts []store.AlertRecord) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, alert := range alerts {
		f.calls = append(f.calls, alert.ID)
	}
}

func (f *fakeNotifier) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

// fakeConnector returns a fixed snapshot regardless of config.
type fakeConnector struct {
	snapshot *connector.ServiceSnapshot
}

func (f *fakeConnector) Name() string     { return "fake" }
func (f *fakeConnector) Type() string     { return "sync_test_fake" }
func (f *fakeConnector) Category() string { return "test" }
func (f *fakeConnector) Fetch(_ context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	return f.snapshot, nil
}
func (f *fakeConnector) Validate(_ context.Context, _ map[string]any) error { return nil }

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

// TestRunSync_NotifiesOnEligibleAlert verifies that RunSync calls the
// AlertNotifier exactly once per newly created (non-info) alert, and not at
// all when the only detected changes are info-severity.
func TestRunSync_NotifiesOnEligibleAlert(t *testing.T) {
	nextSnapshot := &connector.ServiceSnapshot{
		ServiceName: "svc",
		Sections:    []connector.SnapshotSection{}, // section removed vs prev -> "warning" alert
		FetchedAt:   time.Now(),
	}
	connector.Register(
		connector.TypeSchema{Type: "sync_test_fake_a", Category: "test", Name: "Fake A"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: nextSnapshot}, nil
		},
	)

	s := newTestStore(t)
	ctx := context.Background()

	conn := &store.ConnectorRecord{
		Name: "svc", Category: "networking", Type: "sync_test_fake_a", Enabled: true,
	}
	if err := s.CreateConnector(ctx, conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	prevSnapshot := `{"serviceName":"svc","sections":[{"title":"Ports","content":"22,80"}]}`
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{
		ConnectorID: conn.ID, Data: prevSnapshot, FetchedAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	notifier := &fakeNotifier{}
	engine := NewEngine(s, nil, notifier, nil, "")

	result, err := engine.RunSync(ctx, conn.ID, "job1")
	if err != nil {
		t.Fatalf("RunSync: %v", err)
	}
	if result.AlertsCount != 1 {
		t.Fatalf("expected 1 alert, got %d", result.AlertsCount)
	}
	if got := notifier.count(); got != 1 {
		t.Fatalf("expected notifier called once, got %d", got)
	}
}

// TestRunSync_NoNotifyOnInfoOnlyChange verifies info-severity changes (e.g. a
// newly added section) do not create an alert or notify, matching the
// existing "non-info" eligibility rule in RunSync.
func TestRunSync_NoNotifyOnInfoOnlyChange(t *testing.T) {
	nextSnapshot := &connector.ServiceSnapshot{
		ServiceName: "svc",
		Sections:    []connector.SnapshotSection{{Title: "Ports", Content: "22,80"}, {Title: "New", Content: "x"}},
		FetchedAt:   time.Now(),
	}
	connector.Register(
		connector.TypeSchema{Type: "sync_test_fake_b", Category: "test", Name: "Fake B"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: nextSnapshot}, nil
		},
	)

	s := newTestStore(t)
	ctx := context.Background()

	conn := &store.ConnectorRecord{
		Name: "svc", Category: "networking", Type: "sync_test_fake_b", Enabled: true,
	}
	if err := s.CreateConnector(ctx, conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	prevSnapshot := `{"serviceName":"svc","sections":[{"title":"Ports","content":"22,80"}]}`
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{
		ConnectorID: conn.ID, Data: prevSnapshot, FetchedAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	notifier := &fakeNotifier{}
	engine := NewEngine(s, nil, notifier, nil, "")

	result, err := engine.RunSync(ctx, conn.ID, "job1")
	if err != nil {
		t.Fatalf("RunSync: %v", err)
	}
	if result.AlertsCount != 0 {
		t.Fatalf("expected 0 alerts for info-only change, got %d", result.AlertsCount)
	}
	if got := notifier.count(); got != 0 {
		t.Fatalf("expected notifier not called, got %d", got)
	}
}

func TestRunSyncInvokesQualityChecker(t *testing.T) {
	connector.Register(
		connector.TypeSchema{Type: "sync_test_quality", Category: "test", Name: "Quality"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}}, nil
		},
	)

	tests := []struct {
		name       string
		connector  store.ConnectorRecord
		wantCalls  int
		wantRunErr bool
	}{
		{name: "success", connector: store.ConnectorRecord{Name: "success", Category: "networking", Type: "sync_test_quality", Enabled: true}, wantCalls: 1},
		{name: "error", connector: store.ConnectorRecord{Name: "error", Category: "networking", Type: "missing_quality_connector", Enabled: true}, wantCalls: 1, wantRunErr: true},
		{name: "skipped", connector: store.ConnectorRecord{Name: "skipped", Category: "networking", Type: "sync_test_quality", Enabled: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStore(t)
			ctx := context.Background()
			if err := s.CreateConnector(ctx, &tt.connector); err != nil {
				t.Fatalf("CreateConnector: %v", err)
			}

			checker := &fakeQualityChecker{}
			_, err := NewEngine(s, nil, nil, checker, "").RunSync(ctx, tt.connector.ID, "quality-job")
			if (err != nil) != tt.wantRunErr {
				t.Fatalf("RunSync error = %v, want error %v", err, tt.wantRunErr)
			}
			if got := checker.count(); got != tt.wantCalls {
				t.Fatalf("quality checker calls = %d, want %d", got, tt.wantCalls)
			}
		})
	}
}

func TestRunSyncQualityCheckerErrorIsNonFatal(t *testing.T) {
	connector.Register(
		connector.TypeSchema{Type: "sync_test_quality_error", Category: "test", Name: "Quality error"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}}, nil
		},
	)
	s := newTestStore(t)
	record := &store.ConnectorRecord{Name: "quality error", Category: "networking", Type: "sync_test_quality_error", Enabled: true}
	if err := s.CreateConnector(context.Background(), record); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	checker := &fakeQualityChecker{err: errors.New("quality unavailable")}
	result, err := NewEngine(s, nil, nil, checker, "").RunSync(context.Background(), record.ID, "quality-error-job")
	if err != nil || result.Status != "success" {
		t.Fatalf("RunSync() = (%+v, %v), want successful sync", result, err)
	}
	if checker.count() != 1 {
		t.Fatalf("quality checker calls = %d, want 1", checker.count())
	}
}

func TestRunSyncInvokesDocRegeneratorAfterQualityChecker(t *testing.T) {
	connector.Register(
		connector.TypeSchema{Type: "sync_test_doc_regen", Category: "test", Name: "Doc regen"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}}, nil
		},
	)
	s := newTestStore(t)
	record := &store.ConnectorRecord{Name: "doc regen", Category: "networking", Type: "sync_test_doc_regen", Enabled: true}
	if err := s.CreateConnector(context.Background(), record); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}

	var order []string
	checker := &fakeQualityChecker{order: &order}
	regenerator := &fakeDocRegenerator{order: &order}
	engine := NewEngine(s, nil, nil, checker, "")
	engine.SetDocRegenerator(regenerator)

	result, err := engine.RunSync(context.Background(), record.ID, "doc-regen-job")
	if err != nil || result.Status != "success" {
		t.Fatalf("RunSync() = (%+v, %v), want successful sync", result, err)
	}
	if regenerator.count() != 1 {
		t.Fatalf("doc regenerator calls = %d, want 1", regenerator.count())
	}
	if len(order) != 2 || order[0] != "qualityChecker" || order[1] != "docRegenerator" {
		t.Fatalf("call order = %v, want [qualityChecker docRegenerator]", order)
	}
}

func TestRunSyncDocRegeneratorErrorIsNonFatal(t *testing.T) {
	connector.Register(
		connector.TypeSchema{Type: "sync_test_doc_regen_error", Category: "test", Name: "Doc regen error"},
		func(_ map[string]any) (connector.Connector, error) {
			return &fakeConnector{snapshot: &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}}, nil
		},
	)
	s := newTestStore(t)
	record := &store.ConnectorRecord{Name: "doc regen error", Category: "networking", Type: "sync_test_doc_regen_error", Enabled: true}
	if err := s.CreateConnector(context.Background(), record); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}

	regenerator := &fakeDocRegenerator{err: errors.New("regeneration unavailable")}
	engine := NewEngine(s, nil, nil, nil, "")
	engine.SetDocRegenerator(regenerator)

	result, err := engine.RunSync(context.Background(), record.ID, "doc-regen-error-job")
	if err != nil || result.Status != "success" {
		t.Fatalf("RunSync() = (%+v, %v), want successful sync", result, err)
	}
	if regenerator.count() != 1 {
		t.Fatalf("doc regenerator calls = %d, want 1", regenerator.count())
	}
}

// A failure in a later alert batch must also roll back earlier batches and
// the snapshot, so a retry can detect the same drift and notify only once.
func TestRunSyncBatchRollbackAndRetry(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	next := &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}
	connector.Register(connector.TypeSchema{Type: "sync_batch_rollback", Category: "test", Name: "Batch"}, func(_ map[string]any) (connector.Connector, error) {
		return &fakeConnector{snapshot: next}, nil
	})
	rec := &store.ConnectorRecord{Name: "svc", Category: "networking", Type: "sync_batch_rollback", Enabled: true}
	if err := s.CreateConnector(ctx, rec); err != nil {
		t.Fatal(err)
	}
	prev := connector.ServiceSnapshot{ServiceName: "svc"}
	for i := 0; i < 220; i++ {
		prev.Sections = append(prev.Sections, connector.SnapshotSection{Title: fmt.Sprintf("section-%03d", i), Content: "old"})
	}
	data, err := json.Marshal(prev)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := &store.SnapshotRecord{ConnectorID: rec.ID, Data: string(data), FetchedAt: time.Now().Add(-time.Hour).Format(time.RFC3339)}
	if err := s.CreateSnapshot(ctx, snapshot); err != nil {
		t.Fatal(err)
	}
	if _, err := s.DB().ExecContext(ctx, `CREATE TRIGGER reject_alert_batch BEFORE INSERT ON alerts WHEN (SELECT COUNT(*) FROM alerts) >= 120 BEGIN SELECT RAISE(ABORT, 'injected alert failure'); END`); err != nil {
		t.Fatal(err)
	}
	notifier := &fakeNotifier{}
	engine := NewEngine(s, nil, notifier, nil, "")
	result, err := engine.RunSync(ctx, rec.ID, "failed-batch")
	if err == nil || result.Status != "error" {
		t.Fatalf("expected failed sync, got %+v, %v", result, err)
	}
	if result.ChangesCount != 0 || result.AlertsCount != 0 || result.SnapshotID != "" || notifier.count() != 0 {
		t.Fatalf("failed writes reported or notified: %+v, notifications=%d", result, notifier.count())
	}
	for _, table := range []string{"changes", "alerts"} {
		var count int
		if err := s.DB().QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("%s retained %d rows after rollback", table, count)
		}
	}
	latest, err := s.GetLatestSnapshot(ctx, rec.ID)
	if err != nil || latest.ID != snapshot.ID {
		t.Fatalf("snapshot advanced after failure: %+v, %v", latest, err)
	}
	if _, err := s.DB().ExecContext(ctx, "DROP TRIGGER reject_alert_batch"); err != nil {
		t.Fatal(err)
	}
	result, err = engine.RunSync(ctx, rec.ID, "retry-batch")
	if err != nil {
		t.Fatal(err)
	}
	if result.ChangesCount != 220 || result.AlertsCount != 220 || notifier.count() != 220 {
		t.Fatalf("retry did not persist and notify all drift: %+v, notifications=%d", result, notifier.count())
	}
	alerts, total, err := s.ListAlerts(ctx, rec.ID, "", "", "", 0, 250)
	if err != nil || total != 220 {
		t.Fatalf("alerts: count=%d, err=%v", total, err)
	}
	for _, alert := range alerts {
		change, err := s.GetChange(ctx, alert.ChangeID)
		if err != nil || change.ServiceID != rec.ID || change.Severity != "warning" {
			t.Fatalf("invalid committed alert/change link: %+v, %v", change, err)
		}
	}
}
