package docs

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// requireDocOperator 403s unless the caller may mutate a doc scoped to
// connectorID. Lab-wide docs (connectorID == "", e.g. the Lab Topology doc)
// have no connector to check against, so they're gated on instance-admin
// instead of a per-connector grant.
func (h *Handler) requireDocOperator(w http.ResponseWriter, r *http.Request, connectorID string) bool {
	if connectorID == "" {
		if !auth.InstanceAdminFromContext(r.Context()) {
			httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
			return false
		}
		return true
	}
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

// requireDocViewer 404s (not 403, to avoid confirming existence) unless the
// caller may view a doc scoped to connectorID. Lab-wide docs are visible to
// any authenticated user, same as before this migration.
func (h *Handler) requireDocViewer(w http.ResponseWriter, r *http.Request, connectorID string) bool {
	if connectorID == "" {
		return true
	}
	ok, err := h.Store.UserHasConnectorRole(r.Context(), auth.UserIDFromContext(r.Context()), connectorID, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return false
	}
	if !ok {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return false
	}
	return true
}
