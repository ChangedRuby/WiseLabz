package system

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// GetJobs handles GET /api/system/jobs. Operator-only. Returns every
// scheduled job currently registered with the scheduler, combining its live
// next-run time with its persisted health (#384) — see
// scheduler.Runner.ListJobs.
func (h *Handler) GetJobs(w http.ResponseWriter, r *http.Request) {
	if h.Scheduler == nil {
		httputil.JSON(w, http.StatusOK, []any{})
		return
	}
	jobs, err := h.Scheduler.ListJobs(r.Context())
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, jobs)
}
