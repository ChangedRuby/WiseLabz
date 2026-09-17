package changes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/api/settings"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	settingsH := settings.NewHandler(s, &config.Config{}, ai.NewRegistry())
	return NewHandler(s, settingsH, ai.NewRegistry(), nil)
}

// withOperatorGrant seeds a user with an operator grant on connectorID and
// returns req with that user in context, as auth.AuthMiddleware would.
func withOperatorGrant(t *testing.T, h *Handler, req *http.Request, connectorID string) *http.Request {
	t.Helper()
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, connectorID, "operator")
	return req.WithContext(auth.ContextWithUser(req.Context(), userID, false))
}

func TestListEmpty(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/changes", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/changes/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestAcknowledgeNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/changes/missing/ack", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Acknowledge(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestDismissNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/changes/missing/dismiss", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Dismiss(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestAIUpdate(t *testing.T) {
	h := newTestHandler(t)

	t.Run("change not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/changes/missing/ai-update", nil)
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.AIUpdate(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("ai disabled returns conflict", func(t *testing.T) {
		h2 := newTestHandler(t)
		c := &store.ChangeRecord{
			ServiceID:      "svc-1",
			ChangeType:     "config",
			Severity:       "warning",
			Summary:        "Test change",
			Status:         "new",
			Diff:           "[]",
			AffectedDocIDs: "[]",
		}
		if err := h2.Store.CreateChange(context.Background(), c); err != nil {
			t.Fatalf("create change: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/changes/"+c.ID+"/ai-update", nil)
		req.SetPathValue("id", c.ID)
		req = withOperatorGrant(t, h2, req, c.ServiceID)
		rr := httptest.NewRecorder()
		h2.AIUpdate(rr, req)
		if rr.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusConflict, rr.Body.String())
		}
	})
}

func TestGetSuccess(t *testing.T) {
	h := newTestHandler(t)
	c := &store.ChangeRecord{
		ServiceID:      "svc-1",
		ChangeType:     "config",
		Severity:       "info",
		Summary:        "Test change",
		Status:         "new",
		Diff:           "[]",
		AffectedDocIDs: "[]",
	}
	if err := h.Store.CreateChange(context.Background(), c); err != nil {
		t.Fatalf("create change: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/changes/"+c.ID, nil)
	req.SetPathValue("id", c.ID)
	req = withOperatorGrant(t, h, req, c.ServiceID)
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result["id"] != c.ID {
		t.Errorf("id = %v, want %s", result["id"], c.ID)
	}
}

func TestAcknowledgeSuccess(t *testing.T) {
	h := newTestHandler(t)
	c := &store.ChangeRecord{
		ServiceID:      "svc-1",
		ChangeType:     "config",
		Severity:       "info",
		Summary:        "Test change",
		Status:         "new",
		Diff:           "[]",
		AffectedDocIDs: "[]",
	}
	if err := h.Store.CreateChange(context.Background(), c); err != nil {
		t.Fatalf("create change: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/changes/"+c.ID+"/ack", nil)
	req.SetPathValue("id", c.ID)
	req = withOperatorGrant(t, h, req, c.ServiceID)
	rr := httptest.NewRecorder()
	h.Acknowledge(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result["status"] != "acknowledged" {
		t.Errorf("status = %v, want acknowledged", result["status"])
	}
}

func TestDismissSuccess(t *testing.T) {
	h := newTestHandler(t)
	c := &store.ChangeRecord{
		ServiceID:      "svc-1",
		ChangeType:     "config",
		Severity:       "info",
		Summary:        "Test change",
		Status:         "new",
		Diff:           "[]",
		AffectedDocIDs: "[]",
	}
	if err := h.Store.CreateChange(context.Background(), c); err != nil {
		t.Fatalf("create change: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/changes/"+c.ID+"/dismiss", nil)
	req.SetPathValue("id", c.ID)
	req = withOperatorGrant(t, h, req, c.ServiceID)
	rr := httptest.NewRecorder()
	h.Dismiss(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	var result map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result["status"] != "dismissed" {
		t.Errorf("status = %v, want dismissed", result["status"])
	}
}

func TestBulkResolve(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/changes/bulk-resolve", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.BulkResolve(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/changes/bulk-resolve", strings.NewReader(`{"ids":["a"],"status":"deleted"}`))
		rr := httptest.NewRecorder()
		h.BulkResolve(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("empty ids", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/changes/bulk-resolve", strings.NewReader(`{"ids":[],"status":"acknowledged"}`))
		rr := httptest.NewRecorder()
		h.BulkResolve(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("too many ids", func(t *testing.T) {
		ids := make([]string, 501)
		for i := range ids {
			ids[i] = `"id` + strconv.Itoa(i) + `"`
		}
		body := `{"ids":[` + strings.Join(ids, ",") + `],"status":"acknowledged"}`
		req := httptest.NewRequest(http.MethodPost, "/api/changes/bulk-resolve", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.BulkResolve(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown ids reported per-item, not fatal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/changes/bulk-resolve", strings.NewReader(`{"ids":["missing-1"],"status":"acknowledged"}`))
		rr := httptest.NewRecorder()
		h.BulkResolve(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
		var resp map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		results, _ := resp["results"].([]any)
		if len(results) != 1 {
			t.Fatalf("len(results) = %d, want 1", len(results))
		}
		first, _ := results[0].(map[string]any)
		if first["reason"] != "not_found" {
			t.Errorf("reason = %v, want not_found", first["reason"])
		}
	})
}

// countingProvider is a fake ai.Provider that counts Suggest calls, used to
// verify Explain caches its result instead of re-invoking the provider.
type countingProvider struct{ calls int }

func (p *countingProvider) Name() string { return "mock" }
func (p *countingProvider) Suggest(_ context.Context, _ *ai.SuggestRequest) (string, error) {
	p.calls++
	return "This change matters because it affects the firewall.", nil
}
func (p *countingProvider) SuggestStream(_ context.Context, _ *ai.SuggestRequest) (<-chan ai.SuggestChunk, error) {
	return nil, nil
}

func TestExplain(t *testing.T) {
	t.Run("change not found", func(t *testing.T) {
		h := newTestHandler(t)
		req := httptest.NewRequest(http.MethodPost, "/api/changes/missing/explain", nil)
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Explain(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("ai disabled returns conflict", func(t *testing.T) {
		h := newTestHandler(t)
		c := &store.ChangeRecord{
			ServiceID: "svc-1", ChangeType: "config", Severity: "warning",
			Summary: "Test change", Status: "new", Diff: "[]", AffectedDocIDs: "[]",
		}
		if err := h.Store.CreateChange(context.Background(), c); err != nil {
			t.Fatalf("create change: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/api/changes/"+c.ID+"/explain", nil)
		req.SetPathValue("id", c.ID)
		req = withOperatorGrant(t, h, req, c.ServiceID)
		rr := httptest.NewRecorder()
		h.Explain(rr, req)
		if rr.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusConflict, rr.Body.String())
		}
	})

	t.Run("generates once, caches on second call", func(t *testing.T) {
		s := apitest.NewStore(t)
		if _, err := s.DB().ExecContext(context.Background(), `UPDATE ai_config SET enabled = 1, provider = 'mock' WHERE id = 1`); err != nil {
			t.Fatalf("enable ai config: %v", err)
		}
		registry := ai.NewRegistry()
		provider := &countingProvider{}
		registry.Register("mock", func(map[string]any) (ai.Provider, error) { return provider, nil })
		settingsH := settings.NewHandler(s, &config.Config{}, registry)
		h := NewHandler(s, settingsH, registry, nil)

		c := &store.ChangeRecord{
			ServiceID: "svc-1", ChangeType: "firewall.rule.modified", Severity: "warning",
			Summary: "Firewall rule changed", Status: "new", Diff: "[]", AffectedDocIDs: "[]",
		}
		if err := h.Store.CreateChange(context.Background(), c); err != nil {
			t.Fatalf("create change: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/changes/"+c.ID+"/explain", nil)
		req.SetPathValue("id", c.ID)
		req = withOperatorGrant(t, h, req, c.ServiceID)
		rr := httptest.NewRecorder()
		h.Explain(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
		var first map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &first); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if first["narration"] == "" {
			t.Fatalf("narration = %q, want non-empty", first["narration"])
		}
		if provider.calls != 1 {
			t.Fatalf("calls after first Explain = %d, want 1", provider.calls)
		}

		req2 := httptest.NewRequest(http.MethodPost, "/api/changes/"+c.ID+"/explain", nil)
		req2.SetPathValue("id", c.ID)
		req2 = withOperatorGrant(t, h, req2, c.ServiceID)
		rr2 := httptest.NewRecorder()
		h.Explain(rr2, req2)
		if rr2.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr2.Code, http.StatusOK, rr2.Body.String())
		}
		var second map[string]any
		if err := json.Unmarshal(rr2.Body.Bytes(), &second); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if second["narration"] != first["narration"] {
			t.Fatalf("narration changed between calls: %v vs %v", first["narration"], second["narration"])
		}
		if provider.calls != 1 {
			t.Fatalf("calls after second Explain = %d, want 1 (cached)", provider.calls)
		}
	})
}
