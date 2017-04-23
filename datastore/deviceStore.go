package datastore

import "github.com/fairway-corp/swagchat-api/models"

type DeviceStore interface {
	CreateDeviceStore()

	InsertDevice(device *models.Device) StoreChannel
	SelectDevices() StoreChannel
	SelectDevice(userId string, platform int) StoreChannel
	SelectDevicesByUserId(userId string) StoreChannel
	UpdateDevice(device *models.Device) StoreChannel
	DeleteDevice(userId string, platform int) StoreChannel
}
