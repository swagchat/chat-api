package datastore

import "github.com/swagchat/chat-api/model"

type deviceStore interface {
	createDeviceStore()

	InsertDevice(device *model.Device) (*model.Device, error)
	SelectDevices(userID string) ([]*model.Device, error)
	SelectDevice(userID string, platform int) (*model.Device, error)
	SelectDevicesByUserID(userID string) ([]*model.Device, error)
	SelectDevicesByToken(token string) ([]*model.Device, error)
	UpdateDevice(device *model.Device) error
	DeleteDevice(userID string, platform int) error
}
