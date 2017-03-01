package log

import (
	"os"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/wrapp/gokit/env"
)

var jsonFormatter = logrus.JSONFormatter{}

type wrappFormatter struct{}

func (f *wrappFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	name := env.Get("SERVICE_NAME")
	host, err := os.Hostname()
	if err != nil {
		host = err.Error()
	}

	e := entry.WithFields(logrus.Fields{
		"service-name": name,
		"host":         host,
	})

	e.Time = time.Now().UTC()
	e.Level = entry.Level
	e.Message = entry.Message
	return (&jsonFormatter).Format(e)
}

func init() {
	logrus.SetFormatter(&wrappFormatter{})
	logrus.SetOutput(os.Stdout)
}
