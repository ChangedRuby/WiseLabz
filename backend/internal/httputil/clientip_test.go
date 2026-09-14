package httputil

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientIPUntrustedPeerIgnoresHeaders(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "203.0.113.9:5555"
	r.Header.Set("X-Forwarded-For", "1.2.3.4")

	if got := ClientIP(r, "10.0.0.0/8"); got != "203.0.113.9" {
		t.Fatalf("ClientIP = %q, want peer address for untrusted proxy", got)
	}
}

func TestClientIPTrustedPeerUsesForwardedFor(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.5:5555"
	r.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.5")

	if got := ClientIP(r, "10.0.0.0/8"); got != "198.51.100.7" {
		t.Fatalf("ClientIP = %q, want real client from forwarded chain", got)
	}
}

func TestClientIPRejectsNonIPForwardedFor(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.5:5555"
	r.Header.Set("X-Forwarded-For", "not-an-ip")

	if got := ClientIP(r, "10.0.0.0/8"); got != "10.0.0.5" {
		t.Fatalf("ClientIP = %q, want fallback to peer for unparseable header", got)
	}
}

// TestIsSecureRequestUntrustedPeerIgnoresForwardedProto is a regression test
// for GHSA-c753: X-Forwarded-Proto must not be trusted from a peer outside
// the configured trusted-proxy list, or an attacker can force a cookie's
// Secure flag off by spoofing the header directly.
func TestIsSecureRequestUntrustedPeerIgnoresForwardedProto(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "203.0.113.9:5555"
	r.Header.Set("X-Forwarded-Proto", "https")

	if IsSecureRequest(r, "10.0.0.0/8") {
		t.Fatal("IsSecureRequest() = true for untrusted peer's forwarded header, want false")
	}
}

func TestIsSecureRequestTrustedPeerUsesForwardedProto(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.5:5555"
	r.Header.Set("X-Forwarded-Proto", "https")

	if !IsSecureRequest(r, "10.0.0.0/8") {
		t.Fatal("IsSecureRequest() = false for trusted proxy's forwarded header, want true")
	}
}

func TestIsSecureRequestTLS(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.TLS = &tls.ConnectionState{}

	if !IsSecureRequest(r, "") {
		t.Fatal("IsSecureRequest() = false for a TLS connection, want true")
	}
}
