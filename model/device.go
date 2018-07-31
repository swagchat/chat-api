package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	PlatformIOS = iota + 1
	PlatformAndroid
	PlatformEnd
)

type Devices struct {
	Devices []*Device `json:"devices"`
}

type Device struct {
	scpb.Device
}

func (d *Device) UpdateDevice(req *UpdateDeviceRequest) {
	if req.Token != "" {
		d.Token = req.Token
	}
}

type CreateDeviceRequest struct {
	scpb.CreateDeviceRequest
}

func (crur *CreateDeviceRequest) Validate() *ErrorResponse {
	return nil
}

func (cdr *CreateDeviceRequest) GenerateDevice() *Device {
	device := &Device{}
	device.UserID = cdr.UserID
	device.Platform = cdr.Platform
	device.Token = cdr.Token
	return device
}

type GetDevicesRequest struct {
	scpb.GetDevicesRequest
}

type DevicesResponse struct {
	Devices []*Device
}

func (dr *DevicesResponse) ConvertToPbDevices() *scpb.DevicesResponse {
	devices := make([]*scpb.Device, len(dr.Devices))
	for i := 0; i < len(dr.Devices); i++ {
		d := dr.Devices[i]
		device := &scpb.Device{
			UserID:               d.UserID,
			Platform:             d.Platform,
			Token:                d.Token,
			NotificationDeviceID: d.NotificationDeviceID,
		}
		devices[i] = device
	}
	return &scpb.DevicesResponse{
		Devices: devices,
	}
}

type UpdateDeviceRequest struct {
	scpb.UpdateDeviceRequest
}

func (udr *UpdateDeviceRequest) Validate() *ErrorResponse {
	if udr.Platform != scpb.Platform_PlatformIos && udr.Platform != scpb.Platform_PlatformAndroid {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "platform",
				Reason: "platform is invalid. Currently only 1(iOS) and 2(Android) are supported.",
			},
		}
		return NewErrorResponse("Failed to update device.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if udr.Token == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "token",
				Reason: "token is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to update device.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	return nil
}

type DeleteDeviceRequest struct {
	scpb.DeleteDeviceRequest
	Room *Room
}

func (ddr *DeleteDeviceRequest) Validate() *ErrorResponse {
	return nil
}
