package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createBlockUserStore() {
	rdbCreateBlockUserStore(p.ctx, p.database)
}

func (p *gcpSQLProvider) InsertBlockUsers(blockUsers []*model.BlockUser, opts ...InsertBlockUsersOption) error {
	return rdbInsertBlockUsers(p.ctx, p.database, blockUsers, opts...)
}

func (p *gcpSQLProvider) SelectBlockUsers(userID string) ([]*model.MiniUser, error) {
	return rdbSelectBlockUsers(p.ctx, p.database, userID)
}

func (p *gcpSQLProvider) SelectBlockUserIDs(userID string) ([]string, error) {
	return rdbSelectBlockUserIDs(p.ctx, p.database, userID)
}

func (p *gcpSQLProvider) SelectBlockedUsers(userID string) ([]*model.MiniUser, error) {
	return rdbSelectBlockedUsers(p.ctx, p.database, userID)
}

func (p *gcpSQLProvider) SelectBlockedUserIDs(userID string) ([]string, error) {
	return rdbSelectBlockedUserIDs(p.ctx, p.database, userID)
}

func (p *gcpSQLProvider) SelectBlockUser(userID, blockUserID string) (*model.BlockUser, error) {
	return rdbSelectBlockUser(p.ctx, p.database, userID, blockUserID)
}

func (p *gcpSQLProvider) DeleteBlockUsers(opts ...DeleteBlockUsersOption) error {
	return rdbDeleteBlockUsers(p.ctx, p.database, opts...)
}
