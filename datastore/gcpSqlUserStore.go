package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (provider GcpSqlProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (provider GcpSqlProvider) SelectUser(userId string, isWithRooms, isWithDevices bool) StoreChannel {
	return RdbSelectUser(userId, isWithRooms, isWithDevices)
}

func (provider GcpSqlProvider) SelectUsers() StoreChannel {
	return RdbSelectUsers()
}

func (provider GcpSqlProvider) SelectRoomsForUser(userId string) StoreChannel {
	return RdbSelectRoomsForUser(userId)
}

func (provider GcpSqlProvider) SelectUserIdsByUserIds(userIds []string) StoreChannel {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (provider GcpSqlProvider) UpdateUser(user *models.User) StoreChannel {
	return RdbUpdateUser(user)
}

func (provider GcpSqlProvider) UpdateUserDeleted(userId string) StoreChannel {
	return RdbUpdateUserDeleted(userId)
}

//func (provider GcpSqlProvider) UserSelectUserRooms(userId string) StoreChannel {
//	return RdbUserSelectUserRooms(userId)
//}

//func (provider GcpSqlProvider) UserUnreadCountUp(userId string) StoreChannel {
//	return RdbUserUnreadCountUp(userId)
//}
//
//func (provider GcpSqlProvider) UserUnreadCountRecalc(userId string) StoreChannel {
//	return RdbUserUnreadCountRecalc(userId)
//}
