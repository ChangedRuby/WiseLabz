package connector

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/WiseLabz/wiselabz/internal/httpx"
)

// DefaultHTTPTimeout bounds a whole connector request, retries included.
const DefaultHTTPTimeout = 30 * time.Second

// HTTPClientOptions configures NewHTTPClient. The zero value verifies TLS.
type HTTPClientOptions struct {
	// SkipTLSVerify disables certificate verification, for services with
	// self-signed certificates (the connector's "Verify TLS" toggle off).
	SkipTLSVerify bool
	// Jar, when set, stores cookies between requests (session-auth APIs).
	Jar http.CookieJar
}

// NewHTTPClient returns the shared client every HTTP connector uses (#265):
// httpx's hardened transport (TLS 1.2+, bounded timeouts, no redirects)
// dialing through GuardedDialer so loopback and link-local targets stay
// blocked, with idempotent requests retried on transient upstream errors
// (see httpx.RetryTransport).
func NewHTTPClient(o HTTPClientOptions) *http.Client {
	transport := httpx.NewTransport(httpx.Options{
		Timeout:            DefaultHTTPTimeout,
		InsecureSkipVerify: o.SkipTLSVerify,
		DialContext:        GuardedDialer(DefaultHTTPTimeout).DialContext,
	})
	return &http.Client{
		Timeout:       DefaultHTTPTimeout,
		Transport:     httpx.RetryTransport(transport, httpx.RetryPolicy{}),
		Jar:           o.Jar,
		CheckRedirect: httpx.NoRedirect,
	}
}

// MapTransportError wraps an error from http.Client.Do: a deadline or
// network timeout becomes a TimeoutError, anything else a plain
// "request failed" error.
func MapTransportError(err error) error {
	if IsTimeout(err) {
		return NewTimeoutError(fmt.Errorf("request failed: %w", err))
	}
	return fmt.Errorf("request failed: %w", err)
}

// IsTimeout reports whether err represents a request deadline being
// exceeded, covering both an expired context and a net.Error timeout (e.g.
// a dial or read timing out on the underlying transport).
func IsTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// CheckStatus maps an upstream HTTP status to the connector error types:
// 401/403 become AuthError, 502/503/504 ServiceUnavailableError, any other
// status >= 400 a plain error. body is included in the message. It returns
// nil for a status below 400.
func CheckStatus(status int, body []byte) error {
	if status < 400 {
		return nil
	}
	err := fmt.Errorf("API returned %d: %s", status, string(body))
	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return NewAuthError(err)
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return NewServiceUnavailableError(err)
	}
	return err
}
