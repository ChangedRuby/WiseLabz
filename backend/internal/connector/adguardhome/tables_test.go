package adguardhome

import (
	"fmt"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TestBuildersOnMalformedAndEmptyInput checks every builder degrades to a
// placeholder instead of panicking or emitting half a table.
func TestBuildersOnMalformedAndEmptyInput(t *testing.T) {
	tests := []struct {
		name      string
		call      func([]byte) (string, int)
		empty     string
		wantEmpty string
	}{
		{
			name: "status",
			call: func(b []byte) (string, int) {
				c, _ := buildStatus(b)
				return c, 0
			},
			empty: `{}`,
		},
		{
			name: "dns_info",
			call: func(b []byte) (string, int) {
				c, ups, _ := buildDNSInfo(b)
				return c, len(ups)
			},
			empty: `{}`,
		},
		{
			name: "filtering",
			call: func(b []byte) (string, int) {
				c, _, e, _ := buildFiltering(b)
				return c, len(e)
			},
			empty:     `{}`,
			wantEmpty: "_No filter lists configured_",
		},
		{
			name: "rewrites",
			call: func(b []byte) (string, int) {
				c, e := buildRewriteTable(b)
				return c, len(e)
			},
			empty:     `[]`,
			wantEmpty: "_No DNS rewrites configured_",
		},
		{
			name: "clients",
			call: func(b []byte) (string, int) {
				c, e := buildClientTable(b)
				return c, len(e)
			},
			empty:     `{"clients":[]}`,
			wantEmpty: "_No persistent clients configured_",
		},
		{
			name: "dhcp",
			call: func(b []byte) (string, int) {
				c, e, _ := buildDHCP(b)
				return c, len(e)
			},
			empty:     `{}`,
			wantEmpty: "_No DHCP leases_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" empty", func(t *testing.T) {
			content, n := tt.call([]byte(tt.empty))
			if n != 0 {
				t.Fatalf("entities = %d, want 0", n)
			}
			if tt.wantEmpty != "" && !strings.Contains(content, tt.wantEmpty) {
				t.Fatalf("content = %q, want it to contain %q", content, tt.wantEmpty)
			}
		})
		t.Run(tt.name+" invalid JSON", func(t *testing.T) {
			content, n := tt.call([]byte(`not json`))
			if !strings.Contains(content, "malformed response") || n != 0 {
				t.Fatalf("content = %q (%d entities), want a malformed-response placeholder", content, n)
			}
		})
	}
}

// TestBuildStatusMetadata checks the status fields the rest of Fetch relies
// on are parsed.
func TestBuildStatusMetadata(t *testing.T) {
	content, info := buildStatus([]byte(statusJSON))
	if !info.dhcpAvailable {
		t.Error("dhcpAvailable = false, want true")
	}
	want := map[string]string{
		"adguard_version":    "v0.107.52",
		"protection_enabled": "true",
		"running":            "true",
		"dns_port":           "53",
		"dhcp_available":     "true",
	}
	for k, v := range want {
		if info.metadata[k] != v {
			t.Errorf("metadata[%q] = %q, want %q", k, info.metadata[k], v)
		}
	}
	for _, substr := range []string{"| Version | v0.107.52 |", "| DNS addresses | 192.168.1.2 |", "| Language | en |"} {
		if !strings.Contains(content, substr) {
			t.Errorf("content missing %q:\n%s", substr, content)
		}
	}
}

// TestBuildStatusOmitsVersionWhenAbsent keeps an empty version out of the
// metadata rather than recording it as "".
func TestBuildStatusOmitsVersionWhenAbsent(t *testing.T) {
	_, info := buildStatus([]byte(`{"running":true}`))
	if _, ok := info.metadata["adguard_version"]; ok {
		t.Errorf("metadata = %+v, want no adguard_version", info.metadata)
	}
}

func TestBuildDNSInfo(t *testing.T) {
	content, upstreams, metadata := buildDNSInfo([]byte(dnsInfoJSON))
	wantUpstreams := []string{"https://dns.quad9.net/dns-query", "tls://1.1.1.1"}
	if len(upstreams) != len(wantUpstreams) {
		t.Fatalf("upstreams = %v, want %v", upstreams, wantUpstreams)
	}
	for i, u := range wantUpstreams {
		if upstreams[i] != u {
			t.Errorf("upstreams[%d] = %q, want %q", i, upstreams[i], u)
		}
	}
	if metadata["upstream_count"] != "2" || metadata["dnssec_enabled"] != "true" || metadata["blocking_mode"] != "custom_ip" {
		t.Errorf("metadata = %+v", metadata)
	}
	// blocking_ipv4/6 only make sense in custom_ip mode.
	if !strings.Contains(content, "| Blocking IPv4 | 0.0.0.0 |") {
		t.Errorf("content missing the custom blocking IP:\n%s", content)
	}
	other, _, _ := buildDNSInfo([]byte(`{"blocking_mode":"nxdomain"}`))
	if strings.Contains(other, "Blocking IPv4") {
		t.Errorf("content should omit blocking IPs outside custom_ip mode:\n%s", other)
	}
}

func TestBuildFiltering(t *testing.T) {
	lists, rules, entities, metadata := buildFiltering([]byte(filteringJSON))
	if len(entities) != 3 {
		t.Fatalf("entities = %d, want 3", len(entities))
	}
	byName := map[string]connector.SnapshotEntity{}
	for _, e := range entities {
		if e.Kind != "filter_list" {
			t.Errorf("entity %q Kind = %q, want filter_list", e.Name, e.Kind)
		}
		byName[e.Name] = e
	}
	block := byName["AdGuard DNS filter"]
	if block.ExternalID != "1" || block.Attributes["kind"] != "blocklist" || block.Attributes["enabled"] != true {
		t.Errorf("blocklist entity = %+v", block)
	}
	if byName["Local allowlist"].Attributes["kind"] != "allowlist" {
		t.Errorf("allowlist entity = %+v", byName["Local allowlist"])
	}
	if byName["Dead hosts"].Attributes["enabled"] != false {
		t.Errorf("disabled entity = %+v", byName["Dead hosts"])
	}
	if !strings.Contains(lists, "| AdGuard DNS filter | blocklist | yes | 54321 |") {
		t.Errorf("lists content:\n%s", lists)
	}
	// The pipe characters of a rule must not break out of the code block.
	if !strings.Contains(rules, "||ads.example.com^") || !strings.Contains(rules, "@@||cdn.example.com^") {
		t.Errorf("rules content:\n%s", rules)
	}
	want := map[string]string{
		"filtering_enabled":    "true",
		"filter_list_count":    "3",
		"filter_list_enabled":  "2",
		"user_rule_count":      "2",
		"filter_update_hours":  "24",
		"blocklist_list_count": "2",
		"allowlist_list_count": "1",
	}
	for k, v := range want {
		if metadata[k] != v {
			t.Errorf("metadata[%q] = %q, want %q", k, metadata[k], v)
		}
	}
}

// TestBuildUserRulesTruncates keeps an imported rule set from dominating
// the snapshot.
func TestBuildUserRulesTruncates(t *testing.T) {
	rules := make([]string, maxRenderedRules+5)
	for i := range rules {
		rules[i] = fmt.Sprintf("||host%d.example.com^", i)
	}
	content, count := buildUserRules(rules)
	if count != len(rules) {
		t.Errorf("count = %d, want %d", count, len(rules))
	}
	if !strings.Contains(content, "_5 more rule(s) not shown_") {
		t.Errorf("content missing the truncation note:\n%s", content[len(content)-200:])
	}
	if strings.Contains(content, fmt.Sprintf("||host%d.example.com^", maxRenderedRules)) {
		t.Error("content contains a rule past the render cap")
	}
	empty, n := buildUserRules([]string{"", "  "})
	if n != 0 || empty != "_No custom filtering rules configured_" {
		t.Errorf("buildUserRules(blank) = %q (%d)", empty, n)
	}
}

func TestBuildRewriteTable(t *testing.T) {
	content, entities := buildRewriteTable([]byte(rewritesJSON))
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	host := entities[0]
	if host.Kind != "dns_rewrite" || host.Name != "nas.example.com" || host.Hostname != "nas.example.com" || host.IP != "192.168.1.10" {
		t.Errorf("rewrite entity = %+v", host)
	}
	if host.Attributes["wildcard"] != false || host.Attributes["answer"] != "192.168.1.10" {
		t.Errorf("rewrite attributes = %+v", host.Attributes)
	}
	wildcard := entities[1]
	if wildcard.Attributes["wildcard"] != true {
		t.Errorf("wildcard attributes = %+v", wildcard.Attributes)
	}
	// A wildcard pattern is not a hostname anything can be matched against.
	if wildcard.Hostname != "" || wildcard.IP != "" {
		t.Errorf("wildcard entity = %+v, want no hostname/IP", wildcard)
	}
	if wildcard.ExternalID != "*.lab.example.com=proxy.example.com" {
		t.Errorf("wildcard ExternalID = %q", wildcard.ExternalID)
	}
	if !strings.Contains(content, "| nas.example.com | 192.168.1.10 |") {
		t.Errorf("content:\n%s", content)
	}
}

func TestBuildClientTable(t *testing.T) {
	content, entities := buildClientTable([]byte(clientsJSON))
	// auto_clients are runtime discovery, not configuration.
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2 (auto_clients must be skipped)", len(entities))
	}
	laptop := entities[0]
	if laptop.Kind != "client" || laptop.Name != "laptop" || laptop.ExternalID != "laptop" {
		t.Errorf("client entity = %+v", laptop)
	}
	// The MAC comes first in ids, so the IP must still be picked out.
	if laptop.IP != "192.168.1.50" {
		t.Errorf("client IP = %q, want 192.168.1.50", laptop.IP)
	}
	if laptop.Attributes["useGlobalSettings"] != false || laptop.Attributes["filteringEnabled"] != true ||
		laptop.Attributes["safebrowsingEnabled"] != true || laptop.Attributes["parentalEnabled"] != false {
		t.Errorf("client attributes = %+v", laptop.Attributes)
	}
	if _, ok := entities[1].Attributes["tags"]; ok {
		t.Errorf("client without tags should omit the attribute: %+v", entities[1].Attributes)
	}
	if !strings.Contains(content, "| laptop | aa:bb:cc:dd:ee:ff, 192.168.1.50 | no | yes | facebook |") {
		t.Errorf("content:\n%s", content)
	}
	if strings.Contains(content, "phone") {
		t.Errorf("content lists an auto-client:\n%s", content)
	}
}

func TestBuildDHCP(t *testing.T) {
	content, entities, metadata := buildDHCP([]byte(dhcpJSON))
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	static := entities[0]
	if static.Kind != "dhcp_lease" || static.Name != "laptop" || static.IP != "192.168.1.50" ||
		static.Hostname != "laptop" || static.ExternalID != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("static lease entity = %+v", static)
	}
	if static.Attributes["static"] != true || static.Attributes["mac"] != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("static lease attributes = %+v", static.Attributes)
	}
	if entities[1].Attributes["static"] != false {
		t.Errorf("dynamic lease attributes = %+v", entities[1].Attributes)
	}
	// Expiry changes on every renewal and must stay out of the snapshot.
	for _, e := range entities {
		if _, ok := e.Attributes["expires"]; ok {
			t.Errorf("lease %q carries a volatile expiry attribute", e.Name)
		}
	}
	if !strings.Contains(content, "| IPv4 range | 192.168.1.100 – 192.168.1.200 |") {
		t.Errorf("content:\n%s", content)
	}
	want := map[string]string{"dhcp_enabled": "true", "dhcp_lease_count": "1", "dhcp_static_leases": "1"}
	for k, v := range want {
		if metadata[k] != v {
			t.Errorf("metadata[%q] = %q, want %q", k, metadata[k], v)
		}
	}
}

// TestLeaseNameFallsBackToAddress keeps an unnamed lease identifiable.
func TestLeaseNameFallsBackToAddress(t *testing.T) {
	_, entities, _ := buildDHCP([]byte(`{"leases":[{"mac":"aa:bb:cc:dd:ee:ff","ip":"10.0.0.5"},{"mac":"11:22:33:44:55:66"}]}`))
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	if entities[0].Name != "10.0.0.5" || entities[1].Name != "11:22:33:44:55:66" {
		t.Errorf("names = %q, %q", entities[0].Name, entities[1].Name)
	}
}

// TestCellEscapesPipes guards the Markdown tables against filtering syntax
// that contains '|'.
func TestCellEscapesPipes(t *testing.T) {
	if got := cell("||ads.example.com^"); got != `\|\|ads.example.com^` {
		t.Errorf("cell = %q", got)
	}
	if got := cell(""); got != "—" {
		t.Errorf("cell(\"\") = %q, want an em dash", got)
	}
	if got := cell("a\nb"); got != "a b" {
		t.Errorf("cell with newline = %q", got)
	}
}

// TestUpstreamDependenciesAreDedupedAndSorted keeps the snapshot stable
// across syncs regardless of the order the API returns upstreams in.
func TestUpstreamDependenciesAreDedupedAndSorted(t *testing.T) {
	deps := upstreamDependencies([]string{"9.9.9.9", "1.1.1.1", "9.9.9.9", ""})
	want := []connector.ServiceDependency{
		{Kind: "upstream_service", Name: "1.1.1.1"},
		{Kind: "upstream_service", Name: "9.9.9.9"},
	}
	if len(deps) != len(want) {
		t.Fatalf("deps = %+v, want %+v", deps, want)
	}
	for i := range want {
		if deps[i] != want[i] {
			t.Errorf("deps[%d] = %+v, want %+v", i, deps[i], want[i])
		}
	}
	if upstreamDependencies(nil) != nil {
		t.Error("upstreamDependencies(nil) should be nil")
	}
}
