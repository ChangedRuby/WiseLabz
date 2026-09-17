package cloudflare

import (
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildDNSRecordTableAttributes(t *testing.T) {
	data := []byte(`{"result":[
		{"id":"r1","name":"app.example.com","type":"A","content":"203.0.113.5","proxied":true,"ttl":1},
		{"id":"r2","name":"mail.example.com","type":"MX","content":"mail.provider.com","proxied":false,"ttl":3600}
	]}`)
	_, entities := buildDNSRecordTable("zone1", "example.com", data)
	if len(entities) != 2 {
		t.Fatalf("entities = %+v, want 2", entities)
	}
	if entities[0].ExternalID != "zone1/r1" {
		t.Errorf("entities[0].ExternalID = %q, want %q", entities[0].ExternalID, "zone1/r1")
	}
	if entities[0].Attributes["proxied"] != true || entities[0].Attributes["type"] != "A" {
		t.Errorf("entities[0].Attributes = %+v", entities[0].Attributes)
	}
	if entities[1].Attributes["proxied"] != false || entities[1].Attributes["ttl"] != 3600 {
		t.Errorf("entities[1].Attributes = %+v", entities[1].Attributes)
	}
}

func TestBuildTunnelTableAttributes(t *testing.T) {
	data := []byte(`{"result":[{"id":"t1","name":"prod-tunnel","status":"healthy","conns":2}]}`)
	_, entities := buildTunnelTable(data)
	if len(entities) != 1 {
		t.Fatalf("entities = %+v, want 1", entities)
	}
	if entities[0].Attributes["status"] != "healthy" || entities[0].Attributes["connectorCount"] != 2 {
		t.Errorf("entities[0].Attributes = %+v", entities[0].Attributes)
	}
}

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector emits for dns_record/tunnel/policy entities is declared in its
// catalog with a matching type.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	catalog := attributeCatalog
	emitted := map[string]map[string]string{
		"dns_record": {"type": "string", "proxied": "boolean", "ttl": "number"},
		"tunnel":     {"status": "string", "connectorCount": "number"},
		"policy":     {"enabled": "boolean", "decision": "string"},
	}
	for kind, keys := range emitted {
		specs, ok := catalog[kind]
		if !ok {
			t.Fatalf("catalog missing entity kind %q", kind)
		}
		byName := make(map[string]connector.AttributeSpec, len(specs))
		for _, s := range specs {
			byName[s.Name] = s
		}
		for key, typ := range keys {
			spec, ok := byName[key]
			if !ok {
				t.Errorf("catalog[%q] missing emitted attribute %q", kind, key)
				continue
			}
			if spec.Type != typ {
				t.Errorf("catalog[%q][%q].Type = %q, want %q", kind, key, spec.Type, typ)
			}
		}
	}
}
