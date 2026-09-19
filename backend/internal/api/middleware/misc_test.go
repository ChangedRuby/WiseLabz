package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDGeneratesAndPropagates(t *testing.T) {
	var seen string
	h := RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { seen = GetRequestID(r.Context()) }))

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if seen == "" || rr.Header().Get("X-Request-ID") != seen {
		t.Fatalf("generated id = %q, header = %q", seen, rr.Header().Get("X-Request-ID"))
	}

	rr = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "client-id")
	h.ServeHTTP(rr, req)
	if seen != "client-id" || rr.Header().Get("X-Request-ID") != "client-id" {
		t.Fatalf("client id not honoured: ctx=%q header=%q", seen, rr.Header().Get("X-Request-ID"))
	}
}

func TestGetRequestIDMissing(t *testing.T) {
	if got := GetRequestID(httptest.NewRequest(http.MethodGet, "/", nil).Context()); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}

func TestRecovererReturns500OnPanic(t *testing.T) {
	buf := captureLog(t)
	h := Recoverer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/x", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "internal_error") {
		t.Errorf("body = %s", rr.Body.String())
	}
	if strings.Contains(rr.Body.String(), "boom") {
		t.Errorf("panic value leaked to client: %s", rr.Body.String())
	}
	if !strings.Contains(buf.String(), "panic recovered") {
		t.Errorf("panic not logged: %s", buf.String())
	}
}

func TestRecovererPassThrough(t *testing.T) {
	h := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) }))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusTeapot {
		t.Fatalf("status = %d", rr.Code)
	}
}

func TestRateLimit(t *testing.T) {
	// rate 0 => no refill, so the burst is deterministic regardless of timing.
	keyFn := func(r *http.Request) string { return r.Header.Get("X-Key") }
	h := RateLimit(0, 2, keyFn)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))

	do := func(key string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Key", key)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		return rr
	}

	for i := 0; i < 2; i++ {
		if rr := do("a"); rr.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i, rr.Code)
		}
	}
	rr := do("a")
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("over-limit status = %d, want 429", rr.Code)
	}
	if rr.Header().Get("Retry-After") == "" {
		t.Error("missing Retry-After")
	}
	if !strings.Contains(rr.Body.String(), "rate_limited") {
		t.Errorf("body = %s", rr.Body.String())
	}
	// Buckets are per key.
	if rr := do("b"); rr.Code != http.StatusOK {
		t.Fatalf("other key status = %d, want 200", rr.Code)
	}
	// Empty key is rejected.
	if rr := do(""); rr.Code != http.StatusTooManyRequests {
		t.Fatalf("empty key status = %d, want 429", rr.Code)
	}
}
