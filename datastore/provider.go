package datastore

import (
	"github.com/swagchat/chat-api/utils"
)

type provider interface {
	Connect() error
	Init()
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

func Provider() provider {
	cfg := utils.Config()
	var p provider

	switch cfg.Datastore.Provider {
	case "sqlite":
		p = &sqliteProvider{
			sqlitePath: cfg.Datastore.SQLite.Path,
			trace:      false,
		}
	case "mysql":
		p = &mysqlProvider{
			user:              cfg.Datastore.User,
			password:          cfg.Datastore.Password,
			database:          cfg.Datastore.Database,
			masterSi:          cfg.Datastore.Master,
			replicaSis:        cfg.Datastore.Replicas,
			maxIdleConnection: cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: cfg.Datastore.MaxOpenConnection,
			trace:             false,
		}
	case "gcpSql":
		p = &gcpSqlProvider{
			user:              cfg.Datastore.User,
			password:          cfg.Datastore.Password,
			database:          cfg.Datastore.Database,
			masterSi:          cfg.Datastore.Master,
			replicaSis:        cfg.Datastore.Replicas,
			maxIdleConnection: cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: cfg.Datastore.MaxOpenConnection,
			trace:             false,
		}
	}

	return p
}
