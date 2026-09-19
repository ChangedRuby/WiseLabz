package docs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/api/settings"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	settingsH := settings.NewHandler(s, &config.Config{}, ai.NewRegistry())
	return NewHandler(s, doc.NewEngine(s), settingsH, ai.NewRegistry(), nil, nil)
}

// withGrant seeds a user with the given per-connector role on connectorID
// and returns req with that user in context, as auth.AuthMiddleware would.
func withGrant(t *testing.T, s *store.Store, req *http.Request, connectorID, role string) *http.Request {
	t.Helper()
	userID := apitest.NewUser(t, s, "viewer")
	apitest.GrantConnectorRole(t, s, userID, connectorID, role)
	return req.WithContext(auth.ContextWithUser(req.Context(), userID, false))
}

func TestListEmpty(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetRootIsSynthetic(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/root", nil)
	req.SetPathValue("id", "root")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"lab"`) {
		t.Errorf("expected synthetic lab doc, got: %s", rr.Body.String())
	}
}

func TestGetUnknownIDFallsBackToServicePlaceholder(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/unknown-connector", nil)
	req.SetPathValue("id", "unknown-connector")
	req = withGrant(t, h.Store, req, "unknown-connector", "viewer")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"service"`) {
		t.Errorf("expected service placeholder, got: %s", rr.Body.String())
	}
}

func TestByServiceNoDocsYet(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/service/conn-1", nil)
	req.SetPathValue("id", "conn-1")
	req = withGrant(t, h.Store, req, "conn-1", "viewer")
	rr := httptest.NewRecorder()
	h.ByService(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestSave(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/docs/x", strings.NewReader(`{`))
		req.SetPathValue("id", "x")
		rr := httptest.NewRecorder()
		h.Save(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/docs/missing", strings.NewReader(`{"content":"hi"}`))
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Save(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
		}
	})
}

func TestVersionsOfUnknownDoc(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/missing/versions", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Versions(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetLockNoneHeld(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/x/lock", nil)
	req.SetPathValue("id", "x")
	rr := httptest.NewRecorder()
	h.GetLock(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestTree(t *testing.T) {
	s := apitest.NewStore(t)
	ctx := context.Background()

	// Seed a connector and a doc
	conn := &store.ConnectorRecord{
		Name:     "Test Connector",
		Category: "virtualization",
		Type:     "test-type",
		URL:      "https://test.example.com",
	}
	if err := s.CreateConnector(ctx, conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	docRecord := &store.DocRecord{
		Title:     "Test Doc",
		Kind:      "service",
		ServiceID: conn.ID,
		Content:   "test content",
	}
	if err := s.CreateDoc(ctx, docRecord); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	h := NewHandler(s, doc.NewEngine(s), nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/tree", nil)
	req = withGrant(t, s, req, conn.ID, "viewer")
	rr := httptest.NewRecorder()
	h.Tree(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "Lab Documentation") {
		t.Errorf("expected root node with Lab Documentation title")
	}
	if !strings.Contains(rr.Body.String(), conn.Name) {
		t.Errorf("expected connector in tree: %s", rr.Body.String())
	}
}

func TestTreeEmpty(t *testing.T) {
	s := apitest.NewStore(t)

	h := NewHandler(s, doc.NewEngine(s), nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/tree", nil)
	rr := httptest.NewRecorder()
	h.Tree(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	if !strings.Contains(rr.Body.String(), "Lab Documentation") {
		t.Errorf("expected root node")
	}
}

func TestVersion(t *testing.T) {
	s := apitest.NewStore(t)
	ctx := context.Background()

	// Create a doc and versions
	docRecord := &store.DocRecord{
		Title:     "Test Doc",
		Kind:      "service",
		ServiceID: "test-service",
		Content:   "v1 content",
	}
	if err := s.CreateDoc(ctx, docRecord); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	if err := s.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID:   docRecord.ID,
		Rev:     1,
		Content: "v1 content",
		Trigger: "manual",
	}); err != nil {
		t.Fatalf("create version: %v", err)
	}

	h := NewHandler(s, doc.NewEngine(s), nil, nil, nil, nil)

	t.Run("existing version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/docs/"+docRecord.ID+"/versions/1", nil)
		req.SetPathValue("id", docRecord.ID)
		req.SetPathValue("rev", "1")
		rr := httptest.NewRecorder()
		h.Version(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "v1 content") {
			t.Errorf("expected version content in response")
		}
	})

	t.Run("nonexistent version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/docs/"+docRecord.ID+"/versions/999", nil)
		req.SetPathValue("id", docRecord.ID)
		req.SetPathValue("rev", "999")
		rr := httptest.NewRecorder()
		h.Version(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("invalid rev", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/docs/"+docRecord.ID+"/versions/invalid", nil)
		req.SetPathValue("id", docRecord.ID)
		req.SetPathValue("rev", "invalid")
		rr := httptest.NewRecorder()
		h.Version(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestRestore(t *testing.T) {
	s := apitest.NewStore(t)
	ctx := context.Background()

	// Create a doc with versions
	docRecord := &store.DocRecord{
		Title:     "Test Doc",
		Kind:      "service",
		ServiceID: "test-service",
		Content:   "v1 content",
	}
	if err := s.CreateDoc(ctx, docRecord); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	if err := s.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID:   docRecord.ID,
		Rev:     1,
		Content: "v1 content",
		Trigger: "manual",
	}); err != nil {
		t.Fatalf("create v1: %v", err)
	}

	if err := s.UpdateDoc(ctx, docRecord.ID, "v2 content", nil); err != nil {
		t.Fatalf("update to v2: %v", err)
	}

	docRecord, _ = s.GetDoc(ctx, docRecord.ID)
	if err := s.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID:   docRecord.ID,
		Rev:     docRecord.CurrentVersion,
		Content: "v2 content",
		Trigger: "manual",
	}); err != nil {
		t.Fatalf("create v2: %v", err)
	}

	h := NewHandler(s, doc.NewEngine(s), nil, nil, nil, nil)

	t.Run("restore existing version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/docs/"+docRecord.ID+"/versions/1/restore", nil)
		req.SetPathValue("id", docRecord.ID)
		req.SetPathValue("rev", "1")
		req = withGrant(t, s, req, docRecord.ServiceID, "operator")
		rr := httptest.NewRecorder()
		h.Restore(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}

		// Verify the doc was restored
		restored, _ := s.GetDoc(ctx, docRecord.ID)
		if restored.Content != "v1 content" {
			t.Errorf("content after restore = %q, want %q", restored.Content, "v1 content")
		}
	})

	t.Run("restore nonexistent version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/docs/"+docRecord.ID+"/versions/999/restore", nil)
		req.SetPathValue("id", docRecord.ID)
		req.SetPathValue("rev", "999")
		rr := httptest.NewRecorder()
		h.Restore(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("restore with invalid rev", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/docs/"+docRecord.ID+"/versions/invalid/restore", nil)
		req.SetPathValue("id", docRecord.ID)
		req.SetPathValue("rev", "invalid")
		rr := httptest.NewRecorder()
		h.Restore(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestGenerate(t *testing.T) {
	s := apitest.NewStore(t)
	ctx := context.Background()

	// Seed template and connector
	tmpl := &store.TemplateRecord{
		Name:        "Test Template",
		Description: "Template for testing",
		AppliesTo:   "{}",
	}
	if err := s.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("create template: %v", err)
	}

	if err := s.CreateTemplateSection(ctx, &store.TemplateSectionRecord{
		TemplateID: tmpl.ID,
		Title:      "Overview",
		Ord:        1,
		Body:       "{{.ServiceName}} - {{.Type}}",
	}); err != nil {
		t.Fatalf("create template section: %v", err)
	}

	conn := &store.ConnectorRecord{
		Name:     "Test Service",
		Category: "containers_paas",
		Type:     "test-connector",
		URL:      "https://test.example.com",
	}
	if err := s.CreateConnector(ctx, conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	// Create snapshot for connector
	snap := connector.ServiceSnapshot{
		ServiceName: conn.Name,
		Type:        conn.Type,
		Sections: []connector.SnapshotSection{
			{Title: "Status", Content: "operational"},
		},
		Metadata:  map[string]string{"region": "us-east-1"},
		FetchedAt: time.Now(),
	}
	snapData, _ := json.Marshal(snap)
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{
		ConnectorID: conn.ID,
		Data:        string(snapData),
	}); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}

	h := NewHandler(s, doc.NewEngine(s), nil, nil, nil, nil)

	t.Run("success", func(t *testing.T) {
		payload := strings.NewReader(`{"templateId":"` + tmpl.ID + `","connectorId":"` + conn.ID + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/docs/generate", payload)
		req = withGrant(t, s, req, conn.ID, "operator")
		rr := httptest.NewRecorder()
		h.Generate(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "Test Service") {
			t.Errorf("expected generated content with service name")
		}
	})

	t.Run("missing templateId", func(t *testing.T) {
		payload := strings.NewReader(`{"connectorId":"` + conn.ID + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/docs/generate", payload)
		rr := httptest.NewRecorder()
		h.Generate(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing connectorId", func(t *testing.T) {
		payload := strings.NewReader(`{"templateId":"` + tmpl.ID + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/docs/generate", payload)
		rr := httptest.NewRecorder()
		h.Generate(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		payload := strings.NewReader(`{invalid}`)
		req := httptest.NewRequest(http.MethodPost, "/api/docs/generate", payload)
		rr := httptest.NewRecorder()
		h.Generate(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestTemplateSchema(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/template-schema", nil)
	rr := httptest.NewRecorder()
	h.TemplateSchema(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var schema map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &schema); err != nil {
		t.Fatalf("unmarshal schema: %v", err)
	}

	if _, ok := schema["fields"]; !ok {
		t.Errorf("schema missing 'fields' key")
	}
	if _, ok := schema["functions"]; !ok {
		t.Errorf("schema missing 'functions' key")
	}

	// Verify expected functions are documented
	funcs, _ := schema["functions"].([]any)
	funcNames := make(map[string]bool)
	for _, f := range funcs {
		fm, _ := f.(map[string]any)
		if name, ok := fm["name"].(string); ok {
			funcNames[name] = true
		}
	}

	expectedFuncs := []string{"dateFormat", "truncate", "toJSON", "filterByTitle", "join"}
	for _, fname := range expectedFuncs {
		if !funcNames[fname] {
			t.Errorf("expected function %q in schema", fname)
		}
	}
}

func TestAISuggestInvalidJSON(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/docs/d1/ai-suggest", strings.NewReader("{not json"))
	req.SetPathValue("id", "d1")
	rr := httptest.NewRecorder()
	h.AISuggest(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
}
