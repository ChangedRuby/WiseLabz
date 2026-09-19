package middleware

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/WiseLabz/wiselabz/internal/logsafe"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status  int
	written int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Hijack lets the responseWriter satisfy http.Hijacker so WebSocket upgrades work.
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// sharePathPrefixes are the unauthenticated share-link route; the segment after it
// is a bearer token and must never reach the logs.
var sharePathPrefixes = []string{"/api/v1/share/", "/api/share/"}

// loggablePath returns the request path with secrets masked: the share-link
// token segment is replaced by a placeholder, keeping the rest of the route.
func loggablePath(r *http.Request) string {
	p := r.URL.Path
	for _, prefix := range sharePathPrefixes {
		if rest, ok := strings.CutPrefix(p, prefix); ok {
			_, tail, _ := strings.Cut(rest, "/")
			p = prefix + "{token}"
			if tail != "" {
				p += "/" + tail
			}
			break
		}
	}
	return logsafe.Sanitize(p)
}

// loggableQuery masks the one-time WebSocket ticket in the query string.
func loggableQuery(r *http.Request) string {
	q := r.URL.Query()
	if !q.Has("ticket") {
		return logsafe.Sanitize(r.URL.RawQuery)
	}
	q.Set("ticket", "REDACTED")
	return logsafe.Sanitize(q.Encode())
}

// Logger logs each request with method, path, status, duration, and request ID.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rid := GetRequestID(r.Context())

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		slog.Info("request",
			"method", r.Method,
			"path", loggablePath(r),
			"query", loggableQuery(r),
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
			"request_id", rid,
		)
	})
}
