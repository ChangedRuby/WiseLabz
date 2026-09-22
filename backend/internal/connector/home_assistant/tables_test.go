package homeassistant

import (
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

const (
	configJSON = `{
		"version":"2024.6.4",
		"location_name":"Home | Lab",
		"time_zone":"Europe/Lisbon",
		"country":"PT",
		"language":"en",
		"currency":"EUR",
		"elevation":42,
		"latitude":38.7223,
		"longitude":-9.1393,
		"config_dir":"/config",
		"config_source":"storage",
		"state":"RUNNING",
		"safe_mode":false,
		"recovery_mode":false,
		"internal_url":"http://homeassistant.local:8123",
		"external_url":"https://ha.example.com",
		"unit_system":{"length":"km","mass":"g","temperature":"°C","volume":"L","pressure":"Pa","wind_speed":"km/h"},
		"components":["light","light.hue","mqtt","sensor","sensor.mqtt","zha"]
	}`
	statesJSON = `[
		{"entity_id":"light.living_room","state":"on","last_changed":"2024-06-01T10:00:00+00:00","last_updated":"2024-06-01T10:00:00+00:00",
		 "attributes":{"friendly_name":"Living Room","supported_features":44}},
		{"entity_id":"sensor.outdoor_temperature","state":"21.4","last_updated":"2024-06-02T11:00:00+00:00",
		 "attributes":{"friendly_name":"Outdoor Temperature","device_class":"temperature","unit_of_measurement":"°C","entity_category":"diagnostic"}},
		{"entity_id":"sensor.humidity","state":"unknown","attributes":{"friendly_name":"Humidity","device_class":"humidity","unit_of_measurement":"%"}},
		{"entity_id":"binary_sensor.front_door","state":"off","attributes":{"friendly_name":"Front Door","device_class":"door"}},
		{"entity_id":"device_tracker.phone","state":"home","attributes":{"friendly_name":"Phone","ip":"192.168.1.50","host_name":"phone-01","source_type":"router"}},
		{"entity_id":"device_tracker.tablet","state":"not_home","attributes":{"friendly_name":"Tablet","ip":"not-an-ip"}},
		{"entity_id":"switch.pump","state":"unavailable","attributes":{"friendly_name":"Pump"}},
		{"entity_id":"automation.night","state":"on","attributes":{"friendly_name":"Night"}},
		{"entity_id":"sun.sun","state":"above_horizon","attributes":{"friendly_name":"Sun"}},
		{"entity_id":"text.note","state":"2024-06-01T10:00:00Z","attributes":{"friendly_name":"Note"}},
		{"entity_id":"zone.home","state":"1","attributes":{"friendly_name":"Zone"}}
	]`
	servicesJSON = `[
		{"domain":"light","services":{"turn_on":{},"turn_off":{},"toggle":{}}},
		{"domain":"homeassistant","services":{
			"s01":{},"s02":{},"s03":{},"s04":{},"s05":{},"s06":{},"s07":{},"s08":{},"s09":{},
			"s10":{},"s11":{},"s12":{},"s13":{},"s14":{},"s15":{},"s16":{},"s17":{}}}
	]`
)

func TestBuildOverview(t *testing.T) {
	content, metadata := buildOverview([]byte(configJSON))

	for _, want := range []string{
		"| Version | 2024.6.4 |",
		`| Location | Home \| Lab |`,
		"| Time zone | Europe/Lisbon |",
		"| Elevation | 42 m |",
		"| Unit system | °C · km · g |",
		"| Safe mode | false |",
		"| Integrations loaded | 6 |",
		"WebSocket-only",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("overview missing %q:\n%s", want, content)
		}
	}
	// Precise home coordinates must never reach a snapshot.
	for _, leak := range []string{"38.7223", "-9.1393", "latitude"} {
		if strings.Contains(content, leak) {
			t.Errorf("overview leaks %q:\n%s", leak, content)
		}
	}

	want := map[string]string{
		"home_assistant_version": "2024.6.4",
		"home_assistant_state":   "RUNNING",
		"location_name":          "Home | Lab",
		"time_zone":              "Europe/Lisbon",
		"config_source":          "storage",
		"component_count":        "6",
	}
	for k, v := range want {
		if metadata[k] != v {
			t.Errorf("metadata[%q] = %q, want %q", k, metadata[k], v)
		}
	}
}

func TestBuildOverviewMalformed(t *testing.T) {
	content, metadata := buildOverview([]byte("not json"))
	if !strings.Contains(content, "malformed response") {
		t.Errorf("content = %q, want a malformed-response placeholder", content)
	}
	if metadata != nil {
		t.Errorf("metadata = %+v, want nil", metadata)
	}
}

func TestBuildIntegrations(t *testing.T) {
	content, entities := buildIntegrations([]byte(configJSON))

	if len(entities) != 5 {
		t.Fatalf("entities = %d, want 5 (hue, light, mqtt, sensor, zha): %+v", len(entities), entities)
	}
	byName := map[string]connector.SnapshotEntity{}
	for _, e := range entities {
		if e.Kind != "integration" {
			t.Errorf("entity %q kind = %q, want %q", e.Name, e.Kind, "integration")
		}
		if e.ExternalID != e.Name {
			t.Errorf("entity %q ExternalID = %q, want the integration name", e.Name, e.ExternalID)
		}
		byName[e.Name] = e
	}
	// "light.hue" means the hue integration loaded the light platform.
	hue, ok := byName["hue"]
	if !ok {
		t.Fatalf("no entity for integration hue: %+v", byName)
	}
	platforms, ok := hue.Attributes["platforms"].([]string)
	if !ok || len(platforms) != 1 || platforms[0] != "light" {
		t.Errorf("hue platforms = %v, want [light]", hue.Attributes["platforms"])
	}
	if _, ok := byName["zha"].Attributes["platforms"]; ok {
		t.Errorf("zha should carry no platforms: %+v", byName["zha"].Attributes)
	}
	if !strings.Contains(content, "| mqtt | sensor |") {
		t.Errorf("content missing the mqtt row:\n%s", content)
	}
	// Rows must be sorted so an unchanged instance renders identically.
	if i, j := strings.Index(content, "| hue |"), strings.Index(content, "| zha |"); i < 0 || j < 0 || i > j {
		t.Errorf("rows are not sorted:\n%s", content)
	}
}

func TestBuildIntegrationsEdgeCases(t *testing.T) {
	if content, entities := buildIntegrations([]byte(`{"components":[]}`)); entities != nil || !strings.Contains(content, "No integrations") {
		t.Errorf("empty components = %q / %+v", content, entities)
	}
	if content, _ := buildIntegrations([]byte("[")); !strings.Contains(content, "malformed response") {
		t.Errorf("content = %q, want a malformed-response placeholder", content)
	}
}

func TestBuildEntities(t *testing.T) {
	got := buildEntities([]byte(statesJSON), 0)

	wantMeta := map[string]string{
		"entity_count":             "11",
		"entity_domain_count":      "9",
		"unavailable_entity_count": "2", // switch.pump (unavailable) + sensor.humidity (unknown)
	}
	for k, v := range wantMeta {
		if got.metadata[k] != v {
			t.Errorf("metadata[%q] = %q, want %q", k, got.metadata[k], v)
		}
	}
	if !strings.Contains(got.domains, "| device_tracker | 2 | 0 |") {
		t.Errorf("domain summary missing device_tracker row:\n%s", got.domains)
	}
	if !strings.Contains(got.domains, "| sensor | 2 | 1 |") {
		t.Errorf("domain summary missing sensor row:\n%s", got.domains)
	}

	// No timestamp and no live measurement may reach the rendered table.
	for _, leak := range []string{"last_changed", "last_updated", "2024-06-01T10:00:00+00:00", "21.4"} {
		if strings.Contains(got.entityTable, leak) {
			t.Errorf("entity table leaks the volatile value %q:\n%s", leak, got.entityTable)
		}
	}
	if !strings.Contains(got.entityTable, "| sensor.outdoor_temperature | Outdoor Temperature | sensor | available | temperature | °C |") {
		t.Errorf("entity table missing the collapsed sensor row:\n%s", got.entityTable)
	}
	if !strings.Contains(got.entityTable, "| light.living_room | Living Room | light | on | — | — |") {
		t.Errorf("entity table missing the light row:\n%s", got.entityTable)
	}

	byID := map[string]connector.SnapshotEntity{}
	for _, e := range got.entities {
		if e.Kind != "entity" {
			t.Errorf("entity %q kind = %q, want %q", e.ExternalID, e.Kind, "entity")
		}
		byID[e.ExternalID] = e
	}
	if len(byID) != 11 {
		t.Fatalf("entities = %d, want 11", len(byID))
	}

	light := byID["light.living_room"]
	if light.Name != "Living Room" {
		t.Errorf("light Name = %q, want the friendly name", light.Name)
	}
	if light.Attributes["supportedFeatures"] != float64(44) {
		t.Errorf("light supportedFeatures = %v, want 44", light.Attributes["supportedFeatures"])
	}
	if light.Attributes["state"] != "on" || light.Attributes["available"] != true || light.Attributes["domain"] != "light" {
		t.Errorf("light attributes = %+v", light.Attributes)
	}

	sensor := byID["sensor.outdoor_temperature"]
	if sensor.Attributes["state"] != "available" {
		t.Errorf("sensor state attribute = %v, want %q", sensor.Attributes["state"], "available")
	}
	if sensor.Attributes["entityCategory"] != "diagnostic" || sensor.Attributes["deviceClass"] != "temperature" || sensor.Attributes["unitOfMeasurement"] != "°C" {
		t.Errorf("sensor attributes = %+v", sensor.Attributes)
	}

	pump := byID["switch.pump"]
	if pump.Attributes["state"] != "unavailable" || pump.Attributes["available"] != false {
		t.Errorf("pump attributes = %+v, want unavailable reported as-is", pump.Attributes)
	}

	tracker := byID["device_tracker.phone"]
	if tracker.IP != "192.168.1.50" || tracker.Hostname != "phone-01" {
		t.Errorf("device_tracker IP/Hostname = %q/%q, want them filled for cross-connector matching", tracker.IP, tracker.Hostname)
	}
	if ip := byID["device_tracker.tablet"].IP; ip != "" {
		t.Errorf("unparseable ip attribute = %q, want it dropped", ip)
	}
}

// TestBuildEntitiesIsDeterministic guards the property the whole connector
// exists to preserve: re-fetching an idle instance yields an identical
// snapshot even though Home Assistant restamped every entity.
func TestBuildEntitiesIsDeterministic(t *testing.T) {
	restamped := strings.ReplaceAll(statesJSON, "2024-06-01T10:00:00+00:00", "2025-01-09T22:13:05+00:00")
	restamped = strings.ReplaceAll(restamped, `"state":"21.4"`, `"state":"23.9"`)

	first := buildEntities([]byte(statesJSON), 0)
	second := buildEntities([]byte(restamped), 0)
	if first.entityTable != second.entityTable {
		t.Errorf("entity table changed after a restamp:\n%s\n---\n%s", first.entityTable, second.entityTable)
	}
	if first.domains != second.domains {
		t.Errorf("domain summary changed after a restamp")
	}
	for i := range first.entities {
		if first.entities[i].Attributes["state"] != second.entities[i].Attributes["state"] {
			t.Errorf("entity %q state changed after a restamp", first.entities[i].ExternalID)
		}
	}
}

func TestBuildEntitiesRespectsMaxEntities(t *testing.T) {
	got := buildEntities([]byte(statesJSON), 3)
	if len(got.entities) != 3 {
		t.Fatalf("entities = %d, want 3", len(got.entities))
	}
	// Sorted before capping, so the same three survive every time.
	wantIDs := []string{"automation.night", "binary_sensor.front_door", "device_tracker.phone"}
	for i, want := range wantIDs {
		if got.entities[i].ExternalID != want {
			t.Errorf("entities[%d] = %q, want %q", i, got.entities[i].ExternalID, want)
		}
	}
	if !strings.Contains(got.entityTable, "Showing 3 of 11 entities") {
		t.Errorf("entity table missing the truncation note:\n%s", got.entityTable)
	}
	// The per-domain summary still counts every entity.
	if got.metadata["entity_count"] != "11" {
		t.Errorf("entity_count = %q, want 11 even when the table is capped", got.metadata["entity_count"])
	}
}

func TestBuildEntitiesEdgeCases(t *testing.T) {
	empty := buildEntities([]byte(`[]`), 0)
	if !strings.Contains(empty.entityTable, "No entities") || empty.metadata["entity_count"] != "0" {
		t.Errorf("empty states = %+v", empty)
	}
	bad := buildEntities([]byte(`{}`), 0)
	if !strings.Contains(bad.entityTable, "malformed response") || !strings.Contains(bad.domains, "malformed response") {
		t.Errorf("malformed states = %+v", bad)
	}
}

func TestStableState(t *testing.T) {
	tests := []struct {
		domain, state, want string
	}{
		{"light", "on", "on"},
		{"binary_sensor", "off", "off"},
		{"device_tracker", "home", "home"},
		{"switch", "unavailable", "unavailable"},
		{"switch", "unknown", "unknown"},
		{"switch", "", "unknown"},
		{"sensor", "21.4", "available"},
		{"sensor", "sunny", "available"},         // volatile domain, discrete value
		{"sensor", "unavailable", "unavailable"}, // availability still matters
		{"sun", "above_horizon", "available"},
		{"weather", "cloudy", "available"},
		{"zone", "1", "available"},                    // numeric state in a stable domain
		{"text", "2024-06-01T10:00:00Z", "available"}, // timestamp state
		{"climate", "heat", "heat"},
		{"update", "on", "on"},
	}
	for _, tt := range tests {
		if got := stableState(tt.domain, tt.state); got != tt.want {
			t.Errorf("stableState(%q, %q) = %q, want %q", tt.domain, tt.state, got, tt.want)
		}
	}
}

func TestBuildServices(t *testing.T) {
	content, metadata := buildServices([]byte(servicesJSON))

	if metadata["service_domain_count"] != "2" || metadata["service_count"] != "20" {
		t.Errorf("metadata = %+v", metadata)
	}
	if !strings.Contains(content, "| light | 3 | toggle, turn_off, turn_on |") {
		t.Errorf("content missing the light row:\n%s", content)
	}
	if !strings.Contains(content, "+2 more") {
		t.Errorf("content missing the truncation marker:\n%s", content)
	}
	// Domains are sorted, so homeassistant comes before light.
	if i, j := strings.Index(content, "| homeassistant |"), strings.Index(content, "| light |"); i < 0 || j < 0 || i > j {
		t.Errorf("service domains are not sorted:\n%s", content)
	}
}

func TestBuildServicesEdgeCases(t *testing.T) {
	content, metadata := buildServices([]byte(`[]`))
	if !strings.Contains(content, "No services") || metadata["service_count"] != "0" {
		t.Errorf("empty services = %q / %+v", content, metadata)
	}
	if content, metadata := buildServices([]byte("nope")); !strings.Contains(content, "malformed response") || metadata != nil {
		t.Errorf("malformed services = %q / %+v", content, metadata)
	}
}

func TestCellEscapesMarkdown(t *testing.T) {
	if got := cell(""); got != "—" {
		t.Errorf("cell(\"\") = %q, want an em dash", got)
	}
	if got := cell("a|b\nc"); got != `a\|b c` {
		t.Errorf("cell() = %q, want the pipe escaped and the newline flattened", got)
	}
}

func TestEntityDomain(t *testing.T) {
	tests := map[string]string{"light.kitchen": "light", "sun.sun": "sun", "bare": "bare", ".leading": ".leading"}
	for in, want := range tests {
		if got := entityDomain(in); got != want {
			t.Errorf("entityDomain(%q) = %q, want %q", in, got, want)
		}
	}
}
