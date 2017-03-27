// errormw is a middleware that allows http.Handler to return a specific http error code
// with a custom message. This is usually used on a specific endpoint because not all
// endpoints return an error. It abstracts writing the header and code to the http.Response
// from the http.Handler so it is much cleaner to use this middleware when returning http
// errors.
package errormw

import "net/http"

// StatusError is an interface which is implemented by an underlying struct. It includes
// an integer which is http status code and an error which is a description of the error.
type StatusError interface {
	Status() int
	error
}

type statusError struct {
	status  int
	message string
}

// Status returns the http status code of the error. e.g 500
func (e statusError) Status() int {
	return e.status
}

// Error returns the human readable description of the error.
func (e statusError) Error() string {
	return e.message
}

// ErrHandler is a function which can be used in place of http.Handler to return errors
// through this middleware.
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

// Creates a new StatusError with provided http status code and a message.
func NewError(status int, message string) StatusError {
	return statusError{
		status:  status,
		message: message,
	}
}
