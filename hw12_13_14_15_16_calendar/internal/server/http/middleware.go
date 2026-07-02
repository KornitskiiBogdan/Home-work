package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type response struct {
	http.ResponseWriter
	status int
}

func loggingMiddleware(log Logger, next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &response{ResponseWriter: w, status: r.Response.StatusCode}

		next.ServeHTTP(wrapped, r)
		latencyMS := time.Since(start).Milliseconds()
		log.Info(formatAccessLog(r, wrapped.status, latencyMS))
	})
}

func formatAccessLog(r *http.Request, status int, latencyMS int64) string {
	ip := clientIP(r)
	ts := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	ua := r.UserAgent()
	if ua == "" {
		ua = "-"
	}

	return fmt.Sprintf(
		`%s [%s] %s %s %s %d %d "%s"`,
		ip,
		ts,
		r.Method,
		path,
		r.Proto,
		status,
		latencyMS,
		ua,
	)
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
