package datastore

import (
	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type insertDeviceOptions struct {
	beforeClean bool
}

type InsertDeviceOption func(*insertDeviceOptions)

func InsertDeviceOptionBeforeClean(beforeClean bool) InsertDeviceOption {
	return func(ops *insertDeviceOptions) {
		ops.beforeClean = beforeClean
	}
}

type selectDevicesOptions struct {
	deleted  bool
	userID   string
	platform scpb.Platform
	token    string
}

type SelectDevicesOption func(*selectDevicesOptions)

func SelectDevicesOptionFilterByDeleted(deleted bool) SelectDevicesOption {
	return func(ops *selectDevicesOptions) {
		ops.deleted = deleted
	}
}

func SelectDevicesOptionFilterByUserID(userID string) SelectDevicesOption {
	return func(ops *selectDevicesOptions) {
		ops.userID = userID
	}
}

func SelectDevicesOptionFilterByPlatform(platform scpb.Platform) SelectDevicesOption {
	return func(ops *selectDevicesOptions) {
		ops.platform = platform
	}
}

func SelectDevicesOptionFilterByToken(token string) SelectDevicesOption {
	return func(ops *selectDevicesOptions) {
		ops.token = token
	}
}

type deleteDevicesOptions struct {
	logicalDeleted int64
	userID         string
	platform       scpb.Platform
}

type DeleteDevicesOption func(*deleteDevicesOptions)

func DeleteDevicesOptionWithLogicalDeleted(logicalDeleted int64) DeleteDevicesOption {
	return func(ops *deleteDevicesOptions) {
		ops.logicalDeleted = logicalDeleted
	}
}

func DeleteDevicesOptionFilterByUserID(userID string) DeleteDevicesOption {
	return func(ops *deleteDevicesOptions) {
		ops.userID = userID
	}
}

func DeleteDevicesOptionFilterByPlatform(platform scpb.Platform) DeleteDevicesOption {
	return func(ops *deleteDevicesOptions) {
		ops.platform = platform
	}
}

type deviceStore interface {
	createDeviceStore()

	InsertDevice(device *model.Device, opts ...InsertDeviceOption) error
	SelectDevices(opts ...SelectDevicesOption) ([]*model.Device, error)
	SelectDevice(userID string, platform scpb.Platform) (*model.Device, error)
	UpdateDevice(device *model.Device) error
	DeleteDevices(opts ...DeleteDevicesOption) error
}
