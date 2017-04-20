package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) UserCreateStore() {
	RdbUserCreateStore()
}

func (provider SqliteProvider) UserInsert(user *models.User) StoreChannel {
	return RdbUserInsert(user)
}

func (provider SqliteProvider) UserSelect(userId string, isWithRooms, isWithDevices bool) StoreChannel {
	return RdbUserSelect(userId, isWithRooms, isWithDevices)
}

func (provider SqliteProvider) UserUpdate(user *models.User) StoreChannel {
	return RdbUserUpdate(user)
}

func (provider SqliteProvider) UserUpdateDeleted(userId string) StoreChannel {
	return RdbUserUpdateDeleted(userId)
}

func (provider SqliteProvider) UserSelectAll() StoreChannel {
	return RdbUserSelectAll()
}

func (provider SqliteProvider) UserSelectRoomsForUser(userId string) StoreChannel {
	return RdbUserSelectRoomsForUser(userId)
}

//func (provider SqliteProvider) UserSelectUserRooms(userId string) StoreChannel {
//	return RdbUserSelectUserRooms(userId)
//}

func (provider SqliteProvider) UserUnreadCountUp(userId string) StoreChannel {
	return RdbUserUnreadCountUp(userId)
}

func (provider SqliteProvider) UserUnreadCountRecalc(userId string) StoreChannel {
	return RdbUserUnreadCountRecalc(userId)
}

func (provider SqliteProvider) UserUserIdsSelectByUserIds(userIds []string) StoreChannel {
	return RdbUserUserIdsSelectByUserIds(userIds)
}
