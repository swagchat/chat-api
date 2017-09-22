package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (p *sqliteProvider) InsertDevice(device *models.Device) StoreResult {
	return RdbInsertDevice(device)
}

func (p *sqliteProvider) SelectDevices(userId string) StoreResult {
	return RdbSelectDevices(userId)
}

func (p *sqliteProvider) SelectDevice(userId string, platform int) StoreResult {
	return RdbSelectDevice(userId, platform)
}

func (p *sqliteProvider) SelectDevicesByUserId(userId string) StoreResult {
	return RdbSelectDevicesByUserId(userId)
}

func (p *sqliteProvider) SelectDevicesByToken(token string) StoreResult {
	return RdbSelectDevicesByToken(token)
}

func (p *sqliteProvider) UpdateDevice(device *models.Device) StoreResult {
	return RdbUpdateDevice(device)
}

func (p *sqliteProvider) DeleteDevice(userId string, platform int) StoreResult {
	return RdbDeleteDevice(userId, platform)
}
