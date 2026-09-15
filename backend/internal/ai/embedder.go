package ai

import (
	"context"
	"fmt"
	"sync"
)

// Embedder is the interface for text-embedding providers. Separate from
// Provider because not every chat provider supports embeddings (Claude does
// not), so the embedding backend is configured independently.
type Embedder interface {
	Name() string
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}

// EmbedRegistry maps embedder names to factory functions, mirroring Registry.
type EmbedRegistry struct {
	mu        sync.RWMutex
	factories map[string]func(config map[string]any) (Embedder, error)
}

// NewEmbedRegistry creates a new embedder registry.
func NewEmbedRegistry() *EmbedRegistry {
	return &EmbedRegistry{
		factories: make(map[string]func(config map[string]any) (Embedder, error)),
	}
}

// Register adds an embedder factory to the registry.
func (r *EmbedRegistry) Register(name string, factory func(config map[string]any) (Embedder, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = factory
}

// Get returns an embedder instance by name.
func (r *EmbedRegistry) Get(name string, config map[string]any) (Embedder, error) {
	r.mu.RLock()
	factory, ok := r.factories[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown embedder: %q", name)
	}
	return factory(config)
}
