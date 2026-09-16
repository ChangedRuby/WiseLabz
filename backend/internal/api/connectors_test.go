package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

func TestConnectorsListSuccess(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/connectors", nil, viewerToken)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var got []store.ConnectorRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

func TestConnectorsSchemaMarksStubTypes(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodGet, "/api/connectors/schema", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var schemas []struct {
		Type string `json:"type"`
		Stub bool   `json:"stub"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &schemas); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	byType := make(map[string]bool)
	for _, s := range schemas {
		byType[s.Type] = s.Stub
	}

	if stub, ok := byType["docker"]; !ok || stub {
		t.Errorf("docker.stub = %v (present=%v), want false", stub, ok)
	}
	if stub, ok := byType["proxmox"]; !ok || stub {
		t.Errorf("proxmox.stub = %v (present=%v), want false", stub, ok)
	}
}

func TestConnectorsCreateRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodPost, "/api/connectors", map[string]any{
		"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com",
	}, viewerToken)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestConnectorsCreateValidation(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	tests := []struct {
		name string
		body any
	}{
		{"malformed json", "not-json"},
		{"missing required fields", map[string]any{"name": "svc"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/connectors", tt.body, opToken)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestConnectorsCreateSuccess(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors", map[string]any{
		"name": "svc", "category": "virtualization", "type": "proxmox", "url": "https://example.com", "owner": "Platform",
	}, opToken)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body)
	}
	var created store.ConnectorRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created connector to have an ID")
	}

	got, err := app.Store.GetConnector(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if got.Name != "svc" || got.Owner != "Platform" {
		t.Errorf("connector = %+v, want name svc and owner Platform", got)
	}
}

func TestConnectorsUpdateOwner(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	c := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com", Owner: "Platform"}
	if err := app.Store.CreateConnector(context.Background(), c); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	rec := app.req(t, http.MethodPut, "/api/connectors/"+c.ID, map[string]any{"owner": ""}, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	got, err := app.Store.GetConnector(context.Background(), c.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if got.Owner != "" {
		t.Fatalf("Owner = %q, want blank", got.Owner)
	}
}

func TestConnectorsUpdateRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	rec := app.req(t, http.MethodPut, "/api/connectors/some-id", map[string]any{"name": "x"}, viewerToken)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

func TestConnectorsUpdateScheduleSeconds(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}

	t.Run("absent leaves schedule unchanged", func(t *testing.T) {
		rec := app.req(t, http.MethodPatch, "/api/connectors/"+conn.ID, map[string]any{"name": "svc2"}, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}
		got, err := app.Store.GetConnector(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("GetConnector: %v", err)
		}
		if got.ScheduleSeconds != nil {
			t.Fatalf("ScheduleSeconds = %v, want still nil (field absent from request)", got.ScheduleSeconds)
		}
	})

	t.Run("sets schedule", func(t *testing.T) {
		rec := app.req(t, http.MethodPut, "/api/connectors/"+conn.ID, map[string]any{"scheduleSeconds": 1800}, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}
		var got store.ConnectorRecord
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.ScheduleSeconds == nil || *got.ScheduleSeconds != 1800 {
			t.Fatalf("ScheduleSeconds = %v, want 1800", got.ScheduleSeconds)
		}
	})

	t.Run("explicit null clears schedule to manual-only", func(t *testing.T) {
		rec := app.req(t, http.MethodPut, "/api/connectors/"+conn.ID, map[string]any{"scheduleSeconds": nil}, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}
		got, err := app.Store.GetConnector(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("GetConnector: %v", err)
		}
		if got.ScheduleSeconds != nil {
			t.Fatalf("ScheduleSeconds = %v, want nil after explicit null", got.ScheduleSeconds)
		}
	})
}

func TestConnectorsSyncsHistory(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")

	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}

	older := time.Now().Add(-time.Hour).Format(time.RFC3339)
	newer := time.Now().Format(time.RFC3339)
	if err := app.Store.CreateSyncRun(context.Background(), &store.SyncRunRecord{
		ConnectorID: conn.ID, StartedAt: older, Status: store.SyncRunStatusError, Error: "boom", Attempt: 1,
	}); err != nil {
		t.Fatalf("seed sync run 1: %v", err)
	}
	if err := app.Store.CreateSyncRun(context.Background(), &store.SyncRunRecord{
		ConnectorID: conn.ID, StartedAt: newer, Status: store.SyncRunStatusSuccess, Attempt: 2, ChangesCount: 2,
	}); err != nil {
		t.Fatalf("seed sync run 2: %v", err)
	}

	rec := app.req(t, http.MethodGet, "/api/connectors/"+conn.ID+"/syncs", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var got []store.SyncRunRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Status != store.SyncRunStatusSuccess || got[1].Status != store.SyncRunStatusError {
		t.Fatalf("got = %+v, want newest (success) first", got)
	}

	t.Run("limit query param is respected", func(t *testing.T) {
		rec := app.req(t, http.MethodGet, "/api/connectors/"+conn.ID+"/syncs?limit=1", nil, viewerToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}
		var got []store.SyncRunRecord
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
	})

	t.Run("unknown connector 404s", func(t *testing.T) {
		rec := app.req(t, http.MethodGet, "/api/connectors/does-not-exist/syncs", nil, viewerToken)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestConnectorRestartPreview(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")

	conn := &store.ConnectorRecord{
		Name: "svc", Category: "virtualization", Type: "unregistered", URL: "https://example.com",
	}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	snapshot, err := json.Marshal(connector.ServiceSnapshot{
		ServiceName: "service-a",
		Dependencies: []connector.ServiceDependency{
			{Kind: "host", Name: "node-a", Ref: "host-1"},
			{Kind: "upstream_service", Name: "database"},
		},
	})
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := app.Store.CreateSnapshot(context.Background(), &store.SnapshotRecord{
		ConnectorID: conn.ID,
		Data:        string(snapshot),
	}); err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	t.Run("preview uses latest stored snapshot", func(t *testing.T) {
		beforeConnector, err := app.Store.GetConnector(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("get connector before preview: %v", err)
		}
		beforeSnapshots, err := app.Store.CountSnapshotsByConnector(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("count snapshots before preview: %v", err)
		}
		_, beforeAudits, err := app.Store.ListAuditRecords(context.Background(), "", "", "", "", 0, 100)
		if err != nil {
			t.Fatalf("count audits before preview: %v", err)
		}

		rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/restart?dryRun=true", nil, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}

		var got struct {
			TargetService            string `json:"targetService"`
			EstimatedDowntimeSeconds int    `json:"estimatedDowntimeSeconds"`
			DependentServices        []struct {
				Kind string `json:"kind"`
				Name string `json:"name"`
				Ref  string `json:"ref"`
			} `json:"dependentServices"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.TargetService != "service-a" || got.EstimatedDowntimeSeconds != 30 {
			t.Fatalf("preview = %+v, want service-a and 30 seconds", got)
		}
		if len(got.DependentServices) != 2 || got.DependentServices[0].Ref != "host-1" {
			t.Fatalf("dependent services = %+v, want stored dependencies", got.DependentServices)
		}
		if got.DependentServices[1].Ref != "" {
			t.Fatalf("empty dependency ref = %q, want empty", got.DependentServices[1].Ref)
		}

		var shape struct {
			DependentServices []map[string]json.RawMessage `json:"dependentServices"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &shape); err != nil {
			t.Fatalf("decode response shape: %v", err)
		}
		first := shape.DependentServices[0]
		firstHasExpectedFields := len(first) == 3 && first["kind"] != nil && first["name"] != nil && first["ref"] != nil
		if !firstHasExpectedFields {
			t.Fatalf("first dependency fields = %v, want exactly kind, name, ref", first)
		}
		second := shape.DependentServices[1]
		secondHasExpectedFields := len(second) == 2 && second["kind"] != nil && second["name"] != nil && second["ref"] == nil
		if !secondHasExpectedFields {
			t.Fatalf("second dependency fields = %v, want kind and name with unknown ref omitted", second)
		}

		afterConnector, err := app.Store.GetConnector(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("get connector after preview: %v", err)
		}
		if !reflect.DeepEqual(afterConnector, beforeConnector) {
			t.Fatalf("connector changed during preview: before=%+v after=%+v", beforeConnector, afterConnector)
		}
		afterSnapshots, err := app.Store.CountSnapshotsByConnector(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("count snapshots after preview: %v", err)
		}
		if afterSnapshots != beforeSnapshots {
			t.Fatalf("snapshot count = %d after preview, want %d", afterSnapshots, beforeSnapshots)
		}
		_, afterAudits, err := app.Store.ListAuditRecords(context.Background(), "", "", "", "", 0, 100)
		if err != nil {
			t.Fatalf("count audits after preview: %v", err)
		}
		if afterAudits != beforeAudits {
			t.Fatalf("audit count = %d after preview, want %d", afterAudits, beforeAudits)
		}
	})

	t.Run("empty dependencies encode as an empty array", func(t *testing.T) {
		emptySnapshot, err := json.Marshal(connector.ServiceSnapshot{ServiceName: "service-b"})
		if err != nil {
			t.Fatalf("marshal empty snapshot: %v", err)
		}
		if err := app.Store.CreateSnapshot(context.Background(), &store.SnapshotRecord{
			ConnectorID: conn.ID,
			Data:        string(emptySnapshot),
			FetchedAt:   time.Now().UTC().Add(time.Minute).Format(time.RFC3339),
		}); err != nil {
			t.Fatalf("seed empty snapshot: %v", err)
		}

		rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/restart?dryRun=true", nil, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}
		var got struct {
			TargetService     string            `json:"targetService"`
			DependentServices []json.RawMessage `json:"dependentServices"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got.TargetService != "service-b" || got.DependentServices == nil || len(got.DependentServices) != 0 {
			t.Fatalf("preview = %+v, want service-b with a non-null empty dependency array", got)
		}
	})

	t.Run("viewer is forbidden", func(t *testing.T) {
		rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/restart?dryRun=true", nil, viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestConnectorRestartNonDryRunOnMissingConnectorReturnsNotFound(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	// Any value other than an exact "dryRun=true" now takes the real,
	// mutating restart path (gated by elevation) instead of the preview.
	for _, suffix := range []string{"", "?dryRun=false", "?dryRun=", "?dryRun=True", "?dryRun=true&dryRun=false"} {
		t.Run(suffix, func(t *testing.T) {
			rec := app.req(t, http.MethodPost, "/api/connectors/unknown/restart"+suffix, nil, opToken)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
			}
		})
	}
}

func TestConnectorRestartPreviewSnapshotErrors(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}

	t.Run("missing snapshot", func(t *testing.T) {
		rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/restart?dryRun=true", nil, opToken)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("unknown connector", func(t *testing.T) {
		rec := app.req(t, http.MethodPost, "/api/connectors/unknown/restart?dryRun=true", nil, opToken)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
		}
	})

	if err := app.Store.CreateSnapshot(context.Background(), &store.SnapshotRecord{
		ConnectorID: conn.ID,
		Data:        "not-json",
	}); err != nil {
		t.Fatalf("seed malformed snapshot: %v", err)
	}
	t.Run("malformed snapshot", func(t *testing.T) {
		rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/restart?dryRun=true", nil, opToken)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestConnectorsDeleteElevationBoundary(t *testing.T) {
	app := newTestApp(t)
	opID, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")

	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}

	t.Run("viewer forbidden", func(t *testing.T) {
		rec := app.req(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("operator missing elevation token", func(t *testing.T) {
		rec := app.req(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, opToken)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("operator invalid elevation token", func(t *testing.T) {
		req := app.newRequest(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, opToken)
		req.Header.Set("X-Elevation-Token", "garbage")
		rec := app.serve(req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("operator elevation token scoped to wrong action", func(t *testing.T) {
		tok := app.elevationToken(t, opID, "user.delete")
		req := app.newRequest(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, opToken)
		req.Header.Set("X-Elevation-Token", tok)
		rec := app.serve(req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("operator valid elevation token", func(t *testing.T) {
		tok := app.elevationToken(t, opID, "connector.delete")
		req := app.newRequest(t, http.MethodDelete, "/api/connectors/"+conn.ID, nil, opToken)
		req.Header.Set("X-Elevation-Token", tok)
		rec := app.serve(req)
		if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 200 or 204; body = %s", rec.Code, rec.Body)
		}

		if _, err := app.Store.GetConnector(context.Background(), conn.ID); err == nil {
			t.Fatal("expected connector to be deleted")
		}
	})
}
