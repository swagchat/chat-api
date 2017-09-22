package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *gcpSqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (p *gcpSqlProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (p *gcpSqlProvider) SelectUser(userId string, isWithRooms, isWithDevices, isWithBlocks bool) StoreResult {
	return RdbSelectUser(userId, isWithRooms, isWithDevices, isWithBlocks)
}

func (p *gcpSqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
}

func (p *gcpSqlProvider) SelectUsers() StoreResult {
	return RdbSelectUsers()
}

func (p *gcpSqlProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (p *gcpSqlProvider) UpdateUser(user *models.User) StoreResult {
	return RdbUpdateUser(user)
}

func (p *gcpSqlProvider) UpdateUserDeleted(userId string) StoreResult {
	return RdbUpdateUserDeleted(userId)
}

func (p *gcpSqlProvider) SelectContacts(userId string) StoreResult {
	return RdbSelectContacts(userId)
}
