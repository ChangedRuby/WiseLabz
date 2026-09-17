package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedAlert(t *testing.T, app *testApp) *store.AlertRecord {
	t.Helper()
	a := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "info", Title: "alert", Description: "desc", Status: "pending",
	}
	if err := app.Store.CreateAlert(context.Background(), a); err != nil {
		t.Fatalf("seed alert: %v", err)
	}
	return a
}

func TestAlertsListSuccess(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	seedAlert(t, app)

	rec := app.req(t, http.MethodGet, "/api/alerts", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestAlertsResolveRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	a := seedAlert(t, app)

	rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/resolve", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestAlertsResolveSuccess(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/resolve", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	got, err := app.Store.GetAlert(context.Background(), a.ID)
	if err != nil {
		t.Fatalf("GetAlert: %v", err)
	}
	if got.Status != "resolved" {
		t.Errorf("Status = %q, want resolved", got.Status)
	}
}

func TestAlertsSnoozeValidation(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	tests := []struct {
		name string
		body any
	}{
		{"missing until", map[string]any{}},
		{"non-RFC3339 until", map[string]any{"until": "tomorrow"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/snooze", tt.body, opToken)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestAlertsSnoozeSuccess(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/snooze", map[string]any{"until": "2099-01-01T00:00:00Z"}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

// TestAlertsBulkSnoozeRoleBoundary verifies a caller without an operator
// grant on the alert's connector gets a per-item "forbidden" outcome, not a
// blanket 403 — bulk endpoints report per-ID authorization since a batch can
// span connectors the caller has different access to on each.
func TestAlertsBulkSnoozeRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	a := seedAlert(t, app)

	rec := app.req(t, http.MethodPost, "/api/alerts/bulk-snooze",
		map[string]any{"ids": []string{a.ID}, "until": "2099-01-01T00:00:00Z"}, viewerToken)
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

func TestAlertsBulkSnoozeValidation(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	tests := []struct {
		name string
		body any
	}{
		{"missing until", map[string]any{"ids": []string{a.ID}}},
		{"non-RFC3339 until", map[string]any{"ids": []string{a.ID}, "until": "tomorrow"}},
		{"empty ids", map[string]any{"ids": []string{}, "until": "2099-01-01T00:00:00Z"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/alerts/bulk-snooze", tt.body, opToken)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestAlertsBulkSnoozeRejectsTooManyIDs(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	ids := make([]string, 501) // no connector grant needed: rejected before authorization

	rec := app.req(t, http.MethodPost, "/api/alerts/bulk-snooze", map[string]any{
		"ids": ids, "until": "2099-01-01T00:00:00Z",
	}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

// TestAlertsBulkSnoozePartialFailure mirrors
// TestChangesBulkResolvePartialFailure: a mixed batch (a real seeded alert
// and a nonexistent id) must not abort — each item gets its own outcome, and
// only the successful item is snoozed and audited.
func TestAlertsBulkSnoozePartialFailure(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")

	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")
	const missingID = "does-not-exist"

	rec := app.req(t, http.MethodPost, "/api/alerts/bulk-snooze", map[string]any{
		"ids":   []string{a.ID, missingID},
		"until": "2099-01-01T00:00:00Z",
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
	if len(body.Results) != 2 {
		t.Fatalf("got %d results, want 2: %+v", len(body.Results), body.Results)
	}

	outcomes := map[string]struct{ status, reason string }{}
	for _, r := range body.Results {
		outcomes[r.ID] = struct{ status, reason string }{r.Status, r.Reason}
	}

	if o := outcomes[a.ID]; o.status != "success" {
		t.Errorf("seeded alert outcome = %+v, want success", o)
	}
	if o := outcomes[missingID]; o.status != "error" || o.reason != "not_found" {
		t.Errorf("missing outcome = %+v, want error/not_found", o)
	}

	got, err := app.Store.GetAlert(context.Background(), a.ID)
	if err != nil {
		t.Fatalf("GetAlert: %v", err)
	}
	if got.Status != "snoozed" {
		t.Errorf("Status = %q, want snoozed", got.Status)
	}

	auditRec := app.req(t, http.MethodGet, "/api/system/audit?action=alert.bulk_snooze", nil, opToken)
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
	if page.Items[0].TargetID != a.ID {
		t.Errorf("audit TargetID = %q, want %q", page.Items[0].TargetID, a.ID)
	}
	if page.Items[0].ActorUserID != opUserID {
		t.Errorf("audit ActorUserID = %q, want %q", page.Items[0].ActorUserID, opUserID)
	}
}

// TestAlertsListDaysWindow verifies the optional ?days= range gate on
// GET /api/alerts: a recent alert plus a 30-day-old one, filtered down to
// just the recent one under ?days=1, with no days param defaulting back to
// "everything" for back-compat with the full Alerts page.
func TestAlertsListDaysWindow(t *testing.T) {
	app := newTestApp(t)
	viewerUserID, viewerToken := app.user(t, "viewer")
	app.connectorGrant(t, viewerUserID, "svc-1", "viewer")
	ctx := context.Background()
	now := time.Now().UTC()

	recent := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "info", Title: "recent", Description: "desc", Status: "pending",
		CreatedAt: now.Format(time.RFC3339),
	}
	if err := app.Store.CreateAlert(ctx, recent); err != nil {
		t.Fatalf("create recent alert: %v", err)
	}
	old := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "info", Title: "old", Description: "desc", Status: "pending",
		CreatedAt: now.AddDate(0, 0, -30).Format(time.RFC3339),
	}
	if err := app.Store.CreateAlert(ctx, old); err != nil {
		t.Fatalf("create old alert: %v", err)
	}

	tests := []struct {
		name      string
		query     string
		wantTotal int
	}{
		{"days=1 excludes old alert", "?days=1", 1},
		{"no days param returns everything", "", 2},
		{"days=0 ignored, returns everything", "?days=0", 2},
		{"negative days ignored, returns everything", "?days=-1", 2},
		{"malformed days ignored, returns everything", "?days=abc", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodGet, "/api/alerts"+tt.query, nil, viewerToken)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
			}
			var body struct {
				Items []struct {
					ID string `json:"id"`
				} `json:"items"`
				Total int `json:"total"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.Total != tt.wantTotal || len(body.Items) != tt.wantTotal {
				t.Fatalf("total = %d, items len = %d, want %d", body.Total, len(body.Items), tt.wantTotal)
			}
		})
	}
}
