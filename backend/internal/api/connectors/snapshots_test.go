package connectors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func snapshotFixture(t *testing.T) (*Handler, string, string, string) {
	t.Helper()
	h := newTestHandler(t)
	conn := &store.ConnectorRecord{Name: "Snapshot test", Type: "custom", Category: "networking", URL: "https://example.com"}
	if err := h.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatal(err)
	}
	makeSnapshot := func(at, content string) string {
		t.Helper()
		data, err := json.Marshal(map[string]any{
			"serviceName": "Snapshot test", "type": "custom", "fetchedAt": at,
			"sections":     []map[string]string{{"title": "State", "content": content}},
			"entities":     []map[string]any{{"kind": "host", "name": "server", "ip": "10.0.0.1", "attributes": map[string]any{"ready": true}}},
			"dependencies": []map[string]string{{"kind": "network", "name": "lan", "ref": "n1"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		record := &store.SnapshotRecord{ConnectorID: conn.ID, Data: string(data), FetchedAt: at}
		if err := h.Store.CreateSnapshot(context.Background(), record); err != nil {
			t.Fatal(err)
		}
		return record.ID
	}
	return h, conn.ID, makeSnapshot("2026-01-01T00:00:00Z", "old"), makeSnapshot("2026-01-02T00:00:00Z", "new")
}

func snapshotRequest(id, user, path string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	r.SetPathValue("id", id)
	return r.WithContext(auth.ContextWithUser(r.Context(), user, false))
}

func snapshotResponse(t *testing.T, handler http.Handler, r *http.Request, want int) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != want {
		t.Fatalf("%s status %d, want %d: %s", r.URL, w.Code, want, w.Body.String())
	}
	return w
}

func TestSnapshotsViewerAndCursor(t *testing.T) {
	h, id, _, _ := snapshotFixture(t)
	user := apitest.NewUser(t, h.Store, "viewer")
	path := "/api/connectors/" + id + "/snapshots?cursor=&limit=1"
	endpoint := auth.RequireConnectorRole(h.Store, "viewer", "id")(http.HandlerFunc(h.Snapshots))
	snapshotResponse(t, endpoint, snapshotRequest(id, user, path), http.StatusForbidden)
	apitest.GrantConnectorRole(t, h.Store, user, id, "viewer")
	page := snapshotResponse(t, endpoint, snapshotRequest(id, user, path), http.StatusOK)
	var rows []store.SnapshotSummary
	if err := json.Unmarshal(page.Body.Bytes(), &rows); err != nil || len(rows) != 1 || rows[0].SizeBytes == 0 || rows[0].Golden {
		t.Fatalf("first page: %+v, %v", rows, err)
	}
	cursor := page.Header().Get(httputil.NextCursorHeader)
	if cursor == "" {
		t.Fatal("missing next cursor")
	}
	legacy := snapshotResponse(t, endpoint, snapshotRequest(id, user, "/api/connectors/"+id+"/snapshots?limit=1"), http.StatusOK)
	if legacy.Header().Get(httputil.NextCursorHeader) != "" {
		t.Fatal("unexpected cursor without cursor parameter")
	}
	page = snapshotResponse(t, endpoint, snapshotRequest(id, user, "/api/connectors/"+id+"/snapshots?cursor="+cursor+"&limit=1"), http.StatusOK)
	if err := json.Unmarshal(page.Body.Bytes(), &rows); err != nil || len(rows) != 1 {
		t.Fatalf("next page: %+v, %v", rows, err)
	}
	// A full last page may include a cursor; the subsequent page is empty.
}

func TestSnapshotOwnershipAndFullShape(t *testing.T) {
	h, id, from, _ := snapshotFixture(t)
	other := &store.ConnectorRecord{Name: "Other", Type: "custom", Category: "networking", URL: "https://example.com"}
	if err := h.Store.CreateConnector(context.Background(), other); err != nil {
		t.Fatal(err)
	}
	r := snapshotRequest(id, "viewer", "/api/connectors/"+id+"/snapshots/"+from)
	r.SetPathValue("snapshotId", from)
	page := snapshotResponse(t, http.HandlerFunc(h.Snapshot), r, http.StatusOK)
	var full map[string]any
	if err := json.Unmarshal(page.Body.Bytes(), &full); err != nil {
		t.Fatal(err)
	}
	if full["id"] != from || len(full["entities"].([]any)) != 1 || len(full["dependencies"].([]any)) != 1 || full["metadata"] == nil {
		t.Fatalf("full snapshot: %+v", full)
	}
	r.SetPathValue("id", other.ID)
	snapshotResponse(t, http.HandlerFunc(h.Snapshot), r, http.StatusNotFound)
}

func TestSnapshotDiffValidationOwnershipAndAudit(t *testing.T) {
	h, id, from, to := snapshotFixture(t)
	other := &store.ConnectorRecord{Name: "Other", Type: "custom", Category: "networking", URL: "https://example.com"}
	if err := h.Store.CreateConnector(context.Background(), other); err != nil {
		t.Fatal(err)
	}
	foreign := &store.SnapshotRecord{ConnectorID: other.ID, Data: `{"serviceName":"Other","type":"custom","sections":[],"fetchedAt":"2026-01-03T00:00:00Z"}`, FetchedAt: "2026-01-03T00:00:00Z"}
	if err := h.Store.CreateSnapshot(context.Background(), foreign); err != nil {
		t.Fatal(err)
	}
	user := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, user, id, "viewer")
	base := "/api/connectors/" + id + "/snapshots/diff?from=" + from + "&to=" + to
	endpoint := auth.RequireConnectorRole(h.Store, "viewer", "id")(http.HandlerFunc(h.SnapshotDiff))
	call := func(path string, want int) *httptest.ResponseRecorder {
		return snapshotResponse(t, endpoint, snapshotRequest(id, user, path), want)
	}
	bad := call(base+"&format=pdf", http.StatusBadRequest)
	if !strings.Contains(bad.Body.String(), `"invalid_format"`) || !strings.Contains(bad.Body.String(), `"details"`) {
		t.Fatalf("bad format: %s", bad.Body.String())
	}
	call(base+"&format=", http.StatusBadRequest)
	call("/api/connectors/"+id+"/snapshots/diff?to="+to, http.StatusBadRequest)
	call(base[:strings.Index(base, "&to=")]+"&to="+foreign.ID, http.StatusNotFound)
	plain := call(base, http.StatusOK)
	var diff map[string]any
	if err := json.Unmarshal(plain.Body.Bytes(), &diff); err != nil || diff["provenance"] == nil || diff["summary"] == nil {
		t.Fatalf("diff body: %s, %v", plain.Body.String(), err)
	}
	rows, _, err := h.Store.ListAuditRecords(context.Background(), "snapshot.diff.export", "connector", "", "", 0, 10)
	if err != nil || len(rows) != 0 {
		t.Fatalf("plain view audit: %+v, %v", rows, err)
	}
	for _, format := range []string{"json", "csv", "md", "html"} {
		export := call(base+"&format="+format, http.StatusOK)
		if !strings.HasPrefix(export.Header().Get("Content-Disposition"), `attachment; filename="wiselabz-snapshot-diff-`+id) {
			t.Fatalf("%s disposition: %s", format, export.Header().Get("Content-Disposition"))
		}
	}
	rows, _, err = h.Store.ListAuditRecords(context.Background(), "snapshot.diff.export", "connector", "", "", 0, 10)
	if err != nil || len(rows) != 4 || rows[0].ActorUserID != user || rows[0].TargetID != id {
		t.Fatalf("export audit: %+v, %v", rows, err)
	}
}
