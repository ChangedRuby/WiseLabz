package dnsresolver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestDNSResolverConfigPush(t *testing.T) {
	listBody := `{"data":[{"id":1,"host":"web","domain":"example.com","ip":"10.0.0.5","descr":""}]}`
	tests := []struct {
		name      string
		entityRef string
		fieldKey  string
		value     any
		wantErr   bool
	}{
		{name: "success", entityRef: "web.example.com", fieldKey: "ip", value: "10.0.0.9"},
		{name: "empty entityRef errors", entityRef: "", fieldKey: "ip", value: "10.0.0.9", wantErr: true},
		{name: "unsupported field errors", entityRef: "web.example.com", fieldKey: "ttl", value: 300, wantErr: true},
		{name: "not found errors", entityRef: "missing.example.com", fieldKey: "ip", value: "10.0.0.9", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath, gotMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					_, _ = w.Write([]byte(listBody))
				case http.MethodPatch:
					gotPath, gotMethod = r.URL.Path, r.Method
					_, _ = w.Write([]byte(`{"data":null}`))
				default:
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
			}))
			defer server.Close()

			c := &Connector{url: server.URL, apiKey: "key", client: server.Client()}
			err := c.ConfigPush(context.Background(), nil, tt.entityRef, tt.fieldKey, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ConfigPush() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ConfigPush() error = %v", err)
			}
			if gotMethod != "PATCH" || gotPath != "/api/v2/services/dns_resolver/host_override" {
				t.Errorf("request = %s %s, want PATCH /api/v2/services/dns_resolver/host_override", gotMethod, gotPath)
			}
		})
	}
}

func TestDNSResolverWritableFields(t *testing.T) {
	c := &Connector{}
	fields := c.WritableFields()
	if len(fields) != 1 || fields[0].Key != "ip" {
		t.Errorf("WritableFields() = %+v, want one field \"ip\"", fields)
	}
}

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
		{Kind: "dns_record", Hostname: "nas.internal.example.com", IP: "10.0.0.5", Attributes: map[string]any{"description": "NAS", "is_ipv6": false}},
		{Kind: "dns_record", Hostname: "example.com", IP: "10.0.0.1", Attributes: map[string]any{"description": "Root domain", "is_ipv6": false}},
	}
	for i, w := range want {
		if !reflect.DeepEqual(entities[i], w) {
			t.Errorf("entities[%d] = %+v, want %+v", i, entities[i], w)
		}
	}
}
