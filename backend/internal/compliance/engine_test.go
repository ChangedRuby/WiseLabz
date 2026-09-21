package compliance

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func testCatalog() Catalog {
	return Catalog{"docker": {"container": {
		{Name: "name", Type: "string"},
		{Name: "tags", Type: "string_array"},
		{Name: "ports", Type: "number"},
		{Name: "privileged", Type: "boolean"},
	}}}
}

func TestEvaluateOperators(t *testing.T) {
	entity := Entity{Kind: "container", Name: "web", Attributes: map[string]any{
		"name": "web-server", "tags": []any{"prod", "public"}, "ports": json.Number("8080"), "privileged": true,
	}}
	tests := []struct {
		name      string
		condition Condition
		want      bool
	}{
		{"string eq", Condition{"name", "eq", "web-server"}, true},
		{"string neq", Condition{"name", "neq", "db"}, true},
		{"string contains", Condition{"name", "contains", "server"}, true},
		{"string regex", Condition{"name", "regex", "^web-"}, true},
		{"string exists", Condition{"name", "exists", true}, true},
		{"array eq decoded", Condition{"tags", "eq", []string{"prod", "public"}}, true},
		{"array neq", Condition{"tags", "neq", []string{"prod"}}, true},
		{"array contains", Condition{"tags", "contains", "prod"}, true},
		{"array exists", Condition{"tags", "exists", true}, true},
		{"number eq json", Condition{"ports", "eq", float64(8080)}, true},
		{"number neq", Condition{"ports", "neq", 443}, true},
		{"number gt", Condition{"ports", "gt", 1024}, true},
		{"number lt", Condition{"ports", "lt", 9000}, true},
		{"number exists", Condition{"ports", "exists", true}, true},
		{"boolean eq", Condition{"privileged", "eq", true}, true},
		{"boolean neq", Condition{"privileged", "neq", false}, true},
		{"boolean exists", Condition{"privileged", "exists", true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := Rule{EntityKind: "container", Conditions: []Condition{tt.condition}}
			got := Evaluate(rule, Snapshot{Entities: []Entity{entity}})
			if (len(got) == 1) != tt.want {
				t.Fatalf("matches = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateEdgeCases(t *testing.T) {
	entity := Entity{Kind: "container", Name: "web", Attributes: map[string]any{
		"null": nil, "name": "web", "tags": []string{"prod"}, "ports": "8080", "privileged": "true",
	}}
	tests := []struct {
		name      string
		condition Condition
		want      bool
	}{
		{"missing eq", Condition{"missing", "eq", "x"}, false},
		{"missing neq", Condition{"missing", "neq", "x"}, true},
		{"missing exists true", Condition{"missing", "exists", true}, false},
		{"missing exists false", Condition{"missing", "exists", false}, true},
		{"null eq", Condition{"null", "eq", nil}, false},
		{"null neq", Condition{"null", "neq", "x"}, false},
		{"null exists", Condition{"null", "exists", true}, true},
		{"wrong contains value", Condition{"tags", "contains", 2}, false},
		{"string regex mismatch", Condition{"name", "regex", "^db"}, false},
		{"invalid regex", Condition{"name", "regex", "["}, false},
		{"number type mismatch", Condition{"ports", "gt", 80}, false},
		{"boolean type mismatch", Condition{"privileged", "eq", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Evaluate(Rule{EntityKind: "container", Conditions: []Condition{tt.condition}}, Snapshot{Entities: []Entity{entity}})
			if (len(got) == 1) != tt.want {
				t.Fatalf("matches = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluateAndKindAndOrder(t *testing.T) {
	rule := Rule{EntityKind: "container", Conditions: []Condition{
		{Attribute: "name", Op: "contains", Value: "web"},
		{Attribute: "privileged", Op: "eq", Value: true},
	}}
	snapshot := Snapshot{Entities: []Entity{
		{Kind: "vm", Name: "wrong-kind", Attributes: map[string]any{"name": "web", "privileged": true}},
		{Kind: "container", Name: "not-privileged", Attributes: map[string]any{"name": "web", "privileged": false}},
		{Kind: "container", Name: "second", Attributes: map[string]any{"name": "web2", "privileged": true}},
		{Kind: "container", Name: "third", Attributes: map[string]any{"name": "web3", "privileged": true}},
	}}
	got := Evaluate(rule, snapshot)
	if len(got) != 2 || got[0].Name != "second" || got[1].Name != "third" {
		t.Fatalf("matches = %#v", got)
	}
}

func TestEvaluateLargeSnapshot(t *testing.T) {
	entities := make([]Entity, 0, 1001)
	for range 1000 {
		entities = append(entities, Entity{Kind: "container", Name: "no", Attributes: map[string]any{"privileged": false}})
	}
	entities = append(entities, Entity{Kind: "container", Name: "yes", Attributes: map[string]any{"privileged": true}})
	got := Evaluate(Rule{EntityKind: "container", Conditions: []Condition{{Attribute: "privileged", Op: "eq", Value: true}}}, Snapshot{Entities: entities})
	if len(got) != 1 || got[0].Name != "yes" {
		t.Fatalf("matches = %#v", got)
	}
}

func TestValidate(t *testing.T) {
	valid := Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "name", Op: "eq", Value: "web"}}}
	tests := []struct {
		name string
		rule Rule
		want string
	}{
		{"valid", valid, ""},
		{"unknown connector", Rule{ConnectorType: "nope", EntityKind: "container", Conditions: valid.Conditions}, "unknown connector"},
		{"unknown kind", Rule{ConnectorType: "docker", EntityKind: "vm", Conditions: valid.Conditions}, "unknown entity kind"},
		{"unknown attribute", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "nope", Op: "eq", Value: "x"}}}, "unknown attribute"},
		{"invalid string op", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "name", Op: "gt", Value: "x"}}}, "not valid"},
		{"invalid array op", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "tags", Op: "regex", Value: "x"}}}, "not valid"},
		{"invalid number op", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "ports", Op: "contains", Value: "x"}}}, "not valid"},
		{"invalid boolean op", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "privileged", Op: "contains", Value: "x"}}}, "not valid"},
		{"empty", Rule{ConnectorType: "docker", EntityKind: "container"}, "at least"},
		{"invalid regex", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "name", Op: "regex", Value: "["}}}, "invalid regex"},
		{"non-string regex", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "name", Op: "regex", Value: 1}}}, "must be a string"},
		{"long regex", Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "name", Op: "regex", Value: strings.Repeat("a", maxRegexLength+1)}}}, "exceeds"},
	}
	tooMany := valid
	tooMany.Conditions = make([]Condition, maxConditions+1)
	for i := range tooMany.Conditions {
		tooMany.Conditions[i] = valid.Conditions[0]
	}
	tests = append(tests, struct {
		name string
		rule Rule
		want string
	}{"too many", tooMany, "at most"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.rule, testCatalog())
			if tt.want == "" && err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("Validate error = %v, want %q", err, tt.want)
			}
		})
	}
}

// TestValidateFieldAttribution pins the rule field each rejection names, so
// the API layer can surface it as a {field, msg} detail without re-deriving
// it from the message text.
func TestValidateFieldAttribution(t *testing.T) {
	ok := Condition{Attribute: "name", Op: "eq", Value: "web"}
	rule := func(conditions ...Condition) Rule {
		return Rule{ConnectorType: "docker", EntityKind: "container", Conditions: conditions}
	}
	tests := []struct {
		name  string
		rule  Rule
		field string
	}{
		{"unknown connector", Rule{ConnectorType: "nope", EntityKind: "container", Conditions: []Condition{ok}}, "connectorType"},
		{"unknown kind", Rule{ConnectorType: "docker", EntityKind: "vm", Conditions: []Condition{ok}}, "entityKind"},
		{"empty conditions", rule(), "conditions"},
		{"unknown attribute", rule(ok, Condition{Attribute: "nope", Op: "eq", Value: "x"}), "conditions[1].attribute"},
		{"invalid op", rule(ok, Condition{Attribute: "name", Op: "gt", Value: "x"}), "conditions[1].op"},
		{"invalid regex", rule(Condition{Attribute: "name", Op: "regex", Value: "["}), "conditions[0].value"},
		{"non-string regex", rule(ok, ok, Condition{Attribute: "name", Op: "regex", Value: 1}), "conditions[2].value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.rule, testCatalog())

			var invalid *ValidationError
			if !errors.As(err, &invalid) {
				t.Fatalf("Validate error = %v, want *ValidationError", err)
			}
			if len(invalid.Fields) != 1 {
				t.Fatalf("Fields = %+v, want exactly one", invalid.Fields)
			}
			if invalid.Fields[0].Field != tt.field {
				t.Errorf("Field = %q, want %q", invalid.Fields[0].Field, tt.field)
			}
			if invalid.Fields[0].Msg == "" {
				t.Error("Msg is empty")
			}
			// Error() must still read as the sentence callers already log.
			if invalid.Error() != invalid.Message {
				t.Errorf("Error() = %q, want %q", invalid.Error(), invalid.Message)
			}
		})
	}
}

// TestValidateRegexErrorUnwraps keeps the regexp compile failure reachable
// through the ValidationError wrapper.
func TestValidateRegexErrorUnwraps(t *testing.T) {
	err := Validate(
		Rule{ConnectorType: "docker", EntityKind: "container", Conditions: []Condition{{Attribute: "name", Op: "regex", Value: "["}}},
		testCatalog(),
	)
	var invalid *ValidationError
	if !errors.As(err, &invalid) {
		t.Fatalf("Validate error = %v, want *ValidationError", err)
	}
	if errors.Unwrap(invalid) == nil {
		t.Error("Unwrap() = nil, want the regexp compile error")
	}
}
