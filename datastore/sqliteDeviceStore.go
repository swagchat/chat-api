package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.database)
}

func (p *sqliteProvider) InsertDevice(device *model.Device) (*model.Device, error) {
	return rdbInsertDevice(p.database, device)
}

func (p *sqliteProvider) SelectDevices(userID string) ([]*model.Device, error) {
	return rdbSelectDevices(p.database, userID)
}

func (p *sqliteProvider) SelectDevice(userID string, platform int) (*model.Device, error) {
	return rdbSelectDevice(p.database, userID, platform)
}

func (p *sqliteProvider) SelectDevicesByUserID(userID string) ([]*model.Device, error) {
	return rdbSelectDevicesByUserID(p.database, userID)
}

func (p *sqliteProvider) SelectDevicesByToken(token string) ([]*model.Device, error) {
	return rdbSelectDevicesByToken(p.database, token)
}

func (p *sqliteProvider) UpdateDevice(device *model.Device) error {
	return rdbUpdateDevice(p.database, device)
}

func (p *sqliteProvider) DeleteDevice(userID string, platform int) error {
	return rdbDeleteDevice(p.database, userID, platform)
}
