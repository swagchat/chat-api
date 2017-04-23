package datastore

import "github.com/fairway-corp/swagchat-api/models"

type DeviceStore interface {
	CreateDeviceStore()

	InsertDevice(device *models.Device) StoreResult
	SelectDevices() StoreResult
	SelectDevice(userId string, platform int) StoreResult
	SelectDevicesByUserId(userId string) StoreResult
	UpdateDevice(device *models.Device) StoreResult
	DeleteDevice(userId string, platform int) StoreResult
}
