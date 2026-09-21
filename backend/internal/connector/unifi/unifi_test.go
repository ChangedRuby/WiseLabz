package unifi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const (
	sitesJSON = `{"meta":{"rc":"ok"},"data":[
		{"_id":"s1","name":"default","desc":"Home","role":"admin"},
		{"_id":"s2","name":"7xk2lp0q","desc":"Lab","role":"readonly"}
	]}`
	devicesJSON = `{"meta":{"rc":"ok"},"data":[
		{"mac":"aa:bb:cc:00:11:22","name":"Living Room AP","model":"U6LR","type":"uap","version":"6.6.55","ip":"10.0.1.20",
		 "adopted":true,"disabled":false,"state":1,"uptime":98213,"rx_bytes":123456789,"last_seen":1700000000},
		{"mac":"aa:bb:cc:00:11:33","name":"Core Switch","model":"US8P60","type":"usw","version":"6.5.59","ip":"10.0.1.21",
		 "adopted":true,"disabled":true,"state":"0"},
		{"mac":"aa:bb:cc:00:11:44","model":"UDMPRO","type":"udm","version":"3.2.12","ip":"10.0.1.1","adopted":true,"state":42}
	]}`
	networksJSON = `{"meta":{"rc":"ok"},"data":[
		{"_id":"n1","name":"LAN","purpose":"corporate","networkgroup":"LAN","vlan_enabled":false,"ip_subnet":"10.0.1.1/24","dhcpd_enabled":true,"enabled":true},
		{"_id":"n2","name":"IoT","purpose":"corporate","networkgroup":"LAN","vlan_enabled":true,"vlan":"20","ip_subnet":"10.0.20.1/24","dhcpd_enabled":true},
		{"_id":"n3","name":"WAN","purpose":"wan","networkgroup":"WAN","enabled":false}
	]}`
	wlansJSON = `{"meta":{"rc":"ok"},"data":[
		{"_id":"w1","name":"HomeNet","enabled":true,"security":"wpapsk","wpa_mode":"wpa2","is_guest":false,"hide_ssid":false,
		 "mac_filter_enabled":false,"pmf_mode":"optional","x_passphrase":"super-secret"},
		{"_id":"w2","name":"Guest","security":"wpapsk","wpa_mode":"wpa3","is_guest":true,"hide_ssid":true,
		 "mac_filter_enabled":true,"mac_filter_policy":"allow","pmf_mode":"required"}
	]}`
	firewallJSON = `{"meta":{"rc":"ok"},"data":[
		{"_id":"f2","name":"Block IoT to LAN","ruleset":"LAN_IN","rule_index":"2001","action":"drop","protocol":"all",
		 "enabled":true,"logging":true,"src_address":"10.0.20.0/24","dst_address":"10.0.1.0/24"},
		{"_id":"f1","name":"Allow established","ruleset":"LAN_IN","rule_index":2000,"action":"accept","protocol":"all","enabled":true},
		{"_id":"f3","name":"Drop WAN","ruleset":"WAN_IN","rule_index":3000,"action":"reject","protocol":"tcp_udp"}
	]}`
	clientsJSON = `{"meta":{"rc":"ok"},"data":[
		{"mac":"11:11:11:11:11:11","is_wired":false,"essid":"HomeNet","uptime":123,"rx_bytes":9999,"signal":-54},
		{"mac":"22:22:22:22:22:22","is_wired":false,"essid":"HomeNet"},
		{"mac":"33:33:33:33:33:33","is_wired":false,"essid":"Guest"},
		{"mac":"44:44:44:44:44:44","is_wired":true,"network":"LAN"}
	]}`
	sysinfoJSON = `{"meta":{"rc":"ok"},"data":[{"version":"8.0.28","uptime":42,"timezone":"Europe/Lisbon"}]}`
)

// controllerOpts configures the fake controller.
type controllerOpts struct {
	unifiOS  bool              // serve the /proxy/network prefix and the UniFi OS login
	apiKey   string            // when set, accept X-API-KEY instead of a session cookie
	username string            // credentials the login endpoint accepts
	password string            //
	site     string            // site name the Network API is served for, defaults to "default"
	onPath   func(path string) // observes every request path
}

// unifiAPI returns an httptest server mimicking a UniFi controller: a login
// endpoint that hands out a session cookie, and Network API paths that
// require it.
func unifiAPI(t *testing.T, opts controllerOpts) *httptest.Server {
	t.Helper()
	if opts.site == "" {
		opts.site = "default"
	}
	prefix := prefixClassic
	loginPath := pathLoginClass
	if opts.unifiOS {
		prefix = prefixUniFiOS
		loginPath = pathLoginUOS
	}
	body := map[string]string{
		pathSites:                                    sitesJSON,
		"/api/s/" + opts.site + "/stat/device":       devicesJSON,
		"/api/s/" + opts.site + "/rest/networkconf":  networksJSON,
		"/api/s/" + opts.site + "/rest/wlanconf":     wlansJSON,
		"/api/s/" + opts.site + "/rest/firewallrule": firewallJSON,
		"/api/s/" + opts.site + "/stat/sta":          clientsJSON,
		"/api/s/" + opts.site + "/stat/sysinfo":      sysinfoJSON,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if opts.onPath != nil {
			opts.onPath(r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == loginPath {
			var creds struct{ Username, Password string }
			_ = decodeJSONBody(r, &creds)
			if creds.Username != opts.username || creds.Password != opts.password {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.Invalid"}}`))
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "unifises", Value: "sess123", Path: "/"})
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			return
		}
		payload, ok := body[strings.TrimPrefix(r.URL.Path, prefix)]
		if !ok || (prefix != "" && !strings.HasPrefix(r.URL.Path, prefix)) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.NoSuchPath"}}`))
			return
		}
		if !authorized(r, opts) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.LoginRequired"}}`))
			return
		}
		_, _ = w.Write([]byte(payload))
	}))
	t.Cleanup(server.Close)
	return server
}

func authorized(r *http.Request, opts controllerOpts) bool {
	if opts.apiKey != "" {
		return r.Header.Get("X-API-KEY") == opts.apiKey
	}
	cookie, err := r.Cookie("unifises")
	return err == nil && cookie.Value == "sess123"
}

func decodeJSONBody(r *http.Request, out any) error {
	defer func() { _ = r.Body.Close() }()
	return json.NewDecoder(r.Body).Decode(out)
}

// newTestConnector builds the connector through the registered factory so
// the config plumbing (auth mode, site, verify_tls) is exercised too.
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

func passwordConfig(serverURL string) map[string]any {
	return map[string]any{"url": serverURL, "username": "admin", "password": "hunter2"}
}

func TestRegisteredSchema(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if schema.Category != "networking" || schema.Name != "UniFi" || schema.Stub {
		t.Fatalf("schema = %+v", schema)
	}
	want := map[string]bool{
		"url": true, "auth_mode": false, "username": false, "password": false,
		"api_key": false, "site": false, "controller_type": false, "verify_tls": false,
	}
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
	if len(got) != len(want) {
		t.Errorf("schema fields = %v, want exactly %v", got, want)
	}
	if _, ok := connector.AttributeCatalog()[typeName]; !ok {
		t.Errorf("attribute catalog not registered for %q", typeName)
	}
}

// TestSchemaConfigValidation checks the select fields are validated by the
// shared config validator.
func TestSchemaConfigValidation(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	if err := connector.ValidateConfig(*schema, map[string]any{"auth_mode": authAPIKey, "controller_type": controllerUniFiOS}); err != nil {
		t.Errorf("ValidateConfig(api_key/unifios) = %v, want nil", err)
	}
	if err := connector.ValidateConfig(*schema, map[string]any{"auth_mode": "kerberos"}); err == nil {
		t.Error("ValidateConfig(kerberos) = nil, want an error")
	}
	if err := connector.ValidateConfig(*schema, map[string]any{"controller_type": "cloud"}); err == nil {
		t.Error("ValidateConfig(cloud) = nil, want an error")
	}
}

func TestIdentity(t *testing.T) {
	c := &Connector{}
	if c.Name() != "UniFi" || c.Type() != "unifi" || c.Category() != "networking" {
		t.Fatalf("identity = %s/%s/%s", c.Name(), c.Type(), c.Category())
	}
}

// TestFactoryDefaults checks a stored config missing the optional fields
// still yields a usable connector.
func TestFactoryDefaults(t *testing.T) {
	c := newTestConnector(t, map[string]any{"url": "https://unifi.example.com/"})
	if c.url != "https://unifi.example.com" {
		t.Errorf("url = %q, want the trailing slash trimmed", c.url)
	}
	if c.site != "default" || c.authMode != authPassword || c.controllerType != controllerAuto {
		t.Errorf("defaults = %s/%s/%s", c.site, c.authMode, c.controllerType)
	}
	bogus := newTestConnector(t, map[string]any{"url": "x", "auth_mode": "kerberos", "controller_type": "cloud"})
	if bogus.authMode != authPassword || bogus.controllerType != controllerAuto {
		t.Errorf("unrecognised values = %s/%s, want the defaults", bogus.authMode, bogus.controllerType)
	}
}

func TestValidate(t *testing.T) {
	classic := unifiAPI(t, controllerOpts{username: "admin", password: "hunter2"})
	unifiOS := unifiAPI(t, controllerOpts{unifiOS: true, username: "admin", password: "hunter2"})
	keyed := unifiAPI(t, controllerOpts{unifiOS: true, apiKey: "key123"})

	tests := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{name: "classic auto-detected", config: passwordConfig(classic.URL)},
		{name: "unifi os auto-detected", config: passwordConfig(unifiOS.URL)},
		{
			name:   "classic pinned",
			config: map[string]any{"url": classic.URL, "username": "admin", "password": "hunter2", "controller_type": controllerClassic},
		},
		{
			name:   "unifi os pinned",
			config: map[string]any{"url": unifiOS.URL, "username": "admin", "password": "hunter2", "controller_type": controllerUniFiOS},
		},
		{
			name:   "api key",
			config: map[string]any{"url": keyed.URL, "auth_mode": authAPIKey, "api_key": "key123"},
		},
		{name: "missing url", config: map[string]any{}, wantErr: "url is required"},
		{
			name:    "password mode without password",
			config:  map[string]any{"url": classic.URL, "username": "admin"},
			wantErr: "password auth requires",
		},
		{
			name:    "api key mode without key",
			config:  map[string]any{"url": keyed.URL, "auth_mode": authAPIKey},
			wantErr: "requires an api_key",
		},
		{
			name:    "wrong password",
			config:  map[string]any{"url": classic.URL, "username": "admin", "password": "wrong"},
			wantErr: "auth error",
		},
		{
			name:    "wrong api key",
			config:  map[string]any{"url": keyed.URL, "auth_mode": authAPIKey, "api_key": "nope"},
			wantErr: "auth error",
		},
		{
			name:    "pinned to the wrong flavour",
			config:  map[string]any{"url": classic.URL, "username": "admin", "password": "hunter2", "controller_type": controllerUniFiOS},
			wantErr: "controller returned 404",
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

// TestAutoDetectReportsUniFiOSError checks that when neither flavour
// answers, the UniFi OS error is the one surfaced.
func TestAutoDetectReportsUniFiOSError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pathLoginUOS {
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"unifi-os-said-no"}}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"classic-said-no"}}`))
	}))
	defer server.Close()

	c := newTestConnector(t, passwordConfig(server.URL))
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "unifi-os-said-no") {
		t.Fatalf("Validate() error = %v, want the UniFi OS error", err)
	}
}

// TestAPIKeyIsNotSentToLogin guards against leaking the API key on a
// password-mode handshake and vice versa.
func TestAPIKeyIsNotSentInPasswordMode(t *testing.T) {
	var sawKey bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != "" {
			sawKey = true
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == pathLoginUOS {
			http.SetCookie(w, &http.Cookie{Name: "unifises", Value: "sess123", Path: "/"})
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			return
		}
		_, _ = w.Write([]byte(sitesJSON))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "username": "admin", "password": "hunter2", "api_key": "key123"})
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if sawKey {
		t.Error("X-API-KEY was sent in password auth mode")
	}
}

// TestSessionCookieIsReplayed checks the login cookie reaches the Network
// API calls, which is what makes password mode work at all.
func TestSessionCookieIsReplayed(t *testing.T) {
	server := unifiAPI(t, controllerOpts{username: "admin", password: "hunter2"})
	c := newTestConnector(t, passwordConfig(server.URL))

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"sites"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if strings.Contains(snapshot.Sections[0].Content, "unavailable") {
		t.Fatalf("Sites section = %q, want the session cookie to have been replayed", snapshot.Sections[0].Content)
	}
}

func TestFetchHappyPath(t *testing.T) {
	server := unifiAPI(t, controllerOpts{unifiOS: true, username: "admin", password: "hunter2"})
	c := newTestConnector(t, passwordConfig(server.URL+"/"))

	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if snapshot.ServiceName != "UniFi" || snapshot.Type != typeName {
		t.Fatalf("snapshot = %s/%s", snapshot.ServiceName, snapshot.Type)
	}

	wantSections := []string{"Sites", "Devices", "Networks", "WLANs", "Firewall Rules", "Clients"}
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
		"unifi_url":             server.URL,
		"unifi_site":            "default",
		"unifi_controller":      controllerUniFiOS,
		"unifi_version":         "8.0.28",
		"site_count":            "2",
		"device_count":          "3",
		"network_count":         "3",
		"wlan_count":            "2",
		"firewall_rule_count":   "3",
		"client_count":          "4",
		"wired_client_count":    "1",
		"wireless_client_count": "3",
	}
	for k, want := range wantMetadata {
		if got := snapshot.Metadata[k]; got != want {
			t.Errorf("metadata[%q] = %q, want %q", k, got, want)
		}
	}

	wantDeps := []connector.ServiceDependency{
		{Kind: "network", Name: "IoT"},
		{Kind: "network", Name: "LAN"},
		{Kind: "network", Name: "WAN"},
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

// TestFetchUsesConfiguredSite checks the site name lands in the request
// paths rather than being assumed to be "default".
func TestFetchUsesConfiguredSite(t *testing.T) {
	var paths []string
	server := unifiAPI(t, controllerOpts{
		username: "admin", password: "hunter2", site: "7xk2lp0q",
		onPath: func(p string) { paths = append(paths, p) },
	})
	config := passwordConfig(server.URL)
	config["site"] = "7xk2lp0q"
	c := newTestConnector(t, config)

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"devices"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if strings.Contains(snapshot.Sections[0].Content, "unavailable") {
		t.Fatalf("Devices = %q, want the configured site to have been used", snapshot.Sections[0].Content)
	}
	if snapshot.Metadata["unifi_site"] != "7xk2lp0q" {
		t.Errorf("metadata[unifi_site] = %q", snapshot.Metadata["unifi_site"])
	}
	var found bool
	for _, p := range paths {
		if p == "/api/s/7xk2lp0q/stat/device" {
			found = true
		}
	}
	if !found {
		t.Errorf("paths = %v, want a request to /api/s/7xk2lp0q/stat/device", paths)
	}
}

// TestFetchSelectiveFields checks the "fields" hint limits which endpoints
// are hit and which sections come back.
func TestFetchSelectiveFields(t *testing.T) {
	var paths []string
	server := unifiAPI(t, controllerOpts{
		username: "admin", password: "hunter2",
		onPath: func(p string) { paths = append(paths, p) },
	})
	c := newTestConnector(t, passwordConfig(server.URL))

	snapshot, err := c.Fetch(context.Background(), map[string]any{"fields": []string{"wlans", "firewall_rules"}})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 2 ||
		snapshot.Sections[0].Title != "WLANs" || snapshot.Sections[1].Title != "Firewall Rules" {
		t.Fatalf("sections = %+v, want only WLANs and Firewall Rules", snapshot.Sections)
	}
	for _, p := range paths {
		for _, unwanted := range []string{"/stat/device", "/rest/networkconf", "/stat/sta", pathSites} {
			if strings.HasSuffix(p, unwanted) {
				t.Errorf("unexpected request to %q for a wlans+firewall fetch", p)
			}
		}
	}
}

// TestFetchDegradesPerSection checks a failing endpoint only degrades its
// own section instead of failing the whole snapshot.
func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == pathLoginClass:
			http.SetCookie(w, &http.Cookie{Name: "unifises", Value: "sess123", Path: "/"})
			_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
		case r.URL.Path == pathSites:
			_, _ = w.Write([]byte(sitesJSON))
		case strings.HasSuffix(r.URL.Path, "/stat/device"):
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.Boom"}}`))
		case strings.HasSuffix(r.URL.Path, "/rest/networkconf"):
			_, _ = w.Write([]byte("not json"))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.NoSuchPath"}}`))
		}
	}))
	defer server.Close()

	c := newTestConnector(t, passwordConfig(server.URL))
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 6 {
		t.Fatalf("sections = %d, want 6", len(snapshot.Sections))
	}
	byTitle := map[string]string{}
	for _, s := range snapshot.Sections {
		byTitle[s.Title] = s.Content
	}
	if strings.Contains(byTitle["Sites"], "unavailable") {
		t.Errorf("Sites should have succeeded: %q", byTitle["Sites"])
	}
	if !strings.Contains(byTitle["Devices"], "api.err.Boom") {
		t.Errorf("Devices = %q, want the controller error message reported", byTitle["Devices"])
	}
	if !strings.Contains(byTitle["Networks"], "malformed response") {
		t.Errorf("Networks = %q, want a malformed-response placeholder", byTitle["Networks"])
	}
	if !strings.Contains(byTitle["WLANs"], "controller returned 404") {
		t.Errorf("WLANs = %q, want the upstream 404 reported", byTitle["WLANs"])
	}
	if _, ok := snapshot.Metadata["unifi_version"]; ok {
		t.Errorf("unifi_version should be absent when sysinfo 404s: %+v", snapshot.Metadata)
	}
}

// TestFetchWithoutSessionDegradesInsteadOfFailing checks an unreachable
// controller yields a snapshot describing the failure, like the sibling
// connectors, rather than an error.
func TestFetchWithoutSessionDegradesInsteadOfFailing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.LoginRequired"}}`))
	}))
	defer server.Close()

	c := newTestConnector(t, passwordConfig(server.URL))
	snapshot, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(snapshot.Sections) != 1 || !strings.Contains(snapshot.Sections[0].Content, "auth error") {
		t.Fatalf("sections = %+v, want a single auth-error placeholder", snapshot.Sections)
	}
	if len(snapshot.Entities) != 0 {
		t.Errorf("entities = %+v, want none", snapshot.Entities)
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
			name: "400 on login maps to AuthError", status: http.StatusBadRequest,
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
				var unavailableErr *connector.ServiceUnavailableError
				if errors.As(err, &authErr) || errors.As(err, &unavailableErr) {
					t.Fatalf("error = %v (%T), want a plain error", err, err)
				}
				if !strings.Contains(err.Error(), "controller returned 500") {
					t.Fatalf("error = %v, want it to mention the status", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"controller said no"}}`))
			}))
			defer server.Close()

			c := newTestConnector(t, map[string]any{
				"url": server.URL, "username": "admin", "password": "hunter2",
				"controller_type": controllerClassic,
			})
			err := c.Validate(context.Background(), nil)
			if err == nil {
				t.Fatal("Validate() error = nil, want an error")
			}
			tt.assert(t, err)
		})
	}
}

// TestControllerErrorMessageIsSurfaced checks the UniFi envelope's own
// message is reported instead of the raw JSON body.
func TestControllerErrorMessageIsSurfaced(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"meta":{"rc":"error","msg":"api.err.ServerBusy"}}`))
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "username": "a", "password": "b", "controller_type": controllerClassic})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "api.err.ServerBusy") || strings.Contains(err.Error(), `"meta"`) {
		t.Fatalf("Validate() error = %v, want just the controller message", err)
	}
}

func TestExpiredDeadlineMapsToTimeout(t *testing.T) {
	server := unifiAPI(t, controllerOpts{username: "admin", password: "hunter2"})
	c := newTestConnector(t, passwordConfig(server.URL))

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

	c := newTestConnector(t, map[string]any{"url": server.URL, "username": "a", "password": "b", "controller_type": controllerClassic})
	err := c.Validate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("Validate() error = %v, want the response-size cap to trip", err)
	}
}

// TestGuardedDialerBlocksLoopback confirms the SSRF guard is wired in when
// the test override is not active.
func TestGuardedDialerBlocksLoopback(t *testing.T) {
	server := unifiAPI(t, controllerOpts{username: "admin", password: "hunter2"})
	impl, err := connector.Get(typeName, passwordConfig(server.URL))
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
		if r.URL.Path == pathLoginClass {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	c := newTestConnector(t, map[string]any{"url": server.URL, "username": "a", "password": "b", "controller_type": controllerClassic})
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
		_, _ = w.Write([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
	}))
	defer server.Close()

	config := map[string]any{"url": server.URL, "username": "a", "password": "b", "controller_type": controllerClassic, "verify_tls": false}
	c := newTestConnector(t, config)
	if err := c.Validate(context.Background(), nil); err != nil {
		t.Fatalf("Validate() with verify_tls=false error = %v, want nil", err)
	}

	config["verify_tls"] = true
	strict := newTestConnector(t, config)
	if err := strict.Validate(context.Background(), nil); err == nil {
		t.Fatal("Validate() with verify_tls=true error = nil, want a certificate error")
	}
}
