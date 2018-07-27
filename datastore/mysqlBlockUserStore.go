package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*model.BlockUser) error {
	return rdbInsertBlockUsers(p.ctx, p.database, blockUsers)
}

func (p *mysqlProvider) SelectBlockUsers(userID string) ([]string, error) {
	return rdbSelectBlockUsers(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.ctx, p.database, userID, blockUserID)
}

func (p *mysqlProvider) DeleteBlockUsers(userID string, blockUserIDs []string) error {
	return rdbDeleteBlockUsers(p.ctx, p.database, userID, blockUserIDs)
}
