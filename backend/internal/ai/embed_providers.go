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
)

// ollamaEmbedder talks to a local (or self-hosted) Ollama server's batch
// embeddings endpoint.
type ollamaEmbedder struct {
	baseURL string
	model   string
	client  *http.Client
}

// RegisterOllamaEmbedder registers the "ollama" embedder, the default local
// embedding backend.
func RegisterOllamaEmbedder(r *EmbedRegistry) {
	r.Register("ollama", func(config map[string]any) (Embedder, error) {
		baseURL, _ := config["baseUrl"].(string)
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		model, _ := config["model"].(string)
		if model == "" {
			model = "nomic-embed-text"
		}
		return &ollamaEmbedder{
			baseURL: strings.TrimRight(baseURL, "/"),
			model:   model,
			client:  &http.Client{Timeout: 60 * time.Second, CheckRedirect: noRedirect},
		}, nil
	})
}

func (e *ollamaEmbedder) Name() string { return "ollama" }

func (e *ollamaEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	body, err := json.Marshal(map[string]any{"model": e.model, "input": texts})
	if err != nil {
		return nil, fmt.Errorf("marshal embed request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/api/embed", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build embed request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("embed request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, connector.MaxResponseBytes))
		return nil, fmt.Errorf("embedder returned %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var out struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.NewDecoder(connector.LimitedBody(resp.Body)).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode embed response: %w", err)
	}
	if len(out.Embeddings) != len(texts) {
		return nil, fmt.Errorf("embedder returned %d vectors for %d inputs", len(out.Embeddings), len(texts))
	}
	return out.Embeddings, nil
}

// openAIEmbedder talks to OpenAI's (or an OpenAI-compatible) embeddings endpoint.
type openAIEmbedder struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// RegisterOpenAIEmbedder registers the "openai" embedder.
func RegisterOpenAIEmbedder(r *EmbedRegistry) {
	r.Register("openai", func(config map[string]any) (Embedder, error) {
		baseURL, _ := config["baseUrl"].(string)
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		apiKey, _ := config["apiKey"].(string)
		model, _ := config["model"].(string)
		if model == "" {
			model = "text-embedding-3-small"
		}
		return &openAIEmbedder{
			baseURL: strings.TrimRight(baseURL, "/"),
			apiKey:  apiKey,
			model:   model,
			client:  &http.Client{Timeout: 60 * time.Second, CheckRedirect: noRedirect},
		}, nil
	})
}

func (e *openAIEmbedder) Name() string { return "openai" }

func (e *openAIEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	body, err := json.Marshal(map[string]any{"model": e.model, "input": texts})
	if err != nil {
		return nil, fmt.Errorf("marshal embed request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, e.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build embed request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if e.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+e.apiKey)
	}

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("embed request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, connector.MaxResponseBytes))
		return nil, fmt.Errorf("embedder returned %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var out struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
			Index     int       `json:"index"`
		} `json:"data"`
	}
	if err := json.NewDecoder(connector.LimitedBody(resp.Body)).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode embed response: %w", err)
	}
	if len(out.Data) != len(texts) {
		return nil, fmt.Errorf("embedder returned %d vectors for %d inputs", len(out.Data), len(texts))
	}
	vectors := make([][]float32, len(texts))
	for _, d := range out.Data {
		if d.Index < 0 || d.Index >= len(vectors) {
			return nil, fmt.Errorf("embedder returned out-of-range index %d", d.Index)
		}
		vectors[d.Index] = d.Embedding
	}
	return vectors, nil
}
