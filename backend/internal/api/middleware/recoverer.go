package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recoverer catches panics in downstream handlers, logs them, and returns 500.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			slog.Error("panic recovered",
				"panic", rec,
				"stack", string(debug.Stack()),
				"path", loggablePath(r),
				"method", r.Method,
				"request_id", GetRequestID(r.Context()),
			)
			http.Error(w, `{"code":"internal_error","message":"An internal error occurred"}`, http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}
