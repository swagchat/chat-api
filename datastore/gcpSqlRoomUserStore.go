package datastore

import "github.com/swagchat/chat-api/models"

func (p *gcpSQLProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *gcpSQLProvider) DeleteAndInsertRoomUsers(roomUsers []*models.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSQLProvider) InsertRoomUsers(roomUsers []*models.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSQLProvider) SelectRoomUser(roomID, userID string) (*models.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *gcpSQLProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*models.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *gcpSQLProvider) SelectRoomUsersByRoomID(roomID string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *gcpSQLProvider) SelectRoomUsersByUserID(userID string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectRoomUserIDsByRoomID(roomID string, opts ...interface{}) ([]string, error) {
	return rdbSelectRoomUserIDsByRoomID(p.database, roomID, opts...)
}

func (p *gcpSQLProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*models.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *gcpSQLProvider) UpdateRoomUser(roomUser *models.RoomUser) (*models.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *gcpSQLProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
