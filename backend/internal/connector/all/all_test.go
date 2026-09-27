package all

import (
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/connector/connectortest"
)

func TestAllConnectorImplementationsRegister(t *testing.T) {
	tests := []struct {
		typ      string
		category string
		stub     bool
	}{
		{typ: "adguardhome", category: "dns"},
		{typ: "custom", category: "virtualization"},
		{typ: "docker", category: "containers_paas"},
		{typ: "home_assistant", category: "virtualization"},
		{typ: "opnsense", category: "networking"},
		{typ: "pfsense", category: "networking"},
		{typ: "portainer", category: "containers_paas"},
		{typ: "pihole", category: "dns"},
		{typ: "proxmox", category: "virtualization"},
		{typ: "traefik", category: "networking"},
		{typ: "truenas", category: "virtualization"},
		{typ: "unifi", category: "networking"},
	}
	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			schema, err := connector.GetTypeSchema(tt.typ)
			if err != nil {
				t.Fatalf("GetTypeSchema: %v", err)
			}
			if schema.Category != tt.category || schema.Stub != tt.stub {
				t.Fatalf("schema = %+v", schema)
			}
			implementation, err := connector.Get(tt.typ, map[string]any{})
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			if implementation.Type() != tt.typ || implementation.Category() != tt.category {
				t.Fatalf("implementation = %s/%s", implementation.Type(), implementation.Category())
			}
		})
	}
}

func TestConnectorFailureContract(t *testing.T) {
	tests := []struct {
		typ    string
		config map[string]any
		opaque bool
	}{
		{"adguardhome", map[string]any{"auth_mode": "none"}, false},
		{"custom", nil, true},
		{"dnsresolver", map[string]any{"api_key": "bad"}, false},
		{"home_assistant", map[string]any{"access_token": "bad"}, false},
		{"netbird", map[string]any{"api_token": "bad"}, false},
		{"opnsense", map[string]any{"api_key": "bad", "api_secret": "bad"}, false},
		{"pfsense", map[string]any{"api_key": "bad"}, false},
		{"pihole", map[string]any{"password": "bad", "api_version": "v6"}, false},
		{"portainer", map[string]any{"api_key": "bad"}, false},
		{"proxmox", map[string]any{"token_id": "bad", "token_secret": "bad"}, false},
		{"traefik", map[string]any{"auth_mode": "none"}, false},
		{"truenas", map[string]any{"api_key": "bad"}, false},
		{"unifi", map[string]any{"auth_mode": "api_key", "api_key": "bad", "controller_type": "unifi_os"}, false},
	}
	covered := map[string]bool{"cloudflare": true, "docker": true} // suites in their packages need access to private client fields
	for _, tt := range tests {
		covered[tt.typ] = true
		t.Run(tt.typ, func(t *testing.T) {
			connectortest.Run(t, func(serverURL string) (connector.Connector, map[string]any, error) {
				cfg := map[string]any{"url": serverURL}
				for k, v := range tt.config {
					cfg[k] = v
				}
				c, err := connector.Get(tt.typ, cfg)
				return c, cfg, err
			}, tt.opaque)
		})
	}
	for _, schema := range connector.ListSchemas() {
		if !covered[schema.Type] {
			t.Errorf("connector %q has no failure contract suite", schema.Type)
		}
	}
}

func TestConnectorCapabilitiesMatchOptionalInterfaces(t *testing.T) {
	for _, schema := range connector.ListSchemas() {
		t.Run(schema.Type, func(t *testing.T) {
			implementation, err := connector.Get(schema.Type, map[string]any{})
			if err != nil {
				t.Fatal(err)
			}
			got := schema.Capabilities
			_, restart := implementation.(connector.Restarter)
			_, start := implementation.(connector.Starter)
			_, stop := implementation.(connector.Stopper)
			_, push := implementation.(connector.ConfigPusher)
			_, refresh := implementation.(connector.CredentialRefresher)
			want := connector.CapabilityDescriptor{Restart: restart, Start: start, Stop: stop, ConfigPush: push, CredentialRefresh: refresh}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("capabilities = %+v, want %+v", got, want)
			}
			if schema.IsCredentialRefresher != got.CredentialRefresh || len(schema.LifecycleVerbs) != countLifecycle(got) {
				t.Fatalf("legacy schema fields disagree with capabilities: %+v", schema)
			}
		})
	}
}

func countLifecycle(c connector.CapabilityDescriptor) int {
	n := 0
	for _, supported := range []bool{c.Restart, c.Start, c.Stop} {
		if supported {
			n++
		}
	}
	return n
}
