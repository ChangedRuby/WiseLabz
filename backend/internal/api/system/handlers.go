// Package system provides system-level API handlers (health, version, status).
package system

import (
	"net/http"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/robfig/cron/v3"

	"github.com/WiseLabz/wiselabz/internal/config"
	"github.com/WiseLabz/wiselabz/internal/diagnostics"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/scheduler"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// ReadyState is a shared, concurrency-safe readiness flag. main creates one
// and passes it both to the router (whose /readyz handler reads it) and to
// the lifecycle manager (which flips it as the first step of ordered
// shutdown, before anything stops accepting work).
type ReadyState struct {
	notReady atomic.Bool
}

// SetNotReady marks the server as not ready. One-way: once shutdown starts,
// the server never becomes ready again.
func (r *ReadyState) SetNotReady() { r.notReady.Store(true) }

// NotReady reports whether SetNotReady has been called.
func (r *ReadyState) NotReady() bool { return r.notReady.Load() }

// Handler holds dependencies for system endpoints.
type Handler struct {
	DB               store.DBTX
	Config           *config.Config
	Store            *store.Store
	Scheduler        *scheduler.Runner // for re-registering backup/retention jobs
	BackupDir        string            // directory where backups are written
	Ready            *ReadyState       // nil means always-ready (e.g. in tests)
	BackupJobIDMu    sync.Mutex        // protects BackupJobID
	BackupJobID      cron.EntryID      // current backup job entry ID (0 if not registered)
	RetentionJobIDMu sync.Mutex        // protects RetentionJobID
	RetentionJobID   cron.EntryID      // current retention job entry ID (0 if not registered)
}

// NewHandler creates a new system handler. ready may be nil, in which case
// the /readyz endpoint never reports the not-ready-for-shutdown state.
func NewHandler(db store.DBTX, cfg *config.Config, s *store.Store, scheduler *scheduler.Runner, backupDir string, ready *ReadyState) *Handler {
	return &Handler{
		DB:        db,
		Config:    cfg,
		Store:     s,
		Scheduler: scheduler,
		BackupDir: backupDir,
		Ready:     ready,
	}
}

// Health responds with the server health status.
// GET /api/health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	health := diagnostics.CheckHealth(r.Context(), h.DB)
	httputil.JSON(w, http.StatusOK, map[string]any{
		"status":     health.Status,
		"healthy":    health.Status == "ok",
		"components": health.Components,
	})
}

// Liveness responds whenever the HTTP server is running.
// GET /healthz
func (h *Handler) Liveness(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readiness responds once the database is reachable and all migrations are applied.
// GET /readyz
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	if h.Ready != nil && h.Ready.NotReady() {
		httputil.JSON(w, http.StatusServiceUnavailable, map[string]any{
			"status":     "degraded",
			"ready":      false,
			"components": []diagnostics.Component{{Name: "shutdown", Status: "draining"}},
		})
		return
	}

	health := diagnostics.CheckHealth(r.Context(), h.DB)
	components := health.Components
	ready := health.Status == "ok"

	migrationStatus := "unknown"
	if ready {
		status, err := h.Store.MigrationStatus()
		switch {
		case err != nil:
			ready = false
		case status.Dirty:
			migrationStatus = "dirty"
			ready = false
		case status.Pending():
			migrationStatus = "pending"
			ready = false
		default:
			migrationStatus = "ok"
		}
	}
	components = append(components, diagnostics.Component{Name: "migrations", Status: migrationStatus})

	code := http.StatusOK
	overallStatus := "ok"
	if !ready {
		code = http.StatusServiceUnavailable
		overallStatus = "degraded"
	}
	httputil.JSON(w, code, map[string]any{
		"status":     overallStatus,
		"ready":      ready,
		"components": components,
	})
}

// Info responds with instance info (version, sync schedule, integrations).
// GET /api/system/info
func (h *Handler) Info(w http.ResponseWriter, _ *http.Request) {
	version := "dev"
	if buildInfo, ok := debug.ReadBuildInfo(); ok && buildInfo.Main.Version != "" {
		version = buildInfo.Main.Version
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"version":      version,
		"syncSchedule": h.Config.Sync.Schedule,
		"integrations": []map[string]any{},
	})
}
