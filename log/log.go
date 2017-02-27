package log

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

var jsonFormatter = logrus.JSONFormatter{}

type wrappFormatter struct{}

// Format logs according to WEP-007
func (f *wrappFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	jsonBytes, err := (&jsonFormatter).Format(entry)
	prefix := []byte(strings.ToUpper(entry.Level.String()) + " ")
	return append(prefix[:], jsonBytes[:]...), err
}

func init() {
	logrus.SetFormatter(&wrappFormatter{})
	logrus.SetOutput(os.Stdout)
}
