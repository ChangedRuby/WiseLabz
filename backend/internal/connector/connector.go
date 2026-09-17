// Package connector defines the interface for infrastructure data connectors.
package connector

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"
)

// Connector fetches infrastructure data from a specific source.
type Connector interface {
	Name() string
	Type() string
	Category() string
	// Fetch pulls a snapshot. config may carry a "fields" hint (see
	// RequestedFields) requesting a partial fetch; a connector that doesn't
	// support selective fetch may ignore it and always return everything.
	Fetch(ctx context.Context, config map[string]any) (*ServiceSnapshot, error)
	Validate(ctx context.Context, config map[string]any) error
}

// CredentialRefresher is implemented by connectors whose credentials can be
// refreshed without user interaction (e.g. OAuth2 refresh tokens). The sync
// engine calls RefreshCredentials before Fetch when the connector's stored
// credentials have expired, and persists the returned config and expiry.
type CredentialRefresher interface {
	RefreshCredentials(ctx context.Context, config map[string]any) (newConfig map[string]any, expiresAt time.Time, err error)
}

// Restarter is implemented by connectors whose vendor API exposes a restart
// action. entityRef is the target entity's SnapshotEntity.ExternalID (a VM
// ID, container ID, service name, ...), or "" for connectors that manage a
// single implicit service. A connector that doesn't implement Restarter is
// simply not restart-capable.
type Restarter interface {
	Restart(ctx context.Context, config map[string]any, entityRef string) error
}

// Starter is implemented by connectors whose vendor API exposes a start
// action. Same entityRef convention as Restarter.
type Starter interface {
	Start(ctx context.Context, config map[string]any, entityRef string) error
}

// Stopper is implemented by connectors whose vendor API exposes a stop
// action. Same entityRef convention as Restarter.
type Stopper interface {
	Stop(ctx context.Context, config map[string]any, entityRef string) error
}

// ConfigField describes one field a connector exposes for config-push: a
// curated subset of what the connector's config schema could theoretically
// write, deliberately narrower than the full Fetch/Validate config shape.
type ConfigField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"`        // reuse the existing SchemaField Type vocabulary
	EntityScope bool   `json:"entityScope"` // true if per-entity (needs entityRef), false if connector-global
}

// ConfigPusher is implemented by connectors that expose a curated whitelist
// of writable fields for field-level partial config updates (ADR 0003).
// value is a driver value, not user input: the handler assigns fieldKey and
// value only after checking fieldKey against WritableFields.
type ConfigPusher interface {
	WritableFields() []ConfigField
	ConfigPush(ctx context.Context, config map[string]any, entityRef, fieldKey string, value any) error
}

// AuthError indicates a connector rejected credentials (expired, revoked, or
// invalid). Retrying with the same credentials will not help.
type AuthError struct{ Err error }

func (e *AuthError) Error() string { return fmt.Sprintf("auth error: %v", e.Err) }
func (e *AuthError) Unwrap() error { return e.Err }

// NewAuthError wraps err as an AuthError.
func NewAuthError(err error) *AuthError { return &AuthError{Err: err} }

// TimeoutError indicates a connector call exceeded its deadline.
type TimeoutError struct{ Err error }

func (e *TimeoutError) Error() string { return fmt.Sprintf("timeout: %v", e.Err) }
func (e *TimeoutError) Unwrap() error { return e.Err }

// NewTimeoutError wraps err as a TimeoutError.
func NewTimeoutError(err error) *TimeoutError { return &TimeoutError{Err: err} }

// MalformedResponseError indicates the upstream service returned data the
// connector could not parse.
type MalformedResponseError struct{ Err error }

func (e *MalformedResponseError) Error() string { return fmt.Sprintf("malformed response: %v", e.Err) }
func (e *MalformedResponseError) Unwrap() error { return e.Err }

// NewMalformedResponseError wraps err as a MalformedResponseError.
func NewMalformedResponseError(err error) *MalformedResponseError {
	return &MalformedResponseError{Err: err}
}

// ServiceUnavailableError indicates the upstream service is reachable but
// reported itself as unavailable (e.g. HTTP 503, maintenance mode).
type ServiceUnavailableError struct{ Err error }

func (e *ServiceUnavailableError) Error() string {
	return fmt.Sprintf("service unavailable: %v", e.Err)
}
func (e *ServiceUnavailableError) Unwrap() error { return e.Err }

// NewServiceUnavailableError wraps err as a ServiceUnavailableError.
func NewServiceUnavailableError(err error) *ServiceUnavailableError {
	return &ServiceUnavailableError{Err: err}
}

// RequestedFields extracts the optional "fields" selective-fetch hint from a
// connector config map. Returns nil if absent (meaning: fetch everything).
func RequestedFields(config map[string]any) []string {
	raw, ok := config["fields"]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

// WantsField reports whether field should be included given a fields hint.
// An empty/nil fields list means "everything" — every field is wanted.
func WantsField(fields []string, field string) bool {
	if len(fields) == 0 {
		return true
	}
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

// ServiceSnapshot represents a point-in-time view of a service's state.
type ServiceSnapshot struct {
	ServiceName  string              `json:"serviceName"`
	Type         string              `json:"type"`
	Sections     []SnapshotSection   `json:"sections"`
	Dependencies []ServiceDependency `json:"dependencies,omitempty"`
	Entities     []SnapshotEntity    `json:"entities,omitempty"`
	Metadata     map[string]string   `json:"metadata,omitempty"`
	FetchedAt    time.Time           `json:"fetchedAt"`
}

// SnapshotEntity is a structured, addressable object a connector observed
// (a VM, a container, a firewall rule, a DNS record). Cross-connector
// linking matches entities against each other by ExternalID, IP, or
// Hostname — see doc.matchEntities.
type SnapshotEntity struct {
	Kind       string `json:"kind"` // "vm", "container", "rule", "dns_record"
	Name       string `json:"name"`
	IP         string `json:"ip,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	ExternalID string `json:"externalId,omitempty"` // Proxmox VMID, container ID, etc.
	// Attributes carries stable, security-relevant, connector-specific
	// values (enabled, privileged, protocol, ...) for compliance rules to
	// evaluate (see AttributeSpec/RegisterAttributeCatalog). Values must be
	// JSON-friendly (bool, string, number, string arrays) and stable across
	// syncs — no uptime, CPU%, or other counters, since those would make
	// every snapshot look changed.
	Attributes map[string]any `json:"attributes,omitempty"`
}

// SnapshotSection is a named section of infrastructure data.
type SnapshotSection struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ServiceDependency is an operational dependency a service relies on: a
// host it runs on, a network it's reachable through, storage it consumes,
// or an upstream service it calls.
type ServiceDependency struct {
	Kind string `json:"kind"` // host | network | storage | upstream_service
	Name string `json:"name"`
	Ref  string `json:"ref,omitempty"` // optional connector/service ID if known
}

// GuardedDialer returns a *net.Dialer whose Control rejects connections to
// loopback and link-local addresses (the latter covers cloud metadata
// endpoints like 169.254.169.254). Private ranges (RFC 1918) are allowed
// since this is a self-hosted monitoring tool meant to connect to internal
// infrastructure. Enforcing the check in the dialer (rather than
// pre-resolving with net.LookupIP) closes the DNS-rebinding TOCTOU window:
// the address actually dialed is the one validated, even if the hostname
// re-resolves between check and use.
func GuardedDialer(timeout time.Duration) *net.Dialer {
	return &net.Dialer{
		Timeout: timeout,
		Control: func(_, address string, _ syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return fmt.Errorf("split address %q: %w", address, err)
			}
			ip := net.ParseIP(host)
			if ip == nil {
				return fmt.Errorf("unresolvable address %q", host)
			}
			if IsDangerousIP(ip) {
				return fmt.Errorf("connection to blocked address %s denied", ip)
			}
			return nil
		},
	}
}

// IsDangerousIP returns true for loopback and link-local unicast addresses.
func IsDangerousIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsLinkLocalUnicast()
}

// MetadataValue Metadata returns a string metadata value, or the fallback if not set.
func (s *ServiceSnapshot) MetadataValue(key, fallback string) string {
	if s.Metadata == nil {
		return fallback
	}
	if v, ok := s.Metadata[key]; ok {
		return v
	}
	return fallback
}
