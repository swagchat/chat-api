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

func PutRoomUsers(roomId string, put *models.RequestRoomUserIds, dsCfg *utils.Datastore) (*models.RoomUsers, *models.ProblemDetail) {
	room, pd := selectRoom(roomId, dsCfg)
	if pd != nil {
		return nil, pd
	}

	put.RemoveDuplicate()

	userForRooms, err := datastore.Provider(dsCfg).SelectUsersForRoom(roomId)
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

	userIds, pd := getExistUserIds(put.UserIds, dsCfg)
	if pd != nil {
		return nil, pd
	}

	if room.NotificationTopicId == "" {
		notificationTopicId, pd := createTopic(room.RoomId)
		if pd != nil {
			return nil, pd
		}

		room.NotificationTopicId = notificationTopicId
		room.Modified = time.Now().Unix()
		_, err := datastore.Provider(dsCfg).UpdateRoom(room)
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
	for _, userId := range userIds {
		roomUsers = append(roomUsers, &models.RoomUser{
			RoomId:      roomId,
			UserId:      userId,
			UnreadCount: &zero,
			MetaData:    []byte("{}"),
			Created:     nowTimestamp,
			Modified:    nowTimestamp,
		})
	}
	err = datastore.Provider(dsCfg).InsertRoomUsers(roomUsers)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's user list failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	roomUsers, err = datastore.Provider(dsCfg).SelectRoomUsersByRoomId(roomId)
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

	ctx, _ := context.WithCancel(context.Background())
	go subscribeByRoomUsers(ctx, roomUsers, dsCfg)
	go publishUserJoin(roomId, dsCfg)

	return returnRoomUsers, nil
}

func PutRoomUser(put *models.RoomUser, dsCfg *utils.Datastore) (*models.RoomUser, *models.ProblemDetail) {
	roomUser, pd := selectRoomUser(put.RoomId, put.UserId, dsCfg)
	if pd != nil {
		return nil, pd
	}

	roomUser.Put(put)
	if pd := roomUser.IsValid(); pd != nil {
		return nil, pd
	}
	roomUser.BeforeSave()

	roomUser, err := datastore.Provider(dsCfg).UpdateRoomUser(roomUser)
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

func DeleteRoomUsers(roomId string, deleteUserIds *models.RequestRoomUserIds, dsCfg *utils.Datastore) (*models.RoomUsers, *models.ProblemDetail) {
	room, pd := selectRoom(roomId, dsCfg)
	if pd != nil {
		return nil, pd
	}

	deleteUserIds.RemoveDuplicate()

	if pd := deleteUserIds.IsValid("DELETE", room); pd != nil {
		return nil, pd
	}

	userIds, pd := getExistUserIds(deleteUserIds.UserIds, dsCfg)
	if pd != nil {
		return nil, pd
	}

	err := datastore.Provider(dsCfg).DeleteRoomUser(roomId, userIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	roomUsers, err := datastore.Provider(dsCfg).SelectRoomUsersByRoomIdAndUserIds(&roomId, userIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete room's user failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	ctx, _ := context.WithCancel(context.Background())
	go unsubscribeByRoomUsers(ctx, roomUsers, dsCfg)

	roomUsers, err = datastore.Provider(dsCfg).SelectRoomUsersByRoomId(roomId)
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

func selectRoomUser(roomId, userId string, dsCfg *utils.Datastore) (*models.RoomUser, *models.ProblemDetail) {
	roomUser, err := datastore.Provider(dsCfg).SelectRoomUser(roomId, userId)
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

func publishUserJoin(roomId string, dsCfg *utils.Datastore) {
	userForRooms, err := datastore.Provider(dsCfg).SelectUsersForRoom(roomId)
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
		RoomId:    roomId,
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

func subscribeByRoomUsers(ctx context.Context, roomUsers []*models.RoomUser, dsCfg *utils.Datastore) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)

			devices, err := datastore.Provider(dsCfg).SelectDevicesByUserId(ru.UserId)
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
								err := datastore.Provider(dsCfg).UpdateDevice(d)
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
						go subscribe(ctx, []*models.RoomUser{ru}, d, dsCfg)
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

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*models.RoomUser, dsCfg *utils.Datastore) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)
			err := datastore.Provider(dsCfg).DeleteRoomUser(ru.RoomId, []string{ru.UserId})
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Delete room's user failed",
					Status: http.StatusInternalServerError,
				}
				pdChan <- pd
			}

			devices, err := datastore.Provider(dsCfg).SelectDevicesByUserId(ru.UserId)
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
					subscription, err := datastore.Provider(dsCfg).SelectSubscription(ru.RoomId, ru.UserId, d.Platform)
					if err != nil {
						pd := &models.ProblemDetail{
							Title:  "User registration failed",
							Status: http.StatusInternalServerError,
							Error:  err,
						}
						pdChan <- pd
					}
					go unsubscribe(ctx, []*models.Subscription{subscription}, dsCfg)
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

func createTopic(roomId string) (string, *models.ProblemDetail) {
	nRes := <-notification.Provider().CreateTopic(roomId)
	if nRes.ProblemDetail != nil {
		return "", nRes.ProblemDetail
	}
	if nRes.Data == nil {
		return "", nil
	} else {
		return *nRes.Data.(*string), nil
	}
}

func getExistUserIds(requestUserIds []string, dsCfg *utils.Datastore) ([]string, *models.ProblemDetail) {
	existUserIds, err := datastore.Provider(dsCfg).SelectUserIdsByUserIds(requestUserIds)
	if err != nil {
		pd := &models.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Getting userIds failed",
			Error:  err,
		}
		return nil, pd
	}

	if len(existUserIds) != len(requestUserIds) {
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

	return existUserIds, nil
}
