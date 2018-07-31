package model

import (
	"net/http"

	scpb "github.com/swagchat/protobuf/protoc-gen-go"
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

func (crur *CreateRoomUsersRequest) Validate() *ErrorResponse {
	if crur.Room.Type == scpb.RoomType_RoomTypeOneOnOne {
		if len(crur.UserIDs) != 1 {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "userIds",
					Reason: "In case of 1-on-1 room type, only one user can be specified for userIDs.",
				},
			}
			return NewErrorResponse("Failed to create a user.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}
	return nil
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

type GetRoomUsersRequest struct {
	scpb.GetRoomUsersRequest
}

type RoomUsersResponse struct {
	scpb.RoomUsersResponse
}

func (rur *RoomUsersResponse) ConvertToPbRoomUsers() *scpb.RoomUsersResponse {
	return &scpb.RoomUsersResponse{
		Users:   rur.Users,
		UserIDs: rur.UserIDs,
	}
}

type UpdateRoomUserRequest struct {
	scpb.UpdateRoomUserRequest
}

type DeleteRoomUsersRequest struct {
	scpb.DeleteRoomUsersRequest
	Room *Room
}

func (drur *DeleteRoomUsersRequest) Validate() *ErrorResponse {
	if drur.Room.Type == scpb.RoomType_RoomTypeOneOnOne {
		if len(drur.Room.Users)-len(drur.UserIDs) != 1 {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "userIds",
					Reason: "In case of 1-on-1 room type, only one user must be specified.",
				},
			}
			return NewErrorResponse("Failed to delete room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}
	return nil
}
