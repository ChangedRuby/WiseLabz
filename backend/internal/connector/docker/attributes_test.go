package docker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildContainerTableAttributes(t *testing.T) {
	data := []byte(`[
		{"Id":"c1","Names":["/web"],"Image":"nginx:latest","State":"running","Status":"Up",
		 "Ports":[{"PrivatePort":80,"PublicPort":8080,"Type":"tcp"},{"PrivatePort":53,"Type":"udp"}],
		 "HostConfig":{"NetworkMode":"bridge"}},
		{"Id":"c2","Names":["/worker"],"Image":"worker:v1","State":"running","Status":"Up",
		 "HostConfig":{"NetworkMode":"host"}},
		{"Names":["/no-id"],"Image":"scratch","State":"exited","Status":"Exited"}
	]`)
	_, entities := buildContainerTable(data)
	if len(entities) != 3 {
		t.Fatalf("entities = %+v, want 3", entities)
	}

	want0 := map[string]any{
		"image":           "nginx:latest",
		"network_mode":    "bridge",
		"published_ports": []string{"8080:80/tcp"},
	}
	if !reflect.DeepEqual(entities[0].Attributes, want0) {
		t.Errorf("entities[0].Attributes = %+v, want %+v", entities[0].Attributes, want0)
	}

	want1 := map[string]any{"image": "worker:v1", "network_mode": "host"}
	if !reflect.DeepEqual(entities[1].Attributes, want1) {
		t.Errorf("entities[1].Attributes = %+v, want %+v", entities[1].Attributes, want1)
	}
	if _, ok := entities[1].Attributes["published_ports"]; ok {
		t.Errorf("entities[1] should omit published_ports when no port is published: %+v", entities[1].Attributes)
	}

	if entities[2].ExternalID != "" {
		t.Errorf("entities[2].ExternalID = %q, want empty (no Id in list response)", entities[2].ExternalID)
	}
}

func TestEnrichContainerAttributesMergesInspectData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/containers/privileged-1/json":
			_, _ = w.Write([]byte(`{"Config":{"User":"root"},"HostConfig":{"Privileged":true,"ReadonlyRootfs":true,"RestartPolicy":{"Name":"always"}}}`))
		case "/containers/plain-1/json":
			_, _ = w.Write([]byte(`{"Config":{"User":""},"HostConfig":{"Privileged":false,"ReadonlyRootfs":false,"RestartPolicy":{"Name":""}}}`))
		case "/containers/broken-1/json":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	ents := []connector.SnapshotEntity{
		{Kind: "container", Name: "priv", ExternalID: "privileged-1", Attributes: map[string]any{"image": "img1"}},
		{Kind: "container", Name: "plain", ExternalID: "plain-1", Attributes: map[string]any{"image": "img2"}},
		{Kind: "container", Name: "broken", ExternalID: "broken-1", Attributes: map[string]any{"image": "img3"}},
		{Kind: "container", Name: "no-id", ExternalID: ""},
	}
	c.enrichContainerAttributes(context.Background(), ents)

	wantPriv := map[string]any{
		"image":            "img1",
		"privileged":       true,
		"read_only_rootfs": true,
		"restart_policy":   "always",
		"user":             "root",
	}
	if !reflect.DeepEqual(ents[0].Attributes, wantPriv) {
		t.Errorf("privileged container attributes = %+v, want %+v", ents[0].Attributes, wantPriv)
	}

	wantPlain := map[string]any{
		"image":            "img2",
		"privileged":       false,
		"read_only_rootfs": false,
	}
	if !reflect.DeepEqual(ents[1].Attributes, wantPlain) {
		t.Errorf("plain container attributes = %+v, want %+v", ents[1].Attributes, wantPlain)
	}
	if _, ok := ents[1].Attributes["restart_policy"]; ok {
		t.Errorf("plain container should omit restart_policy when empty: %+v", ents[1].Attributes)
	}
	if _, ok := ents[1].Attributes["user"]; ok {
		t.Errorf("plain container should omit user when empty: %+v", ents[1].Attributes)
	}

	// Inspect call failed: original attributes untouched, no partial merge.
	wantBroken := map[string]any{"image": "img3"}
	if !reflect.DeepEqual(ents[2].Attributes, wantBroken) {
		t.Errorf("broken container attributes = %+v, want %+v (inspect failure tolerated)", ents[2].Attributes, wantBroken)
	}

	if ents[3].Attributes != nil {
		t.Errorf("no-id container attributes = %+v, want nil (no inspect call made)", ents[3].Attributes)
	}
}

func TestFetchEnrichesContainersFromInspectEndToEnd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/info":
			_, _ = w.Write([]byte(`{"Name":"docker-host"}`))
		case "/containers/json":
			_, _ = w.Write([]byte(`[{"Id":"abc123","Names":["/web"],"Image":"nginx","State":"running","Status":"Up",
				"Ports":[{"PrivatePort":80,"PublicPort":8080,"Type":"tcp"}],"HostConfig":{"NetworkMode":"bridge"}}]`))
		case "/containers/abc123/json":
			_, _ = w.Write([]byte(`{"Config":{"User":"appuser"},"HostConfig":{"Privileged":false,"ReadonlyRootfs":true,"RestartPolicy":{"Name":"unless-stopped"}}}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{host: "tcp://example", baseURL: server.URL, client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"containers"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(snap.Entities) != 1 {
		t.Fatalf("Entities = %+v, want 1", snap.Entities)
	}
	want := map[string]any{
		"image":            "nginx",
		"network_mode":     "bridge",
		"published_ports":  []string{"8080:80/tcp"},
		"privileged":       false,
		"read_only_rootfs": true,
		"restart_policy":   "unless-stopped",
		"user":             "appuser",
	}
	if !reflect.DeepEqual(snap.Entities[0].Attributes, want) {
		t.Errorf("Entities[0].Attributes = %+v, want %+v", snap.Entities[0].Attributes, want)
	}
}

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector emits for container entities is declared in its catalog with a
// matching type, so PR3's rule engine and the schema endpoint never drift
// from what Fetch actually produces.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	catalog := attributeCatalog
	emitted := map[string]map[string]string{
		"container": {
			"privileged":       "boolean",
			"network_mode":     "string",
			"restart_policy":   "string",
			"published_ports":  "string_array",
			"user":             "string",
			"read_only_rootfs": "boolean",
			"image":            "string",
		},
	}
	for kind, keys := range emitted {
		specs, ok := catalog[kind]
		if !ok {
			t.Fatalf("catalog missing entity kind %q", kind)
		}
		byName := make(map[string]connector.AttributeSpec, len(specs))
		for _, s := range specs {
			byName[s.Name] = s
		}
		for key, typ := range keys {
			spec, ok := byName[key]
			if !ok {
				t.Errorf("catalog[%q] missing emitted attribute %q", kind, key)
				continue
			}
			if spec.Type != typ {
				t.Errorf("catalog[%q][%q].Type = %q, want %q", kind, key, spec.Type, typ)
			}
		}
	}
}
