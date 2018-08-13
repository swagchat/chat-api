package service

import (
	"fmt"
	"testing"

	"github.com/swagchat/chat-api/model"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	TestServiceAddDevice       = "[service] add device test"
	TestServiceRetrieveDevices = "[service] retrieve device test"
	TestServiceDeleteDevice    = "[service] delete device test"
)

func TestDevice(t *testing.T) {
	t.Run(TestServiceAddDevice, func(t *testing.T) {
		req := &model.AddDeviceRequest{}
		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformIos
		req.Token = "service-user-id-token-0001"
		device, errRes := AddDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceAddDevice, errMsg)
		}
		if device == nil {
			t.Fatalf("Failed to %s. Expected device to be not nil, but it was nil", TestServiceAddDevice)
		}

		device, errRes = AddDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceAddDevice, errMsg)
		}
		if device != nil {
			t.Fatalf("Failed to %s. Expected device to be nil, but it was not nil", TestServiceAddDevice)
		}

		req.Token = "service-user-id-token-0002"
		device, errRes = AddDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceAddDevice, errMsg)
		}
		if device == nil {
			t.Fatalf("Failed to %s. Expected device to be not nil, but it was nil", TestServiceAddDevice)
		}

		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformNone
		_, errRes = AddDevice(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceAddDevice)
		}
	})

	t.Run(TestServiceRetrieveDevices, func(t *testing.T) {
		req := &model.RetrieveDevicesRequest{}
		req.UserID = "service-user-id-0001"
		res, errRes := RetrieveDevices(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceRetrieveDevices, errMsg)
		}
		if len(res.Devices) != 1 {
			t.Fatalf("Failed to %s. Expected res.Devices count to be 1, but it was %d", TestServiceRetrieveDevices, len(res.Devices))
		}
	})

	t.Run(TestServiceDeleteDevice, func(t *testing.T) {
		req := &model.DeleteDeviceRequest{}
		req.UserID = "service-user-id-0001"
		req.Platform = scpb.Platform_PlatformIos
		errRes := DeleteDevice(ctx, req)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteDevice, errMsg)
		}

		errRes = DeleteDevice(ctx, req)
		if errRes == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteDevice)
		}

		cReq := &model.AddDeviceRequest{}
		cReq.UserID = "service-user-id-0001"
		cReq.Platform = scpb.Platform_PlatformAndroid
		cReq.Token = "service-user-id-token-0001"
		device, errRes := AddDevice(ctx, cReq)
		if errRes != nil {
			errMsg := ""
			if errRes.Error != nil {
				errMsg = fmt.Sprintf(" [%s]", errRes.Error.Error())
			}
			t.Fatalf("Failed to %s. Expected errRes to be nil, but it was not nil%s", TestServiceDeleteDevice, errMsg)
		}
		if device == nil {
			t.Fatalf("Failed to %s. Expected errRes to be not nil, but it was nil", TestServiceDeleteDevice)
		}
	})
}
