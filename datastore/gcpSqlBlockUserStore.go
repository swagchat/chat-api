package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateBlockUserStore() {
	RdbCreateBlockUserStore(p.database)
}

func (p *gcpSqlProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return RdbInsertBlockUsers(p.database, blockUsers)
}

func (p *gcpSqlProvider) SelectBlockUser(userId, blockUserId string) (*models.BlockUser, error) {
	return RdbSelectBlockUser(p.database, userId, blockUserId)
}

func (p *gcpSqlProvider) SelectBlockUsersByUserId(userId string) ([]string, error) {
	return RdbSelectBlockUsersByUserId(p.database, userId)
}

func (p *gcpSqlProvider) DeleteBlockUser(userId string, blockUserIds []string) error {
	return RdbDeleteBlockUser(p.database, userId, blockUserIds)
}
