package middleware

import (
	"net/http"

	"github.com/WiseLabz/wiselabz/internal/httputil"
)

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
			if httputil.IsSecureRequest(r, trustedProxies) {
				h.Set("Strict-Transport-Security", "max-age=15552000")
			}
			next.ServeHTTP(w, r)
		})
	}
}
