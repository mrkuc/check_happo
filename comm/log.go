package comm

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Out = os.Stderr

	formatter := &logrus.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.DisableSorting = true
	log.Formatter = formatter
}

func Logger() *logrus.Logger {
	return log
}

func LoggerLevelDebug() {
	log.Level = logrus.DebugLevel
}

func LoggerLevelInfo() {
	log.Level = logrus.InfoLevel
}

func LoggerLevelError() {
	log.Level = logrus.ErrorLevel
}
