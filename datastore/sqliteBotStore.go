package datastore

func (p *sqliteProvider) CreateBotStore() {
	RdbCreateBotStore()
}

//func (provider SqliteProvider) InsertUser(user *models.User) StoreResult {
//	return RdbInsertUser(user)
//}

func (p *sqliteProvider) SelectBot(userId string) StoreResult {
	return RdbSelectBot(userId)
}

//func (provider SqliteProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
//	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
//}
//
//func (provider SqliteProvider) SelectUsers() StoreResult {
//	return RdbSelectUsers()
//}
//
//func (provider SqliteProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
//	return RdbSelectUserIdsByUserIds(userIds)
//}
//
//func (provider SqliteProvider) UpdateUser(user *models.User) StoreResult {
//	return RdbUpdateUser(user)
//}
//
//func (provider SqliteProvider) UpdateUserDeleted(userId string) StoreResult {
//	return RdbUpdateUserDeleted(userId)
//}
//
//func (provider SqliteProvider) SelectContacts(userId string) StoreResult {
//	return RdbSelectContacts(userId)
//}
