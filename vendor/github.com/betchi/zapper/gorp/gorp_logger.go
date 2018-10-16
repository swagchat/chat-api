package gorp

import (
	"fmt"

	"github.com/betchi/zapper"
)

// Logger is a logger for gorp
type Logger struct{}

var logger = &Logger{}

// GlobalLogger retrieve global logger for elastic apm
func GlobalLogger() *Logger {
	return logger
}

// Printf logs a message at DebugLevel for gorp
func (l *Logger) Printf(format string, args ...interface{}) {
	if l == nil {
		return
	}

	zapper.GlobalLogger().Info(fmt.Sprintf(format, args...))
}
