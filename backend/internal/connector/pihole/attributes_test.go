package pihole

import (
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildHostsTableAttributes(t *testing.T) {
	data := []byte(`{"config":{"dns":{"hosts":["10.0.0.5 nas.internal.example.com","2001:db8::1 router.internal.example.com"]}}}`)
	_, entities := buildHostsTable(data)
	if len(entities) != 2 {
		t.Fatalf("entities = %+v, want 2", entities)
	}

	want := []map[string]any{
		{"source": "local_dns", "is_ipv6": false},
		{"source": "local_dns", "is_ipv6": true},
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
		"dns_record": {"source": "string", "is_ipv6": "boolean"},
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
