package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateDeviceStore() {
	RdbCreateDeviceStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return RdbInsertDevice(p.sqlitePath, device)
}

func (p *sqliteProvider) SelectDevices(userId string) ([]*models.Device, error) {
	return RdbSelectDevices(p.sqlitePath, userId)
}

func (p *sqliteProvider) SelectDevice(userId string, platform int) (*models.Device, error) {
	return RdbSelectDevice(p.sqlitePath, userId, platform)
}

func (p *sqliteProvider) SelectDevicesByUserId(userId string) ([]*models.Device, error) {
	return RdbSelectDevicesByUserId(p.sqlitePath, userId)
}

func (p *sqliteProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return RdbSelectDevicesByToken(p.sqlitePath, token)
}

func (p *sqliteProvider) UpdateDevice(device *models.Device) error {
	return RdbUpdateDevice(p.sqlitePath, device)
}

func (p *sqliteProvider) DeleteDevice(userId string, platform int) error {
	return RdbDeleteDevice(p.sqlitePath, userId, platform)
}
