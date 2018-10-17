package gorp

import (
	"github.com/betchi/zapper"
)

// Logger is a logger for gorp
type Logger struct{}

var logger = &Logger{}

// GlobalLogger retrieve global logger for elastic apm
func GlobalLogger() *Logger {
	return logger
}

// Error logs a message at ErrorLevel for jaeger
func (l *Logger) Error(msg string) {
	if l == nil {
		return
	}

	zapper.GlobalLogger().Error(msg)
}

// Infof logs a message at InfoLevel for jaeger
func (l *Logger) Infof(msg string, args ...interface{}) {
	if l == nil {
		return
	}

	zapper.GlobalLogger().Infow(msg, args)
}
