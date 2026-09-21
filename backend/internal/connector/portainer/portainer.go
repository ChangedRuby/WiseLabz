// Package portainer implements a Portainer connector, targeting the
// Portainer CE/BE v2 REST API (the same API that backs the web UI).
//
// Portainer fronts one or more Docker environments ("endpoints"), so the
// connector documents the environments themselves, the stacks deployed
// through Portainer, and — via Portainer's Docker proxy — the containers,
// volumes and networks of every Docker-capable environment. Authentication
// is an access token sent as "X-API-Key".
//
// Note: the HTTP plumbing below (guarded dialer, TLS config, status-code to
// connector-error mapping, timeout detection) is deliberately duplicated
// from the sibling connectors rather than shared — #265/#266 will rewrite
// the connector base, and this package stays self-contained until then.
// The same goes for the Docker payload shapes: connector/docker decodes the
// same endpoints, but against a different transport, so the conventions are
// reused rather than the code.
package portainer

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "portainer"

// Portainer API paths this connector reads. The per-environment Docker
// proxy paths are built from pathEndpoints (see dockerPath).
const (
	pathStatus       = "/api/status"
	pathSystemStatus = "/api/system/status"
	pathEndpoints    = "/api/endpoints"
	pathStacks       = "/api/stacks"
)

// Docker proxy suffixes, appended to /api/endpoints/{id}/docker.
const (
	dockerContainers = "/containers/json?all=true"
	dockerVolumes    = "/volumes"
	dockerNetworks   = "/networks"
)

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "containers_paas",
		Name:     "Portainer",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Portainer URL", Type: "text", Required: true, Placeholder: "https://portainer.example.com:9443", Description: "Base URL of the Portainer instance (the host serving /api)."},
			{Key: "api_key", Label: "API Access Token", Type: "password", Required: true, Description: "Portainer access token, sent as \"X-API-Key\". Create one under My account → Access tokens."},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, newConnector)
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on environment/stack/container/volume/network entities, exposed via
// GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"environment": {
		{Name: "type", Type: "string", Description: "Environment kind (docker, agent, azure, edge_agent, kubernetes, kubernetes_agent, kubernetes_edge)"},
		{Name: "status", Type: "string", Description: "Environment reachability as reported by Portainer (up, down)"},
		{Name: "url", Type: "string", Description: "Docker/Kubernetes endpoint URL Portainer connects to"},
		{Name: "publicUrl", Type: "string", Description: "Public URL published for the environment, if configured"},
		{Name: "tls", Type: "boolean", Description: "Whether Portainer talks to the environment over TLS"},
		{Name: "swarm", Type: "boolean", Description: "Whether the environment is a Docker Swarm cluster"},
		{Name: "dockerVersion", Type: "string", Description: "Docker engine version from the environment's latest snapshot"},
		{Name: "groupId", Type: "number", Description: "Portainer environment group the environment belongs to"},
		{Name: "tags", Type: "string_array", Description: "Portainer tag IDs assigned to the environment"},
	},
	"stack": {
		{Name: "type", Type: "string", Description: "Stack kind (swarm, compose)"},
		{Name: "status", Type: "string", Description: "Stack status reported by Portainer (active, inactive)"},
		{Name: "environment", Type: "string", Description: "Name of the environment the stack is deployed to"},
		{Name: "entryPoint", Type: "string", Description: "Compose file the stack is deployed from"},
	},
	"container": {
		{Name: "image", Type: "string", Description: "Image reference the container was created from"},
		{Name: "state", Type: "string", Description: "Container state (running, exited, paused, ...)"},
		{Name: "network_mode", Type: "string", Description: "Docker network mode (bridge, host, none, container:<id>, ...)"},
		{Name: "published_ports", Type: "string_array", Description: "Published host:container/protocol port mappings"},
		{Name: "environment", Type: "string", Description: "Name of the Portainer environment the container runs in"},
		{Name: "stack", Type: "string", Description: "Compose project the container belongs to, from its labels"},
	},
	"volume": {
		{Name: "driver", Type: "string", Description: "Volume driver (local, nfs, ...)"},
		{Name: "mountpoint", Type: "string", Description: "Host path the volume is mounted from"},
		{Name: "environment", Type: "string", Description: "Name of the Portainer environment the volume lives in"},
	},
	"network": {
		{Name: "driver", Type: "string", Description: "Network driver (bridge, overlay, macvlan, ...)"},
		{Name: "scope", Type: "string", Description: "Network scope (local, swarm, global)"},
		{Name: "internal", Type: "boolean", Description: "Whether the network is internal (no external connectivity)"},
		{Name: "environment", Type: "string", Description: "Name of the Portainer environment the network lives in"},
	},
}

// Connector fetches environments, stacks and Docker resources from a
// Portainer instance.
type Connector struct {
	url    string
	apiKey string
	client *http.Client
}

func newConnector(config map[string]any) (connector.Connector, error) {
	rawURL, _ := config["url"].(string)
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
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: !verifyTLS},
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Connector{url: strings.TrimSuffix(rawURL, "/"), apiKey: apiKey, client: client}, nil
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Portainer" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "containers_paas" }

// Validate checks the connector has a URL and a token, and that the token
// is accepted by the API. It probes /api/endpoints rather than /api/status
// because the latter is served unauthenticated and would accept any token.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	if c.url == "" {
		return fmt.Errorf("portainer url is required")
	}
	if c.apiKey == "" {
		return fmt.Errorf("portainer api_key is required")
	}
	_, err := c.doRequest(ctx, pathEndpoints)
	return err
}

// Fetch retrieves the Portainer environments and stacks plus, for every
// Docker-capable environment, its containers, volumes and networks.
// config may carry a "fields" selective-fetch hint naming a subset of
// {"environments","stacks","containers","volumes","networks"}. Each section
// degrades on its own: a failing endpoint (or a single unreachable
// environment) becomes a placeholder rather than failing the whole Fetch.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	metadata := map[string]string{"portainer_url": c.url}
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity

	if version, instanceID, err := c.fetchStatus(ctx); err == nil {
		putMetadata(metadata, "portainer_version", version)
		putMetadata(metadata, "instance_id", instanceID)
	}

	// The environment list is both a section of its own and the fan-out
	// list for the Docker proxy sections, so it is fetched once.
	envs, envErr := c.fetchEnvironments(ctx, fields)
	if connector.WantsField(fields, "environments") {
		if envErr != nil {
			sections = append(sections, unavailable("Environments", envErr))
		} else {
			content, ents := buildEnvironmentTable(envs)
			sections = append(sections, connector.SnapshotSection{Title: "Environments", Content: content})
			entities = append(entities, ents...)
			metadata["environment_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	if connector.WantsField(fields, "stacks") {
		raw, err := c.doRequest(ctx, pathStacks)
		if err != nil {
			sections = append(sections, unavailable("Stacks", err))
		} else {
			content, ents := buildStackTable(raw, environmentNames(envs))
			sections = append(sections, connector.SnapshotSection{Title: "Stacks", Content: content})
			entities = append(entities, ents...)
			metadata["stack_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	for _, proxy := range dockerSections {
		if !connector.WantsField(fields, proxy.field) {
			continue
		}
		if envErr != nil {
			sections = append(sections, unavailable(proxy.title, envErr))
			continue
		}
		section, ents := c.dockerSection(ctx, envs, proxy)
		sections = append(sections, section)
		entities = append(entities, ents...)
		metadata[proxy.field+"_count"] = fmt.Sprintf("%d", len(ents))
	}

	return &connector.ServiceSnapshot{
		ServiceName:  "Portainer",
		Type:         typeName,
		Sections:     sections,
		Entities:     entities,
		Dependencies: environmentDependencies(envs),
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

// dockerSectionSpec describes one section rendered from a Docker proxy call
// made against every Docker-capable environment.
type dockerSectionSpec struct {
	field  string // selective-fetch field name
	title  string // section title
	path   string // suffix under /api/endpoints/{id}/docker
	header string // Markdown table header, written once for all environments
	rows   func(envName string, raw []byte) (string, []connector.SnapshotEntity, error)
}

var dockerSections = []dockerSectionSpec{
	{field: "containers", title: "Containers", path: dockerContainers, header: containerHeader, rows: containerRows},
	{field: "volumes", title: "Volumes", path: dockerVolumes, header: volumeHeader, rows: volumeRows},
	{field: "networks", title: "Networks", path: dockerNetworks, header: networkHeader, rows: networkRows},
}

// dockerSection renders one proxy section across every Docker-capable
// environment. An environment that fails (unreachable agent, edge device
// offline, malformed payload) contributes a note instead of rows, so the
// remaining environments still make it into the snapshot.
func (c *Connector) dockerSection(ctx context.Context, envs []environment, spec dockerSectionSpec) (connector.SnapshotSection, []connector.SnapshotEntity) {
	var body strings.Builder
	var notes strings.Builder
	var entities []connector.SnapshotEntity
	reachable := 0

	for _, env := range envs {
		if !env.dockerCapable() {
			continue
		}
		reachable++
		raw, err := c.doRequest(ctx, dockerPath(env.ID, spec.path))
		if err != nil {
			_, _ = fmt.Fprintf(&notes, "_%s for %s unavailable: %s_\n\n", spec.title, cell(env.Name), err.Error())
			continue
		}
		rows, ents, err := spec.rows(env.Name, raw)
		if err != nil {
			_, _ = fmt.Fprintf(&notes, "_%s for %s unavailable: %s_\n\n", spec.title, cell(env.Name), connector.NewMalformedResponseError(err).Error())
			continue
		}
		body.WriteString(rows)
		entities = append(entities, ents...)
	}

	content := notes.String()
	switch {
	case body.Len() > 0:
		content += spec.header + body.String()
	case reachable == 0:
		content += "_No Docker environments to query_"
	case notes.Len() == 0:
		content += "_No " + strings.ToLower(spec.title) + " returned_"
	}
	return connector.SnapshotSection{Title: spec.title, Content: strings.TrimRight(content, "\n") + "\n"}, entities
}

// fetchEnvironments reads /api/endpoints when any section needs it, and
// reports nil, nil when no requested field does.
func (c *Connector) fetchEnvironments(ctx context.Context, fields []string) ([]environment, error) {
	needed := connector.WantsField(fields, "environments")
	for _, spec := range dockerSections {
		needed = needed || connector.WantsField(fields, spec.field)
	}
	if !needed {
		return nil, nil
	}
	raw, err := c.doRequest(ctx, pathEndpoints)
	if err != nil {
		return nil, err
	}
	var envs []environment
	if err := json.Unmarshal(raw, &envs); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	return envs, nil
}

// fetchStatus reads the Portainer version and instance ID. Portainer 2.18
// moved the endpoint to /api/system/status and kept /api/status as a
// deprecated alias, so the newer path is tried first.
func (c *Connector) fetchStatus(ctx context.Context) (version, instanceID string, err error) {
	raw, err := c.doRequest(ctx, pathSystemStatus)
	if err != nil {
		raw, err = c.doRequest(ctx, pathStatus)
	}
	if err != nil {
		return "", "", err
	}
	var status struct {
		Version    string `json:"Version"`
		InstanceID string `json:"InstanceID"`
	}
	if err := json.Unmarshal(raw, &status); err != nil {
		return "", "", connector.NewMalformedResponseError(err)
	}
	return status.Version, status.InstanceID, nil
}

// environmentDependencies reports each Portainer environment as a host this
// service depends on, deduplicated and sorted for a stable snapshot.
func environmentDependencies(envs []environment) []connector.ServiceDependency {
	seen := make(map[string]struct{}, len(envs))
	names := make([]string, 0, len(envs))
	for _, env := range envs {
		if env.Name == "" {
			continue
		}
		if _, dup := seen[env.Name]; dup {
			continue
		}
		seen[env.Name] = struct{}{}
		names = append(names, env.Name)
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	deps := make([]connector.ServiceDependency, 0, len(names))
	for _, name := range names {
		deps = append(deps, connector.ServiceDependency{Kind: "host", Name: name})
	}
	return deps
}

// dockerPath builds the Docker proxy path for one environment.
func dockerPath(envID int, suffix string) string {
	return fmt.Sprintf("%s/%d/docker%s", pathEndpoints, envID, suffix)
}

func unavailable(title string, err error) connector.SnapshotSection {
	return connector.SnapshotSection{Title: title, Content: "_" + title + " unavailable: " + err.Error() + "_"}
}

func putMetadata(metadata map[string]string, key, value string) {
	if value != "" {
		metadata[key] = value
	}
}

func (c *Connector) doRequest(ctx context.Context, path string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
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
// exceeded, covering both a canceled context and a net.Error timeout.
func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
