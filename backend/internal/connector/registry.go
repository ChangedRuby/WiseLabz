package connector

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

// Factory creates a Connector from its config.
type Factory func(config map[string]any) (Connector, error)

// TypeSchema describes the configuration schema for a connector type.
type TypeSchema struct {
	Type     string        `json:"type"`
	Category string        `json:"category"`
	Name     string        `json:"displayName"`
	Fields   []SchemaField `json:"fields"`
	// Stub is true for connector types with no real implementation yet
	// (Fetch/Validate always fail). The UI hides/disables Test, Sync, and
	// data-viewing actions for these.
	Stub bool `json:"stub,omitempty"`
	// DegradedLatencyThresholdMs overrides the global
	// DegradedLatencyThreshold for health checks on this connector type
	// (e.g. a WiFi-backed host needs a looser SLA than an on-LAN one). Zero
	// means "use the global default".
	DegradedLatencyThresholdMs int `json:"degradedLatencyThresholdMs,omitempty"`
	// IsCredentialRefresher is computed by ListSchemas, not set at
	// registration: true if this type's Connector implementation satisfies
	// CredentialRefresher (see IsCredentialRefresherType).
	IsCredentialRefresher bool `json:"isCredentialRefresher,omitempty"`
}

// DegradedLatencyThreshold returns this type's configured health-check
// threshold, falling back to the package-wide default when unset.
func (s TypeSchema) DegradedLatencyThreshold() time.Duration {
	if s.DegradedLatencyThresholdMs <= 0 {
		return DegradedLatencyThreshold
	}
	return time.Duration(s.DegradedLatencyThresholdMs) * time.Millisecond
}

// SchemaField describes a single configuration field.
type SchemaField struct {
	Key   string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"kind"` // "text", "password", "number", "select", "toggle", "secret"
	// "secret" is a multi-line paste field for PEM certs/keys and similar
	// blobs. It validates like "text" and is encrypted at rest like
	// "password" (see store.IsSecretFieldType); unlike "password" it isn't
	// masked in the UI, since masking a multi-line block isn't useful.
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Description string `json:"description,omitempty"`
	// Pattern, MinLength, MaxLength apply to "text"/"password" fields.
	Pattern   string `json:"pattern,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`
	// Options lists the allowed values for a "select" field.
	Options []string `json:"options,omitempty"`
}

// ConfigValidationError reports a config value that fails its SchemaField's
// validation rules.
type ConfigValidationError struct {
	Field   string
	Message string
}

func (e *ConfigValidationError) Error() string {
	return fmt.Sprintf("field %q: %s", e.Field, e.Message)
}

// ValidateConfig checks config against a connector type's schema, catching
// malformed values (a non-URL in a URL field, an out-of-range string, an
// invalid enum choice) before the connector's own Validate/Fetch ever runs.
// It only validates fields that are present — connectors can still be
// created/updated with required fields filled in later; presence is
// enforced at Fetch/Validate time by the connector itself.
func ValidateConfig(schema TypeSchema, config map[string]any) error {
	for _, f := range schema.Fields {
		raw, present := config[f.Key]
		str, isStr := raw.(string)
		if !present || !isStr || str == "" {
			continue
		}
		if f.MinLength > 0 && len(str) < f.MinLength {
			return &ConfigValidationError{Field: f.Key, Message: fmt.Sprintf("must be at least %d characters", f.MinLength)}
		}
		if f.MaxLength > 0 && len(str) > f.MaxLength {
			return &ConfigValidationError{Field: f.Key, Message: fmt.Sprintf("must be at most %d characters", f.MaxLength)}
		}
		if f.Pattern != "" {
			re, err := regexp.Compile(f.Pattern)
			if err != nil {
				return fmt.Errorf("field %q: invalid schema pattern: %w", f.Key, err)
			}
			if !re.MatchString(str) {
				return &ConfigValidationError{Field: f.Key, Message: "does not match the required format"}
			}
		}
		if f.Type == "select" && len(f.Options) > 0 {
			ok := false
			for _, opt := range f.Options {
				if opt == str {
					ok = true
					break
				}
			}
			if !ok {
				return &ConfigValidationError{Field: f.Key, Message: fmt.Sprintf("must be one of %v", f.Options)}
			}
		}
	}
	return nil
}

// AttributeSpec describes one structured attribute an entity kind may carry
// (see SnapshotEntity.Attributes): its key, JSON value shape, and what it
// means. Connectors declare these next to their Register() call so the
// catalog stays in sync with what Fetch actually emits.
type AttributeSpec struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "string", "number", "boolean", "string_array"
	Description string `json:"description"`
}

var (
	registry         = make(map[string]Factory)
	typeSchema       = make(map[string]TypeSchema)
	attributeCatalog = make(map[string]map[string][]AttributeSpec) // connector type -> entity kind -> specs
	mu               sync.RWMutex
)

// RegisterAttributeCatalog declares, for a connector type, which entity
// kinds carry which Attributes keys. Exposed via GET /api/compliance/schema
// for PR3's rule engine. Call it from the same init() that calls Register.
func RegisterAttributeCatalog(typ string, catalog map[string][]AttributeSpec) {
	mu.Lock()
	defer mu.Unlock()
	attributeCatalog[typ] = catalog
}

// AttributeCatalog returns the full connectorType -> entity kind -> attribute
// spec catalog for every connector type that registered one.
func AttributeCatalog() map[string]map[string][]AttributeSpec {
	mu.RLock()
	defer mu.RUnlock()
	out := make(map[string]map[string][]AttributeSpec, len(attributeCatalog))
	for typ, kinds := range attributeCatalog {
		out[typ] = kinds
	}
	return out
}

// Register registers a connector factory and its type schema.
func Register(schema TypeSchema, factory Factory) {
	mu.Lock()
	defer mu.Unlock()
	registry[schema.Type] = factory
	typeSchema[schema.Type] = schema
}

// Get returns a new connector instance for the given type and config.
func Get(typ string, config map[string]any) (Connector, error) {
	mu.RLock()
	factory, ok := registry[typ]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown connector type: %q", typ)
	}
	return factory(config)
}

// GetTypeSchema returns the configuration schema for a connector type.
func GetTypeSchema(typ string) (*TypeSchema, error) {
	mu.RLock()
	schema, ok := typeSchema[typ]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown connector type: %q", typ)
	}
	return &schema, nil
}

// IsCredentialRefresherType reports whether typ's connector implementation
// satisfies CredentialRefresher (its credentials can be refreshed without
// user interaction, e.g. an OAuth2 refresh token) — such connectors are
// skipped by the credential_rotation quality check since there's nothing
// for a human to rotate. Constructing a connector via its factory is
// expected to be cheap and side-effect-free; only Fetch/Validate/
// RefreshCredentials touch the network, matching how the sync engine
// already probes this interface at runtime.
func IsCredentialRefresherType(typ string) bool {
	inst, err := Get(typ, map[string]any{})
	if err != nil {
		return false
	}
	_, ok := inst.(CredentialRefresher)
	return ok
}

// ListSchemas returns all registered connector type schemas, with
// IsCredentialRefresher computed for each.
func ListSchemas() []TypeSchema {
	mu.RLock()
	defer mu.RUnlock()
	var out []TypeSchema
	for _, s := range typeSchema {
		if factory, ok := registry[s.Type]; ok {
			if inst, err := factory(map[string]any{}); err == nil {
				_, s.IsCredentialRefresher = inst.(CredentialRefresher)
			}
		}
		out = append(out, s)
	}
	return out
}
