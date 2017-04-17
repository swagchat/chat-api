package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) DeviceCreateStore() {
	RdbDeviceCreateStore()
}

func (provider GcpSqlProvider) DeviceInsert(device *models.Device) StoreChannel {
	return RdbDeviceInsert(device)
}

func (provider GcpSqlProvider) DeviceSelect(userId string, platform int) StoreChannel {
	return RdbDeviceSelect(userId, platform)
}

func (provider GcpSqlProvider) DeviceUpdate(device *models.Device) StoreChannel {
	return RdbDeviceUpdate(device)
}

func (provider GcpSqlProvider) DeviceSelectAll() StoreChannel {
	return RdbDeviceSelectAll()
}

func (provider GcpSqlProvider) DeviceDelete(userId string, platform int) StoreChannel {
	return RdbDeviceDelete(userId, platform)
}
