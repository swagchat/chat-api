package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.database)
}

func (p *gcpSQLProvider) InsertDevice(device *model.Device) (*model.Device, error) {
	return rdbInsertDevice(p.database, device)
}

func (p *gcpSQLProvider) SelectDevices(userID string) ([]*model.Device, error) {
	return rdbSelectDevices(p.database, userID)
}

func (p *gcpSQLProvider) SelectDevice(userID string, platform int) (*model.Device, error) {
	return rdbSelectDevice(p.database, userID, platform)
}

func (p *gcpSQLProvider) SelectDevicesByUserID(userID string) ([]*model.Device, error) {
	return rdbSelectDevicesByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectDevicesByToken(token string) ([]*model.Device, error) {
	return rdbSelectDevicesByToken(p.database, token)
}

func (p *gcpSQLProvider) UpdateDevice(device *model.Device) error {
	return rdbUpdateDevice(p.database, device)
}

func (p *gcpSQLProvider) DeleteDevice(userID string, platform int) error {
	return rdbDeleteDevice(p.database, userID, platform)
}
