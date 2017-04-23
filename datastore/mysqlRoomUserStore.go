package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider MysqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (provider MysqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (provider MysqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel {
	return RdbInsertRoomUsers(roomUsers)
}

func (provider MysqlProvider) SelectRoomUser(roomId, userId string) StoreChannel {
	return RdbSelectRoomUser(roomId, userId)
}

func (provider MysqlProvider) SelectRoomUsersByRoomId(roomId string) StoreChannel {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (provider MysqlProvider) SelectRoomUsersByUserId(userId string) StoreChannel {
	return RdbSelectRoomUsersByUserId(userId)
}

func (provider MysqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (provider MysqlProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreChannel {
	return RdbUpdateRoomUser(roomUser)
}

func (provider MysqlProvider) DeleteRoomUser(roomId string, userIds []string) StoreChannel {
	return RdbDeleteRoomUser(roomId, userIds)
}

//func (provider MysqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
//	return RdbRoomUserInsert(roomUser)
//}
/*
func (provider MysqlProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider MysqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}
*/
/*
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
*/
