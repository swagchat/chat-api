package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider MysqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider MysqlProvider) RoomUsersInsert(roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	return RdbRoomUsersInsert(roomUsers, isDeleteFirst)
}

func (provider MysqlProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
}

func (provider MysqlProvider) RoomUsersSelectByUserId(userId string) StoreChannel {
	return RdbRoomUsersSelectByUserId(userId)
}

func (provider MysqlProvider) RoomUserUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUserUsersSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectIds(roomId, userIds)
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

func (provider MysqlProvider) RoomUserDeleteByRoomId(roomId string) StoreChannel {
	return RdbRoomUserDeleteByRoomId(roomId)
}

func (provider MysqlProvider) RoomUserDeleteByUserId(userId string) StoreChannel {
	return RdbRoomUserDeleteByUserId(userId)
}

func (provider MysqlProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider MysqlProvider) RoomUserMarkAllAsRead(userId string) StoreChannel {
	return RdbRoomUserMarkAllAsRead(userId)
}
