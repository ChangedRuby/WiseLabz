// Package proxmox implements a real Proxmox VE API connector.
package proxmox

import (
	"context"
	"net/http"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "proxmox"

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "virtualization",
		Name:     "Proxmox VE",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "API URL", Type: "text", Required: true, Placeholder: "https://pve.example.com:8006/api2/json"},
			{Key: "token_id", Label: "API Token ID", Type: "text", Required: true, Placeholder: "root@pam!monitoring"},
			{Key: "token_secret", Label: "API Token Secret", Type: "password", Required: true},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		url, _ := config["url"].(string)
		tokenID, _ := config["token_id"].(string)
		tokenSecret, _ := config["token_secret"].(string)
		verifyTLS := true
		if v, ok := config["verify_tls"]; ok {
			if b, ok := v.(bool); ok {
				verifyTLS = b
			}
		}
		client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS})
		return &Connector{
			url:         strings.TrimSuffix(url, "/"),
			tokenID:     tokenID,
			tokenSecret: tokenSecret,
			client:      client,
		}, nil
	})
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on "vm" and "container" entities (see Fetch), exposed via
// GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"vm": {
		{Name: "status", Type: "string", Description: "Guest power state (running, stopped, ...)"},
		{Name: "firewall_enabled", Type: "boolean", Description: "Whether the per-guest firewall is enabled"},
		{Name: "onboot", Type: "boolean", Description: "Whether the guest starts automatically on host boot"},
		{Name: "agent_enabled", Type: "boolean", Description: "Whether the QEMU guest agent is enabled"},
		{Name: "protection", Type: "boolean", Description: "Whether removal/disk-wipe protection is enabled"},
		{Name: "template", Type: "boolean", Description: "Whether the guest is a template"},
		{Name: "os_type", Type: "string", Description: "Configured guest OS type"},
	},
	"container": {
		{Name: "status", Type: "string", Description: "Guest power state (running, stopped, ...)"},
		{Name: "firewall_enabled", Type: "boolean", Description: "Whether the per-guest firewall is enabled"},
		{Name: "onboot", Type: "boolean", Description: "Whether the guest starts automatically on host boot"},
		{Name: "agent_enabled", Type: "boolean", Description: "Whether the QEMU guest agent is enabled"},
		{Name: "protection", Type: "boolean", Description: "Whether removal/disk-wipe protection is enabled"},
		{Name: "template", Type: "boolean", Description: "Whether the guest is a template"},
		{Name: "os_type", Type: "string", Description: "Configured guest OS type"},
		{Name: "unprivileged", Type: "boolean", Description: "Whether the container runs unprivileged"},
	},
}

// Connector fetches data from a Proxmox VE API.
type Connector struct {
	url         string
	tokenID     string
	tokenSecret string
	client      *http.Client
}

// Name returns the connector display name.
func (p *Connector) Name() string { return "Proxmox VE" }

// Type returns the connector type.
func (p *Connector) Type() string { return typeName }

// Category returns the connector category.
func (p *Connector) Category() string { return "virtualization" }

// Validate tests the connection to the Proxmox API.
func (p *Connector) Validate(ctx context.Context, _ map[string]any) error {
	_, err := p.doRequest(ctx, "GET", "/nodes", nil)
	return err
}
