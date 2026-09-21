package unifi

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// deviceStates maps the controller's numeric device state onto the label the
// UniFi UI shows. Unknown codes render as "state-<n>" rather than being
// dropped, so a firmware that adds a state is still legible.
var deviceStates = map[int]string{
	0:  "offline",
	1:  "connected",
	2:  "pending adoption",
	4:  "updating",
	5:  "provisioning",
	6:  "unreachable",
	7:  "adopting",
	9:  "adoption error",
	11: "isolated",
}

// flexInt decodes a controller field that is a number in some firmware
// versions and a decimal string in others (vlan, rule_index).
type flexInt struct {
	Value int
	Set   bool
}

// UnmarshalJSON accepts a JSON number, a decimal string or null.
func (f *flexInt) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("parse %q as an integer: %w", s, err)
	}
	f.Value, f.Set = n, true
	return nil
}

// cell escapes a value for use inside a Markdown table cell; controller
// fields such as firewall addresses may legitimately contain '|'.
func cell(s string) string {
	if s == "" {
		return "—"
	}
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", `\|`)
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// boolOr resolves an optional controller boolean, which is absent on older
// firmware, to its documented default.
func boolOr(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

func malformed(title string, err error) string {
	return "_" + title + " unavailable: " + connector.NewMalformedResponseError(err).Error() + "_"
}

func empty(noun string) string {
	return "_No " + noun + " returned_"
}

// decode unwraps the UniFi response envelope ({"meta":..., "data":[...]})
// into out.
func decode(raw []byte, out any) error {
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return err
	}
	if len(envelope.Data) == 0 {
		return nil
	}
	return json.Unmarshal(envelope.Data, out)
}

func putString(attrs map[string]any, key, value string) {
	if value != "" {
		attrs[key] = value
	}
}

// buildSiteTable renders /api/self/sites: every site the account can read,
// which is how an operator finds the internal name to configure.
func buildSiteTable(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var sites []struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
		Role string `json:"role"`
	}
	if err := decode(raw, &sites); err != nil {
		return malformed("Sites", err), nil, nil
	}
	if len(sites) == 0 {
		return empty("sites"), nil, nil
	}

	var b strings.Builder
	b.WriteString("| Site | Name | Role |\n")
	b.WriteString("|------|------|------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(sites))
	for _, s := range sites {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s |\n", cell(s.Desc), cell(s.Name), cell(s.Role))
		attrs := map[string]any{}
		putString(attrs, "role", s.Role)
		putString(attrs, "displayName", s.Desc)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "site", Name: s.Name, ExternalID: s.Name, Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

// buildDeviceTable renders /api/s/<site>/stat/device. The endpoint returns
// per-device counters (uptime, tx/rx bytes, load) alongside the inventory;
// only the inventory fields are decoded so an idle controller keeps
// producing an identical section.
func buildDeviceTable(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var devices []struct {
		MAC      string  `json:"mac"`
		Name     string  `json:"name"`
		Model    string  `json:"model"`
		Type     string  `json:"type"`
		Version  string  `json:"version"`
		IP       string  `json:"ip"`
		Adopted  bool    `json:"adopted"`
		Disabled bool    `json:"disabled"`
		State    flexInt `json:"state"`
	}
	if err := decode(raw, &devices); err != nil {
		return malformed("Devices", err), nil, nil
	}
	if len(devices) == 0 {
		return empty("devices"), nil, nil
	}

	var b strings.Builder
	b.WriteString("| Device | Model | Type | IP | Firmware | Adopted | State |\n")
	b.WriteString("|--------|-------|------|----|----------|---------|-------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(devices))
	for _, d := range devices {
		state := deviceState(d.State)
		name := d.Name
		if name == "" {
			name = d.MAC
		}
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n",
			cell(name), cell(d.Model), cell(d.Type), cell(d.IP), cell(d.Version),
			yesNo(d.Adopted), cell(state))

		attrs := map[string]any{"adopted": d.Adopted, "disabled": d.Disabled}
		putString(attrs, "model", d.Model)
		putString(attrs, "deviceType", d.Type)
		putString(attrs, "firmwareVersion", d.Version)
		putString(attrs, "state", state)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "device", Name: name, IP: d.IP, Hostname: d.Name, ExternalID: d.MAC, Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

func deviceState(state flexInt) string {
	if !state.Set {
		return ""
	}
	if label, ok := deviceStates[state.Value]; ok {
		return label
	}
	return "state-" + strconv.Itoa(state.Value)
}

// buildNetworkTable renders /api/s/<site>/rest/networkconf: the configured
// LANs, VLANs and WAN interfaces.
func buildNetworkTable(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var networks []struct {
		ID           string  `json:"_id"`
		Name         string  `json:"name"`
		Purpose      string  `json:"purpose"`
		NetworkGroup string  `json:"networkgroup"`
		VLANEnabled  bool    `json:"vlan_enabled"`
		VLAN         flexInt `json:"vlan"`
		Subnet       string  `json:"ip_subnet"`
		DHCPEnabled  bool    `json:"dhcpd_enabled"`
		Enabled      *bool   `json:"enabled"`
	}
	if err := decode(raw, &networks); err != nil {
		return malformed("Networks", err), nil, nil
	}
	if len(networks) == 0 {
		return empty("networks"), nil, nil
	}

	var b strings.Builder
	b.WriteString("| Network | Purpose | VLAN | Subnet | DHCP | Enabled |\n")
	b.WriteString("|---------|---------|------|--------|------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(networks))
	for _, n := range networks {
		vlan := "—"
		if n.VLAN.Set {
			vlan = strconv.Itoa(n.VLAN.Value)
		}
		enabled := boolOr(n.Enabled, true)
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n",
			cell(n.Name), cell(n.Purpose), vlan, cell(n.Subnet), yesNo(n.DHCPEnabled), yesNo(enabled))

		attrs := map[string]any{
			"enabled":     enabled,
			"vlanEnabled": n.VLANEnabled,
			"dhcpEnabled": n.DHCPEnabled,
		}
		putString(attrs, "purpose", n.Purpose)
		putString(attrs, "subnet", n.Subnet)
		putString(attrs, "networkGroup", n.NetworkGroup)
		if n.VLAN.Set {
			attrs["vlan"] = n.VLAN.Value
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind: "network", Name: n.Name, ExternalID: n.ID, Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

// buildWLANTable renders /api/s/<site>/rest/wlanconf. The payload carries
// the pre-shared key (x_passphrase); it is never decoded or rendered.
func buildWLANTable(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var wlans []struct {
		ID               string `json:"_id"`
		Name             string `json:"name"`
		Enabled          *bool  `json:"enabled"`
		Security         string `json:"security"`
		WPAMode          string `json:"wpa_mode"`
		IsGuest          bool   `json:"is_guest"`
		HideSSID         bool   `json:"hide_ssid"`
		MACFilterEnabled bool   `json:"mac_filter_enabled"`
		MACFilterPolicy  string `json:"mac_filter_policy"`
		PMFMode          string `json:"pmf_mode"`
	}
	if err := decode(raw, &wlans); err != nil {
		return malformed("WLANs", err), nil, nil
	}
	if len(wlans) == 0 {
		return empty("WLANs"), nil, nil
	}

	var b strings.Builder
	b.WriteString("| SSID | Security | WPA Mode | Guest | Hidden | MAC Filter | Enabled |\n")
	b.WriteString("|------|----------|----------|-------|--------|------------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(wlans))
	for _, w := range wlans {
		enabled := boolOr(w.Enabled, true)
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n",
			cell(w.Name), cell(w.Security), cell(w.WPAMode), yesNo(w.IsGuest),
			yesNo(w.HideSSID), yesNo(w.MACFilterEnabled), yesNo(enabled))

		attrs := map[string]any{
			"enabled":          enabled,
			"guest":            w.IsGuest,
			"hideSsid":         w.HideSSID,
			"macFilterEnabled": w.MACFilterEnabled,
		}
		putString(attrs, "security", w.Security)
		putString(attrs, "wpaMode", w.WPAMode)
		putString(attrs, "macFilterPolicy", w.MACFilterPolicy)
		putString(attrs, "pmfMode", w.PMFMode)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "wlan", Name: w.Name, ExternalID: w.ID, Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

// buildFirewallTable renders /api/s/<site>/rest/firewallrule, sorted by
// ruleset and evaluation order so the section reads like the controller's
// own rule list regardless of the order the API returns.
func buildFirewallTable(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var rules []struct {
		ID         string  `json:"_id"`
		Name       string  `json:"name"`
		Ruleset    string  `json:"ruleset"`
		RuleIndex  flexInt `json:"rule_index"`
		Action     string  `json:"action"`
		Protocol   string  `json:"protocol"`
		Enabled    *bool   `json:"enabled"`
		Logging    bool    `json:"logging"`
		SrcAddress string  `json:"src_address"`
		DstAddress string  `json:"dst_address"`
	}
	if err := decode(raw, &rules); err != nil {
		return malformed("Firewall Rules", err), nil, nil
	}
	if len(rules) == 0 {
		return empty("firewall rules"), nil, nil
	}
	sort.SliceStable(rules, func(i, j int) bool {
		if rules[i].Ruleset != rules[j].Ruleset {
			return rules[i].Ruleset < rules[j].Ruleset
		}
		return rules[i].RuleIndex.Value < rules[j].RuleIndex.Value
	})

	var b strings.Builder
	b.WriteString("| Ruleset | # | Rule | Action | Protocol | Source | Destination | Enabled |\n")
	b.WriteString("|---------|---|------|--------|----------|--------|-------------|---------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(rules))
	for _, r := range rules {
		index := "—"
		if r.RuleIndex.Set {
			index = strconv.Itoa(r.RuleIndex.Value)
		}
		enabled := boolOr(r.Enabled, true)
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			cell(r.Ruleset), index, cell(r.Name), cell(r.Action), cell(r.Protocol),
			cell(r.SrcAddress), cell(r.DstAddress), yesNo(enabled))

		attrs := map[string]any{"enabled": enabled, "logging": r.Logging}
		putString(attrs, "action", r.Action)
		putString(attrs, "ruleset", r.Ruleset)
		putString(attrs, "protocol", r.Protocol)
		putString(attrs, "sourceAddress", r.SrcAddress)
		putString(attrs, "destinationAddress", r.DstAddress)
		if r.RuleIndex.Set {
			attrs["ruleIndex"] = r.RuleIndex.Value
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind: "firewall_rule", Name: r.Name, ExternalID: r.ID, Attributes: attrs,
		})
	}
	return b.String(), entities, nil
}

// buildClientSummary renders /api/s/<site>/stat/sta as counts only. The
// endpoint is the one genuinely volatile source on the controller (per
// client signal, traffic counters and last_seen), so no client entities are
// emitted and no per-client row is rendered: the section reports how many
// clients sit on each SSID or wired network.
func buildClientSummary(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	var clients []struct {
		IsWired bool   `json:"is_wired"`
		ESSID   string `json:"essid"`
		Network string `json:"network"`
	}
	if err := decode(raw, &clients); err != nil {
		return malformed("Clients", err), nil, nil
	}
	if len(clients) == 0 {
		return empty("clients"), nil, map[string]string{
			"client_count": "0", "wired_client_count": "0", "wireless_client_count": "0",
		}
	}

	wired := 0
	byConnection := map[string]int{}
	for _, cl := range clients {
		group := cl.ESSID
		if cl.IsWired {
			wired++
			group = cl.Network
		}
		if group == "" {
			group = "(unknown)"
		}
		byConnection[group]++
	}
	groups := make([]string, 0, len(byConnection))
	for g := range byConnection {
		groups = append(groups, g)
	}
	sort.Strings(groups)

	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "%d clients: %d wireless, %d wired\n\n", len(clients), len(clients)-wired, wired)
	b.WriteString("| SSID / Network | Clients |\n")
	b.WriteString("|----------------|---------|\n")
	for _, g := range groups {
		_, _ = fmt.Fprintf(&b, "| %s | %d |\n", cell(g), byConnection[g])
	}

	metadata := map[string]string{
		"client_count":          strconv.Itoa(len(clients)),
		"wired_client_count":    strconv.Itoa(wired),
		"wireless_client_count": strconv.Itoa(len(clients) - wired),
	}
	return b.String(), nil, metadata
}
