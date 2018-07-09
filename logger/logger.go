package logger

import (
	"github.com/betchi/zapper"
)

var (
	logger = newLogger()

	// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Debug = logger.Debug
	// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Info = logger.Info
	// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Warn = logger.Warn
	// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Error = logger.Error
)

func newLogger() *zapper.Logger {
	return zapper.NewLogger(&zapper.Config{
		EnableConsole: true,
		ConsoleFormat: "text",
		ConsoleLevel:  "debug",
		EnableFile:    false,
	})
}
