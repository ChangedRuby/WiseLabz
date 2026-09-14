package ai

import (
	"context"
	"testing"
)

// testProvider is a minimal test provider for registry testing.
type testProvider struct {
	name string
}

func (t *testProvider) Name() string { return t.name }

func (t *testProvider) Suggest(_ context.Context, _ *SuggestRequest) (string, error) {
	return "test suggestion", nil
}

func (t *testProvider) SuggestStream(_ context.Context, _ *SuggestRequest) (<-chan SuggestChunk, error) {
	ch := make(chan SuggestChunk, 1)
	go func() {
		defer close(ch)
		ch <- SuggestChunk{ContentDelta: "test", Done: true}
	}()
	return ch, nil
}

func TestRegistryList(t *testing.T) {
	t.Run("empty registry", func(t *testing.T) {
		r := NewRegistry()
		names := r.List()
		if len(names) != 0 {
			t.Errorf("empty registry returned %d names, want 0", len(names))
		}
	})

	t.Run("multiple providers", func(t *testing.T) {
		r := NewRegistry()
		r.Register("provider1", func(_ map[string]any) (Provider, error) {
			return &testProvider{name: "provider1"}, nil
		})
		r.Register("provider2", func(_ map[string]any) (Provider, error) {
			return &testProvider{name: "provider2"}, nil
		})

		names := r.List()
		if len(names) != 2 {
			t.Errorf("registry returned %d names, want 2", len(names))
		}
	})
}

func TestRegistryGet(t *testing.T) {
	t.Run("provider not found", func(t *testing.T) {
		r := NewRegistry()
		_, err := r.Get("unknown", map[string]any{})
		if err == nil {
			t.Errorf("Get unknown provider should return error")
		}
	})

	t.Run("provider found", func(t *testing.T) {
		r := NewRegistry()
		r.Register("test", func(_ map[string]any) (Provider, error) {
			return &testProvider{name: "test"}, nil
		})

		p, err := r.Get("test", map[string]any{})
		if err != nil {
			t.Errorf("Get test provider failed: %v", err)
		}
		if p == nil {
			t.Errorf("Get test provider returned nil")
		}
		if p.Name() != "test" {
			t.Errorf("provider name = %q, want test", p.Name())
		}
	})
}

func TestStubProviderName(t *testing.T) {
	p := &StubProvider{}
	if p.Name() != "stub" {
		t.Errorf("Name() = %q, want stub", p.Name())
	}
}

func TestStubProviderSuggest(t *testing.T) {
	p := &StubProvider{}
	result, err := p.Suggest(context.Background(), &SuggestRequest{})
	if err != nil {
		t.Errorf("Suggest() returned error: %v", err)
	}
	if result == "" {
		t.Errorf("Suggest() returned empty string")
	}
}

func TestStubProviderSuggestStream(t *testing.T) {
	p := &StubProvider{}
	ch, err := p.SuggestStream(context.Background(), &SuggestRequest{})
	if err != nil {
		t.Errorf("SuggestStream() returned error: %v", err)
	}
	if ch == nil {
		t.Errorf("SuggestStream() returned nil channel")
	}

	chunk := <-ch
	if !chunk.Done {
		t.Errorf("chunk.Done = %v, want true", chunk.Done)
	}
	if chunk.ContentDelta == "" {
		t.Errorf("chunk.ContentDelta is empty")
	}
}
