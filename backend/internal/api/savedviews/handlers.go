// Package savedviews provides API handlers for per-user saved operational
// views (named, reusable filter sets on the services/changes/alerts list
// surfaces). See docs/SAVED_VIEWS.md.
package savedviews

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for saved-view endpoints.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a new saved-views handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

// surfaces are the list pages that support saved views. Kept in sync with
// the `surface` enum in docs/openapi.yaml.
var surfaces = map[string]bool{
	"services": true,
	"changes":  true,
	"alerts":   true,
}

// List handles GET /api/saved-views?surface=X. Scoped to the calling user.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	surface := r.URL.Query().Get("surface")
	if !surfaces[surface] {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "surface must be one of: services, changes, alerts", []httputil.FieldError{{Field: "surface", Msg: "must be one of: services, changes, alerts"}})
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	views, err := h.Store.ListSavedViews(r.Context(), userID, surface)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, views)
}

// Create handles POST /api/saved-views.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		Surface string          `json:"surface"`
		Name    string          `json:"name"`
		Filters json.RawMessage `json:"filters"`
	}](w, r)
	if !ok {
		return
	}
	if !surfaces[req.Surface] {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "surface must be one of: services, changes, alerts", []httputil.FieldError{{Field: "surface", Msg: "must be one of: services, changes, alerts"}})
		return
	}
	if req.Name == "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "name is required", []httputil.FieldError{{Field: "name", Msg: "is required"}})
		return
	}

	filters := string(req.Filters)
	if filters == "" || filters == "null" {
		filters = "{}"
	}

	v := &store.SavedView{
		UserID:  auth.UserIDFromContext(r.Context()),
		Surface: req.Surface,
		Name:    req.Name,
		Filters: filters,
	}
	if err := h.Store.CreateSavedView(r.Context(), v); err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, v)
}

// Delete handles DELETE /api/saved-views/{id}. A user may only delete their
// own saved views (same ownership-boundary pattern as DELETE /me/sessions/{id}).
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := auth.UserIDFromContext(r.Context())

	v, err := h.Store.GetSavedView(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Saved view not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if v.UserID != userID {
		httputil.Error(w, http.StatusForbidden, "forbidden", "Cannot delete another user's saved view")
		return
	}

	if err := h.Store.DeleteSavedView(r.Context(), id); err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.NoContent(w)
}
