package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
	gorp "gopkg.in/gorp.v2"
)

type sqliteProvider struct {
	dirPath  string
	database string
	trace    bool
}

func (p *sqliteProvider) Connect(dsCfg *utils.Datastore) error {
	if _, ok := rdbStores[dsCfg.Database]; ok {
		return nil
	}

	rs := &rdbStore{}
	if rs.master() != nil {
		return nil
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", p.dirPath, p.database))
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

	rdbStores[dsCfg.Database] = rs
	p.init()
	return nil
}

func (p *sqliteProvider) init() {
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
		return errors.Wrap(err, "Drop database failure")
	}
	return nil
}
