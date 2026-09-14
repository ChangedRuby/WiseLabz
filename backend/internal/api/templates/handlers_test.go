package templates

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/doc"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	return NewHandler(s, doc.NewEngine(s))
}

func TestListEmpty(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestGetNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/templates/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Get(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestDeleteNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/templates/missing", nil)
	req.SetPathValue("id", "missing")
	rr := httptest.NewRecorder()
	h.Delete(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestCreate(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/templates", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/templates", strings.NewReader(`{"name":""}`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path with sections", func(t *testing.T) {
		body := `{"name":"Runbook Template","description":"d","sections":[{"title":"Overview","order":1,"body":"# hi"}]}`
		req := httptest.NewRequest(http.MethodPost, "/api/templates", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}

		var created map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		id, _ := created["id"].(string)
		if id == "" {
			t.Fatalf("no id in create response: %s", rr.Body.String())
		}

		getReq := httptest.NewRequest(http.MethodGet, "/api/templates/"+id, nil)
		getReq.SetPathValue("id", id)
		getRR := httptest.NewRecorder()
		h.Get(getRR, getReq)
		if getRR.Code != http.StatusOK {
			t.Fatalf("Get() status = %d, want %d; body=%s", getRR.Code, http.StatusOK, getRR.Body.String())
		}

		versionsReq := httptest.NewRequest(http.MethodGet, "/api/templates/"+id+"/versions", nil)
		versionsReq.SetPathValue("id", id)
		versionsRR := httptest.NewRecorder()
		h.Versions(versionsRR, versionsReq)
		if versionsRR.Code != http.StatusOK {
			t.Fatalf("Versions() status = %d, want %d; body=%s", versionsRR.Code, http.StatusOK, versionsRR.Body.String())
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/templates/"+id, nil)
		deleteReq.SetPathValue("id", id)
		deleteRR := httptest.NewRecorder()
		h.Delete(deleteRR, deleteReq)
		if deleteRR.Code != http.StatusNoContent {
			t.Fatalf("Delete() status = %d, want %d; body=%s", deleteRR.Code, http.StatusNoContent, deleteRR.Body.String())
		}
	})
}

func TestVersionNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/templates/missing/versions/1", nil)
	req.SetPathValue("id", "missing")
	req.SetPathValue("rev", "1")
	rr := httptest.NewRecorder()
	h.Version(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNotFound, rr.Body.String())
	}
}
