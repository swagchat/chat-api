package protobuf

import (
	"net/http"

	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/utils"
)

func (ru *RoomUser) IsValid() *model.ProblemDetail {
	if ru.RoomID != "" && !utils.IsValidID(ru.RoomID) {
		return &model.ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []model.InvalidParam{
				model.InvalidParam{
					Name:   "roomId",
					Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if ru.UserID != "" && !utils.IsValidID(ru.UserID) {
		return &model.ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []model.InvalidParam{
				model.InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	return nil
}

func (ru *RoomUser) BeforeSave() {
	// nowTimestamp := time.Now().Unix()
	// if ru.Created == 0 {
	// 	ru.Created = nowTimestamp
	// }
	// ru.Modified = nowTimestamp
}

func (ru *RoomUser) Put(put *RoomUser) {
	// if put.UnreadCount != nil {
	ru.UnreadCount = put.UnreadCount
	// }
	// if put.MetaData != nil {
	// 	ru.MetaData = put.MetaData
	// }
}

type ErrorRoomUser struct {
	UserId string               `json:"userId,omitempty"`
	Error  *model.ProblemDetail `json:"error"`
}

type ResponseRoomUser struct {
	RoomUsers []RoomUser      `json:"roomUsers,omitempty"`
	Errors    []ErrorRoomUser `json:"errors,omitempty"`
}

type RequestRoomUserIDs struct {
	UserIDs []string `json:"userIds,omitempty" db:"-"`
}

type RoomUsers struct {
	RoomUsers []*RoomUser `json:"roomUsers"`
}

func (rus *PostRoomUserReq) IsValid(method string, r *model.Room) *model.ProblemDetail {
	if len(rus.UserIDs) == 0 {
		return &model.ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []model.InvalidParam{
				model.InvalidParam{
					Name:   "userIds",
					Reason: "Not set.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	// if *r.Type == OneOnOne {
	// 	for _, userId := range rus.UserIds {
	// 		if userId == r.UserId {
	// 			return &model.ProblemDetail{
	// 				Message:  "Invalid params",
	// 				Status: http.StatusBadRequest,
	// 				model.InvalidParams: []model.InvalidParam{
	// 					model.InvalidParam{
	// 						Name:   "userIds",
	// 						Reason: "In case of 1-on-1 room type, it must always set one userId different from this room's userId.",
	// 					},
	// 				},
	// 			}
	// 		}
	// 	}
	// }

	if method == "POST" && r.Type == model.OneOnOne {
		if len(rus.UserIDs) == 2 {
			return &model.ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []model.InvalidParam{
					model.InvalidParam{
						Name:   "userIds",
						Reason: "In case of 1-on-1 room type, It can only update once.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}
	}

	if method == "PUT" && r.Type == model.OneOnOne {
		if len(r.Users) == 2 {
			return &model.ProblemDetail{
				Message: "Invalid params",
				InvalidParams: []model.InvalidParam{
					model.InvalidParam{
						Name:   "userIds",
						Reason: "In case of 1-on-1 room type, It can only update once.",
					},
				},
				Status: http.StatusBadRequest,
			}
		}
	}

	// if method == "DELETE" && *r.Type == OneOnOne {
	// 	return &model.ProblemDetail{
	// 		Title:     "Request error",
	// 		Status:    http.StatusBadRequest,
	// 		ErrorName: ERROR_NAME_OPERATION_NOT_PERMITTED,
	// 	}
	// }

	return nil
}

func (prur *PostRoomUserReq) RemoveDuplicate() {
	if prur != nil {
		prur.UserIDs = utils.RemoveDuplicate(prur.UserIDs)
	}
}

func (drur *DeleteRoomUserReq) RemoveDuplicate() {
	if drur != nil {
		drur.UserIDs = utils.RemoveDuplicate(drur.UserIDs)
	}
}
