package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (provider MysqlProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (provider MysqlProvider) SelectUser(userId string, isWithRooms, isWithDevices bool) StoreChannel {
	return RdbSelectUser(userId, isWithRooms, isWithDevices)
}

func (provider MysqlProvider) SelectUsers() StoreChannel {
	return RdbSelectUsers()
}

func (provider MysqlProvider) SelectRoomsForUser(userId string) StoreChannel {
	return RdbSelectRoomsForUser(userId)
}

func (provider MysqlProvider) SelectUserIdsByUserIds(userIds []string) StoreChannel {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (provider MysqlProvider) UpdateUser(user *models.User) StoreChannel {
	return RdbUpdateUser(user)
}

func (provider MysqlProvider) UpdateUserDeleted(userId string) StoreChannel {
	return RdbUpdateUserDeleted(userId)
}

//func (provider MysqlProvider) UserSelectUserRooms(userId string) StoreChannel {
//	return RdbUserSelectUserRooms(userId)
//}

//func (provider MysqlProvider) UserUnreadCountUp(userId string) StoreChannel {
//	return RdbUserUnreadCountUp(userId)
//}
//
//func (provider MysqlProvider) UserUnreadCountRecalc(userId string) StoreChannel {
//	return RdbUserUnreadCountRecalc(userId)
//}
