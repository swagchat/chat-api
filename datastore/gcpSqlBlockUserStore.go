package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (p *gcpSqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) StoreResult {
	return RdbInsertBlockUsers(blockUsers)
}

func (p *gcpSqlProvider) SelectBlockUser(userId, blockUserId string) StoreResult {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (p *gcpSqlProvider) SelectBlockUsersByUserId(userId string) StoreResult {
	return RdbSelectBlockUsersByUserId(userId)
}

func (p *gcpSqlProvider) DeleteBlockUser(userId string, blockUserIds []string) StoreResult {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
