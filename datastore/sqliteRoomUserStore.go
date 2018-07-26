package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *sqliteProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *sqliteProvider) InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	return rdbInsertRoomUsers(p.database, roomUsers, opts...)
}

func (p *sqliteProvider) SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsers(p.database, opts...)
}

func (p *sqliteProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *sqliteProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *sqliteProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *sqliteProvider) UpdateRoomUser(roomUser *model.RoomUser) error {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *sqliteProvider) DeleteRoomUsers(roomID string, userIDs []string) error {
	return rdbDeleteRoomUsers(p.database, roomID, userIDs)
}
