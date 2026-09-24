package connector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return false }

func TestIsTimeout(t *testing.T) {
	if !IsTimeout(context.DeadlineExceeded) {
		t.Error("IsTimeout(context.DeadlineExceeded) = false, want true")
	}
	if !IsTimeout(fmt.Errorf("wrapped: %w", context.DeadlineExceeded)) {
		t.Error("IsTimeout(wrapped deadline) = false, want true")
	}
	if !IsTimeout(&net.OpError{Op: "read", Net: "tcp", Err: timeoutError{}}) {
		t.Error("IsTimeout(net timeout) = false, want true")
	}
	if IsTimeout(errors.New("connection refused")) {
		t.Error("IsTimeout(plain error) = true, want false")
	}
}

func TestMapTransportError(t *testing.T) {
	var te *TimeoutError
	if err := MapTransportError(context.DeadlineExceeded); !errors.As(err, &te) {
		t.Errorf("MapTransportError(deadline) = %v, want TimeoutError", err)
	}
	if err := MapTransportError(errors.New("refused")); errors.As(err, &te) {
		t.Errorf("MapTransportError(refused) = %v, want plain error", err)
	}
}

func TestCheckStatus(t *testing.T) {
	var authErr *AuthError
	var unavailable *ServiceUnavailableError
	tests := []struct {
		status int
		check  func(error) bool
	}{
		{200, func(err error) bool { return err == nil }},
		{304, func(err error) bool { return err == nil }},
		{401, func(err error) bool { return errors.As(err, &authErr) }},
		{403, func(err error) bool { return errors.As(err, &authErr) }},
		{502, func(err error) bool { return errors.As(err, &unavailable) }},
		{503, func(err error) bool { return errors.As(err, &unavailable) }},
		{504, func(err error) bool { return errors.As(err, &unavailable) }},
		{404, func(err error) bool {
			return err != nil && !errors.As(err, &authErr) && !errors.As(err, &unavailable)
		}},
		{500, func(err error) bool {
			return err != nil && !errors.As(err, &authErr) && !errors.As(err, &unavailable)
		}},
	}
	for _, tt := range tests {
		if err := CheckStatus(tt.status, []byte("body")); !tt.check(err) {
			t.Errorf("CheckStatus(%d) = %v", tt.status, err)
		}
	}
}

func TestNewHTTPClientRetriesAndBlocksRedirects(t *testing.T) {
	AllowLoopbackForTest(t)
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/flaky":
			if calls.Add(1) == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		case "/redirect":
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
		}
	}))
	defer srv.Close()

	client := NewHTTPClient(HTTPClientOptions{})
	resp, err := client.Get(srv.URL + "/flaky")
	if err != nil {
		t.Fatalf("GET /flaky: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK || calls.Load() != 2 {
		t.Errorf("GET /flaky = %d after %d calls, want 200 after 2", resp.StatusCode, calls.Load())
	}

	resp, err = client.Get(srv.URL + "/redirect")
	if err != nil {
		t.Fatalf("GET /redirect: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("GET /redirect = %d, want 302 (redirect not followed)", resp.StatusCode)
	}
}

func TestNewHTTPClientBlocksLoopbackByDefault(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	defer srv.Close()
	if _, err := NewHTTPClient(HTTPClientOptions{}).Get(srv.URL); err == nil {
		t.Error("GET loopback succeeded, want the guarded dialer to refuse it")
	}
}
