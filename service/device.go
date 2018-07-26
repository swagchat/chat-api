package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
)

// GetDevices is get devices
func GetDevices(ctx context.Context, userID string) (*model.Devices, *model.ProblemDetail) {
	devices, err := datastore.Provider(ctx).SelectDevices(userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get device failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	return &model.Devices{
		Devices: devices,
	}, nil
}

// GetDevice is get device
func GetDevice(ctx context.Context, userID string, platform int32) (*model.Device, *model.ProblemDetail) {
	device, pd := selectDevice(ctx, userID, platform)
	if pd != nil {
		return nil, pd
	}

	return device, nil
}

// PutDevice is put device
func PutDevice(ctx context.Context, put *model.Device) (*model.Device, *model.ErrorResponse) {
	// if pd := put.IsValid(); pd != nil {
	// 	return nil, pd
	// }

	// // User existence check
	// _, pd := selectUser(ctx, put.UserID)
	// if pd != nil {
	// 	return nil, pd
	// }

	isExist := true
	device, pd := selectDevice(ctx, put.UserID, put.Platform)
	if pd != nil {
		errRes := &model.ErrorResponse{}
		errRes.Message = "Failed to update device."
		errRes.Status = http.StatusInternalServerError
		errRes.Error = pd.Error
		return nil, errRes
	}
	if device == nil {
		isExist = false
	}

	if !isExist || (device.Token != put.Token) {
		// When using another user on the same device, delete the notification information
		// of the olderuser in order to avoid duplication of the device token
		deleteDevices, err := datastore.Provider(ctx).SelectDevicesByToken(put.Token)
		if err != nil {
			errRes := &model.ErrorResponse{}
			errRes.Message = "Failed to update device."
			errRes.Status = http.StatusInternalServerError
			errRes.Error = err
			return nil, errRes
		}
		if deleteDevices != nil {
			wg := &sync.WaitGroup{}
			for _, deleteDevice := range deleteDevices {
				nRes := <-notification.Provider().DeleteEndpoint(deleteDevice.NotificationDeviceID)
				if nRes.Error != nil {
					errRes := &model.ErrorResponse{}
					errRes.Error = nRes.Error
					return nil, errRes
				}
				err := datastore.Provider(ctx).DeleteDevice(deleteDevice.UserID, deleteDevice.Platform)
				if err != nil {
					errRes := &model.ErrorResponse{}
					errRes.Message = "Failed to update device."
					errRes.Status = http.StatusInternalServerError
					errRes.Error = err
					return nil, errRes
				}
				wg.Add(1)
				go unsubscribeByDevice(ctx, deleteDevice, wg)
			}
			wg.Wait()
		}

		nRes := <-notification.Provider().CreateEndpoint(put.UserID, put.Platform, put.Token)
		if nRes.Error != nil {
			errRes := &model.ErrorResponse{}
			errRes.Message = "Failed to update device."
			errRes.Status = http.StatusInternalServerError
			errRes.Error = nRes.Error
			return nil, errRes
		}
		put.NotificationDeviceID = put.Token
		if nRes.Data != nil {
			put.NotificationDeviceID = *nRes.Data.(*string)
		}

		if isExist {
			err := datastore.Provider(ctx).UpdateDevice(put)
			if err != nil {
				errRes := &model.ErrorResponse{}
				errRes.Message = "Failed to update device."
				errRes.Status = http.StatusInternalServerError
				errRes.Error = err
				return nil, errRes
			}
			nRes = <-notification.Provider().DeleteEndpoint(device.NotificationDeviceID)
			if nRes.Error != nil {
				errRes := &model.ErrorResponse{}
				errRes.Message = "Failed to update device."
				errRes.Status = http.StatusInternalServerError
				errRes.Error = nRes.Error
				return nil, errRes
			}
			go func() {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go unsubscribeByDevice(ctx, device, wg)
				wg.Wait()
				go subscribeByDevice(ctx, put, nil)
			}()
		} else {
			device, err = datastore.Provider(ctx).InsertDevice(put)
			if err != nil {
				errRes := &model.ErrorResponse{}
				errRes.Message = "Failed to update device."
				errRes.Status = http.StatusInternalServerError
				errRes.Error = err
				return nil, errRes
			}
			go subscribeByDevice(ctx, device, nil)
		}
		return device, nil
	}

	return nil, nil
}

// DeleteDevice is delete device
func DeleteDevice(ctx context.Context, userID string, platform int32) *model.ErrorResponse {
	// User existence check
	// _, pd := selectUser(ctx, userID)
	// if pd != nil {
	// 	return pd
	// }

	device, pd := selectDevice(ctx, userID, platform)
	if pd != nil {
		errRes := &model.ErrorResponse{}
		errRes.Message = "Failed to delete device."
		errRes.Status = http.StatusInternalServerError
		errRes.Error = pd.Error
		return errRes
	}

	np := notification.Provider()
	nRes := <-np.DeleteEndpoint(device.NotificationDeviceID)
	if nRes.Error != nil {
		errRes := &model.ErrorResponse{}
		errRes.Message = "Failed to delete device."
		errRes.Status = http.StatusInternalServerError
		errRes.Error = nRes.Error
		return errRes
	}

	err := datastore.Provider(ctx).DeleteDevice(userID, platform)
	if err != nil {
		errRes := &model.ErrorResponse{}
		errRes.Message = "Failed to delete device."
		errRes.Status = http.StatusInternalServerError
		errRes.Error = err
		return errRes
	}

	go unsubscribeByDevice(ctx, device, nil)

	return nil
}

func selectDevice(ctx context.Context, userID string, platform int32) (*model.Device, *model.ProblemDetail) {
	device, err := datastore.Provider(ctx).SelectDevice(userID, platform)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get device failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if device == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
	}
	return device, nil
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
	np := notification.Provider()
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
					notificationTopicID, errRes := createTopic(room.RoomID)
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
	np := notification.Provider()
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
