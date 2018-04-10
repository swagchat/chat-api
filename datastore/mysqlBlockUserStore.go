package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore(p.database)
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return RdbInsertBlockUsers(p.database, blockUsers)
}

func (p *mysqlProvider) SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	return RdbSelectBlockUser(p.database, userId, blockUserId)
}

func (p *mysqlProvider) SelectBlockUsersByUserId(userId string) ([]string, error) {
	return RdbSelectBlockUsersByUserId(p.database, userId)
}

func (p *mysqlProvider) DeleteBlockUser(userId string, blockUserIds []string) error {
	return RdbDeleteBlockUser(p.database, userId, blockUserIds)
}
