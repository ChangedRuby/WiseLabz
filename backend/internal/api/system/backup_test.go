package system

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"

	// Register connector implementations so the bundle export uses their
	// TypeSchemas to correctly identify and redact secret fields.
	_ "github.com/WiseLabz/wiselabz/internal/connector/all"
)

const testEncKey = "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="

// TestExportBackupRedactsConnectorSecrets verifies the ExportBackup handler's
// exported bundle never carries a connector's secret config fields (API keys,
// tokens, passwords, etc.), matching the redaction pattern used in
// internal/backup/backup.go. This mirrors the test in PR #221 diagnostics.
func TestExportBackupRedactsConnectorSecrets(t *testing.T) {
	s := apitest.NewStore(t)
	ctx := context.Background()

	// Create a Proxmox connector with secrets
	proxmoxCfg, err := store.MarshalConnectorConfig("proxmox", map[string]any{
		"url":          "https://pve.example.com:8006",
		"token_id":     "root@pam!monitoring",
		"token_secret": "super-secret-token",
	}, testEncKey)
	if err != nil {
		t.Fatalf("marshal proxmox config: %v", err)
	}
	if err := s.CreateConnector(ctx, &store.ConnectorRecord{
		Name:       "pve1",
		Category:   "virtualization",
		Type:       "proxmox",
		URL:        "https://pve.example.com:8006",
		ConfigData: proxmoxCfg,
	}); err != nil {
		t.Fatalf("create proxmox connector: %v", err)
	}

	// Create an OPNsense connector with secrets
	opnsenseCfg, err := store.MarshalConnectorConfig("opnsense", map[string]any{
		"url":        "https://fw.example.com",
		"api_key":    "some-key",
		"api_secret": "super-secret-api-secret",
	}, testEncKey)
	if err != nil {
		t.Fatalf("marshal opnsense config: %v", err)
	}
	if err := s.CreateConnector(ctx, &store.ConnectorRecord{
		Name:       "fw1",
		Category:   "networking",
		Type:       "opnsense",
		URL:        "https://fw.example.com",
		ConfigData: opnsenseCfg,
	}); err != nil {
		t.Fatalf("create opnsense connector: %v", err)
	}

	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)
	req := httptest.NewRequest(http.MethodGet, "/api/system/backup/export", nil)
	rr := httptest.NewRecorder()
	h.ExportBackup(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("ExportBackup() status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	// Verify secrets are not in the exported JSON
	exportedBody := rr.Body.String()
	if strings.Contains(exportedBody, "super-secret-token") {
		t.Error("exported backup contains the proxmox token_secret value")
	}
	if strings.Contains(exportedBody, "super-secret-api-secret") {
		t.Error("exported backup contains the opnsense api_secret value")
	}

	// Parse the exported bundle to verify structure
	var b backup.Bundle
	if err := json.Unmarshal(rr.Body.Bytes(), &b); err != nil {
		t.Fatalf("unmarshal exported bundle: %v", err)
	}

	if b.Version != backup.BundleVersion {
		t.Errorf("bundle version = %d, want %d", b.Version, backup.BundleVersion)
	}

	// Verify each connector's config is redacted
	for _, c := range b.Connectors {
		cfg, err := store.ParseConnectorConfig(c.Type, c.ConfigData, testEncKey)
		if err != nil {
			t.Fatalf("parse connector %s config: %v", c.Type, err)
		}

		if c.Type == "proxmox" {
			if _, ok := cfg["token_secret"]; ok {
				t.Error("proxmox connector ConfigData still contains token_secret in export")
			}
			// Non-secret fields should be preserved
			if _, ok := cfg["token_id"]; !ok {
				t.Error("proxmox connector lost non-secret field token_id in export")
			}
		}

		if c.Type == "opnsense" {
			if _, ok := cfg["api_secret"]; ok {
				t.Error("opnsense connector ConfigData still contains api_secret in export")
			}
			if _, ok := cfg["api_key"]; ok {
				t.Error("opnsense connector ConfigData still contains api_key in export")
			}
			// Non-secret fields should be preserved
			if _, ok := cfg["url"]; !ok {
				t.Error("opnsense connector lost non-secret field url in export")
			}
		}
	}
}

// TestImportBackupCorruptJSON verifies ImportBackup rejects invalid JSON gracefully.
func TestImportBackupCorruptJSON(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", strings.NewReader(`{`))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rr.Body.String(), "Invalid JSON") {
		t.Errorf("expected 'Invalid JSON' error, got: %s", rr.Body.String())
	}
}

// TestImportBackupTruncatedJSON verifies ImportBackup handles truncated data.
func TestImportBackupTruncatedJSON(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	// Valid JSON structure but incomplete (missing closing braces)
	truncated := `{"version":1,"exportedAt":"2024-01-01","connectors":[{"id":"c1","name":"test","type":"proxmox","url":"https://example.com","category":"virtualization","configData":"{}`
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", strings.NewReader(truncated))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// TestImportBackupRequestTooLarge verifies ImportBackup enforces the 10 MiB limit.
func TestImportBackupRequestTooLarge(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	// Create a valid JSON structure with a large array of connectors to exceed 10 MiB.
	// Each connector record is roughly 200 bytes, so 50k connectors = ~10 MiB+.
	conns := make([]map[string]any, 0)
	for i := 0; i < 51000; i++ {
		conns = append(conns, map[string]any{
			"id":         "c" + string(rune(i)),
			"name":       "connector-" + string(rune(i)),
			"type":       "proxmox",
			"category":   "virtualization",
			"url":        "https://example.com",
			"configData": strings.Repeat("x", 100),
		})
	}
	b := map[string]any{
		"version":    1,
		"exportedAt": "2024-01-01T00:00:00Z",
		"connectors": conns,
	}
	payload, _ := json.Marshal(b)
	if int64(len(payload)) <= MaxImportBytes {
		t.Skipf("test payload only %d bytes, need > %d to trigger limit", len(payload), MaxImportBytes)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(payload))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusRequestEntityTooLarge, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "10 MiB") {
		t.Errorf("expected '10 MiB' error message, got: %s", rr.Body.String())
	}
}

// TestImportBackupWrongVersion verifies ImportBackup rejects incompatible bundle versions.
func TestImportBackupWrongVersion(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	b := backup.Bundle{
		Version:    backup.BundleVersion + 1,
		ExportedAt: "2024-01-01T00:00:00Z",
	}
	body, _ := json.Marshal(b)
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rr.Body.String(), "unsupported backup version") {
		t.Errorf("expected version error, got: %s", rr.Body.String())
	}
}

// TestImportBackupInvalidCategory verifies ImportBackup rejects connectors with invalid categories.
func TestImportBackupInvalidCategory(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	b := backup.Bundle{
		Version:    backup.BundleVersion,
		ExportedAt: "2024-01-01T00:00:00Z",
		Connectors: []store.ConnectorRecord{
			{
				ID:       "conn-1",
				Name:     "bad-connector",
				Type:     "proxmox",
				Category: "invalid_category",
				URL:      "https://example.com",
			},
		},
	}
	body, _ := json.Marshal(b)
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusBadRequest, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "invalid category") {
		t.Errorf("expected invalid category error, got: %s", rr.Body.String())
	}
}

// TestImportBackupOrphanDocVersion verifies ImportBackup rejects doc versions
// that reference non-existent docs, preventing partial-state application.
func TestImportBackupOrphanDocVersion(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	b := backup.Bundle{
		Version:    backup.BundleVersion,
		ExportedAt: "2024-01-01T00:00:00Z",
		Docs:       []store.DocRecord{{ID: "doc-1", Title: "Doc 1", Kind: "lab", Content: ""}},
		DocVersions: []store.DocVersionRecord{
			{ID: "ver-1", DocID: "doc-does-not-exist", Rev: 1, Content: "orphan", Trigger: "manual"},
		},
	}
	body, _ := json.Marshal(b)
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rr.Body.String(), "unknown doc") {
		t.Errorf("expected orphan doc version error, got: %s", rr.Body.String())
	}
}

// TestImportBackupOrphanTemplateSection verifies ImportBackup rejects template
// sections that reference non-existent templates, preventing partial-state application.
func TestImportBackupOrphanTemplateSection(t *testing.T) {
	s := apitest.NewStore(t)
	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	b := backup.Bundle{
		Version:    backup.BundleVersion,
		ExportedAt: "2024-01-01T00:00:00Z",
		Templates:  []store.TemplateRecord{{ID: "tmpl-1", Name: "Template 1"}},
		TemplateSections: []store.TemplateSectionRecord{
			{ID: "sec-1", TemplateID: "tmpl-does-not-exist", Title: "Section 1", Ord: 0, Body: "orphan"},
		},
	}
	body, _ := json.Marshal(b)
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rr.Body.String(), "unknown template") {
		t.Errorf("expected orphan template section error, got: %s", rr.Body.String())
	}
}

// TestImportBackupPartialFailureRollback verifies that if import fails mid-way
// (e.g., duplicate doc version), the entire transaction is rolled back so no
// partial state is applied.
func TestImportBackupPartialFailureRollback(t *testing.T) {
	s := apitest.NewStore(t)
	ctx := context.Background()

	// Create initial doc and version
	doc := &store.DocRecord{ID: "doc-1", Title: "Doc 1", Kind: "lab", Content: "v1"}
	if err := s.CreateDoc(ctx, doc); err != nil {
		t.Fatalf("create initial doc: %v", err)
	}
	if err := s.CreateDocVersion(ctx, &store.DocVersionRecord{
		ID: "ver-1-orig", DocID: doc.ID, Rev: 1, Content: "v1", Trigger: "manual",
	}); err != nil {
		t.Fatalf("create initial version: %v", err)
	}

	h := NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir(), nil)

	// Import bundle with:
	// - A new connector (should be imported)
	// - A new doc (should be imported)
	// - Two versions with the same ID (second will fail)
	b := backup.Bundle{
		Version:    backup.BundleVersion,
		ExportedAt: "2024-01-01T00:00:00Z",
		Connectors: []store.ConnectorRecord{
			{ID: "new-conn", Name: "new", Type: "proxmox", Category: "virtualization", URL: "https://example.com"},
		},
		Docs: []store.DocRecord{
			{ID: "doc-import", Title: "New Doc", Kind: "lab", Content: "import"},
		},
		DocVersions: []store.DocVersionRecord{
			{ID: "ver-new-1", DocID: "doc-import", Rev: 1, Content: "v1", Trigger: "manual"},
			{ID: "ver-new-1", DocID: "doc-import", Rev: 1, Content: "v2", Trigger: "manual"}, // duplicate ID
		},
	}
	body, _ := json.Marshal(b)
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	h.ImportBackup(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d (expected failure)", rr.Code, http.StatusBadRequest)
	}

	// Verify new connector was NOT persisted (rollback worked)
	_, err := s.GetConnector(ctx, "new-conn")
	if err == nil {
		t.Error("new connector was persisted despite failed import (rollback didn't work)")
	}

	// Verify new doc was NOT persisted
	_, err = s.GetDoc(ctx, "doc-import")
	if err == nil {
		t.Error("new doc was persisted despite failed import (rollback didn't work)")
	}

	// Verify initial data is still intact
	originalDoc, err := s.GetDoc(ctx, doc.ID)
	if err != nil {
		t.Fatalf("original doc lost: %v", err)
	}
	if originalDoc.Title != "Doc 1" {
		t.Errorf("original doc corrupted: title = %q, want 'Doc 1'", originalDoc.Title)
	}
}

// TestImportBackupIdempotence verifies that importing the same bundle twice
// results in the same state: the second import skips all records without error.
func TestImportBackupIdempotence(t *testing.T) {
	ctx := context.Background()

	// Create source store and export
	src := apitest.NewStore(t)
	if err := src.CreateConnector(ctx, &store.ConnectorRecord{
		ID:       "conn-1",
		Name:     "connector1",
		Type:     "proxmox",
		Category: "virtualization",
		URL:      "https://example.com",
	}); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	b, err := backup.Export(ctx, src)
	if err != nil {
		t.Fatalf("Export() error: %v", err)
	}

	// Create destination store and import twice
	dst := apitest.NewStore(t)
	h := NewHandler(dst.DB(), &config.Config{}, dst, nil, t.TempDir(), nil)

	body, _ := json.Marshal(b)

	// First import
	req1 := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr1 := httptest.NewRecorder()
	h.ImportBackup(rr1, req1)
	if rr1.Code != http.StatusOK {
		t.Fatalf("first import status = %d, want %d; body=%s", rr1.Code, http.StatusOK, rr1.Body.String())
	}

	var result1 backup.Result
	if err := json.Unmarshal(rr1.Body.Bytes(), &result1); err != nil {
		t.Fatalf("unmarshal first import result: %v", err)
	}
	if result1.Connectors.Imported != 1 {
		t.Fatalf("first import: imported = %d, want 1", result1.Connectors.Imported)
	}

	// Second import (should skip all)
	req2 := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	rr2 := httptest.NewRecorder()
	h.ImportBackup(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("second import status = %d, want %d; body=%s", rr2.Code, http.StatusOK, rr2.Body.String())
	}

	var result2 backup.Result
	if err := json.Unmarshal(rr2.Body.Bytes(), &result2); err != nil {
		t.Fatalf("unmarshal second import result: %v", err)
	}
	if result2.Connectors.Imported != 0 {
		t.Errorf("second import: imported = %d, want 0", result2.Connectors.Imported)
	}
	if result2.Connectors.Skipped != 1 {
		t.Errorf("second import: skipped = %d, want 1", result2.Connectors.Skipped)
	}
}

// TestExportImportRoundTrip verifies the full export/import cycle preserves
// non-secret data and handles all entity types correctly.
func TestExportImportRoundTrip(t *testing.T) {
	ctx := context.Background()
	src := apitest.NewStore(t)

	// Create a complete test dataset
	conn := &store.ConnectorRecord{
		Name:     "test-conn",
		Type:     "proxmox",
		Category: "virtualization",
		URL:      "https://pve.example.com",
	}
	if err := src.CreateConnector(ctx, conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	doc := &store.DocRecord{Title: "Test Doc", Kind: "lab", Content: "test content"}
	if err := src.CreateDoc(ctx, doc); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	if err := src.CreateDocVersion(ctx, &store.DocVersionRecord{
		DocID: doc.ID, Rev: 1, Content: "v1", Trigger: "manual",
	}); err != nil {
		t.Fatalf("create doc version: %v", err)
	}

	tmpl := &store.TemplateRecord{Name: "Test Template"}
	if err := src.CreateTemplate(ctx, tmpl); err != nil {
		t.Fatalf("create template: %v", err)
	}

	if err := src.CreateTemplateSection(ctx, &store.TemplateSectionRecord{
		TemplateID: tmpl.ID, Title: "Test Section", Ord: 0, Body: "section content",
	}); err != nil {
		t.Fatalf("create template section: %v", err)
	}

	// Export
	srcHandler := NewHandler(src.DB(), &config.Config{}, src, nil, t.TempDir(), nil)
	exportReq := httptest.NewRequest(http.MethodGet, "/api/system/backup/export", nil)
	exportRR := httptest.NewRecorder()
	srcHandler.ExportBackup(exportRR, exportReq)

	if exportRR.Code != http.StatusOK {
		t.Fatalf("export status = %d, want %d", exportRR.Code, http.StatusOK)
	}

	// Import into destination
	dst := apitest.NewStore(t)
	dstHandler := NewHandler(dst.DB(), &config.Config{}, dst, nil, t.TempDir(), nil)
	importReq := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(exportRR.Body.Bytes()))
	importRR := httptest.NewRecorder()
	dstHandler.ImportBackup(importRR, importReq)

	if importRR.Code != http.StatusOK {
		t.Fatalf("import status = %d, want %d; body=%s", importRR.Code, http.StatusOK, importRR.Body.String())
	}

	var result backup.Result
	if err := json.Unmarshal(importRR.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal import result: %v", err)
	}

	// Verify counts
	if result.Connectors.Imported != 1 {
		t.Errorf("connectors imported = %d, want 1", result.Connectors.Imported)
	}
	if result.Docs.Imported != 1 {
		t.Errorf("docs imported = %d, want 1", result.Docs.Imported)
	}
	if result.DocVersions.Imported != 1 {
		t.Errorf("doc versions imported = %d, want 1", result.DocVersions.Imported)
	}
	if result.Templates.Imported != 1 {
		t.Errorf("templates imported = %d, want 1", result.Templates.Imported)
	}
	if result.TemplateSections.Imported != 1 {
		t.Errorf("template sections imported = %d, want 1", result.TemplateSections.Imported)
	}

	// Verify data integrity
	importedDoc, err := dst.GetDoc(ctx, doc.ID)
	if err != nil {
		t.Fatalf("get imported doc: %v", err)
	}
	if importedDoc.Title != "Test Doc" {
		t.Errorf("imported doc title = %q, want 'Test Doc'", importedDoc.Title)
	}
	if importedDoc.Content != "test content" {
		t.Errorf("imported doc content = %q, want 'test content'", importedDoc.Content)
	}
}
