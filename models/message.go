package models

import "github.com/fairway-corp/swagchat-api/utils"

type Messages struct {
	Messages []*Message `json:"messages" db:"-"`
	AllCount int64      `json:"allCount" db:"all_count"`
}

type Message struct {
	Id        uint64         `json:"id,omitempty" db:"id"`
	MessageId string         `json:"messageId,omitempty" db:"message_id"`
	RoomId    string         `json:"roomId" db:"room_id"`
	UserId    string         `json:"userId,omitempty" db:"user_id"`
	Type      string         `json:"type,omitempty" db:"type"`
	EventName string         `json:"eventName" db:"-"`
	Payload   utils.JSONText `json:"payload,omitempty" db:"payload"`
	Created   int64          `json:"created,omitempty" db:"created"`
	Modified  int64          `json:"modified,omitempty" db:"modified"`
	Deleted   int64          `json:"-" db:"deleted"`
}

type ResponseMessages struct {
	MessageIds []string `json:"messageIds"`
}

type PayloadText struct {
	Text string `json:"text"`
}

type PayloadImage struct {
	Mime         string `json:"mime"`
	SourceUrl    string `json:"sourceUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

type PayloadLocation struct {
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PayloadUsers struct {
	Users []string `json:"users"`
}
