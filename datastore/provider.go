package datastore

import (
	"context"

	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
)

type provider interface {
	Connect(dsCfg *utils.Datastore) error
	CreateTables()
	DropDatabase() error
	appClientStore
	assetStore
	blockUserStore
	deviceStore
	messageStore
	roomStore
	roomUserStore
	settingStore
	subscriptionStore
	userStore
	userRoleStore
	webhookStore
}

// Provider is get datastore provider
func Provider(ctx context.Context) provider {
	var p provider

	cfg := utils.Config()
	dsCfg := cfg.Datastore

	if cfg.Datastore.Dynamic {
		dsCfg.Database = ctx.Value(utils.CtxWorkspace).(string)
	}

	switch dsCfg.Provider {
	case "sqlite":
		p = &sqliteProvider{
			onMemory:      dsCfg.SQLite.OnMemory,
			dirPath:       dsCfg.SQLite.DirPath,
			database:      dsCfg.Database,
			enableLogging: dsCfg.EnableLogging,
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
			enableLogging:     dsCfg.EnableLogging,
		}
	case "gcSql":
		p = &gcpSQLProvider{
			user:              dsCfg.User,
			password:          dsCfg.Password,
			database:          dsCfg.Database,
			masterSi:          dsCfg.Master,
			replicaSis:        dsCfg.Replicas,
			maxIdleConnection: dsCfg.MaxIdleConnection,
			maxOpenConnection: dsCfg.MaxOpenConnection,
			enableLogging:     dsCfg.EnableLogging,
		}
	}

	err := p.Connect(dsCfg)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return p
}
