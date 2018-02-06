package middleware

import (
	"fmt"
	"log"

	"net/http"
	"runtime/debug"

	"github.com/justinas/alice"
)

// Recovery is a middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	PrintStack bool
	next       http.Handler
}

// NewRecoveryMW returns a new instance of Recovery middleware which traps panics
func NewRecoveryMW() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return NewRecovery(next)
	}
}

// NewRecovery returns a new instance of Recovery
func NewRecovery(next http.Handler) *Recovery {
	return &Recovery{
		PrintStack: true,
		next:       next,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := debug.Stack()

			log.Printf("%s\n%s", err, stack)

			if rec.PrintStack {
				fmtStr := `{"message":"%s","stack":"%s","type":"error"}`
				fmt.Fprintf(rw, fmtStr, err, stack)
			}
		}
	}()

	rec.next.ServeHTTP(rw, r)
}
