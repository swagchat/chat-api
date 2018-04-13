package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createUserStore() {
	rdbCreateUserStore(p.database)
}

func (p *gcpSQLProvider) InsertUser(user *models.User) (*models.User, error) {
	return rdbInsertUser(p.database, user)
}

func (p *gcpSQLProvider) SelectUser(userID string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return rdbSelectUser(p.database, userID, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *gcpSQLProvider) SelectUserByUserIDAndAccessToken(userID, accessToken string) (*models.User, error) {
	return rdbSelectUserByUserIDAndAccessToken(p.database, userID, accessToken)
}

func (p *gcpSQLProvider) SelectUsers() ([]*models.User, error) {
	return rdbSelectUsers(p.database)
}

func (p *gcpSQLProvider) SelectUserIDsByUserIDs(userIDs []string) ([]string, error) {
	return rdbSelectUserIDsByUserIDs(p.database, userIDs)
}

func (p *gcpSQLProvider) UpdateUser(user *models.User) (*models.User, error) {
	return rdbUpdateUser(p.database, user)
}

func (p *gcpSQLProvider) UpdateUserDeleted(userID string) error {
	return rdbUpdateUserDeleted(p.database, userID)
}

func (p *gcpSQLProvider) SelectContacts(userID string) ([]*models.User, error) {
	return rdbSelectContacts(p.database, userID)
}
