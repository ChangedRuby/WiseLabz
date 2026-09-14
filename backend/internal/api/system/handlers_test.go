package system

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/config"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	s := apitest.NewStore(t)
	return NewHandler(s.DB(), &config.Config{}, s, nil, t.TempDir())
}

func TestHealth(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	h.Health(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestInfo(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/info", nil)
	rr := httptest.NewRecorder()
	h.Info(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestDiagnostics(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/diagnostics", nil)
	rr := httptest.NewRecorder()
	h.Diagnostics(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestListAudit(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/audit", nil)
	rr := httptest.NewRecorder()
	h.ListAudit(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestExportAudit(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/audit/export?format=xml", nil)
		rr := httptest.NewRecorder()
		h.ExportAudit(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("csv default", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/audit/export?format=csv", nil)
		rr := httptest.NewRecorder()
		h.ExportAudit(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
		if !strings.Contains(rr.Body.String(), "id,actorUserId") {
			t.Errorf("expected CSV header, got: %s", rr.Body.String())
		}
	})
}

func TestGetRetentionSettingsDefault(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/settings/retention", nil)
	rr := httptest.NewRecorder()
	h.GetRetentionSettings(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestUpdateRetentionSettings(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/settings/retention", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.UpdateRetentionSettings(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("empty cron", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/settings/retention", strings.NewReader(`{"cronExpr":""}`))
		rr := httptest.NewRecorder()
		h.UpdateRetentionSettings(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid cron", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/settings/retention", strings.NewReader(`{"cronExpr":"not a cron"}`))
		rr := httptest.NewRecorder()
		h.UpdateRetentionSettings(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("negative retention days", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/settings/retention", strings.NewReader(`{"cronExpr":"0 0 * * *","snapshotDays":-1}`))
		rr := httptest.NewRecorder()
		h.UpdateRetentionSettings(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/settings/retention", strings.NewReader(`{"cronExpr":"0 0 * * *","snapshotDays":30,"docVersionDays":60,"alertDays":90,"syncRunDays":30,"auditDays":90}`))
		rr := httptest.NewRecorder()
		h.UpdateRetentionSettings(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})
}

func TestGetBackupScheduleDefault(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/backup/schedule", nil)
	rr := httptest.NewRecorder()
	h.GetBackupSchedule(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}

func TestUpdateBackupSchedule(t *testing.T) {
	h := newTestHandler(t)

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/backup/schedule", strings.NewReader(`{`))
		rr := httptest.NewRecorder()
		h.UpdateBackupSchedule(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("empty cron", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/backup/schedule", strings.NewReader(`{"cronExpr":""}`))
		rr := httptest.NewRecorder()
		h.UpdateBackupSchedule(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid cron", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/backup/schedule", strings.NewReader(`{"cronExpr":"garbage"}`))
		rr := httptest.NewRecorder()
		h.UpdateBackupSchedule(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/system/backup/schedule", strings.NewReader(`{"cronExpr":"0 3 * * *","maxBackups":7,"maxAgeHours":168,"enabled":true}`))
		rr := httptest.NewRecorder()
		h.UpdateBackupSchedule(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
		}
	})
}

func TestExportImportBackupRoundTrip(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/system/backup/export", nil)
	rr := httptest.NewRecorder()
	h.ExportBackup(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("ExportBackup() status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	t.Run("import invalid json", func(t *testing.T) {
		importReq := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", strings.NewReader(`{`))
		importRR := httptest.NewRecorder()
		h.ImportBackup(importRR, importReq)
		if importRR.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", importRR.Code, http.StatusBadRequest)
		}
	})

	t.Run("import the export back in", func(t *testing.T) {
		importReq := httptest.NewRequest(http.MethodPost, "/api/system/backup/import", strings.NewReader(rr.Body.String()))
		importRR := httptest.NewRecorder()
		h.ImportBackup(importRR, importReq)
		if importRR.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body=%s", importRR.Code, http.StatusOK, importRR.Body.String())
		}
	})
}

func TestListBackupRuns(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/backup/runs", nil)
	rr := httptest.NewRecorder()
	h.ListBackupRuns(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}
}
