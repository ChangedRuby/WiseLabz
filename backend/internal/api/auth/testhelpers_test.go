package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	upstreamauth "github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// testHandler wires a Handler to a fresh, migrated SQLite database, matching
// the setup internal/api/testapp_test.go uses for router-level tests.
type testHandler struct {
	H     *Handler
	Store *store.Store
	JWT   *upstreamauth.Service
}

func newTestHandler(t *testing.T) *testHandler {
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

	jwtSvc := upstreamauth.NewService("test-secret", 15*time.Minute, 24*time.Hour)
	cfg := &config.Config{
		Server: config.Server{Origin: "http://localhost:5173"},
		Auth:   config.AuthSettings{Secret: "test-secret"},
	}

	return &testHandler{
		H:     NewHandler(s, jwtSvc, cfg),
		Store: s,
		JWT:   jwtSvc,
	}
}

// createUser seeds a local user with the given role/disabled state and
// returns it with the plaintext password used to hash it.
func (th *testHandler) createUser(t *testing.T, role string, disabled bool) (*store.User, string) {
	t.Helper()
	const password = "password123"
	hash, err := upstreamauth.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	u := &store.User{
		Username:     "user-" + uuid.New().String(),
		DisplayName:  "Test User",
		Role:         role,
		AuthSource:   "local",
		PasswordHash: hash,
		Disabled:     disabled,
	}
	if err := th.Store.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u, password
}

// authedRequest runs req through auth.AuthMiddleware (as production does)
// with an access token for userID/role, so handlers reading
// auth.UserIDFromContext/RoleFromContext behave exactly as in prod.
func (th *testHandler) authedRequest(t *testing.T, req *http.Request, userID, role string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	pair, err := th.JWT.IssuePair(userID, role)
	if err != nil {
		t.Fatalf("issue pair: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)

	rr := httptest.NewRecorder()
	upstreamauth.AuthMiddleware(th.JWT, th.Store)(handler).ServeHTTP(rr, req)
	return rr
}
