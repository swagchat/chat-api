package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (p *sqliteProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return RdbInsertBlockUsers(blockUsers)
}

func (p *sqliteProvider) SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (p *sqliteProvider) SelectBlockUsersByUserId(userId string) ([]string, error) {
	return RdbSelectBlockUsersByUserId(userId)
}

func (p *sqliteProvider) DeleteBlockUser(userId string, blockUserIds []string) error {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
