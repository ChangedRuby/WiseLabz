package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func loginRefreshCookie(t *testing.T, app *testApp, username string) *http.Cookie {
	t.Helper()
	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": username, "password": "correct-password"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("login status = %d: %s", rec.Code, rec.Body)
	}
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "refresh_token" && cookie.Path == "/" {
			return cookie
		}
	}
	t.Fatal("login did not set refresh token cookie")
	return nil
}

func TestLoginExpiresLegacyRefreshCookiePath(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "alice", "correct-password", "viewer")
	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "alice", "password": "correct-password"}, "")
	var legacyExpired bool
	for _, cookie := range rec.Result().Cookies() {
		legacyExpired = legacyExpired || (cookie.Name == "refresh_token" && cookie.Path == "/api/auth" && cookie.MaxAge < 0)
	}
	if !legacyExpired {
		t.Fatal("login did not expire legacy /api/auth refresh cookie")
	}
}

func seedLocalUser(t *testing.T, app *testApp, username, password, role string) *store.User {
	t.Helper()
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	instanceAdminRole := "user"
	if role == "operator" {
		instanceAdminRole = "admin"
	}
	u := &store.User{Username: username, DisplayName: username, InstanceAdminRole: instanceAdminRole, AuthSource: "local", PasswordHash: hash}
	if err := app.Store.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

func TestLoginValidation(t *testing.T) {
	app := newTestApp(t)

	tests := []struct {
		name string
		body any
	}{
		{"malformed json", "not-json"},
		{"missing password", map[string]any{"username": "alice"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/auth/login", tt.body, "")
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestLoginWrongPassword(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "alice", "correct-password", "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "alice", "password": "wrong-password"}, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body)
	}
}

func TestLoginSuccess(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "alice", "correct-password", "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "alice", "password": "correct-password"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestLoginCreatesLocalSessionWithEmptyAuthProvider(t *testing.T) {
	app := newTestApp(t)
	user := seedLocalUser(t, app, "alice", "correct-password", "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{
		"username":   "alice",
		"password":   "correct-password",
		"providerId": "authentik",
	}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	sessions, err := app.Store.ListUserSessions(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("ListUserSessions() error: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}
	if sessions[0].AuthProviderID != "" {
		t.Fatalf("AuthProviderID = %q, want empty local auth source", sessions[0].AuthProviderID)
	}
}

func TestRefreshRotatesOnlyActiveSession(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "alice", "correct-password", "viewer")
	first := loginRefreshCookie(t, app, "alice")
	if first.Path != "/" {
		t.Fatalf("cookie path = %q, want /", first.Path)
	}
	request := app.newRequest(t, http.MethodPost, "/api/auth/refresh", nil, "")
	request.AddCookie(first)
	rotated := app.serve(request)
	if rotated.Code != http.StatusOK {
		t.Fatalf("refresh status = %d: %s", rotated.Code, rotated.Body)
	}
	replay := app.newRequest(t, http.MethodPost, "/api/auth/refresh", nil, "")
	replay.AddCookie(first)
	if rec := app.serve(replay); rec.Code != http.StatusUnauthorized {
		t.Fatalf("replay status = %d, want 401", rec.Code)
	}
}

func TestRefreshRejectsAccessToken(t *testing.T) {
	app := newTestApp(t)
	_, access := app.user(t, "viewer")
	rec := app.req(t, http.MethodPost, "/api/auth/refresh", map[string]any{"refreshToken": access}, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestElevationCannotCrossOperators(t *testing.T) {
	app := newTestApp(t)
	ownerID, _ := app.user(t, "operator")
	_, otherToken := app.user(t, "operator")
	targetID, _ := app.user(t, "viewer")
	req := app.newRequest(t, http.MethodDelete, "/api/users/"+targetID, nil, otherToken)
	req.Header.Set("X-Elevation-Token", app.elevationToken(t, ownerID, "user.delete"))
	if rec := app.serve(req); rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestUpdateMeSuccess(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodPatch, "/api/me", map[string]any{"displayName": "New Name"}, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestChangePasswordWrongCurrentPassword(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")
	_ = userID

	rec := app.req(t, http.MethodPost, "/api/me/password", map[string]any{
		"currentPassword": "wrong", "newPassword": "new-password-123",
	}, token)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body)
	}
}

func TestUsersRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	t.Run("list", func(t *testing.T) {
		rec := app.req(t, http.MethodGet, "/api/users", nil, viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("create", func(t *testing.T) {
		rec := app.req(t, http.MethodPost, "/api/users", map[string]any{"username": "bob", "password": "password123"}, viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestUsersCreateValidation(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/users", map[string]any{"username": "bob"}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestUsersCreateSuccess(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/users", map[string]any{"username": "bob", "password": "password123"}, opToken)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body)
	}
}

func TestUsersDeleteElevationBoundary(t *testing.T) {
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	targetID, _ := app.user(t, "viewer")

	t.Run("missing elevation token", func(t *testing.T) {
		rec := app.req(t, http.MethodDelete, "/api/users/"+targetID, nil, opToken)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("self-delete blocked even with valid elevation", func(t *testing.T) {
		tok := app.elevationToken(t, opID, "user.delete")
		req := app.newRequest(t, http.MethodDelete, "/api/users/"+opID, nil, opToken)
		req.Header.Set("X-Elevation-Token", tok)
		rec := app.serve(req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("valid elevation token succeeds", func(t *testing.T) {
		tok := app.elevationToken(t, opID, "user.delete")
		req := app.newRequest(t, http.MethodDelete, "/api/users/"+targetID, nil, opToken)
		req.Header.Set("X-Elevation-Token", tok)
		rec := app.serve(req)
		if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 200 or 204; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestUsersResetPasswordElevationBoundary(t *testing.T) {
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	targetID, _ := app.user(t, "viewer")

	t.Run("missing elevation token", func(t *testing.T) {
		rec := app.req(t, http.MethodPost, "/api/users/"+targetID+"/reset-password", map[string]any{"newPassword": "new-password-123"}, opToken)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("valid elevation token succeeds", func(t *testing.T) {
		tok := app.elevationToken(t, opID, "user.resetPassword")
		req := app.newRequest(t, http.MethodPost, "/api/users/"+targetID+"/reset-password", map[string]any{"newPassword": "new-password-123"}, opToken)
		req.Header.Set("X-Elevation-Token", tok)
		rec := app.serve(req)
		if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 200 or 204; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestLogoutSuccess(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "testuser", "correct-password", "viewer")
	cookie := loginRefreshCookie(t, app, "testuser")

	// Get access token from login response
	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "testuser", "password": "correct-password"}, "")
	var loginResp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	token := loginResp["accessToken"].(string)

	req := app.newRequest(t, http.MethodPost, "/api/auth/logout", nil, token)
	req.AddCookie(cookie)
	rec = app.serve(req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body)
	}
}

func TestElevateSuccess(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "alice", "correct-password", "viewer")
	_, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/elevate", map[string]any{
		"password": "password123",
		"action":   "connector.delete",
	}, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["token"] == nil || resp["token"] == "" {
		t.Fatalf("expected token in response, got %v", resp)
	}
}

func TestElevateWrongPassword(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/elevate", map[string]any{
		"password": "wrong-password",
		"action":   "connector.delete",
	}, token)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body)
	}
}

func TestProvidersEmpty(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/auth/providers", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	oidc, ok := resp["oidc"]
	if !ok {
		t.Fatalf("expected 'oidc' field in response")
	}
	// Should be an empty array, not null
	oidcList, ok := oidc.([]any)
	if !ok {
		t.Fatalf("expected 'oidc' to be an array, got %T", oidc)
	}
	if len(oidcList) != 0 {
		t.Fatalf("expected empty oidc list, got %d providers", len(oidcList))
	}
}

func TestMeSuccess(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/me", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["id"] != userID {
		t.Fatalf("expected id = %s, got %v", userID, resp["id"])
	}
	if resp["role"] != "user" {
		t.Fatalf("expected role = user, got %v", resp["role"])
	}
}

func TestListSessionsSuccess(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "sessionuser", "correct-password", "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "sessionuser", "password": "correct-password"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("login failed: status = %d, body = %s", rec.Code, rec.Body)
	}

	var loginResp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	token := loginResp["accessToken"].(string)

	rec = app.req(t, http.MethodGet, "/api/me/sessions", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var resp []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	// Should have at least one session (the one created during login)
	if len(resp) == 0 {
		t.Fatalf("expected at least one session")
	}
}

func TestDeleteSessionSuccess(t *testing.T) {
	app := newTestApp(t)
	seedLocalUser(t, app, "deluser", "correct-password", "viewer")

	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "deluser", "password": "correct-password"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("login failed: status = %d, body = %s", rec.Code, rec.Body)
	}

	var loginResp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	token := loginResp["accessToken"].(string)

	// First, get the list of sessions
	rec = app.req(t, http.MethodGet, "/api/me/sessions", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var sessions []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&sessions); err != nil {
		t.Fatalf("decode sessions: %v", err)
	}
	if len(sessions) == 0 {
		t.Fatalf("expected at least one session")
	}

	sessionID := sessions[0]["id"].(string)

	// Delete a session
	rec = app.req(t, http.MethodDelete, "/api/me/sessions/"+sessionID, nil, token)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body)
	}
}

func TestDeleteSessionNotOwner(t *testing.T) {
	app := newTestApp(t)

	// Create and login user1
	seedLocalUser(t, app, "user1", "correct-password", "viewer")
	rec := app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "user1", "password": "correct-password"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("user1 login failed: status = %d", rec.Code)
	}
	var loginResp1 map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&loginResp1); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	user1Token := loginResp1["accessToken"].(string)

	// Create and login user2
	seedLocalUser(t, app, "user2", "correct-password", "viewer")
	rec = app.req(t, http.MethodPost, "/api/auth/login", map[string]any{"username": "user2", "password": "correct-password"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("user2 login failed: status = %d", rec.Code)
	}
	var loginResp2 map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&loginResp2); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	user2Token := loginResp2["accessToken"].(string)

	// Get user2's session ID
	rec = app.req(t, http.MethodGet, "/api/me/sessions", nil, user2Token)
	if rec.Code != http.StatusOK {
		t.Fatalf("list sessions failed: status = %d", rec.Code)
	}

	var sessions []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&sessions); err != nil {
		t.Fatalf("decode sessions response: %v", err)
	}
	if len(sessions) == 0 {
		t.Fatalf("expected at least one session")
	}

	sessionID := sessions[0]["id"].(string)

	// Try to delete user2's session as user1
	rec = app.req(t, http.MethodDelete, "/api/me/sessions/"+sessionID, nil, user1Token)
	if rec.Code != http.StatusForbidden && rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 403 or 404; body = %s", rec.Code, rec.Body)
	}
}

func TestListUsersAsAdmin(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/users", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var resp []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	// Should have at least the created user
	if len(resp) == 0 {
		t.Fatalf("expected at least one user in list")
	}
}

func TestListUsersAsNonAdmin(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/users", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	targetID, _ := app.user(t, "viewer")

	rec := app.req(t, http.MethodPatch, "/api/users/"+targetID, map[string]any{"displayName": "Updated Name"}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["displayName"] != "Updated Name" {
		t.Fatalf("expected displayName = Updated Name, got %v", resp["displayName"])
	}
}

func TestUpdateUserAsNonAdmin(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	targetID, _ := app.user(t, "viewer")

	rec := app.req(t, http.MethodPatch, "/api/users/"+targetID, map[string]any{"displayName": "Updated Name"}, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestOIDCCallbackStateMismatch(t *testing.T) {
	app := newTestApp(t)

	rec := app.req(t, http.MethodPost, "/api/auth/oidc/callback", map[string]any{
		"providerId": "google",
		"code":       "auth_code_123",
		"state":      "invalid_state_token",
	}, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body)
	}
}

func TestOIDCCallbackNoProvider(t *testing.T) {
	app := newTestApp(t)

	// Create a valid state token (this mimics what a real OIDC flow would do)
	// But we'll use an invalid provider ID
	rec := app.req(t, http.MethodPost, "/api/auth/oidc/callback", map[string]any{
		"providerId": "nonexistent-provider",
		"code":       "auth_code_123",
		"state":      "SGFzdGVkU3RhdGU6MTIzNDU2Nzg5MDphYmMxMjM0NTY3OA==", // dummy state
	}, "")
	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 400 or 401; body = %s", rec.Code, rec.Body)
	}
}
