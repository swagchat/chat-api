package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (provider GcpSqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (provider GcpSqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbInsertRoomUsers(roomUsers)
}

func (provider GcpSqlProvider) SelectRoomUser(roomId, userId string) StoreResult {
	return RdbSelectRoomUser(roomId, userId)
}

func (provider GcpSqlProvider) SelectRoomUsersByRoomId(roomId string) StoreResult {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (provider GcpSqlProvider) SelectRoomUsersByUserId(userId string) StoreResult {
	return RdbSelectRoomUsersByUserId(userId)
}

func (provider GcpSqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (provider GcpSqlProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreResult {
	return RdbUpdateRoomUser(roomUser)
}

func (provider GcpSqlProvider) DeleteRoomUser(roomId string, userIds []string) StoreResult {
	return RdbDeleteRoomUser(roomId, userIds)
}

//func (provider GcpSqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreResult {
//	return RdbRoomUserInsert(roomUser)
//}

/*
func (provider GcpSqlProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreResult {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider GcpSqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreResult {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}
*/
/*
func (provider GcpSqlProvider) RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbRoomUsersDeleteByRoomIdAndUserIds(roomId, userIds)
}

func (provider GcpSqlProvider) RoomUserDeleteByRoomId(roomId string) StoreResult {
	return RdbRoomUserDeleteByRoomId(roomId)
}

func (provider GcpSqlProvider) RoomUserDeleteByUserId(userId string) StoreResult {
	return RdbRoomUserDeleteByUserId(userId)
}

func (provider GcpSqlProvider) RoomUserUnreadCountUp(roomId string, currentUserId string) StoreResult {
	return RdbRoomUserUnreadCountUp(roomId, currentUserId)
}

func (provider GcpSqlProvider) RoomUserMarkAllAsRead(userId string) StoreResult {
	return RdbRoomUserMarkAllAsRead(userId)
}
*/
