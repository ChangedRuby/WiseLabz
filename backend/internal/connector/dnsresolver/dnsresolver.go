// Package dnsresolver implements a pfSense/OPNsense DNS Resolver (Unbound)
// connector, targeting the jaredhendrickson13/pfsense-api v2 REST plugin's
// DNS Resolver host-override endpoint.
package dnsresolver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "dnsresolver"

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
		dialer := connector.GuardedDialer(30 * time.Second)
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DialContext:     dialer.DialContext,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyTLS},
			},
		}
		return &Connector{
			url:    strings.TrimSuffix(url, "/"),
			apiKey: apiKey,
			client: client,
		}, nil
	})
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

func (c *Connector) doRequest(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
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

func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
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
		entities = append(entities, connector.SnapshotEntity{
			Kind:     "dns_record",
			Hostname: hostname,
			IP:       o.IP,
		})
	}
	return b.String(), entities
}
