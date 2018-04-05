package datastore

import "github.com/swagchat/chat-api/models"

type DeviceStore interface {
	CreateDeviceStore()

	InsertDevice(device *models.Device) (*models.Device, error)
	SelectDevices(userId string) ([]*models.Device, error)
	SelectDevice(userId string, platform int) (*models.Device, error)
	SelectDevicesByUserId(userId string) ([]*models.Device, error)
	SelectDevicesByToken(token string) ([]*models.Device, error)
	UpdateDevice(device *models.Device) error
	DeleteDevice(userId string, platform int) error
}
