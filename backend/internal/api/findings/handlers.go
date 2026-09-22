// Package findings provides API handlers for documentation quality findings.
package findings

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for quality finding endpoints.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a quality finding handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

func (h *Handler) response(ctx context.Context, finding store.QualityFindingRecord) (map[string]any, error) {
	connector, err := h.Store.GetConnector(ctx, finding.ConnectorID)
	if err != nil {
		return nil, err
	}

	var docID, resolvedAt, ruleID any
	if finding.DocID != "" {
		docID = finding.DocID
	}
	if finding.ResolvedAt != "" {
		resolvedAt = finding.ResolvedAt
	}
	if finding.RuleID != "" {
		ruleID = finding.RuleID
	}

	return map[string]any{
		"id":              finding.ID,
		"connectorId":     finding.ConnectorID,
		"connectorName":   connector.Name,
		"docId":           docID,
		"ruleId":          ruleID,
		"checkType":       finding.CheckType,
		"severity":        finding.Severity,
		"title":           finding.Title,
		"description":     finding.Description,
		"remediationLink": finding.RemediationLink,
		"status":          finding.Status,
		"detectedCount":   finding.DetectedCount,
		"firstDetectedAt": finding.FirstDetectedAt,
		"lastSeenAt":      finding.LastSeenAt,
		"resolvedAt":      resolvedAt,
	}, nil
}

// List handles GET /api/findings.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)
	findings, total, err := h.Store.ListQualityFindings(
		r.Context(),
		r.URL.Query().Get("connectorId"),
		r.URL.Query().Get("checkType"),
		r.URL.Query().Get("status"),
		"",
		offset,
		pageSize,
	)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	connectorIDs := make([]string, len(findings))
	for i, f := range findings {
		connectorIDs[i] = f.ConnectorID
	}
	allowed, err := h.Store.FilterConnectorIDsByGrant(r.Context(), auth.UserIDFromContext(r.Context()), connectorIDs, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	isAllowed := make(map[string]bool, len(allowed))
	for _, id := range allowed {
		isAllowed[id] = true
	}

	items := make([]map[string]any, 0, len(findings))
	for _, finding := range findings {
		if !isAllowed[finding.ConnectorID] {
			continue
		}
		item, err := h.response(r.Context(), finding)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		items = append(items, item)
	}
	httputil.WritePaginated(w, items, page, pageSize, total)
}

// Get handles GET /api/findings/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	finding, err := h.Store.GetQualityFinding(r.Context(), r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Quality finding not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.viewerOrNotFound(w, r, finding.ConnectorID) {
		return
	}
	response, err := h.response(r.Context(), *finding)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, response)
}

// Resolve handles POST /api/findings/{id}/resolve.
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.Store.GetQualityFinding(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Quality finding not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.operatorOrForbidden(w, r, existing.ConnectorID) {
		return
	}
	if err := h.Store.UpdateQualityFindingStatus(r.Context(), id, "resolved"); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Quality finding not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "finding.resolve", "finding", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "finding.resolve", "error", err)
	}
	finding, err := h.Store.GetQualityFinding(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	response, err := h.response(r.Context(), *finding)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, response)
}

// viewerOrNotFound 404s (not 403, to avoid confirming the resource's
// existence) when the caller lacks at least a viewer grant on connectorID.
func (h *Handler) viewerOrNotFound(w http.ResponseWriter, r *http.Request, connectorID string) bool {
	ok, err := h.Store.UserHasConnectorRole(r.Context(), auth.UserIDFromContext(r.Context()), connectorID, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return false
	}
	if !ok {
		httputil.Error(w, http.StatusNotFound, "not_found", "Quality finding not found")
		return false
	}
	return true
}

// operatorOrForbidden 403s when the caller lacks at least an operator grant
// on connectorID.
func (h *Handler) operatorOrForbidden(w http.ResponseWriter, r *http.Request, connectorID string) bool {
	ok, err := h.Store.UserHasConnectorRole(r.Context(), auth.UserIDFromContext(r.Context()), connectorID, "operator")
	if err != nil {
		httputil.Errorf(w, err)
		return false
	}
	if !ok {
		httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
		return false
	}
	return true
}
