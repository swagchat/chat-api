package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) RoomUserCreateStore() {
	RdbRoomUserCreateStore()
}

func (provider SqliteProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
	return RdbRoomUserInsert(roomUser)
}

func (provider SqliteProvider) RoomUsersInsert(roomId string, roomUsers []*models.RoomUser, isDeleteFirst bool) StoreChannel {
	return RdbRoomUsersInsert(roomId, roomUsers, isDeleteFirst)
}

func (provider SqliteProvider) RoomUserUsersSelect(roomId string) StoreChannel {
	return RdbRoomUserUsersSelect(roomId)
}

func (provider SqliteProvider) RoomUsersSelect(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelect(roomId, userIds)
}

func (provider SqliteProvider) RoomUsersSelectUserIds(roomId string) StoreChannel {
	return RdbRoomUsersSelectUserIds(roomId)
}

func (provider SqliteProvider) RoomUsersSelectIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersSelectIds(roomId, userIds)
}

func (provider SqliteProvider) RoomUserSelect(roomId, userId string) StoreChannel {
	return RdbRoomUserSelect(roomId, userId)
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

func (provider SqliteProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreChannel {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider SqliteProvider) RoomUserMarkAllAsRead(userId string) StoreChannel {
	return RdbRoomUserMarkAllAsRead(userId)
}
