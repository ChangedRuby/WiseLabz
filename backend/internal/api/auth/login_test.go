package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func doJSON(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	var r *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		r = bytes.NewReader(data)
	} else {
		r = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestLogin(t *testing.T) {
	th := newTestHandler(t)
	user, password := th.createUser(t, "viewer", false)

	t.Run("happy path", func(t *testing.T) {
		req := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": user.Username,
			"password": password,
		})
		rr := httptest.NewRecorder()
		th.H.Login(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
		var body struct {
			AccessToken string `json:"accessToken"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.AccessToken == "" {
			t.Fatal("expected non-empty accessToken")
		}
		if cookies := rr.Result().Cookies(); len(cookies) == 0 {
			t.Fatal("expected a refresh_token cookie to be set")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		req := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": user.Username,
			"password": "not-the-password",
		})
		rr := httptest.NewRecorder()
		th.H.Login(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		req := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": "does-not-exist",
			"password": "whatever",
		})
		rr := httptest.NewRecorder()
		th.H.Login(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		req := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{"username": user.Username})
		rr := httptest.NewRecorder()
		th.H.Login(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte("{not json")))
		rr := httptest.NewRecorder()
		th.H.Login(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("disabled user", func(t *testing.T) {
		disabledUser, disabledPassword := th.createUser(t, "viewer", true)
		req := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": disabledUser.Username,
			"password": disabledPassword,
		})
		rr := httptest.NewRecorder()
		th.H.Login(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("oidc-only user cannot use local login", func(t *testing.T) {
		oidcUser := &store.User{Username: "oidc-user", DisplayName: "OIDC", Role: "viewer", AuthSource: "oidc"}
		if err := th.Store.CreateUser(req(t).Context(), oidcUser); err != nil {
			t.Fatalf("create oidc user: %v", err)
		}
		r := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": oidcUser.Username,
			"password": "anything",
		})
		rr := httptest.NewRecorder()
		th.H.Login(rr, r)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("lockout after repeated failures", func(t *testing.T) {
		lockUser, _ := th.createUser(t, "viewer", false)
		for i := 0; i < maxFailedLoginAttempts; i++ {
			r := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
				"username": lockUser.Username,
				"password": "wrong-password",
			})
			rr := httptest.NewRecorder()
			th.H.Login(rr, r)
			if rr.Code != http.StatusUnauthorized {
				t.Fatalf("attempt %d: status = %d, want %d", i, rr.Code, http.StatusUnauthorized)
			}
		}

		// Even the correct password must now be rejected: the account is locked.
		_, correctPassword := lockUser, "password123"
		r := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": lockUser.Username,
			"password": correctPassword,
		})
		rr := httptest.NewRecorder()
		th.H.Login(rr, r)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("locked account: status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})
}

// req is a tiny helper for constructing a throwaway request when only its
// context is needed (e.g. as a context.Background() stand-in for store calls).
func req(t *testing.T) *http.Request {
	t.Helper()
	return httptest.NewRequest(http.MethodGet, "/", nil)
}

func TestChangePassword(t *testing.T) {
	th := newTestHandler(t)
	user, password := th.createUser(t, "viewer", false)

	t.Run("happy path", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/me/password", map[string]string{
			"currentPassword": password,
			"newPassword":     "new-password-456",
		})
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.ChangePassword)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}

		// The old password must no longer work.
		loginReq := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": user.Username,
			"password": password,
		})
		loginRR := httptest.NewRecorder()
		th.H.Login(loginRR, loginReq)
		if loginRR.Code != http.StatusUnauthorized {
			t.Fatalf("old password still works: status = %d", loginRR.Code)
		}
	})

	t.Run("wrong current password", func(t *testing.T) {
		other, otherPassword := th.createUser(t, "viewer", false)
		_ = otherPassword
		r := doJSON(t, http.MethodPost, "/api/me/password", map[string]string{
			"currentPassword": "totally-wrong",
			"newPassword":     "new-password-456",
		})
		rr := th.authedRequest(t, r, other.ID, other.Role, th.H.ChangePassword)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		other, _ := th.createUser(t, "viewer", false)
		r := doJSON(t, http.MethodPost, "/api/me/password", map[string]string{"currentPassword": "x"})
		rr := th.authedRequest(t, r, other.ID, other.Role, th.H.ChangePassword)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestResetPassword(t *testing.T) {
	th := newTestHandler(t)

	t.Run("happy path", func(t *testing.T) {
		user, _ := th.createUser(t, "viewer", false)
		r := doJSON(t, http.MethodPost, "/api/users/"+user.ID+"/reset-password", map[string]string{"newPassword": "brand-new-password"})
		r.SetPathValue("id", user.ID)

		rr := httptest.NewRecorder()
		th.H.ResetPassword(rr, r)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}

		loginReq := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{
			"username": user.Username,
			"password": "brand-new-password",
		})
		loginRR := httptest.NewRecorder()
		th.H.Login(loginRR, loginReq)
		if loginRR.Code != http.StatusOK {
			t.Fatalf("login with new password: status = %d, body=%s", loginRR.Code, loginRR.Body.String())
		}
	})

	t.Run("missing id", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/users//reset-password", map[string]string{"newPassword": "x"})
		rr := httptest.NewRecorder()
		th.H.ResetPassword(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/users/nope/reset-password", map[string]string{"newPassword": "brand-new-password"})
		r.SetPathValue("id", "nope")
		rr := httptest.NewRecorder()
		th.H.ResetPassword(rr, r)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("empty new password", func(t *testing.T) {
		user, _ := th.createUser(t, "viewer", false)
		r := doJSON(t, http.MethodPost, "/api/users/"+user.ID+"/reset-password", map[string]string{"newPassword": ""})
		r.SetPathValue("id", user.ID)
		rr := httptest.NewRecorder()
		th.H.ResetPassword(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestCreateUser(t *testing.T) {
	th := newTestHandler(t)

	t.Run("happy path", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/users", map[string]any{
			"username": "newuser1",
			"password": "some-password",
			"role":     "viewer",
		})
		rr := httptest.NewRecorder()
		th.H.CreateUser(rr, r)
		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusCreated, rr.Body.String())
		}
	})

	t.Run("duplicate username", func(t *testing.T) {
		body := map[string]any{"username": "dupuser", "password": "some-password", "role": "viewer"}
		r1 := doJSON(t, http.MethodPost, "/api/users", body)
		th.H.CreateUser(httptest.NewRecorder(), r1)

		r2 := doJSON(t, http.MethodPost, "/api/users", body)
		rr2 := httptest.NewRecorder()
		th.H.CreateUser(rr2, r2)
		if rr2.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", rr2.Code, http.StatusConflict)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/users", map[string]any{"username": "onlyusername"})
		rr := httptest.NewRecorder()
		th.H.CreateUser(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid role", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/users", map[string]any{
			"username": "badrole", "password": "some-password", "role": "superadmin",
		})
		rr := httptest.NewRecorder()
		th.H.CreateUser(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("dashboard defaults requires operator", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/users", map[string]any{
			"username": "viewerwithperm", "password": "some-password", "role": "viewer",
			"canManageDashboardDefaults": true,
		})
		rr := httptest.NewRecorder()
		th.H.CreateUser(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	th := newTestHandler(t)

	t.Run("happy path", func(t *testing.T) {
		user, _ := th.createUser(t, "viewer", false)
		r := doJSON(t, http.MethodPatch, "/api/users/"+user.ID, map[string]any{"displayName": "New Name"})
		r.SetPathValue("id", user.ID)
		rr := httptest.NewRecorder()
		th.H.UpdateUser(rr, r)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})

	t.Run("missing id", func(t *testing.T) {
		r := doJSON(t, http.MethodPatch, "/api/users/", map[string]any{"displayName": "x"})
		rr := httptest.NewRecorder()
		th.H.UpdateUser(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("no fields", func(t *testing.T) {
		user, _ := th.createUser(t, "viewer", false)
		r := doJSON(t, http.MethodPatch, "/api/users/"+user.ID, map[string]any{})
		r.SetPathValue("id", user.ID)
		rr := httptest.NewRecorder()
		th.H.UpdateUser(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid role", func(t *testing.T) {
		user, _ := th.createUser(t, "viewer", false)
		r := doJSON(t, http.MethodPatch, "/api/users/"+user.ID, map[string]any{"role": "superadmin"})
		r.SetPathValue("id", user.ID)
		rr := httptest.NewRecorder()
		th.H.UpdateUser(rr, r)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("disabling revokes sessions", func(t *testing.T) {
		user, password := th.createUser(t, "viewer", false)
		loginReq := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{"username": user.Username, "password": password})
		loginRR := httptest.NewRecorder()
		th.H.Login(loginRR, loginReq)
		if loginRR.Code != http.StatusOK {
			t.Fatalf("login: status = %d", loginRR.Code)
		}

		sessions, err := th.Store.ListUserSessions(req(t).Context(), user.ID)
		if err != nil || len(sessions) == 0 {
			t.Fatalf("expected a session to exist, err=%v sessions=%v", err, sessions)
		}

		r := doJSON(t, http.MethodPatch, "/api/users/"+user.ID, map[string]any{"disabled": true})
		r.SetPathValue("id", user.ID)
		rr := httptest.NewRecorder()
		th.H.UpdateUser(rr, r)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		sessions, err = th.Store.ListUserSessions(req(t).Context(), user.ID)
		if err != nil {
			t.Fatalf("list sessions: %v", err)
		}
		if len(sessions) != 0 {
			t.Fatalf("expected sessions revoked after disable, got %d", len(sessions))
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		r := doJSON(t, http.MethodPatch, "/api/users/does-not-exist", map[string]any{"displayName": "x"})
		r.SetPathValue("id", "does-not-exist")
		rr := httptest.NewRecorder()
		th.H.UpdateUser(rr, r)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	th := newTestHandler(t)

	t.Run("happy path", func(t *testing.T) {
		operator, _ := th.createUser(t, "operator", false)
		user, _ := th.createUser(t, "viewer", false)
		r := httptest.NewRequest(http.MethodDelete, "/api/users/"+user.ID, nil)
		r.SetPathValue("id", user.ID)
		rr := th.authedRequest(t, r, operator.ID, operator.Role, th.H.DeleteUser)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}
	})

	t.Run("cannot delete self", func(t *testing.T) {
		user, _ := th.createUser(t, "operator", false)
		r := httptest.NewRequest(http.MethodDelete, "/api/users/"+user.ID, nil)
		r.SetPathValue("id", user.ID)
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.DeleteUser)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing id", func(t *testing.T) {
		operator, _ := th.createUser(t, "operator", false)
		r := httptest.NewRequest(http.MethodDelete, "/api/users/", nil)
		rr := th.authedRequest(t, r, operator.ID, operator.Role, th.H.DeleteUser)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		operator, _ := th.createUser(t, "operator", false)
		r := httptest.NewRequest(http.MethodDelete, "/api/users/does-not-exist", nil)
		r.SetPathValue("id", "does-not-exist")
		rr := th.authedRequest(t, r, operator.ID, operator.Role, th.H.DeleteUser)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}
