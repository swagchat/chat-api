package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) UserCreateStore() {
	RdbUserCreateStore()
}

func (provider GcpSqlProvider) UserInsert(user *models.User) StoreChannel {
	return RdbUserInsert(user)
}

func (provider GcpSqlProvider) UserSelect(userId string, isWithRooms, isWithDevices bool) StoreChannel {
	return RdbUserSelect(userId, isWithRooms, isWithDevices)
}

func (provider GcpSqlProvider) UserUpdate(user *models.User) StoreChannel {
	return RdbUserUpdate(user)
}

func (provider GcpSqlProvider) UserUpdateDeleted(userId string) StoreChannel {
	return RdbUserUpdateDeleted(userId)
}

func (provider GcpSqlProvider) UserSelectAll() StoreChannel {
	return RdbUserSelectAll()
}

func (provider GcpSqlProvider) UserSelectRoomsForUser(userId string) StoreChannel {
	return RdbUserSelectRoomsForUser(userId)
}

//func (provider GcpSqlProvider) UserSelectUserRooms(userId string) StoreChannel {
//	return RdbUserSelectUserRooms(userId)
//}

func (provider GcpSqlProvider) UserUnreadCountUp(userId string) StoreChannel {
	return RdbUserUnreadCountUp(userId)
}

func (provider GcpSqlProvider) UserUnreadCountRecalc(userId string) StoreChannel {
	return RdbUserUnreadCountRecalc(userId)
}

func (provider GcpSqlProvider) UserUserIdsSelectByUserIds(userIds []string) StoreChannel {
	return RdbUserUserIdsSelectByUserIds(userIds)
}
