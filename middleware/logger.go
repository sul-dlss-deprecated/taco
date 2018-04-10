package middleware

import (
	"fmt"
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
	logRequest(r, lrw.statusCode, lrw.message, start)
}

func logRequest(r *http.Request, statusCode int, message *[]byte, start time.Time) {
	elapsed := time.Since(start).Round(time.Microsecond * 100)
	str := fmt.Sprintf("%s %s Completed %v in %s [%s]", r.Method, r.URL, statusCode, elapsed, getIPAddress(r))
	if message != nil {
		log.Printf("%s\n\t%s", str, string(*message))
	} else {
		log.Println(str)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	message    *[]byte
}

// newLoggingResponseWriter creates a new instance of loggingResponseWriter
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, nil}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(body []byte) (int, error) {
	switch lrw.statusCode {
	case 200, 201, 404:
		// do nothing
	default:
		// log the error
		lrw.message = &body
	}
	return lrw.ResponseWriter.Write(body)
}
