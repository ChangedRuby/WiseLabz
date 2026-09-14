package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// refreshCookie returns the live refresh_token cookie from a response —
// setRefreshCookie always emits two Set-Cookie headers (an empty
// MaxAge:-1 deletion for the old path-scoped cookie, then the real one), and
// a plain replay of every cookie onto a new request lets the deletion
// cookie win, so callers must pick the non-empty one explicitly.
func refreshCookie(t *testing.T, rr *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range rr.Result().Cookies() {
		if c.Name == "refresh_token" && c.Value != "" {
			return c
		}
	}
	t.Fatal("no live refresh_token cookie in response")
	return nil
}

func TestRefresh(t *testing.T) {
	th := newTestHandler(t)
	user, password := th.createUser(t, "viewer", false)

	login := func() *httptest.ResponseRecorder {
		r := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{"username": user.Username, "password": password})
		rr := httptest.NewRecorder()
		th.H.Login(rr, r)
		return rr
	}

	t.Run("happy path via cookie", func(t *testing.T) {
		loginRR := login()
		if loginRR.Code != http.StatusOK {
			t.Fatalf("login: status = %d", loginRR.Code)
		}
		r := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		r.AddCookie(refreshCookie(t, loginRR))
		rr := httptest.NewRecorder()
		th.H.Refresh(rr, r)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})

	t.Run("happy path via body", func(t *testing.T) {
		loginRR := login()
		var refreshToken string
		for _, c := range loginRR.Result().Cookies() {
			if c.Name == "refresh_token" && c.Value != "" {
				refreshToken = c.Value
			}
		}
		if refreshToken == "" {
			t.Fatal("no refresh token cookie from login")
		}
		r := doJSON(t, http.MethodPost, "/api/auth/refresh", map[string]string{"refreshToken": refreshToken})
		rr := httptest.NewRecorder()
		th.H.Refresh(rr, r)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})

	t.Run("missing token", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		rr := httptest.NewRecorder()
		th.H.Refresh(rr, r)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("invalid/expired token", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		r.AddCookie(&http.Cookie{Name: "refresh_token", Value: "not-a-real-token"})
		rr := httptest.NewRecorder()
		th.H.Refresh(rr, r)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("disabled user rejected", func(t *testing.T) {
		loginRR := login()
		if err := th.Store.UpdateUser(context.Background(), user.ID, map[string]any{"disabled": true}); err != nil {
			t.Fatalf("disable user: %v", err)
		}
		r := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		r.AddCookie(refreshCookie(t, loginRR))
		rr := httptest.NewRecorder()
		th.H.Refresh(rr, r)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusForbidden)
		}
	})

	t.Run("token already rotated/revoked cannot be reused", func(t *testing.T) {
		user2, password2 := th.createUser(t, "viewer", false)
		r0 := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{"username": user2.Username, "password": password2})
		loginRR := httptest.NewRecorder()
		th.H.Login(loginRR, r0)

		originalCookie := refreshCookie(t, loginRR)

		r1 := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		r1.AddCookie(originalCookie)
		rr1 := httptest.NewRecorder()
		th.H.Refresh(rr1, r1)
		if rr1.Code != http.StatusOK {
			t.Fatalf("first refresh: status = %d", rr1.Code)
		}

		// Replay the original (now-rotated) refresh token cookie.
		r2 := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		r2.AddCookie(originalCookie)
		rr2 := httptest.NewRecorder()
		th.H.Refresh(rr2, r2)
		if rr2.Code != http.StatusUnauthorized {
			t.Fatalf("replayed refresh: status = %d, want %d", rr2.Code, http.StatusUnauthorized)
		}
	})
}

func TestLogout(t *testing.T) {
	th := newTestHandler(t)
	user, password := th.createUser(t, "viewer", false)

	r0 := doJSON(t, http.MethodPost, "/api/auth/login", map[string]string{"username": user.Username, "password": password})
	loginRR := httptest.NewRecorder()
	th.H.Login(loginRR, r0)

	r := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	r.AddCookie(refreshCookie(t, loginRR))
	rr := th.authedRequest(t, r, user.ID, user.Role, th.H.Logout)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
	}

	sessions, err := th.Store.ListUserSessions(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("expected sessions cleared after logout, got %d", len(sessions))
	}

	var cleared bool
	for _, c := range rr.Result().Cookies() {
		if c.Name == "refresh_token" && c.MaxAge < 0 {
			cleared = true
		}
	}
	if !cleared {
		t.Fatal("expected logout to clear the refresh_token cookie")
	}
}

func TestElevate(t *testing.T) {
	th := newTestHandler(t)
	user, password := th.createUser(t, "operator", false)

	t.Run("happy path", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/auth/elevate", map[string]string{"password": password, "action": "connector.delete"})
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.Elevate)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/auth/elevate", map[string]string{"password": "wrong", "action": "connector.delete"})
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.Elevate)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		r := doJSON(t, http.MethodPost, "/api/auth/elevate", map[string]string{"password": password})
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.Elevate)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})
}

func TestDeleteSession(t *testing.T) {
	th := newTestHandler(t)
	user, _ := th.createUser(t, "viewer", false)
	other, _ := th.createUser(t, "viewer", false)

	session := &store.Session{UserID: user.ID, TokenHash: "hash-a", UserAgent: "test"}
	if err := th.Store.CreateSession(context.Background(), session); err != nil {
		t.Fatalf("create session: %v", err)
	}

	t.Run("cannot delete another user's session", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/me/sessions/"+session.ID, nil)
		r.SetPathValue("id", session.ID)
		rr := th.authedRequest(t, r, other.ID, other.Role, th.H.DeleteSession)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusForbidden)
		}
	})

	t.Run("missing id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/me/sessions/", nil)
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.DeleteSession)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown session", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/me/sessions/does-not-exist", nil)
		r.SetPathValue("id", "does-not-exist")
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.DeleteSession)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/me/sessions/"+session.ID, nil)
		r.SetPathValue("id", session.ID)
		rr := th.authedRequest(t, r, user.ID, user.Role, th.H.DeleteSession)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusNoContent, rr.Body.String())
		}
	})
}
