package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (provider MysqlProvider) InsertDevice(device *models.Device) StoreResult {
	return RdbInsertDevice(device)
}

func (provider MysqlProvider) SelectDevices(userId string) StoreResult {
	return RdbSelectDevices(userId)
}

func (provider MysqlProvider) SelectDevice(userId string, platform int) StoreResult {
	return RdbSelectDevice(userId, platform)
}

func (provider MysqlProvider) SelectDevicesByUserId(userId string) StoreResult {
	return RdbSelectDevicesByUserId(userId)
}

func (provider MysqlProvider) SelectDevicesByToken(token string) StoreResult {
	return RdbSelectDevicesByToken(token)
}

func (provider MysqlProvider) UpdateDevice(device *models.Device) StoreResult {
	return RdbUpdateDevice(device)
}

func (provider MysqlProvider) DeleteDevice(userId string, platform int) StoreResult {
	return RdbDeleteDevice(userId, platform)
}
