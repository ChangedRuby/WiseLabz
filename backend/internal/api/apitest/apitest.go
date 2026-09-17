// Package apitest provides a minimal shared harness for internal/api/*
// subpackage handler tests: a real migrated SQLite-backed store and a real
// auth middleware/JWT chain, so handler tests exercise auth.UserIDFromContext
// and auth.InstanceAdminFromContext exactly as production does instead of
// faking context values that only auth's own unexported keys can set.
package apitest

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite" // sqlite driver registration; this package is test-only despite the non-_test name

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// NewStore builds a fresh, migrated, initialized SQLite-backed store for a test.
func NewStore(t *testing.T) *store.Store {
	t.Helper()

	dir := t.TempDir()
	db, err := sql.Open("sqlite", "file:"+dir+"/test.db?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := store.RunMigrations(db, "sqlite", logger); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	s := store.New(db, "sqlite")
	if err := s.Init(context.Background(), "admin-seed-pw-1234"); err != nil {
		t.Fatalf("store init: %v", err)
	}
	return s
}

// NewUser seeds a local user with the given flat instance role ("operator"
// or "viewer", kept as the caller-facing spelling every existing test uses)
// and returns its ID. "operator" maps to the instance-admin flag; it carries
// no per-connector access — grant that separately with GrantConnectorRole.
func NewUser(t *testing.T, s *store.Store, role string) string {
	t.Helper()

	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	u := &store.User{
		Username:          "user-" + uuid.New().String(),
		DisplayName:       "Test User",
		InstanceAdminRole: instanceAdminRole(role),
		AuthSource:        "local",
		PasswordHash:      hash,
	}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u.ID
}

// instanceAdminRole maps the test-facing "operator"/"viewer" spelling onto
// the stored instance_admin_role domain ("admin"/"user").
func instanceAdminRole(role string) string {
	if role == "operator" {
		return "admin"
	}
	return "user"
}

// GrantConnectorRole gives userID the given per-connector role ("viewer" or
// "operator") on connectorID, for tests exercising auth.RequireConnectorRole
// or in-handler UserHasConnectorRole checks.
func GrantConnectorRole(t *testing.T, s *store.Store, userID, connectorID, role string) {
	t.Helper()
	if _, err := s.UpsertConnectorGrant(context.Background(), userID, connectorID, role); err != nil {
		t.Fatalf("grant connector role: %v", err)
	}
}

// JWTService returns a test JWT service backed by a fixed secret.
func JWTService() *auth.Service {
	return auth.NewService("test-secret", 15*time.Minute, 24*time.Hour)
}

// Token issues a valid access token for userID with the given flat instance
// role ("operator"/"viewer", same spelling as NewUser).
func Token(t *testing.T, jwtSvc *auth.Service, userID, role string) string {
	t.Helper()
	pair, err := jwtSvc.IssuePair(userID, role == "operator")
	if err != nil {
		t.Fatalf("issue pair: %v", err)
	}
	return pair.AccessToken
}

// WithAuth wraps h with the real auth middleware (JWT + API-key lookup
// against s), so handlers reading auth.UserIDFromContext /
// auth.InstanceAdminFromContext see real context values.
func WithAuth(jwtSvc *auth.Service, s *store.Store, h http.Handler) http.Handler {
	return auth.AuthMiddleware(jwtSvc, s)(h)
}

// AuthedUser is a convenience bundling a seeded user with a ready-to-use
// authenticated handler wrapper.
func AuthedUser(t *testing.T, s *store.Store, role string, h http.Handler) (userID, token string, wrapped http.Handler) {
	t.Helper()
	jwtSvc := JWTService()
	userID = NewUser(t, s, role)
	token = Token(t, jwtSvc, userID, role)
	wrapped = WithAuth(jwtSvc, s, h)
	return userID, token, wrapped
}
