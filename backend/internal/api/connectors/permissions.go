package connectors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// ListPermissions handles GET /api/connectors/{id}/permissions. Instance-
// admin only, mounted separately from the per-connector-role routes (see
// router.go) since granting access is an instance-wide action, not
// something a connector's own operators can do to each other.
func (h *Handler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	connectorID := r.PathValue("id")
	if _, err := h.Store.GetConnector(r.Context(), connectorID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	grants, err := h.Store.ListConnectorGrants(r.Context(), connectorID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, grants)
}

// PutPermission handles PUT /api/connectors/{id}/permissions/{userId}.
func (h *Handler) PutPermission(w http.ResponseWriter, r *http.Request) {
	connectorID := r.PathValue("id")
	userID := r.PathValue("userId")

	req, ok := httputil.DecodeJSON[struct {
		Role string `json:"role"`
	}](w, r)
	if !ok {
		return
	}
	if req.Role != "viewer" && req.Role != "operator" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "role must be 'viewer' or 'operator'")
		return
	}

	if _, err := h.Store.GetConnector(r.Context(), connectorID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Connector not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}
	if _, err := h.Store.GetUserByID(r.Context(), userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	grant, err := h.Store.UpsertConnectorGrant(r.Context(), userID, connectorID, req.Role)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.permission.granted", "connector", connectorID, map[string]any{
		"targetUserId": userID, "role": req.Role,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.permission.granted", "error", err)
	}

	httputil.JSON(w, http.StatusOK, grant)
}

// DeletePermission handles DELETE /api/connectors/{id}/permissions/{userId}.
func (h *Handler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	connectorID := r.PathValue("id")
	userID := r.PathValue("userId")

	if err := h.Store.DeleteConnectorGrant(r.Context(), userID, connectorID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "Grant not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "connector.permission.revoked", "connector", connectorID, map[string]any{
		"targetUserId": userID,
	}); err != nil {
		slog.Error("failed to record audit", "action", "connector.permission.revoked", "error", err)
	}

	httputil.NoContent(w)
}
