package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/utils"
)

const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
)

type Messages struct {
	Messages []*Message `json:"messages" db:"-"`
	AllCount int64      `json:"allCount" db:"all_count"`
}

type Message struct {
	ID        uint64         `json:"-" db:"id"`
	MessageID string         `json:"messageId" db:"message_id,notnull"`
	RoomID    string         `json:"roomId" db:"room_id,notnull"`
	UserID    string         `json:"userId" db:"user_id,notnull"`
	Type      string         `json:"type,omitempty" db:"type"`
	EventName string         `json:"eventName,omitempty" db:"-"`
	Payload   utils.JSONText `json:"payload" db:"payload"`
	Created   int64          `json:"created" db:"created,notnull"`
	Modified  int64          `json:"modified" db:"modified,notnull"`
	Deleted   int64          `json:"-" db:"deleted,notnull"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		MessageID string         `json:"messageId"`
		RoomID    string         `json:"roomId"`
		UserID    string         `json:"userId"`
		Type      string         `json:"type"`
		EventName string         `json:"eventName,omitempty"`
		Payload   utils.JSONText `json:"payload"`
		Created   string         `json:"created"`
		Modified  string         `json:"modified"`
	}{
		MessageID: m.MessageID,
		RoomID:    m.RoomID,
		UserID:    m.UserID,
		Type:      m.Type,
		EventName: m.EventName,
		Payload:   m.Payload,
		Created:   time.Unix(m.Created, 0).In(l).Format(time.RFC3339),
		Modified:  time.Unix(m.Modified, 0).In(l).Format(time.RFC3339),
	})
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
	Filename     string `json:"filename"`
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
	if m.MessageID != "" && !utils.IsValidID(m.MessageID) {
		return &ProblemDetail{
			Title:  "Request error",
			Status: http.StatusBadRequest,
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
			Title:  "Request error",
			Status: http.StatusBadRequest,
			InvalidParams: []InvalidParam{
				InvalidParam{
					Name:   "payload",
					Reason: "payload is empty.",
				},
			},
		}
	}

	if m.Type == MessageTypeText {
		var pt PayloadText
		json.Unmarshal(m.Payload, &pt)
		if pt.Text == "" {
			return &ProblemDetail{
				Title:  "Request error",
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Text type needs text.",
					},
				},
			}
		}
	}

	if m.Type == MessageTypeImage {
		var pi PayloadImage
		json.Unmarshal(m.Payload, &pi)
		if pi.Mime == "" {
			return &ProblemDetail{
				Title:  "Request error",
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Image type needs mime.",
					},
				},
			}
		}

		if pi.SourceUrl == "" {
			return &ProblemDetail{
				Title:  "Request error",
				Status: http.StatusBadRequest,
				InvalidParams: []InvalidParam{
					InvalidParam{
						Name:   "payload",
						Reason: "Image type needs sourceUrl.",
					},
				},
			}
		}
	}

	return nil
}

func (m *Message) BeforeSave() {
	if m.MessageID == "" {
		m.MessageID = utils.GenerateUUID()
	}

	nowTimestamp := time.Now().Unix()
	if m.Created == 0 {
		m.Created = nowTimestamp
	}
	m.Modified = nowTimestamp
}
