package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/httpx"
)

// openAICompatibleProvider talks to any OpenAI-compatible chat completions
// endpoint (OpenAI itself, or a self-hosted/Ollama server exposing the same
// /chat/completions shape).
type openAICompatibleProvider struct {
	name    string
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// RegisterOpenAICompatible registers the OpenAI-compatible provider under the
// "openai" and "ollama" names (both expose the same /chat/completions API).
func RegisterOpenAICompatible(r *Registry) {
	factory := func(name string) func(map[string]any) (Provider, error) {
		return func(config map[string]any) (Provider, error) {
			baseURL, _ := config["baseUrl"].(string)
			if baseURL == "" {
				baseURL = "https://api.openai.com/v1"
			}
			apiKey, _ := config["apiKey"].(string)
			model, _ := config["model"].(string)
			if model == "" {
				model = "gpt-4o-mini"
			}
			return &openAICompatibleProvider{
				name:    name,
				baseURL: strings.TrimRight(baseURL, "/"),
				apiKey:  apiKey,
				model:   model,
				client:  httpx.NewClient(httpx.Options{Timeout: 60 * time.Second}),
			}, nil
		}
	}
	r.Register("openai", factory("openai"))
	r.Register("ollama", factory("ollama"))
}

func (p *openAICompatibleProvider) Name() string { return p.name }

// Suggest sends a single (non-streaming) chat completion request.
func (p *openAICompatibleProvider) Suggest(ctx context.Context, req *SuggestRequest) (string, error) {
	body := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": req.SystemPrompt},
			{"role": "user", "content": req.UserPrompt},
		},
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal ai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("build ai request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ai request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, connector.MaxResponseBytes))
		return "", &StatusError{Code: resp.StatusCode, Body: strings.TrimSpace(string(b))}
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(connector.LimitedBody(resp.Body)).Decode(&out); err != nil {
		return "", fmt.Errorf("decode ai response: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("ai provider returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}

// SuggestStream wraps Suggest into a single completed chunk so the Provider
// interface is satisfied without a separate streaming code path.
func (p *openAICompatibleProvider) SuggestStream(ctx context.Context, req *SuggestRequest) (<-chan SuggestChunk, error) {
	ch := make(chan SuggestChunk, 1)
	go func() {
		defer close(ch)
		content, err := p.Suggest(ctx, req)
		if err != nil {
			ch <- SuggestChunk{Done: true, Error: err.Error()}
			return
		}
		ch <- SuggestChunk{ContentDelta: content, Done: true}
	}()
	return ch, nil
}
