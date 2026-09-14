// Package docs provides API handlers for documentation CRUD and AI suggestions.
package docs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

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
	WSHub     *ws.Hub
}

// NewHandler creates a new doc handler.
func NewHandler(s *store.Store, eng *doc.Engine, settingsH *settings.Handler, aiRegistry *ai.Registry, hub *ws.Hub) *Handler {
	return &Handler{Store: s, DocEngine: eng, Settings: settingsH, AI: aiRegistry, WSHub: hub}
}

// Tree handles GET /api/docs/tree.
// Returns docs grouped by service (lab root + per-connector children).
func (h *Handler) Tree(w http.ResponseWriter, r *http.Request) {
	connectors, err := h.Store.ListAllConnectors(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
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
	httputil.WritePaginated(w, docs, page, pageSize, total)
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
	httputil.JSON(w, http.StatusOK, d)
}

// ByService handles GET /api/docs/service/{connectorId}.
func (h *Handler) ByService(w http.ResponseWriter, r *http.Request) {
	connectorID := r.PathValue("id")
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
		_ = h.Store.CreateDocVersion(r.Context(), &store.DocVersionRecord{
			DocID:   id,
			Rev:     d.CurrentVersion,
			Content: req.Content,
			Author:  userID,
			Trigger: trigger,
		})
	}

	httputil.JSON(w, http.StatusOK, d)
}

// Versions handles GET /api/docs/{id}/versions.
func (h *Handler) Versions(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	versions, err := h.Store.GetDocVersions(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, versions)
}

// Version handles GET /api/docs/{id}/versions/{rev}.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")
	rev, err := strconv.Atoi(r.PathValue("rev"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rev")
		return
	}

	versions, err := h.Store.GetDocVersions(r.Context(), docID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	for _, v := range versions {
		if v.Rev == rev {
			httputil.JSON(w, http.StatusOK, v)
			return
		}
	}
	httputil.Error(w, http.StatusNotFound, "not_found", "Version not found")
}

// Restore handles POST /api/docs/{id}/versions/{rev}/restore.
// Restores a doc to a previous version.
func (h *Handler) Restore(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")
	rev, err := strconv.Atoi(r.PathValue("rev"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid rev")
		return
	}

	versions, err := h.Store.GetDocVersions(r.Context(), docID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	var target *store.DocVersionRecord
	for _, v := range versions {
		if v.Rev == rev {
			target = &v
			break
		}
	}
	if target == nil {
		httputil.Error(w, http.StatusNotFound, "not_found", "Version not found")
		return
	}

	if err := h.Store.UpdateDoc(r.Context(), docID, target.Content, nil); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "doc.restore", "doc", docID, map[string]any{
		"rev": rev,
	}); err != nil {
		slog.Error("failed to record audit", "action", "doc.restore", "error", err)
	}

	d, _ := h.Store.GetDoc(r.Context(), docID)
	httputil.JSON(w, http.StatusOK, d)
}

// GetLock handles GET /api/docs/{id}/lock. Returns the current lock, or an
// empty object if the doc is unlocked/expired. Open to any authenticated role.
func (h *Handler) GetLock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	lock, err := h.Store.GetDocLock(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.JSON(w, http.StatusOK, map[string]any{})
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, lock)
}

// AcquireLock handles POST /api/docs/{id}/lock. Acquires or renews the
// advisory editing lock; 409 with the current holder if held by someone
// else. Operator-only, same gate as Save.
func (h *Handler) AcquireLock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := auth.UserIDFromContext(r.Context())

	lock, err := h.Store.AcquireDocLock(r.Context(), id, userID)
	if errors.Is(err, store.ErrLockHeldByOther) {
		httputil.JSON(w, http.StatusConflict, lock)
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if h.WSHub != nil {
		h.WSHub.Broadcast(ws.EventDocLockAcquired, lock)
	}
	httputil.JSON(w, http.StatusOK, lock)
}

// ReleaseLock handles POST /api/docs/{id}/lock/release. Operator-only.
func (h *Handler) ReleaseLock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := auth.UserIDFromContext(r.Context())

	if err := h.Store.ReleaseDocLock(r.Context(), id, userID); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if h.WSHub != nil {
		h.WSHub.Broadcast(ws.EventDocLockReleased, map[string]any{"docId": id, "userId": userID})
	}
	httputil.JSON(w, http.StatusOK, map[string]any{})
}

// Generate handles POST /api/docs/generate.
// Generates a doc from a template and connector snapshot.
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		TemplateID  string `json:"templateId"`
		ConnectorID string `json:"connectorId"`
	}](w, r)
	if !ok {
		return
	}
	if req.TemplateID == "" || req.ConnectorID == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "templateId and connectorId are required")
		return
	}

	result, err := h.DocEngine.GenerateFromTemplate(r.Context(), req.TemplateID, req.ConnectorID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, result)
}

// TemplateSchema handles GET /api/docs/template-schema.
// Returns the JSON schema of templateData (available fields/types) and the list
// of available template functions for use in template autocomplete/reference.
func (h *Handler) TemplateSchema(w http.ResponseWriter, _ *http.Request) {
	schema := map[string]any{
		"fields": map[string]map[string]any{
			"ServiceName": {
				"type":        "string",
				"description": "The name of the service",
			},
			"Type": {
				"type":        "string",
				"description": "The type of service",
			},
			"Sections": {
				"type":        "array",
				"description": "Array of snapshot sections",
				"items": map[string]any{
					"type": "object",
					"fields": map[string]map[string]any{
						"Title": {
							"type":        "string",
							"description": "Section title",
						},
						"Content": {
							"type":        "string",
							"description": "Section content",
						},
					},
				},
			},
			"Dependencies": {
				"type":        "array",
				"description": "Array of service dependencies",
				"items": map[string]any{
					"type": "object",
					"fields": map[string]map[string]any{
						"Kind": {
							"type":        "string",
							"description": "Dependency kind (host|network|storage|upstream_service)",
						},
						"Name": {
							"type":        "string",
							"description": "Dependency name",
						},
						"Ref": {
							"type":        "string",
							"description": "Optional connector/service ID reference",
						},
					},
				},
			},
			"Metadata": {
				"type":        "object",
				"description": "Key-value metadata map (map[string]string)",
			},
			"GeneratedAt": {
				"type":        "string",
				"description": "ISO-8601 timestamp when the template was rendered",
			},
		},
		"functions": []map[string]string{
			{
				"name":        "dateFormat",
				"description": "Format a time value or RFC3339 string using a Go layout string. Usage: .GeneratedAt | dateFormat \"2006-01-02\"",
			},
			{
				"name":        "truncate",
				"description": "Truncate a string to N characters, appending '...' if truncated. Usage: .SomeText | truncate 50",
			},
			{
				"name":        "toJSON",
				"description": "Marshal a value to a JSON string for embedding structured data. Usage: .Metadata | toJSON",
			},
			{
				"name":        "filterByTitle",
				"description": "Filter sections by exact title match. Usage: .Sections | filterByTitle \"Configuration\"",
			},
			{
				"name":        "join",
				"description": "Join an array of strings with a separator. Usage: .SomeStrings | join \", \"",
			},
		},
	}
	httputil.JSON(w, http.StatusOK, schema)
}

// AISuggest handles POST /api/docs/{id}/ai-suggest.
// Batched (non-streaming) suggestion: the full result is delivered over the
// doc.ai_suggestion WS event, correlated by the returned requestId.
func (h *Handler) AISuggest(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")

	var req struct {
		Prompt    string `json:"prompt"`
		Selection string `json:"selection"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	d, err := h.Store.GetDoc(r.Context(), docID)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	cfg := h.Settings.LoadAIConfig(r.Context())
	if !cfg.Enabled {
		httputil.Error(w, http.StatusConflict, "ai_disabled", "AI module is not enabled")
		return
	}

	provider, err := h.AI.Get(cfg.Provider, map[string]any{
		"apiKey": cfg.APIKey, "model": cfg.Model, "baseUrl": cfg.BaseURL,
	})
	if err != nil {
		httputil.Error(w, http.StatusConflict, "ai_disabled", fmt.Sprintf("AI provider unavailable: %v", err))
		return
	}

	prompt := req.Prompt
	if req.Selection != "" {
		prompt = fmt.Sprintf("%s\n\nSelected text:\n%s", prompt, req.Selection)
	}

	userID := auth.UserIDFromContext(r.Context())
	requestID := uuid.New().String()

	go func() {
		content, err := provider.Suggest(context.Background(), &ai.SuggestRequest{
			SystemPrompt: "You maintain internal infrastructure documentation. Suggest an improved version of the document based on the request.",
			UserPrompt:   prompt,
			DocContent:   d.Content,
		})
		payload := map[string]any{"docId": docID, "requestId": requestID}
		if err != nil {
			payload["status"] = "error"
			payload["error"] = err.Error()
		} else {
			payload["status"] = "complete"
			payload["fullContent"] = content
		}
		if h.WSHub != nil {
			h.WSHub.BroadcastToUser(userID, ws.EventDocAISuggestion, payload)
		}
	}()

	httputil.JSON(w, http.StatusAccepted, map[string]any{"requestId": requestID})
}
