package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func (p *mysqlProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertDevice(device *model.Device) (*model.Device, error) {
	return rdbInsertDevice(p.ctx, p.database, device)
}

func (p *mysqlProvider) SelectDevices(opts ...SelectDevicesOption) ([]*model.Device, error) {
	return rdbSelectDevices(p.ctx, p.database, opts...)
}

func (p *mysqlProvider) SelectDevice(userID string, platform scpb.Platform) (*model.Device, error) {
	return rdbSelectDevice(p.ctx, p.database, userID, platform)
}

func (p *mysqlProvider) SelectDevicesByUserID(userID string) ([]*model.Device, error) {
	return rdbSelectDevicesByUserID(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectDevicesByToken(token string) ([]*model.Device, error) {
	return rdbSelectDevicesByToken(p.ctx, p.database, token)
}

func (p *mysqlProvider) UpdateDevice(device *model.Device) error {
	return rdbUpdateDevice(p.ctx, p.database, device)
}

func (p *mysqlProvider) DeleteDevice(userID string, platform scpb.Platform) error {
	return rdbDeleteDevice(p.ctx, p.database, userID, platform)
}
