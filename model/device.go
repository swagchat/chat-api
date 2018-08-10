package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type Device struct {
	scpb.Device
}

func (d *Device) ConvertToPbDevice() *scpb.Device {
	return &scpb.Device{
		UserID:               d.UserID,
		Platform:             d.Platform,
		Token:                d.Token,
		NotificationDeviceID: d.NotificationDeviceID,
	}
}

type CreateDeviceRequest struct {
	scpb.CreateDeviceRequest
}

func (cdr *CreateDeviceRequest) Validate() *ErrorResponse {
	if !(cdr.Platform == scpb.Platform_PlatformIos || cdr.Platform == scpb.Platform_PlatformAndroid) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "platform",
				Reason: "platform is invalid. Currently only 1(iOS) and 2(Android) are supported.",
			},
		}
		return NewErrorResponse("Failed to update device.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if cdr.Token == "" {
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

func (cdr *CreateDeviceRequest) GenerateDevice() *Device {
	device := &Device{}
	device.UserID = cdr.UserID
	device.Platform = cdr.Platform
	device.Token = cdr.Token
	return device
}

type RetrieveDevicesRequest struct {
	scpb.RetrieveDevicesRequest
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

type DeleteDeviceRequest struct {
	scpb.DeleteDeviceRequest
	Room *Room
}
