// Package dashboard provides the dashboard overview API handler.
package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ttlcache"
)

const (
	overviewDefaultDays = 7
	overviewCacheTTL    = 5 * time.Second
)

// Handler holds dependencies for dashboard endpoints.
type Handler struct {
	Store *store.Store
	cache *ttlcache.Cache[map[string]any]
}

// NewHandler creates a new dashboard handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s, cache: ttlcache.New[map[string]any](overviewCacheTTL)}
}

// Overview handles GET /api/dashboard/overview.
// Returns aggregated dashboard data, cached briefly per ?days= window.
func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	days := r.URL.Query().Get("days")
	if cached, ok := h.cache.Get(days); ok {
		httputil.JSON(w, http.StatusOK, cached)
		return
	}
	since := store.SinceFromDays(days, overviewDefaultDays)

	statusCounts, _ := h.Store.CountConnectorsByStatus(ctx)
	pendingAlerts, _ := h.Store.CountAlertsPending(ctx)
	latestChanges, _ := h.Store.GetLatestChangeSummaries(ctx, 5, since)
	lastSync, _ := h.Store.GetLastSyncTimestamp(ctx)

	recentChanges := make([]map[string]any, len(latestChanges))
	for i, c := range latestChanges {
		recentChanges[i] = map[string]any{
			"id":            c.ID,
			"serviceId":     c.ServiceID,
			"serviceName":   c.ServiceName,
			"changeType":    c.ChangeType,
			"severity":      c.Severity,
			"summary":       c.Summary,
			"willTriggerAi": false,
			"detectedAt":    c.DetectedAt,
		}
	}

	resp := map[string]any{
		"statusCounts":  statusCounts,
		"pendingAlerts": pendingAlerts,
		"recentChanges": recentChanges,
		"lastSyncAt":    lastSync,
	}
	h.cache.Set(days, resp)
	httputil.JSON(w, http.StatusOK, resp)
}

// GetLayout handles GET /api/dashboard/layout.
func (h *Handler) GetLayout(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	// Query dashboard_layouts table
	var widgets string
	err := h.Store.DB().QueryRowContext(r.Context(),
		`SELECT widgets FROM dashboard_layouts WHERE user_id = ?`, userID,
	).Scan(&widgets)

	if err != nil {
		// No saved layout yet — a new user starts on the administrator default.
		widgets, err = h.getAdminDefaultWidgets(r.Context())
		if err != nil {
			widgets = "[]"
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"widgets": json.RawMessage(widgets),
	})
}

// getAdminDefaultWidgets returns the raw JSON widgets column of the single
// dashboard_admin_default row (id=1), seeded by migration 000010.
func (h *Handler) getAdminDefaultWidgets(ctx context.Context) (string, error) {
	var widgets string
	err := h.Store.DB().QueryRowContext(ctx,
		`SELECT widgets FROM dashboard_admin_default WHERE id = 1`,
	).Scan(&widgets)
	if err != nil {
		return "", err
	}
	return widgets, nil
}

// GetAdminDefault handles GET /api/dashboard/layout/admin-default.
func (h *Handler) GetAdminDefault(w http.ResponseWriter, r *http.Request) {
	widgets, err := h.getAdminDefaultWidgets(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"widgets": json.RawMessage(widgets),
	})
}

// PutAdminDefault handles PUT /api/dashboard/layout/admin-default.
func (h *Handler) PutAdminDefault(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		Widgets json.RawMessage `json:"widgets"`
	}](w, r)
	if !ok {
		return
	}

	widgets := string(req.Widgets)
	if widgets == "" || widgets == "null" {
		widgets = "[]"
	}

	_, err := h.Store.DB().ExecContext(r.Context(),
		`UPDATE dashboard_admin_default SET widgets = ? WHERE id = 1`, widgets,
	)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"widgets": json.RawMessage(widgets),
	})
}

// ResetLayout handles POST /api/dashboard/layout/reset.
// Overwrites the caller's saved layout with the current administrator default.
func (h *Handler) ResetLayout(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	widgets, err := h.getAdminDefaultWidgets(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.saveLayoutWidgets(r.Context(), userID, widgets); err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"widgets": json.RawMessage(widgets),
	})
}

// SaveLayout handles PUT /api/dashboard/layout.
func (h *Handler) SaveLayout(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	req, ok := httputil.DecodeJSON[struct {
		Widgets json.RawMessage `json:"widgets"`
	}](w, r)
	if !ok {
		return
	}

	widgets := string(req.Widgets)
	if widgets == "" || widgets == "null" {
		widgets = "[]"
	}

	if err := h.saveLayoutWidgets(r.Context(), userID, widgets); err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"widgets": json.RawMessage(widgets),
	})
}

// saveLayoutWidgets upserts a user's layout row.
func (h *Handler) saveLayoutWidgets(ctx context.Context, userID, widgets string) error {
	_, err := h.Store.DB().ExecContext(ctx, `
		INSERT INTO dashboard_layouts (user_id, widgets) VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET widgets = ?
	`, userID, widgets, widgets)
	if err != nil {
		// SQLite doesn't support ON CONFLICT — use INSERT OR REPLACE
		_, err = h.Store.DB().ExecContext(ctx, `
			INSERT OR REPLACE INTO dashboard_layouts (user_id, widgets) VALUES (?, ?)
		`, userID, widgets)
	}
	return err
}
