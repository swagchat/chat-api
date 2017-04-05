package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider GcpSqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider GcpSqlProvider) RoomUsersInsert(roomId string, roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	return RdbRoomUsersInsert(roomId, roomUsers, isDeleteFirst)
}

func (provider GcpSqlProvider) RoomUserUsersSelect(roomId string) StoreChannel {
	return RdbRoomUserUsersSelect(roomId)
}

func (provider GcpSqlProvider) RoomUsersSelect(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelect(roomId, userIds)
}

func (provider GcpSqlProvider) RoomUsersSelectUserIds(roomId string) StoreChannel {
	return RdbRoomUsersSelectUserIds(roomId)
}

func (provider GcpSqlProvider) RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectIds(roomId, userIds)
}

func (provider GcpSqlProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
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

func (provider GcpSqlProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider GcpSqlProvider) RoomUserMarkAllAsRead(userId string) StoreChannel {
	return RdbRoomUserMarkAllAsRead(userId)
}
