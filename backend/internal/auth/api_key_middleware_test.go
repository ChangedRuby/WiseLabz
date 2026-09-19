package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testAPIKeyChecker struct {
	claims *APIKeyClaims
	called bool
}

func (c *testAPIKeyChecker) LookupAPIKey(context.Context, string) (*APIKeyClaims, error) {
	return c.claims, nil
}

func (c *testAPIKeyChecker) TouchAPIKeyLastUsed(context.Context, string) error {
	c.called = true
	return nil
}

func TestAuthMiddlewareAPIKeyLifecycle(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	checker := &testAPIKeyChecker{claims: &APIKeyClaims{KeyID: "key-1", UserID: "user-1", InstanceAdmin: true}}
	handler := AuthMiddleware(svc, checker)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer wlz_secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !checker.called {
		t.Fatalf("API key auth status = %d, touched = %v; want 200 and touch", rec.Code, checker.called)
	}
}

// TestAuthMiddlewareAcceptsNonAdminAPIKey guards against a regression to the
// old string-role check: InstanceAdmin's zero value (false) is a valid,
// ordinary non-admin key, not an invalid/rejected state.
func TestAuthMiddlewareAcceptsNonAdminAPIKey(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	checker := &testAPIKeyChecker{claims: &APIKeyClaims{KeyID: "key", UserID: "user", InstanceAdmin: false}}
	handler := AuthMiddleware(svc, checker)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer wlz_secret")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestAuthMiddlewareRejectsExpiredAndRevokedAPIKeys(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	for _, tc := range []struct {
		name   string
		claims *APIKeyClaims
	}{
		{name: "expired", claims: &APIKeyClaims{KeyID: "key", UserID: "user", InstanceAdmin: false, ExpiresAt: time.Now().UTC().Add(-time.Minute).Format(time.RFC3339)}},
		{name: "revoked", claims: &APIKeyClaims{KeyID: "key", UserID: "user", InstanceAdmin: false, RevokedAt: time.Now().UTC().Format(time.RFC3339)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checker := &testAPIKeyChecker{claims: tc.claims}
			handler := AuthMiddleware(svc, checker)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				t.Error("handler should not be called")
			}))
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer wlz_secret")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", rec.Code)
			}
		})
	}
}

func TestAuthMiddlewareThrottlesAPIKeyLastUsed(t *testing.T) {
	svc := NewService("test-secret", time.Minute, time.Hour)
	for _, tc := range []struct {
		name     string
		lastUsed string
		touch    bool
	}{
		{"unused", "", true},
		{"recent", time.Now().UTC().Format(time.RFC3339), false},
		{"old", time.Now().UTC().Add(-2 * time.Minute).Format(time.RFC3339), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checker := &testAPIKeyChecker{claims: &APIKeyClaims{KeyID: "key", UserID: "user", LastUsedAt: tc.lastUsed}}
			handler := AuthMiddleware(svc, checker)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer wlz_secret")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK || checker.called != tc.touch {
				t.Fatalf("status = %d, touched = %v, want 200, %v", rec.Code, checker.called, tc.touch)
			}
			// A warm last-used timestamp must never bypass revocation.
			checker.claims.RevokedAt = time.Now().UTC().Format(time.RFC3339)
			rec = httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("revoked key status = %d", rec.Code)
			}
		})
	}
}
