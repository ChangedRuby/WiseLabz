package apikeys

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
)

func TestCreate(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	_, token, wrapped := apitest.AuthedUser(t, s, "operator", http.HandlerFunc(h.Create))

	t.Run("happy path", func(t *testing.T) {
		body := `{"name":"ci-key"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}
		var resp map[string]any
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp["name"] != "ci-key" {
			t.Errorf("name = %v, want ci-key", resp["name"])
		}
		rawToken, _ := resp["token"].(string)
		if !strings.HasPrefix(rawToken, "wlz_") {
			t.Errorf("token = %q, want wlz_ prefix", rawToken)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(`{"name":"  "}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(`not-json`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("expiresAt in the past", func(t *testing.T) {
		body := `{"name":"expired","expiresAt":"2000-01-01T00:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("expiresAt not RFC3339", func(t *testing.T) {
		body := `{"name":"bad-date","expiresAt":"tomorrow"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("valid future expiresAt", func(t *testing.T) {
		body := `{"name":"future-key","expiresAt":"` + time.Now().Add(24*time.Hour).UTC().Format(time.RFC3339) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}
	})

	t.Run("unauthenticated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(`{"name":"x"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		wrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})
}

func TestList(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	jwtSvc := apitest.JWTService()
	userID := apitest.NewUser(t, s, "operator")
	tok := apitest.Token(t, jwtSvc, userID, "operator")

	createReq := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(`{"name":"listed-key"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+tok)
	createRR := httptest.NewRecorder()
	apitest.WithAuth(jwtSvc, s, http.HandlerFunc(h.Create)).ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("seed create status = %d, body=%s", createRR.Code, createRR.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/api-keys", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rr := httptest.NewRecorder()
	apitest.WithAuth(jwtSvc, s, http.HandlerFunc(h.List)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var out []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("len(out) = %d, want 1", len(out))
	}
	if _, hasToken := out[0]["token"]; hasToken {
		t.Errorf("List response leaked raw token field")
	}
	if _, hasHash := out[0]["tokenHash"]; hasHash {
		t.Errorf("List response leaked tokenHash field")
	}
}

func TestRevoke(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s)
	jwtSvc := apitest.JWTService()

	ownerID := apitest.NewUser(t, s, "operator")
	otherID := apitest.NewUser(t, s, "operator")
	ownerToken := apitest.Token(t, jwtSvc, ownerID, "operator")
	otherToken := apitest.Token(t, jwtSvc, otherID, "operator")

	createWrapped := apitest.WithAuth(jwtSvc, s, http.HandlerFunc(h.Create))
	createReq := httptest.NewRequest(http.MethodPost, "/api/auth/api-keys", strings.NewReader(`{"name":"to-revoke"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+ownerToken)
	createRR := httptest.NewRecorder()
	createWrapped.ServeHTTP(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	keyID, _ := created["id"].(string)
	if keyID == "" {
		t.Fatalf("no id in create response: %s", createRR.Body.String())
	}

	revokeWrapped := apitest.WithAuth(jwtSvc, s, http.HandlerFunc(h.Revoke))

	t.Run("missing id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/auth/api-keys/", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()
		revokeWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/auth/api-keys/missing", nil)
		req.SetPathValue("id", "missing-id")
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()
		revokeWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("forbidden for a different user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/auth/api-keys/"+keyID, nil)
		req.SetPathValue("id", keyID)
		req.Header.Set("Authorization", "Bearer "+otherToken)
		rr := httptest.NewRecorder()
		revokeWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusForbidden, rr.Body.String())
		}
	})

	t.Run("owner can revoke", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/auth/api-keys/"+keyID, nil)
		req.SetPathValue("id", keyID)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		rr := httptest.NewRecorder()
		revokeWrapped.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}
	})
}
