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

func (provider SqliteProvider) RoomUsersSelectByUserId(userId string) StoreChannel {
	return RdbRoomUsersSelectByUserId(userId)
}

func (provider SqliteProvider) RoomUserUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUserUsersSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectIds(roomId, userIds)
}

func (provider SqliteProvider) RoomUserUpdate(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserUpdate(roomUser)
}

func (provider SqliteProvider) RoomUserDelete(roomId string, userIds []string) StoreChannel {
	return RdbRoomUserDelete(roomId, userIds)
}

func (provider SqliteProvider) RoomUsersDeleteByUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersDeleteByUserIds(roomId, userIds)
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
