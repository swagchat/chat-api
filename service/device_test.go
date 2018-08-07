package service

import (
	"fmt"
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
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateDevice, errMsg)
		}
		if device == nil {
			t.Fatalf("Failed to %s. Expected device to be not nil, but it was nil", TestNameCreateDevice)
		}

		device, errRes = CreateDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateDevice, errMsg)
		}
		if device != nil {
			t.Fatalf("Failed to %s. Expected device to be nil, but it was not nil", TestNameCreateDevice)
		}

		req.Token = "service-user-id-token-0002"
		device, errRes = CreateDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameCreateDevice, errMsg)
		}
		if device == nil {
			t.Fatalf("Failed to %s. Expected device to be not nil, but it was nil", TestNameCreateDevice)
		}

		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformNone
		_, errRes = CreateDevice(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameCreateDevice)
		}
	})

	t.Run(TestNameGetDevices, func(t *testing.T) {
		req := &model.GetDevicesRequest{}
		req.UserID = "service-user-id-0001"
		res, errRes := GetDevices(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameGetDevices, errMsg)
		}
		if len(res.Devices) != 1 {
			t.Fatalf("Failed to %s. Expected res.Devices count to be 1, but it was %d", TestNameGetDevices, len(res.Devices))
		}
	})

	t.Run(TestNameDeleteDevice, func(t *testing.T) {
		req := &model.DeleteDeviceRequest{}
		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformIos
		errRes := DeleteDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteDevice, errMsg)
		}

		errRes = DeleteDevice(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteDevice)
		}

		cReq := &model.CreateDeviceRequest{}
		cReq.UserID = "service-user-id-0001"
		cReq.Platform = scpb.Platform_PlatformAndroid
		cReq.Token = "service-user-id-token-0001"
		device, errRes := CreateDevice(ctx, cReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestNameDeleteDevice, errMsg)
		}
		if device == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestNameDeleteDevice)
		}
	})
}
