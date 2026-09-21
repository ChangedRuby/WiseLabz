package all

import (
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
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
		{typ: "opnsense", category: "networking"},
		{typ: "pfsense", category: "networking"},
		{typ: "portainer", category: "containers_paas"},
		{typ: "pihole", category: "dns"},
		{typ: "proxmox", category: "virtualization"},
		{typ: "traefik", category: "networking"},
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
