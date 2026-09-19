// Package scheduler provides a shared cron-based scheduling primitive for
// running background jobs at configured intervals.
package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/WiseLabz/wiselabz/internal/logsafe"
)

// Runner manages a set of cron-scheduled jobs.
type Runner struct {
	c      *cron.Cron
	logger *slog.Logger
	mu     sync.Mutex // protects access to the cron instance and ctx
	ctx    context.Context
}

// New creates a new Runner with support for both standard 5-field cron
// expressions and 6-field expressions with a leading seconds field (needed
// for sub-minute cadence), matching config.validateCronExpressions.
func New(logger *slog.Logger) *Runner {
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)
	return &Runner{
		c:      cron.New(cron.WithParser(parser), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))),
		logger: logger,
		ctx:    context.Background(),
	}
}

// AddJob registers a new cron job with the given name and expression.
// The function fn is wrapped to recover panics and log entry/exit.
// Returns the cron entry ID for later removal, or an error if the cron
// expression is invalid.
func (r *Runner) AddJob(name, cronExpr string, fn func(ctx context.Context)) (cron.EntryID, error) {
	// Wrap the function to recover panics and log entry/exit
	wrapped := func() {
		r.logger.Debug("job started", "job", name)
		defer func() {
			if rec := recover(); rec != nil {
				r.logger.Error("job panicked", "job", name, "panic", rec)
			}
		}()

		// Use the context passed to Start so job bodies observe cancellation
		// on server shutdown, the same as the old ticker loops did.
		fn(r.jobContext())
		r.logger.Debug("job completed", "job", name)
	}

	// Use AddFunc which takes the cron expression directly and validates it
	entryID, err := r.c.AddFunc(cronExpr, wrapped)
	if err != nil {
		return 0, fmt.Errorf("add job %q (cron %q): %w", name, cronExpr, err)
	}

	r.logger.Debug("job registered", "job", name, "cron", logsafe.Sanitize(cronExpr))
	return entryID, nil
}

// RemoveJob removes a previously registered job by its entry ID.
func (r *Runner) RemoveJob(id cron.EntryID) {
	r.c.Remove(id)
	r.logger.Debug("job removed", "entry_id", id)
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
