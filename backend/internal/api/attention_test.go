package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestAttentionEmptyList(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Data     []any `json:"data"`
		Total    int   `json:"total"`
		Page     int   `json:"page"`
		PageSize int   `json:"pageSize"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Total != 0 || len(body.Data) != 0 {
		t.Errorf("expected empty list, got total=%d data len=%d", body.Total, len(body.Data))
	}
}

func TestAttentionMergesAlertsAndFindings(t *testing.T) {
	app := newTestApp(t)
	viewerID, viewerToken := app.user(t, "viewer")
	app.connectorGrant(t, viewerID, "svc-1", "viewer")

	alert := seedAlert(t, app)
	finding := seedQualityFinding(t, app)
	app.connectorGrant(t, viewerID, finding.ConnectorID, "viewer")

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Data []struct {
			ID          string `json:"id"`
			Kind        string `json:"kind"`
			Severity    string `json:"severity"`
			Title       string `json:"title"`
			ConnectorID string `json:"connectorId"`
			DetectedAt  string `json:"detectedAt"`
		} `json:"data"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Total != 2 || len(body.Data) != 2 {
		t.Fatalf("expected 2 items, got total=%d data len=%d", body.Total, len(body.Data))
	}

	foundAlert := false
	foundFinding := false
	for _, item := range body.Data {
		if item.ID == alert.ID {
			foundAlert = true
			if item.Kind != "alert" {
				t.Errorf("alert kind = %q, want alert", item.Kind)
			}
			if item.ConnectorID != alert.ServiceID {
				t.Errorf("alert connectorId = %q, want %q", item.ConnectorID, alert.ServiceID)
			}
			if item.Title != alert.Title {
				t.Errorf("alert title = %q, want %q", item.Title, alert.Title)
			}
		}
		if item.ID == finding.ID {
			foundFinding = true
			if item.Kind != "finding" {
				t.Errorf("finding kind = %q, want finding", item.Kind)
			}
			if item.ConnectorID != finding.ConnectorID {
				t.Errorf("finding connectorId = %q, want %q", item.ConnectorID, finding.ConnectorID)
			}
			if item.Title != finding.Title {
				t.Errorf("finding title = %q, want %q", item.Title, finding.Title)
			}
		}
	}

	if !foundAlert {
		t.Errorf("alert %s not found in response", alert.ID)
	}
	if !foundFinding {
		t.Errorf("finding %s not found in response", finding.ID)
	}
}

func TestAttentionSeverityOrdering(t *testing.T) {
	app := newTestApp(t)
	viewerID, viewerToken := app.user(t, "viewer")
	app.connectorGrant(t, viewerID, "svc-1", "viewer")
	ctx := context.Background()

	criticalAlert := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "critical", Title: "critical alert", Description: "desc", Status: "pending",
	}
	if err := app.Store.CreateAlert(ctx, criticalAlert); err != nil {
		t.Fatalf("create critical alert: %v", err)
	}

	warningAlert := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "warning", Title: "warning alert", Description: "desc", Status: "pending",
	}
	if err := app.Store.CreateAlert(ctx, warningAlert); err != nil {
		t.Fatalf("create warning alert: %v", err)
	}

	infoAlert := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "info", Title: "info alert", Description: "desc", Status: "pending",
	}
	if err := app.Store.CreateAlert(ctx, infoAlert); err != nil {
		t.Fatalf("create info alert: %v", err)
	}

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Data []struct {
			Severity string `json:"severity"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Data) != 3 {
		t.Fatalf("expected 3 items, got %d", len(body.Data))
	}

	if body.Data[0].Severity != "critical" {
		t.Errorf("first severity = %q, want critical", body.Data[0].Severity)
	}
	if body.Data[1].Severity != "warning" {
		t.Errorf("second severity = %q, want warning", body.Data[1].Severity)
	}
	if body.Data[2].Severity != "info" {
		t.Errorf("third severity = %q, want info", body.Data[2].Severity)
	}
}

func TestAttentionRunbookLinkForAlert(t *testing.T) {
	app := newTestApp(t)
	viewerID, viewerToken := app.user(t, "viewer")
	app.connectorGrant(t, viewerID, "svc-1", "viewer")

	alert := seedAlert(t, app)
	rb := seedRunbook(t, app, "alert_severity", alert.Severity)

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Data []struct {
			ID        string `json:"id"`
			RunbookID string `json:"runbookId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(body.Data))
	}

	if body.Data[0].ID != alert.ID {
		t.Errorf("alert id = %q, want %q", body.Data[0].ID, alert.ID)
	}
	if body.Data[0].RunbookID != rb.ID {
		t.Errorf("runbook id = %q, want %q", body.Data[0].RunbookID, rb.ID)
	}
}

func TestAttentionRunbookLinkForFinding(t *testing.T) {
	app := newTestApp(t)
	viewerID, viewerToken := app.user(t, "viewer")
	app.connectorGrant(t, viewerID, "svc-1", "viewer")

	finding := seedQualityFinding(t, app)
	app.connectorGrant(t, viewerID, finding.ConnectorID, "viewer")
	rb := seedRunbook(t, app, "finding_check_type", finding.CheckType)

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Data []struct {
			ID        string `json:"id"`
			RunbookID string `json:"runbookId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(body.Data))
	}

	if body.Data[0].ID != finding.ID {
		t.Errorf("finding id = %q, want %q", body.Data[0].ID, finding.ID)
	}
	if body.Data[0].RunbookID != rb.ID {
		t.Errorf("runbook id = %q, want %q", body.Data[0].RunbookID, rb.ID)
	}
}

// TestAttentionDaysWindow verifies the optional ?days= range gate on
// GET /api/attention: a recent+old alert and a recent+old finding, filtered
// down to just the two recent items under ?days=1, with no days param
// defaulting back to "everything" for back-compat.
func TestAttentionDaysWindow(t *testing.T) {
	app := newTestApp(t)
	viewerID, viewerToken := app.user(t, "viewer")
	app.connectorGrant(t, viewerID, "svc-1", "viewer")
	ctx := context.Background()
	now := time.Now().UTC()

	recentAlert := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "info", Title: "recent alert", Description: "desc", Status: "pending",
		CreatedAt: now.Format(time.RFC3339),
	}
	if err := app.Store.CreateAlert(ctx, recentAlert); err != nil {
		t.Fatalf("create recent alert: %v", err)
	}
	oldAlert := &store.AlertRecord{
		ServiceID: "svc-1", Severity: "info", Title: "old alert", Description: "desc", Status: "pending",
		CreatedAt: now.AddDate(0, 0, -30).Format(time.RFC3339),
	}
	if err := app.Store.CreateAlert(ctx, oldAlert); err != nil {
		t.Fatalf("create old alert: %v", err)
	}

	connector := &store.ConnectorRecord{Name: "attention days connector", Category: "networking", Type: "test", Enabled: true}
	if err := app.Store.CreateConnector(ctx, connector); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	app.connectorGrant(t, viewerID, connector.ID, "viewer")
	recentFinding := &store.QualityFindingRecord{
		ConnectorID: connector.ID, CheckType: "stale", Severity: "warning",
		Title: "recent finding", Description: "desc",
		FirstDetectedAt: now.Format(time.RFC3339), LastSeenAt: now.Format(time.RFC3339),
	}
	if err := app.Store.UpsertQualityFinding(ctx, recentFinding); err != nil {
		t.Fatalf("upsert recent finding: %v", err)
	}
	old := now.AddDate(0, 0, -30).Format(time.RFC3339)
	oldFinding := &store.QualityFindingRecord{
		ConnectorID: connector.ID, CheckType: "empty", Severity: "warning",
		Title: "old finding", Description: "desc",
		FirstDetectedAt: old, LastSeenAt: old,
	}
	if err := app.Store.UpsertQualityFinding(ctx, oldFinding); err != nil {
		t.Fatalf("upsert old finding: %v", err)
	}

	tests := []struct {
		name      string
		query     string
		wantTotal int
	}{
		{"days=1 excludes old items across both origins", "?days=1", 2},
		{"no days param returns everything", "", 4},
		{"days=0 ignored, returns everything", "?days=0", 4},
		{"negative days ignored, returns everything", "?days=-1", 4},
		{"malformed days ignored, returns everything", "?days=abc", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodGet, "/api/attention"+tt.query, nil, viewerToken)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
			}
			var body struct {
				Data  []any `json:"data"`
				Total int   `json:"total"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.Total != tt.wantTotal || len(body.Data) != tt.wantTotal {
				t.Fatalf("total = %d, data len = %d, want %d", body.Total, len(body.Data), tt.wantTotal)
			}
		})
	}
}

func TestAttentionAuthenticatedAccess(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 for viewer role; body = %s", rec.Code, rec.Body)
	}
}

func TestAttentionHidesUngrantedConnectors(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	seedAlert(t, app)
	seedQualityFinding(t, app)

	rec := app.req(t, http.MethodGet, "/api/attention", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Total != 0 {
		t.Errorf("user without grants saw %d items, want 0", body.Total)
	}
}
