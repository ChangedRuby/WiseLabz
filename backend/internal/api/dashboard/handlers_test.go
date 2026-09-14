package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	return NewHandler(apitest.NewStore(t))
}

func TestOverview(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/overview?days=30", nil)
	rr := httptest.NewRecorder()
	h.Overview(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetLayoutFallsBackToAdminDefault(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/layout", nil)
	rr := httptest.NewRecorder()
	h.GetLayout(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetAdminDefault(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/layout/admin-default", nil)
	rr := httptest.NewRecorder()
	h.GetAdminDefault(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestPutAdminDefault(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/dashboard/layout/admin-default", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.PutAdminDefault(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path persists and is read back by GetAdminDefault", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/dashboard/layout/admin-default", strings.NewReader(`{"widgets":[{"id":"w1"}]}`))
		rr := httptest.NewRecorder()
		h.PutAdminDefault(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}

		getReq := httptest.NewRequest(http.MethodGet, "/api/dashboard/layout/admin-default", nil)
		getRR := httptest.NewRecorder()
		h.GetAdminDefault(getRR, getReq)
		if !strings.Contains(getRR.Body.String(), `"w1"`) {
			t.Errorf("GetAdminDefault() did not reflect update: %s", getRR.Body.String())
		}
	})

	t.Run("null widgets defaults to empty array", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/dashboard/layout/admin-default", strings.NewReader(`{"widgets":null}`))
		rr := httptest.NewRecorder()
		h.PutAdminDefault(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
		var resp map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		widgets, ok := resp["widgets"].([]any)
		if !ok || len(widgets) != 0 {
			t.Errorf("widgets = %v, want empty array", resp["widgets"])
		}
	})
}

func TestSaveAndResetLayout(t *testing.T) {
	h := newTestHandler(t)

	t.Run("save invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/dashboard/layout", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.SaveLayout(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("save then get reflects the saved layout", func(t *testing.T) {
		saveReq := httptest.NewRequest(http.MethodPut, "/api/dashboard/layout", strings.NewReader(`{"widgets":[{"id":"mine"}]}`))
		saveRR := httptest.NewRecorder()
		h.SaveLayout(saveRR, saveReq)
		if saveRR.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", saveRR.Code, http.StatusOK, saveRR.Body.String())
		}

		getReq := httptest.NewRequest(http.MethodGet, "/api/dashboard/layout", nil)
		getRR := httptest.NewRecorder()
		h.GetLayout(getRR, getReq)
		if !strings.Contains(getRR.Body.String(), `"mine"`) {
			t.Errorf("GetLayout() did not reflect saved layout: %s", getRR.Body.String())
		}
	})

	t.Run("reset overwrites saved layout with admin default", func(t *testing.T) {
		putReq := httptest.NewRequest(http.MethodPut, "/api/dashboard/layout/admin-default", strings.NewReader(`{"widgets":[{"id":"admin-default"}]}`))
		h.PutAdminDefault(httptest.NewRecorder(), putReq)

		resetReq := httptest.NewRequest(http.MethodPost, "/api/dashboard/layout/reset", nil)
		resetRR := httptest.NewRecorder()
		h.ResetLayout(resetRR, resetReq)
		if resetRR.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", resetRR.Code, http.StatusOK, resetRR.Body.String())
		}
		if !strings.Contains(resetRR.Body.String(), `"admin-default"`) {
			t.Errorf("ResetLayout() did not apply admin default: %s", resetRR.Body.String())
		}
	})
}
