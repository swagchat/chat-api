package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/fairway-corp/swagchat-api/datastore"
	"github.com/fairway-corp/swagchat-api/messaging"
	"github.com/fairway-corp/swagchat-api/models"
	"github.com/fairway-corp/swagchat-api/notification"
	"github.com/fairway-corp/swagchat-api/utils"
)

func PostRoomUsers(roomId string, requestRoomUsers *models.RoomUsers) (*models.ResponseRoomUser, *models.ProblemDetail) {
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Create room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	var dRes datastore.StoreResult

	dRes = <-dp.RoomSelect(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	room := dRes.Data.(*models.Room)

	if len(requestRoomUsers.Users) == 0 {
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

	requestUserIds := utils.RemoveDuplicate(requestRoomUsers.Users)

	dRes = <-dp.UserSelectByUserIds(requestUserIds)
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

	np := notification.GetProvider()
	if np != nil && room.NotificationTopicId == nil {
		nRes := <-np.CreateTopic(roomId)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}

		room.NotificationTopicId = nRes.Data.(*string)
		room.Modified = time.Now().UnixNano()
		dRes = <-dp.RoomUpdate(room)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}
	}

	ruRes := addRoomUsers(roomId, *room.NotificationTopicId, existUsers)

	dRes = <-dp.RoomUsersSelectUserIds(roomId)
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
			dRes := <-dp.RoomUsersSelect(&roomId, deleteUserIds)
			if dRes.ProblemDetail != nil {
				return nil, dRes.ProblemDetail
			}
			delRuRes := deleteRoomUsers(dRes.Data.([]*models.RoomUser))
			ruRes.Errors = append(ruRes.Errors, delRuRes.Errors...)
		}
	}

	return ruRes, nil
}

/*
func GetRoomUsers(roomId string) (*models.Room, *models.ProblemDetail) {
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Get room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomUserUsersSelect(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	room := &models.Room{
		Users: dRes.Data.([]*models.User),
	}
	return room, nil
}
*/

func PutRoomUsers(roomId string, requestRoomUsers *models.RoomUsers) (*models.ResponseRoomUser, *models.ProblemDetail) {
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Add room's user list)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	var dRes datastore.StoreResult

	dRes = <-dp.RoomSelect(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	room := dRes.Data.(*models.Room)

	if len(requestRoomUsers.Users) == 0 {
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

	requestUserIds := utils.RemoveDuplicate(requestRoomUsers.Users)

	dRes = <-dp.UserSelectByUserIds(requestUserIds)
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

	np := notification.GetProvider()
	if np != nil && room.NotificationTopicId == nil {
		nRes := <-np.CreateTopic(roomId)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}

		room.NotificationTopicId = nRes.Data.(*string)
		room.Modified = time.Now().UnixNano()
		dRes = <-dp.RoomUpdate(room)
		if dRes.ProblemDetail != nil {
			return nil, dRes.ProblemDetail
		}
	}

	ruRes := addRoomUsers(roomId, *room.NotificationTopicId, existUsers)
	go publishUserJoin(roomId)
	return ruRes, nil
}

func DeleteRoomUsers(roomId string, requestRoomUsers *models.RoomUsers) (*models.ResponseRoomUser, *models.ProblemDetail) {
	if roomId == "" {
		return nil, &models.ProblemDetail{
			Title:     "Request parameter error. (Delete room's user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomSelect(roomId)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}

	deleteUserIds := utils.RemoveDuplicate(requestRoomUsers.Users)
	dRes = <-dp.RoomUsersSelect(&roomId, deleteUserIds)
	if dRes.ProblemDetail != nil {
		return nil, dRes.ProblemDetail
	}
	ruRes := deleteRoomUsers(dRes.Data.([]*models.RoomUser))

	return ruRes, nil
}

func PutRoomUser(roomId, userId string, requestRoomUser *models.RoomUser) *models.ProblemDetail {
	if roomId == "" {
		return &models.ProblemDetail{
			Title:     "Request parameter error. (Update room's user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is required, but it's empty.",
				},
			},
		}
	}
	if userId == "" {
		return &models.ProblemDetail{
			Title:     "Request parameter error. (Update room's user item)",
			Status:    http.StatusBadRequest,
			ErrorName: models.ERROR_NAME_INVALID_PARAM,
			InvalidParams: []models.InvalidParam{
				models.InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	dp := datastore.GetProvider()
	dRes := <-dp.RoomUserSelect(roomId, userId)
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

	dRes = <-dp.RoomUserUpdate(roomUser)
	if dRes.ProblemDetail != nil {
		return dRes.ProblemDetail
	}

	if requestRoomUser.UnreadCount != nil {
		dRes = <-dp.UserUnreadCountRecalc(userId)
		return dRes.ProblemDetail
	}

	return nil
}

func publishUserJoin(roomId string) {
	dp := datastore.GetProvider()
	dRes := <-dp.RoomSelectUsersForRoom(roomId)
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

func addRoomUsers(roomId, notificationTopicId string, existUsers []*models.User) *models.ResponseRoomUser {
	np := notification.GetProvider()
	dp := datastore.GetProvider()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ruChan := make(chan *models.RoomUser, 1)
	errRuChan := make(chan *models.ErrorRoomUser, 1)

	var zero int64
	zero = 0
	ruRes := &models.ResponseRoomUser{
		RoomUsers: make([]models.RoomUser, 0),
		Errors:    make([]models.ErrorRoomUser, 0),
	}

	d := utils.NewDispatcher(10)
	for _, user := range existUsers {
		ctx = context.WithValue(ctx, "user", user)
		d.Work(ctx, func(ctx context.Context) {
			targetUser := ctx.Value("user").(*models.User)
			if np != nil && targetUser.NotificationDeviceId == nil && targetUser.DeviceToken != nil {
				nRes := <-np.CreateEndpoint("", models.PLATFORM_IOS, *targetUser.DeviceToken)
				if nRes.ProblemDetail != nil {
					errRuChan <- &models.ErrorRoomUser{
						UserId: targetUser.UserId,
						Error:  nRes.ProblemDetail,
					}
				} else {
					targetUser.NotificationDeviceId = nRes.Data.(*string)
					targetUser.Modified = time.Now().UnixNano()
					dRes := <-dp.UserUpdate(targetUser)
					if dRes.ProblemDetail != nil {
						errRuChan <- &models.ErrorRoomUser{
							UserId: targetUser.UserId,
							Error:  dRes.ProblemDetail,
						}
					}
				}
			}

			dRes := <-dp.RoomUserSelect(roomId, targetUser.UserId)
			if dRes.ProblemDetail != nil {
				errRuChan <- &models.ErrorRoomUser{
					UserId: targetUser.UserId,
					Error:  dRes.ProblemDetail,
				}
			}

			var roomUser *models.RoomUser
			if np != nil {
				if dRes.Data == nil {
					roomUser = &models.RoomUser{
						RoomId:      roomId,
						UserId:      targetUser.UserId,
						UnreadCount: &zero,
						MetaData:    []byte("{}"),
					}
					if targetUser.NotificationDeviceId != nil {
						nRes := <-np.Subscribe(notificationTopicId, *targetUser.NotificationDeviceId)
						if nRes.ProblemDetail != nil {
							errRuChan <- &models.ErrorRoomUser{
								UserId: targetUser.UserId,
								Error:  nRes.ProblemDetail,
							}
						} else {
							roomUser.NotificationSubscribeId = nRes.Data.(*string)
						}
					}
					dRes = <-dp.RoomUserInsert(roomUser)
					if dRes.ProblemDetail != nil {
						errRuChan <- &models.ErrorRoomUser{
							UserId: targetUser.UserId,
							Error:  dRes.ProblemDetail,
						}
					}
				} else {
					roomUser = dRes.Data.(*models.RoomUser)
					if roomUser.NotificationSubscribeId == nil {
						if targetUser.NotificationDeviceId != nil {
							nRes := <-np.Subscribe(notificationTopicId, *targetUser.NotificationDeviceId)
							if nRes.ProblemDetail != nil {
								errRuChan <- &models.ErrorRoomUser{
									UserId: targetUser.UserId,
									Error:  nRes.ProblemDetail,
								}
							} else {
								roomUser.NotificationSubscribeId = nRes.Data.(*string)
							}
						}
						dRes = <-dp.RoomUserUpdate(roomUser)
						if dRes.ProblemDetail != nil {
							errRuChan <- &models.ErrorRoomUser{
								UserId: targetUser.UserId,
								Error:  dRes.ProblemDetail,
							}
						}
					}
				}
			}
			if roomUser != nil {
				ruChan <- roomUser
			}

			select {
			case <-ctx.Done():
				return
			case ru := <-ruChan:
				ruRes.RoomUsers = append(ruRes.RoomUsers, *ru)
				return
			case errRu := <-errRuChan:
				ruRes.Errors = append(ruRes.Errors, *errRu)
				return
			}
		})
	}
	d.Wait()
	return ruRes
}

func deleteRoomUsers(roomUsers []*models.RoomUser) *models.ResponseRoomUser {
	np := notification.GetProvider()
	dp := datastore.GetProvider()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	delDoneChan := make(chan bool, 1)
	errRuChan := make(chan *models.ErrorRoomUser, 1)
	ruRes := &models.ResponseRoomUser{
		Errors: make([]models.ErrorRoomUser, 0),
	}

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, "roomUser", roomUser)
		d.Work(ctx, func(ctx context.Context) {
			targetRoomUser := ctx.Value("roomUser").(*models.RoomUser)
			if targetRoomUser.NotificationSubscribeId != nil {
				nRes := <-np.Unsubscribe(*targetRoomUser.NotificationSubscribeId)
				if nRes.ProblemDetail != nil {
					errRuChan <- &models.ErrorRoomUser{
						UserId: targetRoomUser.UserId,
						Error:  nRes.ProblemDetail,
					}
				} else {
					dRes := <-dp.RoomUserDelete(targetRoomUser.RoomId, []string{targetRoomUser.UserId})
					if dRes.ProblemDetail != nil {
						errRuChan <- &models.ErrorRoomUser{
							UserId: targetRoomUser.UserId,
							Error:  dRes.ProblemDetail,
						}
					}
				}
			} else {
				dRes := <-dp.RoomUserDelete(targetRoomUser.RoomId, []string{targetRoomUser.UserId})
				if dRes.ProblemDetail != nil {
					errRuChan <- &models.ErrorRoomUser{
						UserId: targetRoomUser.UserId,
						Error:  dRes.ProblemDetail,
					}
				}
			}
			delDoneChan <- true

			select {
			case <-ctx.Done():
				return
			case <-delDoneChan:
				return
			case errRu := <-errRuChan:
				ruRes.Errors = append(ruRes.Errors, *errRu)
				return
			}
		})
	}
	d.Wait()
	return ruRes
}
