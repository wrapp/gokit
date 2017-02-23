package errormw

import (
	"net/http"

	"github.com/wrapp/gokit/wrpctx"
)

type HttpError interface {
	Status() int
	error
}

type statusError struct {
	status  int
	message string
}

func (e statusError) Status() int {
	return e.status
}

func (e statusError) Error() string {
	return e.message
}

type ErrorHandler struct{}

func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	err, ok := wrpctx.Get(r.Context(), "error").(error)
	if err == nil || !ok {
		return
	}

	status := http.StatusInternalServerError
	if herr, ok := err.(HttpError); ok {
		status = herr.Status()
	}
	http.Error(w, err.Error(), status)
}

func NewError(status int, message string) HttpError {
	return statusError{
		status:  status,
		message: message,
	}
}

func New() *ErrorHandler {
	return &ErrorHandler{}
}
