package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// oidcFlowCookie is set on the browser that starts an OIDC login (from
// Providers) and read back on OIDCCallback, binding the state/nonce to that
// specific browser instead of trusting the caller-supplied state alone.
// GET /api/auth/providers is unauthenticated, so a self-contained signed
// state token can be handed to anyone and replayed from a different
// browser/session (CSRF / authorization-code injection); the cookie closes
// that gap.
const oidcFlowCookie = "oidc_flow"

// OIDCCallback handles POST /api/auth/oidc/callback.
// Exchanges an OIDC authorization code for identity, creates or finds the user,
// and returns a JWT token pair.
func (h *Handler) OIDCCallback(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		ProviderID string `json:"providerId"`
		Code       string `json:"code"`
		State      string `json:"state"`
	}](w, r)
	if !ok {
		return
	}
	if req.ProviderID == "" || req.Code == "" || req.State == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "providerId, code, and state are required")
		return
	}

	nonce, ok := h.verifyOIDCFlowState(w, r, req.ProviderID, req.State)
	if !ok {
		return
	}

	provCfg, prov, ok := h.resolveOIDCProvider(w, r, req.ProviderID)
	if !ok {
		return
	}

	claims, ok := h.exchangeOIDCCode(w, r, prov, provCfg, req.ProviderID, req.Code, nonce)
	if !ok {
		return
	}
	role := oidcRoleForGroups(claims.Groups, provCfg.GroupRoleMapping)

	user, isNewUser, ok := h.resolveOIDCUser(w, r, claims, role)
	if !ok {
		return
	}

	h.completeOIDCLogin(w, r, user, req.ProviderID, isNewUser)
}

// verifyOIDCFlowState reads back the flow cookie this browser was issued when
// it started the login, clears it so it cannot be replayed, and constant-time
// compares the cookie's state against the caller-supplied state. Returns the
// nonce bound to this browser's attempt.
func (h *Handler) verifyOIDCFlowState(w http.ResponseWriter, r *http.Request, providerID, state string) (string, bool) {
	cookieState, nonce, ok := readOIDCFlowCookie(r, providerID)
	clearOIDCFlowCookie(w, r, h.Config.Server.TrustedProxies, providerID)
	if !ok || subtle.ConstantTimeCompare([]byte(cookieState), []byte(state)) != 1 {
		httputil.Error(w, http.StatusUnauthorized, "oidc_error", "Invalid or expired state")
		return "", false
	}
	return nonce, true
}

// resolveOIDCProvider finds the provider configuration, rejects it if an admin
// has disabled it, and initializes the provider if needed.
func (h *Handler) resolveOIDCProvider(w http.ResponseWriter, r *http.Request, providerID string) (*config.OIDCProvider, *auth.OIDCProvider, bool) {
	provCfg := h.findOIDCProvider(providerID)
	if provCfg == nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_provider", "Unknown OIDC provider")
		return nil, nil, false
	}
	if !h.oidcProviderEnabled(r.Context(), providerID) {
		httputil.Error(w, http.StatusForbidden, "oidc_error", "OIDC provider is disabled")
		return nil, nil, false
	}
	prov := h.getOrInitOIDCProvider(r.Context(), provCfg)
	if prov == nil {
		httputil.Error(w, http.StatusInternalServerError, "oidc_error", "Failed to initialize OIDC provider")
		return nil, nil, false
	}
	return provCfg, prov, true
}

// exchangeOIDCCode trades the authorization code for identity claims against
// the pinned redirect URL, then enforces the verified-email requirement and
// the provider's email domain allowlist.
func (h *Handler) exchangeOIDCCode(w http.ResponseWriter, r *http.Request, prov *auth.OIDCProvider, provCfg *config.OIDCProvider, providerID, code, nonce string) (*auth.OIDCClaims, bool) {
	claims, err := prov.Exchange(r.Context(), code, nonce, h.oidcRedirectURL(r))
	if err != nil {
		slog.Error("OIDC exchange failed", "error", err, "provider", logsafe.Sanitize(providerID))
		httputil.Error(w, http.StatusUnauthorized, "oidc_error", "Failed to authenticate with provider")
		return nil, false
	}
	if claims.Subject == "" || claims.Issuer == "" || !claims.EmailVerified {
		httputil.Error(w, http.StatusUnauthorized, "oidc_error", "OIDC identity requires a verified email")
		return nil, false
	}
	if len(provCfg.EmailDomainAllowlist) > 0 && !emailDomainAllowed(claims.Email, provCfg.EmailDomainAllowlist) {
		httputil.Error(w, http.StatusForbidden, "oidc_error", "OIDC email domain is not allowed")
		return nil, false
	}
	return claims, true
}

// resolveOIDCUser finds or creates the local user for an OIDC identity,
// rejects disabled accounts, and keeps the instance-admin role in sync with
// the provider's group mapping.
func (h *Handler) resolveOIDCUser(w http.ResponseWriter, r *http.Request, claims *auth.OIDCClaims, role string) (*store.User, bool, bool) {
	user, err := h.Store.GetUserByOIDCIdentity(r.Context(), claims.Issuer, claims.Subject)
	isNewUser := false
	if errors.Is(err, store.ErrNotFound) {
		user = newOIDCUser(claims, role)
		created, err := h.Store.CreateOIDCUser(r.Context(), user, claims.Issuer, claims.Subject)
		if errors.Is(err, store.ErrConflict) {
			// Lost a race with a concurrent create for the same identity; use the row that won.
			user, err = h.Store.GetUserByOIDCIdentity(r.Context(), claims.Issuer, claims.Subject)
		}
		if err != nil {
			httputil.Errorf(w, err)
			return nil, false, false
		}
		isNewUser = created
	} else if err != nil {
		httputil.Errorf(w, err)
		return nil, false, false
	}
	if user.Disabled {
		httputil.Error(w, http.StatusForbidden, "forbidden", "Account is disabled")
		return nil, false, false
	}
	if !isNewUser && user.AuthSource == "oidc" && user.InstanceAdminRole != role {
		if err := h.Store.UpdateUser(r.Context(), user.ID, map[string]any{"instance_admin_role": role}); err != nil {
			httputil.Errorf(w, err)
			return nil, false, false
		}
		user.InstanceAdminRole = role
	}
	return user, isNewUser, true
}

// newOIDCUser builds the local record for a first-time OIDC identity. The
// username is derived from a hash of issuer+subject so it is stable across
// logins without leaking the provider's identifiers.
func newOIDCUser(claims *auth.OIDCClaims, role string) *store.User {
	displayName := claims.Name
	if displayName == "" {
		displayName = claims.PreferredName
	}
	if displayName == "" {
		displayName = claims.Email
	}
	identityHash := sha256.Sum256([]byte(claims.Issuer + "\x00" + claims.Subject))
	return &store.User{
		Username:          fmt.Sprintf("oidc_%x", identityHash[:]),
		DisplayName:       displayName,
		Email:             claims.Email,
		InstanceAdminRole: role,
		AuthSource:        "oidc",
	}
}

// completeOIDCLogin issues the token pair, records the refresh session and
// writes the login response.
func (h *Handler) completeOIDCLogin(w http.ResponseWriter, r *http.Request, user *store.User, providerID string, isNewUser bool) {
	pair, err := h.JWT.IssuePair(user.ID, user.InstanceAdminRole == "admin")
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	session := &store.Session{
		UserID:         user.ID,
		TokenHash:      store.HashToken(pair.RefreshToken),
		AuthProviderID: providerID,
		UserAgent:      r.UserAgent(),
		IP:             httputil.ClientIP(r, h.Config.Server.TrustedProxies),
	}
	if err := h.Store.CreateSession(r.Context(), session); err != nil {
		httputil.Errorf(w, err)
		return
	}

	setRefreshCookie(w, r, h.Config.Server.TrustedProxies, pair.RefreshToken, h.Config.Auth.RefreshTokenTTLDuration())

	httputil.JSON(w, http.StatusOK, map[string]any{
		"accessToken": pair.AccessToken,
		"expiresIn":   pair.ExpiresIn,
		"user":        sanitizeUser(user),
		"isNewUser":   isNewUser,
	})
}

// Providers handles GET /api/auth/providers.
// Returns OIDC providers merged from config file (secrets hidden) and DB enable flags.
func (h *Handler) Providers(w http.ResponseWriter, r *http.Request) {
	type providerInfo struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		AuthURL     string `json:"authUrl"`
	}

	redirectURL := h.oidcRedirectURL(r)

	var oidc []providerInfo
	for i := range h.Config.Auth.OIDC {
		p := &h.Config.Auth.OIDC[i]
		if !h.oidcProviderEnabled(r.Context(), p.ID) {
			continue
		}
		prov := h.getOrInitOIDCProvider(r.Context(), p)
		if prov == nil {
			slog.Warn("skipping OIDC provider in list due to initialization failure", "id", p.ID)
			continue
		}

		state, err := randomOIDCToken()
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		nonce, err := randomOIDCToken()
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		setOIDCFlowCookie(w, r, h.Config.Server.TrustedProxies, p.ID, state, nonce)
		authURL := prov.AuthURL(state, nonce, redirectURL)
		oidc = append(oidc, providerInfo{
			ID:          p.ID,
			DisplayName: p.DisplayName,
			AuthURL:     authURL,
		})
	}

	if oidc == nil {
		oidc = []providerInfo{}
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"oidc":         oidc,
		"localEnabled": true,
	})
}

// --- OIDC helpers ---

// validHostPort reports whether s is a syntactically plain host[:port] with no
// slashes, userinfo, whitespace or other URL metacharacters.
func validHostPort(s string) bool {
	if s == "" || len(s) > 255 {
		return false
	}
	if strings.ContainsAny(s, "/\\@?#% \t\r\n") {
		return false
	}
	u, err := url.Parse("http://" + s)
	if err != nil || u.Host != s || u.Hostname() == "" {
		return false
	}
	return u.User == nil && u.Path == "" && u.RawQuery == "" && u.Fragment == ""
}

// oidcRedirectURL builds the callback URL. r.Host is attacker-influenceable,
// so configured Server.Origin entries take precedence: a matching origin is
// used as-is, a single configured origin is used unconditionally, and only
// otherwise is r.Host used (if syntactically valid).
func (h *Handler) oidcRedirectURL(r *http.Request) string {
	var origins []*url.URL
	for _, o := range strings.Split(h.Config.Server.Origin, ",") {
		u, err := url.Parse(strings.TrimSpace(o))
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			continue
		}
		origins = append(origins, u)
	}
	for _, u := range origins {
		if strings.EqualFold(u.Host, r.Host) {
			return u.Scheme + "://" + u.Host + "/auth/callback"
		}
	}
	if len(origins) == 1 {
		return origins[0].Scheme + "://" + origins[0].Host + "/auth/callback"
	}
	scheme := "http"
	if httputil.IsSecureRequest(r, h.Config.Server.TrustedProxies) {
		scheme = "https"
	}
	host := r.Host
	if !validHostPort(host) {
		host = "localhost"
	}
	return fmt.Sprintf("%s://%s/auth/callback", scheme, host)
}

func (h *Handler) findOIDCProvider(id string) *config.OIDCProvider {
	for i := range h.Config.Auth.OIDC {
		if h.Config.Auth.OIDC[i].ID == id {
			return &h.Config.Auth.OIDC[i]
		}
	}
	return nil
}

func (h *Handler) getOrInitOIDCProvider(ctx context.Context, cfg *config.OIDCProvider) *auth.OIDCProvider {
	h.oidcMu.Lock()
	defer h.oidcMu.Unlock()
	if h.oidcProv == nil {
		h.oidcProv = make(map[string]*auth.OIDCProvider)
	}
	if p, ok := h.oidcProv[cfg.ID]; ok && p.IsInitialized() {
		return p
	}

	prov := &auth.OIDCProvider{
		ID:           cfg.ID,
		DisplayName:  cfg.DisplayName,
		IssuerURL:    cfg.IssuerURL,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       cfg.Scopes,
		GroupsClaim:  cfg.GroupsClaim,
	}
	if err := prov.Initialize(ctx); err != nil {
		slog.Error("failed to init OIDC provider", "id", cfg.ID, "error", err)
		return nil
	}
	h.oidcProv[cfg.ID] = prov
	return prov
}

func (h *Handler) oidcProviderEnabled(ctx context.Context, providerID string) bool {
	flags, err := h.Store.GetOIDCProviderFlags(ctx)
	if err != nil {
		slog.Error("failed to get OIDC provider flags", "error", err)
		return false
	}
	enabled, configured := flags[providerID]
	return !configured || enabled
}

func emailDomainAllowed(email string, allowlist []string) bool {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(email)), "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	domain := parts[1]
	for _, allowed := range allowlist {
		if strings.EqualFold(strings.TrimPrefix(strings.TrimSpace(allowed), "@"), domain) {
			return true
		}
	}
	return false
}

// oidcRoleForGroups maps OIDC groups to the flat instance-admin role
// ("admin"/"user") via the admin-configured GroupRoleMapping. Accepts the
// mapping's legacy "operator"/"viewer" values too, so existing provider
// configs don't need updating alongside this migration.
func oidcRoleForGroups(groups []string, mapping map[string]string) string {
	role := "user"
	for _, group := range groups {
		mapped := mapping[group]
		if strings.EqualFold(mapped, "operator") || strings.EqualFold(mapped, "admin") {
			return "admin"
		}
		if strings.EqualFold(mapped, "viewer") || strings.EqualFold(mapped, "user") {
			role = "user"
		}
	}
	return role
}

// randomOIDCToken returns a random URL-safe token used as an OIDC state or
// nonce value.
func randomOIDCToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate oidc token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// oidcFlowCookieName derives a per-provider cookie name so concurrent login
// attempts against different providers don't clobber each other's cookie.
func oidcFlowCookieName(providerID string) string {
	sum := sha256.Sum256([]byte(providerID))
	return oidcFlowCookie + "_" + base64.RawURLEncoding.EncodeToString(sum[:12])
}

// setOIDCFlowCookie stores the state/nonce generated for this browser's login
// attempt in a short-lived HttpOnly cookie, scoped to the auth endpoints.
func setOIDCFlowCookie(w http.ResponseWriter, r *http.Request, trustedProxies, providerID, state, nonce string) {
	http.SetCookie(w, &http.Cookie{
		Name:     oidcFlowCookieName(providerID),
		Value:    providerID + "." + state + "." + nonce,
		Path:     "/api/auth",
		MaxAge:   300,
		HttpOnly: true,
		// Secure is derived, not literal true, so CodeQL can't verify it;
		// IsSecureRequest returns true for both direct TLS and a trusted
		// TLS-terminating proxy.
		Secure:   httputil.IsSecureRequest(r, trustedProxies), // codeql[go/cookie-secure-not-set]
		SameSite: http.SameSiteLaxMode,
	})
}

// readOIDCFlowCookie returns the state and nonce this browser was issued for
// providerID, if any.
func readOIDCFlowCookie(r *http.Request, providerID string) (state, nonce string, ok bool) {
	cookie, err := r.Cookie(oidcFlowCookieName(providerID))
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(cookie.Value, ".", 3)
	if len(parts) != 3 || parts[0] != providerID || parts[1] == "" || parts[2] == "" {
		return "", "", false
	}
	return parts[1], parts[2], true
}

// clearOIDCFlowCookie deletes the flow cookie so it cannot be replayed.
func clearOIDCFlowCookie(w http.ResponseWriter, r *http.Request, trustedProxies, providerID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     oidcFlowCookieName(providerID),
		Value:    "",
		Path:     "/api/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   httputil.IsSecureRequest(r, trustedProxies), // codeql[go/cookie-secure-not-set]
		SameSite: http.SameSiteLaxMode,
	})
}
