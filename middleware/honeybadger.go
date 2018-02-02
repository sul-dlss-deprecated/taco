package middleware

import (
	"net/http"

	"github.com/honeybadger-io/honeybadger-go"
	"github.com/justinas/alice"
)

// NewHoneyBadgerMW returns a new instance of HoneyBadger
// middleware which logs requests
func NewHoneyBadgerMW() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return honeybadger.Handler(next)
	}
}
