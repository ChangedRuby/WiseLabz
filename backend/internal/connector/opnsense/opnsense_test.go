package opnsense

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestValidateUsesBasicAuthAndSurfacesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || user != "key" || password != "secret" || r.URL.Path != "/api/core/firmware/status" {
			t.Fatalf("request auth/path invalid")
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("denied"))
	}))
	defer server.Close()
	c := &Connector{url: server.URL, apiKey: "key", apiSecret: "secret", client: server.Client()}
	err := c.Validate(context.Background(), nil)
	var authErr *connector.AuthError
	if !errors.As(err, &authErr) || !strings.Contains(err.Error(), "API returned 401: denied") {
		t.Fatalf("Validate() error = %v, want *connector.AuthError", err)
	}
}

func TestFetchSurfacesMalformedSystemResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/core/firmware/status":
			_, _ = w.Write([]byte(`not json`))
		case "/api/diagnostics/interface/getInterfaces", "/api/firewall/filter/searchRule", "/api/routes/gateway/status":
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, apiKey: "key", apiSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if !strings.Contains(snap.Sections[0].Content, "malformed response") {
		t.Fatalf("System section = %q, want malformed response placeholder", snap.Sections[0].Content)
	}
}

func TestFetchSurfacesWANAndUpstreamDependencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/core/firmware/status":
			_, _ = w.Write([]byte(`{"product_name":"OPNsense","product_version":"24.1"}`))
		case "/api/diagnostics/interface/getInterfaces":
			_, _ = w.Write([]byte(`{"rows":[
				{"identifier":"wan","device":"igb0","ipaddr":"203.0.113.5","status":"up","media":"1000baseT"},
				{"identifier":"lan","device":"igb1","ipaddr":"10.0.0.1","status":"up","media":"1000baseT"}
			]}`))
		case "/api/firewall/filter/searchRule":
			_, _ = w.Write([]byte(`{"rows":[]}`))
		case "/api/routes/gateway/status":
			_, _ = w.Write([]byte(`{"items":[{"name":"WAN_GW","address":"203.0.113.1","status":"online","rtt":"5ms","loss":"0%"}]}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, apiKey: "key", apiSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if len(snap.Sections) != 4 {
		t.Fatalf("Sections = %d, want 4", len(snap.Sections))
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "network", Name: "igb0"},
		{Kind: "upstream_service", Name: "WAN_GW"},
	}
	for _, want := range wantDeps {
		found := false
		for _, got := range snap.Dependencies {
			if got.Kind == want.Kind && got.Name == want.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Dependencies missing %+v, got %+v", want, snap.Dependencies)
		}
	}
}

func TestDoRequestErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		checkAuthError bool
		checkUnavail   bool
	}{
		{
			name:           "401 Unauthorized returns AuthError",
			statusCode:     http.StatusUnauthorized,
			checkAuthError: true,
		},
		{
			name:           "403 Forbidden returns AuthError",
			statusCode:     http.StatusForbidden,
			checkAuthError: true,
		},
		{
			name:         "502 BadGateway returns ServiceUnavailableError",
			statusCode:   http.StatusBadGateway,
			checkUnavail: true,
		},
		{
			name:         "503 ServiceUnavailable returns ServiceUnavailableError",
			statusCode:   http.StatusServiceUnavailable,
			checkUnavail: true,
		},
		{
			name:         "504 GatewayTimeout returns ServiceUnavailableError",
			statusCode:   http.StatusGatewayTimeout,
			checkUnavail: true,
		},
		{
			name:       "500 InternalServerError returns generic error",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("error response"))
			}))
			defer server.Close()

			c := &Connector{url: server.URL, apiKey: "key", apiSecret: "secret", client: server.Client()}
			_, err := c.doRequest(context.Background(), "GET", "/api/test")

			if err == nil {
				t.Errorf("doRequest() error = nil, want error")
				return
			}

			if tt.checkAuthError {
				var authErr *connector.AuthError
				if !errors.As(err, &authErr) {
					t.Errorf("doRequest() error = %T, want *connector.AuthError", err)
				}
			}
			if tt.checkUnavail {
				var unavailErr *connector.ServiceUnavailableError
				if !errors.As(err, &unavailErr) {
					t.Errorf("doRequest() error = %T, want *connector.ServiceUnavailableError", err)
				}
			}
		})
	}
}

func TestDoRequestContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(1 * time.Second)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := &Connector{url: server.URL, apiKey: "key", apiSecret: "secret", client: server.Client()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := c.doRequest(ctx, "GET", "/api/test")
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Errorf("doRequest() error = %v, want *connector.TimeoutError", err)
	}
}

func TestBuildRuleTableMalformedCases(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{
			name: "empty JSON object returns placeholder",
			data: []byte(`{}`),
			want: "_No firewall rules returned_",
		},
		{
			name: "JSON with empty rows returns placeholder",
			data: []byte(`{"rows":[]}`),
			want: "_No firewall rules returned_",
		},
		{
			name: "invalid JSON returns placeholder",
			data: []byte(`not json`),
			want: "_No firewall rules returned_",
		},
		{
			name: "partial JSON returns placeholder",
			data: []byte(`{"rows":[{"description":"test"`),
			want: "_No firewall rules returned_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := buildRuleTable(tt.data)
			if result != tt.want {
				t.Errorf("buildRuleTable() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestBuildRuleTableValidRules(t *testing.T) {
	data := []byte(`{
		"rows":[
			{"description":"Allow SSH","action":"pass","protocol":"tcp","source_net":"any","destination_net":"any","enabled":"1"},
			{"description":"Block DNS","action":"block","protocol":"udp","source_net":"10.0.0.0/8","destination_net":"any","enabled":""}
		]
	}`)
	result, _ := buildRuleTable(data)
	if !strings.Contains(result, "Allow SSH") || !strings.Contains(result, "Block DNS") {
		t.Errorf("buildRuleTable() missing expected rules in: %q", result)
	}
}
