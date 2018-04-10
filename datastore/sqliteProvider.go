package datastore

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
	gorp "gopkg.in/gorp.v2"
)

type sqliteProvider struct {
	sqlitePath string
	trace      bool
}

func (p *sqliteProvider) Connect(dsCfg *utils.Datastore) error {
	if _, ok := rdbStores[dsCfg.Database]; ok {
		return nil
	}

	rs := &rdbStore{}
	if rs.master() != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", p.sqlitePath)
	if err != nil {
		logging.Log(zapcore.FatalLevel, &logging.AppLog{
			Message: "SQLite connect error",
			Error:   err,
		})
	}
	var master *gorp.DbMap
	master = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	if p.trace {
		master.TraceOn("", log.New(os.Stdout, "sql-trace:", log.Lmicroseconds))
	}
	rs.setMaster(master)

	rdbStores[dsCfg.SQLite.Path] = rs
	p.init()
	return nil
}

func (p *sqliteProvider) init() {
	p.CreateApiStore()
	p.CreateAssetStore()
	p.CreateUserStore()
	p.CreateBlockUserStore()
	p.CreateBotStore()
	p.CreateRoomStore()
	p.CreateRoomUserStore()
	p.CreateMessageStore()
	p.CreateDeviceStore()
	p.CreateSettingStore()
	p.CreateSubscriptionStore()
}

func (p *sqliteProvider) DropDatabase() error {
	if err := os.Remove(p.sqlitePath); err != nil {
		return err
	}
	return nil
}
