package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.database)
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*model.BlockUser) error {
	return rdbInsertBlockUsers(p.database, blockUsers)
}

func (p *mysqlProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.database, userID, blockUserID)
}

func (p *mysqlProvider) SelectBlockUsersByUserID(userID string) ([]string, error) {
	return rdbSelectBlockUsersByUserID(p.database, userID)
}

func (p *mysqlProvider) DeleteBlockUser(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUser(p.database, userID, blockUserIDs)
}
