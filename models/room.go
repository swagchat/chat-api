package models

import (
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/fairway-corp/swagchat-api/utils"
)

type RoomType int

const (
	ONE_ON_ONE RoomType = iota + 1
	PRIVATE_ROOM
	PUBLIC_ROOM
	NOTICE_ROOM
	ROOM_TYPE_END
)

func (rt RoomType) String() string {
	switch rt {
	case PRIVATE_ROOM:
		return "PRIVATE_ROOM"
	case PUBLIC_ROOM:
		return "PUBLIC_ROOM"
	case ONE_ON_ONE:
		return "ONE_ON_ONE"
	default:
		return "Unknown"
	}
}

type Rooms struct {
	Rooms    []*Room `json:"rooms" db:"-"`
	AllCount int64   `json:"allCount" db:"all_count"`
}

type Room struct {
	Id                    uint64         `json:"-" db:"id"`
	RoomId                string         `json:"roomId" db:"room_id,notnull"`
	UserId                string         `json:"userId" db:"user_id,notnull"`
	Name                  string         `json:"name" db:"name,notnull"`
	PictureUrl            string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl        string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData              utils.JSONText `json:"metaData" db:"meta_data"`
	AvailableMessageTypes string         `json:"availableMessageTypes,omitempty" db:"available_message_types"`
	Type                  *RoomType      `json:"type,omitempty" db:"type,notnull"`
	LastMessage           string         `json:"lastMessage" db:"last_message"`
	LastMessageUpdated    int64          `json:"lastMessageUpdated" db:"last_message_updated,notnull"`
	MessageCount          int64          `json:"messageCount" db:"-"`
	NotificationTopicId   string         `json:"notificationTopicId,omitempty" db:"notification_topic_id"`
	IsCanLeft             *bool          `json:"isCanLeft,omitempty" db:"is_can_left,notnull"`
	IsShowUsers           *bool          `json:"isShowUsers,omitempty" db:"is_show_users,notnull"`
	Created               int64          `json:"created" db:"created,notnull"`
	Modified              int64          `json:"modified" db:"modified,notnull"`
	Deleted               int64          `json:"-" db:"deleted,notnull"`

	Users []*UserForRoom `json:"users,omitempty" db:"-"`
}

type UserForRoom struct {
	// from User
	UserId         string         `json:"userId" db:"user_id"`
	Name           string         `json:"name" db:"name"`
	PictureUrl     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData       utils.JSONText `json:"metaData" db:"meta_data"`
	IsCanLeft      *bool          `json:"isCanBlock,omitempty" db:"is_can_block,notnull"`
	IsShowUsers    *bool          `json:"isShowUsers,omitempty" db:"is_show_users,notnull"`
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
		RoomId                string         `json:"roomId"`
		UserId                string         `json:"userId"`
		Name                  string         `json:"name"`
		PictureUrl            string         `json:"pictureUrl,omitempty"`
		InformationUrl        string         `json:"informationUrl,omitempty"`
		MetaData              utils.JSONText `json:"metaData"`
		AvailableMessageTypes []string       `json:"availableMessageTypes,omitempty"`
		Type                  *RoomType      `json:"type"`
		LastMessage           string         `json:"lastMessage"`
		LastMessageUpdated    string         `json:"lastMessageUpdated"`
		MessageCount          int64          `json:"messageCount"`
		NotificationTopicId   string         `json:"notificationTopicId,omitempty"`
		IsCanLeft             *bool          `json:"isCanLeft,omitempty"`
		IsShowUsers           *bool          `json:"isShowUsers,omitempty"`
		Created               string         `json:"created"`
		Modified              string         `json:"modified"`
		Users                 []*UserForRoom `json:"users,omitempty"`
	}{
		RoomId:                r.RoomId,
		UserId:                r.UserId,
		Name:                  r.Name,
		PictureUrl:            r.PictureUrl,
		InformationUrl:        r.InformationUrl,
		MetaData:              r.MetaData,
		AvailableMessageTypes: availableMessageTypesSlice,
		Type:               r.Type,
		LastMessage:        r.LastMessage,
		LastMessageUpdated: lmu,
		MessageCount:       r.MessageCount,
		IsCanLeft:          r.IsCanLeft,
		IsShowUsers:        r.IsShowUsers,
		Created:            time.Unix(r.Created, 0).In(l).Format(time.RFC3339),
		Modified:           time.Unix(r.Modified, 0).In(l).Format(time.RFC3339),
		Users:              r.Users,
	})
}

func (ufr *UserForRoom) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserId         string         `json:"userId"`
		Name           string         `json:"name"`
		PictureUrl     string         `json:"pictureUrl,omitempty"`
		InformationUrl string         `json:"informationUrl,omitempty"`
		MetaData       utils.JSONText `json:"metaData"`
		IsCanLeft      *bool          `json:"isCanBlock,omitempty"`
		IsShowUsers    *bool          `json:"isShowUsers,omitempty"`
		Created        string         `json:"created"`
		Modified       string         `json:"modified"`
		RuUnreadCount  int64          `json:"ruUnreadCount"`
		RuMetaData     utils.JSONText `json:"ruMetaData"`
		RuCreated      string         `json:"ruCreated"`
		RuModified     string         `json:"ruModified"`
	}{
		UserId:         ufr.UserId,
		Name:           ufr.Name,
		PictureUrl:     ufr.PictureUrl,
		InformationUrl: ufr.InformationUrl,
		MetaData:       ufr.MetaData,
		IsCanLeft:      ufr.IsCanLeft,
		IsShowUsers:    ufr.IsShowUsers,
		Created:        time.Unix(ufr.Created, 0).In(l).Format(time.RFC3339),
		Modified:       time.Unix(ufr.Modified, 0).In(l).Format(time.RFC3339),
		RuUnreadCount:  ufr.RuUnreadCount,
		RuMetaData:     ufr.RuMetaData,
		RuCreated:      time.Unix(ufr.RuCreated, 0).In(l).Format(time.RFC3339),
		RuModified:     time.Unix(ufr.RuModified, 0).In(l).Format(time.RFC3339),
	})
}

func (r *Room) IsValid() *ProblemDetail {
	if r.RoomId != "" && !utils.IsValidId(r.RoomId) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "roomId",
					Reason: "roomId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}

	if r.UserId == "" {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is required, but it's empty.",
				},
			},
		}
	}

	if r.UserId != "" && !utils.IsValidId(r.UserId) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}

	if r.Type == nil {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "type",
					Reason: "type is required, but it's empty.",
				},
			},
		}
	}

	if !(*r.Type > 0 && *r.Type < ROOM_TYPE_END) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "type",
					Reason: "type is incorrect.",
				},
			},
		}
	}

	if *r.Type != ONE_ON_ONE && r.Name == "" {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create room item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "name",
					Reason: "name is required, but it's empty.",
				},
			},
		}
	}

	return nil
}

func (r *Room) BeforeSave() {
	if r.RoomId == "" {
		r.RoomId = utils.CreateUuid()
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

	nowTimestamp := time.Now().Unix()
	if r.Created == 0 {
		r.Created = nowTimestamp
	}
	r.Modified = nowTimestamp
}

func (r *Room) Put(put *Room) *ProblemDetail {
	if put.Name != "" {
		r.Name = put.Name
	}
	if put.PictureUrl != "" {
		r.PictureUrl = put.PictureUrl
	}
	if put.InformationUrl != "" {
		r.InformationUrl = put.InformationUrl
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
	if put.Type != nil {
		if *r.Type == ONE_ON_ONE && *put.Type != ONE_ON_ONE {
			return &ProblemDetail{
				Title:     "Request parameter error. (Update room item)",
				Status:    http.StatusBadRequest,
				ErrorName: ERROR_NAME_INVALID_PARAM,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "type",
						Reason: "In case of 1-on-1 room type, type can not be changed.",
					},
				},
			}
		} else if *r.Type != ONE_ON_ONE && *put.Type == ONE_ON_ONE {
			return &ProblemDetail{
				Title:     "Request parameter error. (Update room item)",
				Status:    http.StatusBadRequest,
				ErrorName: ERROR_NAME_INVALID_PARAM,
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
