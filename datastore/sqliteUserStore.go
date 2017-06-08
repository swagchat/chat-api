package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (provider SqliteProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (provider SqliteProvider) SelectUser(userId string, isWithRooms, isWithDevices bool) StoreResult {
	return RdbSelectUser(userId, isWithRooms, isWithDevices)
}

func (provider SqliteProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (provider SqliteProvider) SelectUsers() StoreResult {
	return RdbSelectUsers()
}

func (provider SqliteProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (provider SqliteProvider) UpdateUser(user *models.User) StoreResult {
	return RdbUpdateUser(user)
}

func (provider SqliteProvider) UpdateUserDeleted(userId string) StoreResult {
	return RdbUpdateUserDeleted(userId)
}

func (provider SqliteProvider) SelectContacts(userId string) StoreResult {
	return RdbSelectContacts(userId)
}
