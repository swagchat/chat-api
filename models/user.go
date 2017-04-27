package models

import (
	"net/http"
	"time"

	"github.com/fairway-corp/swagchat-api/utils"
)

type Users struct {
	Users []*User `json:"users"`
}

type User struct {
	Id             uint64         `json:"-" db:"id"`
	UserId         string         `json:"userId" db:"user_id"`
	Name           string         `json:"name" db:"name"`
	PictureUrl     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl string         `json:"informationUrl,omitempty" db:"information_url"`
	UnreadCount    *uint64        `json:"unreadCount" db:"unread_count"`
	MetaData       utils.JSONText `json:"metaData" db:"meta_data"`
	Created        int64          `json:"created" db:"created"`
	Modified       int64          `json:"modified" db:"modified"`
	Deleted        int64          `json:"-" db:"deleted"`

	Rooms   []*RoomForUser `json:"rooms,omitempty" db:"-"`
	Devices []*Device      `json:"devices,omitempty" db:"-"`
}

type RoomForUser struct {
	// from room
	RoomId             string         `json:"roomId,omitempty" db:"room_id"`
	UserId             string         `json:"userId,omitempty" db:"user_id"`
	Name               string         `json:"name,omitempty" db:"name"`
	PictureUrl         string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl     string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData           utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
	LastMessage        string         `json:"lastMessage,omitempty" db:"last_message"`
	LastMessageUpdated int64          `json:"lastMessageUpdated,omitempty" db:"last_message_updated"`
	Created            int64          `json:"created" db:"created"`
	Modified           int64          `json:"modified,omitempty" db:"modified"`

	// from RoomUser
	RuUnreadCount *int64         `json:"ruUnreadCount,omitempty" db:"ru_unread_count"`
	RuMetaData    utils.JSONText `json:"ruMetaData,omitempty" db:"ru_meta_data"`
	RuCreated     int64          `json:"ruCreated,omitempty" db:"ru_created"`
	RuModified    int64          `json:"ruModified,omitempty" db:"ru_modified"`
}

func (u *User) IsValid() *ProblemDetail {
	if u.UserId != "" && !utils.IsValidId(u.UserId) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create user item)",
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

	if u.Name == "" {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create user item)",
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

func (u *User) BeforeSave() {
	if u.UserId == "" {
		u.UserId = utils.CreateUuid()
	}

	if u.MetaData == nil {
		u.MetaData = []byte("{}")
	}

	if u.UnreadCount == nil {
		unreadCount := uint64(0)
		u.UnreadCount = &unreadCount
	}

	nowTimestamp := time.Now().UnixNano()
	if u.Created == 0 {
		u.Created = nowTimestamp
	}
	u.Modified = nowTimestamp
}

func (u *User) Put(put *User) {
	if put.Name != "" {
		u.Name = put.Name
	}
	if put.PictureUrl != "" {
		u.PictureUrl = put.PictureUrl
	}
	if put.InformationUrl != "" {
		u.InformationUrl = put.InformationUrl
	}
	if put.UnreadCount != nil {
		u.UnreadCount = put.UnreadCount
	}
	if put.MetaData != nil {
		u.MetaData = put.MetaData
	}
}
