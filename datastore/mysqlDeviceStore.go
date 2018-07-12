package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.database)
}

func (p *mysqlProvider) InsertDevice(device *model.Device) (*model.Device, error) {
	return rdbInsertDevice(p.database, device)
}

func (p *mysqlProvider) SelectDevices(userID string) ([]*model.Device, error) {
	return rdbSelectDevices(p.database, userID)
}

func (p *mysqlProvider) SelectDevice(userID string, platform int32) (*model.Device, error) {
	return rdbSelectDevice(p.database, userID, platform)
}

func (p *mysqlProvider) SelectDevicesByUserID(userID string) ([]*model.Device, error) {
	return rdbSelectDevicesByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectDevicesByToken(token string) ([]*model.Device, error) {
	return rdbSelectDevicesByToken(p.database, token)
}

func (p *mysqlProvider) UpdateDevice(device *model.Device) error {
	return rdbUpdateDevice(p.database, device)
}

func (p *mysqlProvider) DeleteDevice(userID string, platform int32) error {
	return rdbDeleteDevice(p.database, userID, platform)
}
