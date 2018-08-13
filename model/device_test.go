package model

import (
	"testing"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestModelDevice           = "[model] Device test"
	TestModelAddDeviceRequest = "[model] AddDeviceRequest test"
	TestModelDevicesResponse  = "[model] DevicesResponse test"
)

func TestDevice(t *testing.T) {
	t.Run(TestModelDevice, func(t *testing.T) {
		d := &Device{}
		d.UserID = "model-user-id-0001"
		d.Platform = scpb.Platform_PlatformIos
		d.Token = "model-token-0001"
		d.NotificationDeviceID = "model-device-id-0001"

		pbd := d.ConvertToPbDevice()
		if pbd.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected pbd.UserID to be \"model-user-id-0001\", but it was %s", TestModelDevice, pbd.UserID)
		}
		if pbd.Platform != scpb.Platform_PlatformIos {
			t.Fatalf("Failed to %s. Expected pbd.Platform to be 1, but it was %d", TestModelDevice, pbd.Platform)
		}
		if pbd.Token != "model-token-0001" {
			t.Fatalf("Failed to %s. Expected pbd.Token to be \"model-token-0001\", but it was %s", TestModelDevice, pbd.Token)
		}
		if pbd.NotificationDeviceID != "model-device-id-0001" {
			t.Fatalf("Failed to %s. Expected pbd.NotificationDeviceID to be \"model-device-id-0001\", but it was %s", TestModelDevice, pbd.NotificationDeviceID)
		}
	})

	t.Run(TestModelAddDeviceRequest, func(t *testing.T) {
		cdr := &AddDeviceRequest{}
		cdr.Platform = scpb.Platform_PlatformNone
		errRes := cdr.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelAddDeviceRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelAddDeviceRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "platform" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"platform\", but it was %s", TestModelAddDeviceRequest, errRes.InvalidParams[0].Name)
		}

		cdr.Platform = scpb.Platform_PlatformIos
		errRes = cdr.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestModelAddDeviceRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams count to be 1, but it was %d", TestModelAddDeviceRequest, len(errRes.InvalidParams))
		}
		if errRes.InvalidParams[0].Name != "token" {
			t.Fatalf("Failed to %s. Expected errRes.InvalidParams[0].Name is \"token\", but it was %s", TestModelAddDeviceRequest, errRes.InvalidParams[0].Name)
		}

		cdr.UserID = "model-user-id-0001"
		cdr.Platform = scpb.Platform_PlatformIos
		cdr.Token = "model-token-0001"
		errRes = cdr.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil. %s is invalid", TestModelAddDeviceRequest, errRes.InvalidParams[0].Name)
		}

		pbd := cdr.GenerateDevice()
		if pbd.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected pbd.UserID to be \"model-user-id-0001\", but it was %s", TestModelAddDeviceRequest, pbd.UserID)
		}
		if pbd.Platform != scpb.Platform_PlatformIos {
			t.Fatalf("Failed to %s. Expected pbd.Platform to be 1, but it was %d", TestModelAddDeviceRequest, pbd.Platform)
		}
		if pbd.Token != "model-token-0001" {
			t.Fatalf("Failed to %s. Expected pbd.Token to be \"model-user-id-0001\", but it was %s", TestModelAddDeviceRequest, pbd.Token)
		}
		if pbd.NotificationDeviceID != "" {
			t.Fatalf("Failed to %s. Expected pbd.NotificationDeviceID to be \"\", but it was %s", TestModelAddDeviceRequest, pbd.NotificationDeviceID)
		}
	})

	t.Run(TestModelDevicesResponse, func(t *testing.T) {
		d1 := &Device{}
		d1.UserID = "model-user-id-0001"
		d1.Platform = scpb.Platform_PlatformIos
		d1.Token = "model-token-0001"
		d1.NotificationDeviceID = "model-device-id-0001"

		d2 := &Device{}
		d2.UserID = "model-user-id-0001"
		d2.Platform = scpb.Platform_PlatformAndroid
		d2.Token = "model-token-0002"
		d2.NotificationDeviceID = "model-device-id-0002"

		res := &DevicesResponse{}
		res.Devices = []*Device{d1, d2}
		pbRes := res.ConvertToPbDevices()
		if len(pbRes.Devices) != 2 {
			t.Fatalf("Failed to %s. Expected pbRes.Devices count to be 2, but it was %d", TestModelDevicesResponse, len(pbRes.Devices))
		}
		if pbRes.Devices[0].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[0].UserID to be \"model-user-id-0001\", but it was %s", TestModelDevicesResponse, pbRes.Devices[0].UserID)
		}
		if pbRes.Devices[0].Platform != scpb.Platform_PlatformIos {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[0].Platform to be 1, but it was %d", TestModelDevicesResponse, pbRes.Devices[0].Platform)
		}
		if pbRes.Devices[0].Token != "model-token-0001" {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[0].Token to be \"model-token-0001\", but it was %s", TestModelDevicesResponse, pbRes.Devices[0].Token)
		}
		if pbRes.Devices[0].NotificationDeviceID != "model-device-id-0001" {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[0].NotificationDeviceID to be \"model-device-id-0001\", but it was %s", TestModelDevicesResponse, pbRes.Devices[0].NotificationDeviceID)
		}
		if pbRes.Devices[1].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[1].UserID to be \"model-user-id-0001\", but it was %s", TestModelDevicesResponse, pbRes.Devices[1].UserID)
		}
		if pbRes.Devices[1].Platform != scpb.Platform_PlatformAndroid {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[1].Platform to be 2, but it was %d", TestModelDevicesResponse, pbRes.Devices[1].Platform)
		}
		if pbRes.Devices[1].Token != "model-token-0002" {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[1].Token to be \"model-token-0002\", but it was %s", TestModelDevicesResponse, pbRes.Devices[1].Token)
		}
		if pbRes.Devices[1].NotificationDeviceID != "model-device-id-0002" {
			t.Fatalf("Failed to %s. Expected pbRes.Devices[1].NotificationDeviceID to be \"model-device-id-0002\", but it was %s", TestModelDevicesResponse, pbRes.Devices[1].NotificationDeviceID)
		}
	})
}
