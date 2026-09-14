package httputil

import (
	"net"
	"net/http"
	"strings"
)

// ClientIP resolves the request's real client IP. X-Forwarded-For and
// X-Real-IP are attacker-controlled on any request that didn't pass through
// a trusted reverse proxy, so they are only trusted when the immediate peer
// (RemoteAddr) falls inside trustedProxyCIDRs; otherwise the peer address is
// used as-is. trustedProxyCIDRs is comma-separated, e.g. "10.0.0.0/8".
func ClientIP(r *http.Request, trustedProxyCIDRs string) string {
	peer := hostOnly(r.RemoteAddr)

	if !isTrustedProxy(peer, trustedProxyCIDRs) {
		return peer
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		// Walk from the right: each trusted proxy appends its peer's address,
		// so the first entry past our own trusted hops is the real client.
		for i := len(parts) - 1; i >= 0; i-- {
			candidate := strings.TrimSpace(parts[i])
			if net.ParseIP(candidate) == nil {
				continue
			}
			if isTrustedProxy(candidate, trustedProxyCIDRs) {
				continue
			}
			return candidate
		}
	}

	if xri := strings.TrimSpace(r.Header.Get("X-Real-IP")); xri != "" && net.ParseIP(xri) != nil {
		return xri
	}

	return peer
}

// IsSecureRequest reports whether the request arrived over HTTPS, for
// deciding a response cookie's Secure flag. X-Forwarded-Proto is
// attacker-controlled on any request that didn't pass through a trusted
// reverse proxy, so it is only trusted under the same peer check ClientIP
// uses for X-Forwarded-For/X-Real-IP.
func IsSecureRequest(r *http.Request, trustedProxyCIDRs string) bool {
	if r.TLS != nil {
		return true
	}
	if !isTrustedProxy(hostOnly(r.RemoteAddr), trustedProxyCIDRs) {
		return false
	}
	return r.Header.Get("X-Forwarded-Proto") == "https"
}

func hostOnly(addr string) string {
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}

func isTrustedProxy(ip, cidrs string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, c := range strings.Split(cidrs, ",") {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		_, network, err := net.ParseCIDR(c)
		if err != nil {
			continue
		}
		if network.Contains(parsed) {
			return true
		}
	}
	return false
}
