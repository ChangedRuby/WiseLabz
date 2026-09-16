package pfsense

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestValidateUsesBearerTokenAndSurfacesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer key" || r.URL.Path != "/api/v2/system/version" {
			t.Fatalf("request auth/path invalid")
		}
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("denied"))
	}))
	defer server.Close()
	c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
	err := c.Validate(context.Background(), nil)
	var authErr *connector.AuthError
	if !errors.As(err, &authErr) || !strings.Contains(err.Error(), "API returned 401: denied") {
		t.Fatalf("Validate() error = %v, want *connector.AuthError", err)
	}
}

func TestFetchSurfacesWANAndUpstreamDependencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/system/version":
			_, _ = w.Write([]byte(`{"data":{"platform":"pfSense","config_version":"2.7.2"}}`))
		case "/api/v2/interfaces":
			_, _ = w.Write([]byte(`{"data":[
				{"id":"wan","if":"em0","ipaddr":"203.0.113.5","status":"up","enable":true},
				{"id":"lan","if":"em1","ipaddr":"10.0.0.1","status":"up","enable":true}
			]}`))
		case "/api/v2/firewall/rules":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/api/v2/routing/gateways":
			_, _ = w.Write([]byte(`{"data":[{"name":"WAN_GW","gateway":"203.0.113.1","status":"online","monitor_rtt":"5ms","monitor_loss":"0%"}]}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if len(snap.Sections) != 4 {
		t.Fatalf("Sections = %d, want 4", len(snap.Sections))
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "network", Name: "em0"},
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

func TestFetchSurfacesMalformedSystemResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/system/version":
			_, _ = w.Write([]byte(`not json`))
		case "/api/v2/interfaces", "/api/v2/firewall/rules", "/api/v2/routing/gateways":
			_, _ = w.Write([]byte(`{}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if !strings.Contains(snap.Sections[0].Content, "malformed response") {
		t.Fatalf("System section = %q, want malformed response placeholder", snap.Sections[0].Content)
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

			c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
			_, err := c.doRequest(context.Background(), "/api/test")

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

	c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := c.doRequest(ctx, "/api/test")
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
			name: "JSON with empty data returns placeholder",
			data: []byte(`{"data":[]}`),
			want: "_No firewall rules returned_",
		},
		{
			name: "invalid JSON returns placeholder",
			data: []byte(`not json`),
			want: "_No firewall rules returned_",
		},
		{
			name: "partial JSON returns placeholder",
			data: []byte(`{"data":[{"descr":"test"`),
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
		"data":[
			{"descr":"Allow SSH","type":"pass","protocol":"tcp","source":"any","destination":"any","disabled":false},
			{"descr":"Block DNS","type":"block","protocol":"udp","source":"10.0.0.0/8","destination":"any","disabled":true}
		]
	}`)
	result, _ := buildRuleTable(data)
	if !strings.Contains(result, "Allow SSH") || !strings.Contains(result, "Block DNS") {
		t.Errorf("buildRuleTable() missing expected rules in: %q", result)
	}
}

func TestRestart(t *testing.T) {
	tests := []struct {
		name       string
		entityRef  string
		statusCode int
		wantErr    bool
	}{
		{name: "success", entityRef: "unbound", statusCode: http.StatusOK},
		{name: "empty entityRef errors", entityRef: "", wantErr: true},
		{name: "http error", entityRef: "unbound", statusCode: http.StatusInternalServerError, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotMethod, gotBody string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotMethod = r.URL.Path, r.Method
				b, _ := io.ReadAll(r.Body)
				gotBody = string(b)
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer server.Close()

			c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
			err := c.Restart(context.Background(), nil, tt.entityRef)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Restart() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Restart() error = %v", err)
			}
			if gotMethod != "POST" || gotPath != "/api/v2/status/service" {
				t.Errorf("request = %s %s, want POST /api/v2/status/service", gotMethod, gotPath)
			}
			if !strings.Contains(gotBody, tt.entityRef) || !strings.Contains(gotBody, "restart") {
				t.Errorf("body = %q, want it to reference %q and restart action", gotBody, tt.entityRef)
			}
		})
	}
}
