package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createUserStore() {
	rdbCreateUserStore(p.ctx, p.database)
}

func (p *mysqlProvider) InsertUser(user *model.User, opts ...InsertUserOption) error {
	return rdbInsertUser(p.ctx, p.database, user, opts...)
}

func (p *mysqlProvider) SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	return rdbSelectUsers(p.ctx, p.database, limit, offset, opts...)
}

func (p *mysqlProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	return rdbSelectUser(p.ctx, p.database, userID, opts...)
}

func (p *mysqlProvider) SelectCountUsers(opts ...SelectUsersOption) (int64, error) {
	return rdbSelectCountUsers(p.ctx, p.database, opts...)
}

func (p *mysqlProvider) SelectUserIDsOfUser(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsOfUser(p.ctx, p.database, userIDs)
}

func (p *mysqlProvider) UpdateUser(user *model.User) error {
	return rdbUpdateUser(p.ctx, p.database, user)
}

func (p *mysqlProvider) SelectContacts(userID string, limit, offset int32, opts ...SelectContactsOption) ([]*model.User, error) {
	return rdbSelectContacts(p.ctx, p.database, userID, limit, offset, opts...)
}
