package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf"
)

type RoomUser struct {
	scpb.RoomUser
}

func (ru *RoomUser) UpdateRoomUser(req *UpdateRoomUserRequest) {
	if req.UnreadCount != nil {
		ru.UnreadCount = *req.UnreadCount
	}

	if req.Display != nil {
		ru.Display = *req.Display
	}
}

type CreateRoomUsersRequest struct {
	scpb.CreateRoomUsersRequest
	Room *Room
}

func (crur *CreateRoomUsersRequest) GenerateRoomUsers() []*RoomUser {
	roomUsers := make([]*RoomUser, len(crur.UserIDs))
	for i, userID := range crur.UserIDs {
		ru := &RoomUser{}
		ru.RoomID = crur.RoomID
		ru.UserID = userID
		ru.UnreadCount = int32(0)
		ru.Display = crur.Display
		roomUsers[i] = ru
	}
	return roomUsers
}

func (crur *CreateRoomUsersRequest) Validate() *ErrorResponse {
	if crur.Room.Type == scpb.RoomType_OneOnOne {
		if len(crur.UserIDs) != 1 {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "userIds",
					Reason: "In case of 1-on-1 room type, only one user can be specified for userIDs.",
				},
			}
			return NewErrorResponse("Failed to create a user.", invalidParams, http.StatusBadRequest, nil)
		}
	}
	return nil
}

type GetUserIdsOfRoomUserRequest struct {
	scpb.GetUserIdsOfRoomUserRequest
}

type UpdateRoomUserRequest struct {
	scpb.UpdateRoomUserRequest
}

type DeleteRoomUsersRequest struct {
	scpb.DeleteRoomUsersRequest
	Room *Room
}

func (drur *DeleteRoomUsersRequest) Validate() *ProblemDetail {
	if drur.Room.Type == scpb.RoomType_OneOnOne {
		if len(drur.Room.Users)-len(drur.UserIDs) != 1 {
			return &ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []*InvalidParam{
					&InvalidParam{
						Name:   "userIds",
						Reason: "In case of 1-on-1 room type, only one user must be specified.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}
	}
	return nil
}
