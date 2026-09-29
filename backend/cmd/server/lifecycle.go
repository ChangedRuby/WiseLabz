package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	syshandler "github.com/WiseLabz/wiselabz/internal/api/system"
	"github.com/WiseLabz/wiselabz/internal/leader"
	"github.com/WiseLabz/wiselabz/internal/notifications"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

type lifecycleScheduler interface {
	Start(context.Context)
	Stop()
}

// lifecycleDeps holds everything the lifecycle manager needs to start and
// stop the server's long-running goroutines.
type lifecycleDeps struct {
	Logger          *slog.Logger
	HTTPServer      *http.Server
	WSHub           *ws.Hub
	Scheduler       lifecycleScheduler
	Dispatcher      *notifications.Dispatcher
	Store           *store.Store
	Ready           *syshandler.ReadyState
	Elector         leader.Election
	LeaderElection  bool
	ShutdownTimeout time.Duration
}

// lifecycleManager starts every long-running server goroutine (HTTP server,
// WebSocket hub, scheduler, notification delivery retries, doc lock sweep)
// under one errgroup with a shared "work" context, and on Shutdown runs an
// ORDERED stop:
//
//  1. mark readiness as not-ready (so a load balancer stops routing here)
//  2. stop accepting new HTTP/WebSocket work and drain in-flight requests
//  3. stop the scheduler (blocks until any in-flight cron job finishes)
//  4. cancel the work context and wait for every remaining goroutine to exit
//  5. wait for in-flight notification dispatch goroutines
//  6. close the DB last
//
// The work context is deliberately separate from the process's signal
// context: goroutines must keep running while Shutdown works through steps
// 2-3 in order, not all cancel out simultaneously the instant a signal
// arrives.
type lifecycleManager struct {
	deps lifecycleDeps

	workCtx    context.Context
	workCancel context.CancelFunc
	group      *errgroup.Group
	errors     chan error
	mu         sync.Mutex
	stopping   bool
}

// newLifecycleManager builds a manager. Call Start to launch the goroutines,
// then Shutdown to stop them in order.
func newLifecycleManager(deps lifecycleDeps) *lifecycleManager {
	workCtx, workCancel := context.WithCancel(context.Background())
	group := new(errgroup.Group)
	return &lifecycleManager{
		deps:       deps,
		workCtx:    workCtx,
		workCancel: workCancel,
		group:      group,
		errors:     make(chan error, 1),
	}
}

// Errors reports a fatal leadership loss so main can shut down and exit.
func (m *lifecycleManager) Errors() <-chan error { return m.errors }

// Start launches every long-running goroutine under the errgroup. It returns
// immediately; the goroutines run until Shutdown is called.
func (m *lifecycleManager) Start() {
	logger := m.deps.Logger

	m.group.Go(func() error {
		m.deps.WSHub.Run(m.workCtx)
		return nil
	})

	if m.deps.LeaderElection {
		m.deps.Ready.RequireLeader()
		m.group.Go(func() error {
			if err := m.deps.Elector.Campaign(m.workCtx); err != nil {
				if m.workCtx.Err() == nil {
					m.fail(err)
					return err
				}
				return nil
			}
			m.mu.Lock()
			if m.stopping {
				m.mu.Unlock()
				return nil
			}
			m.startLeaderWorkers()
			m.deps.Ready.SetLeaderHeld(true)
			m.mu.Unlock()
			logger.Info("leader acquired; background workers started")
			select {
			case <-m.workCtx.Done():
				return nil
			case err := <-m.deps.Elector.Watch(m.workCtx):
				if err != nil {
					m.deps.Ready.SetLeaderHeld(false)
					logger.Error("leader lock lost", "error", err)
					m.fail(err)
					return err
				}
				return nil
			}
		})
	} else {
		m.startLeaderWorkers()
	}

	m.group.Go(func() error {
		logger.Info("HTTP server listening", "addr", m.deps.HTTPServer.Addr)
		if err := m.deps.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
		}
		return nil
	})
}

func (m *lifecycleManager) fail(err error) {
	select {
	case m.errors <- err:
	default:
	}
}

func (m *lifecycleManager) startLeaderWorkers() {
	logger := m.deps.Logger
	m.deps.Scheduler.Start(m.workCtx)
	m.group.Go(func() error {
		notifications.RunDeliveryRetries(m.workCtx, m.deps.Dispatcher, logger)
		return nil
	})
	m.group.Go(func() error {
		store.RunDocLockSweep(m.workCtx, m.deps.Store, m.deps.WSHub, store.DocLockHeartbeat, logger)
		return nil
	})
}

// Shutdown runs the ordered stop described on lifecycleManager and blocks
// until every goroutine started by Start has exited and the DB is closed.
func (m *lifecycleManager) Shutdown() error {
	logger := m.deps.Logger

	// 1. Mark not-ready first so /readyz starts failing before anything else
	// changes, giving a load balancer a head start on draining traffic.
	if m.deps.Ready != nil {
		m.deps.Ready.SetNotReady()
	}
	m.mu.Lock()
	m.stopping = true
	m.mu.Unlock()

	// 2. Stop accepting new HTTP/WebSocket upgrade requests and drain
	// in-flight ones.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), m.deps.ShutdownTimeout)
	defer shutdownCancel()
	if err := m.deps.HTTPServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}

	// 3. Stop the scheduler; blocks until any in-flight job (backup,
	// retention, quality, sync, digest, alert expiry) finishes.
	m.deps.Scheduler.Stop()

	// 4. Cancel the work context so the WS hub, delivery retrier, and doc
	// lock sweep goroutines return, then wait for every goroutine started by
	// Start (including the two above) to actually exit.
	m.workCancel()
	if err := m.group.Wait(); err != nil {
		logger.Error("lifecycle group error", "error", err)
	}

	// 5. Wait for in-flight notification dispatch goroutines (e.g. an alert
	// created just before shutdown) so they don't touch a closed DB.
	m.deps.Dispatcher.Wait()
	// Release the session lock before closing the pool.
	if m.deps.Elector != nil {
		if err := m.deps.Elector.Close(); err != nil {
			logger.Error("leader session close error", "error", err)
		}
	}

	// 6. Close the DB last, now that nothing above can still be using it.
	return m.deps.Store.Close()
}
