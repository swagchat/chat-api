package utils

import (
	"os"

	"go.uber.org/zap"
)

// AppLogger is global application logger
var AppLogger *zap.Logger

func SetupLogger() {
	c := GetConfig()

	if AppLogger == nil {
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
		AppLogger = logger.WithOptions(zap.Fields(
			zap.String("appName", AppName),
			zap.String("apiVersion", APIVersion),
			zap.String("buildVersion", BuildVersion),
		))
	}
}
