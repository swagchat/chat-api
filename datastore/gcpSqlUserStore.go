package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (p *gcpSqlProvider) InsertUser(user *models.User) (*models.User, error) {
	return RdbInsertUser(user)
}

func (p *gcpSqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *gcpSqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (p *gcpSqlProvider) SelectUsers() ([]*models.User, error) {
	return RdbSelectUsers()
}

func (p *gcpSqlProvider) SelectUserIdsByUserIds(userIds []string) ([]string, error) {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (p *gcpSqlProvider) UpdateUser(user *models.User) (*models.User, error) {
	return RdbUpdateUser(user)
}

func (p *gcpSqlProvider) UpdateUserDeleted(userId string) error {
	return RdbUpdateUserDeleted(userId)
}

func (p *gcpSqlProvider) SelectContacts(userId string) ([]*models.User, error) {
	return RdbSelectContacts(userId)
}
