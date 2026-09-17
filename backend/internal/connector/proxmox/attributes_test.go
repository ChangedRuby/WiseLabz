package proxmox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// TestFetchEntityAttributes covers running vs stopped VMs, a template VM,
// and privileged vs unprivileged containers, asserting the exact Attributes
// map built from /config and /firewall/options responses.
func TestFetchEntityAttributes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nodes":
			_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
		case "/nodes/pve1/qemu":
			_, _ = w.Write([]byte(`{"data":[
				{"vmid":100,"name":"web1","status":"running","cpus":2,"mem":1024,"uptime":10},
				{"vmid":101,"name":"tmpl1","status":"stopped","cpus":1,"mem":512,"uptime":0}
			]}`))
		case "/nodes/pve1/lxc":
			_, _ = w.Write([]byte(`{"data":[
				{"vmid":200,"name":"ct1","status":"running","cpus":1,"mem":256,"uptime":5},
				{"vmid":201,"name":"ct2","status":"stopped","cpus":1,"mem":256,"uptime":0}
			]}`))
		case "/nodes/pve1/qemu/100/config":
			_, _ = w.Write([]byte(`{"data":{"onboot":1,"protection":0,"agent":"enabled=1,fstrim_cloned_disks=1","template":0,"ostype":"l26"}}`))
		case "/nodes/pve1/qemu/100/agent/network-get-interfaces":
			_, _ = w.Write([]byte(`{"data":{"result":[]}}`))
		case "/nodes/pve1/qemu/100/firewall/options":
			_, _ = w.Write([]byte(`{"data":{"enable":1}}`))
		case "/nodes/pve1/qemu/101/config":
			_, _ = w.Write([]byte(`{"data":{"onboot":0,"protection":1,"template":1,"ostype":"win10"}}`))
		case "/nodes/pve1/qemu/101/firewall/options":
			_, _ = w.Write([]byte(`{"data":{"enable":0}}`))
		case "/nodes/pve1/lxc/200/config":
			_, _ = w.Write([]byte(`{"data":{"onboot":1,"protection":0,"unprivileged":1,"ostype":"debian"}}`))
		case "/nodes/pve1/lxc/200/interfaces":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/nodes/pve1/lxc/200/firewall/options":
			_, _ = w.Write([]byte(`{"data":{"enable":1}}`))
		case "/nodes/pve1/lxc/201/config":
			_, _ = w.Write([]byte(`{"data":{"onboot":0,"protection":0,"unprivileged":0,"ostype":"alpine"}}`))
		case "/nodes/pve1/lxc/201/firewall/options":
			_, _ = w.Write([]byte(`{"data":{"enable":0}}`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"vms", "containers", "entities"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	byRef := make(map[string]connector.SnapshotEntity, len(snap.Entities))
	for _, e := range snap.Entities {
		byRef[e.Kind+"/"+e.ExternalID] = e
	}

	wantWeb1 := map[string]any{
		"status":           "running",
		"onboot":           true,
		"protection":       false,
		"agent_enabled":    true,
		"template":         false,
		"os_type":          "l26",
		"firewall_enabled": true,
	}
	if got := byRef["vm/100"].Attributes; !reflect.DeepEqual(got, wantWeb1) {
		t.Errorf("vm/100 Attributes = %+v, want %+v", got, wantWeb1)
	}

	wantTmpl1 := map[string]any{
		"status":           "stopped",
		"onboot":           false,
		"protection":       true,
		"agent_enabled":    false,
		"template":         true,
		"os_type":          "win10",
		"firewall_enabled": false,
	}
	if got := byRef["vm/101"].Attributes; !reflect.DeepEqual(got, wantTmpl1) {
		t.Errorf("vm/101 (template) Attributes = %+v, want %+v", got, wantTmpl1)
	}

	wantCt1 := map[string]any{
		"status":           "running",
		"onboot":           true,
		"protection":       false,
		"template":         false,
		"agent_enabled":    false,
		"unprivileged":     true,
		"os_type":          "debian",
		"firewall_enabled": true,
	}
	if got := byRef["container/200"].Attributes; !reflect.DeepEqual(got, wantCt1) {
		t.Errorf("container/200 Attributes = %+v, want %+v", got, wantCt1)
	}

	wantCt2 := map[string]any{
		"status":           "stopped",
		"onboot":           false,
		"protection":       false,
		"template":         false,
		"agent_enabled":    false,
		"unprivileged":     false,
		"os_type":          "alpine",
		"firewall_enabled": false,
	}
	if got := byRef["container/201"].Attributes; !reflect.DeepEqual(got, wantCt2) {
		t.Errorf("container/201 (privileged) Attributes = %+v, want %+v", got, wantCt2)
	}
}

// TestFetchEntityAttributesOmittedWithoutEntitiesField ensures that without
// the "entities" fields hint, only the always-available "status" attribute
// is set — no per-guest /config or /firewall/options calls are made.
func TestFetchEntityAttributesOmittedWithoutEntitiesField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nodes":
			_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
		case "/nodes/pve1/qemu":
			_, _ = w.Write([]byte(`{"data":[{"vmid":100,"name":"web1","status":"running","cpus":2,"mem":1024,"uptime":10}]}`))
		case "/nodes/pve1/lxc":
			_, _ = w.Write([]byte(`{"data":[]}`))
		default:
			t.Fatalf("unexpected request path %s: entities not requested, no per-guest calls expected", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"vms", "containers"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(snap.Entities) != 1 {
		t.Fatalf("Entities = %+v, want 1", snap.Entities)
	}
	want := map[string]any{"status": "running"}
	if got := snap.Entities[0].Attributes; !reflect.DeepEqual(got, want) {
		t.Errorf("Attributes = %+v, want %+v", got, want)
	}
}

// TestFetchEntityAttributesDegradesOnConfigError ensures a /config or
// /firewall/options failure only drops those attributes, not the whole
// Fetch or entity.
func TestFetchEntityAttributesDegradesOnConfigError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/nodes":
			_, _ = w.Write([]byte(`{"data":[{"node":"pve1","status":"online","uptime":100,"cpu":0.1,"mem":{"used":1,"total":2}}]}`))
		case "/nodes/pve1/qemu":
			_, _ = w.Write([]byte(`{"data":[{"vmid":100,"name":"web1","status":"stopped","cpus":2,"mem":1024,"uptime":0}]}`))
		case "/nodes/pve1/lxc":
			_, _ = w.Write([]byte(`{"data":[]}`))
		case "/nodes/pve1/qemu/100/config":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`error`))
		case "/nodes/pve1/qemu/100/firewall/options":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`error`))
		default:
			t.Fatalf("unexpected request path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c := &Connector{url: server.URL, tokenID: "user@pam!token", tokenSecret: "secret", client: server.Client()}
	snap, err := c.Fetch(context.Background(), map[string]any{"fields": []any{"vms", "containers", "entities"}})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(snap.Entities) != 1 {
		t.Fatalf("Entities = %+v, want 1", snap.Entities)
	}
	want := map[string]any{"status": "stopped"}
	if got := snap.Entities[0].Attributes; !reflect.DeepEqual(got, want) {
		t.Errorf("Attributes = %+v, want %+v (config/firewall errors should be omitted, not fail Fetch)", got, want)
	}
}

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector emits for vm/container entities is declared in its catalog with
// a matching type, so PR3's rule engine and the schema endpoint never drift
// from what Fetch actually produces.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	catalog := attributeCatalog
	emitted := map[string]map[string]string{
		"vm": {
			"status": "string", "firewall_enabled": "boolean", "onboot": "boolean",
			"agent_enabled": "boolean", "protection": "boolean", "template": "boolean", "os_type": "string",
		},
		"container": {
			"status": "string", "firewall_enabled": "boolean", "onboot": "boolean",
			"agent_enabled": "boolean", "protection": "boolean", "template": "boolean", "os_type": "string",
			"unprivileged": "boolean",
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
