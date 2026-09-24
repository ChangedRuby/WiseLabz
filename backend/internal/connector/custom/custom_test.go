package custom

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestValidateCustomURL(t *testing.T) {
	for _, raw := range []string{"", "ftp://example.test", "http:///missing-host"} {
		if err := validateCustomURL(raw); err == nil {
			t.Errorf("validateCustomURL(%q) error = nil", raw)
		}
	}
	if err := validateCustomURL("https://example.test/status"); err != nil {
		t.Fatalf("validateCustomURL(valid): %v", err)
	}
}

func TestFetchUsesConfiguredMethodAndHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.Header.Get("X-Test") != "value" {
			t.Fatalf("request = %s %q", r.Method, r.Header.Get("X-Test"))
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()
	c := &Connector{client: server.Client()}
	snapshot, err := c.Fetch(context.Background(), map[string]any{
		"url": server.URL, "method": http.MethodPatch, "headers": `{"X-Test":"value"}`,
	})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.Metadata["status_code"] != "200" || !strings.Contains(snapshot.Sections[0].Content, `{"status":"ok"}`) {
		t.Fatalf("snapshot = %+v", snapshot)
	}
}

func TestValidateSurfacesAuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()
	c := &Connector{client: server.Client()}
	err := c.Validate(context.Background(), map[string]any{"url": server.URL})
	var authErr *connector.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("Validate() error = %v, want *connector.AuthError", err)
	}
}

func TestValidateAndFetchTimeoutError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := &Connector{client: &http.Client{Timeout: 10 * time.Millisecond}}
	err := c.Validate(context.Background(), map[string]any{"url": server.URL})
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Errorf("Validate() error = %v, want *connector.TimeoutError", err)
	}

	err = c.Validate(context.Background(), map[string]any{"url": server.URL})
	if !errors.As(err, &timeoutErr) {
		t.Errorf("Fetch() error = %v, want *connector.TimeoutError", err)
	}
}

func TestValidateAndFetchServiceUnavailableError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"502 BadGateway", http.StatusBadGateway},
		{"503 ServiceUnavailable", http.StatusServiceUnavailable},
		{"504 GatewayTimeout", http.StatusGatewayTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("unavailable"))
			}))
			defer server.Close()

			c := &Connector{client: server.Client()}
			err := c.Validate(context.Background(), map[string]any{"url": server.URL})
			var unavailErr *connector.ServiceUnavailableError
			if !errors.As(err, &unavailErr) {
				t.Errorf("Validate() error = %v, want *connector.ServiceUnavailableError", err)
			}

			// Fetch doesn't check status codes - it just returns the response as-is
			snap, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
			if err != nil {
				t.Errorf("Fetch() error = %v, want nil (Fetch tolerates any status code)", err)
				return
			}
			if snap.Metadata["status_code"] != fmt.Sprintf("%d", tt.statusCode) {
				t.Errorf("Fetch() status_code = %s, want %d", snap.Metadata["status_code"], tt.statusCode)
			}
		})
	}
}

func TestGuardedClientRejectsLoopback(t *testing.T) {
	c := &Connector{client: newGuardedClient()}
	err := c.Validate(context.Background(), map[string]any{"url": "http://127.0.0.1:8080/status"})
	if err == nil {
		t.Error("Validate(loopback) error = nil, want rejection")
	}
	if !strings.Contains(err.Error(), "blocked") && !strings.Contains(err.Error(), "Loopback") {
		t.Errorf("Validate(loopback) error = %v, want 'blocked' or 'Loopback' message", err)
	}
}

func TestGuardedClientRejectsLinkLocal(t *testing.T) {
	c := &Connector{client: newGuardedClient()}
	err := c.Validate(context.Background(), map[string]any{"url": "http://169.254.169.254/metadata"})
	if err == nil {
		t.Error("Validate(link-local) error = nil, want rejection")
	}
	if !strings.Contains(err.Error(), "blocked") && !strings.Contains(err.Error(), "link") {
		t.Errorf("Validate(link-local) error = %v, want 'blocked' or 'link' message", err)
	}
}

func TestFetchMalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	c := &Connector{client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
	if err != nil {
		t.Errorf("Fetch(malformed) error = %v, want nil (tolerates invalid JSON)", err)
		return
	}
	if !strings.Contains(snap.Sections[0].Content, "not json") {
		t.Errorf("Fetch(malformed) content = %q, want malformed JSON in output", snap.Sections[0].Content)
	}
}

func TestFetchStructuredEntitiesWithAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
  "entities": [
    {"kind": "vm", "name": "web-1", "ip": "192.168.1.10", "attributes": {"enabled": true, "cpu_count": 4}},
    {"kind": "vm", "name": "db-1", "ip": "192.168.1.20", "attributes": {"enabled": false, "cpu_count": 8}}
  ]
}`))
	}))
	defer server.Close()

	c := &Connector{client: server.Client()}
	snapshot, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	if len(snapshot.Entities) != 2 {
		t.Fatalf("snapshot.Entities = %d entities, want 2", len(snapshot.Entities))
	}

	// Check first entity
	if snapshot.Entities[0].Kind != "vm" || snapshot.Entities[0].Name != "web-1" || snapshot.Entities[0].IP != "192.168.1.10" {
		t.Fatalf("Entities[0] = %+v, want kind=vm name=web-1 ip=192.168.1.10", snapshot.Entities[0])
	}

	if enabled, ok := snapshot.Entities[0].Attributes["enabled"].(bool); !ok || !enabled {
		t.Errorf("Entities[0].Attributes[enabled] = %v, want true", snapshot.Entities[0].Attributes["enabled"])
	}

	if cpuCount, ok := snapshot.Entities[0].Attributes["cpu_count"].(float64); !ok || cpuCount != 4 {
		t.Errorf("Entities[0].Attributes[cpu_count] = %v, want 4", snapshot.Entities[0].Attributes["cpu_count"])
	}

	// Check second entity
	if snapshot.Entities[1].Kind != "vm" || snapshot.Entities[1].Name != "db-1" {
		t.Fatalf("Entities[1] = %+v, want kind=vm name=db-1", snapshot.Entities[1])
	}

	if enabled, ok := snapshot.Entities[1].Attributes["enabled"].(bool); !ok || enabled {
		t.Errorf("Entities[1].Attributes[enabled] = %v, want false", snapshot.Entities[1].Attributes["enabled"])
	}
}

func TestFetchStructuredEntitiesWithoutAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
  "entities": [
    {"kind": "rule", "name": "allow-ssh"},
    {"kind": "rule", "name": "deny-http"}
  ]
}`))
	}))
	defer server.Close()

	c := &Connector{client: server.Client()}
	snapshot, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	if len(snapshot.Entities) != 2 {
		t.Fatalf("snapshot.Entities = %d entities, want 2", len(snapshot.Entities))
	}

	if snapshot.Entities[0].Kind != "rule" || snapshot.Entities[0].Name != "allow-ssh" {
		t.Fatalf("Entities[0] = %+v", snapshot.Entities[0])
	}

	if snapshot.Entities[0].Attributes != nil {
		t.Errorf("Entities[0].Attributes = %v, want nil", snapshot.Entities[0].Attributes)
	}
}

func TestFetchPlainJSONWithoutEntities(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"status": "ok", "uptime": 12345}`))
	}))
	defer server.Close()

	c := &Connector{client: server.Client()}
	snapshot, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	// Should have no structured entities, just the Response section
	if len(snapshot.Entities) != 0 {
		t.Fatalf("snapshot.Entities = %d entities, want 0 (no entities array)", len(snapshot.Entities))
	}

	// Response section should contain the original JSON
	if !strings.Contains(snapshot.Sections[0].Content, `"status"`) {
		t.Errorf("Sections[0].Content missing original JSON")
	}
}

func TestFetchNonJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Service is running"))
	}))
	defer server.Close()

	c := &Connector{client: server.Client()}
	snapshot, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	// Should have no structured entities, just the Response section
	if len(snapshot.Entities) != 0 {
		t.Fatalf("snapshot.Entities = %d entities, want 0 (non-JSON response)", len(snapshot.Entities))
	}

	// Response section should contain the original content
	if !strings.Contains(snapshot.Sections[0].Content, "Service is running") {
		t.Errorf("Sections[0].Content = %q, want plain text", snapshot.Sections[0].Content)
	}
}

func TestFetchEntitiesWithMissingKind(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{
  "entities": [
    {"name": "no-kind", "ip": "10.0.0.1"},
    {"kind": "vm", "name": "valid", "ip": "10.0.0.2"}
  ]
}`))
	}))
	defer server.Close()

	c := &Connector{client: server.Client()}
	snapshot, err := c.Fetch(context.Background(), map[string]any{"url": server.URL})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	// Should only have the entity with a kind
	if len(snapshot.Entities) != 1 {
		t.Fatalf("snapshot.Entities = %d entities, want 1 (entity without kind should be skipped)", len(snapshot.Entities))
	}

	if snapshot.Entities[0].Kind != "vm" || snapshot.Entities[0].Name != "valid" {
		t.Fatalf("Entities[0] = %+v", snapshot.Entities[0])
	}
}
