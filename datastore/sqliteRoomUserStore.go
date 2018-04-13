package datastore

import "github.com/swagchat/chat-api/models"

func (p *sqliteProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.sqlitePath)
}

func (p *sqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.sqlitePath, roomUsers)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return rdbInsertRoomUsers(p.sqlitePath, roomUsers)
}

func (p *sqliteProvider) SelectRoomUser(roomID, userID string) (*models.RoomUser, error) {
	return rdbSelectRoomUser(p.sqlitePath, roomID, userID)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*models.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.sqlitePath, myUserID, opponentUserID)
}

func (p *sqliteProvider) SelectRoomUsersByRoomID(roomID string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.sqlitePath, roomID)
}

func (p *sqliteProvider) SelectRoomUsersByUserID(userID string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.sqlitePath, userID)
}

func (p *sqliteProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.sqlitePath, roomID, userIDs)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return rdbUpdateRoomUser(p.sqlitePath, roomUser)
}

func (p *sqliteProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.sqlitePath, roomID, userIDs)
}
