package datastore

import "github.com/swagchat/chat-api/models"

type MessageStore interface {
	CreateMessageStore()

	InsertMessage(message *models.Message) (string, error)
	SelectMessage(messageId string) (*models.Message, error)
	SelectMessages(roomId string, limit, offset int, order string) ([]*models.Message, error)
	SelectCountMessagesByRoomId(roomId string) (int64, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
}
