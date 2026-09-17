package api_test

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestAuditListRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/system/audit", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

// TestConnectorCreateProducesAuditRecord exercises one instrumented action
// end-to-end: creating a connector must leave a matching row retrievable
// through the audit list endpoint, attributed to the acting operator.
func TestConnectorCreateProducesAuditRecord(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")

	createRec := app.req(t, http.MethodPost, "/api/connectors", map[string]any{
		"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com",
	}, opToken)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", createRec.Code, createRec.Body)
	}
	var created store.ConnectorRecord
	if err := json.Unmarshal(createRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create body: %v", err)
	}

	listRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200; body = %s", listRec.Code, listRec.Body)
	}

	var page struct {
		Items []store.AuditRecord `json:"items"`
		Total int                 `json:"total"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode list body: %v", err)
	}

	var found *store.AuditRecord
	for i := range page.Items {
		if page.Items[i].Action == "connector.create" && page.Items[i].TargetID == created.ID {
			found = &page.Items[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("no audit record for connector.create/%s in %+v", created.ID, page.Items)
	}
	if found.ActorUserID != opUserID {
		t.Errorf("ActorUserID = %q, want %q", found.ActorUserID, opUserID)
	}
	if found.ActorRole != "admin" {
		t.Errorf("ActorRole = %q, want admin", found.ActorRole)
	}
	if found.TargetType != "connector" {
		t.Errorf("TargetType = %q, want connector", found.TargetType)
	}
}

func TestAuditListFiltersByAction(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	for _, name := range []string{"svc-a", "svc-b"} {
		rec := app.req(t, http.MethodPost, "/api/connectors", map[string]any{
			"name": name, "category": "virtualization", "type": "proxmox", "url": "https://example.com",
		}, opToken)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d, want 201; body = %s", name, rec.Code, rec.Body)
		}
	}

	rec := app.req(t, http.MethodGet, "/api/system/audit?action=connector.create", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var page struct {
		Items []store.AuditRecord `json:"items"`
		Total int                 `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if page.Total != 2 {
		t.Fatalf("Total = %d, want 2", page.Total)
	}
	for _, r := range page.Items {
		if r.Action != "connector.create" {
			t.Errorf("item action = %q, want connector.create", r.Action)
		}
	}
}

func TestElevationAttemptAuditThroughHTTPRoutes(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")

	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	missing := app.req(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, opToken)
	if missing.Code != http.StatusBadRequest {
		t.Fatalf("missing status = %d, want 400; body = %s", missing.Code, missing.Body)
	}

	req := app.newRequest(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, opToken)
	req.Header.Set("X-Elevation-Token", "garbage")
	invalid := app.serve(req)
	if invalid.Code != http.StatusUnauthorized {
		t.Fatalf("invalid status = %d, want 401; body = %s", invalid.Code, invalid.Body)
	}

	listRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
	if listRec.Code != http.StatusOK {
		t.Fatalf("audit list status = %d, want 200; body = %s", listRec.Code, listRec.Body)
	}
	var page struct {
		Items []store.AuditRecord `json:"items"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode list body: %v", err)
	}

	var requested, denied *store.AuditRecord
	for i := range page.Items {
		switch page.Items[i].Action {
		case "auth.elevation_requested":
			requested = &page.Items[i]
		case "auth.elevation_denied":
			denied = &page.Items[i]
		}
	}
	if requested == nil {
		t.Fatal("no auth.elevation_requested audit record")
	}
	if denied == nil {
		t.Fatal("no auth.elevation_denied audit record")
	}
	if requested.TargetType != "action" || requested.TargetID != "connector.delete" {
		t.Fatalf("requested target = %s/%s, want action/connector.delete", requested.TargetType, requested.TargetID)
	}
	var detail map[string]any
	if err := json.Unmarshal([]byte(denied.Detail), &detail); err != nil {
		t.Fatalf("unmarshal denied detail: %v", err)
	}
	if detail["action"] != "connector.delete" || detail["reason"] != "invalid" {
		t.Fatalf("denied detail = %#v, want action connector.delete reason invalid", detail)
	}
}

func TestAlertResolveProducesAuditRecord(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/resolve", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("resolve status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	auditRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
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

	var found *store.AuditRecord
	for i := range page.Items {
		if page.Items[i].Action == "alert.resolve" && page.Items[i].TargetID == a.ID {
			found = &page.Items[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("no audit record for alert.resolve/%s", a.ID)
	}
	if found.TargetType != "alert" {
		t.Errorf("TargetType = %q, want alert", found.TargetType)
	}
	if found.ActorUserID != opUserID {
		t.Errorf("ActorUserID = %q, want %q", found.ActorUserID, opUserID)
	}
	if found.ActorRole != "admin" {
		t.Errorf("ActorRole = %q, want admin", found.ActorRole)
	}
}

func TestAlertDismissProducesAuditRecord(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/dismiss", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("dismiss status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	auditRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
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

	var found *store.AuditRecord
	for i := range page.Items {
		if page.Items[i].Action == "alert.dismiss" && page.Items[i].TargetID == a.ID {
			found = &page.Items[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("no audit record for alert.dismiss/%s", a.ID)
	}
	if found.TargetType != "alert" {
		t.Errorf("TargetType = %q, want alert", found.TargetType)
	}
	if found.ActorUserID != opUserID {
		t.Errorf("ActorUserID = %q, want %q", found.ActorUserID, opUserID)
	}
	if found.ActorRole != "admin" {
		t.Errorf("ActorRole = %q, want admin", found.ActorRole)
	}
}

func TestAlertSnoozeProducesAuditRecord(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	a := seedAlert(t, app)
	app.connectorGrant(t, opUserID, a.ServiceID, "operator")

	rec := app.req(t, http.MethodPost, "/api/alerts/"+a.ID+"/snooze", map[string]any{"until": "2099-01-01T00:00:00Z"}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("snooze status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	auditRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
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

	var found *store.AuditRecord
	for i := range page.Items {
		if page.Items[i].Action == "alert.snooze" && page.Items[i].TargetID == a.ID {
			found = &page.Items[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("no audit record for alert.snooze/%s", a.ID)
	}
	if found.TargetType != "alert" {
		t.Errorf("TargetType = %q, want alert", found.TargetType)
	}
	if found.ActorUserID != opUserID {
		t.Errorf("ActorUserID = %q, want %q", found.ActorUserID, opUserID)
	}
	if found.ActorRole != "admin" {
		t.Errorf("ActorRole = %q, want admin", found.ActorRole)
	}
}

func TestFindingResolveProducesAuditRecord(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	finding := seedQualityFinding(t, app)
	app.connectorGrant(t, opUserID, finding.ConnectorID, "operator")

	rec := app.req(t, http.MethodPost, "/api/findings/"+finding.ID+"/resolve", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("resolve status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	auditRec := app.req(t, http.MethodGet, "/api/system/audit", nil, opToken)
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

	var found *store.AuditRecord
	for i := range page.Items {
		if page.Items[i].Action == "finding.resolve" && page.Items[i].TargetID == finding.ID {
			found = &page.Items[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("no audit record for finding.resolve/%s", finding.ID)
	}
	if found.TargetType != "finding" {
		t.Errorf("TargetType = %q, want finding", found.TargetType)
	}
	if found.ActorUserID != opUserID {
		t.Errorf("ActorUserID = %q, want %q", found.ActorUserID, opUserID)
	}
	if found.ActorRole != "admin" {
		t.Errorf("ActorRole = %q, want admin", found.ActorRole)
	}
}

func TestAuditExportRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/system/audit/export", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

// seedAuditRecords writes n audit records directly via the store, bypassing
// the API, so export tests aren't coupled to any particular audited action.
func seedAuditRecords(t *testing.T, app *testApp, action, targetType string, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		if err := app.Store.CreateAuditRecord(context.Background(), &store.AuditRecord{
			ActorUserID: "seed-user", ActorRole: "operator", Action: action, TargetType: targetType, TargetID: "t",
		}); err != nil {
			t.Fatalf("CreateAuditRecord() error: %v", err)
		}
	}
}

func TestAuditExportJSON(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	seedAuditRecords(t, app, "export.test", "widget", 3)

	rec := app.req(t, http.MethodGet, "/api/system/audit/export?format=json", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	if got := rec.Header().Get("Content-Disposition"); got == "" {
		t.Errorf("Content-Disposition header missing")
	}

	var records []store.AuditRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &records); err != nil {
		t.Fatalf("decode export body: %v", err)
	}
	count := 0
	for _, r := range records {
		if r.Action == "export.test" {
			count++
		}
	}
	if count != 3 {
		t.Errorf("export.test records in export = %d, want 3", count)
	}
}

func TestAuditExportCSV(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	seedAuditRecords(t, app, "export.csv.test", "widget", 2)

	rec := app.req(t, http.MethodGet, "/api/system/audit/export?format=csv", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/csv; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/csv; charset=utf-8", ct)
	}
	if got := rec.Header().Get("Content-Disposition"); got == "" {
		t.Errorf("Content-Disposition header missing")
	}

	reader := csv.NewReader(strings.NewReader(rec.Body.String()))
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("parse CSV body: %v", err)
	}
	if len(rows) < 1 {
		t.Fatalf("CSV body has no rows")
	}
	wantHeader := []string{"id", "actorUserId", "actorRole", "action", "targetType", "targetId", "detail", "createdAt"}
	if !slicesEqual(rows[0], wantHeader) {
		t.Errorf("CSV header = %v, want %v", rows[0], wantHeader)
	}
	count := 0
	for _, row := range rows[1:] {
		if len(row) > 3 && row[3] == "export.csv.test" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("export.csv.test rows in CSV = %d, want 2", count)
	}
}

func TestAuditExportInvalidFormat(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodGet, "/api/system/audit/export?format=xml", nil, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestAuditExportFilters(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	seedAuditRecords(t, app, "export.filter.match", "gadget", 2)
	seedAuditRecords(t, app, "export.filter.other", "widget", 5)

	rec := app.req(t, http.MethodGet, "/api/system/audit/export?format=json&action=export.filter.match&targetType=gadget", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var records []store.AuditRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &records); err != nil {
		t.Fatalf("decode export body: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("filtered export returned %d records, want 2", len(records))
	}
	for _, r := range records {
		if r.Action != "export.filter.match" || r.TargetType != "gadget" {
			t.Errorf("unexpected record in filtered export: %+v", r)
		}
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
