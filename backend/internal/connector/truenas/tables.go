package truenas

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// zfsProp is a ZFS property as the TrueNAS API reports it: an object with
// "value"/"rawvalue"/"source" on current releases, but a bare string on
// some older builds and on a few properties. Both shapes decode here so a
// single odd property cannot turn an entire section into a parse failure.
type zfsProp struct {
	Value    string `json:"value"`
	RawValue string `json:"rawvalue"`
}

// UnmarshalJSON accepts either the object form or a bare string/null.
func (p *zfsProp) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" {
		return nil
	}
	if strings.HasPrefix(trimmed, `"`) {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		p.Value = s
		return nil
	}
	type alias zfsProp
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*p = zfsProp(a)
	return nil
}

// String returns the human-readable property value, falling back to the raw
// value when TrueNAS only fills that one.
func (p zfsProp) String() string {
	if p.Value != "" {
		return p.Value
	}
	return p.RawValue
}

// cell escapes a value for use inside a Markdown table cell: dataset paths
// and share comments can carry '|' or newlines, which would otherwise
// split a row into extra columns.
func cell(s string) string {
	if s == "" {
		return "—"
	}
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", `\|`)
}

func malformed(title string, err error) string {
	return "_" + title + " unavailable: " + connector.NewMalformedResponseError(err).Error() + "_"
}

func empty(noun string) string { return "_No " + noun + " returned_" }

func countMeta(key string, n int) map[string]string {
	return map[string]string{key: fmt.Sprintf("%d", n)}
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

// humanBytes renders a byte count in binary units. Capacities are static,
// so this never makes an unchanged snapshot look different.
func humanBytes(n int64) string {
	if n <= 0 {
		return ""
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit && exp < 5; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

// buildSystem renders /system/info. Uptime, load averages and the current
// datetime are deliberately dropped: they change on every poll.
func buildSystem(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var info struct {
		Version       string `json:"version"`
		Hostname      string `json:"hostname"`
		Model         string `json:"model"`
		Cores         int    `json:"cores"`
		PhysicalCores int    `json:"physical_cores"`
		PhysMem       int64  `json:"physmem"`
		Product       string `json:"system_product"`
		Manufacturer  string `json:"system_manufacturer"`
		Serial        string `json:"system_serial"`
		License       any    `json:"license"`
		ECCMemory     bool   `json:"ecc_memory"`
	}
	if err := json.Unmarshal(raw, &info); err != nil {
		return malformed("System", err), nil, nil
	}

	var b strings.Builder
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	rows := []struct{ label, value string }{
		{"Version", info.Version},
		{"Hostname", info.Hostname},
		{"Product", info.Product},
		{"Manufacturer", info.Manufacturer},
		{"Serial", info.Serial},
		{"CPU", info.Model},
		{"Cores", coreSummary(info.Cores, info.PhysicalCores)},
		{"Memory", humanBytes(info.PhysMem)},
		{"ECC memory", fmt.Sprintf("%t", info.ECCMemory)},
		{"Licensed", fmt.Sprintf("%t", info.License != nil)},
	}
	for _, r := range rows {
		if r.value == "" {
			continue
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s |\n", r.label, cell(r.value))
	}

	metadata := map[string]string{}
	if info.Version != "" {
		metadata["truenas_version"] = info.Version
		metadata["truenas_flavor"] = flavor(info.Version)
	}
	if info.Hostname != "" {
		metadata["truenas_hostname"] = info.Hostname
	}
	return b.String(), nil, metadata
}

func coreSummary(logical, physical int) string {
	switch {
	case logical <= 0:
		return ""
	case physical > 0 && physical != logical:
		return fmt.Sprintf("%d (%d physical)", logical, physical)
	default:
		return fmt.Sprintf("%d", logical)
	}
}

// flavor reports which TrueNAS edition a version string names, so a
// snapshot records SCALE vs CORE without a second API call.
func flavor(version string) string {
	upper := strings.ToUpper(version)
	switch {
	case strings.Contains(upper, "SCALE"):
		return "SCALE"
	case strings.Contains(upper, "CORE"):
		return "CORE"
	default:
		return "unknown"
	}
}

// buildPools renders /pool. Pool size and allocation are omitted: they move
// with every write and would churn the snapshot.
func buildPools(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var pools []struct {
		Name     string  `json:"name"`
		Status   string  `json:"status"`
		Healthy  bool    `json:"healthy"`
		Path     string  `json:"path"`
		Encrypt  int     `json:"encrypt"`
		AutoTrim zfsProp `json:"autotrim"`
		Topology struct {
			Data    []vdev `json:"data"`
			Cache   []vdev `json:"cache"`
			Log     []vdev `json:"log"`
			Spare   []vdev `json:"spare"`
			Special []vdev `json:"special"`
		} `json:"topology"`
	}
	if err := json.Unmarshal(raw, &pools); err != nil {
		return malformed("Pools", err), nil, nil
	}
	if len(pools) == 0 {
		return empty("pools"), nil, countMeta("pool_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Pool | Status | Healthy | Encrypted | Topology | Path |\n")
	b.WriteString("|------|--------|---------|-----------|----------|------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(pools))
	for _, p := range pools {
		groups := []struct {
			label string
			vdevs []vdev
		}{
			{"data", p.Topology.Data}, {"special", p.Topology.Special},
			{"log", p.Topology.Log}, {"cache", p.Topology.Cache}, {"spare", p.Topology.Spare},
		}
		var topology []string
		for _, g := range groups {
			topology = append(topology, summarizeVdevs(g.label, g.vdevs)...)
		}
		encrypted := p.Encrypt > 0
		_, _ = fmt.Fprintf(&b, "| %s | %s | %t | %t | %s | %s |\n",
			cell(p.Name), cell(p.Status), p.Healthy, encrypted,
			cell(strings.Join(topology, "<br>")), cell(p.Path))

		attrs := map[string]any{
			"healthy":       p.Healthy,
			"encrypted":     encrypted,
			"dataVdevCount": len(p.Topology.Data),
		}
		putString(attrs, "status", p.Status)
		putString(attrs, "path", p.Path)
		putString(attrs, "autotrim", p.AutoTrim.String())
		putStrings(attrs, "topology", topology)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "pool", Name: p.Name, ExternalID: p.Name, Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("pool_count", len(pools))
}

// vdev is one virtual device in a pool topology group.
type vdev struct {
	Type     string `json:"type"`
	Children []struct {
		Type string `json:"type"`
	} `json:"children"`
}

// summarizeVdevs describes a topology group as one stable line per vdev,
// e.g. "data: MIRROR (2 disks)".
func summarizeVdevs(label string, vdevs []vdev) []string {
	out := make([]string, 0, len(vdevs))
	for _, v := range vdevs {
		typ := v.Type
		if typ == "" {
			typ = "UNKNOWN"
		}
		if n := len(v.Children); n > 0 {
			out = append(out, fmt.Sprintf("%s: %s (%d disks)", label, typ, n))
			continue
		}
		out = append(out, fmt.Sprintf("%s: %s", label, typ))
	}
	return out
}

// buildDatasets renders /pool/dataset. Used/available space is omitted for
// the same reason pool allocation is.
func buildDatasets(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var datasets []dataset
	if err := json.Unmarshal(raw, &datasets); err != nil {
		return malformed("Datasets", err), nil, nil
	}
	flat := flattenDatasets(datasets)
	if len(flat) == 0 {
		return empty("datasets"), nil, countMeta("dataset_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Dataset | Type | Encrypted | Compression | Quota | Read-only | Mountpoint |\n")
	b.WriteString("|---------|------|-----------|-------------|-------|-----------|------------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(flat))
	for _, d := range flat {
		quota := d.Quota.String()
		if quota == "" || quota == "0" {
			quota = "none"
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s | %t | %s | %s | %s | %s |\n",
			cell(d.Name), cell(d.Type), d.Encrypted, cell(d.Compression.String()),
			cell(quota), cell(d.ReadOnly.String()), cell(d.MountPoint))

		attrs := map[string]any{"encrypted": d.Encrypted}
		putString(attrs, "datasetType", d.Type)
		putString(attrs, "pool", d.Pool)
		putString(attrs, "encryptionAlgorithm", d.EncryptionAlgorithm)
		putString(attrs, "compression", d.Compression.String())
		putString(attrs, "deduplication", d.Deduplication.String())
		putString(attrs, "atime", d.ATime.String())
		putString(attrs, "readonly", d.ReadOnly.String())
		putString(attrs, "sync", d.Sync.String())
		putString(attrs, "quota", d.Quota.String())
		putString(attrs, "mountpoint", d.MountPoint)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "dataset", Name: d.Name, ExternalID: d.ID, Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("dataset_count", len(flat))
}

// dataset mirrors the fields of /pool/dataset this connector reads.
// Children are present when the appliance returns the nested tree form.
type dataset struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Pool                string    `json:"pool"`
	Encrypted           bool      `json:"encrypted"`
	EncryptionAlgorithm string    `json:"encryption_algorithm"`
	MountPoint          string    `json:"mountpoint"`
	Compression         zfsProp   `json:"compression"`
	Deduplication       zfsProp   `json:"deduplication"`
	ATime               zfsProp   `json:"atime"`
	ReadOnly            zfsProp   `json:"readonly"`
	Sync                zfsProp   `json:"sync"`
	Quota               zfsProp   `json:"quota"`
	Children            []dataset `json:"children"`
}

// flattenDatasets walks the nested tree TrueNAS returns for /pool/dataset
// and yields every dataset once, sorted by name so the ordering is stable
// regardless of how the API nests them.
func flattenDatasets(roots []dataset) []dataset {
	var out []dataset
	seen := make(map[string]struct{})
	var walk func([]dataset)
	walk = func(items []dataset) {
		for _, d := range items {
			key := d.ID
			if key == "" {
				key = d.Name
			}
			if _, dup := seen[key]; !dup {
				seen[key] = struct{}{}
				children := d.Children
				d.Children = nil
				out = append(out, d)
				walk(children)
				continue
			}
			walk(d.Children)
		}
	}
	walk(roots)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// buildDisks renders /disk. Temperatures are not read: they change
// constantly and belong in metrics, not a configuration snapshot.
func buildDisks(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var disks []struct {
		Name        string `json:"name"`
		Serial      string `json:"serial"`
		Model       string `json:"model"`
		Type        string `json:"type"`
		Size        int64  `json:"size"`
		Pool        string `json:"pool"`
		Description string `json:"description"`
		ToggleSMART bool   `json:"togglesmart"`
	}
	if err := json.Unmarshal(raw, &disks); err != nil {
		return malformed("Disks", err), nil, nil
	}
	if len(disks) == 0 {
		return empty("disks"), nil, countMeta("disk_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Disk | Type | Size | Model | Serial | Pool | S.M.A.R.T. |\n")
	b.WriteString("|------|------|------|-------|--------|------|------------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(disks))
	for _, d := range disks {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %t |\n",
			cell(d.Name), cell(d.Type), cell(humanBytes(d.Size)), cell(d.Model),
			cell(d.Serial), cell(d.Pool), d.ToggleSMART)

		attrs := map[string]any{"smartEnabled": d.ToggleSMART}
		if d.Size > 0 {
			attrs["size"] = d.Size
		}
		putString(attrs, "serial", d.Serial)
		putString(attrs, "model", d.Model)
		putString(attrs, "diskType", d.Type)
		putString(attrs, "pool", d.Pool)
		putString(attrs, "description", d.Description)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "disk", Name: d.Name, ExternalID: d.Name, Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("disk_count", len(disks))
}

// buildSMBShares renders /sharing/smb.
func buildSMBShares(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var shares []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Path      string `json:"path"`
		Enabled   bool   `json:"enabled"`
		ReadOnly  bool   `json:"ro"`
		GuestOK   bool   `json:"guestok"`
		Browsable bool   `json:"browsable"`
		Purpose   string `json:"purpose"`
		Comment   string `json:"comment"`
	}
	if err := json.Unmarshal(raw, &shares); err != nil {
		return malformed("SMB Shares", err), nil, nil
	}
	if len(shares) == 0 {
		return empty("SMB shares"), nil, countMeta("smb_share_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Share | Path | Enabled | Read-only | Guest | Browsable | Purpose |\n")
	b.WriteString("|-------|------|---------|-----------|-------|-----------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(shares))
	for _, s := range shares {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %t | %t | %t | %t | %s |\n",
			cell(s.Name), cell(s.Path), s.Enabled, s.ReadOnly, s.GuestOK, s.Browsable, cell(s.Purpose))

		attrs := map[string]any{
			"protocol":    "smb",
			"enabled":     s.Enabled,
			"readonly":    s.ReadOnly,
			"guestAccess": s.GuestOK,
			"browsable":   s.Browsable,
		}
		putString(attrs, "path", s.Path)
		putString(attrs, "purpose", s.Purpose)
		putString(attrs, "comment", s.Comment)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "share", Name: s.Name, ExternalID: fmt.Sprintf("smb:%d", s.ID), Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("smb_share_count", len(shares))
}

// buildNFSShares renders /sharing/nfs. TrueNAS 24.10 replaced the "paths"
// list with a single "path"; both shapes are accepted.
func buildNFSShares(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var shares []struct {
		ID          int      `json:"id"`
		Path        string   `json:"path"`
		Paths       []string `json:"paths"`
		Enabled     bool     `json:"enabled"`
		ReadOnly    bool     `json:"ro"`
		MaprootUser string   `json:"maproot_user"`
		Networks    []string `json:"networks"`
		Hosts       []string `json:"hosts"`
		Comment     string   `json:"comment"`
	}
	if err := json.Unmarshal(raw, &shares); err != nil {
		return malformed("NFS Shares", err), nil, nil
	}
	if len(shares) == 0 {
		return empty("NFS shares"), nil, countMeta("nfs_share_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Export | Enabled | Read-only | Networks | Hosts | Maproot |\n")
	b.WriteString("|--------|---------|-----------|----------|-------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(shares))
	for _, s := range shares {
		paths := s.Paths
		if len(paths) == 0 && s.Path != "" {
			paths = []string{s.Path}
		}
		name := strings.Join(paths, ", ")
		_, _ = fmt.Fprintf(&b, "| %s | %t | %t | %s | %s | %s |\n",
			cell(name), s.Enabled, s.ReadOnly, cell(strings.Join(s.Networks, ", ")),
			cell(strings.Join(s.Hosts, ", ")), cell(s.MaprootUser))

		attrs := map[string]any{"protocol": "nfs", "enabled": s.Enabled, "readonly": s.ReadOnly}
		putString(attrs, "path", strings.Join(paths, ", "))
		putString(attrs, "maprootUser", s.MaprootUser)
		putString(attrs, "comment", s.Comment)
		putStrings(attrs, "networks", s.Networks)
		putStrings(attrs, "hosts", s.Hosts)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "share", Name: name, ExternalID: fmt.Sprintf("nfs:%d", s.ID), Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("nfs_share_count", len(shares))
}

// buildServices renders /service.
func buildServices(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var services []struct {
		ID      int    `json:"id"`
		Service string `json:"service"`
		Enable  bool   `json:"enable"`
		State   string `json:"state"`
	}
	if err := json.Unmarshal(raw, &services); err != nil {
		return malformed("Services", err), nil, nil
	}
	if len(services) == 0 {
		return empty("services"), nil, countMeta("service_count", 0)
	}
	sort.Slice(services, func(i, j int) bool { return services[i].Service < services[j].Service })

	var b strings.Builder
	b.WriteString("| Service | State | Start on boot |\n")
	b.WriteString("|---------|-------|---------------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(services))
	running := 0
	for _, s := range services {
		if strings.EqualFold(s.State, "RUNNING") {
			running++
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s | %t |\n", cell(s.Service), cell(s.State), s.Enable)

		attrs := map[string]any{"startOnBoot": s.Enable}
		putString(attrs, "state", s.State)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "service", Name: s.Service, ExternalID: s.Service, Attributes: attrs,
		})
	}
	metadata := countMeta("service_count", len(services))
	metadata["services_running"] = fmt.Sprintf("%d", running)
	return b.String(), entities, metadata
}

// buildInterfaces renders /interface. Only the configured addresses are
// read; link state and traffic counters are left out so an idle appliance
// keeps producing an identical snapshot.
func buildInterfaces(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var ifaces []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		IPv4DHCP    bool   `json:"ipv4_dhcp"`
		MTU         int    `json:"mtu"`
		Aliases     []struct {
			Type    string `json:"type"`
			Address string `json:"address"`
			Netmask int    `json:"netmask"`
		} `json:"aliases"`
	}
	if err := json.Unmarshal(raw, &ifaces); err != nil {
		return malformed("Network Interfaces", err), nil, nil
	}
	if len(ifaces) == 0 {
		return empty("network interfaces"), nil, countMeta("interface_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Interface | Type | Addresses | DHCP | MTU | Description |\n")
	b.WriteString("|-----------|------|-----------|------|-----|-------------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(ifaces))
	for _, i := range ifaces {
		var addresses []string
		for _, a := range i.Aliases {
			if a.Address == "" {
				continue
			}
			if a.Netmask > 0 {
				addresses = append(addresses, fmt.Sprintf("%s/%d", a.Address, a.Netmask))
				continue
			}
			addresses = append(addresses, a.Address)
		}
		mtu := "default"
		if i.MTU > 0 {
			mtu = fmt.Sprintf("%d", i.MTU)
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %t | %s | %s |\n",
			cell(i.Name), cell(i.Type), cell(strings.Join(addresses, ", ")), i.IPv4DHCP, mtu, cell(i.Description))

		attrs := map[string]any{"dhcp": i.IPv4DHCP}
		if i.MTU > 0 {
			attrs["mtu"] = i.MTU
		}
		putString(attrs, "interfaceType", i.Type)
		putString(attrs, "description", i.Description)
		putStrings(attrs, "addresses", addresses)
		name := i.Name
		if name == "" {
			name = i.ID
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind: "interface", Name: name, ExternalID: name, IP: firstIP(addresses), Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("interface_count", len(ifaces))
}

// firstIP strips the prefix length off the first configured address so
// cross-connector matching can join on a bare IP.
func firstIP(addresses []string) string {
	if len(addresses) == 0 {
		return ""
	}
	addr, _, found := strings.Cut(addresses[0], "/")
	if !found {
		return addresses[0]
	}
	return addr
}

// buildSnapshotTasks renders /pool/snapshottask. The last-run state is
// skipped; only the task definition is snapshotted.
func buildSnapshotTasks(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var tasks []struct {
		ID            int      `json:"id"`
		Dataset       string   `json:"dataset"`
		Recursive     bool     `json:"recursive"`
		Enabled       bool     `json:"enabled"`
		LifetimeValue int      `json:"lifetime_value"`
		LifetimeUnit  string   `json:"lifetime_unit"`
		NamingSchema  string   `json:"naming_schema"`
		Exclude       []string `json:"exclude"`
		Schedule      struct {
			Minute string `json:"minute"`
			Hour   string `json:"hour"`
			DOM    string `json:"dom"`
			Month  string `json:"month"`
			DOW    string `json:"dow"`
		} `json:"schedule"`
	}
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return malformed("Snapshot Tasks", err), nil, nil
	}
	if len(tasks) == 0 {
		return empty("snapshot tasks"), nil, countMeta("snapshot_task_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Dataset | Recursive | Schedule | Retention | Naming Schema | Enabled |\n")
	b.WriteString("|---------|-----------|----------|-----------|---------------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(tasks))
	for _, t := range tasks {
		schedule := strings.TrimSpace(strings.Join([]string{
			t.Schedule.Minute, t.Schedule.Hour, t.Schedule.DOM, t.Schedule.Month, t.Schedule.DOW,
		}, " "))
		lifetime := ""
		if t.LifetimeValue > 0 {
			lifetime = fmt.Sprintf("%d %s", t.LifetimeValue, t.LifetimeUnit)
		}
		_, _ = fmt.Fprintf(&b, "| %s | %t | %s | %s | %s | %t |\n",
			cell(t.Dataset), t.Recursive, cell(schedule), cell(lifetime), cell(t.NamingSchema), t.Enabled)

		attrs := map[string]any{"recursive": t.Recursive, "enabled": t.Enabled}
		putString(attrs, "dataset", t.Dataset)
		putString(attrs, "schedule", schedule)
		putString(attrs, "lifetime", lifetime)
		putString(attrs, "namingSchema", t.NamingSchema)
		putStrings(attrs, "exclude", t.Exclude)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "snapshot_task", Name: t.Dataset, ExternalID: fmt.Sprintf("snapshottask:%d", t.ID), Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("snapshot_task_count", len(tasks))
}

// buildReplicationTasks renders /replication.
func buildReplicationTasks(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var tasks []struct {
		ID              int      `json:"id"`
		Name            string   `json:"name"`
		Direction       string   `json:"direction"`
		Transport       string   `json:"transport"`
		SourceDatasets  []string `json:"source_datasets"`
		TargetDataset   string   `json:"target_dataset"`
		Recursive       bool     `json:"recursive"`
		Enabled         bool     `json:"enabled"`
		Auto            bool     `json:"auto"`
		RetentionPolicy string   `json:"retention_policy"`
	}
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return malformed("Replication Tasks", err), nil, nil
	}
	if len(tasks) == 0 {
		return empty("replication tasks"), nil, countMeta("replication_task_count", 0)
	}

	var b strings.Builder
	b.WriteString("| Task | Direction | Transport | Source | Target | Auto | Enabled |\n")
	b.WriteString("|------|-----------|-----------|--------|--------|------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(tasks))
	for _, t := range tasks {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %t | %t |\n",
			cell(t.Name), cell(t.Direction), cell(t.Transport), cell(strings.Join(t.SourceDatasets, ", ")),
			cell(t.TargetDataset), t.Auto, t.Enabled)

		attrs := map[string]any{"recursive": t.Recursive, "enabled": t.Enabled, "auto": t.Auto}
		putString(attrs, "direction", t.Direction)
		putString(attrs, "transport", t.Transport)
		putString(attrs, "targetDataset", t.TargetDataset)
		putString(attrs, "retentionPolicy", t.RetentionPolicy)
		putStrings(attrs, "sourceDatasets", t.SourceDatasets)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "replication_task", Name: t.Name, ExternalID: fmt.Sprintf("replication:%d", t.ID), Attributes: attrs,
		})
	}
	return b.String(), entities, countMeta("replication_task_count", len(tasks))
}
