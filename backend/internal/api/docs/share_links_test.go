package docs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedConnector(t *testing.T, s *store.Store) *store.ConnectorRecord {
	t.Helper()
	conn := &store.ConnectorRecord{
		Name:     "Test Connector",
		Category: "virtualization",
		Type:     "test-type",
		URL:      "https://test.example.com",
	}
	if err := s.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	return conn
}

func seedDoc(t *testing.T, s *store.Store, connectorID string) *store.DocRecord {
	t.Helper()
	d := &store.DocRecord{
		Title:     "Test Doc",
		Kind:      "service",
		ServiceID: connectorID,
		Content:   "test content",
	}
	if err := s.CreateDoc(context.Background(), d); err != nil {
		t.Fatalf("create doc: %v", err)
	}
	return d
}

func asUser(req *http.Request, userID string, instanceAdmin bool) *http.Request {
	return req.WithContext(auth.ContextWithUser(req.Context(), userID, instanceAdmin))
}

func futureExpiry() string {
	return time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)
}

// --- Creation ---

func TestCreateShareLinkRequiresOperatorOnConnector(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")

	body := `{"docTreeRoot":"` + conn.ID + `","expiresAt":"` + futureExpiry() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateShareLinkSucceedsForOperator(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, conn.ID, "operator")

	body := `{"docTreeRoot":"` + conn.ID + `","expiresAt":"` + futureExpiry() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	token, _ := resp["token"].(string)
	if !strings.HasPrefix(token, "wlz_share_") {
		t.Errorf("token = %q, want wlz_share_ prefix", token)
	}
	if _, leaked := resp["tokenHash"]; leaked {
		t.Error("response leaked tokenHash")
	}
}

func TestCreateShareLinkRootRequiresOperatorOnEveryConnector(t *testing.T) {
	h := newTestHandler(t)
	connA := seedConnector(t, h.Store)
	connB := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")
	// Operator on A only, not B — root subtree covers both.
	apitest.GrantConnectorRole(t, h.Store, userID, connA.ID, "operator")

	body := `{"docTreeRoot":"root","expiresAt":"` + futureExpiry() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 (missing operator on connB=%s); body=%s", rr.Code, connB.ID, rr.Body.String())
	}
}

func TestCreateShareLinkRootSucceedsWhenOperatorOnEveryConnector(t *testing.T) {
	h := newTestHandler(t)
	connA := seedConnector(t, h.Store)
	connB := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, connA.ID, "operator")
	apitest.GrantConnectorRole(t, h.Store, userID, connB.ID, "operator")

	body := `{"docTreeRoot":"root","expiresAt":"` + futureExpiry() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateShareLinkMissingExpiryRejected(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, conn.ID, "operator")

	body := `{"docTreeRoot":"` + conn.ID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateShareLinkPastExpiryRejected(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, conn.ID, "operator")

	past := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	body := `{"docTreeRoot":"` + conn.ID + `","expiresAt":"` + past + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateShareLinkUnknownDocTreeRootRejected(t *testing.T) {
	h := newTestHandler(t)
	userID := apitest.NewUser(t, h.Store, "operator")

	body := `{"docTreeRoot":"does-not-exist","expiresAt":"` + futureExpiry() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, true)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateShareLinkForDocResolvesToItsConnector(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	d := seedDoc(t, h.Store, conn.ID)
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, conn.ID, "operator")

	body := `{"docTreeRoot":"` + d.ID + `","expiresAt":"` + futureExpiry() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/docs/share-links", strings.NewReader(body))
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.CreateShareLink(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rr.Code, rr.Body.String())
	}
}

// --- List / Revoke ---

func TestListShareLinksNeverLeaksTokenHash(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	userID := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, userID, conn.ID, "operator")

	createReq := httptest.NewRequest(http.MethodPost, "/api/docs/share-links",
		strings.NewReader(`{"docTreeRoot":"`+conn.ID+`","expiresAt":"`+futureExpiry()+`"}`))
	createReq = asUser(createReq, userID, false)
	h.CreateShareLink(httptest.NewRecorder(), createReq)

	listReq := httptest.NewRequest(http.MethodGet, "/api/docs/share-links", nil)
	listReq = asUser(listReq, userID, false)
	rr := httptest.NewRecorder()
	h.ListShareLinks(rr, listReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "tokenHash") || strings.Contains(rr.Body.String(), "wlz_share_") {
		t.Errorf("list response leaked token/hash: %s", rr.Body.String())
	}
}

func TestListShareLinksOnlyReturnsOwnLinks(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	owner := apitest.NewUser(t, h.Store, "viewer")
	other := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, owner, conn.ID, "operator")

	createReq := httptest.NewRequest(http.MethodPost, "/api/docs/share-links",
		strings.NewReader(`{"docTreeRoot":"`+conn.ID+`","expiresAt":"`+futureExpiry()+`"}`))
	createReq = asUser(createReq, owner, false)
	h.CreateShareLink(httptest.NewRecorder(), createReq)

	listReq := httptest.NewRequest(http.MethodGet, "/api/docs/share-links", nil)
	listReq = asUser(listReq, other, false)
	rr := httptest.NewRecorder()
	h.ListShareLinks(rr, listReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	var links []map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &links); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("got %d links for other user, want 0", len(links))
	}
}

func TestRevokeShareLinkByNonOwnerNotFound(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	owner := apitest.NewUser(t, h.Store, "viewer")
	other := apitest.NewUser(t, h.Store, "viewer")
	apitest.GrantConnectorRole(t, h.Store, owner, conn.ID, "operator")

	createReq := httptest.NewRequest(http.MethodPost, "/api/docs/share-links",
		strings.NewReader(`{"docTreeRoot":"`+conn.ID+`","expiresAt":"`+futureExpiry()+`"}`))
	createReq = asUser(createReq, owner, false)
	createRR := httptest.NewRecorder()
	h.CreateShareLink(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	id := created["id"].(string)

	revokeReq := httptest.NewRequest(http.MethodDelete, "/api/docs/share-links/"+id, nil)
	revokeReq.SetPathValue("id", id)
	revokeReq = asUser(revokeReq, other, false)
	rr := httptest.NewRecorder()
	h.RevokeShareLink(rr, revokeReq)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rr.Code, rr.Body.String())
	}
}

func TestRevokeShareLinkByInstanceAdminAllowed(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	owner := apitest.NewUser(t, h.Store, "viewer")
	admin := apitest.NewUser(t, h.Store, "operator")
	apitest.GrantConnectorRole(t, h.Store, owner, conn.ID, "operator")

	createReq := httptest.NewRequest(http.MethodPost, "/api/docs/share-links",
		strings.NewReader(`{"docTreeRoot":"`+conn.ID+`","expiresAt":"`+futureExpiry()+`"}`))
	createReq = asUser(createReq, owner, false)
	createRR := httptest.NewRecorder()
	h.CreateShareLink(createRR, createReq)
	var created map[string]any
	if err := json.Unmarshal(createRR.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	id := created["id"].(string)

	revokeReq := httptest.NewRequest(http.MethodDelete, "/api/docs/share-links/"+id, nil)
	revokeReq.SetPathValue("id", id)
	revokeReq = asUser(revokeReq, admin, true)
	rr := httptest.NewRecorder()
	h.RevokeShareLink(rr, revokeReq)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rr.Code, rr.Body.String())
	}
}

func TestRevokeShareLinkUnknownID(t *testing.T) {
	h := newTestHandler(t)
	userID := apitest.NewUser(t, h.Store, "viewer")
	req := httptest.NewRequest(http.MethodDelete, "/api/docs/share-links/missing", nil)
	req.SetPathValue("id", "missing")
	req = asUser(req, userID, false)
	rr := httptest.NewRecorder()
	h.RevokeShareLink(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rr.Code, rr.Body.String())
	}
}

// --- Unauthenticated view route ---

func createTestShareLink(t *testing.T, h *Handler, docTreeRoot, expiresAt string) string {
	t.Helper()
	link := &store.ShareLink{
		TokenHash:   store.HashToken("test-token-" + docTreeRoot),
		DocTreeRoot: docTreeRoot,
		CreatedBy:   apitest.NewUser(t, h.Store, "viewer"),
		ExpiresAt:   expiresAt,
	}
	if err := h.Store.CreateShareLink(context.Background(), link); err != nil {
		t.Fatalf("create share link: %v", err)
	}
	return "test-token-" + docTreeRoot
}

func TestResolveShareLinkUnknownToken(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/share/nope/tree", nil)
	req.SetPathValue("token", "nope")
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rr.Code, rr.Body.String())
	}
}

func TestResolveShareLinkRevoked(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	token := createTestShareLink(t, h, conn.ID, futureExpiry())

	link, err := h.Store.GetShareLinkByHash(context.Background(), store.HashToken(token))
	if err != nil {
		t.Fatalf("get share link: %v", err)
	}
	if err := h.Store.RevokeShareLink(context.Background(), link.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/tree", nil)
	req.SetPathValue("token", token)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rr, req)
	if rr.Code != http.StatusGone {
		t.Fatalf("status = %d, want 410; body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "share_link_revoked") {
		t.Errorf("expected share_link_revoked code, got: %s", rr.Body.String())
	}
}

func TestResolveShareLinkExpired(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	past := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	token := createTestShareLink(t, h, conn.ID, past)

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/tree", nil)
	req.SetPathValue("token", token)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rr, req)
	if rr.Code != http.StatusGone {
		t.Fatalf("status = %d, want 410; body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "share_link_expired") {
		t.Errorf("expected share_link_expired code, got: %s", rr.Body.String())
	}
}

func TestResolveShareLinkTouchesLastAccessed(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	token := createTestShareLink(t, h, conn.ID, futureExpiry())

	link, err := h.Store.GetShareLinkByHash(context.Background(), store.HashToken(token))
	if err != nil {
		t.Fatalf("get share link: %v", err)
	}
	if link.LastAccessedAt != "" {
		t.Fatalf("expected empty LastAccessedAt before first view, got %q", link.LastAccessedAt)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/tree", nil)
	req.SetPathValue("token", token)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkTree)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}

	after, err := h.Store.GetShareLinkByID(context.Background(), link.ID)
	if err != nil {
		t.Fatalf("get share link by id: %v", err)
	}
	if after.LastAccessedAt == "" {
		t.Error("expected LastAccessedAt to be set after a successful view")
	}
}

func TestShareLinkTreeOnlyContainsCoveredConnector(t *testing.T) {
	h := newTestHandler(t)
	connA := seedConnector(t, h.Store)
	connB := seedConnector(t, h.Store)
	seedDoc(t, h.Store, connA.ID)
	seedDoc(t, h.Store, connB.ID)
	token := createTestShareLink(t, h, connA.ID, futureExpiry())

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/tree", nil)
	req.SetPathValue("token", token)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkTree)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), connA.ID) {
		t.Errorf("expected covered connector %s in tree, body=%s", connA.ID, rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), connB.ID) {
		t.Errorf("uncovered connector %s leaked into scoped tree, body=%s", connB.ID, rr.Body.String())
	}
}

func TestShareLinkDocOutsideSubtreeNotFound(t *testing.T) {
	h := newTestHandler(t)
	connA := seedConnector(t, h.Store)
	connB := seedConnector(t, h.Store)
	docB := seedDoc(t, h.Store, connB.ID)
	token := createTestShareLink(t, h, connA.ID, futureExpiry())

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/docs/"+docB.ID, nil)
	req.SetPathValue("token", token)
	req.SetPathValue("docId", docB.ID)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkDoc)).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for doc outside shared subtree; body=%s", rr.Code, rr.Body.String())
	}
}

func TestShareLinkDocInsideSubtreeServed(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	d := seedDoc(t, h.Store, conn.ID)
	token := createTestShareLink(t, h, conn.ID, futureExpiry())

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/docs/"+d.ID, nil)
	req.SetPathValue("token", token)
	req.SetPathValue("docId", d.ID)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkDoc)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), d.Title) {
		t.Errorf("expected doc content in response, got: %s", rr.Body.String())
	}
}

func TestShareLinkDocScopedToSingleDoc(t *testing.T) {
	h := newTestHandler(t)
	conn := seedConnector(t, h.Store)
	d1 := seedDoc(t, h.Store, conn.ID)
	d2 := seedDoc(t, h.Store, conn.ID)
	// Share link scoped to a single doc ID, not the whole connector.
	token := createTestShareLink(t, h, d1.ID, futureExpiry())

	req := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/docs/"+d1.ID, nil)
	req.SetPathValue("token", token)
	req.SetPathValue("docId", d1.ID)
	rr := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkDoc)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 for the shared doc itself; body=%s", rr.Code, rr.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/share/"+token+"/docs/"+d2.ID, nil)
	req2.SetPathValue("token", token)
	req2.SetPathValue("docId", d2.ID)
	rr2 := httptest.NewRecorder()
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkDoc)).ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for a sibling doc not covered by the doc-scoped link; body=%s", rr2.Code, rr2.Body.String())
	}
}

func newDocsEngineHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	return NewHandler(s, doc.NewEngine(s), nil, nil, nil, nil)
}

func TestShareLinkHandlerHasNoWriteMethods(t *testing.T) {
	// Regression guard: the unauthenticated share-link route group in
	// router.go must only ever wire ShareLinkTree/ShareLinkDoc (both
	// read-only). This test doesn't inspect router.go directly (no access
	// to it from this package) — it documents the invariant the router
	// wiring must uphold: no Save/AcquireLock/AISuggest-style write handler
	// may be reachable through docH.ResolveShareLink.
	h := newDocsEngineHandler(t)
	conn := seedConnector(t, h.Store)
	d := seedDoc(t, h.Store, conn.ID)
	token := createTestShareLink(t, h, conn.ID, futureExpiry())

	req := httptest.NewRequest(http.MethodPut, "/api/share/"+token+"/docs/"+d.ID, strings.NewReader(`{"content":"hacked"}`))
	req.SetPathValue("token", token)
	req.SetPathValue("docId", d.ID)
	rr := httptest.NewRecorder()
	// ShareLinkDoc is read-only regardless of HTTP method (router.go only
	// ever wires it to GET); calling it directly here proves it never
	// mutates even if invoked out of band.
	h.ResolveShareLink(http.HandlerFunc(h.ShareLinkDoc)).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rr.Code, rr.Body.String())
	}
	after, err := h.Store.GetDoc(context.Background(), d.ID)
	if err != nil {
		t.Fatalf("get doc: %v", err)
	}
	if after.Content != "test content" {
		t.Errorf("doc content changed via share-link route: got %q", after.Content)
	}
}
