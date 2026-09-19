package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.yaml.in/yaml/v3"
)

// Known drift between docs/openapi.yaml and the router. Each entry is
// "METHOD /path" (path relative to the /api server prefix, with {param}
// placeholders normalised). Keep these lists small: fix the spec or the
// router instead of growing them.
var (
	// routerOnly: served by the router but deliberately absent from the spec.
	routerOnly = map[string]string{
		"GET /ws":               "WebSocket upgrade, not a REST operation; documented in docs/WS_CONTRACT.md",
		"PATCH /connectors/{p}": "legacy alias of PUT for clients predating the spec (see router.go)",
	}
	// specOnly: documented in the spec but not implemented by the router.
	specOnly = map[string]string{}
)

func normalizeParams(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			parts[i] = "{p}"
		}
	}
	return strings.Join(parts, "/")
}

func routerOperations(t *testing.T, prefix string, routes chi.Routes) map[string]bool {
	t.Helper()
	ops := map[string]bool{}
	err := chi.Walk(routes, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/") // chi Mount wildcard artefacts
		if !strings.HasPrefix(route, prefix+"/") {
			return nil
		}
		route = strings.TrimPrefix(route, prefix)
		if route != "/" {
			route = strings.TrimSuffix(route, "/")
		}
		ops[method+" "+normalizeParams(route)] = true
		return nil
	})
	if err != nil {
		t.Fatalf("walk router: %v", err)
	}
	return ops
}

func specOperations(t *testing.T) map[string]bool {
	t.Helper()
	raw, err := os.ReadFile("../../../docs/openapi.yaml")
	if err != nil {
		t.Fatalf("read spec: %v", err)
	}
	var doc struct {
		Paths map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parse spec: %v", err)
	}
	methods := map[string]bool{"get": true, "put": true, "post": true, "delete": true, "patch": true}
	ops := map[string]bool{}
	for path, item := range doc.Paths {
		for m := range item {
			if methods[m] {
				ops[strings.ToUpper(m)+" "+normalizeParams(path)] = true
			}
		}
	}
	return ops
}

func TestOpenAPIMatchesRouter(t *testing.T) {
	app := newTestApp(t)
	routes, ok := app.Router.(chi.Routes)
	if !ok {
		t.Fatalf("router is %T, want chi.Routes", app.Router)
	}
	legacy := routerOperations(t, "/api", routes)
	// /api/v1 is mounted as an alias of /api; the two must not diverge.
	v1 := routerOperations(t, "/api/v1", routes)
	for op := range legacy {
		if strings.Contains(op, " /v1/") {
			delete(legacy, op) // /api/v1/... is seen through the /api prefix too
		}
	}
	for op := range legacy {
		if !v1[op] {
			t.Errorf("%s served under /api but not /api/v1", op)
		}
	}
	for op := range v1 {
		if !legacy[op] {
			t.Errorf("%s served under /api/v1 but not /api", op)
		}
	}

	spec := specOperations(t)
	var problems []string
	for op := range legacy {
		if !spec[op] && routerOnly[op] == "" {
			problems = append(problems, "in router, missing from openapi.yaml: "+op)
		}
	}
	for op := range spec {
		if !legacy[op] && specOnly[op] == "" {
			problems = append(problems, "in openapi.yaml, missing from router: "+op)
		}
	}
	for op := range routerOnly {
		if !legacy[op] || spec[op] {
			problems = append(problems, "stale routerOnly allowlist entry: "+op)
		}
	}
	for op := range specOnly {
		if legacy[op] || !spec[op] {
			problems = append(problems, "stale specOnly allowlist entry: "+op)
		}
	}
	sort.Strings(problems)
	for _, p := range problems {
		t.Error(p)
	}
}

func TestAPIV1AliasServesSameHandlers(t *testing.T) {
	app := newTestApp(t)
	for _, path := range []string{"/api/health", "/api/v1/health"} {
		rec := httptest.NewRecorder()
		app.Router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, rec.Code)
		}
	}
	// Auth middleware must gate the alias exactly like the legacy prefix.
	for _, path := range []string{"/api/me", "/api/v1/me", "/api/v1/connectors"} {
		rec := httptest.NewRecorder()
		app.Router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("unauthenticated GET %s = %d, want 401", path, rec.Code)
		}
	}
}
