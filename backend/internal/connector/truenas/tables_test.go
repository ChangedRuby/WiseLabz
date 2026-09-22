package truenas

import (
	"reflect"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// builders indexes every section builder by the empty-state text it renders
// when the appliance returns nothing.
var builders = []struct {
	name      string
	build     func([]byte) (string, []connector.SnapshotEntity, map[string]string)
	wantEmpty string
	countKey  string
}{
	{"pools", buildPools, "_No pools returned_", "pool_count"},
	{"datasets", buildDatasets, "_No datasets returned_", "dataset_count"},
	{"disks", buildDisks, "_No disks returned_", "disk_count"},
	{"smb shares", buildSMBShares, "_No SMB shares returned_", "smb_share_count"},
	{"nfs shares", buildNFSShares, "_No NFS shares returned_", "nfs_share_count"},
	{"services", buildServices, "_No services returned_", "service_count"},
	{"interfaces", buildInterfaces, "_No network interfaces returned_", "interface_count"},
	{"snapshot tasks", buildSnapshotTasks, "_No snapshot tasks returned_", "snapshot_task_count"},
	{"replication tasks", buildReplicationTasks, "_No replication tasks returned_", "replication_task_count"},
}

func TestBuildersOnEmptyAndMalformedInput(t *testing.T) {
	for _, b := range builders {
		t.Run(b.name+" empty array", func(t *testing.T) {
			content, entities, meta := b.build([]byte(`[]`))
			if content != b.wantEmpty || len(entities) != 0 {
				t.Fatalf("content = %q (%d entities), want %q (0)", content, len(entities), b.wantEmpty)
			}
			if meta[b.countKey] != "0" {
				t.Fatalf("meta[%q] = %q, want \"0\"", b.countKey, meta[b.countKey])
			}
		})
		t.Run(b.name+" invalid JSON", func(t *testing.T) {
			content, entities, meta := b.build([]byte(`not json`))
			if !strings.Contains(content, "malformed response") || len(entities) != 0 || meta != nil {
				t.Fatalf("content = %q (%d entities, meta %v), want a malformed-response placeholder",
					content, len(entities), meta)
			}
		})
	}

	content, entities, meta := buildSystem([]byte(`not json`))
	if !strings.Contains(content, "malformed response") || entities != nil || meta != nil {
		t.Fatalf("buildSystem(bad) = %q / %v / %v, want a malformed-response placeholder", content, entities, meta)
	}
}

func TestBuildSystem(t *testing.T) {
	content, entities, meta := buildSystem([]byte(systemInfoJSON))
	if entities != nil {
		t.Errorf("buildSystem emitted %d entities, want none", len(entities))
	}
	for _, want := range []string{
		"| Version | TrueNAS-SCALE-24.04.2 |",
		"| Hostname | nas.lan |",
		"| Serial | ABC1234 |",
		"| Cores | 12 (6 physical) |",
		"| Memory | 64.0 GiB |",
		"| ECC memory | true |",
		"| Licensed | false |",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing %q:\n%s", want, content)
		}
	}
	wantMeta := map[string]string{
		"truenas_version":  "TrueNAS-SCALE-24.04.2",
		"truenas_flavor":   "SCALE",
		"truenas_hostname": "nas.lan",
	}
	if !reflect.DeepEqual(meta, wantMeta) {
		t.Errorf("meta = %v, want %v", meta, wantMeta)
	}
}

func TestFlavor(t *testing.T) {
	tests := map[string]string{
		"TrueNAS-SCALE-24.04.2": "SCALE",
		"TrueNAS-13.0-U6.1":     "unknown",
		"FreeNAS-CORE-12":       "CORE",
		"":                      "unknown",
	}
	for version, want := range tests {
		if got := flavor(version); got != want {
			t.Errorf("flavor(%q) = %q, want %q", version, got, want)
		}
	}
}

func TestHumanBytes(t *testing.T) {
	tests := map[int64]string{
		0:             "",
		-1:            "",
		512:           "512 B",
		1024:          "1.0 KiB",
		68719476736:   "64.0 GiB",
		4000787030016: "3.6 TiB",
	}
	for n, want := range tests {
		if got := humanBytes(n); got != want {
			t.Errorf("humanBytes(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestBuildPools(t *testing.T) {
	content, entities, meta := buildPools([]byte(poolsJSON))
	if len(entities) != 2 || meta["pool_count"] != "2" {
		t.Fatalf("entities = %d, meta = %v", len(entities), meta)
	}
	tank := entities[0]
	if tank.Kind != "pool" || tank.Name != "tank" || tank.ExternalID != "tank" {
		t.Fatalf("tank = %+v", tank)
	}
	wantAttrs := map[string]any{
		"status":        "ONLINE",
		"healthy":       true,
		"encrypted":     false,
		"autotrim":      "off",
		"path":          "/mnt/tank",
		"dataVdevCount": 1,
		"topology":      []string{"data: MIRROR (2 disks)", "log: DISK"},
	}
	if !reflect.DeepEqual(tank.Attributes, wantAttrs) {
		t.Errorf("tank attributes = %+v, want %+v", tank.Attributes, wantAttrs)
	}

	vault := entities[1]
	if vault.Attributes["encrypted"] != true || vault.Attributes["healthy"] != false {
		t.Errorf("vault attributes = %+v, want encrypted and unhealthy", vault.Attributes)
	}
	// "autotrim" arrives as a bare string on this pool: both shapes decode.
	if vault.Attributes["autotrim"] != "on" {
		t.Errorf("vault autotrim = %v, want \"on\" from the bare-string form", vault.Attributes["autotrim"])
	}
	if !strings.Contains(content, "data: RAIDZ2 (4 disks)") {
		t.Errorf("content missing the RAIDZ2 topology:\n%s", content)
	}
}

// TestBuildDatasetsFlattensTree checks the nested children TrueNAS returns
// are hoisted into one sorted list, with each dataset counted once.
func TestBuildDatasetsFlattensTree(t *testing.T) {
	content, entities, meta := buildDatasets([]byte(datasetsJSON))
	if meta["dataset_count"] != "3" {
		t.Fatalf("meta = %v, want 3 datasets", meta)
	}
	var names []string
	for _, e := range entities {
		names = append(names, e.Name)
	}
	want := []string{"tank", "tank/media", "tank/vm-disk"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("dataset names = %v, want %v", names, want)
	}

	media := entities[1]
	if media.ExternalID != "tank/media" {
		t.Errorf("media ExternalID = %q", media.ExternalID)
	}
	wantAttrs := map[string]any{
		"datasetType":         "FILESYSTEM",
		"pool":                "tank",
		"encrypted":           true,
		"encryptionAlgorithm": "AES-256-GCM",
		"compression":         "ZSTD",
		"deduplication":       "OFF",
		"atime":               "OFF",
		"readonly":            "ON",
		"sync":                "DISABLED",
		"quota":               "500G",
		"mountpoint":          "/mnt/tank/media",
	}
	if !reflect.DeepEqual(media.Attributes, wantAttrs) {
		t.Errorf("media attributes = %+v, want %+v", media.Attributes, wantAttrs)
	}
	// A null property must not blow up, and a zero quota reads as "none".
	if _, ok := entities[2].Attributes["atime"]; ok {
		t.Errorf("vm-disk atime = %v, want it omitted for a null property", entities[2].Attributes["atime"])
	}
	if !strings.Contains(content, "| tank | FILESYSTEM | false | LZ4 | none |") {
		t.Errorf("content missing the tank row with an unlimited quota:\n%s", content)
	}
}

func TestBuildDisks(t *testing.T) {
	content, entities, meta := buildDisks([]byte(disksJSON))
	if len(entities) != 2 || meta["disk_count"] != "2" {
		t.Fatalf("entities = %d, meta = %v", len(entities), meta)
	}
	wantAttrs := map[string]any{
		"serial":       "WD-A1",
		"model":        "WDC WD40EFRX",
		"diskType":     "HDD",
		"size":         int64(4000787030016),
		"pool":         "tank",
		"smartEnabled": true,
		"description":  "bay 1",
	}
	if !reflect.DeepEqual(entities[0].Attributes, wantAttrs) {
		t.Errorf("sda attributes = %+v, want %+v", entities[0].Attributes, wantAttrs)
	}
	if _, ok := entities[1].Attributes["pool"]; ok {
		t.Errorf("unassigned disk carries a pool attribute: %+v", entities[1].Attributes)
	}
	if !strings.Contains(content, "3.6 TiB") {
		t.Errorf("content missing a rendered capacity:\n%s", content)
	}
}

func TestBuildShares(t *testing.T) {
	_, smb, smbMeta := buildSMBShares([]byte(smbSharesJSON))
	if len(smb) != 2 || smbMeta["smb_share_count"] != "2" {
		t.Fatalf("smb = %d, meta = %v", len(smb), smbMeta)
	}
	if smb[0].Kind != "share" || smb[0].ExternalID != "smb:1" || smb[0].Attributes["protocol"] != "smb" {
		t.Errorf("smb[0] = %+v", smb[0])
	}
	if smb[1].Attributes["guestAccess"] != true || smb[1].Attributes["enabled"] != false {
		t.Errorf("public share attributes = %+v, want guest access and disabled", smb[1].Attributes)
	}

	content, nfs, nfsMeta := buildNFSShares([]byte(nfsSharesJSON))
	if len(nfs) != 2 || nfsMeta["nfs_share_count"] != "2" {
		t.Fatalf("nfs = %d, meta = %v", len(nfs), nfsMeta)
	}
	if nfs[0].ExternalID != "nfs:3" || nfs[0].Name != "/mnt/tank/backups" {
		t.Errorf("nfs[0] = %+v, want the \"paths\" list honoured", nfs[0])
	}
	// TrueNAS 24.10 replaced "paths" with a single "path".
	if nfs[1].Name != "/mnt/vault/archive" || nfs[1].Attributes["path"] != "/mnt/vault/archive" {
		t.Errorf("nfs[1] = %+v, want the single \"path\" form honoured", nfs[1])
	}
	if !strings.Contains(content, "10.0.0.0/24") {
		t.Errorf("content missing the authorised network:\n%s", content)
	}
}

// TestBuildSharesEscapesPipes guards the Markdown table against comments
// and paths containing '|'.
func TestBuildSharesEscapesPipes(t *testing.T) {
	content, _, _ := buildSMBShares([]byte(`[{"id":1,"name":"a|b","path":"/mnt/x","purpose":"c|d"}]`))
	if strings.Contains(content, "a|b") || !strings.Contains(content, `a\|b`) {
		t.Errorf("content did not escape the pipe:\n%s", content)
	}
}

// TestBuildServicesSortsByName keeps the section stable no matter what
// order the API lists services in.
func TestBuildServicesSortsByName(t *testing.T) {
	content, entities, meta := buildServices([]byte(servicesJSON))
	var names []string
	for _, e := range entities {
		names = append(names, e.Name)
	}
	if want := []string{"cifs", "nfs", "ssh"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("services = %v, want %v", names, want)
	}
	if meta["service_count"] != "3" || meta["services_running"] != "2" {
		t.Errorf("meta = %v", meta)
	}
	if !strings.Contains(content, "| nfs | STOPPED | false |") {
		t.Errorf("content missing the stopped nfs row:\n%s", content)
	}
	if entities[2].Attributes["state"] != "RUNNING" || entities[2].Attributes["startOnBoot"] != true {
		t.Errorf("ssh attributes = %+v", entities[2].Attributes)
	}
}

func TestBuildInterfaces(t *testing.T) {
	content, entities, meta := buildInterfaces([]byte(interfacesJSON))
	if len(entities) != 2 || meta["interface_count"] != "2" {
		t.Fatalf("entities = %d, meta = %v", len(entities), meta)
	}
	eno1 := entities[0]
	if eno1.IP != "10.0.0.10" {
		t.Errorf("eno1 IP = %q, want the prefix stripped for cross-connector matching", eno1.IP)
	}
	wantAttrs := map[string]any{
		"interfaceType": "PHYSICAL",
		"description":   "lan uplink",
		"addresses":     []string{"10.0.0.10/24"},
		"dhcp":          false,
		"mtu":           9000,
	}
	if !reflect.DeepEqual(eno1.Attributes, wantAttrs) {
		t.Errorf("eno1 attributes = %+v, want %+v", eno1.Attributes, wantAttrs)
	}
	br0 := entities[1]
	if br0.IP != "" || br0.Attributes["dhcp"] != true {
		t.Errorf("br0 = %+v, want no static address and dhcp true", br0)
	}
	if _, ok := br0.Attributes["mtu"]; ok {
		t.Errorf("br0 carries an mtu attribute for the default MTU: %+v", br0.Attributes)
	}
	if !strings.Contains(content, "| br0 | BRIDGE | — | true | default | — |") {
		t.Errorf("content missing the br0 row:\n%s", content)
	}
}

func TestBuildSnapshotTasks(t *testing.T) {
	content, entities, meta := buildSnapshotTasks([]byte(snapshotTasksJSON))
	if len(entities) != 2 || meta["snapshot_task_count"] != "2" {
		t.Fatalf("entities = %d, meta = %v", len(entities), meta)
	}
	wantAttrs := map[string]any{
		"dataset":      "tank/media",
		"recursive":    true,
		"enabled":      true,
		"schedule":     "0 2 * * *",
		"lifetime":     "2 WEEK",
		"namingSchema": "auto-%Y%m%d.%H%M",
		"exclude":      []string{"tank/media/tmp"},
	}
	if !reflect.DeepEqual(entities[0].Attributes, wantAttrs) {
		t.Errorf("task attributes = %+v, want %+v", entities[0].Attributes, wantAttrs)
	}
	if entities[0].ExternalID != "snapshottask:1" {
		t.Errorf("ExternalID = %q", entities[0].ExternalID)
	}
	if _, ok := entities[1].Attributes["lifetime"]; ok {
		t.Errorf("task without a lifetime carries one: %+v", entities[1].Attributes)
	}
	if !strings.Contains(content, "*/15 * * * *") {
		t.Errorf("content missing the second task's schedule:\n%s", content)
	}
}

func TestBuildReplicationTasks(t *testing.T) {
	content, entities, meta := buildReplicationTasks([]byte(replicationTasksJSON))
	if len(entities) != 1 || meta["replication_task_count"] != "1" {
		t.Fatalf("entities = %d, meta = %v", len(entities), meta)
	}
	wantAttrs := map[string]any{
		"direction":       "PUSH",
		"transport":       "SSH",
		"sourceDatasets":  []string{"tank/media"},
		"targetDataset":   "backup/media",
		"recursive":       true,
		"enabled":         true,
		"auto":            true,
		"retentionPolicy": "SOURCE",
	}
	if !reflect.DeepEqual(entities[0].Attributes, wantAttrs) {
		t.Errorf("attributes = %+v, want %+v", entities[0].Attributes, wantAttrs)
	}
	if entities[0].ExternalID != "replication:1" || entities[0].Kind != "replication_task" {
		t.Errorf("entity = %+v", entities[0])
	}
	if !strings.Contains(content, "| offsite | PUSH | SSH | tank/media | backup/media | true | true |") {
		t.Errorf("content missing the task row:\n%s", content)
	}
}

func TestZFSPropDecodesBothShapes(t *testing.T) {
	tests := []struct {
		json string
		want string
	}{
		{`{"value":"LZ4","rawvalue":"lz4"}`, "LZ4"},
		{`{"rawvalue":"lz4"}`, "lz4"},
		{`"on"`, "on"},
		{`null`, ""},
		{`{}`, ""},
	}
	for _, tt := range tests {
		var p zfsProp
		if err := p.UnmarshalJSON([]byte(tt.json)); err != nil {
			t.Fatalf("UnmarshalJSON(%s): %v", tt.json, err)
		}
		if got := p.String(); got != tt.want {
			t.Errorf("zfsProp(%s).String() = %q, want %q", tt.json, got, tt.want)
		}
	}
	var p zfsProp
	if err := p.UnmarshalJSON([]byte(`[1,2]`)); err == nil {
		t.Error("UnmarshalJSON([1,2]) = nil, want an error")
	}
}

func TestCell(t *testing.T) {
	if got := cell(""); got != "—" {
		t.Errorf("cell(\"\") = %q, want an em dash", got)
	}
	if got := cell("a|b\nc"); got != `a\|b c` {
		t.Errorf("cell = %q", got)
	}
}
