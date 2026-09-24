package report

import (
	"context"
	"log/slog"
	"strings"

	"github.com/WiseLabz/wiselabz/internal/store"
	"github.com/robfig/cron/v3"
)

// Scheduler is the small portion of scheduler.Runner report schedules need.
type Scheduler interface {
	AddJob(string, string, func(context.Context) error) (cron.EntryID, error)
	RemoveJob(cron.EntryID)
}

// Manager owns registration and replacement of report cron jobs.
type Manager struct {
	Store     *store.Store
	Generator *Generator
	Scheduler Scheduler
	Notify    func(context.Context, store.ReportRecord, store.ReportDefinitionRecord)
	entries   map[string]cron.EntryID
}

// NewManager creates the owner for report schedules and their generated snapshots.
func NewManager(s *store.Store, g *Generator, scheduler Scheduler, notify func(context.Context, store.ReportRecord, store.ReportDefinitionRecord)) *Manager {
	return &Manager{Store: s, Generator: g, Scheduler: scheduler, Notify: notify, entries: map[string]cron.EntryID{}}
}

// Init registers persisted schedules and generates an initial report when needed.
func (m *Manager) Init(ctx context.Context) error {
	defs, err := m.Store.ListReportDefinitions(ctx)
	if err != nil {
		return err
	}
	for _, d := range defs {
		if err = m.Register(d); err != nil {
			return err
		}
		if d.Enabled {
			if _, e := m.Store.LatestScheduledReport(ctx, d.ID); e == store.ErrNotFound {
				r, e := m.Generator.Generate(ctx, d, "scheduled")
				if m.Notify != nil {
					m.Notify(ctx, r, d)
				}
				LogPartial(e)
			}
		}
	}
	return nil
}

// Register replaces the scheduled job for a report definition.
func (m *Manager) Register(d store.ReportDefinitionRecord) error {
	m.Unregister(d.Slug)
	if !d.Enabled {
		return nil
	}
	id, err := m.Scheduler.AddJob("report:"+d.Slug, "CRON_TZ="+d.Timezone+" "+d.CronExpr, func(ctx context.Context) error {
		r, e := m.Generator.Generate(ctx, d, "scheduled")
		if m.Notify != nil {
			m.Notify(ctx, r, d)
		}
		return e
	})
	if err == nil {
		m.entries[d.Slug] = id
	}
	return err
}

// Unregister removes a report schedule by slug.
func (m *Manager) Unregister(slug string) {
	if id, ok := m.entries[slug]; ok {
		m.Scheduler.RemoveJob(id)
		delete(m.entries, slug)
	}
}

// Run generates and notifies for a report on demand.
func (m *Manager) Run(ctx context.Context, d store.ReportDefinitionRecord) (store.ReportRecord, error) {
	r, e := m.Generator.Generate(ctx, d, "manual")
	if m.Notify != nil {
		m.Notify(ctx, r, d)
	}
	return r, e
}

// JobName returns the scheduler's stable name for a report slug.
func JobName(slug string) string { return "report:" + strings.TrimSpace(slug) }

// LogPartial records section-level generation errors without discarding the report.
func LogPartial(err error) {
	if err != nil {
		slog.Error("report generated partially", "error", err)
	}
}
