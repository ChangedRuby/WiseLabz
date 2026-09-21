package pihole

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// The v5 (api.php/groups.php) and v6 (REST) APIs return the same concepts in
// different shapes. Each parser normalises its version's payload into these
// rows so the table builders — and the entities they emit — stay identical
// across versions.

// groupRow is one Pi-hole group (the unit clients and lists are assigned to).
type groupRow struct {
	ID      int
	Name    string
	Comment string
	Enabled bool
}

// adlistRow is one subscribed blocklist/allowlist (adlist) URL.
type adlistRow struct {
	Address string
	Kind    string // "block" or "allow"
	Comment string
	Enabled bool
	Groups  []string
}

// clientRow is one configured client (IP, MAC, hostname or interface).
type clientRow struct {
	Client  string
	Comment string
	Groups  []string
}

// domainRow is one entry of the domain allow/deny lists.
type domainRow struct {
	Domain  string
	Rule    string // "allow" or "deny"
	Match   string // "exact" or "regex"
	Comment string
	Enabled bool
	Groups  []string
}

// groupNames indexes group rows by id so adlist/client/domain payloads —
// which reference groups by numeric id in both API versions — can render and
// expose group names instead of opaque numbers.
func groupNames(groups []groupRow) map[int]string {
	names := make(map[int]string, len(groups))
	for _, g := range groups {
		names[g.ID] = g.Name
	}
	return names
}

// resolveGroups maps group ids to names, falling back to the raw id when a
// group is unknown (deleted between calls, or a partial response).
func resolveGroups(ids []int, names map[int]string) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if name, ok := names[id]; ok && name != "" {
			out = append(out, name)
			continue
		}
		out = append(out, fmt.Sprintf("%d", id))
	}
	sort.Strings(out)
	return out
}

// cell escapes a value for use inside a markdown table cell.
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	if s == "" {
		return "-"
	}
	return s
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func buildGroupTable(rows []groupRow) (content string, entities []connector.SnapshotEntity) {
	if len(rows) == 0 {
		return "_No groups returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Group | Enabled | Comment |\n")
	b.WriteString("|-------|---------|---------|\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", cell(r.Name), yesNo(r.Enabled), cell(r.Comment))
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "dns_group",
			Name:       r.Name,
			ExternalID: r.Name,
			Attributes: map[string]any{
				"enabled": r.Enabled,
				"comment": r.Comment,
			},
		})
	}
	return b.String(), entities
}

func buildAdlistTable(rows []adlistRow) (content string, entities []connector.SnapshotEntity) {
	if len(rows) == 0 {
		return "_No blocklists returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Address | Type | Enabled | Groups | Comment |\n")
	b.WriteString("|---------|------|---------|--------|---------|\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
			cell(r.Address), cell(r.Kind), yesNo(r.Enabled), cell(strings.Join(r.Groups, ", ")), cell(r.Comment))
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "blocklist",
			Name:       r.Address,
			ExternalID: r.Address,
			Attributes: map[string]any{
				"enabled":   r.Enabled,
				"list_type": r.Kind,
				"groups":    r.Groups,
			},
		})
	}
	return b.String(), entities
}

func buildClientTable(rows []clientRow) (content string, entities []connector.SnapshotEntity) {
	if len(rows) == 0 {
		return "_No clients returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Client | Groups | Comment |\n")
	b.WriteString("|--------|--------|---------|\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", cell(r.Client), cell(strings.Join(r.Groups, ", ")), cell(r.Comment))
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "dns_client",
			Name:       r.Client,
			IP:         clientIP(r.Client),
			ExternalID: r.Client,
			Attributes: map[string]any{
				"groups":  r.Groups,
				"comment": r.Comment,
			},
		})
	}
	return b.String(), entities
}

func buildDomainTable(rows []domainRow) (content string, entities []connector.SnapshotEntity) {
	if len(rows) == 0 {
		return "_No domain rules returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Domain | Rule | Match | Enabled | Groups | Comment |\n")
	b.WriteString("|--------|------|-------|---------|--------|---------|\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n",
			cell(r.Domain), cell(r.Rule), cell(r.Match), yesNo(r.Enabled),
			cell(strings.Join(r.Groups, ", ")), cell(r.Comment))
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "domain_rule",
			Name:       r.Domain,
			ExternalID: r.Rule + "/" + r.Match + "/" + r.Domain,
			Attributes: map[string]any{
				"enabled":    r.Enabled,
				"rule_type":  r.Rule,
				"match_kind": r.Match,
				"groups":     r.Groups,
			},
		})
	}
	return b.String(), entities
}

// --- v6 parsers -------------------------------------------------------------

func parseGroupsV6(raw []byte) ([]groupRow, error) {
	var resp struct {
		Groups []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Comment string `json:"comment"`
			Enabled bool   `json:"enabled"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]groupRow, 0, len(resp.Groups))
	for _, g := range resp.Groups {
		rows = append(rows, groupRow{ID: g.ID, Name: g.Name, Comment: g.Comment, Enabled: g.Enabled})
	}
	return rows, nil
}

func parseAdlistsV6(raw []byte, names map[int]string) ([]adlistRow, error) {
	var resp struct {
		Lists []struct {
			Address string `json:"address"`
			Type    string `json:"type"`
			Comment string `json:"comment"`
			Enabled bool   `json:"enabled"`
			Groups  []int  `json:"groups"`
		} `json:"lists"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]adlistRow, 0, len(resp.Lists))
	for _, l := range resp.Lists {
		kind := l.Type
		if kind == "" {
			kind = "block"
		}
		rows = append(rows, adlistRow{
			Address: l.Address, Kind: kind, Comment: l.Comment,
			Enabled: l.Enabled, Groups: resolveGroups(l.Groups, names),
		})
	}
	return rows, nil
}

func parseClientsV6(raw []byte, names map[int]string) ([]clientRow, error) {
	var resp struct {
		Clients []struct {
			Client  string `json:"client"`
			Comment string `json:"comment"`
			Groups  []int  `json:"groups"`
		} `json:"clients"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]clientRow, 0, len(resp.Clients))
	for _, cl := range resp.Clients {
		rows = append(rows, clientRow{Client: cl.Client, Comment: cl.Comment, Groups: resolveGroups(cl.Groups, names)})
	}
	return rows, nil
}

func parseDomainsV6(raw []byte, names map[int]string) ([]domainRow, error) {
	var resp struct {
		Domains []struct {
			Domain  string `json:"domain"`
			Type    string `json:"type"`
			Kind    string `json:"kind"`
			Comment string `json:"comment"`
			Enabled bool   `json:"enabled"`
			Groups  []int  `json:"groups"`
		} `json:"domains"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]domainRow, 0, len(resp.Domains))
	for _, d := range resp.Domains {
		rows = append(rows, domainRow{
			Domain: d.Domain, Rule: d.Type, Match: d.Kind, Comment: d.Comment,
			Enabled: d.Enabled, Groups: resolveGroups(d.Groups, names),
		})
	}
	return rows, nil
}

// clientIP returns the client identifier when it is a bare IP address, so
// cross-connector entity matching can link Pi-hole clients to hosts seen by
// other connectors. MAC addresses, hostnames and interface names yield "".
func clientIP(client string) string {
	if net.ParseIP(client) != nil {
		return client
	}
	return ""
}
