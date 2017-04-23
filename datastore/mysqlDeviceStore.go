package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (provider MysqlProvider) InsertDevice(device *models.Device) StoreChannel {
	return RdbInsertDevice(device)
}

func (provider MysqlProvider) SelectDevices() StoreChannel {
	return RdbSelectDevices()
}

func (provider MysqlProvider) SelectDevice(userId string, platform int) StoreChannel {
	return RdbSelectDevice(userId, platform)
}

func (provider MysqlProvider) SelectDevicesByUserId(userId string) StoreChannel {
	return RdbSelectDevicesByUserId(userId)
}

func (provider MysqlProvider) UpdateDevice(device *models.Device) StoreChannel {
	return RdbUpdateDevice(device)
}

func (provider MysqlProvider) DeleteDevice(userId string, platform int) StoreChannel {
	return RdbDeleteDevice(userId, platform)
}
