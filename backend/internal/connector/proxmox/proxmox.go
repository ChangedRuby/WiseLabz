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
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifyTLS},
			},
		}
		return &Connector{
			url:         strings.TrimSuffix(url, "/"),
			tokenID:     tokenID,
			tokenSecret: tokenSecret,
			client:      client,
		}, nil
	})
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
					if wantEntities && vm.Status == "running" {
						ent.IP = p.fetchQemuIP(ctx, node.Node, vm.VMID)
					}
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
					if wantEntities && ct.Status == "running" {
						ent.IP = p.fetchLxcIP(ctx, node.Node, ct.VMID)
					}
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

func (p *Connector) doRequest(ctx context.Context, method, path string, _ io.Reader) ([]byte, error) {
	url := p.url + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
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
