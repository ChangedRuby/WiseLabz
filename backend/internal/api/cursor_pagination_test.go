package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// cursorPage is the additive pagination envelope: nextCursor must only appear
// once a request has opted into keyset mode.
type cursorPage struct {
	Items      []map[string]any `json:"items"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	NextCursor string           `json:"nextCursor"`
}

func decodeCursorPage(t *testing.T, body []byte) cursorPage {
	t.Helper()
	var page cursorPage
	if err := json.Unmarshal(body, &page); err != nil {
		t.Fatalf("decode page: %v (body %s)", err, body)
	}
	return page
}

// walkCursorPages follows nextCursor from the first keyset page to the last and
// returns every item id it saw, failing on a duplicate or a non-terminating
// traversal.
func walkCursorPages(t *testing.T, app *testApp, token, path string, pageSize int) []string {
	t.Helper()

	var ids []string
	seen := map[string]bool{}
	cursor := ""
	for page := 0; ; page++ {
		if page > 50 {
			t.Fatalf("cursor traversal did not terminate after %d pages", page)
		}
		u := fmt.Sprintf("%s&pageSize=%d&cursor=%s", path, pageSize, url.QueryEscape(cursor))
		rec := app.req(t, http.MethodGet, u, nil, token)
		if rec.Code != http.StatusOK {
			t.Fatalf("page %d status = %d, want 200; body = %s", page, rec.Code, rec.Body)
		}
		got := decodeCursorPage(t, rec.Body.Bytes())
		for _, item := range got.Items {
			id, _ := item["id"].(string)
			if id == "" {
				t.Fatalf("page %d item has no id: %+v", page, item)
			}
			if seen[id] {
				t.Fatalf("id %s returned twice across cursor pages (page %d)", id, page)
			}
			seen[id] = true
			ids = append(ids, id)
		}
		if got.NextCursor == "" {
			return ids
		}
		cursor = got.NextCursor
	}
}

// TestAuditCursorPaginationTraversal walks /api/system/audit by cursor and
// checks the pages tile the full result set exactly: no duplicates, no gaps,
// including across runs of records that share created_at.
func TestAuditCursorPaginationTraversal(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	ctx := context.Background()

	// 13 records over 3 timestamps, so pages of 4 split ties.
	base := time.Now().UTC().Truncate(time.Second)
	for i := range 13 {
		if err := app.Store.CreateAuditRecord(ctx, &store.AuditRecord{
			ActorUserID: "u", Action: "connector.sync", TargetType: "connector",
			TargetID:  fmt.Sprintf("c-%02d", i),
			CreatedAt: base.Add(-time.Duration(i/5) * time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("CreateAuditRecord(%d): %v", i, err)
		}
	}

	ids := walkCursorPages(t, app, opToken, "/api/system/audit?action=connector.sync", 4)
	if len(ids) != 13 {
		t.Fatalf("cursor traversal visited %d records, want 13", len(ids))
	}

	// The same rows, in the same order, as one big keyset page.
	rec := app.req(t, http.MethodGet, "/api/system/audit?action=connector.sync&pageSize=100&cursor=", nil, opToken)
	all := decodeCursorPage(t, rec.Body.Bytes())
	if all.Total != 13 {
		t.Errorf("total = %d, want 13", all.Total)
	}
	if all.NextCursor != "" {
		t.Errorf("nextCursor = %q on the last page, want empty", all.NextCursor)
	}
	for i, item := range all.Items {
		if item["id"] != ids[i] {
			t.Errorf("row %d = %v, cursor traversal had %s", i, item["id"], ids[i])
		}
	}
}

// TestAuditOffsetPaginationUnchanged pins the backward-compatible half of the
// contract: a request that does not send cursor gets the historical envelope,
// with no nextCursor key at all.
func TestAuditOffsetPaginationUnchanged(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	ctx := context.Background()

	for i := range 5 {
		if err := app.Store.CreateAuditRecord(ctx, &store.AuditRecord{
			ActorUserID: "u", Action: "connector.sync", TargetType: "connector",
			TargetID: fmt.Sprintf("c-%d", i),
		}); err != nil {
			t.Fatalf("CreateAuditRecord(%d): %v", i, err)
		}
	}

	rec := app.req(t, http.MethodGet, "/api/system/audit?action=connector.sync&pageSize=2&page=2", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var raw map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if _, present := raw["nextCursor"]; present {
		t.Errorf("offset response must not carry nextCursor, got %v", raw)
	}
	for _, key := range []string{"items", "total", "page", "pageSize"} {
		if _, present := raw[key]; !present {
			t.Errorf("offset response missing %q; keys = %v", key, raw)
		}
	}
	if raw["page"] != float64(2) || raw["pageSize"] != float64(2) || raw["total"] != float64(5) {
		t.Errorf("envelope = %v, want page 2, pageSize 2, total 5", raw)
	}
}

func TestAuditRejectsMalformedCursor(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/system/audit?cursor=not-a-cursor", nil, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

// TestChangesCursorPaginationTraversal walks /api/changes by cursor. The
// handler filters the DB page by connector grant after fetching it, so this
// also pins that the cursor advances past the whole DB page rather than the
// filtered remainder.
func TestChangesCursorPaginationTraversal(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	ctx := context.Background()

	connector := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(ctx, connector); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	app.connectorGrant(t, userID, connector.ID, "viewer")

	base := time.Now().UTC().Truncate(time.Second)
	for i := range 11 {
		if err := app.Store.CreateChange(ctx, &store.ChangeRecord{
			ServiceID: connector.ID, ChangeType: "modified", Severity: "warning",
			Summary:    fmt.Sprintf("change %02d", i),
			DetectedAt: base.Add(-time.Duration(i/4) * time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("CreateChange(%d): %v", i, err)
		}
	}

	ids := walkCursorPages(t, app, token, "/api/changes?serviceId="+connector.ID, 3)
	if len(ids) != 11 {
		t.Fatalf("cursor traversal visited %d changes, want 11", len(ids))
	}
}

// TestSyncsCursorPaginationUsesHeader covers the one endpoint whose published
// body is a bare array: the cursor rides on X-Next-Cursor and the JSON stays a
// plain array.
func TestSyncsCursorPaginationUsesHeader(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	ctx := context.Background()

	connector := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(ctx, connector); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	app.connectorGrant(t, opUserID, connector.ID, "viewer")

	base := time.Now().UTC().Truncate(time.Second)
	for i := range 7 {
		if err := app.Store.CreateSyncRun(ctx, &store.SyncRunRecord{
			ConnectorID: connector.ID, Status: store.SyncRunStatusSuccess,
			StartedAt: base.Add(-time.Duration(i/3) * time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("CreateSyncRun(%d): %v", i, err)
		}
	}

	// Without cursor: unchanged bare array, no cursor header.
	rec := app.req(t, http.MethodGet, "/api/connectors/"+connector.ID+"/syncs?limit=3", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var runs []store.SyncRunRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &runs); err != nil {
		t.Fatalf("offset body is not a bare array: %v (%s)", err, rec.Body)
	}
	if got := rec.Header().Get(httputil.NextCursorHeader); got != "" {
		t.Errorf("%s = %q on an offset request, want empty", httputil.NextCursorHeader, got)
	}

	var ids []string
	seen := map[string]bool{}
	cursor := ""
	for page := 0; page < 10; page++ {
		rec := app.req(t, http.MethodGet,
			fmt.Sprintf("/api/connectors/%s/syncs?limit=3&cursor=%s", connector.ID, url.QueryEscape(cursor)), nil, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("page %d status = %d; body = %s", page, rec.Code, rec.Body)
		}
		var pageRuns []store.SyncRunRecord
		if err := json.Unmarshal(rec.Body.Bytes(), &pageRuns); err != nil {
			t.Fatalf("page %d body is not a bare array: %v (%s)", page, err, rec.Body)
		}
		for _, run := range pageRuns {
			if seen[run.ID] {
				t.Fatalf("sync run %s returned twice across cursor pages", run.ID)
			}
			seen[run.ID] = true
			ids = append(ids, run.ID)
		}
		cursor = rec.Header().Get(httputil.NextCursorHeader)
		if cursor == "" {
			break
		}
	}
	if len(ids) != 7 {
		t.Fatalf("cursor traversal visited %d sync runs, want 7", len(ids))
	}
}
