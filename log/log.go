package log

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/wrapp/gokit/env"
)

var jsonFormatter = log.JSONFormatter{}

type wrappFormatter struct{}

func (f *wrappFormatter) Format(entry *log.Entry) ([]byte, error) {
	name := env.ServiceName()
	host, err := os.Hostname()
	if err != nil {
		host = err.Error()
	}

	e := entry.WithFields(log.Fields{
		"service-name": name,
		"host":         host,
	})

	e.Time = time.Now().UTC()
	e.Level = entry.Level
	e.Message = entry.Message
	return (&jsonFormatter).Format(e)
}

func init() {
	log.SetFormatter(&wrappFormatter{})
	log.SetOutput(os.Stdout)
}
