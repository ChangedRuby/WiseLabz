package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/WiseLabz/wiselabz/internal/connector"
)

func testSnapshotWithAttributes() *connector.ServiceSnapshot {
	return &connector.ServiceSnapshot{
		ServiceName: "pve1",
		Type:        "proxmox",
		Entities: []connector.SnapshotEntity{
			{
				Kind:       "vm",
				Name:       "web-01",
				ExternalID: "100",
				Attributes: map[string]any{
					"status":           "running",
					"firewall_enabled": true,
					"protection":       false,
					"onboot":           true,
				},
			},
		},
	}
}

func TestSnapshotAttributesRoundTripSQLite(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	c := &ConnectorRecord{Name: "pve1", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}

	want := testSnapshotWithAttributes()
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	sn := &SnapshotRecord{ConnectorID: c.ID, Data: string(data)}
	if err := s.CreateSnapshot(ctx, sn); err != nil {
		t.Fatalf("CreateSnapshot() error: %v", err)
	}

	got, err := s.GetLatestSnapshot(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetLatestSnapshot() error: %v", err)
	}
	var roundTripped connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(got.Data), &roundTripped); err != nil {
		t.Fatalf("unmarshal round-tripped snapshot: %v", err)
	}
	if len(roundTripped.Entities) != 1 {
		t.Fatalf("Entities = %+v, want 1 entity", roundTripped.Entities)
	}
	attrs := roundTripped.Entities[0].Attributes
	if attrs["status"] != "running" || attrs["firewall_enabled"] != true || attrs["protection"] != false || attrs["onboot"] != true {
		t.Errorf("Attributes = %+v, want status/firewall_enabled/protection/onboot preserved", attrs)
	}
}

func TestSnapshotAttributesRoundTripPostgres(t *testing.T) {
	dsn := os.Getenv("WISELABZ_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("WISELABZ_TEST_POSTGRES_DSN not set; skipping postgres snapshot round-trip test")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer db.Close() //nolint:errcheck
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := RunMigrations(db, "postgres", logger); err != nil {
		t.Fatalf("RunMigrations(postgres) error: %v", err)
	}
	s := New(db, "postgres")
	ctx := context.Background()
	c := &ConnectorRecord{
		ID: uuid.New().String(), Name: "pve1-pg", Category: "virtualization", Type: "proxmox", URL: "https://example.com",
	}
	if err := s.CreateConnector(ctx, c); err != nil {
		t.Fatalf("CreateConnector(postgres) error: %v", err)
	}

	want := testSnapshotWithAttributes()
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	sn := &SnapshotRecord{ConnectorID: c.ID, Data: string(data)}
	if err := s.CreateSnapshot(ctx, sn); err != nil {
		t.Fatalf("CreateSnapshot(postgres) error: %v", err)
	}

	got, err := s.GetLatestSnapshot(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetLatestSnapshot(postgres) error: %v", err)
	}
	var roundTripped connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(got.Data), &roundTripped); err != nil {
		t.Fatalf("unmarshal round-tripped snapshot: %v", err)
	}
	if len(roundTripped.Entities) != 1 || roundTripped.Entities[0].Attributes["status"] != "running" {
		t.Errorf("Entities = %+v, want status=running preserved", roundTripped.Entities)
	}
}

// TestUnmarshalOldSnapshotWithoutAttributes is a backward-compatibility
// regression test: snapshots stored before this field existed must still
// unmarshal cleanly, with Attributes left nil rather than failing to decode.
func TestUnmarshalOldSnapshotWithoutAttributes(t *testing.T) {
	old := `{
		"serviceName": "pve1",
		"type": "proxmox",
		"fetchedAt": "2025-01-01T00:00:00Z",
		"entities": [
			{"kind": "vm", "name": "web-01", "externalId": "100"}
		]
	}`
	var snap connector.ServiceSnapshot
	if err := json.Unmarshal([]byte(old), &snap); err != nil {
		t.Fatalf("unmarshal old snapshot: %v", err)
	}
	if len(snap.Entities) != 1 {
		t.Fatalf("Entities = %+v, want 1 entity", snap.Entities)
	}
	if snap.Entities[0].Attributes != nil {
		t.Errorf("Attributes = %+v, want nil for a pre-attributes snapshot", snap.Entities[0].Attributes)
	}
}
