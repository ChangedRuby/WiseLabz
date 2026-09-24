// Package cloudflare implements a Cloudflare API connector covering DNS
// zones, Tunnels, and Zero Trust Access policies.
package cloudflare

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

const (
	typeName = "cloudflare"
	apiBase  = "https://api.cloudflare.com/client/v4"
)

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "dns",
		Name:     "Cloudflare",
		// External API: looser latency SLA than an on-LAN device.
		DegradedLatencyThresholdMs: 5000,
		Fields: []connector.SchemaField{
			{Key: "api_token", Label: "API Token", Type: "password", Required: true},
			{Key: "account_id", Label: "Account ID", Type: "text", Required: true},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		apiToken, _ := config["api_token"].(string)
		accountID, _ := config["account_id"].(string)
		if accountID != "" {
			if err := connector.ValidateRefSegment(accountID); err != nil {
				return nil, fmt.Errorf("invalid cloudflare account_id: %w", err)
			}
		}
		client := connector.NewHTTPClient(connector.HTTPClientOptions{})
		return &Connector{
			apiToken:  apiToken,
			accountID: accountID,
			client:    client,
			baseURL:   apiBase,
		}, nil
	})
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on "dns_record", "tunnel", and "policy" entities, exposed via
// GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"dns_record": {
		{Name: "type", Type: "string", Description: "DNS record type (A, AAAA, CNAME, TXT, ...)"},
		{Name: "proxied", Type: "boolean", Description: "Whether the record is proxied through Cloudflare"},
		{Name: "ttl", Type: "number", Description: "Record TTL in seconds"},
	},
	"tunnel": {
		{Name: "status", Type: "string", Description: "Tunnel connection status"},
		{Name: "connectorCount", Type: "number", Description: "Number of active cloudflared connectors for this tunnel"},
	},
	"policy": {
		{Name: "enabled", Type: "boolean", Description: "Whether the Access policy is active"},
		{Name: "decision", Type: "string", Description: "Policy decision (allow, deny, bypass, ...)"},
	},
}

// Connector fetches data from the Cloudflare API.
type Connector struct {
	apiToken  string
	accountID string
	client    *http.Client
	baseURL   string
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Cloudflare" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "dns" }

// Validate tests the connection to the Cloudflare API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	_, err := c.doRequest(ctx, "GET", "/zones")
	return err
}

// Fetch retrieves DNS zones/records, Tunnels, and Access policies.
func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"account_id": c.accountID}

	if content, zoneEntities, err := c.fetchDNSSection(ctx); err == nil {
		sections = append(sections, connector.SnapshotSection{Title: "DNS Zones", Content: content})
		entities = append(entities, zoneEntities...)
	} else {
		sections = append(sections, connector.SnapshotSection{Title: "DNS Zones", Content: "_DNS zones unavailable: " + err.Error() + "_"})
	}

	if raw, err := c.doRequest(ctx, "GET", "/accounts/"+c.accountID+"/cfd_tunnel"); err == nil {
		content, tunnelEntities := buildTunnelTable(raw)
		sections = append(sections, connector.SnapshotSection{Title: "Tunnels", Content: content})
		entities = append(entities, tunnelEntities...)
	} else {
		sections = append(sections, connector.SnapshotSection{Title: "Tunnels", Content: "_Tunnels unavailable: " + err.Error() + "_"})
	}

	if content, policyEntities, err := c.fetchAccessSection(ctx); err == nil {
		sections = append(sections, connector.SnapshotSection{Title: "Access Policies", Content: content})
		entities = append(entities, policyEntities...)
	} else {
		sections = append(sections, connector.SnapshotSection{Title: "Access Policies", Content: "_Access policies unavailable: " + err.Error() + "_"})
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Cloudflare",
		Type:        typeName,
		Sections:    sections,
		Entities:    entities,
		Metadata:    metadata,
		FetchedAt:   start,
	}, nil
}

// fetchDNSSection lists zones under the account, then fetches DNS records
// per zone. Each entity's ExternalID is "<zoneID>/<recordID>" so ConfigPush
// can recover the zone the record belongs to.
func (c *Connector) fetchDNSSection(ctx context.Context) (string, []connector.SnapshotEntity, error) {
	raw, err := c.doRequest(ctx, "GET", "/zones?account.id="+c.accountID)
	if err != nil {
		return "", nil, err
	}
	var zonesResp struct {
		Result []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &zonesResp); err != nil {
		return "", nil, connector.NewMalformedResponseError(err)
	}

	var b strings.Builder
	b.WriteString("| Zone | Name | Type | Content | Proxied | TTL |\n")
	b.WriteString("|------|------|------|---------|---------|-----|\n")
	var entities []connector.SnapshotEntity
	for _, zone := range zonesResp.Result {
		recRaw, err := c.doRequest(ctx, "GET", "/zones/"+zone.ID+"/dns_records")
		if err != nil {
			continue
		}
		rows, recEntities := buildDNSRecordTable(zone.ID, zone.Name, recRaw)
		b.WriteString(rows)
		entities = append(entities, recEntities...)
	}
	return b.String(), entities, nil
}

// fetchAccessSection lists Access applications under the account and their
// nested policies. Each policy entity's ExternalID is "<appID>/<policyID>".
func (c *Connector) fetchAccessSection(ctx context.Context) (string, []connector.SnapshotEntity, error) {
	raw, err := c.doRequest(ctx, "GET", "/accounts/"+c.accountID+"/access/apps")
	if err != nil {
		return "", nil, err
	}
	var appsResp struct {
		Result []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Policies []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Decision string `json:"decision"`
			} `json:"policies"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &appsResp); err != nil {
		return "", nil, connector.NewMalformedResponseError(err)
	}

	var b strings.Builder
	b.WriteString("| Application | Policy | Decision |\n")
	b.WriteString("|-------------|--------|----------|\n")
	var entities []connector.SnapshotEntity
	for _, app := range appsResp.Result {
		for _, p := range app.Policies {
			if _, err := fmt.Fprintf(&b, "| %s | %s | %s |\n", app.Name, p.Name, p.Decision); err != nil {
				return "", nil, err
			}
			enabled := p.Decision != "" && !strings.EqualFold(p.Decision, "deny")
			attrs := map[string]any{"enabled": enabled}
			if p.Decision != "" {
				attrs["decision"] = p.Decision
			}
			entities = append(entities, connector.SnapshotEntity{
				Kind:       "policy",
				Name:       p.Name,
				ExternalID: app.ID + "/" + p.ID,
				Attributes: attrs,
			})
		}
	}
	return b.String(), entities, nil
}

// WritableFields lists the config-push-eligible DNS record and Access
// policy fields.
func (c *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "proxied", Label: "DNS Record Proxied", Type: "toggle", EntityScope: true},
		{Key: "enabled", Label: "Access Policy Enabled", Type: "toggle", EntityScope: true},
	}
}

// ConfigPush toggles a writable field on the entity identified by entityRef.
// entityRef is "<zoneID>/<recordID>" for "proxied" and "<appID>/<policyID>"
// for "enabled", matching the ExternalID convention set when the entities
// were built in Fetch.
func (c *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	parts := strings.SplitN(entityRef, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("cloudflare config-push requires an entityRef of the form \"<parentID>/<id>\", got %q", entityRef)
	}
	if err := connector.ValidateCompositeRef(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}

	switch fieldKey {
	case "proxied":
		zoneID, recordID := parts[0], parts[1]
		proxied, _ := value.(bool)
		body, err := json.Marshal(map[string]any{"proxied": proxied})
		if err != nil {
			return err
		}
		_, err = c.doRequestBody(ctx, "PATCH", "/zones/"+zoneID+"/dns_records/"+recordID, body)
		return err
	case "enabled":
		appID, policyID := parts[0], parts[1]
		enabled, _ := value.(bool)
		decision := "allow"
		if !enabled {
			decision = "deny"
		}
		body, err := json.Marshal(map[string]any{"decision": decision})
		if err != nil {
			return err
		}
		_, err = c.doRequestBody(ctx, "PATCH", "/accounts/"+c.accountID+"/access/apps/"+appID+"/policies/"+policyID, body)
		return err
	default:
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
}

func (c *Connector) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	return c.doRequestBody(ctx, method, path, nil)
}

func (c *Connector) doRequestBody(ctx context.Context, method, path string, body []byte) (data []byte, err error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
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

func buildDNSRecordTable(zoneID, zoneName string, raw []byte) (string, []connector.SnapshotEntity) {
	var resp struct {
		Result []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
			Proxied bool   `json:"proxied"`
			TTL     int    `json:"ttl"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Result) == 0 {
		return "", nil
	}
	var b strings.Builder
	var entities []connector.SnapshotEntity
	for _, r := range resp.Result {
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s | %s | %t | %d |\n", zoneName, r.Name, r.Type, r.Content, r.Proxied, r.TTL); err != nil {
			return "", nil
		}
		attrs := map[string]any{"proxied": r.Proxied}
		if r.Type != "" {
			attrs["type"] = r.Type
		}
		if r.TTL != 0 {
			attrs["ttl"] = r.TTL
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "dns_record",
			Name:       r.Name,
			Hostname:   r.Name,
			ExternalID: zoneID + "/" + r.ID,
			Attributes: attrs,
		})
	}
	return b.String(), entities
}

func buildTunnelTable(raw []byte) (string, []connector.SnapshotEntity) {
	var resp struct {
		Result []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			Status         string `json:"status"`
			ConnsActiveAt  string `json:"conns_active_at"`
			ConnectorCount int    `json:"conns"`
		} `json:"result"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Result) == 0 {
		return "_No tunnel data returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Name | Status |\n")
	b.WriteString("|------|--------|\n")
	var entities []connector.SnapshotEntity
	for _, tun := range resp.Result {
		if _, err := fmt.Fprintf(&b, "| %s | %s |\n", tun.Name, tun.Status); err != nil {
			return "", nil
		}
		attrs := map[string]any{}
		if tun.Status != "" {
			attrs["status"] = tun.Status
		}
		if tun.ConnectorCount != 0 {
			attrs["connectorCount"] = tun.ConnectorCount
		}
		entities = append(entities, connector.SnapshotEntity{Kind: "tunnel", Name: tun.Name, ExternalID: tun.ID, Attributes: attrs})
	}
	return b.String(), entities
}
