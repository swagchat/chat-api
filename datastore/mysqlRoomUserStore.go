package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore(p.database)
}

func (p *mysqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbInsertRoomUsers(p.database, roomUsers)
}

func (p *mysqlProvider) SelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	return RdbSelectRoomUser(p.database, roomId, userId)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	return RdbSelectRoomUserOfOneOnOne(p.database, myUserId, opponentUserId)
}

func (p *mysqlProvider) SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomId(p.database, roomId)
}

func (p *mysqlProvider) SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByUserId(p.database, userId)
}

func (p *mysqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomIdAndUserIds(p.database, roomId, userIds)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return RdbUpdateRoomUser(p.database, roomUser)
}

func (p *mysqlProvider) DeleteRoomUser(roomId string, userIds []string) error {
	return RdbDeleteRoomUser(p.database, roomId, userIds)
}
