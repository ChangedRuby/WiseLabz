// Package attention provides the attention queue API handler (aggregated alerts + findings).
package attention

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ttlcache"
)

// cacheTTL bounds how stale a cached attention page may be.
const cacheTTL = 5 * time.Second

// Handler holds dependencies for attention queue endpoints.
type Handler struct {
	Store *store.Store
	cache *ttlcache.Cache[httputil.DataPaginatedResponse]
}

// NewHandler creates a new attention handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s, cache: ttlcache.New[httputil.DataPaginatedResponse](cacheTTL)}
}

// List handles GET /api/attention.
// Returns merged, sorted list of open alerts and findings.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize, offset := httputil.Paginate(r)

	// days is optional; unset keeps today's back-compat behavior of showing everything.
	since := store.SinceFromDays(r.URL.Query().Get("days"), 0)

	userID := auth.UserIDFromContext(r.Context())
	key := fmt.Sprintf("%s|%s|%d|%d", userID, r.URL.Query().Get("days"), offset, pageSize)
	if cached, ok := h.cache.Get(key); ok {
		httputil.WriteDataPaginated(w, cached.Data, cached.Page, cached.PageSize, cached.Total)
		return
	}

	items, total, err := h.Store.MergedAttentionItems(r.Context(), userID, since, offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	resp := httputil.DataPaginatedResponse{
		Data:     items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	h.cache.Set(key, resp)
	httputil.WriteDataPaginated(w, resp.Data, resp.Page, resp.PageSize, resp.Total)
}
