package traefik

import (
	"reflect"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildersOnMalformedAndEmptyInput(t *testing.T) {
	tests := []struct {
		name string
		call func([]byte) (string, int)
		want string
	}{
		{
			name: "routers",
			call: func(b []byte) (string, int) {
				c, e, _ := buildRouterTable(b)
				return c, len(e)
			},
			want: "_No routers returned_",
		},
		{
			name: "services",
			call: func(b []byte) (string, int) {
				c, e := buildServiceTable(b)
				return c, len(e)
			},
			want: "_No services returned_",
		},
		{
			name: "middlewares",
			call: func(b []byte) (string, int) {
				c, e := buildMiddlewareTable(b)
				return c, len(e)
			},
			want: "_No middlewares returned_",
		},
		{
			name: "entrypoints",
			call: func(b []byte) (string, int) {
				c, e := buildEntryPointTable(b)
				return c, len(e)
			},
			want: "_No entry points returned_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+" empty array", func(t *testing.T) {
			content, n := tt.call([]byte(`[]`))
			if content != tt.want || n != 0 {
				t.Fatalf("content = %q (%d entities), want %q (0)", content, n, tt.want)
			}
		})
		t.Run(tt.name+" invalid JSON", func(t *testing.T) {
			content, n := tt.call([]byte(`not json`))
			if !strings.Contains(content, "malformed response") || n != 0 {
				t.Fatalf("content = %q (%d entities), want a malformed-response placeholder", content, n)
			}
		})
		t.Run(tt.name+" truncated JSON", func(t *testing.T) {
			content, n := tt.call([]byte(`[{"name":`))
			if !strings.Contains(content, "malformed response") || n != 0 {
				t.Fatalf("content = %q (%d entities), want a malformed-response placeholder", content, n)
			}
		})
	}

	t.Run("overview invalid JSON", func(t *testing.T) {
		content, metadata := buildOverview([]byte(`not json`))
		if !strings.Contains(content, "malformed response") || metadata != nil {
			t.Fatalf("content = %q, metadata = %+v", content, metadata)
		}
	})
}

func TestBuildOverview(t *testing.T) {
	content, metadata := buildOverview([]byte(overviewJSON))

	for _, want := range []string{"| HTTP | 2 | 2 | 1 | 1 |", "| TCP | 1 | 1 | 0 | 0 |", "Providers: Docker, File", "Access log: true", "Metrics: Prometheus", "Tracing: Jaeger"} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing %q:\n%s", want, content)
		}
	}
	want := map[string]string{
		"http_routers_total":     "2",
		"http_services_total":    "2",
		"http_middlewares_total": "1",
		"tcp_routers_total":      "1",
		"udp_routers_total":      "0",
	}
	if !reflect.DeepEqual(metadata, want) {
		t.Errorf("metadata = %+v, want %+v", metadata, want)
	}
}

func TestBuildRouterTable(t *testing.T) {
	content, entities, services := buildRouterTable([]byte(routersJSON))

	if !strings.Contains(content, "yes (letsencrypt)") {
		t.Errorf("content missing the cert resolver:\n%s", content)
	}
	if want := []string{"app@docker", "api@internal"}; !reflect.DeepEqual(services, want) {
		t.Errorf("services = %+v, want %+v", services, want)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}

	app := entities[0]
	if app.Kind != "router" || app.Name != "app@docker" || app.ExternalID != "app@docker" {
		t.Errorf("entities[0] identity = %+v", app)
	}
	if app.Hostname != "app.example.com" {
		t.Errorf("entities[0].Hostname = %q, want %q", app.Hostname, "app.example.com")
	}
	wantAttrs := map[string]any{
		"status":       "enabled",
		"rule":         "Host(`app.example.com`)",
		"service":      "app@docker",
		"provider":     "docker",
		"entryPoints":  []string{"websecure"},
		"middlewares":  []string{"auth@file"},
		"tls":          true,
		"certResolver": "letsencrypt",
	}
	if !reflect.DeepEqual(app.Attributes, wantAttrs) {
		t.Errorf("entities[0].Attributes = %+v, want %+v", app.Attributes, wantAttrs)
	}

	api := entities[1]
	if api.Hostname != "" {
		t.Errorf("entities[1].Hostname = %q, want empty for a PathPrefix rule", api.Hostname)
	}
	if api.Attributes["tls"] != false {
		t.Errorf("entities[1].Attributes[tls] = %v, want false", api.Attributes["tls"])
	}
	if _, ok := api.Attributes["certResolver"]; ok {
		t.Errorf("entities[1] should omit certResolver: %+v", api.Attributes)
	}
	if _, ok := api.Attributes["middlewares"]; ok {
		t.Errorf("entities[1] should omit empty middlewares: %+v", api.Attributes)
	}
}

// TestBuildRouterTableEscapesPipes guards the Markdown table against
// Traefik's "||" rule syntax splitting a row into extra columns.
func TestBuildRouterTableEscapesPipes(t *testing.T) {
	raw := []byte("[{\"name\":\"multi@file\",\"rule\":\"Host(`a.example.com`) || Host(`b.example.com`)\",\"service\":\"s@file\",\"status\":\"enabled\"}]")
	content, entities, _ := buildRouterTable(raw)

	if strings.Contains(content, "`) || Host(") {
		t.Errorf("content has an unescaped pipe:\n%s", content)
	}
	if !strings.Contains(content, `\|\|`) {
		t.Errorf("content missing the escaped pipes:\n%s", content)
	}
	row := ""
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "multi@file") {
			row = line
		}
	}
	if got := strings.Count(row, "|") - strings.Count(row, `\|`); got != 8 {
		t.Errorf("row %q has %d unescaped pipes, want 8 (7 columns)", row, got)
	}
	// The raw rule is preserved in attributes, unescaped.
	if entities[0].Attributes["rule"] != "Host(`a.example.com`) || Host(`b.example.com`)" {
		t.Errorf("attributes[rule] = %v, want the raw rule", entities[0].Attributes["rule"])
	}
	if entities[0].Hostname != "a.example.com" {
		t.Errorf("Hostname = %q, want the first literal host", entities[0].Hostname)
	}
}

func TestHostFromRule(t *testing.T) {
	tests := []struct {
		rule string
		want string
	}{
		{"Host(`app.example.com`)", "app.example.com"},
		{"Host(`a.example.com`) && PathPrefix(`/x`)", "a.example.com"},
		{"PathPrefix(`/api`)", ""},
		{"HostRegexp(`^.+\\.example\\.com$`)", ""},
		{"", ""},
		{"Host(``)", ""},
		{"Host(`unterminated", ""},
	}
	for _, tt := range tests {
		if got := hostFromRule(tt.rule); got != tt.want {
			t.Errorf("hostFromRule(%q) = %q, want %q", tt.rule, got, tt.want)
		}
	}
}

func TestBuildServiceTable(t *testing.T) {
	content, entities := buildServiceTable([]byte(servicesJSON))

	if !strings.Contains(content, "1/2 UP") {
		t.Errorf("content missing the health summary:\n%s", content)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}

	wantApp := map[string]any{
		"status":         "enabled",
		"provider":       "docker",
		"serverCount":    2,
		"servers":        []string{"http://10.0.0.5:8080", "http://10.0.0.6:8080"},
		"passHostHeader": true,
	}
	if !reflect.DeepEqual(entities[0].Attributes, wantApp) {
		t.Errorf("entities[0].Attributes = %+v, want %+v", entities[0].Attributes, wantApp)
	}

	// A service without a loadBalancer block still yields a usable entity.
	wantAPI := map[string]any{"status": "enabled", "provider": "internal", "serverCount": 0, "passHostHeader": false}
	if !reflect.DeepEqual(entities[1].Attributes, wantAPI) {
		t.Errorf("entities[1].Attributes = %+v, want %+v", entities[1].Attributes, wantAPI)
	}
}

// TestBuildServiceTableServerVariants covers TCP-style "address" servers and
// the passHostHeader default (absent means true in Traefik).
func TestBuildServiceTableServerVariants(t *testing.T) {
	raw := []byte(`[
		{"name":"tcp@file","status":"enabled","loadBalancer":{"servers":[{"address":"10.0.0.9:5432"}]}},
		{"name":"explicit-false@file","status":"enabled","loadBalancer":{"passHostHeader":false,"servers":[{"url":"http://10.0.0.1"}]}}
	]`)
	_, entities := buildServiceTable(raw)
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	if !reflect.DeepEqual(entities[0].Attributes["servers"], []string{"10.0.0.9:5432"}) {
		t.Errorf("tcp servers = %+v", entities[0].Attributes["servers"])
	}
	if entities[0].Attributes["passHostHeader"] != true {
		t.Errorf("passHostHeader = %v, want true when the field is absent", entities[0].Attributes["passHostHeader"])
	}
	if entities[1].Attributes["passHostHeader"] != false {
		t.Errorf("passHostHeader = %v, want false when explicitly disabled", entities[1].Attributes["passHostHeader"])
	}
}

func TestHealthSummary(t *testing.T) {
	tests := []struct {
		name   string
		status map[string]string
		want   string
	}{
		{name: "none reported", status: nil, want: ""},
		{name: "all up", status: map[string]string{"a": "UP", "b": "up"}, want: "2/2 UP"},
		{name: "mixed", status: map[string]string{"a": "UP", "b": "DOWN"}, want: "1/2 UP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := healthSummary(tt.status); got != tt.want {
				t.Errorf("healthSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestBuildMiddlewareTable covers both the Traefik v2 shape (kind inferred
// from the extra config key) and the v3 shape (explicit "type").
func TestBuildMiddlewareTable(t *testing.T) {
	content, entities := buildMiddlewareTable([]byte(middlewaresJSON))

	if !strings.Contains(content, "basicAuth") || !strings.Contains(content, "redirectScheme") {
		t.Errorf("content missing middleware types:\n%s", content)
	}
	if !strings.Contains(content, "app@docker") {
		t.Errorf("content missing the usedBy column:\n%s", content)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	wantV2 := map[string]any{"status": "enabled", "provider": "file", "types": []string{"basicAuth"}}
	if !reflect.DeepEqual(entities[0].Attributes, wantV2) {
		t.Errorf("v2 middleware attributes = %+v, want %+v", entities[0].Attributes, wantV2)
	}
	wantV3 := map[string]any{"status": "enabled", "provider": "docker", "types": []string{"redirectScheme"}}
	if !reflect.DeepEqual(entities[1].Attributes, wantV3) {
		t.Errorf("v3 middleware attributes = %+v, want %+v", entities[1].Attributes, wantV3)
	}
	if entities[0].Kind != "middleware" || entities[0].ExternalID != "auth@file" {
		t.Errorf("entities[0] identity = %+v", entities[0])
	}
}

// TestBuildMiddlewareTableMultipleKinds checks the v2 fallback sorts the
// inferred kinds so the snapshot is stable across syncs.
func TestBuildMiddlewareTableMultipleKinds(t *testing.T) {
	raw := []byte(`[{"name":"chain@file","status":"enabled","provider":"file","headers":{},"basicAuth":{},"error":""}]`)
	_, entities := buildMiddlewareTable(raw)
	if len(entities) != 1 {
		t.Fatalf("entities = %d, want 1", len(entities))
	}
	if !reflect.DeepEqual(entities[0].Attributes["types"], []string{"basicAuth", "headers"}) {
		t.Errorf("types = %+v, want [basicAuth headers]", entities[0].Attributes["types"])
	}
}

func TestBuildEntryPointTable(t *testing.T) {
	content, entities := buildEntryPointTable([]byte(entryPointsJSON))

	if !strings.Contains(content, ":443") || !strings.Contains(content, "yes (letsencrypt)") {
		t.Errorf("content missing entry point details:\n%s", content)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	wantWeb := map[string]any{"address": ":80", "asDefault": false, "tls": false}
	if !reflect.DeepEqual(entities[0].Attributes, wantWeb) {
		t.Errorf("web attributes = %+v, want %+v", entities[0].Attributes, wantWeb)
	}
	wantSecure := map[string]any{"address": ":443", "asDefault": true, "tls": true}
	if !reflect.DeepEqual(entities[1].Attributes, wantSecure) {
		t.Errorf("websecure attributes = %+v, want %+v", entities[1].Attributes, wantSecure)
	}
}

func TestServiceDependencies(t *testing.T) {
	got := serviceDependencies([]string{"b@docker", "a@file", "b@docker", ""})
	want := []connector.ServiceDependency{
		{Kind: "upstream_service", Name: "a@file"},
		{Kind: "upstream_service", Name: "b@docker"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("serviceDependencies() = %+v, want %+v", got, want)
	}
	if got := serviceDependencies(nil); got != nil {
		t.Errorf("serviceDependencies(nil) = %+v, want nil", got)
	}
	if got := serviceDependencies([]string{"", ""}); got != nil {
		t.Errorf("serviceDependencies(empty names) = %+v, want nil", got)
	}
}

func TestCell(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", "—"},
		{"plain", "plain"},
		{"a|b", `a\|b`},
		{"line\nbreak", "line break"},
	}
	for _, tt := range tests {
		if got := cell(tt.in); got != tt.want {
			t.Errorf("cell(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
