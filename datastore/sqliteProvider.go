package datastore

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/utils"
	gorp "gopkg.in/gorp.v2"
)

type sqliteProvider struct {
	onMemory      bool
	dirPath       string
	database      string
	enableLogging bool
}

func (p *sqliteProvider) Connect(dsCfg *utils.Datastore) error {
	if _, ok := rdbStores[dsCfg.Database]; ok {
		return nil
	}

	rs := &rdbStore{}
	if rs.master() != nil {
		return nil
	}

	var db *sql.DB
	var err error
	if p.onMemory {
		db, err = sql.Open("sqlite3", fmt.Sprintf("%s/%s.db?cache=shared&mode=memory", p.dirPath, p.database))
	} else {
		db, err = sql.Open("sqlite3", fmt.Sprintf("%s/%s.db?cache=shared", p.dirPath, p.database))
	}
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	var master *gorp.DbMap
	master = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	if p.enableLogging {
		master.TraceOn("[master]", logger.Logger())
	}
	rs.setMaster(master)

	rdbStores[dsCfg.Database] = rs
	p.CreateTables()
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
	if err := os.Remove(p.database); err != nil {
		err = errors.Wrap(err, "Drop database failure")
		logger.Error(err.Error())
		return err
	}
	return nil
}
