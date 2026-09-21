package traefik

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const (
	overviewJSON = `{
		"http":{"routers":{"total":2,"warnings":0,"errors":1},"services":{"total":2,"warnings":0,"errors":0},"middlewares":{"total":1,"warnings":0,"errors":0}},
		"tcp":{"routers":{"total":1,"warnings":0,"errors":0},"services":{"total":1,"warnings":0,"errors":0},"middlewares":{"total":0,"warnings":0,"errors":0}},
		"udp":{"routers":{"total":0,"warnings":0,"errors":0},"services":{"total":0,"warnings":0,"errors":0},"middlewares":{"total":0,"warnings":0,"errors":0}},
		"features":{"tracing":"Jaeger","metrics":"Prometheus","accessLog":true},
		"providers":["Docker","File"]
	}`
	routersJSON = `[
		{"name":"app@docker","rule":"Host(` + "`app.example.com`" + `)","service":"app@docker","status":"enabled","provider":"docker","entryPoints":["websecure"],"middlewares":["auth@file"],"tls":{"certResolver":"letsencrypt"}},
		{"name":"api@internal","rule":"PathPrefix(` + "`/api`" + `)","service":"api@internal","status":"disabled","provider":"internal","entryPoints":["traefik"]}
	]`
	servicesJSON = `[
		{"name":"app@docker","status":"enabled","provider":"docker","serverStatus":{"http://10.0.0.5:8080":"UP","http://10.0.0.6:8080":"DOWN"},
		 "loadBalancer":{"passHostHeader":true,"servers":[{"url":"http://10.0.0.5:8080"},{"url":"http://10.0.0.6:8080"}]}},
		{"name":"api@internal","status":"enabled","provider":"internal"}
	]`
	middlewaresJSON = `[
		{"name":"auth@file","status":"enabled","provider":"file","usedBy":["app@docker"],"basicAuth":{"users":["u:hash"]}},
		{"name":"redirect@docker","status":"enabled","provider":"docker","type":"redirectScheme","redirectScheme":{"scheme":"https"}}
	]`
	entryPointsJSON = `[
		{"name":"web","address":":80","asDefault":false},
		{"name":"websecure","address":":443","asDefault":true,"http":{"tls":{"certResolver":"letsencrypt"}}}
	]`
	versionJSON = `{"Version":"3.0.0","Codename":"beaufort"}`
)

// traefikAPI returns an httptest server serving the full happy-path API. If
// authCheck is non-nil it gates every request; returning false makes the
// handler reply 401.
func traefikAPI(t *testing.T, authCheck func(*http.Request) bool) *httptest.Server {
	t.Helper()
	body := map[string]string{
		pathOverview:    overviewJSON,
		pathVersion:     versionJSON,
		pathRouters:     routersJSON,
		pathServices:    servicesJSON,
		pathMiddlewares: middlewaresJSON,
		pathEntryPoints: entryPointsJSON,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authCheck != nil && !authCheck(r) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		payload, ok := body[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	t.Cleanup(server.Close)
	return server
}

// newTestConnector builds the connector through the registered factory so
// the config plumbing (auth mode, verify_tls) is exercised too.
func newTestConnector(t *testing.T, config map[string]any) *Connector {
	t.Helper()
	connector.AllowLoopbackForTest(t)
	impl, err := connector.Get(typeName, config)
	if err != nil {
		t.Fatalf("connector.Get: %v", err)
	}
	c, ok := impl.(*Connector)
	if !ok {
		t.Fatalf("connector.Get returned %T, want *Connector", impl)
	}
	return c
}

func TestRegisteredSchema(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if schema.Category != "networking" || schema.Name != "Traefik" || schema.Stub {
		t.Fatalf("schema = %+v", schema)
	}
	want := map[string]bool{"url": true, "auth_mode": false, "username": false, "password": false, "api_token": false, "verify_tls": false}
	got := make(map[string]bool, len(schema.Fields))
	for _, f := range schema.Fields {
		got[f.Key] = f.Required
	}
	for key, required := range want {
		gotRequired, ok := got[key]
		if !ok {
			t.Errorf("schema missing field %q", key)
			continue
		}
		if gotRequired != required {
			t.Errorf("field %q Required = %t, want %t", key, gotRequired, required)
		}
	}
	if _, ok := connector.AttributeCatalog()[typeName]; !ok {
		t.Errorf("attribute catalog not registered for %q", typeName)
	}
}

// TestSchemaConfigValidation checks auth_mode is validated as a select by
// the shared config validator.
func TestSchemaConfigValidation(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if err := connector.ValidateConfig(*schema, map[string]any{"auth_mode": authBasic}); err != nil {
		t.Errorf("ValidateConfig(basic) = %v, want nil", err)
	}
	if err := connector.ValidateConfig(*schema, map[string]any{"auth_mode": "kerberos"}); err == nil {
		t.Error("ValidateConfig(kerberos) = nil, want an error")
	}
}

func TestIdentity(t *testing.T) {
	c := &Connector{}
	if c.Name() != "Traefik" || c.Type() != "traefik" || c.Category() != "networking" {
		t.Fatalf("identity = %s/%s/%s", c.Name(), c.Type(), c.Category())
	}
}

func TestValidateAuthModes(t *testing.T) {
	server := traefikAPI(t, func(r *http.Request) bool {
		switch {
		case r.Header.Get("Authorization") == "":
			return true // unauthenticated deployments are allowed
		case r.Header.Get("Authorization") == "Bearer tok123":
			return true
		default:
			user, pass, ok := r.BasicAuth()
			return ok && user == "admin" && pass == "hunter2"
		}
	})

	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{name: "no auth", config: map[string]any{"url": server.URL}},
		{name: "explicit none", config: map[string]any{"url": server.URL, "auth_mode": authNone}},
		{name: "basic auth", config: map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "admin", "password": "hunter2"}},
		{name: "bearer token", config: map[string]any{"url": server.URL, "auth_mode": authToken, "api_token": "tok123"}},
		{
			name:    "basic auth missing password",
			config:  map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "admin"},
			wantErr: "basic auth requires",
		},
		{
			name:    "token mode missing token",
			config:  map[string]any{"url": server.URL, "auth_mode": authToken},
			wantErr: "requires an api_token",
		},
		{
			name:    "unknown auth mode",
			config:  map[string]any{"url": server.URL, "auth_mode": "kerberos"},
			wantErr: "unknown auth mode",
		},
		{name: "missing url", config: map[string]any{}, wantErr: "url is required"},
		{
			name:    "wrong basic credentials",
			config:  map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "admin", "password": "wrong"},
			wantErr: "auth error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestConnector(t, tt.config)
			err := c.Validate(context.Background(), nil)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

// TestBasicAuthHeaderIsSent asserts the exact wire format of both
// authenticated modes.
func TestBasicAuthHeaderIsSent(t *testing.T) {
	var gotAuth string
	server := traefikAPI(t, func(r *http.Request) bool {
		gotAuth = r.Header.Get("Authorization")
		return true
	})

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "ad:min", "password": "p@ss"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("ad:min:p@ss"))
	if gotAuth != want {
		t.Errorf("basic Authorization = %q, want %q", gotAuth, want)
	}

	c = newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authToken, "api_token": "tok123"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if gotAuth != "Bearer tok123" {
		t.Errorf("token Authorization = %q, want %q", gotAuth, "Bearer tok123")
	}
}

// TestNoAuthModeSendsNoAuthorizationHeader guards against leaking stray
// credentials when the API is unauthenticated.
func TestNoAuthModeSendsNoAuthorizationHeader(t *testing.T) {
	var seen []string
	server := traefikAPI(t, func(r *http.Request) bool {
		seen = append(seen, r.Header.Get("Authorization"))
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "username": "admin", "password": "hunter2", "api_token": "tok"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	for _, h := range seen {
		if h != "" {
			t.Fatalf("Authorization header = %q, want empty in auth mode none", h)
		}
	}
}

func TestFetchHappyPath(t *testing.T) {
	server := traefikAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL + "/"})

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.ServiceName != "Traefik" || snapshot.Type != typeName {
		t.Fatalf("snapshot = %s/%s", snapshot.ServiceName, snapshot.Type)
	}

	wantSections := []string{"Overview", "HTTP Routers", "HTTP Services", "HTTP Middlewares", "Entry Points"}
	if len(snapshot.Sections) != len(wantSections) {
		t.Fatalf("sections = %d, want %d", len(snapshot.Sections), len(wantSections))
	}
	for i, title := range wantSections {
		if snapshot.Sections[i].Title != title {
			t.Errorf("sections[%d].Title = %q, want %q", i, snapshot.Sections[i].Title, title)
		}
		if strings.Contains(snapshot.Sections[i].Content, "unavailable") {
			t.Errorf("sections[%d] unexpectedly unavailable: %q", i, snapshot.Sections[i].Content)
		}
	}

	wantMetadata := map[string]string{
		"traefik_url":            server.URL,
		"traefik_version":        "3.0.0",
		"http_routers_total":     "2",
		"http_services_total":    "2",
		"http_middlewares_total": "1",
		"tcp_routers_total":      "1",
		"udp_routers_total":      "0",
		"router_count":           "2",
		"service_count":          "2",
		"middleware_count":       "2",
	}
	for k, want := range wantMetadata {
		if got := snapshot.Metadata[k]; got != want {
			t.Errorf("metadata[%q] = %q, want %q", k, got, want)
		}
	}

	byKind := map[string]int{}
	for _, e := range snapshot.Entities {
		byKind[e.Kind]++
	}
	want := map[string]int{"router": 2, "service": 2, "middleware": 2, "entrypoint": 2}
	for kind, n := range want {
		if byKind[kind] != n {
			t.Errorf("entities of kind %q = %d, want %d", kind, byKind[kind], n)
		}
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "upstream_service", Name: "api@internal"},
		{Kind: "upstream_service", Name: "app@docker"},
	}
	if len(snapshot.Dependencies) != len(wantDeps) {
		t.Fatalf("dependencies = %+v, want %+v", snapshot.Dependencies, wantDeps)
	}
	for i, d := range wantDeps {
		if snapshot.Dependencies[i] != d {
			t.Errorf("dependencies[%d] = %+v, want %+v", i, snapshot.Dependencies[i], d)
		}
	}
}

// TestFetchSelectiveFields checks the "fields" hint limits which endpoints
// are hit and which sections come back.
func TestFetchSelectiveFields(t *testing.T) {
	var hit []string
	server := traefikAPI(t, func(r *http.Request) bool {
		hit = append(hit, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"routers"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "HTTP Routers" {
		t.Fatalf("sections = %+v, want only HTTP Routers", snapshot.Sections)
	}
	for _, p := range hit {
		if p == pathServices || p == pathMiddlewares || p == pathEntryPoints || p == pathOverview {
			t.Errorf("unexpected request to %q for a routers-only fetch", p)
		}
	}
}

// TestFetchDegradesPerSection checks a failing endpoint only degrades its
// own section instead of failing the whole snapshot.
func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathOverview:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(overviewJSON))
		case pathRouters:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		case pathServices:
			_, _ = w.Write([]byte("not json"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	connector.AllowLoopbackForTest(t)
	c := newTestConnector(t, map[string]any{"url": server.URL})
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 5 {
		t.Fatalf("sections = %d, want 5", len(snapshot.Sections))
	}
	byTitle := map[string]string{}
	for _, s := range snapshot.Sections {
		byTitle[s.Title] = s.Content
	}
	if strings.Contains(byTitle["Overview"], "unavailable") {
		t.Errorf("Overview should have succeeded: %q", byTitle["Overview"])
	}
	if !strings.Contains(byTitle["HTTP Routers"], "API returned 500") {
		t.Errorf("HTTP Routers = %q, want the upstream 500 reported", byTitle["HTTP Routers"])
	}
	if !strings.Contains(byTitle["HTTP Services"], "malformed response") {
		t.Errorf("HTTP Services = %q, want a malformed-response placeholder", byTitle["HTTP Services"])
	}
	if !strings.Contains(byTitle["Entry Points"], "API returned 404") {
		t.Errorf("Entry Points = %q, want the upstream 404 reported", byTitle["Entry Points"])
	}
	if _, ok := snapshot.Metadata["traefik_version"]; ok {
		t.Errorf("traefik_version should be absent when /api/version 404s: %+v", snapshot.Metadata)
	}
}

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		name   string
		status int
		assert func(*testing.T, error)
	}{
		{
			name: "401 maps to AuthError", status: http.StatusUnauthorized,
			assert: func(t *testing.T, err error) {
				var target *connector.AuthError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.AuthError", err, err)
				}
			},
		},
		{
			name: "403 maps to AuthError", status: http.StatusForbidden,
			assert: func(t *testing.T, err error) {
				var target *connector.AuthError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.AuthError", err, err)
				}
			},
		},
		{
			name: "502 maps to ServiceUnavailableError", status: http.StatusBadGateway,
			assert: func(t *testing.T, err error) {
				var target *connector.ServiceUnavailableError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.ServiceUnavailableError", err, err)
				}
			},
		},
		{
			name: "503 maps to ServiceUnavailableError", status: http.StatusServiceUnavailable,
			assert: func(t *testing.T, err error) {
				var target *connector.ServiceUnavailableError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.ServiceUnavailableError", err, err)
				}
			},
		},
		{
			name: "504 maps to ServiceUnavailableError", status: http.StatusGatewayTimeout,
			assert: func(t *testing.T, err error) {
				var target *connector.ServiceUnavailableError
				if !errors.As(err, &target) {
					t.Fatalf("error = %v (%T), want *connector.ServiceUnavailableError", err, err)
				}
			},
		},
		{
			name: "500 stays a plain error", status: http.StatusInternalServerError,
			assert: func(t *testing.T, err error) {
				var authErr *connector.AuthError
				var unavailable *connector.ServiceUnavailableError
				if errors.As(err, &authErr) || errors.As(err, &unavailable) {
					t.Fatalf("error = %v (%T), want a plain error", err, err)
				}
				if !strings.Contains(err.Error(), "API returned 500") {
					t.Fatalf("error = %v, want it to mention the status", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte("upstream said no"))
			}))
			defer server.Close()

			c := newTestConnector(t, map[string]any{"url": server.URL})
			err := c.Validate(context.Background(), nil)
			if err == nil {
				t.Fatal("Validate() error = nil, want an error")
			}
			tt.assert(t, err)
		})
	}
}

func TestExpiredDeadlineMapsToTimeout(t *testing.T) {
	server := traefikAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL})

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	err := c.Validate(ctx, nil)
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("Validate() error = %v (%T), want *connector.TimeoutError", err, err)
	}
}

// TestOversizedResponseIsRejected checks the shared body cap is enforced.
func TestOversizedResponseIsRejected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		chunk := strings.Repeat("a", 1<<20)
		for written := 0; written <= connector.MaxResponseBytes; written += len(chunk) {
			if _, err := w.Write([]byte(chunk)); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("Validate() error = %v, want the response-size cap to trip", err)
	}
}

// TestGuardedDialerBlocksLoopback confirms the SSRF guard is wired in when
// the test override is not active.
func TestGuardedDialerBlocksLoopback(t *testing.T) {
	server := traefikAPI(t, nil)
	impl, err := connector.Get(typeName, map[string]any{"url": server.URL})
	if err != nil {
		t.Fatalf("connector.Get: %v", err)
	}
	if err := impl.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() error = nil, want the guarded dialer to block loopback")
	}
}

func TestRedirectsAreNotFollowed(t *testing.T) {
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Path == pathOverview {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() error = %v, want the 302 to be returned as-is", err)
	}
	if hits != 1 {
		t.Fatalf("hits = %d, want 1 (redirect must not be followed)", hits)
	}
}

func TestVerifyTLSDisabledAllowsSelfSignedCert(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(overviewJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "verify_tls": false})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() with verify_tls=false error = %v, want nil", err)
	}

	strict := newTestConnector(t, map[string]any{"url": server.URL, "verify_tls": true})
	if err := strict.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() with verify_tls=true error = nil, want a certificate error")
	}
}
