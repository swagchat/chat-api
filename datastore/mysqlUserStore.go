package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateUserStore() {
	RdbCreateUserStore(p.database)
}

func (p *mysqlProvider) InsertUser(user *models.User) (*models.User, error) {
	return RdbInsertUser(p.database, user)
}

func (p *mysqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return RdbSelectUser(p.database, userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *mysqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	return RdbSelectUserByUserIdAndAccessToken(p.database, userId, accessToken)
}

func (p *mysqlProvider) SelectUsers() ([]*models.User, error) {
	return RdbSelectUsers(p.database)
}

func (p *mysqlProvider) SelectUserIdsByUserIds(userIds []string) ([]string, error) {
	return RdbSelectUserIdsByUserIds(p.database, userIds)
}

func (p *mysqlProvider) UpdateUser(user *models.User) (*models.User, error) {
	return RdbUpdateUser(p.database, user)
}

func (p *mysqlProvider) UpdateUserDeleted(userId string) error {
	return RdbUpdateUserDeleted(p.database, userId)
}

func (p *mysqlProvider) SelectContacts(userId string) ([]*models.User, error) {
	return RdbSelectContacts(p.database, userId)
}
