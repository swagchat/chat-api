package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.database)
}

func (p *sqliteProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return rdbInsertBlockUsers(p.database, blockUsers)
}

func (p *sqliteProvider) SelectBlockUser(userID, blockUserID string) (*models.BlockUser, error) {
	return rdbSelectBlockUser(p.database, userID, blockUserID)
}

func (p *sqliteProvider) SelectBlockUsersByUserID(userID string) ([]string, error) {
	return rdbSelectBlockUsersByUserID(p.database, userID)
}

func (p *sqliteProvider) DeleteBlockUser(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUser(p.database, userID, blockUserIDs)
}
