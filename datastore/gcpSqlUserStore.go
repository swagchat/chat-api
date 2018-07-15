package datastore

import "github.com/swagchat/chat-api/model"

func (p *gcpSQLProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *gcpSQLProvider) InsertUser(user *model.User, opts ...interface{}) (*model.User, error) {
	return rdbInsertUser(p.database, user, opts...)
}

func (p *gcpSQLProvider) SelectUser(userID string, opts ...UserOption) (*model.User, error) {
	return rdbSelectUser(p.database, userID, opts...)
}

func (p *gcpSQLProvider) SelectUserByUserIDAndAccessToken(userID, accessToken string) (*model.User, error) {
	return rdbSelectUserByUserIDAndAccessToken(p.database, userID, accessToken)
}

func (p *gcpSQLProvider) SelectUsers() ([]*model.User, error) {
	return rdbSelectUsers(p.database)
}

func (p *gcpSQLProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *gcpSQLProvider) UpdateUser(user *model.User) (*model.User, error) {
	return rdbUpdateUser(p.database, user)
}

func (p *gcpSQLProvider) UpdateUserDeleted(userID string) error {
	return rdbUpdateUserDeleted(p.database, userID)
}

func (p *gcpSQLProvider) SelectContacts(userID string) ([]*model.User, error) {
	return rdbSelectContacts(p.database, userID)
}
