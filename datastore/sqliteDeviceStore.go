package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (provider SqliteProvider) InsertDevice(device *models.Device) StoreChannel {
	return RdbInsertDevice(device)
}

func (provider SqliteProvider) SelectDevices() StoreChannel {
	return RdbSelectDevices()
}

func (provider SqliteProvider) SelectDevice(userId string, platform int) StoreChannel {
	return RdbSelectDevice(userId, platform)
}

func (provider SqliteProvider) SelectDevicesByUserId(userId string) StoreChannel {
	return RdbSelectDevicesByUserId(userId)
}

func (provider SqliteProvider) UpdateDevice(device *models.Device) StoreChannel {
	return RdbUpdateDevice(device)
}

func (provider SqliteProvider) DeleteDevice(userId string, platform int) StoreChannel {
	return RdbDeleteDevice(userId, platform)
}
