package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/pbroker"
	"github.com/swagchat/chat-api/utils"
)

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
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptions(
		datastore.SelectDeletedSubscriptionsOptionFilterByUserID(device.UserID),
		datastore.SelectDeletedSubscriptionsOptionFilterByPlatform(device.Platform),
	)
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
		ctx = context.WithValue(ctx, config.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(config.CtxRoomUser).(*model.RoomUser)
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
		ctx = context.WithValue(ctx, config.CtxSubscription, subscription)
		d.Work(ctx, func(ctx context.Context) {
			targetSubscription := ctx.Value(config.CtxSubscription).(*model.Subscription)
			nRes := <-np.Unsubscribe(targetSubscription.NotificationSubscriptionID)
			if nRes.Error != nil {
				errCh <- nRes.Error
			}
			err := dp.DeleteSubscriptions(
				datastore.DeleteSubscriptionsOptionFilterByRoomID(targetSubscription.RoomID),
				datastore.DeleteSubscriptionsOptionFilterByUserID(targetSubscription.UserID),
				datastore.DeleteSubscriptionsOptionFilterByPlatform(targetSubscription.Platform),
			)
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

func unsubscribeByUserID(ctx context.Context, userID string) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptions(datastore.SelectDeletedSubscriptionsOptionFilterByUserID(userID))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	unsubscribe(ctx, subscriptions)
}

func unsubscribeByRoomID(ctx context.Context, roomID string, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptions(datastore.SelectDeletedSubscriptionsOptionFilterByRoomID(roomID))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	<-unsubscribe(ctx, subscriptions)
	if wg != nil {
		wg.Done()
	}
}

func publishUserJoin(ctx context.Context, roomID string) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID, datastore.SelectRoomOptionWithUsers(true))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	go func() {
		userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
			datastore.SelectUserIDsOfRoomUserOptionWithRoomID(roomID),
		)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		buffer := new(bytes.Buffer)
		json.NewEncoder(buffer).Encode(room.Users)
		rtmEvent := &pbroker.RTMEvent{
			Type:    pbroker.UserJoin,
			Payload: buffer.Bytes(),
			UserIDs: userIDs,
		}
		err = pbroker.Provider(ctx).PublishMessage(rtmEvent)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}()
}

func subscribeByRoomUsers(ctx context.Context, roomUsers []*model.RoomUser) {
	doneChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, config.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(config.CtxRoomUser).(*model.RoomUser)

			devices, err := datastore.Provider(ctx).SelectDevices(datastore.SelectDevicesOptionFilterByUserID(ru.UserID))
			if err != nil {
				errChan <- err
			}
			if devices != nil {
				for _, d := range devices {
					if d.Token != "" {
						if d.NotificationDeviceID == "" {
							nRes := <-notification.Provider(ctx).CreateEndpoint("", d.Platform, d.Token)
							if nRes.Error != nil {
								errChan <- nRes.Error
							} else {
								d.NotificationDeviceID = *nRes.Data.(*string)
								err := datastore.Provider(ctx).UpdateDevice(d)
								if err != nil {
									errChan <- err
								}
							}
						}
						go subscribe(ctx, roomUsers, d)
					}
				}
			}
			doneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case err := <-errChan:
				logger.Error(err.Error())
				return
			}
		})
	}
	d.Wait()
	return
}

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*model.RoomUser) {
	doneChan := make(chan bool, 1)
	errResChan := make(chan *model.ErrorResponse, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, config.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(config.CtxRoomUser).(*model.RoomUser)
			err := datastore.Provider(ctx).DeleteRoomUsers(
				datastore.DeleteRoomUsersOptionFilterByRoomIDs([]string{ru.RoomID}),
				datastore.DeleteRoomUsersOptionFilterByUserIDs([]string{ru.UserID}),
			)
			if err != nil {
				errRes := model.NewErrorResponse("Failed to unsubscribe.", http.StatusInternalServerError, model.WithError(err))
				errResChan <- errRes
			}

			devices, err := datastore.Provider(ctx).SelectDevices(datastore.SelectDevicesOptionFilterByUserID(ru.UserID))
			if err != nil {
				errRes := model.NewErrorResponse("Failed to unsubscribe.", http.StatusInternalServerError, model.WithError(err))
				errResChan <- errRes
			}
			if devices != nil {
				for _, d := range devices {
					subscription, err := datastore.Provider(ctx).SelectSubscription(ru.RoomID, ru.UserID, d.Platform)
					if err != nil {
						errRes := model.NewErrorResponse("Failed to unsubscribe.", http.StatusInternalServerError, model.WithError(err))
						errResChan <- errRes
					}
					go unsubscribe(ctx, []*model.Subscription{subscription})
				}
			}
			doneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case <-errResChan:
				return
			}
		})
	}
	d.Wait()
	return
}

func createTopic(ctx context.Context, roomID string) (string, *model.ErrorResponse) {
	nRes := <-notification.Provider(ctx).CreateTopic(roomID)
	if nRes.Error != nil {
		return "", model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(nRes.Error))
	}
	if nRes.Data == nil {
		return "", nil
	}
	return *nRes.Data.(*string), nil
}
