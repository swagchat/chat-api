package datastore

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v2"
)

type sqliteProvider struct {
	sqlitePath string
	trace      bool
}

func (p *sqliteProvider) Connect() error {
	rs := RdbStoreInstance()
	if rs.master() != nil {
		return nil
	}
	if p.sqlitePath == "" {
		return errors.New("not key sqlitePath")
	}

	db, err := sql.Open("sqlite3", p.sqlitePath)
	if err != nil {
		fatal(err)
	}
	var master *gorp.DbMap
	master = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	if p.trace {
		master.TraceOn("", log.New(os.Stdout, "sql-trace:", log.Lmicroseconds))
	}
	rs.setMaster(master)
	return nil
}

func (p *sqliteProvider) Init() {
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
