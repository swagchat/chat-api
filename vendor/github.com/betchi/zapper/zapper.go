package zapper

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Field is an alias for zap.Field
// type Field = zapcore.Field

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
	levelFatal = "fatal"
)

var (
	// Int constructs a field with the given key and value.
	Int = zap.Int
	// String constructs a field with the given key and value.
	String = zap.String
	logger *Logger

	// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Debug = logger.Debug
	// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Info = logger.Info
	// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Warn = logger.Warn
	// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
	Error = logger.Error

	encoderConfig = newEncoderConfig()
)

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
	case levelFatal:
		return zapcore.FatalLevel
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

// InitGlobalLogger initialize global logger
func InitGlobalLogger(config *Config) {
	logger = NewLogger(&Config{
		EnableConsole: config.EnableConsole,
		ConsoleFormat: config.ConsoleFormat,
		ConsoleLevel:  config.ConsoleLevel,
		EnableFile:    config.EnableFile,
		FileFormat:    config.FileFormat,
		FileLevel:     config.FileLevel,
		FilePath:      config.FilePath,
	})
	Debug = logger.Debug
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
}

// GlobalLogger retrieve global logger
func GlobalLogger() *Logger {
	return logger
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
			Filename: config.FilePath,
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

// Logger is struct for logging.
type Logger struct {
	zap          *zap.Logger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(message string, fields ...zapcore.Field) {
	if l == nil {
		return
	}
	l.zap.Debug(message, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(message string, fields ...zapcore.Field) {
	if l == nil {
		return
	}
	l.zap.Info(message, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(message string, fields ...zapcore.Field) {
	if l == nil {
		return
	}
	l.zap.Warn(message, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(message string, fields ...zapcore.Field) {
	if l == nil {
		return
	}
	l.zap.Error(message, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Fatal(message string, fields ...zapcore.Field) {
	if l == nil {
		return
	}
	l.zap.Fatal(message, fields...)
}

// Debugf logs a message at DebugLevel.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Debug(fmt.Sprintf(format, args...))
}

// Infof logs a message at InfoLevel.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Info(fmt.Sprintf(format, args...))
}

// Warnf logs a message at WarnLevel.
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Warn(fmt.Sprintf(format, args...))
}

// Errorf logs a message at ErrorLevel.
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Error(fmt.Sprintf(format, args...))
}

// Fatalf logs a message at FatalLevel.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Fatal(fmt.Sprintf(format, args...))
}

// Printf logs a message at DebugLevel. For gorp trace message
func (l *Logger) Printf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Debug(fmt.Sprintf(format, args...))
}
