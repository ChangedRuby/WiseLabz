package doc

import (
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func TestRenderMermaid(t *testing.T) {
	links := []EntityLink{
		{
			Entity:        connector.SnapshotEntity{Kind: "container", Name: "web-01", ExternalID: "abc123"},
			ConnectorID:   "docker-1",
			ConnectorName: "Docker",
			Reason:        "IP address",
		},
		{
			Entity:        connector.SnapshotEntity{Kind: "dns_record", Name: "web-01", Hostname: "web-01.lab.local"},
			ConnectorID:   "pihole-1",
			ConnectorName: "Pi-hole",
			Reason:        "hostname",
		},
	}

	got := renderMermaid("web-01", links)
	want := "graph LR\n" +
		"    n" + shortHash("center|web-01") + "[\"web-01\"]\n" +
		"    n" + shortHash("docker-1|container|abc123") + "[\"web-01 (container)\"]\n" +
		"    n" + shortHash("center|web-01") + " -->|IP address| n" + shortHash("docker-1|container|abc123") + "\n" +
		"    n" + shortHash("pihole-1|dns_record|web-01") + "[\"web-01 (dns_record)\"]\n" +
		"    n" + shortHash("center|web-01") + " -->|hostname| n" + shortHash("pihole-1|dns_record|web-01") + "\n"

	if got != want {
		t.Fatalf("renderMermaid() =\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderMermaidNoLinks(t *testing.T) {
	got := renderMermaid("web-01", nil)
	want := "graph LR\n    n" + shortHash("center|web-01") + "[\"web-01\"]\n"
	if got != want {
		t.Fatalf("renderMermaid() = %q, want %q", got, want)
	}
}
