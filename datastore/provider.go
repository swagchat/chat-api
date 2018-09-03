package datastore

import (
	"context"
	"os"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
)

type provider interface {
	Connect(dsCfg *config.Datastore) error
	CreateTables()
	DropDatabase() error
	Close()
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

	cfg := config.Config()
	dsCfg := cfg.Datastore

	if cfg.Datastore.Dynamic {
		dsCfg.Database = ctx.Value(config.CtxWorkspace).(string)
	}

	switch dsCfg.Provider {
	case "sqlite":
		p = &sqliteProvider{
			ctx:               ctx,
			onMemory:          dsCfg.SQLite.OnMemory,
			dirPath:           dsCfg.SQLite.DirPath,
			database:          dsCfg.Database,
			maxIdleConnection: dsCfg.MaxIdleConnection,
			maxOpenConnection: dsCfg.MaxOpenConnection,
			enableLogging:     dsCfg.EnableLogging,
		}
	case "mysql":
		p = &mysqlProvider{
			ctx:               ctx,
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
			ctx:               ctx,
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
		os.Exit(1)
	}

	return p
}
