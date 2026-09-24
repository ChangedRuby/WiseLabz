// Package dnsresolver implements a pfSense/OPNsense DNS Resolver (Unbound)
// connector, targeting the jaredhendrickson13/pfsense-api v2 REST plugin's
// DNS Resolver host-override endpoint.
package dnsresolver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "dnsresolver"

// attributeCatalog declares the structured Attributes this connector fills
// on "dns_record" entities (see buildHostOverrideTable), exposed via GET
// /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"dns_record": {
		{Name: "description", Type: "string", Description: "Description of the DNS record"},
		{Name: "is_ipv6", Type: "boolean", Description: "Whether the IP address is IPv6"},
	},
}

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "dns",
		Name:     "DNS Resolver",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Appliance URL", Type: "text", Required: true, Placeholder: "https://pfsense.example.com"},
			{Key: "api_key", Label: "API Key", Type: "password", Required: true},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		url, _ := config["url"].(string)
		apiKey, _ := config["api_key"].(string)
		verifyTLS := true
		if v, ok := config["verify_tls"]; ok {
			if b, ok := v.(bool); ok {
				verifyTLS = b
			}
		}
		client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS})
		return &Connector{
			url:    strings.TrimSuffix(url, "/"),
			apiKey: apiKey,
			client: client,
		}, nil
	})
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// Connector fetches DNS Resolver (Unbound) host overrides from a
// pfSense/OPNsense appliance API.
type Connector struct {
	url    string
	apiKey string
	client *http.Client
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "DNS Resolver" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "dns" }

// Validate tests the connection to the DNS Resolver API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	_, err := c.doRequest(ctx, "/api/v2/services/dns_resolver/host_override")
	return err
}

// Fetch retrieves DNS Resolver host overrides.
func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"dnsresolver_url": c.url}

	if raw, err := c.doRequest(ctx, "/api/v2/services/dns_resolver/host_override"); err != nil {
		sections = append(sections, connector.SnapshotSection{Title: "Host Overrides", Content: "_Host overrides unavailable: " + err.Error() + "_"})
	} else {
		content, ents := buildHostOverrideTable(raw)
		sections = append(sections, connector.SnapshotSection{Title: "Host Overrides", Content: content})
		entities = ents
	}

	return &connector.ServiceSnapshot{
		ServiceName: "DNS Resolver",
		Type:        typeName,
		Sections:    sections,
		Entities:    entities,
		Metadata:    metadata,
		FetchedAt:   start,
	}, nil
}

// WritableFields lists the config-push-eligible host override field. No
// Starter/Stopper here: like PR1's restart, this connector monitors a DNS
// Resolver service it doesn't own the lifecycle of (see the pfSense/OPNsense
// connectors for that role) — same scope decision, mirrored.
func (c *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "ip", Label: "Host Override IP", Type: "text", EntityScope: true},
	}
}

// ConfigPush repoints the host override identified by entityRef (its
// "host.domain" hostname, as built by buildHostOverrideTable) to a new IP.
// ponytail: one field (ip) per call, matching the handler's one-field
// revert contract.
func (c *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("dnsresolver config-push requires a target hostname")
	}
	if fieldKey != "ip" {
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
	raw, err := c.doRequest(ctx, "/api/v2/services/dns_resolver/host_override")
	if err != nil {
		return fmt.Errorf("resolve host override id: %w", err)
	}
	var resp struct {
		Data []struct {
			ID     json.Number `json:"id"`
			Host   string      `json:"host"`
			Domain string      `json:"domain"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return connector.NewMalformedResponseError(fmt.Errorf("decode host overrides: %w", err))
	}
	for _, o := range resp.Data {
		hostname := o.Domain
		if o.Host != "" {
			hostname = fmt.Sprintf("%s.%s", o.Host, o.Domain)
		}
		if hostname != entityRef {
			continue
		}
		body, err := json.Marshal(map[string]any{"id": o.ID, "ip": value})
		if err != nil {
			return err
		}
		return c.doRequestBody(ctx, "PATCH", "/api/v2/services/dns_resolver/host_override", body)
	}
	return fmt.Errorf("host override %q not found", entityRef)
}

func (c *Connector) doRequest(ctx context.Context, path string) ([]byte, error) {
	return c.doRequestBodyRaw(ctx, "GET", path, nil)
}

func (c *Connector) doRequestBody(ctx context.Context, method, path string, body []byte) error {
	_, err := c.doRequestBodyRaw(ctx, method, path, body)
	return err
}

func (c *Connector) doRequestBodyRaw(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.url+path, reqBody)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, connector.MapTransportError(err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := connector.ReadBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if statusErr := connector.CheckStatus(resp.StatusCode, data); statusErr != nil {
		return nil, statusErr
	}

	return data, nil
}

// isIPv6 returns true if the IP string is an IPv6 address.
func isIPv6(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.To4() == nil
}

// buildHostOverrideTable renders the DNS Resolver host overrides as a
// markdown table and extracts one SnapshotEntity per override.
func buildHostOverrideTable(raw []byte) (content string, entities []connector.SnapshotEntity) {
	var resp struct {
		Data []struct {
			Host        string `json:"host"`
			Domain      string `json:"domain"`
			IP          string `json:"ip"`
			Description string `json:"descr"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "_Host overrides unavailable: " + connector.NewMalformedResponseError(err).Error() + "_", nil
	}
	if len(resp.Data) == 0 {
		return "_No host overrides returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Host | Domain | IP | Description |\n")
	b.WriteString("|------|--------|----|--------------|\n")
	for _, o := range resp.Data {
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", o.Host, o.Domain, o.IP, o.Description); err != nil {
			return "", nil
		}
		hostname := o.Domain
		if o.Host != "" {
			hostname = fmt.Sprintf("%s.%s", o.Host, o.Domain)
		}
		attrs := map[string]any{
			"is_ipv6": isIPv6(o.IP),
		}
		if o.Description != "" {
			attrs["description"] = o.Description
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "dns_record",
			Hostname:   hostname,
			IP:         o.IP,
			Attributes: attrs,
		})
	}
	return b.String(), entities
}
