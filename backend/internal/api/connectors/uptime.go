package connectors

import (
	"errors"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// uptimeWindows are the fixed lookback windows reported by Uptime, keyed by
// the label used in the response.
var uptimeWindows = []struct {
	label string
	d     time.Duration
}{
	{"24h", 24 * time.Hour},
	{"7d", 7 * 24 * time.Hour},
	{"30d", 30 * 24 * time.Hour},
}

// Uptime handles GET /api/connectors/{id}/uptime.
// Returns availability % and MTTR computed from the health_checks time
// series (populated by Health, see connector/health.go) over the 24h, 7d,
// and 30d windows.
func (h *Handler) Uptime(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := h.Store.GetConnector(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	now := time.Now().UTC()
	windows := make(map[string]any, len(uptimeWindows))
	for _, win := range uptimeWindows {
		stats, err := h.Store.GetConnectorUptime(r.Context(), id, now.Add(-win.d), now)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		windows[win.label] = map[string]any{
			"windowStart":     stats.WindowStart,
			"windowEnd":       stats.WindowEnd,
			"checkCount":      stats.CheckCount,
			"availabilityPct": stats.AvailabilityPct,
			"mttrSeconds":     stats.MTTRSeconds,
			"outageCount":     stats.OutageCount,
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"connectorId": id,
		"windows":     windows,
	})
}
