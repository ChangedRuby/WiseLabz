package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRetentionSettingsRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/system/settings/retention", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("GET status = %d, want 403; body = %s", rec.Code, rec.Body)
	}

	rec = app.req(t, http.MethodPut, "/api/system/settings/retention", map[string]any{"cronExpr": "0 0 * * *"}, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("PUT status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestRetentionSettingsGetPutRoundTrip(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	// GET before any PUT returns the config-seeded defaults (fallback path,
	// since newTestApp doesn't seed a retention_settings row).
	getRec := app.req(t, http.MethodGet, "/api/system/settings/retention", nil, opToken)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want 200; body = %s", getRec.Code, getRec.Body)
	}

	putBody := map[string]any{
		"snapshotDays":    30,
		"docVersionDays":  60,
		"alertDays":       90,
		"syncRunDays":     14,
		"auditDays":       120,
		"healthCheckDays": 45,
		"cronExpr":        "0 1 * * *",
	}
	putRec := app.req(t, http.MethodPut, "/api/system/settings/retention", putBody, opToken)
	if putRec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want 200; body = %s", putRec.Code, putRec.Body)
	}
	var putResp map[string]any
	if err := json.Unmarshal(putRec.Body.Bytes(), &putResp); err != nil {
		t.Fatalf("decode PUT response: %v", err)
	}
	if putResp["auditDays"].(float64) != 120 {
		t.Errorf("PUT response auditDays = %v, want 120", putResp["auditDays"])
	}

	getRec2 := app.req(t, http.MethodGet, "/api/system/settings/retention", nil, opToken)
	if getRec2.Code != http.StatusOK {
		t.Fatalf("GET (after PUT) status = %d, want 200; body = %s", getRec2.Code, getRec2.Body)
	}
	var getResp map[string]any
	if err := json.Unmarshal(getRec2.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if getResp["cronExpr"] != "0 1 * * *" {
		t.Errorf("GET cronExpr = %v, want %q", getResp["cronExpr"], "0 1 * * *")
	}
	if getResp["syncRunDays"].(float64) != 14 {
		t.Errorf("GET syncRunDays = %v, want 14", getResp["syncRunDays"])
	}
	if getResp["healthCheckDays"].(float64) != 45 {
		t.Errorf("GET healthCheckDays = %v, want 45", getResp["healthCheckDays"])
	}
}

func TestRetentionSettingsValidation(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	negRec := app.req(t, http.MethodPut, "/api/system/settings/retention", map[string]any{
		"snapshotDays": -1, "cronExpr": "0 0 * * *",
	}, opToken)
	if negRec.Code != http.StatusBadRequest {
		t.Fatalf("negative days status = %d, want 400; body = %s", negRec.Code, negRec.Body)
	}

	badCronRec := app.req(t, http.MethodPut, "/api/system/settings/retention", map[string]any{
		"cronExpr": "not a cron expression",
	}, opToken)
	if badCronRec.Code != http.StatusBadRequest {
		t.Fatalf("bad cron status = %d, want 400; body = %s", badCronRec.Code, badCronRec.Body)
	}

	emptyCronRec := app.req(t, http.MethodPut, "/api/system/settings/retention", map[string]any{
		"cronExpr": "",
	}, opToken)
	if emptyCronRec.Code != http.StatusBadRequest {
		t.Fatalf("empty cron status = %d, want 400; body = %s", emptyCronRec.Code, emptyCronRec.Body)
	}
}

func TestRetentionSettingsUpdateRecordsAudit(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/system/settings/retention", map[string]any{
		"snapshotDays": 10, "cronExpr": "0 2 * * *",
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	records, total, err := app.Store.ListAuditRecords(context.Background(), "retention.settings.update", "retention_settings", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if total != 1 {
		t.Fatalf("audit records for retention.settings.update = %d, want 1", total)
	}
	if records[0].TargetID != "default" {
		t.Errorf("audit record TargetID = %q, want %q", records[0].TargetID, "default")
	}
}
