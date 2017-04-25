package kit

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/negroni"

	"github.com/wrapp/gokit/env"
	kitlog "github.com/wrapp/gokit/log"
	"github.com/wrapp/gokit/middleware/recoverymw"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/middleware/wrpctxmw"
)

type ShutdownHandlerFunc func()

// Service interface provides the functionality of any service. It allows you to
// define your implementation of a service if you need to.
type Service interface {
	Handler() http.Handler
	DrainConnections(bool, time.Duration)
	SetPreShutdownHandler(ShutdownHandlerFunc)
	SetPostShutdownHandler(ShutdownHandlerFunc)
	SetServiceName(string)
	ListenAndServe(string) error
}

type service struct {
	drainConn    bool
	timeout      time.Duration
	preShutdown  ShutdownHandlerFunc
	postShutdown ShutdownHandlerFunc
	handler      *negroni.Negroni
}

// Handler returns the http.Handler of the service. When a service is started this handler is
// used to serve the requests over HTTP.
func (s *service) Handler() http.Handler {
	return s.handler
}

// DrainConnections allows you to enable or disable graceful shutdowns. This
// functionality was used on go 1.8. It allows you to wait for in-flight connections
// before a service can be shutdown. The advantage is the requests will not be
// killed immediately and gives those requests a chance to finish properly.

// If `drain` is set to true then it will enable graceful shutdowns. Service will
// wait for `timeout` before shutting down forcefully. If `drain`is set to false
// then `timeout` is ignored.
func (s *service) DrainConnections(drain bool, timeout time.Duration) {
	s.drainConn = drain
	s.timeout = timeout
}

// SetServiceName sets the name of the service for all default components. If there are
// custom components then programmer has the responsibility to set those properly.
func (s *service) SetServiceName(name string) {
	kitlog.SetServiceName(name)
}

// SetPreShutdownHandler sets a custom `handler` function which is called just before service
// starts the shutdown process when in connection draining is set. This function has no effect
// if connection draining is not set. See DrainConnections for more information.
func (s *service) SetPreShutdownHandler(handler ShutdownHandlerFunc) {
	s.preShutdown = handler
}

// SetPostShutdownHandler sets a custom `handler` function which is called right after service
// shuts down when in connection draining is set. This handler will be called even if there
// was an error while shutting down the service. This function has no effect if connection
// draining is not set. See DrainConnections for more information.
func (s *service) SetPostShutdownHandler(handler ShutdownHandlerFunc) {
	s.postShutdown = handler
}

// ListenAndServe starts the service on given address. `addr` contains the ip of the
// interface and port in the form `ip-addr:port` e.g `0.0.0.0:8080`. This is a blocking
// call unless there is an error which will be returned when function exits. By default
// graceful shutdowns are enabled and the service will wait for in-flight requests
// when an OS interrupt is received before shutting down. This behaviour can be turned
// off with DrainConnections function.

// By default all the timeouts (ReadTimeout, WriteTimeout, IdleTimeout, ReadHeaderTimeout)
// are set to 60s. These timeouts are set to avoid memory leaks.
func (s *service) ListenAndServe(addr string) error {
	srv := http.Server{
		Addr:              addr,
		Handler:           s.handler,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	errorChan := make(chan error)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorChan <- err
		}
	}()

	var err error
	select {
	case <-stopChan:
		if s.drainConn {
			if s.preShutdown != nil {
				s.preShutdown()
			}

			ctx, _ := context.WithTimeout(context.Background(), s.timeout)
			err = srv.Shutdown(ctx)

			if s.postShutdown != nil {
				s.postShutdown()
			}
		} else {
			err = srv.Close()
		}
	case err = <-errorChan:
	}

	return err
}

// NewService creates a new service with all the custom handlers provided in the arguments.
// This will not add any default handlers in the service.
func NewService(handlers ...negroni.Handler) Service {
	return &service{
		drainConn: true,
		timeout:   25 * time.Second,
		handler:   negroni.New(handlers...),
	}
}

// SimpleService initializes the service with some default middlewares. The `http.Handler`
// provided will be used as the last handler in the service. `http.Handler` usually contains
// the endpoints and business logic of the service.
// Following middlewares are initialized (in order) when calling this function.
/*
	- Wrapp Context (wrpctxmw) is a wrapper around `context.Context`.
	- Request ID (requestidmw) adds a unique id for each incoming request.
	- Recovery (recoverymw) provides functionality to recover from panics in the http.Handler.
*/
func SimpleService(handler http.Handler) Service {
	s := NewService(
		wrpctxmw.New(),
		requestidmw.New(),
		recoverymw.New(),
		negroni.Wrap(handler),
	)
	s.SetServiceName(env.ServiceName())
	return s
}
