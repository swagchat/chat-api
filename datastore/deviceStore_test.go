package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestStoreInsertDevice  = "[store] insert device test"
	TestStoreSelectDevices = "[store] select devices test"
	TestStoreSelectDevice  = "[store] select device test"
	TestStoreUpdateDevice  = "[store] select update test"
	TestStoreDeleteDevice  = "[store] delete block user test"
)

func TestDeviceStore(t *testing.T) {
	var device *model.Device
	var err error

	t.Run(TestStoreInsertDevice, func(t *testing.T) {
		newDevice1 := &model.Device{}
		newDevice1.UserID = "datastore-user-id-0001"
		newDevice1.Platform = scpb.Platform_PlatformIos
		newDevice1.Token = "user-id-token-0001"
		newDevice1.NotificationDeviceID = "user-id-device-id-0001"
		err := Provider(ctx).InsertDevice(newDevice1)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertDevice, err.Error())
		}

		newDevice2 := &model.Device{}
		newDevice2.UserID = "datastore-user-id-0001"
		newDevice2.Platform = scpb.Platform_PlatformAndroid
		newDevice2.Token = "user-id-token-0002"
		newDevice2.NotificationDeviceID = "user-id-device-id-0002"
		err = Provider(ctx).InsertDevice(
			newDevice2,
			InsertDeviceOptionBeforeClean(true),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertDevice, err.Error())
		}

		newDevice3 := &model.Device{}
		newDevice3.UserID = "datastore-user-id-0002"
		newDevice3.Platform = scpb.Platform_PlatformIos
		newDevice3.Token = "user-id-token-0003"
		newDevice3.NotificationDeviceID = "user-id-device-id-0003"
		err = Provider(ctx).InsertDevice(newDevice3)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertDevice, err.Error())
		}

		newDevice4 := &model.Device{}
		newDevice4.UserID = "datastore-user-id-0002"
		newDevice4.Platform = scpb.Platform_PlatformAndroid
		newDevice4.Token = "user-id-token-0004"
		newDevice4.NotificationDeviceID = "user-id-device-id-0004"
		err = Provider(ctx).InsertDevice(newDevice4)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreInsertDevice, err.Error())
		}
	})

	t.Run(TestStoreSelectDevices, func(t *testing.T) {
		_, err = Provider(ctx).SelectDevices()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectDevices)
		}
		errMsg := "An error occurred while getting devices. Be sure to specify either userId or platform or token"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectDevices, errMsg, err.Error())
		}

		devices, err := Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevices, err.Error())
		}
		if len(devices) != 2 {
			t.Fatalf("Failed to %s. Expected devices count to be 2, but it was %d", TestStoreSelectDevices, len(devices))
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByPlatform(scpb.Platform_PlatformIos),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevices, err.Error())
		}
		if len(devices) != 2 {
			t.Fatalf("Failed to %s. Expected devices count to be 2, but it was %d", TestStoreSelectDevices, len(devices))
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByToken("user-id-token-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevices, err.Error())
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s. Expected devices count to be 1, but it was %d", TestStoreSelectDevices, len(devices))
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
			SelectDevicesOptionFilterByPlatform(scpb.Platform_PlatformAndroid),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevices, err.Error())
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s. Expected devices count to be 1, but it was %d", TestStoreSelectDevices, len(devices))
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
			SelectDevicesOptionFilterByToken("user-id-token-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevices, err.Error())
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s. Expected devices count to be 1, but it was %d", TestStoreSelectDevices, len(devices))
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByPlatform(scpb.Platform_PlatformIos),
			SelectDevicesOptionFilterByToken("user-id-token-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevices, err.Error())
		}
		if len(devices) != 0 {
			t.Fatalf("Failed to %s. Expected devices count to be 0, but it was %d", TestStoreSelectDevices, len(devices))
		}
	})

	t.Run(TestStoreSelectDevice, func(t *testing.T) {
		device, err = Provider(ctx).SelectDevice("not-exist-user", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevice, err.Error())
		}
		if device != nil {
			t.Fatalf("Failed to %s. Expected device is nil, but it was not nil", TestStoreSelectDevice)
		}

		device, err = Provider(ctx).SelectDevice("datastore-user-id-0001", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreSelectDevice, err.Error())
		}
	})

	t.Run(TestStoreUpdateDevice, func(t *testing.T) {
		device.Token = "update"
		device.NotificationDeviceID = "update"
		err = Provider(ctx).UpdateDevice(device)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}

		device, err = Provider(ctx).SelectDevice("datastore-user-id-0001", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}
		if device.Token != "update" {
			t.Fatalf("Failed to %s. Expected device.Token to be \"update\", but it was %s", TestStoreUpdateDevice, device.Token)
		}
		if device.NotificationDeviceID != "update" {
			t.Fatalf("Failed to %s. Expected device.NotificationDeviceID to be \"update\", but it was %s", TestStoreUpdateDevice, device.Token)
		}
	})

	t.Run(TestStoreDeleteDevice, func(t *testing.T) {
		err = Provider(ctx).DeleteDevices()
		if err == nil {
			t.Fatalf("Failed to %s. Expected err to be not nil, but it was nil", TestStoreSelectDevices)
		}
		errMsg := "An error occurred while deleting devices. Be sure to specify either userID or platform"
		if err.Error() != errMsg {
			t.Fatalf("Failed to %s. Expected err message to be \"%s\", but it was %s", TestStoreSelectDevices, errMsg, err.Error())
		}

		err = Provider(ctx).DeleteDevices(
			DeleteDevicesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}
		devices, err := Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}
		if len(devices) != 0 {
			t.Fatalf("Failed to %s. Expected devices count to be 0, but it was %d", TestStoreUpdateDevice, len(devices))
		}

		err = Provider(ctx).DeleteDevices(
			DeleteDevicesOptionWithLogicalDeleted(1),
			DeleteDevicesOptionFilterByUserID("datastore-user-id-0002"),
			DeleteDevicesOptionFilterByPlatform(scpb.Platform_PlatformIos),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}
		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s. Expected devices count to be 1, but it was %d", TestStoreUpdateDevice, len(devices))
		}
		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByDeleted(true),
			SelectDevicesOptionFilterByUserID("datastore-user-id-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s. Expected err to be nil, but it was not nil [%s]", TestStoreUpdateDevice, err.Error())
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s. Expected devices count to be 1, but it was %d", TestStoreUpdateDevice, len(devices))
		}
	})
}
