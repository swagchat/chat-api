package datastore

import "github.com/swagchat/chat-api/models"

type messageStore interface {
	createMessageStore()

	InsertMessage(message *models.Message) (string, error)
	SelectMessage(messageID string) (*models.Message, error)
	SelectMessages(roleIds []models.Role, roomID string, limit, offset int, order string) ([]*models.Message, error)
	SelectCountMessagesByRoomID(roleIDs []models.Role, roomID string) (int64, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
}
