// Package scheduler provides a shared cron-based scheduling primitive for
// running background jobs at configured intervals. It also centralizes
// per-job health tracking (#384): every job registered via AddJob returns an
// error, which the Runner persists (when a HealthStore is set) and turns
// into a system.job_failed notification on ok<->failing transitions only
// (when a Notifier is set), instead of leaving each job to hand-roll its own
// in-memory failing flag.
package scheduler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/WiseLabz/wiselabz/internal/logsafe"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// EventJobFailed is the notification event type sent when a scheduled job
// starts failing (severity warning) or recovers (severity info). It must
// match notifications.EventSystemJobFailed — kept as a separate constant
// there to avoid an import cycle (notifications doesn't depend on
// scheduler).
const EventJobFailed = "system.job_failed"

// HealthStore is the slice of store.Store the scheduler uses to persist job
// health across restarts. *store.Store satisfies this directly.
type HealthStore interface {
	GetJobHealth(ctx context.Context, name string) (store.JobHealthRecord, error)
	UpsertJobHealth(ctx context.Context, rec store.JobHealthRecord) error
	ListJobHealth(ctx context.Context) ([]store.JobHealthRecord, error)
}

// Notifier is the slice of notifications.Dispatcher the scheduler uses to
// announce job failure/recovery transitions. Same shape as
// docexport.Notifier, which it replaces.
type Notifier interface {
	NotifySystemEvent(ctx context.Context, eventType, severity, title, message string)
}

// Runner manages a set of cron-scheduled jobs.
type Runner struct {
	c      *cron.Cron
	logger *slog.Logger
	mu     sync.Mutex // protects access to the cron instance, ctx, and entries
	ctx    context.Context
	// entries maps a registered job's name to bookkeeping ListJobs needs:
	// its cron expression (for display) and cron entry ID (to ask the cron
	// library for the next scheduled run).
	entries map[string]jobEntry

	healthStore HealthStore
	notifier    Notifier
}

type jobEntry struct {
	id       cron.EntryID
	cronExpr string
}

// New creates a new Runner with support for both standard 5-field cron
// expressions and 6-field expressions with a leading seconds field (needed
// for sub-minute cadence), matching config.validateCronExpressions.
func New(logger *slog.Logger) *Runner {
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)
	return &Runner{
		c:       cron.New(cron.WithParser(parser), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))),
		logger:  logger,
		ctx:     context.Background(),
		entries: make(map[string]jobEntry),
	}
}

// SetHealthTracking installs the store and notifier used to persist job
// health and announce ok<->failing transitions. Either may be nil: with no
// store, health is only logged (never persisted); with no notifier, no
// system.job_failed events are sent. Must be called before Start; it is not
// safe to call concurrently with running jobs.
func (r *Runner) SetHealthTracking(store HealthStore, notifier Notifier) {
	r.healthStore = store
	r.notifier = notifier
}

// AddJob registers a new cron job with the given name and expression. The
// function fn is wrapped to recover panics (counted as a failure) and log
// entry/exit; its returned error is used to update job_health and, on an
// ok->failing or failing->ok transition, fire a system.job_failed
// notification. Returns the cron entry ID for later removal, or an error if
// the cron expression is invalid.
func (r *Runner) AddJob(name, cronExpr string, fn func(ctx context.Context) error) (cron.EntryID, error) {
	// Wrap the function to recover panics, log entry/exit, and track health.
	wrapped := func() {
		r.logger.Debug("job started", "job", name)
		var jobErr error
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					r.logger.Error("job panicked", "job", name, "panic", rec)
					jobErr = fmt.Errorf("panic: %v", rec)
				}
			}()
			// Use the context passed to Start so job bodies observe
			// cancellation on server shutdown, the same as the old ticker
			// loops did.
			jobErr = fn(r.jobContext())
		}()
		if jobErr != nil {
			r.logger.Error("job failed", "job", name, "error", jobErr)
		}
		r.trackHealth(name, cronExpr, jobErr)
		r.logger.Debug("job completed", "job", name)
	}

	// Use AddFunc which takes the cron expression directly and validates it
	entryID, err := r.c.AddFunc(cronExpr, wrapped)
	if err != nil {
		return 0, fmt.Errorf("add job %q (cron %q): %w", name, cronExpr, err)
	}

	r.mu.Lock()
	r.entries[name] = jobEntry{id: entryID, cronExpr: cronExpr}
	r.mu.Unlock()

	r.logger.Debug("job registered", "job", name, "cron", logsafe.Sanitize(cronExpr))
	return entryID, nil
}

// trackHealth persists the outcome of one job run and fires a
// system.job_failed notification on an ok<->failing transition only. Store
// errors are logged and never block the job or its caller — health tracking
// is best-effort.
func (r *Runner) trackHealth(name, cronExpr string, jobErr error) {
	if r.healthStore == nil {
		return
	}
	ctx := r.jobContext()

	prevStatus := "ok" // no prior row (first run ever) behaves like a healthy baseline
	var prev store.JobHealthRecord
	prev, err := r.healthStore.GetJobHealth(ctx, name)
	switch {
	case err == nil:
		prevStatus = prev.LastStatus
	case errors.Is(err, sql.ErrNoRows):
		// No prior row: first run ever, prevStatus stays "ok" so a first
		// failure still notifies.
	default:
		// Any other error is unexpected but must not block the job; just
		// log and fall back to prevStatus == "ok" so we never suppress a
		// legitimate first "failing" notification.
		r.logger.Error("job health: read previous status", "job", name, "error", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	rec := store.JobHealthRecord{
		Name:          name,
		CronExpr:      cronExpr,
		LastRunAt:     now,
		UpdatedAt:     now,
		LastSuccessAt: prev.LastSuccessAt,
		LastFailureAt: prev.LastFailureAt,
	}

	if jobErr != nil {
		rec.LastStatus = "failing"
		rec.LastError = jobErr.Error()
		rec.LastFailureAt = now
	} else {
		rec.LastStatus = "ok"
		rec.LastSuccessAt = now
	}

	if err := r.healthStore.UpsertJobHealth(ctx, rec); err != nil {
		r.logger.Error("job health: persist status", "job", name, "error", err)
	}

	if r.notifier == nil || prevStatus == rec.LastStatus {
		return
	}
	switch rec.LastStatus {
	case "failing":
		r.notifier.NotifySystemEvent(ctx, EventJobFailed, "warning",
			fmt.Sprintf("Scheduled job %s failing", name),
			fmt.Sprintf("The scheduled job %q failed and will retry on its next run: %s", name, rec.LastError))
	case "ok":
		r.notifier.NotifySystemEvent(ctx, EventJobFailed, "info",
			fmt.Sprintf("Scheduled job %s recovered", name),
			fmt.Sprintf("The scheduled job %q succeeded again.", name))
	}
}

// RemoveJob removes a previously registered job by its entry ID.
func (r *Runner) RemoveJob(id cron.EntryID) {
	r.c.Remove(id)
	r.mu.Lock()
	for name, e := range r.entries {
		if e.id == id {
			delete(r.entries, name)
			break
		}
	}
	r.mu.Unlock()
	r.logger.Debug("job removed", "entry_id", id)
}

// JobInfo describes one registered job for GET /api/system/jobs, combining
// its live cron schedule (name, next run) with its persisted health.
type JobInfo struct {
	Name          string `json:"name"`
	CronExpr      string `json:"cronExpr"`
	LastStatus    string `json:"lastStatus,omitempty"`
	LastError     string `json:"lastError,omitempty"`
	LastRunAt     string `json:"lastRunAt,omitempty"`
	LastSuccessAt string `json:"lastSuccessAt,omitempty"`
	LastFailureAt string `json:"lastFailureAt,omitempty"`
	NextRunAt     string `json:"nextRunAt,omitempty"`
}

// ListJobs returns one JobInfo per currently-registered job, combining its
// live next-run time from the cron entry with its persisted health record
// (if any — a job that has never run yet has no health row). Jobs with no
// health row yet (never fired) still appear, with only Name/CronExpr/NextRunAt set.
func (r *Runner) ListJobs(ctx context.Context) ([]JobInfo, error) {
	r.mu.Lock()
	snapshot := make(map[string]jobEntry, len(r.entries))
	for name, e := range r.entries {
		snapshot[name] = e
	}
	r.mu.Unlock()

	var health map[string]store.JobHealthRecord
	if r.healthStore != nil {
		recs, err := r.healthStore.ListJobHealth(ctx)
		if err != nil {
			return nil, fmt.Errorf("list job health: %w", err)
		}
		health = make(map[string]store.JobHealthRecord, len(recs))
		for _, rec := range recs {
			health[rec.Name] = rec
		}
	}

	infos := make([]JobInfo, 0, len(snapshot))
	for name, e := range snapshot {
		info := JobInfo{Name: name, CronExpr: e.cronExpr}
		if rec, ok := health[name]; ok {
			info.LastStatus = rec.LastStatus
			info.LastError = rec.LastError
			info.LastRunAt = rec.LastRunAt
			info.LastSuccessAt = rec.LastSuccessAt
			info.LastFailureAt = rec.LastFailureAt
		}
		if entry := r.c.Entry(e.id); entry.ID != 0 && !entry.Next.IsZero() {
			info.NextRunAt = entry.Next.UTC().Format(time.RFC3339)
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// EntryCount returns the number of currently registered jobs. Exported
// mainly for tests that need to assert a re-registration replaced a job
// instead of stacking a duplicate alongside it.
func (r *Runner) EntryCount() int {
	return len(r.c.Entries())
}

// jobContext returns the context most recently passed to Start, so job
// bodies can observe cancellation from the server shutdown context.
func (r *Runner) jobContext() context.Context {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ctx
}

// Start starts the cron runner and watches the given context for cancellation.
// When ctx.Done() is triggered, Stop() is called automatically. Jobs added
// via AddJob receive this ctx (or a later one from a subsequent Start call).
func (r *Runner) Start(ctx context.Context) {
	r.mu.Lock()
	r.ctx = ctx
	r.logger.Info("Scheduler started")
	r.c.Start()
	r.mu.Unlock()

	// Watch for context cancellation in a separate goroutine
	go func() {
		<-ctx.Done()
		r.Stop()
	}()
}

// Stop gracefully stops the cron runner, blocking until any in-flight job
// finishes. This lets callers (e.g. server shutdown) safely close shared
// resources like the DB right after Stop returns.
func (r *Runner) Stop() {
	r.mu.Lock()
	c := r.c
	r.mu.Unlock()

	// c.Stop() returns a context that is done when all in-flight jobs have
	// completed.
	<-c.Stop().Done()
	r.logger.Info("Scheduler stopped")
}
