package attention

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
)

func TestList(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/attention", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestListWithDaysFilter(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/api/attention?days=7&page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}
