package portainer

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

const testAPIKey = "ptr_testtoken"

const (
	systemStatusJSON = `{"Version":"2.19.4","InstanceID":"b1c2d3e4-0000-4000-8000-abcdefabcdef"}`
	endpointsJSON    = `[
		{"Id":1,"Name":"local","Type":1,"URL":"unix:///var/run/docker.sock","PublicURL":"portainer.example.com","Status":1,"GroupId":1,"TLS":false,"TagIds":[3,4],
		 "Snapshots":[{"DockerVersion":"24.0.7","Swarm":false}]},
		{"Id":2,"Name":"edge-nas","Type":2,"URL":"tcp://10.0.0.5:9001/","Status":1,"GroupId":2,"TLS":true,"TagIds":[],
		 "Snapshots":[{"DockerVersion":"25.0.3","Swarm":true}]},
		{"Id":3,"Name":"offline-agent","Type":2,"URL":"tcp://10.0.0.9:9001","Status":2,"GroupId":1},
		{"Id":4,"Name":"k8s-prod","Type":5,"URL":"https://10.0.0.20:6443","Status":1,"GroupId":1}
	]`
	stacksJSON = `[
		{"Id":7,"Name":"blog","Type":2,"EndpointId":1,"Status":1,"EntryPoint":"docker-compose.yml"},
		{"Id":8,"Name":"metrics","Type":1,"EndpointId":2,"Status":2,"EntryPoint":"stack.yml"},
		{"Id":9,"Name":"orphan","Type":2,"EndpointId":99,"Status":1,"EntryPoint":""}
	]`
	containersJSON = `[
		{"Id":"c0ffee11","Names":["/blog-web"],"Image":"ghcr.io/example/blog:1.2","State":"running","Status":"Up 3 days",
		 "Labels":{"com.docker.compose.project":"blog"},
		 "Ports":[{"PrivatePort":80,"PublicPort":8080,"Type":"tcp"},{"PrivatePort":9000,"PublicPort":0,"Type":"tcp"}],
		 "HostConfig":{"NetworkMode":"blog_default"},
		 "NetworkSettings":{"Networks":{"blog_default":{"IPAddress":"172.20.0.4"}}}},
		{"Id":"deadbeef","Names":["/standalone"],"Image":"alpine:3.19","State":"exited","Status":"Exited (0) 2 hours ago"}
	]`
	volumesJSON = `{"Volumes":[
		{"Name":"blog_data","Driver":"local","Mountpoint":"/var/lib/docker/volumes/blog_data/_data"},
		{"Name":"cache","Driver":"local","Mountpoint":"/var/lib/docker/volumes/cache/_data"}
	]}`
	networksJSON = `[
		{"Id":"net1","Name":"bridge","Driver":"bridge","Scope":"local","Internal":false},
		{"Id":"net2","Name":"blog_default","Driver":"bridge","Scope":"local","Internal":true}
	]`
)

// portainerAPI returns an httptest server serving the full happy-path API,
// including the Docker proxy for the two reachable Docker environments. If
// observe is non-nil it is called for every request and may reject it,
// which makes the handler reply 401.
func portainerAPI(t *testing.T, observe func(*http.Request) bool) *httptest.Server {
	t.Helper()
	body := map[string]string{
		pathSystemStatus: systemStatusJSON,
		pathEndpoints:    endpointsJSON,
		pathStacks:       stacksJSON,
	}
	for _, id := range []int{1, 2} {
		body[dockerPath(id, dockerContainers)] = containersJSON
		body[dockerPath(id, dockerVolumes)] = volumesJSON
		body[dockerPath(id, dockerNetworks)] = networksJSON
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if observe != nil && !observe(r) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"Unauthorized"}`))
			return
		}
		payload, ok := body[r.URL.RequestURI()]
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
// the config plumbing (api_key, verify_tls) is exercised too.
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
	if schema.Category != "containers_paas" || schema.Name != "Portainer" || schema.Stub {
		t.Fatalf("schema = %+v", schema)
	}
	want := map[string]bool{"url": true, "api_key": true, "verify_tls": false}
	got := make(map[string]bool, len(schema.Fields))
	for _, f := range schema.Fields {
		got[f.Key] = f.Required
	}
	if len(got) != len(want) {
		t.Errorf("schema fields = %+v, want exactly %+v", got, want)
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

// TestAPIKeyIsStoredAsPassword guards the token against being rendered as a
// plain-text field in the UI or stored unencrypted.
func TestAPIKeyIsStoredAsPassword(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	for _, f := range schema.Fields {
		if f.Key == "api_key" && f.Type != "password" {
			t.Fatalf("api_key field kind = %q, want %q", f.Type, "password")
		}
	}
}

func TestIdentity(t *testing.T) {
	c := &Connector{}
	if c.Name() != "Portainer" || c.Type() != "portainer" || c.Category() != "containers_paas" {
		t.Fatalf("identity = %s/%s/%s", c.Name(), c.Type(), c.Category())
	}
}

func TestValidate(t *testing.T) {
	server := portainerAPI(t, func(r *http.Request) bool {
		return r.Header.Get("X-API-Key") == testAPIKey
	})

	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{name: "valid", config: map[string]any{"url": server.URL, "api_key": testAPIKey}},
		{name: "trailing slash is trimmed", config: map[string]any{"url": server.URL + "/", "api_key": testAPIKey}},
		{name: "missing url", config: map[string]any{"api_key": testAPIKey}, wantErr: "url is required"},
		{name: "missing api key", config: map[string]any{"url": server.URL}, wantErr: "api_key is required"},
		{name: "wrong api key", config: map[string]any{"url": server.URL, "api_key": "nope"}, wantErr: "auth error"},
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

// TestValidateProbesAnAuthenticatedEndpoint pins Validate to /api/endpoints:
// /api/status is served unauthenticated and would accept any token.
func TestValidateProbesAnAuthenticatedEndpoint(t *testing.T) {
	var paths []string
	server := portainerAPI(t, func(r *http.Request) bool {
		paths = append(paths, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if len(paths) != 1 || paths[0] != pathEndpoints {
		t.Fatalf("Validate requested %v, want only %q", paths, pathEndpoints)
	}
}

func TestAPIKeyHeaderIsSent(t *testing.T) {
	var seen []string
	server := portainerAPI(t, func(r *http.Request) bool {
		seen = append(seen, r.Header.Get("X-API-Key"))
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
	if _, err := c.Fetch(context.Background(), nil); err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(seen) == 0 {
		t.Fatal("no requests observed")
	}
	for _, key := range seen {
		if key != testAPIKey {
			t.Fatalf("X-API-Key = %q, want %q on every request", key, testAPIKey)
		}
	}
}

func TestFetchHappyPath(t *testing.T) {
	server := portainerAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL + "/", "api_key": testAPIKey})

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.ServiceName != "Portainer" || snapshot.Type != typeName {
		t.Fatalf("snapshot = %s/%s", snapshot.ServiceName, snapshot.Type)
	}

	wantSections := []string{"Environments", "Stacks", "Containers", "Volumes", "Networks"}
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
		"portainer_url":     server.URL,
		"portainer_version": "2.19.4",
		"instance_id":       "b1c2d3e4-0000-4000-8000-abcdefabcdef",
		"environment_count": "4",
		"stack_count":       "3",
		"containers_count":  "4", // 2 containers across the 2 reachable environments
		"volumes_count":     "4",
		"networks_count":    "4",
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
	want := map[string]int{"environment": 4, "stack": 3, "container": 4, "volume": 4, "network": 4}
	for kind, n := range want {
		if byKind[kind] != n {
			t.Errorf("entities of kind %q = %d, want %d", kind, byKind[kind], n)
		}
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "host", Name: "edge-nas"},
		{Kind: "host", Name: "k8s-prod"},
		{Kind: "host", Name: "local"},
		{Kind: "host", Name: "offline-agent"},
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

// TestFetchQueriesOnlyDockerEnvironments checks the Docker proxy is not
// called for Kubernetes or unreachable environments.
func TestFetchQueriesOnlyDockerEnvironments(t *testing.T) {
	var hit []string
	server := portainerAPI(t, func(r *http.Request) bool {
		hit = append(hit, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
	if _, err := c.Fetch(context.Background(), nil); err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	for _, p := range hit {
		for _, forbidden := range []string{"/api/endpoints/3/docker", "/api/endpoints/4/docker"} {
			if strings.HasPrefix(p, forbidden) {
				t.Errorf("unexpected Docker proxy request to %q", p)
			}
		}
	}
	for _, wantPath := range []string{"/api/endpoints/1/docker/containers/json", "/api/endpoints/2/docker/networks"} {
		found := false
		for _, p := range hit {
			if p == wantPath {
				found = true
			}
		}
		if !found {
			t.Errorf("missing request to %q, got %v", wantPath, hit)
		}
	}
}

// TestFetchSelectiveFields checks the "fields" hint limits which endpoints
// are hit and which sections come back.
func TestFetchSelectiveFields(t *testing.T) {
	var hit []string
	server := portainerAPI(t, func(r *http.Request) bool {
		hit = append(hit, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"stacks"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "Stacks" {
		t.Fatalf("sections = %+v, want only Stacks", snapshot.Sections)
	}
	for _, p := range hit {
		if p == pathEndpoints || strings.Contains(p, "/docker/") {
			t.Errorf("unexpected request to %q for a stacks-only fetch", p)
		}
	}
	if len(snapshot.Dependencies) != 0 {
		t.Errorf("dependencies = %+v, want none when environments are not fetched", snapshot.Dependencies)
	}
}

// TestFetchContainersStillFetchesEnvironments checks the environment list is
// pulled as the fan-out source even when its own section wasn't requested.
func TestFetchContainersStillFetchesEnvironments(t *testing.T) {
	server := portainerAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"containers"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "Containers" {
		t.Fatalf("sections = %+v, want only Containers", snapshot.Sections)
	}
	if !strings.Contains(snapshot.Sections[0].Content, "blog-web") {
		t.Errorf("Containers section = %q, want the containers of both environments", snapshot.Sections[0].Content)
	}
}

// TestFetchDegradesPerSection checks a failing endpoint only degrades its
// own section instead of failing the whole snapshot.
func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pathEndpoints:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(endpointsJSON))
		case r.URL.Path == pathStacks:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		case r.URL.RequestURI() == dockerPath(1, dockerContainers):
			_, _ = w.Write([]byte(containersJSON))
		case r.URL.RequestURI() == dockerPath(2, dockerContainers):
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("agent unreachable"))
		case strings.HasSuffix(r.URL.Path, "/volumes"):
			_, _ = w.Write([]byte("not json"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
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
	if strings.Contains(byTitle["Environments"], "unavailable") {
		t.Errorf("Environments should have succeeded: %q", byTitle["Environments"])
	}
	if !strings.Contains(byTitle["Stacks"], "API returned 500") {
		t.Errorf("Stacks = %q, want the upstream 500 reported", byTitle["Stacks"])
	}
	if !strings.Contains(byTitle["Containers"], "service unavailable") || !strings.Contains(byTitle["Containers"], "edge-nas") {
		t.Errorf("Containers = %q, want the edge-nas failure noted", byTitle["Containers"])
	}
	if !strings.Contains(byTitle["Containers"], "blog-web") {
		t.Errorf("Containers = %q, want the reachable environment's rows kept", byTitle["Containers"])
	}
	if !strings.Contains(byTitle["Volumes"], "malformed response") {
		t.Errorf("Volumes = %q, want a malformed-response note", byTitle["Volumes"])
	}
	if !strings.Contains(byTitle["Networks"], "API returned 404") {
		t.Errorf("Networks = %q, want the upstream 404 reported", byTitle["Networks"])
	}
	if _, ok := snapshot.Metadata["portainer_version"]; ok {
		t.Errorf("portainer_version should be absent when the status endpoints 404: %+v", snapshot.Metadata)
	}
}

// TestFetchDegradesWhenEnvironmentListFails checks every environment-derived
// section degrades together, without failing Fetch.
func TestFetchDegradesWhenEnvironmentListFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pathStacks {
			_, _ = w.Write([]byte(stacksJSON))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	for _, s := range snapshot.Sections {
		if s.Title == "Stacks" {
			continue
		}
		if !strings.Contains(s.Content, "API returned 500") {
			t.Errorf("section %q = %q, want the environment-list failure reported", s.Title, s.Content)
		}
	}
	if len(snapshot.Dependencies) != 0 {
		t.Errorf("dependencies = %+v, want none", snapshot.Dependencies)
	}
}

// TestStatusFallsBackToLegacyPath covers Portainer < 2.18, which only
// serves /api/status.
func TestStatusFallsBackToLegacyPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathSystemStatus:
			w.WriteHeader(http.StatusNotFound)
		case pathStatus:
			_, _ = w.Write([]byte(`{"Version":"2.11.1","InstanceID":"legacy"}`))
		default:
			_, _ = w.Write([]byte("[]"))
		}
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.Metadata["portainer_version"] != "2.11.1" || snapshot.Metadata["instance_id"] != "legacy" {
		t.Fatalf("metadata = %+v, want the legacy /api/status values", snapshot.Metadata)
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

			c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
			err := c.Validate(context.Background(), nil)
			if err == nil {
				t.Fatal("Validate() error = nil, want an error")
			}
			tt.assert(t, err)
		})
	}
}

func TestExpiredDeadlineMapsToTimeout(t *testing.T) {
	server := portainerAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})

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

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("Validate() error = %v, want the response-size cap to trip", err)
	}
}

// TestGuardedDialerBlocksLoopback confirms the SSRF guard is wired in when
// the test override is not active.
func TestGuardedDialerBlocksLoopback(t *testing.T) {
	server := portainerAPI(t, nil)
	impl, err := connector.Get(typeName, map[string]any{"url": server.URL, "api_key": testAPIKey})
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
		if r.URL.Path == pathEndpoints {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey})
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
		_, _ = w.Write([]byte(endpointsJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey, "verify_tls": false})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() with verify_tls=false error = %v, want nil", err)
	}

	strict := newTestConnector(t, map[string]any{"url": server.URL, "api_key": testAPIKey, "verify_tls": true})
	if err := strict.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() with verify_tls=true error = nil, want a certificate error")
	}
}
