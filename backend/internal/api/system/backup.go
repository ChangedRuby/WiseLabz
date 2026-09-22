package system

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/httputil"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// MaxImportBytes bounds a backup upload before decoding to prevent a single
// request from consuming unbounded server memory.
const MaxImportBytes = 10 << 20

// ExportBackup handles GET /api/system/backup/export. Operator-only.
func (h *Handler) ExportBackup(w http.ResponseWriter, r *http.Request) {
	b, err := backup.Export(r.Context(), h.Store)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	filename := fmt.Sprintf("wiselabz-backup-%s.json", time.Now().UTC().Format("20060102-150405"))
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	httputil.JSON(w, http.StatusOK, b)
}

// ImportBackup handles POST /api/system/backup/import. Operator-only.
func (h *Handler) ImportBackup(w http.ResponseWriter, r *http.Request) {
	var b backup.Bundle
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, MaxImportBytes))
	if err := decoder.Decode(&b); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httputil.Error(w, http.StatusRequestEntityTooLarge, "request_too_large", "Backup upload exceeds the 10 MiB limit")
			return
		}
		httputil.Error(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
		return
	}

	result, err := backup.Import(r.Context(), h.Store, &b)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid_backup", err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, result)
}

// GetBackupSchedule handles GET /api/system/backup/schedule. Operator-only.
// Returns the current backup schedule configuration.
func (h *Handler) GetBackupSchedule(w http.ResponseWriter, r *http.Request) {
	sched, err := h.Store.GetBackupSchedule(r.Context())
	if err != nil {
		// If schedule doesn't exist (shouldn't happen after init), return a sensible default
		slog.Warn("backup schedule not found, returning default", "error", err)
		sched = store.BackupSchedule{
			CronExpr:    "0 3 * * *",
			MaxBackups:  14,
			MaxAgeHours: 720,
			Enabled:     true,
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"cronExpr":    sched.CronExpr,
		"maxBackups":  sched.MaxBackups,
		"maxAgeHours": sched.MaxAgeHours,
		"enabled":     sched.Enabled,
		"updatedAt":   sched.UpdatedAt,
	})
}

// BackupScheduleRequest is the request body for PUT /api/system/backup/schedule.
type BackupScheduleRequest struct {
	CronExpr    string `json:"cronExpr"`
	MaxBackups  int    `json:"maxBackups"`
	MaxAgeHours int    `json:"maxAgeHours"`
	Enabled     bool   `json:"enabled"`
}

// UpdateBackupSchedule handles PUT /api/system/backup/schedule. Operator-only.
// Updates the backup schedule and re-registers the cron job.
func (h *Handler) UpdateBackupSchedule(w http.ResponseWriter, r *http.Request) {
	req, ok := httputil.DecodeJSON[BackupScheduleRequest](w, r)
	if !ok {
		return
	}

	// Validate cron expression (try both 5-field and 6-field parsers)
	if req.CronExpr == "" {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_cron", "cronExpr must not be empty", []httputil.FieldError{{Field: "cronExpr", Msg: "must not be empty"}})
		return
	}
	parser5Field := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	parser6Field := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err5 := parser5Field.Parse(req.CronExpr)
	_, err6 := parser6Field.Parse(req.CronExpr)
	if err5 != nil && err6 != nil {
		httputil.ErrorWithDetails(w, http.StatusBadRequest, "invalid_cron", "Invalid cron expression: must be valid 5-field or 6-field cron format", []httputil.FieldError{{Field: "cronExpr", Msg: "must be a valid 5-field or 6-field cron expression"}})
		return
	}

	// Update database
	newSched := store.BackupSchedule{
		CronExpr:    req.CronExpr,
		MaxBackups:  req.MaxBackups,
		MaxAgeHours: req.MaxAgeHours,
		Enabled:     req.Enabled,
	}
	if err := h.Store.UpsertBackupSchedule(r.Context(), newSched); err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Audit the change
	if err := h.Store.RecordAuditFromContext(r.Context(), "backup.schedule.update", "backup_schedule", "default", req); err != nil {
		slog.Error("audit backup schedule update", "error", err)
	}

	// Re-register the job with the scheduler
	h.reregisterBackupJob(newSched)

	httputil.JSON(w, http.StatusOK, map[string]any{
		"cronExpr":    newSched.CronExpr,
		"maxBackups":  newSched.MaxBackups,
		"maxAgeHours": newSched.MaxAgeHours,
		"enabled":     newSched.Enabled,
		"updatedAt":   newSched.UpdatedAt,
	})
}

// maxBackupRunsLimit caps the page size for backup run listings.
const maxBackupRunsLimit = 200

// backupRunPage is GET /api/system/backup/runs' published envelope
// (BackupRunPage in docs/openapi.yaml) — limit/offset rather than the
// page/pageSize shape httputil.PaginatedResponse carries.
type backupRunPage struct {
	Runs   any `json:"runs"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ListBackupRuns handles GET /api/system/backup/runs. Operator-only.
// Returns a paginated list of past backups.
func (h *Handler) ListBackupRuns(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	if limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
			limit = n
		}
		if limit > maxBackupRunsLimit {
			limit = maxBackupRunsLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		if n, err := strconv.Atoi(offsetStr); err == nil && n >= 0 {
			offset = n
		}
	}

	runs, total, err := h.Store.ListBackupRuns(r.Context(), limit, offset)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Convert to response format
	type RunResponse struct {
		ID          string `json:"id"`
		TriggeredBy string `json:"triggeredBy"`
		FilePath    string `json:"filePath"`
		SizeBytes   int64  `json:"sizeBytes"`
		CreatedAt   string `json:"createdAt"`
	}
	var items []RunResponse
	for _, run := range runs {
		items = append(items, RunResponse{
			ID:          run.ID,
			TriggeredBy: run.TriggeredBy,
			FilePath:    run.FilePath,
			SizeBytes:   run.SizeBytes,
			CreatedAt:   run.CreatedAt,
		})
	}

	// Not httputil.WritePaginated: the published BackupRunPage schema names the
	// list `runs` and pages with limit/offset rather than page/pageSize, so
	// moving it onto the shared envelope would break existing clients. Kept as
	// a declared type so the shape lives in one place.
	httputil.JSON(w, http.StatusOK, backupRunPage{
		Runs:   items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// CreateBackupRun handles POST /api/system/backup/run. Operator-only.
// Triggers a manual backup now.
func (h *Handler) CreateBackupRun(w http.ResponseWriter, r *http.Request) {
	// Export to file
	run, err := backup.ExportToFile(r.Context(), h.Store, h.BackupDir)
	if err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Mark it as manual
	run.TriggeredBy = "manual"

	// Store the run in the database
	dbRun := store.BackupRun{
		ID:          run.ID,
		TriggeredBy: run.TriggeredBy,
		FilePath:    run.FilePath,
		SizeBytes:   run.SizeBytes,
		CreatedAt:   run.CreatedAt,
	}
	if err := h.Store.CreateBackupRun(r.Context(), dbRun); err != nil {
		httputil.Errorf(w, err)
		return
	}

	// Audit the manual backup
	if err := h.Store.RecordAuditFromContext(r.Context(), "backup.run.created", "backup_run", run.ID, map[string]any{
		"triggeredBy": "manual",
		"sizeBytes":   run.SizeBytes,
	}); err != nil {
		slog.Error("audit backup run creation", "error", err)
	}

	// Apply pruning (same as scheduled backups)
	h.pruneBackups(r.Context())

	httputil.JSON(w, http.StatusOK, map[string]any{
		"id":          run.ID,
		"triggeredBy": run.TriggeredBy,
		"filePath":    run.FilePath,
		"sizeBytes":   run.SizeBytes,
		"createdAt":   run.CreatedAt,
	})
}

// pruneBackups applies the current backup retention policy, deleting old/excess files.
// Errors are logged but don't fail the request.
func (h *Handler) pruneBackups(ctx context.Context) {
	sched, err := h.Store.GetBackupSchedule(ctx)
	if err != nil {
		slog.Error("get backup schedule for pruning", "error", err)
		return
	}

	// Compute cutoff time
	cutoffTime := time.Now().UTC().Add(time.Duration(-sched.MaxAgeHours) * time.Hour)
	cutoff := cutoffTime.Format(time.RFC3339)

	// Find and delete old runs
	pruned, err := h.Store.PruneBackupRuns(ctx, sched.MaxBackups, cutoff)
	if err != nil {
		slog.Error("prune backup runs", "error", err)
		return
	}

	// Delete the actual files
	for _, run := range pruned {
		if err := os.Remove(run.FilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			slog.Error("delete backup file", "path", run.FilePath, "error", err)
		}
	}

	// Audit if anything was pruned
	if len(pruned) > 0 {
		if err := h.Store.RecordAuditFromContext(ctx, "backup.run.pruned", "backup_schedule", "default", map[string]any{
			"prunedCount": len(pruned),
		}); err != nil {
			slog.Error("audit backup pruning", "error", err)
		}
	}
}

// InitBackupJob registers the backup cron job at startup from the persisted
// schedule, routing through the same BackupJobID bookkeeping as
// reregisterBackupJob uses. This must be the only place that ever calls
// h.Scheduler.AddJob("backup", ...) directly — if the initial registration
// happens anywhere else (e.g. main.go calling the scheduler directly), the
// Handler never learns that job's entry ID, so the first PUT /schedule call
// can't remove it and a duplicate job accumulates instead of replacing it.
// Seeds the schedule from config defaults if no row exists yet (first boot).
func (h *Handler) InitBackupJob(ctx context.Context) {
	sched, err := h.Store.GetBackupSchedule(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("failed to read backup schedule", "error", err)
			return
		}
		slog.Info("initializing backup schedule with defaults")
		sched = store.BackupSchedule{
			CronExpr:    h.Config.Backup.CronExpr,
			MaxBackups:  h.Config.Backup.MaxBackups,
			MaxAgeHours: h.Config.Backup.MaxAgeHours,
			Enabled:     h.Config.Backup.Enabled,
		}
		if err := h.Store.UpsertBackupSchedule(ctx, sched); err != nil {
			slog.Error("failed to initialize backup schedule", "error", err)
			return
		}
	}
	h.reregisterBackupJob(sched)
}

// reregisterBackupJob removes the old job (if any) and registers a new one with the updated schedule.
func (h *Handler) reregisterBackupJob(sched store.BackupSchedule) {
	if h.Scheduler == nil {
		return
	}

	h.BackupJobIDMu.Lock()
	defer h.BackupJobIDMu.Unlock()

	// Remove the old job if one exists
	if h.BackupJobID != 0 {
		h.Scheduler.RemoveJob(h.BackupJobID)
	}

	// Register the new job (if enabled)
	if sched.Enabled {
		id, err := h.Scheduler.AddJob("backup", sched.CronExpr, func(jobCtx context.Context) {
			h.runBackupJob(jobCtx, sched)
		})
		if err != nil {
			// err wraps sched.CronExpr (user-controlled via PUT /schedule); strip
			// line breaks before logging so a crafted cron string can't forge
			// additional log entries (CWE-117).
			slog.Error("register backup job", "error", stripLogControlChars(err.Error()))
			h.BackupJobID = 0
		} else {
			h.BackupJobID = id
		}
	} else {
		h.BackupJobID = 0
	}
}

// runBackupJob is the actual backup job that runs on schedule.
func (h *Handler) runBackupJob(ctx context.Context, sched store.BackupSchedule) {
	run, err := backup.ExportToFile(ctx, h.Store, h.BackupDir)
	if err != nil {
		slog.Error("backup export to file", "error", err)
		return
	}

	run.TriggeredBy = "schedule"

	// Store in database
	dbRun := store.BackupRun{
		ID:          run.ID,
		TriggeredBy: run.TriggeredBy,
		FilePath:    run.FilePath,
		SizeBytes:   run.SizeBytes,
		CreatedAt:   run.CreatedAt,
	}
	if err := h.Store.CreateBackupRun(ctx, dbRun); err != nil {
		slog.Error("create backup run", "error", err)
		return
	}

	// Apply pruning
	cutoffTime := time.Now().UTC().Add(time.Duration(-sched.MaxAgeHours) * time.Hour)
	cutoff := cutoffTime.Format(time.RFC3339)
	pruned, err := h.Store.PruneBackupRuns(ctx, sched.MaxBackups, cutoff)
	if err != nil {
		slog.Error("prune backup runs", "error", err)
	} else {
		for _, prun := range pruned {
			if err := os.Remove(prun.FilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
				slog.Error("delete backup file", "path", prun.FilePath, "error", err)
			}
		}
	}

	slog.Info("Backup created", "id", run.ID, "size", run.SizeBytes)
}

// stripLogControlChars removes CR/LF from a string before it's logged, so a
// value that embeds user input (e.g. an error wrapping a submitted cron
// expression) can't forge additional log entries.
func stripLogControlChars(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	return strings.ReplaceAll(s, "\r", "")
}
