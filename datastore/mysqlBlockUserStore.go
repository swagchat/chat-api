package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.database)
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*model.BlockUser) error {
	return rdbInsertBlockUsers(p.database, blockUsers)
}

func (p *mysqlProvider) SelectBlockUsers(userID string) ([]string, error) {
	return rdbSelectBlockUsers(p.database, userID)
}

func (p *mysqlProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.database, userID, blockUserID)
}

func (p *mysqlProvider) DeleteBlockUsers(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUsers(p.database, userID, blockUserIDs)
}
