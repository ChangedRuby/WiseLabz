package httpx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

// RetryPolicy configures RetryTransport. The zero value is valid.
type RetryPolicy struct {
	// MaxRetries is the number of retries after the first attempt. Zero
	// means DefaultMaxRetries; negative disables retries.
	MaxRetries int
	// BaseDelay is the first backoff delay, doubled on each retry. Zero
	// means DefaultBaseDelay.
	BaseDelay time.Duration
	// MaxDelay caps a single backoff delay, including one taken from a
	// Retry-After header. Zero means DefaultMaxDelay.
	MaxDelay time.Duration
}

// Retry defaults, sized for homelab APIs that briefly return 502/503 while
// a service restarts: two retries finish within about a second.
const (
	DefaultMaxRetries = 2
	DefaultBaseDelay  = 250 * time.Millisecond
	DefaultMaxDelay   = 5 * time.Second
)

// RetryTransport wraps next so idempotent requests (GET, HEAD, OPTIONS) are
// retried with exponential backoff on transport errors and on 429, 502, 503
// and 504 responses. Other methods are never retried: a POST that timed out
// may already have taken effect upstream. The request context bounds the
// whole sequence, and a canceled or expired context is never retried.
func RetryTransport(next http.RoundTripper, p RetryPolicy) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	if p.MaxRetries == 0 {
		p.MaxRetries = DefaultMaxRetries
	}
	if p.BaseDelay <= 0 {
		p.BaseDelay = DefaultBaseDelay
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = DefaultMaxDelay
	}
	return &retryTransport{next: next, policy: p}
}

type retryTransport struct {
	next   http.RoundTripper
	policy RetryPolicy
}

// Unwrap returns the RoundTripper wrapped by RetryTransport, or rt itself
// when it isn't one. It exists so tests can inspect the underlying
// transport (e.g. its TLS config) without reaching into unexported state.
func Unwrap(rt http.RoundTripper) http.RoundTripper {
	if r, ok := rt.(*retryTransport); ok {
		return r.next
	}
	return rt
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.policy.MaxRetries < 0 || !idempotent(req.Method) || (req.Body != nil && req.Body != http.NoBody && req.GetBody == nil) {
		return t.next.RoundTrip(req)
	}
	for attempt := 0; ; attempt++ {
		resp, err := t.next.RoundTrip(req)
		if attempt >= t.policy.MaxRetries || !retryable(req.Context(), resp, err) {
			return resp, err
		}
		delay := t.backoff(attempt, resp)
		if resp != nil {
			// Drain a little so the connection can be reused, then close.
			_, _ = io.CopyN(io.Discard, resp.Body, 4<<10)
			_ = resp.Body.Close()
		}
		if err := sleep(req.Context(), delay); err != nil {
			return nil, err
		}
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			req = req.Clone(req.Context())
			req.Body = body
		}
	}
}

func (t *retryTransport) backoff(attempt int, resp *http.Response) time.Duration {
	delay := t.policy.BaseDelay << attempt
	if resp != nil {
		if secs, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && secs >= 0 {
			delay = time.Duration(secs) * time.Second
		}
	}
	return min(delay, t.policy.MaxDelay)
}

func idempotent(method string) bool {
	return method == "" || IsSafeMethod(method)
}

func retryable(ctx context.Context, resp *http.Response, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	if err != nil {
		return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
