// Package changes provides API handlers for infrastructure change records.
package changes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/api/settings"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/sync"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// Handler holds dependencies for change endpoints.
type Handler struct {
	Store    *store.Store
	Settings *settings.Handler
	AI       *ai.Registry
	WSHub    *ws.Hub
}

const maxBulkIDs = 500

// NewHandler creates a new change handler.
func NewHandler(s *store.Store, settingsH *settings.Handler, aiRegistry *ai.Registry, hub *ws.Hub) *Handler {
	return &Handler{Store: s, Settings: settingsH, AI: aiRegistry, WSHub: hub}
}

// diffToSpec converts the stored []sync.DiffPatch JSON into the spec's Diff{format,hunks} shape.
func diffToSpec(raw string) map[string]any {
	var patches []sync.DiffPatch
	_ = json.Unmarshal([]byte(raw), &patches)

	hunks := make([]map[string]any, 0, len(patches))
	for _, p := range patches {
		var before, after any
		if p.Old != "" {
			before = p.Old
		}
		if p.New != "" {
			after = p.New
		}
		hunks = append(hunks, map[string]any{
			"path":   p.Section,
			"before": before,
			"after":  after,
		})
	}
	return map[string]any{
		"format": "infra",
		"hunks":  hunks,
	}
}

// changeDetail builds the spec's ChangeDetail response for a change record.
func (h *Handler) changeDetail(ctx context.Context, id string) (map[string]any, error) {
	c, err := h.Store.GetChange(ctx, id)
	if err != nil {
		return nil, err
	}

	serviceName := ""
	if conn, err := h.Store.GetConnector(ctx, c.ServiceID); err == nil {
		serviceName = conn.Name
	}

	var affectedDocIDs []string
	_ = json.Unmarshal([]byte(c.AffectedDocIDs), &affectedDocIDs)
	if affectedDocIDs == nil {
		affectedDocIDs = []string{}
	}

	return map[string]any{
		"id":             c.ID,
		"serviceId":      c.ServiceID,
		"serviceName":    serviceName,
		"changeType":     c.ChangeType,
		"severity":       c.Severity,
		"summary":        c.Summary,
		"willTriggerAi":  false,
		"detectedAt":     c.DetectedAt,
		"status":         c.Status,
		"diff":           diffToSpec(c.Diff),
		"affectedDocIds": affectedDocIDs,
		"narration":      c.Narration,
	}, nil
}

// List handles GET /api/changes.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)
	serviceID := r.URL.Query().Get("serviceId")
	severity := r.URL.Query().Get("severity")

	changes, total, err := h.Store.ListChanges(r.Context(), serviceID, severity, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.WritePaginated(w, changes, page, pageSize, total)
}

// Get handles GET /api/changes/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	detail, err := h.changeDetail(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Change not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, detail)
}

// Acknowledge handles POST /api/changes/{id}/ack.
func (h *Handler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Store.UpdateChangeStatus(r.Context(), id, "acknowledged"); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Change not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "change.ack", "change", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "change.ack", "error", err)
	}
	detail, err := h.changeDetail(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, detail)
}

// Dismiss handles POST /api/changes/{id}/dismiss.
func (h *Handler) Dismiss(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.Store.UpdateChangeStatus(r.Context(), id, "dismissed"); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Change not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "change.dismiss", "change", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "change.dismiss", "error", err)
	}
	detail, err := h.changeDetail(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, detail)
}

// lowRiskBulkStatuses maps the bulk-resolve request's target status to the
// audit action name recorded per successfully-resolved item.
var lowRiskBulkStatuses = map[string]string{
	"acknowledged": "change.bulk_ack",
	"dismissed":    "change.bulk_dismiss",
}

// bulkResolveRequest is the body of POST /api/changes/bulk-resolve.
type bulkResolveRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

// bulkResolveItemResult is the per-item outcome in the bulk-resolve response.
type bulkResolveItemResult struct {
	ID     string `json:"id"`
	Status string `json:"status"` // "success" | "error"
	Reason string `json:"reason,omitempty"`
}

// BulkResolve handles POST /api/changes/bulk-resolve. It acknowledges or
// dismisses an explicit, caller-supplied list of change IDs in one request.
//
// Design decision (issue #27): "low-risk" is defined server-side as
// severity != critical (info and warning qualify). The client's selection is
// never trusted — each ID is independently re-checked for existence and
// severity here, since this is a safety-relevant control. One bad ID never
// aborts the batch: every item gets its own success/error outcome, and one
// audit record is written per successfully-resolved item (not one for the
// whole batch) so provenance stays per-change.
func (h *Handler) BulkResolve(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[bulkResolveRequest](w, r)
	if !ok {
		return
	}
	auditAction, ok := lowRiskBulkStatuses[req.Status]
	if !ok {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "status must be acknowledged or dismissed")
		return
	}
	if len(req.IDs) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "ids must be a non-empty array")
		return
	}
	if len(req.IDs) > maxBulkIDs {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "ids must contain at most 500 items")
		return
	}

	results := make([]bulkResolveItemResult, 0, len(req.IDs))
	changes, err := h.Store.ListChangesByID(r.Context(), req.IDs)
	if err != nil {
		for _, id := range req.IDs {
			results = append(results, bulkResolveItemResult{ID: id, Status: "error", Reason: "internal_error"})
		}
		httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
		return
	}
	resolvedIDs := make([]string, 0, len(changes))
	auditRecords := make([]store.AuditRecord, 0, len(changes))
	eligible := make(map[string]bool, len(changes))
	for _, id := range req.IDs {
		c, ok := changes[id]
		if !ok || c.Severity == "critical" {
			continue
		}
		resolvedIDs = append(resolvedIDs, id)
		eligible[id] = true
		auditRecords = append(auditRecords, store.AuditRecord{TargetID: id, Detail: fmt.Sprintf(`{"severity":%q}`, c.Severity)})
	}
	updateErr := h.Store.UpdateChangeStatuses(r.Context(), resolvedIDs, req.Status)
	if updateErr == nil {
		if err := h.Store.RecordAuditBatchFromContext(r.Context(), auditAction, "change", auditRecords); err != nil {
			slog.Error("failed to record audit", "action", auditAction, "error", err)
		}
	}
	for _, id := range req.IDs {
		c, ok := changes[id]
		switch {
		case !ok:
			results = append(results, bulkResolveItemResult{ID: id, Status: "error", Reason: "not_found"})
		case c.Severity == "critical":
			results = append(results, bulkResolveItemResult{ID: id, Status: "error", Reason: "not_low_risk"})
		case updateErr != nil && eligible[id]:
			results = append(results, bulkResolveItemResult{ID: id, Status: "error", Reason: "internal_error"})
		default:
			results = append(results, bulkResolveItemResult{ID: id, Status: "success"})
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
}

// AIUpdate handles POST /api/changes/{id}/ai-update.
// Batched (non-streaming) suggestion: the full result is delivered over the
// doc.ai_suggestion WS event, correlated by the returned requestId.
func (h *Handler) AIUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	c, err := h.Store.GetChange(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Change not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	cfg := h.Settings.LoadAIConfig(r.Context())
	if !cfg.Enabled || len(cfg.Providers) == 0 {
		httputil.Error(w, http.StatusConflict, "ai_disabled", "AI module is not enabled")
		return
	}

	var affectedDocIDs []string
	_ = json.Unmarshal([]byte(c.AffectedDocIDs), &affectedDocIDs)
	docID := ""
	if len(affectedDocIDs) > 0 {
		docID = affectedDocIDs[0]
	}

	userID := auth.UserIDFromContext(r.Context())
	requestID := uuid.New().String()

	go func() {
		result, err := ai.SuggestWithFallback(context.Background(), h.AI, cfg.Providers, &ai.SuggestRequest{
			SystemPrompt: "Summarize this infrastructure change and suggest an updated documentation snippet.",
			UserPrompt:   fmt.Sprintf("Change summary: %s\n\nDiff:\n%s", c.Summary, c.Diff),
		})
		payload := map[string]any{"requestId": requestID}
		if docID != "" {
			payload["docId"] = docID
		}
		if err != nil {
			payload["status"] = "error"
			payload["error"] = err.Error()
		} else {
			payload["status"] = "complete"
			payload["fullContent"] = result.Content
			payload["provider"] = result.Provider
			payload["fallbackUsed"] = result.FallbackUsed
		}
		if h.WSHub != nil {
			h.WSHub.BroadcastToUser(userID, ws.EventDocAISuggestion, payload)
		}
	}()

	httputil.JSON(w, http.StatusAccepted, map[string]any{"requestId": requestID})
}

// Explain handles POST /api/changes/{id}/explain. On-demand, plain-English
// narration of why the change matters — generated once and cached on the
// Change record; later calls reuse the stored text instead of re-invoking
// the AI provider.
func (h *Handler) Explain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	c, err := h.Store.GetChange(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Change not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	if c.Narration != "" {
		detail, err := h.changeDetail(r.Context(), id)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		httputil.JSON(w, http.StatusOK, detail)
		return
	}

	cfg := h.Settings.LoadAIConfig(r.Context())
	if !cfg.Enabled || len(cfg.Providers) == 0 {
		httputil.Error(w, http.StatusConflict, "ai_disabled", "AI module is not enabled")
		return
	}

	result, err := ai.SuggestWithFallback(r.Context(), h.AI, cfg.Providers, &ai.SuggestRequest{
		SystemPrompt: "You explain infrastructure changes to engineers in plain English. " +
			"In 2-4 sentences, explain why this change matters and what its practical impact is. " +
			"Do not restate the mechanical diff line by line.",
		UserPrompt: fmt.Sprintf("Change summary: %s\n\nDiff:\n%s", c.Summary, c.Diff),
	})
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, "ai_error", fmt.Sprintf("AI provider failed: %v", err))
		return
	}

	if err := h.Store.UpdateChangeNarration(r.Context(), id, result.Content); err != nil {
		httputil.Errorf(w, err)
		return
	}

	detail, err := h.changeDetail(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	detail["provider"] = result.Provider
	detail["fallbackUsed"] = result.FallbackUsed
	httputil.JSON(w, http.StatusOK, detail)
}
