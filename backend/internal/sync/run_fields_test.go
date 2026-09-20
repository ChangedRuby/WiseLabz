package sync

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

type fieldsRecordingConnector struct {
	snapshot *connector.ServiceSnapshot
	config   map[string]any
}

func (c *fieldsRecordingConnector) Name() string     { return "fields recording" }
func (c *fieldsRecordingConnector) Type() string     { return "sync_test_fields_behavior" }
func (c *fieldsRecordingConnector) Category() string { return "test" }
func (c *fieldsRecordingConnector) Validate(_ context.Context, _ map[string]any) error {
	return nil
}
func (c *fieldsRecordingConnector) Fetch(_ context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	c.config = config
	return c.snapshot, nil
}

func TestRunSyncFieldsPersistsPartialFetchAndChange(t *testing.T) {
	fetched := &connector.ServiceSnapshot{
		ServiceName: "svc",
		Sections:    []connector.SnapshotSection{{Title: "VMs", Content: "vm-2"}},
		FetchedAt:   time.Now(),
	}
	recorder := &fieldsRecordingConnector{snapshot: fetched}
	connector.Register(
		connector.TypeSchema{Type: "sync_test_fields_behavior", Category: "test", Name: "Fields Behavior"},
		func(_ map[string]any) (connector.Connector, error) { return recorder, nil },
	)

	s := newTestStore(t)
	ctx := context.Background()
	rec := &store.ConnectorRecord{Name: "svc", Category: "networking", Type: "sync_test_fields_behavior", Enabled: true}
	if err := s.CreateConnector(ctx, rec); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{
		ConnectorID: rec.ID,
		Data:        `{"serviceName":"svc","sections":[{"title":"VMs","content":"vm-1"}]}`,
		FetchedAt:   time.Now().Add(-time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("CreateSnapshot: %v", err)
	}

	result, err := NewEngine(s, nil, nil, nil, "").RunSyncFields(ctx, rec.ID, "partial-job", []string{"vms"})
	if err != nil {
		t.Fatalf("RunSyncFields: %v", err)
	}
	if got := connector.RequestedFields(recorder.config); len(got) != 1 || got[0] != "vms" {
		t.Fatalf("requested fields = %v, want [vms]", got)
	}
	if result.Status != "success" || result.ChangesCount != 1 {
		t.Fatalf("result = %#v, want successful sync with one change", result)
	}

	latest, err := s.GetLatestSnapshot(ctx, rec.ID)
	if err != nil {
		t.Fatalf("GetLatestSnapshot: %v", err)
	}
	if !strings.Contains(latest.Data, `"vm-2"`) {
		t.Fatalf("latest snapshot = %s, want partial fetch data", latest.Data)
	}
	changes, _, err := s.ListChanges(ctx, rec.ID, "", 0, 10)
	if err != nil {
		t.Fatalf("ListChanges: %v", err)
	}
	if len(changes) != 1 || changes[0].Summary != "Section modified: VMs" {
		t.Fatalf("changes = %#v, want VMs modification", changes)
	}
}

func TestRunSyncFieldsWithoutHintDoesNotSetFieldsConfig(t *testing.T) {
	recorder := &fieldsRecordingConnector{snapshot: &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}}
	connector.Register(
		connector.TypeSchema{Type: "sync_test_fields_behavior_empty", Category: "test", Name: "Fields Behavior Empty"},
		func(_ map[string]any) (connector.Connector, error) { return recorder, nil },
	)

	s := newTestStore(t)
	rec := &store.ConnectorRecord{Name: "svc", Category: "networking", Type: "sync_test_fields_behavior_empty", Enabled: true}
	if err := s.CreateConnector(context.Background(), rec); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	if _, err := NewEngine(s, nil, nil, nil, "").RunSyncFields(context.Background(), rec.ID, "full-job", nil); err != nil {
		t.Fatalf("RunSyncFields: %v", err)
	}
	if _, ok := recorder.config["fields"]; ok {
		t.Fatalf("connector config = %v, want no fields hint", recorder.config)
	}
}
