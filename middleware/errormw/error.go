package errormw

import "net/http"

type StatusError interface {
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

type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	err, ok := err.(error)
	if err == nil || !ok {
		return
	}

	status := http.StatusInternalServerError
	if herr, ok := err.(StatusError); ok {
		status = herr.Status()
	}
	http.Error(w, err.Error(), status)
}

func NewError(status int, message string) StatusError {
	return statusError{
		status:  status,
		message: message,
	}
}
