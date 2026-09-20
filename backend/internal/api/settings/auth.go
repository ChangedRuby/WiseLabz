package settings

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// GetAuthConfig handles GET /api/auth/config.
func (h *Handler) GetAuthConfig(w http.ResponseWriter, r *http.Request) {
	var rec struct {
		LocalEnabled         int
		AccessTokenTTL       int
		RefreshTokenTTL      int
		StepUpForDestructive int
	}
	err := h.Store.DB().QueryRowContext(r.Context(), `
		SELECT local_enabled, access_token_ttl, refresh_token_ttl, step_up_for_destructive
		FROM auth_config WHERE id = 1
	`).Scan(&rec.LocalEnabled, &rec.AccessTokenTTL, &rec.RefreshTokenTTL, &rec.StepUpForDestructive)

	localEnabled := true
	accessTTL := int(h.Config.Auth.AccessTokenTTLDuration().Seconds())
	refreshTTL := int(h.Config.Auth.RefreshTokenTTLDuration().Seconds())
	stepUp := h.Config.Auth.StepUpForDestructive
	if err == nil {
		localEnabled = rec.LocalEnabled != 0
		accessTTL = rec.AccessTokenTTL
		refreshTTL = rec.RefreshTokenTTL
		stepUp = rec.StepUpForDestructive != 0
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"localEnabled":         localEnabled,
		"accessTokenTtl":       accessTTL,
		"refreshTokenTtl":      refreshTTL,
		"stepUpForDestructive": stepUp,
		"oidcProviders":        h.oidcProviders(r.Context()),
	})
}

// UpdateAuthConfig handles PUT /api/auth/config.
func (h *Handler) UpdateAuthConfig(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		LocalEnabled         *bool `json:"localEnabled"`
		AccessTokenTTL       *int  `json:"accessTokenTtl"`
		RefreshTokenTTL      *int  `json:"refreshTokenTtl"`
		StepUpForDestructive *bool `json:"stepUpForDestructive"`
	}](w, r)
	if !ok {
		return
	}

	var args []any
	var parts []string

	if req.LocalEnabled != nil {
		parts = append(parts, "local_enabled = ?")
		args = append(args, boolToInt(*req.LocalEnabled))
	}
	if req.AccessTokenTTL != nil {
		parts = append(parts, "access_token_ttl = ?")
		args = append(args, *req.AccessTokenTTL)
	}
	if req.RefreshTokenTTL != nil {
		parts = append(parts, "refresh_token_ttl = ?")
		args = append(args, *req.RefreshTokenTTL)
	}
	if req.StepUpForDestructive != nil {
		parts = append(parts, "step_up_for_destructive = ?")
		args = append(args, boolToInt(*req.StepUpForDestructive))
	}

	if len(parts) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "No fields to update")
		return
	}

	query := "UPDATE auth_config SET " + strings.Join(parts, ", ") + " WHERE id = 1"

	_, err := h.Store.DB().ExecContext(r.Context(), query, args...)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Fields only, not values — TTLs and toggles aren't secret, but this
	// keeps the shape consistent with other config-mutating audit entries.
	fields := make([]string, len(parts))
	for i, p := range parts {
		fields[i] = strings.SplitN(p, " = ", 2)[0]
	}
	if err := h.Store.RecordAuditFromContext(r.Context(), "auth.config.update", "auth_config", "", map[string]any{
		"fields": fields,
	}); err != nil {
		slog.Error("failed to record audit", "action", "auth.config.update", "error", err)
	}

	h.GetAuthConfig(w, r)
}

// UpdateProviderEnabled handles PUT /api/auth/providers/{providerId}/enabled.
func (h *Handler) UpdateProviderEnabled(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("providerId")

	var found *config.OIDCProvider
	for i := range h.Config.Auth.OIDC {
		if h.Config.Auth.OIDC[i].ID == id {
			found = &h.Config.Auth.OIDC[i]
			break
		}
	}
	if found == nil {
		httputil.Error(w, http.StatusNotFound, "not_found", "Unknown OIDC provider")
		return
	}

	req, ok := httputil.DecodeJSON[struct {
		Enabled bool `json:"enabled"`
	}](w, r)
	if !ok {
		return
	}

	if err := h.Store.SetOIDCProviderEnabled(r.Context(), id, req.Enabled); err != nil {
		httputil.Errorf(w, err)
		return
	}

	detail := map[string]any{
		"enabled": req.Enabled,
	}
	if !req.Enabled {
		revoked, err := h.Store.RevokeSessionsByAuthSource(r.Context(), id)
		if err != nil {
			httputil.Errorf(w, err)
			return
		}
		detail["revokedSessions"] = revoked
	}

	if err := h.Store.RecordAuditFromContext(r.Context(), "auth.provider.enabled", "oidc_provider", id, detail); err != nil {
		slog.Error("failed to record audit", "action", "auth.provider.enabled", "error", err)
	}

	httputil.JSON(w, http.StatusOK, oidcProviderJSON(found, req.Enabled))
}

func (h *Handler) oidcProviders(ctx context.Context) []map[string]any {
	flags, err := h.Store.GetOIDCProviderFlags(ctx)
	if err != nil {
		slog.Error("failed to get OIDC provider flags", "error", err)
	}
	out := make([]map[string]any, 0, len(h.Config.Auth.OIDC))
	for i := range h.Config.Auth.OIDC {
		p := &h.Config.Auth.OIDC[i]
		enabled := true
		if v, ok := flags[p.ID]; ok {
			enabled = v
		}
		out = append(out, oidcProviderJSON(p, enabled))
	}
	return out
}

func oidcProviderJSON(p *config.OIDCProvider, enabled bool) map[string]any {
	return map[string]any{
		"id":               p.ID,
		"displayName":      p.DisplayName,
		"issuerUrl":        p.IssuerURL,
		"clientId":         p.ClientID,
		"secretConfigured": p.ClientSecret != "",
		"enabled":          enabled,
		"source":           "file",
		"scopes":           p.Scopes,
	}
}
