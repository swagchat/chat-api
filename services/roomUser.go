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
	"github.com/fairway-corp/swagchat-api/messaging"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func PostRoomUsers(roomId string, rus *models.RoomUsers) (*models.ResponseRoomUser, *models.ProblemDetail) {
	room, pd := selectRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	if len(rus.Users) == 0 {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "users",
					Reason: "Not set.",
				},
			},
		}
	}

	requestUserIds := utils.RemoveDuplicate(rus.Users)

	dRes := <-datastore.GetProvider().UserSelectByUserIds(requestUserIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	existUsers := dRes.Data.([]*models.User)
	if len(existUsers) != len(requestUserIds) {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "users",
					Reason: "It contains a userId that does not exist.",
				},
			},
		}
	}

	if room.NotificationTopicId == "" {
		pd = createTopic(room)
		if pd != nil {
			return nil, pd
		}
	}
	/*
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		SubscribeByRoomUsers(ctx, existUsers)

		dRes = <-datastore.GetProvider().RoomUsersSelectUserIds(roomId)
		if dRes.ProblemDetail == nil && dRes.Data != nil {
			currentUserIds := dRes.Data.([]string)
			deleteUserIds := make([]string, 0)
			var isHit bool
			for _, currentUserId := range currentUserIds {
				isHit = false
				for _, requestUserId := range requestUserIds {
					if currentUserId == requestUserId {
						isHit = true
					}
				}
				if !isHit {
					deleteUserIds = append(deleteUserIds, currentUserId)
				}
			}
			if len(deleteUserIds) > 0 {
				dRes := <-datastore.GetProvider().RoomUsersSelect(&roomId, deleteUserIds)
				if dRes.ProblemDetail != nil {
					return nil, dRes.ProblemDetail
				}
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				delRuRes := UnsubscribeByRoomUsers(ctx, dRes.Data.([]*models.RoomUser))
				ruRes.Errors = append(ruRes.Errors, delRuRes.Errors...)
			}
		}
	*/

	return nil, nil
}

func PutRoomUsers(roomId string, rus *models.RoomUsers) ([]*models.RoomUser, *models.ProblemDetail) {
	room, pd := selectRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	if pd := rus.IsValid(); pd != nil {
		return nil, pd
	}

	rus.RemoveDuplicate()

	users, pd := getUsersByUserIds(rus.Users)
	if pd != nil {
		return nil, pd
	}

	if room.NotificationTopicId == "" {
		pd = createTopic(room)
		if pd != nil {
			return nil, pd
		}
	}

	var zero int64
	zero = 0
	roomUsers := make([]*models.RoomUser, 0)
	for _, u := range users {
		roomUsers = append(roomUsers, &models.RoomUser{
			RoomId:      roomId,
			UserId:      u.UserId,
			UnreadCount: &zero,
			MetaData:    []byte("{}"),
			Created:     time.Now().UnixNano(),
		})
	}
	dRes := <-datastore.GetProvider().RoomUsersInsert(roomUsers, false)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	dRes = <-datastore.GetProvider().RoomUsersSelectIds(&roomId, nil)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	roomUsers = dRes.Data.([]*models.RoomUser)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go subscribeByRoomUsers(ctx, roomUsers)

	go publishUserJoin(roomId)
	return roomUsers, nil
}

func DeleteRoomUsers(roomId string, rus *models.RoomUsers) (*models.ResponseRoomUser, *models.ProblemDetail) {
	// Room existence check
	_, pd := selectRoom(roomId)
	if pd != nil {
		return nil, pd
	}

	deleteUserIds := utils.RemoveDuplicate(rus.Users)
	dRes := <-datastore.GetProvider().RoomUsersSelectIds(&roomId, deleteUserIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go unsubscribeByRoomUsers(ctx, dRes.Data.([]*models.RoomUser))

	return nil, nil
}

func PutRoomUser(roomId, userId string, requestRoomUser *models.RoomUser) *models.ProblemDetail {
	dRes := <-datastore.GetProvider().RoomUserSelect(roomId, userId)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	roomUser := dRes.Data.(*models.RoomUser)
	if roomUser == nil {
		return &models.ProblemDetail{
			Title:     "Request parameter error. (Update room's user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId and userId",
					Reason: "room's user item is not exist.",
				},
			},
		}
	}
	if requestRoomUser.UnreadCount != nil {
		roomUser.UnreadCount = requestRoomUser.UnreadCount
	}
	if requestRoomUser.MetaData != nil {
		roomUser.MetaData = requestRoomUser.MetaData
	}

	dRes = <-datastore.GetProvider().RoomUserUpdate(roomUser)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	if requestRoomUser.UnreadCount != nil {
		dRes = <-datastore.GetProvider().UserUnreadCountRecalc(userId)
		return dRes.ProblemDetail
	}

	return nil
}

func publishUserJoin(roomId string) {
	dRes := <-datastore.GetProvider().RoomSelectUsersForRoom(roomId)
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

	if utils.Cfg.RealtimeServer.Endpoint != "" {
		message := &models.Message{
			RoomId:    roomId,
			EventName: "userJoin",
			Payload:   utils.JSONText(buf.String()),
		}
		bytes, _ := json.Marshal(message)
		messagingInfo := &messaging.MessagingInfo{
			Message: string(bytes),
		}
		messagingProvider := messaging.GetMessagingProvider()
		err := messagingProvider.PublishMessage(messagingInfo)
		if err != nil {
			utils.AppLogger.Error("",
				zap.String("msg", "Publish error. (Add room's user list)"),
				zap.String("detail", err.Error()),
			)
		}
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

			dRes := <-datastore.GetProvider().DeviceSelectByUserId(ru.UserId)
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
								dRes := <-datastore.GetProvider().DeviceUpdate(d)
								if dRes.ProblemDetail != nil {
									pdChan <- dRes.ProblemDetail
								}
							}
						}
					}
					go subscribe(ctx, []*models.RoomUser{ru}, d)
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
			targetRoomUser := ctx.Value("roomUser").(*models.RoomUser)
			if targetRoomUser.NotificationSubscribeId != nil {
				nRes := <-notification.GetProvider().Unsubscribe(*targetRoomUser.NotificationSubscribeId)
				if nRes.ProblemDetail != nil {
					pdChan <- nRes.ProblemDetail
				} else {
					dRes := <-datastore.GetProvider().RoomUserDelete(targetRoomUser.RoomId, []string{targetRoomUser.UserId})
					if dRes.ProblemDetail != nil {
						pdChan <- dRes.ProblemDetail
					}
				}
			} else {
				dRes := <-datastore.GetProvider().RoomUserDelete(targetRoomUser.RoomId, []string{targetRoomUser.UserId})
				if dRes.ProblemDetail != nil {
					pdChan <- dRes.ProblemDetail
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

func createTopic(room *models.Room) *models.ProblemDetail {
	nRes := <-notification.GetProvider().CreateTopic(room.RoomId)
	if nRes.ProblemDetail != nil {
		return nRes.ProblemDetail
	}

	if nRes.Data != nil {
		room.NotificationTopicId = *nRes.Data.(*string)
		room.Modified = time.Now().UnixNano()
		dRes := <-datastore.GetProvider().RoomUpdate(room)
		if dRes.ProblemDetail != nil {
			return dRes.ProblemDetail
		}
	}

	return nil
}

func getUsersByUserIds(requestUserIds []string) ([]*models.User, *models.ProblemDetail) {
	dRes := <-datastore.GetProvider().UserSelectByUserIds(requestUserIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	existUsers := dRes.Data.([]*models.User)
	if len(existUsers) != len(requestUserIds) {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "users",
					Reason: "It contains a userId that does not exist.",
				},
			},
		}
	}

	return existUsers, nil
}
