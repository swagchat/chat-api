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
		newDevice := &model.Device{}
		newDevice.UserID = "datastore-user-id-0001"
		newDevice.Platform = scpb.Platform_PlatformIos
		newDevice.Token = "insert"
		err := Provider(ctx).InsertDevice(newDevice)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameInsertDevice)
		}
	})

	t.Run(TestNameSelectDevice, func(t *testing.T) {
		device, err = Provider(ctx).SelectDevice("datastore-user-id-0001", scpb.Platform_PlatformIos)
		if err != nil {
			t.Fatalf("Failed to %s", TestNameSelectDevice)
		}
	})

	t.Run(TestNameUpdateDevice, func(t *testing.T) {
		device.Token = "update"
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
	})
}
