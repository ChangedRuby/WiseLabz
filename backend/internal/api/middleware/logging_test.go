package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func TestLoggerRedactsShareToken(t *testing.T) {
	buf := captureLog(t)
	h := Logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))

	for _, path := range []string{"/api/share/SECRETTOKEN123/tree", "/api/share/SECRETTOKEN123/docs/d1", "/api/share/SECRETTOKEN123"} {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))
	}
	out := buf.String()
	if strings.Contains(out, "SECRETTOKEN123") {
		t.Fatalf("log output contains share token: %s", out)
	}
	if !strings.Contains(out, "/api/share/{token}/docs/d1") {
		t.Errorf("expected masked route in log, got: %s", out)
	}
}

func TestLoggerRedactsWSTicket(t *testing.T) {
	buf := captureLog(t)
	h := Logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/ws?ticket=ONETIMEVALUE", nil))
	if strings.Contains(buf.String(), "ONETIMEVALUE") {
		t.Fatalf("log output contains ws ticket: %s", buf.String())
	}
}
