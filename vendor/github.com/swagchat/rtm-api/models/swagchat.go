package models

import "github.com/swagchat/rtm-api/utils"

type Message struct {
	Id        uint64         `json:"id,omitempty" db:"id"`
	MessageId string         `json:"messageId,omitempty" db:"message_id"`
	RoomId    string         `json:"roomId" db:"room_id"`
	UserId    string         `json:"userId" db:"user_id"`
	Type      string         `json:"type" db:"type"`
	EventName string         `json:"eventName" db:"-"`
	Payload   utils.JSONText `json:"payload" db:"payload"`
	Created   int64          `json:"created,omitempty" db:"created"`
	Modified  int64          `json:"modified,omitempty" db:"modified"`
	Deleted   int64          `json:"-" db:"deleted"`
}

type PayloadText struct {
	Text string `json:"text"`
}

type PayloadImage struct {
	Mime         string `json:"mime"`
	SourceUrl    string `json:"sourceUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type PayloadUsers struct {
	Users []UserForRoom `json:"users"`
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
