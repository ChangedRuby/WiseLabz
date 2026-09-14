package savedviews

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
)

func TestCreateAndList(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	_, token, createWrapped := apitest.AuthedUser(t, s, "operator", http.HandlerFunc(h.Create))

	t.Run("invalid surface", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/saved-views", strings.NewReader(`{"surface":"nope","name":"x"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/saved-views", strings.NewReader(`{"surface":"alerts","name":""}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/saved-views", strings.NewReader(`{`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.Create(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/api/saved-views", strings.NewReader(`{"surface":"alerts","name":"my view","filters":{"severity":"high"}}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	createWrapped.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
	}

	t.Run("list requires known surface", func(t *testing.T) {
		listReq := httptest.NewRequest(http.MethodGet, "/api/saved-views?surface=bogus", nil)
		listRR := httptest.NewRecorder()
		h.List(listRR, listReq)
		if listRR.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", listRR.Code, http.StatusBadRequest)
		}
	})
}

func TestDelete(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	jwtSvc := apitest.JWTService()

	ownerID := apitest.NewUser(t, s, "operator")
	otherID := apitest.NewUser(t, s, "operator")
	ownerToken := apitest.Token(t, jwtSvc, ownerID, "operator")
	otherToken := apitest.Token(t, jwtSvc, otherID, "operator")

	createWrapped := apitest.WithAuth(jwtSvc, s, http.HandlerFunc(h.Create))
	createReq := httptest.NewRequest(http.MethodPost, "/api/saved-views", strings.NewReader(`{"surface":"changes","name":"mine"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+ownerToken)
	createRR := httptest.NewRecorder()
	createWrapped.ServeHTTP(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("no id in create response: %s", createRR.Body.String())
	}

	deleteWrapped := apitest.WithAuth(jwtSvc, s, http.HandlerFunc(h.Delete))

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/saved-views/missing", nil)
		req.SetPathValue("id", "missing")
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()
		deleteWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("forbidden for a different user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/saved-views/"+id, nil)
		req.SetPathValue("id", id)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		rr := httptest.NewRecorder()
		deleteWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusForbidden, rr.Body.String())
		}
	})

	t.Run("owner can delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/saved-views/"+id, nil)
		req.SetPathValue("id", id)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()
		deleteWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}
	})
}
