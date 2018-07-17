package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *sqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*model.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*model.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *sqliteProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *sqliteProvider) SelectRoomUsersByRoomID(roomID string) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *sqliteProvider) SelectRoomUsersByUserID(userID string) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsOfRoomUser(roomID string, opts ...RoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *sqliteProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *model.RoomUser) (*model.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *sqliteProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
