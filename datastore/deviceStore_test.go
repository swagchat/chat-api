package datastore

import (
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameInsertDevice  = "insert device test"
	TestNameSelectDevices = "select devices test"
	TestNameSelectDevice  = "select device test"
	TestNameUpdateDevice  = "select update test"
	TestNameDeleteDevice  = "delete block user test"
)

func TestDeviceStore(t *testing.T) {
	var device *model.Device
	var err error

	t.Run(TestNameInsertDevice, func(t *testing.T) {
		newDevice1 := &model.Device{}
		newDevice1.UserID = "datastore-user-id-0001"
		newDevice1.Platform = scpb.Platform_PlatformIos
		newDevice1.Token = "user-id-token-0001"
		newDevice1.NotificationDeviceID = "user-id-device-id-0001"
		err := Provider(ctx).InsertDevice(newDevice1)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
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
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		newDevice3 := &model.Device{}
		newDevice3.UserID = "datastore-user-id-0002"
		newDevice3.Platform = scpb.Platform_PlatformIos
		newDevice3.Token = "user-id-token-0003"
		newDevice3.NotificationDeviceID = "user-id-device-id-0003"
		err = Provider(ctx).InsertDevice(newDevice3)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		newDevice4 := &model.Device{}
		newDevice4.UserID = "datastore-user-id-0002"
		newDevice4.Platform = scpb.Platform_PlatformAndroid
		newDevice4.Token = "user-id-token-0004"
		newDevice4.NotificationDeviceID = "user-id-device-id-0004"
		err = Provider(ctx).InsertDevice(newDevice4)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
	})

	t.Run(TestNameSelectDevices, func(t *testing.T) {
		_, err = Provider(ctx).SelectDevices()
		if err == nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if err.Error() != "An error occurred while getting devices. Be sure to specify either userId or platform or token" {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		devices, err := Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if len(devices) != 2 {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByPlatform(scpb.Platform_PlatformIos),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if len(devices) != 2 {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByToken("user-id-token-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
			SelectDevicesOptionFilterByPlatform(scpb.Platform_PlatformAndroid),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
			SelectDevicesOptionFilterByToken("user-id-token-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}

		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByPlatform(scpb.Platform_PlatformIos),
			SelectDevicesOptionFilterByToken("user-id-token-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
		if len(devices) != 0 {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
	})

	t.Run(TestNameSelectDevice, func(t *testing.T) {
		device, err = Provider(ctx).SelectDevice("not-exist-user", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectDevice)
		}
		if device != nil {
			t.Fatalf("Failed to %s", TestNameSelectDevice)
		}

		device, err = Provider(ctx).SelectDevice("datastore-user-id-0001", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectDevice)
		}
	})

	t.Run(TestNameUpdateDevice, func(t *testing.T) {
		device.Token = "update"
		device.NotificationDeviceID = "update"
		err = Provider(ctx).UpdateDevice(device)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateDevice)
		}

		device, err = Provider(ctx).SelectDevice("datastore-user-id-0001", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameUpdateDevice)
		}
		if device.Token != "update" {
			t.Fatalf("Failed to %s", TestNameUpdateDevice)
		}
		if device.NotificationDeviceID != "update" {
			t.Fatalf("Failed to %s", TestNameUpdateDevice)
		}
	})

	t.Run(TestNameDeleteDevice, func(t *testing.T) {
		err = Provider(ctx).DeleteDevices()
		if err == nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		if err.Error() != "An error occurred while deleting devices. Be sure to specify either userID or platform" {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}

		err = Provider(ctx).DeleteDevices(
			DeleteDevicesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		devices, err := Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0001"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		if len(devices) != 0 {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}

		err = Provider(ctx).DeleteDevices(
			DeleteDevicesOptionWithLogicalDeleted(1),
			DeleteDevicesOptionFilterByUserID("datastore-user-id-0002"),
			DeleteDevicesOptionFilterByPlatform(scpb.Platform_PlatformIos),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByUserID("datastore-user-id-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		devices, err = Provider(ctx).SelectDevices(
			SelectDevicesOptionFilterByDeleted(true),
			SelectDevicesOptionFilterByUserID("datastore-user-id-0002"),
		)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		if len(devices) != 1 {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
	})
}
