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
	req, ok := httputil.DecodeJSON[struct {
		Name      string `json:"name"`
		ExpiresAt string `json:"expiresAt"`
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

	rawToken, err := newToken()
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	// Role is a display-only snapshot of the creator's flat instance-admin
	// role at creation time; auth actually derives a key's access from the
	// owning user's current role at lookup time (store.LookupAPIKey), not
	// from this column, so per-connector-scoped API keys are out of scope
	// here (#240 PR1 keeps API keys flat admin/user). Stored as
	// "operator"/"viewer" — the api_keys.role CHECK constraint predates this
	// migration and wasn't touched, only users.role was renamed/reworked.
	role := "viewer"
	if auth.InstanceAdminFromContext(r.Context()) {
		role = "operator"
	}
	key := &store.APIKey{
		UserID:    auth.UserIDFromContext(r.Context()),
		Name:      req.Name,
		TokenHash: store.HashToken(rawToken),
		Role:      role,
		ExpiresAt: req.ExpiresAt,
	}
	if err := h.Store.CreateAPIKey(r.Context(), key); err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "auth.api_key.create", "api_key", key.ID, map[string]any{
		"name": key.Name,
	}); err != nil {
		slog.Error("failed to record audit", "action", "auth.api_key.create", "error", err)
	}

	httputil.JSON(w, http.StatusCreated, map[string]any{
		"id":         key.ID,
		"name":       key.Name,
		"role":       key.Role,
		"createdAt":  key.CreatedAt,
		"expiresAt":  key.ExpiresAt,
		"lastUsedAt": key.LastUsedAt,
		"revokedAt":  key.RevokedAt,
		"token":      rawToken,
	})
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
		"id":         key.ID,
		"name":       key.Name,
		"role":       key.Role,
		"createdAt":  key.CreatedAt,
		"expiresAt":  key.ExpiresAt,
		"lastUsedAt": key.LastUsedAt,
		"revokedAt":  key.RevokedAt,
	}
}
