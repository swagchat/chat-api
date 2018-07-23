package datastore

import (
	"fmt"
	"sync/atomic"

	"github.com/fairway-corp/operator-api/utils"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbStores                = make(map[string]*rdbStore)
	tableNameBot             = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "bot")
	tableNameGuestSetting    = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "guest_setting")
	tableNameOperatorSetting = fmt.Sprintf("%s%s", utils.Config().Datastore.TableNamePrefix, "operator_setting")
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
