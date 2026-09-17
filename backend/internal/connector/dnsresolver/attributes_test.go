package dnsresolver

import (
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildHostOverrideTableAttributes(t *testing.T) {
	data := []byte(`{"data":[
		{"host":"nas","domain":"internal.example.com","ip":"10.0.0.5","descr":"NAS Storage"},
		{"host":"","domain":"example.com","ip":"2001:db8::1","descr":"Root domain"},
		{"host":"web","domain":"example.com","ip":"192.168.1.10","descr":""}
	]}`)
	_, entities := buildHostOverrideTable(data)
	if len(entities) != 3 {
		t.Fatalf("entities = %+v, want 3", entities)
	}

	want := []map[string]any{
		{"description": "NAS Storage", "is_ipv6": false},
		{"description": "Root domain", "is_ipv6": true},
		{"is_ipv6": false},
	}
	for i, w := range want {
		if !reflect.DeepEqual(entities[i].Attributes, w) {
			t.Errorf("entities[%d].Attributes = %+v, want %+v", i, entities[i].Attributes, w)
		}
	}
}

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector emits for dns_record entities is declared in its catalog
// with a matching type, so PR3's rule engine and the schema endpoint never
// drift from what Fetch actually produces.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	catalog := attributeCatalog
	emitted := map[string]map[string]string{
		"dns_record": {
			"description": "string",
			"is_ipv6":     "boolean",
		},
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
