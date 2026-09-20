package docs

import (
	"errors"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// GetLock handles GET /api/docs/{id}/lock. Returns the current lock, or an
// empty object if the doc is unlocked/expired. Open to any authenticated role.
func (h *Handler) GetLock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	lock, err := h.Store.GetDocLock(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.JSON(w, http.StatusOK, map[string]any{})
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, lock)
}

// AcquireLock handles POST /api/docs/{id}/lock. Acquires or renews the
// advisory editing lock; 409 with the current holder if held by someone
// else. Operator-only, same gate as Save.
func (h *Handler) AcquireLock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := auth.UserIDFromContext(r.Context())

	existing, err := h.Store.GetDoc(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.requireDocOperator(w, r, existing.ServiceID) {
		return
	}

	lock, err := h.Store.AcquireDocLock(r.Context(), id, userID)
	if errors.Is(err, store.ErrLockHeldByOther) {
		httputil.JSON(w, http.StatusConflict, lock)
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if h.WSHub != nil {
		h.WSHub.Broadcast(ws.EventDocLockAcquired, lock)
	}
	httputil.JSON(w, http.StatusOK, lock)
}

// ReleaseLock handles POST /api/docs/{id}/lock/release. Operator-only.
func (h *Handler) ReleaseLock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := auth.UserIDFromContext(r.Context())

	existing, err := h.Store.GetDoc(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.requireDocOperator(w, r, existing.ServiceID) {
		return
	}

	if err := h.Store.ReleaseDocLock(r.Context(), id, userID); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if h.WSHub != nil {
		h.WSHub.Broadcast(ws.EventDocLockReleased, map[string]any{"docId": id, "userId": userID})
	}
	httputil.JSON(w, http.StatusOK, map[string]any{})
}
