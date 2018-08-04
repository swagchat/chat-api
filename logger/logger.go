package logger

import (
	"github.com/betchi/zapper"
	"github.com/swagchat/chat-api/config"
)

var (
	logger *zapper.Logger

	// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Debug = logger.Debug
	// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Info = logger.Info
	// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Warn = logger.Warn
	// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Error = logger.Error
)

func InitLogger(config *config.Logger) {
	logger = zapper.NewLogger(&zapper.Config{
		EnableConsole: config.EnableConsole,
		ConsoleFormat: config.ConsoleFormat,
		ConsoleLevel:  config.ConsoleLevel,
		EnableFile:    config.EnableFile,
		FileFormat:    config.FileFormat,
		FileLevel:     config.FileLevel,
		Filepath:      config.FilePath,
	})
	Debug = logger.Debug
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
}

func Logger() *zapper.Logger {
	return logger
}
