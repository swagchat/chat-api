package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func CreateDevice(userId string, platform int, post *models.Device) (*models.Device, *models.ProblemDetail) {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	post.UserId = userId
	post.Platform = platform
	if pd := post.IsValid(); pd != nil {
		return nil, pd
	}

	nRes := <-notification.GetProvider().CreateEndpoint(userId, platform, post.Token)
	if nRes.ProblemDetail != nil {
		return nil, nRes.ProblemDetail
	}
	notificationDeviceId := post.Token
	if nRes.Data != nil {
		notificationDeviceId = *nRes.Data.(*string)
	}

	post.NotificationDeviceId = notificationDeviceId
	dRes := <-datastore.GetProvider().InsertDevice(post)
	device := dRes.Data.(*models.Device)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go subscribeByDevice(ctx, device)

	return device, dRes.ProblemDetail
}

func GetDevices() (*models.Devices, *models.ProblemDetail) {
	dRes := <-datastore.GetProvider().SelectDevices()
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

func PutDevice(userId string, platform int, put *models.Device) (*models.Device, *models.ProblemDetail) {
	// User existence check
	_, pd := selectUser(userId)
	if pd != nil {
		return nil, pd
	}

	device, pd := SelectDevice(userId, platform)
	if pd != nil {
		return nil, pd
	}

	if put.Token != "" && device.Token != put.Token {
		np := notification.GetProvider()
		nRes := <-np.DeleteEndpoint(device.NotificationDeviceId)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}
		nRes = <-np.CreateEndpoint(userId, platform, put.Token)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}
		notificationDeviceId := put.Token
		if nRes.Data != nil {
			notificationDeviceId = *nRes.Data.(*string)
		}

		newDevice := &models.Device{
			UserId:               userId,
			Platform:             platform,
			Token:                put.Token,
			NotificationDeviceId: notificationDeviceId,
		}
		dRes := <-datastore.GetProvider().UpdateDevice(newDevice)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go unsubscribeByDevice(ctx, device)
		go subscribeByDevice(ctx, newDevice)
		return newDevice, nil
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

	dRes := <-datastore.GetProvider().DeleteDevice(userId, platform)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go unsubscribeByDevice(ctx, device)

	return nil
}

func SelectDevice(userId string, platform int) (*models.Device, *models.ProblemDetail) {
	dp := datastore.GetProvider()
	dRes := <-dp.SelectDevice(userId, platform)
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

func subscribeByDevice(ctx context.Context, device *models.Device) {
	dRes := <-datastore.GetProvider().SelectRoomUsersByUserId(device.UserId)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	if dRes.Data != nil {
		subscribe(ctx, dRes.Data.([]*models.RoomUser), device)
	}
}

func unsubscribeByDevice(ctx context.Context, device *models.Device) {
	dRes := <-datastore.GetProvider().SelectSubscriptionsByUserIdAndPlatform(device.UserId, device.Platform)
	if dRes.ProblemDetail != nil {
		pdBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("problemDetail", string(pdBytes)),
			zap.String("err", fmt.Sprintf("%+v", dRes.ProblemDetail.Error)),
		)
	}
	unsubscribe(ctx, dRes.Data.([]*models.Subscription))
}

func subscribe(ctx context.Context, roomUsers []*models.RoomUser, device *models.Device) {
	np := notification.GetProvider()
	dp := datastore.GetProvider()
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)
			dRes := <-dp.SelectRoom(ru.RoomId)
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			} else {
				room := dRes.Data.(*models.Room)
				nRes := <-np.Subscribe(room.NotificationTopicId, device.NotificationDeviceId)
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
						dRes := <-dp.InsertSubscription(subscription)
						if dRes.ProblemDetail != nil {
							pdChan <- dRes.ProblemDetail
						}
					}
				}
			}
			doneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case pd := <-pdChan:
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
}

func unsubscribe(ctx context.Context, subscriptions []*models.Subscription) {
	np := notification.GetProvider()
	dp := datastore.GetProvider()
	doneChan := make(chan bool, 1)
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
			dRes := <-dp.DeleteSubscription(subscription)
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			}
			doneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case pd := <-pdChan:
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
}
