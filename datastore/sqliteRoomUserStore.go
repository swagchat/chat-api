package datastore

import (
	"github.com/swagchat/chat-api/protobuf"
)

func (p *sqliteProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *sqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*protobuf.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*protobuf.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *sqliteProvider) SelectRoomUser(roomID, userID string) (*protobuf.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*protobuf.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *sqliteProvider) SelectRoomUsersByRoomID(roomID string) ([]*protobuf.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *sqliteProvider) SelectRoomUsersByUserID(userID string) ([]*protobuf.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *sqliteProvider) SelectRoomUserIDsByRoomID(roomID string, opts ...interface{}) ([]string, error) {
	return rdbSelectRoomUserIDsByRoomID(p.database, roomID, opts...)
}

func (p *sqliteProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*protobuf.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *protobuf.RoomUser) (*protobuf.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *sqliteProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
