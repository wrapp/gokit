package main

import (
	"fmt"
	"net/http"

	"github.com/wrapp/gokit/kit"
	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/wrpctx"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
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
	err := http.ListenAndServe(":8080", service.Handler())
	fmt.Println(err)
}
