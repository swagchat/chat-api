package service

import (
	"context"
	"net/http"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logger"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/tracer"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

// CreateRoomUsers creates room users
func CreateRoomUsers(ctx context.Context, req *model.CreateRoomUsersRequest) *model.ErrorResponse {
	span := tracer.Provider(ctx).StartSpan("CreateRoomUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	room, errRes := confirmRoomExist(ctx, req.RoomID, datastore.SelectRoomOptionWithUsers(true))
	if errRes != nil {
		errRes.Message = "Failed to create room users."
		return errRes
	}

	req.Room = room

	userIDs, errRes := getExistUserIDs(ctx, req.UserIDs)
	if errRes != nil {
		errRes.Message = "Failed to create room users."
		return errRes
	}
	req.UserIDs = userIDs

	errRes = req.Validate()
	if errRes != nil {
		return errRes
	}

	// if room.NotificationTopicID == "" {
	// 	notificationTopicID, pd := createTopicOld(room.RoomID)
	// 	if pd != nil {
	// 		return pd
	// 	}

	// 	room.NotificationTopicID = notificationTopicID
	// 	room.Modified = time.Now().Unix()
	// 	_, err := datastore.Provider(ctx).UpdateRoom(room)
	// 	if err != nil {
	// 		pd := &model.ProblemDetail{
	// 			Message: "Failed to create room users.",
	// 			Status:  http.StatusInternalServerError,
	// 			Error:   err,
	// 		}
	// 		return pd
	// 	}
	// }

	roomUsers := req.GenerateRoomUsers()
	err := datastore.Provider(ctx).InsertRoomUsers(roomUsers)
	if err != nil {
		return model.NewErrorResponse("Failed to create room users.", http.StatusInternalServerError, model.WithError(err))
	}

	go subscribeByRoomUsers(ctx, roomUsers)
	go publishUserJoin(ctx, req.RoomID)

	return nil
}

func GetRoomUsers(ctx context.Context, req *model.GetRoomUsersRequest) (*model.RoomUsersResponse, *model.ErrorResponse) {
	span := tracer.Provider(ctx).StartSpan("GetRoomUsers", "service")
	defer tracer.Provider(ctx).Finish(span)

	res := &model.RoomUsersResponse{}

	if req.ResponseType == scpb.ResponseType_UserIdList {
		userIDs, err := datastore.Provider(ctx).SelectUserIDsOfRoomUser(
			req.RoomID,
			datastore.SelectUserIDsOfRoomUserOptionWithRoles(req.RoleIDs),
		)
		if err != nil {
			return nil, model.NewErrorResponse("Failed to get userIds.", http.StatusInternalServerError, model.WithError(err))
		}

		res.UserIDs = userIDs
		return res, nil
	}

	// TODO
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

	err := datastore.Provider(ctx).DeleteRoomUsers(req.RoomID, req.UserIDs)
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
