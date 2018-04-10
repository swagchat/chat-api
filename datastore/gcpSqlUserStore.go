package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateUserStore() {
	RdbCreateUserStore(p.database)
}

func (p *gcpSqlProvider) InsertUser(user *models.User) (*models.User, error) {
	return RdbInsertUser(p.database, user)
}

func (p *gcpSqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) (*models.User, error) {
	return RdbSelectUser(p.database, userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *gcpSqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) (*models.User, error) {
	return RdbSelectUserByUserIdAndAccessToken(p.database, userId, accessToken)
}

func (p *gcpSqlProvider) SelectUsers() ([]*models.User, error) {
	return RdbSelectUsers(p.database)
}

func (p *gcpSqlProvider) SelectUserIdsByUserIds(userIds []string) ([]string, error) {
	return RdbSelectUserIdsByUserIds(p.database, userIds)
}

func (p *gcpSqlProvider) UpdateUser(user *models.User) (*models.User, error) {
	return RdbUpdateUser(p.database, user)
}

func (p *gcpSqlProvider) UpdateUserDeleted(userId string) error {
	return RdbUpdateUserDeleted(p.database, userId)
}

func (p *gcpSqlProvider) SelectContacts(userId string) ([]*models.User, error) {
	return RdbSelectContacts(p.database, userId)
}
