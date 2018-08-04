package datastore

import (
	"github.com/pkg/errors"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func (p *sqliteProvider) createDeviceStore() {
	master := RdbStore(p.database).master()
	rdbCreateDeviceStore(p.ctx, master)
}

func (p *sqliteProvider) InsertDevice(device *model.Device, opts ...InsertDeviceOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while inserting device")
		logger.Error(err.Error())
		return err
	}

	err = rdbInsertDevice(p.ctx, master, tx, device, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while inserting device")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *sqliteProvider) SelectDevices(opts ...SelectDevicesOption) ([]*model.Device, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectDevices(p.ctx, replica, opts...)
}

func (p *sqliteProvider) SelectDevice(userID string, platform scpb.Platform) (*model.Device, error) {
	replica := RdbStore(p.database).replica()
	return rdbSelectDevice(p.ctx, replica, userID, platform)
}

func (p *sqliteProvider) UpdateDevice(device *model.Device) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while updating device")
		logger.Error(err.Error())
		return err
	}

	err = rdbUpdateDevice(p.ctx, master, tx, device)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while updating device")
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *sqliteProvider) DeleteDevices(opts ...DeleteDevicesOption) error {
	master := RdbStore(p.database).master()
	tx, err := master.Begin()
	if err != nil {
		err = errors.Wrap(err, "An error occurred while deleting device")
		logger.Error(err.Error())
		return err
	}

	err = rdbDeleteDevices(p.ctx, master, tx, opts...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		err = errors.Wrap(err, "An error occurred while deleting device")
		logger.Error(err.Error())
		return err
	}

	return nil
}
