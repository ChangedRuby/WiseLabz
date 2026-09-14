package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Refresh handles POST /api/auth/refresh.
// Reads the refresh token from an HTTP-only cookie and issues a new access token.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	token := ""
	if err != nil {
		// Also check body for backwards compat
		var req struct {
			RefreshToken string `json:"refreshToken"`
		}
		if json.NewDecoder(r.Body).Decode(&req) == nil && req.RefreshToken != "" {
			token = req.RefreshToken
		}
		if token == "" {
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Missing refresh token")
			return
		}
	} else {
		token = cookie.Value
	}

	claims, err := h.JWT.ValidateRefresh(token)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired refresh token")
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "User not found")
		return
	}
	if user.Disabled {
		httputil.Error(w, http.StatusForbidden, "forbidden", "Account is disabled")
		return
	}

	// Issue new pair (rotate refresh token)
	pair, err := h.JWT.IssuePair(user.ID, user.Role)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.RotateSessionToken(r.Context(), user.ID, store.HashToken(token), store.HashToken(pair.RefreshToken)); err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Refresh token is no longer active")
		return
	}

	setRefreshCookie(w, r, h.Config.Server.TrustedProxies, pair.RefreshToken, h.Config.Auth.RefreshTokenTTLDuration())

	httputil.JSON(w, http.StatusOK, map[string]any{
		"accessToken": pair.AccessToken,
		"expiresIn":   pair.ExpiresIn,
	})
}

// Logout handles POST /api/auth/logout.
// Clears the refresh token cookie and deletes the session.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	cookie, _ := r.Cookie("refresh_token")
	if cookie != nil {
		for _, path := range []string{"/", "/api/auth"} {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     path,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   httputil.IsSecureRequest(r, h.Config.Server.TrustedProxies),
				SameSite: http.SameSiteLaxMode,
			})
		}
	}

	// Delete all sessions for user
	if err := h.Store.DeleteUserSessions(r.Context(), userID); err != nil {
		h.logError("failed to delete user sessions", err)
	}

	httputil.NoContent(w)
}

// Elevate handles POST /api/auth/elevate.
// Requires password re-verification and issues a short-lived elevation token.
func (h *Handler) Elevate(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	var req struct {
		Password string `json:"password"`
		Action   string `json:"action"` // e.g. "connector.delete"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}
	if req.Password == "" || req.Action == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "password and action are required")
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "User not found")
		return
	}

	if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid password")
		return
	}

	token, err := h.JWT.IssueElevation(userID, req.Action)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "auth.elevate", "action", req.Action, nil); err != nil {
		h.logError("failed to record audit", err)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"token":     token.Token,
		"expiresAt": token.ExpiresAt.Format(time.RFC3339),
	})
}

// ListSessions handles GET /api/me/sessions.
func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	sessions, err := h.Store.ListUserSessions(r.Context(), userID)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	var currentHash string
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		currentHash = store.HashToken(cookie.Value)
	}

	// Spec: GET /me/sessions returns a bare Session[] (see openapi.yaml).
	httputil.JSON(w, http.StatusOK, sanitizeSessions(sessions, currentHash))
}

// DeleteSession handles DELETE /api/me/sessions/{id}.
func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	sessionID := r.PathValue("id")
	if sessionID == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Session ID is required")
		return
	}

	session, err := h.Store.GetSession(r.Context(), sessionID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "not_found", "Session not found")
		return
	}
	if session.UserID != userID {
		httputil.Error(w, http.StatusForbidden, "forbidden", "Cannot delete another user's session")
		return
	}

	if err := h.Store.DeleteSession(r.Context(), sessionID); err != nil {
		httputil.Errorf(w, err)
		return
	}

	httputil.NoContent(w)
}
