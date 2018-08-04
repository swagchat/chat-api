package service

import (
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestNameCreateDevice = "create device test"
	TestNameGetDevices   = "get device test"
	TestNameDeleteDevice = "delete device test"
)

func TestDevice(t *testing.T) {
	t.Run(TestNameCreateDevice, func(t *testing.T) {
		req := &model.CreateDeviceRequest{}
		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformIos
		req.Token = "service-user-id-token-0001"
		device, errRes := CreateDevice(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}
		if device == nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}

		device, errRes = CreateDevice(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}
		if device != nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}

		req.Token = "service-user-id-token-0002"
		device, errRes = CreateDevice(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}
		if device == nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}

		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformNone
		_, errRes = CreateDevice(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameCreateDevice)
		}
	})

	t.Run(TestNameGetDevices, func(t *testing.T) {
		req := &model.GetDevicesRequest{}
		req.UserID = "service-user-id-0001"
		res, errRes := GetDevices(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameGetDevices)
		}
		if len(res.Devices) != 1 {
			t.Fatalf("Failed to %s", TestNameGetDevices)
		}
	})

	t.Run(TestNameDeleteDevice, func(t *testing.T) {
		req := &model.DeleteDeviceRequest{}
		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformIos
		errRes := DeleteDevice(ctx, req)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}

		errRes = DeleteDevice(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}

		cReq := &model.CreateDeviceRequest{}
		cReq.UserID = "service-user-id-0001"
		cReq.Platform = scpb.Platform_PlatformAndroid
		cReq.Token = "service-user-id-token-0001"
		device, errRes := CreateDevice(ctx, cReq)
		if errRes != nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
		if device == nil {
			t.Fatalf("Failed to %s", TestNameDeleteDevice)
		}
	})
}
