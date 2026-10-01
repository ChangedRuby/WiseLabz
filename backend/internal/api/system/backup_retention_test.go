package system

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/backup"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// seedBackupRuns stores n runs (oldest first, each ageDays older than the
// next) with a bundle and manifest file on disk, and returns the bundle paths.
func seedBackupRuns(t *testing.T, h *Handler, n int) []string {
	t.Helper()
	ctx := context.Background()
	var paths []string
	for i := 0; i < n; i++ {
		path := filepath.Join(h.BackupDir, fmt.Sprintf("bundle-%d.json", i))
		for _, p := range []string{path, backup.ManifestPath(path)} {
			if err := os.WriteFile(p, []byte("{}"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		created := time.Now().UTC().Add(-time.Duration(n-i) * 48 * time.Hour).Format(time.RFC3339)
		run := store.BackupRun{ID: fmt.Sprintf("run-%d", i), TriggeredBy: "schedule", FilePath: path, SizeBytes: 2, CreatedAt: created}
		if err := h.Store.CreateBackupRun(ctx, run); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}
	return paths
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }

func TestPruneBackupsAgeLimit(t *testing.T) {
	for _, tc := range []struct {
		name        string
		maxBackups  int
		maxAgeHours int
		wantKept    int // newest N remain
	}{
		{"zero age limit keeps everything", 0, 0, 5},
		{"stored negative age limit keeps everything", 0, -5, 5},
		{"zero age limit still applies count limit", 3, 0, 3},
		{"positive age limit prunes old backups", 0, 72, 1}, // run ages: 240,192,144,96,48h
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t)
			paths := seedBackupRuns(t, h, 5)
			sched := store.BackupSchedule{CronExpr: "0 3 * * *", MaxBackups: tc.maxBackups, MaxAgeHours: tc.maxAgeHours, Enabled: true}
			if err := h.Store.UpsertBackupSchedule(context.Background(), sched); err != nil {
				t.Fatal(err)
			}

			h.pruneBackups(context.Background())

			for i, p := range paths {
				keep := i >= len(paths)-tc.wantKept
				if exists(p) != keep || exists(backup.ManifestPath(p)) != keep {
					t.Errorf("run %d: bundle exists=%v manifest exists=%v, want both %v", i, exists(p), exists(backup.ManifestPath(p)), keep)
				}
			}
		})
	}
}

func TestUpdateBackupScheduleRejectsNegativeMaxAge(t *testing.T) {
	h := newTestHandler(t)
	ctx := context.Background()
	orig := store.BackupSchedule{CronExpr: "0 3 * * *", MaxBackups: 7, MaxAgeHours: 168, Enabled: true}
	if err := h.Store.UpsertBackupSchedule(ctx, orig); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/system/backup/schedule", strings.NewReader(`{"cronExpr":"0 4 * * *","maxBackups":7,"maxAgeHours":-1,"enabled":true}`))
	rr := httptest.NewRecorder()
	h.UpdateBackupSchedule(rr, req)

	if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "maxAgeHours") {
		t.Fatalf("status=%d body=%s, want 400 naming maxAgeHours", rr.Code, rr.Body.String())
	}
	got, err := h.Store.GetBackupSchedule(ctx)
	if err != nil || got.CronExpr != orig.CronExpr || got.MaxAgeHours != orig.MaxAgeHours {
		t.Fatalf("schedule changed: %+v err=%v", got, err)
	}
}
