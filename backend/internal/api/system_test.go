package api_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthReportsHealthy(t *testing.T) {
	app := newTestApp(t)
	rec := app.req(t, http.MethodGet, "/api/health", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body)
	}

	var body struct {
		Healthy bool `json:"healthy"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.Healthy {
		t.Error("healthy = false, want true")
	}
}

func TestHealthProbeRoutes(t *testing.T) {
	app := newTestApp(t)
	for _, path := range []string{"/healthz", "/readyz"} {
		rec := app.req(t, http.MethodGet, path, nil, "")
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200; body = %s", path, rec.Code, rec.Body)
		}
	}
}

func TestVersionRouteIsNotExposed(t *testing.T) {
	app := newTestApp(t)
	rec := app.req(t, http.MethodGet, "/api/version", nil, "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body)
	}
}
