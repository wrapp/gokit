package requestidmw

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/wrapp/gokit/wrpctx"
)

const (
	defaultHeaderKey = "X-Request-Id"
	defaultCtxKey    = "Request-Id"
)

type XRequestIDHandler struct {
	HeaderKey    string
	CtxKey       string
	GenerateFunc func() string
}

func (h XRequestIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, _ := h.getOrGenerate(r)
	if h.CtxKey != "" {
		w.Header().Set(h.HeaderKey, id)
		wrpctx.Set(r.Context(), h.CtxKey, id)
	}
	next(w, r)
}

func (h XRequestIDHandler) getOrGenerate(r *http.Request) (string, bool) {
	id := r.Header.Get(h.HeaderKey)
	if id == "" {
		id = h.GenerateFunc()
		return id, true
	}
	return id, false
}

func generateUUID() string {
	return uuid.NewV4().String()
}

func New() XRequestIDHandler {
	return XRequestIDHandler{
		CtxKey:       defaultCtxKey,
		HeaderKey:    defaultHeaderKey,
		GenerateFunc: generateUUID,
	}
}
