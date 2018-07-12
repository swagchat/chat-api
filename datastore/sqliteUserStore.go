package datastore

import "github.com/swagchat/chat-api/model"

func (p *sqliteProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *sqliteProvider) InsertUser(user *model.User, opts ...interface{}) (*model.User, error) {
	return rdbInsertUser(p.database, user, opts...)
}

func (p *sqliteProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	return rdbSelectUser(p.database, userID, opts...)
}

func (p *sqliteProvider) SelectUserByUserIDAndAccessToken(userID, accessToken string) (*model.User, error) {
	return rdbSelectUserByUserIDAndAccessToken(p.database, userID, accessToken)
}

func (p *sqliteProvider) SelectUsers() ([]*model.User, error) {
	return rdbSelectUsers(p.database)
}

func (p *sqliteProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *sqliteProvider) UpdateUser(user *model.User) (*model.User, error) {
	return rdbUpdateUser(p.database, user)
}

func (p *sqliteProvider) UpdateUserDeleted(userID string) error {
	return rdbUpdateUserDeleted(p.database, userID)
}

func (p *sqliteProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.database, userID)
}
