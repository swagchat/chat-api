package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (provider SqliteProvider) InsertBlockUsers(blockUsers []*models.BlockUser) StoreResult {
	return RdbInsertBlockUsers(blockUsers)
}

func (provider SqliteProvider) SelectBlockUser(userId, blockUserId string) StoreResult {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (provider SqliteProvider) SelectBlockUsersByUserId(userId string) StoreResult {
	return RdbSelectBlockUsersByUserId(userId)
}

func (provider SqliteProvider) DeleteBlockUser(userId string, blockUserIds []string) StoreResult {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
