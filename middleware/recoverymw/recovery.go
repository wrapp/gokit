package recoverymw

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

type PanicHandlerFunc func(interface{}, []byte)

type Recovery struct {
	PanicHandlerFunc PanicHandlerFunc
	StackSize        int
	PrintStack       bool
}

func New() Recovery {
	return Recovery{defaultPanicHandler, 1024 * 50, false}
}

func defaultPanicHandler(err interface{}, stack []byte) {
	log.WithFields(log.Fields{
		"panic":      err,
		"stacktrace": string(stack),
	}).Error("PANIC! in http handler")
}

func (rec Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
