package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertBlockUsers(blockUsers []*model.BlockUser, opts ...InsertBlockUsersOption) error {
	return rdbInsertBlockUsers(p.ctx, p.database, blockUsers, opts...)
}

func (p *mysqlProvider) SelectBlockUsers(userID string) ([]*model.MiniUser, error) {
	return rdbSelectBlockUsers(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectBlockUserIDs(userID string) ([]string, error) {
	return rdbSelectBlockUserIDs(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectBlockedUsers(userID string) ([]*model.MiniUser, error) {
	return rdbSelectBlockedUsers(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectBlockedUserIDs(userID string) ([]string, error) {
	return rdbSelectBlockedUserIDs(p.ctx, p.database, userID)
}

func (p *mysqlProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.ctx, p.database, userID, blockUserID)
}

func (p *mysqlProvider) DeleteBlockUsers(opts ...DeleteBlockUsersOption) error {
	return rdbDeleteBlockUsers(p.ctx, p.database, opts...)
}
