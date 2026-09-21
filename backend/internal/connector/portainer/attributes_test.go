package portainer

import (
	"encoding/json"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector actually emits is declared in its catalog with a matching type,
// and that the catalog declares nothing the connector never emits.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	var envs []environment
	if err := json.Unmarshal([]byte(endpointsJSON), &envs); err != nil {
		t.Fatalf("decode endpoints fixture: %v", err)
	}
	_, environments := buildEnvironmentTable(envs)
	_, stacks := buildStackTable([]byte(stacksJSON), environmentNames(envs))
	_, containers, err := containerRows("local", []byte(containersJSON))
	if err != nil {
		t.Fatalf("containerRows: %v", err)
	}
	_, volumes, err := volumeRows("local", []byte(volumesJSON))
	if err != nil {
		t.Fatalf("volumeRows: %v", err)
	}
	_, networks, err := networkRows("local", []byte(networksJSON))
	if err != nil {
		t.Fatalf("networkRows: %v", err)
	}

	emitted := map[string]map[string]string{}
	for _, group := range [][]connector.SnapshotEntity{environments, stacks, containers, volumes, networks} {
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
