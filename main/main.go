package main

import (
	"fmt"
	"net/http"

	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/wrapp/gokit/kit"
	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/trace"
	"github.com/wrapp/gokit/wrpctx"
)

func requestIDGetter(ctx context.Context) func() string {
	return func() string {
		return requestidmw.GetID(ctx)
	}
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		wrpctx.Set(ctx, "key", "value")
		fmt.Fprintf(w, "(%s) %s", requestidmw.GetID(ctx), "Welcome to the home page!")
		log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Log context...")

		c := trace.NewClient("My Service", requestIDGetter(ctx))
		c.Get("http://localhost:8080/err")
	})

	mux.HandleFunc("/err", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		wrpctx.Set(req.Context(), "error", errormw.NewError(http.StatusServiceUnavailable, "Error"))
		log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Error handler")
	})

	service := kit.Classic("My service", mux)
	//service := kit.NewService("My service", []negroni.Handler{
	//	error.NewErrorMiddleware(),
	//	// one more
	//	// another one
	//	negroni.Wrap(mux),
	//}...)

	fmt.Printf("Starting service '%s'...\n", service.Name())
	err := service.ListenAndServe("localhost:8080")
	fmt.Println(err)
}
