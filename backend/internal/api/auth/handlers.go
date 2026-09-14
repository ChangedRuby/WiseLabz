// Package auth provides authentication and user management API handlers.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for auth API endpoints.
type Handler struct {
	Store    *store.Store
	JWT      *auth.Service
	Config   *config.Config
	oidcProv map[string]*auth.OIDCProvider // initialized on first use
}

// NewHandler creates a new auth handler.
func NewHandler(s *store.Store, jwtSvc *auth.Service, cfg *config.Config) *Handler {
	return &Handler{
		Store:  s,
		JWT:    jwtSvc,
		Config: cfg,
	}
}

// --- /api/users endpoints (operator only) ---

// ListUsers handles GET /api/users.
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	_, pageSize, offset := httputil.Paginate(r)
	users, _, err := h.Store.ListUsers(r.Context(), offset, pageSize)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	sanitized := make([]map[string]any, len(users))
	for i, u := range users {
		sanitized[i] = sanitizeUser(&u)
	}

	// Spec: GET /users returns a bare User[] (see openapi.yaml).
	httputil.JSON(w, http.StatusOK, sanitized)
}

// CreateUser handles POST /api/users.
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username                   string `json:"username"`
		Password                   string `json:"password"`
		Email                      string `json:"email"`
		Role                       string `json:"role"`
		CanManageDashboardDefaults bool   `json:"canManageDashboardDefaults"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}
	if req.Username == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "username and password are required")
		return
	}
	if req.Role == "" {
		req.Role = "viewer"
	}
	if req.Role != "viewer" && req.Role != "operator" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "role must be 'viewer' or 'operator'")
		return
	}
	if req.CanManageDashboardDefaults && req.Role != "operator" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "canManageDashboardDefaults requires role 'operator'")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", fmt.Sprintf("Invalid password: %v", err))
		return
	}

	user := &store.User{
		Username:                   req.Username,
		DisplayName:                req.Username,
		Email:                      req.Email,
		Role:                       req.Role,
		AuthSource:                 "local",
		PasswordHash:               hash,
		CanManageDashboardDefaults: req.CanManageDashboardDefaults,
	}
	if err := h.Store.CreateUser(r.Context(), user); err != nil {
		if errors.Is(err, store.ErrConflict) {
			httputil.Error(w, http.StatusConflict, "conflict", "Username already exists")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, sanitizeUser(user))
}

// UpdateUser handles PATCH /api/users/{id}.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "User ID is required")
		return
	}

	var req struct {
		Username                   *string `json:"username"`
		DisplayName                *string `json:"displayName"`
		Email                      *string `json:"email"`
		Role                       *string `json:"role"`
		Disabled                   *bool   `json:"disabled"`
		CanManageDashboardDefaults *bool   `json:"canManageDashboardDefaults"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	updates := make(map[string]any)
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Role != nil {
		if *req.Role != "viewer" && *req.Role != "operator" {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "role must be 'viewer' or 'operator'")
			return
		}
		updates["role"] = *req.Role
	}
	if req.Disabled != nil {
		updates["disabled"] = *req.Disabled
	}
	if req.CanManageDashboardDefaults != nil {
		effectiveRole := req.Role
		if effectiveRole == nil {
			existing, err := h.Store.GetUserByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					httputil.Error(w, http.StatusNotFound, "not_found", "User not found")
					return
				}
				httputil.Errorf(w, err)
				return
			}
			effectiveRole = &existing.Role
		}
		if *req.CanManageDashboardDefaults && *effectiveRole != "operator" {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "canManageDashboardDefaults requires role 'operator'")
			return
		}
		updates["can_manage_dashboard_defaults"] = *req.CanManageDashboardDefaults
	}

	if len(updates) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "No fields to update")
		return
	}

	if err := h.Store.UpdateUser(r.Context(), userID, updates); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		if errors.Is(err, store.ErrConflict) {
			httputil.Error(w, http.StatusConflict, "conflict", "Username already exists")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	// A role demotion or disable must take effect immediately, not just once
	// the caller's still-valid access token expires.
	if _, ok := updates["role"]; ok {
		if err := h.Store.DeleteUserSessions(r.Context(), userID); err != nil {
			h.logError("failed to delete user sessions after role change", err)
		}
	} else if disabled, ok := updates["disabled"].(bool); ok && disabled {
		if err := h.Store.DeleteUserSessions(r.Context(), userID); err != nil {
			h.logError("failed to delete user sessions after disable", err)
		}
	}
	if disabled, ok := updates["disabled"].(bool); ok && disabled {
		if err := h.Store.RevokeAllAPIKeysForUser(r.Context(), userID); err != nil {
			h.logError("failed to revoke api keys after disable", err)
		}
	}

	user, _ := h.Store.GetUserByID(r.Context(), userID)
	httputil.JSON(w, http.StatusOK, sanitizeUser(user))
}

// DeleteUser handles DELETE /api/users/{id}.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "User ID is required")
		return
	}

	// Prevent self-deletion
	if auth.UserIDFromContext(r.Context()) == userID {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Cannot delete your own account")
		return
	}

	if err := h.Store.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	httputil.NoContent(w)
}

// ResetPassword handles POST /api/users/{id}/reset-password.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "User ID is required")
		return
	}

	var req struct {
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}
	if req.NewPassword == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "newPassword is required")
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", fmt.Sprintf("Invalid password: %v", err))
		return
	}

	if err := h.Store.UpdateUser(r.Context(), userID, map[string]any{
		"password_hash": hash,
	}); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		httputil.Errorf(w, err)
		return
	}

	// A forced reset is the standard response to a suspected compromise; it
	// must evict any session an attacker is already holding.
	if err := h.Store.DeleteUserSessions(r.Context(), userID); err != nil {
		h.logError("failed to delete user sessions after reset", err)
	}

	httputil.NoContent(w)
}

// --- Helpers ---

func sanitizeUser(u *store.User) map[string]any {
	return map[string]any{
		"id":                         u.ID,
		"username":                   u.Username,
		"displayName":                u.DisplayName,
		"email":                      u.Email,
		"role":                       u.Role,
		"authSource":                 u.AuthSource,
		"disabled":                   u.Disabled,
		"canManageDashboardDefaults": u.CanManageDashboardDefaults,
		"createdAt":                  u.CreatedAt,
	}
}

func sanitizeSessions(sessions []store.Session, currentHash string) []map[string]any {
	out := make([]map[string]any, len(sessions))
	for i, s := range sessions {
		out[i] = map[string]any{
			"id":         s.ID,
			"userAgent":  s.UserAgent,
			"ip":         s.IP,
			"createdAt":  s.CreatedAt,
			"lastSeenAt": s.LastSeenAt,
			"current":    currentHash != "" && s.TokenHash == currentHash,
		}
	}
	return out
}

func setRefreshCookie(w http.ResponseWriter, r *http.Request, trustedProxies, token string, maxAge time.Duration) {
	secure := httputil.IsSecureRequest(r, trustedProxies)
	// Remove the pre-#83 path-scoped cookie so a browser cannot replay it first.
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/api/auth", MaxAge: -1, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) logError(msg string, err error) {
	slog.Error(msg, "error", err)
}
