package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// Change to a temp directory with no config file
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)          //nolint:errcheck
	defer os.Chdir(oldDir) //nolint:errcheck

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("server.host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("server.port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.DB.Driver != "sqlite" {
		t.Errorf("db.driver = %q, want sqlite", cfg.DB.Driver)
	}
	if cfg.DB.MaxOpenConns != 20 || cfg.DB.MaxIdleConns != 5 ||
		cfg.DB.ConnMaxLifetime() != 30*time.Minute || cfg.DB.ConnMaxIdleTime() != 5*time.Minute {
		t.Errorf("unexpected db pool defaults: %+v", cfg.DB)
	}
	if cfg.Auth.AccessTokenTTL != 900 {
		t.Errorf("auth.access_token_ttl = %d, want 900", cfg.Auth.AccessTokenTTL)
	}
	if cfg.Sync.Schedule != "0 */6 * * *" {
		t.Errorf("sync.schedule = %q, want %q", cfg.Sync.Schedule, "0 */6 * * *")
	}
	if cfg.Log.Level != "info" {
		t.Errorf("log.level = %q, want info", cfg.Log.Level)
	}
	if cfg.Retention.SnapshotDays != 90 {
		t.Errorf("retention.snapshot_days = %d, want 90", cfg.Retention.SnapshotDays)
	}
	if cfg.Retention.DocVersionDays != 365 {
		t.Errorf("retention.doc_version_days = %d, want 365", cfg.Retention.DocVersionDays)
	}
	if cfg.Retention.AlertDays != 180 {
		t.Errorf("retention.alert_days = %d, want 180", cfg.Retention.AlertDays)
	}
	if cfg.Retention.SyncRunDays != 90 {
		t.Errorf("retention.sync_run_days = %d, want 90", cfg.Retention.SyncRunDays)
	}
	if cfg.Retention.AuditDays != 180 {
		t.Errorf("retention.audit_days = %d, want 180", cfg.Retention.AuditDays)
	}
	if cfg.Retention.HealthCheckDays != 90 {
		t.Errorf("retention.health_check_days = %d, want 90", cfg.Retention.HealthCheckDays)
	}
	if cfg.Retention.CronExpr != "0 0 * * *" {
		t.Errorf("retention.cron_expr = %q, want 0 0 * * *", cfg.Retention.CronExpr)
	}
	if cfg.Quality.CronExpr != "0 0 * * *" {
		t.Errorf("quality.cron_expr = %q, want 0 0 * * *", cfg.Quality.CronExpr)
	}
	if cfg.Sync.PollCronExpr != "*/30 * * * * *" {
		t.Errorf("sync.poll_cron_expr = %q, want */30 * * * * *", cfg.Sync.PollCronExpr)
	}
	if cfg.DocExport.Dir != "./data/docexport" {
		t.Errorf("doc_export.dir = %q, want ./data/docexport", cfg.DocExport.Dir)
	}
	if cfg.DocExport.CronExpr != "0 2 * * *" {
		t.Errorf("doc_export.cron_expr = %q, want 0 2 * * *", cfg.DocExport.CronExpr)
	}
	if cfg.DocExport.Enabled {
		t.Errorf("doc_export.enabled = true, want false")
	}
}

func TestLoadFromYAML(t *testing.T) {
	dir := t.TempDir()

	yamlContent := `
db:
  driver: postgres
  dsn: postgres://user:pass@localhost/wiselabz
server:
  host: 127.0.0.1
  port: 9090
auth:
  secret: test-secret-key
  oidc:
    - id: authentik
      display_name: Authentik
      issuer_url: https://auth.example.com
      client_id: abc123
      client_secret: secret123
      scopes:
        - openid
        - profile
      groups_claim: groups
      group_role_mapping:
        admins: operator
      email_domain_allowlist:
        - example.com
ai:
  enabled: true
  provider: ollama
  model: llama3
  base_url: http://localhost:11434
  mode: suggest_only
sync:
  schedule: "0 */2 * * *"
log:
  level: debug
  format: json
`
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0o644); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	os.Chdir(dir)          //nolint:errcheck
	defer os.Chdir(oldDir) //nolint:errcheck

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.DB.Driver != "postgres" {
		t.Errorf("db.driver = %q, want postgres", cfg.DB.Driver)
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("server.port = %d, want 9090", cfg.Server.Port)
	}
	if cfg.Auth.Secret != "test-secret-key" {
		t.Errorf("auth.secret = %q, want test-secret-key", cfg.Auth.Secret)
	}
	if len(cfg.Auth.OIDC) != 1 {
		t.Fatalf("len(auth.oidc) = %d, want 1", len(cfg.Auth.OIDC))
	}
	if cfg.Auth.OIDC[0].ID != "authentik" {
		t.Errorf("auth.oidc[0].id = %q, want authentik", cfg.Auth.OIDC[0].ID)
	}
	if cfg.Auth.OIDC[0].GroupsClaim != "groups" || cfg.Auth.OIDC[0].GroupRoleMapping["admins"] != "operator" || len(cfg.Auth.OIDC[0].EmailDomainAllowlist) != 1 {
		t.Errorf("auth.oidc[0] claim mapping = %#v", cfg.Auth.OIDC[0])
	}
	if !cfg.AI.Enabled {
		t.Error("ai.enabled = false, want true")
	}
	if cfg.Log.Format != "json" {
		t.Errorf("log.format = %q, want json", cfg.Log.Format)
	}
}

func TestLoadEnvOverride(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)          //nolint:errcheck
	defer os.Chdir(oldDir) //nolint:errcheck

	t.Setenv("WISELABZ_SERVER_PORT", "7070")
	t.Setenv("WISELABZ_DB_DRIVER", "postgres")
	t.Setenv("WISELABZ_AUTH_SECRET", "env-secret")
	t.Setenv("WISELABZ_RETENTION_SNAPSHOT_DAYS", "30")
	t.Setenv("WISELABZ_RETENTION_DOC_VERSION_DAYS", "0")
	t.Setenv("WISELABZ_RETENTION_ALERT_DAYS", "60")
	t.Setenv("WISELABZ_RETENTION_SYNC_RUN_DAYS", "14")
	t.Setenv("WISELABZ_RETENTION_CRON_EXPR", "0 2 * * *")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 7070 {
		t.Errorf("server.port = %d, want 7070", cfg.Server.Port)
	}
	if cfg.DB.Driver != "postgres" {
		t.Errorf("db.driver = %q, want postgres", cfg.DB.Driver)
	}
	if cfg.Auth.Secret != "env-secret" {
		t.Errorf("auth.secret = %q, want env-secret", cfg.Auth.Secret)
	}
	if cfg.Retention.SnapshotDays != 30 {
		t.Errorf("retention.snapshot_days = %d, want 30", cfg.Retention.SnapshotDays)
	}
	if cfg.Retention.DocVersionDays != 0 {
		t.Errorf("retention.doc_version_days = %d, want 0", cfg.Retention.DocVersionDays)
	}
	if cfg.Retention.AlertDays != 60 {
		t.Errorf("retention.alert_days = %d, want 60", cfg.Retention.AlertDays)
	}
	if cfg.Retention.SyncRunDays != 14 {
		t.Errorf("retention.sync_run_days = %d, want 14", cfg.Retention.SyncRunDays)
	}
	if cfg.Retention.CronExpr != "0 2 * * *" {
		t.Errorf("retention.cron_expr = %q, want 0 2 * * *", cfg.Retention.CronExpr)
	}
}

// TestLoadEnvOverrideAllFields asserts every WISELABZ_ env var wired up in
// Load has a working override, including server.{read,write,shutdown}_timeout_seconds
// which previously had no override at all (issue #147).
func TestLoadEnvOverrideAllFields(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)          //nolint:errcheck
	defer os.Chdir(oldDir) //nolint:errcheck

	env := map[string]string{
		"WISELABZ_DB_DRIVER":                       "postgres",
		"WISELABZ_DB_DSN":                          "postgres://x",
		"WISELABZ_DB_MAX_OPEN_CONNS":               "7",
		"WISELABZ_DB_MAX_IDLE_CONNS":               "3",
		"WISELABZ_DB_CONN_MAX_LIFETIME_SECONDS":    "60",
		"WISELABZ_DB_CONN_MAX_IDLE_TIME_SECONDS":   "30",
		"WISELABZ_SERVER_HOST":                     "127.0.0.1",
		"WISELABZ_SERVER_PORT":                     "9090",
		"WISELABZ_SERVER_ORIGIN":                   "https://example.com",
		"WISELABZ_SERVER_TRUSTED_PROXIES":          "10.0.0.0/8",
		"WISELABZ_SERVER_EMBED":                    "true",
		"WISELABZ_SERVER_READ_TIMEOUT_SECONDS":     "5",
		"WISELABZ_SERVER_WRITE_TIMEOUT_SECONDS":    "6",
		"WISELABZ_SERVER_SHUTDOWN_TIMEOUT_SECONDS": "7",
		"WISELABZ_ENCRYPTION_KEY":                  "env-key",
		"WISELABZ_AUTH_SECRET":                     "env-secret",
		"WISELABZ_AUTH_ACCESS_TOKEN_TTL":           "60",
		"WISELABZ_AUTH_REFRESH_TOKEN_TTL":          "120",
		"WISELABZ_AUTH_STEP_UP_FOR_DESTRUCTIVE":    "false",
		"WISELABZ_AI_ENABLED":                      "true",
		"WISELABZ_AI_PROVIDER":                     "openai",
		"WISELABZ_AI_MODEL":                        "gpt-x",
		"WISELABZ_AI_API_KEY":                      "key",
		"WISELABZ_AI_BASE_URL":                     "http://localhost",
		"WISELABZ_AI_MODE":                         "auto_update",
		"WISELABZ_AI_EMBED_PROVIDER":               "openai",
		"WISELABZ_AI_EMBED_MODEL":                  "text-embedding-3-small",
		"WISELABZ_AI_EMBED_API_KEY":                "embed-key",
		"WISELABZ_AI_EMBED_BASE_URL":               "http://embed-host",
		"WISELABZ_SYNC_SCHEDULE":                   "* * * * *",
		"WISELABZ_SYNC_POLL_CRON_EXPR":             "*/5 * * * * *",
		"WISELABZ_QUALITY_CRON_EXPR":               "0 1 * * *",
		"WISELABZ_ROTATION_MAX_AGE_DAYS":           "45",
		"WISELABZ_ROTATION_WARN_DAYS":              "7",
		"WISELABZ_LOG_LEVEL":                       "debug",
		"WISELABZ_LOG_FORMAT":                      "json",
		"WISELABZ_RETENTION_SNAPSHOT_DAYS":         "1",
		"WISELABZ_RETENTION_DOC_VERSION_DAYS":      "2",
		"WISELABZ_RETENTION_ALERT_DAYS":            "3",
		"WISELABZ_RETENTION_SYNC_RUN_DAYS":         "4",
		"WISELABZ_RETENTION_AUDIT_DAYS":            "5",
		"WISELABZ_RETENTION_HEALTH_CHECK_DAYS":     "6",
		"WISELABZ_RETENTION_CRON_EXPR":             "0 4 * * *",
		"WISELABZ_BACKUP_DIR":                      "/tmp/backups",
		"WISELABZ_BACKUP_CRON_EXPR":                "0 5 * * *",
		"WISELABZ_BACKUP_MAX_BACKUPS":              "1",
		"WISELABZ_BACKUP_MAX_AGE_HOURS":            "2",
		"WISELABZ_BACKUP_ENABLED":                  "false",
		"WISELABZ_DOC_EXPORT_DIR":                  "/tmp/docexport",
		"WISELABZ_DOC_EXPORT_CRON_EXPR":            "0 6 * * *",
		"WISELABZ_DOC_EXPORT_ENABLED":              "true",
	}
	for k, v := range env {
		t.Setenv(k, v)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	want := Config{
		DB:         Database{Driver: "postgres", DSN: "postgres://x", MaxOpenConns: 7, MaxIdleConns: 3, ConnMaxLifetimeSeconds: 60, ConnMaxIdleTimeSeconds: 30},
		Server:     Server{Host: "127.0.0.1", Port: 9090, Origin: "https://example.com", TrustedProxies: "10.0.0.0/8", Embed: true, ReadTimeoutSeconds: 5, WriteTimeoutSeconds: 6, ShutdownTimeoutSeconds: 7},
		Encryption: EncryptionSettings{Key: "env-key"},
		Auth:       AuthSettings{Secret: "env-secret", AccessTokenTTL: 60, RefreshTokenTTL: 120, StepUpForDestructive: false},
		AI: AISettings{
			Enabled: true, Provider: "openai", Model: "gpt-x", APIKey: "key", BaseURL: "http://localhost", Mode: "auto_update",
			EmbedProvider: "openai", EmbedModel: "text-embedding-3-small", EmbedAPIKey: "embed-key", EmbedBaseURL: "http://embed-host",
		},
		Sync:      SyncSettings{Schedule: "* * * * *", PollCronExpr: "*/5 * * * * *"},
		Quality:   QualitySettings{CronExpr: "0 1 * * *"},
		Rotation:  RotationSettings{MaxAgeDays: 45, WarnDays: 7},
		Log:       LogSettings{Level: "debug", Format: "json"},
		Retention: RetentionSettings{SnapshotDays: 1, DocVersionDays: 2, AlertDays: 3, SyncRunDays: 4, AuditDays: 5, HealthCheckDays: 6, CronExpr: "0 4 * * *"},
		Backup:    BackupSettings{Dir: "/tmp/backups", CronExpr: "0 5 * * *", MaxBackups: 1, MaxAgeHours: 2, Enabled: false},
		DocExport: DocExportSettings{Dir: "/tmp/docexport", CronExpr: "0 6 * * *", Enabled: true},
	}

	got := *cfg
	got.Auth.OIDC = nil
	want.Auth.OIDC = nil
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() with all env vars set =\n%+v\nwant\n%+v", got, want)
	}
}

func TestLoadRejectsInvalidCronExpr(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)          //nolint:errcheck
	defer os.Chdir(oldDir) //nolint:errcheck

	t.Setenv("WISELABZ_AUTH_SECRET", "env-secret")
	t.Setenv("WISELABZ_RETENTION_CRON_EXPR", "not a cron expression")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error for invalid retention.cron_expr")
	}
}

func TestLoadRejectsEmptyCronExpr(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)          //nolint:errcheck
	defer os.Chdir(oldDir) //nolint:errcheck

	t.Setenv("WISELABZ_AUTH_SECRET", "env-secret")
	t.Setenv("WISELABZ_RETENTION_CRON_EXPR", "")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want error for empty retention.cron_expr")
	}
}

func TestServerAddr(t *testing.T) {
	s := Server{Host: "0.0.0.0", Port: 8080}
	if addr := s.Addr(); addr != "0.0.0.0:8080" {
		t.Errorf("Addr() = %q, want 0.0.0.0:8080", addr)
	}
}

func TestReadTimeoutDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected time.Duration
	}{
		{"positive value", 15, 15 * time.Second},
		{"zero value defaults to 10s", 0, 10 * time.Second},
		{"negative value defaults to 10s", -1, 10 * time.Second},
		{"large positive value", 300, 300 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{ReadTimeoutSeconds: tt.seconds}
			if got := s.ReadTimeoutDuration(); got != tt.expected {
				t.Errorf("ReadTimeoutDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWriteTimeoutDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected time.Duration
	}{
		{"positive value", 45, 45 * time.Second},
		{"zero value defaults to 30s", 0, 30 * time.Second},
		{"negative value defaults to 30s", -5, 30 * time.Second},
		{"large positive value", 600, 600 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{WriteTimeoutSeconds: tt.seconds}
			if got := s.WriteTimeoutDuration(); got != tt.expected {
				t.Errorf("WriteTimeoutDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShutdownTimeoutDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected time.Duration
	}{
		{"positive value", 20, 20 * time.Second},
		{"zero value defaults to 10s", 0, 10 * time.Second},
		{"negative value defaults to 10s", -10, 10 * time.Second},
		{"large positive value", 120, 120 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{ShutdownTimeoutSeconds: tt.seconds}
			if got := s.ShutdownTimeoutDuration(); got != tt.expected {
				t.Errorf("ShutdownTimeoutDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccessTokenTTLDuration(t *testing.T) {
	tests := []struct {
		name     string
		ttl      int
		expected time.Duration
	}{
		{"positive value", 1800, 1800 * time.Second},
		{"zero value defaults to 15min", 0, 15 * time.Minute},
		{"negative value defaults to 15min", -1, 15 * time.Minute},
		{"config default 900 seconds (15min)", 900, 900 * time.Second},
		{"custom large value", 3600, 3600 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AuthSettings{AccessTokenTTL: tt.ttl}
			if got := a.AccessTokenTTLDuration(); got != tt.expected {
				t.Errorf("AccessTokenTTLDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRefreshTokenTTLDuration(t *testing.T) {
	tests := []struct {
		name     string
		ttl      int
		expected time.Duration
	}{
		{"positive value", 172800, 172800 * time.Second},
		{"zero value defaults to 7 days", 0, 7 * 24 * time.Hour},
		{"negative value defaults to 7 days", -100, 7 * 24 * time.Hour},
		{"config default 604800 seconds (7 days)", 604800, 604800 * time.Second},
		{"custom value 14 days in seconds", 1209600, 1209600 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AuthSettings{RefreshTokenTTL: tt.ttl}
			if got := a.RefreshTokenTTLDuration(); got != tt.expected {
				t.Errorf("RefreshTokenTTLDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}
