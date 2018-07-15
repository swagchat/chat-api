package datastore

import (
	scpb "github.com/swagchat/protobuf"
)

func (p *sqliteProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *sqliteProvider) DeleteAndInsertRoomUsers(roomUsers []*scpb.RoomUser) error {
	return rdbDeleteAndInsertRoomUsers(p.database, roomUsers)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*scpb.RoomUser) error {
	return rdbInsertRoomUsers(p.database, roomUsers)
}

func (p *sqliteProvider) SelectRoomUser(roomID, userID string) (*scpb.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*scpb.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *sqliteProvider) SelectRoomUsersByRoomID(roomID string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByRoomID(p.database, roomID)
}

func (p *sqliteProvider) SelectRoomUsersByUserID(userID string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByUserID(p.database, userID)
}

func (p *sqliteProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *sqliteProvider) SelectRoomUsersByRoomIDAndUserIDs(roomID *string, userIDs []string) ([]*scpb.RoomUser, error) {
	return rdbSelectRoomUsersByRoomIDAndUserIDs(p.database, roomID, userIDs)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *scpb.RoomUser) (*scpb.RoomUser, error) {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *sqliteProvider) DeleteRoomUser(roomID string, userIDs []string) error {
	return rdbDeleteRoomUser(p.database, roomID, userIDs)
}
