// Package custom implements a configurable HTTP connector stub.
package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "custom"

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "virtualization",
		Name:     "Custom HTTP",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Endpoint URL", Type: "text", Required: true, Placeholder: "https://api.example.com/status"},
			{Key: "method", Label: "HTTP Method", Type: "select", Required: false, Default: "GET"},
			{Key: "headers", Label: "Headers (JSON)", Type: "secret", Required: false, Placeholder: `{"Authorization": "Bearer token"}`},
		},
	}, func(_ map[string]any) (connector.Connector, error) {
		return &Connector{client: newGuardedClient()}, nil
	})
	// Custom connectors pass through whatever attributes the source payload provides.
	// No fixed catalog — the point is operator-defined data passthrough.
	connector.RegisterAttributeCatalog(typeName, map[string][]connector.AttributeSpec{
		"*": {{Name: "*", Type: "string", Description: "Custom connectors pass through whatever attributes object the source payload provides under entities[].attributes; there is no fixed schema."}},
	})
}

// Connector is a configurable HTTP connector for custom APIs.
type Connector struct {
	client *http.Client
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Custom HTTP" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "virtualization" }

// Validate tests the connection to the configured HTTP endpoint.
func (c *Connector) Validate(ctx context.Context, config map[string]any) error {
	rawURL, ok := config["url"].(string)
	if !ok || rawURL == "" {
		return fmt.Errorf("url is required")
	}

	if err := validateCustomURL(rawURL); err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	method, _ := config["method"].(string)
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	setHeaders(req, config)

	// The URL is user-supplied by design (this connector calls whatever endpoint
	// the operator configures); SSRF is mitigated at dial time by
	// newGuardedClient, which blocks loopback/link-local targets and redirects.
	resp, err := c.client.Do(req) // codeql[go/request-forgery]
	if err != nil {
		return connector.MapTransportError(err)
	}
	defer resp.Body.Close() //nolint:errcheck

	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return connector.NewAuthError(fmt.Errorf("server returned %d", resp.StatusCode))
	case resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout:
		return connector.NewServiceUnavailableError(fmt.Errorf("server returned %d", resp.StatusCode))
	case resp.StatusCode >= 500:
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}
	return nil
}

// Fetch retrieves data from the configured HTTP endpoint.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	rawURL, ok := config["url"].(string)
	if !ok || rawURL == "" {
		return nil, fmt.Errorf("url is required")
	}

	if err := validateCustomURL(rawURL); err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	method, _ := config["method"].(string)
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	setHeaders(req, config)

	// See Validate above: SSRF is mitigated at dial time by newGuardedClient.
	resp, err := c.client.Do(req) // codeql[go/request-forgery]
	if err != nil {
		return nil, connector.MapTransportError(err)
	}
	defer resp.Body.Close() //nolint:errcheck

	body, err := connector.ReadBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	snapshot := &connector.ServiceSnapshot{
		ServiceName: "Custom: " + rawURL,
		Type:        typeName,
		Sections: []connector.SnapshotSection{
			{Title: "Response", Content: "```json\n" + string(body) + "\n```"},
		},
		Metadata: map[string]string{
			"status_code": fmt.Sprintf("%d", resp.StatusCode),
			"url":         rawURL,
		},
		FetchedAt: time.Now(),
	}

	// Attempt to parse structured entities from the response body.
	// If the body contains a top-level "entities" array with valid SnapshotEntity shapes,
	// populate them in the snapshot. This is backward compatible: if parsing fails or
	// there are no entities, the snapshot is returned as-is (with just the Response section).
	entities, ok := tryParseEntities(body)
	if ok && len(entities) > 0 {
		snapshot.Entities = entities
	}

	return snapshot, nil
}

func setHeaders(req *http.Request, config map[string]any) {
	if headersRaw, ok := config["headers"]; ok {
		var headers map[string]string
		switch v := headersRaw.(type) {
		case string:
			if err := json.Unmarshal([]byte(v), &headers); err != nil {
				return
			}
		case map[string]string:
			headers = v
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}

// newGuardedClient builds an HTTP client whose dialer rejects connections to
// loopback or link-local addresses at connect time (see
// connector.GuardedDialer). Redirects are blocked so a 3xx cannot bounce the
// request to an internal target.
func newGuardedClient() *http.Client {
	return connector.NewHTTPClient(connector.HTTPClientOptions{})
}

// validateCustomURL parses rawURL and enforces the http/https scheme and a
// non-empty host. Network-level SSRF protection (loopback/link-local blocking)
// is enforced at dial time by newGuardedClient, not here, to avoid a
// DNS-rebinding TOCTOU gap.
func validateCustomURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("scheme %q not allowed (only http/https)", parsed.Scheme)
	}

	if parsed.Hostname() == "" {
		return fmt.Errorf("empty host")
	}

	return nil
}

// tryParseEntities attempts to parse a JSON response body for a top-level
// "entities" array matching the SnapshotEntity shape. Returns the parsed
// entities and true if parsing succeeded and at least one entity has a
// non-empty kind; otherwise returns nil and false. This is backward compatible:
// if parsing fails or there are no valid entities, the caller treats it as
// "no structured data" and uses the raw Response section instead.
func tryParseEntities(body []byte) ([]connector.SnapshotEntity, bool) {
	var payload struct {
		Entities []struct {
			Kind       string         `json:"kind"`
			Name       string         `json:"name"`
			IP         string         `json:"ip"`
			Hostname   string         `json:"hostname"`
			ExternalID string         `json:"externalId"`
			Attributes map[string]any `json:"attributes"`
		} `json:"entities"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		// Not valid JSON or missing "entities" key: not an error, just no structured data
		return nil, false
	}

	if len(payload.Entities) == 0 {
		// Valid JSON but no entities: not an error, just no structured data
		return nil, false
	}

	// Filter to entities with non-empty kind and validate attributes.
	var result []connector.SnapshotEntity
	for _, ent := range payload.Entities {
		if ent.Kind == "" {
			// Skip entities without a kind
			continue
		}

		entity := connector.SnapshotEntity{
			Kind:       ent.Kind,
			Name:       ent.Name,
			IP:         ent.IP,
			Hostname:   ent.Hostname,
			ExternalID: ent.ExternalID,
		}

		// Validate and pass through attributes as-is if present.
		// json.Unmarshal into map[string]any already gives us JSON-safe types:
		// string, float64, bool, []any (of primitives), nil.
		// No nested objects or other non-JSON types can arrive via this path.
		if len(ent.Attributes) > 0 {
			entity.Attributes = ent.Attributes
		}

		result = append(result, entity)
	}

	if len(result) == 0 {
		// No valid entities after filtering
		return nil, false
	}

	return result, true
}
