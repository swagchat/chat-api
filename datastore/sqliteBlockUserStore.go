package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return RdbInsertBlockUsers(p.sqlitePath, blockUsers)
}

func (p *sqliteProvider) SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	return RdbSelectBlockUser(p.sqlitePath, userId, blockUserId)
}

func (p *sqliteProvider) SelectBlockUsersByUserId(userId string) ([]string, error) {
	return RdbSelectBlockUsersByUserId(p.sqlitePath, userId)
}

func (p *sqliteProvider) DeleteBlockUser(userId string, blockUserIds []string) error {
	return RdbDeleteBlockUser(p.sqlitePath, userId, blockUserIds)
}
