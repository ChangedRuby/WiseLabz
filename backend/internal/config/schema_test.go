package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestSchemaMatchesConfig(t *testing.T) {
	schema, err := Schema()
	if err != nil {
		t.Fatal(err)
	}
	// Walk every field, including file-only OIDC collections, so added fields,
	// missing keys and unsupported types cannot silently disappear from the schema.
	var check func(reflect.Type, map[string]any, string, bool)
	check = func(typ reflect.Type, node map[string]any, path string, env bool) {
		t.Helper()
		switch typ.Kind() {
		case reflect.Struct:
			props, ok := node["properties"].(map[string]any)
			if node["type"] != "object" || !ok || len(props) != typ.NumField() {
				t.Fatalf("%s: properties do not match %s: %v", path, typ, node)
			}
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				key := field.Tag.Get("mapstructure")
				child, ok := props[key].(map[string]any)
				if !ok {
					t.Fatalf("%s: missing property %s", path, key)
				}
				check(field.Type, child, strings.TrimPrefix(path+"_"+strings.ToUpper(key), "_"), env)
			}
		case reflect.Slice, reflect.Map:
			wantType, childKey := "array", "items"
			if typ.Kind() == reflect.Map {
				wantType, childKey = "object", "additionalProperties"
			}
			child, ok := node[childKey].(map[string]any)
			if node["type"] != wantType || !ok || node["x-config-file-only"] != true {
				t.Fatalf("%s: invalid collection schema: %v", path, node)
			}
			check(typ.Elem(), child, path, false)
		default:
			want := map[reflect.Kind]string{reflect.String: "string", reflect.Int: "integer", reflect.Bool: "boolean"}[typ.Kind()]
			if want == "" || node["type"] != want {
				t.Fatalf("%s: schema type %v does not match %s", path, node["type"], typ)
			}
			if env && node["x-env"] != "WISELABZ_"+path {
				t.Errorf("%s: env = %v", path, node["x-env"])
			}
			if !env && node["x-env"] != nil {
				t.Errorf("%s: file-only field advertises an env override", path)
			}
		}
		for _, key := range []string{"default", "const", "examples"} {
			if _, ok := node[key]; ok {
				t.Errorf("%s: schema unexpectedly includes values: %s", path, key)
			}
		}
	}
	check(reflect.TypeOf(Config{}), schema, "", true)
}
