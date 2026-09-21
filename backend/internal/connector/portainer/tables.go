package portainer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// Portainer environment (endpoint) types, as returned in the "Type" field.
const (
	envTypeDocker          = 1
	envTypeAgent           = 2
	envTypeAzure           = 3
	envTypeEdgeAgent       = 4
	envTypeKubernetesLocal = 5
	envTypeKubernetesAgent = 6
	envTypeKubernetesEdge  = 7
)

// Portainer status enums: environments and stacks both use 1 for the
// healthy state and 2 for the unhealthy/inactive one.
const (
	statusUp   = 1
	statusDown = 2
)

// Portainer stack types.
const (
	stackTypeSwarm   = 1
	stackTypeCompose = 2
)

// Markdown table headers for the per-environment Docker proxy sections.
// They are emitted once per section, ahead of every environment's rows.
const (
	containerHeader = "| Environment | Name | Image | State | Status | Ports |\n" +
		"|-------------|------|-------|-------|--------|-------|\n"
	volumeHeader = "| Environment | Name | Driver | Mountpoint |\n" +
		"|-------------|------|--------|------------|\n"
	networkHeader = "| Environment | Name | Driver | Scope | Internal |\n" +
		"|-------------|------|--------|-------|----------|\n"
)

// composeProjectLabel is the label Docker Compose (and Portainer stacks)
// set on every container of a stack.
const composeProjectLabel = "com.docker.compose.project"

// environment is one Portainer environment ("endpoint"): a Docker engine,
// an agent, an Azure ACI account or a Kubernetes cluster.
type environment struct {
	ID        int    `json:"Id"`
	Name      string `json:"Name"`
	Type      int    `json:"Type"`
	URL       string `json:"URL"`
	PublicURL string `json:"PublicURL"`
	Status    int    `json:"Status"`
	GroupID   int    `json:"GroupId"`
	TLS       bool   `json:"TLS"`
	TagIDs    []int  `json:"TagIds"`
	Snapshots []struct {
		DockerVersion string `json:"DockerVersion"`
		Swarm         bool   `json:"Swarm"`
	} `json:"Snapshots"`
}

// dockerCapable reports whether the environment can be queried through
// Portainer's Docker proxy (/api/endpoints/{id}/docker/...). Azure ACI and
// Kubernetes environments speak different APIs and are skipped.
func (e environment) dockerCapable() bool {
	switch e.Type {
	case envTypeDocker, envTypeAgent, envTypeEdgeAgent:
		return e.Status != statusDown
	default:
		return false
	}
}

// dockerVersion returns the engine version from the environment's latest
// snapshot, or "" when Portainer has not snapshotted it yet.
func (e environment) dockerVersion() string {
	if len(e.Snapshots) == 0 {
		return ""
	}
	return e.Snapshots[0].DockerVersion
}

// swarm reports whether the environment's latest snapshot saw a Swarm.
func (e environment) swarm() bool {
	return len(e.Snapshots) > 0 && e.Snapshots[0].Swarm
}

// cell escapes a value for use inside a Markdown table cell; container
// commands, image references and mountpoints can all contain '|'.
func cell(s string) string {
	if s == "" {
		return "—"
	}
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", `\|`)
}

// environmentTypeName maps Portainer's numeric environment type to a stable
// string for attributes and tables.
func environmentTypeName(t int) string {
	switch t {
	case envTypeDocker:
		return "docker"
	case envTypeAgent:
		return "agent"
	case envTypeAzure:
		return "azure"
	case envTypeEdgeAgent:
		return "edge_agent"
	case envTypeKubernetesLocal:
		return "kubernetes"
	case envTypeKubernetesAgent:
		return "kubernetes_agent"
	case envTypeKubernetesEdge:
		return "kubernetes_edge"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}

// statusName maps a Portainer environment status to up/down.
func statusName(s int) string {
	switch s {
	case statusUp:
		return "up"
	case statusDown:
		return "down"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// stackTypeName maps a Portainer stack type to swarm/compose.
func stackTypeName(t int) string {
	switch t {
	case stackTypeSwarm:
		return "swarm"
	case stackTypeCompose:
		return "compose"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}

// stackStatusName maps a Portainer stack status to active/inactive.
func stackStatusName(s int) string {
	switch s {
	case statusUp:
		return "active"
	case statusDown:
		return "inactive"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// environmentNames indexes environment names by ID so stacks can name the
// environment they are deployed to.
func environmentNames(envs []environment) map[int]string {
	names := make(map[int]string, len(envs))
	for _, env := range envs {
		names[env.ID] = env.Name
	}
	return names
}

// buildEnvironmentTable renders the decoded /api/endpoints payload.
func buildEnvironmentTable(envs []environment) (string, []connector.SnapshotEntity) {
	if len(envs) == 0 {
		return "_No environments returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Environment | Type | Status | URL | Docker | Swarm |\n")
	b.WriteString("|-------------|------|--------|-----|--------|-------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(envs))
	for _, env := range envs {
		typeName := environmentTypeName(env.Type)
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %t |\n",
			cell(env.Name), typeName, statusName(env.Status), cell(env.URL),
			cell(env.dockerVersion()), env.swarm())

		attrs := map[string]any{
			"type":    typeName,
			"status":  statusName(env.Status),
			"tls":     env.TLS,
			"swarm":   env.swarm(),
			"groupId": env.GroupID,
		}
		putString(attrs, "url", env.URL)
		putString(attrs, "publicUrl", env.PublicURL)
		putString(attrs, "dockerVersion", env.dockerVersion())
		putStrings(attrs, "tags", tagNames(env.TagIDs))
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "environment",
			Name:       env.Name,
			ExternalID: fmt.Sprintf("%d", env.ID),
			Hostname:   hostFromURL(env.URL),
			Attributes: attrs,
		})
	}
	return b.String(), entities
}

// tagNames renders Portainer's numeric tag IDs as strings; the tag names
// themselves live behind a separate endpoint this connector doesn't read.
func tagNames(ids []int) []string {
	if len(ids) == 0 {
		return nil
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, fmt.Sprintf("%d", id))
	}
	return out
}

// hostFromURL extracts the host part of an environment URL
// ("tcp://10.0.0.5:2376" -> "10.0.0.5") so environments cross-link with
// hosts documented by other connectors. Socket URLs yield "".
func hostFromURL(raw string) string {
	if raw == "" || strings.HasPrefix(raw, "unix://") || strings.HasPrefix(raw, "npipe://") {
		return ""
	}
	host := raw
	if i := strings.Index(host, "://"); i >= 0 {
		host = host[i+3:]
	}
	host = strings.TrimSuffix(host, "/")
	if i := strings.IndexByte(host, '/'); i >= 0 {
		host = host[:i]
	}
	if i := strings.LastIndexByte(host, ':'); i > 0 {
		host = host[:i]
	}
	return host
}

// buildStackTable renders /api/stacks, naming each stack's environment.
func buildStackTable(raw []byte, envNames map[int]string) (string, []connector.SnapshotEntity) {
	var stacks []struct {
		ID         int    `json:"Id"`
		Name       string `json:"Name"`
		Type       int    `json:"Type"`
		EndpointID int    `json:"EndpointId"`
		Status     int    `json:"Status"`
		EntryPoint string `json:"EntryPoint"`
	}
	if err := json.Unmarshal(raw, &stacks); err != nil {
		return malformed("Stacks", err), nil
	}
	if len(stacks) == 0 {
		return "_No stacks returned_", nil
	}

	var b strings.Builder
	b.WriteString("| Stack | Environment | Type | Status | Entry Point |\n")
	b.WriteString("|-------|-------------|------|--------|-------------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(stacks))
	for _, s := range stacks {
		envName := envNames[s.EndpointID]
		if envName == "" {
			envName = fmt.Sprintf("endpoint %d", s.EndpointID)
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			cell(s.Name), cell(envName), stackTypeName(s.Type), stackStatusName(s.Status), cell(s.EntryPoint))

		attrs := map[string]any{
			"type":   stackTypeName(s.Type),
			"status": stackStatusName(s.Status),
		}
		putString(attrs, "environment", envName)
		putString(attrs, "entryPoint", s.EntryPoint)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "stack", Name: s.Name, ExternalID: fmt.Sprintf("%d", s.ID), Attributes: attrs,
		})
	}
	return b.String(), entities
}

// containerRows renders one environment's /containers/json payload as table
// rows (without a header) plus its container entities.
func containerRows(envName string, raw []byte) (string, []connector.SnapshotEntity, error) {
	var containers []struct {
		ID     string            `json:"Id"`
		Names  []string          `json:"Names"`
		Image  string            `json:"Image"`
		State  string            `json:"State"`
		Status string            `json:"Status"`
		Labels map[string]string `json:"Labels"`
		Ports  []struct {
			PrivatePort int    `json:"PrivatePort"`
			PublicPort  int    `json:"PublicPort"`
			Type        string `json:"Type"`
		} `json:"Ports"`
		HostConfig struct {
			NetworkMode string `json:"NetworkMode"`
		} `json:"HostConfig"`
		NetworkSettings struct {
			Networks map[string]struct {
				IPAddress string `json:"IPAddress"`
			} `json:"Networks"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(raw, &containers); err != nil {
		return "", nil, err
	}

	var b strings.Builder
	entities := make([]connector.SnapshotEntity, 0, len(containers))
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		var ports []string
		for _, p := range c.Ports {
			if p.PublicPort == 0 {
				continue
			}
			ports = append(ports, fmt.Sprintf("%d:%d/%s", p.PublicPort, p.PrivatePort, p.Type))
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n",
			cell(envName), cell(name), cell(c.Image), cell(c.State), cell(c.Status), cell(strings.Join(ports, ", ")))

		attrs := map[string]any{}
		putString(attrs, "image", c.Image)
		putString(attrs, "state", c.State)
		putString(attrs, "network_mode", c.HostConfig.NetworkMode)
		putString(attrs, "environment", envName)
		putString(attrs, "stack", c.Labels[composeProjectLabel])
		putStrings(attrs, "published_ports", ports)

		ent := connector.SnapshotEntity{Kind: "container", Name: name, ExternalID: c.ID, Attributes: attrs}
		for _, net := range c.NetworkSettings.Networks {
			if net.IPAddress != "" {
				ent.IP = net.IPAddress
				break
			}
		}
		entities = append(entities, ent)
	}
	return b.String(), entities, nil
}

// volumeRows renders one environment's /volumes payload.
func volumeRows(envName string, raw []byte) (string, []connector.SnapshotEntity, error) {
	var resp struct {
		Volumes []struct {
			Name       string `json:"Name"`
			Driver     string `json:"Driver"`
			Mountpoint string `json:"Mountpoint"`
		} `json:"Volumes"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", nil, err
	}

	var b strings.Builder
	entities := make([]connector.SnapshotEntity, 0, len(resp.Volumes))
	for _, v := range resp.Volumes {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s |\n",
			cell(envName), cell(v.Name), cell(v.Driver), cell(v.Mountpoint))

		attrs := map[string]any{}
		putString(attrs, "driver", v.Driver)
		putString(attrs, "mountpoint", v.Mountpoint)
		putString(attrs, "environment", envName)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "volume", Name: v.Name, ExternalID: scopedID(envName, v.Name), Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

// networkRows renders one environment's /networks payload.
func networkRows(envName string, raw []byte) (string, []connector.SnapshotEntity, error) {
	var networks []struct {
		ID       string `json:"Id"`
		Name     string `json:"Name"`
		Driver   string `json:"Driver"`
		Scope    string `json:"Scope"`
		Internal bool   `json:"Internal"`
	}
	if err := json.Unmarshal(raw, &networks); err != nil {
		return "", nil, err
	}

	var b strings.Builder
	entities := make([]connector.SnapshotEntity, 0, len(networks))
	for _, n := range networks {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %t |\n",
			cell(envName), cell(n.Name), cell(n.Driver), cell(n.Scope), n.Internal)

		attrs := map[string]any{"internal": n.Internal}
		putString(attrs, "driver", n.Driver)
		putString(attrs, "scope", n.Scope)
		putString(attrs, "environment", envName)
		externalID := n.ID
		if externalID == "" {
			externalID = scopedID(envName, n.Name)
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind: "network", Name: n.Name, ExternalID: externalID, Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

// scopedID qualifies a resource name with its environment, since volume and
// network names are only unique within one Docker engine.
func scopedID(envName, name string) string { return envName + "/" + name }

func malformed(title string, err error) string {
	return "_" + title + " unavailable: " + connector.NewMalformedResponseError(err).Error() + "_"
}

func putString(attrs map[string]any, key, value string) {
	if value != "" {
		attrs[key] = value
	}
}

func putStrings(attrs map[string]any, key string, values []string) {
	if len(values) > 0 {
		attrs[key] = values
	}
}
