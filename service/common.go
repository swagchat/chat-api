package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/pbroker"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

// func selectUser(ctx context.Context, userID string, opts ...datastore.SelectUserOption) (*model.User, *model.ProblemDetail) {
// 	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
// 	if err != nil {
// 		pd := &model.ProblemDetail{
// 			Message: "Get user failed",
// 			Status:  http.StatusInternalServerError,
// 			Error:   err,
// 		}
// 		return nil, pd
// 	}
// 	if user == nil {
// 		logger.Error(fmt.Sprintf("User does not exist. UserId[%s]", userID))
// 		return nil, &model.ProblemDetail{
// 			Status: http.StatusNotFound,
// 			Error:  errors.New("Not found"),
// 		}
// 	}
// 	return user, nil
// }

func confirmUserNotExist(ctx context.Context, userID string, opts ...datastore.SelectUserOption) (*model.User, *model.ErrorResponse) {
	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if user != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: fmt.Sprintf("That user already exist. userId[%s]", userID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return user, nil
}

func confirmUserExist(ctx context.Context, userID string, opts ...datastore.SelectUserOption) (*model.User, *model.ErrorResponse) {
	user, err := datastore.Provider(ctx).SelectUser(userID, opts...)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if user == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: fmt.Sprintf("That user is not exist. userId[%s]", userID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return user, nil
}

func unsubscribeByUserID(ctx context.Context, userID string) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByUserID(userID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	unsubscribe(ctx, subscriptions)
}

// ContactsAuthz is contacts authorize
func ContactsAuthz(ctx context.Context, requestUserID, resourceUserID string) *model.ErrorResponse {
	req := &model.GetContactsRequest{}
	req.UserID = requestUserID

	contacts, errRes := GetContacts(ctx, req)
	if errRes != nil {
		return errRes
	}

	isAuthorized := false
	for _, contact := range contacts.Users {
		if contact.UserID == resourceUserID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return model.NewErrorResponse("You do not have permission", http.StatusUnauthorized)
	}

	return nil
}

func confirmRoomExist(ctx context.Context, roomID string, opts ...datastore.SelectRoomOption) (*model.Room, *model.ErrorResponse) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID, opts...)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if room == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: fmt.Sprintf("That room is not exist. roomId[%s]", roomID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return room, nil
}

func unsubscribeByRoomID(ctx context.Context, roomID string, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByRoomID(roomID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	<-unsubscribe(ctx, subscriptions)
	if wg != nil {
		wg.Done()
	}
}

// RoomAuthz is room authorize
func RoomAuthz(ctx context.Context, roomID, userID string) *model.ErrorResponse {
	room, errRes := confirmRoomExist(ctx, roomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		return errRes
	}

	if room.Type == scpb.RoomType_RoomTypePublicRoom {
		return nil
	}

	isAuthorized := false
	for _, user := range room.Users {
		if user.UserID == userID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return model.NewErrorResponse("You are not this room member", http.StatusUnauthorized)
	}

	return nil
}

func updateLastAccessRoomID(ctx context.Context, roomID string) {
	userID := ctx.Value(utils.CtxUserID).(string)
	user, _ := confirmUserExist(ctx, userID)
	user.LastAccessRoomID = roomID
	datastore.Provider(ctx).UpdateUser(user)
}

func confirmMessageNotExist(ctx context.Context, messageID string) (*model.Message, *model.ErrorResponse) {
	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if message != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "messageId",
				Reason: fmt.Sprintf("That message already exist. messageId[%s]", messageID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return message, nil
}

func confirmMessageExist(ctx context.Context, messageID string) (*model.Message, *model.ErrorResponse) {
	message, err := datastore.Provider(ctx).SelectMessage(messageID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if message == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "messageId",
				Reason: fmt.Sprintf("That message is not exist. messageId[%s]", messageID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return message, nil
}

func publishUserJoin(ctx context.Context, roomID string) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID, datastore.SelectRoomOptionWithUsers(true))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	go func() {
		userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(roomID)
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
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*model.RoomUser)

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
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*model.RoomUser)
			err := datastore.Provider(ctx).DeleteRoomUsers(ru.RoomID, []string{ru.UserID})
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
		errRes := &model.ErrorResponse{}
		errRes.Status = http.StatusInternalServerError
		errRes.Error = nRes.Error
		return "", errRes
	}
	if nRes.Data == nil {
		return "", nil
	}
	return *nRes.Data.(*string), nil
}

// func getExistUserIDsOld(ctx context.Context, requestUserIDs []string) ([]string, *model.ProblemDetail) {
// 	existUserIDs, err := datastore.Provider(ctx).SelectUserIDsOfUser(requestUserIDs)
// 	if err != nil {
// 		pd := &model.ProblemDetail{
// 			Message: "Getting userIds failed",
// 			Status:  http.StatusInternalServerError,
// 			Error:   err,
// 		}
// 		return nil, pd
// 	}

// 	if len(existUserIDs) != len(requestUserIDs) {
// 		pd := &model.ProblemDetail{
// 			Message: "Invalid params",
// 			InvalidParams: []*model.InvalidParam{
// 				&model.InvalidParam{
// 					Name:   "userIds",
// 					Reason: "It contains a userId that does not exist.",
// 				},
// 			},
// 			Status: http.StatusBadRequest,
// 		}
// 		return nil, pd
// 	}

// 	return existUserIDs, nil
// }

func getExistUserIDs(ctx context.Context, requestUserIDs []string) ([]string, *model.ErrorResponse) {
	existUserIDs, err := datastore.Provider(ctx).SelectUserIDsOfUser(requestUserIDs)
	if err != nil {
		errRes := &model.ErrorResponse{}
		errRes.Status = http.StatusBadRequest
		errRes.Error = err
		return nil, errRes
	}

	if len(existUserIDs) != len(requestUserIDs) {
		errRes := &model.ErrorResponse{}
		errRes.InvalidParams = []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userIds",
				Reason: "It contains a userId that does not exist.",
			},
		}
		errRes.Status = http.StatusBadRequest
		return nil, errRes
	}

	return existUserIDs, nil
}

func confirmRoomUserExist(ctx context.Context, roomID, userID string) (*model.RoomUser, *model.ErrorResponse) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUser(roomID, userID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if roomUser == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId | userId",
				Reason: fmt.Sprintf("That room user is not exist. roomId[%s] userId[%s]", roomID, userID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return roomUser, nil
}

func confirmDeviceNotExist(ctx context.Context, userID string, platform scpb.Platform) (*model.Device, *model.ErrorResponse) {
	device, err := datastore.Provider(ctx).SelectDevice(userID, platform)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if device != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId, platform",
				Reason: fmt.Sprintf("That device already exist. userId[%s] platform[%d]", userID, platform),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return device, nil
}

func confirmDeviceExist(ctx context.Context, userID string, platform scpb.Platform) (*model.Device, *model.ErrorResponse) {
	device, err := datastore.Provider(ctx).SelectDevice(userID, platform)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if device == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId, platform",
				Reason: fmt.Sprintf("That device is not exist. userId[%s] platform[%d]", userID, platform),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}

	return device, nil
}

func confirmAssetExist(ctx context.Context, assetID string) (*model.Asset, *model.ErrorResponse) {
	asset, err := datastore.Provider(ctx).SelectAsset(assetID)
	if err != nil {
		return nil, model.NewErrorResponse("", http.StatusInternalServerError, model.WithError(err))
	}
	if asset == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId, platform",
				Reason: fmt.Sprintf("That asset is not exist. assetId[%s]", assetID),
			},
		}
		return nil, model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
	}
	return asset, nil
}
