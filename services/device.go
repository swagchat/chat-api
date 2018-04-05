package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
)

func GetDevices(userId string) (*models.Devices, *models.ProblemDetail) {
	devices, err := datastore.Provider().SelectDevices(userId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get device failed",
			Status: http.StatusInternalServerError,
		}
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			ProblemDetail: pd,
			Error:         err,
		})
		return nil, pd
	}

	return &models.Devices{
		Devices: devices,
	}, nil
}

func GetDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	device, pd := selectDevice(userId, platform)
	if pd != nil {
		return nil, pd
	}

	return device, nil
}

func PutDevice(put *models.Device) (*models.Device, *models.ProblemDetail) {
	if pd := put.IsValid(); pd != nil {
		return nil, pd
	}

	// User existence check
	_, pd := selectUser(put.UserId)
	if pd != nil {
		return nil, pd
	}

	isExist := true
	device, pd := selectDevice(put.UserId, put.Platform)
	if device == nil {
		isExist = false
	}

	if !isExist || (device.Token != put.Token) {
		ctx, _ := context.WithCancel(context.Background())

		// When using another user on the same device, delete the notification information
		// of the olderuser in order to avoid duplication of the device token
		deleteDevices, err := datastore.Provider().SelectDevicesByToken(put.Token)
		if err != nil {
			pd := &models.ProblemDetail{
				Title:  "Update device failed",
				Status: http.StatusInternalServerError,
			}
			logging.Log(zapcore.ErrorLevel, &logging.AppLog{
				ProblemDetail: pd,
				Error:         err,
			})
			return nil, pd
		}
		if deleteDevices != nil {
			wg := &sync.WaitGroup{}
			for _, deleteDevice := range deleteDevices {
				nRes := <-notification.Provider().DeleteEndpoint(deleteDevice.NotificationDeviceId)
				if nRes.ProblemDetail != nil {
					return nil, nRes.ProblemDetail
				}
				err := datastore.Provider().DeleteDevice(deleteDevice.UserId, deleteDevice.Platform)
				if err != nil {
					pd := &models.ProblemDetail{
						Title:  "Update device failed",
						Status: http.StatusInternalServerError,
					}
					logging.Log(zapcore.ErrorLevel, &logging.AppLog{
						ProblemDetail: pd,
						Error:         err,
					})
					return nil, pd
				}
				wg.Add(1)
				go unsubscribeByDevice(ctx, deleteDevice, wg)
			}
			wg.Wait()
		}

		nRes := <-notification.Provider().CreateEndpoint(put.UserId, put.Platform, put.Token)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}
		put.NotificationDeviceId = put.Token
		if nRes.Data != nil {
			put.NotificationDeviceId = *nRes.Data.(*string)
		}

		if isExist {
			err := datastore.Provider().UpdateDevice(put)
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Update device failed",
					Status: http.StatusInternalServerError,
				}
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
					Error:         err,
				})
				return nil, pd
			}
			nRes = <-notification.Provider().DeleteEndpoint(device.NotificationDeviceId)
			if nRes.ProblemDetail != nil {
				return nil, nRes.ProblemDetail
			}
			go func() {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go unsubscribeByDevice(ctx, device, wg)
				wg.Wait()
				go subscribeByDevice(ctx, put, nil)
			}()
		} else {
			device, err = datastore.Provider().InsertDevice(put)
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Update device failed",
					Status: http.StatusInternalServerError,
				}
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
					Error:         err,
				})
				return nil, pd
			}
			go subscribeByDevice(ctx, device, nil)
		}
		return device, nil
	}

	return nil, nil
}

func DeleteDevice(userId string, platform int) *models.ProblemDetail {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return pd
	}

	device, pd := selectDevice(userId, platform)
	if pd != nil {
		return pd
	}

	np := notification.Provider()
	nRes := <-np.DeleteEndpoint(device.NotificationDeviceId)
	if nRes.ProblemDetail != nil {
		return nRes.ProblemDetail
	}

	err := datastore.Provider().DeleteDevice(userId, platform)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete device failed",
			Status: http.StatusInternalServerError,
		}
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			ProblemDetail: pd,
			Error:         err,
		})
		return pd
	}

	ctx, _ := context.WithCancel(context.Background())
	go unsubscribeByDevice(ctx, device, nil)

	return nil
}

func selectDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	device, err := datastore.Provider().SelectDevice(userId, platform)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get device failed",
			Status: http.StatusInternalServerError,
		}
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			ProblemDetail: pd,
			Error:         err,
		})
		return nil, pd
	}
	if device == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}
	return device, nil
}

func subscribeByDevice(ctx context.Context, device *models.Device, wg *sync.WaitGroup) {
	roomUser, err := datastore.Provider().SelectRoomUsersByUserId(device.UserId)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	if roomUser != nil {
		<-subscribe(ctx, roomUser, device)
	}
	if wg != nil {
		wg.Done()
	}
}

func unsubscribeByDevice(ctx context.Context, device *models.Device, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider().SelectDeletedSubscriptionsByUserIdAndPlatform(device.UserId, device.Platform)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	<-unsubscribe(ctx, subscriptions)
	if wg != nil {
		wg.Done()
	}
}

func subscribe(ctx context.Context, roomUsers []*models.RoomUser, device *models.Device) chan bool {
	np := notification.Provider()
	dp := datastore.Provider()
	doneCh := make(chan bool, 1)
	pdCh := make(chan *models.ProblemDetail, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)
			room, pd := selectRoom(ru.RoomId)
			if pd != nil {
				pdCh <- pd
			} else {
				if room.NotificationTopicId == "" {
					notificationTopicId, pd := createTopic(room.RoomId)
					if pd != nil {
						pdCh <- pd
					}

					room.NotificationTopicId = notificationTopicId
					room.Modified = time.Now().Unix()
					_, err := datastore.Provider().UpdateRoom(room)
					if err != nil {
						pd := &models.ProblemDetail{
							Status: http.StatusInternalServerError,
							Title:  "Update room failed",
						}
						pdCh <- pd
					}
				}
				nRes := <-np.Subscribe(room.NotificationTopicId, device.NotificationDeviceId)
				if nRes.ProblemDetail != nil {
					pdCh <- nRes.ProblemDetail
				} else {
					if nRes.Data != nil {
						notificationSubscriptionId := nRes.Data.(*string)
						subscription := &models.Subscription{
							RoomId:                     ru.RoomId,
							UserId:                     ru.UserId,
							Platform:                   device.Platform,
							NotificationSubscriptionId: *notificationSubscriptionId,
						}
						subscription, err := dp.InsertSubscription(subscription)
						if err != nil {
							pd := &models.ProblemDetail{
								Title:  "User registration failed",
								Status: http.StatusInternalServerError,
								Error:  err,
							}
							pdCh <- pd
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
			case pd := <-pdCh:
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
				})
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}

func unsubscribe(ctx context.Context, subscriptions []*models.Subscription) chan bool {
	np := notification.Provider()
	dp := datastore.Provider()
	doneCh := make(chan bool, 1)
	pdCh := make(chan *models.ProblemDetail, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, subscription := range subscriptions {
		ctx = context.WithValue(ctx, "subscription", subscription)
		d.Work(ctx, func(ctx context.Context) {
			targetSubscription := ctx.Value("subscription").(*models.Subscription)
			nRes := <-np.Unsubscribe(targetSubscription.NotificationSubscriptionId)
			if nRes.ProblemDetail != nil {
				pdCh <- nRes.ProblemDetail
			}
			err := dp.DeleteSubscription(targetSubscription)
			if err != nil {
				pd := &models.ProblemDetail{
					Error: err,
				}
				pdCh <- pd
			} else {
				doneCh <- true
			}

			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case pd := <-pdCh:
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
				})
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}
