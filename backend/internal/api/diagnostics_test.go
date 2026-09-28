package api_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestDiagnosticsRoleBoundary(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/system/diagnostics", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestDiagnosticsRedactsSecrets(t *testing.T) {
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

	notifCfg := map[string]any{
		"channels": []map[string]any{{"type": "smtp", "password": "super-secret-smtp-password"}},
		"routing":  []map[string]any{},
	}
	if rec := app.req(t, http.MethodPut, "/api/notifications/config", notifCfg, opToken); rec.Code != http.StatusOK {
		t.Fatalf("update notifications config: status = %d, body = %s", rec.Code, rec.Body)
	}

	rec := app.req(t, http.MethodGet, "/api/system/diagnostics", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	body := rec.Body.String()
	if strings.Contains(body, "token_secret") || strings.Contains(body, "super-secret-token") {
		t.Errorf("diagnostics body contains connector secret: %s", body)
	}
	if strings.Contains(body, "super-secret-smtp-password") {
		t.Errorf("diagnostics body contains notification channel secret: %s", body)
	}
	if strings.Contains(body, "\"channels\"") {
		t.Errorf("diagnostics body includes notification channel config, which should be excluded entirely: %s", body)
	}
}

func TestDiagnosticsIncludesHealthVersionsAndFailures(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	if err := app.Store.CreateConnector(context.Background(), &store.ConnectorRecord{
		Name: "pve1", Category: "virtualization", Type: "proxmox", URL: "https://pve.example.com",
	}); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	connectors, err := app.Store.ListAllConnectors(context.Background())
	if err != nil || len(connectors) == 0 {
		t.Fatalf("list connectors: %v", err)
	}
	if err := app.Store.CreateSyncRun(context.Background(), &store.SyncRunRecord{
		ConnectorID: connectors[0].ID, Status: store.SyncRunStatusError, Error: "connection refused",
	}); err != nil {
		t.Fatalf("create sync run: %v", err)
	}

	rec := app.req(t, http.MethodGet, "/api/system/diagnostics", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	if got := rec.Header().Get("Content-Disposition"); !strings.Contains(got, "wiselabz-diagnostics-") {
		t.Errorf("Content-Disposition = %q, want attachment filename", got)
	}
	body := rec.Body.String()
	for _, want := range []string{`"health"`, `"versions"`, `"sanitizedConfig"`, `"recentFailures"`, "connection refused"} {
		if !strings.Contains(body, want) {
			t.Errorf("diagnostics body missing %q: %s", want, body)
		}
	}
}
