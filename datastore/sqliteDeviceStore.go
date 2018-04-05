package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (p *sqliteProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return RdbInsertDevice(device)
}

func (p *sqliteProvider) SelectDevices(userId string) ([]*models.Device, error) {
	return RdbSelectDevices(userId)
}

func (p *sqliteProvider) SelectDevice(userId string, platform int) (*models.Device, error) {
	return RdbSelectDevice(userId, platform)
}

func (p *sqliteProvider) SelectDevicesByUserId(userId string) ([]*models.Device, error) {
	return RdbSelectDevicesByUserId(userId)
}

func (p *sqliteProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return RdbSelectDevicesByToken(token)
}

func (p *sqliteProvider) UpdateDevice(device *models.Device) error {
	return RdbUpdateDevice(device)
}

func (p *sqliteProvider) DeleteDevice(userId string, platform int) error {
	return RdbDeleteDevice(userId, platform)
}
