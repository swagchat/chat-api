package datastore

import (
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

type provider interface {
	Connect(dsCfg *utils.Datastore) error
	DropDatabase() error
	ApiStore
	AssetStore
	BlockUserStore
	BotStore
	DeviceStore
	MessageStore
	RoomStore
	RoomUserStore
	SettingStore
	SubscriptionStore
	UserStore
}

func Provider(dsCfg *utils.Datastore) provider {
	var p provider

	switch dsCfg.Provider {
	case "sqlite":
		p = &sqliteProvider{
			sqlitePath: dsCfg.SQLite.Path,
			trace:      false,
		}
	case "mysql":
		p = &mysqlProvider{
			user:              dsCfg.User,
			password:          dsCfg.Password,
			database:          dsCfg.Database,
			masterSi:          dsCfg.Master,
			replicaSis:        dsCfg.Replicas,
			maxIdleConnection: dsCfg.MaxIdleConnection,
			maxOpenConnection: dsCfg.MaxOpenConnection,
			trace:             false,
		}
	case "gcSql":
		p = &gcpSqlProvider{
			user:              dsCfg.User,
			password:          dsCfg.Password,
			database:          dsCfg.Database,
			masterSi:          dsCfg.Master,
			replicaSis:        dsCfg.Replicas,
			maxIdleConnection: dsCfg.MaxIdleConnection,
			maxOpenConnection: dsCfg.MaxOpenConnection,
			trace:             false,
		}
	}

	err := p.Connect(dsCfg)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Message: "Database connect error",
			Error:   err,
		})
	}

	return p
}
