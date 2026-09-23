// Package docexport writes every generated doc to a local directory as
// Markdown on a schedule, so documentation survives outside the tool (and
// doubles as a lightweight secondary backup). It reuses the content the doc
// engine already renders and persists (store.DocRecord.Content) rather than
// re-rendering anything itself.
//
// Optionally (see ConfigureGit) the directory is a persistent clone of a Git
// remote: each run fetches, hard-resets to the remote branch, re-exports,
// commits the difference and pushes it (never forcing).
package docexport

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/WiseLabz/wiselabz/internal/store"
)

// DefaultCronExpr is the cron expression the scheduled export job registers
// with by default (daily, off-hours). Operators can override it via
// config.DocExportSettings.CronExpr.
const DefaultCronExpr = "0 2 * * *"

// exportPageSize matches internal/backup's paginated doc export page size.
const exportPageSize = 1000

// Exporter renders every doc in the store to Markdown files in a target
// directory.
type Exporter struct {
	store *store.Store

	git      *gitTarget // nil: local-directory mode
	notifier Notifier   // nil: failures are only logged

	mu      sync.Mutex // serializes runs (the Git worktree isn't concurrency-safe)
	failing bool       // last run failed; in memory only, so a restart re-notifies
}

// NewExporter creates a new Exporter backed by s.
func NewExporter(s *store.Store) *Exporter {
	return &Exporter{store: s}
}

// Result summarizes a completed export run.
type Result struct {
	Dir     string   `json:"dir"`
	Count   int      `json:"count"`
	Files   []string `json:"files"`
	Removed []string `json:"removed,omitempty"`
}

// ExportAll fetches every doc from the store and writes each as a Markdown
// file under dir, creating dir if needed. Filenames are derived from the
// doc's title (slugified) plus a short ID suffix to keep them unique and
// stable across runs even if two docs share a title. Any *.md file left
// over in dir from a previous run whose doc was renamed or deleted is
// removed, so dir always mirrors the current set of docs.
func (e *Exporter) ExportAll(ctx context.Context, dir string) (Result, error) {
	if dir == "" {
		return Result{}, fmt.Errorf("export directory must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Result{}, fmt.Errorf("create export directory: %w", err)
	}

	docs, err := fetchAllDocs(ctx, e.store)
	if err != nil {
		return Result{}, fmt.Errorf("fetch docs: %w", err)
	}

	files := make([]string, 0, len(docs))
	seen := make(map[string]bool, len(docs))
	for _, d := range docs {
		name := fileName(d)
		if seen[name] {
			// Extremely unlikely (would require a colliding slug+ID
			// prefix); disambiguate rather than overwrite silently.
			name = d.ID + ".md"
		}
		seen[name] = true

		if err := os.WriteFile(filepath.Join(dir, name), []byte(d.Content), 0o644); err != nil {
			return Result{}, fmt.Errorf("write doc %q: %w", d.ID, err)
		}
		files = append(files, name)
	}

	removed, err := pruneStale(dir, seen)
	if err != nil {
		return Result{}, fmt.Errorf("prune stale files: %w", err)
	}

	return Result{Dir: dir, Count: len(files), Files: files, Removed: removed}, nil
}

// fetchAllDocs pages through every doc in the store, mirroring
// internal/backup.exportDocs.
func fetchAllDocs(ctx context.Context, s *store.Store) ([]store.DocRecord, error) {
	docs := []store.DocRecord{}
	for offset := 0; ; offset += exportPageSize {
		page, total, err := s.ListAllDocsWithContent(ctx, "", offset, exportPageSize)
		if err != nil {
			return nil, err
		}
		if len(page) == 0 {
			return docs, nil
		}
		docs = append(docs, page...)
		if len(docs) >= total {
			return docs, nil
		}
	}
}

// generatedName matches the filenames fileName (and its collision fallback,
// a bare UUID) produces. Only these are ever pruned, so anything else an
// operator keeps next to the export (README.md, notes.md, …) survives.
var generatedName = regexp.MustCompile(`^(?:[a-z0-9-]+-[0-9a-f]{8}|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})\.md$`)

// IsGeneratedName reports whether name looks like a file ExportAll writes,
// i.e. one pruneStale is allowed to delete.
func IsGeneratedName(name string) bool { return generatedName.MatchString(name) }

// pruneStale removes any top-level generated doc file in dir that isn't in
// keep, leaving other files (e.g. a README the operator dropped in) untouched.
func pruneStale(dir string, keep map[string]bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var removed []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !IsGeneratedName(name) || keep[name] {
			continue
		}
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			return removed, fmt.Errorf("remove stale file %q: %w", name, err)
		}
		removed = append(removed, name)
	}
	return removed, nil
}

// fileName derives a stable, filesystem-safe Markdown filename for a doc:
// a slug of its title, plus the first 8 characters of its ID so renamed
// titles don't collide and unrelated docs with the same title don't
// overwrite each other.
func fileName(d store.DocRecord) string {
	slug := slugify(d.Title)
	if slug == "" {
		slug = "untitled"
	}
	idSuffix := d.ID
	if len(idSuffix) > 8 {
		idSuffix = idSuffix[:8]
	}
	return fmt.Sprintf("%s-%s.md", slug, idSuffix)
}

// slugify lowercases s and replaces every run of characters that aren't
// ASCII letters, digits, or '-' with a single '-', trimming leading and
// trailing separators.
func slugify(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// RunExportOnce runs a single export pass and logs the outcome. It's the
// function the scheduled "docexport" cron job (wired in cmd/server/main.go)
// calls, following the same log-on-failure convention as other scheduled
// jobs (e.g. internal/backup.RunVerifyOnce). In Git mode dir is the
// persistent clone and the docs go to its configured subdirectory. When a
// Notifier is set, a system.job_failed event is sent on the ok→failing and
// failing→ok transitions only.
func RunExportOnce(ctx context.Context, e *Exporter, dir string, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	var err error
	if e.git != nil {
		err = e.runGit(ctx, dir, logger)
	} else {
		var res Result
		res, err = e.ExportAll(ctx, dir)
		if err == nil {
			logger.Info("doc export: completed", "dir", res.Dir, "count", res.Count, "removed", len(res.Removed))
		}
	}

	if err != nil {
		logger.Error("doc export: failed", "dir", dir, "error", err)
		if !e.failing && e.notifier != nil {
			e.notifier.NotifySystemEvent(ctx, EventJobFailed, "warning",
				"Doc export failing", fmt.Sprintf("The scheduled doc export job failed and will retry on its next run: %v", err))
		}
		e.failing = true
		return
	}
	if e.failing && e.notifier != nil {
		e.notifier.NotifySystemEvent(ctx, EventJobFailed, "info",
			"Doc export recovered", "The scheduled doc export job succeeded again.")
	}
	e.failing = false
}

// EventJobFailed is the notification event type sent when a scheduled job
// starts failing (severity warning) or recovers (severity info). It must
// match notifications.EventSystemJobFailed.
const EventJobFailed = "system.job_failed"

// Notifier is the slice of notifications.Dispatcher the exporter uses.
type Notifier interface {
	NotifySystemEvent(ctx context.Context, eventType, severity, title, message string)
}

// SetNotifier installs the notifier used for failure/recovery events.
func (e *Exporter) SetNotifier(n Notifier) { e.notifier = n }
