package kit

import log "github.com/Sirupsen/logrus"

func recoveryHandlerFunc(error interface{}) {
	log.WithField("panic", error).Error("Recovered from panic")
}
