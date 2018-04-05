package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSqlProvider) CreateRoomUserStore() {
	RdbCreateRoomUserStore()
}

func (p *gcpSqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbDeleteAndInsertRoomUsers(roomUsers)
}

func (p *gcpSqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return RdbInsertRoomUsers(roomUsers)
}

func (p *gcpSqlProvider) SelectRoomUser(roomId, userId string) (*models.RoomUser, error) {
	return RdbSelectRoomUser(roomId, userId)
}

func (p *gcpSqlProvider) SelectRoomUserOfOneOnOne(myUserId, opponentUserId string) (*models.RoomUser, error) {
	return RdbSelectRoomUserOfOneOnOne(myUserId, opponentUserId)
}

func (p *gcpSqlProvider) SelectRoomUsersByRoomId(roomId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomId(roomId)
}

func (p *gcpSqlProvider) SelectRoomUsersByUserId(userId string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByUserId(userId)
}

func (p *gcpSqlProvider) SelectRoomUsersByRoomIdAndUserIds(roomId *string, userIds []string) ([]*models.RoomUser, error) {
	return RdbSelectRoomUsersByRoomIdAndUserIds(roomId, userIds)
}

func (p *gcpSqlProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return RdbUpdateRoomUser(roomUser)
}

func (p *gcpSqlProvider) DeleteRoomUser(roomId string, userIds []string) error {
	return RdbDeleteRoomUser(roomId, userIds)
}
