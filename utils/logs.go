package utils

import "github.com/sirupsen/logrus"

// LogLevel defines a logging level that is passed to logging functions
// to define what type of log this is
type LogLevel string

const (
	// LogLevelDebug should be used to record things that only show in debug mode
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo should be used record non-failure non-important actions
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn should be used to record failure non-important actions
	LogLevelWarn LogLevel = "warn"
	// LogLevelError should be used to record failure important actions
	LogLevelError LogLevel = "error"
	// LogLevelFatal should be used log and then exit the program with an abnormal exit code
	LogLevelFatal LogLevel = "fatal"
)

// StructuredLog allows the application to log to stdout, json formatted,
//  levels accepted are debug, info, warn, error, and fatal
func StructuredLog(level LogLevel, description string, object interface{}) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logger := logrus.WithFields(logrus.Fields{
		"data": object,
	})

	switch level {
	case LogLevelDebug:
		logger.Debugln(description)
	case LogLevelInfo:
		logger.Infoln(description)
	case LogLevelWarn:
		logger.Warnln(description)
	case LogLevelError:
		logger.Errorln(description)
	case LogLevelFatal:
		logger.Fatalln(description)
	default:
		logger.Infoln(description)
	}
}
