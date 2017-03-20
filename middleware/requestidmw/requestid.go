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

// XRequestIDHandler contains the generator function of request id. A custom generator
// function can be used to generate new request ids.
type XRequestIDHandler struct {
	GenerateFunc func() string
}

func (h XRequestIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, _ := h.getOrGenerate(r)
	w.Header().Set(headerKey, id)
	wrpctx.Set(r.Context(), ctxKey, id)
	next(w, r)
}

func (h XRequestIDHandler) getOrGenerate(r *http.Request) (string, bool) {
	id := r.Header.Get(headerKey)
	if id == "" {
		id = h.GenerateFunc()
		return id, true
	}
	return id, false
}

func generateUUID() string {
	return uuid.NewV4().String()
}

// GetID returns the request-id from a context.Context.
func GetID(ctx context.Context) string {
	return wrpctx.Get(ctx, ctxKey).(string)
}

// SetID sets the request id in an http.Header
func SetID(h *http.Header, id string) {
	h.Set(headerKey, id)
}

// New creates a new XRequestIDHandler middleware.
func New() XRequestIDHandler {
	return XRequestIDHandler{
		GenerateFunc: generateUUID,
	}
}
