package datastore

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v2"
)

type SqliteProvider struct {
	databasePath string
}

func (provider SqliteProvider) Connect() error {
	if dbMap == nil {
		if provider.databasePath == "" {
			return errors.New("not key databasePath")
		} else {
			db, err := sql.Open("sqlite3", provider.databasePath)
			if err != nil {
				return err
			}
			dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
		}
	}
	return nil
}

func (provider SqliteProvider) Init() {
	provider.UserCreateStore()
	provider.RoomCreateStore()
	provider.RoomUserCreateStore()
	provider.MessageCreateStore()
	provider.DeviceCreateStore()
	provider.SubscriptionCreateStore()
}

func (provider SqliteProvider) DropDatabase() error {
	if err := os.Remove(provider.databasePath); err != nil {
		return err
	}
	return nil
}
