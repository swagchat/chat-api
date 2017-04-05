package utils

import (
	"fmt"

	"go.uber.org/zap"
)

func init() {
	setupConfig()
	settingFlag()
	setupLogger()

	AppLogger.Info("",
		zap.String("configName", "ApiServer"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.ApiServer)),
	)

	AppLogger.Info("",
		zap.String("configName", "RealtimeServer"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.RealtimeServer)),
	)

	switch Cfg.ApiServer.Storage {
	case "local":
		AppLogger.Info("",
			zap.String("configName", "LocalStorage"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.LocalStorage)),
		)
	case "gcpStorage":
		AppLogger.Info("",
			zap.String("configName", "GcpStorage"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.GcpStorage)),
		)
	case "awsS3":
		AppLogger.Info("",
			zap.String("configName", "AwsS3"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.AwsS3)),
		)
	}

	switch Cfg.ApiServer.Datastore {
	case "sqlite":
		AppLogger.Info("",
			zap.String("configName", "Sqlite"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.Sqlite)),
		)
	case "mysql":
		AppLogger.Info("",
			zap.String("configName", "Mysql"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.Mysql)),
		)
	case "gcpSql":
		AppLogger.Info("",
			zap.String("configName", "GcpSql"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.GcpSql)),
		)
	}

	switch Cfg.ApiServer.Messaging {
	case "gcpPubsub":
		AppLogger.Info("",
			zap.String("configName", "GcpPubsub"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.GcpPubsub)),
		)
	}

	switch Cfg.ApiServer.Notification {
	case "awsSns":
		AppLogger.Info("",
			zap.String("configName", "AwsSns"),
			zap.String("configValue", fmt.Sprintf("%#v", Cfg.AwsSns)),
		)
	}
}
