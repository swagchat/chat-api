package datastore

import (
	"github.com/swagchat/chat-api/utils"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbStoreInstance        *rdbStore = nil
	TABLE_NAME_API                    = utils.Cfg.Datastore.TableNamePrefix + "api"
	TABLE_NAME_USER                   = utils.Cfg.Datastore.TableNamePrefix + "user"
	TABLE_NAME_BLOCK_USER             = utils.Cfg.Datastore.TableNamePrefix + "block_user"
	TABLE_NAME_ROOM                   = utils.Cfg.Datastore.TableNamePrefix + "room"
	TABLE_NAME_ROOM_USER              = utils.Cfg.Datastore.TableNamePrefix + "room_user"
	TABLE_NAME_MESSAGE                = utils.Cfg.Datastore.TableNamePrefix + "message"
	TABLE_NAME_DEVICE                 = utils.Cfg.Datastore.TableNamePrefix + "device"
	TABLE_NAME_SUBSCRIPTION           = utils.Cfg.Datastore.TableNamePrefix + "subscription"
)

type rdbStore struct {
	master *gorp.DbMap
	slave  *gorp.DbMap
}

func RdbStoreInstance() *rdbStore {
	if rdbStoreInstance == nil {
		rdbStoreInstance = &rdbStore{}
	}
	return rdbStoreInstance
}

func (rs *rdbStore) Master() *gorp.DbMap {
	return rs.master
}

func (rs *rdbStore) SetMaster(m *gorp.DbMap) {
	rs.master = m
}

func (rs *rdbStore) Slave() *gorp.DbMap {
	if rs.slave == nil {
		return rs.master
	}
	return rs.slave
}

func (rs *rdbStore) SetSlave(s *gorp.DbMap) {
	rs.slave = s
}
