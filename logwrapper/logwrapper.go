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

// StandardLogger enforces specific log message formats
var StandardLogger *logrus.Logger
var once sync.Once

// NewLogger initializes the standard logger
func Init() {
  once.Do(func() {
    StandardLogger = logrus.New()
    StandardLogger.Formatter = &logrus.JSONFormatter{}
  })
}

// Declare variables to store log messages as new Events
var (
  invalidArgMessage      = Event{1, "Invalid arg: %s"}
  invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
  missingArgMessage      = Event{3, "Missing arg: %s"}
)

// InvalidArg is a standard error message
func InvalidArg(argumentName string) {
  StandardLogger.Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func InvalidArgValue(argumentName string, argumentValue string) {
  StandardLogger.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func  MissingArg(argumentName string) {
  StandardLogger.Errorf(missingArgMessage.message, argumentName)
}