package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

type templateBody struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	CurrentVersion int    `json:"currentVersion"`
	Sections       []struct {
		Title string `json:"title"`
		Order int    `json:"order"`
		Body  string `json:"body"`
	} `json:"sections"`
}

func seedTemplate(t *testing.T, app *testApp, token string, appliesTo any) templateBody {
	t.Helper()
	rec := app.req(t, http.MethodPost, "/api/templates", map[string]any{
		"name": "Service template", "description": "Template description", "appliesTo": appliesTo,
		"sections": []map[string]any{{"title": "Overview", "order": 0, "body": "Type: {{.Type}}"}},
	}, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create template status = %d, want 201; body = %s", rec.Code, rec.Body)
	}
	var body templateBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode template: %v", err)
	}
	return body
}

func seedPreviewConnector(t *testing.T, app *testApp, name, category, connectorType string, snapshot bool) string {
	t.Helper()
	connector := &store.ConnectorRecord{Name: name, Category: category, Type: connectorType}
	if err := app.Store.CreateConnector(context.Background(), connector); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	if snapshot {
		data, err := json.Marshal(map[string]any{
			"serviceName": name,
			"type":        connectorType,
			"sections":    []any{},
			"metadata":    map[string]string{},
			"fetchedAt":   "2026-09-05T12:00:00Z",
		})
		if err != nil {
			t.Fatalf("marshal snapshot: %v", err)
		}
		if err := app.Store.CreateSnapshot(context.Background(), &store.SnapshotRecord{
			ConnectorID: connector.ID, Data: string(data),
		}); err != nil {
			t.Fatalf("create snapshot: %v", err)
		}
	}
	return connector.ID
}

func TestTemplatesPreviewDoesNotCreateDoc(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	viewerID, viewerToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization", "type": "proxmox"})
	connectorID := seedPreviewConnector(t, app, "pve-1", "virtualization", "proxmox", true)
	app.connectorGrant(t, viewerID, connectorID, "viewer")
	app.connectorGrant(t, opID, connectorID, "operator")

	before, err := app.Store.CountDocs(context.Background())
	if err != nil {
		t.Fatalf("count docs before: %v", err)
	}
	rec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/preview",
		map[string]any{"connectorId": connectorID}, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	if opRec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/preview",
		map[string]any{"connectorId": connectorID}, opToken); opRec.Code != http.StatusOK {
		t.Fatalf("operator preview status = %d, want 200; body = %s", opRec.Code, opRec.Body)
	}
	after, err := app.Store.CountDocs(context.Background())
	if err != nil {
		t.Fatalf("count docs after: %v", err)
	}
	if after != before {
		t.Fatalf("doc count changed from %d to %d", before, after)
	}

	var body struct {
		Affected []any `json:"affected"`
		Detail   *struct {
			DocID   string `json:"docId"`
			Content string `json:"content"`
		} `json:"detail"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode preview: %v", err)
	}
	if len(body.Affected) != 1 || body.Detail == nil || body.Detail.DocID != "" || body.Detail.Content == "" {
		t.Fatalf("unexpected preview: %+v", body)
	}
}

func TestTemplatesPreviewAffectedConnectors(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	viewerID, viewerToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization", "type": "proxmox"})
	app.connectorGrant(t, viewerID, seedPreviewConnector(t, app, "pve-1", "virtualization", "proxmox", true), "viewer")
	app.connectorGrant(t, viewerID, seedPreviewConnector(t, app, "pve-2", "virtualization", "proxmox", true), "viewer")
	seedPreviewConnector(t, app, "docker-1", "containers_paas", "docker", true)

	rec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/preview", map[string]any{}, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Affected []struct {
			ConnectorID    string  `json:"connectorId"`
			HasExistingDoc bool    `json:"hasExistingDoc"`
			WouldChange    bool    `json:"wouldChange"`
			RenderError    *string `json:"renderError"`
		} `json:"affected"`
		Detail any `json:"detail"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode preview: %v", err)
	}
	if len(body.Affected) != 2 {
		t.Fatalf("affected length = %d, want 2: %s", len(body.Affected), rec.Body)
	}
	for _, item := range body.Affected {
		if item.HasExistingDoc || !item.WouldChange || item.RenderError != nil {
			t.Errorf("unexpected affected item: %+v", item)
		}
	}
	if body.Detail != nil {
		t.Errorf("detail = %#v, want null", body.Detail)
	}
}

func TestTemplatesPreviewCapturesMissingSnapshot(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	viewerID, viewerToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization"})
	connectorID := seedPreviewConnector(t, app, "pve-1", "virtualization", "proxmox", false)
	app.connectorGrant(t, viewerID, connectorID, "viewer")
	if err := app.Store.CreateDoc(context.Background(), &store.DocRecord{
		Title: "Existing", Kind: "service", ServiceID: connectorID, Content: "old",
	}); err != nil {
		t.Fatalf("create existing doc: %v", err)
	}

	rec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/preview", map[string]any{}, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Affected []struct {
			HasExistingDoc bool    `json:"hasExistingDoc"`
			RenderError    *string `json:"renderError"`
		} `json:"affected"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode preview: %v", err)
	}
	if len(body.Affected) != 1 || !body.Affected[0].HasExistingDoc || body.Affected[0].RenderError == nil {
		t.Fatalf("expected per-connector render error: %s", rec.Body)
	}
}

func TestTemplatesPreviewZeroMatches(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization"})
	seedPreviewConnector(t, app, "docker-1", "containers_paas", "docker", true)

	rec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/preview", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("preview status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Affected []any `json:"affected"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode preview: %v", err)
	}
	if len(body.Affected) != 0 {
		t.Fatalf("affected = %#v, want empty", body.Affected)
	}
}

func TestTemplatesUpdateCreatesVersion(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization"})

	rec := app.req(t, http.MethodPut, "/api/templates/"+template.ID, map[string]any{
		"name":     "Updated template",
		"sections": []map[string]any{{"title": "Résumé", "order": 0, "body": "ação 🚀"}},
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var updated templateBody
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	if updated.CurrentVersion != 2 {
		t.Fatalf("currentVersion = %d, want 2", updated.CurrentVersion)
	}
	versions, err := app.Store.GetTemplateVersions(context.Background(), template.ID)
	if err != nil {
		t.Fatalf("get versions: %v", err)
	}
	if len(versions) != 2 || versions[0].Rev != 2 {
		t.Fatalf("versions = %+v, want revs 2 and 1", versions)
	}
}

func TestTemplatesConcurrentUpdatesCreateDistinctVersions(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	template := seedTemplate(t, app, opToken, nil)

	start := make(chan struct{})
	statuses := make(chan int, 2)
	var wg sync.WaitGroup
	for _, name := range []string{"Concurrent A", "Concurrent B"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			statuses <- app.req(t, http.MethodPut, "/api/templates/"+template.ID,
				map[string]any{"name": name}, opToken).Code
		}()
	}
	close(start)
	wg.Wait()
	close(statuses)
	for status := range statuses {
		if status != http.StatusOK {
			t.Fatalf("concurrent update status = %d, want 200", status)
		}
	}

	versions, err := app.Store.GetTemplateVersions(context.Background(), template.ID)
	if err != nil {
		t.Fatalf("get versions: %v", err)
	}
	if len(versions) != 3 || versions[0].Rev != 3 || versions[1].Rev != 2 || versions[2].Rev != 1 {
		t.Fatalf("versions = %+v, want distinct revs 3, 2, 1", versions)
	}
}

func TestTemplatesVersionsListAndGetRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization"})

	list := app.req(t, http.MethodGet, "/api/templates/"+template.ID+"/versions", nil, viewerToken)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200; body = %s", list.Code, list.Body)
	}
	if opList := app.req(t, http.MethodGet, "/api/templates/"+template.ID+"/versions", nil, opToken); opList.Code != http.StatusOK {
		t.Fatalf("operator list status = %d, want 200; body = %s", opList.Code, opList.Body)
	}
	var metas []struct {
		Rev int `json:"rev"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &metas); err != nil {
		t.Fatalf("decode versions: %v", err)
	}
	if len(metas) != 1 || metas[0].Rev != 1 {
		t.Fatalf("versions = %+v, want rev 1", metas)
	}

	get := app.req(t, http.MethodGet, "/api/templates/"+template.ID+"/versions/1", nil, viewerToken)
	if get.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", get.Code, get.Body)
	}
	if opGet := app.req(t, http.MethodGet, "/api/templates/"+template.ID+"/versions/1", nil, opToken); opGet.Code != http.StatusOK {
		t.Fatalf("operator get status = %d, want 200; body = %s", opGet.Code, opGet.Body)
	}
	var version struct {
		Rev      int `json:"rev"`
		Sections []struct {
			Title string `json:"title"`
			Order int    `json:"order"`
			Body  string `json:"body"`
		} `json:"sections"`
	}
	if err := json.Unmarshal(get.Body.Bytes(), &version); err != nil {
		t.Fatalf("decode version: %v", err)
	}
	if version.Rev != 1 || len(version.Sections) != 1 || version.Sections[0].Order != 0 {
		t.Fatalf("unexpected version: %+v", version)
	}
}

func TestTemplatesRestoreCreatesNewVersionAndUpdatesLiveContent(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	operatorID, opToken := app.user(t, "operator")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization"})
	update := app.req(t, http.MethodPut, "/api/templates/"+template.ID, map[string]any{
		"name":     "Version two",
		"sections": []map[string]any{{"title": "Changed", "order": 0, "body": "v2"}},
	}, opToken)
	if update.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", update.Code, update.Body)
	}

	restore := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/versions/1/restore", nil, opToken)
	if restore.Code != http.StatusOK {
		t.Fatalf("restore status = %d, want 200; body = %s", restore.Code, restore.Body)
	}
	var restored templateBody
	if err := json.Unmarshal(restore.Body.Bytes(), &restored); err != nil {
		t.Fatalf("decode restore: %v", err)
	}
	if restored.CurrentVersion != 3 || restored.Name != "Service template" ||
		len(restored.Sections) != 1 || restored.Sections[0].Body != "Type: {{.Type}}" {
		t.Fatalf("unexpected restored template: %+v", restored)
	}
	versions, err := app.Store.GetTemplateVersions(context.Background(), template.ID)
	if err != nil {
		t.Fatalf("get versions: %v", err)
	}
	if len(versions) != 3 || versions[0].Rev != 3 || versions[0].Trigger != "restore" {
		t.Fatalf("versions = %+v, want restore rev 3", versions)
	}

	audits, _, err := app.Store.ListAuditRecords(context.Background(), "template.restore", "template", "", "", 0, 10)
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(audits) != 1 || audits[0].TargetID != template.ID || audits[0].ActorUserID != operatorID {
		t.Fatalf("unexpected audit records: %+v", audits)
	}
}

func TestTemplatesRestoreCurrentRevisionCreatesNewVersion(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	template := seedTemplate(t, app, opToken, nil)

	rec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/versions/1/restore", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("restore status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	versions, err := app.Store.GetTemplateVersions(context.Background(), template.ID)
	if err != nil {
		t.Fatalf("get versions: %v", err)
	}
	if len(versions) != 2 || versions[0].Rev != 2 || versions[0].Trigger != "restore" {
		t.Fatalf("versions = %+v, want new restore rev 2", versions)
	}
}

func TestTemplatesRestoreRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, nil)

	rec := app.req(t, http.MethodPost, "/api/templates/"+template.ID+"/versions/1/restore", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestTemplatesPreviewRequiresViewerGrant(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	viewerID, viewerToken := app.user(t, "viewer")
	_, noGrantToken := app.user(t, "viewer")
	template := seedTemplate(t, app, opToken, map[string]any{"category": "virtualization", "type": "proxmox"})
	a := seedPreviewConnector(t, app, "pve-a", "virtualization", "proxmox", true)
	b := seedPreviewConnector(t, app, "pve-b", "virtualization", "proxmox", true)
	app.connectorGrant(t, viewerID, a, "viewer")
	app.connectorGrant(t, opID, a, "operator")
	app.connectorGrant(t, opID, b, "operator")
	path := "/api/templates/" + template.ID + "/preview"

	type previewBody struct {
		Affected []struct {
			ConnectorID string `json:"connectorId"`
		} `json:"affected"`
		Detail any `json:"detail"`
	}
	preview := func(token string, connectorID string) (int, previewBody) {
		rec := app.req(t, http.MethodPost, path, map[string]any{"connectorId": connectorID}, token)
		var body previewBody
		if rec.Code == http.StatusOK {
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode preview: %v", err)
			}
		}
		return rec.Code, body
	}

	// No grants: empty affected, detail null, 404 for any connectorId.
	if code, body := preview(noGrantToken, ""); code != http.StatusOK || len(body.Affected) != 0 || body.Detail != nil {
		t.Errorf("grantless preview = %d %+v, want 200 with empty affected and null detail", code, body)
	}
	if code, _ := preview(noGrantToken, a); code != http.StatusNotFound {
		t.Errorf("grantless detail = %d, want 404", code)
	}

	// Grant on A only: B hidden from affected and detail; A visible.
	if code, body := preview(viewerToken, ""); code != http.StatusOK || len(body.Affected) != 1 || body.Affected[0].ConnectorID != a {
		t.Errorf("partial-grant preview = %d %+v, want only %s", code, body, a)
	}
	if code, _ := preview(viewerToken, b); code != http.StatusNotFound {
		t.Errorf("detail for ungranted connector = %d, want 404", code)
	}
	if code, _ := preview(viewerToken, "does-not-exist"); code != http.StatusNotFound {
		t.Errorf("detail for nonexistent connector = %d, want 404", code)
	}
	if code, body := preview(viewerToken, a); code != http.StatusOK || body.Detail == nil {
		t.Errorf("detail for granted connector = %d %+v, want 200 with detail", code, body)
	}

	// Connector-restricted API key: its owner can view both, the key only A.
	key := mintAPIKey(t, app, opID, "", []string{a})
	if code, body := preview(key, ""); code != http.StatusOK || len(body.Affected) != 1 || body.Affected[0].ConnectorID != a {
		t.Errorf("restricted key preview = %d %+v, want only %s", code, body, a)
	}
	if code, _ := preview(key, b); code != http.StatusNotFound {
		t.Errorf("restricted key detail outside restriction = %d, want 404", code)
	}
}
