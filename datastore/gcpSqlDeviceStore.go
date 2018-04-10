package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore(p.database)
}

func (p *gcpSqlProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return RdbInsertDevice(p.database, device)
}

func (p *gcpSqlProvider) SelectDevices(userId string) ([]*models.Device, error) {
	return RdbSelectDevices(p.database, userId)
}

func (p *gcpSqlProvider) SelectDevice(userId string, platform int) (*models.Device, error) {
	return RdbSelectDevice(p.database, userId, platform)
}

func (p *gcpSqlProvider) SelectDevicesByUserId(userId string) ([]*models.Device, error) {
	return RdbSelectDevicesByUserId(p.database, userId)
}

func (p *gcpSqlProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return RdbSelectDevicesByToken(p.database, token)
}

func (p *gcpSqlProvider) UpdateDevice(device *models.Device) error {
	return RdbUpdateDevice(p.database, device)
}

func (p *gcpSqlProvider) DeleteDevice(userId string, platform int) error {
	return RdbDeleteDevice(p.database, userId, platform)
}
