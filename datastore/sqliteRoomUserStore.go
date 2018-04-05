package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (p *sqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbInsertRoomUsers(roomUsers)
}

func (p *sqliteProvider) SelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	return RdbSelectRoomUser(roomId, userId)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	return RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId)
}

func (p *sqliteProvider) SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (p *sqliteProvider) SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByUserId(userId)
}

func (p *sqliteProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return RdbUpdateRoomUser(roomUser)
}

func (p *sqliteProvider) DeleteRoomUser(roomId string, userIds []string) error {
	return RdbDeleteRoomUser(roomId, userIds)
}
