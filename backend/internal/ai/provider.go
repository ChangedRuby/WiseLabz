// Package ai provides AI provider abstraction for doc suggestions.
package ai

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
)

// Provider is the interface for AI suggestion providers.
type Provider interface {
	Name() string
	Suggest(_ context.Context, _ *SuggestRequest) (string, error)
	SuggestStream(_ context.Context, _ *SuggestRequest) (<-chan SuggestChunk, error)
}

// SuggestRequest contains the prompt and context for an AI suggestion.
type SuggestRequest struct {
	SystemPrompt string `json:"systemPrompt"`
	UserPrompt   string `json:"userPrompt"`
	DocContent   string `json:"docContent"`
	MaxTokens    int    `json:"maxTokens"`
}

// SuggestChunk is a streaming response chunk from an AI provider.
type SuggestChunk struct {
	ContentDelta string `json:"contentDelta"`
	Done         bool   `json:"done"`
	Error        string `json:"error,omitempty"`
}

// Registry maps provider names to factory functions.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]func(config map[string]any) (Provider, error)
}

// NewRegistry creates a new AI provider registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]func(config map[string]any) (Provider, error)),
	}
}

// Register adds a provider factory to the registry.
func (r *Registry) Register(name string, factory func(config map[string]any) (Provider, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = factory
}

// Get returns a provider instance by name.
func (r *Registry) Get(name string, config map[string]any) (Provider, error) {
	r.mu.RLock()
	factory, ok := r.factories[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown AI provider: %q", name)
	}
	return factory(config)
}

// List returns registered provider names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var names []string
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// StatusError wraps an HTTP status code an AI provider returned, so callers
// (fallback routing) can tell a retryable failure (429/5xx) from others.
type StatusError struct {
	Code int
	Body string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("ai provider returned %d: %s", e.Code, e.Body)
}

// isRetryable reports whether err should advance to the next configured
// provider: HTTP 429/5xx, or a timeout.
func isRetryable(err error) bool {
	var se *StatusError
	if errors.As(err, &se) {
		return se.Code == 429 || se.Code >= 500
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// ProviderConfig names one entry in an ordered fallback chain: which
// registered provider to instantiate and its connection settings.
type ProviderConfig struct {
	Name    string
	Model   string
	APIKey  string
	BaseURL string
}

// SuggestResult carries a successful Suggest call's answer plus which
// provider produced it, so callers can surface fallback provenance.
type SuggestResult struct {
	Content      string
	Provider     string
	FallbackUsed bool
}

// SuggestWithFallback tries each configured provider in priority order,
// advancing to the next on a retryable failure (429, 5xx, timeout). If every
// provider fails, it returns a combined error naming each attempted provider.
func SuggestWithFallback(ctx context.Context, registry *Registry, configs []ProviderConfig, req *SuggestRequest) (*SuggestResult, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no AI providers configured")
	}

	var errs []error
	for i, c := range configs {
		provider, err := registry.Get(c.Name, map[string]any{
			"apiKey": c.APIKey, "model": c.Model, "baseUrl": c.BaseURL,
		})
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", c.Name, err))
			continue
		}

		content, err := provider.Suggest(ctx, req)
		if err == nil {
			return &SuggestResult{Content: content, Provider: provider.Name(), FallbackUsed: i > 0}, nil
		}
		errs = append(errs, fmt.Errorf("%s: %w", c.Name, err))
		if !isRetryable(err) {
			break
		}
	}
	return nil, errors.Join(errs...)
}

// StubProvider is a stub AI provider for when AI is not configured.
type StubProvider struct{}

// Name returns "stub".
func (s *StubProvider) Name() string { return "stub" }

// Suggest returns a stub message.
func (s *StubProvider) Suggest(_ context.Context, _ *SuggestRequest) (string, error) {
	return "AI suggestions are not configured. Please enable an AI provider in settings.", nil
}

// SuggestStream returns a single chunk with the stub message.
func (s *StubProvider) SuggestStream(_ context.Context, _ *SuggestRequest) (<-chan SuggestChunk, error) {
	ch := make(chan SuggestChunk, 1)
	go func() {
		defer close(ch)
		ch <- SuggestChunk{
			ContentDelta: "AI suggestions are not configured. Please enable an AI provider in settings.",
			Done:         true,
		}
	}()
	return ch, nil
}

// noRedirect blocks redirects so a configured base URL cannot bounce
// requests (and credentials) to another host.
func noRedirect(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
