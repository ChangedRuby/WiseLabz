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
)

// claudeProvider talks to Anthropic's Messages API.
type claudeProvider struct {
	name    string
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// RegisterClaude registers the Claude provider under the "claude" name.
func RegisterClaude(r *Registry) {
	factory := func(config map[string]any) (Provider, error) {
		baseURL, _ := config["baseUrl"].(string)
		if baseURL == "" {
			baseURL = "https://api.anthropic.com"
		}
		apiKey, _ := config["apiKey"].(string)
		model, _ := config["model"].(string)
		if model == "" {
			// ponytail: claude-sonnet-5 is current default; update if a more recent stable model exists
			model = "claude-sonnet-5"
		}
		return &claudeProvider{
			name:    "claude",
			baseURL: strings.TrimRight(baseURL, "/"),
			apiKey:  apiKey,
			model:   model,
			client:  &http.Client{Timeout: 60 * time.Second},
		}, nil
	}
	r.Register("claude", factory)
}

func (p *claudeProvider) Name() string { return p.name }

// Suggest sends a single (non-streaming) message to the Claude API.
func (p *claudeProvider) Suggest(ctx context.Context, req *SuggestRequest) (string, error) {
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	body := map[string]any{
		"model":      p.model,
		"max_tokens": maxTokens,
		"system":     req.SystemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": req.UserPrompt},
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal ai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("build ai request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ai request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", &StatusError{Code: resp.StatusCode, Body: strings.TrimSpace(string(b))}
	}

	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode ai response: %w", err)
	}
	if len(out.Content) == 0 {
		return "", fmt.Errorf("ai provider returned no content")
	}

	var result strings.Builder
	for _, block := range out.Content {
		if block.Type == "text" {
			result.WriteString(block.Text)
		}
	}
	return result.String(), nil
}

// SuggestStream wraps Suggest into a single completed chunk so the Provider
// interface is satisfied without a separate streaming code path.
func (p *claudeProvider) SuggestStream(ctx context.Context, req *SuggestRequest) (<-chan SuggestChunk, error) {
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
