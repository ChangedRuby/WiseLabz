// Package httpx builds hardened outbound HTTP clients shared by every
// non-inbound caller (AI providers, notification webhooks, doc export Git
// remotes, and eventually connectors — see #265): TLS 1.2+, bounded
// timeouts, and no redirect following.
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
	// block loopback/link-local targets. Nil uses a plain net.Dialer.
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// NewTransport returns an *http.Transport with TLS 1.2+ and bounded
// dial/handshake/header timeouts.
func NewTransport(o Options) *http.Transport {
	dial := o.DialContext
	if dial == nil {
		dial = (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dial,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: o.InsecureSkipVerify}, //nolint:gosec // opt-in per caller
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
	}
}

// NewClient returns an *http.Client using NewTransport(o) that never follows
// redirects.
func NewClient(o Options) *http.Client {
	timeout := o.Timeout
	switch {
	case timeout == 0:
		timeout = DefaultTimeout
	case timeout < 0:
		timeout = 0
	}
	return &http.Client{
		Timeout:       timeout,
		Transport:     NewTransport(o),
		CheckRedirect: NoRedirect,
	}
}

// NoRedirect is an http.Client.CheckRedirect func that returns the redirect
// response itself instead of following it.
func NoRedirect(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
