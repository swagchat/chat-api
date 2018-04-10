package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore(p.sqlitePath)
}

func (p *sqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbDeleteAndInsertRoomUsers(p.sqlitePath, roomUsers)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbInsertRoomUsers(p.sqlitePath, roomUsers)
}

func (p *sqliteProvider) SelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	return RdbSelectRoomUser(p.sqlitePath, roomId, userId)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	return RdbSelectRoomUserOfOneOnOne(p.sqlitePath, myUserId, opponentUserId)
}

func (p *sqliteProvider) SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomId(p.sqlitePath, roomId)
}

func (p *sqliteProvider) SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByUserId(p.sqlitePath, userId)
}

func (p *sqliteProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomIdAndUserIds(p.sqlitePath, roomId, userIds)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return RdbUpdateRoomUser(p.sqlitePath, roomUser)
}

func (p *sqliteProvider) DeleteRoomUser(roomId string, userIds []string) error {
	return RdbDeleteRoomUser(p.sqlitePath, roomId, userIds)
}
