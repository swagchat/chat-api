package datastore

func (p *gcpSqlProvider) CreateBotStore() {
	RdbCreateBotStore()
}

//func (provider GcpSqlProvider) InsertUser(user *models.User) StoreResult {
//	return RdbInsertUser(user)
//}

func (p *gcpSqlProvider) SelectBot(userId string) StoreResult {
	return RdbSelectBot(userId)
}

//func (provider GcpSqlProvider) SelectUserByUserIdAndAccessToken(userId, accessToken string) StoreResult {
//	return RdbSelectUserByUserIdAndAccessToken(userId, accessToken)
//}
//
//func (provider GcpSqlProvider) SelectUsers() StoreResult {
//	return RdbSelectUsers()
//}
//
//func (provider GcpSqlProvider) SelectUserIdsByUserIds(userIds []string) StoreResult {
//	return RdbSelectUserIdsByUserIds(userIds)
//}
//
//func (provider GcpSqlProvider) UpdateUser(user *models.User) StoreResult {
//	return RdbUpdateUser(user)
//}
//
//func (provider GcpSqlProvider) UpdateUserDeleted(userId string) StoreResult {
//	return RdbUpdateUserDeleted(userId)
//}
//
//func (provider GcpSqlProvider) SelectContacts(userId string) StoreResult {
//	return RdbSelectContacts(userId)
//}
