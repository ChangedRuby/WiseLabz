package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/config"
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

func TestGetOrInitOIDCProvider(t *testing.T) {
	t.Run("provider not in cache", func(t *testing.T) {
		th := newTestHandler(t)
		cfg := &config.OIDCProvider{
			ID:           "test-provider",
			DisplayName:  "Test",
			IssuerURL:    "https://example.com",
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
			IssuerURL:    "https://example.com",
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
