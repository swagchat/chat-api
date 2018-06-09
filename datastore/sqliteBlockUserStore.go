package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return rdbInsertBlockUsers(p.sqlitePath, blockUsers)
}

func (p *sqliteProvider) SelectBlockUser(userID, blockUserID string) (*models.BlockUser, error) {
	return rdbSelectBlockUser(p.sqlitePath, userID, blockUserID)
}

func (p *sqliteProvider) SelectBlockUsersByUserID(userID string) ([]string, error) {
	return rdbSelectBlockUsersByUserID(p.sqlitePath, userID)
}

func (p *sqliteProvider) DeleteBlockUser(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUser(p.sqlitePath, userID, blockUserIDs)
}
