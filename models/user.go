package models

import "github.com/fairway-corp/swagchat-api/utils"

type Users struct {
	Users []*User `json:"users"`
}

type User struct {
	Id                   uint64         `json:"-" db:"id"`
	UserId               string         `json:"userId" db:"user_id"`
	Name                 string         `json:"name" db:"name"`
	PictureUrl           string         `json:"pictureUrl,omitempty" db:"picture_url"`
	InformationUrl       string         `json:"informationUrl,omitempty" db:"information_url"`
	UnreadCount          *int64         `json:"unreadCount" db:"unread_count"`
	MetaData             utils.JSONText `json:"metaData" db:"meta_data"`
	DeviceToken          *string        `json:"deviceToken,omitempty" db:"device_token"`
	NotificationDeviceId *string        `json:"-" db:"notification_device_id"`
	Created              int64          `json:"created" db:"created"`
	Modified             int64          `json:"modified" db:"modified"`
	Deleted              int64          `json:"-" db:"deleted"`

	Rooms []*RoomForUser `json:"rooms,omitempty" db:"-"`
}

type RoomForUser struct {
	// from room
	RoomId             string         `json:"roomId,omitempty" db:"room_id"`
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
