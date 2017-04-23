package models

import (
	"net/http"
	"time"

	"github.com/fairway-corp/swagchat-api/utils"
)

type Rooms struct {
	Rooms    []*Room `json:"rooms" db:"-"`
	AllCount int64   `json:"allCount" db:"all_count"`
}

type Room struct {
	Id                  uint64         `json:"-" db:"id"`
	RoomId              string         `json:"roomId" db:"room_id"`
	Name                string         `json:"name" db:"name"`
	PictureUrl          string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl      string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData            utils.JSONText `json:"metaData" db:"meta_data"`
	IsPublic            *bool          `json:"isPublic,omitempty" db:"is_public"`
	LastMessage         string         `json:"lastMessage,omitempty" db:"last_message"`
	LastMessageUpdated  int64          `json:"lastMessageUpdated,omitempty" db:"last_message_updated"`
	NotificationTopicId string         `json:"notificationTopicId,omitempty" db:"notification_topic_id"`
	Created             int64          `json:"created" db:"created"`
	Modified            int64          `json:"modified" db:"modified"`
	Deleted             int64          `json:"-" db:"deleted"`

	Users []*UserForRoom `json:"users,omitempty" db:"-"`
}

type UserForRoom struct {
	// from User
	UserId         string         `json:"userId,omitempty" db:"user_id"`
	Name           string         `json:"name,omitempty" db:"name"`
	PictureUrl     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData       utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
	Created        int64          `json:"created" db:"created"`
	Modified       int64          `json:"modified,omitempty" db:"modified"`

	// from RoomUser
	RuUnreadCount *int64         `json:"ruUnreadCount,omitempty" db:"ru_unread_count"`
	RuMetaData    utils.JSONText `json:"ruMetaData,omitempty" db:"ru_meta_data"`
	RuCreated     int64          `json:"ruCreated,omitempty" db:"ru_created"`
	RuModified    int64          `json:"ruModified,omitempty" db:"ru_modified"`
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

	if r.Name == "" {
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

	if r.IsPublic == nil {
		isPublic := false
		r.IsPublic = &isPublic
	}

	nowTimestamp := time.Now().UnixNano()
	if r.Created == 0 {
		r.Created = nowTimestamp
	}
	r.Modified = nowTimestamp
}

func (r *Room) Put(put *Room) {
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
	if put.IsPublic != nil {
		r.IsPublic = put.IsPublic
	}
}
