package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore(p.database)
}

func (p *gcpSqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSqlProvider) SelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	return RdbSelectRoomUser(p.database, roomId, userId)
}

func (p *gcpSqlProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	return RdbSelectRoomUserOfOneOnOne(p.database, myUserId, opponentUserId)
}

func (p *gcpSqlProvider) SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomId(p.database, roomId)
}

func (p *gcpSqlProvider) SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByUserId(p.database, userId)
}

func (p *gcpSqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomIdAndUserIds(p.database, roomId, userIds)
}

func (p *gcpSqlProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return RdbUpdateRoomUser(p.database, roomUser)
}

func (p *gcpSqlProvider) DeleteRoomUser(roomId string, userIds []string) error {
	return RdbDeleteRoomUser(p.database, roomId, userIds)
}
