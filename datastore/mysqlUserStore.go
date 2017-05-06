package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (provider MysqlProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (provider MysqlProvider) SelectUser(userId string, isWithRooms, isWithDevices bool) StoreResult {
	return RdbSelectUser(userId, isWithRooms, isWithDevices)
}

func (provider MysqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (provider MysqlProvider) SelectUsers() StoreResult {
	return RdbSelectUsers()
}

func (provider MysqlProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (provider MysqlProvider) UpdateUser(user *models.User) StoreResult {
	return RdbUpdateUser(user)
}

func (provider MysqlProvider) UpdateUserDeleted(userId string) StoreResult {
	return RdbUpdateUserDeleted(userId)
}
