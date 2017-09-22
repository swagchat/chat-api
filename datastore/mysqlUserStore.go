package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (p *mysqlProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (p *mysqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) StoreResult {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *mysqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (p *mysqlProvider) SelectUsers() StoreResult {
	return RdbSelectUsers()
}

func (p *mysqlProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (p *mysqlProvider) UpdateUser(user *models.User) StoreResult {
	return RdbUpdateUser(user)
}

func (p *mysqlProvider) UpdateUserDeleted(userId string) StoreResult {
	return RdbUpdateUserDeleted(userId)
}

func (p *mysqlProvider) SelectContacts(userId string) StoreResult {
	return RdbSelectContacts(userId)
}
