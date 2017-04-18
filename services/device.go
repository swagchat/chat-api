package services

import (
	"context"
	"log"
	"net/http"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func CreateDevice(userId string, post *models.Device) (*models.Device, *models.ProblemDetail) {
	_, pd := SelectUser(userId)
	if pd != nil {
		return nil, pd
	}

	if pd := post.IsValid(); pd != nil {
		return nil, pd
	}

	np := notification.GetProvider()
	nRes := <-np.CreateEndpoint(userId, post.Platform, post.Token)
	if nRes.ProblemDetail != nil {
		return nil, nRes.ProblemDetail
	}
	notificationDeviceId := ""
	if nRes.Data != nil {
		notificationDeviceId = *nRes.Data.(*string)
	}

	post.BeforeSave(userId, notificationDeviceId)
	dp := datastore.GetProvider()
	dRes := <-dp.DeviceInsert(post)

	device := dRes.Data.(*models.Device)
	pd = dRes.ProblemDetail

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go subscribeAllRoomByDevice(ctx, device)

	return device, pd
}

func DeleteDevice(userId string, platform int) *models.ProblemDetail {
	_, pd := SelectUser(userId)
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

	dp := datastore.GetProvider()
	dRes := <-dp.DeviceDelete(userId, platform)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go unsubscribeAllRoomByDevice(ctx, device)

	return nil
}

func subscribeAllRoomByDevice(ctx context.Context, device *models.Device) {
	dp := datastore.GetProvider()
	dRes := <-dp.RoomUsersSelect(nil, []string{device.UserId})
	if dRes.ProblemDetail != nil {
		// TODO log
	}
	Subscribe(dRes.Data.([]*models.RoomUser), device)
}

func unsubscribeAllRoomByDevice(ctx context.Context, device *models.Device) {
	dp := datastore.GetProvider()
	dRes := <-dp.SubscriptionSelectByUserIdAndPlatform(device.UserId, device.Platform)
	if dRes.ProblemDetail != nil {
		// TODO log
	}
	subscriptions := dRes.Data.([]*models.Subscription)

	dp.SubscriptionUpdateDeletedByUserIdAndPlatform(device.UserId, device.Platform)
	if dRes.ProblemDetail != nil {
		// TODO log
	}
	Unsubscribe(subscriptions)
}

//	allDeleteFlg := 0
//	np := notification.GetProvider()
//	if put.Devices != nil {
//		for _, requestDevice := range put.Devices {
//			isExistDevice := false
//			for _, currentDevice := range user.Devices {
//				if requestDevice.Platform == currentDevice.Platform {
//					isExistDevice = true
//
//					if requestDevice.Token != nil && currentDevice.Token != requestDevice.Token {
//						log.Println("----------------> デバイストークン変更処理")
//
//						if !models.IsValidDevicePlatform(requestDevice.Platform) {
//							return nil, &models.ProblemDetail{
//								Title:     "Request parameter error. (Create user item)",
//								Status:    http.StatusBadRequest,
//								ErrorName: models.ERROR_NAME_INVALID_PARAM,
//								InvalidParams: []models.InvalidParam{
//									models.InvalidParam{
//										Name:   "device.platform",
//										Reason: "platform is invalid. Currently only 1(iOS) and 2(Android) are supported.",
//									},
//								},
//							}
//						}
//						nRes := <-np.CreateEndpoint(*requestDevice.Token)
//						if nRes.ProblemDetail != nil {
//							return nil, &models.ProblemDetail{
//								Title:     "Updating user item error. (Update user item)",
//								Status:    http.StatusInternalServerError,
//								ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//							}
//						}
//						if nRes.Data == nil {
//							return nil, &models.ProblemDetail{
//								Title:     "Creating notification endpoint. (Update user item)",
//								Status:    http.StatusInternalServerError,
//								ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//							}
//						}
//						notificationDeviceId := nRes.Data.(*string)
//
//						nowDatetime := time.Now().UnixNano()
//						tmpDevice := &models.Device{
//							UserId:               userId,
//							Platform:             requestDevice.Platform,
//							Token:                requestDevice.Token,
//							NotificationDeviceId: notificationDeviceId,
//							Created:              nowDatetime,
//							Modified:             nowDatetime,
//						}
//						user.Devices = []*models.Device{tmpDevice}
//						dRes := <-dp.UserInsert(user)
//						if dRes.ProblemDetail != nil {
//							return nil, dRes.ProblemDetail
//						}
//
//						dRes = <-dp.RoomUsersSelect(nil, []string{userId})
//						if dRes.ProblemDetail != nil {
//							return nil, dRes.ProblemDetail
//						}
//						ruRes := subscribeAllRoom(dRes.Data.([]*models.RoomUser), *notificationDeviceId)
//						if len(ruRes.Errors) > 0 {
//							return nil, &models.ProblemDetail{
//								Title:     "Updating user item error. (Update user item)",
//								Status:    http.StatusInternalServerError,
//								ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//							}
//						}
//					}
//
//					if requestDevice.Token == nil {
//						log.Println("----------------> デバイストークン削除処理")
//						if requestDevice.NotificationDeviceId != nil {
//							nRes := <-np.DeleteEndpoint(*requestDevice.NotificationDeviceId)
//							if nRes.ProblemDetail != nil {
//								return nil, nRes.ProblemDetail
//							}
//						}
//						dRes := <-dp.DeviceDelete(user.UserId, requestDevice.Platform)
//						if dRes.ProblemDetail != nil {
//							return nil, dRes.ProblemDetail
//						}
//						allDeleteFlg++
//					}
//				}
//			}
//			if isExistDevice {
//				log.Println("----------------> デバイストークン新規追加処理")
//				/*
//					dRes := <-dp.RoomUsersSelect(nil, []string{userId})
//					if dRes.ProblemDetail != nil {
//						return nil, dRes.ProblemDetail
//					}
//					ruRes := subscribeAllRoom(dRes.Data.([]*models.RoomUser), *notificationDeviceId)
//					if len(ruRes.Errors) > 0 {
//						return nil, &models.ProblemDetail{
//							Title:     "Updating user item error. (Update user item)",
//							Status:    http.StatusInternalServerError,
//							ErrorName: models.ERROR_NAME_NOTIFICATION_ERROR,
//						}
//					}
//				*/
//			}
//		}
//	}
//
//	if allDeleteFlg == len(user.Devices) {
//		dRes = <-dp.RoomUsersSelect(nil, []string{userId})
//		if dRes.ProblemDetail != nil {
//			return nil, dRes.ProblemDetail
//		}
//		go deleteRoomUsers(dRes.Data.([]*models.RoomUser))
//	}
//
//	if user.DeviceToken != nil && *put.DeviceToken == "" {
//		if user.NotificationDeviceId != nil {
//			nRes := <-np.DeleteEndpoint(*user.NotificationDeviceId)
//			if nRes.ProblemDetail != nil {
//				return nil, nRes.ProblemDetail
//			}
//		}
//	} else if (user.DeviceToken == nil && *put.DeviceToken != "") ||
//		(user.DeviceToken != nil && (user.DeviceToken != put.DeviceToken)) {
//	}

func SelectDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.DeviceSelect(userId, platform)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	return dRes.Data.(*models.Device), nil
}

func Subscribe(roomUsers []*models.RoomUser, device *models.Device) {
	np := notification.GetProvider()
	dp := datastore.GetProvider()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subscribeDoneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)
	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)
			dRes := <-dp.RoomSelect(ru.RoomId)
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			} else {
				room := dRes.Data.(*models.Room)
				nRes := <-np.Subscribe(*room.NotificationTopicId, device.NotificationDeviceId)
				if nRes.ProblemDetail != nil {
					pdChan <- nRes.ProblemDetail
				} else {
					if nRes.Data != nil {
						notificationSubscriptionId := nRes.Data.(*string)
						subscription := &models.Subscription{
							RoomId:                     ru.RoomId,
							UserId:                     ru.UserId,
							Platform:                   device.Platform,
							NotificationSubscriptionId: notificationSubscriptionId,
						}
						dRes := <-dp.SubscriptionInsert(subscription)
						if dRes.ProblemDetail != nil {
							pdChan <- dRes.ProblemDetail
						}
					}
				}
			}
			subscribeDoneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-subscribeDoneChan:
				return
			case pd := <-pdChan:
				// TODO log
				log.Println(pd)
				return
			}
		})
	}
	d.Wait()
}

func Unsubscribe(subscriptions []*models.Subscription) {
	np := notification.GetProvider()
	dp := datastore.GetProvider()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	unsubscribeDoneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, subscription := range subscriptions {
		ctx = context.WithValue(ctx, "subscription", subscription)
		d.Work(ctx, func(ctx context.Context) {
			targetSubscription := ctx.Value("subscription").(*models.Subscription)
			nRes := <-np.Unsubscribe(*targetSubscription.NotificationSubscriptionId)
			if nRes.ProblemDetail != nil {
				pdChan <- nRes.ProblemDetail
			}
			dRes := <-dp.SubscriptionDelete(subscription)
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			}
			unsubscribeDoneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-unsubscribeDoneChan:
				return
			case pd := <-pdChan:
				// TODO log
				log.Println(pd)
				return
			}
		})
	}
	d.Wait()
}
