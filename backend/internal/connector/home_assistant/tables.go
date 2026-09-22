package homeassistant

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

// volatileDomains are entity domains whose state is a continuously changing
// measurement or timestamp rather than a discrete mode. Their values are
// collapsed to availability so an idle instance keeps producing identical
// snapshots.
var volatileDomains = map[string]bool{
	"counter": true, "date": true, "datetime": true, "image": true,
	"input_datetime": true, "input_number": true, "number": true,
	"sensor": true, "sun": true, "time": true, "weather": true,
}

// maxServiceNamesPerRow caps how many service names a single domain row
// lists before it is summarised; the homeassistant and script domains alone
// can carry dozens.
const maxServiceNamesPerRow = 15

// cell escapes a value for use inside a Markdown table cell. Friendly names
// and service descriptions are free-form user text and may contain '|'.
func cell(s string) string {
	if s == "" {
		return "—"
	}
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", `\|`)
}

func malformed(title string, err error) string {
	return "_" + title + " unavailable: " + connector.NewMalformedResponseError(err).Error() + "_"
}

// configPayload is the subset of /api/config this connector reports. The
// latitude/longitude the endpoint also returns are deliberately not
// recorded: they are precise home coordinates and add nothing to a
// configuration snapshot.
type configPayload struct {
	Version      string   `json:"version"`
	LocationName string   `json:"location_name"`
	TimeZone     string   `json:"time_zone"`
	Country      string   `json:"country"`
	Language     string   `json:"language"`
	Currency     string   `json:"currency"`
	Elevation    int      `json:"elevation"`
	ConfigDir    string   `json:"config_dir"`
	ConfigSource string   `json:"config_source"`
	State        string   `json:"state"`
	SafeMode     bool     `json:"safe_mode"`
	RecoveryMode bool     `json:"recovery_mode"`
	InternalURL  string   `json:"internal_url"`
	ExternalURL  string   `json:"external_url"`
	Components   []string `json:"components"`
	UnitSystem   struct {
		Length      string `json:"length"`
		Mass        string `json:"mass"`
		Temperature string `json:"temperature"`
		Volume      string `json:"volume"`
		Pressure    string `json:"pressure"`
		WindSpeed   string `json:"wind_speed"`
	} `json:"unit_system"`
}

// buildOverview renders the instance configuration from /api/config and
// returns the snapshot metadata it carries alongside.
func buildOverview(raw []byte) (string, map[string]string) {
	var cfg configPayload
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return malformed("Overview", err), nil
	}

	rows := []struct{ key, value string }{
		{"Version", cfg.Version},
		{"Location", cfg.LocationName},
		{"Time zone", cfg.TimeZone},
		{"Country", cfg.Country},
		{"Language", cfg.Language},
		{"Currency", cfg.Currency},
		{"Elevation", fmt.Sprintf("%d m", cfg.Elevation)},
		{"Unit system", unitSystemSummary(cfg)},
		{"State", cfg.State},
		{"Config source", cfg.ConfigSource},
		{"Config directory", cfg.ConfigDir},
		{"Internal URL", cfg.InternalURL},
		{"External URL", cfg.ExternalURL},
		{"Safe mode", strconv.FormatBool(cfg.SafeMode)},
		{"Recovery mode", strconv.FormatBool(cfg.RecoveryMode)},
		{"Integrations loaded", strconv.Itoa(len(cfg.Components))},
	}

	var b strings.Builder
	b.WriteString("| Setting | Value |\n")
	b.WriteString("|---------|-------|\n")
	for _, r := range rows {
		_, _ = fmt.Fprintf(&b, "| %s | %s |\n", r.key, cell(r.value))
	}
	b.WriteString("\n_Areas, devices and config entries are not exposed by the Home Assistant REST API (they are WebSocket-only), so this snapshot describes integrations and entities instead._\n")

	metadata := map[string]string{
		"component_count": strconv.Itoa(len(cfg.Components)),
	}
	putMeta(metadata, "home_assistant_version", cfg.Version)
	putMeta(metadata, "location_name", cfg.LocationName)
	putMeta(metadata, "time_zone", cfg.TimeZone)
	putMeta(metadata, "config_source", cfg.ConfigSource)
	putMeta(metadata, "home_assistant_state", cfg.State)
	return b.String(), metadata
}

func unitSystemSummary(cfg configPayload) string {
	parts := make([]string, 0, 3)
	for _, v := range []string{cfg.UnitSystem.Temperature, cfg.UnitSystem.Length, cfg.UnitSystem.Mass} {
		if v != "" {
			parts = append(parts, v)
		}
	}
	return strings.Join(parts, " · ")
}

// buildIntegrations renders the "components" list from /api/config. Entries
// are either a bare integration ("mqtt") or "<platform>.<integration>"
// ("sensor.mqtt"), so they are folded into one row per integration.
func buildIntegrations(raw []byte) (string, []connector.SnapshotEntity) {
	var cfg configPayload
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return malformed("Integrations", err), nil
	}
	if len(cfg.Components) == 0 {
		return "_No integrations reported_", nil
	}

	platforms := map[string][]string{}
	for _, component := range cfg.Components {
		name, platform := component, ""
		if dot := strings.IndexByte(component, '.'); dot > 0 {
			platform, name = component[:dot], component[dot+1:]
		}
		if _, ok := platforms[name]; !ok {
			platforms[name] = nil
		}
		if platform != "" {
			platforms[name] = append(platforms[name], platform)
		}
	}

	names := make([]string, 0, len(platforms))
	for name := range platforms {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	b.WriteString("| Integration | Platforms |\n")
	b.WriteString("|-------------|-----------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(names))
	for _, name := range names {
		loaded := platforms[name]
		sort.Strings(loaded)
		_, _ = fmt.Fprintf(&b, "| %s | %s |\n", cell(name), cell(strings.Join(loaded, ", ")))

		attrs := map[string]any{}
		putStrings(attrs, "platforms", loaded)
		entities = append(entities, connector.SnapshotEntity{
			Kind: "integration", Name: name, ExternalID: name, Attributes: attrs,
		})
	}
	return b.String(), entities
}

// entityResult carries everything buildEntities derives from /api/states.
type entityResult struct {
	domains     string
	entityTable string
	entities    []connector.SnapshotEntity
	metadata    map[string]string
}

type statePayload struct {
	EntityID   string         `json:"entity_id"`
	State      string         `json:"state"`
	Attributes map[string]any `json:"attributes"`
}

// buildEntities renders /api/states as a per-domain summary plus an entity
// table, and emits one "entity" snapshot entity per row. maxEntities caps
// both the table and the emitted entities (0 means no cap); entities are
// sorted by entity ID first, so the cap always keeps the same set.
func buildEntities(raw []byte, maxEntities int) entityResult {
	var states []statePayload
	if err := json.Unmarshal(raw, &states); err != nil {
		msg := malformed("Entities", err)
		return entityResult{domains: msg, entityTable: msg}
	}
	if len(states) == 0 {
		return entityResult{
			domains:     "_No entities returned_",
			entityTable: "_No entities returned_",
			metadata:    map[string]string{"entity_count": "0"},
		}
	}

	sort.Slice(states, func(i, j int) bool { return states[i].EntityID < states[j].EntityID })

	type domainCount struct{ total, unavailable int }
	counts := map[string]*domainCount{}
	unavailable := 0
	for _, s := range states {
		domain := entityDomain(s.EntityID)
		c, ok := counts[domain]
		if !ok {
			c = &domainCount{}
			counts[domain] = c
		}
		c.total++
		if !isAvailable(s.State) {
			c.unavailable++
			unavailable++
		}
	}

	domains := make([]string, 0, len(counts))
	for d := range counts {
		domains = append(domains, d)
	}
	sort.Strings(domains)

	var summary strings.Builder
	summary.WriteString("| Domain | Entities | Unavailable |\n")
	summary.WriteString("|--------|----------|-------------|\n")
	for _, d := range domains {
		_, _ = fmt.Fprintf(&summary, "| %s | %d | %d |\n", cell(d), counts[d].total, counts[d].unavailable)
	}

	shown := states
	if maxEntities > 0 && len(shown) > maxEntities {
		shown = shown[:maxEntities]
	}

	var table strings.Builder
	table.WriteString("| Entity ID | Name | Domain | State | Device Class | Unit |\n")
	table.WriteString("|-----------|------|--------|-------|--------------|------|\n")
	entities := make([]connector.SnapshotEntity, 0, len(shown))
	for _, s := range shown {
		domain := entityDomain(s.EntityID)
		state := stableState(domain, s.State)
		name := attrString(s.Attributes, "friendly_name")
		deviceClass := attrString(s.Attributes, "device_class")
		unit := attrString(s.Attributes, "unit_of_measurement")

		_, _ = fmt.Fprintf(&table, "| %s | %s | %s | %s | %s | %s |\n",
			cell(s.EntityID), cell(name), cell(domain), cell(state), cell(deviceClass), cell(unit))

		attrs := map[string]any{"domain": domain, "state": state, "available": isAvailable(s.State)}
		putString(attrs, "deviceClass", deviceClass)
		putString(attrs, "unitOfMeasurement", unit)
		putString(attrs, "entityCategory", attrString(s.Attributes, "entity_category"))
		if features, ok := attrNumber(s.Attributes, "supported_features"); ok {
			attrs["supportedFeatures"] = features
		}
		entities = append(entities, connector.SnapshotEntity{
			Kind:       "entity",
			Name:       nameOrID(name, s.EntityID),
			IP:         attrIP(s.Attributes, "ip"),
			Hostname:   attrString(s.Attributes, "host_name"),
			ExternalID: s.EntityID,
			Attributes: attrs,
		})
	}
	if len(shown) < len(states) {
		_, _ = fmt.Fprintf(&table, "\n_Showing %d of %d entities; raise max_entities to include more._\n", len(shown), len(states))
	}

	return entityResult{
		domains:     summary.String(),
		entityTable: table.String(),
		entities:    entities,
		metadata: map[string]string{
			"entity_count":             strconv.Itoa(len(states)),
			"entity_domain_count":      strconv.Itoa(len(domains)),
			"unavailable_entity_count": strconv.Itoa(unavailable),
		},
	}
}

// entityDomain returns the part of an entity ID before the dot.
func entityDomain(entityID string) string {
	if dot := strings.IndexByte(entityID, '.'); dot > 0 {
		return entityID[:dot]
	}
	return entityID
}

func nameOrID(name, entityID string) string {
	if name != "" {
		return name
	}
	return entityID
}

// isAvailable reports whether Home Assistant currently holds a value for an
// entity.
func isAvailable(state string) bool {
	return state != "" && state != "unavailable" && state != "unknown"
}

// stableState reduces an entity state to a value that only changes when the
// instance's configuration or availability changes. Measurements (any
// numeric state, any state in a volatile domain) and timestamps are reported
// as "available" so that re-syncing an idle instance yields an identical
// snapshot instead of a diff on every sensor reading.
func stableState(domain, state string) string {
	if !isAvailable(state) {
		if state == "" {
			return "unknown"
		}
		return state
	}
	if volatileDomains[domain] || isNumeric(state) || isTimestamp(state) {
		return "available"
	}
	return state
}

func isNumeric(state string) bool {
	_, err := strconv.ParseFloat(state, 64)
	return err == nil
}

func isTimestamp(state string) bool {
	_, err := time.Parse(time.RFC3339, state)
	return err == nil
}

// buildServices renders /api/services: one row per service domain with the
// service names it exposes.
func buildServices(raw []byte) (string, map[string]string) {
	var domains []struct {
		Domain   string                     `json:"domain"`
		Services map[string]json.RawMessage `json:"services"`
	}
	if err := json.Unmarshal(raw, &domains); err != nil {
		return malformed("Services", err), nil
	}
	if len(domains) == 0 {
		return "_No services returned_", map[string]string{"service_domain_count": "0", "service_count": "0"}
	}

	sort.Slice(domains, func(i, j int) bool { return domains[i].Domain < domains[j].Domain })

	var b strings.Builder
	b.WriteString("| Domain | Services | Names |\n")
	b.WriteString("|--------|----------|-------|\n")
	total := 0
	for _, d := range domains {
		names := make([]string, 0, len(d.Services))
		for name := range d.Services {
			names = append(names, name)
		}
		sort.Strings(names)
		total += len(names)
		_, _ = fmt.Fprintf(&b, "| %s | %d | %s |\n", cell(d.Domain), len(names), cell(truncateNames(names)))
	}

	return b.String(), map[string]string{
		"service_domain_count": strconv.Itoa(len(domains)),
		"service_count":        strconv.Itoa(total),
	}
}

func truncateNames(names []string) string {
	if len(names) <= maxServiceNamesPerRow {
		return strings.Join(names, ", ")
	}
	return fmt.Sprintf("%s, +%d more", strings.Join(names[:maxServiceNamesPerRow], ", "), len(names)-maxServiceNamesPerRow)
}

func attrString(attrs map[string]any, key string) string {
	if v, ok := attrs[key].(string); ok {
		return v
	}
	return ""
}

func attrNumber(attrs map[string]any, key string) (float64, bool) {
	v, ok := attrs[key].(float64)
	return v, ok
}

// attrIP returns the attribute only when it parses as an IP address, so a
// bogus value never ends up in cross-connector entity matching.
func attrIP(attrs map[string]any, key string) string {
	v := attrString(attrs, key)
	if v == "" || net.ParseIP(v) == nil {
		return ""
	}
	return v
}

func putMeta(metadata map[string]string, key, value string) {
	if value != "" {
		metadata[key] = value
	}
}

func putString(attrs map[string]any, key, value string) {
	if value != "" {
		attrs[key] = value
	}
}

func putStrings(attrs map[string]any, key string, values []string) {
	if len(values) > 0 {
		attrs[key] = values
	}
}
