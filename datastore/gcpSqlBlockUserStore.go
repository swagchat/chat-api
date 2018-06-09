package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.database)
}

func (p *gcpSQLProvider) InsertBlockUsers(blockUsers []*models.BlockUser) error {
	return rdbInsertBlockUsers(p.database, blockUsers)
}

func (p *gcpSQLProvider) SelectBlockUser(userID, blockUserID string) (*models.BlockUser, error) {
	return rdbSelectBlockUser(p.database, userID, blockUserID)
}

func (p *gcpSQLProvider) SelectBlockUsersByUserID(userID string) ([]string, error) {
	return rdbSelectBlockUsersByUserID(p.database, userID)
}

func (p *gcpSQLProvider) DeleteBlockUser(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUser(p.database, userID, blockUserIDs)
}
