package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (provider GcpSqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (provider GcpSqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (provider GcpSqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreChannel {
	return RdbInsertRoomUsers(roomUsers)
}

func (provider GcpSqlProvider) SelectRoomUser(roomId, userId string) StoreChannel {
	return RdbSelectRoomUser(roomId, userId)
}

func (provider GcpSqlProvider) SelectRoomUsersByRoomId(roomId string) StoreChannel {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (provider GcpSqlProvider) SelectRoomUsersByUserId(userId string) StoreChannel {
	return RdbSelectRoomUsersByUserId(userId)
}

func (provider GcpSqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (provider GcpSqlProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreChannel {
	return RdbUpdateRoomUser(roomUser)
}

func (provider GcpSqlProvider) DeleteRoomUser(roomId string, userIds []string) StoreChannel {
	return RdbDeleteRoomUser(roomId, userIds)
}

//func (provider GcpSqlProvider) RoomUserInsert(roomUser *models.RoomUser) StoreChannel {
//	return RdbRoomUserInsert(roomUser)
//}

/*
func (provider GcpSqlProvider) RoomUsersUsersSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUsersSelectByRoomId(roomId)
}

func (provider GcpSqlProvider) RoomUsersUserIdsSelectByRoomId(roomId string) StoreChannel {
	return RdbRoomUsersUserIdsSelectByRoomId(roomId)
}
*/
/*
func (provider GcpSqlProvider) RoomUsersDeleteByRoomIdAndUserIds(roomId *string, userIds []string) StoreChannel {
	return RdbRoomUsersDeleteByRoomIdAndUserIds(roomId, userIds)
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
*/
