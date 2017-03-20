// jsonrqmw provides a handler to receive JSON request from the client. It allows the possibility
// to receive JSON objects in request body (using POST as method), validates the JSON objects
// against a JSON schema, Unmarshal the JSON object and then pass it on to the http.Handler.
// http.Handler then can get the object from context.Context and uses it for processing. If
// any of these steps fails the middleware returns a 400 BadRequest to the response without
// forwarding the JSON to the http.Handler.
package jsonrqmw

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/wrapp/gokit/wrpctx"
)

// A factory function that must be provided to create JSON objects. Middleware uses this
// function to creates JSON objects and unmarshal the passed json to that object. If it is
// successful this object is added to the context.Context.
type JsonObjectFactory func() interface{}

type jsonRequestHandler struct {
	handler    http.HandlerFunc
	schema     *gojsonschema.Schema
	objFactory JsonObjectFactory
}

const jsonKey = "json"

func (j *jsonRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ctype := r.Header.Get("Content-Type"); !strings.Contains(ctype, "application/json") {
		http.Error(w, "Content-Type is not application/json", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strBody := string(body)
	res, err := j.schema.Validate(gojsonschema.NewStringLoader(strBody))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !res.Valid() {
		http.Error(w, toReadable(res), http.StatusBadRequest)
		return
	}

	obj := j.objFactory()
	err = json.Unmarshal(body, obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := wrpctx.NewWithValue(r.Context(), jsonKey, obj)
	cr := r.WithContext(ctx)
	j.handler(w, cr)
}

// New creates an http.Handler which wraps the passed http.HandlerFunc. The http.Handler
// adds the functionality to verify the incoming request with passed JSON schema and
// unmarshales the object into the object created by passed JSON object factory.
func New(h http.HandlerFunc, schemaPath string, objf JsonObjectFactory) http.Handler {
	return &jsonRequestHandler{
		handler:    h,
		schema:     loadSchema(schemaPath),
		objFactory: objf,
	}
}

// Get returns the unmarshaled JSON object from the context.Context. It will be added
// to the context only when it was successfully validated and parsed to the JSON object
// returned from the factory.
func Get(ctx context.Context) interface{} {
	return wrpctx.GetCtxValue(ctx, jsonKey)
}

func loadSchema(schemaFile string) *gojsonschema.Schema {
	loader := gojsonschema.NewReferenceLoader(schemaFile)
	schema, err := gojsonschema.NewSchema(loader)

	if err != nil {
		panic(err)
	}
	return schema
}

func toReadable(result *gojsonschema.Result) string {
	var errStr []string
	for _, e := range result.Errors() {
		errStr = append(errStr, e.Description())
	}
	return strings.Join(errStr, "\n")
}
