package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateDeviceStore() {
	RdbCreateDeviceStore()
}

func (provider GcpSqlProvider) InsertDevice(device *models.Device) StoreChannel {
	return RdbInsertDevice(device)
}

func (provider GcpSqlProvider) SelectDevices() StoreChannel {
	return RdbSelectDevices()
}

func (provider GcpSqlProvider) SelectDevice(userId string, platform int) StoreChannel {
	return RdbSelectDevice(userId, platform)
}

func (provider GcpSqlProvider) SelectDevicesByUserId(userId string) StoreChannel {
	return RdbSelectDevicesByUserId(userId)
}

func (provider GcpSqlProvider) UpdateDevice(device *models.Device) StoreChannel {
	return RdbUpdateDevice(device)
}

func (provider GcpSqlProvider) DeleteDevice(userId string, platform int) StoreChannel {
	return RdbDeleteDevice(userId, platform)
}
