package datastore

import "github.com/swagchat/chat-api/model"

func (p *mysqlProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *mysqlProvider) InsertUser(user *model.User, opts ...interface{}) (*model.User, error) {
	return rdbInsertUser(p.database, user, opts...)
}

func (p *mysqlProvider) SelectUser(userID string, opts ...SelectUserOption) (*model.User, error) {
	return rdbSelectUser(p.database, userID, opts...)
}

func (p *mysqlProvider) SelectUserByUserIDAndAccessToken(userID, accessToken string) (*model.User, error) {
	return rdbSelectUserByUserIDAndAccessToken(p.database, userID, accessToken)
}

func (p *mysqlProvider) SelectUsers() ([]*model.User, error) {
	return rdbSelectUsers(p.database)
}

func (p *mysqlProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *mysqlProvider) UpdateUser(user *model.User) (*model.User, error) {
	return rdbUpdateUser(p.database, user)
}

func (p *mysqlProvider) UpdateUserDeleted(userID string) error {
	return rdbUpdateUserDeleted(p.database, userID)
}

func (p *mysqlProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.database, userID)
}
