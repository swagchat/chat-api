package datastore

func (p *mysqlProvider) CreateBotStore() {
	RdbCreateBotStore()
}

//func (provider MysqlProvider) InsertUser(user *models.User) StoreResult {
//	return RdbInsertUser(user)
//}

func (p *mysqlProvider) SelectBot(userId string) StoreResult {
	return RdbSelectBot(userId)
}

//func (provider MysqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
//	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
//}
//
//func (provider MysqlProvider) SelectUsers() StoreResult {
//	return RdbSelectUsers()
//}
//
//func (provider MysqlProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
//	return RdbSelectUserIdsByUserIds(userIds)
//}
//
//func (provider MysqlProvider) UpdateUser(user *models.User) StoreResult {
//	return RdbUpdateUser(user)
//}
//
//func (provider MysqlProvider) UpdateUserDeleted(userId string) StoreResult {
//	return RdbUpdateUserDeleted(userId)
//}
//
//func (provider MysqlProvider) SelectContacts(userId string) StoreResult {
//	return RdbSelectContacts(userId)
//}
