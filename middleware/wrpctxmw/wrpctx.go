package wrpctxmw

import (
	"net/http"

	"github.com/wrapp/gokit/wrpctx"
)

type WrpCtxHandler struct{}

func (h WrpCtxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cr := r.WithContext(wrpctx.New(r.Context()))
	next(w, cr)
}

func New() WrpCtxHandler {
	return WrpCtxHandler{}
}
