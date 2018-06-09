package datastore

import "github.com/swagchat/chat-api/models"

type deviceStore interface {
	createDeviceStore()

	InsertDevice(device *models.Device) (*models.Device, error)
	SelectDevices(userID string) ([]*models.Device, error)
	SelectDevice(userID string, platform int) (*models.Device, error)
	SelectDevicesByUserID(userID string) ([]*models.Device, error)
	SelectDevicesByToken(token string) ([]*models.Device, error)
	UpdateDevice(device *models.Device) error
	DeleteDevice(userID string, platform int) error
}
