package netbird

import (
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildPeerTableAttributes(t *testing.T) {
	data := []byte(`[
		{"id":"p1","name":"laptop","ip":"100.64.0.1","os":"linux","approval_required":false,"connected":true,"groups":[{"name":"engineering"}]},
		{"id":"p2","name":"phone","ip":"100.64.0.2","os":"android","approval_required":true,"connected":false}
	]`)
	_, entities := buildPeerTable(data)
	if len(entities) != 2 {
		t.Fatalf("entities = %+v, want 2", entities)
	}
	want := []map[string]any{
		{"approved": true, "connected": true, "os": "linux", "groups": []string{"engineering"}},
		{"approved": false, "connected": false, "os": "android"},
	}
	for i, w := range want {
		if !reflect.DeepEqual(entities[i].Attributes, w) {
			t.Errorf("entities[%d].Attributes = %+v, want %+v", i, entities[i].Attributes, w)
		}
	}
}

func TestBuildPolicyTableAttributes(t *testing.T) {
	data := []byte(`[
		{"id":"pol1","name":"allow-eng","enabled":true,"rules":[{"protocol":"tcp","sources":[{"name":"engineering"}],"destinations":[{"name":"prod"}]}]},
		{"id":"pol2","name":"deny-all","enabled":false}
	]`)
	_, entities := buildPolicyTable(data)
	if len(entities) != 2 {
		t.Fatalf("entities = %+v, want 2", entities)
	}

	allowAttrs := entities[0].Attributes
	if allowAttrs["enabled"] != true || allowAttrs["protocol"] != "tcp" {
		t.Errorf("allow-eng attributes = %+v", allowAttrs)
	}
	if !reflect.DeepEqual(allowAttrs["sourceGroups"], []string{"engineering"}) {
		t.Errorf("allow-eng sourceGroups = %+v", allowAttrs["sourceGroups"])
	}
	if !reflect.DeepEqual(allowAttrs["destinationGroups"], []string{"prod"}) {
		t.Errorf("allow-eng destinationGroups = %+v", allowAttrs["destinationGroups"])
	}

	denyAttrs := entities[1].Attributes
	if denyAttrs["enabled"] != false {
		t.Errorf("deny-all attributes = %+v", denyAttrs)
	}
	if _, ok := denyAttrs["protocol"]; ok {
		t.Errorf("deny-all attributes should omit protocol when no rules: %+v", denyAttrs)
	}
}

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector emits for peer/policy entities is declared in its catalog with a
// matching type.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	catalog := attributeCatalog
	emitted := map[string]map[string]string{
		"peer":   {"approved": "boolean", "connected": "boolean", "os": "string", "groups": "string_array"},
		"policy": {"enabled": "boolean", "protocol": "string", "sourceGroups": "string_array", "destinationGroups": "string_array"},
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
