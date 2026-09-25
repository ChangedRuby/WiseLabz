// Package httpx builds hardened outbound HTTP clients shared by every
// non-inbound caller (AI providers, notification webhooks, doc export Git
// remotes, and connectors via connector.NewHTTPClient — see #265): TLS 1.2+,
// bounded timeouts, and no redirect following. RetryTransport adds bounded
// retries for idempotent requests.
package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// DefaultTimeout is the whole-request timeout used when Options.Timeout is zero.
const DefaultTimeout = 30 * time.Second

// Options configures NewTransport / NewClient. The zero value is valid.
type Options struct {
	// Timeout bounds the whole request (http.Client.Timeout). Zero means
	// DefaultTimeout; negative means no client-level timeout (for long
	// streaming transfers such as a Git clone).
	Timeout time.Duration
	// InsecureSkipVerify disables TLS certificate verification.
	InsecureSkipVerify bool
	// DialContext overrides the dialer, e.g. connector.GuardedDialer to
	// block loopback/link-local targets. Nil uses a plain net.Dialer and
	// honors HTTP(S)_PROXY. When set, proxy environment variables are ignored
	// so the dialer always sees the real target (a proxy would bypass a guard).
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// defaultResponseHeaderTimeout bounds the wait for response headers when the
// client has no overall timeout.
const defaultResponseHeaderTimeout = 30 * time.Second

// NewTransport returns an *http.Transport with TLS 1.2+ and bounded
// dial/handshake/header timeouts. The response-header timeout matches the
// effective client timeout so slow non-streaming responses (e.g. LLM
// completions) are not cut off before the client deadline.
func NewTransport(o Options) *http.Transport {
	dial := o.DialContext
	proxy := http.ProxyFromEnvironment
	if dial == nil {
		dial = (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext
	} else {
		proxy = nil
	}
	headerTimeout := clientTimeout(o.Timeout)
	if headerTimeout == 0 {
		headerTimeout = defaultResponseHeaderTimeout
	}
	return &http.Transport{
		Proxy:                 proxy,
		DialContext:           dial,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: o.InsecureSkipVerify}, //nolint:gosec // opt-in per caller
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: headerTimeout,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
	}
}

// NewClient returns an *http.Client using NewTransport(o) that never follows
// redirects.
func NewClient(o Options) *http.Client {
	return &http.Client{
		Timeout:       clientTimeout(o.Timeout),
		Transport:     NewTransport(o),
		CheckRedirect: NoRedirect,
	}
}

// clientTimeout resolves Options.Timeout to an http.Client.Timeout:
// zero becomes DefaultTimeout and negative becomes 0 (no timeout).
func clientTimeout(t time.Duration) time.Duration {
	switch {
	case t == 0:
		return DefaultTimeout
	case t < 0:
		return 0
	}
	return t
}

// NoRedirect is an http.Client.CheckRedirect func that returns the redirect
// response itself instead of following it.
func NoRedirect(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }

// IsSafeMethod reports whether method is a safe (read-only, idempotent) HTTP
// method: GET, HEAD, or OPTIONS. Such requests may be safely retried and, for
// API keys, may be allowed under read-only restrictions.
func IsSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}
