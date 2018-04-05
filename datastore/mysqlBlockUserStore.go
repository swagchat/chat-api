package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return RdbInsertBlockUsers(blockUsers)
}

func (p *mysqlProvider) SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (p *mysqlProvider) SelectBlockUsersByUserId(userId string) ([]string, error) {
	return RdbSelectBlockUsersByUserId(userId)
}

func (p *mysqlProvider) DeleteBlockUser(userId string, blockUserIds []string) error {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
