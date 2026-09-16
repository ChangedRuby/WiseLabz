package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/WiseLabz/wiselabz/internal/store"
)

func seedMaintenanceConnector(t *testing.T, app *testApp) *store.ConnectorRecord {
	t.Helper()
	conn := &store.ConnectorRecord{Name: "svc", Category: "virtualization", Type: "proxmox", URL: "https://example.com"}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	return conn
}

func TestOpenMaintenanceWindowRoleBoundary(t *testing.T) {
	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	conn := seedMaintenanceConnector(t, app)

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/maintenance-window",
		map[string]any{"durationMinutes": 30}, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

// TestOpenMaintenanceWindowNoElevationRequired verifies that, unlike
// connector.delete/start/stop, opening a maintenance window succeeds for an
// operator with no X-Elevation-Token at all — it's reversible and time-boxed.
func TestOpenMaintenanceWindowNoElevationRequired(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	conn := seedMaintenanceConnector(t, app)

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/maintenance-window",
		map[string]any{"durationMinutes": 30}, opToken)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 (no elevation should be required); body = %s", rec.Code, rec.Body)
	}

	var got store.MaintenanceWindowRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.ConnectorID != conn.ID {
		t.Fatalf("ConnectorID = %q, want %q", got.ConnectorID, conn.ID)
	}
	if got.EndsAt <= got.StartsAt {
		t.Fatalf("EndsAt %q should be after StartsAt %q", got.EndsAt, got.StartsAt)
	}

	records, _, err := app.Store.ListAuditRecords(context.Background(), "connector.maintenanceWindow.open", "connector", "", "", 0, 10)
	if err != nil {
		t.Fatalf("ListAuditRecords() error: %v", err)
	}
	if len(records) != 1 || records[0].TargetID != conn.ID {
		t.Fatalf("audit records = %+v, want one row for connector %s", records, conn.ID)
	}
}

func TestOpenMaintenanceWindowInvalidDuration(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	conn := seedMaintenanceConnector(t, app)

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/maintenance-window",
		map[string]any{"durationMinutes": 0}, opToken)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body)
	}
}

func TestOpenMaintenanceWindowConnectorNotFound(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/does-not-exist/maintenance-window",
		map[string]any{"durationMinutes": 30}, opToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

func TestGetMaintenanceWindowAnyAuthenticatedUser(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")
	conn := seedMaintenanceConnector(t, app)

	rec := app.req(t, http.MethodGet, "/api/connectors/"+conn.ID+"/maintenance-window", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	if strings.TrimSpace(rec.Body.String()) != "null" {
		t.Fatalf("body = %s, want null (no window yet)", rec.Body)
	}

	openRec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/maintenance-window",
		map[string]any{"durationMinutes": 30}, opToken)
	if openRec.Code != http.StatusCreated {
		t.Fatalf("open status = %d, want 201; body = %s", openRec.Code, openRec.Body)
	}

	rec = app.req(t, http.MethodGet, "/api/connectors/"+conn.ID+"/maintenance-window", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var got store.MaintenanceWindowRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.ConnectorID != conn.ID {
		t.Fatalf("ConnectorID = %q, want %q", got.ConnectorID, conn.ID)
	}
}

func TestCloseMaintenanceWindowRoleBoundaryAndNoElevation(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")
	conn := seedMaintenanceConnector(t, app)

	openRec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/maintenance-window",
		map[string]any{"durationMinutes": 60}, opToken)
	if openRec.Code != http.StatusCreated {
		t.Fatalf("open status = %d, want 201; body = %s", openRec.Code, openRec.Body)
	}

	t.Run("viewer forbidden", func(t *testing.T) {
		rec := app.req(t, http.MethodDelete, "/api/connectors/"+conn.ID+"/maintenance-window", nil, viewerToken)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
		}
	})

	t.Run("operator closes with no elevation token", func(t *testing.T) {
		rec := app.req(t, http.MethodDelete, "/api/connectors/"+conn.ID+"/maintenance-window", nil, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 (no elevation should be required); body = %s", rec.Code, rec.Body)
		}

		active, err := app.Store.GetActiveMaintenanceWindow(context.Background(), conn.ID)
		if err != nil {
			t.Fatalf("GetActiveMaintenanceWindow() error: %v", err)
		}
		if active != nil {
			t.Fatalf("GetActiveMaintenanceWindow() = %+v, want nil after close", active)
		}

		records, _, err := app.Store.ListAuditRecords(context.Background(), "connector.maintenanceWindow.close", "connector", "", "", 0, 10)
		if err != nil {
			t.Fatalf("ListAuditRecords() error: %v", err)
		}
		if len(records) != 1 || records[0].TargetID != conn.ID {
			t.Fatalf("audit records = %+v, want one row for connector %s", records, conn.ID)
		}
	})

	t.Run("closing again is a no-op", func(t *testing.T) {
		rec := app.req(t, http.MethodDelete, "/api/connectors/"+conn.ID+"/maintenance-window", nil, opToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
		}
	})
}

func TestListActiveMaintenanceWindowsEndpoint(t *testing.T) {
	app := newTestApp(t)
	_, opToken := app.user(t, "operator")
	_, viewerToken := app.user(t, "viewer")
	connA := seedMaintenanceConnector(t, app)

	rec := app.req(t, http.MethodGet, "/api/connectors/maintenance-windows", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	var got []store.MaintenanceWindowRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("len(got) = %d, want 0", len(got))
	}

	openRec := app.req(t, http.MethodPost, "/api/connectors/"+connA.ID+"/maintenance-window",
		map[string]any{"durationMinutes": 30}, opToken)
	if openRec.Code != http.StatusCreated {
		t.Fatalf("open status = %d, want 201; body = %s", openRec.Code, openRec.Body)
	}

	rec = app.req(t, http.MethodGet, "/api/connectors/maintenance-windows", nil, viewerToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(got) != 1 || got[0].ConnectorID != connA.ID {
		t.Fatalf("got = %+v, want exactly one window for connector %s", got, connA.ID)
	}
}
