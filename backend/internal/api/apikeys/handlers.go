// Package apikeys provides self-service opaque API-key management.
package apikeys

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for API-key endpoints.
type Handler struct {
	Store *store.Store
}

// NewHandler creates a new API-key handler.
func NewHandler(s *store.Store) *Handler {
	return &Handler{Store: s}
}

// Create handles POST /api/auth/api-keys. The raw token is returned only in
// this response and is never persisted or included in list responses.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// A restricted key must not mint a broader one.
	if auth.RejectRestrictedAPIKey(w, r) {
		return
	}
	req, ok := httputil.DecodeJSON[struct {
		Name         string   `json:"name"`
		ExpiresAt    string   `json:"expiresAt"`
		Scope        string   `json:"scope"`
		ConnectorIDs []string `json:"connectorIds"`
	}](w, r)
	if !ok {
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "name is required", []httputil.FieldError{{Field: "name", Msg: "is required"}})
		return
	}
	if req.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil || !expiresAt.After(time.Now().UTC()) {
			httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "expiresAt must be a future RFC3339 timestamp", []httputil.FieldError{{Field: "expiresAt", Msg: "must be a future RFC3339 timestamp"}})
			return
		}
	}

	if req.Scope == "" {
		req.Scope = auth.APIKeyScopeFull
	}
	if req.Scope != auth.APIKeyScopeFull && req.Scope != auth.APIKeyScopeRead {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "scope must be full or read", []httputil.FieldError{{Field: "scope", Msg: "must be full or read"}})
		return
	}
	connectorIDs, ok := h.validConnectorIDs(w, r, req.ConnectorIDs)
	if !ok {
		return
	}

	rawToken, err := newToken()
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	// Role is a display-only snapshot of the creator's flat instance-admin
	// role at creation time; auth actually derives a key's access from the
	// owning user's current role at lookup time (store.LookupAPIKey), not
	// from this column; Scope and ConnectorIDs (#278) only narrow that
	// access, never widen it. Stored as
	// "operator"/"viewer" — the api_keys.role CHECK constraint predates this
	// migration and wasn't touched, only users.role was renamed/reworked.
	role := "viewer"
	if auth.InstanceAdminFromContext(r.Context()) {
		role = "operator"
	}
	key := &store.APIKey{
		UserID:       auth.UserIDFromContext(r.Context()),
		Name:         req.Name,
		TokenHash:    store.HashToken(rawToken),
		Role:         role,
		ExpiresAt:    req.ExpiresAt,
		Scope:        req.Scope,
		ConnectorIDs: connectorIDs,
	}
	if err := h.Store.CreateAPIKey(r.Context(), key); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "auth.api_key.create", "api_key", key.ID, map[string]any{
		"name":         key.Name,
		"scope":        key.Scope,
		"connectorIds": key.ConnectorIDs,
	}); err != nil {
		slog.Error("failed to record audit", "action", "auth.api_key.create", "error", err)
	}

	out := sanitize(*key)
	out["token"] = rawToken
	httputil.JSON(w, http.StatusCreated, out)
}

// maxKeyConnectors bounds a key's connector allow-list.
const maxKeyConnectors = 100

// validConnectorIDs dedupes a requested connector allow-list and checks the
// caller holds at least viewer on each connector, so a key can only be
// restricted to connectors its owner can already reach. Unknown IDs and IDs
// without a grant get the same error, so the response doesn't reveal which
// connectors exist.
func (h *Handler) validConnectorIDs(w http.ResponseWriter, r *http.Request, ids []string) ([]string, bool) {
	out := make([]string, 0, len(ids))
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	if len(out) > maxKeyConnectors {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", fmt.Sprintf("at most %d connectors", maxKeyConnectors), []httputil.FieldError{{Field: "connectorIds", Msg: fmt.Sprintf("must list at most %d connectors", maxKeyConnectors)}})
		return nil, false
	}
	userID := auth.UserIDFromContext(r.Context())
	allowed, err := h.Store.FilterConnectorIDsByGrant(r.Context(), userID, out, "viewer")
	if err != nil {
		httputil.Errorf(w, err)
		return nil, false
	}
	if len(allowed) != len(out) {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "connectorIds contains a connector you have no access to", []httputil.FieldError{{Field: "connectorIds", Msg: "must only list connectors you have access to"}})
		return nil, false
	}
	return out, true
}

// List handles GET /api/auth/api-keys. It returns only keys owned by the
// calling user and never exposes a token or hash.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	keys, err := h.Store.ListAPIKeysForUser(r.Context(), auth.UserIDFromContext(r.Context()))
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	out := make([]map[string]any, len(keys))
	for i := range keys {
		out[i] = sanitize(keys[i])
	}
	httputil.JSON(w, http.StatusOK, out)
}

// Revoke handles DELETE /api/auth/api-keys/{id}.
func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "API key ID is required")
		return
	}
	key, err := h.Store.GetAPIKeyByID(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "API key not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if key.UserID != auth.UserIDFromContext(r.Context()) {
		httputil.Error(w, http.StatusForbidden, "forbidden", "Cannot revoke another user's API key")
		return
	}
	if err := h.Store.RevokeAPIKey(r.Context(), id); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "auth.api_key.revoke", "api_key", id, nil); err != nil {
		slog.Error("failed to record audit", "action", "auth.api_key.revoke", "error", err)
	}
	httputil.NoContent(w)
}

func newToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate API key: %w", err)
	}
	return "wlz_" + base64.RawURLEncoding.EncodeToString(buf), nil
}

func sanitize(key store.APIKey) map[string]any {
	return map[string]any{
		"id":           key.ID,
		"name":         key.Name,
		"role":         key.Role,
		"createdAt":    key.CreatedAt,
		"expiresAt":    key.ExpiresAt,
		"lastUsedAt":   key.LastUsedAt,
		"revokedAt":    key.RevokedAt,
		"scope":        key.Scope,
		"connectorIds": key.ConnectorIDs,
	}
}
