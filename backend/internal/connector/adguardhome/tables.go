package adguardhome

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// maxRenderedRules caps how many custom filtering rules are rendered in the
// section body. A hand-maintained rule set is small, but an imported one can
// run to thousands of lines, which would bloat every snapshot and diff.
const maxRenderedRules = 200

// cell escapes a value for use inside a Markdown table cell. Filtering
// rules and upstream specs routinely contain '|' (as in "[/example.org/]"
// syntax or "|ads.example.com^") and would otherwise split the row.
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

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// statusInfo carries the parts of /control/status the rest of Fetch needs:
// the metadata to merge into the snapshot and whether this build even has
// DHCP support.
type statusInfo struct {
	metadata      map[string]string
	dhcpAvailable bool
}

// buildStatus renders /control/status and returns the derived status info.
func buildStatus(raw []byte) (string, statusInfo) {
	var status struct {
		Version           string   `json:"version"`
		Language          string   `json:"language"`
		DNSAddresses      []string `json:"dns_addresses"`
		DNSPort           int      `json:"dns_port"`
		HTTPPort          int      `json:"http_port"`
		ProtectionEnabled bool     `json:"protection_enabled"`
		Running           bool     `json:"running"`
		DHCPAvailable     bool     `json:"dhcp_available"`
	}
	if err := json.Unmarshal(raw, &status); err != nil {
		return malformed("Status", err), statusInfo{}
	}

	var b strings.Builder
	b.WriteString("| Setting | Value |\n")
	b.WriteString("|---------|-------|\n")
	_, _ = fmt.Fprintf(&b, "| Version | %s |\n", cell(status.Version))
	_, _ = fmt.Fprintf(&b, "| Running | %s |\n", yesNo(status.Running))
	_, _ = fmt.Fprintf(&b, "| Protection enabled | %s |\n", yesNo(status.ProtectionEnabled))
	_, _ = fmt.Fprintf(&b, "| DNS addresses | %s |\n", cell(strings.Join(status.DNSAddresses, ", ")))
	_, _ = fmt.Fprintf(&b, "| DNS port | %d |\n", status.DNSPort)
	_, _ = fmt.Fprintf(&b, "| Web port | %d |\n", status.HTTPPort)
	_, _ = fmt.Fprintf(&b, "| DHCP available | %s |\n", yesNo(status.DHCPAvailable))
	if status.Language != "" {
		_, _ = fmt.Fprintf(&b, "| Language | %s |\n", cell(status.Language))
	}

	metadata := map[string]string{
		"protection_enabled": fmt.Sprintf("%t", status.ProtectionEnabled),
		"running":            fmt.Sprintf("%t", status.Running),
		"dns_port":           fmt.Sprintf("%d", status.DNSPort),
		"dhcp_available":     fmt.Sprintf("%t", status.DHCPAvailable),
	}
	if status.Version != "" {
		metadata["adguard_version"] = status.Version
	}
	return b.String(), statusInfo{metadata: metadata, dhcpAvailable: status.DHCPAvailable}
}

// buildDNSInfo renders /control/dns_info. It returns the section content,
// the configured upstream resolvers (for snapshot dependencies) and count
// metadata.
func buildDNSInfo(raw []byte) (string, []string, map[string]string) {
	var info struct {
		UpstreamDNS       []string `json:"upstream_dns"`
		BootstrapDNS      []string `json:"bootstrap_dns"`
		FallbackDNS       []string `json:"fallback_dns"`
		LocalPTRUpstreams []string `json:"local_ptr_upstreams"`
		UpstreamMode      string   `json:"upstream_mode"`
		BlockingMode      string   `json:"blocking_mode"`
		BlockingIPv4      string   `json:"blocking_ipv4"`
		BlockingIPv6      string   `json:"blocking_ipv6"`
		RateLimit         int      `json:"ratelimit"`
		CacheSize         int      `json:"cache_size"`
		EDNSCSEnabled     bool     `json:"edns_cs_enabled"`
		DNSSECEnabled     bool     `json:"dnssec_enabled"`
		DisableIPv6       bool     `json:"disable_ipv6"`
		ResolveClients    bool     `json:"resolve_clients"`
		ProtectionEnabled bool     `json:"protection_enabled"`
	}
	if err := json.Unmarshal(raw, &info); err != nil {
		return malformed("DNS Configuration", err), nil, nil
	}

	var b strings.Builder
	b.WriteString("| Setting | Value |\n")
	b.WriteString("|---------|-------|\n")
	_, _ = fmt.Fprintf(&b, "| Upstream DNS | %s |\n", cell(strings.Join(info.UpstreamDNS, "<br>")))
	_, _ = fmt.Fprintf(&b, "| Bootstrap DNS | %s |\n", cell(strings.Join(info.BootstrapDNS, "<br>")))
	_, _ = fmt.Fprintf(&b, "| Fallback DNS | %s |\n", cell(strings.Join(info.FallbackDNS, "<br>")))
	_, _ = fmt.Fprintf(&b, "| Private PTR upstreams | %s |\n", cell(strings.Join(info.LocalPTRUpstreams, "<br>")))
	_, _ = fmt.Fprintf(&b, "| Upstream mode | %s |\n", cell(info.UpstreamMode))
	_, _ = fmt.Fprintf(&b, "| Blocking mode | %s |\n", cell(info.BlockingMode))
	if info.BlockingMode == "custom_ip" {
		_, _ = fmt.Fprintf(&b, "| Blocking IPv4 | %s |\n", cell(info.BlockingIPv4))
		_, _ = fmt.Fprintf(&b, "| Blocking IPv6 | %s |\n", cell(info.BlockingIPv6))
	}
	_, _ = fmt.Fprintf(&b, "| DNSSEC | %s |\n", yesNo(info.DNSSECEnabled))
	_, _ = fmt.Fprintf(&b, "| EDNS client subnet | %s |\n", yesNo(info.EDNSCSEnabled))
	_, _ = fmt.Fprintf(&b, "| IPv6 disabled | %s |\n", yesNo(info.DisableIPv6))
	_, _ = fmt.Fprintf(&b, "| Resolve clients | %s |\n", yesNo(info.ResolveClients))
	_, _ = fmt.Fprintf(&b, "| Rate limit | %d |\n", info.RateLimit)
	_, _ = fmt.Fprintf(&b, "| Cache size | %d |\n", info.CacheSize)

	metadata := map[string]string{
		"upstream_count": fmt.Sprintf("%d", len(info.UpstreamDNS)),
		"dnssec_enabled": fmt.Sprintf("%t", info.DNSSECEnabled),
	}
	if info.BlockingMode != "" {
		metadata["blocking_mode"] = info.BlockingMode
	}
	if info.UpstreamMode != "" {
		metadata["upstream_mode"] = info.UpstreamMode
	}
	return b.String(), info.UpstreamDNS, metadata
}

// buildFiltering renders /control/filtering/status into two sections: the
// block/allow lists (one "filter_list" entity each) and the custom user
// rules. It returns both bodies, the entities and count metadata.
func buildFiltering(raw []byte) (lists string, rules string, entities []connector.SnapshotEntity, metadata map[string]string) {
	type filter struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		URL        string `json:"url"`
		RulesCount int    `json:"rules_count"`
		Enabled    bool   `json:"enabled"`
	}
	var status struct {
		Enabled          bool     `json:"enabled"`
		Interval         int      `json:"interval"`
		Filters          []filter `json:"filters"`
		WhitelistFilters []filter `json:"whitelist_filters"`
		UserRules        []string `json:"user_rules"`
	}
	if err := json.Unmarshal(raw, &status); err != nil {
		return malformed("Filter Lists", err), malformed("Custom Filtering Rules", err), nil, nil
	}

	groups := []struct {
		kind    string
		filters []filter
	}{
		{"blocklist", status.Filters},
		{"allowlist", status.WhitelistFilters},
	}
	total := len(status.Filters) + len(status.WhitelistFilters)
	enabledLists := 0
	var b strings.Builder
	if total == 0 {
		b.WriteString("_No filter lists configured_")
	} else {
		b.WriteString("| List | Kind | Enabled | Rules | URL |\n")
		b.WriteString("|------|------|---------|-------|-----|\n")
		entities = make([]connector.SnapshotEntity, 0, total)
		for _, g := range groups {
			for _, f := range g.filters {
				if f.Enabled {
					enabledLists++
				}
				_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %d | %s |\n",
					cell(f.Name), g.kind, yesNo(f.Enabled), f.RulesCount, cell(f.URL))
				attrs := map[string]any{"enabled": f.Enabled, "kind": g.kind}
				putString(attrs, "url", f.URL)
				entities = append(entities, connector.SnapshotEntity{
					Kind:       "filter_list",
					Name:       f.Name,
					ExternalID: fmt.Sprintf("%d", f.ID),
					Attributes: attrs,
				})
			}
		}
	}

	rulesContent, ruleCount := buildUserRules(status.UserRules)
	metadata = map[string]string{
		"filtering_enabled":    fmt.Sprintf("%t", status.Enabled),
		"filter_list_count":    fmt.Sprintf("%d", total),
		"filter_list_enabled":  fmt.Sprintf("%d", enabledLists),
		"user_rule_count":      fmt.Sprintf("%d", ruleCount),
		"filter_update_hours":  fmt.Sprintf("%d", status.Interval),
		"blocklist_list_count": fmt.Sprintf("%d", len(status.Filters)),
		"allowlist_list_count": fmt.Sprintf("%d", len(status.WhitelistFilters)),
	}
	return b.String(), rulesContent, entities, metadata
}

// buildUserRules renders the custom filtering rules as a fenced block,
// truncated at maxRenderedRules so an imported rule set can't dominate the
// snapshot. It also returns how many non-empty rules were configured.
func buildUserRules(userRules []string) (string, int) {
	rules := make([]string, 0, len(userRules))
	for _, r := range userRules {
		if strings.TrimSpace(r) != "" {
			rules = append(rules, r)
		}
	}
	if len(rules) == 0 {
		return "_No custom filtering rules configured_", 0
	}
	count := len(rules)
	truncated := 0
	if len(rules) > maxRenderedRules {
		truncated = len(rules) - maxRenderedRules
		rules = rules[:maxRenderedRules]
	}
	var b strings.Builder
	b.WriteString("```\n")
	for _, r := range rules {
		b.WriteString(r)
		b.WriteString("\n")
	}
	b.WriteString("```\n")
	if truncated > 0 {
		_, _ = fmt.Fprintf(&b, "\n_%d more rule(s) not shown_\n", truncated)
	}
	return b.String(), count
}

// buildRewriteTable renders /control/rewrite/list, one "dns_rewrite" entity
// per rewrite. Rewrites pointing at a literal address carry it as the
// entity IP so they cross-link with hosts and DNS records.
func buildRewriteTable(raw []byte) (string, []connector.SnapshotEntity) {
	var rewrites []struct {
		Domain string `json:"domain"`
		Answer string `json:"answer"`
	}
	if err := json.Unmarshal(raw, &rewrites); err != nil {
		return malformed("DNS Rewrites", err), nil
	}
	if len(rewrites) == 0 {
		return "_No DNS rewrites configured_", nil
	}

	var b strings.Builder
	b.WriteString("| Domain | Answer |\n")
	b.WriteString("|--------|--------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(rewrites))
	for _, r := range rewrites {
		_, _ = fmt.Fprintf(&b, "| %s | %s |\n", cell(r.Domain), cell(r.Answer))
		wildcard := strings.HasPrefix(r.Domain, "*.")
		attrs := map[string]any{"wildcard": wildcard}
		putString(attrs, "answer", r.Answer)
		entity := connector.SnapshotEntity{
			Kind:       "dns_rewrite",
			Name:       r.Domain,
			ExternalID: r.Domain + "=" + r.Answer,
			Attributes: attrs,
		}
		if !wildcard {
			entity.Hostname = r.Domain
		}
		if net.ParseIP(r.Answer) != nil {
			entity.IP = r.Answer
		}
		entities = append(entities, entity)
	}
	return b.String(), entities
}

// buildClientTable renders the persistent clients of /control/clients. The
// runtime "auto_clients" list is deliberately skipped: it is discovered
// traffic, not configuration, and would make every snapshot look changed.
func buildClientTable(raw []byte) (string, []connector.SnapshotEntity) {
	var payload struct {
		Clients []struct {
			Name                string   `json:"name"`
			IDs                 []string `json:"ids"`
			Tags                []string `json:"tags"`
			BlockedServices     []string `json:"blocked_services"`
			Upstreams           []string `json:"upstreams"`
			UseGlobalSettings   bool     `json:"use_global_settings"`
			FilteringEnabled    bool     `json:"filtering_enabled"`
			SafebrowsingEnabled bool     `json:"safebrowsing_enabled"`
			ParentalEnabled     bool     `json:"parental_enabled"`
		} `json:"clients"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return malformed("Clients", err), nil
	}
	if len(payload.Clients) == 0 {
		return "_No persistent clients configured_", nil
	}

	var b strings.Builder
	b.WriteString("| Client | Identifiers | Global Settings | Filtering | Blocked Services | Upstreams | Tags |\n")
	b.WriteString("|--------|-------------|-----------------|-----------|------------------|-----------|------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(payload.Clients))
	for _, cl := range payload.Clients {
		_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n",
			cell(cl.Name), cell(strings.Join(cl.IDs, ", ")), yesNo(cl.UseGlobalSettings),
			yesNo(cl.FilteringEnabled), cell(strings.Join(cl.BlockedServices, ", ")),
			cell(strings.Join(cl.Upstreams, "<br>")), cell(strings.Join(cl.Tags, ", ")))

		attrs := map[string]any{
			"useGlobalSettings":   cl.UseGlobalSettings,
			"filteringEnabled":    cl.FilteringEnabled,
			"safebrowsingEnabled": cl.SafebrowsingEnabled,
			"parentalEnabled":     cl.ParentalEnabled,
		}
		putStrings(attrs, "ids", cl.IDs)
		putStrings(attrs, "tags", cl.Tags)
		putStrings(attrs, "blockedServices", cl.BlockedServices)
		putStrings(attrs, "upstreams", cl.Upstreams)
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "client",
			Name:       cl.Name,
			ExternalID: cl.Name,
			IP:         firstIP(cl.IDs),
			Attributes: attrs,
		})
	}
	return b.String(), entities
}

// firstIP returns the first client identifier that is a literal IP, so
// clients cross-link with hosts, VMs and DHCP leases. CIDRs, MACs and
// ClientIDs yield "".
func firstIP(ids []string) string {
	for _, id := range ids {
		if net.ParseIP(id) != nil {
			return id
		}
	}
	return ""
}

// buildDHCP renders /control/dhcp/status: the server settings plus static
// and dynamic leases, one "dhcp_lease" entity each. Lease expiry is
// deliberately left out of the attributes — it changes on every renewal.
func buildDHCP(raw []byte) (string, []connector.SnapshotEntity, map[string]string) {
	type lease struct {
		MAC      string `json:"mac"`
		IP       string `json:"ip"`
		Hostname string `json:"hostname"`
	}
	var status struct {
		Enabled       bool   `json:"enabled"`
		InterfaceName string `json:"interface_name"`
		V4            struct {
			GatewayIP     string `json:"gateway_ip"`
			SubnetMask    string `json:"subnet_mask"`
			RangeStart    string `json:"range_start"`
			RangeEnd      string `json:"range_end"`
			LeaseDuration int    `json:"lease_duration"`
		} `json:"v4"`
		Leases       []lease `json:"leases"`
		StaticLeases []lease `json:"static_leases"`
	}
	if err := json.Unmarshal(raw, &status); err != nil {
		return malformed("DHCP", err), nil, nil
	}

	var b strings.Builder
	b.WriteString("| Setting | Value |\n")
	b.WriteString("|---------|-------|\n")
	_, _ = fmt.Fprintf(&b, "| Enabled | %s |\n", yesNo(status.Enabled))
	_, _ = fmt.Fprintf(&b, "| Interface | %s |\n", cell(status.InterfaceName))
	_, _ = fmt.Fprintf(&b, "| IPv4 range | %s |\n", cell(rangeText(status.V4.RangeStart, status.V4.RangeEnd)))
	_, _ = fmt.Fprintf(&b, "| Gateway | %s |\n", cell(status.V4.GatewayIP))
	_, _ = fmt.Fprintf(&b, "| Subnet mask | %s |\n", cell(status.V4.SubnetMask))
	_, _ = fmt.Fprintf(&b, "| Lease duration | %ds |\n", status.V4.LeaseDuration)

	groups := []struct {
		static bool
		leases []lease
	}{{true, status.StaticLeases}, {false, status.Leases}}
	var entities []connector.SnapshotEntity
	if len(status.Leases)+len(status.StaticLeases) == 0 {
		b.WriteString("\n_No DHCP leases_\n")
	} else {
		b.WriteString("\n| Hostname | IP | MAC | Static |\n")
		b.WriteString("|----------|----|-----|--------|\n")
		for _, g := range groups {
			for _, l := range g.leases {
				_, _ = fmt.Fprintf(&b, "| %s | %s | %s | %s |\n",
					cell(l.Hostname), cell(l.IP), cell(l.MAC), yesNo(g.static))
				attrs := map[string]any{"static": g.static}
				putString(attrs, "mac", l.MAC)
				entities = append(entities, connector.SnapshotEntity{
					Kind:       "dhcp_lease",
					Name:       leaseName(l.Hostname, l.IP, l.MAC),
					Hostname:   l.Hostname,
					IP:         l.IP,
					ExternalID: l.MAC,
					Attributes: attrs,
				})
			}
		}
	}

	metadata := map[string]string{
		"dhcp_enabled":       fmt.Sprintf("%t", status.Enabled),
		"dhcp_lease_count":   fmt.Sprintf("%d", len(status.Leases)),
		"dhcp_static_leases": fmt.Sprintf("%d", len(status.StaticLeases)),
	}
	return b.String(), entities, metadata
}

func rangeText(start, end string) string {
	if start == "" && end == "" {
		return ""
	}
	return start + " – " + end
}

// leaseName prefers the hostname a lease advertises, falling back to its
// address so an unnamed lease is still identifiable.
func leaseName(hostname, ip, mac string) string {
	switch {
	case hostname != "":
		return hostname
	case ip != "":
		return ip
	default:
		return mac
	}
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
