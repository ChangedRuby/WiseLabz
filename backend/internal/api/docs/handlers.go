// Package docs provides API handlers for documentation CRUD and AI suggestions.
package docs

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api/settings"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/doc"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// Handler holds dependencies for doc endpoints.
type Handler struct {
	Store     *store.Store
	DocEngine *doc.Engine
	Settings  *settings.Handler
	AI        *ai.Registry
	Embed     *ai.EmbedRegistry
	WSHub     *ws.Hub
}

// NewHandler creates a new doc handler.
func NewHandler(s *store.Store, eng *doc.Engine, settingsH *settings.Handler, aiRegistry *ai.Registry, embedRegistry *ai.EmbedRegistry, hub *ws.Hub) *Handler {
	return &Handler{Store: s, DocEngine: eng, Settings: settingsH, AI: aiRegistry, Embed: embedRegistry, WSHub: hub}
}

// Tree handles GET /api/docs/tree.
// Returns docs grouped by service (lab root + per-connector children).
func (h *Handler) Tree(w http.ResponseWriter, r *http.Request) {
	connectors, err := h.Store.ListConnectorNames(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	connectorIDs := make([]string, len(connectors))
	for i, c := range connectors {
		connectorIDs[i] = c.ID
	}
	allowedIDs, err := h.Store.FilterConnectorIDsByGrant(r.Context(), auth.UserIDFromContext(r.Context()), connectorIDs, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	isAllowed := make(map[string]bool, len(allowedIDs))
	for _, id := range allowedIDs {
		isAllowed[id] = true
	}
	docsByService, err := h.Store.ListDocsGroupedByService(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	type TreeNode struct {
		ID       string     `json:"docId"`
		Title    string     `json:"title"`
		Kind     string     `json:"kind"`
		Children []TreeNode `json:"children,omitempty"`
	}

	root := TreeNode{
		ID:    "root",
		Title: "Lab Documentation",
		Kind:  "lab",
	}

	for _, c := range connectors {
		if !isAllowed[c.ID] {
			continue
		}
		connNode := TreeNode{
			ID:    c.ID,
			Title: c.Name,
			Kind:  "service",
		}
		for _, d := range docsByService[c.ID] {
			connNode.Children = append(connNode.Children, TreeNode{
				ID:    d.ID,
				Title: d.Title,
				Kind:  d.Kind,
			})
		}
		root.Children = append(root.Children, connNode)
	}

	if root.Children == nil {
		root.Children = []TreeNode{}
	}

	httputil.JSON(w, http.StatusOK, root)
}

// List handles GET /api/docs.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)
	search := r.URL.Query().Get("search")

	docs, total, err := h.Store.ListAllDocs(r.Context(), search, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	serviceIDs := make([]string, 0, len(docs))
	for _, d := range docs {
		if d.ServiceID != "" {
			serviceIDs = append(serviceIDs, d.ServiceID)
		}
	}
	allowed, err := h.Store.FilterConnectorIDsByGrant(r.Context(), auth.UserIDFromContext(r.Context()), serviceIDs, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	isAllowed := make(map[string]bool, len(allowed))
	for _, id := range allowed {
		isAllowed[id] = true
	}
	filtered := make([]store.DocRecord, 0, len(docs))
	for _, d := range docs {
		if d.ServiceID == "" || isAllowed[d.ServiceID] {
			filtered = append(filtered, d)
		}
	}
	httputil.WritePaginated(w, filtered, page, pageSize, total)
}

// Get handles GET /api/docs/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// The tree returns docId "root" for the lab overview. There is no real doc
	// with that ID — return a synthetic placeholder.
	if id == "root" {
		httputil.JSON(w, http.StatusOK, map[string]any{
			"docId":          "root",
			"title":          "Lab Documentation",
			"content":        "_Welcome to your lab documentation. Sync a service to populate this page._",
			"currentVersion": 1,
			"kind":           "lab",
			"serviceId":      "",
		})
		return
	}

	d, err := h.Store.GetDoc(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		// Not a doc ID — maybe it's a connector ID. Fall back to service lookup.
		if !h.requireDocViewer(w, r, id) {
			return
		}
		docs, svcErr := h.Store.ListDocsByService(r.Context(), id)
		if svcErr == nil && len(docs) > 0 {
			httputil.JSON(w, http.StatusOK, docs[0])
			return
		}
		if svcErr == nil {
			httputil.JSON(w, http.StatusOK, map[string]any{
				"docId":          "",
				"title":          "No documentation yet",
				"content":        "Run a sync to generate documentation for this service.",
				"currentVersion": 0,
				"kind":           "service",
				"serviceId":      id,
			})
			return
		}
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.requireDocViewer(w, r, d.ServiceID) {
		return
	}
	httputil.JSON(w, http.StatusOK, d)
}

// ByService handles GET /api/docs/service/{connectorId}.
func (h *Handler) ByService(w http.ResponseWriter, r *http.Request) {
	connectorID := r.PathValue("id")
	if !h.requireDocViewer(w, r, connectorID) {
		return
	}
	docs, err := h.Store.ListDocsByService(r.Context(), connectorID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if len(docs) == 0 {
		httputil.JSON(w, http.StatusOK, map[string]any{
			"docId":          "",
			"title":          "No documentation yet",
			"content":        "Run a sync to generate documentation for this service.",
			"currentVersion": 0,
			"kind":           "service",
			"serviceId":      connectorID,
		})
		return
	}
	httputil.JSON(w, http.StatusOK, docs[0])
}

// Save handles PUT /api/docs/{id}.
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, ok := httputil.DecodeJSON[struct {
		Content     string `json:"content"`
		BaseVersion *int   `json:"baseVersion"`
		Trigger     string `json:"trigger"`
	}](w, r)
	if !ok {
		return
	}

	existing, err := h.Store.GetDoc(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.requireDocOperator(w, r, existing.ServiceID) {
		return
	}

	if err := h.Store.UpdateDoc(r.Context(), id, req.Content, req.BaseVersion); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
			return
		}
		if errors.Is(err, store.ErrVersionConflict) {
			current, getErr := h.Store.GetDoc(r.Context(), id)
			if getErr != nil {
				httputil.Errorf(w, getErr)
				return
			}
			httputil.JSON(w, http.StatusConflict, current)
			return
		}
		httputil.Errorf(w, err)
		return
	}

	// Create version
	d, _ := h.Store.GetDoc(r.Context(), id)
	userID := auth.UserIDFromContext(r.Context())
	if d != nil {
		trigger := req.Trigger
		if trigger == "" {
			trigger = "manual"
		}
		if err := h.Store.CreateDocVersion(r.Context(), &store.DocVersionRecord{
			DocID:   id,
			Rev:     d.CurrentVersion,
			Content: req.Content,
			Author:  userID,
			Trigger: trigger,
		}); err != nil {
			slog.Error("failed to record doc version", "docId", id, "rev", d.CurrentVersion, "error", err)
		}
	}

	if d != nil {
		h.syncDocEmbeddings(r.Context(), d.ID, d.Content)
	}
	httputil.JSON(w, http.StatusOK, d)
}
