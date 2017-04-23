package datastore

import "github.com/fairway-corp/swagchat-api/models"

type MessageStore interface {
	CreateMessageStore()

	InsertMessage(message *models.Message) StoreChannel
	SelectMessage(messageId string) StoreChannel
	SelectMessages(roomId string, limit, offset int) StoreChannel
	SelectCountMessagesByRoomId(roomId string) StoreChannel
	UpdateMessage(message *models.Message) StoreChannel
}
