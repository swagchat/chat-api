package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
	gorp "gopkg.in/gorp.v2"
)

type sqliteProvider struct {
	ctx               context.Context
	onMemory          bool
	dirPath           string
	database          string
	maxIdleConnection int
	maxOpenConnection int
	enableLogging     bool
}

func (p *sqliteProvider) Connect(dsCfg *config.Datastore) error {
	if _, ok := rdbStores[dsCfg.Database]; ok {
		return nil
	}

	rs := &rdbStore{}
	if rs.master() != nil {
		return nil
	}

	var ds string
	if p.onMemory {
		ds = fmt.Sprintf("%s/%s.db?cache=shared&mode=memory", p.dirPath, p.database)
	} else {
		ds = fmt.Sprintf("%s/%s.db?cache=shared", p.dirPath, p.database)
	}
	db, err := sql.Open("sqlite3", ds)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Failed to connect database. %s %s", dsCfg.Provider, ds))
		logger.Error(err.Error())
		return err
	}
	logger.Info(fmt.Sprintf("Connected database. %s %s", dsCfg.Provider, ds))

	db.SetMaxIdleConns(p.maxIdleConnection)
	db.SetMaxOpenConns(p.maxOpenConnection)

	var master *gorp.DbMap
	master = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	if p.enableLogging {
		master.TraceOn("[master]", logger.Logger())
	}
	rs.setMaster(master)

	rdbStores[dsCfg.Database] = rs
	return nil
}

func (p *sqliteProvider) CreateTables() {
	p.createAppClientStore()
	p.createAssetStore()
	p.createBlockUserStore()
	p.createDeviceStore()
	p.createMessageStore()
	p.createRoomStore()
	p.createRoomUserStore()
	p.createSettingStore()
	p.createSubscriptionStore()
	p.createUserStore()
	p.createUserRoleStore()
	p.createWebhookStore()
}

func (p *sqliteProvider) DropDatabase() error {
	dbPath := fmt.Sprintf("%s/%s.db", p.dirPath, p.database)
	if err := os.Remove(dbPath); err != nil {
		err = errors.Wrap(err, "Drop database failure")
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (p *sqliteProvider) Close() {
	for database, rdbStore := range rdbStores {
		if rdbStore != nil {
			master := rdbStore.masterDbMap
			if master != nil {
				close(database, master.Db)
			}
			delete(rdbStores, database)
		}
	}
}
