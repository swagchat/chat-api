package models

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/swagchat/chat-api/utils"
)

// RoomType is room type
type RoomType int

const (
	OneOnOne RoomType = iota + 1
	PrivateRoom
	PublicRoom
	NoticeRoom
	CustomerRoom
	RoomTypeEnd
)

// SpeechMode is speech mode
type SpeechMode int

const (
	SpeechModeNone SpeechMode = iota + 1
	SpeechModeWakeupWebToWeb
	SpeechModeWakeupWebToCloud
	SpeechModeWakeupCloudToCloud
	SpeechModeAlways
	SpeechModeManual
	SpeechModeEnd
)

func (rt RoomType) String() string {
	switch rt {
	case OneOnOne:
		return "OneOnOne"
	case PrivateRoom:
		return "PrivateRoom"
	case PublicRoom:
		return "PublicRoom"
	case NoticeRoom:
		return "NoticeRoom"
	case CustomerRoom:
		return "CustomerRoom"
	default:
		return "Unknown"
	}
}

type Rooms struct {
	Rooms    []*Room `json:"rooms" db:"-"`
	AllCount int64   `json:"allCount" db:"all_count"`
}

type Room struct {
	ID                    uint64         `json:"-" db:"id"`
	RoomID                string         `json:"roomId" db:"room_id,notnull"`
	UserID                string         `json:"userId" db:"user_id,notnull"`
	Name                  string         `json:"name" db:"name,notnull"`
	PictureURL            string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationURL        string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData              utils.JSONText `json:"metaData" db:"meta_data"`
	AvailableMessageTypes string         `json:"availableMessageTypes,omitempty" db:"available_message_types"`
	Type                  RoomType       `json:"type,omitempty" db:"type,notnull"`
	LastMessage           string         `json:"lastMessage" db:"last_message"`
	LastMessageUpdated    int64          `json:"lastMessageUpdated" db:"last_message_updated,notnull"`
	MessageCount          int64          `json:"messageCount" db:"-"`
	NotificationTopicID   string         `json:"notificationTopicId,omitempty" db:"notification_topic_id"`
	IsCanLeft             *bool          `json:"isCanLeft,omitempty" db:"is_can_left,notnull"`
	IsShowUsers           *bool          `json:"isShowUsers,omitempty" db:"is_show_users,notnull"`
	SpeechMode            *SpeechMode    `json:"speechMode,omitempty" db:"speech_mode,notnull"`
	Created               int64          `json:"created" db:"created,notnull"`
	Modified              int64          `json:"modified" db:"modified,notnull"`
	Deleted               int64          `json:"-" db:"deleted,notnull"`

	Users []*UserForRoom `json:"users,omitempty" db:"-"`
	RequestRoomUserIDs
}

type UserForRoom struct {
	// from User
	UserID         string         `json:"userId" db:"user_id"`
	Name           string         `json:"name" db:"name"`
	PictureURL     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationURL string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData       utils.JSONText `json:"metaData" db:"meta_data"`
	Role           *Role          `json:"role,omitempty" db:"role,notnull"`
	IsBot          *bool          `json:"isBot,omitempty" db:"is_bot,notnull"`
	IsCanBlock     *bool          `json:"isCanBlock,omitempty" db:"is_can_block,notnull"`
	IsShowUsers    *bool          `json:"isShowUsers,omitempty" db:"is_show_users,notnull"`
	LastAccessed   int64          `json:"lastAccessed" db:"last_accessed"`
	Created        int64          `json:"created" db:"created"`
	Modified       int64          `json:"modified" db:"modified"`

	// from RoomUser
	RuUnreadCount int64          `json:"ruUnreadCount" db:"ru_unread_count"`
	RuMetaData    utils.JSONText `json:"ruMetaData" db:"ru_meta_data"`
	RuCreated     int64          `json:"ruCreated" db:"ru_created"`
	RuModified    int64          `json:"ruModified" db:"ru_modified"`
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
		RoomID                string         `json:"roomId"`
		UserID                string         `json:"userId"`
		Name                  string         `json:"name"`
		PictureURL            string         `json:"pictureUrl,omitempty"`
		InformationURL        string         `json:"informationUrl,omitempty"`
		MetaData              utils.JSONText `json:"metaData"`
		AvailableMessageTypes []string       `json:"availableMessageTypes,omitempty"`
		Type                  RoomType       `json:"type"`
		LastMessage           string         `json:"lastMessage"`
		LastMessageUpdated    string         `json:"lastMessageUpdated"`
		MessageCount          int64          `json:"messageCount"`
		NotificationTopicID   string         `json:"notificationTopicId,omitempty"`
		IsCanLeft             *bool          `json:"isCanLeft,omitempty"`
		IsShowUsers           *bool          `json:"isShowUsers,omitempty"`
		SpeechMode            *SpeechMode    `json:"speechMode,omitempty"`
		Created               string         `json:"created"`
		Modified              string         `json:"modified"`
		Users                 []*UserForRoom `json:"users,omitempty"`
	}{
		RoomID:                r.RoomID,
		UserID:                r.UserID,
		Name:                  r.Name,
		PictureURL:            r.PictureURL,
		InformationURL:        r.InformationURL,
		MetaData:              r.MetaData,
		AvailableMessageTypes: availableMessageTypesSlice,
		Type:               r.Type,
		LastMessage:        r.LastMessage,
		LastMessageUpdated: lmu,
		MessageCount:       r.MessageCount,
		IsCanLeft:          r.IsCanLeft,
		IsShowUsers:        r.IsShowUsers,
		SpeechMode:         r.SpeechMode,
		Created:            time.Unix(r.Created, 0).In(l).Format(time.RFC3339),
		Modified:           time.Unix(r.Modified, 0).In(l).Format(time.RFC3339),
		Users:              r.Users,
	})
}

func (ufr *UserForRoom) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserID         string         `json:"userId"`
		Name           string         `json:"name"`
		PictureURL     string         `json:"pictureUrl,omitempty"`
		InformationURL string         `json:"informationUrl,omitempty"`
		MetaData       utils.JSONText `json:"metaData"`
		Role           *Role          `json:"role,omitempty`
		IsBot          *bool          `json:"isBot,omitempty"`
		IsCanBlock     *bool          `json:"isCanBlock,omitempty"`
		IsShowUsers    *bool          `json:"isShowUsers,omitempty"`
		LastAccessed   string         `json:"lastAccessed"`
		Created        string         `json:"created"`
		Modified       string         `json:"modified"`
		RuUnreadCount  int64          `json:"ruUnreadCount"`
		RuMetaData     utils.JSONText `json:"ruMetaData"`
		RuCreated      string         `json:"ruCreated"`
		RuModified     string         `json:"ruModified"`
	}{
		UserID:         ufr.UserID,
		Name:           ufr.Name,
		PictureURL:     ufr.PictureURL,
		InformationURL: ufr.InformationURL,
		MetaData:       ufr.MetaData,
		Role:           ufr.Role,
		IsBot:          ufr.IsBot,
		IsCanBlock:     ufr.IsCanBlock,
		IsShowUsers:    ufr.IsShowUsers,
		LastAccessed:   time.Unix(ufr.LastAccessed, 0).In(l).Format(time.RFC3339),
		Created:        time.Unix(ufr.Created, 0).In(l).Format(time.RFC3339),
		Modified:       time.Unix(ufr.Modified, 0).In(l).Format(time.RFC3339),
		RuUnreadCount:  ufr.RuUnreadCount,
		RuMetaData:     ufr.RuMetaData,
		RuCreated:      time.Unix(ufr.RuCreated, 0).In(l).Format(time.RFC3339),
		RuModified:     time.Unix(ufr.RuModified, 0).In(l).Format(time.RFC3339),
	})
}

func (r *Room) IsValidPost() *ProblemDetail {
	if r.RoomID != "" && !utils.IsValidID(r.RoomID) {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "roomId",
					Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}

	if r.UserID == "" {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	if !utils.IsValidID(r.UserID) {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}

	if r.Type == 0 {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "type",
					Reason: "type is required, but it's empty.",
				},
			},
		}
	}

	if !(r.Type > 0 && r.Type < RoomTypeEnd) {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "type",
					Reason: "type is incorrect.",
				},
			},
		}
	}

	if r.Type != CustomerRoom && r.UserIDs == nil {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userIds",
					Reason: fmt.Sprintf("When room type is %d, userIds is a required. ", r.Type),
				},
			},
		}
	}

	if r.SpeechMode != nil && !(*r.SpeechMode > 0 && *r.SpeechMode < SpeechModeEnd) {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "speechMode",
					Reason: "speechMode is incorrect.",
				},
			},
		}
	}

	return nil
}

func (r *Room) IsValidPut() *ProblemDetail {
	return nil
}

func (r *Room) BeforePost() {
	if r.RoomID == "" {
		r.RoomID = utils.GenerateUUID()
	}

	if r.MetaData == nil {
		r.MetaData = []byte("{}")
	}

	if r.IsCanLeft == nil {
		isCanLeft := true
		r.IsCanLeft = &isCanLeft
	}

	if r.IsShowUsers == nil {
		isShowUsers := true
		r.IsShowUsers = &isShowUsers
	}

	if r.SpeechMode == nil {
		speechMode := SpeechMode(SpeechModeNone)
		r.SpeechMode = &speechMode
	}

	nowTimestamp := time.Now().Unix()
	r.LastMessageUpdated = nowTimestamp
	r.Created = nowTimestamp
	r.Modified = nowTimestamp
}

func (r *Room) BeforePut(put *Room) *ProblemDetail {
	if put.Name != "" {
		r.Name = put.Name
	}
	if put.PictureURL != "" {
		r.PictureURL = put.PictureURL
	}
	if put.InformationURL != "" {
		r.InformationURL = put.InformationURL
	}
	if put.MetaData != nil {
		r.MetaData = put.MetaData
	}
	if put.IsCanLeft != nil {
		r.IsCanLeft = put.IsCanLeft
	}
	if put.IsShowUsers != nil {
		r.IsShowUsers = put.IsShowUsers
	}
	if put.Type != 0 {
		if r.Type == OneOnOne && put.Type != OneOnOne {
			return &ProblemDetail{
				Title:  "Request error",
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "type",
						Reason: "In case of 1-on-1 room type, type can not be changed.",
					},
				},
			}
		} else if r.Type != OneOnOne && put.Type == OneOnOne {
			return &ProblemDetail{
				Title:  "Request error",
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "type",
						Reason: "In case of not 1-on-1 room type, type can not change to 1-on-1 room type.",
					},
				},
			}
		} else {
			r.Type = put.Type
		}
	}
	return nil
}
