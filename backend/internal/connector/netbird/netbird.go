// Package netbird implements a Netbird mesh VPN API connector.
package netbird

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "netbird"

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "networking",
		Name:     "Netbird",
		// External API (self-hosted or api.netbird.io): looser latency SLA
		// than an on-LAN device.
		DegradedLatencyThresholdMs: 5000,
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Netbird API URL", Type: "text", Required: true, Default: "https://api.netbird.io", Placeholder: "https://api.netbird.io"},
			{Key: "api_token", Label: "API Token", Type: "password", Required: true},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Default: "true"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		url, _ := config["url"].(string)
		apiToken, _ := config["api_token"].(string)
		verifyTLS := true
		if v, ok := config["verify_tls"]; ok {
			if b, ok := v.(bool); ok {
				verifyTLS = b
			}
		}
		client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS})
		return &Connector{
			url:      strings.TrimSuffix(url, "/"),
			apiToken: apiToken,
			client:   client,
		}, nil
	})
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on "peer" and "policy" entities (see buildPeerTable/buildPolicyTable),
// exposed via GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"peer": {
		{Name: "approved", Type: "boolean", Description: "Whether the peer has been approved to join the network"},
		{Name: "connected", Type: "boolean", Description: "Whether the peer is currently connected"},
		{Name: "os", Type: "string", Description: "Operating system reported by the peer"},
		{Name: "groups", Type: "string_array", Description: "Group names the peer belongs to"},
	},
	"policy": {
		{Name: "enabled", Type: "boolean", Description: "Whether the access policy is active"},
		{Name: "protocol", Type: "string", Description: "Protocol matched by the policy (tcp, udp, all, ...)"},
		{Name: "sourceGroups", Type: "string_array", Description: "Source group names the policy matches"},
		{Name: "destinationGroups", Type: "string_array", Description: "Destination group names the policy matches"},
	},
}

// Connector fetches data from the Netbird management API.
type Connector struct {
	url      string
	apiToken string
	client   *http.Client
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Netbird" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "networking" }

// Validate tests the connection to the Netbird API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	_, err := c.doRequest(ctx, "GET", "/api/peers")
	return err
}

// Fetch retrieves peers, network routes, and access policies.
func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"netbird_url": c.url}

	if raw, err := c.doRequest(ctx, "GET", "/api/peers"); err == nil {
		content, peerEntities := buildPeerTable(raw)
		sections = append(sections, connector.SnapshotSection{Title: "Peers", Content: content})
		entities = append(entities, peerEntities...)
		metadata["peer_count"] = fmt.Sprintf("%d", len(peerEntities))
	} else {
		sections = append(sections, connector.SnapshotSection{Title: "Peers", Content: "_Peers unavailable: " + err.Error() + "_"})
	}

	if raw, err := c.doRequest(ctx, "GET", "/api/routes"); err == nil {
		sections = append(sections, connector.SnapshotSection{Title: "Network Routes", Content: buildRouteTable(raw)})
	} else {
		sections = append(sections, connector.SnapshotSection{Title: "Network Routes", Content: "_Routes unavailable: " + err.Error() + "_"})
	}

	if raw, err := c.doRequest(ctx, "GET", "/api/policies"); err == nil {
		content, policyEntities := buildPolicyTable(raw)
		sections = append(sections, connector.SnapshotSection{Title: "Access Policies", Content: content})
		entities = append(entities, policyEntities...)
	} else {
		sections = append(sections, connector.SnapshotSection{Title: "Access Policies", Content: "_Policies unavailable: " + err.Error() + "_"})
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Netbird",
		Type:        typeName,
		Sections:    sections,
		Entities:    entities,
		Metadata:    metadata,
		FetchedAt:   start,
	}, nil
}

// WritableFields lists the config-push-eligible peer approval field.
func (c *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "approved", Label: "Peer Approved", Type: "toggle", EntityScope: true},
	}
}

// ConfigPush toggles the "approved" state of the peer identified by
// entityRef (the Netbird peer ID).
func (c *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("netbird config-push requires a target peer ID")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	if fieldKey != "approved" {
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
	approved, _ := value.(bool)
	body, err := json.Marshal(map[string]any{"approval_required": !approved})
	if err != nil {
		return err
	}
	_, err = c.doRequestBody(ctx, "PUT", "/api/peers/"+entityRef, body)
	return err
}

func (c *Connector) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	return c.doRequestBody(ctx, method, path, nil)
}

func (c *Connector) doRequestBody(ctx context.Context, method, path string, body []byte) (data []byte, err error) {
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
	req.Header.Set("Authorization", "Token "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, connector.MapTransportError(err)
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

	if statusErr := connector.CheckStatus(resp.StatusCode, data); statusErr != nil {
		return nil, statusErr
	}

	return data, nil
}

func buildPeerTable(raw []byte) (string, []connector.SnapshotEntity) {
	var peers []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		IP        string `json:"ip"`
		OS        string `json:"os"`
		Approved  bool   `json:"approval_required"`
		Connected bool   `json:"connected"`
		Groups    []struct {
			Name string `json:"name"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(raw, &peers); err != nil || len(peers) == 0 {
		return "_No peer data returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Name | IP | OS | Connected | Approved |\n")
	b.WriteString("|------|----|----|-----------|----------|\n")
	var entities []connector.SnapshotEntity
	for _, p := range peers {
		approved := !p.Approved // API field is "approval_required"; approved means it's no longer required
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s | %t | %t |\n", p.Name, p.IP, p.OS, p.Connected, approved); err != nil {
			return "", nil
		}
		var groupNames []string
		for _, g := range p.Groups {
			groupNames = append(groupNames, g.Name)
		}
		attrs := map[string]any{
			"approved":  approved,
			"connected": p.Connected,
		}
		if p.OS != "" {
			attrs["os"] = p.OS
		}
		if len(groupNames) > 0 {
			attrs["groups"] = groupNames
		}
		entities = append(entities, connector.SnapshotEntity{Kind: "peer", Name: p.Name, IP: p.IP, ExternalID: p.ID, Attributes: attrs})
	}
	return b.String(), entities
}

func buildRouteTable(raw []byte) string {
	var routes []struct {
		ID      string `json:"id"`
		Network string `json:"network"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.Unmarshal(raw, &routes); err != nil || len(routes) == 0 {
		return "_No route data returned_"
	}
	var b strings.Builder
	b.WriteString("| Network | Enabled |\n")
	b.WriteString("|---------|---------|\n")
	for _, r := range routes {
		if _, err := fmt.Fprintf(&b, "| %s | %t |\n", r.Network, r.Enabled); err != nil {
			return ""
		}
	}
	return b.String()
}

func buildPolicyTable(raw []byte) (string, []connector.SnapshotEntity) {
	var policies []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
		Rules   []struct {
			Protocol string `json:"protocol"`
			Sources  []struct {
				Name string `json:"name"`
			} `json:"sources"`
			Destinations []struct {
				Name string `json:"name"`
			} `json:"destinations"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(raw, &policies); err != nil || len(policies) == 0 {
		return "_No access policy data returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Name | Enabled |\n")
	b.WriteString("|------|---------|\n")
	var entities []connector.SnapshotEntity
	for _, p := range policies {
		if _, err := fmt.Fprintf(&b, "| %s | %t |\n", p.Name, p.Enabled); err != nil {
			return "", nil
		}
		attrs := map[string]any{"enabled": p.Enabled}
		if len(p.Rules) > 0 {
			r := p.Rules[0]
			if r.Protocol != "" {
				attrs["protocol"] = r.Protocol
			}
			var sourceGroups, destGroups []string
			for _, s := range r.Sources {
				sourceGroups = append(sourceGroups, s.Name)
			}
			for _, d := range r.Destinations {
				destGroups = append(destGroups, d.Name)
			}
			if len(sourceGroups) > 0 {
				attrs["sourceGroups"] = sourceGroups
			}
			if len(destGroups) > 0 {
				attrs["destinationGroups"] = destGroups
			}
		}
		entities = append(entities, connector.SnapshotEntity{Kind: "policy", Name: p.Name, ExternalID: p.ID, Attributes: attrs})
	}
	return b.String(), entities
}
