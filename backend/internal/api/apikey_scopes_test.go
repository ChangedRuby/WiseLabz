package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// createKey creates an API key through the HTTP API and returns the decoded
// response.
func createKey(t *testing.T, app *testApp, token string, body map[string]any) map[string]any {
	t.Helper()
	rec := app.req(t, http.MethodPost, "/api/auth/api-keys", body, token)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create API key status = %d: %s", rec.Code, rec.Body)
	}
	var created map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	return created
}

func newConnector(t *testing.T, app *testApp, name string) string {
	t.Helper()
	conn := &store.ConnectorRecord{Name: name, Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return conn.ID
}

func TestAPIKeyDefaultsToFullScope(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "operator")

	created := createKey(t, app, token, map[string]any{"name": "legacy"})
	if created["scope"] != "full" {
		t.Errorf("scope = %v, want full", created["scope"])
	}
	if ids, ok := created["connectorIds"].([]any); !ok || len(ids) != 0 {
		t.Errorf("connectorIds = %#v, want empty list", created["connectorIds"])
	}
	raw := created["token"].(string)
	if rec := app.req(t, http.MethodGet, "/api/users", nil, raw); rec.Code != http.StatusOK {
		t.Errorf("full key GET /api/users = %d, want 200: %s", rec.Code, rec.Body)
	}
}

func TestReadOnlyAPIKey(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	connID := newConnector(t, app, "pve")
	app.connectorGrant(t, userID, connID, "operator")

	raw := createKey(t, app, token, map[string]any{"name": "ro", "scope": "read"})["token"].(string)

	if rec := app.req(t, http.MethodGet, "/api/connectors/"+connID, nil, raw); rec.Code != http.StatusOK {
		t.Errorf("read key GET connector = %d, want 200: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodPost, "/api/connectors/"+connID+"/test", nil, raw); rec.Code != http.StatusForbidden {
		t.Errorf("read key POST connector test = %d, want 403: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodPost, "/api/auth/api-keys", map[string]any{"name": "x"}, raw); rec.Code != http.StatusForbidden {
		t.Errorf("read key minting a key = %d, want 403: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodPost, "/api/ws/ticket", nil, raw); rec.Code != http.StatusForbidden {
		t.Errorf("read key WS ticket = %d, want 403: %s", rec.Code, rec.Body)
	}
}

func TestReadOnlyAPIKeyCapsConnectorRoleAtViewer(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	connID := newConnector(t, app, "pve")
	app.connectorGrant(t, userID, connID, "operator")
	raw := createKey(t, app, token, map[string]any{"name": "ro", "scope": "read"})["token"].(string)

	rec := app.req(t, http.MethodGet, "/api/connectors/"+connID, nil, raw)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET connector = %d: %s", rec.Code, rec.Body)
	}
	var got map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["myRole"] != "viewer" {
		t.Errorf("connector role via read key = %v, want viewer", got["myRole"])
	}
}

func TestConnectorRestrictedAPIKey(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	allowed := newConnector(t, app, "allowed")
	other := newConnector(t, app, "other")
	app.connectorGrant(t, userID, allowed, "operator")
	app.connectorGrant(t, userID, other, "operator")

	created := createKey(t, app, token, map[string]any{"name": "one", "connectorIds": []string{allowed, allowed}})
	if ids, _ := created["connectorIds"].([]any); len(ids) != 1 || ids[0] != allowed {
		t.Fatalf("connectorIds = %#v, want deduped [%s]", created["connectorIds"], allowed)
	}
	raw := created["token"].(string)

	if rec := app.req(t, http.MethodGet, "/api/connectors/"+allowed, nil, raw); rec.Code != http.StatusOK {
		t.Errorf("GET allowed connector = %d, want 200: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodGet, "/api/connectors/"+other, nil, raw); rec.Code != http.StatusNotFound {
		t.Errorf("GET other connector = %d, want 404: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodGet, "/api/connectors/"+other+"/syncs", nil, raw); rec.Code != http.StatusForbidden {
		t.Errorf("GET other connector syncs = %d, want 403: %s", rec.Code, rec.Body)
	}

	list := app.req(t, http.MethodGet, "/api/connectors", nil, raw)
	if list.Code != http.StatusOK {
		t.Fatalf("list connectors = %d: %s", list.Code, list.Body)
	}
	var items []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(list.Body).Decode(&items); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(items) != 1 || items[0].ID != allowed || list.Header().Get("X-Total-Count") != "1" {
		t.Errorf("list via restricted key = %+v (total %s), want only %s", items, list.Header().Get("X-Total-Count"), allowed)
	}

	// Admin endpoints aren't connector-scoped, so a restricted key loses
	// instance admin even though its owner is one.
	if rec := app.req(t, http.MethodGet, "/api/users", nil, raw); rec.Code != http.StatusForbidden {
		t.Errorf("restricted key GET /api/users = %d, want 403: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodPost, "/api/auth/api-keys", map[string]any{"name": "wider"}, raw); rec.Code != http.StatusForbidden {
		t.Errorf("restricted key minting a key = %d, want 403: %s", rec.Code, rec.Body)
	}
	if rec := app.req(t, http.MethodPost, "/api/ws/ticket", nil, raw); rec.Code != http.StatusForbidden {
		t.Errorf("restricted key WS ticket = %d, want 403: %s", rec.Code, rec.Body)
	}
}

func TestAPIKeyCreateValidation(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	granted := newConnector(t, app, "granted")
	ungranted := newConnector(t, app, "ungranted")
	app.connectorGrant(t, userID, granted, "viewer")

	for name, body := range map[string]map[string]any{
		"unknown scope":       {"name": "k", "scope": "admin"},
		"ungranted connector": {"name": "k", "connectorIds": []string{granted, ungranted}},
		"missing connector":   {"name": "k", "connectorIds": []string{"does-not-exist"}},
	} {
		t.Run(name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/auth/api-keys", body, token)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400: %s", rec.Code, rec.Body)
			}
		})
	}
}
