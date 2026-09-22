package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	syshandler "github.com/WiseLabz/wiselabz/internal/api/system"
	"github.com/WiseLabz/wiselabz/internal/notifications"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

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

	srv := &http.Server{Addr: "127.0.0.1:0", Handler: http.NewServeMux()}

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
	return lc, ready
}

// TestLifecycleManagerOrderedShutdown starts every managed goroutine under
// the lifecycle manager, then shuts it down and asserts: Shutdown returns
// (proving every goroutine it started actually exited, since it blocks on
// the shared errgroup), readiness flips to not-ready as part of the call,
// and the DB is closed last — a query against the store fails only after
// Shutdown has returned.
func TestLifecycleManagerOrderedShutdown(t *testing.T) {
	lc, ready := newTestLifecycle(t)
	lc.Start()

	// Give the goroutines a moment to actually start (HTTP listener bound,
	// scheduler running) before tearing down.
	time.Sleep(50 * time.Millisecond)

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
	case <-time.After(10 * time.Second):
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
	lc.Start()
	time.Sleep(20 * time.Millisecond)

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
