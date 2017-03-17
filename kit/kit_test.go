package kit

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/urfave/negroni"

	"github.com/wrapp/gokit/middleware/errormw"
	"github.com/wrapp/gokit/middleware/jsonrqmw"
	"github.com/wrapp/gokit/middleware/recoverymw"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/middleware/wrpctxmw"
	"github.com/wrapp/gokit/wrpctx"
)

type wrpctxTestHandler struct{}

func (h wrpctxTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	wrpctx.Set(ctx, "key", "value")
	fmt.Fprintf(w, "%s", wrpctx.Get(ctx, "key"))
}

type reqidTestHandler struct{}

func (h reqidTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", requestidmw.GetID(r.Context()))
}

func TestWrappContext(t *testing.T) {
	t.Parallel()
	ctxTestHandler := negroni.Wrap(wrpctxTestHandler{})
	service := NewService(wrpctxmw.New(), ctxTestHandler)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if b := w.Body.String(); b != "value" {
		t.Errorf(`body = %q wanted \"value\"`, b)
	}
}

func testGenerateRequestID(t *testing.T) {
	t.Parallel()
	reqidTestHandler := negroni.Wrap(reqidTestHandler{})
	service := NewService(wrpctxmw.New(), requestidmw.New(), reqidTestHandler)

	// test when request-id is not present
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if id := w.Header().Get("X-Request-Id"); id == "" {
		t.Errorf("X-Request-Id was not set ", id)
	}

	if b := w.Body.String(); b == "" {
		t.Errorf("body = %q wanted non blank body", b)
	}
}

func testPassExistingRequestID(t *testing.T) {
	t.Parallel()
	reqidTestHandler := negroni.Wrap(reqidTestHandler{})
	service := NewService(wrpctxmw.New(), requestidmw.New(), reqidTestHandler)

	// test when request-id is present
	testid := "existing-request-id"
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("X-Request-Id", testid)
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if id := w.Header().Get("X-Request-Id"); id != testid {
		t.Errorf("X-Request-Id was %q expected %q", id, testid)
	}

	if b := w.Body.String(); b == "" {
		t.Errorf("body = %q wanted %q", b, testid)
	}
}

func TestRequestID(t *testing.T) {
	t.Parallel()
	t.Run("TestGenerateRequestID", testGenerateRequestID)
	t.Run("TestPassExistingRequestID", testPassExistingRequestID)
}

func getJsonHandler() http.Handler {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	schema := fmt.Sprintf("file://%s/schema_test.json", wd)

	jsonFactory := func() interface{} {
		return &struct {
			Type string `json:"type"`
		}{}
	}
	return jsonrqmw.New(func(w http.ResponseWriter, r *http.Request) {}, schema, jsonFactory)

}

func testJsonContentType(t *testing.T) {
	t.Parallel()

	service := NewService(wrpctxmw.New(), negroni.Wrap(getJsonHandler()))

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func testJsonBody(t *testing.T) {
	t.Parallel()

	service := NewService(wrpctxmw.New(), negroni.Wrap(getJsonHandler()))

	r, _ := http.NewRequest("POST", "/", strings.NewReader(""))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func testJsonSchema(t *testing.T) {
	t.Parallel()

	service := NewService(wrpctxmw.New(), negroni.Wrap(getJsonHandler()))

	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"type":"J"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400 got %d", w.Code)
	}
}

func testValidJsonSchema(t *testing.T) {
	t.Parallel()

	service := NewService(wrpctxmw.New(), negroni.Wrap(getJsonHandler()))

	r, _ := http.NewRequest("POST", "/", strings.NewReader(`{"type":"JSON"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Expected 200 got %d", w.Code)
	}
}

func TestJsonRequest(t *testing.T) {
	t.Parallel()
	t.Run("TestContentType", testJsonContentType)
	t.Run("TestJsonBody", testJsonBody)
	t.Run("TestJsonSchema", testJsonSchema)
	t.Run("TestValidJsonSchema", testValidJsonSchema)
}

func TestErrorMW(t *testing.T) {
	t.Parallel()
	withStatusErr := func(w http.ResponseWriter, r *http.Request) error {
		return errormw.NewError(500, "custom error")
	}

	withAnyErr := func(w http.ResponseWriter, r *http.Request) error {
		return errors.New("any error")
	}

	withNoErr := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	t.Run("WithError", func(t *testing.T) {
		t.Parallel()
		service := NewService(negroni.Wrap(errormw.ErrorHandler(withStatusErr)))

		r, _ := http.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		service.Handler().ServeHTTP(w, r)

		if w.Code != 500 {
			t.Errorf("Expected 500 got %d", w.Code)
		}

		if w.Body.String() != "custom error\n" {
			t.Errorf("Expected 'custom error' got %q", w.Body.String())
		}
	})

	t.Run("WithAnyError", func(t *testing.T) {
		t.Parallel()
		service := NewService(negroni.Wrap(errormw.ErrorHandler(withAnyErr)))

		r, _ := http.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		service.Handler().ServeHTTP(w, r)

		if w.Code != 500 {
			t.Errorf("Expected 500 got %d", w.Code)
		}

		if w.Body.String() != "any error\n" {
			t.Errorf("Expected 'any error' got %q", w.Body.String())
		}
	})

	t.Run("WithNoError", func(t *testing.T) {
		t.Parallel()
		service := NewService(negroni.Wrap(errormw.ErrorHandler(withNoErr)))

		r, _ := http.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		service.Handler().ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Expected 200 got %d", w.Code)
		}
	})
}

type panicHandler struct{}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("do panic")
}

func TestRecoveryMW(t *testing.T) {
	t.Parallel()

	recovery := recoverymw.Recovery{
		PanicHandlerFunc: func(err interface{}, stack []byte) {},
		StackSize:        1024 * 8,
		PrintStack:       true,
	}
	service := NewService(recovery, negroni.Wrap(panicHandler{}))

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	service.Handler().ServeHTTP(w, r)

	if w.Code != 500 {
		t.Errorf("Expected 500 got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "PANIC!: do panic") {
		t.Errorf("Expected 'PANIC!: do panic' in body")
	}
}
