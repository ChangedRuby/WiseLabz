package pihole

import (
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildHostsTableMalformedCases(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{
			name: "empty JSON object returns placeholder",
			data: []byte(`{}`),
			want: "_No local DNS records returned_",
		},
		{
			name: "JSON with empty hosts returns placeholder",
			data: []byte(`{"config":{"dns":{"hosts":[]}}}`),
			want: "_No local DNS records returned_",
		},
		{
			name: "invalid JSON returns malformed placeholder",
			data: []byte(`not json`),
			want: "malformed response",
		},
		{
			name: "partial JSON returns malformed placeholder",
			data: []byte(`{"config":{"dns":{"hosts":[`),
			want: "malformed response",
		},
		{
			name: "entries without a hostname produce no rows",
			data: []byte(`{"config":{"dns":{"hosts":["10.0.0.5"]}}}`),
			want: "_No local DNS records returned_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, entities := buildHostsTable(tt.data)
			if !strings.Contains(content, tt.want) {
				t.Errorf("buildHostsTable() content = %q, want to contain %q", content, tt.want)
			}
			if entities != nil {
				t.Errorf("buildHostsTable() entities = %+v, want nil", entities)
			}
		})
	}
}

func TestBuildHostsTableValidRecords(t *testing.T) {
	data := []byte(`{"config":{"dns":{"hosts":["10.0.0.5 nas.internal.example.com","10.0.0.6 printer.internal.example.com"]}}}`)
	content, entities := buildHostsTable(data)

	if !strings.Contains(content, "nas.internal.example.com") || !strings.Contains(content, "10.0.0.5") {
		t.Errorf("buildHostsTable() content missing expected record: %q", content)
	}

	if len(entities) != 2 {
		t.Fatalf("buildHostsTable() entities len = %d, want 2", len(entities))
	}

	want := []connector.SnapshotEntity{
		{Kind: "dns_record", Hostname: "nas.internal.example.com", IP: "10.0.0.5"},
		{Kind: "dns_record", Hostname: "printer.internal.example.com", IP: "10.0.0.6"},
	}
	for i, w := range want {
		if entities[i] != w {
			t.Errorf("entities[%d] = %+v, want %+v", i, entities[i], w)
		}
	}
}
