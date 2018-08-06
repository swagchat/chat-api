package datastore

import (
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/logger"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbStores             = make(map[string]*rdbStore)
	tableNameAppClient    = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "app_client")
	tableNameAsset        = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "asset")
	tableNameBlockUser    = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "block_user")
	tableNameBot          = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "bot")
	tableNameDevice       = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "device")
	tableNameMessage      = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "message")
	tableNameRoom         = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "room")
	tableNameRoomUser     = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "room_user")
	tableNameSetting      = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "setting")
	tableNameSubscription = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "subscription")
	tableNameUser         = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "user")
	tableNameUserRole     = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "user_role")
	tableNameWebhook      = fmt.Sprintf("%s%s", config.Config().Datastore.TableNamePrefix, "webhook")
)

type rdbStore struct {
	masterDbMap    *gorp.DbMap
	replicaDbMaps  []*gorp.DbMap
	replicaCounter int64
}

func RdbStore(db string) *rdbStore {
	if rs, ok := rdbStores[db]; ok {
		return rs
	}
	return nil
}

func (rs *rdbStore) master() *gorp.DbMap {
	return rs.masterDbMap
}

func (rs *rdbStore) setMaster(m *gorp.DbMap) {
	rs.masterDbMap = m
}

func (rs *rdbStore) replica() *gorp.DbMap {
	if rs.replicaDbMaps == nil {
		return rs.masterDbMap
	}
	replicaCounter := atomic.AddInt64(&rs.replicaCounter, 1) % int64(len(rs.replicaDbMaps))
	return rs.replicaDbMaps[replicaCounter]
}

func (rs *rdbStore) setReplica(r *gorp.DbMap) {
	rs.replicaDbMaps = append(rs.replicaDbMaps, r)
}

func close(database string, db *sql.DB) {
	if db == nil {
		return
	}

	err := db.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to close database. %s %s %s", config.Config().Datastore.Provider, database, err.Error()))
	} else {
		logger.Info(fmt.Sprintf("Closing database. %s %s", config.Config().Datastore.Provider, database))
	}
}
