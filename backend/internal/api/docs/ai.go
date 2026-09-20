package docs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/ai"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/chat"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// syncDocEmbeddings recomputes chat-retrieval embeddings for a doc so "ask
// your lab" search stays current after every generate/save. Best-effort: an
// embedding-backend failure is logged, not surfaced, so it never blocks the
// doc write it's attached to.
func (h *Handler) syncDocEmbeddings(ctx context.Context, docID, content string) {
	if h.Embed == nil || docID == "" {
		return
	}
	cfg := h.Settings.LoadAIConfig(ctx)
	if !cfg.Enabled || cfg.EmbedProvider == "" {
		return
	}
	embedder, err := h.Embed.Get(cfg.EmbedProvider, map[string]any{
		"apiKey": cfg.EmbedAPIKey, "model": cfg.EmbedModel, "baseUrl": cfg.EmbedBaseURL,
	})
	if err != nil {
		slog.Warn("chat: embedding backend unavailable, skipping doc embedding sync", "docId", docID, "error", err)
		return
	}
	if err := chat.SyncDocEmbeddings(ctx, h.Store, embedder, cfg.EmbedModel, docID, content); err != nil {
		slog.Warn("chat: failed to sync doc embeddings", "docId", docID, "error", err)
	}
}

// aiSuggestTimeout bounds detached AI suggestion calls.
const aiSuggestTimeout = 2 * time.Minute

// AISuggest handles POST /api/docs/{id}/ai-suggest.
// Batched (non-streaming) suggestion: the full result is delivered over the
// doc.ai_suggestion WS event, correlated by the returned requestId.
func (h *Handler) AISuggest(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")

	var req struct {
		Prompt    string `json:"prompt"`
		Selection string `json:"selection"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	d, err := h.Store.GetDoc(r.Context(), docID)
	if errors.Is(err, store.ErrNotFound) {
		httputil.Error(w, http.StatusNotFound, "not_found", "Doc not found")
		return
	}
	if err != nil {
		httputil.Errorf(w, err)
		return
	}
	if !h.requireDocOperator(w, r, d.ServiceID) {
		return
	}

	cfg := h.Settings.LoadAIConfig(r.Context())
	if !cfg.Enabled {
		httputil.Error(w, http.StatusConflict, "ai_disabled", "AI module is not enabled")
		return
	}

	provider, err := h.AI.Get(cfg.Provider, map[string]any{
		"apiKey": cfg.APIKey, "model": cfg.Model, "baseUrl": cfg.BaseURL,
	})
	if err != nil {
		httputil.Error(w, http.StatusConflict, "ai_disabled", fmt.Sprintf("AI provider unavailable: %v", err))
		return
	}

	prompt := req.Prompt
	if req.Selection != "" {
		prompt = fmt.Sprintf("%s\n\nSelected text:\n%s", prompt, req.Selection)
	}

	userID := auth.UserIDFromContext(r.Context())
	requestID := uuid.New().String()

	// The suggestion is delivered over WS after the 202 returns, so it must
	// outlive the request; bound it so a hung provider can't leak the goroutine.
	aiCtx, cancelAI := context.WithTimeout(context.WithoutCancel(r.Context()), aiSuggestTimeout)
	go func() {
		defer cancelAI()
		content, err := provider.Suggest(aiCtx, &ai.SuggestRequest{
			SystemPrompt: "You maintain internal infrastructure documentation. Suggest an improved version of the document based on the request.",
			UserPrompt:   prompt,
			DocContent:   d.Content,
		})
		payload := map[string]any{"docId": docID, "requestId": requestID}
		if err != nil {
			payload["status"] = "error"
			payload["error"] = err.Error()
		} else {
			payload["status"] = "complete"
			payload["fullContent"] = content
		}
		if h.WSHub != nil {
			h.WSHub.BroadcastToUser(userID, ws.EventDocAISuggestion, payload)
		}
	}()

	httputil.JSON(w, http.StatusAccepted, map[string]any{"requestId": requestID})
}
