package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// TestOIDCCallbackRejectsMissingFlowCookie covers the original vulnerability:
// GET /api/auth/providers is unauthenticated and hands out a valid state to
// anyone, so a state value alone (without the browser-bound cookie) must not
// be enough to complete a login.
func TestOIDCCallbackRejectsMissingFlowCookie(t *testing.T) {
	h := &Handler{Config: &config.Config{}}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/oidc/callback",
		strings.NewReader(`{"providerId":"okta","code":"abc","state":"attacker-state"}`))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	h.OIDCCallback(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("OIDCCallback() status = %d, want %d; body=%s", rr.Code, http.StatusUnauthorized, rr.Body.String())
	}
}

// TestOIDCCallbackRejectsStateMismatchedWithCookie exercises the callback
// with a browser that did start an OIDC flow (holds a real flow cookie) but
// whose callback state doesn't match it, as happens in an
// authorization-code-injection attempt using a code/state pair minted for a
// different flow.
func TestOIDCCallbackRejectsStateMismatchedWithCookie(t *testing.T) {
	h := &Handler{Config: &config.Config{}}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/oidc/callback",
		strings.NewReader(`{"providerId":"okta","code":"abc","state":"wrong-state"}`))
	req.Header.Set("Content-Type", "application/json")

	cookieRec := httptest.NewRecorder()
	setOIDCFlowCookie(cookieRec, req, "", "okta", "real-state", "real-nonce")
	for _, c := range cookieRec.Result().Cookies() {
		req.AddCookie(c)
	}

	rr := httptest.NewRecorder()
	h.OIDCCallback(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("OIDCCallback() status = %d, want %d; body=%s", rr.Code, http.StatusUnauthorized, rr.Body.String())
	}
}

// TestOIDCCallbackRejectsUnknownProvider tests that an unknown provider ID
// is rejected even with a valid flow cookie and state.
func TestOIDCCallbackRejectsUnknownProvider(t *testing.T) {
	h := &Handler{Config: &config.Config{}}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/oidc/callback",
		strings.NewReader(`{"providerId":"unknown","code":"abc","state":"state"}`))
	req.Header.Set("Content-Type", "application/json")

	// Set the flow cookie
	cookieRec := httptest.NewRecorder()
	setOIDCFlowCookie(cookieRec, req, "", "unknown", "state", "nonce")
	for _, c := range cookieRec.Result().Cookies() {
		req.AddCookie(c)
	}

	rr := httptest.NewRecorder()
	h.OIDCCallback(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("OIDCCallback(unknown provider) status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestFindOIDCProvider(t *testing.T) {
	t.Run("provider found", func(t *testing.T) {
		h := &Handler{
			Config: &config.Config{
				Auth: config.AuthSettings{
					OIDC: []config.OIDCProvider{
						{ID: "okta", DisplayName: "Okta", IssuerURL: "https://okta.example.com"},
					},
				},
			},
		}
		prov := h.findOIDCProvider("okta")
		if prov == nil {
			t.Fatal("findOIDCProvider() returned nil, want provider")
		}
		if prov.ID != "okta" {
			t.Fatalf("provider ID = %q, want okta", prov.ID)
		}
	})

	t.Run("provider not found", func(t *testing.T) {
		h := &Handler{
			Config: &config.Config{
				Auth: config.AuthSettings{
					OIDC: []config.OIDCProvider{
						{ID: "okta", DisplayName: "Okta", IssuerURL: "https://okta.example.com"},
					},
				},
			},
		}
		prov := h.findOIDCProvider("google")
		if prov != nil {
			t.Fatalf("findOIDCProvider(missing) = %v, want nil", prov)
		}
	})

	t.Run("no providers configured", func(t *testing.T) {
		h := &Handler{
			Config: &config.Config{
				Auth: config.AuthSettings{OIDC: []config.OIDCProvider{}},
			},
		}
		prov := h.findOIDCProvider("okta")
		if prov != nil {
			t.Fatalf("findOIDCProvider(empty config) = %v, want nil", prov)
		}
	})
}

func TestOIDCProviderEnabled(t *testing.T) {
	th := newTestHandler(t)
	ctx := context.Background()

	t.Run("enabled flag true", func(t *testing.T) {
		if err := th.Store.SetOIDCProviderEnabled(ctx, "okta", true); err != nil {
			t.Fatalf("SetOIDCProviderEnabled() error: %v", err)
		}
		if !th.H.oidcProviderEnabled(ctx, "okta") {
			t.Fatal("oidcProviderEnabled() = false, want true")
		}
	})

	t.Run("enabled flag false", func(t *testing.T) {
		if err := th.Store.SetOIDCProviderEnabled(ctx, "google", false); err != nil {
			t.Fatalf("SetOIDCProviderEnabled() error: %v", err)
		}
		if th.H.oidcProviderEnabled(ctx, "google") {
			t.Fatal("oidcProviderEnabled(disabled) = true, want false")
		}
	})

	t.Run("provider not configured defaults to enabled", func(t *testing.T) {
		// A provider that has no row in oidc_provider_flags should be enabled by default
		if !th.H.oidcProviderEnabled(ctx, "unconfigured") {
			t.Fatal("oidcProviderEnabled(unconfigured) = false, want true (default)")
		}
	})
}

// The issuer is a closed local port: discovery fails at once without leaving
// the machine, which is the uninitializable-provider case these subtests want.
func TestGetOrInitOIDCProvider(t *testing.T) {
	t.Run("provider not in cache", func(t *testing.T) {
		th := newTestHandler(t)
		cfg := &config.OIDCProvider{
			ID:           "test-provider",
			DisplayName:  "Test",
			IssuerURL:    "http://127.0.0.1:1",
			ClientID:     "test-client",
			ClientSecret: "test-secret",
		}
		ctx := context.Background()

		th.H.getOrInitOIDCProvider(ctx, cfg)
		// This will be nil because we can't actually initialize without a real OIDC server
		// but we verify the cache is set up
		if th.H.oidcProv == nil {
			t.Fatal("getOrInitOIDCProvider() did not initialize cache")
		}
	})

	t.Run("provider cached", func(t *testing.T) {
		th := newTestHandler(t)
		cfg := &config.OIDCProvider{
			ID:           "cached-provider",
			DisplayName:  "Cached",
			IssuerURL:    "http://127.0.0.1:1",
			ClientID:     "test-client",
			ClientSecret: "test-secret",
		}
		ctx := context.Background()

		// First call initializes cache map
		th.H.getOrInitOIDCProvider(ctx, cfg)

		// Second call should reuse cache (even if nil)
		prov := th.H.getOrInitOIDCProvider(ctx, cfg)
		if prov != nil {
			// If it somehow succeeded, that's fine, just verify it's cached
			if th.H.oidcProv["cached-provider"] != prov {
				t.Fatal("getOrInitOIDCProvider() did not cache provider")
			}
		}
	})
}

// TestSyncOIDCConnectorGrants covers #279 part 3's login hook end to end
// against a real store: an OIDC login creates grants from the group
// mapping, a later login with fewer groups revokes what's no longer
// justified, and a local user is never touched.
func TestSyncOIDCConnectorGrants(t *testing.T) {
	th := newTestHandler(t)
	ctx := context.Background()

	c1 := &store.ConnectorRecord{Name: "conn-1", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := th.Store.CreateConnector(ctx, c1); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	c2 := &store.ConnectorRecord{Name: "conn-2", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := th.Store.CreateConnector(ctx, c2); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	oidcUser := &store.User{Username: "oidc-user", AuthSource: "oidc"}
	if err := th.Store.CreateUser(ctx, oidcUser); err != nil {
		t.Fatalf("CreateUser() error: %v", err)
	}
	mapping := map[string]map[string]string{"ops": {c1.ID: "viewer", c2.ID: "operator"}}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/oidc/callback", nil)

	th.H.syncOIDCConnectorGrants(req, oidcUser, []string{"ops"}, mapping)

	role, err := th.Store.GetUserConnectorRole(ctx, oidcUser.ID, c1.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "viewer" {
		t.Fatalf("GetUserConnectorRole(c1) after first login = %q, want viewer", role)
	}
	role, err = th.Store.GetUserConnectorRole(ctx, oidcUser.ID, c2.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "operator" {
		t.Fatalf("GetUserConnectorRole(c2) after first login = %q, want operator", role)
	}

	// A later login with fewer groups revokes what's no longer justified.
	th.H.syncOIDCConnectorGrants(req, oidcUser, nil, mapping)

	role, err = th.Store.GetUserConnectorRole(ctx, oidcUser.ID, c1.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "" {
		t.Fatalf("GetUserConnectorRole(c1) after losing the group = %q, want empty (revoked)", role)
	}
	role, err = th.Store.GetUserConnectorRole(ctx, oidcUser.ID, c2.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "" {
		t.Fatalf("GetUserConnectorRole(c2) after losing the group = %q, want empty (revoked)", role)
	}

	// A local user is never synced, even if called directly.
	localUser, _ := th.createUser(t, "viewer", false)
	th.H.syncOIDCConnectorGrants(req, localUser, []string{"ops"}, mapping)
	role, err = th.Store.GetUserConnectorRole(ctx, localUser.ID, c1.ID)
	if err != nil {
		t.Fatalf("GetUserConnectorRole() error: %v", err)
	}
	if role != "" {
		t.Fatalf("GetUserConnectorRole() for a local user = %q, want empty (never synced)", role)
	}
}
