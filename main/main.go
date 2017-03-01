package main

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/wrapp/gokit/kit"
	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/wrpctx"
)

func requestIDGetter(ctx context.Context) func() string {
	return func() string {
		return requestidmw.GetID(ctx)
	}
}

type App struct {
	controller Controller
	router     http.Handler
	service    kit.Service
}

type Controller struct {
}

func (a *App) indexHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	wrpctx.Set(ctx, "key", "value")
	fmt.Fprintf(w, "(%s) %s", requestidmw.GetID(ctx), "Welcome to the home page!")
	log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Log context...")

	//c := trace.NewClient(requestIDGetter(ctx))
	//c.Get("http://localhost:8080/err")
}

func (a *App) errHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	wrpctx.Set(req.Context(), "error", errormw.NewError(http.StatusServiceUnavailable, "Error"))
	log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Error handler")
}

func (a *App) init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.indexHandler)
	mux.HandleFunc("/err", a.errHandler)

	a.router = mux
	a.controller = Controller{}
}
func main() {
	app := &App{}
	app.init()

	srv := kit.Classic(app.router)
	app.service = srv
	//service := kit.NewService(
	//	error.NewErrorMiddleware(),
	//	// one more
	//	// another one
	//	negroni.Wrap(mux),
	//...)

	fmt.Printf("Starting service '%s'...\n", srv.Name())
	err := srv.ListenAndServe("localhost:8080")
	if err != nil {
		log.WithField("error", err.Error()).Error("Service stopped")
	}
}
