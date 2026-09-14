package auth

import (
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
