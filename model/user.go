package model

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
	ID               uint64         `json:"-" db:"id"`
	UserID           string         `json:"userId" db:"user_id,notnull"`
	Name             string         `json:"name" db:"name,notnull"`
	PictureURL       string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationURL   string         `json:"informationUrl,omitempty" db:"information_url"`
	UnreadCount      *uint64        `json:"unreadCount,omitempty" db:"unread_count,notnull"`
	MetaData         utils.JSONText `json:"metaData,omitempty" db:"meta_data"`
	Public           *bool          `json:"public,omitempty" db:"public,notnull"`
	CanBlock         *bool          `json:"canBlock,omitempty" db:"can_block,notnull"`
	Lang             string         `json:"lang,omitempty" db:"lang,notnull"`
	AccessToken      string         `json:"accessToken,omitempty" db:"-"`
	LastAccessRoomID string         `json:"lastAccessRoomId,omitempty" db:"last_access_room_id"`
	LastAccessed     int64          `json:"lastAccessed,omitempty" db:"last_accessed,notnull"`
	Created          int64          `json:"created,omitempty" db:"created,notnull"`
	Modified         int64          `json:"modified,omitempty" db:"modified,notnull"`
	Deleted          int64          `json:"-" db:"deleted,notnull"`

	Roles   []int32        `json:"roles,omitempty" db:"-"`
	Rooms   []*RoomForUser `json:"rooms,omitempty" db:"-"`
	Devices []*Device      `json:"devices,omitempty" db:"-"`
	Blocks  []string       `json:"blocks,omitempty" db:"-"`
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
	CanLeft            *bool          `json:"canLeft,omitempty" db:"can_left,notnull"`
	Created            int64          `json:"created" db:"created"`
	Modified           int64          `json:"modified" db:"modified"`

	Users []*UserForRoom `json:"users" db:"-"`

	// from RoomUser
	RuUnreadCount int64 `json:"ruUnreadCount" db:"ru_unread_count"`
}

type UserUnreadCount struct {
	UnreadCount *uint64 `json:"unreadCount" db:"unread_count"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")

	public := false
	if u.Public != nil {
		public = *u.Public
	}

	canBlock := false
	if u.CanBlock != nil {
		canBlock = *u.CanBlock
	}

	return json.Marshal(&struct {
		UserID           string         `json:"userId"`
		Name             string         `json:"name"`
		PictureURL       string         `json:"pictureUrl"`
		InformationURL   string         `json:"informationUrl"`
		UnreadCount      *uint64        `json:"unreadCount"`
		MetaData         utils.JSONText `json:"metaData"`
		Public           bool           `json:"public"`
		CanBlock         bool           `json:"canBlock"`
		Lang             string         `json:"lang"`
		AccessToken      string         `json:"accessToken"`
		LastAccessRoomID string         `json:"lastAccessRoomId"`
		LastAccessed     string         `json:"lastAccessed"`
		Created          string         `json:"created"`
		Modified         string         `json:"modified"`
		Roles            []int32        `json:"roles"`
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
		Public:           public,
		CanBlock:         canBlock,
		Lang:             u.Lang,
		AccessToken:      u.AccessToken,
		LastAccessRoomID: u.LastAccessRoomID,
		LastAccessed:     time.Unix(u.LastAccessed, 0).In(l).Format(time.RFC3339),
		Created:          time.Unix(u.Created, 0).In(l).Format(time.RFC3339),
		Modified:         time.Unix(u.Modified, 0).In(l).Format(time.RFC3339),
		Roles:            u.Roles,
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
		CanLeft            *bool          `json:"canLeft,omitempty"`
		Created            string         `json:"created"`
		Modified           string         `json:"modified"`
		Users              []*UserForRoom `json:"users"`
		RuUnreadCount      int64          `json:"ruUnreadCount"`
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
		CanLeft:            rfu.CanLeft,
		Created:            time.Unix(rfu.Created, 0).In(l).Format(time.RFC3339),
		Modified:           time.Unix(rfu.Modified, 0).In(l).Format(time.RFC3339),
		Users:              rfu.Users,
		RuUnreadCount:      rfu.RuUnreadCount,
	})
}

func (u *User) IsValidPost() *ProblemDetail {
	if u.UserID != "" && !utils.IsValidID(u.UserID) {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if len(u.UserID) > 36 {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "userId",
					Reason: "userId is invalid. A string up to 36 symbols long.",
				},
			},
			Status: http.StatusBadRequest,
		}
	}

	if u.Name == "" {
		return &ProblemDetail{
			Message: "Invalid params",
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "name",
					Reason: "name is required, but it's empty.",
				},
			},
			Status: http.StatusBadRequest,
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

	if u.Public == nil {
		public := true
		u.Public = &public
	}

	if u.CanBlock == nil {
		canBlock := true
		u.CanBlock = &canBlock
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
	if put.Public != nil {
		u.Public = put.Public
	}
	if put.CanBlock != nil {
		u.CanBlock = put.CanBlock
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

	if u.Public == nil {
		public := true
		u.Public = &public
	}

	if u.CanBlock == nil {
		canBlock := true
		u.CanBlock = &canBlock
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

func (u *User) IsRole(role int32) bool {
	if u.Roles == nil {
		return false
	}

	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}

	return false
}
