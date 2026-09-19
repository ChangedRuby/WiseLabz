// Package settings provides API handlers for auth, AI, and notification configuration.
package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/crypto"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// Handler holds dependencies for settings endpoints.
type Handler struct {
	Store  *store.Store
	Config *config.Config
	AI     *ai.Registry
}

// NewHandler creates a new settings handler.
func NewHandler(s *store.Store, cfg *config.Config, aiRegistry *ai.Registry) *Handler {
	return &Handler{Store: s, Config: cfg, AI: aiRegistry}
}

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

// AIConfigValues is the resolved AI configuration (DB row, falling back to
// static config), including the decrypted API key for provider construction.
type AIConfigValues struct {
	Enabled  bool
	Provider string
	Model    string
	BaseURL  string
	APIKey   string
	Mode     string

	// Embedding backend for "ask your lab" chat retrieval, independent of
	// Provider/Model/APIKey/BaseURL above.
	EmbedProvider string
	EmbedModel    string
	EmbedAPIKey   string
	EmbedBaseURL  string

	// Providers is the ordered fallback chain: the primary provider above
	// (if configured) followed by any rows in ai_config_providers, in
	// priority order. Callers pass this to ai.SuggestWithFallback.
	Providers []ai.ProviderConfig
}

// LoadAIConfig reads the current AI configuration, decrypting the stored API key.
func (h *Handler) LoadAIConfig(ctx context.Context) AIConfigValues {
	var rec struct {
		Enabled       int
		Provider      sql.NullString
		Model         sql.NullString
		BaseURL       sql.NullString
		Mode          string
		EmbedProvider sql.NullString
		EmbedModel    sql.NullString
		EmbedBaseURL  sql.NullString
	}
	err := h.Store.DB().QueryRowContext(ctx, `
		SELECT enabled, provider, model, base_url, mode, embed_provider, embed_model, embed_base_url
		FROM ai_config WHERE id = 1
	`).Scan(&rec.Enabled, &rec.Provider, &rec.Model, &rec.BaseURL, &rec.Mode,
		&rec.EmbedProvider, &rec.EmbedModel, &rec.EmbedBaseURL)
	if err != nil {
		cfg := AIConfigValues{
			Enabled: h.Config.AI.Enabled, Provider: h.Config.AI.Provider, Model: h.Config.AI.Model,
			BaseURL: h.Config.AI.BaseURL, APIKey: h.Config.AI.APIKey, Mode: h.Config.AI.Mode,
			EmbedProvider: h.Config.AI.EmbedProvider, EmbedModel: h.Config.AI.EmbedModel,
			EmbedAPIKey: h.Config.AI.EmbedAPIKey, EmbedBaseURL: h.Config.AI.EmbedBaseURL,
		}
		cfg.Providers = primaryProviderConfig(cfg)
		return cfg
	}
	embedProvider := rec.EmbedProvider.String
	if embedProvider == "" {
		embedProvider = h.Config.AI.EmbedProvider
	}
	embedModel := rec.EmbedModel.String
	if embedModel == "" {
		embedModel = h.Config.AI.EmbedModel
	}
	cfg := AIConfigValues{
		Enabled: rec.Enabled != 0, Provider: rec.Provider.String, Model: rec.Model.String,
		BaseURL: rec.BaseURL.String, APIKey: h.GetDecryptedAPIKey(), Mode: rec.Mode,
		EmbedProvider: embedProvider, EmbedModel: embedModel,
		EmbedAPIKey: h.GetDecryptedEmbedAPIKey(), EmbedBaseURL: rec.EmbedBaseURL.String,
	}
	cfg.Providers = append(primaryProviderConfig(cfg), h.loadFallbackProviders(ctx)...)
	return cfg
}

// primaryProviderConfig returns cfg's primary provider as a single-entry
// fallback chain, or none if no provider is configured.
func primaryProviderConfig(cfg AIConfigValues) []ai.ProviderConfig {
	if cfg.Provider == "" {
		return nil
	}
	return []ai.ProviderConfig{{Name: cfg.Provider, Model: cfg.Model, APIKey: cfg.APIKey, BaseURL: cfg.BaseURL}}
}

// loadFallbackProviders reads the ai_config_providers table (priority 2+),
// decrypting each row's API key.
func (h *Handler) loadFallbackProviders(ctx context.Context) []ai.ProviderConfig {
	rows, err := h.Store.DB().QueryContext(ctx, `
		SELECT provider, model, api_key_encrypted, base_url
		FROM ai_config_providers WHERE config_id = 1 ORDER BY priority ASC
	`)
	if err != nil {
		slog.Error("failed to load fallback AI providers", "error", err)
		return nil
	}
	defer rows.Close() //nolint:errcheck

	key, keyErr := crypto.DecodeKey(h.Config.Encryption.Key)
	var out []ai.ProviderConfig
	for rows.Next() {
		var provider, model, encryptedKey, baseURL string
		if err := rows.Scan(&provider, &model, &encryptedKey, &baseURL); err != nil {
			slog.Error("failed to scan fallback AI provider row", "error", err)
			continue
		}
		apiKey := ""
		if encryptedKey != "" && keyErr == nil {
			if plaintext, err := crypto.Decrypt(encryptedKey, key); err == nil {
				apiKey = plaintext
			} else {
				slog.Error("failed to decrypt fallback AI provider API key", "error", err)
			}
		}
		out = append(out, ai.ProviderConfig{Name: provider, Model: model, APIKey: apiKey, BaseURL: baseURL})
	}
	return out
}

// GetAIConfig handles GET /api/ai/config.
func (h *Handler) GetAIConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.LoadAIConfig(r.Context())
	httputil.JSON(w, http.StatusOK, map[string]any{
		"enabled":       cfg.Enabled,
		"provider":      cfg.Provider,
		"model":         cfg.Model,
		"baseUrl":       cfg.BaseURL,
		"mode":          cfg.Mode,
		"embedProvider": cfg.EmbedProvider,
		"embedModel":    cfg.EmbedModel,
		"embedBaseUrl":  cfg.EmbedBaseURL,
	})
}

// UpdateAIConfig handles PUT /api/ai/config.
func (h *Handler) UpdateAIConfig(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[struct {
		Enabled       *bool   `json:"enabled"`
		Provider      *string `json:"provider"`
		Model         *string `json:"model"`
		APIKey        *string `json:"apiKey"`
		BaseURL       *string `json:"baseUrl"`
		Mode          *string `json:"mode"`
		EmbedProvider *string `json:"embedProvider"`
		EmbedModel    *string `json:"embedModel"`
		EmbedAPIKey   *string `json:"embedApiKey"`
		EmbedBaseURL  *string `json:"embedBaseUrl"`
	}](w, r)
	if !ok {
		return
	}

	var args []any
	var parts []string

	if req.Enabled != nil {
		parts = append(parts, "enabled = ?")
		args = append(args, boolToInt(*req.Enabled))
	}
	if req.Provider != nil {
		parts = append(parts, "provider = ?")
		args = append(args, *req.Provider)
	}
	if req.Model != nil {
		parts = append(parts, "model = ?")
		args = append(args, *req.Model)
	}
	if req.APIKey != nil {
		key, err := crypto.DecodeKey(h.Config.Encryption.Key)
		if err != nil {
			slog.Error("Failed to load encryption key", "error", err)
			httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
			return
		}
		encrypted, err := crypto.Encrypt(*req.APIKey, key)
		if err != nil {
			slog.Error("Failed to encrypt API key", "error", err)
			httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
			return
		}
		parts = append(parts, "api_key_encrypted = ?")
		args = append(args, encrypted)
	}
	if req.BaseURL != nil {
		parts = append(parts, "base_url = ?")
		args = append(args, *req.BaseURL)
	}
	if req.Mode != nil {
		parts = append(parts, "mode = ?")
		args = append(args, *req.Mode)
	}
	if req.EmbedProvider != nil {
		parts = append(parts, "embed_provider = ?")
		args = append(args, *req.EmbedProvider)
	}
	if req.EmbedModel != nil {
		parts = append(parts, "embed_model = ?")
		args = append(args, *req.EmbedModel)
	}
	if req.EmbedAPIKey != nil {
		key, err := crypto.DecodeKey(h.Config.Encryption.Key)
		if err != nil {
			slog.Error("Failed to load encryption key", "error", err)
			httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
			return
		}
		encrypted, err := crypto.Encrypt(*req.EmbedAPIKey, key)
		if err != nil {
			slog.Error("Failed to encrypt embed API key", "error", err)
			httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
			return
		}
		parts = append(parts, "embed_api_key_encrypted = ?")
		args = append(args, encrypted)
	}
	if req.EmbedBaseURL != nil {
		parts = append(parts, "embed_base_url = ?")
		args = append(args, *req.EmbedBaseURL)
	}

	if len(parts) == 0 {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "No fields to update")
		return
	}

	query := "UPDATE ai_config SET " + strings.Join(parts, ", ") + " WHERE id = 1"

	_, err := h.Store.DB().ExecContext(r.Context(), query, args...)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	h.GetAIConfig(w, r)
}

// TestAIConfig handles POST /api/ai/config/test.
func (h *Handler) TestAIConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.LoadAIConfig(r.Context())
	if !cfg.Enabled || cfg.Provider == "" {
		httputil.JSON(w, http.StatusOK, map[string]any{"ok": false, "message": "AI module is not enabled"})
		return
	}

	provider, err := h.AI.Get(cfg.Provider, map[string]any{
		"apiKey": cfg.APIKey, "model": cfg.Model, "baseUrl": cfg.BaseURL,
	})
	if err != nil {
		httputil.JSON(w, http.StatusOK, map[string]any{"ok": false, "message": err.Error()})
		return
	}

	start := time.Now()
	_, err = provider.Suggest(r.Context(), &ai.SuggestRequest{
		SystemPrompt: "Reply with the single word OK.",
		UserPrompt:   "ping",
		MaxTokens:    8,
	})
	latencyMs := int(time.Since(start).Milliseconds())
	if err != nil {
		httputil.JSON(w, http.StatusOK, map[string]any{"ok": false, "message": err.Error(), "latencyMs": latencyMs})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]any{"ok": true, "message": "Connection successful", "latencyMs": latencyMs})
}

// fallbackProviderDoc is one entry in the ordered fallback chain, as
// exchanged with the client. APIKey is write-only (never returned).
type fallbackProviderDoc struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	APIKey   string `json:"apiKey,omitempty"`
	BaseURL  string `json:"baseUrl"`
}

// GetAIFallbackProviders handles GET /api/ai/config/fallback-providers:
// the ordered list of extra providers tried after the primary one fails.
func (h *Handler) GetAIFallbackProviders(w http.ResponseWriter, r *http.Request) {
	providers := h.loadFallbackProviders(r.Context())
	out := make([]fallbackProviderDoc, len(providers))
	for i, p := range providers {
		out[i] = fallbackProviderDoc{Provider: p.Name, Model: p.Model, BaseURL: p.BaseURL}
	}
	httputil.JSON(w, http.StatusOK, out)
}

// UpdateAIFallbackProviders handles PUT /api/ai/config/fallback-providers:
// replaces the whole ordered fallback list (priority 2+) in one call.
func (h *Handler) UpdateAIFallbackProviders(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[[]fallbackProviderDoc](w, r)
	if !ok {
		return
	}

	for _, p := range req {
		if p.Provider == "" {
			httputil.Error(w, http.StatusBadRequest, "invalid_request", "provider is required for every fallback entry")
			return
		}
	}

	var key []byte
	if len(req) > 0 {
		var err error
		key, err = crypto.DecodeKey(h.Config.Encryption.Key)
		if err != nil {
			slog.Error("Failed to load encryption key", "error", err)
			httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
			return
		}
	}

	err := h.Store.WithinTransaction(r.Context(), func(tx *store.Store) error {
		if _, err := tx.DB().ExecContext(r.Context(), `DELETE FROM ai_config_providers WHERE config_id = 1`); err != nil {
			return err
		}
		for i, p := range req {
			encryptedKey := ""
			if p.APIKey != "" {
				var err error
				encryptedKey, err = crypto.Encrypt(p.APIKey, key)
				if err != nil {
					return fmt.Errorf("encrypt fallback API key: %w", err)
				}
			}
			_, err := tx.DB().ExecContext(r.Context(), `
				INSERT INTO ai_config_providers (id, config_id, priority, provider, model, api_key_encrypted, base_url)
				VALUES (?, 1, ?, ?, ?, ?, ?)
			`, uuid.New().String(), i+2, p.Provider, p.Model, encryptedKey, p.BaseURL)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	h.GetAIFallbackProviders(w, r)
}

// notificationConfigDoc mirrors the NotificationConfig OpenAPI schema.
type notificationConfigDoc struct {
	Channels []map[string]any `json:"channels"`
	Routing  []map[string]any `json:"routing"`
}

func (h *Handler) loadNotificationConfig(ctx context.Context) notificationConfigDoc {
	cfg := notificationConfigDoc{Channels: []map[string]any{}, Routing: []map[string]any{}}
	var raw string
	if err := h.Store.DB().QueryRowContext(ctx, `SELECT config_json FROM notification_config WHERE id = 1`).Scan(&raw); err == nil && raw != "" {
		if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
			slog.Warn("settings: invalid stored notification config, using defaults", "error", err)
		}
		if cfg.Channels == nil {
			cfg.Channels = []map[string]any{}
		}
		if cfg.Routing == nil {
			cfg.Routing = []map[string]any{}
		}
	}
	return cfg
}

// GetNotificationsConfig handles GET /api/notifications/config. Signing secrets are write-only:
// the encrypted value is replaced by a secretSet flag.
func (h *Handler) GetNotificationsConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.loadNotificationConfig(r.Context())
	for _, ch := range cfg.Channels {
		conf, _ := ch["config"].(map[string]any)
		if conf == nil {
			continue
		}
		if enc, _ := conf["secretEncrypted"].(string); enc != "" {
			conf["secretSet"] = true
		}
		delete(conf, "secretEncrypted")
	}
	httputil.JSON(w, http.StatusOK, cfg)
}

// applyChannelSecrets turns each channel's write-only config.secret into config.secretEncrypted.
// An absent secret keeps the previously stored one for that channel type; an empty string clears it.
func (h *Handler) applyChannelSecrets(ctx context.Context, cfg *notificationConfigDoc) error {
	prev := map[string]string{}
	for _, ch := range h.loadNotificationConfig(ctx).Channels {
		typ, _ := ch["type"].(string)
		if conf, _ := ch["config"].(map[string]any); conf != nil {
			prev[typ], _ = conf["secretEncrypted"].(string)
		}
	}
	var key []byte
	for _, ch := range cfg.Channels {
		conf, _ := ch["config"].(map[string]any)
		if conf == nil {
			continue
		}
		typ, _ := ch["type"].(string)
		delete(conf, "secretSet")
		delete(conf, "secretEncrypted")
		secret, provided := conf["secret"].(string)
		delete(conf, "secret")
		switch {
		case !provided:
			if prev[typ] != "" {
				conf["secretEncrypted"] = prev[typ]
			}
		case secret != "":
			if key == nil {
				var err error
				if key, err = crypto.DecodeKey(h.Config.Encryption.Key); err != nil {
					return fmt.Errorf("load encryption key: %w", err)
				}
			}
			enc, err := crypto.Encrypt(secret, key)
			if err != nil {
				return fmt.Errorf("encrypt signing secret: %w", err)
			}
			conf["secretEncrypted"] = enc
		}
	}
	return nil
}

// UpdateNotificationsConfig handles PUT /api/notifications/config.
func (h *Handler) UpdateNotificationsConfig(w http.ResponseWriter, r *http.Request) {
	cfg, ok := httputil.DecodeJSON[notificationConfigDoc](w, r)
	if !ok {
		return
	}
	if cfg.Channels == nil {
		cfg.Channels = []map[string]any{}
	}
	if cfg.Routing == nil {
		cfg.Routing = []map[string]any{}
	}

	if err := h.applyChannelSecrets(r.Context(), &cfg); err != nil {
		httputil.Errorf(w, err)
		return
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if _, err := h.Store.DB().ExecContext(r.Context(), `UPDATE notification_config SET config_json = ? WHERE id = 1`, string(data)); err != nil {
		httputil.Errorf(w, err)
		return
	}

	h.GetNotificationsConfig(w, r)
}

// TestNotificationsConfig handles POST /api/notifications/config/test.
// in_app has no external dependency and always succeeds; smtp/webhook are
// stubs (see internal/notifications.Dispatcher) that log and report success.
func (h *Handler) TestNotificationsConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Channel == "" {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "channel is required")
		return
	}
	slog.Info("test notification requested", "channel", logsafe.Sanitize(req.Channel))
	httputil.JSON(w, http.StatusOK, map[string]any{"ok": true, "message": "Test notification sent"})
}

// GetDecryptedAPIKey reads the stored encrypted API key and returns it decrypted.
// Returns an empty string if no key is stored or if decryption fails.
func (h *Handler) GetDecryptedAPIKey() string {
	var encrypted string
	err := h.Store.DB().QueryRow(`SELECT api_key_encrypted FROM ai_config WHERE id = 1`).Scan(&encrypted)
	if err != nil || encrypted == "" {
		return ""
	}

	key, err := crypto.DecodeKey(h.Config.Encryption.Key)
	if err != nil {
		slog.Error("Failed to load encryption key", "error", err)
		return ""
	}
	plaintext, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		// Legacy path: the row may predate WISELABZ_ENCRYPTION_KEY, when the
		// key was derived from the JWT signing secret instead.
		legacyKey := crypto.DeriveKey(h.Config.Auth.Secret) //nolint:staticcheck // intentional legacy fallback for rows encrypted before WISELABZ_ENCRYPTION_KEY
		plaintext, err = crypto.Decrypt(encrypted, legacyKey)
		if err != nil {
			slog.Error("Failed to decrypt stored API key", "error", err)
			return ""
		}
		slog.Warn("Decrypted API key using legacy derived key; re-save the AI config to migrate it to WISELABZ_ENCRYPTION_KEY")
	}
	return plaintext
}

// GetDecryptedEmbedAPIKey reads the stored encrypted embedding-provider API
// key and returns it decrypted. Returns an empty string if no key is stored
// or if decryption fails.
func (h *Handler) GetDecryptedEmbedAPIKey() string {
	var encrypted string
	err := h.Store.DB().QueryRow(`SELECT embed_api_key_encrypted FROM ai_config WHERE id = 1`).Scan(&encrypted)
	if err != nil || encrypted == "" {
		return ""
	}

	key, err := crypto.DecodeKey(h.Config.Encryption.Key)
	if err != nil {
		slog.Error("Failed to load encryption key", "error", err)
		return ""
	}
	plaintext, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		slog.Error("Failed to decrypt stored embed API key", "error", err)
		return ""
	}
	return plaintext
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
