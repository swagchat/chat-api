package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return rdbInsertDevice(p.sqlitePath, device)
}

func (p *sqliteProvider) SelectDevices(userID string) ([]*models.Device, error) {
	return rdbSelectDevices(p.sqlitePath, userID)
}

func (p *sqliteProvider) SelectDevice(userID string, platform int) (*models.Device, error) {
	return rdbSelectDevice(p.sqlitePath, userID, platform)
}

func (p *sqliteProvider) SelectDevicesByUserID(userID string) ([]*models.Device, error) {
	return rdbSelectDevicesByUserID(p.sqlitePath, userID)
}

func (p *sqliteProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return rdbSelectDevicesByToken(p.sqlitePath, token)
}

func (p *sqliteProvider) UpdateDevice(device *models.Device) error {
	return rdbUpdateDevice(p.sqlitePath, device)
}

func (p *sqliteProvider) DeleteDevice(userID string, platform int) error {
	return rdbDeleteDevice(p.sqlitePath, userID, platform)
}
