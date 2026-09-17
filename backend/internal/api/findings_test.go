package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedQualityFinding(t *testing.T, app *testApp) *store.QualityFindingRecord {
	t.Helper()
	ctx := context.Background()
	connector := &store.ConnectorRecord{Name: "quality connector", Category: "networking", Type: "test", Enabled: true}
	if err := app.Store.CreateConnector(ctx, connector); err != nil {
		t.Fatalf("CreateConnector: %v", err)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	finding := &store.QualityFindingRecord{
		ConnectorID: connector.ID, CheckType: "ownership_incomplete", Severity: "warning",
		Title: "Connector has no owner", Description: "Assign an owner", RemediationLink: "/connectors/" + connector.ID + "/edit",
		FirstDetectedAt: now, LastSeenAt: now,
	}
	if err := app.Store.UpsertQualityFinding(ctx, finding); err != nil {
		t.Fatalf("UpsertQualityFinding: %v", err)
	}
	return finding
}

func TestFindingsListSuccess(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")
	finding := seedQualityFinding(t, app)
	app.connectorGrant(t, userID, finding.ConnectorID, "viewer")

	rec := app.req(t, http.MethodGet, "/api/findings?status=open&checkType=ownership_incomplete", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		Items []struct {
			ID            string `json:"id"`
			ConnectorName string `json:"connectorName"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Total != 1 || len(body.Items) != 1 || body.Items[0].ID != finding.ID || body.Items[0].ConnectorName != "quality connector" {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFindingsGetSuccess(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "viewer")
	finding := seedQualityFinding(t, app)
	app.connectorGrant(t, userID, finding.ConnectorID, "viewer")

	rec := app.req(t, http.MethodGet, "/api/findings/"+finding.ID, nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var body struct {
		ID              string  `json:"id"`
		ConnectorName   string  `json:"connectorName"`
		DocID           *string `json:"docId"`
		RemediationLink string  `json:"remediationLink"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ID != finding.ID || body.ConnectorName != "quality connector" || body.DocID != nil || body.RemediationLink != finding.RemediationLink {
		t.Fatalf("unexpected response: %+v", body)
	}
}

func TestFindingsResolveRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, token := app.user(t, "viewer")
	finding := seedQualityFinding(t, app)

	rec := app.req(t, http.MethodPost, "/api/findings/"+finding.ID+"/resolve", nil, token)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestFindingsResolveSuccess(t *testing.T) {
	app := newTestApp(t)
	userID, token := app.user(t, "operator")
	finding := seedQualityFinding(t, app)
	app.connectorGrant(t, userID, finding.ConnectorID, "operator")

	rec := app.req(t, http.MethodPost, "/api/findings/"+finding.ID+"/resolve", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	got, err := app.Store.GetQualityFinding(context.Background(), finding.ID)
	if err != nil {
		t.Fatalf("GetQualityFinding: %v", err)
	}
	if got.Status != "resolved" || got.ResolvedAt == "" {
		t.Fatalf("resolved finding = %+v", got)
	}
}
