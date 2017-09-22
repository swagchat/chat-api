package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

type Users struct {
	Users []*User `json:"users"`
}

type User struct {
	Id             uint64         `json:"-" db:"id"`
	UserId         string         `json:"userId" db:"user_id,notnull"`
	Name           string         `json:"name" db:"name,notnull"`
	PictureUrl     string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl string         `json:"informationUrl,omitempty" db:"information_url"`
	UnreadCount    *uint64        `json:"unreadCount,omitempty" db:"unread_count,notnull"`
	MetaData       utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
	IsPublic       *bool          `json:"isPublic,omitempty" db:"is_public,notnull"`
	IsCanBlock     *bool          `json:"isCanBlock,omitempty" db:"is_can_block,notnull"`
	IsShowUsers    *bool          `json:"isShowUsers,omitempty" db:"is_show_users,notnull"`
	AccessToken    string         `json:"accessToken,omitempty" db:"access_token"`
	Created        int64          `json:"created,omitempty" db:"created,notnull"`
	Modified       int64          `json:"modified,omitempty" db:"modified,notnull"`
	Deleted        int64          `json:"-" db:"deleted,notnull"`

	Rooms   []*RoomForUser `json:"rooms,omitempty" db:"-"`
	Devices []*Device      `json:"devices,omitempty" db:"-"`
	Blocks  []string       `json:"blocks,omitempty" db:"-"`
}

type UserMini struct {
	RoomId      string `json:"roomId" db:"room_id"`
	UserId      string `json:"userId" db:"user_id"`
	Name        string `json:"name" db:"name"`
	PictureUrl  string `json:"pictureUrl,omitempty" db:"picture_url"`
	IsShowUsers *bool  `json:"isShowUsers,omitempty" db:"is_show_users"`
}

type RoomForUser struct {
	// from room
	RoomId             string         `json:"roomId" db:"room_id"`
	UserId             string         `json:"userId" db:"user_id"`
	Name               string         `json:"name" db:"name"`
	PictureUrl         string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl     string         `json:"informationUrl,omitempty" db:"information_url"`
	MetaData           utils.JSONText `json:"metaData" db:"meta_data"`
	Type               *RoomType      `json:"type,omitempty" db:"type"`
	LastMessage        string         `json:"lastMessage" db:"last_message"`
	LastMessageUpdated int64          `json:"lastMessageUpdated" db:"last_message_updated"`
	IsCanLeft          *bool          `json:"isCanLeft,omitempty" db:"is_can_left,notnull"`
	Created            int64          `json:"created" db:"created"`
	Modified           int64          `json:"modified" db:"modified"`

	Users []*UserMini `json:"users" db:"-"`

	// from RoomUser
	RuUnreadCount int64          `json:"ruUnreadCount" db:"ru_unread_count"`
	RuMetaData    utils.JSONText `json:"ruMetaData" db:"ru_meta_data"`
	RuCreated     int64          `json:"ruCreated" db:"ru_created"`
	RuModified    int64          `json:"ruModified" db:"ru_modified"`
}

type UserUnreadCount struct {
	UnreadCount *uint64 `json:"unreadCount" db:"unread_count"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		UserId         string         `json:"userId"`
		Name           string         `json:"name"`
		PictureUrl     string         `json:"pictureUrl,omitempty"`
		InformationUrl string         `json:"informationUrl,omitempty"`
		UnreadCount    *uint64        `json:"unreadCount"`
		MetaData       utils.JSONText `json:"metaData"`
		IsPublic       *bool          `json:"isPublic,omitempty"`
		IsCanBlock     *bool          `json:"isCanBlock,omitempty"`
		IsShowUsers    *bool          `json:"isShowUsers,omitempty"`
		AccessToken    string         `json:"accessToken,omitempty"`
		Created        string         `json:"created"`
		Modified       string         `json:"modified"`
		Rooms          []*RoomForUser `json:"rooms,omitempty"`
		Devices        []*Device      `json:"devices,omitempty"`
		Blocks         []string       `json:"blocks,omitempty"`
	}{
		UserId:         u.UserId,
		Name:           u.Name,
		PictureUrl:     u.PictureUrl,
		InformationUrl: u.InformationUrl,
		UnreadCount:    u.UnreadCount,
		MetaData:       u.MetaData,
		IsPublic:       u.IsPublic,
		IsCanBlock:     u.IsCanBlock,
		IsShowUsers:    u.IsShowUsers,
		AccessToken:    u.AccessToken,
		Created:        time.Unix(u.Created, 0).In(l).Format(time.RFC3339),
		Modified:       time.Unix(u.Modified, 0).In(l).Format(time.RFC3339),
		Rooms:          u.Rooms,
		Devices:        u.Devices,
		Blocks:         u.Blocks,
	})
}

func (rfu *RoomForUser) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	lmu := ""
	if rfu.LastMessageUpdated != 0 {
		lmu = time.Unix(rfu.LastMessageUpdated, 0).In(l).Format(time.RFC3339)
	}
	return json.Marshal(&struct {
		RoomId             string         `json:"roomId"`
		UserId             string         `json:"userId"`
		Name               string         `json:"name"`
		PictureUrl         string         `json:"pictureUrl,omitempty"`
		InformationUrl     string         `json:"informationUrl,omitempty"`
		MetaData           utils.JSONText `json:"metaData"`
		Type               *RoomType      `json:"type,omitempty"`
		LastMessage        string         `json:"lastMessage"`
		LastMessageUpdated string         `json:"lastMessageUpdated"`
		IsCanLeft          *bool          `json:"isCanLeft,omitempty"`
		Created            string         `json:"created"`
		Modified           string         `json:"modified"`
		Users              []*UserMini    `json:"users"`
		RuUnreadCount      int64          `json:"ruUnreadCount"`
		RuMetaData         utils.JSONText `json:"ruMetaData"`
		RuCreated          string         `json:"ruCreated"`
		RuModified         string         `json:"ruModified"`
	}{
		RoomId:             rfu.RoomId,
		UserId:             rfu.UserId,
		Name:               rfu.Name,
		PictureUrl:         rfu.PictureUrl,
		InformationUrl:     rfu.InformationUrl,
		MetaData:           rfu.MetaData,
		Type:               rfu.Type,
		LastMessage:        rfu.LastMessage,
		LastMessageUpdated: lmu,
		IsCanLeft:          rfu.IsCanLeft,
		Created:            time.Unix(rfu.Created, 0).In(l).Format(time.RFC3339),
		Modified:           time.Unix(rfu.Modified, 0).In(l).Format(time.RFC3339),
		Users:              rfu.Users,
		RuUnreadCount:      rfu.RuUnreadCount,
		RuMetaData:         rfu.RuMetaData,
		RuCreated:          time.Unix(rfu.RuCreated, 0).In(l).Format(time.RFC3339),
		RuModified:         time.Unix(rfu.RuModified, 0).In(l).Format(time.RFC3339),
	})
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

	if u.IsPublic == nil {
		isPublic := false
		u.IsPublic = &isPublic
	}

	if u.IsCanBlock == nil {
		isCanBlock := true
		u.IsCanBlock = &isCanBlock
	}

	if u.IsShowUsers == nil {
		isShowUsers := true
		u.IsShowUsers = &isShowUsers
	}

	if u.UnreadCount == nil {
		unreadCount := uint64(0)
		u.UnreadCount = &unreadCount
	}

	nowTimestamp := time.Now().Unix()
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
	if put.IsPublic != nil {
		u.IsPublic = put.IsPublic
	}
	if put.IsCanBlock != nil {
		u.IsCanBlock = put.IsCanBlock
	}
}
