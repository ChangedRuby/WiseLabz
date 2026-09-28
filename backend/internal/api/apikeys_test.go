package api_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAPIKeyRoutesEndToEnd(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	expiresAt := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)

	rec := app.req(t, http.MethodPost, "/api/auth/api-keys", map[string]any{
		"name":      "CI",
		"expiresAt": expiresAt,
	}, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create API key status = %d: %s", rec.Code, rec.Body)
	}
	var created map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	rawToken, ok := created["token"].(string)
	if !ok || !strings.HasPrefix(rawToken, "wlz_") {
		t.Fatalf("created token = %v, want wlz_ prefix", created["token"])
	}
	keyID, _ := created["id"].(string)
	if keyID == "" {
		t.Fatal("created response omitted id")
	}
	for _, field := range []string{"lastUsedAt", "revokedAt"} {
		if _, ok := created[field]; !ok {
			t.Errorf("created response omitted %s", field)
		}
	}

	list := app.req(t, http.MethodGet, "/api/auth/api-keys", nil, token)
	if list.Code != http.StatusOK {
		t.Fatalf("list API keys status = %d: %s", list.Code, list.Body)
	}
	var keys []map[string]any
	if err := json.NewDecoder(list.Body).Decode(&keys); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(keys) != 1 || keys[0]["name"] != "CI" {
		t.Fatalf("list response = %#v, want CI key", keys)
	}
	if _, ok := keys[0]["token"]; ok {
		t.Error("list response exposed raw token")
	}
	if _, ok := keys[0]["tokenHash"]; ok {
		t.Error("list response exposed token hash")
	}

	me := app.req(t, http.MethodGet, "/api/me", nil, rawToken)
	if me.Code != http.StatusOK {
		t.Fatalf("API key authenticated request status = %d: %s", me.Code, me.Body)
	}

	otherID, otherToken := app.user(t, "viewer")
	_ = otherID
	forbidden := app.req(t, http.MethodDelete, "/api/auth/api-keys/"+keyID, nil, otherToken)
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("cross-user revoke status = %d, want 403", forbidden.Code)
	}

	revoke := app.req(t, http.MethodDelete, "/api/auth/api-keys/"+keyID, nil, token)
	if revoke.Code != http.StatusNoContent {
		t.Fatalf("revoke status = %d: %s", revoke.Code, revoke.Body)
	}
	if after := app.req(t, http.MethodGet, "/api/me", nil, rawToken); after.Code != http.StatusUnauthorized {
		t.Fatalf("revoked API key status = %d, want 401", after.Code)
	}
	_ = userID
}

func TestAPIKeyCreateRejectsInvalidExpiryAndEmptyName(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, token := app.user(t, "viewer")
	for name, body := range map[string]any{
		"empty name":        map[string]any{"name": "  "},
		"expired timestamp": map[string]any{"name": "CI", "expiresAt": time.Now().UTC().Add(-time.Minute).Format(time.RFC3339)},
		"invalid timestamp": map[string]any{"name": "CI", "expiresAt": "not-a-time"},
	} {
		t.Run(name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/auth/api-keys", body, token)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400: %s", rec.Code, rec.Body)
			}
		})
	}
}
