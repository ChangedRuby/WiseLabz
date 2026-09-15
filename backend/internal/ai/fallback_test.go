package ai

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func registerFailThenSucceed(r *Registry, failErr error) {
	r.Register("bad", func(_ map[string]any) (Provider, error) {
		return &testProvider{name: "bad", err: failErr}, nil
	})
	r.Register("good", func(_ map[string]any) (Provider, error) {
		return &testProvider{name: "good"}, nil
	})
}

func TestSuggestWithFallbackAdvancesOnRetryableError(t *testing.T) {
	r := NewRegistry()
	registerFailThenSucceed(r, &StatusError{Code: 503})
	configs := []ProviderConfig{{Name: "bad"}, {Name: "good"}}

	result, err := SuggestWithFallback(context.Background(), r, configs, &SuggestRequest{})
	if err != nil {
		t.Fatalf("SuggestWithFallback() error = %v", err)
	}
	if result.Provider != "good" || !result.FallbackUsed {
		t.Fatalf("result = %+v, want provider=good fallbackUsed=true", result)
	}
}

func TestSuggestWithFallbackFirstProviderSucceeds(t *testing.T) {
	r := NewRegistry()
	registerFailThenSucceed(r, &StatusError{Code: 503})
	configs := []ProviderConfig{{Name: "good"}, {Name: "bad"}}

	result, err := SuggestWithFallback(context.Background(), r, configs, &SuggestRequest{})
	if err != nil {
		t.Fatalf("SuggestWithFallback() error = %v", err)
	}
	if result.Provider != "good" || result.FallbackUsed {
		t.Fatalf("result = %+v, want provider=good fallbackUsed=false", result)
	}
}

func TestSuggestWithFallbackAllFail(t *testing.T) {
	r := NewRegistry()
	r.Register("bad1", func(_ map[string]any) (Provider, error) {
		return &testProvider{name: "bad1", err: &StatusError{Code: 429}}, nil
	})
	r.Register("bad2", func(_ map[string]any) (Provider, error) {
		return &testProvider{name: "bad2", err: &StatusError{Code: 500}}, nil
	})
	configs := []ProviderConfig{{Name: "bad1"}, {Name: "bad2"}}

	_, err := SuggestWithFallback(context.Background(), r, configs, &SuggestRequest{})
	if err == nil {
		t.Fatal("SuggestWithFallback() error = nil, want combined error")
	}
	if !strings.Contains(err.Error(), "bad1") || !strings.Contains(err.Error(), "bad2") {
		t.Fatalf("error = %v, want both provider names", err)
	}
}

func TestSuggestWithFallbackStopsOnNonRetryableError(t *testing.T) {
	r := NewRegistry()
	r.Register("bad", func(_ map[string]any) (Provider, error) {
		return &testProvider{name: "bad", err: &StatusError{Code: 400}}, nil
	})
	r.Register("good", func(_ map[string]any) (Provider, error) {
		return &testProvider{name: "good"}, nil
	})
	configs := []ProviderConfig{{Name: "bad"}, {Name: "good"}}

	_, err := SuggestWithFallback(context.Background(), r, configs, &SuggestRequest{})
	if err == nil || strings.Contains(err.Error(), "good") {
		t.Fatalf("error = %v, want only bad attempted", err)
	}
}

func TestSuggestWithFallbackNoProviders(t *testing.T) {
	_, err := SuggestWithFallback(context.Background(), NewRegistry(), nil, &SuggestRequest{})
	if err == nil {
		t.Fatal("SuggestWithFallback() error = nil, want error for empty provider list")
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"429", &StatusError{Code: 429}, true},
		{"500", &StatusError{Code: 500}, true},
		{"400", &StatusError{Code: 400}, false},
		{"deadline exceeded", context.DeadlineExceeded, true},
		{"other error", errors.New("boom"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryable(tt.err); got != tt.want {
				t.Errorf("isRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
