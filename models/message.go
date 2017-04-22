package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fairway-corp/swagchat-api/utils"
)

const (
	MESSAGE_TYPE_TEXT  = "text"
	MESSAGE_TYPE_IMAGE = "image"
)

type Messages struct {
	Messages []*Message `json:"messages" db:"-"`
	AllCount int64      `json:"allCount" db:"all_count"`
}

type Message struct {
	Id        uint64         `json:"-" db:"id"`
	MessageId string         `json:"messageId,omitempty" db:"message_id"`
	RoomId    string         `json:"roomId" db:"room_id"`
	UserId    string         `json:"userId,omitempty" db:"user_id"`
	Type      string         `json:"type,omitempty" db:"type"`
	EventName string         `json:"eventName,omitempty" db:"-"`
	Payload   utils.JSONText `json:"payload,omitempty" db:"payload"`
	Created   int64          `json:"created,omitempty" db:"created"`
	Modified  int64          `json:"modified,omitempty" db:"modified"`
	Deleted   int64          `json:"-" db:"deleted"`
}

type ResponseMessages struct {
	MessageIds []string         `json:"messageIds,omitempty"`
	Errors     []*ProblemDetail `json:"errors,omitempty"`
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

func (m *Message) IsValid() *ProblemDetail {
	if m.MessageId != "" && !utils.IsValidId(m.MessageId) {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create message item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "messageId",
					Reason: "messageId is invalid. Available characters are alphabets, numbers and hyphens.",
				},
			},
		}
	}

	if m.Payload == nil {
		return &ProblemDetail{
			Title:     "Request parameter error. (Create message item)",
			Status:    http.StatusBadRequest,
			ErrorName: ERROR_NAME_INVALID_PARAM,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "payload",
					Reason: "payload is empty.",
				},
			},
		}
	}

	if m.Type == MESSAGE_TYPE_TEXT {
		var pt PayloadText
		json.Unmarshal(m.Payload, &pt)
		if pt.Text == "" {
			return &ProblemDetail{
				Title:     "Request parameter error. (Create message item)",
				Status:    http.StatusBadRequest,
				ErrorName: ERROR_NAME_INVALID_PARAM,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Text type needs text.",
					},
				},
			}
		}
	}

	if m.Type == MESSAGE_TYPE_IMAGE {
		var pi PayloadImage
		json.Unmarshal(m.Payload, &pi)
		if pi.Mime == "" || pi.SourceUrl == "" {
			return &ProblemDetail{
				Title:     "Request parameter error. (Create message item)",
				Status:    http.StatusBadRequest,
				ErrorName: ERROR_NAME_INVALID_PARAM,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Image type needs mime and sourceUrl.",
					},
				},
			}
		}
	}

	return nil
}

func (m *Message) BeforeSave() {
	if m.MessageId == "" {
		m.MessageId = utils.CreateUuid()
	}

	nowDatetime := time.Now().UnixNano()
	if m.Created == 0 {
		m.Created = nowDatetime
	}
	m.Modified = nowDatetime
}
