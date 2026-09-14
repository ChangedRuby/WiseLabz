package diagnostics_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/diagnostics"
	"github.com/WiseLabz/wiselabz/internal/store"

	// Register connector implementations so RedactConnectorConfig resolves
	// their secret fields the same way it does in production.
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

const testEncKey = "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	dsn := "file:" + dir + "/test.db?cache=shared"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
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

// TestCollectRedactsConnectorSecrets verifies Collect's bundle never carries
// a connector's secret config fields, matching backup.Export's redaction.
func TestCollectRedactsConnectorSecrets(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	proxmoxCfg, err := store.MarshalConnectorConfig("proxmox", map[string]any{
		"url":          "https://pve.example.com:8006",
		"token_id":     "root@pam!monitoring",
		"token_secret": "super-secret-token",
	}, testEncKey)
	if err != nil {
		t.Fatalf("marshal proxmox config: %v", err)
	}
	if err := s.CreateConnector(ctx, &store.ConnectorRecord{
		Name: "pve1", Category: "virtualization", Type: "proxmox", URL: "https://pve.example.com:8006",
		ConfigData: proxmoxCfg,
	}); err != nil {
		t.Fatalf("create proxmox connector: %v", err)
	}

	cfg := &config.Config{Sync: config.SyncSettings{Schedule: "@daily"}}

	b, err := diagnostics.Collect(ctx, s, cfg)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	raw, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	if strings.Contains(string(raw), "super-secret-token") {
		t.Error("diagnostics bundle contains the proxmox token_secret value")
	}

	if len(b.SanitizedConfig.Connectors) != 1 {
		t.Fatalf("len(connectors) = %d, want 1", len(b.SanitizedConfig.Connectors))
	}
	got := b.SanitizedConfig.Connectors[0]
	if got.Name != "pve1" || got.Type != "proxmox" {
		t.Errorf("connector = %+v, want pve1/proxmox", got)
	}
	if strings.Contains(got.ConfigData, "super-secret-token") {
		t.Error("connector ConfigData still contains token_secret")
	}
}

// TestCollectIncludesHealthVersionsAndSchedule verifies the non-connector
// parts of the bundle (health, versions, sync schedule, OIDC summaries) are
// populated from the running config and store.
func TestCollectIncludesHealthVersionsAndSchedule(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	cfg := &config.Config{
		Sync: config.SyncSettings{Schedule: "0 3 * * *"},
		Auth: config.AuthSettings{
			OIDC: []config.OIDCProvider{
				{ID: "authentik", DisplayName: "Authentik", ClientSecret: "should-not-appear"},
			},
		},
	}

	b, err := diagnostics.Collect(ctx, s, cfg)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if b.Health.Status != "ok" {
		t.Errorf("Health.Status = %q, want ok", b.Health.Status)
	}
	if b.SanitizedConfig.SyncSchedule != "0 3 * * *" {
		t.Errorf("SyncSchedule = %q, want 0 3 * * *", b.SanitizedConfig.SyncSchedule)
	}
	if len(b.SanitizedConfig.AuthProviders.OIDC) != 1 || b.SanitizedConfig.AuthProviders.OIDC[0].ID != "authentik" {
		t.Errorf("OIDC summaries = %+v, want one entry with id=authentik", b.SanitizedConfig.AuthProviders.OIDC)
	}
	if !b.SanitizedConfig.AuthProviders.Local {
		t.Error("AuthProviders.Local = false, want true")
	}

	raw, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	if strings.Contains(string(raw), "should-not-appear") {
		t.Error("diagnostics bundle leaked an OIDC client secret")
	}
	if b.GeneratedAt == "" {
		t.Error("GeneratedAt is empty")
	}
}

// TestCollectListsRecentFailures verifies failed sync runs and deliveries
// surface in the bundle so support can see recent failure history.
func TestCollectListsRecentFailures(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	c := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	if err := s.CreateSyncRun(ctx, &store.SyncRunRecord{ConnectorID: c.ID, StartedAt: "2024-01-01T00:00:00Z", Status: store.SyncRunStatusError, Error: "boom"}); err != nil {
		t.Fatalf("CreateSyncRun() error: %v", err)
	}

	cfg := &config.Config{}
	b, err := diagnostics.Collect(ctx, s, cfg)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	if len(b.RecentFailures.SyncRuns) != 1 {
		t.Fatalf("len(RecentFailures.SyncRuns) = %d, want 1", len(b.RecentFailures.SyncRuns))
	}
	if b.RecentFailures.SyncRuns[0].ConnectorID != c.ID {
		t.Errorf("failed sync run connector = %q, want %q", b.RecentFailures.SyncRuns[0].ConnectorID, c.ID)
	}
}

// TestCheckHealthReportsDegradedOnClosedDB verifies CheckHealth reports
// "degraded" when the database ping fails, the branch Collect relies on to
// surface DB outages in the bundle.
func TestCheckHealthReportsDegradedOnClosedDB(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	if err := s.DB().Close(); err != nil {
		t.Fatalf("close db: %v", err)
	}

	h := diagnostics.CheckHealth(ctx, s.DB())
	if h.Status != "degraded" {
		t.Errorf("Health.Status = %q, want degraded", h.Status)
	}
	if len(h.Components) != 1 || h.Components[0].Status != "down" {
		t.Errorf("Components = %+v, want database/down", h.Components)
	}
}
