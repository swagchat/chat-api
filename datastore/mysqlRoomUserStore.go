package datastore

import (
	scpb "github.com/swagchat/protobuf"
)

func (p *mysqlProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *mysqlProvider) DeleteAndInsertRoomUsers(roomUsers []*scpb.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*scpb.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *mysqlProvider) SelectRoomUser(roomID, userID string) (*scpb.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*scpb.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *mysqlProvider) SelectRoomUsersByRoomID(roomID string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *mysqlProvider) SelectRoomUsersByUserID(userID string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *mysqlProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *mysqlProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *scpb.RoomUser) (*scpb.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *mysqlProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
