package api_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetJobsRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/system/jobs", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

// TestGetJobsListsRegisteredJobs confirms the endpoint reflects whatever the
// scheduler currently has registered — including "backup", which
// api.NewRouter registers itself via InitBackupJob — without needing any
// job to have actually run yet (#384).
func TestGetJobsListsRegisteredJobs(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/system/jobs", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var jobs []struct {
		Name      string `json:"name"`
		CronExpr  string `json:"cronExpr"`
		NextRunAt string `json:"nextRunAt"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &jobs); err != nil {
		t.Fatalf("unmarshal: %v; body = %s", err, rec.Body)
	}

	names := map[string]bool{}
	for _, j := range jobs {
		names[j.Name] = true
		if j.CronExpr == "" {
			t.Errorf("job %q has empty cronExpr", j.Name)
		}
	}
	if !names["backup"] {
		t.Fatalf("jobs = %+v, want at least backup registered", jobs)
	}
}
