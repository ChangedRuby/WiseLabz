// Package alerts provides API handlers for alert lifecycle management.
package alerts

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for alert endpoints.
type Handler struct {
	Store *store.Store
}

const maxBulkIDs = 500

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

// List handles GET /api/alerts.
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
	httputil.WritePaginated(w, alerts, page, pageSize, total)
}

// Get handles GET /api/alerts/{id}.
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
	httputil.JSON(w, http.StatusOK, a)
}

// Resolve handles POST /api/alerts/{id}/resolve.
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
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
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "until is required")
		return
	}
	if _, err := time.Parse(time.RFC3339, req.Until); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "until must be an RFC3339 timestamp")
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
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "until is required")
		return
	}
	if _, err := time.Parse(time.RFC3339, req.Until); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "until must be an RFC3339 timestamp")
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

	results := make([]bulkSnoozeItemResult, 0, len(req.IDs))
	found, err := h.Store.UpdateAlertStatuses(r.Context(), req.IDs, "snoozed", req.Until)
	if err != nil {
		for _, id := range req.IDs {
			results = append(results, bulkSnoozeItemResult{ID: id, Status: "error", Reason: "internal_error"})
		}
		httputil.JSON(w, http.StatusOK, map[string]any{"results": results})
		return
	}
	auditRecords := make([]store.AuditRecord, 0, len(found))
	for _, id := range req.IDs {
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
