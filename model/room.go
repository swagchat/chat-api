package model

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

type Room struct {
	scpb.Room
	MetaData JSONText `json:"metaData" db:"meta_data"`
	// Users    []*MiniUser `json:"users" db:"-"`
	// ID                    uint64         `json:"-" db:"id"`
	// RoomID                string         `json:"roomId" db:"room_id,notnull"`
	// UserID                string         `json:"userId" db:"user_id,notnull"`
	// Name                  string         `json:"name" db:"name,notnull"`
	// PictureURL            string         `json:"pictureUrl,omitempty" db:"picture_url"`
	// InformationURL        string         `json:"informationUrl,omitempty" db:"information_url"`
	// Type                  RoomType       `json:"type,omitempty" db:"type,notnull"`
	// CanLeft               *bool          `json:"canLeft,omitempty" db:"can_left,notnull"`
	// SpeechMode            *SpeechMode    `json:"speechMode,omitempty" db:"speech_mode,notnull"`
	// MetaData              JSONText `json:"metaData" db:"meta_data"`
	// AvailableMessageTypes string         `json:"availableMessageTypes,omitempty" db:"available_message_types"`
	// LastMessage           string         `json:"lastMessage" db:"last_message"`
	// LastMessageUpdated    int64          `json:"lastMessageUpdated" db:"last_message_updated,notnull"`
	// MessageCount          int64          `json:"messageCount" db:"-"`
	// NotificationTopicID   string         `json:"notificationTopicId,omitempty" db:"notification_topic_id"`
	// Created               int64          `json:"created" db:"created,notnull"`
	// Modified              int64          `json:"modified" db:"modified,notnull"`
	// Deleted               int64          `json:"-" db:"deleted,notnull"`

	// Users   []*MiniUser `json:"users,omitempty" db:"-"`
	// UserIDs []string       `json:"userIds,omitempty" db:"-"`
}

func (r *Room) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	lmu := ""
	if r.LastMessageUpdated != 0 {
		lmu = time.Unix(r.LastMessageUpdated, 0).In(l).Format(time.RFC3339)
	}
	return json.Marshal(&struct {
		RoomID                string           `json:"roomId"`
		UserID                string           `json:"userId"`
		Name                  string           `json:"name"`
		PictureURL            string           `json:"pictureUrl"`
		InformationURL        string           `json:"informationUrl"`
		Type                  scpb.RoomType    `json:"type"`
		CanLeft               bool             `json:"canLeft"`
		SpeechMode            scpb.SpeechMode  `json:"speechMode"`
		MetaData              JSONText         `json:"metaData"`
		AvailableMessageTypes string           `json:"availableMessageTypes"`
		LastMessage           string           `json:"lastMessage"`
		LastMessageUpdated    string           `json:"lastMessageUpdated"`
		MessageCount          int64            `json:"messageCount"`
		NotificationTopicID   string           `json:"notificationTopicId"`
		Created               string           `json:"created"`
		Modified              string           `json:"modified"`
		Users                 []*scpb.MiniUser `json:"users,omitempty"`
	}{
		RoomID:                r.RoomID,
		UserID:                r.UserID,
		Name:                  r.Name,
		PictureURL:            r.PictureURL,
		InformationURL:        r.InformationURL,
		Type:                  r.Type,
		CanLeft:               r.CanLeft,
		SpeechMode:            r.SpeechMode,
		MetaData:              r.MetaData,
		AvailableMessageTypes: r.AvailableMessageTypes,
		LastMessage:           r.LastMessage,
		LastMessageUpdated:    lmu,
		MessageCount:          r.MessageCount,
		Created:               time.Unix(r.Created, 0).In(l).Format(time.RFC3339),
		Modified:              time.Unix(r.Modified, 0).In(l).Format(time.RFC3339),
		Users:                 r.Users,
	})
}

func (r *Room) ConvertToPbRoom() *scpb.Room {
	// TODO
	pbRoom := &scpb.Room{
		RoomID:         r.RoomID,
		UserID:         r.UserID,
		Name:           r.Name,
		PictureURL:     r.PictureURL,
		InformationURL: r.InformationURL,
		MetaData:       r.MetaData,
	}
	return pbRoom
}

func (r *Room) UpdateRoom(req *UpdateRoomRequest) {
	// TODO
	if req.Name != nil {
		r.Name = *req.Name
	}

	if req.PictureURL != nil {
		r.PictureURL = *req.PictureURL
	}

	if req.InformationURL != nil {
		r.InformationURL = *req.InformationURL
	}

	if req.Type != nil {
		r.Type = *req.Type
	}

	if req.CanLeft != nil {
		r.CanLeft = *req.CanLeft
	}

	if req.SpeechMode != nil {
		r.SpeechMode = *req.SpeechMode
	}

	if req.MetaData != nil {
		r.MetaData = req.MetaData
	}

	if req.AvailableMessageTypes != nil {
		r.AvailableMessageTypes = *req.AvailableMessageTypes
	}

	nowTimestamp := time.Now().Unix()
	r.Modified = nowTimestamp
}

type MiniUser struct {
	scpb.MiniUser
}

func (ufr *MiniUser) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserID         string   `json:"userId"`
		Name           string   `json:"name"`
		PictureURL     string   `json:"pictureUrl,omitempty"`
		InformationURL string   `json:"informationUrl,omitempty"`
		MetaData       JSONText `json:"metaData"`
		CanBlock       *bool    `json:"canBlock,omitempty"`
		LastAccessed   string   `json:"lastAccessed"`
		Created        string   `json:"created"`
		Modified       string   `json:"modified"`
		RuDisplay      *bool    `json:"ruDisplay,omitempty"`
	}{
		UserID:         ufr.UserID,
		Name:           ufr.Name,
		PictureURL:     ufr.PictureURL,
		InformationURL: ufr.InformationURL,
		MetaData:       ufr.MetaData,
		CanBlock:       ufr.CanBlock,
		LastAccessed:   time.Unix(ufr.LastAccessed, 0).In(l).Format(time.RFC3339),
		Created:        time.Unix(ufr.Created, 0).In(l).Format(time.RFC3339),
		Modified:       time.Unix(ufr.Modified, 0).In(l).Format(time.RFC3339),
		RuDisplay:      ufr.RuDisplay,
	})
}

type CreateRoomRequest struct {
	scpb.CreateRoomRequest
	MetaData JSONText `json:"metaData,omitempty" db:"meta_data"`
}

func (r *CreateRoomRequest) Validate() *ErrorResponse {
	if r.RoomID != nil && *r.RoomID != "" && !isValidID(*r.RoomID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if r.UserID == nil || *r.UserID == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if !isValidID(*r.UserID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if r.Type == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "type is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	roomType := scpb.RoomType.String(*r.Type)
	if roomType == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "type is incorrect.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if *r.Type == scpb.RoomType_OneOnOneRoom && len(r.UserIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "In case of 1on1 type, it is necessary to set userIds.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if r.MetaData != nil && !isJSON(r.MetaData.String()) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "metaData",
				Reason: "metaData is not json format.",
			},
		}
		return NewErrorResponse("Failed to create room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	// if r.SpeechMode != nil && !(*r.SpeechMode > 0 && *r.SpeechMode < SpeechModeEnd) {
	// 	return &ProblemDetail{
	// 		Message: "Invalid params",
	// 		InvalidParams: []*InvalidParam{
	// 			&InvalidParam{
	// 				Name:   "speechMode",
	// 				Reason: "speechMode is incorrect.",
	// 			},
	// 		},
	// 		Status: http.StatusBadRequest,
	// 	}
	// }

	return nil
}

func (crr *CreateRoomRequest) GenerateRoom() *Room {
	r := &Room{}

	if crr.RoomID == nil || *crr.RoomID == "" {
		r.RoomID = utils.GenerateUUID()
	} else {
		r.RoomID = *crr.RoomID
	}

	r.UserID = *crr.UserID

	if crr.Name == nil {
		r.Name = ""
	} else {
		r.Name = *crr.Name
	}

	if crr.PictureURL == nil {
		r.PictureURL = ""
	} else {
		r.PictureURL = *crr.PictureURL
	}

	if crr.InformationURL == nil {
		r.InformationURL = ""
	} else {
		r.InformationURL = *crr.InformationURL
	}

	if crr.Type == nil {
		r.Type = scpb.RoomType_PublicRoom
	} else {
		r.Type = *crr.Type
	}

	if crr.CanLeft == nil {
		r.CanLeft = true
	} else {
		r.CanLeft = *crr.CanLeft
	}

	if crr.SpeechMode == nil {
		r.SpeechMode = scpb.SpeechMode_SpeechModeNone
	} else {
		r.SpeechMode = *crr.SpeechMode
	}

	if crr.MetaData == nil {
		r.MetaData = []byte("{}")
	} else {
		r.MetaData = crr.MetaData
	}

	if crr.AvailableMessageTypes == nil {
		r.AvailableMessageTypes = ""
	} else {
		r.AvailableMessageTypes = *crr.AvailableMessageTypes
	}

	nowTimestamp := time.Now().Unix()
	r.LastMessageUpdated = nowTimestamp
	r.Created = nowTimestamp
	r.Modified = nowTimestamp

	return r
}

func (crr *CreateRoomRequest) GenerateRoomUsers() []*RoomUser {
	rus := make([]*RoomUser, len(crr.UserIDs)+1)
	me := &RoomUser{}
	me.RoomID = *crr.RoomID
	me.UserID = *crr.UserID
	me.UnreadCount = int32(0)
	me.Display = true

	rus[0] = me
	for i := 0; i < len(crr.UserIDs); i++ {
		ru := &RoomUser{}
		ru.RoomID = *crr.RoomID
		ru.UserID = crr.UserIDs[i]
		ru.UnreadCount = int32(0)
		ru.Display = true
		rus[i+1] = ru
	}
	return rus
}

type RetrieveRoomsRequest struct {
	scpb.RetrieveRoomsRequest
}

type RoomsResponse struct {
	scpb.RoomsResponse
	Rooms []*Room `json:"rooms"`
}

func (rr *RoomsResponse) ConvertToPbRooms() *scpb.RoomsResponse {
	rooms := make([]*scpb.Room, len(rr.Rooms))
	for i, v := range rr.Rooms {
		metaData, _ := v.MetaData.MarshalJSON()
		rooms[i] = &scpb.Room{
			RoomID:         v.RoomID,
			UserID:         v.UserID,
			Name:           v.Name,
			PictureURL:     v.PictureURL,
			InformationURL: v.InformationURL,
			MetaData:       metaData,
			Created:        v.Created,
			Modified:       v.Modified,
		}
	}
	return &scpb.RoomsResponse{
		Rooms: rooms,
	}
}

type RetrieveRoomRequest struct {
	scpb.RetrieveRoomRequest
}

type UpdateRoomRequest struct {
	scpb.UpdateRoomRequest
	MetaData JSONText `json:"metaData,omitempty" db:"meta_data"`
}

func (uur *UpdateRoomRequest) Validate(room *Room) *ErrorResponse {
	// TODO
	if uur.Type != nil {
		if room.Type == scpb.RoomType_OneOnOneRoom && *uur.Type != scpb.RoomType_OneOnOneRoom {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "type",
					Reason: "In case of 1-on-1 room type, type can not be changed.",
				},
			}
			return NewErrorResponse("Failed to update room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		} else if room.Type != scpb.RoomType_OneOnOneRoom && *uur.Type == scpb.RoomType_OneOnOneRoom {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "type",
					Reason: "In case of not 1-on-1 room type, type can not change to 1-on-1 room type.",
				},
			}
			return NewErrorResponse("Failed to update room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}

	if uur.MetaData != nil && !isJSON(uur.MetaData.String()) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "metaData",
				Reason: "metaData is not json format.",
			},
		}
		return NewErrorResponse("Failed to update room.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	return nil
}

func (uur *UpdateRoomRequest) GenerateRoomUsers(room *Room) []*RoomUser {
	rus := make([]*RoomUser, len(uur.UserIDs)+1)
	me := &RoomUser{}
	me.RoomID = room.RoomID
	me.UserID = room.UserID
	me.UnreadCount = int32(0)
	me.Display = true

	rus[0] = me
	for i := 0; i < len(uur.UserIDs); i++ {
		ru := &RoomUser{}
		ru.RoomID = room.RoomID
		ru.UserID = uur.UserIDs[i]
		ru.UnreadCount = int32(0)
		ru.Display = true
		rus[i+1] = ru
	}
	return rus
}

type DeleteRoomRequest struct {
	scpb.DeleteRoomRequest
}

type RetrieveRoomMessagesRequest struct {
	scpb.RetrieveRoomMessagesRequest
}

type RoomMessagesResponse struct {
	scpb.RoomMessagesResponse
	Messages []*Message `json:"messages"`
}

func (rmr *RoomMessagesResponse) ConvertToPbRoomMessages() *scpb.RoomMessagesResponse {
	pbRoomMessages := &scpb.RoomMessagesResponse{}

	messages := make([]*scpb.Message, len(rmr.Messages))
	for i, v := range rmr.Messages {
		payload, _ := v.Payload.MarshalJSON()
		messages[i] = &scpb.Message{
			MessageID: v.MessageID,
			RoomID:    v.RoomID,
			UserID:    v.UserID,
			Type:      v.Type,
			Payload:   payload,
			Role:      v.Role,
			Created:   v.Created,
			Modified:  v.Modified,
			Deleted:   v.Deleted,
			EventName: v.EventName,
			UserIDs:   v.UserIDs,
		}
	}
	pbRoomMessages.Messages = messages
	pbRoomMessages.AllCount = rmr.AllCount
	pbRoomMessages.Limit = rmr.Limit
	pbRoomMessages.Offset = rmr.Offset
	pbRoomMessages.Orders = rmr.Orders
	pbRoomMessages.RoomID = rmr.RoomID
	pbRoomMessages.RoleIDs = rmr.RoleIDs
	return pbRoomMessages
}
