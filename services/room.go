package services

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
	"go.uber.org/zap/zapcore"
)

// PostRoom is post room
func PostRoom(ctx context.Context, post *models.Room) (*models.Room, *models.ProblemDetail) {
	if pd := post.IsValidPost(); pd != nil {
		return nil, pd
	}
	post.BeforePost()
	post.RequestRoomUserIds.RemoveDuplicate()

	if *post.Type == models.ONE_ON_ONE {
		roomUser, err := datastore.Provider(ctx).SelectRoomUserOfOneOnOne(post.UserId, post.RequestRoomUserIds.UserIds[0])
		if err != nil {
			pd := &models.ProblemDetail{
				Title:  "Room registration failed",
				Status: http.StatusInternalServerError,
				Error:  err,
			}
			return nil, pd
		}
		if roomUser != nil {
			return nil, &models.ProblemDetail{
				Title:  "Resource already exists",
				Status: http.StatusConflict,
			}
		}
	}

	if pd := post.RequestRoomUserIds.IsValid("POST", post); pd != nil {
		return nil, pd
	}

	if post.RequestRoomUserIds.UserIds != nil {
		notificationTopicID, pd := createTopic(post.RoomId)
		if pd != nil {
			return nil, pd
		}
		post.NotificationTopicId = notificationTopicID
	}

	room, err := datastore.Provider(ctx).InsertRoom(post)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Room registration failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(room.RoomId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	room.Users = userForRooms

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsersByRoomID(room.RoomId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, room.RoomId)

	return room, nil
}

// GetRooms is get rooms
func GetRooms(ctx context.Context, values url.Values) (*models.Rooms, *models.ProblemDetail) {
	rooms, err := datastore.Provider(ctx).SelectRooms()
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get rooms failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	count, err := datastore.Provider(ctx).SelectCountRooms()
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room count failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	return &models.Rooms{
		Rooms:    rooms,
		AllCount: count,
	}, nil
}

// GetRoom is get room
func GetRoom(ctx context.Context, roomID string) (*models.Room, *models.ProblemDetail) {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return nil, pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	room.Users = userForRooms

	count, err := datastore.Provider(ctx).SelectCountMessagesByRoomID(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	room.MessageCount = count
	return room, nil
}

// PutRoom is put room
func PutRoom(ctx context.Context, put *models.Room) (*models.Room, *models.ProblemDetail) {
	room, pd := selectRoom(ctx, put.RoomId)
	if pd != nil {
		return nil, pd
	}

	if pd := room.IsValidPut(); pd != nil {
		return nil, pd
	}

	if pd := room.BeforePut(put); pd != nil {
		return nil, pd
	}

	room, err := datastore.Provider(ctx).UpdateRoom(room)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Update room failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(room.RoomId)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room's users failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	room.Users = userForRooms
	return room, nil
}

// DeleteRoom is delete room
func DeleteRoom(ctx context.Context, roomID string) *models.ProblemDetail {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return pd
	}

	if room.NotificationTopicId != "" {
		nRes := <-notification.Provider().DeleteTopic(room.NotificationTopicId)
		if nRes.ProblemDetail != nil {
			return nRes.ProblemDetail
		}
	}

	room.Deleted = time.Now().Unix()
	err := datastore.Provider(ctx).UpdateRoomDeleted(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete room failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return pd
	}

	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go unsubscribeByRoomID(ctx, roomID, wg)
		wg.Wait()
		room.NotificationTopicId = ""
		datastore.Provider(ctx).UpdateRoom(room)
	}()

	return nil
}

// GetRoomMessages is get room messages
func GetRoomMessages(ctx context.Context, roomID string, params url.Values) (*models.Messages, *models.ProblemDetail) {
	limit, offset, order, pd := setPagingParams(params)
	if pd != nil {
		return nil, pd
	}

	messages, err := datastore.Provider(ctx).SelectMessages(roomID, limit, offset, order)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room messages failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	returnMessages := &models.Messages{
		Messages: messages,
	}

	count, err := datastore.Provider(ctx).SelectCountMessagesByRoomID(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room messages failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	returnMessages.AllCount = count
	return returnMessages, nil
}

func selectRoom(ctx context.Context, roomID string) (*models.Room, *models.ProblemDetail) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get room failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if room == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}
	return room, nil
}

func unsubscribeByRoomID(ctx context.Context, roomID string, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByRoomID(roomID)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	<-unsubscribe(ctx, subscriptions)
	if wg != nil {
		wg.Done()
	}
}

func setPagingParams(params url.Values) (int, int, string, *models.ProblemDetail) {
	var err error
	limit := 10
	offset := 0
	order := "ASC"
	if limitArray, ok := params["limit"]; ok {
		limit, err = strconv.Atoi(limitArray[0])
		if err != nil {
			return limit, offset, order, &models.ProblemDetail{
				Title:     "Request parameter error.",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "limit",
						Reason: "limit is incorrect.",
					},
				},
			}
		}
	}
	if offsetArray, ok := params["offset"]; ok {
		offset, err = strconv.Atoi(offsetArray[0])
		if err != nil {
			return limit, offset, order, &models.ProblemDetail{
				Title:     "Request parameter error.",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "offset",
						Reason: "offset is incorrect.",
					},
				},
			}
		}
	}
	if orderArray, ok := params["order"]; ok {
		order := orderArray[0]
		allowedOrders := []string{
			"DESC",
			"desc",
			"ASC",
			"asc",
		}
		if utils.SearchStringValueInSlice(allowedOrders, order) {
			return limit, offset, order, &models.ProblemDetail{
				Title:     "Request parameter error.",
				Status:    http.StatusBadRequest,
				ErrorName: models.ERROR_NAME_INVALID_PARAM,
				InvalidParams: []models.InvalidParam{
					models.InvalidParam{
						Name:   "order",
						Reason: "order is incorrect.",
					},
				},
			}
		}
	}
	return limit, offset, order, nil
}

// RoomAuthz is room authorize
func RoomAuthz(ctx context.Context, roomID, userID string) *models.ProblemDetail {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return pd
	}

	if *room.Type == models.PUBLIC_ROOM {
		return nil
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
	if err != nil {
		pd := &models.ProblemDetail{
			Status: http.StatusInternalServerError,
			Title:  "Get room users failed",
			Error:  err,
		}
		return pd
	}

	isAuthorized := false
	for _, userForRoom := range userForRooms {
		if userForRoom.UserId == userID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return &models.ProblemDetail{
			Title:     "You are not this room member",
			Status:    http.StatusUnauthorized,
			ErrorName: models.ERROR_NAME_UNAUTHORIZED,
		}
	}

	return nil
}
