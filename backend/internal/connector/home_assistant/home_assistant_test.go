package homeassistant

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

const testToken = "llat-abc123"

// homeAssistantAPI returns an httptest server serving the full happy-path
// REST API. Every request must carry the bearer token; anything else gets a
// 401, exactly like Home Assistant.
func homeAssistantAPI(t *testing.T, onRequest func(*http.Request)) *httptest.Server {
	t.Helper()
	body := map[string]string{
		pathConfig:   configJSON,
		pathStates:   statesJSON,
		pathServices: servicesJSON,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if onRequest != nil {
			onRequest(r)
		}
		if r.Header.Get("Authorization") != "Bearer "+testToken {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"Unauthorized"}`))
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
// the config plumbing (token, max_entities, verify_tls) is exercised too.
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
	if schema.Category != category || schema.Name != "Home Assistant" || schema.Stub {
		t.Fatalf("schema = %+v", schema)
	}
	want := map[string]bool{"url": true, "access_token": true, "max_entities": false, "verify_tls": false}
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
	for _, f := range schema.Fields {
		if f.Key == "access_token" && f.Type != "password" {
			t.Errorf("access_token Type = %q, want password so it is encrypted at rest", f.Type)
		}
	}
	if _, ok := connector.AttributeCatalog()[typeName]; !ok {
		t.Errorf("attribute catalog not registered for %q", typeName)
	}
}

func TestIdentity(t *testing.T) {
	c := &Connector{}
	if c.Name() != "Home Assistant" || c.Type() != "home_assistant" || c.Category() != "virtualization" {
		t.Fatalf("identity = %s/%s/%s", c.Name(), c.Type(), c.Category())
	}
}

func TestValidate(t *testing.T) {
	server := homeAssistantAPI(t, nil)

	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{name: "valid", config: map[string]any{"url": server.URL, "access_token": testToken}},
		{name: "trailing slash is trimmed", config: map[string]any{"url": server.URL + "/", "access_token": testToken}},
		{name: "missing url", config: map[string]any{"access_token": testToken}, wantErr: "url is required"},
		{name: "missing token", config: map[string]any{"url": server.URL}, wantErr: "access_token is required"},
		{name: "wrong token", config: map[string]any{"url": server.URL, "access_token": "nope"}, wantErr: "auth error"},
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

// TestBearerHeaderIsSent asserts the exact wire format Home Assistant
// expects for a long-lived access token.
func TestBearerHeaderIsSent(t *testing.T) {
	var gotAuth, gotAccept string
	server := homeAssistantAPI(t, func(r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAccept = r.Header.Get("Accept")
	})

	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if gotAuth != "Bearer "+testToken {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer "+testToken)
	}
	if gotAccept != "application/json" {
		t.Errorf("Accept = %q, want application/json", gotAccept)
	}
}

func TestFetchHappyPath(t *testing.T) {
	server := homeAssistantAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.ServiceName != "Home Assistant" || snapshot.Type != typeName {
		t.Fatalf("snapshot = %s/%s", snapshot.ServiceName, snapshot.Type)
	}

	wantSections := []string{"Overview", "Integrations", "Entity Domains", "Entities", "Services"}
	if len(snapshot.Sections) != len(wantSections) {
		t.Fatalf("sections = %d, want %d", len(snapshot.Sections), len(wantSections))
	}
	for i, title := range wantSections {
		if snapshot.Sections[i].Title != title {
			t.Errorf("sections[%d].Title = %q, want %q", i, snapshot.Sections[i].Title, title)
		}
		if strings.Contains(snapshot.Sections[i].Content, "unavailable:") {
			t.Errorf("sections[%d] unexpectedly unavailable: %q", i, snapshot.Sections[i].Content)
		}
	}

	wantMetadata := map[string]string{
		"home_assistant_url":       server.URL,
		"home_assistant_version":   "2024.6.4",
		"home_assistant_state":     "RUNNING",
		"location_name":            "Home | Lab",
		"time_zone":                "Europe/Lisbon",
		"config_source":            "storage",
		"component_count":          "6",
		"integration_count":        "5",
		"entity_count":             "11",
		"entity_domain_count":      "9",
		"unavailable_entity_count": "2",
		"service_domain_count":     "2",
		"service_count":            "20",
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
	if want := map[string]int{"integration": 5, "entity": 11}; byKind["integration"] != want["integration"] || byKind["entity"] != want["entity"] {
		t.Errorf("entities by kind = %v, want %v", byKind, want)
	}
	if snapshot.FetchedAt.IsZero() {
		t.Error("FetchedAt is zero")
	}
}

// TestFetchIsStableAcrossSyncs is the regression guard for the whole point
// of this connector: an idle instance that only restamped its entities must
// produce an identical snapshot.
func TestFetchIsStableAcrossSyncs(t *testing.T) {
	states := statesJSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathConfig:
			_, _ = w.Write([]byte(configJSON))
		case pathStates:
			_, _ = w.Write([]byte(states))
		case pathServices:
			_, _ = w.Write([]byte(servicesJSON))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	connector.AllowLoopbackForTest(t)
	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
	first, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	states = strings.ReplaceAll(states, "2024-06-01T10:00:00+00:00", "2025-01-09T22:13:05+00:00")
	states = strings.ReplaceAll(states, `"state":"21.4"`, `"state":"23.9"`)
	second, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	for i := range first.Sections {
		if first.Sections[i].Content != second.Sections[i].Content {
			t.Errorf("section %q changed between syncs:\n%s\n---\n%s",
				first.Sections[i].Title, first.Sections[i].Content, second.Sections[i].Content)
		}
	}
}

// TestFetchSelectiveFields checks the "fields" hint limits which endpoints
// are hit and which sections come back.
func TestFetchSelectiveFields(t *testing.T) {
	tests := []struct {
		field        string
		wantSections []string
		wantPaths    []string
	}{
		{field: "config", wantSections: []string{"Overview"}, wantPaths: []string{pathConfig}},
		{field: "integrations", wantSections: []string{"Integrations"}, wantPaths: []string{pathConfig}},
		{field: "entities", wantSections: []string{"Entity Domains", "Entities"}, wantPaths: []string{pathStates}},
		{field: "services", wantSections: []string{"Services"}, wantPaths: []string{pathServices}},
	}
	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			var hit []string
			server := homeAssistantAPI(t, func(r *http.Request) { hit = append(hit, r.URL.Path) })
			c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})

			snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []any{tt.field}})
			if err != nil {
				t.Fatalf("Fetch: %v", err)
			}
			var titles []string
			for _, s := range snapshot.Sections {
				titles = append(titles, s.Title)
			}
			if strings.Join(titles, ",") != strings.Join(tt.wantSections, ",") {
				t.Fatalf("sections = %v, want %v", titles, tt.wantSections)
			}
			if strings.Join(hit, ",") != strings.Join(tt.wantPaths, ",") {
				t.Fatalf("requested %v, want %v", hit, tt.wantPaths)
			}
		})
	}
}

// TestFetchConfigIsRequestedOnce checks the two sections backed by
// /api/config share a single request.
func TestFetchConfigIsRequestedOnce(t *testing.T) {
	var configHits int
	server := homeAssistantAPI(t, func(r *http.Request) {
		if r.URL.Path == pathConfig {
			configHits++
		}
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
	if _, err := c.Fetch(context.Background(), nil); err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if configHits != 1 {
		t.Errorf("/api/config hits = %d, want 1", configHits)
	}
}

// TestFetchDegradesPerSection checks a failing endpoint only degrades its
// own sections instead of failing the whole snapshot.
func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathConfig:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		case pathStates:
			_, _ = w.Write([]byte("not json"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	connector.AllowLoopbackForTest(t)
	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
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
	for _, title := range []string{"Overview", "Integrations"} {
		if !strings.Contains(byTitle[title], "API returned 500") {
			t.Errorf("%s = %q, want the upstream 500 reported", title, byTitle[title])
		}
	}
	for _, title := range []string{"Entity Domains", "Entities"} {
		if !strings.Contains(byTitle[title], "malformed response") {
			t.Errorf("%s = %q, want a malformed-response placeholder", title, byTitle[title])
		}
	}
	if !strings.Contains(byTitle["Services"], "API returned 404") {
		t.Errorf("Services = %q, want the upstream 404 reported", byTitle["Services"])
	}
	if _, ok := snapshot.Metadata["home_assistant_version"]; ok {
		t.Errorf("version metadata should be absent when /api/config fails: %+v", snapshot.Metadata)
	}
}

func TestMaxEntitiesConfig(t *testing.T) {
	tests := []struct {
		name string
		raw  any
		want int
	}{
		{name: "absent", raw: nil, want: defaultMaxEntities},
		{name: "json number", raw: float64(25), want: 25},
		{name: "int", raw: 7, want: 7},
		{name: "schema default string", raw: "500", want: 500},
		{name: "empty string", raw: "", want: defaultMaxEntities},
		{name: "unlimited", raw: float64(0), want: 0},
		{name: "negative", raw: float64(-3), want: 0},
		{name: "garbage", raw: "many", want: defaultMaxEntities},
		{name: "wrong type", raw: true, want: defaultMaxEntities},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := map[string]any{"url": "http://ha", "access_token": testToken}
			if tt.raw != nil {
				config["max_entities"] = tt.raw
			}
			c := newTestConnector(t, config)
			if c.maxEntities != tt.want {
				t.Fatalf("maxEntities = %d, want %d", c.maxEntities, tt.want)
			}
		})
	}
}

func TestFetchAppliesMaxEntities(t *testing.T) {
	server := homeAssistantAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken, "max_entities": float64(2)})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"entities"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(snapshot.Entities))
	}
	if snapshot.Metadata["entity_count"] != "11" {
		t.Errorf("entity_count = %q, want the untruncated total", snapshot.Metadata["entity_count"])
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

			c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
			err := c.Validate(context.Background(), nil)
			if err == nil {
				t.Fatal("Validate() error = nil, want an error")
			}
			tt.assert(t, err)
		})
	}
}

func TestExpiredDeadlineMapsToTimeout(t *testing.T) {
	server := homeAssistantAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	err := c.Validate(ctx, nil)
	var timeoutErr *connector.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("Validate() error = %v (%T), want *connector.TimeoutError", err, err)
	}
}

// TestOversizedResponseIsRejected checks the shared body cap is enforced —
// /api/states on a large instance is the realistic way to hit it.
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

	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("Validate() error = %v, want the response-size cap to trip", err)
	}
}

// TestGuardedDialerBlocksLoopback confirms the SSRF guard is wired in when
// the test override is not active.
func TestGuardedDialerBlocksLoopback(t *testing.T) {
	server := homeAssistantAPI(t, nil)
	impl, err := connector.Get(typeName, map[string]any{"url": server.URL, "access_token": testToken})
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
		if r.URL.Path == pathConfig {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken})
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
		_, _ = w.Write([]byte(configJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken, "verify_tls": false})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() with verify_tls=false error = %v, want nil", err)
	}

	strict := newTestConnector(t, map[string]any{"url": server.URL, "access_token": testToken, "verify_tls": true})
	if err := strict.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() with verify_tls=true error = nil, want a certificate error")
	}
}
