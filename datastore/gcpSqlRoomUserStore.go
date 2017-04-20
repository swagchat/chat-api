package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider GcpSqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider GcpSqlProvider) RoomUsersInsert(roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	return RdbRoomUsersInsert(roomUsers, isDeleteFirst)
}

func (provider GcpSqlProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
}

func (provider GcpSqlProvider) RoomUsersSelectByUserId(userId string) StoreChannel {
	return RdbRoomUsersSelectByUserId(userId)
}

func (provider GcpSqlProvider) RoomUserUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUserUsersSelectByRoomId(roomId)
}

func (provider GcpSqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}

func (provider GcpSqlProvider) RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectIds(roomId, userIds)
}

func (provider GcpSqlProvider) RoomUserUpdate(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserUpdate(roomUser)
}

func (provider GcpSqlProvider) RoomUserDelete(roomId string, userIds []string) StoreChannel {
	return RdbRoomUserDelete(roomId, userIds)
}

func (provider GcpSqlProvider) RoomUsersDeleteByUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersDeleteByUserIds(roomId, userIds)
}

func (provider GcpSqlProvider) RoomUserDeleteByRoomId(roomId string) StoreChannel {
	return RdbRoomUserDeleteByRoomId(roomId)
}

func (provider GcpSqlProvider) RoomUserDeleteByUserId(userId string) StoreChannel {
	return RdbRoomUserDeleteByUserId(userId)
}

func (provider GcpSqlProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider GcpSqlProvider) RoomUserMarkAllAsRead(userId string) StoreChannel {
	return RdbRoomUserMarkAllAsRead(userId)
}
