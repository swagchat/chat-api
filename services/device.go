package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
)

func GetDevices(userId string) (*models.Devices, *models.ProblemDetail) {
	dRes := datastore.GetProvider().SelectDevices(userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	devices := &models.Devices{
		Devices: dRes.Data.([]*models.Device),
	}
	return devices, nil
}

func GetDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	user, pd := SelectDevice(userId, platform)
	return user, pd
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
	device, pd := SelectDevice(put.UserId, put.Platform)
	if device == nil {
		isExist = false
	}

	if !isExist || (device.Token != put.Token) {
		ctx, _ := context.WithCancel(context.Background())

		// When using another user on the same device, delete the notification information
		// of the olderuser in order to avoid duplication of the device token
		dRes := datastore.GetProvider().SelectDevicesByToken(put.Token)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}
		if dRes.Data != nil {
			wg := &sync.WaitGroup{}
			deleteDevices := dRes.Data.([]*models.Device)
			for _, deleteDevice := range deleteDevices {
				nRes := <-notification.GetProvider().DeleteEndpoint(deleteDevice.NotificationDeviceId)
				if nRes.ProblemDetail != nil {
					return nil, nRes.ProblemDetail
				}
				dRes := datastore.GetProvider().DeleteDevice(deleteDevice.UserId, deleteDevice.Platform)
				if dRes.ProblemDetail != nil {
					return nil, dRes.ProblemDetail
				}
				wg.Add(1)
				go unsubscribeByDevice(ctx, deleteDevice, wg)
			}
			wg.Wait()
		}

		nRes := <-notification.GetProvider().CreateEndpoint(put.UserId, put.Platform, put.Token)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}
		put.NotificationDeviceId = put.Token
		if nRes.Data != nil {
			put.NotificationDeviceId = *nRes.Data.(*string)
		}

		if isExist {
			dRes := datastore.GetProvider().UpdateDevice(put)
			if dRes.ProblemDetail != nil {
				return nil, dRes.ProblemDetail
			}
			nRes = <-notification.GetProvider().DeleteEndpoint(device.NotificationDeviceId)
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
			dRes := datastore.GetProvider().InsertDevice(put)
			if dRes.ProblemDetail != nil {
				return nil, dRes.ProblemDetail
			}
			go subscribeByDevice(ctx, put, nil)
		}
		return put, nil
	} else {
		return nil, nil
	}
}

func DeleteDevice(userId string, platform int) *models.ProblemDetail {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return pd
	}

	device, pd := SelectDevice(userId, platform)
	if pd != nil {
		return pd
	}

	np := notification.GetProvider()
	nRes := <-np.DeleteEndpoint(device.NotificationDeviceId)
	if nRes.ProblemDetail != nil {
		return nRes.ProblemDetail
	}

	dRes := datastore.GetProvider().DeleteDevice(userId, platform)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	ctx, _ := context.WithCancel(context.Background())
	go unsubscribeByDevice(ctx, device, nil)

	return nil
}

func SelectDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	dRes := datastore.GetProvider().SelectDevice(userId, platform)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, nil
	}
	return dRes.Data.(*models.Device), nil
}

func subscribeByDevice(ctx context.Context, device *models.Device, wg *sync.WaitGroup) {
	dRes := datastore.GetProvider().SelectRoomUsersByUserId(device.UserId)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	if dRes.Data != nil {
		<-subscribe(ctx, dRes.Data.([]*models.RoomUser), device)
	}
	if wg != nil {
		wg.Done()
	}
}

func unsubscribeByDevice(ctx context.Context, device *models.Device, wg *sync.WaitGroup) {
	dRes := datastore.GetProvider().SelectDeletedSubscriptionsByUserIdAndPlatform(device.UserId, device.Platform)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	<-unsubscribe(ctx, dRes.Data.([]*models.Subscription))
	if wg != nil {
		wg.Done()
	}
}

func subscribe(ctx context.Context, roomUsers []*models.RoomUser, device *models.Device) chan bool {
	np := notification.GetProvider()
	dp := datastore.GetProvider()
	doneCh := make(chan bool, 1)
	pdCh := make(chan *models.ProblemDetail, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)
			dRes := dp.SelectRoom(ru.RoomId)
			if dRes.ProblemDetail != nil {
				pdCh <- dRes.ProblemDetail
			} else {
				room := dRes.Data.(*models.Room)
				if room.NotificationTopicId == "" {
					notificationTopicId, pd := createTopic(room.RoomId)
					if pd != nil {
						pdCh <- pd
					}

					room.NotificationTopicId = notificationTopicId
					room.Modified = time.Now().Unix()
					dRes := datastore.GetProvider().UpdateRoom(room)
					if dRes.ProblemDetail != nil {
						pdCh <- dRes.ProblemDetail
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
						dRes := dp.InsertSubscription(subscription)
						if dRes.ProblemDetail != nil {
							pdCh <- dRes.ProblemDetail
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
				pdBytes, _ := json.Marshal(pd)
				utils.AppLogger.Error("",
					zap.String("problemDetail", string(pdBytes)),
					zap.String("err", fmt.Sprintf("%+v", pd.Error)),
				)
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}

func unsubscribe(ctx context.Context, subscriptions []*models.Subscription) chan bool {
	np := notification.GetProvider()
	dp := datastore.GetProvider()
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
			dRes := dp.DeleteSubscription(targetSubscription)
			if dRes.ProblemDetail != nil {
				pdCh <- dRes.ProblemDetail
			} else {
				doneCh <- true
			}

			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case pd := <-pdCh:
				pdBytes, _ := json.Marshal(pd)
				utils.AppLogger.Error("",
					zap.String("problemDetail", string(pdBytes)),
					zap.String("err", fmt.Sprintf("%+v", pd.Error)),
				)
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}
