package connectors

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// OpenMaintenanceWindow handles POST /api/connectors/{id}/maintenance-window.
// Opens a time-boxed maintenance window on a connector: while active, the
// scheduler skips it and drift alerting/change-record creation from any sync
// (scheduled or manual) is suppressed, though snapshots keep being saved.
// Operator-only; unlike the other mutating connector actions, this one is
// reversible and time-boxed so it doesn't require step-up elevation.
func (h *Handler) OpenMaintenanceWindow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	req, ok := httputil.DecodeJSON[struct {
		DurationMinutes int `json:"durationMinutes"`
	}](w, r)
	if !ok {
		return
	}
	if req.DurationMinutes <= 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "durationMinutes must be positive")
		return
	}

	now := time.Now().UTC()
	m := &store.MaintenanceWindowRecord{
		ConnectorID: id,
		StartsAt:    now.Format(time.RFC3339),
		EndsAt:      now.Add(time.Duration(req.DurationMinutes) * time.Minute).Format(time.RFC3339),
		CreatedBy:   auth.UserIDFromContext(r.Context()),
	}
	if err := h.Store.CreateMaintenanceWindow(r.Context(), m); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.maintenanceWindow.open", "connector", id, map[string]any{
		"durationMinutes": req.DurationMinutes, "endsAt": m.EndsAt,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.maintenanceWindow.open", "error", err)
	}

	httputil.JSON(w, http.StatusCreated, m)
}

// CloseMaintenanceWindow handles DELETE /api/connectors/{id}/maintenance-window.
// Ends the connector's active maintenance window early. A no-op (200) when
// there is none.
func (h *Handler) CloseMaintenanceWindow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	active, err := h.Store.GetActiveMaintenanceWindow(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if active == nil {
		httputil.JSON(w, http.StatusOK, map[string]any{"closed": false})
		return
	}

	if err := h.Store.CloseMaintenanceWindow(r.Context(), active.ID); err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.maintenanceWindow.close", "connector", id, map[string]any{
		"windowId": active.ID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.maintenanceWindow.close", "error", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"closed": true})
}

// GetMaintenanceWindow handles GET /api/connectors/{id}/maintenance-window.
// Returns the connector's active maintenance window, or null if none.
func (h *Handler) GetMaintenanceWindow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	m, err := h.Store.GetActiveMaintenanceWindow(r.Context(), id)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, m)
}

// ListActiveMaintenance handles GET /api/connectors/maintenance-windows.
// Returns every currently active maintenance window, across all connectors —
// used by the ServicesPage badge to avoid an N+1 fetch per row.
func (h *Handler) ListActiveMaintenance(w http.ResponseWriter, r *http.Request) {
	windows, err := h.Store.ListActiveMaintenanceWindows(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	ids := make([]string, len(windows))
	for i, m := range windows {
		ids[i] = m.ConnectorID
	}
	allowed, err := h.Store.FilterConnectorIDsByGrant(r.Context(), auth.UserIDFromContext(r.Context()), ids, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	isAllowed := make(map[string]bool, len(allowed))
	for _, id := range allowed {
		isAllowed[id] = true
	}
	filtered := make([]store.MaintenanceWindowRecord, 0, len(windows))
	for _, m := range windows {
		if isAllowed[m.ConnectorID] {
			filtered = append(filtered, m)
		}
	}
	httputil.JSON(w, http.StatusOK, filtered)
}
