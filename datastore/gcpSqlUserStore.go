package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *gcpSQLProvider) InsertUser(user *model.User, opts ...UserOption) error {
	return rdbInsertUser(p.database, user, opts...)
}

func (p *gcpSQLProvider) SelectUsers(limit, offset int32, opts ...UserOption) ([]*model.User, error) {
	return rdbSelectUsers(p.database, limit, offset, opts...)
}

func (p *gcpSQLProvider) SelectUser(userID string, opts ...UserOption) (*model.User, error) {
	return rdbSelectUser(p.database, userID, opts...)
}

func (p *gcpSQLProvider) SelectCountUsers(opts ...UserOption) (int64, error) {
	return rdbSelectCountUsers(p.database, opts...)
}

func (p *gcpSQLProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *gcpSQLProvider) UpdateUser(user *model.User) error {
	return rdbUpdateUser(p.database, user)
}

func (p *gcpSQLProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.database, userID)
}
