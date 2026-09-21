// Package adguardhome implements an AdGuard Home connector, targeting the
// /control REST API that backs the AdGuard Home web UI.
//
// The API is protected by HTTP basic auth when a user is configured, and is
// reachable unauthenticated on installs that never set one up, so the
// connector supports an unauthenticated and a basic-auth mode.
//
// Note: the HTTP plumbing below (guarded dialer, TLS config, status-code to
// connector-error mapping, timeout detection) is deliberately duplicated
// from the sibling connectors rather than shared — #265/#266 will rewrite
// the connector base, and this package stays self-contained until then.
package adguardhome

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "adguardhome"

// Auth modes accepted by the "auth_mode" config field.
const (
	authNone  = "none"
	authBasic = "basic"
)

// AdGuard Home API paths this connector reads.
const (
	pathStatus          = "/control/status"
	pathDNSInfo         = "/control/dns_info"
	pathFilteringStatus = "/control/filtering/status"
	pathRewriteList     = "/control/rewrite/list"
	pathClients         = "/control/clients"
	pathDHCPStatus      = "/control/dhcp/status"
)

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "dns",
		Name:     "AdGuard Home",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "AdGuard Home URL", Type: "text", Required: true, Placeholder: "http://adguard.example.com:3000", Description: "Base URL of the AdGuard Home web interface (the host serving /control/status)."},
			{Key: "auth_mode", Label: "Authentication", Type: "select", Required: false, Default: authBasic, Options: []string{authNone, authBasic}, Description: "How the API is protected: basic auth with the web-UI user, or none when no user is configured."},
			{Key: "username", Label: "Username", Type: "text", Required: false, Description: "Web-interface username (auth mode \"basic\")."},
			{Key: "password", Label: "Password", Type: "password", Required: false, Description: "Web-interface password (auth mode \"basic\")."},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, newConnector)
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on filter-list/rewrite/client/lease entities, exposed via
// GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"filter_list": {
		{Name: "enabled", Type: "boolean", Description: "Whether the filter list is applied"},
		{Name: "kind", Type: "string", Description: "Whether the list blocks or allows (blocklist, allowlist)"},
		{Name: "url", Type: "string", Description: "Source URL the filter list is downloaded from"},
	},
	"dns_rewrite": {
		{Name: "answer", Type: "string", Description: "Address or hostname the rewrite resolves the domain to"},
		{Name: "wildcard", Type: "boolean", Description: "Whether the rewritten domain is a wildcard pattern"},
	},
	"client": {
		{Name: "ids", Type: "string_array", Description: "Identifiers matching the client (IPs, CIDRs, MACs, ClientIDs)"},
		{Name: "useGlobalSettings", Type: "boolean", Description: "Whether the client inherits the global filtering settings"},
		{Name: "filteringEnabled", Type: "boolean", Description: "Whether filtering is enabled for this client"},
		{Name: "safebrowsingEnabled", Type: "boolean", Description: "Whether browsing security is enabled for this client"},
		{Name: "parentalEnabled", Type: "boolean", Description: "Whether parental control is enabled for this client"},
		{Name: "blockedServices", Type: "string_array", Description: "Services blocked for this client"},
		{Name: "upstreams", Type: "string_array", Description: "Client-specific upstream DNS servers"},
		{Name: "tags", Type: "string_array", Description: "Tags assigned to the client"},
	},
	"dhcp_lease": {
		{Name: "mac", Type: "string", Description: "MAC address the lease is bound to"},
		{Name: "static", Type: "boolean", Description: "Whether the lease is statically reserved rather than dynamic"},
	},
}

// Connector fetches DNS, filtering and DHCP configuration from an AdGuard
// Home instance.
type Connector struct {
	url      string
	authMode string
	username string
	password string
	client   *http.Client
}

func newConnector(config map[string]any) (connector.Connector, error) {
	rawURL, _ := config["url"].(string)
	authMode, _ := config["auth_mode"].(string)
	if authMode == "" {
		authMode = authBasic
	}
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
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
		url:      strings.TrimSuffix(rawURL, "/"),
		authMode: authMode,
		username: username,
		password: password,
		client:   client,
	}, nil
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "AdGuard Home" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "dns" }

// Validate checks the credentials make sense for the selected auth mode and
// tests the connection against the AdGuard Home API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	if c.url == "" {
		return fmt.Errorf("adguard home url is required")
	}
	switch c.authMode {
	case authNone:
	case authBasic:
		if c.username == "" || c.password == "" {
			return fmt.Errorf("basic auth requires both username and password")
		}
	default:
		return fmt.Errorf("unknown auth mode %q", c.authMode)
	}
	_, err := c.doRequest(ctx, pathStatus)
	return err
}

// Fetch retrieves the AdGuard Home status plus DNS settings, filter lists,
// custom rules, DNS rewrites, clients and DHCP configuration.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	metadata := map[string]string{"adguard_url": c.url}
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	var upstreams []string

	// The status response also tells us whether DHCP is even available, so
	// it is read whenever DHCP is wanted, even if "status" itself isn't.
	var status statusInfo
	if connector.WantsField(fields, "status") || connector.WantsField(fields, "dhcp") {
		raw, err := c.doRequest(ctx, pathStatus)
		switch {
		case err != nil:
			if connector.WantsField(fields, "status") {
				sections = append(sections, unavailable("Status", err))
			}
		default:
			var content string
			content, status = buildStatus(raw)
			if connector.WantsField(fields, "status") {
				sections = append(sections, connector.SnapshotSection{Title: "Status", Content: content})
			}
			for k, v := range status.metadata {
				metadata[k] = v
			}
		}
	}

	if connector.WantsField(fields, "dns_config") {
		raw, err := c.doRequest(ctx, pathDNSInfo)
		if err != nil {
			sections = append(sections, unavailable("DNS Configuration", err))
		} else {
			content, ups, meta := buildDNSInfo(raw)
			sections = append(sections, connector.SnapshotSection{Title: "DNS Configuration", Content: content})
			upstreams = ups
			for k, v := range meta {
				metadata[k] = v
			}
		}
	}

	if connector.WantsField(fields, "filtering") {
		raw, err := c.doRequest(ctx, pathFilteringStatus)
		if err != nil {
			sections = append(sections,
				unavailable("Filter Lists", err),
				unavailable("Custom Filtering Rules", err))
		} else {
			listContent, rulesContent, ents, meta := buildFiltering(raw)
			sections = append(sections,
				connector.SnapshotSection{Title: "Filter Lists", Content: listContent},
				connector.SnapshotSection{Title: "Custom Filtering Rules", Content: rulesContent})
			entities = append(entities, ents...)
			for k, v := range meta {
				metadata[k] = v
			}
		}
	}

	if connector.WantsField(fields, "rewrites") {
		raw, err := c.doRequest(ctx, pathRewriteList)
		if err != nil {
			sections = append(sections, unavailable("DNS Rewrites", err))
		} else {
			content, ents := buildRewriteTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "DNS Rewrites", Content: content})
			entities = append(entities, ents...)
			metadata["rewrite_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	if connector.WantsField(fields, "clients") {
		raw, err := c.doRequest(ctx, pathClients)
		if err != nil {
			sections = append(sections, unavailable("Clients", err))
		} else {
			content, ents := buildClientTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "Clients", Content: content})
			entities = append(entities, ents...)
			metadata["client_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	// DHCP is optional: AdGuard Home builds without DHCP support report
	// dhcp_available=false and answer /control/dhcp/status with an error,
	// so the section is skipped entirely rather than reported as degraded.
	if connector.WantsField(fields, "dhcp") && status.dhcpAvailable {
		raw, err := c.doRequest(ctx, pathDHCPStatus)
		if err != nil {
			sections = append(sections, unavailable("DHCP", err))
		} else {
			content, ents, meta := buildDHCP(raw)
			sections = append(sections, connector.SnapshotSection{Title: "DHCP", Content: content})
			entities = append(entities, ents...)
			for k, v := range meta {
				metadata[k] = v
			}
		}
	}

	return &connector.ServiceSnapshot{
		ServiceName:  "AdGuard Home",
		Type:         typeName,
		Sections:     sections,
		Entities:     entities,
		Dependencies: upstreamDependencies(upstreams),
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

// upstreamDependencies turns the configured upstream resolvers into
// upstream_service dependencies, deduplicated and sorted for a stable
// snapshot.
func upstreamDependencies(upstreams []string) []connector.ServiceDependency {
	seen := make(map[string]struct{}, len(upstreams))
	var names []string
	for _, u := range upstreams {
		if u == "" {
			continue
		}
		if _, dup := seen[u]; dup {
			continue
		}
		seen[u] = struct{}{}
		names = append(names, u)
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	deps := make([]connector.ServiceDependency, 0, len(names))
	for _, n := range names {
		deps = append(deps, connector.ServiceDependency{Kind: "upstream_service", Name: n})
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
	if c.authMode == authBasic {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, connector.NewTimeoutError(fmt.Errorf("request failed: %w", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
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

// isTimeout reports whether err represents a request deadline being
// exceeded, covering both a canceled context and a net.Error timeout.
func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
