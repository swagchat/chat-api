package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createRoomStore() {
	master := RdbStore(p.database).master()
	rdbCreateRoomStore(p.ctx, master)
}

func (p *gcpSQLProvider) InsertRoom(room *model.Room, opts ...InsertRoomOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertRoom(p.ctx, master, tx, room, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *gcpSQLProvider) SelectRooms(limit, offset int32, opts ...SelectRoomsOption) ([]*model.Room, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRooms(p.ctx, replica, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectRoom(roomID string, opts ...SelectRoomOption) (*model.Room, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectRoom(p.ctx, replica, roomID, opts...)
}

func (p *gcpSQLProvider) SelectCountRooms(opts ...SelectRoomsOption) (int64, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectCountRooms(p.ctx, replica, opts...)
}

func (p *gcpSQLProvider) UpdateRoom(room *model.Room, opts ...UpdateRoomOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	err = rdbUpdateRoom(p.ctx, master, tx, room, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting user roles")
		logger.Error(err.Error())
		return err
	}

	return nil
}
