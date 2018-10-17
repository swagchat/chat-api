package elasticapm

import (
	"fmt"

	"github.com/betchi/zapper"
)

// Config is settings of elastic apm logger
type Config struct {
	Noop bool
}

// Logger is a logger for elastic apm
type Logger struct {
	Noop bool
}

var logger *Logger

// InitGlobalLogger initialize global logger for elastic apm
func InitGlobalLogger(config *Config) {
	logger = &Logger{
		Noop: config.Noop,
	}
}

// GlobalLogger retrieve global logger for elastic apm
func GlobalLogger() *Logger {
	return logger
}

// Debugf logs a message at DebugLevel for elastic apm
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l == nil || l.Noop {
		return
	}

	zapper.GlobalLogger().Debug(fmt.Sprintf(format, args...))
}

// Errorf logs a message at ErrorLevel for elastic apm
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l == nil || l.Noop {
		return
	}

	zapper.GlobalLogger().Error(fmt.Sprintf(format, args...))
}
