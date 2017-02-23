package kit

import (
	"net/http"

	"github.com/urfave/negroni"
	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/middleware/wrpctxmw"
)

type Service interface {
	Name() string
	Handler() http.Handler
}

type service struct {
	name    string
	handler *negroni.Negroni
}

func (s service) Name() string {
	return s.name
}

func (s service) Handler() http.Handler {
	return s.handler
}

func NewService(name string, handlers ...negroni.Handler) Service {
	return service{
		name:    name,
		handler: negroni.New(handlers...),
	}
}

func Classic(name string, handler http.Handler) Service {
	return NewService(
		name,
		wrpctxmw.New(),
		errormw.New(),
		negroni.Wrap(handler),
	)
}
