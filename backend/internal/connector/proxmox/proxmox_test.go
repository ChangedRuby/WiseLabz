package proxmox

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestFetchIncludesHostAndStorageDependencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nodes":
			_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
		case "/nodes/pve1/qemu":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/nodes/pve1/lxc":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/nodes/pve1/storage":
			_, _ = w.Write([]byte(`{"data":[{"storage":"local-zfs","type":"zfspool","used":10,"total":100,"avail":90}]}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if len(snap.Sections) != 1 || snap.Sections[0].Title != "pve1" {
		t.Fatalf("Sections = %+v, want one section for pve1", snap.Sections)
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "host", Name: "pve1"},
		{Kind: "storage", Name: "local-zfs"},
	}
	if !reflect.DeepEqual(snap.Dependencies, wantDeps) {
		t.Fatalf("Dependencies = %+v, want %+v", snap.Dependencies, wantDeps)
	}
}

func TestFetchWithFieldsHintSkipsUnrequestedCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nodes":
			_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
		case "/nodes/pve1/qemu":
			_, _ = w.Write([]byte(`{"data":[{"vmid":100,"name":"vm1","status":"running","cpus":2,"mem":1024,"uptime":10}]}`))
		default:
			t.Fatalf("unexpected request path %s: selective fetch should only hit /nodes and /nodes/pve1/qemu", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"vms"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if !strings.Contains(snap.Sections[0].Content, "Virtual Machines") {
		t.Errorf("requested field 'vms' missing from output: %q", snap.Sections[0].Content)
	}
	if strings.Contains(snap.Sections[0].Content, "Storage") || strings.Contains(snap.Sections[0].Content, "Containers") {
		t.Errorf("unrequested sections present in output: %q", snap.Sections[0].Content)
	}
}

func TestValidateUsesTokenAndSurfacesStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "PVEAPIToken=user@pam!token=secret" || r.URL.Path != "/nodes" {
			t.Fatalf("request auth/path invalid")
		}
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("denied"))
	}))
	defer server.Close()
	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	err := c.Validate(context.Background(), nil)
	var authErr *connector.AuthError
	if !errors.As(err, &authErr) || !strings.Contains(err.Error(), "API returned 403: denied") {
		t.Fatalf("Validate() error = %v, want *connector.AuthError", err)
	}
}

func TestFetchSurfacesMalformedNodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nodes" {
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	_, err := c.Fetch(context.Background(), nil)
	var malformedErr *connector.MalformedResponseError
	if !errors.As(err, &malformedErr) {
		t.Fatalf("Fetch() error = %v, want *connector.MalformedResponseError", err)
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

			c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
			_, err := c.doRequest(context.Background(), "GET", "/nodes", nil)

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

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := c.doRequest(ctx, "GET", "/nodes", nil)
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Errorf("doRequest() error = %v, want *connector.TimeoutError", err)
	}
}
