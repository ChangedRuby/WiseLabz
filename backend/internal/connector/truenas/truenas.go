// Package truenas implements a TrueNAS connector, targeting the TrueNAS
// SCALE/CORE REST API v2.0 (the API behind the web UI, authenticated with
// an API key issued from Credentials -> Local Users -> API Keys).
//
// The snapshot deliberately carries configuration rather than telemetry:
// pool and dataset usage, uptime, load averages and disk temperatures all
// move on every poll and would make an otherwise unchanged appliance look
// different on each sync, so they are left out.
//
// Note: the HTTP plumbing below (guarded dialer, TLS config, status-code to
// connector-error mapping, timeout detection) is deliberately duplicated
// from the sibling connectors rather than shared — #265/#266 will rewrite
// the connector base, and this package stays self-contained until then.
package truenas

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const (
	typeName = "truenas"
	// category is the closest fit among the categories the connectors table
	// allows (virtualization|containers_paas|networking|dns); there is no
	// storage category and adding one would need a migration.
	category = "virtualization"
)

// TrueNAS API v2.0 paths this connector reads.
const (
	pathSystemInfo      = "/api/v2.0/system/info"
	pathPools           = "/api/v2.0/pool"
	pathDatasets        = "/api/v2.0/pool/dataset"
	pathDisks           = "/api/v2.0/disk"
	pathSMBShares       = "/api/v2.0/sharing/smb"
	pathNFSShares       = "/api/v2.0/sharing/nfs"
	pathServices        = "/api/v2.0/service"
	pathInterfaces      = "/api/v2.0/interface"
	pathSnapshotTasks   = "/api/v2.0/pool/snapshottask"
	pathReplicationTask = "/api/v2.0/replication"
)

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: category,
		Name:     "TrueNAS",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "TrueNAS URL", Type: "text", Required: true, Placeholder: "https://truenas.example.com", Description: "Base URL of the TrueNAS web UI (the host serving /api/v2.0)."},
			{Key: "api_key", Label: "API Key", Type: "password", Required: true, Description: "API key from Credentials → Local Users → API Keys, sent as \"Authorization: Bearer <key>\"."},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, newConnector)
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on the entities it emits, exposed via GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"pool": {
		{Name: "status", Type: "string", Description: "Pool status reported by ZFS (ONLINE, DEGRADED, FAULTED, ...)"},
		{Name: "healthy", Type: "boolean", Description: "Whether TrueNAS considers the pool healthy"},
		{Name: "encrypted", Type: "boolean", Description: "Whether the pool root is encrypted"},
		{Name: "autotrim", Type: "string", Description: "ZFS autotrim setting of the pool (on/off)"},
		{Name: "path", Type: "string", Description: "Mount path of the pool"},
		{Name: "dataVdevCount", Type: "number", Description: "Number of data vdevs in the pool topology"},
		{Name: "topology", Type: "string_array", Description: "Per-vdev topology summary, e.g. \"data: MIRROR (2 disks)\""},
	},
	"dataset": {
		{Name: "datasetType", Type: "string", Description: "Dataset type (FILESYSTEM or VOLUME/zvol)"},
		{Name: "pool", Type: "string", Description: "Pool the dataset belongs to"},
		{Name: "encrypted", Type: "boolean", Description: "Whether the dataset is encrypted"},
		{Name: "encryptionAlgorithm", Type: "string", Description: "Encryption algorithm used by the dataset, if any"},
		{Name: "compression", Type: "string", Description: "Configured compression algorithm"},
		{Name: "deduplication", Type: "string", Description: "Configured deduplication setting"},
		{Name: "atime", Type: "string", Description: "Configured atime setting"},
		{Name: "readonly", Type: "string", Description: "Configured readonly setting"},
		{Name: "sync", Type: "string", Description: "Configured sync setting"},
		{Name: "quota", Type: "string", Description: "Configured quota, or 0 when unlimited"},
		{Name: "mountpoint", Type: "string", Description: "Mount point of the dataset"},
	},
	"disk": {
		{Name: "serial", Type: "string", Description: "Disk serial number"},
		{Name: "model", Type: "string", Description: "Disk model as reported by the controller"},
		{Name: "diskType", Type: "string", Description: "Disk type (HDD or SSD)"},
		{Name: "size", Type: "number", Description: "Disk capacity in bytes"},
		{Name: "pool", Type: "string", Description: "Pool the disk is assigned to, if any"},
		{Name: "smartEnabled", Type: "boolean", Description: "Whether S.M.A.R.T. monitoring is enabled for the disk"},
		{Name: "description", Type: "string", Description: "Operator-set description of the disk"},
	},
	"share": {
		{Name: "protocol", Type: "string", Description: "Sharing protocol (smb or nfs)"},
		{Name: "path", Type: "string", Description: "Filesystem path exported by the share"},
		{Name: "enabled", Type: "boolean", Description: "Whether the share is enabled"},
		{Name: "readonly", Type: "boolean", Description: "Whether the share is exported read-only"},
		{Name: "guestAccess", Type: "boolean", Description: "Whether unauthenticated guest access is allowed (SMB)"},
		{Name: "browsable", Type: "boolean", Description: "Whether the share is browsable in network neighbourhood (SMB)"},
		{Name: "purpose", Type: "string", Description: "SMB share preset/purpose"},
		{Name: "comment", Type: "string", Description: "Operator-set comment on the share"},
		{Name: "networks", Type: "string_array", Description: "Networks authorised to mount the share (NFS)"},
		{Name: "hosts", Type: "string_array", Description: "Hosts authorised to mount the share (NFS)"},
		{Name: "maprootUser", Type: "string", Description: "User root is mapped to on the export (NFS)"},
	},
	"service": {
		{Name: "state", Type: "string", Description: "Service state reported by TrueNAS (RUNNING, STOPPED)"},
		{Name: "startOnBoot", Type: "boolean", Description: "Whether the service is configured to start on boot"},
	},
	"interface": {
		{Name: "interfaceType", Type: "string", Description: "Interface type (PHYSICAL, BRIDGE, LINK_AGGREGATION, VLAN)"},
		{Name: "addresses", Type: "string_array", Description: "Statically configured addresses in CIDR form"},
		{Name: "dhcp", Type: "boolean", Description: "Whether the interface obtains its IPv4 address over DHCP"},
		{Name: "mtu", Type: "number", Description: "Configured MTU, when not the default"},
		{Name: "description", Type: "string", Description: "Operator-set description of the interface"},
	},
	"snapshot_task": {
		{Name: "dataset", Type: "string", Description: "Dataset the periodic snapshot task covers"},
		{Name: "recursive", Type: "boolean", Description: "Whether child datasets are snapshotted too"},
		{Name: "enabled", Type: "boolean", Description: "Whether the task is enabled"},
		{Name: "schedule", Type: "string", Description: "Cron schedule of the task (minute hour dom month dow)"},
		{Name: "lifetime", Type: "string", Description: "Retention of the snapshots the task creates, e.g. \"2 WEEK\""},
		{Name: "namingSchema", Type: "string", Description: "Naming schema applied to created snapshots"},
		{Name: "exclude", Type: "string_array", Description: "Child datasets excluded from the task"},
	},
	"replication_task": {
		{Name: "direction", Type: "string", Description: "Replication direction (PUSH or PULL)"},
		{Name: "transport", Type: "string", Description: "Replication transport (SSH, SSH+NETCAT, LOCAL)"},
		{Name: "sourceDatasets", Type: "string_array", Description: "Datasets replicated from"},
		{Name: "targetDataset", Type: "string", Description: "Dataset replicated to"},
		{Name: "recursive", Type: "boolean", Description: "Whether child datasets are replicated too"},
		{Name: "enabled", Type: "boolean", Description: "Whether the task is enabled"},
		{Name: "auto", Type: "boolean", Description: "Whether the task runs automatically on its schedule"},
		{Name: "retentionPolicy", Type: "string", Description: "Retention policy applied on the target (SOURCE, CUSTOM, NONE)"},
	},
}

// Connector fetches storage configuration from a TrueNAS appliance.
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
	client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS})
	return &Connector{url: strings.TrimSuffix(rawURL, "/"), apiKey: apiKey, client: client}, nil
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "TrueNAS" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return category }

// Validate checks the configuration is complete and reaches the API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	if c.url == "" {
		return fmt.Errorf("truenas url is required")
	}
	if c.apiKey == "" {
		return fmt.Errorf("truenas api_key is required")
	}
	_, err := c.doRequest(ctx, pathSystemInfo)
	return err
}

// section describes one snapshot section: the selective-fetch field that
// gates it, the API path it reads, and how the payload is rendered.
type section struct {
	field string
	title string
	path  string
	build func(raw []byte) (string, []connector.SnapshotEntity, map[string]string)
}

// sections lists every section the connector emits, in snapshot order.
var sections = []section{
	{field: "system", title: "System", path: pathSystemInfo, build: buildSystem},
	{field: "pools", title: "Pools", path: pathPools, build: buildPools},
	{field: "datasets", title: "Datasets", path: pathDatasets, build: buildDatasets},
	{field: "disks", title: "Disks", path: pathDisks, build: buildDisks},
	{field: "shares", title: "SMB Shares", path: pathSMBShares, build: buildSMBShares},
	{field: "shares", title: "NFS Shares", path: pathNFSShares, build: buildNFSShares},
	{field: "services", title: "Services", path: pathServices, build: buildServices},
	{field: "interfaces", title: "Network Interfaces", path: pathInterfaces, build: buildInterfaces},
	{field: "snapshot_tasks", title: "Snapshot Tasks", path: pathSnapshotTasks, build: buildSnapshotTasks},
	{field: "replication_tasks", title: "Replication Tasks", path: pathReplicationTask, build: buildReplicationTasks},
}

// Fetch retrieves system info, pools, datasets, disks, shares, services,
// network interfaces and the periodic snapshot/replication tasks. Each
// section degrades on its own: an endpoint an appliance does not serve (or
// an API key without the matching privilege) leaves a placeholder instead
// of failing the whole snapshot.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	metadata := map[string]string{"truenas_url": c.url}
	var snapshotSections []connector.SnapshotSection
	var entities []connector.SnapshotEntity

	for _, s := range sections {
		if !connector.WantsField(fields, s.field) {
			continue
		}
		raw, err := c.doRequest(ctx, s.path)
		if err != nil {
			snapshotSections = append(snapshotSections, unavailable(s.title, err))
			continue
		}
		content, ents, meta := s.build(raw)
		snapshotSections = append(snapshotSections, connector.SnapshotSection{Title: s.title, Content: content})
		entities = append(entities, ents...)
		for k, v := range meta {
			metadata[k] = v
		}
	}

	return &connector.ServiceSnapshot{
		ServiceName:  "TrueNAS",
		Type:         typeName,
		Sections:     snapshotSections,
		Entities:     entities,
		Dependencies: poolDependencies(entities),
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

// poolDependencies turns the discovered pools into storage dependencies,
// deduplicated and sorted so an unchanged appliance diffs identically.
func poolDependencies(entities []connector.SnapshotEntity) []connector.ServiceDependency {
	seen := make(map[string]struct{})
	var names []string
	for _, e := range entities {
		if e.Kind != "pool" || e.Name == "" {
			continue
		}
		if _, dup := seen[e.Name]; dup {
			continue
		}
		seen[e.Name] = struct{}{}
		names = append(names, e.Name)
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	deps := make([]connector.ServiceDependency, 0, len(names))
	for _, n := range names {
		deps = append(deps, connector.ServiceDependency{Kind: "storage", Name: n})
	}
	return deps
}

func unavailable(title string, err error) connector.SnapshotSection {
	return connector.SnapshotSection{Title: title, Content: "_" + title + " unavailable: " + err.Error() + "_"}
}

func (c *Connector) doRequest(ctx context.Context, path string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
