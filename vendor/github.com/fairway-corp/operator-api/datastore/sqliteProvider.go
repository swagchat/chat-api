package datastore

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/fairway-corp/operator-api/logger"
	"github.com/fairway-corp/operator-api/utils"
	"github.com/pkg/errors"
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
		db, err = sql.Open("sqlite3", ":memory:")
	} else {
		db, err = sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", p.dirPath, p.database))
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
	// p.init()
	return nil
}

func (p *sqliteProvider) Init() {
	p.createBotStore()
	p.CreateGuestSettingStore()
	p.CreateOperatorSettingStore()
}

func (p *sqliteProvider) DropDatabase() error {
	if err := os.Remove(p.database); err != nil {
		err = errors.Wrap(err, "Drop database failure")
		logger.Error(err.Error())
		return err
	}
	return nil
}
