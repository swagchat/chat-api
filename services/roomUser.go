package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/rtm"
	"github.com/fairway-corp/swagchat-api/utils"
)

func PutRoomUsers(roomId string, put *models.RequestRoomUserIds) (*models.RoomUsers, *models.ProblemDetail) {
	room, pd := selectRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	put.RemoveDuplicate()

	dRes := datastore.GetProvider().SelectUsersForRoom(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	room.Users = dRes.Data.([]*models.UserForRoom)

	if pd := put.IsValid("PUT", room); pd != nil {
		return nil, pd
	}

	userIds, pd := getExistUserIds(put.UserIds)
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
		dRes := datastore.GetProvider().UpdateRoom(room)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
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
	dRes = datastore.GetProvider().InsertRoomUsers(roomUsers)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	dRes = datastore.GetProvider().SelectRoomUsersByRoomId(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	returnRoomUsers := &models.RoomUsers{
		RoomUsers: dRes.Data.([]*models.RoomUser),
	}

	ctx, _ := context.WithCancel(context.Background())
	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(roomId)

	return returnRoomUsers, nil
}

func PutRoomUser(put *models.RoomUser) (*models.RoomUser, *models.ProblemDetail) {
	roomUser, pd := selectRoomUser(put.RoomId, put.UserId)
	if pd != nil {
		return nil, pd
	}

	roomUser.Put(put)
	if pd := roomUser.IsValid(); pd != nil {
		return nil, pd
	}
	roomUser.BeforeSave()

	dRes := datastore.GetProvider().UpdateRoomUser(roomUser)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	return dRes.Data.(*models.RoomUser), nil
}

func DeleteRoomUsers(roomId string, deleteUserIds *models.RequestRoomUserIds) (*models.RoomUsers, *models.ProblemDetail) {
	room, pd := selectRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	deleteUserIds.RemoveDuplicate()

	if pd := deleteUserIds.IsValid("DELETE", room); pd != nil {
		return nil, pd
	}

	userIds, pd := getExistUserIds(deleteUserIds.UserIds)
	if pd != nil {
		return nil, pd
	}

	dRes := datastore.GetProvider().DeleteRoomUser(roomId, userIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	dRes = datastore.GetProvider().SelectRoomUsersByRoomIdAndUserIds(&roomId, userIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	ctx, _ := context.WithCancel(context.Background())
	go unsubscribeByRoomUsers(ctx, dRes.Data.([]*models.RoomUser))

	dRes = datastore.GetProvider().SelectRoomUsersByRoomId(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	returnRoomUsers := &models.RoomUsers{
		RoomUsers: dRes.Data.([]*models.RoomUser),
	}

	return returnRoomUsers, nil
}

func selectRoomUser(roomId, userId string) (*models.RoomUser, *models.ProblemDetail) {
	dRes := datastore.GetProvider().SelectRoomUser(roomId, userId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	if dRes.Data == nil {
		return nil, &models.ProblemDetail{
			Status: http.StatusNotFound,
		}
	}
	return dRes.Data.(*models.RoomUser), nil
}

func publishUserJoin(roomId string) {
	dRes := datastore.GetProvider().SelectUsersForRoom(roomId)
	if dRes.ProblemDetail != nil {
		problemDetailBytes, _ := json.Marshal(dRes.ProblemDetail)
		utils.AppLogger.Error("",
			zap.String("msg", "Publish error. (Add room's user list)"),
			zap.String("problemDetail", string(problemDetailBytes)),
		)
		return
	}
	users := dRes.Data.([]*models.UserForRoom)
	b, _ := json.Marshal(users)
	buf := new(bytes.Buffer)
	buf.Write(b)

	message := &models.Message{
		RoomId:    roomId,
		EventName: "userJoin",
		Payload:   utils.JSONText(buf.String()),
	}
	bytes, err := json.Marshal(message)
	if err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", err.Error()),
		)
	}
	mi := &rtm.MessagingInfo{
		Message: string(bytes),
	}
	err = rtm.GetMessagingProvider().PublishMessage(mi)
	if err != nil {
		utils.AppLogger.Error("",
			zap.String("msg", "Publish error. (Add room's user list)"),
			zap.String("detail", err.Error()),
		)
	}
}

func subscribeByRoomUsers(ctx context.Context, roomUsers []*models.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)

			dRes := datastore.GetProvider().SelectDevicesByUserId(ru.UserId)
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			}
			if dRes.Data != nil {
				devices := dRes.Data.([]*models.Device)
				for _, d := range devices {
					if d.Token != "" {
						if d.NotificationDeviceId == "" {
							nRes := <-notification.GetProvider().CreateEndpoint("", d.Platform, d.Token)
							if nRes.ProblemDetail != nil {
								pdChan <- dRes.ProblemDetail
							} else {
								d.NotificationDeviceId = *nRes.Data.(*string)
								dRes := datastore.GetProvider().UpdateDevice(d)
								if dRes.ProblemDetail != nil {
									pdChan <- dRes.ProblemDetail
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
	return
}

func unsubscribeByRoomUsers(ctx context.Context, roomUsers []*models.RoomUser) {
	doneChan := make(chan bool, 1)
	pdChan := make(chan *models.ProblemDetail, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value("roomUser").(*models.RoomUser)
			dRes := datastore.GetProvider().DeleteRoomUser(ru.RoomId, []string{ru.UserId})
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			}

			dRes = datastore.GetProvider().SelectDevicesByUserId(ru.UserId)
			if dRes.ProblemDetail != nil {
				pdChan <- dRes.ProblemDetail
			}
			if dRes.Data != nil {
				devices := dRes.Data.([]*models.Device)
				for _, d := range devices {
					dRes = datastore.GetProvider().SelectSubscription(ru.RoomId, ru.UserId, d.Platform)
					if dRes.ProblemDetail != nil {
						pdChan <- dRes.ProblemDetail
					}
					go unsubscribe(ctx, dRes.Data.([]*models.Subscription))
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
	return
}

func createTopic(roomId string) (string, *models.ProblemDetail) {
	nRes := <-notification.GetProvider().CreateTopic(roomId)
	if nRes.ProblemDetail != nil {
		return "", nRes.ProblemDetail
	}
	if nRes.Data == nil {
		return "", nil
	} else {
		return *nRes.Data.(*string), nil
	}
}

func getExistUserIds(requestUserIds []string) ([]string, *models.ProblemDetail) {
	dRes := datastore.GetProvider().SelectUserIdsByUserIds(requestUserIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	existUserIds := dRes.Data.([]string)
	if len(existUserIds) != len(requestUserIds) {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userIds",
					Reason: "It contains a userId that does not exist.",
				},
			},
		}
	}

	return existUserIds, nil
}
