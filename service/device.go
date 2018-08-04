package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/tracer"
)

// CreateDevice creates device
func CreateDevice(ctx context.Context, req *model.CreateDeviceRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("CreateDevice", "service")
	defer tracer.Provider(ctx).Finish(span)

	return nil
}

// GetDevices gets devices
func GetDevices(ctx context.Context, req *model.GetDevicesRequest) (*model.DevicesResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("GetDevices", "service")
	defer tracer.Provider(ctx).Finish(span)

	devices, err := datastore.Provider(ctx).SelectDevices(datastore.SelectDevicesOptionFilterByUserID(req.UserID))
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get devices.", http.StatusInternalServerError, model.WithError(err))
	}

	return &model.DevicesResponse{
		Devices: devices,
	}, nil
}

// UpdateDevice updates device
func UpdateDevice(ctx context.Context, req *model.UpdateDeviceRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("UpdateDevice", "service")
	defer tracer.Provider(ctx).Finish(span)

	errRes := req.Validate()
	if errRes != nil {
		return errRes
	}

	device, errRes := confirmDeviceExist(ctx, req.UserID, req.Platform)
	if errRes != nil {
		errRes.Message = "Failed to update device."
		return errRes
	}

	if device == nil || (device.Token != req.Token) {
		// When using another user on the same device, delete the notification information
		// of the olderuser in order to avoid duplication of the device token
		deleteDevices, err := datastore.Provider(ctx).SelectDevices(datastore.SelectDevicesOptionFilterByToken(req.Token))
		if err != nil {
			return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(err))
		}
		if deleteDevices != nil {
			wg := &sync.WaitGroup{}
			for _, deleteDevice := range deleteDevices {
				nRes := <-notification.Provider(ctx).DeleteEndpoint(deleteDevice.NotificationDeviceID)
				if nRes.Error != nil {
					return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(nRes.Error))
				}
				deleted := time.Now().Unix()
				err := datastore.Provider(ctx).DeleteDevices(
					datastore.DeleteDevicesOptionWithLogicalDeleted(deleted),
					datastore.DeleteDevicesOptionFilterByUserID(deleteDevice.UserID),
					datastore.DeleteDevicesOptionFilterByPlatform(deleteDevice.Platform),
				)
				if err != nil {
					return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(err))
				}
				wg.Add(1)
				go unsubscribeByDevice(ctx, deleteDevice, wg)
			}
			wg.Wait()
		}

		nRes := <-notification.Provider(ctx).CreateEndpoint(req.UserID, req.Platform, req.Token)
		if nRes.Error != nil {
			return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(nRes.Error))
		}
		device.NotificationDeviceID = req.Token
		if nRes.Data != nil {
			device.NotificationDeviceID = *nRes.Data.(*string)
		}

		if device != nil {
			err := datastore.Provider(ctx).UpdateDevice(device)
			if err != nil {
				return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(err))
			}
			nRes = <-notification.Provider(ctx).DeleteEndpoint(device.NotificationDeviceID)
			if nRes.Error != nil {
				return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(nRes.Error))
			}
			go func() {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go unsubscribeByDevice(ctx, device, wg)
				wg.Wait()
				go subscribeByDevice(ctx, device, nil)
			}()
		} else {
			err = datastore.Provider(ctx).InsertDevice(device)
			if err != nil {
				return model.NewErrorResponse("Failed to update device.", http.StatusInternalServerError, model.WithError(err))
			}
			go subscribeByDevice(ctx, device, nil)
		}
		return nil
	}

	return nil
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

	deleted := time.Now().Unix()
	err := datastore.Provider(ctx).DeleteDevices(
		datastore.DeleteDevicesOptionWithLogicalDeleted(deleted),
		datastore.DeleteDevicesOptionFilterByUserID(req.UserID),
		datastore.DeleteDevicesOptionFilterByPlatform(req.Platform),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete devices.", http.StatusInternalServerError, model.WithError(err))
	}

	go unsubscribeByDevice(ctx, device, nil)

	return nil
}
