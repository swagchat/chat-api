package datastore

import "github.com/swagchat/chat-api/models"

func (p *mysqlProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *mysqlProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *mysqlProvider) SelectRoomUser(roomID, userID string) (*models.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*models.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *mysqlProvider) SelectRoomUsersByRoomID(roomID string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *mysqlProvider) SelectRoomUsersByUserID(userID string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectRoomUserIDsByRoomID(roomID string) ([]string, error) {
	return rdbSelectRoomUserIDsByRoomID(p.database, roomID)
}

func (p *mysqlProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *mysqlProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
