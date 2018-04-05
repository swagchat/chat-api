package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (p *gcpSqlProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return RdbInsertDevice(device)
}

func (p *gcpSqlProvider) SelectDevices(userId string) ([]*models.Device, error) {
	return RdbSelectDevices(userId)
}

func (p *gcpSqlProvider) SelectDevice(userId string, platform int) (*models.Device, error) {
	return RdbSelectDevice(userId, platform)
}

func (p *gcpSqlProvider) SelectDevicesByUserId(userId string) ([]*models.Device, error) {
	return RdbSelectDevicesByUserId(userId)
}

func (p *gcpSqlProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return RdbSelectDevicesByToken(token)
}

func (p *gcpSqlProvider) UpdateDevice(device *models.Device) error {
	return RdbUpdateDevice(device)
}

func (p *gcpSqlProvider) DeleteDevice(userId string, platform int) error {
	return RdbDeleteDevice(userId, platform)
}
