package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// Fetch retrieves engine info, containers, images, volumes, and networks.
// config may carry a "fields" selective-fetch hint naming a subset of
// {"containers","images","volumes","networks"} to skip the other calls.
// The System section always runs. Each section fetch failure is tolerated
// as a placeholder rather than failing the whole Fetch.
func (d *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)

	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	metadata := map[string]string{"docker_host": d.host}

	if raw, err := d.doRequest(ctx, "/info"); err != nil {
		sections = append(sections, connector.SnapshotSection{Title: "System", Content: "_System info unavailable: " + err.Error() + "_"})
	} else {
		var info struct {
			Name              string `json:"Name"`
			ServerVersion     string `json:"ServerVersion"`
			Containers        int    `json:"Containers"`
			ContainersRunning int    `json:"ContainersRunning"`
			Images            int    `json:"Images"`
			NCPU              int    `json:"NCPU"`
			MemTotal          int64  `json:"MemTotal"`
		}
		if err := json.Unmarshal(raw, &info); err != nil {
			sections = append(sections, connector.SnapshotSection{
				Title:   "System",
				Content: "_System info unavailable: " + connector.NewMalformedResponseError(err).Error() + "_",
			})
		} else {
			content := fmt.Sprintf("**Host**: %s\n**Engine Version**: %s\n**Containers**: %d (%d running)\n**Images**: %d\n**CPUs**: %d\n**Memory**: %d bytes\n",
				info.Name, info.ServerVersion, info.Containers, info.ContainersRunning, info.Images, info.NCPU, info.MemTotal)
			sections = append(sections, connector.SnapshotSection{Title: "System", Content: content})
			metadata["engine_version"] = info.ServerVersion
		}
	}

	if connector.WantsField(fields, "containers") {
		if raw, err := d.doRequest(ctx, "/containers/json?all=true"); err != nil {
			sections = append(sections, connector.SnapshotSection{Title: "Containers", Content: "_Containers unavailable: " + err.Error() + "_"})
		} else {
			content, ents := buildContainerTable(raw)
			d.enrichContainerAttributes(ctx, ents)
			sections = append(sections, connector.SnapshotSection{Title: "Containers", Content: content})
			entities = append(entities, ents...)
		}
	}
	if connector.WantsField(fields, "images") {
		sections = append(sections, d.fetchSection(ctx, "Images", "/images/json", buildImageTable))
	}
	if connector.WantsField(fields, "volumes") {
		sections = append(sections, d.fetchSection(ctx, "Volumes", "/volumes", buildVolumeTable))
	}
	if connector.WantsField(fields, "networks") {
		sections = append(sections, d.fetchSection(ctx, "Networks", "/networks", buildNetworkTable))
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Docker",
		Type:        typeName,
		Sections:    sections,
		Dependencies: []connector.ServiceDependency{
			{Kind: "host", Name: d.host},
		},
		Entities:  entities,
		Metadata:  metadata,
		FetchedAt: start,
	}, nil
}

// fetchSection runs a single endpoint fetch and renders it with build,
// tolerating failure as a placeholder section rather than failing Fetch.
func (d *Connector) fetchSection(ctx context.Context, title, path string, build func([]byte) string) connector.SnapshotSection {
	raw, err := d.doRequest(ctx, path)
	if err != nil {
		return connector.SnapshotSection{Title: title, Content: "_" + title + " unavailable: " + err.Error() + "_"}
	}
	return connector.SnapshotSection{Title: title, Content: build(raw)}
}

func buildContainerTable(raw []byte) (string, []connector.SnapshotEntity) {
	var containers []struct {
		ID     string   `json:"Id"`
		Names  []string `json:"Names"`
		Image  string   `json:"Image"`
		State  string   `json:"State"`
		Status string   `json:"Status"`
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
	if err := json.Unmarshal(raw, &containers); err != nil || len(containers) == 0 {
		return "_No containers returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Name | Image | State | Status |\n")
	b.WriteString("|------|-------|-------|--------|\n")
	var entities []connector.SnapshotEntity
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", name, c.Image, c.State, c.Status); err != nil {
			return "", nil
		}
		attrs := map[string]any{"image": c.Image}
		if c.HostConfig.NetworkMode != "" {
			attrs["network_mode"] = c.HostConfig.NetworkMode
		}
		var ports []string
		for _, p := range c.Ports {
			if p.PublicPort == 0 {
				continue
			}
			ports = append(ports, fmt.Sprintf("%d:%d/%s", p.PublicPort, p.PrivatePort, p.Type))
		}
		if len(ports) > 0 {
			attrs["published_ports"] = ports
		}
		ent := connector.SnapshotEntity{Kind: "container", Name: name, ExternalID: c.ID, Attributes: attrs}
		for _, net := range c.NetworkSettings.Networks {
			if net.IPAddress != "" {
				ent.IP = net.IPAddress
				break
			}
		}
		entities = append(entities, ent)
	}
	return b.String(), entities
}

// enrichContainerAttributes fetches the detailed GET /containers/{id}/json
// (inspect) response for each container and merges the privileged,
// restart_policy, user, and read_only_rootfs attributes into it — fields the
// list endpoint (/containers/json) doesn't return. It mutates ents in place.
// A container with no ID (e.g. an unrealistic/malformed list entry) or a
// failed/malformed inspect call is skipped so the rest of Fetch still
// succeeds; those specific attributes are simply omitted.
func (d *Connector) enrichContainerAttributes(ctx context.Context, ents []connector.SnapshotEntity) {
	for i := range ents {
		if ents[i].ExternalID == "" {
			continue
		}
		raw, err := d.doRequest(ctx, "/containers/"+ents[i].ExternalID+"/json")
		if err != nil {
			continue
		}
		var inspect struct {
			Config struct {
				User string `json:"User"`
			} `json:"Config"`
			HostConfig struct {
				Privileged     bool `json:"Privileged"`
				ReadonlyRootfs bool `json:"ReadonlyRootfs"`
				RestartPolicy  struct {
					Name string `json:"Name"`
				} `json:"RestartPolicy"`
			} `json:"HostConfig"`
		}
		if err := json.Unmarshal(raw, &inspect); err != nil {
			continue
		}
		if ents[i].Attributes == nil {
			ents[i].Attributes = map[string]any{}
		}
		ents[i].Attributes["privileged"] = inspect.HostConfig.Privileged
		ents[i].Attributes["read_only_rootfs"] = inspect.HostConfig.ReadonlyRootfs
		if inspect.HostConfig.RestartPolicy.Name != "" {
			ents[i].Attributes["restart_policy"] = inspect.HostConfig.RestartPolicy.Name
		}
		if inspect.Config.User != "" {
			ents[i].Attributes["user"] = inspect.Config.User
		}
	}
}

func buildImageTable(raw []byte) string {
	var images []struct {
		RepoTags []string `json:"RepoTags"`
		Size     int64    `json:"Size"`
	}
	if err := json.Unmarshal(raw, &images); err != nil || len(images) == 0 {
		return "_No images returned_"
	}
	var b strings.Builder
	b.WriteString("| Tags | Size (bytes) |\n")
	b.WriteString("|------|---------------|\n")
	for _, img := range images {
		tags := strings.Join(img.RepoTags, ", ")
		if tags == "" {
			tags = "<none>"
		}
		if _, err := fmt.Fprintf(&b, "| %s | %d |\n", tags, img.Size); err != nil {
			return ""
		}
	}
	return b.String()
}

func buildVolumeTable(raw []byte) string {
	var resp struct {
		Volumes []struct {
			Name       string `json:"Name"`
			Driver     string `json:"Driver"`
			Mountpoint string `json:"Mountpoint"`
		} `json:"Volumes"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Volumes) == 0 {
		return "_No volumes returned_"
	}
	var b strings.Builder
	b.WriteString("| Name | Driver | Mountpoint |\n")
	b.WriteString("|------|--------|------------|\n")
	for _, v := range resp.Volumes {
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s |\n", v.Name, v.Driver, v.Mountpoint); err != nil {
			return ""
		}
	}
	return b.String()
}

func buildNetworkTable(raw []byte) string {
	var networks []struct {
		Name   string `json:"Name"`
		Driver string `json:"Driver"`
		Scope  string `json:"Scope"`
	}
	if err := json.Unmarshal(raw, &networks); err != nil || len(networks) == 0 {
		return "_No networks returned_"
	}
	var b strings.Builder
	b.WriteString("| Name | Driver | Scope |\n")
	b.WriteString("|------|--------|-------|\n")
	for _, n := range networks {
		if _, err := fmt.Fprintf(&b, "| %s | %s | %s |\n", n.Name, n.Driver, n.Scope); err != nil {
			return ""
		}
	}
	return b.String()
}
