package portainer

import (
	"encoding/json"
	"strings"
	"testing"
)

func decodeEnvironments(t *testing.T, raw string) []environment {
	t.Helper()
	var envs []environment
	if err := json.Unmarshal([]byte(raw), &envs); err != nil {
		t.Fatalf("decode environments: %v", err)
	}
	return envs
}

func TestBuildEnvironmentTable(t *testing.T) {
	content, entities := buildEnvironmentTable(decodeEnvironments(t, endpointsJSON))

	for _, want := range []string{
		"| local | docker | up | unix:///var/run/docker.sock | 24.0.7 | false |",
		"| edge-nas | agent | up | tcp://10.0.0.5:9001/ | 25.0.3 | true |",
		"| offline-agent | agent | down | tcp://10.0.0.9:9001 | — | false |",
		"| k8s-prod | kubernetes | up | https://10.0.0.20:6443 | — | false |",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing row %q:\n%s", want, content)
		}
	}

	if len(entities) != 4 {
		t.Fatalf("entities = %d, want 4", len(entities))
	}
	local := entities[0]
	if local.ExternalID != "1" || local.Name != "local" {
		t.Errorf("entities[0] = %+v, want the local environment", local)
	}
	if local.Hostname != "" {
		t.Errorf("entities[0].Hostname = %q, want empty for a unix socket URL", local.Hostname)
	}
	wantAttrs := map[string]any{
		"type": "docker", "status": "up", "tls": false, "swarm": false, "groupId": 1,
		"url": "unix:///var/run/docker.sock", "publicUrl": "portainer.example.com", "dockerVersion": "24.0.7",
	}
	for key, want := range wantAttrs {
		if got := local.Attributes[key]; got != want {
			t.Errorf("entities[0].Attributes[%q] = %v, want %v", key, got, want)
		}
	}
	tags, ok := local.Attributes["tags"].([]string)
	if !ok || len(tags) != 2 || tags[0] != "3" || tags[1] != "4" {
		t.Errorf("entities[0].Attributes[\"tags\"] = %v, want [3 4]", local.Attributes["tags"])
	}
	if entities[1].Hostname != "10.0.0.5" {
		t.Errorf("entities[1].Hostname = %q, want %q", entities[1].Hostname, "10.0.0.5")
	}
	if _, ok := entities[2].Attributes["dockerVersion"]; ok {
		t.Errorf("offline-agent should carry no dockerVersion: %+v", entities[2].Attributes)
	}
	if _, ok := entities[2].Attributes["tags"]; ok {
		t.Errorf("offline-agent should carry no tags: %+v", entities[2].Attributes)
	}
}

func TestBuildEnvironmentTableEmpty(t *testing.T) {
	content, entities := buildEnvironmentTable(nil)
	if content != "_No environments returned_" || entities != nil {
		t.Fatalf("buildEnvironmentTable(nil) = %q, %+v", content, entities)
	}
}

func TestEnvironmentDockerCapability(t *testing.T) {
	tests := []struct {
		name string
		env  environment
		want bool
	}{
		{name: "local docker", env: environment{Type: envTypeDocker, Status: statusUp}, want: true},
		{name: "agent", env: environment{Type: envTypeAgent, Status: statusUp}, want: true},
		{name: "edge agent", env: environment{Type: envTypeEdgeAgent, Status: statusUp}, want: true},
		{name: "down docker", env: environment{Type: envTypeDocker, Status: statusDown}},
		{name: "azure", env: environment{Type: envTypeAzure, Status: statusUp}},
		{name: "kubernetes", env: environment{Type: envTypeKubernetesLocal, Status: statusUp}},
		{name: "kubernetes agent", env: environment{Type: envTypeKubernetesAgent, Status: statusUp}},
		{name: "kubernetes edge", env: environment{Type: envTypeKubernetesEdge, Status: statusUp}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.dockerCapable(); got != tt.want {
				t.Fatalf("dockerCapable() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestHostFromURL(t *testing.T) {
	tests := map[string]string{
		"tcp://10.0.0.5:9001":            "10.0.0.5",
		"tcp://10.0.0.5:9001/":           "10.0.0.5",
		"https://portainer.example.com":  "portainer.example.com",
		"https://10.0.0.20:6443/api":     "10.0.0.20",
		"unix:///var/run/docker.sock":    "",
		"npipe:////./pipe/docker_engine": "",
		"":                               "",
	}
	for in, want := range tests {
		if got := hostFromURL(in); got != want {
			t.Errorf("hostFromURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEnumNames(t *testing.T) {
	if got := environmentTypeName(42); got != "unknown(42)" {
		t.Errorf("environmentTypeName(42) = %q", got)
	}
	if got := statusName(9); got != "unknown(9)" {
		t.Errorf("statusName(9) = %q", got)
	}
	if got := stackTypeName(9); got != "unknown(9)" {
		t.Errorf("stackTypeName(9) = %q", got)
	}
	if got := stackStatusName(9); got != "unknown(9)" {
		t.Errorf("stackStatusName(9) = %q", got)
	}
}

func TestBuildStackTable(t *testing.T) {
	envs := decodeEnvironments(t, endpointsJSON)
	content, entities := buildStackTable([]byte(stacksJSON), environmentNames(envs))

	for _, want := range []string{
		"| blog | local | compose | active | docker-compose.yml |",
		"| metrics | edge-nas | swarm | inactive | stack.yml |",
		"| orphan | endpoint 99 | compose | active | — |",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing row %q:\n%s", want, content)
		}
	}
	if len(entities) != 3 {
		t.Fatalf("entities = %d, want 3", len(entities))
	}
	if entities[0].ExternalID != "7" || entities[0].Attributes["environment"] != "local" {
		t.Errorf("entities[0] = %+v", entities[0])
	}
	if _, ok := entities[2].Attributes["entryPoint"]; ok {
		t.Errorf("orphan stack should carry no entryPoint: %+v", entities[2].Attributes)
	}
}

func TestBuildStackTableEdgeCases(t *testing.T) {
	content, entities := buildStackTable([]byte("[]"), nil)
	if content != "_No stacks returned_" || entities != nil {
		t.Fatalf("empty stacks = %q, %+v", content, entities)
	}
	content, entities = buildStackTable([]byte("not json"), nil)
	if !strings.Contains(content, "malformed response") || entities != nil {
		t.Fatalf("malformed stacks = %q, %+v", content, entities)
	}
}

func TestContainerRows(t *testing.T) {
	rows, entities, err := containerRows("local", []byte(containersJSON))
	if err != nil {
		t.Fatalf("containerRows: %v", err)
	}
	if !strings.Contains(rows, "| local | blog-web | ghcr.io/example/blog:1.2 | running | Up 3 days | 8080:80/tcp |") {
		t.Errorf("rows missing the blog-web row:\n%s", rows)
	}
	if !strings.Contains(rows, "| local | standalone | alpine:3.19 | exited | Exited (0) 2 hours ago | — |") {
		t.Errorf("rows missing the standalone row:\n%s", rows)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	web := entities[0]
	if web.ExternalID != "c0ffee11" || web.IP != "172.20.0.4" {
		t.Errorf("entities[0] = %+v", web)
	}
	if web.Attributes["stack"] != "blog" || web.Attributes["network_mode"] != "blog_default" {
		t.Errorf("entities[0].Attributes = %+v", web.Attributes)
	}
	ports, ok := web.Attributes["published_ports"].([]string)
	if !ok || len(ports) != 1 || ports[0] != "8080:80/tcp" {
		t.Errorf("published_ports = %v, want only the published mapping", web.Attributes["published_ports"])
	}
	if _, ok := entities[1].Attributes["stack"]; ok {
		t.Errorf("unlabelled container should carry no stack: %+v", entities[1].Attributes)
	}
	if entities[1].IP != "" {
		t.Errorf("entities[1].IP = %q, want empty", entities[1].IP)
	}
}

func TestContainerRowsMalformed(t *testing.T) {
	if _, _, err := containerRows("local", []byte("not json")); err == nil {
		t.Fatal("containerRows(not json) error = nil, want a decode error")
	}
}

func TestVolumeRows(t *testing.T) {
	rows, entities, err := volumeRows("edge-nas", []byte(volumesJSON))
	if err != nil {
		t.Fatalf("volumeRows: %v", err)
	}
	if !strings.Contains(rows, "| edge-nas | blog_data | local | /var/lib/docker/volumes/blog_data/_data |") {
		t.Errorf("rows missing the blog_data row:\n%s", rows)
	}
	if len(entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(entities))
	}
	if entities[0].ExternalID != "edge-nas/blog_data" {
		t.Errorf("entities[0].ExternalID = %q, want the environment-scoped id", entities[0].ExternalID)
	}
	if entities[0].Attributes["environment"] != "edge-nas" {
		t.Errorf("entities[0].Attributes = %+v", entities[0].Attributes)
	}
	if _, _, err := volumeRows("edge-nas", []byte("not json")); err == nil {
		t.Fatal("volumeRows(not json) error = nil, want a decode error")
	}
}

func TestNetworkRows(t *testing.T) {
	rows, entities, err := networkRows("local", []byte(networksJSON))
	if err != nil {
		t.Fatalf("networkRows: %v", err)
	}
	if !strings.Contains(rows, "| local | blog_default | bridge | local | true |") {
		t.Errorf("rows missing the blog_default row:\n%s", rows)
	}
	if len(entities) != 2 || entities[0].ExternalID != "net1" {
		t.Fatalf("entities = %+v", entities)
	}
	if entities[1].Attributes["internal"] != true {
		t.Errorf("entities[1].Attributes = %+v", entities[1].Attributes)
	}

	// A network without an Id (some older engines) falls back to an
	// environment-scoped name so it still has a stable identity.
	_, fallback, err := networkRows("local", []byte(`[{"Name":"legacy","Driver":"bridge"}]`))
	if err != nil {
		t.Fatalf("networkRows: %v", err)
	}
	if len(fallback) != 1 || fallback[0].ExternalID != "local/legacy" {
		t.Fatalf("fallback entities = %+v", fallback)
	}
	if _, _, err := networkRows("local", []byte("not json")); err == nil {
		t.Fatal("networkRows(not json) error = nil, want a decode error")
	}
}

// TestCellEscapesPipes guards table rows against values containing '|',
// which would otherwise split a row into extra columns.
func TestCellEscapesPipes(t *testing.T) {
	if got := cell("sh -c 'a | b'"); got != `sh -c 'a \| b'` {
		t.Errorf("cell() = %q", got)
	}
	if got := cell("two\nlines"); got != "two lines" {
		t.Errorf("cell() = %q", got)
	}
	if got := cell(""); got != "—" {
		t.Errorf("cell(\"\") = %q", got)
	}
}
