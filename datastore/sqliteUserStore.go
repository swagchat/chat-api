package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateUserStore() {
	RdbCreateUserStore()
}

func (provider SqliteProvider) InsertUser(user *models.User) StoreResult {
	return RdbInsertUser(user)
}

func (provider SqliteProvider) SelectUser(userId string, isWithRooms, isWithDevices bool) StoreChannel {
	return RdbSelectUser(userId, isWithRooms, isWithDevices)
}

func (provider SqliteProvider) SelectUsers() StoreChannel {
	return RdbSelectUsers()
}

func (provider SqliteProvider) SelectRoomsForUser(userId string) StoreChannel {
	return RdbSelectRoomsForUser(userId)
}

func (provider SqliteProvider) SelectUserIdsByUserIds(userIds []string) StoreChannel {
	return RdbSelectUserIdsByUserIds(userIds)
}

func (provider SqliteProvider) UpdateUser(user *models.User) StoreChannel {
	return RdbUpdateUser(user)
}

func (provider SqliteProvider) UpdateUserDeleted(userId string) StoreChannel {
	return RdbUpdateUserDeleted(userId)
}

//func (provider SqliteProvider) UserSelectUserRooms(userId string) StoreChannel {
//	return RdbUserSelectUserRooms(userId)
//}

//func (provider SqliteProvider) UserUnreadCountUp(userId string) StoreChannel {
//	return RdbUserUnreadCountUp(userId)
//}
//
//func (provider SqliteProvider) UserUnreadCountRecalc(userId string) StoreChannel {
//	return RdbUserUnreadCountRecalc(userId)
//}
