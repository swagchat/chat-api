package datastore

import (
	"fmt"
	"sync/atomic"

	"github.com/swagchat/chat-api/utils"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbStores             = make(map[string]*rdbStore)
	tableNameAppClient    = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "app_client")
	tableNameAsset        = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "asset")
	tableNameBlockUser    = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "block_user")
	tableNameBot          = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "bot")
	tableNameDevice       = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "device")
	tableNameMessage      = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "message")
	tableNameRoom         = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "room")
	tableNameRoomUser     = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "room_user")
	tableNameSetting      = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "setting")
	tableNameSubscription = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "subscription")
	tableNameUser         = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "user")
	tableNameUserRole     = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "user_role")
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
