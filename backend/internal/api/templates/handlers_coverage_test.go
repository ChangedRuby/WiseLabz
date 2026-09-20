package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func templateRequest(t *testing.T, fn http.HandlerFunc, id, rev, body string, want int) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.SetPathValue("id", id)
	req.SetPathValue("rev", rev)
	rr := httptest.NewRecorder()
	fn(rr, req)
	if rr.Code != want {
		t.Fatalf("status=%d, want %d: %s", rr.Code, want, rr.Body.String())
	}
	return rr
}

func TestVersionLifecycle(t *testing.T) {
	h := newTestHandler(t)
	created := templateRequest(t, h.Create, "", "", `{"name":"Original","description":"before","appliesTo":{"type":"custom"},"sections":[{"title":"Overview","order":2,"body":"Original body"}]}`, 201)
	var initial struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(created.Body.Bytes(), &initial); err != nil {
		t.Fatal(err)
	}
	if initial.ID == "" {
		t.Fatal("missing template ID")
	}
	templateRequest(t, h.Update, initial.ID, "", `{"name":"Updated","description":"after","appliesTo":{"category":"networking"},"sections":[{"title":"Changed","order":1,"body":"Updated body"}]}`, 200)
	version := templateRequest(t, h.Version, initial.ID, "1", "", 200)
	var snapshot struct {
		Name      string                         `json:"name"`
		Sections  []store.TemplateVersionSection `json:"sections"`
		AppliesTo map[string]string              `json:"appliesTo"`
	}
	if err := json.Unmarshal(version.Body.Bytes(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.Name != "Original" || len(snapshot.Sections) != 1 || snapshot.Sections[0].Body != "Original body" || snapshot.AppliesTo["type"] != "custom" {
		t.Fatalf("original version changed: %+v", snapshot)
	}
	templateRequest(t, h.Restore, initial.ID, "1", "", 200)
	live, err := h.Store.GetTemplate(context.Background(), initial.ID)
	if err != nil {
		t.Fatal(err)
	}
	sections, err := h.Store.GetTemplateSections(context.Background(), initial.ID)
	if err != nil {
		t.Fatal(err)
	}
	if live.Name != "Original" || live.Description != "before" || live.CurrentVersion != 3 || len(sections) != 1 || sections[0].Body != "Original body" {
		t.Fatalf("restore: %+v, sections=%+v", live, sections)
	}
	versions := templateRequest(t, h.Versions, initial.ID, "", "", 200)
	var history []map[string]any
	if err := json.Unmarshal(versions.Body.Bytes(), &history); err != nil {
		t.Fatal(err)
	}
	if len(history) != 3 || history[0]["trigger"] != "restore" || history[0]["rev"] != float64(3) {
		t.Fatalf("history=%s", versions.Body.String())
	}
	if _, ok := history[0]["sections"]; ok {
		t.Fatal("list leaked full snapshot")
	}
	audits, _, err := h.Store.ListAuditRecords(context.Background(), "template.restore", "template", "", "", 0, 10)
	if err != nil || len(audits) != 1 || audits[0].TargetID != initial.ID {
		t.Fatalf("audit=%+v, err=%v", audits, err)
	}
}

func TestListPagination(t *testing.T) {
	h := newTestHandler(t)
	for i := 1; i <= 3; i++ {
		rec := &store.TemplateRecord{Name: fmt.Sprintf("template-%d", i), CreatedAt: fmt.Sprintf("2025-01-0%dT00:00:00Z", i)}
		if err := h.Store.CreateTemplate(context.Background(), rec); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		query string
		names []string
	}{
		{"?page=2&pageSize=1", []string{"template-2"}},
		{"?page=4&pageSize=1", []string{}},
		{"?page=bad&pageSize=-1", []string{"template-3", "template-2", "template-1"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h.List(rr, httptest.NewRequest("GET", "/"+tc.query, nil))
			var rows []store.TemplateRecord
			if rr.Code != 200 {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
			if err := json.Unmarshal(rr.Body.Bytes(), &rows); err != nil {
				t.Fatal(err)
			}
			if rows == nil || len(rows) != len(tc.names) {
				t.Fatalf("rows=%s", rr.Body.String())
			}
			for i, name := range tc.names {
				if rows[i].Name != name {
					t.Fatalf("row %d=%q, want %q", i, rows[i].Name, name)
				}
			}
		})
	}
}

func TestTemplateErrorPaths(t *testing.T) {
	h := newTestHandler(t)
	for _, tc := range []struct {
		name      string
		fn        http.HandlerFunc
		rev, body string
		status    int
	}{
		{"update malformed", h.Update, "", "{", 400},
		{"update missing", h.Update, "", `{"name":"new"}`, 404},
		{"version malformed", h.Version, "bad", "", 400},
		{"version missing", h.Version, "1", "", 404},
		{"restore malformed", h.Restore, "bad", "", 400},
		{"restore missing", h.Restore, "1", "", 404},
		{"preview malformed", h.Preview, "", "{", 400},
	} {
		t.Run(tc.name, func(t *testing.T) { templateRequest(t, tc.fn, "missing", tc.rev, tc.body, tc.status) })
	}
	for name, fn := range map[string]http.HandlerFunc{"list": h.List, "get": h.Get, "create": h.Create, "update": h.Update, "delete": h.Delete, "versions": h.Versions, "version": h.Version, "restore": h.Restore, "preview": h.Preview} {
		t.Run(name+" canceled store", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`)).WithContext(ctx)
			req.SetPathValue("id", "missing")
			req.SetPathValue("rev", "1")
			rr := httptest.NewRecorder()
			fn(rr, req)
			if rr.Code != 500 {
				t.Fatalf("status=%d: %s", rr.Code, rr.Body.String())
			}
		})
	}
}

func TestPreviewDoesNotPersist(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()
	tmpl := &store.TemplateRecord{Name: "Preview"}
	if err := h.Store.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatal(err)
	}
	if err := h.Store.CreateTemplateSection(ctx, &store.TemplateSectionRecord{TemplateID: tmpl.ID, Title: "Overview", Body: "{{.ServiceName}}"}); err != nil {
		t.Fatal(err)
	}
	conn := &store.ConnectorRecord{Name: "Service", Type: "custom", Category: "networking", URL: "https://example.com", ConfigData: "{}"}
	if err := h.Store.CreateConnector(ctx, conn); err != nil {
		t.Fatal(err)
	}
	// A matching connector without a snapshot produces a per-connector render error.
	rr := templateRequest(t, h.Preview, tmpl.ID, "", "", 200)
	if !strings.Contains(rr.Body.String(), `"renderError":"`) {
		t.Fatalf("missing render error: %s", rr.Body.String())
	}
	if err := h.Store.CreateSnapshot(ctx, &store.SnapshotRecord{ConnectorID: conn.ID, Data: `{"serviceName":"Service","type":"custom","sections":[]}`, FetchedAt: "2025-01-01T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	for _, existing := range []bool{false, true} {
		if existing {
			if _, err := h.DocEngine.GenerateFromTemplate(ctx, tmpl.ID, conn.ID); err != nil {
				t.Fatal(err)
			}
		}
		rr = templateRequest(t, h.Preview, tmpl.ID, "", fmt.Sprintf(`{"connectorId":%q}`, conn.ID), 200)
		var result struct {
			Affected []struct {
				HasExistingDoc bool    `json:"hasExistingDoc"`
				WouldChange    bool    `json:"wouldChange"`
				RenderError    *string `json:"renderError"`
			} `json:"affected"`
			Detail struct {
				Content string `json:"content"`
			} `json:"detail"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if len(result.Affected) != 1 || result.Affected[0].HasExistingDoc != existing || result.Affected[0].WouldChange == existing || result.Affected[0].RenderError != nil || !strings.Contains(result.Detail.Content, "Service") {
			t.Fatalf("preview=%s", rr.Body.String())
		}
		docs, err := h.Store.ListDocsByService(ctx, conn.ID)
		if err != nil {
			t.Fatal(err)
		}
		want := 0
		if existing {
			want = 1
		}
		if len(docs) != want {
			t.Fatalf("preview persisted docs: %d", len(docs))
		}
	}
	versions, err := h.Store.GetTemplateVersions(ctx, tmpl.ID)
	if err != nil || len(versions) != 0 {
		t.Fatalf("preview persisted versions: %+v, %v", versions, err)
	}
}
