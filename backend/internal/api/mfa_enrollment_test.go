package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/auth"
)

func postWithElevation(t *testing.T, app *testApp, path, token, elevation string) int {
	t.Helper()
	req := app.newRequest(t, http.MethodPost, path, map[string]any{"name": "x"}, token)
	if elevation != "" {
		req.Header.Set("X-Elevation-Token", elevation)
	}
	return app.serve(req).Code
}

func TestMFAEnrollmentRequiresStepUp(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")

	for _, path := range []string{"/api/me/mfa/totp", "/api/me/mfa/webauthn/register/begin"} {
		if got := postWithElevation(t, app, path, token, ""); got != http.StatusBadRequest {
			t.Errorf("%s without elevation = %d, want 400", path, got)
		}
		if got := postWithElevation(t, app, path, token, app.elevationToken(t, userID, "user.delete")); got != http.StatusUnauthorized {
			t.Errorf("%s with wrong-action elevation = %d, want 401", path, got)
		}
	}
	if got := postWithElevation(t, app, "/api/me/mfa/totp", token, app.elevationToken(t, userID, "mfa.manage")); got != http.StatusOK && got != http.StatusCreated {
		t.Errorf("totp with mfa.manage elevation = %d, want 2xx", got)
	}
}

func TestMFAEnrollOnlySessionSkipsStepUp(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, _ := app.user(t, "viewer")
	pair, err := app.JWT.IssuePairWithOptions(userID, false, auth.IssuePairOptions{MFAEnrollOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := postWithElevation(t, app, "/api/me/mfa/totp", pair.AccessToken, ""); got != http.StatusOK && got != http.StatusCreated {
		t.Fatalf("enroll-only totp without elevation = %d, want 2xx", got)
	}

	// Once a factor is confirmed the exemption ends.
	f, err := app.Store.CreatePendingTOTP(context.Background(), userID, "t", "s")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.Store.ConfirmFactor(context.Background(), f.ID); err != nil {
		t.Fatal(err)
	}
	if got := postWithElevation(t, app, "/api/me/mfa/totp", pair.AccessToken, ""); got != http.StatusBadRequest {
		t.Fatalf("enroll-only with confirmed factor, no elevation = %d, want 400", got)
	}
}

func TestAPIKeysRejectedOnAccountSecurityRoutes(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	elev := app.elevationToken(t, userID, "mfa.manage")

	for _, body := range []map[string]any{
		{"name": "full"},
		{"name": "ro", "scope": "read"},
	} {
		raw := createKey(t, app, token, body)["token"].(string)
		for _, c := range []struct{ method, path string }{
			{http.MethodGet, "/api/me/mfa"},
			{http.MethodPost, "/api/me/mfa/totp"},
			{http.MethodPost, "/api/me/mfa/webauthn/register/begin"},
			{http.MethodPost, "/api/me/password"},
			{http.MethodPost, "/api/auth/elevate"},
			{http.MethodGet, "/api/auth/elevate/methods"},
		} {
			req := app.newRequest(t, c.method, c.path, map[string]any{"name": "x"}, raw)
			req.Header.Set("X-Elevation-Token", elev)
			if rec := app.serve(req); rec.Code != http.StatusForbidden {
				t.Errorf("%v key %s %s = %d, want 403: %s", body["name"], c.method, c.path, rec.Code, rec.Body)
			}
		}
	}
	if rec := app.req(t, http.MethodGet, "/api/me/mfa", nil, token); rec.Code != http.StatusOK {
		t.Errorf("session GET /me/mfa = %d, want 200", rec.Code)
	}
}
