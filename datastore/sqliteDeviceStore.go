package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) DeviceCreateStore() {
	RdbDeviceCreateStore()
}

func (provider SqliteProvider) DeviceInsert(device *models.Device) StoreChannel {
	return RdbDeviceInsert(device)
}

func (provider SqliteProvider) DeviceSelect(userId string, platform int) StoreChannel {
	return RdbDeviceSelect(userId, platform)
}

func (provider SqliteProvider) DeviceUpdate(device *models.Device) StoreChannel {
	return RdbDeviceUpdate(device)
}

func (provider SqliteProvider) DeviceSelectAll() StoreChannel {
	return RdbDeviceSelectAll()
}

func (provider SqliteProvider) DeviceDelete(userId string, platform int) StoreChannel {
	return RdbDeviceDelete(userId, platform)
}
