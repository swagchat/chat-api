package datastore

import "github.com/swagchat/chat-api/models"

type DeviceStore interface {
	CreateDeviceStore()

	InsertDevice(device *models.Device) StoreResult
	SelectDevices(userId string) StoreResult
	SelectDevice(userId string, platform int) StoreResult
	SelectDevicesByUserId(userId string) StoreResult
	SelectDevicesByToken(token string) StoreResult
	UpdateDevice(device *models.Device) StoreResult
	DeleteDevice(userId string, platform int) StoreResult
}
