package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

// Role is user role
type Role int

const (
	General Role = iota + 1
	Guest
	Operator
	External
	RoleEnd
)

type Users struct {
	Users []*User `json:"users"`
}

type User struct {
	ID               uint64         `json:"-" db:"id"`
	UserID           string         `json:"userId" db:"user_id,notnull"`
	Name             string         `json:"name" db:"name,notnull"`
	PictureURL       string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationURL   string         `json:"informationUrl,omitempty" db:"information_url"`
	UnreadCount      *uint64        `json:"unreadCount,omitempty" db:"unread_count,notnull"`
	MetaData         utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
	Role             *Role          `json:"role,omitempty" db:"role,notnull"`
	IsBot            *bool          `json:"isBot,omitempty" db:"is_bot,notnull"`
	IsPublic         *bool          `json:"isPublic,omitempty" db:"is_public,notnull"`
	IsCanBlock       *bool          `json:"isCanBlock,omitempty" db:"is_can_block,notnull"`
	IsShowUsers      *bool          `json:"isShowUsers,omitempty" db:"is_show_users,notnull"`
	Lang             string         `json:"lang,omitempty" db:"lang,notnull"`
	AccessToken      string         `json:"accessToken,omitempty" db:"-"`
	LastAccessRoomID string         `json:"lastAccessRoomId,omitempty" db:"last_access_room_id"`
	LastAccessed     int64          `json:"lastAccessed,omitempty" db:"last_accessed,notnull"`
	Created          int64          `json:"created,omitempty" db:"created,notnull"`
	Modified         int64          `json:"modified,omitempty" db:"modified,notnull"`
	Deleted          int64          `json:"-" db:"deleted,notnull"`

	Rooms   []*RoomForUser `json:"rooms,omitempty" db:"-"`
	Devices []*Device      `json:"devices,omitempty" db:"-"`
	Blocks  []string       `json:"blocks,omitempty" db:"-"`
}

type UserMini struct {
	RoomID       string `json:"roomId" db:"room_id"`
	UserID       string `json:"userId" db:"user_id"`
	Name         string `json:"name" db:"name"`
	PictureURL   string `json:"pictureUrl,omitempty" db:"picture_url"`
	Role         *Role  `json:"role,omitempty" db:"role"`
	IsShowUsers  *bool  `json:"isShowUsers,omitempty" db:"is_show_users"`
	LastAccessed int64  `json:"lastAccessed,omitempty" db:"last_accessed"`
}

func (um *UserMini) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")

	isShowUsers := false
	if um.IsShowUsers != nil {
		isShowUsers = *um.IsShowUsers
	}

	return json.Marshal(&struct {
		RoomID       string `json:"roomId"`
		UserID       string `json:"userId"`
		Name         string `json:"name"`
		PictureURL   string `json:"pictureUrl"`
		Role         *Role  `json:"role,omitempty"`
		IsShowUsers  bool   `json:"isShowUsers"`
		LastAccessed string `json:"lastAccessed"`
	}{
		RoomID:       um.RoomID,
		UserID:       um.UserID,
		Name:         um.Name,
		PictureURL:   um.PictureURL,
		Role: um.Role,
		IsShowUsers:  isShowUsers,
		LastAccessed: time.Unix(um.LastAccessed, 0).In(l).Format(time.RFC3339),
	})
}

type RoomForUser struct {
	// from room
	RoomID             string         `json:"roomId" db:"room_id"`
	UserID             string         `json:"userId" db:"user_id"`
	Name               string         `json:"name" db:"name"`
	PictureURL         string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationURL     string         `json:"informationUrl,omitempty" db:"information_url"`
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

	isBot := false
	if u.IsBot != nil {
		isBot = *u.IsBot
	}

	isPublic := false
	if u.IsPublic != nil {
		isPublic = *u.IsPublic
	}

	isCanBlock := false
	if u.IsCanBlock != nil {
		isCanBlock = *u.IsCanBlock
	}

	isShowUsers := false
	if u.IsShowUsers != nil {
		isShowUsers = *u.IsShowUsers
	}

	return json.Marshal(&struct {
		UserID           string         `json:"userId"`
		Name             string         `json:"name"`
		PictureURL       string         `json:"pictureUrl"`
		InformationURL   string         `json:"informationUrl"`
		UnreadCount      *uint64        `json:"unreadCount"`
		MetaData         utils.JSONText `json:"metaData"`
		Role             *Role          `json:"role"`
		IsBot            bool           `json:"isBot"`
		IsPublic         bool           `json:"isPublic"`
		IsCanBlock       bool           `json:"isCanBlock"`
		IsShowUsers      bool           `json:"isShowUsers"`
		Lang             string         `json:"lang"`
		AccessToken      string         `json:"accessToken"`
		LastAccessRoomID string         `json:"lastAccessRoomId"`
		LastAccessed     string         `json:"lastAccessed"`
		Created          string         `json:"created"`
		Modified         string         `json:"modified"`
		Rooms            []*RoomForUser `json:"rooms"`
		Devices          []*Device      `json:"devices"`
		Blocks           []string       `json:"blocks"`
	}{
		UserID:           u.UserID,
		Name:             u.Name,
		PictureURL:       u.PictureURL,
		InformationURL:   u.InformationURL,
		UnreadCount:      u.UnreadCount,
		MetaData:         u.MetaData,
		Role:             u.Role,
		IsBot:            isBot,
		IsPublic:         isPublic,
		IsCanBlock:       isCanBlock,
		IsShowUsers:      isShowUsers,
		Lang:             u.Lang,
		AccessToken:      u.AccessToken,
		LastAccessRoomID: u.LastAccessRoomID,
		LastAccessed:     time.Unix(u.LastAccessed, 0).In(l).Format(time.RFC3339),
		Created:          time.Unix(u.Created, 0).In(l).Format(time.RFC3339),
		Modified:         time.Unix(u.Modified, 0).In(l).Format(time.RFC3339),
		Rooms:            u.Rooms,
		Devices:          u.Devices,
		Blocks:           u.Blocks,
	})
}

func (rfu *RoomForUser) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	lmu := ""
	if rfu.LastMessageUpdated != 0 {
		lmu = time.Unix(rfu.LastMessageUpdated, 0).In(l).Format(time.RFC3339)
	}
	return json.Marshal(&struct {
		RoomID             string         `json:"roomId"`
		UserID             string         `json:"userId"`
		Name               string         `json:"name"`
		PictureURL         string         `json:"pictureUrl,omitempty"`
		InformationURL     string         `json:"informationUrl,omitempty"`
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
		RoomID:             rfu.RoomID,
		UserID:             rfu.UserID,
		Name:               rfu.Name,
		PictureURL:         rfu.PictureURL,
		InformationURL:     rfu.InformationURL,
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

func (u *User) IsValidPost() *ProblemDetail {
	if u.UserID != "" && !utils.IsValidID(u.UserID) {
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

	if len(u.UserID) > 36 {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. A string up to 36 symbols long.",
				},
			},
		}
	}

	if u.Name == "" {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
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

func (u *User) IsValidPut() *ProblemDetail {
	return nil
}

func (u *User) BeforePost() {
	if u.UserID == "" {
		u.UserID = utils.GenerateUUID()
	}

	if u.MetaData == nil {
		u.MetaData = []byte("{}")
	}

	if u.Role == nil {
		general := General
		u.Role = &general
	}

	if u.IsBot == nil {
		isBot := false
		u.IsBot = &isBot
	}

	if u.IsPublic == nil {
		isPublic := true
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

	u.Rooms = make([]*RoomForUser, 0)
	u.Devices = make([]*Device, 0)
	u.Blocks = make([]string, 0)
	nowTimestamp := time.Now().Unix()
	u.Created = nowTimestamp
	u.Modified = nowTimestamp
}

func (u *User) BeforePut(put *User) {
	if put.Name != "" {
		u.Name = put.Name
	}
	if put.PictureURL != "" {
		u.PictureURL = put.PictureURL
	}
	if put.InformationURL != "" {
		u.InformationURL = put.InformationURL
	}
	if put.UnreadCount != nil {
		u.UnreadCount = put.UnreadCount
	}
	if put.MetaData != nil {
		u.MetaData = put.MetaData
	}
	if put.IsBot != nil {
		u.IsBot = put.IsBot
	}
	if put.IsPublic != nil {
		u.IsPublic = put.IsPublic
	}
	if put.IsCanBlock != nil {
		u.IsCanBlock = put.IsCanBlock
	}
	if put.IsShowUsers != nil {
		u.IsShowUsers = put.IsShowUsers
	}
	if put.Lang != "" {
		u.Lang = put.Lang
	}
}

func (u *User) BeforeInsertGuest() {
	if u.UserID == "" {
		u.UserID = utils.GenerateUUID()
	}

	if u.MetaData == nil {
		u.MetaData = []byte("{}")
	}

	guest := Guest
	u.Role = &guest

	isBot := false
	u.IsBot = &isBot

	if u.IsPublic == nil {
		isPublic := true
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

	unreadCount := uint64(0)
	u.UnreadCount = &unreadCount

	u.Rooms = make([]*RoomForUser, 0)
	u.Devices = make([]*Device, 0)
	u.Blocks = make([]string, 0)
	nowTimestamp := time.Now().Unix()
	u.Created = nowTimestamp
	u.Modified = nowTimestamp
}
