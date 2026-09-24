// Package opnsense implements an OPNSense firewall API connector.
package opnsense

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

const typeName = "opnsense"

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "networking",
		Name:     "OPNSense",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "OPNSense URL", Type: "text", Required: true, Placeholder: "https://opnsense.example.com"},
			{Key: "api_key", Label: "API Key", Type: "password", Required: true},
			{Key: "api_secret", Label: "API Secret", Type: "password", Required: true},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Default: "true"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		url, _ := config["url"].(string)
		apiKey, _ := config["api_key"].(string)
		apiSecret, _ := config["api_secret"].(string)
		verifyTLS := true
		if v, ok := config["verify_tls"]; ok {
			if b, ok := v.(bool); ok {
				verifyTLS = b
			}
		}
		client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS})
		return &Connector{
			url:       strings.TrimSuffix(url, "/"),
			apiKey:    apiKey,
			apiSecret: apiSecret,
			client:    client,
		}, nil
	})
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on "interface" and "rule" entities (see buildInterfaceTable/buildRuleTable),
// exposed via GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"interface": {
		{Name: "enabled", Type: "boolean", Description: "Whether the interface is enabled"},
		{Name: "type", Type: "string", Description: "Interface addressing type (static, dhcp, ppoe, ...)"},
		{Name: "ipv4", Type: "string", Description: "IPv4 address assigned to the interface"},
		{Name: "ipv6", Type: "string", Description: "IPv6 address assigned to the interface"},
		{Name: "gateway", Type: "string", Description: "Gateway configured for the interface"},
	},
	"rule": {
		{Name: "enabled", Type: "boolean", Description: "Whether the firewall rule is enabled"},
		{Name: "action", Type: "string", Description: "Rule action: pass, block, or reject"},
		{Name: "interface", Type: "string", Description: "Interface the rule applies to"},
		{Name: "direction", Type: "string", Description: "Traffic direction the rule matches (in/out)"},
		{Name: "protocol", Type: "string", Description: "Protocol matched by the rule (tcp, udp, any, ...)"},
		{Name: "source", Type: "string", Description: "Source address/network matched by the rule"},
		{Name: "destination", Type: "string", Description: "Destination address/network matched by the rule"},
		{Name: "destination_port", Type: "string", Description: "Destination port or port range matched by the rule"},
		{Name: "log", Type: "boolean", Description: "Whether matching packets are logged"},
		{Name: "disabled_reason", Type: "string", Description: "Reason the rule was auto-disabled, if any"},
	},
}

// Connector fetches data from an OPNSense firewall API.
type Connector struct {
	url       string
	apiKey    string
	apiSecret string
	client    *http.Client
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "OPNSense" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "networking" }

// Validate tests the connection to the OPNSense API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	_, err := c.doRequest(ctx, "GET", "/api/core/firmware/status")
	return err
}

// Fetch retrieves firewall rules, interfaces, gateways, and system health.
func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	var sections []connector.SnapshotSection
	var dependencies []connector.ServiceDependency
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"opnsense_url": c.url}

	// --- System info ---
	if raw, err := c.doRequest(ctx, "GET", "/api/core/firmware/status"); err == nil {
		var info struct {
			Version     string `json:"product_version"`
			ProductName string `json:"product_name"`
		}
		if err := json.Unmarshal(raw, &info); err != nil {
			sections = append(sections, connector.SnapshotSection{
				Title:   "System",
				Content: "_System info unavailable: " + connector.NewMalformedResponseError(err).Error() + "_",
			})
		} else {
			content := fmt.Sprintf("**Product**: %s\n**Version**: %s\n", info.ProductName, info.Version)
			sections = append(sections, connector.SnapshotSection{
				Title:   "System",
				Content: content,
			})
			metadata["version"] = info.Version
		}
	} else {
		sections = append(sections, connector.SnapshotSection{
			Title:   "System",
			Content: "_System info unavailable: " + err.Error() + "_",
		})
	}

	// --- Interfaces ---
	if raw, err := c.doRequest(ctx, "GET", "/api/diagnostics/interface/getInterfaces"); err == nil {
		content, ifaceEntities := buildInterfaceTable(raw)
		sections = append(sections, connector.SnapshotSection{
			Title:   "Interfaces",
			Content: content,
		})
		entities = append(entities, ifaceEntities...)
		if wan := wanInterfaceName(raw); wan != "" {
			dependencies = append(dependencies, connector.ServiceDependency{Kind: "network", Name: wan})
		}
	} else {
		sections = append(sections, connector.SnapshotSection{
			Title:   "Interfaces",
			Content: "_Interfaces unavailable: " + err.Error() + "_",
		})
	}

	// --- Firewall rules ---
	if raw, err := c.doRequest(ctx, "GET", "/api/firewall/filter/searchRule"); err == nil {
		content, ruleEntities := buildRuleTable(raw)
		sections = append(sections, connector.SnapshotSection{
			Title:   "Firewall Rules",
			Content: content,
		})
		entities = append(entities, ruleEntities...)
	} else {
		sections = append(sections, connector.SnapshotSection{
			Title:   "Firewall Rules",
			Content: "_Rules unavailable: " + err.Error() + "_",
		})
	}

	// --- Gateways ---
	if raw, err := c.doRequest(ctx, "GET", "/api/routes/gateway/status"); err == nil {
		content := buildGatewayTable(raw)
		sections = append(sections, connector.SnapshotSection{
			Title:   "Gateways",
			Content: content,
		})
		if upstream := primaryGatewayName(raw); upstream != "" {
			dependencies = append(dependencies, connector.ServiceDependency{Kind: "upstream_service", Name: upstream})
		}
	} else {
		sections = append(sections, connector.SnapshotSection{
			Title:   "Gateways",
			Content: "_Gateways unavailable: " + err.Error() + "_",
		})
	}

	return &connector.ServiceSnapshot{
		ServiceName:  "OPNSense",
		Type:         typeName,
		Sections:     sections,
		Dependencies: dependencies,
		Entities:     entities,
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

// Restart restarts the service identified by entityRef (an OPNSense service
// name, e.g. "unbound" or "dpinger") via the core service-control API.
func (c *Connector) Restart(ctx context.Context, _ map[string]any, entityRef string) error {
	if entityRef == "" {
		return fmt.Errorf("opnsense restart requires a target service name")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	raw, err := c.doRequest(ctx, "POST", "/api/core/service/restart/"+entityRef)
	if err != nil {
		return err
	}
	var resp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return connector.NewMalformedResponseError(fmt.Errorf("decode restart response: %w", err))
	}
	if resp.Status != "ok" {
		return fmt.Errorf("restart failed: status %q", resp.Status)
	}
	return nil
}

// Start starts the service identified by entityRef via the core
// service-control API. Idempotent-safe: OPNSense returns "ok" for an
// already-running service.
func (c *Connector) Start(ctx context.Context, _ map[string]any, entityRef string) error {
	return c.serviceAction(ctx, entityRef, "start")
}

// Stop stops the service identified by entityRef via the core
// service-control API.
func (c *Connector) Stop(ctx context.Context, _ map[string]any, entityRef string) error {
	return c.serviceAction(ctx, entityRef, "stop")
}

func (c *Connector) serviceAction(ctx context.Context, entityRef, action string) error {
	if entityRef == "" {
		return fmt.Errorf("opnsense %s requires a target service name", action)
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	raw, err := c.doRequest(ctx, "POST", "/api/core/service/"+action+"/"+entityRef)
	if err != nil {
		return err
	}
	var resp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return connector.NewMalformedResponseError(fmt.Errorf("decode %s response: %w", action, err))
	}
	if resp.Status != "ok" {
		return fmt.Errorf("%s failed: status %q", action, resp.Status)
	}
	return nil
}

// WritableFields lists the config-push-eligible firewall rule field.
func (c *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "enabled", Label: "Rule Enabled", Type: "toggle", EntityScope: true},
	}
}

// ConfigPush toggles the "enabled" state of the firewall rule identified by
// entityRef (the rule UUID) and applies the change.
func (c *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("opnsense config-push requires a target rule UUID")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	if fieldKey != "enabled" {
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
	enabled := "0"
	if b, _ := value.(bool); b {
		enabled = "1"
	}
	body, err := json.Marshal(map[string]any{"rule": map[string]string{"enabled": enabled}})
	if err != nil {
		return err
	}
	if _, err := c.doRequestBody(ctx, "POST", "/api/firewall/filter/setRule/"+entityRef, body); err != nil {
		return err
	}
	_, err = c.doRequest(ctx, "POST", "/api/firewall/filter/apply")
	return err
}

func (c *Connector) doRequest(ctx context.Context, method, path string) (data []byte, err error) {
	return c.doRequestBody(ctx, method, path, nil)
}

func (c *Connector) doRequestBody(ctx context.Context, method, path string, body []byte) (data []byte, err error) {
	url := c.url + path
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetBasicAuth(c.apiKey, c.apiSecret)
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

func buildInterfaceTable(raw []byte) (string, []connector.SnapshotEntity) {
	var resp struct {
		Rows []struct {
			Device    string `json:"device"`
			IPAddress string `json:"ipaddr"`
			IPv6      string `json:"ipv6"`
			Status    string `json:"status"`
			Media     string `json:"media"`
			Enabled   bool   `json:"enabled"`
			Type      string `json:"type"`
			Gateway   string `json:"gateway"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Rows) == 0 {
		return "_No interface data returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Device | IP Address | Status | Media |\n")
	b.WriteString("|--------|------------|--------|-------|\n")
	var entities []connector.SnapshotEntity
	for _, iface := range resp.Rows {
		_, err := fmt.Fprintf(&b, "| %s | %s | %s | %s |\n",
			iface.Device, iface.IPAddress, iface.Status, iface.Media)
		if err != nil {
			return "", nil
		}
		attrs := map[string]any{
			"enabled": iface.Enabled,
			"ipv4":    iface.IPAddress,
			"ipv6":    iface.IPv6,
		}
		if iface.Type != "" {
			attrs["type"] = iface.Type
		}
		if iface.Gateway != "" {
			attrs["gateway"] = iface.Gateway
		}
		entities = append(entities, connector.SnapshotEntity{Kind: "interface", Name: iface.Device, IP: iface.IPAddress, Attributes: attrs})
	}
	return b.String(), entities
}

// wanInterfaceName returns the device name of the interface identified as
// "wan" in the interfaces response, reusing data already fetched for the
// Interfaces section. Falls back to the first interface if none is
// explicitly identified as WAN.
// ponytail: identifier-based match with a first-row fallback; revisit if
// multi-WAN setups need every WAN link surfaced.
func wanInterfaceName(raw []byte) string {
	var resp struct {
		Rows []struct {
			Identifier string `json:"identifier"`
			Device     string `json:"device"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Rows) == 0 {
		return ""
	}
	for _, iface := range resp.Rows {
		if strings.EqualFold(iface.Identifier, "wan") {
			return iface.Device
		}
	}
	return resp.Rows[0].Device
}

// primaryGatewayName returns the name (or address, if unnamed) of the first
// configured gateway, reusing data already fetched for the Gateways
// section, to surface as the upstream dependency.
func primaryGatewayName(raw []byte) string {
	var resp struct {
		Items []struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Items) == 0 {
		return ""
	}
	if resp.Items[0].Name != "" {
		return resp.Items[0].Name
	}
	return resp.Items[0].Address
}

func buildRuleTable(raw []byte) (string, []connector.SnapshotEntity) {
	var resp struct {
		Rows []struct {
			Description     string `json:"description"`
			Action          string `json:"action"`
			Protocol        string `json:"protocol"`
			Source          string `json:"source_net"`
			Destination     string `json:"destination_net"`
			DestinationPort string `json:"destination_port"`
			Interface       string `json:"interface"`
			Direction       string `json:"direction"`
			Enabled         string `json:"enabled"`
			Log             string `json:"log"`
			DisabledReason  string `json:"disabled_reason"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Rows) == 0 {
		return "_No firewall rules returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Description | Action | Protocol | Source | Destination | Enabled |\n")
	b.WriteString("|-------------|--------|----------|--------|-------------|--------|\n")
	var entities []connector.SnapshotEntity
	count := 0
	for _, r := range resp.Rows {
		if count >= 50 {
			_, err := fmt.Fprintf(&b, "\n_...and %d more rules_", len(resp.Rows)-50)
			if err != nil {
				return "", nil
			}
			break
		}
		enabledStr := r.Enabled
		if enabledStr == "" {
			enabledStr = "1"
		}
		_, err := fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n",
			r.Description, r.Action, r.Protocol, r.Source, r.Destination, enabledStr)
		if err != nil {
			return "", nil
		}
		attrs := map[string]any{
			"enabled":     enabledStr == "1",
			"action":      r.Action,
			"protocol":    r.Protocol,
			"source":      r.Source,
			"destination": r.Destination,
			"log":         r.Log == "1",
		}
		if r.Interface != "" {
			attrs["interface"] = r.Interface
		}
		if r.Direction != "" {
			attrs["direction"] = r.Direction
		}
		if r.DestinationPort != "" {
			attrs["destination_port"] = r.DestinationPort
		}
		if r.DisabledReason != "" {
			attrs["disabled_reason"] = r.DisabledReason
		}
		entities = append(entities, connector.SnapshotEntity{Kind: "rule", Name: r.Description, Attributes: attrs})
		count++
	}
	return b.String(), entities
}

func buildGatewayTable(raw []byte) string {
	var resp struct {
		Items []struct {
			Name    string `json:"name"`
			Address string `json:"address"`
			Status  string `json:"status"`
			RTT     string `json:"rtt"`
			Loss    string `json:"loss"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Items) == 0 {
		return "_No gateway data returned_"
	}
	var b strings.Builder
	b.WriteString("| Gateway | Address | Status | RTT | Loss |\n")
	b.WriteString("|---------|---------|--------|-----|------|\n")
	for _, gw := range resp.Items {
		_, err := fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			gw.Name, gw.Address, gw.Status, gw.RTT, gw.Loss)
		if err != nil {
			return ""
		}
	}
	return b.String()
}
