package connectors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/api/apitest"
	"github.com/WiseLabz/wiselabz/internal/auth"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func actionRequest(id, body string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.SetPathValue("id", id)
	return r
}

func actionResponse(t *testing.T, fn http.HandlerFunc, r *http.Request, status int) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	fn(rr, r)
	if rr.Code != status {
		t.Fatalf("status=%d, want %d: %s", rr.Code, status, rr.Body.String())
	}
	return rr
}

func TestActionMaintenanceLifecycle(t *testing.T) {
	h := newTestHandler(t)
	user := apitest.NewUser(t, h.Store, "operator")
	visible := seedCoverageConnector(t, h, "visible", "networking")
	hidden := seedCoverageConnector(t, h, "hidden", "networking")
	apitest.GrantConnectorRole(t, h.Store, user, visible.ID, "viewer")
	request := func(id, body string) *http.Request {
		r := actionRequest(id, body)
		return r.WithContext(auth.ContextWithUser(r.Context(), user, false))
	}
	for _, id := range []string{visible.ID, hidden.ID} {
		rr := actionResponse(t, h.OpenMaintenanceWindow, request(id, `{"durationMinutes":30}`), 201)
		var window store.MaintenanceWindowRecord
		if err := json.Unmarshal(rr.Body.Bytes(), &window); err != nil {
			t.Fatal(err)
		}
		start, err := time.Parse(time.RFC3339, window.StartsAt)
		if err != nil {
			t.Fatal(err)
		}
		end, err := time.Parse(time.RFC3339, window.EndsAt)
		if err != nil {
			t.Fatal(err)
		}
		if window.ID == "" || window.ConnectorID != id || window.CreatedBy != user || end.Sub(start) != 30*time.Minute {
			t.Fatalf("window=%+v", window)
		}
		rr = actionResponse(t, h.GetMaintenanceWindow, request(id, ""), 200)
		if !strings.Contains(rr.Body.String(), window.ID) {
			t.Fatalf("get=%s", rr.Body.String())
		}
	}
	rr := actionResponse(t, h.ListActiveMaintenance, request("", ""), 200)
	var windows []store.MaintenanceWindowRecord
	if err := json.Unmarshal(rr.Body.Bytes(), &windows); err != nil {
		t.Fatal(err)
	}
	if len(windows) != 1 || windows[0].ConnectorID != visible.ID {
		t.Fatalf("visible windows=%+v", windows)
	}
	for _, want := range []string{`{"closed":true}`, `{"closed":false}`} {
		rr = actionResponse(t, h.CloseMaintenanceWindow, request(visible.ID, ""), 200)
		if strings.TrimSpace(rr.Body.String()) != want {
			t.Fatalf("close=%s", rr.Body.String())
		}
	}
	rr = actionResponse(t, h.GetMaintenanceWindow, request(visible.ID, ""), 200)
	if strings.TrimSpace(rr.Body.String()) != "null" {
		t.Fatalf("closed window=%s", rr.Body.String())
	}
	rr = actionResponse(t, h.ListActiveMaintenance, request("", ""), 200)
	if strings.TrimSpace(rr.Body.String()) != "[]" {
		t.Fatalf("list=%s", rr.Body.String())
	}
	for action, count := range map[string]int{"connector.maintenanceWindow.open": 2, "connector.maintenanceWindow.close": 1} {
		rows, _, err := h.Store.ListAuditRecords(context.Background(), action, "connector", "", "", 0, 10)
		if err != nil || len(rows) != count {
			t.Fatalf("audit %s=%+v err=%v", action, rows, err)
		}
	}
	for _, body := range []string{`{`, `{}`, `{"durationMinutes":-1}`} {
		actionResponse(t, h.OpenMaintenanceWindow, request(visible.ID, body), 400)
	}
	actionResponse(t, h.OpenMaintenanceWindow, request("missing", `{"durationMinutes":30}`), 404)
}

func TestActionPermissions(t *testing.T) {
	h := newTestHandler(t)
	c := seedCoverageConnector(t, h, "permissions", "networking")
	user := apitest.NewUser(t, h.Store, "viewer")
	request := func(id, uid, body string) *http.Request {
		r := actionRequest(id, body)
		r.SetPathValue("userId", uid)
		return r
	}
	for _, role := range []string{"viewer", "operator"} {
		actionResponse(t, h.PutPermission, request(c.ID, user, `{"role":"`+role+`"}`), 200)
		got, err := h.Store.GetUserConnectorRole(context.Background(), user, c.ID)
		if err != nil || got != role {
			t.Fatalf("role=%q err=%v", got, err)
		}
		rr := actionResponse(t, h.ListPermissions, request(c.ID, user, ""), 200)
		var grants []store.ConnectorGrant
		if err := json.Unmarshal(rr.Body.Bytes(), &grants); err != nil {
			t.Fatal(err)
		}
		if len(grants) != 1 || grants[0].Role != role || grants[0].UserID != user {
			t.Fatalf("grants=%+v", grants)
		}
	}
	for _, body := range []string{`{`, `{"role":"admin"}`} {
		actionResponse(t, h.PutPermission, request(c.ID, user, body), 400)
	}
	actionResponse(t, h.PutPermission, request("missing", user, `{"role":"viewer"}`), 404)
	actionResponse(t, h.PutPermission, request(c.ID, "missing", `{"role":"viewer"}`), 404)
	actionResponse(t, h.ListPermissions, request("missing", user, ""), 404)
	actionResponse(t, h.DeletePermission, request(c.ID, user, ""), 204)
	got, err := h.Store.GetUserConnectorRole(context.Background(), user, c.ID)
	if err != nil || got != "" {
		t.Fatalf("revoked role=%q err=%v", got, err)
	}
	actionResponse(t, h.DeletePermission, request(c.ID, user, ""), 404)
	for action, count := range map[string]int{"connector.permission.granted": 2, "connector.permission.revoked": 1} {
		rows, _, err := h.Store.ListAuditRecords(context.Background(), action, "connector", "", "", 0, 10)
		if err != nil || len(rows) != count {
			t.Fatalf("audit %s=%+v err=%v", action, rows, err)
		}
	}
}

func TestActionLifecyclePreviews(t *testing.T) {
	for _, verb := range []string{"restart", "start", "stop"} {
		for _, data := range []string{"", `{`, `{"serviceName":"lab"}`, `{"serviceName":"lab","dependencies":[{"serviceName":"database"}]}`} {
			t.Run(verb+"/"+data, func(t *testing.T) {
				h := newTestHandler(t)
				fn := map[string]http.HandlerFunc{"restart": h.RestartPreview, "start": h.StartPreview, "stop": h.StopPreview}[verb]
				c := seedCoverageConnector(t, h, "lab", "networking")
				status := 200
				switch data {
				case "":
					status = 404
				case "{":
					status = 500
				}
				if data != "" {
					if err := h.Store.CreateSnapshot(context.Background(), &store.SnapshotRecord{ConnectorID: c.ID, Data: data, FetchedAt: "2025-01-01T00:00:00Z"}); err != nil {
						t.Fatal(err)
					}
				}
				r := actionRequest(c.ID, "")
				r.URL.RawQuery = "dryRun=true"
				rr := actionResponse(t, fn, r, status)
				if status == 200 {
					var result struct {
						Target       string            `json:"targetService"`
						Downtime     int               `json:"estimatedDowntimeSeconds"`
						Dependencies []json.RawMessage `json:"dependentServices"`
					}
					if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
						t.Fatal(err)
					}
					downtime := 30
					if verb == "stop" {
						downtime = 0
					}
					count := 0
					if strings.Contains(data, "dependencies") {
						count = 1
					}
					if result.Target != "lab" || result.Downtime != downtime || result.Dependencies == nil || len(result.Dependencies) != count {
						t.Fatalf("preview=%s", rr.Body.String())
					}
				}
				r = actionRequest("missing", "")
				r.URL.RawQuery = "dryRun=true"
				actionResponse(t, fn, r, 404)
			})
		}
	}
}

func TestActionStoreFailures(t *testing.T) {
	h := newTestHandler(t)
	for name, fn := range map[string]http.HandlerFunc{
		"sync": h.Sync, "bulk sync": h.BulkSync, "bulk restart": h.BulkRestart, "bulk reauth": h.BulkReauth,
		"restart": h.RestartPreview, "start": h.StartPreview, "stop": h.StopPreview, "config push": h.ConfigPush,
		"open maintenance": h.OpenMaintenanceWindow, "close maintenance": h.CloseMaintenanceWindow, "get maintenance": h.GetMaintenanceWindow, "list maintenance": h.ListActiveMaintenance,
		"list permissions": h.ListPermissions, "put permission": h.PutPermission, "delete permission": h.DeletePermission,
	} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			r := actionRequest("missing", `{"ids":["missing"],"fieldKey":"memory","role":"viewer","durationMinutes":30}`).WithContext(ctx)
			actionResponse(t, fn, r, 500)
			if name == "restart" || name == "start" || name == "stop" {
				r.URL.RawQuery = "dryRun=true"
				actionResponse(t, fn, r, 500)
			}
		})
	}
}

func TestActionBulkGrantBoundaries(t *testing.T) {
	h := newTestHandler(t)
	user := apitest.NewUser(t, h.Store, "viewer")
	viewer := seedCoverageConnector(t, h, "viewer", "networking")
	hidden := seedCoverageConnector(t, h, "hidden", "networking")
	apitest.GrantConnectorRole(t, h.Store, user, viewer.ID, "viewer")
	body, _ := json.Marshal(map[string]any{"ids": []string{viewer.ID, hidden.ID, "missing"}})
	for name, fn := range map[string]http.HandlerFunc{"sync": h.BulkSync, "reauth": h.BulkReauth, "restart": h.BulkRestart} {
		t.Run(name, func(t *testing.T) {
			r := actionRequest("", string(body))
			r = r.WithContext(auth.ContextWithUser(r.Context(), user, true))
			rr := actionResponse(t, fn, r, 200)
			results := bulkResults(t, rr)
			if len(results) != 3 {
				t.Fatalf("results=%+v", results)
			}
			for _, item := range results {
				want := "forbidden"
				if item.ID == "missing" {
					want = "not_found"
				}
				if item.Status != "error" || item.Reason != want || item.JobID != "" {
					t.Fatalf("result=%+v", item)
				}
			}
			actionResponse(t, fn, actionRequest("", `{`), 400)
		})
	}
	rows, _, err := h.Store.ListAuditRecords(context.Background(), "", "connector", "", "", 0, 10)
	if err != nil || len(rows) != 0 {
		t.Fatalf("unauthorized audit=%+v err=%v", rows, err)
	}
}

func TestActionInvalidConnectorConfig(t *testing.T) {
	for _, cfg := range []string{"{", "{}"} {
		h := newTestHandler(t)
		c := seedCoverageConnector(t, h, "invalid", "networking")
		if err := h.Store.UpdateConnector(context.Background(), c.ID, map[string]any{"config_data": cfg, "type": "unregistered-action-type"}); err != nil {
			t.Fatal(err)
		}
		for name, fn := range map[string]http.HandlerFunc{"restart": h.RestartPreview, "start": h.StartPreview, "stop": h.StopPreview, "config push": h.ConfigPush, "config fields": h.ConfigFields} {
			t.Run(name+cfg, func(t *testing.T) { actionResponse(t, fn, actionRequest(c.ID, `{"fieldKey":"memory"}`), 500) })
		}
	}
	h := newTestHandler(t)
	for _, body := range []string{`{`, `{}`, `{"fieldKey":"memory","entityRef":"../bad"}`} {
		actionResponse(t, h.ConfigPush, actionRequest("missing", body), 400)
	}
	actionResponse(t, h.ConfigPush, actionRequest("missing", `{"fieldKey":"memory"}`), 404)
	actionResponse(t, h.Sync, actionRequest("missing", `{}`), 404)
	c := seedCoverageConnector(t, h, "sync", "networking")
	actionResponse(t, h.Sync, actionRequest(c.ID, `{`), 400)
}
