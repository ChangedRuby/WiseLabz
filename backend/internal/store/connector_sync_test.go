package store

import (
	"context"
	"errors"
	"testing"
)

func TestSnapshotKeysetSummaryAndConnectorOwnership(t *testing.T) {
	ctx := context.Background()
	s := newDocTestStore(t)
	connectorID := createTestConnector(ctx, t, s)
	otherID := createTestConnector(ctx, t, s)
	old := &SnapshotRecord{ConnectorID: connectorID, FetchedAt: "2026-01-01T00:00:00Z", Data: "é"}
	newer := &SnapshotRecord{ConnectorID: connectorID, FetchedAt: old.FetchedAt, Data: "new"}
	for _, sn := range []*SnapshotRecord{old, newer, {ConnectorID: otherID, Data: "elsewhere"}} {
		if err := s.CreateSnapshot(ctx, sn); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.PinGoldenSnapshot(ctx, &GoldenSnapshotRecord{ConnectorID: connectorID, SnapshotID: old.ID, PinnedBy: "tester"}); err != nil {
		t.Fatal(err)
	}
	all, total, err := s.ListSnapshotsByConnectorKeyset(ctx, connectorID, Keyset{}, 10)
	if err != nil || total != 2 || len(all) != 2 {
		t.Fatalf("list = %+v, total %d, err %v", all, total, err)
	}
	for _, item := range all {
		if item.ID == old.ID && (!item.Golden || item.SizeBytes != 2) {
			t.Fatalf("old summary = %+v, want golden and 2 UTF-8 bytes", item)
		}
	}
	page, _, err := s.ListSnapshotsByConnectorKeyset(ctx, connectorID, Keyset{}, 1)
	if err != nil || len(page) != 1 || page[0].ID != all[0].ID {
		t.Fatalf("first page = %+v, err %v", page, err)
	}
	page, total, err = s.ListSnapshotsByConnectorKeyset(ctx, connectorID, Keyset{Sort: page[0].FetchedAt, ID: page[0].ID}, 1)
	if err != nil || total != 2 || len(page) != 1 || page[0].ID != all[1].ID {
		t.Fatalf("second page = %+v, total %d, err %v", page, total, err)
	}
	if _, err := s.GetSnapshotForConnector(ctx, otherID, old.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-connector lookup = %v, want ErrNotFound", err)
	}
	if got, err := s.GetSnapshotForConnector(ctx, connectorID, old.ID); err != nil || got.Data != old.Data {
		t.Fatalf("owned lookup = %+v, err %v", got, err)
	}
}
