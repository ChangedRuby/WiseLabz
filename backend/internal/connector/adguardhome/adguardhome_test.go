package adguardhome

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
	statusJSON = `{
		"version":"v0.107.52","language":"en","dns_addresses":["192.168.1.2"],
		"dns_port":53,"http_port":3000,"protection_enabled":true,"running":true,"dhcp_available":true
	}`
	statusNoDHCPJSON = `{"version":"v0.107.52","dns_port":53,"http_port":3000,"protection_enabled":true,"running":true,"dhcp_available":false}`
	dnsInfoJSON      = `{
		"upstream_dns":["https://dns.quad9.net/dns-query","tls://1.1.1.1"],
		"bootstrap_dns":["9.9.9.10"],"fallback_dns":["8.8.8.8"],
		"local_ptr_upstreams":["192.168.1.1"],
		"upstream_mode":"parallel","blocking_mode":"custom_ip","blocking_ipv4":"0.0.0.0","blocking_ipv6":"::",
		"ratelimit":20,"cache_size":4194304,"edns_cs_enabled":false,"dnssec_enabled":true,
		"disable_ipv6":false,"resolve_clients":true,"protection_enabled":true
	}`
	filteringJSON = `{
		"enabled":true,"interval":24,
		"filters":[
			{"id":1,"name":"AdGuard DNS filter","url":"https://adguardteam.github.io/HostlistsRegistry/assets/filter_1.txt","rules_count":54321,"enabled":true},
			{"id":2,"name":"Dead hosts","url":"https://example.com/dead.txt","rules_count":0,"enabled":false}
		],
		"whitelist_filters":[
			{"id":3,"name":"Local allowlist","url":"https://example.com/allow.txt","rules_count":12,"enabled":true}
		],
		"user_rules":["||ads.example.com^","@@||cdn.example.com^",""]
	}`
	rewritesJSON = `[
		{"domain":"nas.example.com","answer":"192.168.1.10"},
		{"domain":"*.lab.example.com","answer":"proxy.example.com"}
	]`
	clientsJSON = `{
		"clients":[
			{"name":"laptop","ids":["aa:bb:cc:dd:ee:ff","192.168.1.50"],"tags":["user_admin"],
			 "blocked_services":["facebook"],"upstreams":["1.1.1.1"],
			 "use_global_settings":false,"filtering_enabled":true,"safebrowsing_enabled":true,"parental_enabled":false},
			{"name":"tv","ids":["192.168.1.60"],"use_global_settings":true,
			 "filtering_enabled":false,"safebrowsing_enabled":false,"parental_enabled":false}
		],
		"auto_clients":[{"name":"phone","ip":"192.168.1.77","source":"DHCP"}]
	}`
	dhcpJSON = `{
		"enabled":true,"interface_name":"eth0",
		"v4":{"gateway_ip":"192.168.1.1","subnet_mask":"255.255.255.0","range_start":"192.168.1.100","range_end":"192.168.1.200","lease_duration":86400},
		"leases":[{"mac":"11:22:33:44:55:66","ip":"192.168.1.101","hostname":"phone","expires":"2026-09-22T10:00:00Z"}],
		"static_leases":[{"mac":"aa:bb:cc:dd:ee:ff","ip":"192.168.1.50","hostname":"laptop"}]
	}`
)

// adguardAPI returns an httptest server serving the full happy-path API. If
// authCheck is non-nil it gates every request; returning false makes the
// handler reply 401.
func adguardAPI(t *testing.T, authCheck func(*http.Request) bool) *httptest.Server {
	t.Helper()
	body := map[string]string{
		pathStatus:          statusJSON,
		pathDNSInfo:         dnsInfoJSON,
		pathFilteringStatus: filteringJSON,
		pathRewriteList:     rewritesJSON,
		pathClients:         clientsJSON,
		pathDHCPStatus:      dhcpJSON,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authCheck != nil && !authCheck(r) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"forbidden"}`))
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
	if schema.Category != "dns" || schema.Name != "AdGuard Home" || schema.Stub {
		t.Fatalf("schema = %+v", schema)
	}
	want := map[string]bool{"url": true, "auth_mode": false, "username": false, "password": false, "verify_tls": false}
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
	if err := connector.ValidateConfig(*schema, map[string]any{"auth_mode": "token"}); err == nil {
		t.Error("ValidateConfig(token) = nil, want an error")
	}
}

func TestIdentity(t *testing.T) {
	c := &Connector{}
	if c.Name() != "AdGuard Home" || c.Type() != "adguardhome" || c.Category() != "dns" {
		t.Fatalf("identity = %s/%s/%s", c.Name(), c.Type(), c.Category())
	}
}

func TestValidateAuthModes(t *testing.T) {
	server := adguardAPI(t, func(r *http.Request) bool {
		user, pass, ok := r.BasicAuth()
		if !ok {
			return true // installs without a configured user are allowed
		}
		return user == "admin" && pass == "hunter2"
	})

	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{name: "basic auth", config: map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "admin", "password": "hunter2"}},
		{name: "explicit none", config: map[string]any{"url": server.URL, "auth_mode": authNone}},
		{
			name:    "basic is the default and needs credentials",
			config:  map[string]any{"url": server.URL},
			wantErr: "basic auth requires",
		},
		{
			name:    "basic auth missing password",
			config:  map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "admin"},
			wantErr: "basic auth requires",
		},
		{
			name:    "unknown auth mode",
			config:  map[string]any{"url": server.URL, "auth_mode": "token"},
			wantErr: "unknown auth mode",
		},
		{name: "missing url", config: map[string]any{"auth_mode": authNone}, wantErr: "url is required"},
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

// TestBasicAuthHeaderIsSent asserts the exact wire format of the
// authenticated mode.
func TestBasicAuthHeaderIsSent(t *testing.T) {
	var gotAuth string
	server := adguardAPI(t, func(r *http.Request) bool {
		gotAuth = r.Header.Get("Authorization")
		return true
	})

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authBasic, "username": "ad:min", "password": "p@ss"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("ad:min:p@ss"))
	if gotAuth != want {
		t.Errorf("Authorization = %q, want %q", gotAuth, want)
	}
}

// TestNoAuthModeSendsNoAuthorizationHeader guards against leaking stray
// credentials when the instance has no user configured.
func TestNoAuthModeSendsNoAuthorizationHeader(t *testing.T) {
	var seen []string
	server := adguardAPI(t, func(r *http.Request) bool {
		seen = append(seen, r.Header.Get("Authorization"))
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone, "username": "admin", "password": "hunter2"})
	if _, err := c.Fetch(context.Background(), nil); err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(seen) == 0 {
		t.Fatal("no requests recorded")
	}
	for _, h := range seen {
		if h != "" {
			t.Fatalf("Authorization header = %q, want empty in auth mode none", h)
		}
	}
}

func TestFetchHappyPath(t *testing.T) {
	server := adguardAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL + "/", "auth_mode": authNone})

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.ServiceName != "AdGuard Home" || snapshot.Type != typeName {
		t.Fatalf("snapshot = %s/%s", snapshot.ServiceName, snapshot.Type)
	}

	wantSections := []string{"Status", "DNS Configuration", "Filter Lists", "Custom Filtering Rules", "DNS Rewrites", "Clients", "DHCP"}
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
		"adguard_url":        server.URL,
		"adguard_version":    "v0.107.52",
		"protection_enabled": "true",
		"running":            "true",
		"dns_port":           "53",
		"dhcp_available":     "true",
		"upstream_count":     "2",
		"dnssec_enabled":     "true",
		"blocking_mode":      "custom_ip",
		"upstream_mode":      "parallel",
		"filtering_enabled":  "true",
		"filter_list_count":  "3",
		"user_rule_count":    "2",
		"rewrite_count":      "2",
		"client_count":       "2",
		"dhcp_enabled":       "true",
		"dhcp_lease_count":   "1",
		"dhcp_static_leases": "1",
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
	want := map[string]int{"filter_list": 3, "dns_rewrite": 2, "client": 2, "dhcp_lease": 2}
	for kind, n := range want {
		if byKind[kind] != n {
			t.Errorf("entities of kind %q = %d, want %d", kind, byKind[kind], n)
		}
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "upstream_service", Name: "https://dns.quad9.net/dns-query"},
		{Kind: "upstream_service", Name: "tls://1.1.1.1"},
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
	server := adguardAPI(t, func(r *http.Request) bool {
		hit = append(hit, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"rewrites"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "DNS Rewrites" {
		t.Fatalf("sections = %+v, want only DNS Rewrites", snapshot.Sections)
	}
	for _, p := range hit {
		if p != pathRewriteList {
			t.Errorf("unexpected request to %q for a rewrites-only fetch", p)
		}
	}
}

// TestFetchDHCPFieldStillReadsStatus checks the DHCP section works on its
// own: /control/status is consulted for dhcp_available even when the Status
// section itself was not requested.
func TestFetchDHCPFieldStillReadsStatus(t *testing.T) {
	var hit []string
	server := adguardAPI(t, func(r *http.Request) bool {
		hit = append(hit, r.URL.Path)
		return true
	})
	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"dhcp"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "DHCP" {
		t.Fatalf("sections = %+v, want only DHCP", snapshot.Sections)
	}
	want := map[string]bool{pathStatus: true, pathDHCPStatus: true}
	for _, p := range hit {
		if !want[p] {
			t.Errorf("unexpected request to %q for a dhcp-only fetch", p)
		}
	}
}

// TestFetchSkipsDHCPWhenUnavailable checks builds without DHCP support are
// not reported as a degraded section.
func TestFetchSkipsDHCPWhenUnavailable(t *testing.T) {
	var hitDHCP bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pathDHCPStatus {
			hitDHCP = true
		}
		if r.URL.Path != pathStatus {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(statusNoDHCPJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"status", "dhcp"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if hitDHCP {
		t.Error("/control/dhcp/status was requested although dhcp_available is false")
	}
	if len(snapshot.Sections) != 1 || snapshot.Sections[0].Title != "Status" {
		t.Fatalf("sections = %+v, want only Status", snapshot.Sections)
	}
	if snapshot.Metadata["dhcp_available"] != "false" {
		t.Errorf("metadata[dhcp_available] = %q, want \"false\"", snapshot.Metadata["dhcp_available"])
	}
}

// TestFetchDegradesPerSection checks a failing endpoint only degrades its
// own section instead of failing the whole snapshot.
func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathStatus:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(statusJSON))
		case pathDNSInfo:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
		case pathFilteringStatus:
			_, _ = w.Write([]byte("not json"))
		case pathDHCPStatus:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(dhcpJSON))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	byTitle := map[string]string{}
	for _, s := range snapshot.Sections {
		byTitle[s.Title] = s.Content
	}
	if len(snapshot.Sections) != 7 {
		t.Fatalf("sections = %d, want 7: %+v", len(snapshot.Sections), byTitle)
	}
	if strings.Contains(byTitle["Status"], "unavailable") {
		t.Errorf("Status should have succeeded: %q", byTitle["Status"])
	}
	if !strings.Contains(byTitle["DNS Configuration"], "API returned 500") {
		t.Errorf("DNS Configuration = %q, want the upstream 500 reported", byTitle["DNS Configuration"])
	}
	if !strings.Contains(byTitle["Filter Lists"], "malformed response") {
		t.Errorf("Filter Lists = %q, want a malformed-response placeholder", byTitle["Filter Lists"])
	}
	if !strings.Contains(byTitle["Custom Filtering Rules"], "malformed response") {
		t.Errorf("Custom Filtering Rules = %q, want a malformed-response placeholder", byTitle["Custom Filtering Rules"])
	}
	if !strings.Contains(byTitle["Clients"], "API returned 404") {
		t.Errorf("Clients = %q, want the upstream 404 reported", byTitle["Clients"])
	}
	if strings.Contains(byTitle["DHCP"], "unavailable") {
		t.Errorf("DHCP should have succeeded: %q", byTitle["DHCP"])
	}
	if len(snapshot.Dependencies) != 0 {
		t.Errorf("dependencies = %+v, want none when dns_info fails", snapshot.Dependencies)
	}
}

// TestFetchDegradesWhenStatusFails checks a failing /control/status
// degrades only the Status section and disables the DHCP section, since
// DHCP availability is then unknown.
func TestFetchDegradesWhenStatusFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pathStatus {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rewritesJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"status", "dhcp", "rewrites"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 2 {
		t.Fatalf("sections = %+v, want Status and DNS Rewrites", snapshot.Sections)
	}
	if !strings.Contains(snapshot.Sections[0].Content, "unavailable") {
		t.Errorf("Status = %q, want it reported unavailable", snapshot.Sections[0].Content)
	}
	if snapshot.Sections[1].Title != "DNS Rewrites" {
		t.Errorf("sections[1].Title = %q, want DNS Rewrites", snapshot.Sections[1].Title)
	}
}

// TestFetchOnMalformedStatus checks a status body that isn't JSON degrades
// the section without tripping up the rest of the snapshot.
func TestFetchOnMalformedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"status", "dhcp"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || !strings.Contains(snapshot.Sections[0].Content, "malformed response") {
		t.Fatalf("sections = %+v, want a single malformed Status section", snapshot.Sections)
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

			c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
			err := c.Validate(context.Background(), nil)
			if err == nil {
				t.Fatal("Validate() error = nil, want an error")
			}
			tt.assert(t, err)
		})
	}
}

func TestExpiredDeadlineMapsToTimeout(t *testing.T) {
	server := adguardAPI(t, nil)
	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})

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

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("Validate() error = %v, want the response-size cap to trip", err)
	}
}

// TestGuardedDialerBlocksLoopback confirms the SSRF guard is wired in when
// the test override is not active.
func TestGuardedDialerBlocksLoopback(t *testing.T) {
	server := adguardAPI(t, nil)
	impl, err := connector.Get(typeName, map[string]any{"url": server.URL, "auth_mode": authNone})
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
		if r.URL.Path == pathStatus {
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone})
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
		_, _ = w.Write([]byte(statusJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone, "verify_tls": false})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() with verify_tls=false error = %v, want nil", err)
	}

	strict := newTestConnector(t, map[string]any{"url": server.URL, "auth_mode": authNone, "verify_tls": true})
	if err := strict.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() with verify_tls=true error = nil, want a certificate error")
	}
}
