package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *sqliteProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (p *sqliteProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (p *sqliteProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) StoreResult {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *sqliteProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (p *sqliteProvider) SelectUsers() StoreResult {
	return RdbSelectUsers()
}

func (p *sqliteProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (p *sqliteProvider) UpdateUser(user *models.User) StoreResult {
	return RdbUpdateUser(user)
}

func (p *sqliteProvider) UpdateUserDeleted(userId string) StoreResult {
	return RdbUpdateUserDeleted(userId)
}

func (p *sqliteProvider) SelectContacts(userId string) StoreResult {
	return RdbSelectContacts(userId)
}
