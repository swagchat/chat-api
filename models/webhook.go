package models

import (
	"time"

	"encoding/json"
)

// Protocol is protocol
type WebhookProtocol int

const (
	WebhookProtocolHTTP WebhookProtocol = iota + 1
	WebhookProtocolGRPC
)

// WebhookEventType is webhook event
type WebhookEventType int

const (
	WebhookEventTypeRoom WebhookEventType = iota + 1
	WebhookEventTypeMessage
)

type Webhook struct {
	ID          uint64           `json:"-" db:"id"`
	WebhookID   string           `json:"webhookId" db:"webhook_id,notnull"`
	Event       WebhookEventType `json:"event" db:"event,notnull"`
	RoomID      string           `json:"roomId" db:"room_id,notnull"`
	RoleID      int32            `json:"roleId" db:"role_id,notnull"`
	TriggerWord string           `json:"triggerWord,omitempty" db:"trigger_word"`
	Protocol    WebhookProtocol  `json:"protocol" db:"protocol,notnull"`
	Endpoint    string           `json:"endpoint" db:"endpoint,notnull"`
	Token       string           `json:"token,omitempty" db:"token,notnull"`
	Created     int64            `json:"created" db:"created,notnull"`
	Modified    int64            `json:"modified" db:"modified,notnull"`
	Deleted     int64            `json:"-" db:"deleted,notnull"`
}

func (w *Webhook) MarshalJSON() ([]byte, error) {
	l, _ := time.LoadLocation("Etc/GMT")
	return json.Marshal(&struct {
		WebhookID   string           `json:"webhookId"`
		Event       WebhookEventType `json:"event"`
		RoomID      string           `json:"roomId"`
		RoleID      int32            `json:"roleId,omitempty"`
		TriggerWord string           `json:"triggerWord,omitempty"`
		Protocol    WebhookProtocol  `json:"protocol"`
		Endpoint    string           `json:"endpoint"`
		Token       string           `json:"token,omitempty"`
		Created     string           `json:"created"`
		Modified    string           `json:"modified"`
	}{
		WebhookID:   w.WebhookID,
		Event:       w.Event,
		RoomID:      w.RoomID,
		RoleID:      w.RoleID,
		TriggerWord: w.TriggerWord,
		Protocol:    w.Protocol,
		Endpoint:    w.Endpoint,
		Token:       w.Token,
		Created:     time.Unix(w.Created, 0).In(l).Format(time.RFC3339),
		Modified:    time.Unix(w.Modified, 0).In(l).Format(time.RFC3339),
	})
}
