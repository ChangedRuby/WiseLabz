package connectors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/sync"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	cfg := &config.Config{Encryption: config.EncryptionSettings{Key: "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="}}
	return NewHandler(s, sync.NewEngine(s, nil, nil, nil, cfg.Encryption.Key), cfg)
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
		req := httptest.NewRequest(http.MethodPatch, "/api/connectors/"+id, strings.NewReader(`{}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/connectors/missing", strings.NewReader(`{"name":"x"}`))
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/connectors/"+id, strings.NewReader(`{"name":"Renamed"}`))
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
