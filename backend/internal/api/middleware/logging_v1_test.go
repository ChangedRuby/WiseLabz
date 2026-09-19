package middleware

import (
	"net/http/httptest"
	"testing"
)

func TestLoggablePathMasksShareTokenUnderV1(t *testing.T) {
	for in, want := range map[string]string{
		"/api/share/secret/tree":    "/api/share/{token}/tree",
		"/api/v1/share/secret/tree": "/api/v1/share/{token}/tree",
		"/api/v1/share/secret":      "/api/v1/share/{token}",
	} {
		if got := loggablePath(httptest.NewRequest("GET", in, nil)); got != want {
			t.Errorf("loggablePath(%q) = %q, want %q", in, got, want)
		}
	}
}
