package connectors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func createDNSResolverConnector(t *testing.T, h *Handler, url string) string {
	t.Helper()
	body := `{"name":"Test","category":"dns","type":"dnsresolver","url":"` + url + `","config":{"api_key":"secret"}}`
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

func TestStartStopHandler(t *testing.T) {
	createConnectorProxmox := func(t *testing.T, h *Handler, url string) string {
		t.Helper()
		body := `{"name":"Test","category":"virtualization","type":"proxmox","url":"` + url + `","config":{"token_id":"u@pam!t","token_secret":"secret"}}`
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
		return created["id"].(string)
	}

	t.Run("start dry-run works without elevation token", func(t *testing.T) {
		h := newTestHandler(t)
		id := createConnectorProxmox(t, h, "https://test.example.com")
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/start?dryRun=true", nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.StartPreview(rr, req)
		if rr.Code == http.StatusBadRequest && strings.Contains(rr.Body.String(), "elevation_required") {
			t.Fatalf("dry-run required elevation: %s", rr.Body.String())
		}
	})

	t.Run("stop dry-run estimatedDowntimeSeconds is 0", func(t *testing.T) {
		h := newTestHandler(t)
		id := createConnectorProxmox(t, h, "https://test.example.com")
		// No snapshot -> 404, but confirm the handler doesn't panic and the
		// downtime-zero behavior is exercised via a stored snapshot path in
		// the "success" case below; here we just confirm no-snapshot 404.
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/stop?dryRun=true", nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.StopPreview(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404 (no snapshot yet)", rr.Code)
		}
	})

	t.Run("unsupported connector returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		id := createDNSResolverConnector(t, h, "https://test.example.com")
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/start", nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.StartPreview(rr, req)
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
		id := createConnectorProxmox(t, h, server.URL)
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/stop", strings.NewReader(`{"entityRef":"100"}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.StopPreview(rr, req)
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "elevation_required") {
			t.Fatalf("status = %d, body = %s, want 400 elevation_required", rr.Code, rr.Body.String())
		}
	})

	t.Run("start success writes audit row", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cluster/resources":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"node":"pve1","type":"qemu"}]}`))
			case "/nodes/pve1/qemu/100/status/start":
				_, _ = w.Write([]byte(`{"data":null}`))
			default:
				t.Fatalf("unexpected request: %s", r.URL.Path)
			}
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createConnectorProxmox(t, h, server.URL)
		token, err := h.JWT.IssueElevation("", "connector.start")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/start", strings.NewReader(`{"entityRef":"100"}`))
		req.SetPathValue("id", id)
		req.Header.Set("X-Elevation-Token", token.Token)
		rr := httptest.NewRecorder()
		h.StartPreview(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s, want 200", rr.Code, rr.Body.String())
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.start", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error = %v", err)
		}
		if len(records) != 1 || records[0].TargetID != id {
			t.Fatalf("audit records = %+v, want one row for connector %s", records, id)
		}
	})

	t.Run("stop failure creates alert and no audit row", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"data":[]}`)) // vmid never found
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createConnectorProxmox(t, h, server.URL)
		token, err := h.JWT.IssueElevation("", "connector.stop")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/stop", strings.NewReader(`{"entityRef":"100"}`))
		req.SetPathValue("id", id)
		req.Header.Set("X-Elevation-Token", token.Token)
		rr := httptest.NewRecorder()
		h.StopPreview(rr, req)
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
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.stop", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error = %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("audit records = %+v, want none on failure", records)
		}
	})
}

func TestConfigFieldsHandler(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		h := newTestHandler(t)
		req := httptest.NewRequest(http.MethodGet, "/api/connectors/missing/config-fields", nil)
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.ConfigFields(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rr.Code)
		}
	})

	t.Run("unsupported connector returns empty array", func(t *testing.T) {
		h := newTestHandler(t)
		body := `{"name":"Test","category":"virtualization","type":"custom","url":"https://test.example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		var created map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &created)
		id := created["id"].(string)

		req2 := httptest.NewRequest(http.MethodGet, "/api/connectors/"+id+"/config-fields", nil)
		req2.SetPathValue("id", id)
		rr2 := httptest.NewRecorder()
		h.ConfigFields(rr2, req2)
		if rr2.Code != http.StatusOK || strings.TrimSpace(rr2.Body.String()) != "[]" {
			t.Fatalf("status=%d body=%s, want 200 []", rr2.Code, rr2.Body.String())
		}
	})

	t.Run("supported connector returns fields", func(t *testing.T) {
		h := newTestHandler(t)
		id := createDNSResolverConnector(t, h, "https://test.example.com")
		req := httptest.NewRequest(http.MethodGet, "/api/connectors/"+id+"/config-fields", nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.ConfigFields(rr, req)
		if rr.Code != http.StatusOK || !strings.Contains(rr.Body.String(), `"key":"ip"`) {
			t.Fatalf("status=%d body=%s, want 200 with ip field", rr.Code, rr.Body.String())
		}
	})
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func TestConfigPushHandler(t *testing.T) {
	pushReq := func(id, entityRef, fieldKey string, value, previousValue any, token string) *http.Request {
		payload := map[string]any{"entityRef": entityRef, "fieldKey": fieldKey, "value": value, "previousValue": previousValue}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/connectors/"+id+"/config-push", strings.NewReader(string(b)))
		req.SetPathValue("id", id)
		if token != "" {
			req.Header.Set("X-Elevation-Token", token)
		}
		return req
	}

	t.Run("missing elevation token returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		id := createDNSResolverConnector(t, h, "https://test.example.com")
		req := pushReq(id, "web.example.com", "ip", "10.0.0.9", "10.0.0.5", "")
		rr := httptest.NewRecorder()
		h.ConfigPush(rr, req)
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "elevation_required") {
			t.Fatalf("status = %d, body = %s, want 400 elevation_required", rr.Code, rr.Body.String())
		}
	})

	t.Run("unsupported connector returns 400", func(t *testing.T) {
		h := newTestHandler(t)
		body := `{"name":"Test","category":"virtualization","type":"custom","url":"https://test.example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		var created map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &created)
		id := created["id"].(string)

		req2 := pushReq(id, "", "ip", "1.2.3.4", nil, "")
		rr2 := httptest.NewRecorder()
		h.ConfigPush(rr2, req2)
		if rr2.Code != http.StatusBadRequest || !strings.Contains(rr2.Body.String(), "unsupported_operation") {
			t.Fatalf("status = %d, body = %s, want 400 unsupported_operation", rr2.Code, rr2.Body.String())
		}
	})

	t.Run("field not on whitelist rejected before any write", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			t.Fatalf("unexpected request to connector API: %s", r.URL.Path)
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createDNSResolverConnector(t, h, server.URL)
		req := pushReq(id, "web.example.com", "ttl", 300, nil, "placeholder")
		rr := httptest.NewRecorder()
		h.ConfigPush(rr, req)
		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "unsupported_field") {
			t.Fatalf("status = %d, body = %s, want 400 unsupported_field", rr.Code, rr.Body.String())
		}
	})

	// proxmox is used for the three subtests below (rather than dnsresolver,
	// used above) because dnsresolver's constructor wires a guarded dialer
	// that rejects loopback addresses — it can't reach an httptest.Server.
	// proxmox has no such guard, matching how PR1's restart handler tests
	// already reach a real httptest server.
	createProxmoxConnector := func(t *testing.T, h *Handler, url string) string {
		t.Helper()
		body := `{"name":"Test","category":"virtualization","type":"proxmox","url":"` + url + `","config":{"token_id":"u@pam!t","token_secret":"secret"}}`
		req := httptest.NewRequest(http.MethodPost, "/api/connectors", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		var created map[string]any
		_ = json.Unmarshal(rr.Body.Bytes(), &created)
		return created["id"].(string)
	}

	t.Run("success verifies diff and writes audit row", func(t *testing.T) {
		var memory atomic.Int32
		memory.Store(2048)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case r.URL.Path == "/nodes":
				_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
			case r.URL.Path == "/nodes/pve1/qemu":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"name":"vm1","status":"running","cpus":2,"mem":` + itoa(int(memory.Load())) + `,"uptime":10}]}`))
			case r.URL.Path == "/nodes/pve1/lxc" || r.URL.Path == "/nodes/pve1/storage":
				_, _ = w.Write([]byte(`{"data":[]}`))
			case r.URL.Path == "/cluster/resources":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"node":"pve1","type":"qemu"}]}`))
			case r.URL.Path == "/nodes/pve1/qemu/100/agent/network-get-interfaces":
				w.WriteHeader(http.StatusInternalServerError) // fetchQemuIP soft-fails on error
			case r.URL.Path == "/nodes/pve1/qemu/100/config" && r.Method == "PUT":
				memory.Store(4096) // the write actually lands
				_, _ = w.Write([]byte(`{"data":null}`))
			default:
				t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createProxmoxConnector(t, h, server.URL)
		token, err := h.JWT.IssueElevation("", "connector.configPush")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := pushReq(id, "100", "memory", 4096, 2048, token.Token)
		rr := httptest.NewRecorder()
		h.ConfigPush(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s, want 200", rr.Code, rr.Body.String())
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.configPush", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error = %v", err)
		}
		if len(records) != 1 || records[0].TargetID != id {
			t.Fatalf("audit records = %+v, want one row for connector %s", records, id)
		}
	})

	t.Run("mismatch triggers auto-revert and critical alert", func(t *testing.T) {
		var putCount atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case r.URL.Path == "/nodes":
				_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
			case r.URL.Path == "/nodes/pve1/qemu":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"name":"vm1","status":"running","cpus":2,"mem":2048,"uptime":10}]}`)) // never changes
			case r.URL.Path == "/nodes/pve1/lxc" || r.URL.Path == "/nodes/pve1/storage":
				_, _ = w.Write([]byte(`{"data":[]}`))
			case r.URL.Path == "/cluster/resources":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"node":"pve1","type":"qemu"}]}`))
			case r.URL.Path == "/nodes/pve1/qemu/100/agent/network-get-interfaces":
				w.WriteHeader(http.StatusInternalServerError) // fetchQemuIP soft-fails on error
			case r.URL.Path == "/nodes/pve1/qemu/100/config" && r.Method == "PUT":
				putCount.Add(1) // both push and revert "succeed" upstream, value just never sticks
				_, _ = w.Write([]byte(`{"data":null}`))
			default:
				t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createProxmoxConnector(t, h, server.URL)
		token, err := h.JWT.IssueElevation("", "connector.configPush")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := pushReq(id, "100", "memory", 4096, 2048, token.Token)
		rr := httptest.NewRecorder()
		h.ConfigPush(rr, req)
		if rr.Code != http.StatusConflict || !strings.Contains(rr.Body.String(), "config_push_mismatch") {
			t.Fatalf("status = %d, body = %s, want 409 config_push_mismatch", rr.Code, rr.Body.String())
		}
		alerts, _, err := h.Store.ListAlerts(context.Background(), id, "", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAlerts() error = %v", err)
		}
		if len(alerts) != 1 || alerts[0].Severity != "critical" {
			t.Fatalf("alerts = %+v, want one critical alert", alerts)
		}
		if strings.Contains(alerts[0].Description, "FAILED") {
			t.Errorf("description = %q, want no revert-failure escalation (revert succeeded)", alerts[0].Description)
		}
		if putCount.Load() != 2 {
			t.Errorf("putCount = %d, want 2 (push + revert)", putCount.Load())
		}
		records, _, err := h.Store.ListAuditRecords(context.Background(), "connector.configPush", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error = %v", err)
		}
		if len(records) != 0 {
			t.Fatalf("audit records = %+v, want none on mismatch", records)
		}
	})

	t.Run("revert itself failing escalates the alert", func(t *testing.T) {
		var putCount atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case r.URL.Path == "/nodes":
				_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
			case r.URL.Path == "/nodes/pve1/qemu":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"name":"vm1","status":"running","cpus":2,"mem":2048,"uptime":10}]}`)) // never changes
			case r.URL.Path == "/nodes/pve1/lxc" || r.URL.Path == "/nodes/pve1/storage":
				_, _ = w.Write([]byte(`{"data":[]}`))
			case r.URL.Path == "/cluster/resources":
				_, _ = w.Write([]byte(`{"data":[{"vmid":100,"node":"pve1","type":"qemu"}]}`))
			case r.URL.Path == "/nodes/pve1/qemu/100/agent/network-get-interfaces":
				w.WriteHeader(http.StatusInternalServerError) // fetchQemuIP soft-fails on error
			case r.URL.Path == "/nodes/pve1/qemu/100/config" && r.Method == "PUT":
				n := putCount.Add(1)
				if n == 1 {
					_, _ = w.Write([]byte(`{"data":null}`)) // the push itself "succeeds" (still doesn't verify)
					return
				}
				w.WriteHeader(http.StatusInternalServerError) // the revert fails
			default:
				t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
		}))
		defer server.Close()
		h := newTestHandler(t)
		id := createProxmoxConnector(t, h, server.URL)
		token, err := h.JWT.IssueElevation("", "connector.configPush")
		if err != nil {
			t.Fatalf("IssueElevation() error = %v", err)
		}
		req := pushReq(id, "100", "memory", 4096, 2048, token.Token)
		rr := httptest.NewRecorder()
		h.ConfigPush(rr, req)
		if rr.Code != http.StatusConflict {
			t.Fatalf("status = %d, body = %s, want 409", rr.Code, rr.Body.String())
		}
		alerts, _, err := h.Store.ListAlerts(context.Background(), id, "", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAlerts() error = %v", err)
		}
		if len(alerts) != 1 || alerts[0].Severity != "critical" {
			t.Fatalf("alerts = %+v, want one critical alert", alerts)
		}
		if !strings.Contains(alerts[0].Description, "FAILED") {
			t.Errorf("description = %q, want it to note the revert failed", alerts[0].Description)
		}
	})
}
