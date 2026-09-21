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

// TestSchemaExposesAPIVersion pins the version-awareness contract: the
// connector schema offers auto/v6/v5, defaults to auto so existing configs
// keep working, and the factory falls back to auto for anything else.
func TestSchemaExposesAPIVersion(t *testing.T) {
	schema, err := connector.GetTypeSchema(typeName)
	if err != nil {
		t.Fatalf("GetTypeSchema: %v", err)
	}
	var field *connector.SchemaField
	for i := range schema.Fields {
		if schema.Fields[i].Key == "api_version" {
			field = &schema.Fields[i]
		}
	}
	if field == nil {
		t.Fatal("schema has no api_version field")
	}
	if field.Default != versionAuto {
		t.Errorf("api_version default = %q, want %q", field.Default, versionAuto)
	}
	if !reflect.DeepEqual(field.Options, []string{versionAuto, version6, version5}) {
		t.Errorf("api_version options = %v", field.Options)
	}
	if err := connector.ValidateConfig(*schema, map[string]any{"api_version": "v7"}); err == nil {
		t.Error("ValidateConfig accepted an unknown api_version")
	}

	for in, want := range map[string]string{"": versionAuto, "v5": version5, "v6": version6, "bogus": versionAuto} {
		impl, err := connector.Get(typeName, map[string]any{"url": "https://pihole.example.com", "password": "x", "api_version": in})
		if err != nil {
			t.Fatalf("Get(%q): %v", in, err)
		}
		if got := impl.(*Connector).apiVersion; got != want {
			t.Errorf("api_version %q resolved to %q, want %q", in, got, want)
		}
	}
}

// TestAttributeCatalogCoversNewEntityKinds keeps the catalog in step with the
// blocklist/group/client/domain entities the extended Fetch now emits.
func TestAttributeCatalogCoversNewEntityKinds(t *testing.T) {
	emitted := map[string]map[string]string{
		"blocklist":   {"enabled": "boolean", "list_type": "string", "groups": "string_array"},
		"dns_group":   {"enabled": "boolean", "comment": "string"},
		"dns_client":  {"groups": "string_array", "comment": "string"},
		"domain_rule": {"enabled": "boolean", "rule_type": "string", "match_kind": "string", "groups": "string_array"},
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
		}
	}
}
