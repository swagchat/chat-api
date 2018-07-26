package datastore

import "github.com/swagchat/chat-api/model"

const (
	RoomIDAll = "ALL"
)

type selectWebhooksOptions struct {
	roomID      string
	triggerWord string
	roleID      int32
}

type SelectWebhooksOption func(*selectWebhooksOptions)

func SelectWebhooksOptionWithRoomID(roomID string) SelectWebhooksOption {
	return func(ops *selectWebhooksOptions) {
		ops.roomID = roomID
	}
}

func SelectWebhooksOptionWithTriggerWord(triggerWord string) SelectWebhooksOption {
	return func(ops *selectWebhooksOptions) {
		ops.triggerWord = triggerWord
	}
}

func SelectWebhooksOptionWithRole(roleID int32) SelectWebhooksOption {
	return func(ops *selectWebhooksOptions) {
		ops.roleID = roleID
	}
}

type webhookStore interface {
	createWebhookStore()

	SelectWebhooks(event model.WebhookEventType, opts ...SelectWebhooksOption) ([]*model.Webhook, error)
}
