package model

import (
	"testing"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameDevice              = "Device test"
	TestNameCreateDeviceRequest = "CreateDeviceRequest test"
	TestDevicesResponse         = "DevicesResponse test"
)

func TestDevice(t *testing.T) {
	t.Run(TestNameCreateDeviceRequest, func(t *testing.T) {
		d := &Device{}
		d.UserID = "model-user-id-0001"
		d.Platform = scpb.Platform_PlatformIos
		d.Token = "model-token-0001"
		d.NotificationDeviceID = "model-device-id-0001"

		pbd := d.ConvertToPbDevice()
		if pbd.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if pbd.Platform != scpb.Platform_PlatformIos {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if pbd.Token != "model-token-0001" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if pbd.NotificationDeviceID != "model-device-id-0001" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
	})

	t.Run(TestNameCreateDeviceRequest, func(t *testing.T) {
		cdr := &CreateDeviceRequest{}
		cdr.Platform = scpb.Platform_PlatformNone
		errRes := cdr.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if errRes.InvalidParams[0].Name != "platform" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}

		cdr.Platform = scpb.Platform_PlatformIos
		errRes = cdr.Validate()
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if len(errRes.InvalidParams) != 1 {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if errRes.InvalidParams[0].Name != "token" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}

		cdr.UserID = "model-user-id-0001"
		cdr.Platform = scpb.Platform_PlatformIos
		cdr.Token = "model-token-0001"
		errRes = cdr.Validate()
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}

		pbd := cdr.GenerateDevice()
		if pbd.UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if pbd.Platform != scpb.Platform_PlatformIos {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if pbd.Token != "model-token-0001" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
		if pbd.NotificationDeviceID != "" {
			t.Fatalf("Failed to %s", TestNameCreateDeviceRequest)
		}
	})

	t.Run(TestDevicesResponse, func(t *testing.T) {
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
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[0].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[0].Platform != scpb.Platform_PlatformIos {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[0].Token != "model-token-0001" {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[0].NotificationDeviceID != "model-device-id-0001" {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[1].UserID != "model-user-id-0001" {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[1].Platform != scpb.Platform_PlatformAndroid {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[1].Token != "model-token-0002" {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
		if pbRes.Devices[1].NotificationDeviceID != "model-device-id-0002" {
			t.Fatalf("Failed to %s", TestDevicesResponse)
		}
	})
}
