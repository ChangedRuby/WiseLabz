package adguardhome

import (
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector actually emits is declared in its catalog with a matching type,
// and that the catalog declares nothing the connector never emits.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	_, _, filters, _ := buildFiltering([]byte(filteringJSON))
	_, rewrites := buildRewriteTable([]byte(rewritesJSON))
	_, clients := buildClientTable([]byte(clientsJSON))
	_, leases, _ := buildDHCP([]byte(dhcpJSON))

	emitted := map[string]map[string]string{}
	for _, group := range [][]connector.SnapshotEntity{filters, rewrites, clients, leases} {
		for _, e := range group {
			if emitted[e.Kind] == nil {
				emitted[e.Kind] = map[string]string{}
			}
			for key, value := range e.Attributes {
				emitted[e.Kind][key] = jsonType(value)
			}
		}
	}

	for kind, keys := range emitted {
		specs, ok := attributeCatalog[kind]
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
			if spec.Description == "" {
				t.Errorf("catalog[%q][%q] has no description", kind, key)
			}
		}
		for name := range byName {
			if _, ok := keys[name]; !ok {
				t.Errorf("catalog[%q] declares %q, which the fixtures never emit", kind, name)
			}
		}
	}

	for kind := range attributeCatalog {
		if _, ok := emitted[kind]; !ok {
			t.Errorf("catalog declares entity kind %q, which the connector never emits", kind)
		}
	}
}

// jsonType maps a Go attribute value to the AttributeSpec type vocabulary.
func jsonType(v any) string {
	switch v.(type) {
	case bool:
		return "boolean"
	case string:
		return "string"
	case int, int64, float64:
		return "number"
	case []string:
		return "string_array"
	default:
		return "unknown"
	}
}
