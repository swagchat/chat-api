package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createUserStore() {
	rdbCreateUserStore(p.ctx, p.database)
}

func (p *gcpSQLProvider) InsertUser(user *model.User, opts ...InsertUserOption) error {
	return rdbInsertUser(p.ctx, p.database, user, opts...)
}

func (p *gcpSQLProvider) SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	return rdbSelectUsers(p.ctx, p.database, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	return rdbSelectUser(p.ctx, p.database, userID, opts...)
}

func (p *gcpSQLProvider) SelectCountUsers(opts ...SelectUsersOption) (int64, error) {
	return rdbSelectCountUsers(p.ctx, p.database, opts...)
}

func (p *gcpSQLProvider) SelectUserIDsOfUser(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsOfUser(p.ctx, p.database, userIDs)
}

func (p *gcpSQLProvider) UpdateUser(user *model.User) error {
	return rdbUpdateUser(p.ctx, p.database, user)
}

func (p *gcpSQLProvider) SelectContacts(userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	return rdbSelectContacts(p.ctx, p.database, userID, limit, offset, opts...)
}
