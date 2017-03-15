package kit

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/wrapp/gokit/env"
	kitlog "github.com/wrapp/gokit/log"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/middleware/wrpctxmw"
)

type Service interface {
	Handler() http.Handler
	DrainConnections(bool, time.Duration)
	SetServiceName(string)
	ListenAndServe(string) error
}

type service struct {
	drainConn bool
	timeout   time.Duration
	handler   *negroni.Negroni
}

func (s *service) Handler() http.Handler {
	return s.handler
}

func (s *service) DrainConnections(drain bool, timeout time.Duration) {
	s.drainConn = drain
	s.timeout = timeout
}

func (s *service) SetServiceName(name string) {
	kitlog.SetServiceName(name)
}

func (s *service) ListenAndServe(addr string) error {
	srv := http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

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
			ctx, _ := context.WithTimeout(context.Background(), s.timeout)
			err = srv.Shutdown(ctx)
		} else {
			err = srv.Close()
		}
	case err = <-errorChan:
	}

	return err
}

func NewService(handlers ...negroni.Handler) Service {
	return &service{
		drainConn: true,
		timeout:   25 * time.Second,
		handler:   negroni.New(handlers...),
	}
}

func Classic(handler http.Handler) Service {
	recoverymw := negroni.NewRecovery()
	recoverymw.Logger = log.StandardLogger()
	recoverymw.PrintStack = false

	s := NewService(
		wrpctxmw.New(),
		requestidmw.New(),
		recoverymw,
		negroni.Wrap(handler),
	)
	s.SetServiceName(env.ServiceName())
	return s
}
