// Package compliance evaluates user-defined rules against connector snapshots.
package compliance

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	maxConditions  = 20
	maxRegexLength = 256
)

// Rule is a compliance rule stored by the application.
type Rule struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	ConnectorType   string      `json:"connectorType"`
	EntityKind      string      `json:"entityKind"`
	Conditions      []Condition `json:"conditions"`
	Severity        string      `json:"severity"`
	Title           string      `json:"title"`
	RemediationLink string      `json:"remediationLink"`
	Enabled         bool        `json:"enabled"`
}

// Condition compares one entity attribute with a value.
type Condition struct {
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     any    `json:"value"`
}

// Catalog describes the attributes emitted by each connector and entity kind.
type Catalog map[string]map[string][]AttributeSpec

// AttributeSpec describes one entity attribute.
type AttributeSpec struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Snapshot is the rule-engine view of a connector snapshot.
type Snapshot struct {
	Entities []Entity `json:"entities"`
}

// Entity is a structured item in a snapshot.
type Entity struct {
	Kind       string         `json:"kind"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"attributes"`
}

// FieldError names one rule field that failed validation, and why.
type FieldError struct {
	Field string
	Msg   string
}

// ValidationError reports the rule fields that failed validation. Message is
// the human-readable sentence on its own; Fields names the offending rule
// fields, so an API layer can return them per-field instead of re-deriving
// them from the message text.
type ValidationError struct {
	Message string
	Fields  []FieldError
	// Err is the underlying cause when a check has one (a regexp compile
	// failure), so errors.Is/As still reach it through the wrapper.
	Err error
}

func (e *ValidationError) Error() string { return e.Message }

func (e *ValidationError) Unwrap() error { return e.Err }

// invalidField builds a single-field ValidationError. format/args produce
// Message, so the sentence each check already returned is preserved exactly.
func invalidField(field, msg, format string, args ...any) *ValidationError {
	return &ValidationError{
		Message: fmt.Sprintf(format, args...),
		Fields:  []FieldError{{Field: field, Msg: msg}},
	}
}

// Validate checks that a rule can be evaluated with catalog. Empty conditions
// are rejected: matching every entity is not a useful compliance rule.
//
// Every rejection is a *ValidationError naming the rule field at fault; the
// per-condition ones use an indexed path (conditions[2].op) so a caller can
// point at the exact row of a multi-condition rule.
func Validate(rule Rule, catalog Catalog) error {
	kinds, ok := catalog[rule.ConnectorType]
	if !ok {
		return invalidField("connectorType", "is not a known connector type",
			"unknown connector type %q", rule.ConnectorType)
	}
	attributes, ok := kinds[rule.EntityKind]
	if !ok {
		return invalidField("entityKind", "is not a known entity kind for this connector type",
			"unknown entity kind %q for connector type %q", rule.EntityKind, rule.ConnectorType)
	}
	if len(rule.Conditions) == 0 {
		return invalidField("conditions", "must contain at least one condition",
			"at least one condition is required")
	}
	if len(rule.Conditions) > maxConditions {
		return invalidField("conditions", fmt.Sprintf("must contain at most %d conditions", maxConditions),
			"at most %d conditions are allowed", maxConditions)
	}

	for i, condition := range rule.Conditions {
		spec, ok := findAttribute(attributes, condition.Attribute)
		if !ok {
			return invalidField(fmt.Sprintf("conditions[%d].attribute", i), "is not a known attribute for this entity kind",
				"unknown attribute %q", condition.Attribute)
		}
		if !validOp(spec.Type, condition.Op) {
			return invalidField(fmt.Sprintf("conditions[%d].op", i), fmt.Sprintf("is not a valid operator for %s attributes", spec.Type),
				"operator %q is not valid for %s attribute %q", condition.Op, spec.Type, condition.Attribute)
		}
		if condition.Op != "regex" {
			continue
		}
		pattern, ok := condition.Value.(string)
		if !ok {
			return invalidField(fmt.Sprintf("conditions[%d].value", i), "must be a string when op is regex",
				"regex value for %q must be a string", condition.Attribute)
		}
		if len(pattern) > maxRegexLength {
			return invalidField(fmt.Sprintf("conditions[%d].value", i), fmt.Sprintf("must be at most %d characters", maxRegexLength),
				"regex for %q exceeds %d characters", condition.Attribute, maxRegexLength)
		}
		if _, err := regexp.Compile(pattern); err != nil {
			invalid := invalidField(fmt.Sprintf("conditions[%d].value", i), "is not a valid regular expression",
				"invalid regex for %q: %v", condition.Attribute, err)
			invalid.Err = err
			return invalid
		}
	}
	return nil
}

// Evaluate returns matching entities in snapshot order. A missing attribute is
// false for every operator except neq and exists. exists uses a boolean value:
// true requires presence and false requires absence.
func Evaluate(rule Rule, snapshot Snapshot) []Entity {
	matches := make([]Entity, 0)
	for _, entity := range snapshot.Entities {
		if entity.Kind != rule.EntityKind || !matchesAll(entity, rule.Conditions) {
			continue
		}
		matches = append(matches, entity)
	}
	return matches
}

func findAttribute(attributes []AttributeSpec, name string) (AttributeSpec, bool) {
	for _, attribute := range attributes {
		if attribute.Name == name {
			return attribute, true
		}
	}
	return AttributeSpec{}, false
}

func validOp(attributeType, op string) bool {
	switch attributeType {
	case "string":
		return op == "eq" || op == "neq" || op == "contains" || op == "regex" || op == "exists"
	case "string_array":
		return op == "eq" || op == "neq" || op == "contains" || op == "exists"
	case "number":
		return op == "eq" || op == "neq" || op == "gt" || op == "lt" || op == "exists"
	case "boolean":
		return op == "eq" || op == "neq" || op == "exists"
	default:
		return false
	}
}

func matchesAll(entity Entity, conditions []Condition) bool {
	for _, condition := range conditions {
		value, present := entity.Attributes[condition.Attribute]
		if !matches(value, present, condition) {
			return false
		}
	}
	return true
}

func matches(value any, present bool, condition Condition) bool {
	if condition.Op == "exists" {
		want, ok := condition.Value.(bool)
		return ok && present == want
	}
	if !present {
		return condition.Op == "neq"
	}
	if value == nil {
		return false
	}

	switch condition.Op {
	case "eq":
		return equal(value, condition.Value)
	case "neq":
		return !equal(value, condition.Value)
	case "contains":
		return contains(value, condition.Value)
	case "regex":
		pattern, ok := condition.Value.(string)
		if !ok {
			return false
		}
		s, ok := value.(string)
		if !ok {
			return false
		}
		re, err := regexp.Compile(pattern)
		return err == nil && re.MatchString(s)
	case "gt", "lt":
		left, leftOK := number(value)
		right, rightOK := number(condition.Value)
		if !leftOK || !rightOK {
			return false
		}
		return condition.Op == "gt" && left > right || condition.Op == "lt" && left < right
	default:
		return false
	}
}

func equal(left, right any) bool {
	leftNumber, leftOK := number(left)
	rightNumber, rightOK := number(right)
	if leftOK || rightOK {
		return leftOK && rightOK && leftNumber == rightNumber
	}
	leftStrings := stringSlice(left)
	rightStrings := stringSlice(right)
	if leftStrings != nil || rightStrings != nil {
		return leftStrings != nil && rightStrings != nil && reflect.DeepEqual(leftStrings, rightStrings)
	}
	return reflect.DeepEqual(left, right)
}

func contains(value, want any) bool {
	needle, ok := want.(string)
	if !ok {
		return false
	}
	switch actual := value.(type) {
	case string:
		return strings.Contains(actual, needle)
	case []string:
		for _, item := range actual {
			if item == needle {
				return true
			}
		}
	case []any:
		for _, item := range actual {
			if item == needle {
				return true
			}
		}
	}
	return false
}

func stringSlice(value any) []string {
	switch items := value.(type) {
	case []string:
		return items
	case []any:
		result := make([]string, 0, len(items))
		for _, item := range items {
			text, ok := item.(string)
			if !ok {
				return nil
			}
			result = append(result, text)
		}
		return result
	default:
		return nil
	}
}

func number(value any) (float64, bool) {
	switch n := value.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		parsed, err := n.Float64()
		return parsed, err == nil
	default:
		return 0, false
	}
}
