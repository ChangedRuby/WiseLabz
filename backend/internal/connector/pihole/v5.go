package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// Pi-hole v5 paths. v5 has no REST API: the admin UI drives PHP endpoints
// that accept the API token directly, which is what third-party integrations
// have always used. groups.php serves the group-management tables (groups,
// clients, adlists, domain allow/deny lists), customdns.php serves the local
// DNS records, and api.php serves the blocking on/off toggle.
const (
	v5GroupsPHP    = "/admin/scripts/pi-hole/php/groups.php"
	v5CustomDNSPHP = "/admin/scripts/pi-hole/php/customdns.php"
	v5APIPHP       = "/admin/api.php"
)

// v5Actions maps a normalised resource name to its groups.php action.
var v5Actions = map[string]string{
	resourceGroups:  "get_groups",
	resourceLists:   "get_adlists",
	resourceClients: "get_clients",
	resourceDomains: "get_domains",
}

// v5URL builds the absolute URL for a normalised resource on a v5 instance.
func (c *Connector) v5URL(resource string) (string, error) {
	if resource == resourceHosts {
		return c.url + v5CustomDNSPHP + "?action=get&auth=" + url.QueryEscape(c.password), nil
	}
	action, ok := v5Actions[resource]
	if !ok {
		return "", fmt.Errorf("unsupported v5 resource %q", resource)
	}
	return c.url + v5GroupsPHP + "?action=" + action + "&token=" + url.QueryEscape(c.password), nil
}

// doRequestV5 issues a GET against a v5 PHP endpoint and applies the same
// error mapping as the v6 path. v5 answers an unauthenticated read with HTTP
// 200 and an empty JSON array, so that shape is mapped to an AuthError
// rather than silently rendering as "no data".
func (c *Connector) doRequestV5(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	data, err := c.send(req)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(data)) == "[]" {
		return nil, connector.NewAuthError(fmt.Errorf("v5 API rejected the token (empty response)"))
	}
	return data, nil
}

// pingV5 probes whether the configured host speaks the v5 API with the
// configured token, used by automatic version detection.
func (c *Connector) pingV5(ctx context.Context) error {
	target := c.url + v5APIPHP + "?summaryRaw&auth=" + url.QueryEscape(c.password)
	data, err := c.doRequestV5(ctx, target)
	if err != nil {
		return err
	}
	var summary map[string]json.RawMessage
	if err := json.Unmarshal(data, &summary); err != nil {
		return connector.NewMalformedResponseError(err)
	}
	if len(summary) == 0 {
		return connector.NewAuthError(fmt.Errorf("v5 API returned an empty summary"))
	}
	return nil
}

// v5Bool decodes the 0/1, "0"/"1" and true/false spellings v5 uses for its
// boolean columns depending on the endpoint and PHP/SQLite coercion.
type v5Bool bool

// UnmarshalJSON implements json.Unmarshaler for v5Bool.
func (b *v5Bool) UnmarshalJSON(data []byte) error {
	s := strings.Trim(strings.TrimSpace(string(data)), `"`)
	switch s {
	case "true", "1":
		*b = true
	case "false", "0", "", "null":
		*b = false
	default:
		return fmt.Errorf("unexpected boolean value %q", s)
	}
	return nil
}

// --- v5 parsers -------------------------------------------------------------

func parseGroupsV5(raw []byte) ([]groupRow, error) {
	var resp struct {
		Data []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Enabled     v5Bool `json:"enabled"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]groupRow, 0, len(resp.Data))
	for _, g := range resp.Data {
		rows = append(rows, groupRow{ID: g.ID, Name: g.Name, Comment: g.Description, Enabled: bool(g.Enabled)})
	}
	return rows, nil
}

func parseAdlistsV5(raw []byte, names map[int]string) ([]adlistRow, error) {
	var resp struct {
		Data []struct {
			Address string `json:"address"`
			Comment string `json:"comment"`
			Enabled v5Bool `json:"enabled"`
			Groups  []int  `json:"groups"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]adlistRow, 0, len(resp.Data))
	for _, l := range resp.Data {
		// v5 adlists are block-only; the allow/deny split arrived in v6.
		rows = append(rows, adlistRow{
			Address: l.Address, Kind: "block", Comment: l.Comment,
			Enabled: bool(l.Enabled), Groups: resolveGroups(l.Groups, names),
		})
	}
	return rows, nil
}

func parseClientsV5(raw []byte, names map[int]string) ([]clientRow, error) {
	var resp struct {
		Data []struct {
			IP      string `json:"ip"`
			Comment string `json:"comment"`
			Groups  []int  `json:"groups"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]clientRow, 0, len(resp.Data))
	for _, cl := range resp.Data {
		rows = append(rows, clientRow{Client: cl.IP, Comment: cl.Comment, Groups: resolveGroups(cl.Groups, names)})
	}
	return rows, nil
}

// v5DomainType decodes the domainlist.type column: 0 exact allow, 1 exact
// deny, 2 regex allow, 3 regex deny.
func v5DomainType(t int) (rule, match string) {
	switch t {
	case 0:
		return "allow", "exact"
	case 1:
		return "deny", "exact"
	case 2:
		return "allow", "regex"
	case 3:
		return "deny", "regex"
	default:
		return "unknown", "unknown"
	}
}

func parseDomainsV5(raw []byte, names map[int]string) ([]domainRow, error) {
	var resp struct {
		Data []struct {
			Domain  string `json:"domain"`
			Type    int    `json:"type"`
			Comment string `json:"comment"`
			Enabled v5Bool `json:"enabled"`
			Groups  []int  `json:"groups"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, connector.NewMalformedResponseError(err)
	}
	rows := make([]domainRow, 0, len(resp.Data))
	for _, d := range resp.Data {
		rule, match := v5DomainType(d.Type)
		rows = append(rows, domainRow{
			Domain: d.Domain, Rule: rule, Match: match, Comment: d.Comment,
			Enabled: bool(d.Enabled), Groups: resolveGroups(d.Groups, names),
		})
	}
	return rows, nil
}

// buildHostsTableV5 renders customdns.php's `{"data":[["ip","hostname"]]}`
// payload into the same table and dns_record entities buildHostsTable emits
// for v6, so local DNS records diff identically across versions.
func buildHostsTableV5(raw []byte) (content string, entities []connector.SnapshotEntity) {
	var resp struct {
		Data [][]string `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "_Local DNS records unavailable: " + connector.NewMalformedResponseError(err).Error() + "_", nil
	}
	hosts := make([]string, 0, len(resp.Data))
	for _, pair := range resp.Data {
		if len(pair) < 2 {
			continue
		}
		hosts = append(hosts, pair[0]+" "+pair[1])
	}
	normalised, err := json.Marshal(map[string]any{"config": map[string]any{"dns": map[string]any{"hosts": hosts}}})
	if err != nil {
		return "_Local DNS records unavailable: " + err.Error() + "_", nil
	}
	return buildHostsTable(normalised)
}

// setBlockingV5 toggles blocking via api.php's enable/disable actions, the
// v5 equivalent of v6's POST /api/dns/blocking.
func (c *Connector) setBlockingV5(ctx context.Context, blocking bool) error {
	action := "disable"
	if blocking {
		action = "enable"
	}
	_, err := c.doRequestV5(ctx, c.url+v5APIPHP+"?"+action+"&auth="+url.QueryEscape(c.password))
	return err
}

// hostsItemV5 adds or removes a local DNS record via customdns.php.
func (c *Connector) hostsItemV5(ctx context.Context, action, ip, hostname string) error {
	query := url.Values{
		"action": {action},
		"ip":     {ip},
		"domain": {hostname},
		"auth":   {c.password},
	}
	_, err := c.doRequestV5(ctx, c.url+v5CustomDNSPHP+"?"+query.Encode())
	return err
}
