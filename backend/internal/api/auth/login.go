package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

const (
	maxFailedLoginAttempts = 5
	loginLockoutDuration   = 15 * time.Minute
)

// dummyPasswordHash is a bcrypt hash (generated once at the package's normal
// cost factor) verified against on every unknown-user or otherwise-rejected
// login. It exists purely so the unknown-user path costs the same one bcrypt
// comparison as the known-user path, closing the timing side channel that
// let an attacker enumerate usernames by response time.
var dummyPasswordHash = mustHashDummyPassword()

func mustHashDummyPassword() string {
	hash, err := auth.HashPassword("dummy-password-not-a-real-account")
	if err != nil {
		panic("auth: failed to precompute dummy password hash: " + err.Error())
	}
	return hash
}

// Login handles POST /api/auth/login.
// Validates local credentials and returns a JWT token pair.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}
	if req.Username == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Username and password are required")
		return
	}

	user, err := h.Store.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		// Unknown user: still pay for one bcrypt comparison so this path is
		// not distinguishable by timing from a known user with a wrong password.
		_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(req.Password))
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
		return
	}

	if user.LockedUntil != "" {
		if lockedUntil, parseErr := time.Parse(time.RFC3339, user.LockedUntil); parseErr == nil && time.Now().UTC().Before(lockedUntil) {
			_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(req.Password))
			httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
			return
		}
	}

	// Verify the password before checking account status, so a disabled or
	// OIDC-only account still costs one bcrypt comparison like any other
	// rejected login and is not distinguishable by timing.
	verifyErr := auth.VerifyPassword(user.PasswordHash, req.Password)

	if user.Disabled || user.AuthSource != "local" || verifyErr != nil {
		if verifyErr != nil && !user.Disabled && user.AuthSource == "local" {
			if locked, lockErr := h.Store.RegisterFailedLogin(r.Context(), user.ID, maxFailedLoginAttempts, loginLockoutDuration); lockErr != nil {
				h.logError("failed to register failed login", lockErr)
			} else if locked {
				if auditErr := h.Store.RecordAuditFromContext(r.Context(), "auth.account_locked", "user", user.ID, map[string]any{"username": user.Username}); auditErr != nil {
					h.logError("failed to record audit", auditErr)
				}
			}
		}
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Invalid credentials")
		return
	}

	if err := h.Store.ClearFailedLogins(r.Context(), user.ID); err != nil {
		h.logError("failed to clear failed logins", err)
	}

	// Issue token pair
	pair, err := h.JWT.IssuePair(user.ID, user.Role)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Create session
	session := &store.Session{
		UserID:    user.ID,
		TokenHash: store.HashToken(pair.RefreshToken),
		UserAgent: r.UserAgent(),
		IP:        httputil.ClientIP(r, h.Config.Server.TrustedProxies),
	}
	if err := h.Store.CreateSession(r.Context(), session); err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Set refresh token as HTTP-only cookie
	setRefreshCookie(w, r, h.Config.Server.TrustedProxies, pair.RefreshToken, h.Config.Auth.RefreshTokenTTLDuration())

	httputil.JSON(w, http.StatusOK, map[string]any{
		"accessToken": pair.AccessToken,
		"expiresIn":   pair.ExpiresIn,
		"user":        sanitizeUser(user),
	})
}

// Me handles GET /api/me.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	user, err := h.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "User not found")
		return
	}
	httputil.JSON(w, http.StatusOK, sanitizeUser(user))
}

// UpdateMe handles PATCH /api/me.
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	var req struct {
		DisplayName *string `json:"displayName"`
		Email       *string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	updates := make(map[string]any)
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if len(updates) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "No fields to update")
		return
	}

	if err := h.Store.UpdateUser(r.Context(), userID, updates); err != nil {
		httputil.Errorf(w, err)
		return
	}

	user, _ := h.Store.GetUserByID(r.Context(), userID)
	httputil.JSON(w, http.StatusOK, sanitizeUser(user))
}

// ChangePassword handles POST /api/me/password.
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())

	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "currentPassword and newPassword are required")
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	if err := auth.VerifyPassword(user.PasswordHash, req.CurrentPassword); err != nil {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized", "Current password is incorrect")
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid password")
		return
	}

	if err := h.Store.UpdateUser(r.Context(), userID, map[string]any{
		"password_hash": hash,
	}); err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Evict every session (and thus every outstanding refresh token) issued
	// under the old password, then re-issue a fresh pair so the caller isn't
	// logged out of the request they just made.
	if err := h.Store.DeleteUserSessions(r.Context(), userID); err != nil {
		httputil.Errorf(w, err)
		return
	}

	pair, err := h.JWT.IssuePair(userID, auth.RoleFromContext(r.Context()))
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if err := h.Store.CreateSession(r.Context(), &store.Session{
		UserID:    userID,
		TokenHash: store.HashToken(pair.RefreshToken),
		UserAgent: r.UserAgent(),
		IP:        httputil.ClientIP(r, h.Config.Server.TrustedProxies),
	}); err != nil {
		httputil.Errorf(w, err)
		return
	}
	setRefreshCookie(w, r, h.Config.Server.TrustedProxies, pair.RefreshToken, h.Config.Auth.RefreshTokenTTLDuration())

	httputil.JSON(w, http.StatusOK, map[string]any{
		"accessToken": pair.AccessToken,
		"expiresIn":   pair.ExpiresIn,
	})
}
