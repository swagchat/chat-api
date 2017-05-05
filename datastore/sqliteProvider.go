package datastore

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v2"
)

type SqliteProvider struct {
	sqlitePath string
}

func (provider SqliteProvider) Connect() error {
	if dbMap == nil {
		if provider.sqlitePath == "" {
			return errors.New("not key sqlitePath")
		} else {
			db, err := sql.Open("sqlite3", provider.sqlitePath)
			if err != nil {
				return err
			}
			dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
		}
	}
	return nil
}

func (provider SqliteProvider) Init() {
	provider.CreateUserStore()
	provider.CreateRoomStore()
	provider.CreateRoomUserStore()
	provider.CreateMessageStore()
	provider.CreateDeviceStore()
	provider.CreateSubscriptionStore()
}

func (provider SqliteProvider) DropDatabase() error {
	if err := os.Remove(provider.sqlitePath); err != nil {
		return err
	}
	return nil
}
