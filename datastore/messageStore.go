package datastore

import "github.com/swagchat/chat-api/models"

type MessageStore interface {
	CreateMessageStore()

	InsertMessage(message *models.Message) StoreResult
	SelectMessage(messageId string) StoreResult
	SelectMessages(roomId string, limit, offset int, order string) StoreResult
	SelectCountMessagesByRoomId(roomId string) StoreResult
	UpdateMessage(message *models.Message) StoreResult
}
