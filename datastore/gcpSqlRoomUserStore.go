package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *gcpSQLProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *gcpSQLProvider) DeleteAndInsertRoomUsers(roomUsers []*model.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSQLProvider) InsertRoomUsers(roomUsers []*model.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *gcpSQLProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *gcpSQLProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *gcpSQLProvider) SelectRoomUsersByRoomID(roomID string) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *gcpSQLProvider) SelectRoomUsersByUserID(userID string) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *gcpSQLProvider) SelectUserIDsOfRoomUser(roomID string, opts ...RoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *gcpSQLProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *gcpSQLProvider) UpdateRoomUser(roomUser *model.RoomUser) (*model.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *gcpSQLProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
