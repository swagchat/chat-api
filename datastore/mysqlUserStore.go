package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) UserCreateStore() {
	RdbUserCreateStore()
}

func (provider MysqlProvider) UserInsert(user *models.User) StoreChannel {
	return RdbUserInsert(user)
}

func (provider MysqlProvider) UserSelect(userId string, isWithRooms, isWithDevices bool) StoreChannel {
	return RdbUserSelect(userId, isWithRooms, isWithDevices)
}

func (provider MysqlProvider) UserUpdate(user *models.User) StoreChannel {
	return RdbUserUpdate(user)
}

func (provider MysqlProvider) UserUpdateDeleted(userId string) StoreChannel {
	return RdbUserUpdateDeleted(userId)
}

func (provider MysqlProvider) UserSelectAll() StoreChannel {
	return RdbUserSelectAll()
}

func (provider MysqlProvider) UserSelectRoomsForUser(userId string) StoreChannel {
	return RdbUserSelectRoomsForUser(userId)
}

//func (provider MysqlProvider) UserSelectUserRooms(userId string) StoreChannel {
//	return RdbUserSelectUserRooms(userId)
//}

func (provider MysqlProvider) UserUnreadCountUp(userId string) StoreChannel {
	return RdbUserUnreadCountUp(userId)
}

func (provider MysqlProvider) UserUnreadCountRecalc(userId string) StoreChannel {
	return RdbUserUnreadCountRecalc(userId)
}

func (provider MysqlProvider) UserUserIdsSelectByUserIds(userIds []string) StoreChannel {
	return RdbUserUserIdsSelectByUserIds(userIds)
}
