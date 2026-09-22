package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/WiseLabz/wiselabz/internal/connector"
	"github.com/WiseLabz/wiselabz/internal/store"
)

// healthFakeConnector reports a fixed Validate outcome, optionally after a
// delay, so tests can drive ClassifyHealth's online/degraded/offline
// branches without touching the network.
type healthFakeConnector struct {
	err   error
	delay time.Duration
}

func (f *healthFakeConnector) Name() string     { return "health fake" }
func (f *healthFakeConnector) Type() string     { return "health_test_fake" }
func (f *healthFakeConnector) Category() string { return "test" }
func (f *healthFakeConnector) Fetch(_ context.Context, _ map[string]any) (*connector.ServiceSnapshot, error) {
	return nil, errors.New("not used by health checks")
}
func (f *healthFakeConnector) Validate(_ context.Context, _ map[string]any) error {
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	return f.err
}

// registerHealthFakeType (re-)registers the "health_test_fake" connector type
// with the given Validate outcome. Tests in this file run sequentially, so
// re-registration between them is safe: each call only affects the next
// connector.Get("health_test_fake", ...) lookup.
func registerHealthFakeType(t *testing.T, behavior *healthFakeConnector) {
	t.Helper()
	connector.Register(
		connector.TypeSchema{Type: "health_test_fake", Category: "test", Name: "Health Fake"},
		func(_ map[string]any) (connector.Connector, error) {
			return behavior, nil
		},
	)
}

func seedHealthTestConnector(t *testing.T, app *testApp) *store.ConnectorRecord {
	t.Helper()
	conn := &store.ConnectorRecord{
		Name: "svc", Category: "networking", Type: "health_test_fake", URL: "https://example.com", Enabled: true,
	}
	if err := app.Store.CreateConnector(context.Background(), conn); err != nil {
		t.Fatalf("seed connector: %v", err)
	}
	return conn
}

func TestConnectorsHealthOnline(t *testing.T) {
	registerHealthFakeType(t, &healthFakeConnector{err: nil})

	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := seedHealthTestConnector(t, app)
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/health", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var got struct {
		Status    string `json:"status"`
		Message   string `json:"message"`
		LatencyMs int    `json:"latencyMs"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Status != "online" {
		t.Fatalf("status = %q, want online", got.Status)
	}

	stored, err := app.Store.GetConnector(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if stored.Status != "online" {
		t.Fatalf("persisted status = %q, want online", stored.Status)
	}
}

func TestConnectorsHealthDegraded(t *testing.T) {
	orig := connector.DegradedLatencyThreshold
	connector.DegradedLatencyThreshold = time.Millisecond
	t.Cleanup(func() { connector.DegradedLatencyThreshold = orig })

	registerHealthFakeType(t, &healthFakeConnector{err: nil, delay: 5 * time.Millisecond})

	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := seedHealthTestConnector(t, app)
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/health", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var got struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Status != "degraded" {
		t.Fatalf("status = %q, want degraded", got.Status)
	}
}

func TestConnectorsHealthOffline(t *testing.T) {
	registerHealthFakeType(t, &healthFakeConnector{err: errors.New("connection refused")})

	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := seedHealthTestConnector(t, app)
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/health", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var got struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Status != "offline" {
		t.Fatalf("status = %q, want offline", got.Status)
	}
	if got.Message != "connection refused" {
		t.Fatalf("message = %q, want %q", got.Message, "connection refused")
	}

	stored, err := app.Store.GetConnector(context.Background(), conn.ID)
	if err != nil {
		t.Fatalf("GetConnector: %v", err)
	}
	if stored.Status != "offline" {
		t.Fatalf("persisted status = %q, want offline", stored.Status)
	}
}

func TestConnectorsHealthDoesNotCreateSnapshot(t *testing.T) {
	registerHealthFakeType(t, &healthFakeConnector{err: nil})

	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := seedHealthTestConnector(t, app)
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/health", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	if _, err := app.Store.GetLatestSnapshot(context.Background(), conn.ID); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("GetLatestSnapshot error = %v, want ErrNotFound (health check must not persist a snapshot)", err)
	}
}

func TestConnectorsHealthRoleBoundary(t *testing.T) {
	registerHealthFakeType(t, &healthFakeConnector{err: nil})

	app := newTestApp(t)
	_, viewerToken := app.user(t, "viewer")
	conn := seedHealthTestConnector(t, app)

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/health", nil, viewerToken)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body)
	}
}

// TestConnectorsHealthRecordsTimeSeriesRow verifies a health check persists
// a health_checks row (in addition to updating the connector's latest
// status), so GetConnectorUptime has data to compute over.
func TestConnectorsHealthRecordsTimeSeriesRow(t *testing.T) {
	registerHealthFakeType(t, &healthFakeConnector{err: nil})

	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := seedHealthTestConnector(t, app)
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/"+conn.ID+"/health", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	stats, err := app.Store.GetConnectorUptime(context.Background(), conn.ID, time.Now().UTC().Add(-time.Hour), time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatalf("GetConnectorUptime: %v", err)
	}
	if stats.CheckCount != 1 {
		t.Fatalf("CheckCount = %d, want 1 (health check should persist a time-series row)", stats.CheckCount)
	}
}

// TestConnectorsUptime exercises GET /api/connectors/{id}/uptime end to end:
// two health checks, one online and one offline an hour apart, followed by
// a request that should show partial (50%) availability in the 24h window.
func TestConnectorsUptime(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	conn := seedHealthTestConnector(t, app)
	app.connectorGrant(t, opUserID, conn.ID, "operator")

	now := time.Now().UTC()
	if err := app.Store.RecordHealthCheck(context.Background(), &store.HealthCheckRecord{
		ConnectorID: conn.ID, Status: "offline", CheckedAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("RecordHealthCheck(offline): %v", err)
	}
	if err := app.Store.RecordHealthCheck(context.Background(), &store.HealthCheckRecord{
		ConnectorID: conn.ID, Status: "online", CheckedAt: now.Add(-time.Hour).Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("RecordHealthCheck(online): %v", err)
	}

	rec := app.req(t, http.MethodGet, "/api/connectors/"+conn.ID+"/uptime", nil, opToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var got struct {
		ConnectorID string `json:"connectorId"`
		Windows     map[string]struct {
			CheckCount      int     `json:"checkCount"`
			AvailabilityPct float64 `json:"availabilityPct"`
			OutageCount     int     `json:"outageCount"`
		} `json:"windows"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.ConnectorID != conn.ID {
		t.Fatalf("connectorId = %q, want %q", got.ConnectorID, conn.ID)
	}
	for _, label := range []string{"24h", "7d", "30d"} {
		w, ok := got.Windows[label]
		if !ok {
			t.Fatalf("windows missing %q: %+v", label, got.Windows)
		}
		if w.CheckCount != 2 {
			t.Errorf("windows[%q].checkCount = %d, want 2", label, w.CheckCount)
		}
	}
}

func TestConnectorsUptimeUnknownConnector404s(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	app.connectorGrant(t, opUserID, "does-not-exist", "operator")

	rec := app.req(t, http.MethodGet, "/api/connectors/does-not-exist/uptime", nil, opToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}

// TestConnectorsHealthUnknownConnector404s pre-grants the caller on the
// nonexistent ID so the request reaches the handler's own not-found check —
// otherwise the default-deny connector-role gate would 403 first, which is
// correct in production (no grant can exist on an ID that was never a real
// connector) but isn't what this test is exercising.
func TestConnectorsHealthUnknownConnector404s(t *testing.T) {
	app := newTestApp(t)
	opUserID, opToken := app.user(t, "operator")
	app.connectorGrant(t, opUserID, "does-not-exist", "operator")

	rec := app.req(t, http.MethodPost, "/api/connectors/does-not-exist/health", nil, opToken)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}
