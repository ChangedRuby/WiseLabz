// Package proxmox implements a real Proxmox VE API connector.
package proxmox

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

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
		dialer := connector.GuardedDialer(30 * time.Second)
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DialContext:     dialer.DialContext,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyTLS},
			},
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
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

// Fetch retrieves the current state of nodes, VMs, containers, and storage.
// config may carry a "fields" selective-fetch hint (see
// connector.RequestedFields) naming a subset of {"vms","containers","storage"}
// to skip the other upstream calls — e.g. a dashboard quick-check that only
// needs VM state doesn't need to also hit /storage on every node.
func (p *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	wantVMs := connector.WantsField(fields, "vms")
	wantContainers := connector.WantsField(fields, "containers")
	wantStorage := connector.WantsField(fields, "storage")
	wantEntities := connector.WantsField(fields, "entities")

	// Fetch nodes
	nodesRaw, err := p.doRequest(ctx, "GET", "/nodes", nil)
	if err != nil {
		return nil, fmt.Errorf("fetch nodes: %w", err)
	}

	var nodesResponse struct {
		Data []struct {
			Node   string  `json:"node"`
			Status string  `json:"status"`
			Uptime int64   `json:"uptime"`
			CPU    float64 `json:"cpu"`
			Memory struct {
				Used  int64 `json:"used"`
				Total int64 `json:"total"`
			} `json:"mem"`
		} `json:"data"`
	}
	if err := json.Unmarshal(nodesRaw, &nodesResponse); err != nil {
		return nil, connector.NewMalformedResponseError(fmt.Errorf("decode nodes: %w", err))
	}

	var sections []connector.SnapshotSection
	var dependencies []connector.ServiceDependency
	var entities []connector.SnapshotEntity
	metadata := map[string]string{
		"node_count": fmt.Sprintf("%d", len(nodesResponse.Data)),
	}

	// For each node, fetch VMs, containers, and storage
	totalVMs := 0
	totalCTs := 0
	totalStorage := 0

	for _, node := range nodesResponse.Data {
		nodeSection := fmt.Sprintf("## Node: %s\n\n", node.Node)
		nodeSection += fmt.Sprintf("- **Status**: %s\n", node.Status)
		nodeSection += fmt.Sprintf("- **Uptime**: %d seconds\n", node.Uptime)
		nodeSection += fmt.Sprintf("- **CPU**: %.2f%%\n", node.CPU*100)
		nodeSection += fmt.Sprintf("- **Memory**: %d / %d bytes\n\n", node.Memory.Used, node.Memory.Total)

		dependencies = append(dependencies, connector.ServiceDependency{Kind: "host", Name: node.Node})

		// Fetch VMs
		if wantVMs {
			vmsRaw, err := p.doRequest(ctx, "GET", "/nodes/"+node.Node+"/qemu", nil)
			if err != nil {
				sections = append(sections, connector.SnapshotSection{
					Title:   node.Node,
					Content: nodeSection + "_VM fetch error: " + err.Error() + "_",
				})
				continue
			}
			var vmsResponse struct {
				Data []struct {
					VMID   int    `json:"vmid"`
					Name   string `json:"name"`
					Status string `json:"status"`
					CPU    int    `json:"cpus"`
					Memory int64  `json:"mem"`
					Uptime int64  `json:"uptime"`
				} `json:"data"`
			}
			if err := json.Unmarshal(vmsRaw, &vmsResponse); err != nil {
				sections = append(sections, connector.SnapshotSection{
					Title:   node.Node,
					Content: nodeSection + "_VM decode error: " + err.Error() + "_",
				})
				continue
			}

			if len(vmsResponse.Data) > 0 {
				nodeSection += "### Virtual Machines\n\n"
				nodeSection += "| VMID | Name | Status | CPUs | Memory (MB) | Uptime |\n"
				nodeSection += "|------|------|--------|------|-------------|--------|\n"
				for _, vm := range vmsResponse.Data {
					nodeSection += fmt.Sprintf("| %d | %s | %s | %d | %d | %d |\n",
						vm.VMID, vm.Name, vm.Status, vm.CPU, vm.Memory, vm.Uptime)
					ent := connector.SnapshotEntity{
						Kind:       "vm",
						Name:       vm.Name,
						ExternalID: fmt.Sprintf("%d", vm.VMID),
					}
					attrs := map[string]any{"status": vm.Status}
					if wantEntities {
						if vm.Status == "running" {
							ent.IP = p.fetchQemuIP(ctx, node.Node, vm.VMID)
						}
						if cfg, ok := p.fetchQemuConfig(ctx, node.Node, vm.VMID); ok {
							attrs["onboot"] = cfg.Onboot != 0
							attrs["protection"] = cfg.Protection != 0
							attrs["template"] = cfg.Template != 0
							attrs["agent_enabled"] = agentEnabled(cfg.Agent)
							if cfg.OSType != "" {
								attrs["os_type"] = cfg.OSType
							}
						}
						if enabled, ok := p.fetchFirewallEnabled(ctx, "qemu", node.Node, vm.VMID); ok {
							attrs["firewall_enabled"] = enabled
						}
					}
					ent.Attributes = attrs
					entities = append(entities, ent)
				}
				nodeSection += "\n"
				totalVMs += len(vmsResponse.Data)
			}
		}

		// Fetch containers
		if wantContainers {
			ctsRaw, err := p.doRequest(ctx, "GET", "/nodes/"+node.Node+"/lxc", nil)
			if err != nil {
				sections = append(sections, connector.SnapshotSection{
					Title:   node.Node,
					Content: nodeSection + "_CT fetch error: " + err.Error() + "_",
				})
				continue
			}
			var ctsResponse struct {
				Data []struct {
					VMID   int    `json:"vmid"`
					Name   string `json:"name"`
					Status string `json:"status"`
					CPU    int    `json:"cpus"`
					Memory int64  `json:"mem"`
					Uptime int64  `json:"uptime"`
				} `json:"data"`
			}
			if err := json.Unmarshal(ctsRaw, &ctsResponse); err != nil {
				sections = append(sections, connector.SnapshotSection{
					Title:   node.Node,
					Content: nodeSection + "_CT decode error: " + err.Error() + "_",
				})
				continue
			}

			if len(ctsResponse.Data) > 0 {
				nodeSection += "### Containers\n\n"
				nodeSection += "| VMID | Name | Status | CPUs | Memory (MB) | Uptime |\n"
				nodeSection += "|------|------|--------|------|-------------|--------|\n"
				for _, ct := range ctsResponse.Data {
					nodeSection += fmt.Sprintf("| %d | %s | %s | %d | %d | %d |\n",
						ct.VMID, ct.Name, ct.Status, ct.CPU, ct.Memory, ct.Uptime)
					ent := connector.SnapshotEntity{
						Kind:       "container",
						Name:       ct.Name,
						ExternalID: fmt.Sprintf("%d", ct.VMID),
					}
					attrs := map[string]any{"status": ct.Status}
					if wantEntities {
						if ct.Status == "running" {
							ent.IP = p.fetchLxcIP(ctx, node.Node, ct.VMID)
						}
						if cfg, ok := p.fetchLxcConfig(ctx, node.Node, ct.VMID); ok {
							attrs["onboot"] = cfg.Onboot != 0
							attrs["protection"] = cfg.Protection != 0
							attrs["template"] = cfg.Template != 0
							attrs["agent_enabled"] = agentEnabled(cfg.Agent)
							attrs["unprivileged"] = cfg.Unprivileged != 0
							if cfg.OSType != "" {
								attrs["os_type"] = cfg.OSType
							}
						}
						if enabled, ok := p.fetchFirewallEnabled(ctx, "lxc", node.Node, ct.VMID); ok {
							attrs["firewall_enabled"] = enabled
						}
					}
					ent.Attributes = attrs
					entities = append(entities, ent)
				}
				nodeSection += "\n"
				totalCTs += len(ctsResponse.Data)
			}
		}

		// Fetch storage
		if wantStorage {
			storageRaw, err := p.doRequest(ctx, "GET", "/nodes/"+node.Node+"/storage", nil)
			if err != nil {
				sections = append(sections, connector.SnapshotSection{
					Title:   node.Node,
					Content: nodeSection + "_Storage fetch error: " + err.Error() + "_",
				})
				continue
			}
			var storageResponse struct {
				Data []struct {
					Storage string `json:"storage"`
					Type    string `json:"type"`
					Used    int64  `json:"used"`
					Total   int64  `json:"total"`
					Avail   int64  `json:"avail"`
				} `json:"data"`
			}
			if err := json.Unmarshal(storageRaw, &storageResponse); err != nil {
				sections = append(sections, connector.SnapshotSection{
					Title:   node.Node,
					Content: nodeSection + "_Storage decode error: " + err.Error() + "_",
				})
				continue
			}

			if len(storageResponse.Data) > 0 {
				nodeSection += "### Storage\n\n"
				nodeSection += "| Storage | Type | Used (bytes) | Total (bytes) | Avail (bytes) |\n"
				nodeSection += "|---------|------|---------------|----------------|----------------|\n"
				for _, st := range storageResponse.Data {
					nodeSection += fmt.Sprintf("| %s | %s | %d | %d | %d |\n",
						st.Storage, st.Type, st.Used, st.Total, st.Avail)
					dependencies = append(dependencies, connector.ServiceDependency{Kind: "storage", Name: st.Storage})
				}
				nodeSection += "\n"
				totalStorage += len(storageResponse.Data)
			}
		}

		sections = append(sections, connector.SnapshotSection{
			Title:   node.Node,
			Content: nodeSection,
		})
	}

	metadata["total_vms"] = fmt.Sprintf("%d", totalVMs)
	metadata["total_cts"] = fmt.Sprintf("%d", totalCTs)
	metadata["total_storage"] = fmt.Sprintf("%d", totalStorage)
	metadata["proxmox_url"] = p.url

	return &connector.ServiceSnapshot{
		ServiceName:  "Proxmox VE",
		Type:         typeName,
		Sections:     sections,
		Dependencies: dependencies,
		Entities:     entities,
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

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

// fetchQemuIP returns the first non-loopback IPv4 address reported by the
// QEMU guest agent, or "" if the agent isn't installed/running (most labs
// won't have it on every VM) or reports nothing usable. Soft-fails: any
// error here is not a Fetch failure.
func (p *Connector) fetchQemuIP(ctx context.Context, node string, vmid int) string {
	raw, err := p.doRequest(ctx, "GET", fmt.Sprintf("/nodes/%s/qemu/%d/agent/network-get-interfaces", node, vmid), nil)
	if err != nil {
		return ""
	}
	var resp struct {
		Data struct {
			Result []struct {
				Name        string `json:"name"`
				IPAddresses []struct {
					IPAddress     string `json:"ip-address"`
					IPAddressType string `json:"ip-address-type"`
				} `json:"ip-addresses"`
			} `json:"result"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return ""
	}
	for _, iface := range resp.Data.Result {
		if iface.Name == "lo" {
			continue
		}
		for _, addr := range iface.IPAddresses {
			if addr.IPAddressType == "ipv4" && addr.IPAddress != "" {
				return addr.IPAddress
			}
		}
	}
	return ""
}

// fetchLxcIP returns the first non-loopback IPv4 address reported for the
// container. Unlike the QEMU path this needs no guest agent, but still
// soft-fails since older Proxmox versions lack this endpoint.
func (p *Connector) fetchLxcIP(ctx context.Context, node string, vmid int) string {
	raw, err := p.doRequest(ctx, "GET", fmt.Sprintf("/nodes/%s/lxc/%d/interfaces", node, vmid), nil)
	if err != nil {
		return ""
	}
	var resp struct {
		Data []struct {
			Name  string `json:"name"`
			Inet  string `json:"inet"`
			Inet6 string `json:"inet6"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return ""
	}
	for _, iface := range resp.Data {
		if iface.Name == "lo" || iface.Inet == "" {
			continue
		}
		ip, _, err := net.ParseCIDR(iface.Inet)
		if err == nil {
			return ip.String()
		}
		return iface.Inet
	}
	return ""
}

// qemuConfig holds the subset of VM /config fields we surface as
// Attributes. Fields absent from the Proxmox response decode to their zero
// value, which is the correct "disabled"/"unset" reading for each of these
// flags.
type qemuConfig struct {
	Onboot     int    `json:"onboot"`
	Protection int    `json:"protection"`
	Agent      string `json:"agent"`
	Template   int    `json:"template"`
	OSType     string `json:"ostype"`
}

// fetchQemuConfig fetches a VM's /config and returns the fields relevant to
// Attributes, or ok=false if the request/decode failed (soft-fail: any
// error here is not a Fetch failure, the caller just omits the attribute).
func (p *Connector) fetchQemuConfig(ctx context.Context, node string, vmid int) (qemuConfig, bool) {
	raw, err := p.doRequest(ctx, "GET", fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmid), nil)
	if err != nil {
		return qemuConfig{}, false
	}
	var resp struct {
		Data qemuConfig `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return qemuConfig{}, false
	}
	return resp.Data, true
}

// lxcConfig holds the subset of container /config fields we surface as
// Attributes.
type lxcConfig struct {
	Onboot       int    `json:"onboot"`
	Protection   int    `json:"protection"`
	Template     int    `json:"template"`
	Agent        string `json:"agent"`
	Unprivileged int    `json:"unprivileged"`
	OSType       string `json:"ostype"`
}

// fetchLxcConfig fetches a container's /config and returns the fields
// relevant to Attributes, or ok=false on request/decode failure.
func (p *Connector) fetchLxcConfig(ctx context.Context, node string, vmid int) (lxcConfig, bool) {
	raw, err := p.doRequest(ctx, "GET", fmt.Sprintf("/nodes/%s/lxc/%d/config", node, vmid), nil)
	if err != nil {
		return lxcConfig{}, false
	}
	var resp struct {
		Data lxcConfig `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return lxcConfig{}, false
	}
	return resp.Data, true
}

// fetchFirewallEnabled fetches a guest's firewall/options and reports
// whether the per-guest firewall is enabled. guestType is "qemu" or "lxc".
// Soft-fails to ok=false on any request/decode error.
func (p *Connector) fetchFirewallEnabled(ctx context.Context, guestType, node string, vmid int) (bool, bool) {
	raw, err := p.doRequest(ctx, "GET", fmt.Sprintf("/nodes/%s/%s/%d/firewall/options", node, guestType, vmid), nil)
	if err != nil {
		return false, false
	}
	var resp struct {
		Data struct {
			Enable int `json:"enable"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return false, false
	}
	return resp.Data.Enable != 0, true
}

// agentEnabled reports whether a QEMU config's "agent" field indicates the
// guest agent is enabled. Proxmox stores it either as a bare "1"/"0" or a
// comma-separated option string like "enabled=1,fstrim_cloned_disks=1", so
// presence of a leading "1" is treated as enabled.
func agentEnabled(agent string) bool {
	if agent == "" {
		return false
	}
	first := strings.SplitN(agent, ",", 2)[0]
	first = strings.TrimPrefix(first, "enabled=")
	return strings.TrimSpace(first) == "1"
}

func (p *Connector) doRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	url := p.url + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Authorization", "PVEAPIToken="+p.tokenID+"="+p.tokenSecret)
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
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

// isTimeout reports whether err represents a request deadline being
// exceeded, covering both a canceled context and a net.Error timeout
// (e.g. a dial or read timing out on the underlying transport).
func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
