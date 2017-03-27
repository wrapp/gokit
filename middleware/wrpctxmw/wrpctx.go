// wrpctxmw provides a middleware to add some data to context.Context. Each request received
// will have some default fields added to the context. This middleware uses wrpctx to interact
// with those values. This middleware adds a standard way to add and retrieve data from the
// context.Context. In addition, it also allows to iterate over the data which was added which
// is not available in context.Context.
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

// New creates a new WrpCtxHandler middleware.
func New() WrpCtxHandler {
	return WrpCtxHandler{}
}
