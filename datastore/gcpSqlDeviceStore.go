package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createDeviceStore() {
	rdbCreateDeviceStore(p.database)
}

func (p *gcpSQLProvider) InsertDevice(device *models.Device) (*models.Device, error) {
	return rdbInsertDevice(p.database, device)
}

func (p *gcpSQLProvider) SelectDevices(userID string) ([]*models.Device, error) {
	return rdbSelectDevices(p.database, userID)
}

func (p *gcpSQLProvider) SelectDevice(userID string, platform int) (*models.Device, error) {
	return rdbSelectDevice(p.database, userID, platform)
}

func (p *gcpSQLProvider) SelectDevicesByUserID(userID string) ([]*models.Device, error) {
	return rdbSelectDevicesByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectDevicesByToken(token string) ([]*models.Device, error) {
	return rdbSelectDevicesByToken(p.database, token)
}

func (p *gcpSQLProvider) UpdateDevice(device *models.Device) error {
	return rdbUpdateDevice(p.database, device)
}

func (p *gcpSQLProvider) DeleteDevice(userID string, platform int) error {
	return rdbDeleteDevice(p.database, userID, platform)
}
