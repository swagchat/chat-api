package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

// CreateRoom creates room
func CreateRoom(ctx context.Context, req *model.CreateRoomRequest) (*model.Room, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start CreateRoom. Request=[%#v]", req))

	pd := req.Validate()
	if pd != nil {
		return nil, pd
	}

	_, pd = selectUser(ctx, req.UserID)
	if pd != nil {
		return nil, pd
	}

	if req.Type == scpb.RoomType_OneOnOne {
		roomUser, err := datastore.Provider(ctx).SelectRoomUserOfOneOnOne(req.UserID, req.UserIDs[0])
		if err != nil {
			pd := &model.ProblemDetail{
				Message: "Failed to create room.",
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

	r := req.GenerateRoom()
	req.RoomID = r.RoomID

	if len(req.UserIDs) > 0 {
		userIDs, pd := getExistUserIDs(ctx, req.UserIDs)
		if pd != nil {
			return nil, pd
		}
		req.UserIDs = userIDs
	}
	rus := req.GenerateRoomUsers()

	if req.UserIDs != nil {
		notificationTopicID, pd := createTopic(req.RoomID)
		if pd != nil {
			return nil, pd
		}
		r.NotificationTopicID = notificationTopicID
	}

	err := datastore.Provider(ctx).InsertRoom(r, datastore.InsertRoomOptionWithRoomUser(rus))
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create room.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsersByRoomID(r.RoomID)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create room.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	go webhookRoom(ctx, r)
	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, r.RoomID)

	room, err := datastore.Provider(ctx).SelectRoom(
		r.RoomID,
		datastore.SelectRoomOptionWithUsers(true),
	)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Failed to create room.",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	logger.Info("Finish CreateRoom.")
	return room, nil
}

// GetRooms gets rooms
func GetRooms(ctx context.Context, req *model.GetRoomsRequest) (*model.RoomsResponse, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start GetRooms. Request[%#v]", req))

	rooms, err := datastore.Provider(ctx).SelectRooms(
		req.Limit,
		req.Offset,
		datastore.SelectRoomsOptionWithOrders(req.Orders),
	)
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
			Message: "Get rooms failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	res := &model.RoomsResponse{}
	res.Rooms = rooms
	res.AllCount = count
	res.Limit = req.Limit
	res.Offset = req.Offset

	logger.Info(fmt.Sprintf("Finish GetRooms."))
	return res, nil
}

// GetRoom gets a room
func GetRoom(ctx context.Context, req *model.GetRoomRequest) (*model.Room, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start  GetRoom. Request[%#v]", req))

	userID := ctx.Value(utils.CtxUserID).(string)
	user, pd := selectUser(ctx, userID, datastore.UserOptionWithRoles(true))
	if pd != nil {
		return nil, pd
	}

	room, pd := selectRoom(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if pd != nil {
		return nil, pd
	}

	count, err := datastore.Provider(ctx).SelectCountMessages(
		datastore.MessageOptionFilterByRoomID(req.RoomID),
		datastore.MessageOptionFilterByRoleIDs(user.Roles),
	)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	room.MessageCount = count

	logger.Info(fmt.Sprintf("Finish GetRoom. Response[%#v]", room))
	return room, nil
}

// UpdateRoom updates room
func UpdateRoom(ctx context.Context, req *model.UpdateRoomRequest) (*model.Room, *model.ProblemDetail) {
	logger.Info(fmt.Sprintf("Start UpdateRoom. Request[%#v]", req))

	room, pd := selectRoom(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if pd != nil {
		return nil, pd
	}

	pd = req.Validate(room)
	if pd != nil {
		return nil, pd
	}

	room.UpdateRoom(req)

	err := datastore.Provider(ctx).UpdateRoom(room)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Update room failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}

	logger.Info("Finish UpdateRoom.")
	return room, nil
}

// DeleteRoom deletes room
func DeleteRoom(ctx context.Context, req *model.DeleteRoomRequest) *model.ProblemDetail {
	logger.Info(fmt.Sprintf("Start DeleteRoom. Request[%#v]", req))

	room, pd := selectRoom(ctx, req.RoomID)
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
	err := datastore.Provider(ctx).UpdateRoom(room)
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
		go unsubscribeByRoomID(ctx, req.RoomID, wg)
		wg.Wait()
		room.NotificationTopicID = ""
		datastore.Provider(ctx).UpdateRoom(room)
	}()

	logger.Info("Finish DeleteRoom.")
	return nil
}

// GetRoomMessages gets room messages
func GetRoomMessages(ctx context.Context, req *model.GetRoomMessagesRequest) (*model.Messages, *model.ProblemDetail) {
	userID := ctx.Value(utils.CtxUserID).(string)
	user, pd := selectUser(ctx, userID, datastore.UserOptionWithRoles(true))
	if pd != nil {
		return nil, pd
	}

	var roleIDs []int32
	if req.RoleIDs == nil {
		roleIDs = user.Roles
	} else {
		roleIDs = req.RoleIDs
	}

	messages, err := datastore.Provider(ctx).SelectMessages(
		req.Limit,
		req.Offset,
		datastore.MessageOptionOrders(req.Orders),
		datastore.MessageOptionFilterByRoomID(req.RoomID),
		datastore.MessageOptionFilterByRoleIDs(roleIDs),
	)
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

	count, err := datastore.Provider(ctx).SelectCountMessages(
		datastore.MessageOptionFilterByRoomID(req.RoomID),
		datastore.MessageOptionFilterByRoleIDs(req.RoleIDs),
	)
	if err != nil {
		pd := &model.ProblemDetail{
			Message: "Get room messages failed",
			Status:  http.StatusInternalServerError,
			Error:   err,
		}
		return nil, pd
	}
	returnMessages.AllCount = count

	updateLastAccessRoomID(ctx, req.RoomID)

	return returnMessages, nil
}

func selectRoom(ctx context.Context, roomID string, opts ...datastore.SelectRoomOption) (*model.Room, *model.ProblemDetail) {
	room, err := datastore.Provider(ctx).SelectRoom(roomID, opts...)
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

// RoomAuthz is room authorize
func RoomAuthz(ctx context.Context, roomID, userID string) *model.ProblemDetail {
	room, pd := selectRoom(ctx, roomID, datastore.SelectRoomOptionWithUsers(true))
	if pd != nil {
		return pd
	}

	if room.Type == scpb.RoomType_PublicRoom {
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
