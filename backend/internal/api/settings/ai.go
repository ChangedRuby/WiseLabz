package settings

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/crypto"
	"github.com/WiseLabz/wiselabz/internal/httputil"
)

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

// encryptAIKey encrypts an AI provider API key with the instance encryption
// key. On failure it logs (encryptLogMsg for the encrypt step), writes a 500
// response and returns ok=false.
func (h *Handler) encryptAIKey(w http.ResponseWriter, plaintext, encryptLogMsg string) (string, bool) {
	key, err := crypto.DecodeKey(h.Config.Encryption.Key)
	if err != nil {
		slog.Error("Failed to load encryption key", "error", err)
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
		return "", false
	}
	encrypted, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		slog.Error(encryptLogMsg, "error", err)
		httputil.Error(w, http.StatusInternalServerError, "internal_error", "Failed to encrypt API key")
		return "", false
	}
	return encrypted, true
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
		encrypted, ok := h.encryptAIKey(w, *req.APIKey, "Failed to encrypt API key")
		if !ok {
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
		encrypted, ok := h.encryptAIKey(w, *req.EmbedAPIKey, "Failed to encrypt embed API key")
		if !ok {
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

	if !h.persistAIConfigUpdate(w, r, parts, args) {
		return
	}

	h.GetAIConfig(w, r)
}

func (h *Handler) persistAIConfigUpdate(w http.ResponseWriter, r *http.Request, parts []string, args []any) bool {
	query := "UPDATE ai_config SET " + strings.Join(parts, ", ") + " WHERE id = 1"
	_, err := h.Store.DB().ExecContext(r.Context(), query, args...)
	if err != nil {
		httputil.Errorf(w, err)
		return false
	}
	return true
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
		//lint:ignore SA1019 intentional legacy fallback for rows encrypted before WISELABZ_ENCRYPTION_KEY
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
