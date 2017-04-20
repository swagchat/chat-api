package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider SqliteProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider SqliteProvider) RoomUsersInsert(roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	return RdbRoomUsersInsert(roomUsers, isDeleteFirst)
}

func (provider SqliteProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
}

func (provider SqliteProvider) RoomUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersSelectByUserId(userId string) StoreChannel {
	return RdbRoomUsersSelectByUserId(userId)
}

func (provider SqliteProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersSelectByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectByRoomIdAndUserIds(roomId, userIds)
}

func (provider SqliteProvider) RoomUserUpdate(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserUpdate(roomUser)
}

func (provider SqliteProvider) RoomUserDelete(roomId string, userIds []string) StoreChannel {
	return RdbRoomUserDelete(roomId, userIds)
}

func (provider SqliteProvider) RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersDeleteByRoomIdAndUserIds(roomId, userIds)
}

func (provider SqliteProvider) RoomUserDeleteByRoomId(roomId string) StoreChannel {
	return RdbRoomUserDeleteByRoomId(roomId)
}

func (provider SqliteProvider) RoomUserDeleteByUserId(userId string) StoreChannel {
	return RdbRoomUserDeleteByUserId(userId)
}

func (provider SqliteProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider SqliteProvider) RoomUserMarkAllAsRead(userId string) StoreChannel {
	return RdbRoomUserMarkAllAsRead(userId)
}
