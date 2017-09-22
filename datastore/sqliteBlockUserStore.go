package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (p *sqliteProvider) InsertBlockUsers(blockUsers []*models.BlockUser) StoreResult {
	return RdbInsertBlockUsers(blockUsers)
}

func (p *sqliteProvider) SelectBlockUser(userId, blockUserId string) StoreResult {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (p *sqliteProvider) SelectBlockUsersByUserId(userId string) StoreResult {
	return RdbSelectBlockUsersByUserId(userId)
}

func (p *sqliteProvider) DeleteBlockUser(userId string, blockUserIds []string) StoreResult {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
