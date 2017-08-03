package utils

import (
	"fmt"

	"go.uber.org/zap"
)

func init() {
	setupConfig()
	if IsShowVersion {
		return
	}
	setupLogger()
	AppLogger.Info("",
		zap.String("configName", "Config"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg)),
	)

	AppLogger.Info("",
		zap.String("configName", "Storage"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.Storage)),
	)

	AppLogger.Info("",
		zap.String("configName", "Datastore"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.Datastore)),
	)

	AppLogger.Info("",
		zap.String("configName", "Messaging"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.Messaging)),
	)

	AppLogger.Info("",
		zap.String("configName", "Messaging.RealtimeQue"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.Messaging.RealtimeQue)),
	)

	AppLogger.Info("",
		zap.String("configName", "Notification"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.Notification)),
	)

	AppLogger.Info("",
		zap.String("configName", "RealtimeServer"),
		zap.String("configValue", fmt.Sprintf("%#v", Cfg.RealtimeServer)),
	)
}
