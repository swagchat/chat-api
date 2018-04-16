package datastore

import (
	"context"

	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

type provider interface {
	Connect(dsCfg *utils.Datastore) error
	DropDatabase() error
	appClientStore
	assetStore
	blockUserStore
	botStore
	deviceStore
	messageStore
	roomStore
	roomUserStore
	settingStore
	subscriptionStore
	userStore
}

// Provider is get datastore provider
func Provider(ctx context.Context) provider {
	var p provider

	// ctxDsCfg := ctx.Value(utils.CtxDsCfg)
	// if ctxDsCfg == nil {
	// 	logging.Log(zapcore.ErrorLevel, &logging.AppLog{
	// 		Message: "Database connect error. Database config is nil",
	// 	})
	// }
	// dsCfg := ctxDsCfg.(*utils.Datastore)

	cfg := utils.Config()
	dsCfg := cfg.Datastore

	if cfg.Datastore.Dynamic {
		dsCfg.Database = ctx.Value(utils.CtxRealm).(string)
	}

	switch dsCfg.Provider {
	case "sqlite":
		p = &sqliteProvider{
			sqlitePath: dsCfg.SQLite.Path,
			database:   dsCfg.Database,
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
		p = &gcpSQLProvider{
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
