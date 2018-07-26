package datastore

import (
	"github.com/swagchat/chat-api/model"
)

func (p *mysqlProvider) createRoomUserStore() {
	rdbCreateRoomUserStore(p.database)
}

func (p *mysqlProvider) InsertRoomUsers(roomUsers []*model.RoomUser, opts ...InsertRoomUsersOption) error {
	return rdbInsertRoomUsers(p.database, roomUsers, opts...)
}

func (p *mysqlProvider) SelectRoomUsers(opts ...SelectRoomUsersOption) ([]*model.RoomUser, error) {
	return rdbSelectRoomUsers(p.database, opts...)
}

func (p *mysqlProvider) SelectRoomUser(roomID, userID string) (*model.RoomUser, error) {
	return rdbSelectRoomUser(p.database, roomID, userID)
}

func (p *mysqlProvider) SelectRoomUserOfOneOnOne(myUserID, opponentUserID string) (*model.RoomUser, error) {
	return rdbSelectRoomUserOfOneOnOne(p.database, myUserID, opponentUserID)
}

func (p *mysqlProvider) SelectUserIDsOfRoomUser(roomID string, opts ...SelectUserIDsOfRoomUserOption) ([]string, error) {
	return rdbSelectUserIDsOfRoomUser(p.database, roomID, opts...)
}

func (p *mysqlProvider) UpdateRoomUser(roomUser *model.RoomUser) error {
	return rdbUpdateRoomUser(p.database, roomUser)
}

func (p *mysqlProvider) DeleteRoomUsers(roomID string, userIDs []string) error {
	return rdbDeleteRoomUsers(p.database, roomID, userIDs)
}
