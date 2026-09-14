package runbooks

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

func TestListMutuallyExclusiveFilters(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/runbooks?changeType=deploy&alertSeverity=high", nil)
	rr := httptest.NewRecorder()
	h.List(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestGetNotFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/runbooks/missing", nil)
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
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(`{"title":""}`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid target type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(`{"title":"t","targetType":"bogus","targetValue":"v"}`))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path then conflict on duplicate target", func(t *testing.T) {
		body := `{"title":"Deploy runbook","body":"steps","targetType":"change_type","targetValue":"deploy"}`
		req := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}

		dupReq := httptest.NewRequest(http.MethodPost, "/api/runbooks", strings.NewReader(body))
		dupRR := httptest.NewRecorder()
		h.Create(dupRR, dupReq)
		if dupRR.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d; body=%s", dupRR.Code, http.StatusConflict, dupRR.Body.String())
		}
	})
}

func TestUpdateAndDelete(t *testing.T) {
	h := newTestHandler(t)

	createReq := httptest.NewRequest(http.MethodPost, "/api/runbooks",
		strings.NewReader(`{"title":"Original","body":"b","targetType":"alert_severity","targetValue":"high"}`))
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

	t.Run("update invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update invalid field type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{"title":123}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update invalid target type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{"targetType":"bogus"}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("update not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/missing", strings.NewReader(`{"title":"x"}`))
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("update happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/runbooks/"+id, strings.NewReader(`{"title":"Updated","docId":null}`))
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Update(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})

	t.Run("delete not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/runbooks/missing", nil)
		req.SetPathValue("id", "missing")
		rr := httptest.NewRecorder()
		h.Delete(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("delete happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/runbooks/"+id, nil)
		req.SetPathValue("id", id)
		rr := httptest.NewRecorder()
		h.Delete(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}
	})
}
