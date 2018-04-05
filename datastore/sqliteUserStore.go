package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (p *sqliteProvider) InsertUser(user *models.User) (*models.User, error) {
	return RdbInsertUser(user)
}

func (p *sqliteProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *sqliteProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (p *sqliteProvider) SelectUsers() ([]*models.User, error) {
	return RdbSelectUsers()
}

func (p *sqliteProvider) SelectUserIdsByUserIds(userIds []string) ([]string, error) {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (p *sqliteProvider) UpdateUser(user *models.User) (*models.User, error) {
	return RdbUpdateUser(user)
}

func (p *sqliteProvider) UpdateUserDeleted(userId string) error {
	return RdbUpdateUserDeleted(userId)
}

func (p *sqliteProvider) SelectContacts(userId string) ([]*models.User, error) {
	return RdbSelectContacts(userId)
}
