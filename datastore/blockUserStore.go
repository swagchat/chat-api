package datastore

import "github.com/swagchat/chat-api/models"

type blockUserStore interface {
	createBlockUserStore()

	InsertBlockUsers(blockUsers []*models.BlockUser) error
	SelectBlockUser(userID, blockUserID string) (*models.BlockUser, error)
	SelectBlockUsersByUserID(userID string) ([]string, error)
	DeleteBlockUser(userID string, blockUserIDs []string) error
}
