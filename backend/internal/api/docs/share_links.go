package docs

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// shareLinkNode is the resolved meaning of a doc_tree_root value: "root"
// covers every connector, "connector" covers one connector's docs, "doc"
// covers exactly one doc (not its whole connector). Mirrors the node shape
// Tree() already builds rather than delegating resolution to store, keeping
// store CRUD-only.
type shareLinkNode struct {
	kind         string // "root" | "connector" | "doc"
	connectorIDs []string
	docID        string // set only when kind == "doc"
}

func (h *Handler) resolveShareLinkNode(w http.ResponseWriter, r *http.Request, docTreeRoot string) (shareLinkNode, bool) {
	connectors, err := h.Store.ListConnectorNames(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return shareLinkNode{}, false
	}

	if docTreeRoot == "root" {
		ids := make([]string, len(connectors))
		for i, c := range connectors {
			ids[i] = c.ID
		}
		return shareLinkNode{kind: "root", connectorIDs: ids}, true
	}

	// Try as a connector ID first: a real doc never collides with a
	// connector ID (both are UUIDs from disjoint tables), so this is safe.
	for _, c := range connectors {
		if c.ID == docTreeRoot {
			return shareLinkNode{kind: "connector", connectorIDs: []string{c.ID}}, true
		}
	}

	// Otherwise it must be a doc ID: resolve to its connector, but the
	// covered subtree is just that one doc, not its whole connector.
	d, err := h.Store.GetDoc(r.Context(), docTreeRoot)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "docTreeRoot does not refer to a known doc, connector, or the lab root")
		return shareLinkNode{}, false
	}
	if err != nil {
		httputil.Errorf(w, err)
		return shareLinkNode{}, false
	}
	if d.ServiceID == "" {
		// Lab-wide doc (e.g. Lab Topology): no connector to check against,
		// share-link creation requires instance-admin instead.
		if !auth.InstanceAdminFromContext(r.Context()) {
			httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
			return shareLinkNode{}, false
		}
		return shareLinkNode{kind: "doc", docID: d.ID}, true
	}
	return shareLinkNode{kind: "doc", connectorIDs: []string{d.ServiceID}, docID: d.ID}, true
}

// requireShareCreateAccess 403s unless the caller has operator on every
// connector covered by docTreeRoot.
func (h *Handler) requireShareCreateAccess(w http.ResponseWriter, r *http.Request, connectorIDs []string) bool {
	userID := auth.UserIDFromContext(r.Context())
	for _, connectorID := range connectorIDs {
		ok, err := h.Store.UserHasConnectorRole(r.Context(), userID, connectorID, "operator")
		if err != nil {
			httputil.Errorf(w, err)
			return false
		}
		if !ok {
			httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
			return false
		}
	}
	return true
}

const shareTokenPrefix = "wlz_share_"

func newShareToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate share token: %w", err)
	}
	return shareTokenPrefix + base64.RawURLEncoding.EncodeToString(buf), nil
}

// CreateShareLink handles POST /api/docs/share-links.
func (h *Handler) CreateShareLink(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		DocTreeRoot string `json:"docTreeRoot"`
		ExpiresAt   string `json:"expiresAt"`
	}](w, r)
	if !ok {
		return
	}
	if req.DocTreeRoot == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "docTreeRoot is required")
		return
	}
	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil || !expiresAt.After(time.Now().UTC()) {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "expiresAt is required and must be a future RFC3339 timestamp")
		return
	}

	node, ok := h.resolveShareLinkNode(w, r, req.DocTreeRoot)
	if !ok {
		return
	}
	if !h.requireShareCreateAccess(w, r, node.connectorIDs) {
		return
	}

	rawToken, err := newShareToken()
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	link := &store.ShareLink{
		TokenHash:   store.HashToken(rawToken),
		DocTreeRoot: req.DocTreeRoot,
		CreatedBy:   auth.UserIDFromContext(r.Context()),
		ExpiresAt:   req.ExpiresAt,
	}
	if err := h.Store.CreateShareLink(r.Context(), link); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "doc.share_link.create", "share_link", link.ID, map[string]any{
		"docTreeRoot": link.DocTreeRoot,
		"expiresAt":   link.ExpiresAt,
	}); err != nil {
		slog.Error("failed to record audit", "action", "doc.share_link.create", "error", err)
	}

	httputil.JSON(w, http.StatusCreated, map[string]any{
		"id":             link.ID,
		"docTreeRoot":    link.DocTreeRoot,
		"createdBy":      link.CreatedBy,
		"createdAt":      link.CreatedAt,
		"expiresAt":      link.ExpiresAt,
		"revokedAt":      link.RevokedAt,
		"lastAccessedAt": link.LastAccessedAt,
		"token":          rawToken,
	})
}

// ListShareLinks handles GET /api/docs/share-links. Returns only links
// created by the calling user and never exposes a token or hash.
func (h *Handler) ListShareLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.Store.ListShareLinksCreatedBy(r.Context(), auth.UserIDFromContext(r.Context()))
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, links)
}

// RevokeShareLink handles DELETE /api/docs/share-links/{id}.
func (h *Handler) RevokeShareLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	link, err := h.Store.GetShareLinkByID(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Share link not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if link.CreatedBy != auth.UserIDFromContext(r.Context()) && !auth.InstanceAdminFromContext(r.Context()) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Share link not found")
		return
	}
	if err := h.Store.RevokeShareLink(r.Context(), id); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "doc.share_link.revoke", "share_link", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "doc.share_link.revoke", "error", err)
	}
	httputil.NoContent(w)
}

// --- Unauthenticated share-link view routes ---

type shareLinkContextKey struct{}

// ResolveShareLink is middleware for the unauthenticated /api/share/{token}
// route group. It looks up the token, rejects revoked/expired links with
// distinct error codes (not a generic 401/404), touches last_accessed_at on
// success, and stores the resolved ShareLink + its covered connector IDs on
// the request context for the read handlers below.
func (h *Handler) ResolveShareLink(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		link, err := h.Store.GetShareLinkByHash(r.Context(), store.HashToken(token))
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Share link not found")
			return
		}
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		if link.RevokedAt != "" {
			httputil.Error(w, http.StatusGone, "share_link_revoked", "This share link has been revoked")
			return
		}
		expiresAt, err := time.Parse(time.RFC3339, link.ExpiresAt)
		if err != nil || !time.Now().UTC().Before(expiresAt) {
			httputil.Error(w, http.StatusGone, "share_link_expired", "This share link has expired")
			return
		}

		node, ok := h.resolveShareLinkNode(w, r, link.DocTreeRoot)
		if !ok {
			return
		}

		if err := h.Store.TouchShareLinkLastAccessed(r.Context(), link.ID); err != nil {
			slog.Error("failed to touch share link last accessed", "id", link.ID, "error", err)
		}

		ctx := contextWithShareLink(r.Context(), &shareLinkScope{link: link, node: node})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type shareLinkScope struct {
	link *store.ShareLink
	node shareLinkNode
}

func contextWithShareLink(ctx context.Context, scope *shareLinkScope) context.Context {
	return context.WithValue(ctx, shareLinkContextKey{}, scope)
}

func shareLinkFromContext(ctx context.Context) *shareLinkScope {
	scope, _ := ctx.Value(shareLinkContextKey{}).(*shareLinkScope)
	return scope
}

// ShareLinkTree handles GET /api/share/{token}/tree. Read-only: the lab
// tree scoped to the share link's covered subtree.
func (h *Handler) ShareLinkTree(w http.ResponseWriter, r *http.Request) {
	scope := shareLinkFromContext(r.Context())
	node := scope.node

	type TreeNode struct {
		ID       string     `json:"docId"`
		Title    string     `json:"title"`
		Kind     string     `json:"kind"`
		Children []TreeNode `json:"children,omitempty"`
	}
	root := TreeNode{ID: "root", Title: "Shared Documentation", Kind: "lab"}

	if node.kind == "doc" {
		// A doc-scoped link covers exactly that one doc, not its whole
		// connector — even if the doc happens to have connector siblings.
		d, err := h.Store.GetDoc(r.Context(), node.docID)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		root.Children = []TreeNode{{ID: d.ID, Title: d.Title, Kind: d.Kind}}
		httputil.JSON(w, http.StatusOK, root)
		return
	}

	allowed := make(map[string]bool, len(node.connectorIDs))
	for _, id := range node.connectorIDs {
		allowed[id] = true
	}
	connectors, err := h.Store.ListConnectorNames(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	docsByService, err := h.Store.ListDocsGroupedByService(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	for _, c := range connectors {
		if !allowed[c.ID] {
			continue
		}
		connNode := TreeNode{ID: c.ID, Title: c.Name, Kind: "service"}
		for _, d := range docsByService[c.ID] {
			connNode.Children = append(connNode.Children, TreeNode{ID: d.ID, Title: d.Title, Kind: d.Kind})
		}
		root.Children = append(root.Children, connNode)
	}
	if root.Children == nil {
		root.Children = []TreeNode{}
	}
	httputil.JSON(w, http.StatusOK, root)
}

// ShareLinkDoc handles GET /api/share/{token}/docs/{docId}. Read-only,
// rejects any doc outside the share link's covered subtree with 404 (not
// 403, to avoid confirming existence of docs outside the shared scope).
func (h *Handler) ShareLinkDoc(w http.ResponseWriter, r *http.Request) {
	scope := shareLinkFromContext(r.Context())
	node := scope.node
	docID := r.PathValue("docId")

	if node.kind == "doc" {
		if docID != node.docID {
			httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
			return
		}
		d, err := h.Store.GetDoc(r.Context(), docID)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		httputil.JSON(w, http.StatusOK, d)
		return
	}

	d, err := h.Store.GetDoc(r.Context(), docID)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	inScope := false
	for _, id := range node.connectorIDs {
		if id == d.ServiceID {
			inScope = true
			break
		}
	}
	if !inScope {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}

	httputil.JSON(w, http.StatusOK, d)
}
