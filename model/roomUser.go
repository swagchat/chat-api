package model

import (
	"net/http"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

type RoomUsers struct {
	RoomID    string
	UserIDs   []string
	Display   bool
	Room      *Room
	RoomUsers []*RoomUser
}

type RoomUser struct {
	RoomID      string `json:"roomId,omitempty" db:"room_id"`
	UserID      string `json:"userId,omitempty" db:"user_id"`
	UnreadCount int32  `json:"unreadCount,omitempty" db:"unread_count"`
	Display     bool   `json:"display,omitempty" db:"display"`
}

func (ru *RoomUser) Validate() *ProblemDetail {
	if ru.RoomID != "" && !utils.IsValidID(ru.RoomID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "roomId",
					Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if ru.UserID != "" && !utils.IsValidID(ru.UserID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	// if mode == Update && rus.Room.Type == OneOnOne {
	// 	if len(rus.Room.Users) == 2 {
	// 		return &ProblemDetail{
	// 			Message: "Invalid params",
	// 			InvalidParams: []*InvalidParam{
	// 				&InvalidParam{
	// 					Name:   "userIds",
	// 					Reason: "In case of 1-on-1 room type, It can only update once.",
	// 				},
	// 			},
	// 			Status: http.StatusBadRequest,
	// 		}
	// 	}
	// }

	return nil
}

type ErrorRoomUser struct {
	UserId string         `json:"userId,omitempty"`
	Error  *ProblemDetail `json:"error"`
}

type ResponseRoomUser struct {
	RoomUsers []RoomUser      `json:"roomUsers,omitempty"`
	Errors    []ErrorRoomUser `json:"errors,omitempty"`
}

type RequestRoomUserIDs struct {
	UserIDs []string `json:"userIds,omitempty" db:"-"`
}

// ImportFromPbCreateUserRolesRequest import from CreateUserRolesRequest proto
func (rus *RoomUsers) ImportFromPbCreateUserRolesRequest(req *scpb.CreateRoomUsersRequest) {
	rus.RoomID = req.RoomId
	rus.UserIDs = utils.RemoveDuplicate(req.UserIds)
	rus.Display = req.Display
}

func (rus *RoomUsers) GenerateRoomUsers() {
	roomUsers := make([]*RoomUser, len(rus.UserIDs))
	for i, userID := range rus.UserIDs {
		roomUsers[i] = &RoomUser{
			RoomID:      rus.RoomID,
			UserID:      userID,
			UnreadCount: 0,
			Display:     rus.Display,
		}
	}
	rus.RoomUsers = roomUsers
}

func (rus *RoomUsers) Validate() *ProblemDetail {
	if len(rus.UserIDs) == 0 {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []*InvalidParam{
				&InvalidParam{
					Name:   "userIds",
					Reason: "Not set.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if rus.Room.Type == OneOnOne {
		if len(rus.UserIDs) == 2 {
			return &ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []*InvalidParam{
					&InvalidParam{
						Name:   "userIds",
						Reason: "In case of 1-on-1 room type, It can only update once.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}
	}

	return nil
}

// ImportFromPbUpdateUserRoleRequest import from UpdateUserRoleRequest proto
func (ru *RoomUser) ImportFromPbUpdateUserRoleRequest(req *scpb.UpdateRoomUserRequest) {
	ru.UserID = req.UserId
}
