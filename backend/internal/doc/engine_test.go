package doc

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

func newEngineTestStore(t *testing.T) *store.Store {
	t.Helper()

	db, err := sql.Open("sqlite", "file:"+t.TempDir()+"/test.db?cache=shared")
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

func seedEngineTemplate(
	t *testing.T,
	s *store.Store,
	appliesTo string,
	sections ...store.TemplateSectionRecord,
) string {
	t.Helper()
	ctx := context.Background()
	tmpl := &store.TemplateRecord{
		Name:        "Service template",
		Description: "Generated from the latest snapshot",
		AppliesTo:   appliesTo,
	}
	if err := s.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("create template: %v", err)
	}
	for i := range sections {
		sections[i].TemplateID = tmpl.ID
		if err := s.CreateTemplateSection(ctx, &sections[i]); err != nil {
			t.Fatalf("create template section: %v", err)
		}
	}
	return tmpl.ID
}

func seedEngineConnector(
	t *testing.T,
	s *store.Store,
	name string,
	category string,
	connectorType string,
	withSnapshot bool,
) string {
	t.Helper()
	ctx := context.Background()
	record := &store.ConnectorRecord{
		Name:     name,
		Category: category,
		Type:     connectorType,
		URL:      "https://example.test",
	}
	if err := s.CreateConnector(ctx, record); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	if !withSnapshot {
		return record.ID
	}

	data, err := json.Marshal(connector.ServiceSnapshot{
		ServiceName: name,
		Type:        connectorType,
		Sections: []connector.SnapshotSection{
			{Title: "Status", Content: "healthy"},
		},
		Metadata:  map[string]string{"region": "lab"},
		FetchedAt: time.Date(2026, time.September, 5, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{ConnectorID: record.ID, Data: string(data)}); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}
	return record.ID
}

func TestPreviewFromTemplateDoesNotPersist(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, `{}`, store.TemplateSectionRecord{
		Title: "Overview",
		Ord:   1,
		Body:  `{{.ServiceName}} is {{.Type}} in {{index .Metadata "region"}}.`,
	})
	connectorID := seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", true)

	result, err := NewEngine(s).PreviewFromTemplate(ctx, templateID, connectorID)
	if err != nil {
		t.Fatalf("PreviewFromTemplate() error: %v", err)
	}
	if result.DocID != "" {
		t.Fatalf("PreviewFromTemplate() DocID = %q, want empty", result.DocID)
	}
	if !strings.Contains(result.Content, "Node one is proxmox in lab.") {
		t.Fatalf("PreviewFromTemplate() content = %q, want rendered snapshot", result.Content)
	}

	docs, err := s.ListDocsByService(ctx, connectorID)
	if err != nil {
		t.Fatalf("ListDocsByService() error: %v", err)
	}
	if len(docs) != 0 {
		t.Fatalf("ListDocsByService() returned %d docs after preview, want 0", len(docs))
	}
}

func TestGenerateFromTemplateStillPersists(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, `{}`, store.TemplateSectionRecord{
		Title: "Overview",
		Ord:   1,
		Body:  `{{.ServiceName}} is healthy.`,
	})
	connectorID := seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", true)
	engine := NewEngine(s)

	first, err := engine.GenerateFromTemplate(ctx, templateID, connectorID)
	if err != nil {
		t.Fatalf("first GenerateFromTemplate() error: %v", err)
	}
	if first.DocID == "" {
		t.Fatal("first GenerateFromTemplate() DocID is empty")
	}

	second, err := engine.GenerateFromTemplate(ctx, templateID, connectorID)
	if err != nil {
		t.Fatalf("second GenerateFromTemplate() error: %v", err)
	}
	if second.DocID != first.DocID {
		t.Fatalf("second GenerateFromTemplate() DocID = %q, want %q", second.DocID, first.DocID)
	}

	docs, err := s.ListDocsByService(ctx, connectorID)
	if err != nil {
		t.Fatalf("ListDocsByService() error: %v", err)
	}
	if len(docs) != 1 || docs[0].CurrentVersion != 2 {
		t.Fatalf("persisted docs = %+v, want one doc at version 2", docs)
	}
	versions, err := s.GetDocVersions(ctx, first.DocID)
	if err != nil {
		t.Fatalf("GetDocVersions() error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("GetDocVersions() returned %d versions, want 2", len(versions))
	}
}

func TestGenerateFromTemplateReturnsVersionPersistenceError(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, `{}`, store.TemplateSectionRecord{
		Title: "Overview",
		Ord:   1,
		Body:  `{{.ServiceName}} is healthy.`,
	})
	connectorID := seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", true)
	engine := NewEngine(s)

	first, err := engine.GenerateFromTemplate(ctx, templateID, connectorID)
	if err != nil {
		t.Fatalf("first GenerateFromTemplate() error: %v", err)
	}
	if err := s.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID:   first.DocID,
		Rev:     2,
		Content: "reserved revision",
		Trigger: "template",
	}); err != nil {
		t.Fatalf("reserve doc version: %v", err)
	}

	_, err = engine.GenerateFromTemplate(ctx, templateID, connectorID)
	if err == nil || !strings.Contains(err.Error(), "create doc version") {
		t.Fatalf("second GenerateFromTemplate() error = %v, want create doc version error", err)
	}
	doc, err := s.GetDoc(ctx, first.DocID)
	if err != nil {
		t.Fatalf("GetDoc() after failed generation: %v", err)
	}
	if doc.CurrentVersion != 1 {
		t.Fatalf("CurrentVersion after failed generation = %d, want 1", doc.CurrentVersion)
	}
}

func TestMatchingConnectorsFiltersByCategoryAndType(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	matchingID := seedEngineConnector(t, s, "Matching", "virtualization", "proxmox", false)
	wrongCategoryID := seedEngineConnector(t, s, "Wrong category", "networking", "proxmox", false)
	wrongTypeID := seedEngineConnector(t, s, "Wrong type", "virtualization", "libvirt", false)
	engine := NewEngine(s)

	tests := []struct {
		name      string
		appliesTo string
		wantIDs   map[string]bool
	}{
		{
			name:      "category and type",
			appliesTo: `{"category":"virtualization","type":"proxmox"}`,
			wantIDs:   map[string]bool{matchingID: true},
		},
		{
			name:      "category wildcard",
			appliesTo: `{"type":"proxmox"}`,
			wantIDs:   map[string]bool{matchingID: true, wrongCategoryID: true},
		},
		{
			name:      "type wildcard",
			appliesTo: `{"category":"virtualization"}`,
			wantIDs:   map[string]bool{matchingID: true, wrongTypeID: true},
		},
		{
			name:      "all wildcard",
			appliesTo: `{}`,
			wantIDs:   map[string]bool{matchingID: true, wrongCategoryID: true, wrongTypeID: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateID := seedEngineTemplate(t, s, tt.appliesTo)
			got, err := engine.MatchingConnectors(ctx, templateID)
			if err != nil {
				t.Fatalf("MatchingConnectors() error: %v", err)
			}
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("MatchingConnectors() returned %d connectors, want %d", len(got), len(tt.wantIDs))
			}
			for _, candidate := range got {
				if !tt.wantIDs[candidate.ID] {
					t.Errorf("MatchingConnectors() unexpectedly returned %q", candidate.Name)
				}
			}
		})
	}
}

func TestMatchingConnectorsEmptyAppliesToIsWildcard(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, "")
	seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", false)
	seedEngineConnector(t, s, "Router", "networking", "pfsense", false)

	got, err := NewEngine(s).MatchingConnectors(ctx, templateID)
	if err != nil {
		t.Fatalf("MatchingConnectors() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("MatchingConnectors() returned %d connectors, want 2", len(got))
	}
}

func TestMatchingConnectorsRejectsMalformedAppliesTo(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, `{not-json}`)

	_, err := NewEngine(s).MatchingConnectors(ctx, templateID)
	if err == nil || !strings.Contains(err.Error(), "unmarshal template applies_to") {
		t.Fatalf("MatchingConnectors() error = %v, want malformed applies_to error", err)
	}
}

func TestPreviewFromTemplateWithZeroSections(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, `{}`)
	connectorID := seedEngineConnector(t, s, "Bare node", "virtualization", "proxmox", true)

	result, err := NewEngine(s).PreviewFromTemplate(ctx, templateID, connectorID)
	if err != nil {
		t.Fatalf("PreviewFromTemplate() error: %v", err)
	}
	want := "# Bare node\n\n> Generated from the latest snapshot\n\n"
	if result.Content != want {
		t.Fatalf("PreviewFromTemplate() content = %q, want %q", result.Content, want)
	}
}

func TestGenerateFromSnapshotIncludesDependencies(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	record := &store.ConnectorRecord{Name: "Node one", Category: "virtualization", Type: "proxmox", URL: "https://example.test"}
	if err := s.CreateConnector(ctx, record); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	data, err := json.Marshal(connector.ServiceSnapshot{
		ServiceName: "Node one",
		Type:        "proxmox",
		Sections: []connector.SnapshotSection{
			{Title: "Status", Content: "healthy"},
		},
		Dependencies: []connector.ServiceDependency{
			{Kind: "host", Name: "nodeA"},
			{Kind: "upstream_service", Name: "auth-api", Ref: "svc-123"},
		},
		FetchedAt: time.Date(2026, time.September, 5, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{ConnectorID: record.ID, Data: string(data)}); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}

	result, err := NewEngine(s).GenerateFromSnapshot(ctx, record.ID)
	if err != nil {
		t.Fatalf("GenerateFromSnapshot() error: %v", err)
	}
	if !strings.Contains(result.Content, "## Dependencies") {
		t.Fatalf("GenerateFromSnapshot() content = %q, want Dependencies heading", result.Content)
	}
	if !strings.Contains(result.Content, "- **host**: nodeA") {
		t.Fatalf("GenerateFromSnapshot() content = %q, want host dependency bullet", result.Content)
	}
	if !strings.Contains(result.Content, "- **upstream_service**: auth-api (svc-123)") {
		t.Fatalf("GenerateFromSnapshot() content = %q, want upstream_service dependency bullet with ref", result.Content)
	}
}

func TestRegenerateForConnectorNoopWhenContentUnchanged(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	connectorID := seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", true)
	engine := NewEngine(s)

	generated, err := engine.GenerateFromSnapshot(ctx, connectorID)
	if err != nil {
		t.Fatalf("GenerateFromSnapshot() error: %v", err)
	}

	if err := engine.RegenerateForConnector(ctx, connectorID); err != nil {
		t.Fatalf("RegenerateForConnector() error: %v", err)
	}

	doc, err := s.GetDoc(ctx, generated.DocID)
	if err != nil {
		t.Fatalf("GetDoc() error: %v", err)
	}
	if doc.CurrentVersion != 1 {
		t.Fatalf("CurrentVersion after no-op regenerate = %d, want 1 (no new version)", doc.CurrentVersion)
	}
	versions, err := s.GetDocVersions(ctx, generated.DocID)
	if err != nil {
		t.Fatalf("GetDocVersions() error: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("GetDocVersions() returned %d versions, want 1", len(versions))
	}
}

func TestRegenerateForConnectorCreatesSyncVersionWhenChanged(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	connectorID := seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", true)
	engine := NewEngine(s)

	generated, err := engine.GenerateFromSnapshot(ctx, connectorID)
	if err != nil {
		t.Fatalf("GenerateFromSnapshot() error: %v", err)
	}

	// Simulate the snapshot having changed since the doc was generated by
	// pushing a newer one for the same connector.
	data, err := json.Marshal(connector.ServiceSnapshot{
		ServiceName: "Node one",
		Type:        "proxmox",
		Sections: []connector.SnapshotSection{
			{Title: "Status", Content: "degraded"},
		},
		FetchedAt: time.Date(2026, time.September, 6, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{
		ConnectorID: connectorID,
		Data:        string(data),
		FetchedAt:   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}

	if err := engine.RegenerateForConnector(ctx, connectorID); err != nil {
		t.Fatalf("RegenerateForConnector() error: %v", err)
	}

	doc, err := s.GetDoc(ctx, generated.DocID)
	if err != nil {
		t.Fatalf("GetDoc() error: %v", err)
	}
	if doc.CurrentVersion != 2 {
		t.Fatalf("CurrentVersion after regenerate = %d, want 2", doc.CurrentVersion)
	}
	if !strings.Contains(doc.Content, "degraded") {
		t.Fatalf("doc content = %q, want updated content", doc.Content)
	}
	versions, err := s.GetDocVersions(ctx, generated.DocID)
	if err != nil {
		t.Fatalf("GetDocVersions() error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("GetDocVersions() returned %d versions, want 2", len(versions))
	}
	if versions[0].Trigger != "sync" {
		t.Fatalf("latest version trigger = %q, want %q", versions[0].Trigger, "sync")
	}
}

func TestRegenerateForConnectorNoopWhenNoExistingDocs(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	connectorID := seedEngineConnector(t, s, "Node one", "virtualization", "proxmox", true)

	if err := NewEngine(s).RegenerateForConnector(ctx, connectorID); err != nil {
		t.Fatalf("RegenerateForConnector() error: %v", err)
	}

	docs, err := s.ListDocsByService(ctx, connectorID)
	if err != nil {
		t.Fatalf("ListDocsByService() error: %v", err)
	}
	if len(docs) != 0 {
		t.Fatalf("ListDocsByService() returned %d docs, want 0", len(docs))
	}
}

func TestPreviewFromTemplateWithoutSnapshotReturnsRenderError(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	templateID := seedEngineTemplate(t, s, `{}`)
	connectorID := seedEngineConnector(t, s, "Unsynced node", "virtualization", "proxmox", false)

	_, err := NewEngine(s).PreviewFromTemplate(ctx, templateID, connectorID)
	if err == nil || !strings.Contains(err.Error(), "get snapshot") {
		t.Fatalf("PreviewFromTemplate() error = %v, want get snapshot error", err)
	}
}
