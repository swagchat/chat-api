package datastore

import "github.com/swagchat/chat-api/model"

type blockUserStore interface {
	createBlockUserStore()

	InsertBlockUsers(blockUsers []*model.BlockUser) error
	SelectBlockUsers(userID string) ([]string, error)
	SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error)
	DeleteBlockUsers(userID string, blockUserIDs []string) error
}
