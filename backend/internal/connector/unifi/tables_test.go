package unifi

import (
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// byExternalID indexes entities for assertions.
func byExternalID(entities []connector.SnapshotEntity) map[string]connector.SnapshotEntity {
	out := make(map[string]connector.SnapshotEntity, len(entities))
	for _, e := range entities {
		out[e.ExternalID] = e
	}
	return out
}

func TestBuildSiteTable(t *testing.T) {
	content, entities, meta := buildSiteTable([]byte(sitesJSON))
	if meta != nil {
		t.Errorf("meta = %v, want nil", meta)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	if !strings.Contains(content, "| Home | default | admin |") {
		t.Errorf("content = %q, want the default site row", content)
	}
	lab := byExternalID(entities)["7xk2lp0q"]
	if lab.Kind != "site" || lab.Name != "7xk2lp0q" {
		t.Errorf("lab site = %+v", lab)
	}
	if lab.Attributes["role"] != "readonly" || lab.Attributes["displayName"] != "Lab" {
		t.Errorf("lab attributes = %+v", lab.Attributes)
	}
}

func TestBuildDeviceTable(t *testing.T) {
	content, entities, _ := buildDeviceTable([]byte(devicesJSON))
	if len(entities) != 3 {
		t.Fatalf("entities = %d, want 3", len(entities))
	}
	ap := byExternalID(entities)["aa:bb:cc:00:11:22"]
	if ap.Kind != "device" || ap.Name != "Living Room AP" || ap.IP != "10.0.1.20" || ap.Hostname != "Living Room AP" {
		t.Errorf("ap = %+v", ap)
	}
	want := map[string]any{
		"model": "U6LR", "deviceType": "uap", "firmwareVersion": "6.6.55",
		"adopted": true, "disabled": false, "state": "connected",
	}
	for k, v := range want {
		if ap.Attributes[k] != v {
			t.Errorf("ap attributes[%q] = %v, want %v", k, ap.Attributes[k], v)
		}
	}

	// A switch reports its state as a string on some firmware.
	sw := byExternalID(entities)["aa:bb:cc:00:11:33"]
	if sw.Attributes["state"] != "offline" || sw.Attributes["disabled"] != true {
		t.Errorf("switch attributes = %+v", sw.Attributes)
	}
	// An unnamed device falls back to its MAC, and an unknown state code is
	// still rendered rather than dropped.
	gw := byExternalID(entities)["aa:bb:cc:00:11:44"]
	if gw.Name != "aa:bb:cc:00:11:44" || gw.Hostname != "" || gw.Attributes["state"] != "state-42" {
		t.Errorf("gateway = %+v", gw)
	}

	// Volatile counters must not leak into the rendered section.
	for _, volatile := range []string{"98213", "123456789", "1700000000"} {
		if strings.Contains(content, volatile) {
			t.Errorf("content contains the volatile value %q:\n%s", volatile, content)
		}
	}
}

func TestBuildNetworkTable(t *testing.T) {
	content, entities, _ := buildNetworkTable([]byte(networksJSON))
	if len(entities) != 3 {
		t.Fatalf("entities = %d, want 3", len(entities))
	}
	index := byExternalID(entities)

	lan := index["n1"]
	if lan.Kind != "network" || lan.Name != "LAN" {
		t.Errorf("lan = %+v", lan)
	}
	if lan.Attributes["enabled"] != true || lan.Attributes["subnet"] != "10.0.1.1/24" || lan.Attributes["dhcpEnabled"] != true {
		t.Errorf("lan attributes = %+v", lan.Attributes)
	}
	if _, ok := lan.Attributes["vlan"]; ok {
		t.Errorf("untagged network should carry no vlan attribute: %+v", lan.Attributes)
	}

	// vlan arrives as a string here; it must decode to a number attribute.
	iot := index["n2"]
	if iot.Attributes["vlan"] != 20 || iot.Attributes["vlanEnabled"] != true {
		t.Errorf("iot attributes = %+v", iot.Attributes)
	}
	// "enabled" is absent on the WAN entry in the fixture, but explicitly
	// false, so the default must not override it.
	if index["n3"].Attributes["enabled"] != false {
		t.Errorf("wan attributes = %+v", index["n3"].Attributes)
	}
	if !strings.Contains(content, "| IoT | corporate | 20 | 10.0.20.1/24 | yes | yes |") {
		t.Errorf("content = %q, want the IoT row", content)
	}
}

// TestBuildNetworkTableDefaultsEnabled checks the documented default for a
// firmware that omits "enabled" entirely.
func TestBuildNetworkTableDefaultsEnabled(t *testing.T) {
	_, entities, _ := buildNetworkTable([]byte(`{"data":[{"_id":"n9","name":"Old","purpose":"corporate"}]}`))
	if len(entities) != 1 || entities[0].Attributes["enabled"] != true {
		t.Fatalf("entities = %+v, want enabled defaulting to true", entities)
	}
}

func TestBuildWLANTable(t *testing.T) {
	content, entities, _ := buildWLANTable([]byte(wlansJSON))
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	index := byExternalID(entities)

	home := index["w1"]
	if home.Kind != "wlan" || home.Name != "HomeNet" {
		t.Errorf("home = %+v", home)
	}
	want := map[string]any{
		"enabled": true, "security": "wpapsk", "wpaMode": "wpa2",
		"guest": false, "hideSsid": false, "macFilterEnabled": false, "pmfMode": "optional",
	}
	for k, v := range want {
		if home.Attributes[k] != v {
			t.Errorf("home attributes[%q] = %v, want %v", k, home.Attributes[k], v)
		}
	}
	guest := index["w2"]
	if guest.Attributes["enabled"] != true || guest.Attributes["guest"] != true || guest.Attributes["macFilterPolicy"] != "allow" {
		t.Errorf("guest attributes = %+v", guest.Attributes)
	}

	// The pre-shared key must never reach a snapshot.
	if strings.Contains(content, "super-secret") {
		t.Errorf("content leaks the WLAN passphrase:\n%s", content)
	}
	for _, e := range entities {
		for k, v := range e.Attributes {
			if s, ok := v.(string); ok && s == "super-secret" {
				t.Errorf("attribute %q leaks the WLAN passphrase", k)
			}
		}
	}
}

func TestBuildFirewallTable(t *testing.T) {
	content, entities, _ := buildFirewallTable([]byte(firewallJSON))
	if len(entities) != 3 {
		t.Fatalf("entities = %d, want 3", len(entities))
	}
	// Rules render sorted by ruleset then evaluation order, whatever order
	// the controller returned them in.
	wantOrder := []string{"Allow established", "Block IoT to LAN", "Drop WAN"}
	for i, name := range wantOrder {
		if entities[i].Name != name {
			t.Errorf("entities[%d].Name = %q, want %q", i, entities[i].Name, name)
		}
	}
	block := byExternalID(entities)["f2"]
	want := map[string]any{
		"enabled": true, "action": "drop", "ruleset": "LAN_IN", "ruleIndex": 2001,
		"protocol": "all", "logging": true,
		"sourceAddress": "10.0.20.0/24", "destinationAddress": "10.0.1.0/24",
	}
	for k, v := range want {
		if block.Attributes[k] != v {
			t.Errorf("block attributes[%q] = %v, want %v", k, block.Attributes[k], v)
		}
	}
	if byExternalID(entities)["f3"].Attributes["enabled"] != true {
		t.Errorf("a rule without an explicit enabled flag should default to enabled")
	}
	if !strings.Contains(content, "| LAN_IN | 2001 | Block IoT to LAN | drop |") {
		t.Errorf("content = %q, want the block rule row", content)
	}
}

func TestBuildClientSummary(t *testing.T) {
	content, entities, meta := buildClientSummary([]byte(clientsJSON))
	if entities != nil {
		t.Errorf("entities = %+v, want none (clients are volatile)", entities)
	}
	wantMeta := map[string]string{"client_count": "4", "wired_client_count": "1", "wireless_client_count": "3"}
	for k, v := range wantMeta {
		if meta[k] != v {
			t.Errorf("meta[%q] = %q, want %q", k, meta[k], v)
		}
	}
	if !strings.Contains(content, "4 clients: 3 wireless, 1 wired") {
		t.Errorf("content = %q, want the headline counts", content)
	}
	for _, row := range []string{"| Guest | 1 |", "| HomeNet | 2 |", "| LAN | 1 |"} {
		if !strings.Contains(content, row) {
			t.Errorf("content = %q, want the row %q", content, row)
		}
	}
	// No MAC address or per-client counter should be rendered.
	for _, leaked := range []string{"11:11:11:11:11:11", "-54", "9999"} {
		if strings.Contains(content, leaked) {
			t.Errorf("content contains per-client detail %q:\n%s", leaked, content)
		}
	}
}

func TestBuildClientSummaryGroupsUnknown(t *testing.T) {
	content, _, meta := buildClientSummary([]byte(`{"data":[{"is_wired":false},{"is_wired":true}]}`))
	if !strings.Contains(content, "| (unknown) | 2 |") {
		t.Errorf("content = %q, want the clients grouped as unknown", content)
	}
	if meta["client_count"] != "2" {
		t.Errorf("meta = %v", meta)
	}
}

// TestBuildersHandleEmptyAndMalformedPayloads checks every builder degrades
// the same way instead of panicking on an unexpected payload.
func TestBuildersHandleEmptyAndMalformedPayloads(t *testing.T) {
	builders := map[string]func([]byte) (string, []connector.SnapshotEntity, map[string]string){
		"Sites":          buildSiteTable,
		"Devices":        buildDeviceTable,
		"Networks":       buildNetworkTable,
		"WLANs":          buildWLANTable,
		"Firewall Rules": buildFirewallTable,
		"Clients":        buildClientSummary,
	}
	for title, build := range builders {
		t.Run(title+"/empty", func(t *testing.T) {
			content, entities, _ := build([]byte(`{"meta":{"rc":"ok"},"data":[]}`))
			if len(entities) != 0 {
				t.Errorf("entities = %+v, want none", entities)
			}
			if !strings.HasPrefix(content, "_No ") {
				t.Errorf("content = %q, want an empty placeholder", content)
			}
		})
		t.Run(title+"/malformed", func(t *testing.T) {
			content, entities, _ := build([]byte(`{"data":{"not":"a list"}}`))
			if len(entities) != 0 {
				t.Errorf("entities = %+v, want none", entities)
			}
			if !strings.Contains(content, "malformed response") {
				t.Errorf("content = %q, want a malformed-response placeholder", content)
			}
		})
	}
}

func TestCellEscapesPipes(t *testing.T) {
	if got := cell("a|b\nc"); got != `a\|b c` {
		t.Errorf("cell = %q", got)
	}
	if got := cell(""); got != "—" {
		t.Errorf("cell(\"\") = %q", got)
	}
}

func TestFlexIntAcceptsNumbersStringsAndNull(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantSet bool
		wantErr bool
	}{
		{in: `20`, want: 20, wantSet: true},
		{in: `"20"`, want: 20, wantSet: true},
		{in: `null`},
		{in: `""`},
		{in: `"abc"`, wantErr: true},
	}
	for _, tt := range tests {
		var f flexInt
		err := f.UnmarshalJSON([]byte(tt.in))
		if (err != nil) != tt.wantErr {
			t.Errorf("UnmarshalJSON(%s) error = %v, wantErr %t", tt.in, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && (f.Value != tt.want || f.Set != tt.wantSet) {
			t.Errorf("UnmarshalJSON(%s) = %+v, want {%d %t}", tt.in, f, tt.want, tt.wantSet)
		}
	}
}
