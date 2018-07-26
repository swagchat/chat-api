package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.database)
}

func (p *sqliteProvider) InsertBlockUsers(blockUsers []*model.BlockUser) error {
	return rdbInsertBlockUsers(p.database, blockUsers)
}

func (p *sqliteProvider) SelectBlockUsers(userID string) ([]string, error) {
	return rdbSelectBlockUsers(p.database, userID)
}

func (p *sqliteProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.database, userID, blockUserID)
}

func (p *sqliteProvider) DeleteBlockUsers(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUsers(p.database, userID, blockUserIDs)
}
