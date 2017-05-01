package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (provider SqliteProvider) InsertDevice(device *models.Device) StoreResult {
	return RdbInsertDevice(device)
}

func (provider SqliteProvider) SelectDevices(userId string) StoreResult {
	return RdbSelectDevices(userId)
}

func (provider SqliteProvider) SelectDevice(userId string, platform int) StoreResult {
	return RdbSelectDevice(userId, platform)
}

func (provider SqliteProvider) SelectDevicesByUserId(userId string) StoreResult {
	return RdbSelectDevicesByUserId(userId)
}

func (provider SqliteProvider) UpdateDevice(device *models.Device) StoreResult {
	return RdbUpdateDevice(device)
}

func (provider SqliteProvider) DeleteDevice(userId string, platform int) StoreResult {
	return RdbDeleteDevice(userId, platform)
}
