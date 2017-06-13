package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (provider MysqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) StoreResult {
	return RdbInsertBlockUsers(blockUsers)
}

func (provider MysqlProvider) SelectBlockUser(userId, blockUserId string) StoreResult {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (provider MysqlProvider) SelectBlockUsersByUserId(userId string) StoreResult {
	return RdbSelectBlockUsersByUserId(userId)
}

func (provider MysqlProvider) DeleteBlockUser(userId string, blockUserIds []string) StoreResult {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
