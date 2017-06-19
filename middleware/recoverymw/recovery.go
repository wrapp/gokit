// recoverymw provides a package to recover from panics in http.Handler. It is useful to send a
// http error to the client instead of crashing the server because of a programming error in
// http.Handler. Recovery handler recovers from the panic and calls a handler func to handle
// panic gracefully. It also allows the possibility to log or send the stacktrace to external
// service.
package recoverymw

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// PanicHandlerFunc is a handler func which is called when middleware recovers from panic in
// http.Handler
type PanicHandlerFunc func(interface{}, []byte)

// RecoveryHandler is struct which holds the PanicHandlerFunc, size of the stacktrace, and
// a field which tells whether to print the stack in the http.Response or not.
type RecoveryHandler struct {
	PanicHandlerFunc PanicHandlerFunc
	StackSize        int
	PrintStack       bool
}

// New generates a default RecoveryHandler. By default panics are logged in stdout with a
// stacktrace size of 50KB. Stacktrace is not logged to http.Response be default.
func New() RecoveryHandler {
	return RecoveryHandler{defaultPanicHandler, 1024 * 50, false}
}

func defaultPanicHandler(err interface{}, stack []byte) {
	log.WithFields(log.Fields{
		"panic": err,
		"data":  log.Fields{"stacktrace": string(stack)},
	}).Error("PANIC! in http handler")
}

func (rec RecoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, true)]

			defer func() {
				if recErr := recover(); recErr != nil {
					log.WithField("panic", recErr).Error("Error in panic handler")
					http.Error(w, "", http.StatusInternalServerError)
				}
			}()
			rec.PanicHandlerFunc(err, stack)

			var body string
			if rec.PrintStack {
				body = fmt.Sprintf("PANIC!: %s\n%s", err, stack)
			}
			http.Error(w, body, http.StatusInternalServerError)
		}
	}()
	next(w, r)
}
