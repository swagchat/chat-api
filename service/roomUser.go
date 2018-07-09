package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/pbroker"
	"github.com/swagchat/chat-api/protobuf"
	"github.com/swagchat/chat-api/utils"
)

// PutRoomUsers is put room users
func PutRoomUsers(ctx context.Context, put *protobuf.PostRoomUserReq) (*protobuf.RoomUsers, *model.ProblemDetail) {
	room, pd := selectRoom(ctx, put.RoomID)
	if pd != nil {
		return nil, pd
	}

	put.RemoveDuplicate()

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(put.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Get users failed",
			Error:  err,
		}
		return nil, pd
	}
	room.Users = userForRooms

	if pd := put.IsValid("PUT", room); pd != nil {
		return nil, pd
	}

	userIds, pd := getExistUserIDs(ctx, put.UserIDs)
	if pd != nil {
		return nil, pd
	}

	if room.NotificationTopicID == "" {
		notificationTopicID, pd := createTopic(room.RoomID)
		if pd != nil {
			return nil, pd
		}

		room.NotificationTopicID = notificationTopicID
		room.Modified = time.Now().Unix()
		_, err := datastore.Provider(ctx).UpdateRoom(room)
		if err != nil {
			pd := &model.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Get user information failed",
				Error:  err,
			}
			return nil, pd
		}
	}

	// var zero int
	// zero = 0
	roomUsers := make([]*protobuf.RoomUser, 0)
	for _, userID := range userIds {
		roomUsers = append(roomUsers, &protobuf.RoomUser{
			RoomID: put.RoomID,
			UserID: userID,
			// UnreadCount: &zero,
			UnreadCount: 0,
			Display:     put.Display,
		})
	}
	err = datastore.Provider(ctx).InsertRoomUsers(roomUsers)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Get room's user list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	roomUsers, err = datastore.Provider(ctx).SelectRoomUsersByRoomID(put.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Get room's user list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	returnRoomUsers := &protobuf.RoomUsers{
		RoomUsers: roomUsers,
	}

	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, put.RoomID)

	return returnRoomUsers, nil
}

// PutRoomUser is put room user
func PutRoomUser(ctx context.Context, put *protobuf.RoomUser) (*protobuf.RoomUser, *model.ProblemDetail) {
	roomUser, pd := selectRoomUser(ctx, put.RoomID, put.UserID)
	if pd != nil {
		return nil, pd
	}

	roomUser.Put(put)
	if pd := roomUser.IsValid(); pd != nil {
		return nil, pd
	}
	roomUser.BeforeSave()

	roomUser, err := datastore.Provider(ctx).UpdateRoomUser(roomUser)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Room's user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	// var p json.RawMessage
	// err = json.Unmarshal([]byte("{}"), &p)
	// m := &model.Message{
	// 	RoomID:    roomUser.RoomID,
	// 	UserID:    roomUser.UserID,
	// 	Type:      model.MessageTypeUpdateRoomUser,
	// 	EventName: "message",
	// 	Payload:   p,
	// }
	// rtmPublish(ctx, m)

	return roomUser, nil
}

// DeleteRoomUsers is delete room users
func DeleteRoomUsers(ctx context.Context, req *protobuf.DeleteRoomUserReq) (*protobuf.RoomUsers, *model.ProblemDetail) {
	_, pd := selectRoom(ctx, req.RoomID)
	if pd != nil {
		return nil, pd
	}

	req.RemoveDuplicate()

	userIds, pd := getExistUserIDs(ctx, req.UserIDs)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider(ctx).DeleteRoomUser(req.RoomID, userIds)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Delete room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsersByRoomIDAndUserIDs(&req.RoomID, userIds)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Delete room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	go unsubscribeByRoomUsers(ctx, roomUsers)

	roomUsers, err = datastore.Provider(ctx).SelectRoomUsersByRoomID(req.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Get room's users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &protobuf.RoomUsers{
		RoomUsers: roomUsers,
	}, nil
}

func selectRoomUser(ctx context.Context, roomID, userID string) (*protobuf.RoomUser, *model.ProblemDetail) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUser(roomID, userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Title:  "Get room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if roomUser == nil {
		return nil, &model.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}
	return roomUser, nil
}

func publishUserJoin(ctx context.Context, roomID string) {
	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
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
		json.NewEncoder(buffer).Encode(userForRooms)
		rtmEvent := &pbroker.RTMEvent{
			Type:    pbroker.UserJoin,
			Payload: buffer.Bytes(),
			UserIDs: userIDs,
		}
		err = pbroker.Provider().PublishMessage(rtmEvent)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}()
}

func subscribeByRoomUsers(ctx context.Context, roomUsers []*protobuf.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *model.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*protobuf.RoomUser)

			devices, err := datastore.Provider(ctx).SelectDevicesByUserID(ru.UserID)
			if err != nil {
				pd := &model.ProblemDetail{
					Title:  "Subscribe failed",
					Status: http.StatusInternalServerError,
					Error:  err,
				}
				pdChan <- pd
			}
			if devices != nil {
				for _, d := range devices {
					if d.Token != "" {
						if d.NotificationDeviceID == "" {
							nRes := <-notification.Provider().CreateEndpoint("", d.Platform, d.Token)
							if nRes.ProblemDetail != nil {
								pdChan <- nRes.ProblemDetail
							} else {
								d.NotificationDeviceID = *nRes.Data.(*string)
								err := datastore.Provider(ctx).UpdateDevice(d)
								if err != nil {
									pd := &model.ProblemDetail{
										Title:  "Subscribe failed",
										Status: http.StatusInternalServerError,
										Error:  err,
									}
									pdChan <- pd
								}
							}
						}
						go subscribe(ctx, []*protobuf.RoomUser{ru}, d)
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
				logger.Error(pd.Error.Error())
				return
			}
		})
	}
	d.Wait()
	return
}

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*protobuf.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *model.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*protobuf.RoomUser)
			err := datastore.Provider(ctx).DeleteRoomUser(ru.RoomID, []string{ru.UserID})
			if err != nil {
				pd := &model.ProblemDetail{
					Title:  "Delete room's user failed",
					Status: http.StatusInternalServerError,
				}
				pdChan <- pd
			}

			devices, err := datastore.Provider(ctx).SelectDevicesByUserID(ru.UserID)
			if err != nil {
				pd := &model.ProblemDetail{
					Title:  "Subscribe failed",
					Status: http.StatusInternalServerError,
					Error:  err,
				}
				pdChan <- pd
			}
			if devices != nil {
				for _, d := range devices {
					subscription, err := datastore.Provider(ctx).SelectSubscription(ru.RoomID, ru.UserID, d.Platform)
					if err != nil {
						pd := &model.ProblemDetail{
							Title:  "User registration failed",
							Status: http.StatusInternalServerError,
							Error:  err,
						}
						pdChan <- pd
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
			case pd := <-pdChan:
				logger.Error(pd.Error.Error())
				return
			}
		})
	}
	d.Wait()
	return
}

func createTopic(roomID string) (string, *model.ProblemDetail) {
	nRes := <-notification.Provider().CreateTopic(roomID)
	if nRes.ProblemDetail != nil {
		return "", nRes.ProblemDetail
	}
	if nRes.Data == nil {
		return "", nil
	}
	return *nRes.Data.(*string), nil
}

func getExistUserIDs(ctx context.Context, requestUserIDs []string) ([]string, *model.ProblemDetail) {
	existUserIDs, err := datastore.Provider(ctx).SelectUserIDsByUserIDs(requestUserIDs)
	if err != nil {
		pd := &model.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Getting userIds failed",
			Error:  err,
		}
		return nil, pd
	}

	if len(existUserIDs) != len(requestUserIDs) {
		pd := &model.ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []model.InvalidParam{
				model.InvalidParam{
					Name:   "userIds",
					Reason: "It contains a userId that does not exist.",
				},
			},
		}
		return nil, pd
	}

	return existUserIDs, nil
}

func selectUserIDsOfRoomUser(ctx context.Context, in *protobuf.GetUserIDsOfRoomUserReq) (*protobuf.UserIDs, error) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(in.RoomID, datastore.WithRoleIDs(in.RoleIDs))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &protobuf.UserIDs{
		UserIDs: userIDs,
	}, nil
}
