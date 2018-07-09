package datastore

import "github.com/swagchat/chat-api/model"

type blockUserStore interface {
	createBlockUserStore()

	InsertBlockUsers(blockUsers []*model.BlockUser) error
	SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error)
	SelectBlockUsersByUserID(userID string) ([]string, error)
	DeleteBlockUser(userID string, blockUserIDs []string) error
}
