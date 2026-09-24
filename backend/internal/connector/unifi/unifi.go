// Package unifi implements a Ubiquiti UniFi connector, targeting the UniFi
// Network Application controller REST API.
//
// Two controller flavours are supported: a classic self-hosted Network
// Application (login at /api/login, endpoints directly under /api) and a
// UniFi OS console such as a UDM or Cloud Key gen2 (login at
// /api/auth/login, endpoints behind the /proxy/network prefix). UniFi OS
// consoles also accept a site-wide API key instead of a user session.
//
// Only configuration-shaped data is read: sites, adopted devices,
// networks/VLANs, WLANs, firewall rules and a client count summary. Volatile
// counters the controller reports alongside them (uptime, tx/rx bytes,
// last_seen, signal) are deliberately dropped so an unchanged controller
// produces an identical snapshot.
//
// Note: the HTTP plumbing below (guarded dialer, TLS config, status-code to
// connector-error mapping, timeout detection) is deliberately duplicated
// from the sibling connectors rather than shared — #265/#266 will rewrite
// the connector base, and this package stays self-contained until then.
package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "unifi"

// Auth modes accepted by the "auth_mode" config field.
const (
	authPassword = "password"
	authAPIKey   = "api_key"
)

// Controller flavours accepted by the "controller_type" config field.
// controllerAuto probes UniFi OS first and falls back to the classic
// Network Application.
const (
	controllerAuto    = "auto"
	controllerUniFiOS = "unifios"
	controllerClassic = "classic"
)

// Path prefix each controller flavour serves the Network API under.
const (
	prefixClassic  = ""
	prefixUniFiOS  = "/proxy/network"
	pathLoginClass = "/api/login"
	pathLoginUOS   = "/api/auth/login"
	pathSites      = "/api/self/sites"
)

// session carries the resolved controller flavour for one call: the path
// prefix every Network API request goes through. The credentials themselves
// live elsewhere — in password mode the login handshake leaves a session
// cookie in the client's jar, and in API-key mode every request carries the
// key header.
type session struct {
	prefix string
}

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "networking",
		Name:     "UniFi",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Controller URL", Type: "text", Required: true, Placeholder: "https://unifi.example.com:8443", Description: "Base URL of the UniFi Network Application or UniFi OS console."},
			{Key: "auth_mode", Label: "Authentication", Type: "select", Required: false, Default: authPassword, Options: []string{authPassword, authAPIKey}, Description: "Local controller account (username/password) or a UniFi OS API key."},
			{Key: "username", Label: "Username", Type: "text", Required: false, Description: "Local controller account with read access (auth mode \"password\")."},
			{Key: "password", Label: "Password", Type: "password", Required: false, Description: "Password for the local controller account (auth mode \"password\")."},
			{Key: "api_key", Label: "API Key", Type: "password", Required: false, Description: "Sent as \"X-API-KEY\"; supported by UniFi OS consoles (auth mode \"api_key\")."},
			{Key: "site", Label: "Site", Type: "text", Required: false, Default: "default", Description: "Internal site name to read, as it appears in the controller URL (not the display name)."},
			{Key: "controller_type", Label: "Controller Type", Type: "select", Required: false, Default: controllerAuto, Options: []string{controllerAuto, controllerUniFiOS, controllerClassic}, Description: "UniFi OS console, classic Network Application, or auto-detect."},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true", Description: "UniFi controllers ship a self-signed certificate; disable to accept it."},
		},
	}, newConnector)
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on its entities, exposed via GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"site": {
		{Name: "role", Type: "string", Description: "Role the authenticated account holds on the site (admin, readonly)"},
		{Name: "displayName", Type: "string", Description: "Human-readable site description shown in the controller UI"},
	},
	"device": {
		{Name: "model", Type: "string", Description: "Hardware model code reported by the controller (U6LR, US8P60, ...)"},
		{Name: "deviceType", Type: "string", Description: "Device family: uap (access point), usw (switch), ugw/udm (gateway)"},
		{Name: "firmwareVersion", Type: "string", Description: "Firmware version currently running on the device"},
		{Name: "adopted", Type: "boolean", Description: "Whether the device has been adopted by this controller"},
		{Name: "disabled", Type: "boolean", Description: "Whether the device is administratively disabled"},
		{Name: "state", Type: "string", Description: "Provisioning state reported by the controller (connected, offline, adopting, ...)"},
	},
	"network": {
		{Name: "enabled", Type: "boolean", Description: "Whether the network is active"},
		{Name: "purpose", Type: "string", Description: "Network purpose (corporate, guest, vlan-only, wan, remote-user-vpn)"},
		{Name: "vlanEnabled", Type: "boolean", Description: "Whether the network is tagged with a VLAN id"},
		{Name: "vlan", Type: "number", Description: "VLAN id the network is tagged with"},
		{Name: "subnet", Type: "string", Description: "IPv4 subnet in CIDR notation"},
		{Name: "dhcpEnabled", Type: "boolean", Description: "Whether the controller runs a DHCP server on the network"},
		{Name: "networkGroup", Type: "string", Description: "Interface group the network is bound to (LAN, WAN, WAN2)"},
	},
	"wlan": {
		{Name: "enabled", Type: "boolean", Description: "Whether the WLAN is broadcast"},
		{Name: "security", Type: "string", Description: "Security mode (open, wpapsk, wpaeap)"},
		{Name: "wpaMode", Type: "string", Description: "WPA mode negotiated with clients (wpa2, wpa3, wpa3-transition)"},
		{Name: "guest", Type: "boolean", Description: "Whether the WLAN is a guest network subject to guest policies"},
		{Name: "hideSsid", Type: "boolean", Description: "Whether the SSID is hidden from beacons"},
		{Name: "macFilterEnabled", Type: "boolean", Description: "Whether a MAC address filter is applied to the WLAN"},
		{Name: "macFilterPolicy", Type: "string", Description: "MAC filter policy (allow, deny)"},
		{Name: "pmfMode", Type: "string", Description: "Protected Management Frames mode (disabled, optional, required)"},
	},
	"firewall_rule": {
		{Name: "enabled", Type: "boolean", Description: "Whether the firewall rule is active"},
		{Name: "action", Type: "string", Description: "Action taken on matching traffic (accept, drop, reject)"},
		{Name: "ruleset", Type: "string", Description: "Ruleset the rule belongs to (WAN_IN, LAN_IN, GUEST_LOCAL, ...)"},
		{Name: "ruleIndex", Type: "number", Description: "Evaluation order of the rule inside its ruleset"},
		{Name: "protocol", Type: "string", Description: "Protocol the rule matches (all, tcp, udp, tcp_udp, icmp)"},
		{Name: "logging", Type: "boolean", Description: "Whether matching traffic is logged"},
		{Name: "sourceAddress", Type: "string", Description: "Source address or network the rule matches"},
		{Name: "destinationAddress", Type: "string", Description: "Destination address or network the rule matches"},
	},
}

// Connector fetches network configuration from a UniFi controller.
type Connector struct {
	url            string
	authMode       string
	username       string
	password       string
	apiKey         string
	site           string
	controllerType string
	client         *http.Client
}

func newConnector(config map[string]any) (connector.Connector, error) {
	rawURL, _ := config["url"].(string)
	authMode, _ := config["auth_mode"].(string)
	if authMode != authAPIKey {
		authMode = authPassword
	}
	controllerType, _ := config["controller_type"].(string)
	switch controllerType {
	case controllerUniFiOS, controllerClassic:
	default:
		controllerType = controllerAuto
	}
	site, _ := config["site"].(string)
	if site == "" {
		site = "default"
	}
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
	apiKey, _ := config["api_key"].(string)
	verifyTLS := true
	if v, ok := config["verify_tls"]; ok {
		if b, ok := v.(bool); ok {
			verifyTLS = b
		}
	}
	// A cookie jar is what carries the UniFi session: both login endpoints
	// answer with a session cookie (unifises on classic, TOKEN on UniFi OS)
	// that every subsequent request must replay.
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}
	client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS, Jar: jar})
	return &Connector{
		url:            strings.TrimSuffix(rawURL, "/"),
		authMode:       authMode,
		username:       username,
		password:       password,
		apiKey:         apiKey,
		site:           site,
		controllerType: controllerType,
		client:         client,
	}, nil
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "UniFi" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "networking" }

// Validate checks the credentials make sense for the selected auth mode and
// tests them against the controller.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	if c.url == "" {
		return fmt.Errorf("unifi url is required")
	}
	switch c.authMode {
	case authPassword:
		if c.username == "" || c.password == "" {
			return fmt.Errorf("password auth requires both username and password")
		}
	case authAPIKey:
		if c.apiKey == "" {
			return fmt.Errorf("api key auth requires an api_key")
		}
	}
	_, err := c.session(ctx)
	return err
}

// session resolves which controller flavour to talk to and, in password
// mode, performs the login handshake. With controller_type pinned only that
// flavour is attempted; on "auto" UniFi OS is probed first and the classic
// controller is the fallback, reporting the UniFi OS error when neither
// answers since that is the likelier deployment.
func (c *Connector) session(ctx context.Context) (session, error) {
	switch c.controllerType {
	case controllerClassic:
		return session{prefix: prefixClassic}, c.open(ctx, prefixClassic)
	case controllerUniFiOS:
		return session{prefix: prefixUniFiOS}, c.open(ctx, prefixUniFiOS)
	default:
		uosErr := c.open(ctx, prefixUniFiOS)
		if uosErr == nil {
			return session{prefix: prefixUniFiOS}, nil
		}
		classicErr := c.open(ctx, prefixClassic)
		if classicErr == nil {
			return session{prefix: prefixClassic}, nil
		}
		// A 404 from the UniFi OS probe only means "not that flavour",
		// while a rejected credential is the actionable answer, so an auth
		// error from the classic probe wins over it.
		var authErr *connector.AuthError
		if !errors.As(uosErr, &authErr) && errors.As(classicErr, &authErr) {
			return session{}, classicErr
		}
		return session{}, uosErr
	}
}

// open establishes usable credentials for one controller flavour: a login
// handshake in password mode, or a probe of the sites endpoint in API-key
// mode, where there is nothing to hand shake.
func (c *Connector) open(ctx context.Context, prefix string) error {
	if c.authMode == authAPIKey {
		_, err := c.get(ctx, session{prefix: prefix}, pathSites)
		return err
	}
	return c.login(ctx, prefix)
}

// login posts the local account credentials to the flavour's login endpoint;
// the session cookie it returns lands in the client's cookie jar.
func (c *Connector) login(ctx context.Context, prefix string) error {
	path := pathLoginClass
	if prefix == prefixUniFiOS {
		path = pathLoginUOS
	}
	body, err := json.Marshal(map[string]any{
		"username": c.username,
		"password": c.password,
		"remember": true,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	data, status, err := c.do(req)
	if err != nil {
		return err
	}
	// A controller rejects bad credentials with 400 and an error envelope,
	// not 401, so the login endpoint maps that to an auth error itself.
	if status == http.StatusBadRequest {
		return connector.NewAuthError(fmt.Errorf("login rejected: %s", apiMessage(data)))
	}
	return statusError(status, data)
}

// get fetches a Network API path through the session's controller prefix.
func (c *Connector) get(ctx context.Context, s session, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+s.prefix+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.authMode == authAPIKey {
		req.Header.Set("X-API-KEY", c.apiKey)
	}
	return c.send(req)
}

// sitePath builds a site-scoped Network API path for the configured site.
func (c *Connector) sitePath(suffix string) string {
	return "/api/s/" + url.PathEscape(c.site) + suffix
}

// send performs a prepared request and maps transport and HTTP failures onto
// the shared connector error types.
func (c *Connector) send(req *http.Request) ([]byte, error) {
	data, status, err := c.do(req)
	if err != nil {
		return nil, err
	}
	return data, statusError(status, data)
}

// do performs a prepared request and returns the (capped) body alongside the
// status code, mapping only transport-level failures to connector errors.
func (c *Connector) do(req *http.Request) (data []byte, status int, err error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, connector.MapTransportError(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	data, err = connector.ReadBody(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read response: %w", err)
	}
	return data, resp.StatusCode, nil
}

// statusError maps a controller HTTP status onto the shared connector error
// types, returning nil for anything below 400.
func statusError(status int, data []byte) error {
	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return connector.NewAuthError(fmt.Errorf("controller returned %d: %s", status, apiMessage(data)))
	case status == http.StatusBadGateway || status == http.StatusServiceUnavailable || status == http.StatusGatewayTimeout:
		return connector.NewServiceUnavailableError(fmt.Errorf("controller returned %d: %s", status, apiMessage(data)))
	case status >= 400:
		return fmt.Errorf("controller returned %d: %s", status, apiMessage(data))
	}
	return nil
}

// apiMessage pulls the controller's own error string out of the UniFi
// envelope ({"meta":{"rc":"error","msg":"api.err.Invalid"}}), falling back to
// the raw body when it is not shaped like one.
func apiMessage(data []byte) string {
	var envelope struct {
		Meta struct {
			Msg string `json:"msg"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(data, &envelope); err == nil && envelope.Meta.Msg != "" {
		return envelope.Meta.Msg
	}
	return string(data)
}

// sectionFetch describes one snapshot section: the selective-fetch field
// name that gates it, its title, the API path it reads and the builder that
// renders it.
type sectionFetch struct {
	field string
	title string
	path  string
	build func([]byte) (string, []connector.SnapshotEntity, map[string]string)
}

// Fetch retrieves sites, devices, networks, WLANs, firewall rules and a
// client summary from the configured site. Each section degrades
// independently: an endpoint the controller does not expose (or the account
// may not read) renders as unavailable instead of failing the snapshot.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	metadata := map[string]string{"unifi_url": c.url, "unifi_site": c.site}

	s, err := c.session(ctx)
	if err != nil {
		return &connector.ServiceSnapshot{
			ServiceName: "UniFi",
			Type:        typeName,
			Sections:    []connector.SnapshotSection{unavailable("Devices", err)},
			Metadata:    metadata,
			FetchedAt:   start,
		}, nil
	}
	metadata["unifi_controller"] = controllerName(s.prefix)
	if v := c.fetchVersion(ctx, s); v != "" {
		metadata["unifi_version"] = v
	}

	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity
	for _, sec := range c.sectionFetches() {
		if !connector.WantsField(fields, sec.field) {
			continue
		}
		raw, err := c.get(ctx, s, sec.path)
		if err != nil {
			sections = append(sections, unavailable(sec.title, err))
			continue
		}
		content, ents, meta := sec.build(raw)
		sections = append(sections, connector.SnapshotSection{Title: sec.title, Content: content})
		entities = append(entities, ents...)
		for k, v := range meta {
			metadata[k] = v
		}
	}
	for kind, n := range countByKind(entities) {
		metadata[kind+"_count"] = fmt.Sprintf("%d", n)
	}

	return &connector.ServiceSnapshot{
		ServiceName:  "UniFi",
		Type:         typeName,
		Sections:     sections,
		Entities:     entities,
		Dependencies: networkDependencies(entities),
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

// sectionFetches lists the snapshot sections in render order.
func (c *Connector) sectionFetches() []sectionFetch {
	return []sectionFetch{
		{field: "sites", title: "Sites", path: pathSites, build: buildSiteTable},
		{field: "devices", title: "Devices", path: c.sitePath("/stat/device"), build: buildDeviceTable},
		{field: "networks", title: "Networks", path: c.sitePath("/rest/networkconf"), build: buildNetworkTable},
		{field: "wlans", title: "WLANs", path: c.sitePath("/rest/wlanconf"), build: buildWLANTable},
		{field: "firewall_rules", title: "Firewall Rules", path: c.sitePath("/rest/firewallrule"), build: buildFirewallTable},
		{field: "clients", title: "Clients", path: c.sitePath("/stat/sta"), build: buildClientSummary},
	}
}

// fetchVersion reads the controller version out of the site sysinfo,
// ignoring the volatile counters reported next to it. A controller that does
// not serve sysinfo simply contributes no version metadata.
func (c *Connector) fetchVersion(ctx context.Context, s session) string {
	raw, err := c.get(ctx, s, c.sitePath("/stat/sysinfo"))
	if err != nil {
		return ""
	}
	var payload struct {
		Data []struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil || len(payload.Data) == 0 {
		return ""
	}
	return payload.Data[0].Version
}

// controllerName labels the resolved controller flavour for snapshot metadata.
func controllerName(prefix string) string {
	if prefix == prefixUniFiOS {
		return controllerUniFiOS
	}
	return controllerClassic
}

func countByKind(entities []connector.SnapshotEntity) map[string]int {
	counts := map[string]int{}
	for _, e := range entities {
		counts[e.Kind]++
	}
	return counts
}

// networkDependencies turns the configured networks into network
// dependencies, sorted for a stable snapshot.
func networkDependencies(entities []connector.SnapshotEntity) []connector.ServiceDependency {
	var names []string
	for _, e := range entities {
		if e.Kind == "network" && e.Name != "" {
			names = append(names, e.Name)
		}
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	deps := make([]connector.ServiceDependency, 0, len(names))
	for _, n := range names {
		deps = append(deps, connector.ServiceDependency{Kind: "network", Name: n})
	}
	return deps
}

// unavailable renders the standard "section could not be fetched" placeholder.
func unavailable(title string, err error) connector.SnapshotSection {
	return connector.SnapshotSection{Title: title, Content: "_" + title + " unavailable: " + err.Error() + "_"}
}
