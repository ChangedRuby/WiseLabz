package sync

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

type blockingConnector struct {
	fakeConnector
	entered chan context.Context
	release chan struct{}
}

func (c *blockingConnector) Fetch(ctx context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	c.entered <- ctx
	select {
	case <-c.release:
		return &connector.ServiceSnapshot{ServiceName: "test", FetchedAt: time.Now()}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func TestSyncExcludesConcurrentRuns(t *testing.T) {
	s := newTestStore(t)
	c := &blockingConnector{entered: make(chan context.Context, 4), release: make(chan struct{})}
	connector.Register(connector.TypeSchema{Type: "blocking_sync", Category: "networking"}, func(map[string]any) (connector.Connector, error) { return c, nil })
	schedule := 30
	rec := &store.ConnectorRecord{Name: "test", Type: "blocking_sync", Category: "networking", Enabled: true, ScheduleSeconds: &schedule}
	if err := s.CreateConnector(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	e := NewEngine(s, nil, nil, nil, "")
	done := make(chan struct{})
	go func() { defer close(done); e.RunDueSyncs(context.Background(), slog.Default()) }()
	select {
	case ctx := <-c.entered:
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > syncTimeout {
			t.Error("sync must have a bounded deadline")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("scheduled fetch did not start")
	}
	for _, partial := range []bool{false, true} {
		var err error
		if partial {
			_, err = e.RunSyncFields(context.Background(), rec.ID, "partial", []string{"vms"})
		} else {
			_, err = e.RunSync(context.Background(), rec.ID, "manual")
		}
		if !errors.Is(err, ErrAlreadyRunning) {
			t.Errorf("concurrent sync error = %v", err)
		}
	}
	// A different connector can still sync while this one is blocked.
	connector.Register(connector.TypeSchema{Type: "independent_sync", Category: "networking"}, func(map[string]any) (connector.Connector, error) {
		return &fakeConnector{snapshot: &connector.ServiceSnapshot{ServiceName: "other", FetchedAt: time.Now()}}, nil
	})
	other := &store.ConnectorRecord{Name: "other", Type: "independent_sync", Category: "networking", Enabled: true}
	if err := s.CreateConnector(context.Background(), other); err != nil {
		t.Fatal(err)
	}
	if _, err := e.RunSync(context.Background(), other.ID, "independent"); err != nil {
		t.Fatal(err)
	}
	// Another engine cannot claim the same scheduled run while its lease is live.
	NewEngine(s, nil, nil, nil, "").RunDueSyncs(context.Background(), slog.Default())
	close(c.release)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("sync did not finish")
	}
	if _, err := e.RunSync(context.Background(), rec.ID, "after"); err != nil {
		t.Fatal(err)
	}
	snapshots, err := s.GetSnapshotsByConnector(context.Background(), rec.ID, 10)
	if err != nil || len(snapshots) != 2 {
		t.Fatalf("snapshots = %d, err = %v; want two actual runs", len(snapshots), err)
	}
}

func TestSyncCancellationRecordsFailureAndReleasesGuard(t *testing.T) {
	s := newTestStore(t)
	c := &blockingConnector{entered: make(chan context.Context, 2), release: make(chan struct{})}
	connector.Register(connector.TypeSchema{Type: "cancel_sync", Category: "networking"}, func(map[string]any) (connector.Connector, error) { return c, nil })
	schedule := 30
	rec := &store.ConnectorRecord{Name: "test", Type: "cancel_sync", Category: "networking", Enabled: true, ScheduleSeconds: &schedule}
	if err := s.CreateConnector(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	e := NewEngine(s, nil, nil, nil, "")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := e.RunSync(ctx, rec.ID, "cancel"); done <- err }()
	select {
	case <-c.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("fetch did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cancellation did not finish")
	}
	saved, err := s.GetConnector(context.Background(), rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if saved.RetryCount != 1 || saved.LastSyncError == "" || saved.NextRunAt == "" {
		t.Fatalf("failure not persisted: %+v", saved)
	}
	close(c.release)
	if _, err := e.RunSync(context.Background(), rec.ID, "retry"); err != nil {
		t.Fatal(err)
	}
}
