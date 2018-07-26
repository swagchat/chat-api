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

// CreateRoom creates a room
func CreateRoom(ctx context.Context, req *model.CreateRoomRequest) (*model.Room, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  CreateRoom. Request=[%#v]", req))

	errRes := req.Validate()
	if errRes != nil {
		return nil, errRes
	}

	_, errRes = confirmUserExist(ctx, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to create room."
		return nil, errRes
	}

	if req.Type == scpb.RoomType_OneOnOne {
		roomUser, err := datastore.Provider(ctx).SelectRoomUserOfOneOnOne(req.UserID, req.UserIDs[0])
		if err != nil {
			return nil, model.NewErrorResponse("Failed to create room.", nil, http.StatusBadRequest, err)
		}
		if roomUser != nil {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "userId",
					Reason: "userId does not exist",
				},
			}
			return nil, model.NewErrorResponse("Failed to create a room.", invalidParams, http.StatusConflict, nil)
		}
	}

	r := req.GenerateRoom()
	req.RoomID = r.RoomID

	if len(req.UserIDs) > 0 {
		userIDs, errRes := getExistUserIDs(ctx, req.UserIDs)
		if errRes != nil {
			errRes.Message = "Failed to create room."
			return nil, errRes
		}
		req.UserIDs = userIDs
	}
	rus := req.GenerateRoomUsers()

	if req.UserIDs != nil {
		notificationTopicID, errRes := createTopic(req.RoomID)
		if errRes != nil {
			return nil, errRes
		}
		r.NotificationTopicID = notificationTopicID
	}

	err := datastore.Provider(ctx).InsertRoom(r, datastore.InsertRoomOptionWithRoomUser(rus))
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create room.", nil, http.StatusInternalServerError, err)
	}

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsers(
		datastore.SelectRoomUsersOptionWithRoomID(r.RoomID),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create room.", nil, http.StatusInternalServerError, err)
	}

	go webhookRoom(ctx, r)
	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, r.RoomID)

	room, err := datastore.Provider(ctx).SelectRoom(
		r.RoomID,
		datastore.SelectRoomOptionWithUsers(true),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to create room.", nil, http.StatusInternalServerError, err)
	}

	logger.Info("Finish CreateRoom.")
	return room, nil
}

// GetRooms gets rooms
func GetRooms(ctx context.Context, req *model.GetRoomsRequest) (*model.RoomsResponse, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetRooms. Request[%#v]", req))

	rooms, err := datastore.Provider(ctx).SelectRooms(
		req.Limit,
		req.Offset,
		datastore.SelectRoomsOptionWithOrders(req.Orders),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get rooms.", nil, http.StatusInternalServerError, err)
	}

	count, err := datastore.Provider(ctx).SelectCountRooms()
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get rooms.", nil, http.StatusInternalServerError, err)
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
func GetRoom(ctx context.Context, req *model.GetRoomRequest) (*model.Room, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  GetRoom. Request[%#v]", req))

	userID := ctx.Value(utils.CtxUserID).(string)
	user, errRes := confirmUserExist(ctx, userID, datastore.SelectUserOptionWithRoles(true))
	if errRes != nil {
		return nil, errRes
	}

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to get room."
		return nil, errRes
	}

	count, err := datastore.Provider(ctx).SelectCountMessages(
		datastore.SelectMessagesOptionFilterByRoomID(req.RoomID),
		datastore.SelectMessagesOptionFilterByRoleIDs(user.Roles),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get room.", nil, http.StatusInternalServerError, err)

	}
	room.MessageCount = count

	logger.Info(fmt.Sprintf("Finish GetRoom. Response[%#v]", room))
	return room, nil
}

// UpdateRoom updates room
func UpdateRoom(ctx context.Context, req *model.UpdateRoomRequest) (*model.Room, *model.ErrorResponse) {
	logger.Info(fmt.Sprintf("Start  UpdateRoom. Request[%#v]", req))

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to update room."
		return nil, errRes
	}

	errRes = req.Validate(room)
	if errRes != nil {
		return nil, errRes
	}

	room.UpdateRoom(req)

	err := datastore.Provider(ctx).UpdateRoom(room)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to update room.", nil, http.StatusInternalServerError, err)
	}

	logger.Info("Finish UpdateRoom.")
	return room, nil
}

// DeleteRoom deletes room
func DeleteRoom(ctx context.Context, req *model.DeleteRoomRequest) *model.ErrorResponse {
	logger.Info(fmt.Sprintf("Start  DeleteRoom. Request[%#v]", req))

	room, errRes := confirmRoomExist(ctx, req.RoomID)
	if errRes != nil {
		errRes.Message = "Failed to delete room."
		return errRes
	}

	if room.NotificationTopicID != "" {
		nRes := <-notification.Provider().DeleteTopic(room.NotificationTopicID)
		if nRes.Error != nil {
			return model.NewErrorResponse("Failed to delete room.", nil, http.StatusInternalServerError, nRes.Error)
		}
	}

	room.Deleted = time.Now().Unix()
	err := datastore.Provider(ctx).UpdateRoom(room)
	if err != nil {
		return model.NewErrorResponse("Failed to delete room.", nil, http.StatusInternalServerError, err)
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
	user, pd := selectUser(ctx, userID, datastore.SelectUserOptionWithRoles(true))
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
		datastore.SelectMessagesOptionOrders(req.Orders),
		datastore.SelectMessagesOptionFilterByRoomID(req.RoomID),
		datastore.SelectMessagesOptionFilterByRoleIDs(roleIDs),
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
		datastore.SelectMessagesOptionFilterByRoomID(req.RoomID),
		datastore.SelectMessagesOptionFilterByRoleIDs(req.RoleIDs),
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
