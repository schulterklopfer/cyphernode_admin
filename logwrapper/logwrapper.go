package logwrapper

import (
  "github.com/sirupsen/logrus"
  "sync"
)

// Event stores messages to log later, from our standard interface
type Event struct {
  id      int
  message string
}

// standardLogger enforces specific log message formats
var standardLogger *logrus.Logger
var once sync.Once

// NewLogger initializes the standard logger
func initOnce() {
  once.Do(func() {
    standardLogger = logrus.New()
    standardLogger.Formatter = &logrus.TextFormatter{}
  })
}

func Logger() *logrus.Logger {
  if standardLogger == nil {
    initOnce()
  }
  standardLogger.SetReportCaller(true)
  return standardLogger
}

// Declare variables to store log messages as new Events
var (
  invalidArgMessage      = Event{1, "Invalid arg: %s"}
  invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
  missingArgMessage      = Event{3, "Missing arg: %s"}
)

// InvalidArg is a standard error message
func InvalidArg(argumentName string) {
  standardLogger.Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func InvalidArgValue(argumentName string, argumentValue string) {
  standardLogger.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func  MissingArg(argumentName string) {
  standardLogger.Errorf(missingArgMessage.message, argumentName)
}