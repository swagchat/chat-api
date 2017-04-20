package datastore

import "github.com/fairway-corp/swagchat-api/models"

type DeviceStore interface {
	DeviceCreateStore()

	DeviceInsert(device *models.Device) StoreChannel
	DeviceSelect(userId string, platform int) StoreChannel
	DeviceSelectByUserId(userId string) StoreChannel
	DeviceUpdate(device *models.Device) StoreChannel
	DeviceSelectAll() StoreChannel
	DeviceDelete(userId string, platform int) StoreChannel
}
