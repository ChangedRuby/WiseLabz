package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	mcpclient "github.com/mark3labs/mcp-go/client"
	mcpsdk "github.com/mark3labs/mcp-go/mcp"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api/settings"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/store/storetest"
)

// testHarness wires an in-memory store plus an in-process MCP client
// straight against newMCPServer, so these tests exercise the real tool
// handlers (store calls, RBAC filtering, JSON encoding) without going
// through a real HTTP round trip - that path is covered separately by the
// end-to-end test in internal/api.
type testHarness struct {
	Store  *store.Store
	client *mcpclient.Client
}

func newTestHarness(t *testing.T) *testHarness {
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

	cfg := &config.Config{Encryption: config.EncryptionSettings{Key: "YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE="}}
	aiRegistry := ai.NewRegistry()
	settingsH := settings.NewHandler(s, cfg, aiRegistry)

	embedRegistry := ai.NewEmbedRegistry()
	embedRegistry.Register("fake", func(map[string]any) (ai.Embedder, error) {
		return fakeEmbedder{}, nil
	})

	mcpServer := newMCPServer(Deps{Store: s, Settings: settingsH, Embed: embedRegistry})
	client, err := mcpclient.NewInProcessClient(mcpServer)
	if err != nil {
		t.Fatalf("new in-process client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	if _, err := client.Initialize(context.Background(), mcpsdk.InitializeRequest{}); err != nil {
		t.Fatalf("initialize client: %v", err)
	}

	return &testHarness{Store: s, client: client}
}

// fakeEmbedder returns a fixed vector regardless of input, for search_docs
// tests that only need embedding + retrieval to run end to end, not to
// return meaningfully ranked results.
type fakeEmbedder struct{}

func (fakeEmbedder) Name() string { return "fake" }
func (fakeEmbedder) Embed(_ context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i := range texts {
		out[i] = []float32{1, 0, 0}
	}
	return out, nil
}

// callTool invokes name with args under ctx (built by userCtx/restrictedCtx
// below) and decodes the single text content block as JSON into out.
func (h *testHarness) callTool(ctx context.Context, t *testing.T, name string, args map[string]any, out any) *mcpsdk.CallToolResult {
	t.Helper()
	res, err := h.client.CallTool(ctx, mcpsdk.CallToolRequest{
		Params: mcpsdk.CallToolParams{Name: name, Arguments: args},
	})
	if err != nil {
		t.Fatalf("call tool %s: %v", name, err)
	}
	if res.IsError {
		t.Fatalf("tool %s returned error result: %s", name, mcpsdk.GetTextFromContent(res.Content[0]))
	}
	if out != nil {
		text := mcpsdk.GetTextFromContent(res.Content[0])
		if err := json.Unmarshal([]byte(text), out); err != nil {
			t.Fatalf("decode tool %s result: %v; body=%s", name, err, text)
		}
	}
	return res
}

// userCtx builds a context as AuthMiddleware would for a JWT-authenticated
// (unrestricted) user.
func userCtx(userID string) context.Context {
	return auth.ContextWithUser(context.Background(), userID, false)
}

// restrictedCtx builds a context as AuthMiddleware would for an API key
// restricted to connectorIDs.
func restrictedCtx(userID string, connectorIDs []string) context.Context {
	ctx := auth.ContextWithUser(context.Background(), userID, false)
	return auth.ContextWithAPIKeyRestriction(ctx, auth.APIKeyRestriction{ConnectorIDs: connectorIDs})
}

// createUser seeds a bare local user (no grants) and returns its ID.
func createUser(t *testing.T, s *store.Store) string {
	t.Helper()
	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	u := &store.User{
		Username: "user-" + uuid.New().String(), DisplayName: "Test User",
		InstanceAdminRole: "user", AuthSource: "local", PasswordHash: hash,
	}
	if err := s.CreateUser(context.Background(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u.ID
}

// createConnector seeds a minimal enabled connector and returns its ID.
func createConnector(t *testing.T, s *store.Store, name, category string) string {
	t.Helper()
	c := &store.ConnectorRecord{Name: name, Category: category, Type: "generic", URL: "http://example.invalid", Enabled: true}
	if err := s.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	return c.ID
}
