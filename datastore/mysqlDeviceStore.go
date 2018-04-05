package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (p *mysqlProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return RdbInsertDevice(device)
}

func (p *mysqlProvider) SelectDevices(userId string) ([]*models.Device, error) {
	return RdbSelectDevices(userId)
}

func (p *mysqlProvider) SelectDevice(userId string, platform int) (*models.Device, error) {
	return RdbSelectDevice(userId, platform)
}

func (p *mysqlProvider) SelectDevicesByUserId(userId string) ([]*models.Device, error) {
	return RdbSelectDevicesByUserId(userId)
}

func (p *mysqlProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return RdbSelectDevicesByToken(token)
}

func (p *mysqlProvider) UpdateDevice(device *models.Device) error {
	return RdbUpdateDevice(device)
}

func (p *mysqlProvider) DeleteDevice(userId string, platform int) error {
	return RdbDeleteDevice(userId, platform)
}
