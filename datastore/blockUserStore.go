package datastore

import "github.com/swagchat/chat-api/models"

type BlockUserStore interface {
	CreateBlockUserStore()

	InsertBlockUsers(blockUsers []*models.BlockUser) error
	SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error)
	SelectBlockUsersByUserId(userId string) ([]string, error)
	DeleteBlockUser(userId string, blockUserIds []string) error
}
