package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// fakeHealthStore is an in-memory HealthStore for tests, so health tracking
// can be exercised without a real database.
type fakeHealthStore struct {
	mu   sync.Mutex
	recs map[string]store.JobHealthRecord
}

func newFakeHealthStore() *fakeHealthStore {
	return &fakeHealthStore{recs: make(map[string]store.JobHealthRecord)}
}

func (f *fakeHealthStore) GetJobHealth(_ context.Context, name string) (store.JobHealthRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	rec, ok := f.recs[name]
	if !ok {
		return store.JobHealthRecord{}, sql.ErrNoRows
	}
	return rec, nil
}

func (f *fakeHealthStore) UpsertJobHealth(_ context.Context, rec store.JobHealthRecord) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.recs[rec.Name] = rec
	return nil
}

func (f *fakeHealthStore) ListJobHealth(_ context.Context) ([]store.JobHealthRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	recs := make([]store.JobHealthRecord, 0, len(f.recs))
	for _, rec := range f.recs {
		recs = append(recs, rec)
	}
	return recs, nil
}

// seed pre-populates a health row, used to simulate state carried over a
// restart (a fresh fakeHealthStore/Runner but a pre-existing row).
func (f *fakeHealthStore) seed(rec store.JobHealthRecord) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.recs[rec.Name] = rec
}

type notifyCall struct{ eventType, severity, title string }

type fakeNotifier struct {
	mu    sync.Mutex
	calls []notifyCall
}

func (f *fakeNotifier) NotifySystemEvent(_ context.Context, eventType, severity, title, _ string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, notifyCall{eventType, severity, title})
}

func (f *fakeNotifier) snapshot() []notifyCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]notifyCall(nil), f.calls...)
}

func TestJobHealthOkToFailingNotifiesOnce(t *testing.T) {
	logger := testLogger()
	r := New(logger)
	hs := newFakeHealthStore()
	notifier := &fakeNotifier{}
	r.SetHealthTracking(hs, notifier)

	fail := true
	id, err := r.AddJob("flaky", "0 0 * * * *", func(context.Context) error {
		if fail {
			return fmt.Errorf("boom")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	r.c.Entry(id).WrappedJob.Run()
	if calls := notifier.snapshot(); len(calls) != 1 || calls[0].severity != "warning" {
		t.Fatalf("after first failure: notifications = %+v, want one warning", calls)
	}
	rec, err := hs.GetJobHealth(context.Background(), "flaky")
	if err != nil || rec.LastStatus != "failing" {
		t.Fatalf("GetJobHealth() = %+v, %v; want status failing", rec, err)
	}

	// A second consecutive failure must not notify again (failing->failing is silent).
	r.c.Entry(id).WrappedJob.Run()
	if calls := notifier.snapshot(); len(calls) != 1 {
		t.Fatalf("after second failure: notifications = %+v, want still exactly one", calls)
	}

	// Recovery notifies exactly once.
	fail = false
	r.c.Entry(id).WrappedJob.Run()
	calls := notifier.snapshot()
	if len(calls) != 2 || calls[1].severity != "info" {
		t.Fatalf("after recovery: notifications = %+v, want warning then one info", calls)
	}
	rec, err = hs.GetJobHealth(context.Background(), "flaky")
	if err != nil || rec.LastStatus != "ok" || rec.LastError != "" {
		t.Fatalf("GetJobHealth() after recovery = %+v, %v; want status ok, no error", rec, err)
	}

	// A second consecutive success must not notify again (ok->ok is silent).
	r.c.Entry(id).WrappedJob.Run()
	if calls := notifier.snapshot(); len(calls) != 2 {
		t.Fatalf("after second success: notifications = %+v, want still exactly two", calls)
	}
}

func TestJobHealthPersistsAcrossRestart(t *testing.T) {
	hs := newFakeHealthStore()
	// Simulate: the job was already failing before the process restarted.
	hs.seed(store.JobHealthRecord{
		Name: "flaky", CronExpr: "0 0 * * * *", LastStatus: "failing",
		LastError: "boom", UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	})

	// A brand-new Runner (as a restarted process would create) with the
	// same persisted store: the in-memory "already notified" state that
	// used to live on docexport.Exporter.failing is gone, but the
	// persisted row means a success doesn't get treated as a no-op.
	r := New(testLogger())
	notifier := &fakeNotifier{}
	r.SetHealthTracking(hs, notifier)

	id, err := r.AddJob("flaky", "0 0 * * * *", func(context.Context) error { return nil })
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	r.c.Entry(id).WrappedJob.Run()
	calls := notifier.snapshot()
	if len(calls) != 1 || calls[0].severity != "info" {
		t.Fatalf("after restart+success: notifications = %+v, want exactly one recovery info (not a duplicate warning)", calls)
	}
}

func TestJobHealthPanicCountsAsFailure(t *testing.T) {
	logger := testLogger()
	r := New(logger)
	hs := newFakeHealthStore()
	notifier := &fakeNotifier{}
	r.SetHealthTracking(hs, notifier)

	id, err := r.AddJob("panicky", "0 0 * * * *", func(context.Context) error {
		panic("kaboom")
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}

	r.c.Entry(id).WrappedJob.Run()

	rec, err := hs.GetJobHealth(context.Background(), "panicky")
	if err != nil || rec.LastStatus != "failing" {
		t.Fatalf("GetJobHealth() = %+v, %v; want status failing after panic", rec, err)
	}
	if calls := notifier.snapshot(); len(calls) != 1 || calls[0].severity != "warning" {
		t.Fatalf("after panic: notifications = %+v, want one warning", calls)
	}
}

func TestJobHealthWithoutStoreDoesNothing(t *testing.T) {
	// No SetHealthTracking call: AddJob/run must still work exactly as
	// before (health tracking is opt-in).
	r := New(testLogger())
	ran := false
	id, err := r.AddJob("plain", "0 0 * * * *", func(context.Context) error {
		ran = true
		return nil
	})
	if err != nil {
		t.Fatalf("AddJob() error: %v", err)
	}
	r.c.Entry(id).WrappedJob.Run()
	if !ran {
		t.Fatal("job did not run")
	}
}
