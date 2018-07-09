package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	appLogger = NewAppLogger()
	locale    = NewTimezone()
)

type AppLog struct {
	// Info/Debug
	Kind string `json:"kind"`
	// Level      string `json:"level"`
	// UserID     string `json:"userId,omitempty"`
	// RoomID     string `json:"roomId,omitempty"`
	// Event      string `json:"event,omitempty"`
	// Client     string `json:"client,omitempty"`
	// Useragent  string `json:"useragent,omitempty"`
	// IPAddress  string `json:"ipAddress,omitempty"`
	// Language   string `json:"language",omitempty`
	// Provider   string `json:"provider,omitempty"`
	Config  string `json:"config,omitempty"`
	Message string `json:"message,omitempty"`

	// Error
	ProblemDetail *models.ProblemDetail `json:"problemDetail,omitempty"`
	Error         error                 `json:"_"`
	Stacktrace    string                `json:"stacktrace,omitempty"`
	Timestamp     string                `json:"timestamp"`
	File          string                `json:"file"`
	Line          int                   `json:"line"`
}

func NewAppLogger() *zap.Logger {
	c := utils.Config()

	var err error
	var logger *zap.Logger
	if c.Logging.Level == "production" {
		logger, err = zap.NewProduction()
	} else if c.Logging.Level == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	hostname, _ := os.Hostname()
	appLogger := logger.WithOptions(zap.Fields(
		zap.String("appName", utils.AppName),
		zap.String("apiVersion", utils.APIVersion),
		zap.String("buildVersion", utils.BuildVersion),
		zap.String("hostname", hostname),
	))

	return appLogger
}

func NewTimezone() *time.Location {
	locale, _ := time.LoadLocation("Etc/GMT")
	return locale
}

func Log(level zapcore.Level, al *AppLog) {
	fields := make([]zapcore.Field, 0)
	if al.Kind != "" {
		fields = append(fields, zap.String("kind", al.Kind))
	}
	if al.Config != "" {
		fields = append(fields, zap.String("config", al.Config))
	}
	if al.Message != "" {
		fields = append(fields, zap.String("message", al.Message))
	}
	if al.ProblemDetail != nil {
		if al.ProblemDetail.Status != 0 {
			fields = append(fields, zap.String("problem.status", fmt.Sprintf("%d", al.ProblemDetail.Status)))
		}
		if al.ProblemDetail.Title != "" {
			fields = append(fields, zap.String("problem.title", al.ProblemDetail.Title))
		}
		if al.ProblemDetail.Message != "" {
			fields = append(fields, zap.String("problem.detail", al.ProblemDetail.Message))
		}
		if al.ProblemDetail.InvalidParams != nil {
			for i, invalidParam := range al.ProblemDetail.InvalidParams {
				fields = append(fields, zap.String(
					fmt.Sprintf("problem.invalid%d", i),
					fmt.Sprintf("%s: %s", invalidParam.Name, invalidParam.Reason),
				))
			}
		}
	}
	if al.Error != nil {
		al.Stacktrace = fmt.Sprintf("%+v\n", al.Error)
	}
	if al.Stacktrace != "" {
		fields = append(fields, zap.String("stacktrace", al.Stacktrace))
	}

	if al.File != "" {
		fields = append(fields, zap.String("file", al.File))
	}
	if al.Line != 0 {
		fields = append(fields, zap.Int("line", al.Line))
	}

	al.Timestamp = time.Unix(time.Now().Unix(), 0).In(locale).Format(time.RFC3339)
	switch level {
	case zapcore.DebugLevel:
		fields = append(fields, zap.String("level", "debug"))
		appLogger.Debug("", fields...)
	case zapcore.InfoLevel:
		fields = append(fields, zap.String("level", "info"))
		appLogger.Info("", fields...)
	case zapcore.WarnLevel:
		fields = append(fields, zap.String("level", "warn"))
		appLogger.Warn("", fields...)
	case zapcore.ErrorLevel:
		fields = append(fields, zap.String("level", "error"))
		appLogger.Error("", fields...)
	case zapcore.DPanicLevel:
		fields = append(fields, zap.String("level", "dpanic"))
		appLogger.DPanic("", fields...)
	case zapcore.PanicLevel:
		fields = append(fields, zap.String("level", "panic"))
		appLogger.Panic("", fields...)
	case zapcore.FatalLevel:
		fields = append(fields, zap.String("level", "fatal"))
		appLogger.Fatal("", fields...)
	}
}
