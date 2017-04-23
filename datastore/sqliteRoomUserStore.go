package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider SqliteProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (provider SqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (provider SqliteProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbInsertRoomUsers(roomUsers)
}

func (provider SqliteProvider) SelectRoomUser(roomId, userId string) StoreResult {
	return RdbSelectRoomUser(roomId, userId)
}

func (provider SqliteProvider) SelectRoomUsersByRoomId(roomId string) StoreResult {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (provider SqliteProvider) SelectRoomUsersByUserId(userId string) StoreResult {
	return RdbSelectRoomUsersByUserId(userId)
}

func (provider SqliteProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (provider SqliteProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreResult {
	return RdbUpdateRoomUser(roomUser)
}

func (provider SqliteProvider) DeleteRoomUser(roomId string, userIds []string) StoreResult {
	return RdbDeleteRoomUser(roomId, userIds)
}

//func (provider SqliteProvider) RoomUserInsert(roomUser *models.RoomUser) StoreResult {
//	return RdbRoomUserInsert(roomUser)
//}
/*
func (provider SqliteProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreResult {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider SqliteProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreResult {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}
*/
/*
func (provider SqliteProvider) RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbRoomUsersDeleteByRoomIdAndUserIds(roomId, userIds)
}

func (provider SqliteProvider) RoomUserDeleteByRoomId(roomId string) StoreResult {
	return RdbRoomUserDeleteByRoomId(roomId)
}

func (provider SqliteProvider) RoomUserDeleteByUserId(userId string) StoreResult {
	return RdbRoomUserDeleteByUserId(userId)
}

func (provider SqliteProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreResult {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider SqliteProvider) RoomUserMarkAllAsRead(userId string) StoreResult {
	return RdbRoomUserMarkAllAsRead(userId)
}
*/
