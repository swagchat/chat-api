package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *mysqlProvider) InsertUser(user *models.User) (*models.User, error) {
	return rdbInsertUser(p.database, user)
}

func (p *mysqlProvider) SelectUser(userID string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return rdbSelectUser(p.database, userID, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *mysqlProvider) SelectUserByUserIDAndAccessToken(userID, accessToken string) (*models.User, error) {
	return rdbSelectUserByUserIDAndAccessToken(p.database, userID, accessToken)
}

func (p *mysqlProvider) SelectUsers() ([]*models.User, error) {
	return rdbSelectUsers(p.database)
}

func (p *mysqlProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *mysqlProvider) UpdateUser(user *models.User) (*models.User, error) {
	return rdbUpdateUser(p.database, user)
}

func (p *mysqlProvider) UpdateUserDeleted(userID string) error {
	return rdbUpdateUserDeleted(p.database, userID)
}

func (p *mysqlProvider) SelectContacts(userID string) ([]*models.User, error) {
	return rdbSelectContacts(p.database, userID)
}
