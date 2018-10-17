package gorp

import (
	"fmt"

	"github.com/betchi/zapper"
)

// Config is settings of jaeger logger
type Config struct {
	Noop bool
}

// Logger is a logger for gorp
type Logger struct {
	Noop bool
}

var logger *Logger

// InitGlobalLogger initialize global logger for jaeger
func InitGlobalLogger(config *Config) {
	logger = &Logger{
		Noop: config.Noop,
	}
}

// GlobalLogger retrieve global logger for jaeger
func GlobalLogger() *Logger {
	return logger
}

// Error logs a message at ErrorLevel for jaeger
func (l *Logger) Error(msg string) {
	if l == nil || l.Noop {
		return
	}

	zapper.GlobalLogger().Error(msg)
}

// Infof logs a message at InfoLevel for jaeger
func (l *Logger) Infof(msg string, args ...interface{}) {
	if l == nil || l.Noop {
		return
	}

	zapper.GlobalLogger().Info(fmt.Sprintf(msg, args...))
}
