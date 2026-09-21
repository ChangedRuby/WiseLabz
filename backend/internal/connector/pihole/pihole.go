// Package pihole implements a Pi-hole local DNS records connector,
// targeting the Pi-hole v6 REST API.
package pihole

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "pihole"

// API versions this connector speaks. versionAuto probes v6 first and falls
// back to v5, so an existing config keeps working unchanged.
const (
	versionAuto = "auto"
	version5    = "v5"
	version6    = "v6"
)

// Normalised resource names. Each version maps these to its own path (see
// (*Connector).get and v5URL).
const (
	resourceHosts   = "config/dns/hosts"
	resourceGroups  = "groups"
	resourceLists   = "lists"
	resourceClients = "clients"
	resourceDomains = "domains"
)

// session carries the resolved API version and, for v6, the session id
// returned by POST /api/auth. v5 authenticates per-request with the token.
type session struct {
	version string
	sid     string
}

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "dns",
		Name:     "Pi-hole",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Pi-hole URL", Type: "text", Required: true, Placeholder: "https://pihole.example.com"},
			{Key: "password", Label: "Password / App Password (v6) or API Token (v5)", Type: "password", Required: true},
			{Key: "api_version", Label: "API Version", Type: "select", Required: false, Default: versionAuto, Options: []string{versionAuto, version6, version5}},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, func(config map[string]any) (connector.Connector, error) {
		url, _ := config["url"].(string)
		password, _ := config["password"].(string)
		apiVersion, _ := config["api_version"].(string)
		switch apiVersion {
		case version5, version6:
		default:
			apiVersion = versionAuto
		}
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
		return &Connector{
			url:        strings.TrimSuffix(url, "/"),
			password:   password,
			apiVersion: apiVersion,
			client:     client,
		}, nil
	})
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills on
// its entities (see buildHostsTable and the table builders in lists.go),
// exposed via GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"dns_record": {
		{Name: "source", Type: "string", Description: "Source of the DNS record (local_dns, dhcp_lease, etc.)"},
		{Name: "is_ipv6", Type: "boolean", Description: "Whether the record resolves to an IPv6 address"},
	},
	"blocklist": {
		{Name: "enabled", Type: "boolean", Description: "Whether the adlist is active"},
		{Name: "list_type", Type: "string", Description: "Whether the adlist blocks or allows (block, allow)"},
		{Name: "groups", Type: "string_array", Description: "Group names the adlist applies to"},
	},
	"dns_group": {
		{Name: "enabled", Type: "boolean", Description: "Whether the group is active"},
		{Name: "comment", Type: "string", Description: "Free-text description of the group"},
	},
	"dns_client": {
		{Name: "groups", Type: "string_array", Description: "Group names the client belongs to"},
		{Name: "comment", Type: "string", Description: "Free-text description of the client"},
	},
	"domain_rule": {
		{Name: "enabled", Type: "boolean", Description: "Whether the domain rule is active"},
		{Name: "rule_type", Type: "string", Description: "Whether the rule allows or denies the domain (allow, deny)"},
		{Name: "match_kind", Type: "string", Description: "How the domain is matched (exact, regex)"},
		{Name: "groups", Type: "string_array", Description: "Group names the rule applies to"},
	},
}

// Connector fetches DNS and blocking configuration from a Pi-hole instance.
type Connector struct {
	url        string
	password   string
	apiVersion string // auto, v5 or v6 as configured
	client     *http.Client
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Pi-hole" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "dns" }

// Validate tests the connection to the Pi-hole API, resolving the API
// version first so a v5 instance validates against its own endpoints.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	s, err := c.session(ctx)
	if err != nil {
		return err
	}
	_, err = c.get(ctx, s, resourceHosts)
	return err
}

// session resolves the API version to use for this call. With api_version
// pinned it authenticates against that version only; on "auto" it tries the
// v6 session endpoint first and falls back to v5's token auth, reporting the
// v6 error when neither answers.
func (c *Connector) session(ctx context.Context) (session, error) {
	switch c.apiVersion {
	case version5:
		return session{version: version5}, c.pingV5(ctx)
	case version6:
		sid, err := c.authenticate(ctx)
		return session{version: version6, sid: sid}, err
	default:
		sid, err := c.authenticate(ctx)
		if err == nil {
			return session{version: version6, sid: sid}, nil
		}
		if v5Err := c.pingV5(ctx); v5Err == nil {
			return session{version: version5}, nil
		}
		return session{}, err
	}
}

// get fetches a normalised resource using whichever API the session resolved.
func (c *Connector) get(ctx context.Context, s session, resource string) ([]byte, error) {
	if s.version == version5 {
		target, err := c.v5URL(resource)
		if err != nil {
			return nil, err
		}
		return c.doRequestV5(ctx, target)
	}
	return c.doRequest(ctx, s.sid, "/api/"+resource)
}

// Fetch retrieves local DNS records, blocklists, groups, clients and the
// domain allow/deny lists. Each section degrades independently: an endpoint
// the instance does not expose renders as unavailable instead of failing the
// whole snapshot.
func (c *Connector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	metadata := map[string]string{"pihole_url": c.url}

	s, err := c.session(ctx)
	if err != nil {
		return &connector.ServiceSnapshot{
			ServiceName: "Pi-hole",
			Type:        typeName,
			Sections:    []connector.SnapshotSection{unavailable("Local DNS Records", err)},
			Metadata:    metadata,
			FetchedAt:   start,
		}, nil
	}
	metadata["pihole_api_version"] = s.version

	section, entities := c.fetchHosts(ctx, s)
	sections := []connector.SnapshotSection{section}

	groups, groupSection, groupEntities := c.fetchGroups(ctx, s)
	sections = append(sections, groupSection)
	entities = append(entities, groupEntities...)
	names := groupNames(groups)

	for _, fetch := range []func(context.Context, session, map[int]string) (connector.SnapshotSection, []connector.SnapshotEntity){
		c.fetchAdlists, c.fetchClients, c.fetchDomains,
	} {
		sec, ents := fetch(ctx, s, names)
		sections = append(sections, sec)
		entities = append(entities, ents...)
	}

	return &connector.ServiceSnapshot{
		ServiceName: "Pi-hole",
		Type:        typeName,
		Sections:    sections,
		Entities:    entities,
		Metadata:    metadata,
		FetchedAt:   start,
	}, nil
}

// unavailable renders the standard "section could not be fetched" placeholder.
func unavailable(title string, err error) connector.SnapshotSection {
	return connector.SnapshotSection{
		Title:   title,
		Content: "_" + title + " unavailable: " + err.Error() + "_",
	}
}

func (c *Connector) fetchHosts(ctx context.Context, s session) (connector.SnapshotSection, []connector.SnapshotEntity) {
	raw, err := c.get(ctx, s, resourceHosts)
	if err != nil {
		return unavailable("Local DNS Records", err), nil
	}
	content, entities := parseHosts(s.version, raw)
	return connector.SnapshotSection{Title: "Local DNS Records", Content: content}, entities
}

// parseHosts renders the local DNS records for the resolved API version.
func parseHosts(version string, raw []byte) (content string, entities []connector.SnapshotEntity) {
	if version == version5 {
		return buildHostsTableV5(raw)
	}
	return buildHostsTable(raw)
}

// fetchGroups returns the parsed group rows alongside its section, because
// adlists, clients and domain rules reference groups by numeric id.
func (c *Connector) fetchGroups(ctx context.Context, s session) ([]groupRow, connector.SnapshotSection, []connector.SnapshotEntity) {
	raw, err := c.get(ctx, s, resourceGroups)
	if err != nil {
		return nil, unavailable("Groups", err), nil
	}
	var rows []groupRow
	if s.version == version5 {
		rows, err = parseGroupsV5(raw)
	} else {
		rows, err = parseGroupsV6(raw)
	}
	if err != nil {
		return nil, unavailable("Groups", err), nil
	}
	content, entities := buildGroupTable(rows)
	return rows, connector.SnapshotSection{Title: "Groups", Content: content}, entities
}

func (c *Connector) fetchAdlists(ctx context.Context, s session, names map[int]string) (connector.SnapshotSection, []connector.SnapshotEntity) {
	raw, err := c.get(ctx, s, resourceLists)
	if err != nil {
		return unavailable("Blocklists", err), nil
	}
	var rows []adlistRow
	if s.version == version5 {
		rows, err = parseAdlistsV5(raw, names)
	} else {
		rows, err = parseAdlistsV6(raw, names)
	}
	if err != nil {
		return unavailable("Blocklists", err), nil
	}
	content, entities := buildAdlistTable(rows)
	return connector.SnapshotSection{Title: "Blocklists", Content: content}, entities
}

func (c *Connector) fetchClients(ctx context.Context, s session, names map[int]string) (connector.SnapshotSection, []connector.SnapshotEntity) {
	raw, err := c.get(ctx, s, resourceClients)
	if err != nil {
		return unavailable("Clients", err), nil
	}
	var rows []clientRow
	if s.version == version5 {
		rows, err = parseClientsV5(raw, names)
	} else {
		rows, err = parseClientsV6(raw, names)
	}
	if err != nil {
		return unavailable("Clients", err), nil
	}
	content, entities := buildClientTable(rows)
	return connector.SnapshotSection{Title: "Clients", Content: content}, entities
}

func (c *Connector) fetchDomains(ctx context.Context, s session, names map[int]string) (connector.SnapshotSection, []connector.SnapshotEntity) {
	raw, err := c.get(ctx, s, resourceDomains)
	if err != nil {
		return unavailable("Domain Rules", err), nil
	}
	var rows []domainRow
	if s.version == version5 {
		rows, err = parseDomainsV5(raw, names)
	} else {
		rows, err = parseDomainsV6(raw, names)
	}
	if err != nil {
		return unavailable("Domain Rules", err), nil
	}
	content, entities := buildDomainTable(rows)
	return connector.SnapshotSection{Title: "Domain Rules", Content: content}, entities
}

// Restart restarts the Pi-hole DNS resolver (FTL). Pi-hole manages a single
// implicit DNS service, so entityRef is ignored. Only the v6 API exposes a
// restart action.
func (c *Connector) Restart(ctx context.Context, _ map[string]any, _ string) error {
	s, err := c.session(ctx)
	if err != nil {
		return err
	}
	if s.version == version5 {
		return fmt.Errorf("restart is not exposed by the Pi-hole v5 API")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/action/restartdns", nil)
	if err != nil {
		return err
	}
	req.Header.Set("sid", s.sid)
	req.Header.Set("Accept", "application/json")
	_, err = c.send(req)
	return err
}

// Start re-enables Pi-hole DNS blocking. Pi-hole's FTL service has no
// start/stop lifecycle exposed via the API (only restart), so Start/Stop
// here map to the closest on/off concept the API actually exposes: the
// blocking toggle.
func (c *Connector) Start(ctx context.Context, _ map[string]any, _ string) error {
	return c.setBlocking(ctx, true)
}

// Stop disables Pi-hole DNS blocking. See Start for why this maps to the
// blocking toggle rather than a literal service stop.
func (c *Connector) Stop(ctx context.Context, _ map[string]any, _ string) error {
	return c.setBlocking(ctx, false)
}

func (c *Connector) setBlocking(ctx context.Context, blocking bool) error {
	s, err := c.session(ctx)
	if err != nil {
		return err
	}
	if s.version == version5 {
		return c.setBlockingV5(ctx, blocking)
	}
	body, err := json.Marshal(map[string]any{"blocking": blocking, "timer": nil})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/dns/blocking", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("sid", s.sid)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	_, err = c.send(req)
	return err
}

// WritableFields lists the config-push-eligible local DNS record field.
func (c *Connector) WritableFields() []connector.ConfigField {
	return []connector.ConfigField{
		{Key: "ip", Label: "Record IP", Type: "text", EntityScope: true},
	}
}

// ConfigPush repoints the local DNS record for the hostname identified by
// entityRef to a new IP. Pi-hole's hosts config is a flat set of "ip
// hostname" strings with no update verb, so this deletes the existing
// entry for entityRef (if any) and adds the new one.
// ponytail: one field (ip) per call, matching the handler's one-field
// revert contract — not a batch hosts-file replace.
func (c *Connector) ConfigPush(ctx context.Context, _ map[string]any, entityRef, fieldKey string, value any) error {
	if entityRef == "" {
		return fmt.Errorf("pihole config-push requires a target hostname")
	}
	if err := connector.ValidateRefSegment(entityRef); err != nil {
		return fmt.Errorf("invalid entityRef: %w", err)
	}
	if fieldKey != "ip" {
		return fmt.Errorf("unsupported field %q", fieldKey)
	}
	newIP, _ := value.(string)
	s, err := c.session(ctx)
	if err != nil {
		return err
	}
	if raw, err := c.get(ctx, s, resourceHosts); err == nil {
		_, entities := parseHosts(s.version, raw)
		for _, e := range entities {
			if e.Hostname == entityRef && e.IP != "" {
				if err := c.hostsItem(ctx, s, "DELETE", e.IP, e.Hostname); err != nil {
					return fmt.Errorf("remove old record: %w", err)
				}
				break
			}
		}
	}
	return c.hostsItem(ctx, s, "PUT", newIP, entityRef)
}

// hostsItem adds ("PUT") or removes ("DELETE") one local DNS record, mapping
// the verb onto whichever API the session resolved.
func (c *Connector) hostsItem(ctx context.Context, s session, method, ip, hostname string) error {
	if s.version == version5 {
		action := "add"
		if method == "DELETE" {
			action = "delete"
		}
		return c.hostsItemV5(ctx, action, ip, hostname)
	}
	item := ip + " " + hostname
	req, err := http.NewRequestWithContext(ctx, method, c.url+"/api/config/dns/hosts/"+strings.TrimSpace(item), nil)
	if err != nil {
		return err
	}
	req.Header.Set("sid", s.sid)
	req.Header.Set("Accept", "application/json")
	_, err = c.send(req)
	return err
}

// authenticate exchanges the configured password for a session id (sid) via
// POST /api/auth.
func (c *Connector) authenticate(ctx context.Context) (sid string, err error) {
	body, err := json.Marshal(map[string]string{"password": c.password})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/api/auth", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return "", connector.NewTimeoutError(fmt.Errorf("auth request failed: %w", err))
		}
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := connector.ReadBody(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read auth response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", connector.NewAuthError(fmt.Errorf("auth returned %d: %s", resp.StatusCode, string(data)))
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("auth returned %d: %s", resp.StatusCode, string(data))
	}

	var authResp struct {
		Session struct {
			SID     string `json:"sid"`
			Valid   bool   `json:"valid"`
			Message string `json:"message"`
		} `json:"session"`
	}
	if err := json.Unmarshal(data, &authResp); err != nil {
		return "", connector.NewMalformedResponseError(err)
	}
	if authResp.Session.SID == "" {
		return "", connector.NewAuthError(fmt.Errorf("no session id returned: %s", authResp.Session.Message))
	}
	return authResp.Session.SID, nil
}

func (c *Connector) doRequest(ctx context.Context, sid, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("sid", sid)
	req.Header.Set("Accept", "application/json")
	return c.send(req)
}

// send issues req and maps transport and HTTP status failures onto the
// connector error taxonomy, reading the body under connector.ReadBody's cap.
func (c *Connector) send(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := connector.ReadBody(resp.Body)
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

func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// buildHostsTable renders the Pi-hole local DNS "IP hostname" entries as a
// markdown table and extracts one SnapshotEntity per record.
func buildHostsTable(raw []byte) (content string, entities []connector.SnapshotEntity) {
	var resp struct {
		Config struct {
			DNS struct {
				Hosts []string `json:"hosts"`
			} `json:"dns"`
		} `json:"config"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "_Local DNS records unavailable: " + connector.NewMalformedResponseError(err).Error() + "_", nil
	}
	if len(resp.Config.DNS.Hosts) == 0 {
		return "_No local DNS records returned_", nil
	}
	var b strings.Builder
	b.WriteString("| Hostname | IP |\n")
	b.WriteString("|----------|----|\n")
	for _, entry := range resp.Config.DNS.Hosts {
		fields := strings.Fields(entry)
		if len(fields) < 2 {
			continue
		}
		ip, hostname := fields[0], fields[1]
		if _, err := fmt.Fprintf(&b, "| %s | %s |\n", hostname, ip); err != nil {
			return "", nil
		}
		attrs := map[string]any{
			"source":  "local_dns",
			"is_ipv6": net.ParseIP(ip).To4() == nil,
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "dns_record",
			Hostname:   hostname,
			IP:         ip,
			Attributes: attrs,
		})
	}
	if len(entities) == 0 {
		return "_No local DNS records returned_", nil
	}
	return b.String(), entities
}
