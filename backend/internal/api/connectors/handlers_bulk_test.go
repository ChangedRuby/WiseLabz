package connectors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// bulkFakeConnector implements Connector, Restarter, and CredentialRefresher
// so the three bulk handlers can be exercised against one fake type. restart
// and refresh both fail for entity ref "broken", to exercise partial-success
// results.
type bulkFakeConnector struct {
	mu           sync.Mutex
	restartCalls []string
}

func (b *bulkFakeConnector) Name() string     { return "bulk-fake" }
func (b *bulkFakeConnector) Type() string     { return "bulk_fake" }
func (b *bulkFakeConnector) Category() string { return "test" }
func (b *bulkFakeConnector) Validate(_ context.Context, _ map[string]any) error {
	return nil
}
func (b *bulkFakeConnector) Fetch(_ context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	return &connector.ServiceSnapshot{ServiceName: "svc", FetchedAt: time.Now()}, nil
}
func (b *bulkFakeConnector) Restart(_ context.Context, config map[string]any, _ string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if config["url"] == "https://broken.example.com" {
		return &connector.AuthError{Err: context.DeadlineExceeded}
	}
	b.restartCalls = append(b.restartCalls, config["url"].(string))
	return nil
}
func (b *bulkFakeConnector) RefreshCredentials(_ context.Context, config map[string]any) (map[string]any, time.Time, error) {
	if config["url"] == "https://broken.example.com" {
		return nil, time.Time{}, context.DeadlineExceeded
	}
	refreshed := make(map[string]any, len(config)+1)
	for k, v := range config {
		refreshed[k] = v
	}
	refreshed["token"] = "refreshed"
	return refreshed, time.Now().Add(24 * time.Hour), nil
}

func registerBulkFakeConnector(t *testing.T) {
	t.Helper()
	connector.Register(
		connector.TypeSchema{Type: "bulk_fake", Category: "networking", Name: "Bulk Fake"},
		func(_ map[string]any) (connector.Connector, error) { return &bulkFakeConnector{}, nil },
	)
}

func createBulkFakeConnector(t *testing.T, h *Handler, url string) string {
	t.Helper()
	body := `{"name":"Test","category":"networking","type":"bulk_fake","url":"` + url + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
	rr := httptest.NewRecorder()
	h.Create(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
	}
	var created map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &created)
	return created["id"].(string)
}

func bulkResults(t *testing.T, rr *httptest.ResponseRecorder) []bulkItemResult {
	t.Helper()
	var body struct {
		Results []bulkItemResult `json:"results"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal results: %v; body=%s", err, rr.Body.String())
	}
	return body.Results
}

func bulkReq(method, path, body string) *http.Request {
	return httptest.NewRequest(method, path, strings.NewReader(body))
}

func TestBulkSync(t *testing.T) {
	registerBulkFakeConnector(t)

	t.Run("empty ids returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		rr := httptest.NewRecorder()
		h.BulkSync(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-sync", `{"ids":[]}`))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
		}
	})

	t.Run("over limit returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		ids := make([]string, 501)
		for i := range ids {
			ids[i] = "x"
		}
		b, _ := json.Marshal(map[string]any{"ids": ids})
		rr := httptest.NewRecorder()
		h.BulkSync(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-sync", string(b)))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
		}
	})

	t.Run("partial success: valid and missing ids", func(t *testing.T) {
		h := newTestHandler(t)
		id := createBulkFakeConnector(t, h, "https://svc.example.com")
		b, _ := json.Marshal(map[string]any{"ids": []string{id, "missing"}})
		rr := httptest.NewRecorder()
		h.BulkSync(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-sync", string(b)))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
		}
		results := bulkResults(t, rr)
		if len(results) != 2 {
			t.Fatalf("results = %+v, want 2 items", results)
		}
		byID := map[string]bulkItemResult{}
		for _, r := range results {
			byID[r.ID] = r
		}
		if byID[id].Status != "success" || byID[id].JobID == "" {
			t.Errorf("valid id result = %+v, want success with jobId", byID[id])
		}
		if byID["missing"].Status != "error" || byID["missing"].Reason != "not_found" {
			t.Errorf("missing id result = %+v, want error/not_found", byID["missing"])
		}

		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.bulk_sync", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords: %v", err)
		}
		if len(records) != 1 || records[0].TargetID != id {
			t.Fatalf("audit records = %+v, want exactly one row for %s", records, id)
		}
	})

	t.Run("all ids missing: all-fail", func(t *testing.T) {
		h := newTestHandler(t)
		b, _ := json.Marshal(map[string]any{"ids": []string{"a", "b"}})
		rr := httptest.NewRecorder()
		h.BulkSync(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-sync", string(b)))
		results := bulkResults(t, rr)
		for _, r := range results {
			if r.Status != "error" {
				t.Errorf("result = %+v, want error", r)
			}
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.bulk_sync", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords: %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("audit records = %+v, want none", records)
		}
	})
}

func TestBulkReauth(t *testing.T) {
	registerBulkFakeConnector(t)

	t.Run("empty ids returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		rr := httptest.NewRecorder()
		h.BulkReauth(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-reauth", `{"ids":[]}`))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
		}
	})

	t.Run("full success refreshes and persists credentials", func(t *testing.T) {
		h := newTestHandler(t)
		id := createBulkFakeConnector(t, h, "https://svc.example.com")
		b, _ := json.Marshal(map[string]any{"ids": []string{id}})
		rr := httptest.NewRecorder()
		h.BulkReauth(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-reauth", string(b)))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
		}
		results := bulkResults(t, rr)
		if len(results) != 1 || results[0].Status != "success" {
			t.Fatalf("results = %+v, want one success", results)
		}
		rec, err := h.Store.GetConnector(context.Background(), id)
		if err != nil {
			t.Fatalf("GetConnector: %v", err)
		}
		if !strings.Contains(rec.ConfigData, "refreshed") {
			t.Errorf("config_data = %s, want refreshed token persisted", rec.ConfigData)
		}
	})

	t.Run("partial success: one broken connector fails refresh", func(t *testing.T) {
		h := newTestHandler(t)
		ok := createBulkFakeConnector(t, h, "https://svc.example.com")
		broken := createBulkFakeConnector(t, h, "https://broken.example.com")
		b, _ := json.Marshal(map[string]any{"ids": []string{ok, broken}})
		rr := httptest.NewRecorder()
		h.BulkReauth(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-reauth", string(b)))
		results := bulkResults(t, rr)
		byID := map[string]bulkItemResult{}
		for _, r := range results {
			byID[r.ID] = r
		}
		if byID[ok].Status != "success" {
			t.Errorf("ok result = %+v, want success", byID[ok])
		}
		if byID[broken].Status != "error" {
			t.Errorf("broken result = %+v, want error", byID[broken])
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.bulk_reauth", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords: %v", err)
		}
		if len(records) != 1 || records[0].TargetID != ok {
			t.Fatalf("audit records = %+v, want exactly one row for %s", records, ok)
		}
	})
}

func TestBulkRestart(t *testing.T) {
	registerBulkFakeConnector(t)

	t.Run("empty ids returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		rr := httptest.NewRecorder()
		h.BulkRestart(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-restart", `{"ids":[]}`))
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
		}
	})

	// BulkRestart's elevation check is enforced by router middleware
	// (auth.RequireElevation), not inside the handler — see
	// TestConnectorsBulkRestartElevationBoundary in the api package for the
	// missing/invalid/wrong-action/expired-token coverage that requires the
	// full router.

	t.Run("partial success: one broken connector fails restart, one audit row per success", func(t *testing.T) {
		h := newTestHandler(t)
		ok := createBulkFakeConnector(t, h, "https://svc.example.com")
		broken := createBulkFakeConnector(t, h, "https://broken.example.com")
		b, _ := json.Marshal(map[string]any{"ids": []string{ok, broken, "missing"}})
		rr := httptest.NewRecorder()
		h.BulkRestart(rr, bulkReq(http.MethodPost, "/api/connectors/bulk-restart", string(b)))
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
		}
		results := bulkResults(t, rr)
		byID := map[string]bulkItemResult{}
		for _, r := range results {
			byID[r.ID] = r
		}
		if byID[ok].Status != "success" {
			t.Errorf("ok result = %+v, want success", byID[ok])
		}
		if byID[broken].Status != "error" {
			t.Errorf("broken result = %+v, want error", byID[broken])
		}
		if byID["missing"].Status != "error" || byID["missing"].Reason != "not_found" {
			t.Errorf("missing result = %+v, want error/not_found", byID["missing"])
		}

		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.bulk_restart", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords: %v", err)
		}
		if len(records) != 1 || records[0].TargetID != ok {
			t.Fatalf("audit records = %+v, want exactly one row for %s", records, ok)
		}

		alerts, _, err := h.Store.ListAlerts(context.Background(), broken, "", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAlerts: %v", err)
		}
		if len(alerts) != 1 {
			t.Fatalf("alerts for broken connector = %+v, want one failure alert", alerts)
		}
	})
}
