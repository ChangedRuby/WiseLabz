package auth

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	_ "modernc.org/sqlite"

	upstreamauth "github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/crypto"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/store/storetest"
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

	db, err := sql.Open("sqlite", "file:"+storetest.MigratedSQLite(t)+"?cache=shared")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })

	s := store.New(db, "sqlite")
	if err := s.Init(context.Background(), "admin-seed-pw-1234"); err != nil {
		t.Fatalf("store init: %v", err)
	}

	jwtSvc := upstreamauth.NewService("test-secret", 15*time.Minute, 24*time.Hour)
	cfg := &config.Config{
		Server:     config.Server{Origin: "http://localhost:5173"},
		Auth:       config.AuthSettings{Secret: "test-secret"},
		Encryption: config.EncryptionSettings{Key: "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="},
	}

	return &testHandler{
		H:     NewHandler(s, jwtSvc, cfg),
		Store: s,
		JWT:   jwtSvc,
	}
}

// createUser seeds a local user with the given flat instance role
// ("operator" or "viewer", kept as the caller-facing spelling every existing
// test uses; "operator" maps to the instance-admin flag) and disabled
// state, and returns it with the plaintext password used to hash it.
func (th *testHandler) createUser(t *testing.T, role string, disabled bool) (*store.User, string) {
	t.Helper()
	const password = "password123"
	hash, err := upstreamauth.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	u := &store.User{
		Username:          "user-" + uuid.New().String(),
		DisplayName:       "Test User",
		InstanceAdminRole: instanceAdminRoleFor(role),
		AuthSource:        "local",
		PasswordHash:      hash,
		Disabled:          disabled,
	}
	if err := th.Store.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u, password
}

// enrollTOTP confirms a TOTP factor for userID directly against the store
// (bypassing the enrollment HTTP endpoints) and returns the raw secret so
// tests can generate valid codes with totpCodeAt.
func (th *testHandler) enrollTOTP(t *testing.T, userID string) string {
	t.Helper()
	secret, _, err := upstreamauth.GenerateTOTPSecret("WiseLabz", userID)
	if err != nil {
		t.Fatalf("GenerateTOTPSecret() error: %v", err)
	}
	key, err := crypto.DecodeKey(th.H.Config.Encryption.Key)
	if err != nil {
		t.Fatalf("DecodeKey() error: %v", err)
	}
	encrypted, err := crypto.Encrypt(secret, key)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}
	factor, err := th.Store.CreatePendingTOTP(context.Background(), userID, "Authenticator", encrypted)
	if err != nil {
		t.Fatalf("CreatePendingTOTP() error: %v", err)
	}
	if _, err := th.Store.ConfirmFactor(context.Background(), factor.ID); err != nil {
		t.Fatalf("ConfirmFactor() error: %v", err)
	}
	return secret
}

// totpCodeAt generates a valid TOTP code for secret at the given time, using
// the same SHA1/6-digit/30s parameters as internal/auth.GenerateTOTPSecret.
func totpCodeAt(t *testing.T, secret string, when time.Time) string {
	t.Helper()
	code, err := totp.GenerateCodeCustom(secret, when, totp.ValidateOpts{
		Period:    30,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		t.Fatalf("GenerateCodeCustom() error: %v", err)
	}
	return code
}

// instanceAdminRoleFor maps the test-facing "operator"/"viewer" spelling
// onto the stored instance_admin_role domain ("admin"/"user").
func instanceAdminRoleFor(role string) string {
	if role == "operator" {
		return "admin"
	}
	return "user"
}

// authedRequest runs req through auth.AuthMiddleware (as production does)
// with an access token for userID and the given flat instance role. Accepts
// either spelling — "operator"/"admin" both mean instance-admin — so callers
// can pass either a literal createUser() role or a seeded user's own
// InstanceAdminRole field back in.
func (th *testHandler) authedRequest(t *testing.T, req *http.Request, userID, role string, handler http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	pair, err := th.JWT.IssuePair(userID, role == "operator" || role == "admin")
	if err != nil {
		t.Fatalf("issue pair: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)

	rr := httptest.NewRecorder()
	upstreamauth.AuthMiddleware(th.JWT, th.Store)(handler).ServeHTTP(rr, req)
	return rr
}
