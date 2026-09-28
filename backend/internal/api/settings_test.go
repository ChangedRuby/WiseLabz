package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestAuthConfigRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/auth/config", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestAuthConfigGetAndUpdateSuccess(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/auth/config", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPut, "/api/auth/config", map[string]any{"stepUpForDestructive": false}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestAuthConfigUpdateValidation(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/auth/config", map[string]any{}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestProviderDisableRevokesMatchingSessionsAndAuditsCount(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	app.Config.Auth.OIDC = []config.OIDCProvider{{
		ID:          "authentik",
		DisplayName: "Authentik",
	}}
	_, opToken := app.user(t, "operator")
	user := &store.User{Username: "oidc-session-user"}
	if err := app.Store.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	for _, sess := range []store.Session{
		{UserID: user.ID, TokenHash: "local"},
		{UserID: user.ID, TokenHash: "authentik-1", AuthProviderID: "authentik"},
		{UserID: user.ID, TokenHash: "authentik-2", AuthProviderID: "authentik"},
		{UserID: user.ID, TokenHash: "google", AuthProviderID: "google"},
	} {
		if err := app.Store.CreateSession(context.Background(), &sess); err != nil {
			t.Fatalf("CreateSession(%q) error: %v", sess.TokenHash, err)
		}
	}

	rec := app.req(t, http.MethodPut, "/api/auth/providers/authentik/enabled", map[string]any{"enabled": false}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	for _, tokenHash := range []string{"local", "google"} {
		active, err := app.Store.HasSessionTokenHash(context.Background(), user.ID, tokenHash)
		if err != nil || !active {
			t.Fatalf("HasSessionTokenHash(%q) = %v, %v; want true, nil", tokenHash, active, err)
		}
	}
	for _, tokenHash := range []string{"authentik-1", "authentik-2"} {
		active, err := app.Store.HasSessionTokenHash(context.Background(), user.ID, tokenHash)
		if err != nil || active {
			t.Fatalf("HasSessionTokenHash(%q) = %v, %v; want false, nil", tokenHash, active, err)
		}
	}

	audits, _, err := app.Store.ListAuditRecords(context.Background(), "auth.provider.enabled", "oidc_provider", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if len(audits) != 1 {
		t.Fatalf("len(audits) = %d, want 1", len(audits))
	}
	var detail map[string]any
	if err := json.Unmarshal([]byte(audits[0].Detail), &detail); err != nil {
		t.Fatalf("unmarshal detail: %v", err)
	}
	if detail["enabled"] != false || detail["revokedSessions"] != float64(2) {
		t.Fatalf("audit detail = %#v, want enabled=false revokedSessions=2", detail)
	}
}

func TestProviderEnableDoesNotRevokeSessions(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	app.Config.Auth.OIDC = []config.OIDCProvider{{
		ID:          "authentik",
		DisplayName: "Authentik",
	}}
	_, opToken := app.user(t, "operator")
	user := &store.User{Username: "enable-provider-user"}
	if err := app.Store.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	if err := app.Store.CreateSession(context.Background(), &store.Session{
		UserID:         user.ID,
		TokenHash:      "authentik",
		AuthProviderID: "authentik",
	}); err != nil {
		t.Fatalf("CreateSession() error: %v", err)
	}

	rec := app.req(t, http.MethodPut, "/api/auth/providers/authentik/enabled", map[string]any{"enabled": true}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	active, err := app.Store.HasSessionTokenHash(context.Background(), user.ID, "authentik")
	if err != nil || !active {
		t.Fatalf("HasSessionTokenHash() = %v, %v; want true, nil", active, err)
	}
	audits, _, err := app.Store.ListAuditRecords(context.Background(), "auth.provider.enabled", "oidc_provider", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	var detail map[string]any
	if err := json.Unmarshal([]byte(audits[0].Detail), &detail); err != nil {
		t.Fatalf("unmarshal detail: %v", err)
	}
	if _, ok := detail["revokedSessions"]; ok {
		t.Fatalf("audit detail = %#v, want no revokedSessions when enabling", detail)
	}
}

func TestAIConfigRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/ai/config", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestAIConfigGetAndUpdateSuccess(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/ai/config", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPut, "/api/ai/config", map[string]any{"enabled": false}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

// TestAIConfigTestDoesNotPanicOnNilRegistry is a regression test for #195:
// api.Config.AIRegistry was built in main.go but never wired into the
// api.Config literal, so cfg.AIRegistry stayed nil and every AI endpoint
// panicked in h.AI.Get(...). It's caught here at the router level: enabling
// AI with a registered provider and hitting /api/ai/config/test must reach a
// normal JSON response, not the Recoverer middleware's 500 from a nil
// *ai.Registry method call.
func TestAIConfigTestDoesNotPanicOnNilRegistry(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	// baseUrl is a closed local port so the provider call fails at once
	// instead of reaching the real OpenAI API.
	rec := app.req(t, http.MethodPut, "/api/ai/config", map[string]any{
		"enabled": true, "provider": "openai", "apiKey": "sk-test", "model": "gpt-4o-mini",
		"baseUrl": "http://127.0.0.1:1",
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPost, "/api/ai/config/test", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("test status = %d, want 200 (nil AIRegistry would panic into a 500); body = %s", rec.Code, rec.Body)
	}
}

func TestNotificationsConfigRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/notifications/config", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestNotificationsConfigGetAndUpdateSuccess(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/notifications/config", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPut, "/api/notifications/config", map[string]any{"channels": []any{}, "routing": []any{}}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestNotificationsConfigUpdateValidation(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/notifications/config", "not-json", opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}
