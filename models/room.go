package models

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/fairway-corp/swagchat-api/utils"
)

type Rooms struct {
	Rooms    []*Room `json:"rooms" db:"-"`
	AllCount int64   `json:"allCount" db:"all_count"`
}

type Room struct {
	Id                  uint64         `json:"-" db:"id"`
	RoomId              string         `json:"roomId" db:"room_id,notnull"`
	UserId              string         `json:"userId" db:"user_id,notnull"`
	Name                string         `json:"name" db:"name,notnull"`
	PictureUrl          string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl      string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData            utils.JSONText `json:"metaData" db:"meta_data"`
	IsPublic            *bool          `json:"isPublic,omitempty" db:"is_public,notnull"`
	LastMessage         string         `json:"lastMessage" db:"last_message"`
	LastMessageUpdated  int64          `json:"lastMessageUpdated" db:"last_message_updated,notnull"`
	NotificationTopicId string         `json:"notificationTopicId,omitempty" db:"notification_topic_id"`
	Created             int64          `json:"created" db:"created,notnull"`
	Modified            int64          `json:"modified" db:"modified,notnull"`
	Deleted             int64          `json:"-" db:"deleted,notnull"`

	Users []*UserForRoom `json:"users,omitempty" db:"-"`
}

type UserForRoom struct {
	// from User
	UserId         string         `json:"userId" db:"user_id"`
	Name           string         `json:"name" db:"name"`
	PictureUrl     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData       utils.JSONText `json:"metaData" db:"meta_data"`
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
	return json.Marshal(&struct {
		RoomId              string         `json:"roomId"`
		UserId              string         `json:"userId"`
		Name                string         `json:"name"`
		PictureUrl          string         `json:"pictureUrl,omitempty"`
		InformationUrl      string         `json:"informationUrl,omitempty"`
		MetaData            utils.JSONText `json:"metaData"`
		IsPublic            *bool          `json:"isPublic"`
		LastMessage         string         `json:"lastMessage"`
		LastMessageUpdated  string         `json:"lastMessageUpdated"`
		NotificationTopicId string         `json:"notificationTopicId,omitempty"`
		Created             string         `json:"created"`
		Modified            string         `json:"modified"`
		Users               []*UserForRoom `json:"users,omitempty"`
	}{
		RoomId:             r.RoomId,
		UserId:             r.UserId,
		Name:               r.Name,
		PictureUrl:         r.PictureUrl,
		InformationUrl:     r.InformationUrl,
		MetaData:           r.MetaData,
		IsPublic:           r.IsPublic,
		LastMessage:        r.LastMessage,
		LastMessageUpdated: lmu,
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

	nowTimestamp := time.Now().Unix()
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
