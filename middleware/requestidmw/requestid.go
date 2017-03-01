package requestidmw

import (
	"net/http"

	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/wrapp/gokit/wrpctx"
)

const (
	defaultHeaderKey = "X-Request-Id"
	ctxKey           = "request-id"
)

type XRequestIDHandler struct {
	HeaderKey    string
	GenerateFunc func() string
}

func (h XRequestIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, _ := h.getOrGenerate(r)
	w.Header().Set(h.HeaderKey, id)
	wrpctx.Set(r.Context(), ctxKey, id)
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

func GetID(ctx context.Context) string {
	return wrpctx.Get(ctx, ctxKey).(string)
}

func New() XRequestIDHandler {
	return XRequestIDHandler{
		HeaderKey:    defaultHeaderKey,
		GenerateFunc: generateUUID,
	}
}
