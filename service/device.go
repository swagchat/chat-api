package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
)

// CreateDevice creates device
func CreateDevice(ctx context.Context, req *model.CreateDeviceRequest) *model.ErrorResponse {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.CreateDevice")
	defer span.Finish()

	return nil
}

// GetDevices gets devices
func GetDevices(ctx context.Context, req *model.GetDevicesRequest) (*model.DevicesResponse, *model.ErrorResponse) {
	span, _ := opentracing.StartSpanFromContext(ctx, "service.GetDevices")
	defer span.Finish()

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
	span, _ := opentracing.StartSpanFromContext(ctx, "service.UpdateDevice")
	defer span.Finish()

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
				err := datastore.Provider(ctx).DeleteDevice(deleteDevice.UserID, deleteDevice.Platform)
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
			device, err = datastore.Provider(ctx).InsertDevice(device)
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
	span, _ := opentracing.StartSpanFromContext(ctx, "service.DeleteDevices")
	defer span.Finish()

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

	err := datastore.Provider(ctx).DeleteDevice(req.UserID, req.Platform)
	if err != nil {
		return model.NewErrorResponse("Failed to delete devices.", http.StatusInternalServerError, model.WithError(err))
	}

	go unsubscribeByDevice(ctx, device, nil)

	return nil
}

func subscribeByDevice(ctx context.Context, device *model.Device, wg *sync.WaitGroup) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUsers(
		datastore.SelectRoomUsersOptionWithUserIDs([]string{device.UserID}),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if roomUser != nil {
		<-subscribe(ctx, roomUser, device)
	}
	if wg != nil {
		wg.Done()
	}
}

func unsubscribeByDevice(ctx context.Context, device *model.Device, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByUserIDAndPlatform(device.UserID, device.Platform)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	<-unsubscribe(ctx, subscriptions)
	if wg != nil {
		wg.Done()
	}
}

func subscribe(ctx context.Context, roomUsers []*model.RoomUser, device *model.Device) chan bool {
	np := notification.Provider(ctx)
	dp := datastore.Provider(ctx)
	doneCh := make(chan bool, 1)
	errCh := make(chan error, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*model.RoomUser)
			room, errRes := confirmRoomExist(ctx, ru.RoomID)
			if errRes != nil {
				errCh <- errRes.Error
			} else {
				if room.NotificationTopicID == "" {
					notificationTopicID, errRes := createTopic(ctx, room.RoomID)
					if errRes != nil {
						errCh <- errRes.Error
					}

					room.NotificationTopicID = notificationTopicID
					room.Modified = time.Now().Unix()
					err := datastore.Provider(ctx).UpdateRoom(room)
					if err != nil {
						errCh <- err
					}
				}
				nRes := <-np.Subscribe(room.NotificationTopicID, device.NotificationDeviceID)
				if nRes.Error != nil {
					errCh <- nRes.Error
				} else {
					if nRes.Data != nil {
						notificationSubscriptionID := nRes.Data.(*string)
						subscription := &model.Subscription{
							RoomID:                     ru.RoomID,
							UserID:                     ru.UserID,
							Platform:                   device.Platform,
							NotificationSubscriptionID: *notificationSubscriptionID,
						}
						subscription, err := dp.InsertSubscription(subscription)
						if err != nil {
							errCh <- err
						} else {
							doneCh <- true
						}
					}
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case err := <-errCh:
				logger.Error(err.Error())
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}

func unsubscribe(ctx context.Context, subscriptions []*model.Subscription) chan bool {
	np := notification.Provider(ctx)
	dp := datastore.Provider(ctx)
	doneCh := make(chan bool, 1)
	errCh := make(chan error, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, subscription := range subscriptions {
		ctx = context.WithValue(ctx, utils.CtxSubscription, subscription)
		d.Work(ctx, func(ctx context.Context) {
			targetSubscription := ctx.Value(utils.CtxSubscription).(*model.Subscription)
			nRes := <-np.Unsubscribe(targetSubscription.NotificationSubscriptionID)
			if nRes.Error != nil {
				errCh <- nRes.Error
			}
			err := dp.DeleteSubscription(targetSubscription)
			if err != nil {
				errCh <- err
			} else {
				doneCh <- true
			}

			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case err := <-errCh:
				logger.Error(err.Error())
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}
