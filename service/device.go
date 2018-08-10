package service

import (
	"context"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/tracer"
)

// CreateDevice create a device
func CreateDevice(ctx context.Context, req *model.CreateDeviceRequest) (*model.Device, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("UpdateDevice", "service")
	defer tracer.Provider(ctx).Finish(span)

	device, err := datastore.Provider(ctx).SelectDevice(req.UserID, req.Platform)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create device.", http.StatusInternalServerError, model.WithError(err))
	}

	if device != nil && device.Token == req.Token {
		return nil, nil
	}

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	newDevice := req.GenerateDevice()

	nRes := <-notification.Provider(ctx).CreateEndpoint(req.UserID, req.Platform, req.Token)
	if nRes.Error != nil {
		return nil, model.NewErrorResponse("Failed to create device.", http.StatusInternalServerError, model.WithError(nRes.Error))
	}

	if nRes.Data != nil {
		newDevice.NotificationDeviceID = *nRes.Data.(*string)
	}

	err = datastore.Provider(ctx).InsertDevice(
		newDevice,
		datastore.InsertDeviceOptionBeforeClean(true),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create device.", http.StatusInternalServerError, model.WithError(err))
	}
	go subscribeByDevice(ctx, newDevice, nil)

	if device != nil {
		nRes = <-notification.Provider(ctx).DeleteEndpoint(device.NotificationDeviceID)
		if nRes.Error != nil {
			return nil, model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(nRes.Error))
		}
		go unsubscribeByDevice(ctx, device, nil)
	}

	return newDevice, nil
}

// RetrieveDevices retrieves devices
func RetrieveDevices(ctx context.Context, req *model.RetrieveDevicesRequest) (*model.DevicesResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("RetrieveDevices", "service")
	defer tracer.Provider(ctx).Finish(span)

	devices, err := datastore.Provider(ctx).SelectDevices(datastore.SelectDevicesOptionFilterByUserID(req.UserID))
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get devices.", http.StatusInternalServerError, model.WithError(err))
	}

	return &model.DevicesResponse{
		Devices: devices,
	}, nil
}

// DeleteDevice deletes device
func DeleteDevice(ctx context.Context, req *model.DeleteDeviceRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("DeleteDevice", "service")
	defer tracer.Provider(ctx).Finish(span)

	device, errRes := confirmDeviceExist(ctx, req.UserID, req.Platform)
	if errRes != nil {
		errRes.Message = "Failed to delete devices."
		return errRes
	}

	np := notification.Provider(ctx)
	nRes := <-np.DeleteEndpoint(device.NotificationDeviceID)
	if nRes.Error != nil {
		return model.NewErrorResponse("Failed to delete devices.", http.StatusInternalServerError, model.WithError(nRes.Error))
	}

	err := datastore.Provider(ctx).DeleteDevices(
		datastore.DeleteDevicesOptionWithLogicalDeleted(time.Now().Unix()),
		datastore.DeleteDevicesOptionFilterByUserID(req.UserID),
		datastore.DeleteDevicesOptionFilterByPlatform(req.Platform),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete devices.", http.StatusInternalServerError, model.WithError(err))
	}

	go unsubscribeByDevice(ctx, device, nil)

	return nil
}
