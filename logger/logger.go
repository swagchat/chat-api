package logger

import (
	"github.com/betchi/zapper"
	"github.com/swagchat/chat-api/utils"
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
	cfg := utils.Config()
	return zapper.NewLogger(&zapper.Config{
		EnableConsole: cfg.Logger.EnableConsole,
		ConsoleFormat: cfg.Logger.ConsoleFormat,
		ConsoleLevel:  cfg.Logger.ConsoleLevel,
		EnableFile:    cfg.Logger.EnableFile,
		FileFormat:    cfg.Logger.FileFormat,
		FileLevel:     cfg.Logger.FileLevel,
		Filepath:      cfg.Logger.FilePath,
	})
}

func Logger() *zapper.Logger {
	return logger
}
