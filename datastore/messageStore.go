package datastore

import "github.com/fairway-corp/swagchat-api/models"

type MessageStore interface {
	MessageCreateStore()

	MessageInsert(message *models.Message) StoreChannel
	MessageSelect(messageId string) StoreChannel
	MessageUpdate(message *models.Message) StoreChannel
	MessageSelectAll(roomId string, limit, offset int) StoreChannel
	MessageCount(roomId string) StoreChannel
}
