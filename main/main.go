package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/wrapp/gokit/kit"
	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/wrpctx"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		wrpctx.Set(ctx, "key", "value")
		fmt.Fprintf(w, "(%s) %s", wrpctx.Get(ctx, "Request-Id"), "Welcome to the home page!")
		log.WithFields(log.Fields(wrpctx.GetMap(ctx))).Info("Log context...")
	})

	mux.HandleFunc("/err", func(w http.ResponseWriter, req *http.Request) {
		wrpctx.Set(req.Context(), "error", errormw.NewError(http.StatusBadRequest, "Error"))
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
