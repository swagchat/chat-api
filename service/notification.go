package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	logger "github.com/betchi/zapper"
	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/producer"
	"github.com/betchi/tracer"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func subscribe(ctx context.Context, roomUsers []*model.RoomUser, device *model.Device) chan bool {
	span := tracer.StartSpan(ctx, "subscribe", "service")
	defer tracer.Finish(span)

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
					room.ModifiedTimestamp = time.Now().Unix()
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

func subscribeByRoomUsers(ctx context.Context, roomUsers []*model.RoomUser) {
	span := tracer.StartSpan(ctx, "subscribeByRoomUsers", "service")
	defer tracer.Finish(span)

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

func unsubscribe(ctx context.Context, subscriptions []*model.Subscription) chan bool {
	span := tracer.StartSpan(ctx, "unsubscribe", "service")
	defer tracer.Finish(span)

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

func unsubscribeByUserID(ctx context.Context, userID string) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptions(
		datastore.SelectDeletedSubscriptionsOptionFilterByUserID(userID),
	)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	unsubscribe(ctx, subscriptions)
}

func unsubscribeByRoomID(ctx context.Context, roomID string, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptions(
		datastore.SelectDeletedSubscriptionsOptionFilterByRoomID(roomID),
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

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*model.RoomUser) {
	span := tracer.StartSpan(ctx, "unsubscribeByRoomUsers", "service")
	defer tracer.Finish(span)

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

func publishUserJoin(ctx context.Context, roomID string) {
	go func() {
		userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
			datastore.SelectUserIDsOfRoomUserOptionWithRoomID(roomID),
		)
		if err != nil {
			logger.Error(err.Error())
			return
		}

		for _, userID := range userIDs {
			miniRoom, err := datastore.Provider(ctx).SelectMiniRoom(roomID, userID)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			buffer := new(bytes.Buffer)
			json.NewEncoder(buffer).Encode(miniRoom)
			event := &scpb.EventData{
				Type:    scpb.EventType_RoomEvent,
				Data:    buffer.Bytes(),
				UserIDs: []string{userID},
			}
			err = producer.Provider(ctx).PublishMessage(event)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}
	}()
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
