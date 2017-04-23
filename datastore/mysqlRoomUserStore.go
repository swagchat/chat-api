package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (provider MysqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (provider MysqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbInsertRoomUsers(roomUsers)
}

func (provider MysqlProvider) SelectRoomUser(roomId, userId string) StoreResult {
	return RdbSelectRoomUser(roomId, userId)
}

func (provider MysqlProvider) SelectRoomUsersByRoomId(roomId string) StoreResult {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (provider MysqlProvider) SelectRoomUsersByUserId(userId string) StoreResult {
	return RdbSelectRoomUsersByUserId(userId)
}

func (provider MysqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (provider MysqlProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreResult {
	return RdbUpdateRoomUser(roomUser)
}

func (provider MysqlProvider) DeleteRoomUser(roomId string, userIds []string) StoreResult {
	return RdbDeleteRoomUser(roomId, userIds)
}

//func (provider MysqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreResult {
//	return RdbRoomUserInsert(roomUser)
//}
/*
func (provider MysqlProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreResult {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreResult {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}
*/
/*
func (provider MysqlProvider) RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbRoomUsersDeleteByRoomIdAndUserIds(roomId, userIds)
}

func (provider MysqlProvider) RoomUserDeleteByRoomId(roomId string) StoreResult {
	return RdbRoomUserDeleteByRoomId(roomId)
}

func (provider MysqlProvider) RoomUserDeleteByUserId(userId string) StoreResult {
	return RdbRoomUserDeleteByUserId(userId)
}

func (provider MysqlProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreResult {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider MysqlProvider) RoomUserMarkAllAsRead(userId string) StoreResult {
	return RdbRoomUserMarkAllAsRead(userId)
}
*/
