package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) == 0 {
		logLevel = "debug"
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithFields(log.Fields{}).Error()
	}
	log.SetLevel(level)
}
