package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore(p.database)
}

func (p *mysqlProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return RdbInsertDevice(p.database, device)
}

func (p *mysqlProvider) SelectDevices(userId string) ([]*models.Device, error) {
	return RdbSelectDevices(p.database, userId)
}

func (p *mysqlProvider) SelectDevice(userId string, platform int) (*models.Device, error) {
	return RdbSelectDevice(p.database, userId, platform)
}

func (p *mysqlProvider) SelectDevicesByUserId(userId string) ([]*models.Device, error) {
	return RdbSelectDevicesByUserId(p.database, userId)
}

func (p *mysqlProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return RdbSelectDevicesByToken(p.database, token)
}

func (p *mysqlProvider) UpdateDevice(device *models.Device) error {
	return RdbUpdateDevice(p.database, device)
}

func (p *mysqlProvider) DeleteDevice(userId string, platform int) error {
	return RdbDeleteDevice(p.database, userId, platform)
}
