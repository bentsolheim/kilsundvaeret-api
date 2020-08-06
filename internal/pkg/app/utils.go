package app

import (
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
)

func ConfigureLogging(logLevel string) error {
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return stacktrace.Propagate(err, "error while parsing log level %s", logLevel)
	}
	log.SetLevel(level)
	return nil
}
