package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	INFO  = "info"
	DEBUG = "debug"
)

type Logger struct {
	*logrus.Logger
}

func New(logLevel string) *Logger {
	logger := logrus.New()

	var logrusLevel logrus.Level
	logrusLevel = logrus.DebugLevel

	switch strings.ToLower(logLevel) {
	case INFO:
		logrusLevel = logrus.InfoLevel
	case DEBUG:
		logrusLevel = logrus.DebugLevel
	}
	logger.SetLevel(logrusLevel)

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})

	return &Logger{
		Logger: logger,
	}
}
