package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore()
}

func (p *gcpSqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return RdbInsertBlockUsers(blockUsers)
}

func (p *gcpSqlProvider) SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	return RdbSelectBlockUser(userId, blockUserId)
}

func (p *gcpSqlProvider) SelectBlockUsersByUserId(userId string) ([]string, error) {
	return RdbSelectBlockUsersByUserId(userId)
}

func (p *gcpSqlProvider) DeleteBlockUser(userId string, blockUserIds []string) error {
	return RdbDeleteBlockUser(userId, blockUserIds)
}
