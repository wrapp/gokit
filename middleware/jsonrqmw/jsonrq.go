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

func New(h http.HandlerFunc, schemaPath string, objf JsonObjectFactory) http.Handler {
	return &jsonRequestHandler{
		handler:    h,
		schema:     loadSchema(schemaPath),
		objFactory: objf,
	}
}

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
