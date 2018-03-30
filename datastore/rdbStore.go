package datastore

import (
	"sync/atomic"

	"github.com/swagchat/chat-api/utils"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbStoreInstance        *rdbStore = nil
	TABLE_NAME_API                    = utils.Config().Datastore.TableNamePrefix + "api"
	TABLE_NAME_ASSET                  = utils.Config().Datastore.TableNamePrefix + "asset"
	TABLE_NAME_BLOCK_USER             = utils.Config().Datastore.TableNamePrefix + "block_user"
	TABLE_NAME_BOT                    = utils.Config().Datastore.TableNamePrefix + "bot"
	TABLE_NAME_DEVICE                 = utils.Config().Datastore.TableNamePrefix + "device"
	TABLE_NAME_MESSAGE                = utils.Config().Datastore.TableNamePrefix + "message"
	TABLE_NAME_ROOM                   = utils.Config().Datastore.TableNamePrefix + "room"
	TABLE_NAME_ROOM_USER              = utils.Config().Datastore.TableNamePrefix + "room_user"
	TABLE_NAME_SETTING                = utils.Config().Datastore.TableNamePrefix + "setting"
	TABLE_NAME_SUBSCRIPTION           = utils.Config().Datastore.TableNamePrefix + "subscription"
	TABLE_NAME_USER                   = utils.Config().Datastore.TableNamePrefix + "user"
)

type rdbStore struct {
	masterDbMap    *gorp.DbMap
	replicaDbMaps  []*gorp.DbMap
	replicaCounter int64
}

func RdbStoreInstance() *rdbStore {
	if rdbStoreInstance == nil {
		rdbStoreInstance = &rdbStore{}
	}
	return rdbStoreInstance
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
