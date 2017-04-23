package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (provider SqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (provider SqliteProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel {
	return RdbInsertRoomUsers(roomUsers)
}

func (provider SqliteProvider) SelectRoomUser(roomId, userId string) StoreChannel {
	return RdbSelectRoomUser(roomId, userId)
}

func (provider SqliteProvider) SelectRoomUsersByRoomId(roomId string) StoreChannel {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (provider SqliteProvider) SelectRoomUsersByUserId(userId string) StoreChannel {
	return RdbSelectRoomUsersByUserId(userId)
}

func (provider SqliteProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (provider SqliteProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreChannel {
	return RdbUpdateRoomUser(roomUser)
}

func (provider SqliteProvider) DeleteRoomUser(roomId string, userIds []string) StoreChannel {
	return RdbDeleteRoomUser(roomId, userIds)
}

//func (provider SqliteProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
//	return RdbRoomUserInsert(roomUser)
//}
/*
func (provider SqliteProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}
*/
/*
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
*/
