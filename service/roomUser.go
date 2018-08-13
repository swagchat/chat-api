package service

import (
	"context"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
)

// AddRoomUsers creates room users
func AddRoomUsers(ctx context.Context, req *model.AddRoomUsersRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("AddRoomUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to create room users."
		return errRes
	}

	req.Room = room

	errRes = confirmUserIDsExist(ctx, req.UserIDs, "userIds")
	if errRes != nil {
		errRes.Message = "Failed to create room users."
		return errRes
	}

	errRes = req.Validate()
	if errRes != nil {
		return errRes
	}

	if room.NotificationTopicID == "" {
		notificationTopicID, errRes := createTopic(ctx, room.RoomID)
		if errRes != nil {
			errRes.Message = "Failed to create room users."
			return errRes
		}

		room.NotificationTopicID = notificationTopicID
		room.Modified = time.Now().Unix()
		err := datastore.Provider(ctx).UpdateRoom(room)
		if err != nil {
			return model.NewErrorResponse("Failed to create room users.", http.StatusInternalServerError, model.WithError(err))
		}
	}

	roomUsers := req.GenerateRoomUsers()
	err := datastore.Provider(ctx).InsertRoomUsers(roomUsers)
	if err != nil {
		return model.NewErrorResponse("Failed to create room users.", http.StatusInternalServerError, model.WithError(err))
	}

	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, req.RoomID)

	return nil
}

// RetrieveRoomUsers retrieves room users
func RetrieveRoomUsers(ctx context.Context, req *model.RetrieveRoomUsersRequest) (*model.RoomUsersResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("RetrieveRoomUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmRoomExist(ctx, req.RoomID)
	if errRes != nil {
		errRes.Message = "Failed to get roomUsers."
		return nil, errRes
	}

	res := &model.RoomUsersResponse{}

	roomUsers, err := datastore.Provider(ctx).SelectRoomUsers(
		datastore.SelectRoomUsersOptionWithRoomID(req.RoomID),
		datastore.SelectRoomUsersOptionWithUserIDs(req.UserIDs),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get roomUsers.", http.StatusInternalServerError, model.WithError(err))
	}

	res.Users = roomUsers
	return res, nil
}

// RetrieveRoomUserIDs retrieves room userIds
func RetrieveRoomUserIDs(ctx context.Context, req *model.RetrieveRoomUsersRequest) (*model.RoomUserIdsResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("RetrieveRoomUserIDs", "service")
	defer tracer.Provider(ctx).Finish(span)

	_, errRes := confirmRoomExist(ctx, req.RoomID)
	if errRes != nil {
		errRes.Message = "Failed to get userIds."
		return nil, errRes
	}

	res := &model.RoomUserIdsResponse{}

	userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
		datastore.SelectUserIDsOfRoomUserOptionWithRoomID(req.RoomID),
		datastore.SelectUserIDsOfRoomUserOptionWithRoles(req.RoleIDs),
	)
	if err != nil {
		return nil, model.NewErrorResponse("Failed to get userIds.", http.StatusInternalServerError, model.WithError(err))
	}

	res.UserIDs = userIDs
	return res, nil
}

// UpdateRoomUser updates room user
func UpdateRoomUser(ctx context.Context, req *model.UpdateRoomUserRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("UpdateRoomUser", "service")
	defer tracer.Provider(ctx).Finish(span)

	ru, errRes := confirmRoomUserExist(ctx, req.RoomID, req.UserID)
	if errRes != nil {
		errRes.Message = "Failed to update room user."
		return errRes
	}

	ru.UpdateRoomUser(req)

	err := datastore.Provider(ctx).UpdateRoomUser(ru)
	if err != nil {
		return model.NewErrorResponse("Failed to update room user.", http.StatusInternalServerError, model.WithError(err))
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

	return nil
}

// DeleteRoomUsers deletes room users
func DeleteRoomUsers(ctx context.Context, req *model.DeleteRoomUsersRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("DeleteRoomUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to delete room users."
		return errRes
	}

	req.Room = room

	err := datastore.Provider(ctx).DeleteRoomUsers(
		datastore.DeleteRoomUsersOptionFilterByRoomIDs([]string{req.RoomID}),
		datastore.DeleteRoomUsersOptionFilterByUserIDs(req.UserIDs),
	)
	if err != nil {
		return model.NewErrorResponse("Failed to delete room users.", http.StatusInternalServerError, model.WithError(err))
	}

	go func() {
		rus, err := datastore.Provider(ctx).SelectRoomUsers(
			datastore.SelectRoomUsersOptionWithRoomID(req.RoomID),
			datastore.SelectRoomUsersOptionWithUserIDs(req.UserIDs),
		)
		if err != nil {
			logger.Error(err.Error())
		}

		unsubscribeByRoomUsers(ctx, rus)
	}()

	return nil
}
