package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createUserStore() {
	rdbCreateUserStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertUser(user *models.User, opts ...interface{}) (*models.User, error) {
	return rdbInsertUser(p.sqlitePath, user, opts...)
}

func (p *sqliteProvider) SelectUser(userID string, opts ...interface{}) (*models.User, error) {
	return rdbSelectUser(p.sqlitePath, userID, opts...)
}

func (p *sqliteProvider) SelectUserByUserIDAndAccessToken(userID, accessToken string) (*models.User, error) {
	return rdbSelectUserByUserIDAndAccessToken(p.sqlitePath, userID, accessToken)
}

func (p *sqliteProvider) SelectUsers() ([]*models.User, error) {
	return rdbSelectUsers(p.sqlitePath)
}

func (p *sqliteProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.sqlitePath, userIDs)
}

func (p *sqliteProvider) UpdateUser(user *models.User) (*models.User, error) {
	return rdbUpdateUser(p.sqlitePath, user)
}

func (p *sqliteProvider) UpdateUserDeleted(userID string) error {
	return rdbUpdateUserDeleted(p.sqlitePath, userID)
}

func (p *sqliteProvider) SelectContacts(userID string) ([]*models.User, error) {
	return rdbSelectContacts(p.sqlitePath, userID)
}
