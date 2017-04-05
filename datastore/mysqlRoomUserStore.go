package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider MysqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider MysqlProvider) RoomUsersInsert(roomId string, roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	return RdbRoomUsersInsert(roomId, roomUsers, isDeleteFirst)
}

func (provider MysqlProvider) RoomUserUsersSelect(roomId string) StoreChannel {
	return RdbRoomUserUsersSelect(roomId)
}

func (provider MysqlProvider) RoomUsersSelect(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelect(roomId, userIds)
}

func (provider MysqlProvider) RoomUsersSelectUserIds(roomId string) StoreChannel {
	return RdbRoomUsersSelectUserIds(roomId)
}

func (provider MysqlProvider) RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectIds(roomId, userIds)
}

func (provider MysqlProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
}

func (provider MysqlProvider) RoomUserUpdate(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserUpdate(roomUser)
}

func (provider MysqlProvider) RoomUserDelete(roomId string, userIds []string) StoreChannel {
	return RdbRoomUserDelete(roomId, userIds)
}

func (provider MysqlProvider) RoomUsersDeleteByUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersDeleteByUserIds(roomId, userIds)
}

func (provider MysqlProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider MysqlProvider) RoomUserMarkAllAsRead(userId string) StoreChannel {
	return RdbRoomUserMarkAllAsRead(userId)
}
