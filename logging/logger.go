package logging

import (
	"os"

	"github.com/swagchat/rtm-api/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var appLogger = NewAppLogger()

type AppLog struct {
	Kind      string `json:"kind"`
	Level     string `json:"level"`
	UserID    string `json:"userId,omitempty"`
	RoomID    string `json:"roomId,omitempty"`
	Event     string `json:"event,omitempty"`
	Client    string `json:"client,omitempty"`
	Useragent string `json:"useragent,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	Language  string `json:"language",omitempty`
	Provider  string `json:"provider,omitempty"`
	Config    string `json:"config,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
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

func Log(level zapcore.Level, al *AppLog) {
	fields := make([]zapcore.Field, 0)
	if al.Kind != "" {
		fields = append(fields, zap.String("kind", al.Kind))
	}
	if al.UserID != "" {
		fields = append(fields, zap.String("userId", al.UserID))
	}
	if al.RoomID != "" {
		fields = append(fields, zap.String("roomId", al.RoomID))
	}
	if al.Event != "" {
		fields = append(fields, zap.String("event", al.Event))
	}
	if al.Client != "" {
		fields = append(fields, zap.String("client", al.Client))
	}
	if al.Useragent != "" {
		fields = append(fields, zap.String("useragent", al.Useragent))
	}
	if al.IPAddress != "" {
		fields = append(fields, zap.String("ipAddress", al.IPAddress))
	}
	if al.Language != "" {
		fields = append(fields, zap.String("language", al.Language))
	}
	if al.Provider != "" {
		fields = append(fields, zap.String("provider", al.Provider))
	}
	if al.Config != "" {
		fields = append(fields, zap.String("config", al.Config))
	}
	if al.Message != "" {
		fields = append(fields, zap.String("message", al.Message))
	}

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
