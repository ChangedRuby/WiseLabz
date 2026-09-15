package doc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedEngineConnectorWithEntities(t *testing.T, s *store.Store, name, category, connectorType string, entities []connector.SnapshotEntity) string {
	t.Helper()
	ctx := context.Background()
	record := &store.ConnectorRecord{Name: name, Category: category, Type: connectorType, URL: "https://example.test"}
	if err := s.CreateConnector(ctx, record); err != nil {
		t.Fatalf("create connector: %v", err)
	}
	data, err := json.Marshal(connector.ServiceSnapshot{
		ServiceName: name,
		Type:        connectorType,
		Entities:    entities,
		FetchedAt:   time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := s.CreateSnapshot(ctx, &store.SnapshotRecord{ConnectorID: record.ID, Data: string(data)}); err != nil {
		t.Fatalf("create snapshot: %v", err)
	}
	return record.ID
}

func TestMatchEntitiesExternalIDPrecedence(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	seedEngineConnectorWithEntities(t, s, "Proxmox", "virtualization", "proxmox", []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", ExternalID: "100", IP: "10.0.0.5"},
	})
	connectorID := seedEngineConnectorWithEntities(t, s, "Docker", "containers_paas", "docker", nil)

	links, err := matchEntities(ctx, s, connectorID, []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", ExternalID: "100", IP: "10.0.0.99"},
	})
	if err != nil {
		t.Fatalf("matchEntities() error: %v", err)
	}
	if len(links) != 1 || links[0].Reason != "external ID" {
		t.Fatalf("matchEntities() = %+v, want single external ID match", links)
	}
}

func TestMatchEntitiesIPPrecedence(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	seedEngineConnectorWithEntities(t, s, "Proxmox", "virtualization", "proxmox", []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", IP: "10.0.0.5"},
	})
	connectorID := seedEngineConnectorWithEntities(t, s, "pfSense", "networking", "pfsense", nil)

	links, err := matchEntities(ctx, s, connectorID, []connector.SnapshotEntity{
		{Kind: "rule", Name: "allow-web", IP: "10.0.0.5"},
	})
	if err != nil {
		t.Fatalf("matchEntities() error: %v", err)
	}
	if len(links) != 1 || links[0].Reason != "IP address" {
		t.Fatalf("matchEntities() = %+v, want single IP match", links)
	}
}

func TestMatchEntitiesHostnamePrecedenceCaseInsensitive(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	seedEngineConnectorWithEntities(t, s, "Pi-hole", "dns", "pihole", []connector.SnapshotEntity{
		{Kind: "dns_record", Name: "web-01", Hostname: "Web-01.lab.local"},
	})
	connectorID := seedEngineConnectorWithEntities(t, s, "Proxmox", "virtualization", "proxmox", nil)

	links, err := matchEntities(ctx, s, connectorID, []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", Hostname: "web-01.lab.local"},
	})
	if err != nil {
		t.Fatalf("matchEntities() error: %v", err)
	}
	if len(links) != 1 || links[0].Reason != "hostname" {
		t.Fatalf("matchEntities() = %+v, want single hostname match", links)
	}
}

func TestMatchEntitiesFansOutAcrossConnectors(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	seedEngineConnectorWithEntities(t, s, "Docker", "containers_paas", "docker", []connector.SnapshotEntity{
		{Kind: "container", Name: "web-01", IP: "10.0.0.5"},
	})
	seedEngineConnectorWithEntities(t, s, "Pi-hole", "dns", "pihole", []connector.SnapshotEntity{
		{Kind: "dns_record", Name: "web-01", Hostname: "web-01.lab.local"},
	})
	connectorID := seedEngineConnectorWithEntities(t, s, "Proxmox", "virtualization", "proxmox", nil)

	links, err := matchEntities(ctx, s, connectorID, []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", IP: "10.0.0.5", Hostname: "web-01.lab.local"},
	})
	if err != nil {
		t.Fatalf("matchEntities() error: %v", err)
	}
	if len(links) != 2 {
		t.Fatalf("matchEntities() returned %d links, want 2 (one per connector)", len(links))
	}
}

func TestMatchEntitiesDedupesExactExternalIDDuplicates(t *testing.T) {
	ctx := context.Background()
	s := newEngineTestStore(t)
	seedEngineConnectorWithEntities(t, s, "Proxmox", "virtualization", "proxmox", []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", ExternalID: "100", IP: "10.0.0.5"},
	})
	connectorID := seedEngineConnectorWithEntities(t, s, "Docker", "containers_paas", "docker", nil)

	// Two "mine" entities both matching the same external entity by
	// ExternalID+Kind should collapse to a single link.
	links, err := matchEntities(ctx, s, connectorID, []connector.SnapshotEntity{
		{Kind: "vm", Name: "web-01", ExternalID: "100"},
		{Kind: "vm", Name: "web-01-alias", ExternalID: "100"},
	})
	if err != nil {
		t.Fatalf("matchEntities() error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("matchEntities() returned %d links, want 1 deduped link", len(links))
	}
}
