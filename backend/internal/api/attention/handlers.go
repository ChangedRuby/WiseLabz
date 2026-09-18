// Package attention provides the attention queue API handler (aggregated alerts + findings).
package attention

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for attention queue endpoints.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a new attention handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

// List handles GET /api/attention.
// Returns merged, sorted list of open alerts and findings.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)

	// days is optional; unset keeps today's back-compat behavior of showing everything.
	since := store.SinceFromDays(r.URL.Query().Get("days"), 0)

	items, total, err := h.Store.MergedAttentionItems(r.Context(), auth.UserIDFromContext(r.Context()), since, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"data":     items,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}
