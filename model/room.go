package model

import (
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf"
)

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

type Room struct {
	scpb.Room
	MetaData utils.JSONText `json:"metaData" db:"meta_data"`
	Users    []*UserForRoom `json:"users" db:"-"`
	// ID                    uint64         `json:"-" db:"id"`
	// RoomID                string         `json:"roomId" db:"room_id,notnull"`
	// UserID                string         `json:"userId" db:"user_id,notnull"`
	// Name                  string         `json:"name" db:"name,notnull"`
	// PictureURL            string         `json:"pictureUrl,omitempty" db:"picture_url"`
	// InformationURL        string         `json:"informationUrl,omitempty" db:"information_url"`
	// Type                  RoomType       `json:"type,omitempty" db:"type,notnull"`
	// CanLeft               *bool          `json:"canLeft,omitempty" db:"can_left,notnull"`
	// SpeechMode            *SpeechMode    `json:"speechMode,omitempty" db:"speech_mode,notnull"`
	// MetaData              utils.JSONText `json:"metaData" db:"meta_data"`
	// AvailableMessageTypes string         `json:"availableMessageTypes,omitempty" db:"available_message_types"`
	// LastMessage           string         `json:"lastMessage" db:"last_message"`
	// LastMessageUpdated    int64          `json:"lastMessageUpdated" db:"last_message_updated,notnull"`
	// MessageCount          int64          `json:"messageCount" db:"-"`
	// NotificationTopicID   string         `json:"notificationTopicId,omitempty" db:"notification_topic_id"`
	// Created               int64          `json:"created" db:"created,notnull"`
	// Modified              int64          `json:"modified" db:"modified,notnull"`
	// Deleted               int64          `json:"-" db:"deleted,notnull"`

	// Users   []*UserForRoom `json:"users,omitempty" db:"-"`
	// UserIDs []string       `json:"userIds,omitempty" db:"-"`
}

func (r *Room) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	lmu := ""
	if r.LastMessageUpdated != 0 {
		lmu = time.Unix(r.LastMessageUpdated, 0).In(l).Format(time.RFC3339)
	}
	var availableMessageTypesSlice []string
	if r.AvailableMessageTypes != "" {
		availableMessageTypesSlice = strings.Split(r.AvailableMessageTypes, ",")
	}
	return json.Marshal(&struct {
		RoomID                string          `json:"roomId"`
		UserID                string          `json:"userId"`
		Name                  string          `json:"name"`
		PictureURL            string          `json:"pictureUrl,omitempty"`
		InformationURL        string          `json:"informationUrl,omitempty"`
		Type                  scpb.RoomType   `json:"type"`
		CanLeft               bool            `json:"canLeft,omitempty"`
		SpeechMode            scpb.SpeechMode `json:"speechMode,omitempty"`
		MetaData              utils.JSONText  `json:"metaData"`
		AvailableMessageTypes []string        `json:"availableMessageTypes,omitempty"`
		LastMessage           string          `json:"lastMessage"`
		LastMessageUpdated    string          `json:"lastMessageUpdated"`
		MessageCount          int64           `json:"messageCount"`
		NotificationTopicID   string          `json:"notificationTopicId,omitempty"`
		Created               string          `json:"created"`
		Modified              string          `json:"modified"`
		Users                 []*UserForRoom  `json:"users,omitempty"`
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
		AvailableMessageTypes: availableMessageTypesSlice,
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

func (u *Room) UpdateRoom(req *UpdateRoomRequest) {
	// TODO
	if req.Name != nil {
		u.Name = *req.Name
	}

	if req.PictureURL != nil {
		u.PictureURL = *req.PictureURL
	}

	if req.InformationURL != nil {
		u.InformationURL = *req.InformationURL
	}

	if req.MetaData != nil {
		u.MetaData = req.MetaData
	}

	nowTimestamp := time.Now().Unix()
	u.Modified = nowTimestamp
}

type UserForRoom struct {
	scpb.UserForRoom
	// // from User
	// RoomID         string         `json:"roomId" db:"room_id"`
	// UserID         string         `json:"userId" db:"user_id"`
	// Name           string         `json:"name" db:"name"`
	// PictureURL     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	// InformationURL string         `json:"informationUrl,omitempty" db:"information_url"`
	// MetaData       utils.JSONText `json:"metaData" db:"meta_data"`
	// CanBlock       *bool          `json:"canBlock,omitempty" db:"can_block"`
	// LastAccessed   int64          `json:"lastAccessed" db:"last_accessed"`
	// Created        int64          `json:"created" db:"created"`
	// Modified       int64          `json:"modified" db:"modified"`

	// // from RoomUser
	// RuDisplay bool `json:"ruDisplay" db:"ru_display"`
}

func (ufr *UserForRoom) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserID         string         `json:"userId"`
		Name           string         `json:"name"`
		PictureURL     string         `json:"pictureUrl,omitempty"`
		InformationURL string         `json:"informationUrl,omitempty"`
		MetaData       utils.JSONText `json:"metaData"`
		CanBlock       *bool          `json:"canBlock,omitempty"`
		LastAccessed   string         `json:"lastAccessed"`
		Created        string         `json:"created"`
		Modified       string         `json:"modified"`
		RuDisplay      *bool          `json:"ruDisplay"`
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
	MetaData utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
}

func (r *CreateRoomRequest) Validate() *ErrorResponse {
	if r.RoomID != "" && !IsValidID(r.RoomID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to create room.", invalidParams, http.StatusBadRequest, nil)
	}

	if r.UserID == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create room.", invalidParams, http.StatusBadRequest, nil)
	}

	if !IsValidID(r.UserID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to create room.", invalidParams, http.StatusBadRequest, nil)
	}

	if r.Type == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "type is required, but it's empty.",
			},
		}
		return NewErrorResponse("Failed to create room.", invalidParams, http.StatusBadRequest, nil)
	}

	roomType := scpb.RoomType.String(r.Type)
	if roomType == "" {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "type is incorrect.",
			},
		}
		return NewErrorResponse("Failed to create room.", invalidParams, http.StatusBadRequest, nil)
	}

	if r.Type == scpb.RoomType_OneOnOne && len(r.UserIDs) == 0 {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "In case of 1on1 type, it is necessary to set userIds.",
			},
		}
		return NewErrorResponse("Failed to create room.", invalidParams, http.StatusBadRequest, nil)
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

	if crr.RoomID == "" {
		r.RoomID = utils.GenerateUUID()
	} else {
		r.RoomID = crr.RoomID
	}

	r.UserID = crr.UserID
	r.Name = crr.Name
	r.PictureURL = crr.PictureURL
	r.InformationURL = crr.InformationURL
	r.Type = crr.Type

	if crr.CanLeft == nil {
		r.CanLeft = true
	} else {
		r.CanLeft = *crr.CanLeft
	}

	r.SpeechMode = crr.SpeechMode

	if crr.MetaData == nil {
		r.MetaData = []byte("{}")
	} else {
		r.MetaData = crr.MetaData
	}

	nowTimestamp := time.Now().Unix()
	r.LastMessageUpdated = nowTimestamp
	r.Created = nowTimestamp
	r.Modified = nowTimestamp

	return r
}

func (crr *CreateRoomRequest) GenerateRoomUsers() []*RoomUser {
	userIDs := crr.UserIDs
	if userIDs == nil {
		userIDs = make([]string, 0)
	}
	userIDs = append(userIDs, crr.UserID)
	userIDs = utils.RemoveDuplicate(userIDs)

	rus := make([]*RoomUser, len(userIDs))

	for i, v := range userIDs {
		ru := &RoomUser{}
		ru.RoomID = crr.RoomID
		ru.UserID = v
		ru.UnreadCount = int32(0)
		ru.Display = true

		rus[i] = ru
	}
	return rus
}

type GetRoomsRequest struct {
	scpb.GetRoomsRequest
}

type GetRoomRequest struct {
	scpb.GetRoomRequest
}

type UpdateRoomRequest struct {
	scpb.UpdateRoomRequest
	MetaData utils.JSONText `db:"meta_data"`
}

func (uur *UpdateRoomRequest) Validate(room *Room) *ErrorResponse {
	if uur.Type != 0 {
		if room.Type == scpb.RoomType_OneOnOne && uur.Type != scpb.RoomType_OneOnOne {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "type",
					Reason: "In case of 1-on-1 room type, type can not be changed.",
				},
			}
			return NewErrorResponse("Failed to update room.", invalidParams, http.StatusBadRequest, nil)
		} else if room.Type != scpb.RoomType_OneOnOne && uur.Type == scpb.RoomType_OneOnOne {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "type",
					Reason: "In case of not 1-on-1 room type, type can not change to 1-on-1 room type.",
				},
			}
			return NewErrorResponse("Failed to update room.", invalidParams, http.StatusBadRequest, nil)
		}
	}
	return nil
}

type DeleteRoomRequest struct {
	scpb.DeleteRoomRequest
}

type GetRoomMessagesRequest struct {
	scpb.GetRoomMessagesRequest
}
