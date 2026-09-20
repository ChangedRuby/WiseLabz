// Package sync provides the sync engine for fetching connector data,
// diffing snapshots, and creating changes/alerts.
package sync

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/WiseLabz/wiselabz/internal/ws"
)

// maxSyncConcurrency bounds how many connectors RunSyncAll fetches at once,
// matching notifications.maxConcurrentNotifications' fanout pattern.
const maxSyncConcurrency = 4

// ErrAlreadyRunning means this connector already has an active sync.
var ErrAlreadyRunning = errors.New("connector sync already running")

// AlertNotifier dispatches notifications for a newly created alert.
type AlertNotifier interface {
	NotifyAlertsCreated(ctx context.Context, alerts []store.AlertRecord)
}

// QualityChecker evaluates documentation quality after a sync attempt.
type QualityChecker interface {
	RunForConnector(ctx context.Context, connectorID string) error
}

// DocRegenerator re-renders a connector's existing docs from its latest
// snapshot after a sync attempt, recording a new "sync"-triggered doc
// version for any doc whose content changed.
type DocRegenerator interface {
	RegenerateForConnector(ctx context.Context, connectorID string) error
}

// Engine runs sync jobs against connectors.
type Engine struct {
	inFlight       sync.Map // connector ID -> active run; entries are removed on every exit
	store          *store.Store
	hub            *ws.Hub
	notifier       AlertNotifier
	qualityChecker QualityChecker
	docRegenerator DocRegenerator
	// encKey is the base64-encoded AES-256 key (config.Encryption.Key) used
	// to decrypt/re-encrypt secret-bearing connector config fields via
	// store.ParseConnectorConfig/MarshalConnectorConfig.
	encKey string
	// baseCtx parents detached (request-triggered) syncs so they outlive the
	// HTTP request but are cancelled on server shutdown. Defaults to
	// context.Background() until SetBaseContext is called.
	baseCtx context.Context
}

// NewEngine creates a new sync engine.
func NewEngine(s *store.Store, h *ws.Hub, notifier AlertNotifier, qualityChecker QualityChecker, encKey string) *Engine {
	return &Engine{store: s, hub: h, notifier: notifier, qualityChecker: qualityChecker, encKey: encKey, baseCtx: context.Background()}
}

// SetDocRegenerator wires a DocRegenerator into the engine after
// construction. Kept as a setter (rather than a NewEngine parameter) so
// existing call sites don't need to change; a nil regenerator (the default)
// simply skips sync-triggered doc regeneration.
func (e *Engine) SetDocRegenerator(dr DocRegenerator) {
	e.docRegenerator = dr
}

// SetBaseContext sets the context that detached syncs derive from. main wires
// the signal-aware server context here so shutdown cancels in-flight syncs.
func (e *Engine) SetBaseContext(ctx context.Context) {
	e.baseCtx = ctx
}

// BaseContext returns the context detached syncs should run under: it survives
// the triggering HTTP request and is cancelled on shutdown. Each run is still
// bounded by syncTimeout inside runSyncFields.
func (e *Engine) BaseContext() context.Context {
	if e.baseCtx == nil {
		return context.Background()
	}
	return e.baseCtx
}

// RunResult holds the outcome of a sync run.
type RunResult struct {
	ConnectorID  string `json:"connectorId"`
	SnapshotID   string `json:"snapshotId"`
	ChangesCount int    `json:"changesCount"`
	AlertsCount  int    `json:"alertsCount"`
	Status       string `json:"status"` // "success", "error"
	Error        string `json:"error,omitempty"`
	Duration     string `json:"duration"`
}

// RunSyncAll runs sync for all enabled connectors.
func (e *Engine) RunSyncAll(ctx context.Context, jobID string) ([]RunResult, error) {
	connectors, err := e.store.ListAllConnectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("list connectors: %w", err)
	}

	// Fetch+write each connector concurrently, bounded by a semaphore, mirroring
	// notifications.Dispatcher's fanoutSem pattern. Results are written to
	// per-index slots so ordering stays deterministic (matching connectors order)
	// despite concurrent completion.
	sem := make(chan struct{}, maxSyncConcurrency)
	slots := make([]*RunResult, len(connectors))
	var wg sync.WaitGroup
	for i, c := range connectors {
		if !c.Enabled {
			continue
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(i int, connectorID string) {
			defer wg.Done()
			defer func() { <-sem }()
			result, err := e.RunSync(ctx, connectorID, jobID)
			if err != nil {
				slog.Error("sync failed", "connector", logsafe.Sanitize(connectorID), "error", logsafe.Sanitize(err.Error()))
				return
			}
			slots[i] = result
		}(i, c.ID)
	}
	wg.Wait()

	var results []RunResult
	for _, r := range slots {
		if r != nil {
			results = append(results, *r)
		}
	}
	return results, nil
}
