package docs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/api/settings"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/doc"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	settingsH := settings.NewHandler(s, &config.Config{}, ai.NewRegistry())
	return NewHandler(s, doc.NewEngine(s), settingsH, ai.NewRegistry(), nil)
}

func TestListEmpty(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetRootIsSynthetic(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/root", nil)
	req.SetPathValue("id", "root")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"lab"`) {
		t.Errorf("expected synthetic lab doc, got: %s", rr.Body.String())
	}
}

func TestGetUnknownIDFallsBackToServicePlaceholder(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/unknown-connector", nil)
	req.SetPathValue("id", "unknown-connector")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"service"`) {
		t.Errorf("expected service placeholder, got: %s", rr.Body.String())
	}
}

func TestByServiceNoDocsYet(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/service/conn-1", nil)
	req.SetPathValue("id", "conn-1")
	rr := httptest.NewRecorder()
	h.ByService(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestSave(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/docs/x", strings.NewReader(`{`))
		req.SetPathValue("id", "x")
		rr := httptest.NewRecorder()
		h.Save(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/docs/missing", strings.NewReader(`{"content":"hi"}`))
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Save(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
		}
	})
}

func TestVersionsOfUnknownDoc(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/missing/versions", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Versions(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetLockNoneHeld(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/docs/x/lock", nil)
	req.SetPathValue("id", "x")
	rr := httptest.NewRecorder()
	h.GetLock(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}
