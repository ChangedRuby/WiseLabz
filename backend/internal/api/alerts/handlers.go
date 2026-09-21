// Package alerts provides API handlers for alert lifecycle management.
package alerts

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for alert endpoints.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a new alert handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

// alertResponse builds the spec's Alert response for an alert record.
func (h *Handler) alertResponse(ctx context.Context, id string) (map[string]any, error) {
	a, err := h.Store.GetAlert(ctx, id)
	if err != nil {
		return nil, err
	}

	serviceName := ""
	if conn, err := h.Store.GetConnector(ctx, a.ServiceID); err == nil {
		serviceName = conn.Name
	}

	var changeID, snoozedUntil any
	if a.ChangeID != "" {
		changeID = a.ChangeID
	}
	if a.SnoozedUntil != "" {
		snoozedUntil = a.SnoozedUntil
	}

	return map[string]any{
		"id":           a.ID,
		"changeId":     changeID,
		"serviceId":    a.ServiceID,
		"serviceName":  serviceName,
		"severity":     a.Severity,
		"title":        a.Title,
		"description":  a.Description,
		"status":       a.Status,
		"createdAt":    a.CreatedAt,
		"snoozedUntil": snoozedUntil,
	}, nil
}

// List handles GET /api/alerts. Default deny: filters to connectors the
// caller has at least a viewer grant on.
// ponytail: filters after the DB page is fetched, so a page can come back
// shorter than pageSize for a caller with few grants; push the grant filter
// into the SQL query if that skew matters at scale.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)
	serviceID := r.URL.Query().Get("serviceId")
	severity := r.URL.Query().Get("severity")
	status := r.URL.Query().Get("status")

	// days is optional here (unlike the dashboard overview): /api/alerts backs
	// the full Alerts page and must keep showing everything when unset.
	since := store.SinceFromDays(r.URL.Query().Get("days"), 0)

	alerts, total, err := h.Store.ListAlerts(r.Context(), serviceID, severity, status, since, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	alerts, err = h.filterByGrant(r.Context(), alerts)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.WritePaginated(w, alerts, page, pageSize, total)
}

// filterByGrant keeps only alerts whose connector the caller has at least a
// viewer grant on.
func (h *Handler) filterByGrant(ctx context.Context, alerts []store.AlertRecord) ([]store.AlertRecord, error) {
	serviceIDs := make([]string, len(alerts))
	for i, a := range alerts {
		serviceIDs[i] = a.ServiceID
	}
	allowed, err := h.Store.FilterConnectorIDsByGrant(ctx, auth.UserIDFromContext(ctx), serviceIDs, "viewer")
	if err != nil {
		return nil, err
	}
	isAllowed := make(map[string]bool, len(allowed))
	for _, id := range allowed {
		isAllowed[id] = true
	}
	out := make([]store.AlertRecord, 0, len(alerts))
	for _, a := range alerts {
		if isAllowed[a.ServiceID] {
			out = append(out, a)
		}
	}
	return out, nil
}

// Get handles GET /api/alerts/{id}. Default deny: 404s if the caller has no
// grant on the alert's connector.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.Store.GetAlert(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.viewerOrNotFound(w, r, a.ServiceID) {
		return
	}
	httputil.JSON(w, http.StatusOK, a)
}

// viewerOrNotFound 404s (not 403, to avoid confirming the resource's
// existence to a caller with no visibility into its connector) when the
// caller lacks at least a viewer grant on connectorID. Returns true if the
// caller may proceed.
func (h *Handler) viewerOrNotFound(w http.ResponseWriter, r *http.Request, connectorID string) bool {
	ok, err := h.Store.UserHasConnectorRole(r.Context(), auth.UserIDFromContext(r.Context()), connectorID, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return false
	}
	if !ok {
		httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
		return false
	}
	return true
}

// operatorOrForbidden 403s when the caller lacks at least an operator grant
// on connectorID. Returns true if the caller may proceed.
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

// Resolve handles POST /api/alerts/{id}/resolve.
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.Store.GetAlert(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.operatorOrForbidden(w, r, a.ServiceID) {
		return
	}
	if err := h.Store.UpdateAlertStatus(r.Context(), id, "resolved", ""); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "alert.resolve", "alert", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "alert.resolve", "error", err)
	}
	resp, err := h.alertResponse(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, resp)
}

// Dismiss handles POST /api/alerts/{id}/dismiss.
func (h *Handler) Dismiss(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.Store.GetAlert(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.operatorOrForbidden(w, r, a.ServiceID) {
		return
	}
	if err := h.Store.UpdateAlertStatus(r.Context(), id, "dismissed", ""); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "alert.dismiss", "alert", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "alert.dismiss", "error", err)
	}
	resp, err := h.alertResponse(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, resp)
}

// Snooze handles POST /api/alerts/{id}/snooze.
func (h *Handler) Snooze(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, ok := httputil.DecodeJSON[struct {
		Until string `json:"until"`
	}](w, r)
	if !ok {
		return
	}
	if req.Until == "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "until is required", []httputil.FieldError{{Field: "until", Msg: "is required"}})
		return
	}
	if _, err := time.Parse(time.RFC3339, req.Until); err != nil {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "until must be an RFC3339 timestamp", []httputil.FieldError{{Field: "until", Msg: "must be an RFC3339 timestamp"}})
		return
	}

	a, err := h.Store.GetAlert(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.operatorOrForbidden(w, r, a.ServiceID) {
		return
	}

	if err := h.Store.UpdateAlertStatus(r.Context(), id, "snoozed", req.Until); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Alert not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "alert.snooze", "alert", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "alert.snooze", "error", err)
	}
	resp, err := h.alertResponse(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, resp)
}

// bulkSnoozeRequest is the body of POST /api/alerts/bulk-snooze.
type bulkSnoozeRequest struct {
	IDs   []string `json:"ids"`
	Until string   `json:"until"`
}

// bulkSnoozeItemResult is the per-item outcome in the bulk-snooze response.
type bulkSnoozeItemResult struct {
	ID     string `json:"id"`
	Status string `json:"status"` // "success" | "error"
	Reason string `json:"reason,omitempty"`
}

// BulkSnooze handles POST /api/alerts/bulk-snooze. It snoozes an explicit,
// caller-supplied list of alert IDs until a common timestamp in one request.
// One bad ID never aborts the batch: every item gets its own success/error
// outcome, and one audit record is written per successfully-snoozed item.
func (h *Handler) BulkSnooze(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[bulkSnoozeRequest](w, r)
	if !ok {
		return
	}
	if req.Until == "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "until is required", []httputil.FieldError{{Field: "until", Msg: "is required"}})
		return
	}
	if _, err := time.Parse(time.RFC3339, req.Until); err != nil {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "until must be an RFC3339 timestamp", []httputil.FieldError{{Field: "until", Msg: "must be an RFC3339 timestamp"}})
		return
	}
	if len(req.IDs) == 0 {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "ids must be a non-empty array", []httputil.FieldError{{Field: "ids", Msg: "must be a non-empty array"}})
		return
	}
	if len(req.IDs) > httputil.MaxBulkIDs {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "ids must contain at most 500 items", []httputil.FieldError{{Field: "ids", Msg: "must contain at most 500 items"}})
		return
	}

	// ponytail: resolves each alert individually to authorize per-connector
	// (max 500 items per request, same cap as the other bulk endpoints);
	// batch this into one query if bulk-snooze throughput becomes a hot path.
	results := make([]bulkSnoozeItemResult, 0, len(req.IDs))
	allowedIDs := make([]string, 0, len(req.IDs))
	userID := auth.UserIDFromContext(r.Context())
	for _, id := range req.IDs {
		a, err := h.Store.GetAlert(r.Context(), id)
		if errors.Is(err, store.ErrNotFound) {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "not_found"})
			continue
		}
		if err != nil {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "internal_error"})
			continue
		}
		ok, err := h.Store.UserHasConnectorRole(r.Context(), userID, a.ServiceID, "operator")
		if err != nil {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "internal_error"})
			continue
		}
		if !ok {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "forbidden"})
			continue
		}
		allowedIDs = append(allowedIDs, id)
	}

	found, err := h.Store.UpdateAlertStatuses(r.Context(), allowedIDs, "snoozed", req.Until)
	if err != nil {
		for _, id := range allowedIDs {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "internal_error"})
		}
		httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
		return
	}
	auditRecords := make([]store.AuditRecord, 0, len(found))
	for _, id := range allowedIDs {
		if !found[id] {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "not_found"})
			continue
		}
		auditRecords = append(auditRecords, store.AuditRecord{TargetID: id})
		results = append(results, bulkSnoozeItemResult{ID: id, Status: "success"})
	}
	if err := h.Store.RecordAuditBatchFromContext(r.Context(), "alert.bulk_snooze", "alert", auditRecords); err != nil {
		slog.Error("failed to record audit", "action", "alert.bulk_snooze", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
}
