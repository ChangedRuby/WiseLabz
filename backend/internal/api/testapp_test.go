package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// testApp wires a real chi router (via api.NewRouter) to a fresh, migrated
// SQLite database, so role-check and elevation middleware run exactly as in
// production instead of being unit-tested in isolation.
type testApp struct {
	Router    http.Handler
	Store     *store.Store
	JWT       *auth.Service
	Config    *config.Config
	Scheduler *scheduler.Runner
	BackupDir string
	WSHub     *ws.Hub
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()
	return newTestAppWithBackupDir(t, t.TempDir())
}

// newTestAppWithBackupDir is like newTestApp but lets the caller pick the
// backup directory — used by tests that need it to be unwritable/uncreatable.
func newTestAppWithBackupDir(t *testing.T, backupDir string) *testApp {
	t.Helper()

	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"
	db, err := sql.Open("sqlite", dsn)
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

	// Initialize backup schedule with defaults
	backupSched := store.BackupSchedule{
		CronExpr:    "0 3 * * *",
		MaxBackups:  14,
		MaxAgeHours: 720,
		Enabled:     true,
	}
	if err := s.UpsertBackupSchedule(context.Background(), backupSched); err != nil {
		t.Fatalf("upsert backup schedule: %v", err)
	}

	jwtSvc := auth.NewService("test-secret", 15*time.Minute, 24*time.Hour)
	cfg := &config.Config{
		Server: config.Server{Origin: "http://localhost:5173"},
		Auth:   config.AuthSettings{Secret: "test-secret"},
		// A fixed valid base64 32-byte AES-256 key, so tests exercise real
		// connector-secret encryption/decryption rather than skipping it.
		Encryption: config.EncryptionSettings{Key: "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="},
	}

	jobRunner := scheduler.New(logger)

	aiRegistry := ai.NewRegistry()
	ai.RegisterOpenAICompatible(aiRegistry)
	ai.RegisterClaude(aiRegistry)

	embedRegistry := ai.NewEmbedRegistry()
	ai.RegisterOllamaEmbedder(embedRegistry)
	ai.RegisterOpenAIEmbedder(embedRegistry)

	wsHub := ws.NewHub(cfg.Server.Origin)
	go wsHub.Run(context.Background())

	router := api.NewRouter(api.Config{
		WSHub:         wsHub,
		Store:         s,
		JWT:           jwtSvc,
		Config:        cfg,
		DocEngine:     doc.NewEngine(s),
		SyncEngine:    sync.NewEngine(s, nil, nil, nil, cfg.Encryption.Key),
		Scheduler:     jobRunner,
		BackupDir:     backupDir,
		AIRegistry:    aiRegistry,
		EmbedRegistry: embedRegistry,
	})

	return &testApp{Router: router, Store: s, JWT: jwtSvc, Config: cfg, Scheduler: jobRunner, BackupDir: backupDir, WSHub: wsHub}
}

// user seeds a local user with the given flat instance role ("operator" or
// "viewer", the spelling every existing test uses — "operator" maps to the
// instance-admin flag) and returns its ID and a valid access token. Carries
// no per-connector access; grant that separately with connectorGrant for
// tests exercising connector-scoped routes.
func (a *testApp) user(t *testing.T, role string) (userID, accessToken string) {
	t.Helper()

	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	instanceAdmin := role == "operator"
	instanceAdminRole := "user"
	if instanceAdmin {
		instanceAdminRole = "admin"
	}
	u := &store.User{
		Username:          "user-" + uuid.New().String(),
		DisplayName:       "Test User",
		InstanceAdminRole: instanceAdminRole,
		AuthSource:        "local",
		PasswordHash:      hash,
	}
	if err := a.Store.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	pair, err := a.JWT.IssuePair(u.ID, instanceAdmin)
	if err != nil {
		t.Fatalf("issue pair: %v", err)
	}
	return u.ID, pair.AccessToken
}

// connectorGrant gives userID the given per-connector role ("viewer" or
// "operator") on connectorID, for tests exercising connector-scoped routes.
func (a *testApp) connectorGrant(t *testing.T, userID, connectorID, role string) {
	t.Helper()
	if _, err := a.Store.UpsertConnectorGrant(context.Background(), userID, connectorID, role); err != nil {
		t.Fatalf("grant connector role: %v", err)
	}
}

func (a *testApp) elevationToken(t *testing.T, userID, action string) string {
	t.Helper()
	tok, err := a.JWT.IssueElevation(userID, action)
	if err != nil {
		t.Fatalf("issue elevation: %v", err)
	}
	return tok.Token
}

// req performs an HTTP request against the router. If body is non-nil it is
// JSON-encoded. If token is non-empty it's sent as a bearer token.
func (a *testApp) req(t *testing.T, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	return a.serve(a.newRequest(t, method, path, body, token))
}

// newRequest builds a request (JSON-encoding body if non-nil) with an
// optional bearer token, without serving it — for tests that need to set
// extra headers (e.g. X-Elevation-Token) before serving.
func (a *testApp) newRequest(t *testing.T, method, path string, body any, token string) *http.Request {
	t.Helper()

	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		r = bytes.NewReader(data)
	}

	request := httptest.NewRequest(method, path, r)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	return request
}

func (a *testApp) serve(r *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, r)
	return rec
}

// Backup API Integration Tests

func TestBackupScheduleGetDefaults(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	rec := app.req(t, "GET", "/api/system/backup/schedule", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /schedule: got %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		CronExpr    string `json:"cronExpr"`
		MaxBackups  int    `json:"maxBackups"`
		MaxAgeHours int    `json:"maxAgeHours"`
		Enabled     bool   `json:"enabled"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.CronExpr == "" {
		t.Error("CronExpr should be initialized")
	}
	if body.MaxBackups <= 0 {
		t.Error("MaxBackups should be positive")
	}
}

func TestBackupScheduleUpdate(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	newSched := map[string]any{
		"cronExpr":    "0 4 * * *",
		"maxBackups":  7,
		"maxAgeHours": 360,
		"enabled":     false,
	}

	rec := app.req(t, "PUT", "/api/system/backup/schedule", newSched, token)
	if rec.Code != http.StatusOK {
		t.Logf("Response body: %s", rec.Body.String())
		t.Fatalf("PUT /schedule: got %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		CronExpr    string `json:"cronExpr"`
		MaxBackups  int    `json:"maxBackups"`
		MaxAgeHours int    `json:"maxAgeHours"`
		Enabled     bool   `json:"enabled"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.CronExpr != "0 4 * * *" {
		t.Errorf("CronExpr: got %s, want 0 4 * * *", body.CronExpr)
	}
	if body.MaxBackups != 7 {
		t.Errorf("MaxBackups: got %d, want 7", body.MaxBackups)
	}
	if body.Enabled != false {
		t.Errorf("Enabled: got %v, want false", body.Enabled)
	}
}

func TestBackupScheduleUpdateInvalidCron(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	newSched := map[string]any{
		"cronExpr":    "invalid cron",
		"maxBackups":  7,
		"maxAgeHours": 360,
		"enabled":     true,
	}

	rec := app.req(t, "PUT", "/api/system/backup/schedule", newSched, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("PUT /schedule with invalid cron: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestBackupListRunsEmpty(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	rec := app.req(t, "GET", "/api/system/backup/runs", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /runs: got %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		Runs   []any `json:"runs"`
		Total  int   `json:"total"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Total != 0 {
		t.Errorf("total: got %d, want 0", body.Total)
	}
	if len(body.Runs) != 0 {
		t.Errorf("runs length: got %d, want 0", len(body.Runs))
	}
}

func TestBackupCreateManualRun(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	rec := app.req(t, "POST", "/api/system/backup/run", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /run: got %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		ID          string `json:"id"`
		TriggeredBy string `json:"triggeredBy"`
		FilePath    string `json:"filePath"`
		SizeBytes   int64  `json:"sizeBytes"`
		CreatedAt   string `json:"createdAt"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.ID == "" {
		t.Error("ID should not be empty")
	}
	if body.TriggeredBy != "manual" {
		t.Errorf("TriggeredBy: got %s, want manual", body.TriggeredBy)
	}
	if body.FilePath == "" {
		t.Error("FilePath should not be empty")
	}
	if body.SizeBytes <= 0 {
		t.Error("SizeBytes should be positive")
	}

	// Verify the run appears in the list
	listRec := app.req(t, "GET", "/api/system/backup/runs", nil, token)
	var listBody struct {
		Runs  []any `json:"runs"`
		Total int   `json:"total"`
	}
	if err := json.NewDecoder(listRec.Body).Decode(&listBody); err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if listBody.Total != 1 {
		t.Errorf("total after manual run: got %d, want 1", listBody.Total)
	}
}

// TestBackupScheduleUpdateDoesNotLeakSchedulerJobs is a regression test: the
// startup registration in api.NewRouter (via Handler.InitBackupJob) and every
// subsequent PUT /schedule must all route through the same BackupJobID
// bookkeeping. If any of them added a job directly against the scheduler
// without going through that bookkeeping, repeated updates would stack
// duplicate "backup" cron entries instead of replacing the one entry.
func TestBackupScheduleUpdateDoesNotLeakSchedulerJobs(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	if got := app.Scheduler.EntryCount(); got != 1 {
		t.Fatalf("scheduler entries after startup: got %d, want 1 (just the seeded backup job)", got)
	}

	for i, expr := range []string{"0 4 * * *", "0 5 * * *", "0 6 * * *"} {
		sched := map[string]any{
			"cronExpr":    expr,
			"maxBackups":  7,
			"maxAgeHours": 360,
			"enabled":     true,
		}
		rec := app.req(t, "PUT", "/api/system/backup/schedule", sched, token)
		if rec.Code != http.StatusOK {
			t.Fatalf("PUT /schedule #%d: got %d, want %d", i, rec.Code, http.StatusOK)
		}
		if got := app.Scheduler.EntryCount(); got != 1 {
			t.Fatalf("scheduler entries after PUT #%d: got %d, want 1 (old job must be removed before adding the new one)", i, got)
		}
	}

	// Disabling the schedule should deregister the job entirely, not leave a
	// stale disabled entry behind.
	rec := app.req(t, "PUT", "/api/system/backup/schedule", map[string]any{
		"cronExpr":    "0 7 * * *",
		"maxBackups":  7,
		"maxAgeHours": 360,
		"enabled":     false,
	}, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT /schedule (disable): got %d, want %d", rec.Code, http.StatusOK)
	}
	if got := app.Scheduler.EntryCount(); got != 0 {
		t.Fatalf("scheduler entries after disabling: got %d, want 0", got)
	}
}

// TestBackupRoutesRequireOperatorRole verifies all four new backup endpoints
// are gated behind the operator role, not merely authentication.
func TestBackupRoutesRequireOperatorRole(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	cases := []struct {
		method string
		path   string
		body   any
	}{
		{"GET", "/api/system/backup/schedule", nil},
		{"PUT", "/api/system/backup/schedule", map[string]any{"cronExpr": "0 3 * * *", "maxBackups": 1, "maxAgeHours": 1, "enabled": true}},
		{"GET", "/api/system/backup/runs", nil},
		{"POST", "/api/system/backup/run", nil},
	}
	for _, c := range cases {
		rec := app.req(t, c.method, c.path, c.body, viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Errorf("%s %s as viewer: got %d, want %d", c.method, c.path, rec.Code, http.StatusForbidden)
		}

		unauth := app.req(t, c.method, c.path, c.body, "")
		if unauth.Code != http.StatusUnauthorized {
			t.Errorf("%s %s unauthenticated: got %d, want %d", c.method, c.path, unauth.Code, http.StatusUnauthorized)
		}
	}
}

// TestBackupCreateManualRunFailsWhenDirNotCreatable verifies POST /run
// surfaces an error (not a panic or a false-success) when the configured
// backup directory can't be created, e.g. because its parent is a file
// instead of a directory.
func TestBackupCreateManualRunFailsWhenDirNotCreatable(t *testing.T) {
	base := t.TempDir()
	blocker := base + "/blocker"
	if err := os.WriteFile(blocker, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("create blocker file: %v", err)
	}
	// backupDir's parent path component is a regular file, so os.MkdirAll
	// inside ExportToFile cannot create it.
	app := newTestAppWithBackupDir(t, blocker+"/backups")
	_, token := app.user(t, "operator")

	rec := app.req(t, "POST", "/api/system/backup/run", nil, token)
	if rec.Code == http.StatusOK {
		t.Fatalf("POST /run with uncreatable dir: got %d, want a non-2xx error", rec.Code)
	}

	// No run should have been recorded in the database.
	runs, total, err := app.Store.ListBackupRuns(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("ListBackupRuns: %v", err)
	}
	if total != 0 || len(runs) != 0 {
		t.Fatalf("expected no backup run recorded after failed export, got total=%d runs=%d", total, len(runs))
	}
}
