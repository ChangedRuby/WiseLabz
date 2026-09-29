package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	syshandler "github.com/WiseLabz/wiselabz/internal/api/system"
	"github.com/WiseLabz/wiselabz/internal/notifications"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

type standbyElector struct{ entered chan struct{} }

func (e standbyElector) Campaign(ctx context.Context) error {
	close(e.entered)
	<-ctx.Done()
	return ctx.Err()
}
func (standbyElector) Watch(context.Context) <-chan error { return nil }
func (standbyElector) Close() error                       { return nil }

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newTestLifecycle(t *testing.T) (*lifecycleManager, *syshandler.ReadyState) {
	t.Helper()
	logger := testLogger()
	s := apitest.NewStore(t)
	hub := ws.NewHub()
	dispatcher := notifications.NewDispatcher(s, hub)
	sched := scheduler.New(logger)
	ready := &syshandler.ReadyState{}

	mux := http.NewServeMux()
	h := syshandler.NewHandler(s.DB(), nil, s, nil, "", ready)
	mux.HandleFunc("/healthz", h.Liveness)
	mux.HandleFunc("/readyz", h.Readiness)
	srv := &http.Server{Addr: "127.0.0.1:0", Handler: mux}

	lc := newLifecycleManager(lifecycleDeps{
		Logger:          logger,
		HTTPServer:      srv,
		WSHub:           hub,
		Scheduler:       sched,
		Dispatcher:      dispatcher,
		Store:           s,
		Ready:           ready,
		ShutdownTimeout: 5 * time.Second,
	})
	t.Cleanup(func() {
		if lc.workCtx.Err() == nil {
			if err := lc.Shutdown(); err != nil {
				t.Errorf("cleanup Shutdown: %v", err)
			}
		}
	})
	return lc, ready
}

// startTestLifecycle waits for the listener to bind, then proves the HTTP
// goroutine serves requests. Real HTTP and SQLite resources use real time.
func startTestLifecycle(t *testing.T, lc *lifecycleManager) string {
	t.Helper()
	listening := make(chan net.Addr, 1)
	lc.deps.HTTPServer.BaseContext = func(l net.Listener) context.Context {
		listening <- l.Addr()
		return context.Background()
	}
	lc.Start()
	var addr net.Addr
	select {
	case addr = <-listening:
	case <-time.After(30 * time.Second):
		t.Fatal("HTTP listener did not start")
	}
	url := "http://" + addr.String()
	client := &http.Client{Timeout: 30 * time.Second}
	t.Cleanup(client.CloseIdleConnections)
	resp, err := client.Get(url + "/healthz")
	if err != nil {
		t.Fatalf("startup request: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("startup request status = %d", resp.StatusCode)
	}
	return url
}

func waitForLifecycleSignal(t *testing.T, signal <-chan struct{}, event string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(30 * time.Second):
		t.Fatalf("timed out waiting for %s", event)
	}
}

// TestLifecycleManagerOrderedShutdown starts every managed goroutine under
// the lifecycle manager, then shuts it down and asserts: Shutdown returns
// (proving every goroutine it started actually exited, since it blocks on
// the shared errgroup), readiness flips to not-ready as part of the call,
// and the DB is closed last — a query against the store fails only after
// Shutdown has returned.
func TestLifecycleManagerOrderedShutdown(t *testing.T) {
	lc, ready := newTestLifecycle(t)
	startTestLifecycle(t, lc)

	if ready.NotReady() {
		t.Fatal("readiness flipped to not-ready before Shutdown was called")
	}
	if err := lc.deps.Store.Ping(context.Background()); err != nil {
		t.Fatalf("store not usable before Shutdown: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- lc.Shutdown() }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Shutdown() = %v, want nil", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Shutdown did not return: a goroutine leaked past the errgroup")
	}

	if !ready.NotReady() {
		t.Fatal("readiness was not marked not-ready by Shutdown")
	}

	// DB closed last: any query against it now must fail.
	if err := lc.deps.Store.Ping(context.Background()); err == nil {
		t.Fatal("store.Ping succeeded after Shutdown; DB was not closed")
	}

	// Belt and suspenders: the errgroup itself must report every goroutine
	// as exited (Shutdown already blocked on this, so this returns instantly).
	if err := lc.group.Wait(); err != nil {
		t.Fatalf("group.Wait() after Shutdown = %v, want nil", err)
	}
}

// TestLifecycleManagerShutdownCancelsWorkContext asserts the work context
// used by the WS hub, delivery retrier, and doc lock sweep is canceled by
// Shutdown, independently of the process-wide signal context.
func TestLifecycleManagerShutdownCancelsWorkContext(t *testing.T) {
	lc, _ := newTestLifecycle(t)
	startTestLifecycle(t, lc)

	select {
	case <-lc.workCtx.Done():
		t.Fatal("work context canceled before Shutdown was called")
	default:
	}

	if err := lc.Shutdown(); err != nil {
		t.Fatalf("Shutdown() = %v, want nil", err)
	}

	select {
	case <-lc.workCtx.Done():
	default:
		t.Fatal("work context was not canceled by Shutdown")
	}
}

// manualScheduler advances one-second intervals only when explicitly asked.
// It never starts goroutines or changes the clock used by real I/O.
type manualScheduler struct {
	mu      sync.Mutex
	ctx     context.Context
	starts  int
	runs    int
	elapsed time.Duration
	job     func(context.Context)
}

func (s *manualScheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.starts++
	s.ctx = ctx
}

func (s *manualScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = nil
}

func (s *manualScheduler) advance(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx == nil {
		return
	}
	s.elapsed += d
	for s.elapsed >= time.Second {
		s.elapsed -= time.Second
		s.runs++
		s.job(s.ctx)
	}
}

func (s *manualScheduler) counts() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.starts, s.runs
}

func TestStandbyIsUnreadyAndRunsNoScheduler(t *testing.T) {
	lc, ready := newTestLifecycle(t)
	elector := standbyElector{entered: make(chan struct{})}
	lc.deps.Elector = elector
	lc.deps.LeaderElection = true
	var callbacks int
	sched := &manualScheduler{job: func(context.Context) { callbacks++ }}
	lc.deps.Scheduler = sched
	url := startTestLifecycle(t, lc)
	waitForLifecycleSignal(t, elector.entered, "standby campaign")
	client := &http.Client{Timeout: 30 * time.Second}
	t.Cleanup(client.CloseIdleConnections)
	assertStandby := func() {
		t.Helper()
		resp, err := client.Get(url + "/readyz")
		if err != nil {
			t.Fatal(err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusServiceUnavailable || !ready.WaitingForLeader() {
			t.Fatalf("standby readiness = %d, waiting = %v", resp.StatusCode, ready.WaitingForLeader())
		}
		if starts, runs := sched.counts(); starts != 0 || runs != 0 || callbacks != 0 {
			t.Fatalf("standby scheduler starts=%d, runs=%d, callbacks=%d", starts, runs, callbacks)
		}
	}
	assertStandby()
	sched.advance(2 * time.Second)
	assertStandby()
}

type acquiredElector struct{ watching chan struct{} }

func (acquiredElector) Campaign(context.Context) error { return nil }
func (e acquiredElector) Watch(context.Context) <-chan error {
	close(e.watching)
	return nil
}
func (acquiredElector) Close() error { return nil }

func TestLeaderStartsSchedulerAndRunsJob(t *testing.T) {
	lc, ready := newTestLifecycle(t)
	elector := acquiredElector{watching: make(chan struct{})}
	lc.deps.Elector = elector
	lc.deps.LeaderElection = true
	var callbacks int
	sched := &manualScheduler{job: func(ctx context.Context) {
		if ctx != lc.workCtx || ctx.Err() != nil {
			t.Error("scheduler job did not receive the live work context")
		}
		callbacks++
	}}
	lc.deps.Scheduler = sched
	startTestLifecycle(t, lc)
	// Watch begins after scheduler startup and SetLeaderHeld.
	waitForLifecycleSignal(t, elector.watching, "leader watch")
	if ready.WaitingForLeader() {
		t.Fatal("leader still waiting for leadership")
	}
	if starts, runs := sched.counts(); starts != 1 || runs != 0 {
		t.Fatalf("before advancement starts=%d, runs=%d", starts, runs)
	}
	sched.advance(time.Second)
	if starts, runs := sched.counts(); starts != 1 || runs != 1 || callbacks != 1 {
		t.Fatalf("leader scheduler starts=%d, runs=%d, callbacks=%d", starts, runs, callbacks)
	}
}

type gatedStopScheduler struct {
	lifecycleScheduler
	entered chan struct{}
	release chan struct{}
}

func (s *gatedStopScheduler) Stop() {
	close(s.entered)
	<-s.release
	s.lifecycleScheduler.Stop()
}

func TestLifecycleManagerWaitsForSchedulerBeforeCancelAndDBClose(t *testing.T) {
	lc, ready := newTestLifecycle(t)
	sched := &gatedStopScheduler{
		lifecycleScheduler: lc.deps.Scheduler,
		entered:            make(chan struct{}), release: make(chan struct{}),
	}
	lc.deps.Scheduler = sched
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(sched.release) }) }
	// Release the gate even if a startup assertion fails, before lc cleanup.
	t.Cleanup(release)
	startTestLifecycle(t, lc)
	done := make(chan struct{})
	var shutdownErr error
	go func() {
		shutdownErr = lc.Shutdown()
		close(done)
	}()
	t.Cleanup(func() {
		release()
		select {
		case <-done:
			if shutdownErr != nil {
				t.Errorf("Shutdown: %v", shutdownErr)
			}
		case <-time.After(30 * time.Second):
			t.Error("Shutdown did not finish")
		}
	})
	waitForLifecycleSignal(t, sched.entered, "scheduler Stop")
	if !ready.NotReady() {
		t.Fatal("readiness stayed true while scheduler Stop was blocked")
	}
	h := syshandler.NewHandler(nil, nil, nil, nil, "", ready)
	r := httptest.NewRecorder()
	h.Readiness(r, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if r.Code != http.StatusServiceUnavailable {
		t.Fatalf("draining readiness = %d", r.Code)
	}
	if err := lc.deps.Store.Ping(context.Background()); err != nil {
		t.Fatalf("database closed before scheduler stopped: %v", err)
	}
	if lc.workCtx.Err() != nil {
		t.Fatal("work context canceled before scheduler stopped")
	}
	release()
	waitForLifecycleSignal(t, lc.workCtx.Done(), "work cancellation")
	waitForLifecycleSignal(t, done, "Shutdown completion")
	if shutdownErr != nil {
		t.Fatal(shutdownErr)
	}
	if err := lc.group.Wait(); err != nil {
		t.Fatalf("group.Wait: %v", err)
	}
	if err := lc.deps.Store.Ping(context.Background()); err == nil {
		t.Fatal("database stayed open after Shutdown")
	}
}
