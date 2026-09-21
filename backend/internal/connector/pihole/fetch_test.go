package pihole

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// v6Server is a fake Pi-hole v6 REST API covering every endpoint Fetch
// touches. Paths not listed answer 404, so a connector that asks for the
// wrong path fails loudly.
func v6Server(t *testing.T) *httptest.Server {
	t.Helper()
	bodies := map[string]string{
		"/api/config/dns/hosts": `{"config":{"dns":{"hosts":["10.0.0.5 nas.internal.example.com"]}}}`,
		"/api/groups":           `{"groups":[{"id":0,"name":"Default","comment":"default","enabled":true},{"id":1,"name":"Kids","enabled":true}]}`,
		"/api/lists":            `{"lists":[{"address":"https://lists.example.com/ads","type":"block","comment":"ads","enabled":true,"groups":[0]}]}`,
		"/api/clients":          `{"clients":[{"client":"10.0.0.20","comment":"tablet","groups":[1]}]}`,
		"/api/domains":          `{"domains":[{"domain":"ads.example.com","type":"deny","kind":"exact","enabled":true,"groups":[0]}]}`,
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/auth" {
			_, _ = w.Write([]byte(`{"session":{"sid":"abc","valid":true}}`))
			return
		}
		body, ok := bodies[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("sid") != "abc" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
}

// v5Server is a fake Pi-hole v5 instance: no /api/auth, PHP endpoints keyed
// by query action, token-checked the way v5 does it (empty array on refusal).
func v5Server(t *testing.T) *httptest.Server {
	t.Helper()
	actions := map[string]string{
		"get_groups":  `{"data":[{"id":0,"name":"Default","description":"default","enabled":1},{"id":1,"name":"Kids","enabled":1}]}`,
		"get_adlists": `{"data":[{"id":1,"address":"https://lists.example.com/ads","comment":"ads","enabled":1,"groups":[0]}]}`,
		"get_clients": `{"data":[{"id":1,"ip":"10.0.0.20","comment":"tablet","groups":[1]}]}`,
		"get_domains": `{"data":[{"id":1,"domain":"ads.example.com","type":1,"enabled":1,"groups":[0]}]}`,
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("auth") != "token" && q.Get("token") != "token" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		switch r.URL.Path {
		case v5APIPHP:
			if _, ok := q["summaryRaw"]; ok {
				_, _ = w.Write([]byte(`{"status":"enabled","domains_being_blocked":1000}`))
				return
			}
			_, _ = w.Write([]byte(`{"status":"enabled"}`))
		case v5CustomDNSPHP:
			_, _ = w.Write([]byte(`{"data":[["10.0.0.5","nas.internal.example.com"]]}`))
		case v5GroupsPHP:
			body, ok := actions[q.Get("action")]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, _ = w.Write([]byte(body))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func sectionByTitle(t *testing.T, snap *connector.ServiceSnapshot, title string) string {
	t.Helper()
	for _, s := range snap.Sections {
		if s.Title == title {
			return s.Content
		}
	}
	t.Fatalf("snapshot has no section %q (sections: %+v)", title, snap.Sections)
	return ""
}

func entityKinds(snap *connector.ServiceSnapshot) map[string]int {
	counts := map[string]int{}
	for _, e := range snap.Entities {
		counts[e.Kind]++
	}
	return counts
}

func TestFetchBothVersions(t *testing.T) {
	tests := []struct {
		name        string
		server      func(*testing.T) *httptest.Server
		apiVersion  string
		password    string
		wantVersion string
	}{
		{name: "v6 pinned", server: v6Server, apiVersion: version6, password: "secret", wantVersion: version6},
		{name: "v6 auto-detected", server: v6Server, apiVersion: versionAuto, password: "secret", wantVersion: version6},
		{name: "v5 pinned", server: v5Server, apiVersion: version5, password: "token", wantVersion: version5},
		{name: "v5 auto-detected", server: v5Server, apiVersion: versionAuto, password: "token", wantVersion: version5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.server(t)
			defer server.Close()

			c := &Connector{url: server.URL, password: tt.password, apiVersion: tt.apiVersion, client: server.Client()}
			snap, err := c.Fetch(context.Background(), nil)
			if err != nil {
				t.Fatalf("Fetch() error = %v", err)
			}
			if got := snap.Metadata["pihole_api_version"]; got != tt.wantVersion {
				t.Errorf("pihole_api_version = %q, want %q", got, tt.wantVersion)
			}

			wantContent := map[string]string{
				"Local DNS Records": "nas.internal.example.com",
				"Groups":            "Default",
				"Blocklists":        "https://lists.example.com/ads",
				"Clients":           "10.0.0.20",
				"Domain Rules":      "ads.example.com",
			}
			for title, want := range wantContent {
				if content := sectionByTitle(t, snap, title); !strings.Contains(content, want) {
					t.Errorf("section %q = %q, want to contain %q", title, content, want)
				}
			}

			// Group ids resolve to names in the referencing sections.
			if content := sectionByTitle(t, snap, "Clients"); !strings.Contains(content, "Kids") {
				t.Errorf("Clients section = %q, want the resolved group name", content)
			}

			want := map[string]int{"dns_record": 1, "dns_group": 2, "blocklist": 1, "dns_client": 1, "domain_rule": 1}
			got := entityKinds(snap)
			for kind, n := range want {
				if got[kind] != n {
					t.Errorf("entity kind %q count = %d, want %d", kind, got[kind], n)
				}
			}
		})
	}
}

func TestFetchDegradesPerSection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/auth":
			_, _ = w.Write([]byte(`{"session":{"sid":"abc","valid":true}}`))
		case "/api/config/dns/hosts":
			_, _ = w.Write([]byte(`{"config":{"dns":{"hosts":["10.0.0.5 nas.internal.example.com"]}}}`))
		case "/api/groups":
			w.WriteHeader(http.StatusServiceUnavailable)
		case "/api/lists":
			_, _ = w.Write([]byte(`not json`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, password: "secret", apiVersion: version6, client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if content := sectionByTitle(t, snap, "Local DNS Records"); !strings.Contains(content, "nas.internal.example.com") {
		t.Errorf("healthy section did not survive a failing sibling: %q", content)
	}
	if content := sectionByTitle(t, snap, "Groups"); !strings.Contains(content, "service unavailable") {
		t.Errorf("Groups section = %q, want a service-unavailable placeholder", content)
	}
	if content := sectionByTitle(t, snap, "Blocklists"); !strings.Contains(content, "malformed response") {
		t.Errorf("Blocklists section = %q, want a malformed-response placeholder", content)
	}
	if content := sectionByTitle(t, snap, "Clients"); !strings.Contains(content, "unavailable") {
		t.Errorf("Clients section = %q, want an unavailable placeholder", content)
	}
}

func TestFetchAuthFailureReturnsPlaceholderSnapshot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"session":{"valid":false,"message":"bad password"}}`))
	}))
	defer server.Close()

	c := &Connector{url: server.URL, password: "nope", client: server.Client()}
	snap, err := c.Fetch(context.Background(), nil)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(snap.Sections) != 1 || !strings.Contains(snap.Sections[0].Content, "auth error") {
		t.Fatalf("sections = %+v, want a single auth-error placeholder", snap.Sections)
	}
	if _, ok := snap.Metadata["pihole_api_version"]; ok {
		t.Error("metadata reports an API version although no version answered")
	}
}

func TestValidateBothVersions(t *testing.T) {
	tests := []struct {
		name       string
		server     func(*testing.T) *httptest.Server
		apiVersion string
		password   string
		wantErr    bool
	}{
		{name: "v6", server: v6Server, apiVersion: version6, password: "secret"},
		{name: "v5", server: v5Server, apiVersion: version5, password: "token"},
		{name: "v5 bad token", server: v5Server, apiVersion: version5, password: "wrong", wantErr: true},
		{name: "v6 pinned against a v5 host", server: v5Server, apiVersion: version6, password: "token", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.server(t)
			defer server.Close()

			c := &Connector{url: server.URL, password: tt.password, apiVersion: tt.apiVersion, client: server.Client()}
			err := c.Validate(context.Background(), nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestV5AuthErrorOnEmptyArray(t *testing.T) {
	server := v5Server(t)
	defer server.Close()

	c := &Connector{url: server.URL, password: "wrong", apiVersion: version5, client: server.Client()}
	_, err := c.get(context.Background(), session{version: version5}, resourceGroups)
	var authErr *connector.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("get() error = %v, want *connector.AuthError", err)
	}
}

func TestV5UnsupportedResource(t *testing.T) {
	c := &Connector{url: "http://example.invalid", password: "token"}
	if _, err := c.v5URL("nope"); err == nil {
		t.Error("v5URL(unknown) error = nil, want error")
	}
}

func TestStartStopV5(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == v5APIPHP && r.URL.Query().Has("summaryRaw") {
			_, _ = w.Write([]byte(`{"status":"enabled"}`))
			return
		}
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"status":"disabled"}`))
	}))
	defer server.Close()

	c := &Connector{url: server.URL, password: "token", apiVersion: version5, client: server.Client()}
	if err := c.Stop(context.Background(), nil, ""); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if !strings.Contains(gotQuery, "disable") {
		t.Errorf("query = %q, want the disable action", gotQuery)
	}
	if err := c.Start(context.Background(), nil, ""); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if !strings.Contains(gotQuery, "enable") {
		t.Errorf("query = %q, want the enable action", gotQuery)
	}
}

func TestRestartUnsupportedOnV5(t *testing.T) {
	server := v5Server(t)
	defer server.Close()

	c := &Connector{url: server.URL, password: "token", apiVersion: version5, client: server.Client()}
	err := c.Restart(context.Background(), nil, "")
	if err == nil || !strings.Contains(err.Error(), "v5") {
		t.Fatalf("Restart() error = %v, want a v5-unsupported error", err)
	}
}

func TestConfigPushV5(t *testing.T) {
	var actions []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		switch {
		case r.URL.Path == v5APIPHP:
			_, _ = w.Write([]byte(`{"status":"enabled"}`))
		case r.URL.Path == v5CustomDNSPHP && q.Get("action") == "get":
			_, _ = w.Write([]byte(`{"data":[["10.0.0.5","host.example.com"]]}`))
		case r.URL.Path == v5CustomDNSPHP:
			actions = append(actions, q.Get("action")+" "+q.Get("ip")+" "+q.Get("domain"))
			_, _ = w.Write([]byte(`{"success":true}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, password: "token", apiVersion: version5, client: server.Client()}
	if err := c.ConfigPush(context.Background(), nil, "host.example.com", "ip", "10.0.0.9"); err != nil {
		t.Fatalf("ConfigPush() error = %v", err)
	}
	want := []string{"delete 10.0.0.5 host.example.com", "add 10.0.0.9 host.example.com"}
	if len(actions) != len(want) {
		t.Fatalf("actions = %v, want %v", actions, want)
	}
	for i, w := range want {
		if actions[i] != w {
			t.Errorf("actions[%d] = %q, want %q", i, actions[i], w)
		}
	}
}
