package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

// PostRoom is post room
func PostRoom(ctx context.Context, post *model.Room) (*model.Room, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start CreateRoom. Room=[%#v]", post))
	_, pd := selectUser(ctx, post.UserID)
	if pd != nil {
		return nil, pd
	}

	if pd := post.IsValidPost(); pd != nil {
		return nil, pd
	}

	post.BeforePost()

	if post.Type == model.OneOnOne {
		roomUser, err := datastore.Provider(ctx).SelectRoomUserOfOneOnOne(post.UserID, post.UserIDs[0])
		if err != nil {
			pd := &model.ProblemDetail{
				Message: "Room registration failed",
				Status:  http.StatusInternalServerError,
				Error:   err,
			}
			return nil, pd
		}
		if roomUser != nil {
			return nil, &model.ProblemDetail{
				Message: "Resource already exists",
				Status:  http.StatusConflict,
			}
		}
	}

	userIDs := post.UserIDs
	if userIDs == nil {
		userIDs = make([]string, 0)
	}
	userIDs = append(userIDs, post.UserID)
	userIDs = utils.RemoveDuplicate(userIDs)

	rus := make([]*scpb.RoomUser, len(userIDs))
	var zeroValue int32 = 0
	trueValue := true
	for _, userID := range userIDs {
		ru := &scpb.RoomUser{
			RoomID:      post.RoomID,
			UserID:      userID,
			UnreadCount: &zeroValue,
			Display:     &trueValue,
		}
		rus = append(rus, ru)
	}

	if post.UserIDs != nil {
		notificationTopicID, pd := createTopic(post.RoomID)
		if pd != nil {
			return nil, pd
		}
		post.NotificationTopicID = notificationTopicID
	}

	room, err := datastore.Provider(ctx).InsertRoom(post, rus)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Room registration failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(room.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room's users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	room.Users = userForRooms

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsersByRoomID(room.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room's users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	go webhookRoom(ctx, room)
	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, room.RoomID)

	logger.Info(fmt.Sprintf("Finish CreateRoom. Room=[%#v]", room))
	return room, nil
}

// GetRooms is get rooms
func GetRooms(ctx context.Context, values url.Values) (*model.Rooms, *model.ProblemDetail) {
	rooms, err := datastore.Provider(ctx).SelectRooms()
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get rooms failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	count, err := datastore.Provider(ctx).SelectCountRooms()
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room count failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	return &model.Rooms{
		Rooms:    rooms,
		AllCount: count,
	}, nil
}

// GetRoom is get room
func GetRoom(ctx context.Context, roomID string) (*model.Room, *model.ProblemDetail) {
	userID := ctx.Value(utils.CtxUserID).(string)
	user, pd := selectUser(ctx, userID, datastore.UserOptionWithRoles(true))
	if pd != nil {
		return nil, pd
	}

	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return nil, pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room's users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	room.Users = userForRooms

	count, err := datastore.Provider(ctx).SelectCountMessagesByRoomID(user.Roles, roomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	room.MessageCount = count
	return room, nil
}

// PutRoom is put room
func PutRoom(ctx context.Context, put *model.Room) (*model.Room, *model.ProblemDetail) {
	room, pd := selectRoom(ctx, put.RoomID)
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
		pd := &model.ProblemDetail{
			Message: "Update room failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(room.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room's users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	room.Users = userForRooms
	return room, nil
}

// DeleteRoom is delete room
func DeleteRoom(ctx context.Context, roomID string) *model.ProblemDetail {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return pd
	}

	if room.NotificationTopicID != "" {
		nRes := <-notification.Provider().DeleteTopic(room.NotificationTopicID)
		if nRes.ProblemDetail != nil {
			return nRes.ProblemDetail
		}
	}

	room.Deleted = time.Now().Unix()
	err := datastore.Provider(ctx).UpdateRoomDeleted(roomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Delete room failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go unsubscribeByRoomID(ctx, roomID, wg)
		wg.Wait()
		room.NotificationTopicID = ""
		datastore.Provider(ctx).UpdateRoom(room)
	}()

	return nil
}

// GetRoomMessages is get room messages
func GetRoomMessages(ctx context.Context, roomID string, params url.Values) (*model.Messages, *model.ProblemDetail) {
	userID := ctx.Value(utils.CtxUserID).(string)
	user, pd := selectUser(ctx, userID, datastore.UserOptionWithRoles(true))
	if pd != nil {
		return nil, pd
	}

	limit, offset, order, pd := setPagingParams(params)
	if pd != nil {
		return nil, pd
	}

	messages, err := datastore.Provider(ctx).SelectMessages(user.Roles, roomID, limit, offset, order)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room messages failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	returnMessages := &model.Messages{
		Messages: messages,
	}

	count, err := datastore.Provider(ctx).SelectCountMessagesByRoomID(user.Roles, roomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room messages failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	returnMessages.AllCount = count

	updateLastAccessRoomID(ctx, roomID)

	return returnMessages, nil
}

func selectRoom(ctx context.Context, roomID string) (*model.Room, *model.ProblemDetail) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	if room == nil {
		return nil, &model.ProblemDetail{
			Message: "Resource not found",
			Status:  http.StatusNotFound,
		}
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

func setPagingParams(params url.Values) (int, int, string, *model.ProblemDetail) {
	var err error
	limit := 10
	offset := 0
	order := "ASC"
	if limitArray, ok := params["limit"]; ok {
		limit, err = strconv.Atoi(limitArray[0])
		if err != nil {
			return limit, offset, order, &model.ProblemDetail{
				Message: "Request parameter error.",
				Status:  http.StatusBadRequest,
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
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
			return limit, offset, order, &model.ProblemDetail{
				Message: "Request parameter error.",
				Status:  http.StatusBadRequest,
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
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
			return limit, offset, order, &model.ProblemDetail{
				Message: "Request parameter error.",
				Status:  http.StatusBadRequest,
				InvalidParams: []*model.InvalidParam{
					&model.InvalidParam{
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
func RoomAuthz(ctx context.Context, roomID, userID string) *model.ProblemDetail {
	room, pd := selectRoom(ctx, roomID)
	if pd != nil {
		return pd
	}

	if room.Type == model.PublicRoom {
		return nil
	}

	userForRooms, err := datastore.Provider(ctx).SelectUsersForRoom(roomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room users failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return pd
	}

	isAuthorized := false
	for _, userForRoom := range userForRooms {
		if userForRoom.UserID == userID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return &model.ProblemDetail{
			Message: "You are not this room member",
			Status:  http.StatusUnauthorized,
		}
	}

	return nil
}

func updateLastAccessRoomID(ctx context.Context, roomID string) {
	userID := ctx.Value(utils.CtxUserID).(string)
	user, _ := selectUser(ctx, userID)
	user.LastAccessRoomID = roomID
	datastore.Provider(ctx).UpdateUser(user)
}
