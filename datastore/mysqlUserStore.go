package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (p *mysqlProvider) InsertUser(user *models.User) (*models.User, error) {
	return RdbInsertUser(user)
}

func (p *mysqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *mysqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (p *mysqlProvider) SelectUsers() ([]*models.User, error) {
	return RdbSelectUsers()
}

func (p *mysqlProvider) SelectUserIdsByUserIds(userIds []string) ([]string, error) {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (p *mysqlProvider) UpdateUser(user *models.User) (*models.User, error) {
	return RdbUpdateUser(user)
}

func (p *mysqlProvider) UpdateUserDeleted(userId string) error {
	return RdbUpdateUserDeleted(userId)
}

func (p *mysqlProvider) SelectContacts(userId string) ([]*models.User, error) {
	return RdbSelectContacts(userId)
}
