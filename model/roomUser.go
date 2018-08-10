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
	if crur.Room == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: "roomId is not exist",
			},
		}
		return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if crur.Room.Type == scpb.RoomType_OneOnOneRoom {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "room.type",
				Reason: "In case of 1-on-1 room type, Can not add room user.",
			},
		}
		return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if len(crur.UserIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userIds",
				Reason: "userIds is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	for _, userID := range crur.UserIDs {
		if userID == crur.Room.UserID {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "userIds",
					Reason: "userIds can not include room's userId.",
				},
			}
			return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
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

type RetrieveRoomUsersRequest struct {
	scpb.RetrieveRoomUsersRequest
}

type RoomUsersResponse struct {
	scpb.RoomUsersResponse
	Users []*RoomUser
}

func (rur *RoomUsersResponse) ConvertToPbRoomUsers() *scpb.RoomUsersResponse {
	res := &scpb.RoomUsersResponse{}

	rus := make([]*scpb.RoomUser, len(rur.Users))
	for i := 0; i < len(rur.Users); i++ {
		ru := rur.Users[i]
		pbRu := &scpb.RoomUser{
			RoomID:      ru.RoomID,
			UserID:      ru.UserID,
			UnreadCount: ru.UnreadCount,
			Display:     ru.Display,
		}
		rus[i] = pbRu
	}
	res.Users = rus

	return res
}

type RoomUserIdsResponse struct {
	scpb.RoomUserIdsResponse
}

func (rur *RoomUserIdsResponse) ConvertToPbRoomUserIDs() *scpb.RoomUserIdsResponse {
	return &scpb.RoomUserIdsResponse{
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
	if drur.Room == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: "roomId is not exist",
			},
		}
		return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if drur.Room.Type == scpb.RoomType_OneOnOneRoom {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "room.type",
				Reason: "In case of 1-on-1 room type, Can not delete room user.",
			},
		}
		return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if len(drur.UserIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userIds",
				Reason: "userIds is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	for _, userID := range drur.UserIDs {
		if userID == drur.Room.UserID {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "userIds",
					Reason: "userIds can not include room's userId.",
				},
			}
			return NewErrorResponse("Failed to create room users.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}

	return nil
}
