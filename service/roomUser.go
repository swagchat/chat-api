package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/pbroker"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

// CreateRoomUsers create room users
func CreateRoomUsers(ctx context.Context, req *scpb.CreateRoomUsersRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start CreateRoomUser. CreateRoomUserRequest[%#v]", req))

	urs := &model.RoomUsers{}
	urs.ImportFromPbCreateUserRolesRequest(req)

	room, pd := selectRoom(ctx, urs.RoomID)
	if pd != nil {
		return pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(req.RoomId)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create room users.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}
	urs.Room.Users = userForRooms

	pd = urs.Validate()
	if pd != nil {
		return pd
	}

	userIDs, pd := getExistUserIDs(ctx, urs.UserIDs)
	if pd != nil {
		return pd
	}
	urs.UserIDs = userIDs

	if room.NotificationTopicID == "" {
		notificationTopicID, pd := createTopic(room.RoomID)
		if pd != nil {
			return pd
		}

		room.NotificationTopicID = notificationTopicID
		room.Modified = time.Now().Unix()
		_, err := datastore.Provider(ctx).UpdateRoom(room)
		if err != nil {
			pd := &model.ProblemDetail{
				Message: "Failed to create room users.",
				Status:  http.StatusInternalServerError,
				Error:   err,
			}
			return pd
		}
	}

	urs.GenerateRoomUsers()
	err = datastore.Provider(ctx).InsertRoomUsers(urs.RoomUsers)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create room users.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	go subscribeByRoomUsers(ctx, urs.RoomUsers)
	go publishUserJoin(ctx, req.RoomId)

	logger.Info("Finish CreateRoomUser")
	return nil
}

// UpdateRoomUser update room user
func UpdateRoomUser(ctx context.Context, req *scpb.UpdateRoomUserRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start UpdateRoomUser. UpdateRoomUserRequest[%#v]", req))

	ur := &model.RoomUser{}
	ur.ImportFromPbUpdateUserRoleRequest(req)

	_, pd := selectRoomUser(ctx, ur.RoomID, ur.UserID)
	if pd != nil {
		return pd
	}

	pd = ur.Validate()
	if pd != nil {
		return pd
	}

	_, err := datastore.Provider(ctx).UpdateRoomUser(ur)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Member registration failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
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

	logger.Info("Finish UpdateRoomUser")
	return nil
}

// DeleteRoomUsers delete room users
func DeleteRoomUsers(ctx context.Context, req *scpb.DeleteRoomUserRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start DeleteRoomUsers. DeleteRoomUserRequest[%#v]", req))

	_, pd := selectRoom(ctx, req.RoomId)
	if pd != nil {
		return pd
	}

	// req.RemoveDuplicate()

	userIds, pd := getExistUserIDs(ctx, req.UserIds)
	if pd != nil {
		return pd
	}

	err := datastore.Provider(ctx).DeleteRoomUser(req.RoomId, userIds)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Delete room's user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsersByRoomIDAndUserIDs(&req.RoomId, userIds)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Delete room's user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	go unsubscribeByRoomUsers(ctx, roomUsers)

	// roomUsers, err = datastore.Provider(ctx).SelectRoomUsersByRoomID(req.RoomId)
	// if err != nil {
	// 	pd := &model.ProblemDetail{
	// 		Message: "Get room's users failed",
	// 		Status:  http.StatusInternalServerError,
	// 		Error:   err,
	// 	}
	// 	return pd
	// }

	logger.Info(fmt.Sprintf("Finish DeleteRoomUsers"))
	return nil
}

func selectRoomUser(ctx context.Context, roomID, userID string) (*model.RoomUser, *model.ProblemDetail) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUser(roomID, userID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room's user failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if roomUser == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
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

func subscribeByRoomUsers(ctx context.Context, roomUsers []*model.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *model.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*model.RoomUser)

			devices, err := datastore.Provider(ctx).SelectDevicesByUserID(ru.UserID)
			if err != nil {
				pd := &model.ProblemDetail{
					Message: "Subscribe failed",
					Status:  http.StatusInternalServerError,
					Error:   err,
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
										Message: "Subscribe failed",
										Status:  http.StatusInternalServerError,
										Error:   err,
									}
									pdChan <- pd
								}
							}
						}
						go subscribe(ctx, []*model.RoomUser{ru}, d)
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

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*model.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *model.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*model.RoomUser)
			err := datastore.Provider(ctx).DeleteRoomUser(ru.RoomID, []string{ru.UserID})
			if err != nil {
				pd := &model.ProblemDetail{
					Message: "Delete room's user failed",
					Status:  http.StatusInternalServerError,
				}
				pdChan <- pd
			}

			devices, err := datastore.Provider(ctx).SelectDevicesByUserID(ru.UserID)
			if err != nil {
				pd := &model.ProblemDetail{
					Message: "Subscribe failed",
					Status:  http.StatusInternalServerError,
					Error:   err,
				}
				pdChan <- pd
			}
			if devices != nil {
				for _, d := range devices {
					subscription, err := datastore.Provider(ctx).SelectSubscription(ru.RoomID, ru.UserID, d.Platform)
					if err != nil {
						pd := &model.ProblemDetail{
							Message: "User registration failed",
							Status:  http.StatusInternalServerError,
							Error:   err,
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
			Message: "Getting userIds failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	if len(existUserIDs) != len(requestUserIDs) {
		pd := &model.ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*model.InvalidParam{
				&model.InvalidParam{
					Name:   "userIds",
					Reason: "It contains a userId that does not exist.",
				},
			},
			Status: http.StatusBadRequest,
		}
		return nil, pd
	}

	return existUserIDs, nil
}

func SelectUserIDsOfRoomUser(ctx context.Context, req *scpb.GetUserIdsOfRoomUserRequest) (*scpb.UserIds, *model.ProblemDetail) {
	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(req.RoomId, datastore.WithRoleIDs(req.RoleIds))
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to getting userIds.",
			Status:  http.StatusInternalServerError,
		}
		return nil, pd
	}

	return &scpb.UserIds{
		UserIds: userIDs,
	}, nil
}
