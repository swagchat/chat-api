package model

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/swagchat/chat-api/config"
	"github.com/swagchat/chat-api/utils"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

const (
	MessageTypeText           = "text"
	MessageTypeImage          = "image"
	MessageTypeFile           = "file"
	MessageTypeIndicatorStart = "indicator-start"
	MessageTypeIndicatorEnd   = "indicator-end"
	MessageTypeUpdateRoomUser = "updateRoomUser"

	EventNameMessage = "message"
)

type Message struct {
	scpb.Message
	Payload JSONText `json:"payload" db:"payload"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		MessageID        string   `json:"messageId"`
		RoomID           string   `json:"roomId"`
		UserID           string   `json:"userId"`
		Type             string   `json:"type"`
		Payload          JSONText `json:"payload"`
		Role             int32    `json:"role"`
		CreatedTimestamp int64    `json:"createdTimestamp"`
		Created          string   `json:"created"`
		Modified         string   `json:"modified"`
	}{
		MessageID:        m.MessageID,
		RoomID:           m.RoomID,
		UserID:           m.UserID,
		Type:             m.Type,
		Payload:          m.Payload,
		Role:             m.Role,
		CreatedTimestamp: m.CreatedTimestamp,
		Created:          time.Unix(m.CreatedTimestamp, 0).In(l).Format(time.RFC3339),
		Modified:         time.Unix(m.ModifiedTimestamp, 0).In(l).Format(time.RFC3339),
	})
}

func (m *Message) ConvertToPbMessage() *scpb.Message {
	pbMessage := &scpb.Message{}
	pbMessage.MessageID = m.MessageID
	pbMessage.RoomID = m.RoomID
	pbMessage.UserID = m.UserID
	pbMessage.Type = m.Type
	pbMessage.Payload = m.Payload
	pbMessage.Role = m.Role
	pbMessage.Created = m.Created
	pbMessage.Modified = m.Modified
	pbMessage.UserIDs = m.UserIDs
	return pbMessage
}

func (m *Message) ConvertToSendMessageRequest() *SendMessageRequest {
	req := &SendMessageRequest{}
	req.MessageID = &m.MessageID
	req.RoomID = &m.RoomID
	req.UserID = &m.UserID
	req.Type = &m.Type
	req.Payload = m.Payload
	req.Role = &m.Role
	return req
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

type PayloadUsers struct {
	Users []string `json:"users"`
}

type SendMessageRequest struct {
	scpb.SendMessageRequest
	Payload JSONText `json:"payload" db:"payload"`
}

func (m *SendMessageRequest) Validate() *ErrorResponse {
	if m.MessageID != nil && *m.MessageID != "" && !isValidID(*m.MessageID) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "messageId",
				Reason: "messageId is invalid. Available characters are alphabets, numbers and hyphens.",
			},
		}
		return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if m.RoomID == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roomId",
				Reason: "roomId is empty.",
			},
		}
		return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if m.UserID == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "userId",
				Reason: "userId is empty.",
			},
		}
		return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if m.Type == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "type",
				Reason: "type is empty.",
			},
		}
		return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if m.Payload == nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "payload",
				Reason: "payload is empty.",
			},
		}
		return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if !isJSON(m.Payload.String()) {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "payload",
				Reason: "payload is not json format.",
			},
		}
		return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
	}

	if *m.Type == MessageTypeText {
		var pt PayloadText
		json.Unmarshal(m.Payload, &pt)
		if pt.Text == "" {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "payload",
					Reason: "Text type needs text.",
				},
			}
			return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}

	if *m.Type == MessageTypeImage {
		var pi PayloadImage
		json.Unmarshal(m.Payload, &pi)
		if pi.Mime == "" {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "payload",
					Reason: "Image type needs mime.",
				},
			}
			return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}

		if pi.SourceUrl == "" {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "payload",
					Reason: "Image type needs sourceUrl.",
				},
			}
			return NewErrorResponse("Failed to create a message.", http.StatusBadRequest, WithInvalidParams(invalidParams))
		}
	}

	return nil
}

func (cmr *SendMessageRequest) GenerateMessage() *Message {
	m := &Message{}

	if cmr.MessageID == nil || *cmr.MessageID == "" {
		m.MessageID = utils.GenerateUUID()
	} else {
		m.MessageID = *cmr.MessageID
	}

	m.RoomID = *cmr.RoomID
	m.UserID = *cmr.UserID
	m.Type = *cmr.Type
	m.Payload = cmr.Payload

	if cmr.Role == nil {
		m.Role = config.RoleGeneral
	} else {
		m.Role = *cmr.Role
	}

	nowTimestamp := time.Now().Unix()
	m.CreatedTimestamp = nowTimestamp
	m.ModifiedTimestamp = nowTimestamp
	m.DeletedTimestamp = 0

	return m
}
