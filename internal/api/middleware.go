package api

import (
	"net/http"
	"time"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", deterministicID(r))
		next.ServeHTTP(w, r)
	})
}
func deterministicID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" {
		return v
	}
	return "request-static"
}
func Timeout(next http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(next, d, "timeout")
}
