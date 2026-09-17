package opnsense

import (
	"reflect"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildInterfaceTableAttributes(t *testing.T) {
	data := []byte(`{"rows":[
		{"device":"igb0","ipaddr":"203.0.113.5","ipv6":"2001:db8::1","status":"up","media":"1000baseT","enabled":true,"type":"static","gateway":"WAN_GW"},
		{"device":"igb1","ipaddr":"10.0.0.1","status":"up","media":"1000baseT","enabled":false,"type":"static"}
	]}`)
	_, entities := buildInterfaceTable(data)
	if len(entities) != 2 {
		t.Fatalf("entities = %+v, want 2", entities)
	}
	want := []map[string]any{
		{"enabled": true, "ipv4": "203.0.113.5", "ipv6": "2001:db8::1", "type": "static", "gateway": "WAN_GW"},
		{"enabled": false, "ipv4": "10.0.0.1", "ipv6": "", "type": "static"},
	}
	for i, w := range want {
		if !reflect.DeepEqual(entities[i].Attributes, w) {
			t.Errorf("entities[%d].Attributes = %+v, want %+v", i, entities[i].Attributes, w)
		}
	}
}

func TestBuildRuleTableAttributes(t *testing.T) {
	data := []byte(`{"rows":[
		{"description":"Allow SSH","action":"pass","protocol":"tcp","source_net":"any","destination_net":"any","destination_port":"22","interface":"wan","direction":"in","enabled":"1","log":"1"},
		{"description":"Block DNS","action":"block","protocol":"udp","source_net":"10.0.0.0/8","destination_net":"any","enabled":"0","disabled_reason":"maintenance"}
	]}`)
	_, entities := buildRuleTable(data)
	if len(entities) != 2 {
		t.Fatalf("entities = %+v, want 2", entities)
	}

	sshAttrs := entities[0].Attributes
	if sshAttrs["enabled"] != true || sshAttrs["action"] != "pass" || sshAttrs["protocol"] != "tcp" ||
		sshAttrs["source"] != "any" || sshAttrs["destination"] != "any" || sshAttrs["destination_port"] != "22" ||
		sshAttrs["interface"] != "wan" || sshAttrs["direction"] != "in" || sshAttrs["log"] != true {
		t.Errorf("Allow SSH attributes = %+v", sshAttrs)
	}
	if _, ok := sshAttrs["disabled_reason"]; ok {
		t.Errorf("Allow SSH attributes should omit disabled_reason when absent: %+v", sshAttrs)
	}

	dnsAttrs := entities[1].Attributes
	if dnsAttrs["enabled"] != false || dnsAttrs["action"] != "block" || dnsAttrs["disabled_reason"] != "maintenance" || dnsAttrs["log"] != false {
		t.Errorf("Block DNS attributes = %+v", dnsAttrs)
	}
	if _, ok := dnsAttrs["interface"]; ok {
		t.Errorf("Block DNS attributes should omit interface when absent: %+v", dnsAttrs)
	}
}

// TestAttributeCatalogCoversEmittedKeys ensures every attribute key this
// connector emits for interface/rule entities is declared in its catalog
// with a matching type.
func TestAttributeCatalogCoversEmittedKeys(t *testing.T) {
	catalog := attributeCatalog
	emitted := map[string]map[string]string{
		"interface": {"enabled": "boolean", "type": "string", "ipv4": "string", "ipv6": "string", "gateway": "string"},
		"rule": {
			"enabled": "boolean", "action": "string", "interface": "string", "direction": "string",
			"protocol": "string", "source": "string", "destination": "string", "destination_port": "string",
			"log": "boolean", "disabled_reason": "string",
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
