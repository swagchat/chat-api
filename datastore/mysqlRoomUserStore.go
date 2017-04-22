package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider MysqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider MysqlProvider) RoomUsersDeleteAndInsert(roomUsers []*models.RoomUser) StoreChannel {
	return RdbRoomUsersDeleteAndInsert(roomUsers)
}

func (provider MysqlProvider) RoomUsersInsert(roomUsers []*models.RoomUser) StoreChannel {
	return RdbRoomUsersInsert(roomUsers)
}

func (provider MysqlProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
}

func (provider MysqlProvider) RoomUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersSelectByUserId(userId string) StoreChannel {
	return RdbRoomUsersSelectByUserId(userId)
}

func (provider MysqlProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersSelectByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectByRoomIdAndUserIds(roomId, userIds)
}

func (provider MysqlProvider) RoomUserUpdate(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserUpdate(roomUser)
}

func (provider MysqlProvider) RoomUserDelete(roomId string, userIds []string) StoreChannel {
	return RdbRoomUserDelete(roomId, userIds)
}

func (provider MysqlProvider) RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersDeleteByRoomIdAndUserIds(roomId, userIds)
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
