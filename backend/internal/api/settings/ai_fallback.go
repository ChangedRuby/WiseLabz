package settings

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/crypto"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

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

	for i, p := range req {
		if p.Provider == "" {
			httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_request", "provider is required for every fallback entry",
				[]httputil.FieldError{{Field: fmt.Sprintf("[%d].provider", i), Msg: "is required"}})
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
