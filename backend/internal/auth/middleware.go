package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/httputil"
)

type contextKey string

const (
	ctxUserID        contextKey = "userID"
	ctxInstanceAdmin contextKey = "instanceAdmin"
	ctxClaims        contextKey = "claims"
)

// UserIDFromContext extracts the authenticated user ID from the request context.
func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxUserID).(string)
	return id
}

// ContextWithUser returns ctx with the same userID/instance-admin values
// AuthMiddleware would set from a validated token, for tests that call a
// handler directly instead of going through the full middleware chain.
func ContextWithUser(ctx context.Context, userID string, instanceAdmin bool) context.Context {
	ctx = context.WithValue(ctx, ctxUserID, userID)
	return context.WithValue(ctx, ctxInstanceAdmin, instanceAdmin)
}

// InstanceAdminFromContext reports whether the authenticated user holds the
// flat, non-connector-scoped instance-admin role. Per-connector access is
// never carried in the request context; look it up per-request via
// store.UserHasConnectorRole instead (see RequireConnectorRole).
func InstanceAdminFromContext(ctx context.Context) bool {
	admin, _ := ctx.Value(ctxInstanceAdmin).(bool)
	return admin
}

// APIKeyChecker looks up opaque API keys without coupling auth to the store
// package. The store adapts its APIKey model to APIKeyClaims.
type APIKeyChecker interface {
	LookupAPIKey(ctx context.Context, tokenHash string) (*APIKeyClaims, error)
	TouchAPIKeyLastUsed(ctx context.Context, keyID string) error
}

// UserStatusChecker returns a user's current instance-admin role (as
// "admin"/"user") and disabled flag, so AuthMiddleware can reject an access
// token whose claims have gone stale (a role change or disable revokes
// access immediately instead of waiting out the access token's TTL).
// Optional: implemented by *store.Store via a type assertion on the
// APIKeyChecker passed to AuthMiddleware, so existing call sites and
// lightweight test doubles keep working unchanged.
type UserStatusChecker interface {
	GetUserRoleStatus(ctx context.Context, userID string) (role string, disabled bool, err error)
}

// AuthMiddleware validates JWT access tokens and opaque API keys, then injects
// userID + role into the request context. The variadic checker preserves the
// lightweight JWT-only call shape used by auth package tests and callers that
// do not have API-key storage.
func AuthMiddleware(jwtSvc *Service, checkers ...APIKeyChecker) func(http.Handler) http.Handler { //nolint:revive
	var checker APIKeyChecker
	if len(checkers) > 0 {
		checker = checkers[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				http.Error(w, `{"code":"unauthorized","message":"missing or malformed Authorization header"}`, http.StatusUnauthorized)
				return
			}

			claims, err := jwtSvc.ValidateAccess(token)
			if err == nil {
				if statusChecker, ok := checker.(UserStatusChecker); ok {
					role, disabled, statusErr := statusChecker.GetUserRoleStatus(r.Context(), claims.UserID)
					if statusErr != nil || disabled || (role == "admin") != claims.InstanceAdmin {
						http.Error(w, `{"code":"unauthorized","message":"session no longer valid"}`, http.StatusUnauthorized)
						return
					}
				}
				ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
				ctx = context.WithValue(ctx, ctxInstanceAdmin, claims.InstanceAdmin)
				ctx = context.WithValue(ctx, ctxClaims, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if checker != nil {
				keyClaims, lookupErr := checker.LookupAPIKey(r.Context(), hashToken(token))
				if lookupErr == nil && validAPIKey(keyClaims) {
					lastUsed, _ := time.Parse(time.RFC3339, keyClaims.LastUsedAt)
					if time.Since(lastUsed) >= time.Minute {
						if touchErr := checker.TouchAPIKeyLastUsed(r.Context(), keyClaims.KeyID); touchErr != nil {
							slog.Error("failed to update API key last-used timestamp", "key_id", keyClaims.KeyID, "error", touchErr)
						}
					}
					restriction := keyClaims.Restriction
					if restriction.ReadOnly && !isSafeMethod(r.Method) {
						httputil.Error(w, http.StatusForbidden, "forbidden", "API key is read-only")
						return
					}
					ctx := context.WithValue(r.Context(), ctxUserID, keyClaims.UserID)
					// A connector-restricted key never carries instance admin:
					// admin endpoints aren't connector-scoped, so the
					// restriction couldn't be honored there.
					ctx = context.WithValue(ctx, ctxInstanceAdmin, keyClaims.InstanceAdmin && len(restriction.ConnectorIDs) == 0)
					ctx = context.WithValue(ctx, ctxClaims, keyClaims)
					ctx = ContextWithAPIKeyRestriction(ctx, restriction)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			http.Error(w, `{"code":"unauthorized","message":"invalid or expired token"}`, http.StatusUnauthorized)
		})
	}
}

func hashToken(token string) string {
	// Keep hashing in auth so the middleware stays independent of store.
	// SHA-256 is also the repository's established token-hash convention.
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func validAPIKey(claims *APIKeyClaims) bool {
	if claims == nil || claims.KeyID == "" || claims.UserID == "" || claims.RevokedAt != "" {
		return false
	}
	if claims.ExpiresAt == "" {
		return true
	}
	expiresAt, err := time.Parse(time.RFC3339, claims.ExpiresAt)
	return err == nil && time.Now().UTC().Before(expiresAt)
}

// RequireInstanceAdmin returns middleware that rejects any request not made
// by an instance admin. Used for actions that aren't connector-scoped: user
// management, API keys, granting/revoking connector permissions. For
// connector-scoped mutations, use RequireConnectorRole instead.
func RequireInstanceAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !InstanceAdminFromContext(r.Context()) {
			httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ConnectorRoleChecker looks up a user's per-connector role. Implemented by
// *store.Store (see internal/store/connector_permission.go) to avoid an
// import cycle between internal/auth and internal/store.
type ConnectorRoleChecker interface {
	UserHasConnectorRole(ctx context.Context, userID, connectorID, minRole string) (bool, error)
}

// RequireConnectorRole returns middleware that reads a connector ID from the
// named chi URL parameter and requires the authenticated user hold at least
// minRole ("viewer" or "operator") on that specific connector. No grant
// means no access (default deny) — an instance admin is not implicitly
// granted access to a connector they haven't been given a role on.
//
// Endpoints whose connector ID isn't in the URL path (a body-supplied ID
// list, or an ID that must be resolved through another resource first, e.g.
// a doc ID) can't use this middleware; they call checker.UserHasConnectorRole
// directly in the handler instead.
func RequireConnectorRole(checker ConnectorRoleChecker, minRole, idParam string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			connectorID := r.PathValue(idParam)
			ok, err := checker.UserHasConnectorRole(r.Context(), UserIDFromContext(r.Context()), connectorID, minRole)
			if err != nil || !ok {
				httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PermissionChecker looks up whether a user has a named boolean permission flag.
// Implemented by *store.Store (see internal/store/user.go) to avoid an import
// cycle between internal/auth and internal/store.
type PermissionChecker interface {
	UserHasPermission(ctx context.Context, userID, permission string) (bool, error)
}

// AuditRecorder records security-sensitive auth events without coupling this
// package to store.
type AuditRecorder interface {
	RecordAuditFromContext(ctx context.Context, action, targetType, targetID string, detail any) error
}

// RequirePermission returns middleware requiring instance-admin AND a named
// per-user boolean permission flag (e.g. "can_manage_dashboard_defaults").
func RequirePermission(checker PermissionChecker, permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !InstanceAdminFromContext(r.Context()) {
				httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
				return
			}
			ok, err := checker.UserHasPermission(r.Context(), UserIDFromContext(r.Context()), permission)
			if err != nil || !ok {
				httputil.Error(w, http.StatusForbidden, "forbidden", "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireElevation checks for a valid elevation token scoped to the given action.
// Destructive endpoints chain this after RequireRole("operator") for step-up auth.
func RequireElevation(jwtSvc *Service, recorder AuditRecorder, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := ValidateElevationHeader(jwtSvc, recorder, action, r); err != nil {
				WriteElevationError(w, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// elevationError carries the HTTP status and body an elevation failure
// should produce, so ValidateElevationHeader can be called both as
// middleware and directly from a handler that only elevates conditionally.
type elevationError struct {
	status int
	code   string
	msg    string
}

func (e *elevationError) Error() string { return e.msg }

// ValidateElevationHeader checks the request's X-Elevation-Token header
// against a token scoped to action, recording the same
// auth.elevation_requested / auth.elevation_denied audit events RequireElevation
// records. Returns nil if the token is valid; otherwise an *elevationError
// describing the HTTP response to send (use WriteElevationError, or inspect
// via errors.As for a custom response).
func ValidateElevationHeader(jwtSvc *Service, recorder AuditRecorder, action string, r *http.Request) error {
	values, ok := r.Header["X-Elevation-Token"]
	if !ok {
		return &elevationError{
			status: http.StatusBadRequest,
			code:   "elevation_required",
			msg:    "X-Elevation-Token header required for " + action,
		}
	}
	token := ""
	if len(values) > 0 {
		token = values[0]
	}
	recordElevationAudit(r.Context(), recorder, "auth.elevation_requested", action, nil)
	_, err := jwtSvc.ValidateElevation(token, action, UserIDFromContext(r.Context()))
	if err != nil {
		recordElevationAudit(r.Context(), recorder, "auth.elevation_denied", action, map[string]any{
			"action": action,
			"reason": elevationFailureReason(err),
		})
		return &elevationError{
			status: http.StatusUnauthorized,
			code:   "unauthorized",
			msg:    "Invalid elevation token: " + err.Error(),
		}
	}
	return nil
}

// WriteElevationError writes the HTTP response for an error returned by
// ValidateElevationHeader.
func WriteElevationError(w http.ResponseWriter, err error) {
	var elevErr *elevationError
	if !errors.As(err, &elevErr) {
		httputil.Errorf(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(elevErr.status)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"code":    elevErr.code,
		"message": elevErr.msg,
	}); err != nil {
		slog.Debug("failed to write elevation error response", "error", err)
	}
}

func recordElevationAudit(ctx context.Context, recorder AuditRecorder, event, action string, detail any) {
	if recorder == nil {
		return
	}
	if err := recorder.RecordAuditFromContext(ctx, event, "action", action, detail); err != nil {
		slog.Error("failed to record audit", "action", event, "error", err)
	}
}

func elevationFailureReason(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "expired"):
		return "expired"
	case strings.Contains(msg, "different user"):
		return "wrong_user"
	case strings.Contains(msg, "not an elevation token"):
		return "wrong_token_type"
	case strings.Contains(msg, "is for action"):
		return "wrong_action"
	default:
		return "invalid"
	}
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return ""
	}
	return strings.TrimPrefix(auth, prefix)
}
