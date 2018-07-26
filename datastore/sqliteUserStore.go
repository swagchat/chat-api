package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createUserStore() {
	rdbCreateUserStore(p.ctx, p.database)
}

func (p *sqliteProvider) InsertUser(user *model.User, opts ...InsertUserOption) error {
	return rdbInsertUser(p.ctx, p.database, user, opts...)
}

func (p *sqliteProvider) SelectUsers(limit, offset int32, opts ...SelectUsersOption) ([]*model.User, error) {
	return rdbSelectUsers(p.ctx, p.database, limit, offset, opts...)
}

func (p *sqliteProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	return rdbSelectUser(p.ctx, p.database, userID, opts...)
}

func (p *sqliteProvider) SelectCountUsers(opts ...SelectUsersOption) (int64, error) {
	return rdbSelectCountUsers(p.ctx, p.database, opts...)
}

func (p *sqliteProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.ctx, p.database, userIDs)
}

func (p *sqliteProvider) UpdateUser(user *model.User) error {
	return rdbUpdateUser(p.ctx, p.database, user)
}

func (p *sqliteProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.ctx, p.database, userID)
}
