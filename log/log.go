package log

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

var jsonFormatter = log.JSONFormatter{
	TimestampFormat: time.RFC3339,
	FieldMap: log.FieldMap{
		log.FieldKeyTime: "timestamp",
	},
}

type wrappFormatter struct {
	service string
}

var formatter = &wrappFormatter{}

func (f wrappFormatter) Format(entry *log.Entry) ([]byte, error) {
	e := entry.WithFields(log.Fields{
		"service": f.service,
	})

	e.Time = time.Now().UTC()
	e.Level = entry.Level
	e.Message = entry.Message
	return (&jsonFormatter).Format(e)
}

func SetServiceName(service string) {
	formatter.service = service
}

func init() {
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}
