package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) DeviceCreateStore() {
	RdbDeviceCreateStore()
}

func (provider MysqlProvider) DeviceInsert(device *models.Device) StoreChannel {
	return RdbDeviceInsert(device)
}

func (provider MysqlProvider) DeviceSelect(userId string, platform int) StoreChannel {
	return RdbDeviceSelect(userId, platform)
}

func (provider MysqlProvider) DeviceSelectByUserId(userId string) StoreChannel {
	return RdbDeviceSelectByUserId(userId)
}

func (provider MysqlProvider) DeviceUpdate(device *models.Device) StoreChannel {
	return RdbDeviceUpdate(device)
}

func (provider MysqlProvider) DeviceSelectAll() StoreChannel {
	return RdbDeviceSelectAll()
}

func (provider MysqlProvider) DeviceDelete(userId string, platform int) StoreChannel {
	return RdbDeviceDelete(userId, platform)
}
