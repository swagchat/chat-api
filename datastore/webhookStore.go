package datastore

import "github.com/swagchat/chat-api/models"

const (
	RoomIDAll = "ALL"
)

type webhookOptions struct {
	roomID      string
	triggerWord string
	roleID      int32
}

type WebhookOption func(*webhookOptions)

func WithRoomID(roomID string) WebhookOption {
	return func(ops *webhookOptions) {
		ops.roomID = roomID
	}
}

func WithTriggerWord(triggerWord string) WebhookOption {
	return func(ops *webhookOptions) {
		ops.triggerWord = triggerWord
	}
}

func WithRole(roleID int32) WebhookOption {
	return func(ops *webhookOptions) {
		ops.roleID = roleID
	}
}

type webhookStore interface {
	createWebhookStore()

	SelectWebhooks(event models.WebhookEventType, opts ...WebhookOption) ([]*models.Webhook, error)
}
