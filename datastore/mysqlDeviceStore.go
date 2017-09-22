package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (p *mysqlProvider) InsertDevice(device *models.Device) StoreResult {
	return RdbInsertDevice(device)
}

func (p *mysqlProvider) SelectDevices(userId string) StoreResult {
	return RdbSelectDevices(userId)
}

func (p *mysqlProvider) SelectDevice(userId string, platform int) StoreResult {
	return RdbSelectDevice(userId, platform)
}

func (p *mysqlProvider) SelectDevicesByUserId(userId string) StoreResult {
	return RdbSelectDevicesByUserId(userId)
}

func (p *mysqlProvider) SelectDevicesByToken(token string) StoreResult {
	return RdbSelectDevicesByToken(token)
}

func (p *mysqlProvider) UpdateDevice(device *models.Device) StoreResult {
	return RdbUpdateDevice(device)
}

func (p *mysqlProvider) DeleteDevice(userId string, platform int) StoreResult {
	return RdbDeleteDevice(userId, platform)
}
