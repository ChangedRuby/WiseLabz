package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedRunbook(t *testing.T, app *testApp, targetType, targetValue string) *store.RunbookRecord {
	t.Helper()
	rb, err := app.Store.CreateRunbook(context.Background(), &store.RunbookRecord{
		Title:       "Runbook for " + targetValue,
		Body:        "Do the thing.",
		TargetType:  targetType,
		TargetValue: targetValue,
	})
	if err != nil {
		t.Fatalf("seed runbook: %v", err)
	}
	return rb
}

func TestRunbooksListSuccess(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	seedRunbook(t, app, "change_type", "vm.created")

	rec := app.req(t, http.MethodGet, "/api/runbooks", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbooksListFilterByChangeType(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	seedRunbook(t, app, "change_type", "vm.created")
	seedRunbook(t, app, "alert_severity", "critical")

	rec := app.req(t, http.MethodGet, "/api/runbooks?changeType=vm.created", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var page struct {
		Items []store.RunbookRecord `json:"items"`
		Total int                   `json:"total"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("total = %d, len(items) = %d, want 1, 1", page.Total, len(page.Items))
	}
	if page.Items[0].TargetValue != "vm.created" {
		t.Errorf("TargetValue = %q, want vm.created", page.Items[0].TargetValue)
	}
}

func TestRunbooksListRejectsBothFilters(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/runbooks?changeType=vm.created&alertSeverity=critical", nil, viewerToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbooksCreateRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	body := map[string]any{"title": "t", "targetType": "change_type", "targetValue": "vm.created"}
	rec := app.req(t, http.MethodPost, "/api/runbooks", body, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbooksCreateSuccess(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	body := map[string]any{"title": "t", "targetType": "change_type", "targetValue": "vm.created"}
	rec := app.req(t, http.MethodPost, "/api/runbooks", body, opToken)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbooksGetNotFound(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/runbooks/unknown-id", nil, viewerToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbooksUpdateNotFound(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPut, "/api/runbooks/unknown-id", map[string]any{"title": "new"}, opToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbooksDeleteNotFound(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodDelete, "/api/runbooks/unknown-id", nil, opToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

func seedProxmoxConnector(t *testing.T, app *testApp, url string) string {
	t.Helper()
	c := &store.ConnectorRecord{Name: "Proxmox", Category: "virtualization", Type: "proxmox", URL: url}
	if err := app.Store.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("CreateConnector() error: %v", err)
	}
	return c.ID
}

type runbookStepResp struct {
	ID                   string `json:"id"`
	ConnectorID          string `json:"connectorId"`
	ConnectorName        string `json:"connectorName"`
	Verb                 string `json:"verb"`
	EntityRef            string `json:"entityRef"`
	CanExecute           bool   `json:"canExecute"`
	ExecuteBlockedReason string `json:"executeBlockedReason"`
}

type runbookResp struct {
	ID    string            `json:"id"`
	Steps []runbookStepResp `json:"steps"`
}

// createRunbookWithStep creates a runbook (as instance admin) with one step
// pointing at connectorID/verb/entityRef and returns the runbook and its
// single step's ID.
func createRunbookWithStep(t *testing.T, app *testApp, opToken, targetValue, connectorID, verb, entityRef string) (runbookResp, string) {
	t.Helper()
	body := map[string]any{
		"title":       "Restart it",
		"targetType":  "change_type",
		"targetValue": targetValue,
		"steps": []map[string]any{
			{"title": "Restart", "connectorId": connectorID, "verb": verb, "entityRef": entityRef},
		},
	}
	rec := app.req(t, http.MethodPost, "/api/runbooks", body, opToken)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create runbook: status = %d, want 201; body = %s", rec.Code, rec.Body)
	}
	var rb runbookResp
	if err := json.Unmarshal(rec.Body.Bytes(), &rb); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(rb.Steps) != 1 {
		t.Fatalf("len(steps) = %d, want 1", len(rb.Steps))
	}
	return rb, rb.Steps[0].ID
}

func TestRunbookExecuteStepForbiddenWithoutOperatorGrant(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	viewerID, viewerToken := app.user(t, "viewer")
	connID := seedProxmoxConnector(t, app, "https://example.com")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.forbidden", connID, "restart", "")
	_ = viewerID

	rec := app.req(t, http.MethodPost, "/api/runbooks/"+rb.ID+"/steps/"+stepID+"/execute", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestRunbookExecuteStepMissingElevationToken(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	userID, userToken := app.user(t, "viewer")
	connID := seedProxmoxConnector(t, app, "https://example.com")
	app.connectorGrant(t, userID, connID, "operator")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.notoken", connID, "restart", "")

	rec := app.req(t, http.MethodPost, "/api/runbooks/"+rb.ID+"/steps/"+stepID+"/execute", nil, userToken)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "elevation_required") {
		t.Fatalf("status = %d, body = %s, want 400 elevation_required", rec.Code, rec.Body)
	}
}

func TestRunbookExecuteStepWrongActionTokenRejected(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	userID, userToken := app.user(t, "viewer")
	connID := seedProxmoxConnector(t, app, "https://example.com")
	app.connectorGrant(t, userID, connID, "operator")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.wrongaction", connID, "restart", "")

	// Elevation token scoped to connector.start, but the step executes restart.
	wrongToken := app.elevationToken(t, userID, "connector.start")
	req := app.newRequest(t, http.MethodPost, "/api/runbooks/"+rb.ID+"/steps/"+stepID+"/execute", nil, userToken)
	req.Header.Set("X-Elevation-Token", wrongToken)
	rec := app.serve(req)
	if rec.Code != http.StatusUnauthorized || !strings.Contains(rec.Body.String(), "connector.restart") {
		t.Fatalf("status = %d, body = %s, want 401 naming the required action connector.restart", rec.Code, rec.Body)
	}
}

func TestRunbookExecuteStepDryRunWithoutTokenWorks(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	userID, userToken := app.user(t, "viewer")
	connID := seedProxmoxConnector(t, app, "https://example.com")
	app.connectorGrant(t, userID, connID, "operator")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.dryrun", connID, "restart", "")

	rec := app.req(t, http.MethodPost, "/api/runbooks/"+rb.ID+"/steps/"+stepID+"/execute?dryRun=true", nil, userToken)
	// No snapshot exists yet, so this 404s on "no snapshot" rather than
	// elevation_required — the point is that dry-run never asked for a token.
	if rec.Code == http.StatusBadRequest && strings.Contains(rec.Body.String(), "elevation_required") {
		t.Fatalf("dry-run required elevation: %s", rec.Body)
	}
}

func TestRunbookExecuteStepSuccessWritesConnectorAuditWithRunbookContext(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	userID, userToken := app.user(t, "viewer")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cluster/resources":
			_, _ = w.Write([]byte(`{"data":[{"vmid":100,"node":"pve1","type":"qemu"}]}`))
		case "/nodes/pve1/qemu/100/status/reboot":
			_, _ = w.Write([]byte(`{"data":null}`))
		default:
			t.Fatalf("unexpected request: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	connID := seedProxmoxConnector(t, app, server.URL)
	app.connectorGrant(t, userID, connID, "operator")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.success", connID, "restart", "100")

	token := app.elevationToken(t, userID, "connector.restart")
	req := app.newRequest(t, http.MethodPost, "/api/runbooks/"+rb.ID+"/steps/"+stepID+"/execute", nil, userToken)
	req.Header.Set("X-Elevation-Token", token)
	rec := app.serve(req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	records, _, err := app.Store.ListAuditRecords(context.Background(), "connector.restart", "connector", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}
	var detail map[string]any
	if err := json.Unmarshal([]byte(records[0].Detail), &detail); err != nil {
		t.Fatalf("unmarshal detail: %v", err)
	}
	if detail["runbookId"] != rb.ID || detail["stepId"] != stepID {
		t.Errorf("detail = %+v, want runbookId=%s stepId=%s", detail, rb.ID, stepID)
	}
	if detail["entityRef"] != "100" {
		t.Errorf("detail[entityRef] = %v, want 100", detail["entityRef"])
	}
}

func TestRunbookExecuteStepFailureCreatesAlert(t *testing.T) {
	connector.AllowLoopbackForTest(t)
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	userID, userToken := app.user(t, "viewer")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"data":[]}`)) // vmid never found -> restart fails
	}))
	defer server.Close()

	connID := seedProxmoxConnector(t, app, server.URL)
	app.connectorGrant(t, userID, connID, "operator")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.failure", connID, "restart", "100")

	token := app.elevationToken(t, userID, "connector.restart")
	req := app.newRequest(t, http.MethodPost, "/api/runbooks/"+rb.ID+"/steps/"+stepID+"/execute", nil, userToken)
	req.Header.Set("X-Elevation-Token", token)
	rec := app.serve(req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502; body = %s", rec.Code, rec.Body)
	}

	alerts, _, err := app.Store.ListAlerts(context.Background(), "", "", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAlerts() error: %v", err)
	}
	found := false
	for _, a := range alerts {
		if a.ServiceID == connID && a.Severity == "critical" {
			found = true
		}
	}
	if !found {
		t.Fatalf("no critical alert for connector %s; alerts = %+v", connID, alerts)
	}
}

func TestRunbookStepsCanExecuteReflectsOperatorGrant(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	userID, userToken := app.user(t, "viewer")
	connID := seedProxmoxConnector(t, app, "https://example.com")

	rb, _ := createRunbookWithStep(t, app, opToken, "vm.canexecute", connID, "restart", "")

	rec := app.req(t, http.MethodGet, "/api/runbooks/"+rb.ID, nil, userToken)
	var got runbookResp
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Steps[0].CanExecute {
		t.Error("canExecute should be false before granting operator")
	}
	if got.Steps[0].ExecuteBlockedReason != "no_operator_grant" {
		t.Errorf("executeBlockedReason = %q, want no_operator_grant", got.Steps[0].ExecuteBlockedReason)
	}

	app.connectorGrant(t, userID, connID, "operator")
	rec2 := app.req(t, http.MethodGet, "/api/runbooks/"+rb.ID, nil, userToken)
	var got2 runbookResp
	if err := json.Unmarshal(rec2.Body.Bytes(), &got2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !got2.Steps[0].CanExecute {
		t.Error("canExecute should be true after granting operator")
	}
	if got2.Steps[0].ExecuteBlockedReason != "" {
		t.Errorf("executeBlockedReason = %q, want empty", got2.Steps[0].ExecuteBlockedReason)
	}
}

func TestRunbooksCreateStepValidationErrors(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	connID := seedProxmoxConnector(t, app, "https://example.com")

	cases := map[string]map[string]any{
		"invalid verb":              {"title": "t", "connectorId": connID, "verb": "reboot"},
		"unsupported for connector": {"title": "t", "connectorId": connID, "verb": "restart", "entityRef": "../etc"},
		"unknown connector":         {"title": "t", "connectorId": "missing", "verb": "restart"},
	}
	for name, step := range cases {
		t.Run(name, func(t *testing.T) {
			body := map[string]any{
				"title": "t", "targetType": "change_type", "targetValue": "invalid-" + name,
				"steps": []map[string]any{step},
			}
			rec := app.req(t, http.MethodPost, "/api/runbooks", body, opToken)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestRunbooksAuditRows(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	connID := seedProxmoxConnector(t, app, "https://example.com")

	rb, stepID := createRunbookWithStep(t, app, opToken, "vm.audit", connID, "restart", "")

	createRecords, _, err := app.Store.ListAuditRecords(context.Background(), "runbook.create", "runbook", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords(create) error: %v", err)
	}
	if len(createRecords) != 1 || createRecords[0].TargetID != rb.ID {
		t.Fatalf("createRecords = %+v, want one row for runbook %s", createRecords, rb.ID)
	}
	var createDetail map[string]any
	if err := json.Unmarshal([]byte(createRecords[0].Detail), &createDetail); err != nil {
		t.Fatalf("unmarshal create detail: %v", err)
	}
	if createDetail["title"] != "Restart it" {
		t.Errorf("createDetail[title] = %v, want Restart it", createDetail["title"])
	}

	updateBody := map[string]any{"title": "Renamed"}
	rec := app.req(t, http.MethodPut, "/api/runbooks/"+rb.ID, updateBody, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d; body = %s", rec.Code, rec.Body)
	}
	updateRecords, _, err := app.Store.ListAuditRecords(context.Background(), "runbook.update", "runbook", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords(update) error: %v", err)
	}
	if len(updateRecords) != 1 {
		t.Fatalf("len(updateRecords) = %d, want 1", len(updateRecords))
	}

	delRec := app.req(t, http.MethodDelete, "/api/runbooks/"+rb.ID, nil, opToken)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d; body = %s", delRec.Code, delRec.Body)
	}
	deleteRecords, _, err := app.Store.ListAuditRecords(context.Background(), "runbook.delete", "runbook", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords(delete) error: %v", err)
	}
	if len(deleteRecords) != 1 {
		t.Fatalf("len(deleteRecords) = %d, want 1", len(deleteRecords))
	}
	_ = stepID
}

func TestRunbooksFindingCheckTypeAccepted(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")

	body := map[string]any{"title": "t", "targetType": "finding_check_type", "targetValue": "stale"}
	rec := app.req(t, http.MethodPost, "/api/runbooks", body, opToken)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body)
	}

	listRec := app.req(t, http.MethodGet, "/api/runbooks?findingCheckType=stale", nil, viewerToken)
	if listRec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", listRec.Code, listRec.Body)
	}
	var page struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if page.Total != 1 {
		t.Fatalf("total = %d, want 1", page.Total)
	}
}
