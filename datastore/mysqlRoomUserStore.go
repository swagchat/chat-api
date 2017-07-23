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

func (provider MysqlProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) StoreResult {
	return RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId)
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
