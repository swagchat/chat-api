package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (provider GcpSqlProvider) InsertDevice(device *models.Device) StoreResult {
	return RdbInsertDevice(device)
}

func (provider GcpSqlProvider) SelectDevices() StoreResult {
	return RdbSelectDevices()
}

func (provider GcpSqlProvider) SelectDevice(userId string, platform int) StoreResult {
	return RdbSelectDevice(userId, platform)
}

func (provider GcpSqlProvider) SelectDevicesByUserId(userId string) StoreResult {
	return RdbSelectDevicesByUserId(userId)
}

func (provider GcpSqlProvider) UpdateDevice(device *models.Device) StoreResult {
	return RdbUpdateDevice(device)
}

func (provider GcpSqlProvider) DeleteDevice(userId string, platform int) StoreResult {
	return RdbDeleteDevice(userId, platform)
}
