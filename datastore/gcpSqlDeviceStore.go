package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *gcpSqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (p *gcpSqlProvider) InsertDevice(device *models.Device) StoreResult {
	return RdbInsertDevice(device)
}

func (p *gcpSqlProvider) SelectDevices(userId string) StoreResult {
	return RdbSelectDevices(userId)
}

func (p *gcpSqlProvider) SelectDevice(userId string, platform int) StoreResult {
	return RdbSelectDevice(userId, platform)
}

func (p *gcpSqlProvider) SelectDevicesByUserId(userId string) StoreResult {
	return RdbSelectDevicesByUserId(userId)
}

func (p *gcpSqlProvider) SelectDevicesByToken(token string) StoreResult {
	return RdbSelectDevicesByToken(token)
}

func (p *gcpSqlProvider) UpdateDevice(device *models.Device) StoreResult {
	return RdbUpdateDevice(device)
}

func (p *gcpSqlProvider) DeleteDevice(userId string, platform int) StoreResult {
	return RdbDeleteDevice(userId, platform)
}
