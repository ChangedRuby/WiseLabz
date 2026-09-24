package httpx

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// scripted is a RoundTripper that returns queued results in order.
type scripted struct {
	results []func() (*http.Response, error)
	calls   int
	bodies  []string
}

func (s *scripted) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		s.bodies = append(s.bodies, string(b))
	}
	r := s.results[min(s.calls, len(s.results)-1)]
	s.calls++
	return r()
}

func status(code int) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		return &http.Response{StatusCode: code, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(""))}, nil
	}
}

func fail(err error) func() (*http.Response, error) {
	return func() (*http.Response, error) { return nil, err }
}

var fast = RetryPolicy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}

func do(t *testing.T, rt http.RoundTripper, method string, body io.Reader) (*http.Response, error) {
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), method, "http://example.test/", body)
	if err != nil {
		t.Fatal(err)
	}
	return rt.RoundTrip(req)
}

func TestRetryTransportRetriesTransientStatus(t *testing.T) {
	for _, code := range []int{429, 502, 503, 504} {
		s := &scripted{results: []func() (*http.Response, error){status(code), status(200)}}
		resp, err := do(t, RetryTransport(s, fast), http.MethodGet, nil)
		if err != nil || resp.StatusCode != 200 || s.calls != 2 {
			t.Errorf("status %d: got %v, %v after %d calls; want 200 after 2", code, resp, err, s.calls)
		}
	}
}

func TestRetryTransportGivesUpAfterMaxRetries(t *testing.T) {
	s := &scripted{results: []func() (*http.Response, error){status(503)}}
	resp, err := do(t, RetryTransport(s, fast), http.MethodGet, nil)
	if err != nil || resp.StatusCode != 503 || s.calls != 1+DefaultMaxRetries {
		t.Errorf("got %v, %v after %d calls; want 503 after %d", resp, err, s.calls, 1+DefaultMaxRetries)
	}
}

func TestRetryTransportRetriesTransportErrors(t *testing.T) {
	s := &scripted{results: []func() (*http.Response, error){fail(errors.New("connection reset")), status(200)}}
	resp, err := do(t, RetryTransport(s, fast), http.MethodHead, nil)
	if err != nil || resp.StatusCode != 200 || s.calls != 2 {
		t.Errorf("got %v, %v after %d calls; want 200 after 2", resp, err, s.calls)
	}
}

func TestRetryTransportDoesNotRetry(t *testing.T) {
	cases := map[string]struct {
		method string
		result func() (*http.Response, error)
	}{
		"non-idempotent method": {http.MethodPost, status(503)},
		"client error":          {http.MethodGet, status(404)},
		"server error":          {http.MethodGet, status(500)},
		"deadline":              {http.MethodGet, fail(context.DeadlineExceeded)},
	}
	for name, c := range cases {
		s := &scripted{results: []func() (*http.Response, error){c.result, status(200)}}
		_, _ = do(t, RetryTransport(s, fast), c.method, nil)
		if s.calls != 1 {
			t.Errorf("%s: %d calls, want 1", name, s.calls)
		}
	}
}

func TestRetryTransportDisabled(t *testing.T) {
	s := &scripted{results: []func() (*http.Response, error){status(503), status(200)}}
	resp, _ := do(t, RetryTransport(s, RetryPolicy{MaxRetries: -1}), http.MethodGet, nil)
	if resp.StatusCode != 503 || s.calls != 1 {
		t.Errorf("got %d after %d calls, want 503 after 1", resp.StatusCode, s.calls)
	}
}

func TestRetryTransportReplaysBody(t *testing.T) {
	s := &scripted{results: []func() (*http.Response, error){status(503), status(200)}}
	resp, err := do(t, RetryTransport(s, fast), http.MethodGet, strings.NewReader("payload"))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("got %v, %v", resp, err)
	}
	if len(s.bodies) != 2 || s.bodies[0] != "payload" || s.bodies[1] != "payload" {
		t.Errorf("bodies sent = %q, want the payload twice", s.bodies)
	}
}

func TestRetryTransportStopsWhenContextCanceled(t *testing.T) {
	s := &scripted{results: []func() (*http.Response, error){status(503), status(200)}}
	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.test/", nil)
	rt := RetryTransport(&cancelAfterFirst{next: s, cancel: cancel}, RetryPolicy{BaseDelay: time.Hour, MaxDelay: time.Hour})
	if _, err := rt.RoundTrip(req); err != nil || s.calls != 1 {
		t.Errorf("err = %v after %d calls, want the first response and no retry once canceled", err, s.calls)
	}
}

type cancelAfterFirst struct {
	next   http.RoundTripper
	cancel context.CancelFunc
}

func (c *cancelAfterFirst) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.next.RoundTrip(req)
	c.cancel()
	return resp, err
}

func TestRetryTransportHonorsRetryAfterWithinCap(t *testing.T) {
	rt := RetryTransport(nil, RetryPolicy{BaseDelay: time.Millisecond, MaxDelay: 2 * time.Second}).(*retryTransport)
	resp := &http.Response{Header: http.Header{"Retry-After": []string{"1"}}}
	if got := rt.backoff(0, resp); got != time.Second {
		t.Errorf("backoff with Retry-After: 1 = %v, want 1s", got)
	}
	resp.Header.Set("Retry-After", "60")
	if got := rt.backoff(0, resp); got != 2*time.Second {
		t.Errorf("backoff with Retry-After: 60 = %v, want capped at 2s", got)
	}
}
