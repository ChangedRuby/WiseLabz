package auth

import (
	"context"
	"net/http"
	"slices"

	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/httpx"
)

// API-key scopes (#278). ScopeFull keys act with their owner's access;
// ScopeRead keys may only make safe (read-only) requests.
const (
	APIKeyScopeFull = "full"
	APIKeyScopeRead = "read"
)

// APIKeyRestriction narrows an API key below its owner's access. A key can
// never do more than its owner; these fields only take access away.
type APIKeyRestriction struct {
	// ReadOnly limits the key to safe HTTP methods and caps its role on any
	// connector at viewer.
	ReadOnly bool
	// ConnectorIDs, when non-empty, is the only set of connectors the key
	// can reach. A connector-restricted key never holds instance admin.
	ConnectorIDs []string
}

// Restricted reports whether the restriction narrows access at all.
func (r APIKeyRestriction) Restricted() bool {
	return r.ReadOnly || len(r.ConnectorIDs) > 0
}

const ctxAPIKeyRestriction contextKey = "apiKeyRestriction"

// ContextWithAPIKeyRestriction returns ctx carrying an API key's
// restriction, as AuthMiddleware sets it for key-authenticated requests.
func ContextWithAPIKeyRestriction(ctx context.Context, r APIKeyRestriction) context.Context {
	return context.WithValue(ctx, ctxAPIKeyRestriction, r)
}

// APIKeyRestrictionFromContext returns the restriction of the API key that
// authenticated the request. JWT sessions and unrestricted keys yield the
// zero value.
func APIKeyRestrictionFromContext(ctx context.Context) APIKeyRestriction {
	r, _ := ctx.Value(ctxAPIKeyRestriction).(APIKeyRestriction)
	return r
}

// ClampConnectorRole applies the request's API-key restriction to a user's
// role on a connector: "" when the key can't reach the connector, "viewer"
// when a read-only key would otherwise hold operator. Store permission
// checks call this so every connector-scoped endpoint honors key
// restrictions without per-handler code.
func ClampConnectorRole(ctx context.Context, connectorID, role string) string {
	r := APIKeyRestrictionFromContext(ctx)
	if role == "" || !r.Restricted() {
		return role
	}
	if len(r.ConnectorIDs) > 0 && !slices.Contains(r.ConnectorIDs, connectorID) {
		return ""
	}
	if r.ReadOnly {
		return "viewer"
	}
	return role
}

// RejectRestrictedAPIKey writes a 403 and returns true when the request was
// authenticated by a restricted API key. Endpoints that would let a key
// escape its restriction (minting new keys, opening a WebSocket that isn't
// filtered per key) call this first.
func RejectRestrictedAPIKey(w http.ResponseWriter, r *http.Request) bool {
	if !APIKeyRestrictionFromContext(r.Context()).Restricted() {
		return false
	}
	httputil.Error(w, http.StatusForbidden, "forbidden", "not available to a restricted API key")
	return true
}

func isSafeMethod(method string) bool {
	return httpx.IsSafeMethod(method)
}
