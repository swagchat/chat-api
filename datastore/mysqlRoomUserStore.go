package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (p *mysqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbInsertRoomUsers(roomUsers)
}

func (p *mysqlProvider) SelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	return RdbSelectRoomUser(roomId, userId)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	return RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId)
}

func (p *mysqlProvider) SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (p *mysqlProvider) SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByUserId(userId)
}

func (p *mysqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return RdbUpdateRoomUser(roomUser)
}

func (p *mysqlProvider) DeleteRoomUser(roomId string, userIds []string) error {
	return RdbDeleteRoomUser(roomId, userIds)
}
