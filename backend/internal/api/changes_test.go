package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedChange(t *testing.T, app *testApp) *store.ChangeRecord {
	t.Helper()
	return seedChangeWithSeverity(t, app, "info")
}

func seedChangeWithSeverity(t *testing.T, app *testApp, severity string) *store.ChangeRecord {
	t.Helper()
	c := &store.ChangeRecord{
		ServiceID: "svc-1", ChangeType: "config", Severity: severity,
		Summary: "something changed", Diff: "[]", Status: "new", AffectedDocIDs: "[]",
	}
	if err := app.Store.CreateChange(context.Background(), c); err != nil {
		t.Fatalf("seed change: %v", err)
	}
	return c
}

func TestChangesListSuccess(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	seedChange(t, app)

	rec := app.req(t, http.MethodGet, "/api/changes", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestChangesAcknowledgeRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	c := seedChange(t, app)

	rec := app.req(t, http.MethodPost, "/api/changes/"+c.ID+"/ack", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestChangesAcknowledgeSuccess(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	c := seedChange(t, app)
	app.connectorGrant(t, opUserID, c.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/changes/"+c.ID+"/ack", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	got, err := app.Store.GetChange(context.Background(), c.ID)
	if err != nil {
		t.Fatalf("GetChange: %v", err)
	}
	if got.Status != "acknowledged" {
		t.Errorf("Status = %q, want acknowledged", got.Status)
	}
}

func TestChangesDismissNotFound(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/changes/does-not-exist/dismiss", nil, opToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

// TestChangesBulkResolveRoleBoundary verifies a caller without an operator
// grant on the change's connector gets a per-item "forbidden" outcome, not a
// blanket 403 — see the matching TestAlertsBulkSnoozeRoleBoundary comment.
func TestChangesBulkResolveRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	c := seedChange(t, app)

	rec := app.req(t, http.MethodPost, "/api/changes/bulk-resolve",
		map[string]any{"ids": []string{c.ID}, "status": "acknowledged"}, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Results []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Reason string `json:"reason"`
		} `json:"results"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(body.Results) != 1 || body.Results[0].Status != "error" || body.Results[0].Reason != "forbidden" {
		t.Fatalf("results = %+v, want one error/forbidden", body.Results)
	}
}

// TestChangesBulkResolvePartialFailure is the core acceptance-criteria test:
// a mixed batch (valid low-risk, a critical id the server must reject on its
// own re-check, and a nonexistent id) must not abort — each item gets its own
// outcome, and only the successful item is acknowledged and audited.
func TestChangesBulkResolvePartialFailure(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")

	lowRisk := seedChangeWithSeverity(t, app, "info")
	critical := seedChangeWithSeverity(t, app, "critical")
	app.connectorGrant(t, opUserID, lowRisk.ServiceID, "operator")
	app.connectorGrant(t, opUserID, critical.ServiceID, "operator")
	const missingID = "does-not-exist"

	rec := app.req(t, http.MethodPost, "/api/changes/bulk-resolve", map[string]any{
		"ids":    []string{lowRisk.ID, critical.ID, missingID},
		"status": "acknowledged",
	}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Results []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Reason string `json:"reason"`
		} `json:"results"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(body.Results) != 3 {
		t.Fatalf("got %d results, want 3: %+v", len(body.Results), body.Results)
	}

	outcomes := map[string]struct{ status, reason string }{}
	for _, r := range body.Results {
		outcomes[r.ID] = struct{ status, reason string }{r.Status, r.Reason}
	}

	if o := outcomes[lowRisk.ID]; o.status != "success" {
		t.Errorf("lowRisk outcome = %+v, want success", o)
	}
	if o := outcomes[critical.ID]; o.status != "error" || o.reason != "not_low_risk" {
		t.Errorf("critical outcome = %+v, want error/not_low_risk", o)
	}
	if o := outcomes[missingID]; o.status != "error" || o.reason != "not_found" {
		t.Errorf("missing outcome = %+v, want error/not_found", o)
	}

	// The valid item actually got acknowledged; the rejected critical item's
	// status must be untouched.
	got, err := app.Store.GetChange(context.Background(), lowRisk.ID)
	if err != nil {
		t.Fatalf("GetChange lowRisk: %v", err)
	}
	if got.Status != "acknowledged" {
		t.Errorf("lowRisk.Status = %q, want acknowledged", got.Status)
	}
	gotCritical, err := app.Store.GetChange(context.Background(), critical.ID)
	if err != nil {
		t.Fatalf("GetChange critical: %v", err)
	}
	if gotCritical.Status != "new" {
		t.Errorf("critical.Status = %q, want unchanged (new)", gotCritical.Status)
	}

	// Audit: exactly one record for the successful item, none for the
	// rejected ones.
	auditRec := app.req(t, http.MethodGet, "/api/system/audit?action=change.bulk_ack", nil, opToken)
	if auditRec.Code != http.StatusOK {
		t.Fatalf("audit list status = %d, want 200; body = %s", auditRec.Code, auditRec.Body)
	}
	var page struct {
		Items []store.AuditRecord `json:"items"`
		Total int                 `json:"total"`
	}
	if err := json.Unmarshal(auditRec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode audit body: %v", err)
	}
	if page.Total != 1 {
		t.Fatalf("audit Total = %d, want 1: %+v", page.Total, page.Items)
	}
	if page.Items[0].TargetID != lowRisk.ID {
		t.Errorf("audit TargetID = %q, want %q", page.Items[0].TargetID, lowRisk.ID)
	}
	if page.Items[0].ActorUserID != opUserID {
		t.Errorf("audit ActorUserID = %q, want %q", page.Items[0].ActorUserID, opUserID)
	}
}

func TestChangesBulkResolveInvalidStatus(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	c := seedChange(t, app)

	rec := app.req(t, http.MethodPost, "/api/changes/bulk-resolve",
		map[string]any{"ids": []string{c.ID}, "status": "deleted"}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestChangesBulkResolveEmptyIDs(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/changes/bulk-resolve",
		map[string]any{"ids": []string{}, "status": "acknowledged"}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestChangesBulkResolveRejectsTooManyIDs(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	ids := make([]string, 501)

	rec := app.req(t, http.MethodPost, "/api/changes/bulk-resolve", map[string]any{
		"ids": ids, "status": "acknowledged",
	}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}
