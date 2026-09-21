package pihole

import (
	"reflect"
	"strings"
	"testing"
)

func TestResolveGroups(t *testing.T) {
	names := map[int]string{0: "Default", 2: "Kids"}
	got := resolveGroups([]int{2, 0, 7}, names)
	want := []string{"7", "Default", "Kids"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("resolveGroups() = %v, want %v", got, want)
	}
	if got := resolveGroups(nil, names); len(got) != 0 {
		t.Errorf("resolveGroups(nil) = %v, want empty", got)
	}
}

func TestCellEscapesTableBreakers(t *testing.T) {
	if got := cell("a|b\nc"); got != "a\\|b c" {
		t.Errorf("cell() = %q, want %q", got, "a\\|b c")
	}
	if got := cell(""); got != "-" {
		t.Errorf("cell(\"\") = %q, want %q", got, "-")
	}
}

func TestParseGroupsBothVersions(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) ([]groupRow, error)
		data  string
	}{
		{"v6", parseGroupsV6, `{"groups":[{"id":0,"name":"Default","comment":"the default group","enabled":true}]}`},
		{"v5", parseGroupsV5, `{"data":[{"id":0,"name":"Default","description":"the default group","enabled":1}]}`},
	}
	want := groupRow{ID: 0, Name: "Default", Comment: "the default group", Enabled: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := tt.parse([]byte(tt.data))
			if err != nil {
				t.Fatalf("parse error = %v", err)
			}
			if len(rows) != 1 || rows[0] != want {
				t.Fatalf("rows = %+v, want [%+v]", rows, want)
			}
			if _, err := tt.parse([]byte("not json")); err == nil {
				t.Error("parse(invalid) error = nil, want malformed response error")
			}
		})
	}
}

func TestParseAdlistsBothVersions(t *testing.T) {
	names := map[int]string{0: "Default"}
	v6, err := parseAdlistsV6([]byte(`{"lists":[{"address":"https://example.com/hosts","type":"block","comment":"main","enabled":true,"groups":[0]}]}`), names)
	if err != nil {
		t.Fatalf("parseAdlistsV6 error = %v", err)
	}
	v5, err := parseAdlistsV5([]byte(`{"data":[{"address":"https://example.com/hosts","comment":"main","enabled":"1","groups":[0]}]}`), names)
	if err != nil {
		t.Fatalf("parseAdlistsV5 error = %v", err)
	}
	want := adlistRow{Address: "https://example.com/hosts", Kind: "block", Comment: "main", Enabled: true, Groups: []string{"Default"}}
	for name, rows := range map[string][]adlistRow{"v6": v6, "v5": v5} {
		if len(rows) != 1 || !reflect.DeepEqual(rows[0], want) {
			t.Errorf("%s rows = %+v, want [%+v]", name, rows, want)
		}
	}
}

func TestParseAdlistsV6DefaultsMissingTypeToBlock(t *testing.T) {
	rows, err := parseAdlistsV6([]byte(`{"lists":[{"address":"https://example.com/hosts","enabled":true}]}`), nil)
	if err != nil {
		t.Fatalf("parseAdlistsV6 error = %v", err)
	}
	if len(rows) != 1 || rows[0].Kind != "block" {
		t.Fatalf("rows = %+v, want one block list", rows)
	}
}

func TestParseClientsBothVersions(t *testing.T) {
	names := map[int]string{1: "Kids"}
	v6, err := parseClientsV6([]byte(`{"clients":[{"client":"10.0.0.20","comment":"tablet","groups":[1]}]}`), names)
	if err != nil {
		t.Fatalf("parseClientsV6 error = %v", err)
	}
	v5, err := parseClientsV5([]byte(`{"data":[{"ip":"10.0.0.20","comment":"tablet","groups":[1]}]}`), names)
	if err != nil {
		t.Fatalf("parseClientsV5 error = %v", err)
	}
	want := clientRow{Client: "10.0.0.20", Comment: "tablet", Groups: []string{"Kids"}}
	for name, rows := range map[string][]clientRow{"v6": v6, "v5": v5} {
		if len(rows) != 1 || !reflect.DeepEqual(rows[0], want) {
			t.Errorf("%s rows = %+v, want [%+v]", name, rows, want)
		}
	}
}

func TestParseDomainsBothVersions(t *testing.T) {
	v6, err := parseDomainsV6([]byte(`{"domains":[{"domain":"ads.example.com","type":"deny","kind":"exact","comment":"","enabled":true,"groups":[]}]}`), nil)
	if err != nil {
		t.Fatalf("parseDomainsV6 error = %v", err)
	}
	v5, err := parseDomainsV5([]byte(`{"data":[{"domain":"ads.example.com","type":1,"comment":"","enabled":1,"groups":[]}]}`), nil)
	if err != nil {
		t.Fatalf("parseDomainsV5 error = %v", err)
	}
	want := domainRow{Domain: "ads.example.com", Rule: "deny", Match: "exact", Enabled: true, Groups: []string{}}
	for name, rows := range map[string][]domainRow{"v6": v6, "v5": v5} {
		if len(rows) != 1 || !reflect.DeepEqual(rows[0], want) {
			t.Errorf("%s rows = %+v, want [%+v]", name, rows, want)
		}
	}
}

func TestV5DomainType(t *testing.T) {
	tests := []struct {
		in    int
		rule  string
		match string
	}{
		{0, "allow", "exact"},
		{1, "deny", "exact"},
		{2, "allow", "regex"},
		{3, "deny", "regex"},
		{9, "unknown", "unknown"},
	}
	for _, tt := range tests {
		rule, match := v5DomainType(tt.in)
		if rule != tt.rule || match != tt.match {
			t.Errorf("v5DomainType(%d) = %s/%s, want %s/%s", tt.in, rule, match, tt.rule, tt.match)
		}
	}
}

func TestV5BoolUnmarshal(t *testing.T) {
	tests := []struct {
		in      string
		want    bool
		wantErr bool
	}{
		{`1`, true, false},
		{`"1"`, true, false},
		{`true`, true, false},
		{`0`, false, false},
		{`"0"`, false, false},
		{`false`, false, false},
		{`null`, false, false},
		{`"maybe"`, false, true},
	}
	for _, tt := range tests {
		var b v5Bool
		err := b.UnmarshalJSON([]byte(tt.in))
		if tt.wantErr {
			if err == nil {
				t.Errorf("UnmarshalJSON(%s) error = nil, want error", tt.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("UnmarshalJSON(%s) error = %v", tt.in, err)
		}
		if bool(b) != tt.want {
			t.Errorf("UnmarshalJSON(%s) = %v, want %v", tt.in, bool(b), tt.want)
		}
	}
}

func TestBuildTablesEmptyPlaceholders(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{"groups", firstOf(buildGroupTable(nil)), "_No groups returned_"},
		{"adlists", firstOf(buildAdlistTable(nil)), "_No blocklists returned_"},
		{"clients", firstOf(buildClientTable(nil)), "_No clients returned_"},
		{"domains", firstOf(buildDomainTable(nil)), "_No domain rules returned_"},
	}
	for _, tt := range tests {
		if tt.content != tt.want {
			t.Errorf("%s placeholder = %q, want %q", tt.name, tt.content, tt.want)
		}
	}
}

func firstOf(content string, _ any) string { return content }

func TestBuildTablesEntities(t *testing.T) {
	content, entities := buildAdlistTable([]adlistRow{{Address: "https://x/h", Kind: "block", Enabled: true, Groups: []string{"Default"}}})
	if !strings.Contains(content, "https://x/h") {
		t.Errorf("content = %q, want the adlist address", content)
	}
	if len(entities) != 1 || entities[0].Kind != "blocklist" || entities[0].ExternalID != "https://x/h" {
		t.Fatalf("entities = %+v", entities)
	}
	if entities[0].Attributes["list_type"] != "block" || entities[0].Attributes["enabled"] != true {
		t.Errorf("attributes = %+v", entities[0].Attributes)
	}

	_, clients := buildClientTable([]clientRow{{Client: "10.0.0.20"}, {Client: "aa:bb:cc:dd:ee:ff"}})
	if clients[0].IP != "10.0.0.20" {
		t.Errorf("client IP = %q, want the bare IP", clients[0].IP)
	}
	if clients[1].IP != "" {
		t.Errorf("MAC client IP = %q, want empty", clients[1].IP)
	}

	_, domains := buildDomainTable([]domainRow{{Domain: "ads.example.com", Rule: "deny", Match: "regex"}})
	if domains[0].ExternalID != "deny/regex/ads.example.com" {
		t.Errorf("domain ExternalID = %q", domains[0].ExternalID)
	}

	_, groups := buildGroupTable([]groupRow{{Name: "Default", Enabled: true}})
	if groups[0].Kind != "dns_group" || groups[0].ExternalID != "Default" {
		t.Errorf("group entity = %+v", groups[0])
	}
}

func TestBuildHostsTableV5(t *testing.T) {
	content, entities := buildHostsTableV5([]byte(`{"data":[["10.0.0.5","nas.internal.example.com"],["bad"]]}`))
	if !strings.Contains(content, "nas.internal.example.com") {
		t.Errorf("content = %q", content)
	}
	if len(entities) != 1 || entities[0].IP != "10.0.0.5" || entities[0].Kind != "dns_record" {
		t.Fatalf("entities = %+v", entities)
	}
	if entities[0].Attributes["source"] != "local_dns" {
		t.Errorf("attributes = %+v", entities[0].Attributes)
	}

	if content, _ := buildHostsTableV5([]byte("not json")); !strings.Contains(content, "malformed response") {
		t.Errorf("content = %q, want a malformed response placeholder", content)
	}
	if content, _ := buildHostsTableV5([]byte(`{"data":[]}`)); !strings.Contains(content, "_No local DNS records returned_") {
		t.Errorf("content = %q, want the empty placeholder", content)
	}
}
