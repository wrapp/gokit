package kit

import (
	"net/http"

	"time"

	"os"
	"os/signal"

	"context"

	"github.com/urfave/negroni"

	"github.com/wrapp/gokit/env"
	_ "github.com/wrapp/gokit/log"
	"github.com/wrapp/gokit/middleware/requestidmw"
	"github.com/wrapp/gokit/middleware/wrpctxmw"
)

type Service interface {
	Name() string
	Handler() http.Handler
	DrainConnections(bool, time.Duration)
	ListenAndServe(string) error
}

type service struct {
	drainConn bool
	timeout   time.Duration
	handler   *negroni.Negroni
}

func (s *service) Name() string {
	return env.Get("SERVICE_NAME")
}

func (s *service) Handler() http.Handler {
	return s.handler
}

func (s *service) DrainConnections(drain bool, timeout time.Duration) {
	s.drainConn = drain
	s.timeout = timeout
}

func (s *service) ListenAndServe(addr string) error {
	srv := http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	var err error
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			close(stopChan)
		}
	}()

	<-stopChan // wait
	if s.drainConn {
		ctx, _ := context.WithTimeout(context.Background(), s.timeout)
		err = srv.Shutdown(ctx)
	} else {
		err = srv.Close()
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
	return NewService(
		wrpctxmw.New(),
		requestidmw.New(),
		negroni.Wrap(handler),
	)
}
