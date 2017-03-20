# Gokit
Gokit is a standard library for building microservices. It provides tools to solve some common problems so that you
can focus on building the business logic.

## Motivation
It is usual to solve the same problem over and over again in the individual services when building microservices.
Multiple services end up having duplicate code and it is hard to maintain them. Gokit tries to solve some of these
common problems so that it is easier to build microservices without duplicating code.

## Installation and usage
Gokit can be installed through `go get` or by a dependency manager such as
[glide](https://github.com/Masterminds/glide).

```bash
go get github.com/wrapp/gokit/kit
```

or:

```bash
glide get github.com/wrapp/gokit/kit
```

Import gokit with:

```go
import "github.com/wrapp/gokit/kit"
```

Rest of the document highlights individual components of gokit in detail. A full example of gokit can be found
[here](example/example.go).

## Creating a service
A service is created by calling `SimpleService` function in `kit` package. This creates a service and some default 
middlewares for convinence.

```go
srv := kit.SimpleService(/* any router that implements http.Handler */)        
```
or if you dont want to use default middlewares then you can create the service through `NewService`:

```go
// gokit uses negroni for adding middlewares
srv := kit.NewService(
       	//      error.NewErrorMiddleware(),
       	//      negroni.Wrap(router), 
        //      ...
       )
```

The name of the service can be set through `SetServiceName` method. This will set the name of the service name for
all the default components (e.g default logger).

```go
srv.SetServiceName("my-service")
```

After a service is created you can start it by calling `ListenAndServe` method and passing in the address and port.

```go
err := srv.ListenAndServe("localhost:8080")
```

## Context
Gokit provides some wrapper functions for `context.Context`. These wrappers are used internally for setting data in
context and passing it around in different modules. It is recommended to use these functions when you want to read
or set data in context. For details have a look at [wrpctx](wrpctx/wrpctx.go) package.

## Logging
Gokit provides a custom formatter for [logrus](https://github.com/sirupsen/logrus). This formatter adds some extra
fields to the log entry. The log entry is then formatted to JSON before it is written to the stdout. It is recommended
to use this default formatter but you can easily override it if it is necessary. See logrus's documentation to see how
to override the default logger. Furthermore, you can use your custom formatter or any other logging library if you want.

## Connection draining
Go 1.8 released a feature called [graceful shutdowns](https://golang.org/doc/go1.8#http_shutdown) or connection 
draining. Gokit uses this feature to drain in flight connections. This is the default behaviour of the service. To 
disable connection draining, call `DrainConnections(bool, time.Duration)` before `ListenAndServe` is called.

```go
srv.DrainConnections(false, nil)
```

or to set a custom timeout:

```go
srv.DrainConnections(true, 1 * time.Minute)
```

## Middlewares
Gokit provides some middlewares out of the box. Some of the middlewares are added by default when creating the service
through `SimpleService`.  To use a custom list of middlewares use `NewService` instead. Gokit uses
[negroni](https://github.com/urfave/negroni) for adding the middlewares.

### Wrap Context
`Default: yes`

Wrapp context is a simple middleware to add some extra functionality around `context.Context`. For example it is not
easy to iterate the keys and values in default context. Wrapp context allows you to do that by getting a copy of map 
data stored in the context.

```go
m := wrpctx.GetMap(ctx) // returns map[string]interface{}
```

For more details look at [wrpctx](wrpctx/wrpctx.go)

### Request ID
`Default: yes`

Request ID adds a unique id for each incoming request. This request-id can be used to pass along in downstream 
services to map an end-to-end path of a request. It uses `X-Request-Id` header in the http request. If this header is
already present then it will use the id from that header. A new id is generated otherwise. The request-id generator
function can also be customised as per user needs by passing `GeneateFunc` func.

To get the request-id from context: 

```go
rid := requestidmw.GetID(ctx)
```

This id can then be used in [tracing](#tracing).

### Recovery
`Default: yes`

Recovery middleware provides means to recover from panic that are raised in `http.Handler`. A default handler can be 
created from `New` or you can create your own handler like:

```go
rh := RecoveryHandler {
        PanicHandlerFunc: func(err interface{}, stack []byte) { // err contains panic, stack contains the stacktrace
               // handle err or print stacktrace.
        }, 
        StackSize: 1024 * 8, // size of the stacktrace
        PrintStack: false    // whether to print stacktrace on http.Response
}
```

### Error
`Default: no`

Error middleware is used for a particular endpoint or a handler. It simplifies code for writing http errors to the 
response body. As an example:

```go
func IndexHandler(w http.ResponseWriter, req *http.Request) error {
        // do some processing
        return errormw.NewError(http.StatusInternalServerError, "Something went wrong")
}
```

and in your router/mux:
```go
mux.Handle("/index", errormw.ErrorHandler(IndexHandler))
```

When this handler is executed it will return a `500 InternalServerError` to the response with a message
`Something went wrong`. Any http error with a relevant message can be used.

### JSON Request
`Default: no`

JSON request is used when a JSON object is expected in a request body. This only works with http POST. The 
middleware validates the schema of expected object. If something is not right then it responds with a 
`400 BadRequest`. If everything is fine then the object is unmarshaled and passed on to handler in context.

```go
func objFactory() interface{} {
	return &JsonRequest{}
}

schema := "file://path/to/schema.json"

func JsonHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	json := jsonrqmw.Get(ctx).(*JsonRequest)
	// process json
}

// add it to mux/router
jh = jsonrqmw.New(JsonHandler, schema, objFactory))
mux.Handle("/json", jh)
```

## Tracing
Tracing is a way to pass on request-id to other services through http client. Gokit provides a tracing
[client](trace/trace.go) which can replace standard `net/http` client. To create a `TraceClient` object:

```go
func requestIDGetter(ctx context.Context) func() string {
        return func() string {
                return requestidmw.GetID(ctx)
        }
}

c := trace.New(requestIDGetter(ctx))
c.Get("http://localhost:8080/index")
// Do|Post|Head|PostForm
```

In the above example `requestIDGetter` func is created to get a request-id from the [middleware](#request-id).
Typically only one `requestIDGetter` is needed. It is also possible to set the `User-Agent` for outgoing requests
going through trace client. It is set to the `env.ServiceName` by default.

```go
c.SetUserAgent("my-agent")
```

## Other
Gokit also provides some extra utilities.

### Env
Some handy methods to get environment variables can be found [here](env/env.go). 

Name of the service can be retrieved through `env.ServiceName()` if `SERVICE_NAME` environment variable is set.  

### Short circuit
Short circuiting errors:
```go
// err will have e1 after the execution and second function will not be executed
// because first function returned an error
err := ShortCircuit(
        func() error { return e1 },
        func() error { return e2 },
)
```
