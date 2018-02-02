package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

// RequestLogger is a middleware that logs the requests
type RequestLogger struct {
	next http.Handler
}

// NewRequestLoggerMW returns a new instance of RequestLogger
// middleware which logs requests
func NewRequestLoggerMW() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return NewRequestLogger(next)
	}
}

// NewRequestLogger returns a new instance of RequestLogger
func NewRequestLogger(next http.Handler) *RequestLogger {
	return &RequestLogger{next: next}
}

func (rl *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lrw := newLoggingResponseWriter(w)
	rl.next.ServeHTTP(lrw, r)
	logRequest(r, lrw.statusCode, start)
}

func logRequest(r *http.Request, statusCode int, start time.Time) {
	elapsed := time.Since(start).Round(time.Microsecond * 100)
	log.Printf("%s %s Completed %v in %s", r.Method, r.URL, statusCode, elapsed)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newLoggingResponseWriter creates a new instance of loggingResponseWriter
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
