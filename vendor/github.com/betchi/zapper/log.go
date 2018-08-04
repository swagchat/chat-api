package zapper

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
)

// Field is an alias for zap.Field
type Field = zapcore.Field

var (
	// Int constructs a field with the given key and value.
	Int = zap.Int
	// String constructs a field with the given key and value.
	String = zap.String

	encoderConfig = newEncoderConfig()
)

// Config is a logging configuration.
type Config struct {
	// EnableConsole is a flag for enable console log.
	EnableConsole bool
	// ConsoleFormat is a format for console log.
	ConsoleFormat string
	// ConsoleLevel is a level for console log.
	ConsoleLevel string
	// EnableFile is a flag for enable file log.
	EnableFile bool
	// FileFormat is a format for file log.
	FileFormat string
	// FileLevel is a log level for file log.
	FileLevel string
	// FilePath is a file path for file log.
	Filepath string
}

// Logger is struct for logging.
type Logger struct {
	zap          *zap.Logger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
}

func zapLevel(level string) zapcore.Level {
	switch level {
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarn:
		return zapcore.WarnLevel
	case levelDebug:
		return zapcore.DebugLevel
	case levelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewLogger builds a Logger
func NewLogger(config *Config) *Logger {
	cores := []zapcore.Core{}
	logger := &Logger{
		consoleLevel: zap.NewAtomicLevelAt(zapLevel(config.ConsoleLevel)),
		fileLevel:    zap.NewAtomicLevelAt(zapLevel(config.FileLevel)),
	}

	if config.EnableConsole {
		writer := zapcore.Lock(os.Stdout)
		var encoder zapcore.Encoder
		switch config.ConsoleFormat {
		case "json":
			encoder = jsonEncoder()
		case "text":
			fallthrough
		default:
			encoder = planeTextEncoder()
		}
		core := zapcore.NewCore(encoder, writer, logger.consoleLevel)
		cores = append(cores, core)
	}

	if config.EnableFile {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename: config.Filepath,
			MaxSize:  100,
			Compress: true,
		})
		var encoder zapcore.Encoder
		switch config.FileFormat {
		case "json":
			encoder = jsonEncoder()
		case "text":
			fallthrough
		default:
			encoder = planeTextEncoder()
		}
		core := zapcore.NewCore(encoder, writer, logger.fileLevel)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger.zap = zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return logger
}

func planeTextEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func jsonEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(encoderConfig)
}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(message string, fields ...Field) {
	l.zap.Debug(message, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(message string, fields ...Field) {
	l.zap.Info(message, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(message string, fields ...Field) {
	l.zap.Warn(message, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(message string, fields ...Field) {
	l.zap.Error(message, fields...)
}

// Debugf logs a message at DebugLevel.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zap.Debug(fmt.Sprintf(format, args))
}

// Infof logs a message at InfoLevel.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.zap.Info(fmt.Sprintf(format, args))
}

// Warnf logs a message at WarnLevel.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zap.Warn(fmt.Sprintf(format, args))
}

// Errorf logs a message at ErrorLevel.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zap.Error(fmt.Sprintf(format, args))
}

// Printf logs a message at DebugLevel. For gorp trace message
func (l *Logger) Printf(format string, args ...interface{}) {
	l.zap.Debug(fmt.Sprintf(format, args))
}
