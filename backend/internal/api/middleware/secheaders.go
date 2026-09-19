package middleware

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/httputil"
)

// contentSecurityPolicy is tuned for the embedded SPA: Vite bundles scripts and
// self-hosted @fontsource fonts (same origin, no external hosts), the app uses a
// same-origin WebSocket, and React/Radix set inline style attributes, which
// requires 'unsafe-inline' for styles only. Scripts stay strictly 'self'.
const contentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self'; " +
	"style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data: blob:; " +
	"font-src 'self' data:; " +
	"connect-src 'self'; " +
	"frame-ancestors 'none'; " +
	"base-uri 'self'; " +
	"object-src 'none'; " +
	"form-action 'self'"

// SecurityHeaders sets baseline browser hardening headers on every response.
// HSTS is only sent when the request arrived over HTTPS (directly or via a
// trusted proxy), so plain-HTTP self-hosted setups are not pinned to TLS.
func SecurityHeaders(trustedProxies string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "same-origin")
			h.Set("Content-Security-Policy", contentSecurityPolicy)
			if httputil.IsSecureRequest(r, trustedProxies) {
				h.Set("Strict-Transport-Security", "max-age=15552000")
			}
			next.ServeHTTP(w, r)
		})
	}
}
