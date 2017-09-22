package datastore

import "github.com/fairway-corp/swagchat-api/models"

func (p *gcpSqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (p *gcpSqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (p *gcpSqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) StoreResult {
	return RdbInsertRoomUsers(roomUsers)
}

func (p *gcpSqlProvider) SelectRoomUser(roomId, userId string) StoreResult {
	return RdbSelectRoomUser(roomId, userId)
}

func (p *gcpSqlProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) StoreResult {
	return RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId)
}

func (p *gcpSqlProvider) SelectRoomUsersByRoomId(roomId string) StoreResult {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (p *gcpSqlProvider) SelectRoomUsersByUserId(userId string) StoreResult {
	return RdbSelectRoomUsersByUserId(userId)
}

func (p *gcpSqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) StoreResult {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (p *gcpSqlProvider) UpdateRoomUser(roomUser *models.RoomUser) StoreResult {
	return RdbUpdateRoomUser(roomUser)
}

func (p *gcpSqlProvider) DeleteRoomUser(roomId string, userIds []string) StoreResult {
	return RdbDeleteRoomUser(roomId, userIds)
}
