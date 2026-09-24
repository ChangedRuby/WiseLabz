// Package traefik implements a Traefik reverse-proxy connector, targeting
// the Traefik v2/v3 HTTP API (the same API that backs the dashboard).
//
// The API is commonly exposed unauthenticated on localhost, or fronted by a
// basic-auth / forward-auth middleware, so the connector supports an
// unauthenticated, a basic-auth and a bearer-token mode.
//
// Note: the HTTP plumbing below (guarded dialer, TLS config, status-code to
// connector-error mapping, timeout detection) is deliberately duplicated
// from the sibling connectors rather than shared — #265/#266 will rewrite
// the connector base, and this package stays self-contained until then.
package traefik

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const typeName = "traefik"

// Auth modes accepted by the "auth_mode" config field.
const (
	authNone  = "none"
	authBasic = "basic"
	authToken = "token"
)

// Traefik API paths this connector reads.
const (
	pathOverview    = "/api/overview"
	pathVersion     = "/api/version"
	pathRouters     = "/api/http/routers"
	pathServices    = "/api/http/services"
	pathMiddlewares = "/api/http/middlewares"
	pathEntryPoints = "/api/entrypoints"
)

func init() {
	connector.Register(connector.TypeSchema{
		Type:     typeName,
		Category: "networking",
		Name:     "Traefik",
		Fields: []connector.SchemaField{
			{Key: "url", Label: "Traefik API URL", Type: "text", Required: true, Placeholder: "http://traefik.example.com:8080", Description: "Base URL of the Traefik API/dashboard (the host serving /api/overview)."},
			{Key: "auth_mode", Label: "Authentication", Type: "select", Required: false, Default: authNone, Options: []string{authNone, authBasic, authToken}, Description: "How the API is protected: none (localhost), basic auth, or a bearer token / API key."},
			{Key: "username", Label: "Username", Type: "text", Required: false, Description: "Basic-auth username (auth mode \"basic\")."},
			{Key: "password", Label: "Password", Type: "password", Required: false, Description: "Basic-auth password (auth mode \"basic\")."},
			{Key: "api_token", Label: "API Token", Type: "password", Required: false, Description: "Token sent as \"Authorization: Bearer <token>\" (auth mode \"token\")."},
			{Key: "verify_tls", Label: "Verify TLS", Type: "toggle", Required: false, Default: "true"},
		},
	}, newConnector)
	connector.RegisterAttributeCatalog(typeName, attributeCatalog)
}

// attributeCatalog declares the structured Attributes this connector fills
// on router/service/middleware/entrypoint entities, exposed via
// GET /api/compliance/schema.
var attributeCatalog = map[string][]connector.AttributeSpec{
	"router": {
		{Name: "status", Type: "string", Description: "Router status reported by Traefik (enabled, disabled, warning)"},
		{Name: "rule", Type: "string", Description: "Matching rule, e.g. Host(`app.example.com`)"},
		{Name: "service", Type: "string", Description: "Name of the service the router forwards to"},
		{Name: "provider", Type: "string", Description: "Configuration provider the router came from (docker, file, internal, ...)"},
		{Name: "entryPoints", Type: "string_array", Description: "Entry points the router listens on"},
		{Name: "middlewares", Type: "string_array", Description: "Middlewares applied to the router"},
		{Name: "tls", Type: "boolean", Description: "Whether the router terminates TLS"},
		{Name: "certResolver", Type: "string", Description: "ACME certificate resolver used by the router, if any"},
	},
	"service": {
		{Name: "status", Type: "string", Description: "Service status reported by Traefik (enabled, disabled, warning)"},
		{Name: "provider", Type: "string", Description: "Configuration provider the service came from"},
		{Name: "serverCount", Type: "number", Description: "Number of backend servers in the load balancer"},
		{Name: "servers", Type: "string_array", Description: "Backend server URLs or addresses"},
		{Name: "passHostHeader", Type: "boolean", Description: "Whether the original Host header is forwarded to the backend"},
	},
	"middleware": {
		{Name: "status", Type: "string", Description: "Middleware status reported by Traefik (enabled, disabled, warning)"},
		{Name: "provider", Type: "string", Description: "Configuration provider the middleware came from"},
		{Name: "types", Type: "string_array", Description: "Middleware kinds configured on this middleware (basicAuth, headers, redirectScheme, ...)"},
	},
	"entrypoint": {
		{Name: "address", Type: "string", Description: "Listen address of the entry point, e.g. \":443\""},
		{Name: "asDefault", Type: "boolean", Description: "Whether the entry point is used by default when a router declares none"},
		{Name: "tls", Type: "boolean", Description: "Whether the entry point has a default TLS configuration"},
	},
}

// Connector fetches routing configuration from a Traefik API.
type Connector struct {
	url      string
	authMode string
	username string
	password string
	apiToken string
	client   *http.Client
}

func newConnector(config map[string]any) (connector.Connector, error) {
	rawURL, _ := config["url"].(string)
	authMode, _ := config["auth_mode"].(string)
	if authMode == "" {
		authMode = authNone
	}
	username, _ := config["username"].(string)
	password, _ := config["password"].(string)
	apiToken, _ := config["api_token"].(string)
	verifyTLS := true
	if v, ok := config["verify_tls"]; ok {
		if b, ok := v.(bool); ok {
			verifyTLS = b
		}
	}
	client := connector.NewHTTPClient(connector.HTTPClientOptions{SkipTLSVerify: !verifyTLS})
	return &Connector{
		url:      strings.TrimSuffix(rawURL, "/"),
		authMode: authMode,
		username: username,
		password: password,
		apiToken: apiToken,
		client:   client,
	}, nil
}

// Name returns the connector display name.
func (c *Connector) Name() string { return "Traefik" }

// Type returns the connector type identifier.
func (c *Connector) Type() string { return typeName }

// Category returns the connector category.
func (c *Connector) Category() string { return "networking" }

// Validate checks the credentials make sense for the selected auth mode and
// tests the connection against the Traefik API.
func (c *Connector) Validate(ctx context.Context, _ map[string]any) error {
	if c.url == "" {
		return fmt.Errorf("traefik url is required")
	}
	switch c.authMode {
	case authNone:
	case authBasic:
		if c.username == "" || c.password == "" {
			return fmt.Errorf("basic auth requires both username and password")
		}
	case authToken:
		if c.apiToken == "" {
			return fmt.Errorf("token auth requires an api_token")
		}
	default:
		return fmt.Errorf("unknown auth mode %q", c.authMode)
	}
	_, err := c.doRequest(ctx, pathOverview)
	return err
}

// Fetch retrieves the Traefik overview plus HTTP routers, services,
// middlewares and entry points.
func (c *Connector) Fetch(ctx context.Context, config map[string]any) (*connector.ServiceSnapshot, error) {
	start := time.Now()
	fields := connector.RequestedFields(config)
	metadata := map[string]string{"traefik_url": c.url}
	var sections []connector.SnapshotSection
	var entities []connector.SnapshotEntity

	if v, err := c.fetchVersion(ctx); err == nil && v != "" {
		metadata["traefik_version"] = v
	}

	if connector.WantsField(fields, "overview") {
		raw, err := c.doRequest(ctx, pathOverview)
		if err != nil {
			sections = append(sections, unavailable("Overview", err))
		} else {
			content, counts := buildOverview(raw)
			sections = append(sections, connector.SnapshotSection{Title: "Overview", Content: content})
			for k, v := range counts {
				metadata[k] = v
			}
		}
	}

	var routerServices []string
	if connector.WantsField(fields, "routers") {
		raw, err := c.doRequest(ctx, pathRouters)
		if err != nil {
			sections = append(sections, unavailable("HTTP Routers", err))
		} else {
			content, ents, svcs := buildRouterTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "HTTP Routers", Content: content})
			entities = append(entities, ents...)
			routerServices = svcs
			metadata["router_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	if connector.WantsField(fields, "services") {
		raw, err := c.doRequest(ctx, pathServices)
		if err != nil {
			sections = append(sections, unavailable("HTTP Services", err))
		} else {
			content, ents := buildServiceTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "HTTP Services", Content: content})
			entities = append(entities, ents...)
			metadata["service_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	if connector.WantsField(fields, "middlewares") {
		raw, err := c.doRequest(ctx, pathMiddlewares)
		if err != nil {
			sections = append(sections, unavailable("HTTP Middlewares", err))
		} else {
			content, ents := buildMiddlewareTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "HTTP Middlewares", Content: content})
			entities = append(entities, ents...)
			metadata["middleware_count"] = fmt.Sprintf("%d", len(ents))
		}
	}

	if connector.WantsField(fields, "entrypoints") {
		raw, err := c.doRequest(ctx, pathEntryPoints)
		if err != nil {
			sections = append(sections, unavailable("Entry Points", err))
		} else {
			content, ents := buildEntryPointTable(raw)
			sections = append(sections, connector.SnapshotSection{Title: "Entry Points", Content: content})
			entities = append(entities, ents...)
		}
	}

	return &connector.ServiceSnapshot{
		ServiceName:  "Traefik",
		Type:         typeName,
		Sections:     sections,
		Entities:     entities,
		Dependencies: serviceDependencies(routerServices),
		Metadata:     metadata,
		FetchedAt:    start,
	}, nil
}

// fetchVersion reads /api/version, which older v2 builds may not expose.
func (c *Connector) fetchVersion(ctx context.Context) (string, error) {
	raw, err := c.doRequest(ctx, pathVersion)
	if err != nil {
		return "", err
	}
	var v struct {
		Version string `json:"Version"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", connector.NewMalformedResponseError(err)
	}
	return v.Version, nil
}

// serviceDependencies turns the set of services referenced by routers into
// upstream_service dependencies, deduplicated and sorted for a stable
// snapshot.
func serviceDependencies(services []string) []connector.ServiceDependency {
	seen := make(map[string]struct{}, len(services))
	var names []string
	for _, s := range services {
		if s == "" {
			continue
		}
		if _, dup := seen[s]; dup {
			continue
		}
		seen[s] = struct{}{}
		names = append(names, s)
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
	switch c.authMode {
	case authBasic:
		req.SetBasicAuth(c.username, c.password)
	case authToken:
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, connector.MapTransportError(err)
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

	if statusErr := connector.CheckStatus(resp.StatusCode, data); statusErr != nil {
		return nil, statusErr
	}

	return data, nil
}
