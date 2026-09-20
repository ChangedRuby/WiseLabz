package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Schema describes config.yaml without loading configuration or including values.
// x-env names scalar environment overrides; collection settings are file-only.
func Schema() (map[string]any, error) {
	schema, err := schemaFor(reflect.TypeOf(Config{}), "", true)
	if err != nil {
		return nil, err
	}
	schema["$schema"] = "https://json-schema.org/draft/2020-12/schema"
	return schema, nil
}

func schemaFor(t reflect.Type, path string, env bool) (map[string]any, error) {
	schema := make(map[string]any)
	switch t.Kind() {
	case reflect.Struct:
		schema["type"] = "object"
		properties := make(map[string]any)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			key := field.Tag.Get("mapstructure")
			if key == "" || key == "-" {
				return nil, fmt.Errorf("config schema: missing key for %s.%s", t.Name(), field.Name)
			}
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			child, err := schemaFor(field.Type, childPath, env)
			if err != nil {
				return nil, err
			}
			properties[key] = child
		}
		schema["properties"] = properties
	case reflect.Slice, reflect.Map:
		child, err := schemaFor(t.Elem(), path, false)
		if err != nil {
			return nil, err
		}
		if t.Kind() == reflect.Slice {
			schema["type"], schema["items"] = "array", child
		} else {
			schema["type"], schema["additionalProperties"] = "object", child
		}
		schema["x-config-file-only"] = true
	case reflect.String, reflect.Int, reflect.Bool:
		schema["type"] = map[reflect.Kind]string{reflect.String: "string", reflect.Int: "integer", reflect.Bool: "boolean"}[t.Kind()]
		if env {
			schema["x-env"] = "WISELABZ_" + strings.ToUpper(strings.ReplaceAll(path, ".", "_"))
		}
	default:
		return nil, fmt.Errorf("config schema: unsupported type %s at %s", t, path)
	}
	return schema, nil
}
