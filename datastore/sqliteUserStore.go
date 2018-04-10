package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateUserStore() {
	RdbCreateUserStore(p.sqlitePath)
}

func (p *sqliteProvider) InsertUser(user *models.User) (*models.User, error) {
	return RdbInsertUser(p.sqlitePath, user)
}

func (p *sqliteProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return RdbSelectUser(p.sqlitePath, userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *sqliteProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	return RdbSelectUserByUserIdAndAccessToken(p.sqlitePath, userId, accessToken)
}

func (p *sqliteProvider) SelectUsers() ([]*models.User, error) {
	return RdbSelectUsers(p.sqlitePath)
}

func (p *sqliteProvider) SelectUserIdsByUserIds(userIds []string) ([]string, error) {
	return RdbSelectUserIdsByUserIds(p.sqlitePath, userIds)
}

func (p *sqliteProvider) UpdateUser(user *models.User) (*models.User, error) {
	return RdbUpdateUser(p.sqlitePath, user)
}

func (p *sqliteProvider) UpdateUserDeleted(userId string) error {
	return RdbUpdateUserDeleted(p.sqlitePath, userId)
}

func (p *sqliteProvider) SelectContacts(userId string) ([]*models.User, error) {
	return RdbSelectContacts(p.sqlitePath, userId)
}
