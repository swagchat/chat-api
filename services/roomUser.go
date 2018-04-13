package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/rtm"
	"github.com/swagchat/chat-api/utils"
)

// PutRoomUsers is put room users
func PutRoomUsers(ctx context.Context, roomID string, put *models.RequestRoomUserIds) (*models.RoomUsers, *models.ProblemDetail) {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return nil, pd
	}

	put.RemoveDuplicate()

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
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

	userIds, pd := getExistUserIDs(ctx, put.UserIds)
	if pd != nil {
		return nil, pd
	}

	if room.NotificationTopicId == "" {
		notificationTopicID, pd := createTopic(room.RoomId)
		if pd != nil {
			return nil, pd
		}

		room.NotificationTopicId = notificationTopicID
		room.Modified = time.Now().Unix()
		_, err := datastore.Provider(ctx).UpdateRoom(room)
		if err != nil {
			pd := &models.ProblemDetail{
				Status: http.StatusInternalServerError,
				Title:  "Get user information failed",
				Error:  err,
			}
			return nil, pd
		}
	}

	var zero int64
	zero = 0
	roomUsers := make([]*models.RoomUser, 0)
	nowTimestamp := time.Now().Unix()
	for _, userID := range userIds {
		roomUsers = append(roomUsers, &models.RoomUser{
			RoomId:      roomID,
			UserId:      userID,
			UnreadCount: &zero,
			MetaData:    []byte("{}"),
			Created:     nowTimestamp,
			Modified:    nowTimestamp,
		})
	}
	err = datastore.Provider(ctx).InsertRoomUsers(roomUsers)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's user list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	roomUsers, err = datastore.Provider(ctx).SelectRoomUsersByRoomID(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's user list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	returnRoomUsers := &models.RoomUsers{
		RoomUsers: roomUsers,
	}

	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, roomID)

	return returnRoomUsers, nil
}

// PutRoomUser is put room user
func PutRoomUser(ctx context.Context, put *models.RoomUser) (*models.RoomUser, *models.ProblemDetail) {
	roomUser, pd := selectRoomUser(ctx, put.RoomId, put.UserId)
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
		pd := &models.ProblemDetail{
			Title:  "Room's user registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	return roomUser, nil
}

// DeleteRoomUsers is delete room users
func DeleteRoomUsers(ctx context.Context, roomID string, deleteUserIds *models.RequestRoomUserIds) (*models.RoomUsers, *models.ProblemDetail) {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return nil, pd
	}

	deleteUserIds.RemoveDuplicate()

	if pd := deleteUserIds.IsValid("DELETE", room); pd != nil {
		return nil, pd
	}

	userIds, pd := getExistUserIDs(ctx, deleteUserIds.UserIds)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider(ctx).DeleteRoomUser(roomID, userIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsersByRoomIDAndUserIDs(&roomID, userIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	go unsubscribeByRoomUsers(ctx, roomUsers)

	roomUsers, err = datastore.Provider(ctx).SelectRoomUsersByRoomID(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.RoomUsers{
		RoomUsers: roomUsers,
	}, nil
}

func selectRoomUser(ctx context.Context, roomID, userID string) (*models.RoomUser, *models.ProblemDetail) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUser(roomID, userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if roomUser == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}
	return roomUser, nil
}

func publishUserJoin(ctx context.Context, roomID string) {
	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
		return
	}

	b, _ := json.Marshal(userForRooms)
	buf := new(bytes.Buffer)
	buf.Write(b)

	message := &models.Message{
		RoomId:    roomID,
		EventName: "userJoin",
		Payload:   utils.JSONText(buf.String()),
	}
	bytes, err := json.Marshal(message)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	mi := &rtm.MessagingInfo{
		Message: string(bytes),
	}
	err = rtm.Provider().PublishMessage(mi)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
}

func subscribeByRoomUsers(ctx context.Context, roomUsers []*models.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*models.RoomUser)

			devices, err := datastore.Provider(ctx).SelectDevicesByUserID(ru.UserId)
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Subscribe failed",
					Status: http.StatusInternalServerError,
					Error:  err,
				}
				pdChan <- pd
			}
			if devices != nil {
				for _, d := range devices {
					if d.Token != "" {
						if d.NotificationDeviceId == "" {
							nRes := <-notification.Provider().CreateEndpoint("", d.Platform, d.Token)
							if nRes.ProblemDetail != nil {
								pdChan <- nRes.ProblemDetail
							} else {
								d.NotificationDeviceId = *nRes.Data.(*string)
								err := datastore.Provider(ctx).UpdateDevice(d)
								if err != nil {
									pd := &models.ProblemDetail{
										Title:  "Subscribe failed",
										Status: http.StatusInternalServerError,
										Error:  err,
									}
									pdChan <- pd
								}
							}
						}
						go subscribe(ctx, []*models.RoomUser{ru}, d)
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
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
					Error:         pd.Error,
				})
				return
			}
		})
	}
	d.Wait()
	return
}

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*models.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*models.RoomUser)
			err := datastore.Provider(ctx).DeleteRoomUser(ru.RoomId, []string{ru.UserId})
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Delete room's user failed",
					Status: http.StatusInternalServerError,
				}
				pdChan <- pd
			}

			devices, err := datastore.Provider(ctx).SelectDevicesByUserID(ru.UserId)
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Subscribe failed",
					Status: http.StatusInternalServerError,
					Error:  err,
				}
				pdChan <- pd
			}
			if devices != nil {
				for _, d := range devices {
					subscription, err := datastore.Provider(ctx).SelectSubscription(ru.RoomId, ru.UserId, d.Platform)
					if err != nil {
						pd := &models.ProblemDetail{
							Title:  "User registration failed",
							Status: http.StatusInternalServerError,
							Error:  err,
						}
						pdChan <- pd
					}
					go unsubscribe(ctx, []*models.Subscription{subscription})
				}
			}
			doneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-doneChan:
				return
			case pd := <-pdChan:
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
					Error:         pd.Error,
				})
				return
			}
		})
	}
	d.Wait()
	return
}

func createTopic(roomID string) (string, *models.ProblemDetail) {
	nRes := <-notification.Provider().CreateTopic(roomID)
	if nRes.ProblemDetail != nil {
		return "", nRes.ProblemDetail
	}
	if nRes.Data == nil {
		return "", nil
	}
	return *nRes.Data.(*string), nil
}

func getExistUserIDs(ctx context.Context, requestUserIDs []string) ([]string, *models.ProblemDetail) {
	existUserIDs, err := datastore.Provider(ctx).SelectUserIDsByUserIDs(requestUserIDs)
	if err != nil {
		pd := &models.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Getting userIds failed",
			Error:  err,
		}
		return nil, pd
	}

	if len(existUserIDs) != len(requestUserIDs) {
		pd := &models.ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userIds",
					Reason: "It contains a userId that does not exist.",
				},
			},
		}
		return nil, pd
	}

	return existUserIDs, nil
}
