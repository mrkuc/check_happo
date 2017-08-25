package util

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
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

// Logger returns custom logger
func Logger() *logrus.Logger {
	return log
}

// LoggerLevelDebug set level to Debug
func LoggerLevelDebug() {
	log.Level = logrus.DebugLevel
}

// LoggerLevelInfo set level to Info
func LoggerLevelInfo() {
	log.Level = logrus.InfoLevel
}

// LoggerLevelError set level to Error
func LoggerLevelError() {
	log.Level = logrus.ErrorLevel
}

// DumpStruct dump struct to one line string
func DumpStruct(a ...interface{}) string {
	spew.Config = spew.ConfigState{
		Indent:                  "\t",
		DisableMethods:          true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
	}
	return spew.Sdump(a)
}
