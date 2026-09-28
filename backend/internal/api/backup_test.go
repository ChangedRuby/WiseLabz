package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/system"
	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestBackupExportRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/system/backup/export", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupExportRedactsSecrets(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	cfg, err := store.MarshalConnectorConfig("proxmox", map[string]any{
		"url": "https://pve.example.com", "token_id": "root@pam!monitoring", "token_secret": "super-secret-token",
	}, app.Config.Encryption.Key)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	if err := app.Store.CreateConnector(context.Background(), &store.ConnectorRecord{
		Name: "pve1", Category: "virtualization", Type: "proxmox", URL: "https://pve.example.com", ConfigData: cfg,
	}); err != nil {
		t.Fatalf("create connector: %v", err)
	}

	rec := app.req(t, http.MethodGet, "/api/system/backup/export", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	if strings.Contains(body, "token_secret") || strings.Contains(body, "super-secret-token") {
		t.Errorf("export body contains connector secret: %s", body)
	}
	if strings.Contains(body, "api_secret") {
		t.Errorf("export body contains opnsense secret field name: %s", body)
	}
}

func TestBackupImportBadVersion(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/system/backup/import", map[string]any{"version": 9999}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupImportMalformedJSON(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", strings.NewReader("{not valid json"))
	req.Header.Set("Authorization", "Bearer "+opToken)
	rec := app.serve(req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupImportRollbackNotVisibleViaAPI(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	// Two doc versions with the same ID for the same doc: the second insert
	// fails partway through the import, and the whole transaction — including
	// the doc itself — must roll back rather than leaving a partial import.
	b := map[string]any{
		"version": backup.BundleVersion,
		"docs":    []map[string]any{{"id": "doc-rollback", "title": "Rollback test", "kind": "lab"}},
		"docVersions": []map[string]any{
			{"id": "version-one", "docId": "doc-rollback", "rev": 1, "content": "first", "trigger": "manual"},
			{"id": "version-two", "docId": "doc-rollback", "rev": 1, "content": "duplicate", "trigger": "manual"},
		},
	}

	rec := app.req(t, http.MethodPost, "/api/system/backup/import", b, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
	if _, err := app.Store.GetDoc(context.Background(), "doc-rollback"); err == nil {
		t.Fatal("doc from the failed import is visible via the store; transaction did not roll back")
	}
}

func TestBackupScheduleRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodPut, "/api/system/backup/schedule", map[string]any{
		"cronExpr": "0 3 * * *", "maxBackups": 14, "maxAgeHours": 720, "enabled": true,
	}, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupUpdateScheduleRejectsInvalidCron(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/system/backup/schedule", map[string]any{
		"cronExpr": "not a cron expression", "maxBackups": 14, "maxAgeHours": 720, "enabled": true,
	}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupUpdateScheduleRejectsEmptyCron(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/system/backup/schedule", map[string]any{
		"cronExpr": "", "maxBackups": 14, "maxAgeHours": 720, "enabled": true,
	}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupUpdateScheduleAcceptsValidCronAndPersists(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/system/backup/schedule", map[string]any{
		"cronExpr": "0 */6 * * *", "maxBackups": 7, "maxAgeHours": 168, "enabled": false,
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	getRec := app.req(t, http.MethodGet, "/api/system/backup/schedule", nil, opToken)
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", getRec.Code, getRec.Body)
	}
	var sched struct {
		CronExpr   string `json:"cronExpr"`
		MaxBackups int    `json:"maxBackups"`
		Enabled    bool   `json:"enabled"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &sched); err != nil {
		t.Fatalf("decode schedule: %v", err)
	}
	if sched.CronExpr != "0 */6 * * *" || sched.MaxBackups != 7 || sched.Enabled {
		t.Fatalf("schedule not persisted correctly: %+v", sched)
	}
}

func TestBackupImportRejectsOversizedBody(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	body := bytes.Repeat([]byte(" "), system.MaxImportBytes+1)
	req := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+opToken)
	rec := app.serve(req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413; body = %s", rec.Code, rec.Body)
	}
}

func TestBackupExportImportRoundTrip(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	if err := app.Store.CreateConnector(context.Background(), &store.ConnectorRecord{
		Name: "pve1", Category: "virtualization", Type: "proxmox", URL: "https://pve.example.com",
	}); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	doc := &store.DocRecord{Title: "Doc 1", Kind: "lab", Content: "hello"}
	if err := app.Store.CreateDoc(context.Background(), doc); err != nil {
		t.Fatalf("create doc: %v", err)
	}
	if err := app.Store.CreateDocVersion(context.Background(), &store.DocVersionRecord{
		DocID: doc.ID, Rev: 1, Content: "hello", Trigger: "manual",
	}); err != nil {
		t.Fatalf("create doc version: %v", err)
	}
	tmpl := &store.TemplateRecord{Name: "Template 1"}
	if err := app.Store.CreateTemplate(context.Background(), tmpl); err != nil {
		t.Fatalf("create template: %v", err)
	}
	if err := app.Store.CreateTemplateSection(context.Background(), &store.TemplateSectionRecord{
		TemplateID: tmpl.ID, Title: "Section 1", Ord: 0, Body: "body",
	}); err != nil {
		t.Fatalf("create template section: %v", err)
	}

	exportRec := app.req(t, http.MethodGet, "/api/system/backup/export", nil, opToken)
	if exportRec.Code != http.StatusOK {
		t.Fatalf("export status = %d, want 200; body = %s", exportRec.Code, exportRec.Body)
	}
	var bundle backup.Bundle
	if err := json.Unmarshal(exportRec.Body.Bytes(), &bundle); err != nil {
		t.Fatalf("decode export body: %v", err)
	}

	// Import into a fresh second app/store so nothing already exists.
	app2 := newTestApp(t)
	_, opToken2 := app2.user(t, "operator")

	importRec := app2.req(t, http.MethodPost, "/api/system/backup/import", bundle, opToken2)
	if importRec.Code != http.StatusOK {
		t.Fatalf("import status = %d, want 200; body = %s", importRec.Code, importRec.Body)
	}
	var result backup.Result
	if err := json.Unmarshal(importRec.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode import body: %v", err)
	}
	if result.Connectors.Imported != 1 {
		t.Errorf("Connectors.Imported = %d, want 1", result.Connectors.Imported)
	}
	if result.Docs.Imported != 1 {
		t.Errorf("Docs.Imported = %d, want 1", result.Docs.Imported)
	}
	if result.DocVersions.Imported != 1 {
		t.Errorf("DocVersions.Imported = %d, want 1", result.DocVersions.Imported)
	}
	if result.Templates.Imported != 1 {
		t.Errorf("Templates.Imported = %d, want 1", result.Templates.Imported)
	}
	if result.TemplateSections.Imported != 1 {
		t.Errorf("TemplateSections.Imported = %d, want 1", result.TemplateSections.Imported)
	}
}
