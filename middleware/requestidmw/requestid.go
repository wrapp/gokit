// requestidmw is a middleware to add a request-id for each request. It uses X-Request-Id
// header to set or get the request-id. If the incoming request has X-Request-Id header
// then the id will be read from them and it will be set in the header of response writer.
// This id can be used to pass along in other http requests. For example TraceClient can
// use this id to make further http request with the same request-id.

// In addition to http headers it also set the `request_id` field in the wrpctx.

// A new id is generated if there was no header set in the incoming request.
package requestidmw

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/wrapp/gokit/wrpctx"
)

const (
	headerKey = "X-Request-Id"
	ctxKey    = "request_id"
)

type RequestIDFunc func() string

var defGenFunc RequestIDFunc = generateUUID

// XRequestIDHandler contains the generator function of request id. A custom generator
// function can be used to generate new request ids.
type XRequestIDHandler struct {
	GenerateFunc RequestIDFunc
}

func (h XRequestIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, _ := h.getOrGenerate(r)
	w.Header().Set(headerKey, id)
	wrpctx.Set(r.Context(), ctxKey, id)
	next(w, r)
}

func (h XRequestIDHandler) getOrGenerate(r *http.Request) (string, bool) {
	id := IDFromHeader(r.Header)
	if id == "" {
		id = h.GenerateFunc()
		return id, true
	}
	return id, false
}

func generateUUID() string {
	return uuid.NewV4().String()
}

// IDFromCtx returns the request-id from a context.Context. If the request-id is not set in the
// context or it cannot be converted to string then the function will return an empty string.
func IDFromCtx(ctx context.Context) string {
	id := wrpctx.Get(ctx, ctxKey)
	if id == nil {
		return ""
	}
	idStr, ok := id.(string)
	if !ok {
		return ""
	}
	return idStr

}

// IDFromHeader returns the request-id from a http.Header.
func IDFromHeader(h http.Header) string {
	return h.Get(headerKey)
}

// SetIDInHeader sets the request id in http.Header
func SetIDInHeader(h *http.Header, id string) {
	h.Set(headerKey, id)
}

// SetIDInContext sets the request-id in a context.Context
func SetIDInContext(ctx context.Context, id string) {
	wrpctx.Set(ctx, ctxKey, id)
}

// SetDefaultGenFunc gets the default request-id generator function
func SetDefaultGenFunc(genFunc RequestIDFunc) {
	defGenFunc = genFunc
}

// DefaultGenFunc sets the default request-id generator function
func DefaultGenFunc() RequestIDFunc {
	return defGenFunc
}

// New creates a new XRequestIDHandler middleware.
func New() XRequestIDHandler {
	return XRequestIDHandler{
		GenerateFunc: defGenFunc,
	}
}
