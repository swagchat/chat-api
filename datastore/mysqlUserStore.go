package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *mysqlProvider) InsertUser(user *model.User, opts ...UserOption) error {
	return rdbInsertUser(p.database, user, opts...)
}

func (p *mysqlProvider) SelectUsers(limit, offset int32, opts ...UserOption) ([]*model.User, error) {
	return rdbSelectUsers(p.database, limit, offset, opts...)
}

func (p *mysqlProvider) SelectUser(userID string, opts ...UserOption) (*model.User, error) {
	return rdbSelectUser(p.database, userID, opts...)
}

func (p *mysqlProvider) SelectCountUsers(opts ...UserOption) (int64, error) {
	return rdbSelectCountUsers(p.database, opts...)
}

func (p *mysqlProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *mysqlProvider) UpdateUser(user *model.User) error {
	return rdbUpdateUser(p.database, user)
}

func (p *mysqlProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.database, userID)
}
