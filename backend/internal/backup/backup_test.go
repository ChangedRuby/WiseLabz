package backup_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/store"

	// Register connector implementations (proxmox, opnsense, ...) so
	// GetTypeSchema resolves their secret fields the same way it does in
	// production, wired via internal/api/router.go's blank import.
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

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

// testEncKey is a fixed valid base64-encoded 32-byte AES-256 key for tests
// that need to round-trip connector config through MarshalConnectorConfig
// (mirrors backend/internal/api/testapp_test.go's test key).
const testEncKey = "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="

func TestExportRedactsConnectorSecrets(t *testing.T) {
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

	opnsenseCfg, err := store.MarshalConnectorConfig("opnsense", map[string]any{
		"url":        "https://fw.example.com",
		"api_key":    "some-key",
		"api_secret": "super-secret-api-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("marshal opnsense config: %v", err)
	}
	if err := s.CreateConnector(ctx, &store.ConnectorRecord{
		Name: "fw1", Category: "networking", Type: "opnsense", URL: "https://fw.example.com",
		ConfigData: opnsenseCfg,
	}); err != nil {
		t.Fatalf("create opnsense connector: %v", err)
	}

	b, err := backup.Export(ctx, s)
	if err != nil {
		t.Fatalf("Export() error: %v", err)
	}

	raw, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal bundle: %v", err)
	}
	rawStr := string(raw)

	if strings.Contains(rawStr, "super-secret-token") {
		t.Error("exported bundle contains the proxmox token_secret value")
	}
	if strings.Contains(rawStr, "super-secret-api-secret") {
		t.Error("exported bundle contains the opnsense api_secret value")
	}
	if strings.Contains(rawStr, "apiKey") || strings.Contains(rawStr, "api_key_encrypted") {
		t.Error("exported bundle contains an AI API key field")
	}

	for _, c := range b.Connectors {
		cfg, err := store.ParseConnectorConfig(c.Type, c.ConfigData, testEncKey)
		if err != nil {
			t.Fatalf("parse connector config: %v", err)
		}
		if c.Type == "proxmox" {
			if _, ok := cfg["token_secret"]; ok {
				t.Error("proxmox connector still has token_secret in exported config")
			}
			if _, ok := cfg["token_id"]; !ok {
				t.Error("proxmox connector lost non-secret field token_id")
			}
		}
		if c.Type == "opnsense" {
			if _, ok := cfg["api_secret"]; ok {
				t.Error("opnsense connector still has api_secret in exported config")
			}
			if _, ok := cfg["api_key"]; ok {
				t.Error("opnsense connector still has api_key in exported config")
			}
			if _, ok := cfg["url"]; !ok {
				t.Error("opnsense connector lost non-secret field url")
			}
		}
	}
}

func TestValidateBundleRejectsWrongVersion(t *testing.T) {
	b := &backup.Bundle{Version: backup.BundleVersion + 1}
	err := backup.ValidateBundle(b)
	if err == nil {
		t.Fatal("expected error for mismatched version, got nil")
	}
}

func TestValidateBundleRejectsOrphanDocVersion(t *testing.T) {
	b := &backup.Bundle{
		Version: backup.BundleVersion,
		Docs:    []store.DocRecord{{ID: "doc-1"}},
		DocVersions: []store.DocVersionRecord{
			{ID: "ver-1", DocID: "doc-does-not-exist", Rev: 1},
		},
	}
	err := backup.ValidateBundle(b)
	if err == nil {
		t.Fatal("expected error for doc version referencing unknown doc, got nil")
	}
}

func TestImportIsIdempotent(t *testing.T) {
	ctx := context.Background()
	src := newTestStore(t)

	if err := src.CreateConnector(ctx, &store.ConnectorRecord{
		Name: "pve1", Category: "virtualization", Type: "proxmox", URL: "https://pve.example.com",
	}); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	doc := &store.DocRecord{Title: "Doc 1", Kind: "lab", Content: "hello"}
	if err := src.CreateDoc(ctx, doc); err != nil {
		t.Fatalf("create doc: %v", err)
	}
	if err := src.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID: doc.ID, Rev: 1, Content: "hello", Trigger: "manual",
	}); err != nil {
		t.Fatalf("create doc version: %v", err)
	}
	tmpl := &store.TemplateRecord{Name: "Template 1"}
	if err := src.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("create template: %v", err)
	}
	if err := src.CreateTemplateSection(ctx, &store.TemplateSectionRecord{
		TemplateID: tmpl.ID, Title: "Section 1", Ord: 0, Body: "body",
	}); err != nil {
		t.Fatalf("create template section: %v", err)
	}

	b, err := backup.Export(ctx, src)
	if err != nil {
		t.Fatalf("Export() error: %v", err)
	}

	dst := newTestStore(t)

	first, err := backup.Import(ctx, dst, b)
	if err != nil {
		t.Fatalf("first Import() error: %v", err)
	}
	if first.Connectors.Imported != 1 || first.Docs.Imported != 1 ||
		first.DocVersions.Imported != 1 || first.Templates.Imported != 1 ||
		first.TemplateSections.Imported != 1 {
		t.Fatalf("expected everything imported on first run, got %+v", first)
	}

	second, err := backup.Import(ctx, dst, b)
	if err != nil {
		t.Fatalf("second Import() error: %v", err)
	}
	if second.Connectors.Imported != 0 || second.Docs.Imported != 0 ||
		second.DocVersions.Imported != 0 || second.Templates.Imported != 0 ||
		second.TemplateSections.Imported != 0 {
		t.Fatalf("expected nothing newly imported on second run, got %+v", second)
	}
	if second.Connectors.Skipped != 1 || second.Docs.Skipped != 1 ||
		second.DocVersions.Skipped != 1 || second.Templates.Skipped != 1 ||
		second.TemplateSections.Skipped != 1 {
		t.Fatalf("expected everything skipped on second run, got %+v", second)
	}
}

func TestImportRollsBackOnLateFailure(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	b := &backup.Bundle{
		Version: backup.BundleVersion,
		Docs:    []store.DocRecord{{ID: "doc-transaction", Title: "Transaction test", Kind: "lab"}},
		DocVersions: []store.DocVersionRecord{
			{ID: "version-one", DocID: "doc-transaction", Rev: 1, Content: "first", Trigger: "manual"},
			{ID: "version-two", DocID: "doc-transaction", Rev: 1, Content: "duplicate", Trigger: "manual"},
		},
	}

	if _, err := backup.Import(ctx, s, b); err == nil {
		t.Fatal("Import() error = nil, want late duplicate version failure")
	}
	if _, err := s.GetDoc(ctx, "doc-transaction"); err == nil {
		t.Fatal("imported doc remained after failed transaction")
	}
}

func TestExportIncludesRecordsBeyondAPage(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	for i := range 1001 {
		if err := s.CreateDoc(ctx, &store.DocRecord{
			ID:      fmt.Sprintf("doc-%04d", i),
			Title:   fmt.Sprintf("Document %04d", i),
			Kind:    "lab",
			Content: "content",
		}); err != nil {
			t.Fatalf("create doc %d: %v", i, err)
		}
	}

	b, err := backup.Export(ctx, s)
	if err != nil {
		t.Fatalf("Export() error: %v", err)
	}
	if got := len(b.Docs); got != 1001 {
		t.Fatalf("exported docs = %d, want 1001", got)
	}
}

func TestExportToFile(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	// Create a test backup
	tmpDir := t.TempDir()
	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	// Verify returned metadata
	if run.ID == "" {
		t.Error("Run.ID is empty")
	}
	if run.FilePath == "" {
		t.Error("Run.FilePath is empty")
	}
	if run.SizeBytes <= 0 {
		t.Error("Run.SizeBytes should be positive")
	}
	if run.CreatedAt == "" {
		t.Error("Run.CreatedAt is empty")
	}

	// Verify file exists and is valid JSON
	if _, err := os.Stat(run.FilePath); err != nil {
		t.Fatalf("backup file not created: %v", err)
	}

	data, err := os.ReadFile(run.FilePath)
	if err != nil {
		t.Fatalf("read backup file: %v", err)
	}

	var b backup.Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("unmarshal backup JSON: %v", err)
	}

	if b.Version != backup.BundleVersion {
		t.Errorf("Bundle version mismatch: got %d, want %d", b.Version, backup.BundleVersion)
	}
}

// TestExportToFileDirNotWritable verifies ExportToFile returns an error
// (rather than silently succeeding or panicking) when the backup directory
// exists but the process can't write to it. Skipped when running as root,
// since root ignores directory write permissions.
func TestExportToFileDirNotWritable(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits don't block writes")
	}

	ctx := context.Background()
	s := newTestStore(t)

	tmpDir := t.TempDir()
	if err := os.Chmod(tmpDir, 0o500); err != nil {
		t.Fatalf("chmod backup dir read-only: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(tmpDir, 0o700) }) // let t.TempDir() clean up

	if _, err := backup.ExportToFile(ctx, s, tmpDir); err == nil {
		t.Fatal("ExportToFile into a read-only directory: got nil error, want an error")
	}
}

func TestExportToFileCreatesDirectory(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)

	// Use a path that doesn't exist yet
	tmpBase := t.TempDir()
	tmpDir := fmt.Sprintf("%s/subdir/backups", tmpBase)

	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(tmpDir); err != nil {
		t.Fatalf("backup directory not created: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(run.FilePath); err != nil {
		t.Fatalf("backup file not found: %v", err)
	}
}

// TestExportToFilePermissions is a regression test for GHSA-c753: the backup
// bundle is a full infrastructure inventory (secrets redacted, but still
// sensitive) and must not be readable by other local users.
func TestExportToFilePermissions(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission bits aren't enforced")
	}

	ctx := context.Background()
	s := newTestStore(t)

	tmpBase := t.TempDir()
	tmpDir := fmt.Sprintf("%s/backups", tmpBase)

	run, err := backup.ExportToFile(ctx, s, tmpDir)
	if err != nil {
		t.Fatalf("ExportToFile: %v", err)
	}

	dirInfo, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("stat backup dir: %v", err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0o700 {
		t.Errorf("backup directory permissions = %o, want 0700", perm)
	}

	fileInfo, err := os.Stat(run.FilePath)
	if err != nil {
		t.Fatalf("stat backup file: %v", err)
	}
	if perm := fileInfo.Mode().Perm(); perm != 0o600 {
		t.Errorf("backup file permissions = %o, want 0600", perm)
	}
}
