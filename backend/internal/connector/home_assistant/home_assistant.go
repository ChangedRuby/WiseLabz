// Package homeassistant implements a Home Assistant connector built on the
// Home Assistant REST API, authenticated with a long-lived access token.
//
// The REST API is deliberately narrower than the WebSocket API: it exposes
// the instance configuration (/api/config, which also carries the loaded
// integrations), every entity and its state (/api/states) and the callable
// service registry (/api/services). Areas, devices and the config-entry
// registry are WebSocket-only, so this connector derives what it can from
// entities and components instead and says so in the generated sections.
//
// Snapshots are deliberately value-stable: Home Assistant reports
// last_changed/last_updated on every entity and a live measurement for every
// sensor, so an unchanged instance would otherwise produce a different
// snapshot on every sync. Volatile values are collapsed to availability —
// see stableState.
//
// Note: the HTTP plumbing below (guarded dialer, TLS config, status-code to
// connector-error mapping, timeout detection) is deliberately duplicated
// from the sibling connectors rather than shared — #265/#266 will rewrite
// the connector base, and this package stays self-contained until then.
package homeassistant

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const (
	typeName = "home_assistant"
	// category is constrained by the connectors.category CHECK in the
	// store migrations; Home Assistant is filed under the same
	// general-purpose bucket the custom connector uses.
	category = "virtualization"
	// defaultMaxEntities caps how many entities a snapshot carries. A
	// mature Home Assistant instance exposes thousands of entities, most
	// of them diagnostic; the cap keeps snapshots (and their diffs)
	// reviewable. 0 means "no cap".
	defaultMaxEntities = 500
)

// Home Assistant REST API paths this connector reads.
const (
	pathConfig   = "/api/config"
	pathStates   = "/api/states"
	pathServices = "/api/services"
)

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: category,
		Name:     "Home Assistant",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Home Assistant URL", Type: "text", Required: true, Placeholder: "http://homeassistant.local:8123", Description: "Base URL of the Home Assistant instance (the host serving /api/config)."},
			{Key: "access_token", Label: "Long-Lived Access Token", Type: "password", Required: true, Description: "Long-lived access token created from a Home Assistant user profile; sent as \"Authorization: Bearer <token>\"."},
			{Key: "max_entities", Label: "Max Entities", Type: "number", Required: false, Default: strconv.Itoa(defaultMaxEntities), Description: "Maximum number of entities to include in a snapshot, sorted by entity ID. 0 includes every entity."},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, newConnector)
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on entity and integration entities, exposed via GET /api/compliance/schema.
// Every value here is stable across syncs — no measurements, no timestamps.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"entity": {
		{Name: "domain", Type: "string", Description: "Entity domain, i.e. the part of the entity ID before the dot (light, switch, sensor, ...)"},
		{Name: "state", Type: "string", Description: "Entity state, reported only when it is a discrete value (on, off, unavailable, ...); continuously changing measurements are omitted"},
		{Name: "available", Type: "boolean", Description: "Whether Home Assistant currently has a value for the entity (state is neither unavailable nor unknown)"},
		{Name: "deviceClass", Type: "string", Description: "Home Assistant device class of the entity (motion, temperature, door, ...)"},
		{Name: "unitOfMeasurement", Type: "string", Description: "Unit the entity reports its value in"},
		{Name: "entityCategory", Type: "string", Description: "Entity category assigned by its integration (config, diagnostic) — empty for primary entities"},
		{Name: "supportedFeatures", Type: "number", Description: "Bitmask of the optional features the entity's platform supports"},
	},
	"integration": {
		{Name: "platforms", Type: "string_array", Description: "Platforms this integration loaded, e.g. \"sensor\" for the component \"sensor.mqtt\""},
	},
}

// Connector fetches configuration, entities and services from a Home
// Assistant instance.
type Connector struct {
	url         string
	accessToken string
	maxEntities int
	client      *http.Client
}

func newConnector(config map[string]any) (connector.Connector, error) {
	rawURL, _ := config["url"].(string)
	accessToken, _ := config["access_token"].(string)
	verifyTLS := true
	if v, ok := config["verify_tls"]; ok {
		if b, ok := v.(bool); ok {
			verifyTLS = b
		}
	}
	dialer := connector.GuardedDialer(30 * time.Second)
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext:     dialer.DialContext,
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: !verifyTLS},
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Connector{
		url:         strings.TrimSuffix(rawURL, "/"),
		accessToken: accessToken,
		maxEntities: intConfig(config, "max_entities", defaultMaxEntities),
		client:      client,
	}, nil
}

// intConfig reads a numeric config value, which reaches the connector as a
// JSON number, a stringified number (the schema Default is a string) or not
// at all. Anything unparseable falls back to fallback.
func intConfig(config map[string]any, key string, fallback int) int {
	raw, ok := config[key]
	if !ok {
		return fallback
	}
	switch v := raw.(type) {
	case int:
		return max(v, 0)
	case float64:
		return max(int(v), 0)
	case string:
		if v == "" {
			return fallback
		}
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil || n < 0 {
			return fallback
		}
		return n
	default:
		return fallback
	}
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Home Assistant" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return category }

// Validate checks the credentials are present and tests them against the
// Home Assistant API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	if c.url == "" {
		return fmt.Errorf("home assistant url is required")
	}
	if c.accessToken == "" {
		return fmt.Errorf("home assistant access_token is required")
	}
	_, err := c.doRequest(ctx, pathConfig)
	return err
}

// Fetch retrieves the instance configuration, the loaded integrations, the
// entity registry with stable state values, and the callable services. Each
// section degrades on its own: one failing endpoint leaves a placeholder
// instead of failing the whole snapshot.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	metadata := map[string]string{"home_assistant_url": c.url}
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity

	wantOverview := connector.WantsField(fields, "config")
	wantIntegrations := connector.WantsField(fields, "integrations")
	if wantOverview || wantIntegrations {
		// Both sections come out of the same payload, so it is fetched once.
		raw, err := c.doRequest(ctx, pathConfig)
		switch {
		case err != nil:
			if wantOverview {
				sections = append(sections, unavailable("Overview", err))
			}
			if wantIntegrations {
				sections = append(sections, unavailable("Integrations", err))
			}
		default:
			if wantOverview {
				content, meta := buildOverview(raw)
				sections = append(sections, connector.SnapshotSection{Title: "Overview", Content: content})
				for k, v := range meta {
					metadata[k] = v
				}
			}
			if wantIntegrations {
				content, ents := buildIntegrations(raw)
				sections = append(sections, connector.SnapshotSection{Title: "Integrations", Content: content})
				entities = append(entities, ents...)
				metadata["integration_count"] = strconv.Itoa(len(ents))
			}
		}
	}

	if connector.WantsField(fields, "entities") {
		raw, err := c.doRequest(ctx, pathStates)
		if err != nil {
			sections = append(sections, unavailable("Entity Domains", err), unavailable("Entities", err))
		} else {
			result := buildEntities(raw, c.maxEntities)
			sections = append(sections,
				connector.SnapshotSection{Title: "Entity Domains", Content: result.domains},
				connector.SnapshotSection{Title: "Entities", Content: result.entityTable},
			)
			entities = append(entities, result.entities...)
			for k, v := range result.metadata {
				metadata[k] = v
			}
		}
	}

	if connector.WantsField(fields, "services") {
		raw, err := c.doRequest(ctx, pathServices)
		if err != nil {
			sections = append(sections, unavailable("Services", err))
		} else {
			content, meta := buildServices(raw)
			sections = append(sections, connector.SnapshotSection{Title: "Services", Content: content})
			for k, v := range meta {
				metadata[k] = v
			}
		}
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Home Assistant",
		Type:        typeName,
		Sections:    sections,
		Entities:    entities,
		Metadata:    metadata,
		FetchedAt:   start,
	}, nil
}

func unavailable(title string, err error) connector.SnapshotSection {
	return connector.SnapshotSection{Title: title, Content: "_" + title + " unavailable: " + err.Error() + "_"}
}

func (c *Connector) doRequest(ctx context.Context, path string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	data, err = connector.ReadBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return nil, connector.NewAuthError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
	case resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout:
		return nil, connector.NewServiceUnavailableError(fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data)))
	case resp.StatusCode >= 400:
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// isTimeout reports whether err represents a request deadline being
// exceeded, covering both a canceled context and a net.Error timeout.
func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
