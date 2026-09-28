package api_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	mcpclient "github.com/mark3labs/mcp-go/client"
	mcptransport "github.com/mark3labs/mcp-go/client/transport"
	mcpsdk "github.com/mark3labs/mcp-go/mcp"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// TestMCPEndToEnd drives a real mcp-go client, over real HTTP, against the
// mounted /api/mcp route: HTTP -> AuthMiddleware -> API-key restriction
// context -> tool handler -> store. Unlike internal/mcp's unit tests (which
// call the tool handlers in-process, bypassing HTTP), this proves the full
// wiring in routes_mcp.go and router.go actually works end to end.
func TestMCPEndToEnd(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, _ := app.user(t, "viewer")
	connectorID := "svc-1"
	if err := app.Store.CreateConnector(context.Background(), &store.ConnectorRecord{
		ID: connectorID, Name: "svc-1", Category: "networking", Type: "test", Enabled: true,
	}); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	app.connectorGrant(t, userID, connectorID, "viewer")

	rawToken := mintAPIKey(t, app, userID, "full", nil)

	server := httptest.NewServer(app.Router)
	t.Cleanup(server.Close)

	client, err := mcpclient.NewStreamableHttpClient(server.URL+"/api/mcp",
		mcptransport.WithHTTPHeaders(map[string]string{"Authorization": "Bearer " + rawToken}),
	)
	if err != nil {
		t.Fatalf("new streamable http client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	ctx := context.Background()
	if _, err := client.Initialize(ctx, mcpsdk.InitializeRequest{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	tools, err := client.ListTools(ctx, mcpsdk.ListToolsRequest{})
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	wantTools := map[string]bool{
		"list_connectors": false, "search_docs": false, "list_findings": false,
		"list_changes": false, "list_attention_items": false,
	}
	for _, tool := range tools.Tools {
		if _, ok := wantTools[tool.Name]; ok {
			wantTools[tool.Name] = true
		}
	}
	for name, found := range wantTools {
		if !found {
			t.Errorf("tool %s not advertised by the mounted MCP server", name)
		}
	}

	res, err := client.CallTool(ctx, mcpsdk.CallToolRequest{
		Params: mcpsdk.CallToolParams{Name: "list_connectors", Arguments: map[string]any{}},
	})
	if err != nil {
		t.Fatalf("call list_connectors: %v", err)
	}
	if res.IsError {
		t.Fatalf("list_connectors returned an error result: %s", mcpsdk.GetTextFromContent(res.Content[0]))
	}
	var out struct {
		Connectors []struct {
			ID string `json:"id"`
		} `json:"connectors"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal([]byte(mcpsdk.GetTextFromContent(res.Content[0])), &out); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if out.Total != 1 || len(out.Connectors) != 1 || out.Connectors[0].ID != connectorID {
		t.Fatalf("got %+v, want the one granted connector %s", out, connectorID)
	}
}

// TestMCPConnectorRestrictedKey confirms a connector-restricted API key's
// list_connectors call narrows to its allow-list over the real HTTP path,
// same as a REST client hitting /api/connectors would see.
func TestMCPConnectorRestrictedKey(t *testing.T) {
	t.Parallel()
	app := newTestApp(t)
	userID, _ := app.user(t, "viewer")
	allowedID, deniedID := "svc-allowed", "svc-denied"
	for _, id := range []string{allowedID, deniedID} {
		if err := app.Store.CreateConnector(context.Background(), &store.ConnectorRecord{
			ID: id, Name: id, Category: "networking", Type: "test", Enabled: true,
		}); err != nil {
			t.Fatalf("create connector %s: %v", id, err)
		}
		app.connectorGrant(t, userID, id, "viewer")
	}

	rawToken := mintAPIKey(t, app, userID, "read", []string{allowedID})

	server := httptest.NewServer(app.Router)
	t.Cleanup(server.Close)

	client, err := mcpclient.NewStreamableHttpClient(server.URL+"/api/mcp",
		mcptransport.WithHTTPHeaders(map[string]string{"Authorization": "Bearer " + rawToken}),
	)
	if err != nil {
		t.Fatalf("new streamable http client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	ctx := context.Background()
	if _, err := client.Initialize(ctx, mcpsdk.InitializeRequest{}); err != nil {
		t.Fatalf("initialize: %v", err)
	}

	res, err := client.CallTool(ctx, mcpsdk.CallToolRequest{
		Params: mcpsdk.CallToolParams{Name: "list_connectors", Arguments: map[string]any{}},
	})
	if err != nil {
		t.Fatalf("call list_connectors: %v", err)
	}
	if res.IsError {
		t.Fatalf("list_connectors returned an error result: %s", mcpsdk.GetTextFromContent(res.Content[0]))
	}
	var out struct {
		Connectors []struct {
			ID string `json:"id"`
		} `json:"connectors"`
	}
	if err := json.Unmarshal([]byte(mcpsdk.GetTextFromContent(res.Content[0])), &out); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if len(out.Connectors) != 1 || out.Connectors[0].ID != allowedID {
		t.Fatalf("got %+v, want only the allow-listed connector %s", out.Connectors, allowedID)
	}
}

// mintAPIKey creates an API key for userID directly through the store (the
// same path POST /api/auth/api-keys uses) and returns its raw bearer token.
func mintAPIKey(t *testing.T, app *testApp, userID, scope string, connectorIDs []string) string {
	t.Helper()
	rawToken := "wlz_test_" + userID + "_" + scope
	key := &store.APIKey{
		UserID: userID, Name: "mcp-test-key", TokenHash: store.HashToken(rawToken),
		Role: "viewer", Scope: scope, ConnectorIDs: connectorIDs,
	}
	if scope == "" {
		key.Scope = auth.APIKeyScopeFull
	}
	if err := app.Store.CreateAPIKey(context.Background(), key); err != nil {
		t.Fatalf("create api key: %v", err)
	}
	return rawToken
}
