package datastore

import "github.com/swagchat/chat-api/models"

type BlockUserStore interface {
	CreateBlockUserStore()

	InsertBlockUsers(blockUsers []*models.BlockUser) StoreResult
	SelectBlockUser(userId, blockUserId string) StoreResult
	SelectBlockUsersByUserId(userId string) StoreResult
	DeleteBlockUser(userId string, blockUserIds []string) StoreResult
}
