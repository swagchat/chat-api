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

func (provider SqliteProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) StoreResult {
	return RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId)
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
