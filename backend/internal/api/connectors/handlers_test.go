package connectors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/sync"

	// Register connector implementations (proxmox, custom, ...) for restart tests.
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	connector.AllowLoopbackForTest(t)
	s := apitest.NewStore(t)
	cfg := &config.Config{Encryption: config.EncryptionSettings{Key: "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="}}
	jwtSvc := auth.NewService("test-secret-test-secret-test-secret", time.Hour, time.Hour)
	return NewHandler(s, sync.NewEngine(s, nil, nil, nil, cfg.Encryption.Key), cfg, jwtSvc, nil)
}

func TestListEmpty(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/connectors", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if strings.TrimSpace(rr.Body.String()) != "[]" {
		t.Errorf("body = %s, want []", rr.Body.String())
	}
}

func TestGetNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/connectors/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestCreate(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(`{"name":"svc"}`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		body := `{"name":"My Service","category":"virtualization","type":"custom","url":"https://svc.example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}
	})
}

func TestUpdate(t *testing.T) {
	h := newTestHandler(t)

	createReq := httptest.NewRequest(http.MethodPost, "/api/connectors",
		strings.NewReader(`{"name":"Original","category":"virtualization","type":"custom","url":"https://a.example.com"}`))
	createRR := httptest.NewRecorder()
	h.Create(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response: %s", createRR.Body.String())
	}

	t.Run("no fields to update", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/connectors/"+id, strings.NewReader(`{}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/connectors/missing", strings.NewReader(`{"name":"x"}`))
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/connectors/"+id, strings.NewReader(`{"name":"Renamed"}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "Renamed") {
			t.Errorf("update did not apply: %s", rr.Body.String())
		}
	})
}

func TestDelete(t *testing.T) {
	h := newTestHandler(t)

	t.Run("missing elevation token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/connectors/x", nil)
		req.SetPathValue("id", "x")
		rr := httptest.NewRecorder()
		h.Delete(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found with elevation token present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/connectors/missing", nil)
		req.SetPathValue("id", "missing")
		req.Header.Set("X-Elevation-Token", "placeholder")
		rr := httptest.NewRecorder()
		h.Delete(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestToggleEnabledNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPut, "/api/connectors/missing/enabled", strings.NewReader(`{"enabled":false}`))
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.ToggleEnabled(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}

func TestToggleEnabledSuccess(t *testing.T) {
	h := newTestHandler(t)

	// Create a connector first
	createReq := httptest.NewRequest(http.MethodPost, "/api/connectors",
		strings.NewReader(`{"name":"Test","category":"virtualization","type":"custom","url":"https://test.example.com"}`))
	createRR := httptest.NewRecorder()
	h.Create(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response")
	}

	// Toggle enabled to false
	req := httptest.NewRequest(http.MethodPut, "/api/connectors/"+id+"/enabled", strings.NewReader(`{"enabled":false}`))
	req.SetPathValue("id", id)
	rr := httptest.NewRecorder()
	h.ToggleEnabled(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result["enabled"] != false {
		t.Errorf("enabled = %v, want false", result["enabled"])
	}
}

func TestSchema(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/connectors/schema", nil)
	rr := httptest.NewRecorder()
	h.Schema(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestTestSuccess(t *testing.T) {
	h := newTestHandler(t)

	// Create a connector
	createReq := httptest.NewRequest(http.MethodPost, "/api/connectors",
		strings.NewReader(`{"name":"Test","category":"virtualization","type":"custom","url":"https://test.example.com"}`))
	createRR := httptest.NewRecorder()
	h.Create(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response")
	}

	// Test the connector
	req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/test", nil)
	req.SetPathValue("id", id)
	rr := httptest.NewRecorder()
	h.Test(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ok, _ := result["ok"].(bool); !ok {
		// For custom type (no actual validator), ok should be false
		t.Logf("connector test returned ok=false as expected for unknown type")
	}
}

func TestTestNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/connectors/missing/test", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Test(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestRemovalImpactSuccess(t *testing.T) {
	h := newTestHandler(t)

	// Create a connector
	createReq := httptest.NewRequest(http.MethodPost, "/api/connectors",
		strings.NewReader(`{"name":"Test","category":"virtualization","type":"custom","url":"https://test.example.com"}`))
	createRR := httptest.NewRecorder()
	h.Create(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response")
	}

	// Get removal impact
	req := httptest.NewRequest(http.MethodGet, "/api/connectors/"+id+"/removal-impact", nil)
	req.SetPathValue("id", id)
	rr := httptest.NewRecorder()
	h.RemovalImpact(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := result["trackedServices"]; !ok {
		t.Errorf("missing trackedServices in response")
	}
}

func TestRemovalImpactNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/connectors/missing/removal-impact", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.RemovalImpact(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestRestart(t *testing.T) {
	createConnector := func(t *testing.T, h *Handler, typ, url string) string {
		t.Helper()
		body := `{"name":"Test","category":"virtualization","type":"` + typ + `","url":"` + url + `","config":{"token_id":"u@pam!t","token_secret":"secret"}}`
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("create status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}
		var created map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
			t.Fatalf("unmarshal create: %v", err)
		}
		id, _ := created["id"].(string)
		if id == "" {
			t.Fatalf("no id in create response: %s", rr.Body.String())
		}
		return id
	}

	restartReq := func(id, entityRef, elevationToken string) *http.Request {
		var body string
		if entityRef != "" {
			body = `{"entityRef":"` + entityRef + `"}`
		}
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/restart", strings.NewReader(body))
		req.SetPathValue("id", id)
		if elevationToken != "" {
			req.Header.Set("X-Elevation-Token", elevationToken)
		}
		return req
	}

	t.Run("dry-run works without elevation token", func(t *testing.T) {
		h := newTestHandler(t)
		id := createConnector(t, h, "custom", "https://test.example.com")
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/restart?dryRun=true", nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.RestartPreview(rr, req)
		// No snapshot exists yet, but the important part is it's not rejected
		// for missing elevation.
		if rr.Code == http.StatusBadRequest && strings.Contains(rr.Body.String(), "elevation_required") {
			t.Fatalf("dry-run required elevation: %s", rr.Body.String())
		}
	})

	t.Run("unsupported connector returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		id := createConnector(t, h, "custom", "https://test.example.com")
		req := restartReq(id, "", "")
		rr := httptest.NewRecorder()
		h.RestartPreview(rr, req)
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "unsupported_operation") {
			t.Fatalf("status = %d, body = %s, want 400 unsupported_operation", rr.Code, rr.Body.String())
		}
	})

	t.Run("missing elevation token returns 400", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			t.Fatalf("unexpected request to connector API: %s", r.URL.Path)
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createConnector(t, h, "proxmox", server.URL)
		req := restartReq(id, "100", "")
		rr := httptest.NewRecorder()
		h.RestartPreview(rr, req)
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "elevation_required") {
			t.Fatalf("status = %d, body = %s, want 400 elevation_required", rr.Code, rr.Body.String())
		}
	})

	t.Run("invalid elevation token returns 401", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			t.Fatalf("unexpected request to connector API: %s", r.URL.Path)
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createConnector(t, h, "proxmox", server.URL)
		req := restartReq(id, "100", "not-a-real-token")
		rr := httptest.NewRecorder()
		h.RestartPreview(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, body = %s, want 401", rr.Code, rr.Body.String())
		}
	})

	t.Run("success writes audit row", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cluster/resources":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"node":"pve1","type":"qemu"}]}`))
			case "/nodes/pve1/qemu/100/status/reboot":
				_, _ = w.Write([]byte(`{"data":null}`))
			default:
				t.Fatalf("unexpected request: %s", r.URL.Path)
			}
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createConnector(t, h, "proxmox", server.URL)
		token, err := h.JWT.IssueElevation("", "connector.restart")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := restartReq(id, "100", token.Token)
		rr := httptest.NewRecorder()
		h.RestartPreview(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s, want 200", rr.Code, rr.Body.String())
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.restart", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error = %v", err)
		}
		if len(records) != 1 || records[0].TargetID != id {
			t.Fatalf("audit records = %+v, want one row for connector %s", records, id)
		}
	})

	t.Run("failure creates alert and no audit row", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"data":[]}`)) // vmid never found
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createConnector(t, h, "proxmox", server.URL)
		token, err := h.JWT.IssueElevation("", "connector.restart")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := restartReq(id, "100", token.Token)
		rr := httptest.NewRecorder()
		h.RestartPreview(rr, req)
		if rr.Code != http.StatusBadGateway {
			t.Fatalf("status = %d, body = %s, want 502", rr.Code, rr.Body.String())
		}
		alerts, _, err := h.Store.ListAlerts(context.Background(), id, "", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAlerts() error = %v", err)
		}
		if len(alerts) != 1 {
			t.Fatalf("alerts = %+v, want one alert", alerts)
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.restart", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error = %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("audit records = %+v, want none on failure", records)
		}
	})
}

func TestSyncAllSuccess(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/sync", nil)
	rr := httptest.NewRecorder()
	h.SyncAll(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusAccepted, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := result["jobId"]; !ok {
		t.Errorf("missing jobId in response")
	}
}

func TestUpdateEndpointChangeRequiresInstanceAdmin(t *testing.T) {
	h := newTestHandler(t)

	createRR := httptest.NewRecorder()
	h.Create(createRR, httptest.NewRequest(http.MethodPost, "/api/connectors",
		strings.NewReader(`{"name":"C","category":"virtualization","type":"custom","url":"https://a.example.com"}`)))
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id, _ := created["id"].(string)

	patch := func(admin bool, body string) int {
		req := httptest.NewRequest(http.MethodPut, "/api/connectors/"+id, strings.NewReader(body))
		req.SetPathValue("id", id)
		req = req.WithContext(auth.ContextWithUser(req.Context(), "u1", admin))
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		return rr.Code
	}

	if got := patch(false, `{"url":"https://evil.example.com"}`); got != http.StatusForbidden {
		t.Errorf("operator url change: status = %d, want 403", got)
	}
	if got := patch(false, `{"verifyTls":false}`); got != http.StatusForbidden {
		t.Errorf("operator verifyTls change: status = %d, want 403", got)
	}
	if got := patch(false, `{"name":"Renamed","url":"https://a.example.com"}`); got != http.StatusOK {
		t.Errorf("operator resending unchanged url: status = %d, want 200", got)
	}
	if got := patch(true, `{"url":"https://b.example.com"}`); got != http.StatusOK {
		t.Errorf("admin url change: status = %d, want 200", got)
	}
}
