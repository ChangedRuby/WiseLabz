package findings

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
)

func TestListEmpty(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/findings", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetNotFound(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/findings/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Get(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestResolveNotFound(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodPost, "/api/findings/missing/resolve", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Resolve(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}
