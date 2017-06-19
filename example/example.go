package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/wrapp/gokit/kit"
	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/middleware/jsonrqmw"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/wrpctx"
)

func requestIDGetter(ctx context.Context) func() string {
	return func() string {
		return requestidmw.IDFromCtx(ctx)
	}
}

type App struct {
	controller Controller
	router     http.Handler
}

type Controller struct {
}

type JsonRequest struct {
	Type string `json:"type"`
}

func (a *App) indexHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	wrpctx.Set(ctx, "key", "value")
	fmt.Fprintf(w, "(%s) %s", requestidmw.IDFromCtx(ctx), "Welcome to the home page!")
	log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Log context...")

	//c := trace.New(requestIDGetter(ctx))
	//c.Get("http://localhost:8080/err")
}

func (a *App) errHandler(w http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Error handler")
	return errormw.NewError(http.StatusInternalServerError, "Error")
}

func (a *App) jsonHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	json := jsonrqmw.Get(ctx).(*JsonRequest)
	log.WithField("val", json).Info("JSON...")
}

func jsonFactory() interface{} {
	return &JsonRequest{}
}

func (a *App) init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	schema := fmt.Sprintf("file://%s/example/schema.json", wd)

	jsonHandler := negroni.New()
	jsonHandler.UseHandler(jsonrqmw.New(a.jsonHandler, schema, jsonFactory))

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.indexHandler)
	mux.Handle("/err", errormw.ErrorHandler(a.errHandler))
	mux.Handle("/json", jsonHandler)

	a.router = mux
	a.controller = Controller{}
}
func main() {
	app := &App{}
	app.init()

	srv := kit.SimpleService(app.router)
	//srv.SetServiceName("My Service")
	//srv.SetPreShutdownHandler(func() { log.Info("Starting shutdown...") })
	//srv.SetPostShutdownHandler(func() { log.Info("Shutdown completed") })
	//service := kit.NewService(
	//	error.NewErrorMiddleware(),
	//	// one more
	//	// another one
	//	negroni.Wrap(mux),
	//...)

	log.Info("Starting service...")
	err := srv.ListenAndServe("localhost:8080")
	if err != nil {
		log.WithField("error", err.Error()).Error("Service stopped")
	}
}
