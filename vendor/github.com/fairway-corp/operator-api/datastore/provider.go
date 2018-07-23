package datastore

import (
	"context"

	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/pkg/errors"
)

type provider interface {
	Connect(dsCfg *utils.Datastore) error
	Init()
	DropDatabase() error
	botStore
	guestSettingStore
	operatorSettingStore
}

func Provider(ctx context.Context) provider {
	var p provider

	cfg := utils.Config()
	dsCfg := cfg.Datastore

	if cfg.Datastore.Dynamic {
		dsCfg.Database = ctx.Value(utils.CtxWorkspace).(string)
	}

	enableLogging := false
	if cfg.Datastore.EnableLogging {
		enableLogging = true
	}

	switch cfg.Datastore.Provider {
	case "sqlite":
		p = &sqliteProvider{
			onMemory:      dsCfg.SQLite.OnMemory,
			dirPath:       dsCfg.SQLite.DirPath,
			database:      dsCfg.Database,
			enableLogging: enableLogging,
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
			enableLogging:     enableLogging,
		}
	case "gcpSql":
		p = &gcpSQLProvider{
			user:              cfg.Datastore.User,
			password:          cfg.Datastore.Password,
			database:          cfg.Datastore.Database,
			masterSi:          cfg.Datastore.Master,
			replicaSis:        cfg.Datastore.Replicas,
			maxIdleConnection: cfg.Datastore.MaxIdleConnection,
			maxOpenConnection: cfg.Datastore.MaxOpenConnection,
			enableLogging:     enableLogging,
		}
	}

	err := p.Connect(dsCfg)
	if err != nil {
		err := errors.Wrap(err, "Database connect error")
		logger.Error(err.Error())
	}

	return p
}
