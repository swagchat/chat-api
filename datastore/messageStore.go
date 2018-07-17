package datastore

import "github.com/swagchat/chat-api/model"

type messageStore interface {
	createMessageStore()

	InsertMessage(message *model.Message) (string, error)
	SelectMessage(messageID string) (*model.Message, error)
	SelectMessages(roleIds []int32, roomID string, limit, offset int32, order string) ([]*model.Message, error)
	SelectCountMessagesByRoomID(roleIDs []int32, roomID string) (int64, error)
	UpdateMessage(message *model.Message) (*model.Message, error)
}
