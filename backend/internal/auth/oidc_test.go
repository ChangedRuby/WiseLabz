package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsInitializedBeforeAndAfter(t *testing.T) {
	provider := &OIDCProvider{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		IssuerURL:    "https://example.com",
	}

	// Before initialization
	if provider.IsInitialized() {
		t.Error("IsInitialized() = true before Initialize, want false")
	}

	// After initialization with valid server
	server := newMockOIDCServer(t)
	defer server.Close()

	provider.IssuerURL = server.URL
	err := provider.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize() error: %v", err)
	}

	// After successful initialization
	if !provider.IsInitialized() {
		t.Error("IsInitialized() = false after successful Initialize, want true")
	}
}

func TestAuthURLBeforeInitialization(t *testing.T) {
	provider := &OIDCProvider{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		IssuerURL:    "https://example.com",
	}

	// AuthURL before initialization will panic because oauth2 is nil
	defer func() {
		if r := recover(); r == nil {
			t.Error("AuthURL() before initialization should panic, but didn't")
		}
	}()
	_ = provider.AuthURL("state123", "nonce123", "http://localhost/callback")
}

func TestExtractGroups(t *testing.T) {
	for name, tc := range map[string]struct {
		value any
		want  []string
	}{
		"array":        {value: []any{"admins", "readers", ""}, want: []string{"admins", "readers"}},
		"string":       {value: "admins", want: []string{"admins"}},
		"invalid":      {value: 42, want: nil},
		"empty string": {value: "", want: nil},
	} {
		t.Run(name, func(t *testing.T) {
			got := extractGroups(tc.value)
			if len(got) != len(tc.want) {
				t.Fatalf("extractGroups(%#v) = %#v, want %#v", tc.value, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("extractGroups(%#v) = %#v, want %#v", tc.value, got, tc.want)
				}
			}
		})
	}
}

func TestInitializeSuccess(t *testing.T) {
	server := newMockOIDCServer(t)
	defer server.Close()

	provider := &OIDCProvider{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		IssuerURL:    server.URL,
	}

	err := provider.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize() error: %v", err)
	}

	if !provider.IsInitialized() {
		t.Error("IsInitialized() = false after successful Initialize, want true")
	}

	if provider.provider == nil {
		t.Error("provider.provider is nil after Initialize")
	}
	if provider.oauth2 == nil {
		t.Error("provider.oauth2 is nil after Initialize")
	}
}

func TestInitializeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	provider := &OIDCProvider{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		IssuerURL:    server.URL,
	}

	err := provider.Initialize(context.Background())
	if err == nil {
		t.Error("Initialize() error = nil, want error")
	}

	if provider.IsInitialized() {
		t.Error("IsInitialized() = true after failed Initialize, want false")
	}
}

func TestInitializeInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("invalid json {")); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	provider := &OIDCProvider{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		IssuerURL:    server.URL,
	}

	err := provider.Initialize(context.Background())
	if err == nil {
		t.Error("Initialize() error = nil for invalid JSON, want error")
	}

	if provider.IsInitialized() {
		t.Error("IsInitialized() = true after Initialize failure, want false")
	}
}

func TestAuthURLAfterInitialization(t *testing.T) {
	server := newMockOIDCServer(t)
	defer server.Close()

	provider := &OIDCProvider{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		IssuerURL:    server.URL,
	}

	err := provider.Initialize(context.Background())
	if err != nil {
		t.Fatalf("Initialize() error: %v", err)
	}

	redirectURL := "http://localhost:8080/callback"
	authURL := provider.AuthURL("state-abc123", "nonce-abc123", redirectURL)

	if authURL == "" {
		t.Error("AuthURL() returned empty string")
	}

	// Check that URL contains expected parameters
	if !strings.Contains(authURL, "client_id=test-client") {
		t.Errorf("AuthURL() missing client_id parameter: %s", authURL)
	}
	if !strings.Contains(authURL, "redirect_uri=") {
		t.Errorf("AuthURL() missing redirect_uri parameter: %s", authURL)
	}
	if !strings.Contains(authURL, "response_type=code") {
		t.Errorf("AuthURL() missing response_type parameter: %s", authURL)
	}
	if !strings.Contains(authURL, "state=state-abc123") {
		t.Errorf("AuthURL() missing state parameter: %s", authURL)
	}
}

// newMockOIDCServer creates a minimal OIDC discovery server for testing.
func newMockOIDCServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/.well-known/openid-configuration" {
			// Return a minimal but valid discovery document
			// httptest typically uses http:// scheme without X-Forwarded headers
			discovery := map[string]interface{}{
				"issuer":                                "http://" + r.Host,
				"authorization_endpoint":                "http://" + r.Host + "/authorize",
				"token_endpoint":                        "http://" + r.Host + "/token",
				"jwks_uri":                              "http://" + r.Host + "/jwks",
				"response_types_supported":              []string{"code"},
				"subject_types_supported":               []string{"public"},
				"id_token_signing_alg_values_supported": []string{"RS256"},
			}
			if err := json.NewEncoder(w).Encode(discovery); err != nil {
				return
			}
			return
		}

		if r.URL.Path == "/jwks" {
			// Return an empty but valid JWKS
			jwks := map[string]interface{}{
				"keys": []interface{}{},
			}
			if err := json.NewEncoder(w).Encode(jwks); err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}
