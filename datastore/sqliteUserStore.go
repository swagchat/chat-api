package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *sqliteProvider) InsertUser(user *model.User, opts ...InsertUserOption) error {
	return rdbInsertUser(p.database, user, opts...)
}

func (p *sqliteProvider) SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	return rdbSelectUsers(p.database, limit, offset, opts...)
}

func (p *sqliteProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	return rdbSelectUser(p.database, userID, opts...)
}

func (p *sqliteProvider) SelectCountUsers(opts ...SelectUsersOption) (int64, error) {
	return rdbSelectCountUsers(p.database, opts...)
}

func (p *sqliteProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *sqliteProvider) UpdateUser(user *model.User) error {
	return rdbUpdateUser(p.database, user)
}

func (p *sqliteProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.database, userID)
}
