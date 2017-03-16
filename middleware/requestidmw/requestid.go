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

func GetID(ctx context.Context) string {
	return wrpctx.Get(ctx, ctxKey).(string)
}

func SetID(h *http.Header, id string) {
	h.Set(headerKey, id)
}

func New() XRequestIDHandler {
	return XRequestIDHandler{
		GenerateFunc: generateUUID,
	}
}
