package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.database)
}

func (p *mysqlProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return rdbInsertDevice(p.database, device)
}

func (p *mysqlProvider) SelectDevices(userID string) ([]*models.Device, error) {
	return rdbSelectDevices(p.database, userID)
}

func (p *mysqlProvider) SelectDevice(userID string, platform int) (*models.Device, error) {
	return rdbSelectDevice(p.database, userID, platform)
}

func (p *mysqlProvider) SelectDevicesByUserID(userID string) ([]*models.Device, error) {
	return rdbSelectDevicesByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return rdbSelectDevicesByToken(p.database, token)
}

func (p *mysqlProvider) UpdateDevice(device *models.Device) error {
	return rdbUpdateDevice(p.database, device)
}

func (p *mysqlProvider) DeleteDevice(userID string, platform int) error {
	return rdbDeleteDevice(p.database, userID, platform)
}
