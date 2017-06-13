package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (provider GcpSqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) StoreResult {
	return RdbInsertBlockUsers(blockUsers)
}

func (provider GcpSqlProvider) SelectBlockUser(userId, blockUserId string) StoreResult {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (provider GcpSqlProvider) SelectBlockUsersByUserId(userId string) StoreResult {
	return RdbSelectBlockUsersByUserId(userId)
}

func (provider GcpSqlProvider) DeleteBlockUser(userId string, blockUserIds []string) StoreResult {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
