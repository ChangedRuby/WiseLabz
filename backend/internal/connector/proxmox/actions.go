package proxmox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// Restart reboots the VM or container identified by entityRef (a Proxmox
// VMID). It resolves the owning node and guest type via /cluster/resources
// since a bare VMID alone doesn't say which.
func (p *Connector) Restart(ctx context.Context, _ map[string]any, entityRef string) error {
	if entityRef == "" {
		return fmt.Errorf("proxmox restart requires a target VMID")
	}
	raw, err := p.doRequest(ctx, "GET", "/cluster/resources?type=vm", nil)
	if err != nil {
		return fmt.Errorf("resolve VM node: %w", err)
	}
	var resp struct {
		Data []struct {
			VMID int    `json:"vmid"`
			Node string `json:"node"`
			Type string `json:"type"` // "qemu" or "lxc"
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return connector.NewMalformedResponseError(fmt.Errorf("decode cluster resources: %w", err))
	}
	for _, res := range resp.Data {
		if fmt.Sprintf("%d", res.VMID) != entityRef {
			continue
		}
		_, err := p.doRequest(ctx, "POST", fmt.Sprintf("/nodes/%s/%s/%d/status/reboot", res.Node, res.Type, res.VMID), nil)
		return err
	}
	return fmt.Errorf("VMID %s not found", entityRef)
}

// Start starts the VM or container identified by entityRef (a Proxmox
// VMID). Idempotent-safe: Proxmox no-ops (200 with a completed task) a
// start on an already-running guest.
func (p *Connector) Start(ctx context.Context, _ map[string]any, entityRef string) error {
	return p.vmStatusAction(ctx, entityRef, "start")
}

// Stop stops the VM or container identified by entityRef (a Proxmox VMID).
func (p *Connector) Stop(ctx context.Context, _ map[string]any, entityRef string) error {
	return p.vmStatusAction(ctx, entityRef, "stop")
}

// vmStatusAction resolves entityRef's owning node/guest type via
// /cluster/resources (same lookup Restart uses) and POSTs the given
// status action ("start", "stop", "reboot").
func (p *Connector) vmStatusAction(ctx context.Context, entityRef, action string) error {
	if entityRef == "" {
		return fmt.Errorf("proxmox %s requires a target VMID", action)
	}
	raw, err := p.doRequest(ctx, "GET", "/cluster/resources?type=vm", nil)
	if err != nil {
		return fmt.Errorf("resolve VM node: %w", err)
	}
	var resp struct {
		Data []struct {
			VMID int    `json:"vmid"`
			Node string `json:"node"`
			Type string `json:"type"` // "qemu" or "lxc"
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return connector.NewMalformedResponseError(fmt.Errorf("decode cluster resources: %w", err))
	}
	for _, res := range resp.Data {
		if fmt.Sprintf("%d", res.VMID) != entityRef {
			continue
		}
		_, err := p.doRequest(ctx, "POST", fmt.Sprintf("/nodes/%s/%s/%d/status/%s", res.Node, res.Type, res.VMID, action), nil)
		return err
	}
	return fmt.Errorf("VMID %s not found", entityRef)
}

// WritableFields lists the config-push-eligible VM/container fields.
func (p *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "memory", Label: "Memory (MB)", Type: "number", EntityScope: true},
		{Key: "cores", Label: "CPU Cores", Type: "number", EntityScope: true},
	}
}

// ConfigPush writes a single whitelisted field (memory or cores) to the VM
// or container identified by entityRef via the Proxmox config endpoint.
func (p *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("proxmox config-push requires a target VMID")
	}
	raw, err := p.doRequest(ctx, "GET", "/cluster/resources?type=vm", nil)
	if err != nil {
		return fmt.Errorf("resolve VM node: %w", err)
	}
	var resp struct {
		Data []struct {
			VMID int    `json:"vmid"`
			Node string `json:"node"`
			Type string `json:"type"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return connector.NewMalformedResponseError(fmt.Errorf("decode cluster resources: %w", err))
	}
	for _, res := range resp.Data {
		if fmt.Sprintf("%d", res.VMID) != entityRef {
			continue
		}
		form := url.Values{fieldKey: {fmt.Sprintf("%v", value)}}
		_, err = p.doRequest(ctx, "PUT", fmt.Sprintf("/nodes/%s/%s/%d/config", res.Node, res.Type, res.VMID), strings.NewReader(form.Encode()))
		return err
	}
	return fmt.Errorf("VMID %s not found", entityRef)
}
