package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.database)
}

func (p *gcpSQLProvider) InsertBlockUsers(blockUsers []*model.BlockUser) error {
	return rdbInsertBlockUsers(p.database, blockUsers)
}

func (p *gcpSQLProvider) SelectBlockUsers(userID string) ([]string, error) {
	return rdbSelectBlockUsers(p.database, userID)
}

func (p *gcpSQLProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.database, userID, blockUserID)
}

func (p *gcpSQLProvider) DeleteBlockUsers(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUsers(p.database, userID, blockUserIDs)
}
