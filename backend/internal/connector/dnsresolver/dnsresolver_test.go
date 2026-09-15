package dnsresolver

import (
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestBuildHostOverrideTableMalformedCases(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{
			name: "empty JSON object returns placeholder",
			data: []byte(`{}`),
			want: "_No host overrides returned_",
		},
		{
			name: "JSON with empty data returns placeholder",
			data: []byte(`{"data":[]}`),
			want: "_No host overrides returned_",
		},
		{
			name: "invalid JSON returns malformed placeholder",
			data: []byte(`not json`),
			want: "malformed response",
		},
		{
			name: "partial JSON returns malformed placeholder",
			data: []byte(`{"data":[{"host":"test"`),
			want: "malformed response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, entities := buildHostOverrideTable(tt.data)
			if !strings.Contains(content, tt.want) {
				t.Errorf("buildHostOverrideTable() content = %q, want to contain %q", content, tt.want)
			}
			if entities != nil {
				t.Errorf("buildHostOverrideTable() entities = %+v, want nil", entities)
			}
		})
	}
}

func TestBuildHostOverrideTableValidOverrides(t *testing.T) {
	data := []byte(`{
		"data":[
			{"host":"nas","domain":"internal.example.com","ip":"10.0.0.5","descr":"NAS"},
			{"host":"","domain":"example.com","ip":"10.0.0.1","descr":"Root domain"}
		]
	}`)
	content, entities := buildHostOverrideTable(data)

	if !strings.Contains(content, "nas") || !strings.Contains(content, "10.0.0.5") {
		t.Errorf("buildHostOverrideTable() content missing expected override: %q", content)
	}

	if len(entities) != 2 {
		t.Fatalf("buildHostOverrideTable() entities len = %d, want 2", len(entities))
	}

	want := []connector.SnapshotEntity{
		{Kind: "dns_record", Hostname: "nas.internal.example.com", IP: "10.0.0.5"},
		{Kind: "dns_record", Hostname: "example.com", IP: "10.0.0.1"},
	}
	for i, w := range want {
		if entities[i] != w {
			t.Errorf("entities[%d] = %+v, want %+v", i, entities[i], w)
		}
	}
}
